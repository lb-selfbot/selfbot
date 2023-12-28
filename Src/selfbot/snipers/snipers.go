package snipers

import (
	"encoding/json"
	"fmt"
	"go/token"
	"go/types"
	"path/filepath"
	"regexp"
	"strconv"
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

type NitroCode struct {
	Code string
	Time time.Time
}

var nitroRegex = regexp.MustCompile(`(discord\.gift|discord\.com/gifts|discordapp\.com/gifts|promos\.discord\.gg/)/(\w{0,30})`)
var privnoteRegex = regexp.MustCompile(`privnote.com/([a-zA-Z0-9_#]{18,18})`)

func getCodesInLast30Seconds(codes []NitroCode) []NitroCode {
	return utils.Filter(codes, func(code NitroCode) bool {
		return time.Since(code.Time) < 30*time.Second
	})
}

func CheckNitroSniper(
	session *discordgo.Session,
	message *discordgo.Message,
	cfg config.Config,
	countryCode string,
	paymentSourceID string,
	alreadyCheckedNitroCodes []string,
	sentNitroCodes map[string][]NitroCode,
	account string,
) (
	[]string,
	map[string][]NitroCode,
) {
	if !cfg.Snipers.NitroSniper {
		return alreadyCheckedNitroCodes, sentNitroCodes
	}

	if utils.SliceContains(cfg.Snipers.ServerBlacklist, message.GuildID) {
		return alreadyCheckedNitroCodes, sentNitroCodes
	}

	nitroLinks := nitroRegex.FindAllString(message.Content, -1)

	nitroLinks = append(nitroLinks, message.GiftCodes...)

	for _, nitroLink := range nitroLinks {
		nitroLinkSplit := strings.Split(nitroLink, "/")

		nitroCode := nitroLinkSplit[len(nitroLinkSplit)-1]

		// Code already checked
		if utils.SliceContains(alreadyCheckedNitroCodes, nitroCode) {
			return alreadyCheckedNitroCodes, sentNitroCodes
		}

		alreadyCheckedNitroCodes = append(alreadyCheckedNitroCodes, nitroCode)

		// Code length is invalid
		if len(nitroCode) < 15 || len(nitroCode) > 26 {
			logging.Logging.NitroSniperFake(session, message, nitroCode, account)
			return alreadyCheckedNitroCodes, sentNitroCodes
		}

		// Author is spamming codes
		if codes, ok := sentNitroCodes[message.Author.ID]; ok {
			if len(getCodesInLast30Seconds(codes)) >= 5 {
				return alreadyCheckedNitroCodes, sentNitroCodes
			}
		}

		url := fmt.Sprintf("https://discord.com/api/v9/entitlements/gift-codes/%s/redeem", nitroCode)

		payload := map[string]any{
			"with_application":       false,
			"with_subscription_plan": true,
		}

		if countryCode != "" {
			payload["country_code"] = countryCode
		}

		if paymentSourceID != "" {
			payload["payment_source_id"] = paymentSourceID
		}

		start := time.Now()

		response, _ := session.Request("POST", url, payload)

		elapsed := time.Since(start)

		redeemed := false

		responseString := string(response)

		if strings.Contains(responseString, "This gift has been redeemed already.") {
			logging.Logging.NitroSniper(session, message, elapsed, "Already Redeemed", nitroCode, account)
		} else if strings.Contains(responseString, "Already purchased") {
			logging.Logging.NitroSniper(session, message, elapsed, "Already Purchased", nitroCode, account)
		} else if strings.Contains(responseString, "subscription_plan") {
			logging.Logging.NitroSniper(session, message, elapsed, "Redeemed Nitro", nitroCode, account)
			redeemed = true
		} else if strings.Contains(responseString, "Unknown Gift Code") {
			logging.Logging.NitroSniper(session, message, elapsed, "Unknown Code", nitroCode, account)
		} else {
			json, err := simplejson.NewJson(response)
			if err != nil {
				json = simplejson.New()
			}

			logging.Logging.NitroSniper(session, message, elapsed, utils.Misc.Or(json.Get("message").MustString(), "Unknown Status"), nitroCode, account)
		}

		// Log the nitro code if it was invalid
		if !redeemed {
			code := NitroCode{
				Code: nitroCode,
				Time: time.Now(),
			}
			_, ok := sentNitroCodes[message.Author.ID]
			if ok {
				sentNitroCodes[message.Author.ID] = append(sentNitroCodes[message.Author.ID], code)
			} else {
				sentNitroCodes[message.Author.ID] = []NitroCode{code}
			}
		}
	}
	return alreadyCheckedNitroCodes, sentNitroCodes
}

func CheckPrivnoteSniper(session *discordgo.Session, message *discordgo.Message, cfg *config.Config, account string) {
	if !cfg.Snipers.PrivnoteSniper {
		return
	}

	if utils.SliceContains(cfg.Snipers.ServerBlacklist, message.GuildID) {
		return
	}

	privnoteCodes := privnoteRegex.FindAllString(message.Content, -1)

	for _, privnoteCode := range privnoteCodes {
		link := "https://" + privnoteCode

		client := privnote.NewClient()

		start := time.Now()

		note, err := client.ReadNoteFromLink(link)
		if err != nil {
			continue
		}

		elapsed := time.Since(start)

		split := strings.Split(link, "/")
		if len(split) < 2 {
			continue
		}

		id := split[len(split)-1]

		fileName := utils.IO.SanitizeFilename("privnote-" + id + ".txt")
		filePath := filepath.Join(utils.IO.GetMainFolder(), fileName)

		err = utils.IO.Write(filePath, note)
		if err != nil {
			continue
		}

		_ = utils.System.OpenURL(filePath)

		logging.Logging.PrivnoteSniper(session, message, elapsed, "Saved Note", filePath, account)
	}
}

var (
	totalPokemon  = 0
	caughtPokemon = 0
	lastCatch     time.Time
)

func CheckPokeTwoSniper(session *discordgo.Session, message *discordgo.Message, cfg *config.Config) {
	if !cfg.Snipers.PokeTwoSniper.Enabled {
		return
	}

	if message.Author.ID != "716390085896962058" {
		return
	}

	if !utils.SliceContains(cfg.Snipers.PokeTwoSniper.Channels, message.ChannelID) {
		return
	}

	if len(message.Embeds) == 0 {
		return
	}

	emb := message.Embeds[0]

	if !strings.Contains(emb.Title, "has appeared!") {
		return
	}

	if emb.Image == nil {
		return
	}

	totalPokemon++

	if time.Since(lastCatch) < time.Second*time.Duration(cfg.Snipers.PokeTwoSniper.Timeout) {
		utils.Logging.Info("Pokémon catch timeout reached! Skipping...")
		return
	}

	// Check if the Pokémon should be caught
	catchRate := int(float64(caughtPokemon) / float64(totalPokemon) * 100)
	if catchRate > cfg.Snipers.PokeTwoSniper.CatchRate {
		utils.Logging.Info("Pokémon catch rate reached! Current catch rate is", catchRate, "%. Skipping...")
		return
	}

	// Start typing
	if cfg.Snipers.PokeTwoSniper.Typing {
		go session.ChannelTyping(message.ChannelID)
	}

	// Check if the Pokémon is already caught
	congratulationMessage := make(chan *discordgo.Message)

	remove := session.AddHandler(func(_ *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.ID != "716390085896962058" {
			return
		}

		if m.ChannelID != message.ChannelID {
			return
		}

		if strings.Contains(m.Content, "Congratulations") {
			congratulationMessage <- m.Message
		}
	})

	// Detect the Pokémon
	httpResponse := http.Get(emb.Image.URL)
	if httpResponse.Error != nil {
		return
	}

	image := utils.Misc.ImageToBase64(httpResponse.BodyBytes)

	before := time.Now()
	label, confidence, err := api.PokeTwo(image, cfg)
	if err != nil {
		utils.Logging.Error("Error while trying to detect pokemon: ", err)
		return
	}
	elapsed := time.Since(before)

	// Wait for the cooldown
	var sleep time.Duration
	if cfg.Snipers.PokeTwoSniper.RandomSleep {
		sleep = time.Second * time.Duration(utils.Random.Number(cfg.Snipers.PokeTwoSniper.Sleep/2, cfg.Snipers.PokeTwoSniper.Sleep))
	} else {
		sleep = time.Second * time.Duration(cfg.Snipers.PokeTwoSniper.Sleep)
	}

	sleep = sleep - elapsed
	if sleep > 0 {
		time.Sleep(sleep)
	}

	// Check if the Pokémon is already caught
	select {
	case <-congratulationMessage:
		utils.Logging.Error("Pokémon already caught!")
		remove()
		return
	default:
	}

	// Catch the Pokémon
	_, err = session.ChannelMessageSend(
		message.ChannelID,
		fmt.Sprintf("%s c %s", message.Author.Mention(), strings.ToLower(label)),
	)
	if err != nil {
		return
	}

	// Wait for the Pokémon to be caught
	select {
	case m := <-congratulationMessage:
		mentionsSelf := utils.Filter(m.Mentions, func(mention *discordgo.User) bool {
			return mention.ID == session.State.User.ID
		})
		if len(mentionsSelf) == 0 {
			utils.Logging.Error("Pokémon caught by someone else!")
			return
		}

		caughtPokemon++
		lastCatch = time.Now()

		// Remove the mention
		_ = session.MessageAck(message.ChannelID, message.ID)

		// Fake activity
		if cfg.Snipers.PokeTwoSniper.FakeActivity && caughtPokemon%5 == 0 {
			go func() {
				// Wait for random delay --> type --> wait for random delay
				time.Sleep(time.Second * time.Duration(utils.Random.Number(1, 2)))
				_ = session.ChannelTyping(message.ChannelID)
				time.Sleep(time.Second * time.Duration(utils.Random.Number(1, 3)))

				commands := []string{
					"ping", "ev mg", "sh", "pokedex", "s latest", "pokedex",
					"d " + strings.ToLower(label),
				}

				_, _ = session.ChannelMessageSend(
					message.ChannelID,
					fmt.Sprintf("%s %s", message.Author.Mention(), utils.RandomChoice(commands)),
				)
			}()
		}

		var (
			pokemon  string
			received string
		)

		pokemonMatches := regexp.MustCompile(`You caught a (.*?)!`).FindStringSubmatch(m.Content)
		if len(pokemonMatches) > 1 {
			pokemon = pokemonMatches[1]
		}

		receivedMatches := regexp.MustCompile(`You received (.*?)!`).FindStringSubmatch(m.Content)
		if len(receivedMatches) > 1 {
			received = receivedMatches[1]
		}

		logging.Logging.PokeTwoSniper(session, message, pokemon, received, confidence, elapsed)
	case <-time.After(time.Second * 5):
		utils.Logging.Error("PokéTwo sniper timed out!")
	}

	remove()
}

func isPrizeBlacklisted(prize string, cfg *config.Config) bool {
	for _, word := range cfg.Snipers.Blacklist {
		if strings.Contains(strings.ToLower(prize), strings.ToLower(word)) {
			return true
		}
	}
	return false
}

func CheckGiveawaySniper(session *discordgo.Session, message *discordgo.Message, cfg *config.Config, account string) {
	if !cfg.Snipers.GiveawaySniper.Enabled {
		return
	}

	if utils.SliceContains(cfg.Snipers.ServerBlacklist, message.GuildID) {
		return
	}

	switch message.Author.ID {
	case "676039466175168543": // ZAP-Hosting Bot#7133
		if !strings.Contains(message.Content, "It's time for a Giveaway!") {
			break
		}

		if len(message.Components) == 0 {
			return
		}

		actionsRow, ok := message.Components[0].(*discordgo.ActionsRow)
		if !ok {
			break
		}

		if len(actionsRow.Components) == 0 {
			return
		}

		button, ok := actionsRow.Components[0].(*discordgo.Button)
		if !ok {
			break
		}

		session.AddHandlerOnce(func(_ *discordgo.Session, modal *discordgo.InteractionModalCreate) {
			actionsRow, ok := modal.Components[0].(*discordgo.ActionsRow)
			if !ok {
				return
			}
			textInput, ok := actionsRow.Components[0].(*discordgo.TextInput)
			if !ok {
				return
			}

			equation := strings.Split(textInput.Label, "\n")[1]
			equation = strings.ReplaceAll(equation, "=", "")
			equation = strings.ReplaceAll(equation, " ", "")

			tv, err := types.Eval(token.NewFileSet(), nil, token.NoPos, equation)
			if err != nil {
				return
			}

			equationSolution := tv.Value.String()

			time.Sleep(5 * time.Second)

			err = session.InteractionSubmitModal(message, modal, textInput.CustomID, equationSolution)
			if err != nil {
				return
			}
		})

		logging.Logging.GiveawayFound(session, message, "€10 worth of ZAP Coins", "ZAP-Hosting", "1 hour", account)

		go func() {
			time.Sleep(time.Second * time.Duration(cfg.Snipers.GiveawaySniper.Sleep))
			_ = session.InteractionClick(button, message)
		}()
	case "294882584201003009", "1041302711750119504": // GiveawayBot#2381, LightningBot#1992
		// Win check

		if message.Type == discordgo.MessageTypeReply {
			prizeList := regexp.MustCompile(fmt.Sprintf(`Congratulations %s! You won the \*\*(.*?)\*\*!`, session.State.User.Mention())).FindStringSubmatch(message.Content)
			if len(prizeList) > 1 {
				logging.Logging.GiveawayWon(session, message, prizeList[1], account)
			}
		}

		// Entering

		if len(message.Embeds) == 0 {
			break
		}

		if len(message.Components) == 0 {
			return
		}

		component, ok := message.Components[0].(*discordgo.ActionsRow)
		if !ok {
			break
		}

		if len(component.Components) == 0 {
			return
		}

		button, ok := component.Components[0].(*discordgo.Button)
		if !ok {
			break
		}

		var endsUnix int64
		var hostID string
		hostIDList := regexp.MustCompile(`Hosted by: <@(.*?)>`).FindStringSubmatch(message.Embeds[0].Description)
		if len(hostIDList) > 1 {
			hostID = hostIDList[1]
		}
		endsUnixList := regexp.MustCompile(`Ends: <t:(.*?):R>`).FindStringSubmatch(message.Embeds[0].Description)
		if len(endsUnixList) > 1 {
			endsUnix, _ = strconv.ParseInt(endsUnixList[1], 10, 64)
		}

		if endsUnix == 0 {
			break
		}

		endsIn := utils.Misc.FormatDuration(time.Until(time.Unix(endsUnix, 0)))
		prize := message.Embeds[0].Title
		if isPrizeBlacklisted(prize, cfg) {
			return
		}
		logging.Logging.GiveawayFound(session, message, prize, hostID, endsIn, account)
		go func() {
			time.Sleep(time.Second * time.Duration(cfg.Snipers.GiveawaySniper.Sleep))
			_ = session.InteractionClick(button, message)
		}()
	case "396464677032427530": // Giveaway#7832
		// Win check

		prizeList := regexp.MustCompile(fmt.Sprintf(`:tada: %s won \*\*(.*?)\*\*`, session.State.User.Mention())).FindStringSubmatch(message.Content)
		if len(prizeList) > 1 {
			logging.Logging.GiveawayWon(session, message, prizeList[1], account)
		}

		// Entering

		if len(message.Embeds) == 0 {
			break
		}

		if !strings.Contains(message.Embeds[0].Title, "Giveaway") {
			break
		}

		var host string
		hostList := regexp.MustCompile(`Hosted by (.*)`).FindStringSubmatch(message.Embeds[0].Author.Name)
		if len(hostList) > 1 {
			host = hostList[1]
		}
		endsInTime, _ := time.Parse(time.RFC3339, message.Embeds[0].Timestamp)
		prize := message.Embeds[0].Fields[0].Value
		if isPrizeBlacklisted(prize, cfg) {
			return
		}
		logging.Logging.GiveawayFound(session, message, prize, host, utils.Misc.FormatDuration(time.Until(endsInTime)), account)
		go func() {
			time.Sleep(time.Second * time.Duration(cfg.Snipers.GiveawaySniper.Sleep))
			_ = session.MessageReactionAdd(message.ChannelID, message.ID, "🎉")
		}()
	case "700743797977514004": // Friskytool#6219
		// Win check

		prizeList := regexp.MustCompile(fmt.Sprintf(`%s has won the giveaway for (.*)\n`, session.State.User.Mention())).FindStringSubmatch(message.Content)
		if len(prizeList) > 1 {
			logging.Logging.GiveawayWon(session, message, prizeList[1], account)
		}

		// Entering

		if len(message.Embeds) == 0 {
			break
		}

		if !strings.Contains(message.Embeds[0].Description, "React with") {
			break
		}

		host := ""

		hostList := regexp.MustCompile(`Hosted by <@(.*?)>`).FindStringSubmatch(message.Embeds[0].Description)
		if len(hostList) > 1 {
			host = hostList[1]
		}

		endsInTime, _ := time.Parse(time.RFC3339, message.Embeds[0].Timestamp)

		prize := message.Embeds[0].Title
		if isPrizeBlacklisted(prize, cfg) {
			return
		}

		logging.Logging.GiveawayFound(session, message, prize, host, utils.Misc.FormatDuration(time.Until(endsInTime)), account)

		go func() {
			time.Sleep(time.Second * time.Duration(cfg.Snipers.GiveawaySniper.Sleep))
			_ = session.MessageReactionAdd(message.ChannelID, message.ID, "🎉")
		}()
	case "530082442967646230": // Giveaway Boat#2911
		if len(message.Embeds) == 0 {
			break
		}

		// Win check

		if message.Content == fmt.Sprintf(`Congratulations to %s 🎉`, session.State.User.Mention()) {
			prizeList := regexp.MustCompile(fmt.Sprintf(`%s won the giveaway of \[(.*?)\]`, session.State.User.Mention())).FindStringSubmatch(message.Embeds[0].Description)

			if len(prizeList) > 1 {
				logging.Logging.GiveawayWon(session, message, prizeList[1], account)
			}
		}

		// Entering

		if !strings.Contains(message.Embeds[0].Description, "button to enter!") {
			break
		}

		if len(message.Components) == 0 {
			return
		}

		component, ok := message.Components[0].(*discordgo.ActionsRow)
		if !ok {
			break
		}

		if len(component.Components) == 0 {
			return
		}

		button, ok := component.Components[0].(*discordgo.Button)
		if !ok {
			break
		}

		endsInTime, _ := time.Parse(time.RFC3339, message.Embeds[0].Timestamp)

		prize := message.Embeds[0].Title
		if isPrizeBlacklisted(prize, cfg) {
			return
		}

		logging.Logging.GiveawayFound(session, message, prize, "", utils.Misc.FormatDuration(time.Until(endsInTime)), account)

		go func() {
			time.Sleep(time.Second * time.Duration(cfg.Snipers.GiveawaySniper.Sleep))
			_ = session.InteractionClick(button, message)
		}()
	case "239631525350604801": // Pancake#3691
		// Win check

		pattern := fmt.Sprintf(
			`__Congratulations:__\n\*\*(.*?)\*\* \(%s\)\nYou won \*\*(.*?)\*\*!`,
			session.State.User.Mention(),
		)

		prizeList := regexp.MustCompile(pattern).FindStringSubmatch(message.Content)
		if len(prizeList) > 2 {
			logging.Logging.GiveawayWon(session, message, prizeList[2], account)
		}

		// Entering

		if len(message.Embeds) == 0 {
			break
		}

		if !strings.Contains(message.Embeds[0].Description, "React with") {
			break
		}

		prize := ""

		prizeList = regexp.MustCompile(`🎉 (.*?) 🎉`).FindStringSubmatch(message.Embeds[0].Title)
		if len(prizeList) > 1 {
			prize = prizeList[1]
		}

		if isPrizeBlacklisted(prize, cfg) {
			return
		}

		var endsUnix int64

		endsUnixList := regexp.MustCompile(`Event ends <t:(.*?):`).FindStringSubmatch(message.Embeds[0].Description)
		if len(endsUnixList) > 1 {
			endsUnix, _ = strconv.ParseInt(endsUnixList[1], 10, 64)
		}

		endsIn := utils.Misc.FormatDuration(time.Until(time.Unix(endsUnix, 0)))

		logging.Logging.GiveawayFound(session, message, prize, "", endsIn, account)

		go func() {
			time.Sleep(time.Second * time.Duration(cfg.Snipers.GiveawaySniper.Sleep))
			_ = session.MessageReactionAdd(message.ChannelID, message.ID, "🎁")
		}()
	case "570338970261782559": // GiftCord#8109
		// Win check

		pattern := fmt.Sprintf(`Congratulations, %s! You won \*\*(.*?)\*\*!`, session.State.User.Mention())

		prizeList := regexp.MustCompile(pattern).FindStringSubmatch(message.Content)
		if len(prizeList) > 1 {
			logging.Logging.GiveawayWon(session, message, prizeList[1], account)
		}

		// Entering

		if len(message.Embeds) == 0 {
			break
		}

		if !strings.Contains(message.Embeds[0].Description, "Click the **Enter** button below to enter!") {
			break
		}

		if len(message.Components) == 0 {
			return
		}

		component, ok := message.Components[0].(*discordgo.ActionsRow)
		if !ok {
			break
		}

		if len(component.Components) == 0 {
			return
		}

		button, ok := component.Components[0].(*discordgo.Button)
		if !ok {
			break
		}

		prize := message.Embeds[0].Title
		if isPrizeBlacklisted(prize, cfg) {
			return
		}

		hostID := ""
		hostIDList := regexp.MustCompile(`Hosted by: <@(.*?)>`).FindStringSubmatch(message.Embeds[0].Description)
		if len(hostIDList) > 1 {
			hostID = hostIDList[1]
		}

		endsInTime, _ := time.Parse(time.RFC3339, message.Embeds[0].Timestamp)
		endsIn := utils.Misc.FormatDuration(time.Until(endsInTime))

		logging.Logging.GiveawayFound(session, message, prize, hostID, endsIn, account)

		go func() {
			time.Sleep(time.Second * time.Duration(cfg.Snipers.GiveawaySniper.Sleep))
			_ = session.InteractionClick(button, message)
		}()
	case "651627261782523926": // Giveaways#8717
		// Win check

		prizeList := regexp.MustCompile(fmt.Sprintf(`Congratulations <@!%s>, you won \*\*(.*?)\*\*`, session.State.User.ID)).FindStringSubmatch(message.Content)
		if len(prizeList) > 1 {
			logging.Logging.GiveawayWon(session, message, prizeList[1], account)
		}

		// Entering

		if len(message.Embeds) == 0 {
			break
		}

		emb := message.Embeds[0]

		if !strings.Contains(emb.Description, "Click the reaction below to enter!") {
			return
		}

		var endsUnix int64
		var hostID string

		hostIDList := regexp.MustCompile(`Host: <@!(.*?)>`).FindStringSubmatch(message.Embeds[0].Description)
		if len(hostIDList) > 1 {
			hostID = hostIDList[1]
		}

		endsUnixList := regexp.MustCompile(`Ends <t:(.*?):R>`).FindStringSubmatch(message.Embeds[0].Description)
		if len(endsUnixList) > 1 {
			endsUnix, _ = strconv.ParseInt(endsUnixList[1], 10, 64)
		}
		endsIn := utils.Misc.FormatDuration(time.Until(time.Unix(endsUnix, 0)))

		prize := message.Embeds[0].Title

		if isPrizeBlacklisted(prize, cfg) {
			return
		}

		logging.Logging.GiveawayFound(session, message, prize, hostID, endsIn, account)

		go func() {
			time.Sleep(time.Second * time.Duration(cfg.Snipers.GiveawaySniper.Sleep))
			_ = session.MessageReactionAdd(message.ChannelID, message.ID, "🎉")
		}()
	case "693167035068317736": // Rumble Royale#6012
		messageJsonBytes, err := json.Marshal(message)
		if err != nil {
			return
		}

		messageJson := string(messageJsonBytes)

		// Win check

		if strings.Contains(messageJson, "**__WINNER!__**") &&
			strings.Contains(message.Content, session.State.User.Mention()) {
			prize := ""
			prizeList := regexp.MustCompile(`\*\*Reward:\*\* (.*?) `).FindStringSubmatch(messageJson)
			if len(prizeList) > 1 {
				prize = prizeList[1]
			}

			logging.Logging.GiveawayWon(session, message, prize, account)

			return
		}
	case "368521195940741122": // Lawliet#5480
		if len(message.Embeds) == 0 {
			break
		}

		emb := message.Embeds[0]

		// Win check

		if strings.Contains(message.Content, session.State.User.Mention()) && // Mention in content
			strings.Contains(emb.Description, "You just won in the giveaway!") { // Winning text

			prize := strings.ReplaceAll(emb.Title, "🎆 ", "")

			logging.Logging.GiveawayWon(session, message, prize, account)
		}

		// Entering

		if !strings.Contains(emb.Description, "React with 🎉 to enter") {
			return
		}

		var endsUnix int64

		endsUnixList := regexp.MustCompile(`<t:(.*?):R>`).FindStringSubmatch(emb.Description)
		if len(endsUnixList) > 1 {
			endsUnix, _ = strconv.ParseInt(endsUnixList[1], 10, 64)
		}

		endsIn := utils.Misc.FormatDuration(time.Until(time.Unix(endsUnix, 0)))

		prize := strings.ReplaceAll(emb.Title, "🎆 ", "")

		if isPrizeBlacklisted(prize, cfg) {
			return
		}

		logging.Logging.GiveawayFound(session, message, prize, "", endsIn, account)

		go func() {
			time.Sleep(time.Second * time.Duration(cfg.Snipers.GiveawaySniper.Sleep))
			_ = session.MessageReactionAdd(message.ChannelID, message.ID, "🎉")
		}()
	case "159985870458322944": // MEE6#4876
		if len(message.Embeds) == 0 {
			break
		}

		emb := message.Embeds[0]

		// Win check

		if strings.Contains(message.Content, session.State.User.Mention()) && // Mention in content
			strings.Contains(emb.Description, "Congratulations") { // Winning text
			prize := ""

			if len(emb.Fields) > 0 {
				prize = emb.Fields[0].Value
			}

			logging.Logging.GiveawayWon(session, message, prize, account)
		}

		// Entering

		if len(message.Components) == 0 {
			break
		}

		component, ok := message.Components[0].(*discordgo.ActionsRow)
		if !ok {
			break
		}

		if len(component.Components) != 2 {
			break
		}

		enterButton, ok := component.Components[0].(*discordgo.Button)
		if !ok {
			break
		}

		timeButton, ok := component.Components[1].(*discordgo.Button)
		if !ok {
			break
		}

		if enterButton.Emoji.Name != "🎉" || !strings.Contains(enterButton.Label, "Participate") {
			break
		}

		// End time
		label := regexp.MustCompile(`\((.*?)\)`).ReplaceAllLiteralString(timeButton.Label, "")

		endsIn := ""

		endTime, err := time.Parse("Ends on Mon _2 Jan 2006 at 15:04 ", label)
		if err == nil {
			endsIn = utils.Misc.FormatDuration(time.Until(endTime))
		}

		logging.Logging.GiveawayFound(session, message, "", "", endsIn, account)

		go func() {
			time.Sleep(time.Second * time.Duration(cfg.Snipers.GiveawaySniper.Sleep))
			_ = session.InteractionClick(enterButton, message)
		}()
	default:
		isCustomBot := utils.Discord.IsUserInList(session, message.Author.ID, cfg.Snipers.GiveawaySniper.CustomBots)

		if !isCustomBot {
			break
		}

		messageJsonBytes, err := json.Marshal(message)
		if err != nil {
			return
		}

		messageJson := strings.ToLower(string(messageJsonBytes))

		isNewGiveaway := strings.Contains(messageJson, "react with") ||
			strings.Contains(messageJson, "to enter") ||
			strings.Contains(messageJson, "to join") ||
			strings.Contains(messageJson, "giveaway")

		if !isNewGiveaway {
			return
		}

		logging.Logging.GiveawayFound(session, message, "", "", "", account)

		// Check if there is a button to enter the giveaway

		hasComponent := func() bool {
			if len(message.Components) == 0 {
				return false
			}

			component, ok := message.Components[0].(*discordgo.ActionsRow)
			if !ok {
				return false
			}

			if len(component.Components) == 0 {
				return false
			}

			button, ok := component.Components[0].(*discordgo.Button)
			if !ok {
				return false
			}

			go func() {
				time.Sleep(time.Second * time.Duration(cfg.Snipers.GiveawaySniper.Sleep))
				_ = session.InteractionClick(button, message)
			}()

			return true
		}()

		if hasComponent {
			return
		}

		go func() {
			// Wait until the bot has probably reacted with an emoji
			time.Sleep(10 * time.Second)

			// Reload the message to get the updated reactions
			newMessage, err := utils.Discord.GetChannelMessage(session, message.ChannelID, message.ID)
			if err != nil {
				return
			}

			// if there are reactions, react with the first one
			if len(newMessage.Reactions) == 0 {
				return
			}

			emoji := newMessage.Reactions[0].Emoji.APIName()

			time.Sleep(time.Second * time.Duration(cfg.Snipers.GiveawaySniper.Sleep))

			_ = session.MessageReactionAdd(message.ChannelID, message.ID, emoji)
		}()
	}

	// Dyno#3861
	isDynoWebhook := message.Author.Username == "Dyno" && message.Author.Discriminator == "0000" && message.WebhookID != ""
	isDynoBot := message.Author.ID == "155149108183695360"
	if isDynoWebhook || isDynoBot {
		// Win check

		pattern := fmt.Sprintf(`Congratulations %s, you won (.*?)!`, session.State.User.Mention())

		prizeList := regexp.MustCompile(pattern).FindStringSubmatch(message.Content)
		if len(prizeList) > 1 {
			logging.Logging.GiveawayWon(session, message, prizeList[1], account)
		}

		// Entering

		if len(message.Embeds) == 0 {
			return
		}

		emb := message.Embeds[0]

		if len(emb.Fields) != 2 {
			return
		}

		if emb.Image == nil {
			return
		}

		if emb.Footer == nil {
			return
		}

		if !strings.HasPrefix(emb.URL, "https://dyno.gg/giveaway/") {
			return
		}

		if emb.Image.URL != "https://cdn.dyno.gg/giveaway/banner.jpeg" {
			return
		}

		if len(message.Components) == 0 {
			return
		}

		component, ok := message.Components[0].(*discordgo.ActionsRow)
		if !ok {
			return
		}

		if len(component.Components) == 0 {
			return
		}

		button, ok := component.Components[0].(*discordgo.Button)
		if !ok {
			return
		}

		prize := emb.Title

		if isPrizeBlacklisted(prize, cfg) {
			return
		}

		endsInList := regexp.MustCompile(`<t:(.*?):R>`).FindStringSubmatch(emb.Fields[0].Value)
		hostIDList := regexp.MustCompile("<@(.*?)>").FindStringSubmatch(emb.Fields[1].Value)

		if len(endsInList) < 2 || len(hostIDList) < 2 {
			return
		}

		endsUnix, _ := strconv.ParseInt(endsInList[1], 10, 64)

		endsIn := utils.Misc.FormatDuration(time.Until(time.Unix(endsUnix, 0)))

		logging.Logging.GiveawayFound(session, message, prize, hostIDList[1], endsIn, account)

		go func() {
			time.Sleep(time.Second * time.Duration(cfg.Snipers.GiveawaySniper.Sleep))
			_ = session.InteractionClick(button, message)
		}()
	}
}

func CheckGiveawaySniperEdit(session *discordgo.Session, message *discordgo.Message, cfg *config.Config, account string) {
	if !cfg.Snipers.GiveawaySniper.Enabled {
		return
	}

	switch message.Author.ID {
	case "676039466175168543": // ZAP-Hosting Bot#7133
		if !strings.Contains(message.Content, "This Giveaway has ended!") {
			break
		}

		if !strings.Contains(message.Content, session.State.User.Mention()) {
			break
		}

		logging.Logging.GiveawayWon(session, message, "€10 worth of ZAP Coins", account)
	case "693167035068317736": // Rumble Royale#6012
		messageJsonBytes, err := json.Marshal(message)
		if err != nil {
			break
		}

		messageJson := strings.ToLower(string(messageJsonBytes))

		// Entering

		if !strings.Contains(messageJson, "click the emoji below to join") {
			break
		}

		host := ""

		if message.Interaction != nil {
			if message.Interaction.User != nil {
				host = message.Interaction.User.Username
			}
		}

		battleEra := ""
		battleEraList := regexp.MustCompile(`era: \\u003c(.*?)\\u003e(.*?) `).FindStringSubmatch(messageJson)
		if len(battleEraList) > 2 {
			battleEra = battleEraList[2]
		}

		startingIn := ""
		startingInList := regexp.MustCompile(`starting in (.*?)!`).FindStringSubmatch(messageJson)
		if len(startingInList) > 1 {
			startingIn = startingInList[1]
		}

		logging.Logging.GiveawayFound(session, message, battleEra, host, startingIn, account)

		go func() {
			time.Sleep(time.Second * 30)

			_ = session.MessageReactionAdd(message.ChannelID, message.ID, "Swrds:872886436012126279")
		}()
	case "368521195940741122": // Lawliet#5480
		// Win check
		if len(message.Embeds) == 0 {
			break
		}

		emb := message.Embeds[0]

		if !strings.Contains(message.Content, session.State.User.Mention()) || // Mention in content
			!strings.Contains(emb.Description, "You just won in the giveaway!") { // Winning text
			break
		}

		prize := strings.ReplaceAll(emb.Title, "🎆 ", "")

		logging.Logging.GiveawayWon(session, message, prize, account)
	}
}
