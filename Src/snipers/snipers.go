package snipers

import (
	"encoding/json"
	"fmt"
	"go/token"
	"go/types"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/LightningDev1/LB-Selfbot/config"
	"github.com/LightningDev1/LB-Selfbot/logging"
	"github.com/LightningDev1/LB-Selfbot/utils"
	"github.com/LightningDev1/discordgo"
	"github.com/bitly/go-simplejson"
)

type NitroCode struct {
	Code string
	Time time.Time
}

var nitroRegex = regexp.MustCompile(`(discord\.gift|discord\.com/gifts|discordapp\.com/gifts|promos\.discord\.gg/)/(\w{0,25})`)

func getCodesInLast10Seconds(codes []NitroCode) []NitroCode {
	return utils.Filter(codes, func(code NitroCode) bool {
		return time.Since(code.Time) < 10*time.Second
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

	nitroLinks := nitroRegex.FindAllString(message.Content, -1)

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
			if len(getCodesInLast10Seconds(codes)) >= 5 {
				return alreadyCheckedNitroCodes, sentNitroCodes
			}
		}

		url := fmt.Sprintf("https://discord.com/api/v10/entitlements/gift-codes/%s/redeem", nitroCode)

		start := time.Now()

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
				sentNitroCodes[message.Author.ID] = append(sentNitroCodes[message.Author.ID], NitroCode{})
			} else {
				sentNitroCodes[message.Author.ID] = []NitroCode{code}
			}
		}
	}
	return alreadyCheckedNitroCodes, sentNitroCodes
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

	switch message.Author.ID {
	case "676039466175168543": // ZAP-Hosting Bot#7133
		if !strings.Contains(message.Content, "It's time for a Giveaway!") {
			break
		}

		actionsRow, ok := message.Components[0].(*discordgo.ActionsRow)
		if !ok {
			break
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
	case "294882584201003009": // GiveawayBot#2381
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

		component, ok := message.Components[0].(*discordgo.ActionsRow)
		if !ok {
			break
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

		if strings.Contains(message.Embeds[0].Description, "React with") {
			endsInTime, _ := time.Parse(time.RFC3339, message.Embeds[0].Timestamp)

			prize := message.Embeds[0].Title
			if isPrizeBlacklisted(prize, cfg) {
				return
			}

			logging.Logging.GiveawayFound(session, message, prize, "", utils.Misc.FormatDuration(time.Until(endsInTime)), account)

			go func() {
				time.Sleep(time.Second * time.Duration(cfg.Snipers.GiveawaySniper.Sleep))
				_ = session.MessageReactionAdd(message.ChannelID, message.ID, "🎉")
			}()
		}
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
	case "570338970261782559": // GiveawayBot#8109
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

		component, ok := message.Components[0].(*discordgo.ActionsRow)
		if !ok {
			break
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
	default:
		isCustomBot := utils.Discord.IsUserInList(session, message.Author.ID, cfg.Snipers.GiveawaySniper.CustomBots)

		if !isCustomBot {
			return
		}

		messageJsonBytes, err := json.Marshal(message)
		if err != nil {
			return
		}

		messageJson := strings.ToLower(string(messageJsonBytes))

		isNewGiveaway := strings.Contains(messageJson, "react with") ||
			strings.Contains(messageJson, "to enter") ||
			strings.Contains(messageJson, "giveaway")

		if !isNewGiveaway {
			return
		}

		// Check if there is a button to enter the giveaway
		if component, ok := message.Components[0].(*discordgo.ActionsRow); ok {
			if button, ok := component.Components[0].(*discordgo.Button); ok {
				logging.Logging.GiveawayFound(session, message, "", "", "", account)

				go func() {
					time.Sleep(time.Second * time.Duration(cfg.Snipers.GiveawaySniper.Sleep))
					_ = session.InteractionClick(button, message)
				}()

				return
			}
		}

		go func() {
			// Wait until the bot has probably reacted with an emoji
			time.Sleep(10 * time.Second)

			// Reload the message to get the updated reactions
			newMessage, err := session.ChannelMessage(message.ChannelID, message.ID)
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
	}
}
