package eventreactions

import (
	"encoding/json"
	"fmt"
	"go/token"
	"go/types"
	"math/rand"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/LightningDev1/LB-Selfbot/internal/random"
	"github.com/LightningDev1/LB-Selfbot/selfbot/api"
	"github.com/LightningDev1/LB-Selfbot/selfbot/config"
	"github.com/LightningDev1/LB-Selfbot/selfbot/data"
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
	MainSession *discordgo.Session
	Account     string

	CheckedGiftCodes  []string
	GiftCodesByAuthor map[string][]GiftCode

	PrivnoteClient *privnote.Client

	PokeTwoCaptchaRequired bool
	TotalPokemon           int
	CaughtPokemon          int
	LastCatch              time.Time
}

func NewSnipers(mainSession *discordgo.Session, account string) *Snipers {
	return &Snipers{
		MainSession:       mainSession,
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

func (s *Snipers) Handler(session *discordgo.Session, event any) {
	cfg, err := config.GetConfig()
	if err != nil {
		utils.Logging.Debug("Snipers.Handler: Error getting config: ", err)
		return
	}

	if s.MainSession == nil {
		s.MainSession = session
	}

	switch event := event.(type) {
	case *discordgo.MessageCreate:
		if slices.Contains(cfg.Snipers.ServerBlacklist, event.GuildID) {
			return
		}

		s.CheckNitroSniper(session, event.Message)
		s.CheckPrivnoteSniper(session, event.Message)
		s.CheckPokeTwoSniper(session, event)
		s.CheckGiveawaySniper(session, event)
	case *discordgo.MessageUpdate:
		if slices.Contains(cfg.Snipers.ServerBlacklist, event.GuildID) {
			return
		}

		s.CheckNitroSniper(session, event.Message)
		s.CheckPrivnoteSniper(session, event.Message)
		s.CheckGiveawaySniperEdit(session, event)
	}
}

func (s *Snipers) CheckNitroSniper(session *discordgo.Session, m *discordgo.Message) {
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
		if slices.Contains(s.CheckedGiftCodes, code) {
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
			logging.Logging.NitroSniperFake(session, m, code, "")

			continue
		}

		// Redeem code
		url := "https://discord.com/api/v9/entitlements/gift-codes/" + code + "/redeem"
		payload := map[string]any{
			"with_application":       false,
			"with_subscription_plan": true,
			"channel_id":             m.ChannelID,
		}

		if data.CountryCode != "" {
			payload["country_code"] = data.CountryCode
		}

		if data.PaymentSourceID != "" {
			payload["payment_source_id"] = data.PaymentSourceID
		}

		headers := make(map[string]string)
		for k, v := range session.Headers {
			headers[k] = v
		}

		headers["Authorization"] = data.Token

		beforeRequest := time.Now()
		resp, _ := s.MainSession.Request("POST", url, payload)
		elapsed := time.Since(beforeRequest)

		respString := string(resp)

		// Code has been redeemed
		if strings.Contains(respString, "subscription_plan") {
			logging.Logging.NitroSniper(session, m, elapsed, "Redeemed Nitro", code, s.Account)
			continue
		}

		// Code has already been redeemed
		if strings.Contains(respString, "This gift has been redeemed already.") {
			logging.Logging.NitroSniper(session, m, elapsed, "Already Redeemed", code, s.Account)
			continue
		}

		// User already has nitro
		if strings.Contains(respString, "Already purchased") {
			logging.Logging.NitroSniper(session, m, elapsed, "Already Purchased", code, s.Account)
			continue
		}

		// Invalid code
		if strings.Contains(respString, "Unknown Gift Code") {
			logging.Logging.NitroSniper(session, m, elapsed, "Unknown Code", code, s.Account)
			continue
		}

		// Unknown status, try to parse JSON
		json, err := simplejson.NewJson(resp)
		if err != nil {
			json = simplejson.New()
		}

		status := utils.Misc.Or(json.Get("message").MustString(), "Unknown Status")

		logging.Logging.NitroSniper(session, m, elapsed, status, code, s.Account)
	}
}

func (s *Snipers) CheckPrivnoteSniper(session *discordgo.Session, m *discordgo.Message) {
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

		logging.Logging.PrivnoteSniper(session, m, elapsed, "Saved Note", filePath, s.Account)
	}
}

func (s *Snipers) getButtonComponent(m *discordgo.Message, i int) (button *discordgo.Button, ok bool) {
	if len(m.Components) == 0 {
		return
	}

	actionsRow, ok := m.Components[0].(*discordgo.ActionsRow)
	if !ok {
		return
	}

	if len(actionsRow.Components) <= i {
		return
	}

	button, ok = actionsRow.Components[i].(*discordgo.Button)
	return
}

