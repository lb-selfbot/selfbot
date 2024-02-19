package eventreactions

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/LightningDev1/LB-Selfbot/selfbot/api"
	"github.com/LightningDev1/LB-Selfbot/selfbot/config"
	"github.com/LightningDev1/LB-Selfbot/selfbot/http"
	"github.com/LightningDev1/LB-Selfbot/selfbot/logging"
	"github.com/LightningDev1/LB-Selfbot/selfbot/utils"
	"github.com/LightningDev1/discordgo"
	"github.com/LightningDev1/go-privnote"
	"github.com/bitly/go-simplejson"
)

var (
	giftCodeRegex = regexp.MustCompile(`(discord\.gift|discord\.com/gifts|discordapp\.com/gifts|promos\.discord\.gg/)/(\w{0,30})`)
	privnoteRegex = regexp.MustCompile(`privnote.com/([a-zA-Z0-9_#]{18,18})`)
)

type GiftCode struct {
	Code string
	Time time.Time
}

type Snipers struct {
	Account       string
	SetupComplete bool

	CheckedGiftCodes  []string
	GiftCodesByAuthor map[string][]GiftCode

	CountryCode     string
	PaymentSourceID string

	PrivnoteClient *privnote.Client

	TotalPokemon  int
	CaughtPokemon int
	LastCatch     time.Time
}

func NewSnipers(account string) *Snipers {
	return &Snipers{
		Account:           account,
		GiftCodesByAuthor: make(map[string][]GiftCode),
		PrivnoteClient:    privnote.NewClient(),
	}
}

func (*Snipers) Type() EventReactionType {
	return SNIPERS
}

func (*Snipers) Description() string {
	return "Snipers (default)"
}

func (*Snipers) Equals(other EventReaction) bool {
	_, ok := other.(*Snipers)

	return ok
}

func (*Snipers) Data() map[string]any {
	return map[string]any{}
}

func (s *Snipers) Setup(session *discordgo.Session) {
	if s.SetupComplete {
		return
	}
	s.SetupComplete = true

	// Get payment source
	resp, err := session.Request("GET", "https://discord.com/api/v9/users/@me/billing/payment-sources", nil)
	if err != nil {
		utils.Logging.Debug("Snipers.Setup: Error getting payment sources: ", err)
		return
	}

	json, err := simplejson.NewJson(resp)
	if err != nil {
		utils.Logging.Debug("Snipers.Setup: Error parsing payment sources: ", err)
		return
	}

	// No payment sources
	if len(json.MustArray()) == 0 {
		return
	}

	for i := 0; i < len(json.MustArray()); i++ {
		source := json.GetIndex(i)

		if i == 0 || source.Get("default").MustBool() {
			s.PaymentSourceID = source.Get("id").MustString()
			s.CountryCode = source.Get("country").MustString()
		}
	}
}

func (s *Snipers) Handler(session *discordgo.Session, event any) {
	s.Setup(session)

	cfg, err := config.GetConfig()
	if err != nil {
		utils.Logging.Debug("Snipers.Handler: Error getting config: ", err)
		return
	}

	switch event := event.(type) {
	case *discordgo.MessageCreate:
		if utils.SliceContains(cfg.Snipers.ServerBlacklist, event.GuildID) {
			return
		}

		s.CheckNitroSniper(session, event)
		s.CheckPrivnoteSniper(session, event)
	}

}