func (s *Snipers) giveawaySniperSleepDuration(cfg *config.Config, endTime *time.Time) time.Duration {
	sleep := time.Duration(cfg.Snipers.GiveawaySniper.Sleep) * time.Second

	if endTime != nil && time.Until(*endTime) < sleep {
		// Multiply the time until endTime by a random number between 0.2 and 0.7
		untilEnd := float64(time.Until(*endTime))
		sleep = time.Duration(untilEnd * (0.2 + 0.5*rand.Float64()))

		logging.UI.Info("Giveaway ends in less than the sleep duration, sleeping for " + utils.Misc.FormatDuration(sleep) + " instead.")
	}

	return sleep
}

func (s *Snipers) isBlacklisted(cfg *config.Config, text string) bool {
	for _, blacklist := range cfg.Snipers.Blacklist {
		if strings.Contains(strings.ToLower(text), strings.ToLower(blacklist)) {
			return true
		}
	}

	return false
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

	switch m.Author.ID {
	case "676039466175168543": // ZAP-Hosting Bot#7133
		// Check if the message is a giveaway
		if len(m.Embeds) == 0 || !strings.Contains(m.Embeds[0].Title, "Giveaway") {
			break
		}

		// Get enter button
		button, ok := s.getButtonComponent(m.Message, 0)
		if !ok {
			break
		}

		// Add handler for submitting the modal
		session.AddHandlerOnce(func(_ *discordgo.Session, modal *discordgo.InteractionModalCreate) {
			actionsRow, ok := modal.Components[0].(*discordgo.ActionsRow)
			if !ok {
				return
			}

			textInput, ok := actionsRow.Components[0].(*discordgo.TextInput)
			if !ok {
				return
			}

			equation := strings.Split(textInput.Label, "\n")[1] // Get second line of the label
			equation = strings.ReplaceAll(equation, "=", "")    // Remove the equals sign
			equation = strings.ReplaceAll(equation, " ", "")    // Remove spaces

			// Evaluate the equation
			typeValue, err := types.Eval(token.NewFileSet(), nil, token.NoPos, equation)
			if err != nil {
				utils.Logging.Debug("Snipers.CheckGiveawaySniper: Error evaluating equation: ", err)
				return
			}

			// Send the solution of the equation
			solution := typeValue.Value.String()

			time.Sleep(time.Duration(random.Int(5, 10)) * time.Second)

			err = session.InteractionSubmitModal(m.Message, modal, textInput.CustomID, solution)
			if err != nil {
				utils.Logging.Debug("Snipers.CheckGiveawaySniper: Error submitting modal: ", err)
			}
		})

		logging.Logging.GiveawayFound(session, m.Message, "€10 worth of ZAP Coins", "ZAP-Hosting", "1 hour", s.Account)

		go func() {
			time.Sleep(s.giveawaySniperSleepDuration(&cfg, nil))

			// Click the enter button
			err = session.InteractionClick(button, m.Message)
			if err != nil {
				utils.Logging.Debug("Snipers.CheckGiveawaySniper: Error clicking button: ", err)
			}
		}()
	case "294882584201003009", "1041302711750119504": // GiveawayBot#2381, LightningBot#1992
		// Win check
		if m.Type == discordgo.MessageTypeReply && utils.Discord.IsUserMentioned(m.Message, session.State.User.ID) {
			winRegex := regexp.MustCompile(`You won the \*\*(.*?)\*\*!`)

			wins := winRegex.FindAllStringSubmatch(m.Content, -1)
			for _, win := range wins {
				prize := win[1]

				logging.Logging.GiveawayWon(session, m.Message, prize, s.Account)
			}
		}

		// Check if the message could be a giveaway
		button, ok := s.getButtonComponent(m.Message, 0)
		if len(m.Embeds) == 0 || !ok {
			break
		}

		// Get giveaway info
		prize := m.Embeds[0].Title
		if s.isBlacklisted(&cfg, prize) {
			break
		}

		hostRegex := regexp.MustCompile(`Hosted by: <@!?(\d+)>`)
		hostMatches := hostRegex.FindStringSubmatch(m.Embeds[0].Description)

		endsRegex := regexp.MustCompile(`Ends: <t:(\d+):R>`)
		endsMatches := endsRegex.FindStringSubmatch(m.Embeds[0].Description)

		if len(hostMatches) == 0 || len(endsMatches) == 0 {
			break
		}

		hostID, ends := hostMatches[1], endsMatches[1]

		// Parse UNIX timestamp
		endsUnix, _ := strconv.ParseInt(ends, 10, 64)
		endsTime := time.Unix(endsUnix, 0)
		endsFormat := utils.Misc.FormatDuration(time.Until(endsTime))

		// If the ends time is zero, the message is not a giveaway
		if endsUnix == 0 {
			break
		}

		logging.Logging.GiveawayFound(session, m.Message, prize, hostID, endsFormat, s.Account)

		go func() {
			time.Sleep(s.giveawaySniperSleepDuration(&cfg, &endsTime))

			// Click the enter button
			err = session.InteractionClick(button, m.Message)
			if err != nil {
				utils.Logging.Debug("Snipers.CheckGiveawaySniper: Error clicking button: ", err)
			}
		}()
	case "700743797977514004": // Friskytool#6219
		// Win check
		winRegex := regexp.MustCompile(`<@!?(\d+)> has won the giveaway for (.*?)\n`)

		wins := winRegex.FindAllStringSubmatch(m.Content, -1)
		for _, win := range wins {
			winnerID, prize := win[1], win[2]

			if winnerID != session.State.User.ID {
				continue
			}

			logging.Logging.GiveawayWon(session, m.Message, prize, s.Account)
		}

		// Check if the message is a giveaway
		if len(m.Embeds) == 0 || !strings.Contains(m.Embeds[0].Description, "React with") {
			break
		}

		// Get giveaway info
		prize := m.Embeds[0].Title
		if s.isBlacklisted(&cfg, prize) {
			break
		}

		hostRegex := regexp.MustCompile(`Hosted by <@!?(\d+)>`)
		hostMatches := hostRegex.FindStringSubmatch(m.Embeds[0].Description)

		endsTime, err := time.Parse(time.RFC3339, m.Embeds[0].Timestamp)
		if err != nil || len(hostMatches) == 0 {
			break
		}

		hostID := hostMatches[1]
		endsFormat := utils.Misc.FormatDuration(time.Until(endsTime))

		logging.Logging.GiveawayFound(session, m.Message, prize, hostID, endsFormat, s.Account)

		go func() {
			time.Sleep(s.giveawaySniperSleepDuration(&cfg, &endsTime))

			// Add the enter reaction
			err = session.MessageReactionAdd(m.ChannelID, m.ID, "🎉")
			if err != nil {
				utils.Logging.Debug("Snipers.CheckGiveawaySniper: Error adding reaction: ", err)
			}
		}()
	case "530082442967646230": // Giveaway Boat#2911
		if len(m.Embeds) == 0 {
			break
		}

		// Win check
		winRegex := regexp.MustCompile(`<@!?(\d+)> won the giveaway of \[(.*?)\]`)

		wins := winRegex.FindAllStringSubmatch(m.Embeds[0].Description, -1)
		for _, win := range wins {
			winnerID, prize := win[1], win[2]

			if winnerID != session.State.User.ID {
				continue
			}

			logging.Logging.GiveawayWon(session, m.Message, prize, s.Account)
		}

		// Check if the message is a giveaway
		button, ok := s.getButtonComponent(m.Message, 0)
		if !strings.Contains(m.Embeds[0].Description, "button to enter!") || !ok {
			break
		}

		// Get giveaway info
		prize := m.Embeds[0].Title
		if s.isBlacklisted(&cfg, prize) {
			break
		}

		hostRegex := regexp.MustCompile(`Hosted by: <@!?(\d+)>`)
		hostMatches := hostRegex.FindStringSubmatch(m.Embeds[0].Description)

		hostID := ""
		if len(hostMatches) > 0 {
			hostID = hostMatches[1]
		}

		endsTime, err := time.Parse(time.RFC3339, m.Embeds[0].Timestamp)
		if err != nil {
			break
		}

		endsFormat := utils.Misc.FormatDuration(time.Until(endsTime))

		logging.Logging.GiveawayFound(session, m.Message, prize, hostID, endsFormat, s.Account)

		go func() {
			time.Sleep(s.giveawaySniperSleepDuration(&cfg, &endsTime))

			// Click the enter button
			err = session.InteractionClick(button, m.Message)
			if err != nil {
				utils.Logging.Debug("Snipers.CheckGiveawaySniper: Error clicking button: ", err)
			}
		}()
	case "239631525350604801": // Pancake#3691
		// Win check
		if utils.Discord.IsUserMentioned(m.Message, session.State.User.ID) {
			winRegex := regexp.MustCompile(`You won \*\*(.*?)\*\*!`)

			wins := winRegex.FindStringSubmatch(m.Content)
			if len(wins) > 0 {
				logging.Logging.GiveawayWon(session, m.Message, wins[1], s.Account)
			}
		}

		// Check if the message is a giveaway
		if len(m.Embeds) == 0 || !strings.Contains(m.Embeds[0].Description, "React with") {
			break
		}

		// Get giveaway info
		prize := strings.ReplaceAll(m.Embeds[0].Title, "🎉", "")
		prize = strings.TrimPrefix(prize, " ")
		prize = strings.TrimSuffix(prize, " ")
		if s.isBlacklisted(&cfg, prize) {
			break
		}

		endsRegex := regexp.MustCompile(`Event ends <t:(\d+):R>`)
		endsMatches := endsRegex.FindStringSubmatch(m.Embeds[0].Description)

		if len(endsMatches) == 0 {
			break
		}

		// Parse UNIX timestamp
		endsUnix, _ := strconv.ParseInt(endsMatches[1], 10, 64)
		endsTime := time.Unix(endsUnix, 0)
		endsFormat := utils.Misc.FormatDuration(time.Until(endsTime))

		logging.Logging.GiveawayFound(session, m.Message, prize, "", endsFormat, s.Account)

		go func() {
			time.Sleep(s.giveawaySniperSleepDuration(&cfg, &endsTime))

			// Add the enter reaction
			err = session.MessageReactionAdd(m.ChannelID, m.ID, "🎁")
			if err != nil {
				utils.Logging.Debug("Snipers.CheckGiveawaySniper: Error adding reaction: ", err)
			}
		}()
	case "651627261782523926": // Giveaways#8717
		// Win check
		winRegex := regexp.MustCompile(`Congratulations <@!?(\d+)>, you won \*\*(.*?)\*\*`)

		wins := winRegex.FindAllStringSubmatch(m.Content, -1)
		for _, win := range wins {
			winnerID, prize := win[1], win[2]

			if winnerID != session.State.User.ID {
				continue
			}

			logging.Logging.GiveawayWon(session, m.Message, prize, s.Account)
		}

		// Check if the message is a giveaway
		if len(m.Embeds) == 0 ||
			!strings.Contains(m.Embeds[0].Description, "Click the reaction below to enter!") {
			break
		}

		// Get giveaway info
		prize := m.Embeds[0].Title
		if s.isBlacklisted(&cfg, prize) {
			break
		}

		hostRegex := regexp.MustCompile(`Host: <@!?(\d+)>`)
		hostMatches := hostRegex.FindStringSubmatch(m.Embeds[0].Description)

		endsRegex := regexp.MustCompile(`Ends <t:(\d+):R>`)
		endsMatches := endsRegex.FindStringSubmatch(m.Embeds[0].Description)

		if len(hostMatches) == 0 || len(endsMatches) == 0 {
			break
		}

		hostID, ends := hostMatches[1], endsMatches[1]

		// Parse UNIX timestamp
		endsUnix, _ := strconv.ParseInt(ends, 10, 64)
		endsTime := time.Unix(endsUnix, 0)
		endsFormat := utils.Misc.FormatDuration(time.Until(endsTime))

		logging.Logging.GiveawayFound(session, m.Message, prize, hostID, endsFormat, s.Account)

		go func() {
			time.Sleep(s.giveawaySniperSleepDuration(&cfg, &endsTime))

			// Add the enter reaction
			err = session.MessageReactionAdd(m.ChannelID, m.ID, "🎉")
			if err != nil {
				utils.Logging.Debug("Snipers.CheckGiveawaySniper: Error adding reaction: ", err)
			}
		}()
	case "368521195940741122": // Lawliet#5480
		if len(m.Embeds) == 0 {
			break
		}

		// Check if the message is a giveaway
		if len(m.Embeds) == 0 || !strings.Contains(m.Embeds[0].Description, "React with") {
			break
		}

		// Get giveaway info
		prize := strings.ReplaceAll(m.Embeds[0].Title, "🎆 ", "")
		if s.isBlacklisted(&cfg, prize) {
			break
		}

		emojiRegex := regexp.MustCompile(`React with (.*?) to enter`)
		emojiMatches := emojiRegex.FindStringSubmatch(m.Embeds[0].Description)

		endsRegex := regexp.MustCompile(`<t:(\d+):R>`)
		endsMatches := endsRegex.FindStringSubmatch(m.Embeds[0].Description)

		if len(emojiMatches) == 0 || len(endsMatches) == 0 {
			break
		}

		emoji, ends := emojiMatches[1], endsMatches[1]

		// Parse UNIX timestamp
		endsUnix, _ := strconv.ParseInt(ends, 10, 64)
		endsTime := time.Unix(endsUnix, 0)
		endsFormat := utils.Misc.FormatDuration(time.Until(endsTime))

		logging.Logging.GiveawayFound(session, m.Message, prize, "", endsFormat, s.Account)

		go func() {
			time.Sleep(s.giveawaySniperSleepDuration(&cfg, &endsTime))

			// Add the enter reaction
			err = session.MessageReactionAdd(m.ChannelID, m.ID, emoji)
			if err != nil {
				utils.Logging.Debug("Snipers.CheckGiveawaySniper: Error adding reaction: ", err)
			}
		}()
	case "159985870458322944": // MEE6#4876
		if len(m.Embeds) == 0 {
			break
		}

		// Win check
		if strings.Contains(m.Embeds[0].Title, "Congratulations") &&
			utils.Discord.IsUserMentioned(m.Message, session.State.User.ID) {
			prize := ""

			if len(m.Embeds[0].Fields) > 0 {
				prize = m.Embeds[0].Fields[0].Value
			}

			logging.Logging.GiveawayWon(session, m.Message, prize, s.Account)
		}

		enterButton, enterOK := s.getButtonComponent(m.Message, 0)
		timeButton, timeOK := s.getButtonComponent(m.Message, 1)
		if !enterOK || !timeOK {
			break
		}

		// Check if the message is a giveaway
		if !strings.Contains(enterButton.Label, "Participate") {
			break
		}

		// Get giveaway info
		timezoneRegex := regexp.MustCompile(` \((.*?)\)`)
		timezoneMatches := timezoneRegex.FindStringSubmatch(timeButton.Label)

		if len(timezoneMatches) == 0 {
			break
		}

		timeButton.Label = timezoneRegex.ReplaceAllString(timeButton.Label, "")

		endsIn := ""
		var giveawayEnd *time.Time

		timezone, err := time.LoadLocation(timezoneMatches[1])
		if err == nil {
			endsTime, err := time.ParseInLocation("Ends on Mon _2 Jan 2006 at 15:04", timeButton.Label, timezone)
			if err == nil {
				giveawayEnd = &endsTime
				endsIn = utils.Misc.FormatDuration(time.Until(endsTime))
			}
		}

		logging.Logging.GiveawayFound(session, m.Message, "MEE6", "", endsIn, s.Account)

		go func() {
			time.Sleep(s.giveawaySniperSleepDuration(&cfg, giveawayEnd))

			// Click the enter button
			err = session.InteractionClick(enterButton, m.Message)
			if err != nil {
				utils.Logging.Debug("Snipers.CheckGiveawaySniper: Error clicking button: ", err)
			}
		}()
	case "824119071556763668": // Apollo#2300
		// Win check
		winRegex := regexp.MustCompile(`Congratulations, <@!?(\d+)>! You won \*\*(.*?)\*\*!`)

		wins := winRegex.FindAllStringSubmatch(m.Content, -1)
		for _, win := range wins {
			winnerID, prize := win[1], win[2]

			if winnerID != session.State.User.ID {
				continue
			}

			logging.Logging.GiveawayWon(session, m.Message, prize, s.Account)
		}

		// Check if the message is a giveaway
		if len(m.Embeds) == 0 || !strings.Contains(m.Content, "GIVEAWAY!") {
			break
		}

		// Get giveaway info
		prize := ""
		if m.Embeds[0].Author != nil {
			prize = m.Embeds[0].Author.Name
		}

		if s.isBlacklisted(&cfg, prize) {
			break
		}

		hostRegex := regexp.MustCompile(`Hosted by <@!?(\d+)>`)
		hostMatches := hostRegex.FindStringSubmatch(m.Embeds[0].Description)

		endsRegex := regexp.MustCompile(`Ends <t:(\d+):R>`)
		endsMatches := endsRegex.FindStringSubmatch(m.Embeds[0].Description)

		if len(hostMatches) == 0 || len(endsMatches) == 0 {
			break
		}

		hostID, ends := hostMatches[1], endsMatches[1]

		// Parse UNIX timestamp
		endsUnix, _ := strconv.ParseInt(ends, 10, 64)
		endsTime := time.Unix(endsUnix, 0)
		endsFormat := utils.Misc.FormatDuration(time.Until(endsTime))

		logging.Logging.GiveawayFound(session, m.Message, prize, hostID, endsFormat, s.Account)

		go func() {
			time.Sleep(s.giveawaySniperSleepDuration(&cfg, &endsTime))

			// Add the enter reaction
			err = session.MessageReactionAdd(m.ChannelID, m.ID, "🎉")
			if err != nil {
				utils.Logging.Debug("Snipers.CheckGiveawaySniper: Error adding reaction: ", err)
			}
		}()
	case "235148962103951360": // Carl-bot#1536
		// Win check
		if utils.Discord.IsUserMentioned(m.Message, session.State.User.ID) {
			winRegex := regexp.MustCompile(`Congratulations to the winner of \*\*(.*?)\*\*`)

			wins := winRegex.FindAllStringSubmatch(m.Content, -1)
			for _, win := range wins {
				prize := win[1]

				logging.Logging.GiveawayWon(session, m.Message, prize, s.Account)
			}
		}

		// Check if the message is a giveaway
		button, ok := s.getButtonComponent(m.Message, 0)
		if len(m.Embeds) == 0 || len(m.Embeds[0].Fields) < 2 || !ok ||
			!strings.Contains(m.Embeds[0].Description, "join the giveaway") {
			break
		}

		// Get giveaway info
		prize := m.Embeds[0].Title
		if s.isBlacklisted(&cfg, prize) {
			break
		}

		endsRegex := regexp.MustCompile(`<t:(\d+):R>`)
		endsMatches := endsRegex.FindStringSubmatch(m.Embeds[0].Fields[0].Value)

		hostRegex := regexp.MustCompile(`<@!?(\d+)>`)
		hostMatches := hostRegex.FindStringSubmatch(m.Embeds[0].Fields[1].Value)

		if len(endsMatches) == 0 || len(hostMatches) == 0 {
			break
		}

		hostID, ends := hostMatches[1], endsMatches[1]

		// Parse UNIX timestamp
		endsUnix, _ := strconv.ParseInt(ends, 10, 64)
		endsTime := time.Unix(endsUnix, 0)
		endsFormat := utils.Misc.FormatDuration(time.Until(endsTime))

		logging.Logging.GiveawayFound(session, m.Message, prize, hostID, endsFormat, s.Account)

		go func() {
			time.Sleep(s.giveawaySniperSleepDuration(&cfg, &endsTime))

			// Click the enter button
			err = session.InteractionClick(button, m.Message)
			if err != nil {
				utils.Logging.Debug("Snipers.CheckGiveawaySniper: Error clicking button: ", err)
			}
		}()
	case "720351927581278219": // Invite Tracker#0478
		// Win check
		if utils.Discord.IsUserMentioned(m.Message, session.State.User.ID) {
			winRegex := regexp.MustCompile(`you won the \*\*(.*?)\*\*`)

			wins := winRegex.FindAllStringSubmatch(m.Content, -1)
			for _, win := range wins {
				prize := win[1]

				logging.Logging.GiveawayWon(session, m.Message, prize, s.Account)
			}
		}

		// ...

		// Check if the message is a giveaway
		button, ok := s.getButtonComponent(m.Message, 0)
		if len(m.Embeds) == 0 || !ok || !strings.Contains(button.Label, "Entry") {
			break
		}

		// Get giveaway info
		prize := m.Embeds[0].Title
		if s.isBlacklisted(&cfg, prize) {
			break
		}

		endsRegex := regexp.MustCompile(`<t:(\d+):R>`)
		endsMatches := endsRegex.FindStringSubmatch(m.Embeds[0].Description)

		hostRegex := regexp.MustCompile(`<@!?(\d+)>`)
		hostMatches := hostRegex.FindStringSubmatch(m.Embeds[0].Description)

		if len(endsMatches) == 0 || len(hostMatches) == 0 {
			break
		}

		hostID, ends := hostMatches[1], endsMatches[1]

		// Parse UNIX timestamp
		endsUnix, _ := strconv.ParseInt(ends, 10, 64)
		endsTime := time.Unix(endsUnix, 0)
		endsFormat := utils.Misc.FormatDuration(time.Until(endsTime))

		logging.Logging.GiveawayFound(session, m.Message, prize, hostID, endsFormat, s.Account)

		go func() {
			time.Sleep(s.giveawaySniperSleepDuration(&cfg, &endsTime))

			// Click the enter button
			err = session.InteractionClick(button, m.Message)
			if err != nil {
				utils.Logging.Debug("Snipers.CheckGiveawaySniper: Error clicking button: ", err)
			}
		}()
	default: // Custom bots
		if !utils.Discord.IsUserInList(session, m.Author.ID, cfg.Snipers.GiveawaySniper.CustomBots) {
			break
		}

		rawMessage, err := json.Marshal(m)
		if err != nil {
			utils.Logging.Debug("Snipers.CheckGiveawaySniper: Error marshalling message: ", err)
			break
		}

		rawMessageString := strings.ToLower(string(rawMessage))

		// Check if the message is a giveaway
		isGiveaway := strings.Contains(rawMessageString, "giveaway") ||
			strings.Contains(rawMessageString, "react with") ||
			strings.Contains(rawMessageString, "to enter") ||
			strings.Contains(rawMessageString, "to join")
		if !isGiveaway {
			break
		}

		// Get giveaway info
		hostRegex := regexp.MustCompile(`\\u003c@!?(.*?)\\u003e`)
		hostMatches := hostRegex.FindStringSubmatch(rawMessageString)

		endsRegex := regexp.MustCompile(`\\u003ct:(.*?):R\\u003e`)
		endsMatches := endsRegex.FindStringSubmatch(rawMessageString)

		hostID := ""
		endsIn := ""
		var giveawayEnd *time.Time

		if len(hostMatches) > 1 {
			hostID = hostMatches[1]
		}

		if len(endsMatches) > 1 {
			// Parse UNIX timestamp
			endsUnix, _ := strconv.ParseInt(endsMatches[1], 10, 64)
			endsTime := time.Unix(endsUnix, 0)
			giveawayEnd = &endsTime
			endsIn = utils.Misc.FormatDuration(time.Until(endsTime))
		}

		logging.Logging.GiveawayFound(session, m.Message, "Unknown (Custom Bot)", hostID, endsIn, s.Account)

		go func() {
			time.Sleep(s.giveawaySniperSleepDuration(&cfg, giveawayEnd))

			button, ok := s.getButtonComponent(m.Message, 0)
			if ok {
				// Click the enter button
				err = session.InteractionClick(button, m.Message)
				if err != nil {
					utils.Logging.Debug("Snipers.CheckGiveawaySniper: Error clicking button: ", err)
				}

				return
			}

			// Reload the message to get the updated reactions
			newMessage, err := utils.Discord.GetChannelMessage(session, m.ChannelID, m.ID)
			if err != nil {
				utils.Logging.Debug("Snipers.CheckGiveawaySniper: Error reloading message: ", err)
				return
			}

			if len(newMessage.Reactions) == 0 {
				return
			}

			// Add the enter reaction
			var emoji string

			reaction := newMessage.Reactions[0]
			if reaction.Emoji != nil {
				emoji = reaction.Emoji.APIName()
			}

			if emoji == "" {
				return
			}

			// Add the enter reaction
			err = session.MessageReactionAdd(m.ChannelID, m.ID, emoji)
			if err != nil {
				utils.Logging.Debug("Snipers.CheckGiveawaySniper: Error adding reaction: ", err)
			}
		}()
	}
}

func (s *Snipers) CheckGiveawaySniperEdit(session *discordgo.Session, m *discordgo.MessageUpdate) {
	cfg, err := config.GetConfig()
	if err != nil {
		utils.Logging.Debug("Snipers.CheckGiveawaySniper: Error getting config: ", err)
		return
	}

	if !cfg.Snipers.GiveawaySniper.Enabled {
		return
	}

	switch m.Author.ID {
	case "693167035068317736": // Rumble Royale#6012
		// Win check
		if len(m.Embeds) >= 2 && strings.Contains(m.Embeds[0].Title, "WINNER") &&
			utils.Discord.IsUserMentioned(m.Message, session.State.User.ID) {
			logging.Logging.GiveawayWon(session, m.Message, "Rumble Royale", s.Account)
		}

		// Check if the message is a giveaway
		if len(m.Embeds) == 0 || !strings.Contains(m.Embeds[0].Title, "Rumble Royale") {
			break
		}

		// Bot uses invisible characters to prevent sniping (U+200E)
		m.Embeds[0].Description = strings.Map(func(r rune) rune {
			if unicode.IsGraphic(r) {
				return r
			}

			return -1
		}, m.Embeds[0].Description)

		if !strings.Contains(m.Embeds[0].Description, "Click the emoji below to join.") {
			break
		}

		// Get giveaway info
		hostRegex := regexp.MustCompile(`hosted by (.*?)$`)
		hostMatches := hostRegex.FindStringSubmatch(m.Embeds[0].Title)

		endsRegex := regexp.MustCompile(`Starting in (.*?)!`)
		endsMatches := endsRegex.FindStringSubmatch(m.Embeds[0].Description)

		if len(hostMatches) == 0 || len(endsMatches) == 0 {
			break
		}

		host, endsIn := hostMatches[1], endsMatches[1]

		logging.Logging.GiveawayFound(session, m.Message, "Rumble Royale", host, endsIn, s.Account)

		go func() {
			// Wait until the bot has reacted
			time.Sleep(time.Duration(random.Int(5, 30)) * time.Second)

			// Reload the message to get the updated reactions
			newMessage, err := utils.Discord.GetChannelMessage(session, m.ChannelID, m.ID)
			if err != nil {
				utils.Logging.Debug("Snipers.CheckGiveawaySniper: Error reloading message: ", err)
				return
			}

			if len(newMessage.Reactions) == 0 {
				return
			}

			// Add the enter reaction
			var emoji string

			reaction := newMessage.Reactions[0]
			if reaction.Emoji != nil {
				emoji = reaction.Emoji.APIName()
			}

			if emoji == "" {
				return
			}

			// Add the enter reaction
			err = session.MessageReactionAdd(m.ChannelID, m.ID, emoji)
			if err != nil {
				utils.Logging.Debug("Snipers.CheckGiveawaySniper: Error adding reaction: ", err)
			}
		}()
	case "676039466175168543": // ZAP-Hosting Bot#7133
		// Win check
		if strings.Contains(m.Content, "Giveaway Winner:") &&
			utils.Discord.IsUserMentioned(m.Message, session.State.User.ID) {
			logging.Logging.GiveawayWon(session, m.Message, "€10 worth of ZAP Coins", s.Account)
		}
	case "368521195940741122": // Lawliet#5480
		// Win check
		if len(m.Embeds) > 0 && strings.Contains(m.Embeds[0].Description, "You just won in the giveaway!") &&
			utils.Discord.IsUserMentioned(m.Message, session.State.User.ID) {
			prize := m.Embeds[0].Title
			if len(prize) < 2 {
				break
			}

			prize = prize[2:]

			logging.Logging.GiveawayWon(session, m.Message, prize, s.Account)
		}
	}
}