func (s *Snipers) CheckNitroSniper(session *discordgo.Session, m *discordgo.MessageCreate) {
	cfg, err := config.GetConfig()
	if err != nil {
		utils.Logging.Debug("Snipers.CheckNitroSniper: Error getting config: ", err)
		return
	}

	if !cfg.Snipers.NitroSniper {
		return
	}

	rawMessage, err := json.Marshal(m)
	if err != nil {
		utils.Logging.Debug("Snipers.CheckNitroSniper: Error marshalling message: ", err)
		return
	}

	giftCodes := append(
		// Gift codes from regex
		giftCodeRegex.FindAllString(string(rawMessage), -1),

		// Embedded gift codes
		m.GiftCodes...,
	)

	for _, code := range giftCodes {
		// Remove URL from code
		splitCode := strings.Split(code, "/")
		code = splitCode[len(splitCode)-1]

		// Check if code has already been checked
		if utils.SliceContains(s.CheckedGiftCodes, code) {
			continue
		}

		s.CheckedGiftCodes = append(s.CheckedGiftCodes, code)

		// Check if the author is spamming codes
		s.GiftCodesByAuthor[m.Author.ID] = utils.Filter(
			s.GiftCodesByAuthor[m.Author.ID],
			func(giftCode GiftCode) bool {
				return time.Since(giftCode.Time) < 30*time.Second
			},
		)

		if len(s.GiftCodesByAuthor[m.Author.ID]) > 5 {
			continue
		}

		s.GiftCodesByAuthor[m.Author.ID] = append(
			s.GiftCodesByAuthor[m.Author.ID],
			GiftCode{Code: code, Time: time.Now()},
		)

		// Check code length
		if len(code) < 15 || len(code) > 26 {
			logging.Logging.NitroSniperFake(session, m.Message, code, "")

			continue
		}

		// Redeem code
		url := "https://discord.com/api/v9/entitlements/gift-codes/" + code + "/redeem"
		payload := map[string]any{
			"with_application":       false,
			"with_subscription_plan": true,
		}

		if s.CountryCode != "" {
			payload["country_code"] = s.CountryCode
		}

		if s.PaymentSourceID != "" {
			payload["payment_source_id"] = s.PaymentSourceID
		}

		beforeRequest := time.Now()
		resp, _ := session.Request("POST", url, payload)
		elapsed := time.Since(beforeRequest)

		respString := string(resp)

		// Code has been redeemed
		if strings.Contains(respString, "subscription_plan") {
			logging.Logging.NitroSniper(session, m.Message, elapsed, "Redeemed Nitro", code, s.Account)
			continue
		}

		// Code has already been redeemed
		if strings.Contains(respString, "This gift has been redeemed already.") {
			logging.Logging.NitroSniper(session, m.Message, elapsed, "Already Redeemed", code, s.Account)
			continue
		}

		// User already has nitro
		if strings.Contains(respString, "Already purchased") {
			logging.Logging.NitroSniper(session, m.Message, elapsed, "Already Purchased", code, s.Account)
			continue
		}

		// Invalid code
		if strings.Contains(respString, "Unknown Gift Code") {
			logging.Logging.NitroSniper(session, m.Message, elapsed, "Unknown Code", code, s.Account)
			continue
		}

		// Unknown status, try to parse JSON
		json, err := simplejson.NewJson(resp)
		if err != nil {
			json = simplejson.New()
		}

		status := utils.Misc.Or(json.Get("message").MustString(), "Unknown Status")

		logging.Logging.NitroSniper(session, m.Message, elapsed, status, code, s.Account)
	}

}

func (s *Snipers) CheckPrivnoteSniper(session *discordgo.Session, m *discordgo.MessageCreate) {
	cfg, err := config.GetConfig()
	if err != nil {
		utils.Logging.Debug("Snipers.CheckPrivnoteSniper: Error getting config: ", err)
		return
	}

	if !cfg.Snipers.PrivnoteSniper {
		return
	}

	rawMessage, err := json.Marshal(m)
	if err != nil {
		utils.Logging.Debug("Snipers.CheckNitroSniper: Error marshalling message: ", err)
		return
	}

	privnoteCodes := privnoteRegex.FindAllString(string(rawMessage), -1)

	for _, code := range privnoteCodes {
		code = "https://" + code

		// Read note
		beforeRequest := time.Now()
		note, err := s.PrivnoteClient.ReadNoteFromLink(code)
		if err != nil {
			utils.Logging.Debug("Snipers.CheckPrivnoteSniper: Error reading note: ", err)
			continue
		}
		elapsed := time.Since(beforeRequest)

		// Remove URL from code
		splitCode := strings.Split(code, "/")
		code = splitCode[len(splitCode)-1]

		// Save note
		fileName := utils.IO.SanitizeFilename("privnote-" + code + ".txt")
		filePath := filepath.Join(utils.IO.GetMainFolder(), "Downloads", fileName)

		err = utils.IO.Write(filePath, note)
		if err != nil {
			utils.Logging.Debug("Snipers.CheckPrivnoteSniper: Error writing note: ", err)
			continue
		}

		logging.Logging.PrivnoteSniper(session, m.Message, elapsed, "Saved Note", filePath, s.Account)
	}
}

func (s *Snipers) CheckGiveawaySniper(session *discordgo.Session, m *discordgo.MessageCreate) {
	cfg, err := config.GetConfig()
	if err != nil {
		utils.Logging.Debug("Snipers.CheckGiveawaySniper: Error getting config: ", err)
		return
	}

	if !cfg.Snipers.GiveawaySniper.Enabled {
		return
	}
}