func (s *Snipers) CheckPokeTwoSniper(session *discordgo.Session, m *discordgo.MessageCreate) {
	cfg, err := config.GetConfig()
	if err != nil {
		utils.Logging.Debug("Snipers.CheckPokeTwoSniper: Error getting config: ", err)
		return
	}

	isPokeTwo := cfg.Snipers.PokeTwoSniper.Enabled && // Sniper is enabled
		m.Author.ID == "716390085896962058" && // Message is from PokeTwo
		slices.Contains(cfg.Snipers.PokeTwoSniper.Channels, m.ChannelID) // Message is in a sniper channel

	// Captcha check
	captchaRegex := regexp.MustCompile(`Whoa there. Please tell us you're human! https://verify.poketwo.net/captcha/(\d+)`)
	captchaMatches := captchaRegex.FindStringSubmatch(m.Content)
	if len(captchaMatches) > 1 && captchaMatches[1] == session.State.User.ID {
		s.PokeTwoCaptchaRequired = true

		logging.UI.Info("PokéTwo Sniper: Captcha required!")

		err = api.PokeTwoVerify(session.State.User.ID)
		if err != nil {
			utils.Logging.Debug("Snipers.CheckPokeTwoSniper: Error solving captcha: ", err)

			logging.UI.Info("PokéTwo Sniper: Error solving captcha! You need to solve it manually at https://verify.poketwo.net/captcha/" + session.State.User.ID)

			// Try again after 5 minutes
			go func() {
				time.Sleep(time.Minute * 5)
				s.PokeTwoCaptchaRequired = false
			}()
		} else {
			logging.UI.Info("PokéTwo Sniper: Captcha completed!")

			s.PokeTwoCaptchaRequired = false
		}

		return
	}

	sniperConditions := isPokeTwo && len(m.Embeds) == 1

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
	if currentCatchRate > cfg.Snipers.PokeTwoSniper.CatchRate {
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

	image := utils.Misc.ImageToBase64(resp.Body)

	// Detect the Pokémon
	beforeRequest := time.Now()
	pokemon, confidence, err := api.PokeTwo(image, &cfg)
	if err != nil {
		utils.Logging.Debug("Snipers.CheckPokeTwoSniper: Error getting Pokémon: ", err)
		return
	}
	elapsed := time.Since(beforeRequest)

	// Calculate the sleep time
	sleep := time.Duration(cfg.Snipers.PokeTwoSniper.Sleep) * time.Second
	if cfg.Snipers.PokeTwoSniper.RandomSleep {
		// Randomize the sleep time
		sleep = time.Duration(float64(sleep) * (0.2 + 0.5*rand.Float64()))
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
	label := pokemon.NameEnglish
	switch cfg.Snipers.PokeTwoSniper.CatchLanguage {
	case "Japanese":
		label = pokemon.NameJapanese
	case "JapaneseTransliterated":
		label = pokemon.NameJapaneseT
	case "German":
		label = pokemon.NameGerman
	case "French":
		label = pokemon.NameFrench
	}

	if label == "" {
		label = pokemon.NameEnglish
	}

	label = strings.ToLower(label)

	_, err = session.ChannelMessageSend(m.ChannelID, "<@716390085896962058> c "+strings.ToLower(label))
	if err != nil {
		utils.Logging.Debug("Snipers.CheckPokeTwoSniper: Error catching Pokémon: ", err)
		return
	}

	// Wait for the Pokémon to be caught
	select {
	case m := <-congratsMessage:
		if !utils.Discord.IsUserMentioned(m, session.State.User.ID) {
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
			iv       string
		)

		pokemonMatches := regexp.MustCompile(`You caught a (.*?)<`).FindStringSubmatch(m.Content)
		if len(pokemonMatches) > 1 {
			pokemon = pokemonMatches[1]
		}

		receivedMatches := regexp.MustCompile(`You received (.*?)!`).FindStringSubmatch(m.Content)
		if len(receivedMatches) > 1 {
			received = receivedMatches[1]
		}

		ivMatches := regexp.MustCompile(`\((.*?)\)`).FindStringSubmatch(m.Content)
		if len(ivMatches) > 1 {
			iv = ivMatches[1]
		}

		logging.Logging.PokeTwoSniper(session, m, pokemon, received, iv, confidence, elapsed, s.Account)

		// Fake activity

		go func() {
			if !cfg.Snipers.PokeTwoSniper.FakeActivity {
				return
			}

			// Wait for random delay --> type --> wait for random delay
			time.Sleep(time.Second * time.Duration(random.Int(1, 2)))
			_ = session.ChannelTyping(m.ChannelID)
			time.Sleep(time.Second * time.Duration(random.Int(1, 3)))

			// Random PokéTwo commands
			commands := []string{
				"ping", "ev mg", "sh", "pokedex", "s latest", "pokedex",
				"d " + strings.ToLower(label),
			}

			_, _ = session.ChannelMessageSend(
				m.ChannelID,
				"<@716390085896962058> "+random.Choice(commands),
			)
		}()

	case <-time.After(5 * time.Second):
		utils.Logging.Debug("Snipers.CheckPokeTwoSniper: Pokémon not caught in time")
	}

	// Remove the handler
	remove()
}