func (s *Snipers) CheckPokeTwoSniper(session *discordgo.Session, m *discordgo.MessageCreate) {
	cfg, err := config.GetConfig()
	if err != nil {
		utils.Logging.Debug("Snipers.CheckPokeTwoSniper: Error getting config: ", err)
		return
	}

	sniperConditions := cfg.Snipers.PokeTwoSniper.Enabled && // Sniper is enabled
		m.Author.ID == "716390085896962058" && // Message is from PokeTwo
		utils.SliceContains(cfg.Snipers.PokeTwoSniper.Channels, m.ChannelID) && // Message is in a sniper channel
		len(m.Embeds) == 1

	if !sniperConditions {
		return
	}

	emb := m.Embeds[0]

	// Check if the message is a wild pokémon
	if !strings.Contains(emb.Title, "has appeared!") || emb.Image == nil {
		return
	}

	s.TotalPokemon++

	// Check if timeout has been reached
	if time.Since(s.LastCatch) < time.Duration(cfg.Snipers.PokeTwoSniper.Timeout)*time.Second {
		logging.UI.Info("PokéTwo Sniper: Catch timeout reached! Skipping...")
		return
	}

	// Check if the catch rate has been reached
	currentCatchRate := int(float64(s.CaughtPokemon) / float64(s.TotalPokemon) * 100)
	if int(currentCatchRate*100) > cfg.Snipers.PokeTwoSniper.CatchRate {
		logging.UI.Info(fmt.Sprintf("PokéTwo Sniper: Catch rate reached! Current catch rate is %d%%. Skipping...", currentCatchRate))
		return
	}

	// Start typing
	if cfg.Snipers.PokeTwoSniper.Typing {
		go func() {
			_ = session.ChannelTyping(m.ChannelID)
		}()
	}

	// Check if the Pokémon is already caught
	congratsMessage := make(chan *discordgo.Message)

	remove := session.AddHandler(func(s *discordgo.Session, newMessage *discordgo.MessageCreate) {
		if newMessage.Author.ID != "716390085896962058" || // Message is from PokeTwo
			newMessage.ChannelID != m.ChannelID || // Message is in the same channel
			!strings.Contains(newMessage.Content, "Congratulations") { // Message is a congratulations message
			return
		}

		congratsMessage <- newMessage.Message
	})

	// Download the Pokémon image
	resp := http.Get(emb.Image.URL)
	if resp.Error != nil {
		utils.Logging.Debug("Snipers.CheckPokeTwoSniper: Error getting image: ", resp.Error)
		return
	}

	image := utils.Misc.ImageToBase64(resp.BodyBytes)

	// Detect the Pokémon
	beforeRequest := time.Now()
	label, confidence, err := api.PokeTwo(image, &cfg)
	if err != nil {
		utils.Logging.Debug("Snipers.CheckPokeTwoSniper: Error getting Pokémon: ", err)
		return
	}
	elapsed := time.Since(beforeRequest)

	// Calculate the sleep time
	sleep := time.Duration(cfg.Snipers.PokeTwoSniper.Sleep) * time.Second
	if cfg.Snipers.PokeTwoSniper.RandomSleep {
		// Multiply the sleep time by a random number between 0.5 and 1
		sleep *= time.Duration(0.5 + 0.5*rand.Float64())
	}

	// Wait for the sleep time
	sleep = sleep - elapsed
	if sleep > 0 {
		time.Sleep(sleep)
	}

	// Check if the Pokémon is already caught
	select {
	case <-congratsMessage:
		logging.UI.Info("PokéTwo Sniper: Pokémon already caught! Skipping...")
		remove()
		return
	default:
	}

	// Catch the Pokémon
	_, err = session.ChannelMessageSend(m.ChannelID, "<@716390085896962058> c "+strings.ToLower(label))
	if err != nil {
		utils.Logging.Debug("Snipers.CheckPokeTwoSniper: Error catching Pokémon: ", err)
		return
	}

	// Wait for the Pokémon to be caught
	select {
	case m := <-congratsMessage:
		var mentioned bool
		for _, user := range m.Mentions {
			if user.ID == session.State.User.ID {
				mentioned = true
				break
			}
		}

		if !mentioned {
			logging.UI.Info("PokéTwo Sniper: Pokémon caught by someone else! Skipping...")
			break
		}

		s.CaughtPokemon++
		s.LastCatch = time.Now()

		// Remove the mention
		_ = session.MessageAck(m.ChannelID, m.ID)

		// Logging
		var (
			pokemon  string
			received string
		)

		pokemonMatches := regexp.MustCompile(`You caught a (.*?) <`).FindStringSubmatch(m.Content)
		if len(pokemonMatches) > 1 {
			pokemon = pokemonMatches[1]
		}

		receivedMatches := regexp.MustCompile(`You received (.*?)!`).FindStringSubmatch(m.Content)
		if len(receivedMatches) > 1 {
			received = receivedMatches[1]
		}

		logging.Logging.PokeTwoSniper(session, m, pokemon, received, confidence, elapsed)

		// Fake activity

		go func() {
			if !cfg.Snipers.PokeTwoSniper.FakeActivity {
				return
			}

			// Wait for random delay --> type --> wait for random delay
			time.Sleep(time.Second * time.Duration(utils.Random.Number(1, 2)))
			_ = session.ChannelTyping(m.ChannelID)
			time.Sleep(time.Second * time.Duration(utils.Random.Number(1, 3)))

			// Random PokéTwo commands
			commands := []string{
				"ping", "ev mg", "sh", "pokedex", "s latest", "pokedex",
				"d " + strings.ToLower(label),
			}

			_, _ = session.ChannelMessageSend(
				m.ChannelID,
				"<@716390085896962058> "+utils.RandomChoice(commands),
			)
		}()

	case <-time.After(5 * time.Second):
		utils.Logging.Debug("Snipers.CheckPokeTwoSniper: Pokémon not caught in time")
	}

	// Remove the handler
	remove()
}
