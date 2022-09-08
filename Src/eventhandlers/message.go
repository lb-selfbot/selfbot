package eventhandlers

import (
	"net/url"
	"regexp"
	"strings"
	"time"
	"unicode"

	"github.com/LightningDev1/LB-Selfbot/api"
	"github.com/LightningDev1/LB-Selfbot/config"
	"github.com/LightningDev1/LB-Selfbot/embed"
	"github.com/LightningDev1/LB-Selfbot/http"
	"github.com/LightningDev1/LB-Selfbot/logging"
	"github.com/LightningDev1/LB-Selfbot/snipers"
	"github.com/LightningDev1/LB-Selfbot/utils"
	"github.com/LightningDev1/dgc"
	"github.com/LightningDev1/discordgo"
)

var inviteRegex = regexp.MustCompile(`(https://)?(www\.)?(((discord(app)?)?\.com/invite)|((discord(app)?)?\.gg))/(.+)`)

func (e *Events) OnMessage(_ *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author == nil || m.Message == nil {
		logging.Warn("Message: Author or Message is nil: ", m.Author, m.Message)
		return
	}
	if m.Type != discordgo.MessageTypeDefault && m.Type != discordgo.MessageTypeReply {
		return
	}
	e.Settings.discordTimeDiff = time.Since(utils.Discord.SnowflakeToTime(m.ID))
	e.Session.State.MessageAdd(m.Message)
	cfg, err := config.GetConfig()
	if err != nil {
		return
	}
	e.Settings.AlreadyCheckedNitroCodes, e.Settings.SentNitroCodes = snipers.CheckNitroSniper(
		e.Session,
		m.Message,
		cfg,
		e.Settings.AlreadyCheckedNitroCodes,
		e.Settings.SentNitroCodes,
		"",
	)
	snipers.CheckGiveawaySniper(e.Session, m.Message, &cfg, "")
	e.checkIgnore(m.Message)
	e.checkPhising(m.Message, &cfg)
	e.checkEventReactions(m.Message)
	e.checkAFK(m.Message, &cfg)
	e.checkProtections(m.Message, &cfg)
	e.checkNotifiers(m.Message, &cfg)
	e.checkSelfbotDetection(m.Message, &cfg)
}

func (e *Events) OnMessageEdit(_ *discordgo.Session, m *discordgo.MessageUpdate) {
	message, _ := e.Session.State.Message(m.Message.ChannelID, m.Message.ID)
	if message != nil {
		m.Message = message
	}
	if m.Author == nil || m.Message == nil {
		logging.Warn("MessageEdit: Author or Message is nil: ", m.Author, m.Message)
		return
	}
	if m.Type != discordgo.MessageTypeDefault {
		return
	}
	cfg, err := config.GetConfig()
	if err != nil {
		return
	}
	e.Settings.AlreadyCheckedNitroCodes, e.Settings.SentNitroCodes = snipers.CheckNitroSniper(
		e.Session,
		m.Message,
		cfg,
		e.Settings.AlreadyCheckedNitroCodes,
		e.Settings.SentNitroCodes,
		"",
	)
	e.checkIgnore(m.Message)
	e.checkPhising(m.Message, &cfg)
	e.checkNotifiersEdit(m.Message, &cfg)
	e.checkSelfbotDetectionEdit(m.Message, &cfg)

	messages, ok := e.Settings.EditedMessages[m.ChannelID]
	if ok {
		e.Settings.EditedMessages[m.ChannelID] = append(messages, m.Message)
	} else {
		e.Settings.EditedMessages[m.ChannelID] = []*discordgo.Message{m.Message}
	}
}

func (e *Events) OnMessageDelete(_ *discordgo.Session, m *discordgo.MessageDelete) {
	if m.BeforeDelete != nil {
		m.Message = m.BeforeDelete
	}
	if m.Author == nil || m.Message == nil {
		//logging.Warn("MessageDelete: Author or Message is nil: ", m.Author, m.Message)
		return
	}
	if m.Type != discordgo.MessageTypeDefault {
		return
	}
	cfg, err := config.GetConfig()
	if err != nil {
		return
	}
	e.checkProtectionsDelete(m.Message, &cfg)
	e.checkNotifiersDelete(m.Message, &cfg)
	e.checkSelfbotDetectionDelete(m.Message, &cfg)

	messages, ok := e.Settings.DeletedMessages[m.ChannelID]
	if ok {
		e.Settings.DeletedMessages[m.ChannelID] = append(messages, m.Message)
	} else {
		e.Settings.DeletedMessages[m.ChannelID] = []*discordgo.Message{m.Message}
	}
}

// func getCodesInLast10Seconds(codes []NitroCode) []NitroCode {
// 	codesInLast10Seconds := []NitroCode{}
// 	for _, code := range codes {
// 		if time.Since(code.Time) < 10*time.Second {
// 			codesInLast10Seconds = append(codesInLast10Seconds, code)
// 		}
// 	}
// 	return codesInLast10Seconds
// }

// func (e *Events) checkNitroSniper(message *discordgo.Message, cfg *config.Config) {
// 	if !cfg.Snipers.NitroSniper {
// 		return
// 	}
// 	nitroLinks := nitroRegex.FindAllString(message.Content, -1)
// 	for _, nitroLink := range nitroLinks {
// 		nitroLinkSplit := strings.Split(nitroLink, "/")
// 		nitroCode := nitroLinkSplit[len(nitroLinkSplit)-1]
// 		if utils.SliceContains(e.Settings.alreadyCheckedNitroCodes, nitroCode) {
// 			return
// 		}
// 		e.Settings.alreadyCheckedNitroCodes = append(e.Settings.alreadyCheckedNitroCodes, nitroCode)
// 		if len(nitroCode) < 15 || len(nitroCode) > 26 {
// 			logging.UI.NitroSniperFake(e.Session, message, nitroCode, "")
// 			return
// 		}
// 		if codes, ok := e.Settings.sentNitroCodes[message.Author.ID]; ok {
// 			if len(getCodesInLast10Seconds(codes)) >= 5 {
// 				return
// 			}
// 		}
// 		start := time.Now()
// 		result, err := http.Post(fmt.Sprintf("https://discord.com/api/v10/entitlements/gift-codes/%s/redeem", nitroCode), "", http.GetDiscordHeaders(cfg.Account.Discord.Token))
// 		if err != nil {
// 			return
// 		}
// 		elapsed := time.Since(start)
// 		redeemed := false
// 		if strings.Contains(result, "This gift has been redeemed already.") {
// 			logging.UI.NitroSniper(e.Session, message, elapsed, "Already Redeemed", nitroCode, "")
// 		} else if strings.Contains(result, "Already purchased") {
// 			logging.UI.NitroSniper(e.Session, message, elapsed, "Already Purchased", nitroCode, "")
// 		} else if strings.Contains(result, "subscription_plan") {
// 			logging.UI.NitroSniper(e.Session, message, elapsed, "Redeemed Nitro", nitroCode, "")
// 			redeemed = true
// 		} else if strings.Contains(result, "Unknown Gift Code") {
// 			logging.UI.NitroSniper(e.Session, message, elapsed, "Unknown Code", nitroCode, "")
// 		} else {
// 			json, err := simplejson.NewJson([]byte(result))
// 			if err == nil {
// 				logging.UI.NitroSniper(e.Session, message, elapsed, utils.Misc.Or(json.Get("message").MustString(), "Unknown Status"), nitroCode, "")
// 			}
// 		}
// 		if !redeemed {
// 			code := NitroCode{
// 				Code: nitroCode,
// 				Time: time.Now(),
// 			}
// 			_, ok := e.Settings.sentNitroCodes[message.Author.ID]
// 			if ok {
// 				e.Settings.sentNitroCodes[message.Author.ID] = append(e.Settings.sentNitroCodes[message.Author.ID], NitroCode{})
// 			} else {
// 				e.Settings.sentNitroCodes[message.Author.ID] = []NitroCode{code}
// 			}
// 		}
// 	}
// }

// var privnoteRegex = regexp.MustCompile(`privnote.com/([a-zA-Z\d_#]{18})`)

// func (e *Events) checkPrivnoteSniper(message *discordgo.Message, cfg *config.Config) {
// 	if !cfg.Snipers.PrivnoteSniper {
// 		return
// 	}
// 	privnoteLinks := privnoteRegex.FindAllString(message.Content, -1)
// 	for _, privnoteLink := range privnoteLinks {
// 		privnoteLink = "https://" + privnoteLink
// 		// try to open
// 		// impossible because of cloudflare?
// 	}
// }

func (e *Events) checkGiveawaySniper(m *discordgo.Message, cfg *config.Config) {
	if !cfg.Snipers.GiveawaySniper.Sniper || !m.Author.Bot {
		return
	}
	lowerContent := strings.ToLower(m.Content)

	winningWords := []string{"congratulations", "new", "**congratulations**", "congratulations,", "__congratulations:__", "won"}
	for _, winningWord := range winningWords {
		if !strings.Contains(lowerContent, winningWord) {
			continue
		}
		for _, mention := range m.Mentions {
			if mention.ID != e.Session.State.User.ID {
				continue
			}
			switch m.Author.ID {
			case "294882584201003009": // GiveawayBot#2381
			case "570338970261782559": // GiveawayBot#8109
			case "606026008109514762": // Giveawaysbot#3411
			case "396464677032427530": // Giveaway#7832
				if strings.Contains(lowerContent, "the new winner") {
					// TODO: Won Unknown Giveaway
				} else {
					prizeSlice := strings.Split(m.Content, "**")
					if len(prizeSlice) < 1 {
						// TODO: Won Unknown Giveaway
					}
					prize := prizeSlice[1]
					_ = prize
					// TODO: Won Giveaway
				}
			case "673918978178940951": // Prism#5621
				if len(m.Embeds) < 1 {
					break
				}
				splitDescription := strings.Split(m.Embeds[0].Description, "[")
				if len(splitDescription) < 2 {
					break
				}
				splitDescription = strings.Split(splitDescription[1], "]")
				if len(splitDescription) < 1 {
					break
				}
				prize := splitDescription[0]
				_ = prize
				// TODO: Won Giveaway
			case "239631525350604801": // Pancake#3691
				if strings.Contains(lowerContent, "you are the new winner!") {
					// TODO: Won Unknown Giveaway
				} else {
					prizeSlice := strings.Split(m.Content, " **")
					if len(prizeSlice) < 2 {
						break
					}
					prize := strings.ReplaceAll(prizeSlice[1], "**!", "")
					_ = prize
					// TODO: Won Giveaway
				}
			case "720351927581278219": // Invite Tracker#0478
				prizeSlice := strings.Split(m.Content, "the **")
				if len(prizeSlice) < 2 {
					break
				}
				prize := strings.ReplaceAll(prizeSlice[1], "**!", "")
				_ = prize
				// TODO: Won Giveaway
			case "530082442967646230": // Giveaway Boat#2911
				if strings.Contains(lowerContent, "new winner") {
					prizeSlice := strings.Split(m.Content, "**")
					if len(prizeSlice) < 2 {
						break
					}
					prize := prizeSlice[1]
					_ = prize
					// TODO: Won Giveaway
				} else if len(m.Embeds) > 0 {
					descriptionSlice := strings.Split(m.Embeds[0].Description, "[")
					if len(descriptionSlice) < 2 {
						break
					}
					descriptionSlice = strings.Split(descriptionSlice[1], "]")
					if len(descriptionSlice) < 1 {
						break
					}
					prize := descriptionSlice[0]
					_ = prize
					// TODO: Won Giveaway
				}
			default:
				if utils.SliceContains(cfg.Snipers.GiveawaySniper.CustomBots, m.Author.ID) {
					// TODO: Won Custom Giveaway
				}
			}
		}
	}

	if len(m.Embeds) < 1 {
		return
	}
	var (
		prize    string
		emoji    string
		isCustom bool
	)

	giveawayWords := []string{"**giveaway**", "**new", "giveaway**", "**giveaway", "giveaway"}
	for _, giveawayWord := range giveawayWords {
		if strings.Contains(lowerContent, giveawayWord) {
			switch m.Author.ID {
			case "294882584201003009": // GiveawayBot#2381
			case "570338970261782559": // GiveawayBot#8109
				prize = m.Embeds[0].Author.Name
				emoji = "🎉"
			case "673918978178940951": // Prism#5621
				prizeSlice := strings.Split(m.Embeds[0].Description, "[")
				if len(prizeSlice) < 2 {
					break
				}
				prize = prizeSlice[1]
				emoji = "🎉"
			case "606026008109514762": // Giveawaysbot#3411
				prize = m.Embeds[0].Title
				emoji = "🎉"
			case "239631525350604801": // Pancake#3691
				prize = m.Embeds[0].Title
				prize = strings.ReplaceAll(prize, "🎉 ", "")
				prize = strings.ReplaceAll(prize, " 🎉", "")
				emoji = "🎁"
			case "720351927581278219": // Invite Tracker#0478
				prize = m.Embeds[0].Title
				prize = strings.ReplaceAll(prize, "<a:blobgift:834855281358930020>  ", "")
				prize = strings.ReplaceAll(prize, "  <a:blobgift:834855281358930020>", "")
				emoji = "🎉"
			default:
				if utils.SliceContains(cfg.Snipers.GiveawaySniper.CustomBots, m.Author.ID) {
					prize = "Unknown"
					emoji = "🎉"
					isCustom = true
				}
			}
		}
	}

	if prize == "" && strings.Contains(m.Embeds[0].Description, "React with 🎉 to enter") {
		switch m.Author.ID {
		case "396464677032427530": // Giveaway#7832
			prize = m.Embeds[0].Fields[0].Value
			emoji = "🎉"
		case "530082442967646230": // Giveaway Boat#2911
			prize = m.Embeds[0].Title
			emoji = "🎉"
		}
	}

	if prize != "" && emoji != "" {
		for _, blacklistedWord := range cfg.Snipers.Blacklist {
			if strings.Contains(prize, blacklistedWord) {
				return
			}
		}
		_ = isCustom
		// TODO: Giveaway Found
		go func() {
			time.Sleep(time.Second * time.Duration(cfg.Snipers.GiveawaySniper.Sleep))
			_ = e.Session.MessageReactionAdd(m.ChannelID, m.ID, emoji)
		}()
	}
}

func (e *Events) checkIgnore(message *discordgo.Message) {
	for _, user := range e.Settings.IgnoredUsers {
		if message.Author.ID == user.ID {
			_ = e.Session.MessageAck(message.ChannelID, message.ID)
		}
	}
	for _, guild := range e.Settings.IgnoredGuilds {
		if message.GuildID == guild.ID {
			_ = e.Session.GuildAck(message.GuildID)
		}
	}
}

func (e *Events) checkPhising(message *discordgo.Message, cfg *config.Config) {
	if !cfg.Protections.AntiPhising {
		return
	}

	for _, part := range strings.Split(message.Content, " ") {
		link, err := url.Parse(part)
		if err == nil {
			hostname := link.Hostname()
			if utils.SliceContains(e.Settings.phisingDomains, hostname) ||
				utils.SliceContains(e.Settings.suspiciousDomains, hostname) {
				_ = e.Session.MessageAck(message.ChannelID, message.ID)
				logging.UI.PhisingDetected(e.Session, message, hostname)
			}
		}
	}
}

func (e *Events) checkEventReactions(message *discordgo.Message) {
	if message.Author.ID == e.Session.State.User.ID || message.GuildID == "" {
		return
	}
	member, err := e.Session.State.Member(message.GuildID, message.Author.ID)
	if err != nil {
		return
	}
	isMentioned := false
	for _, user := range message.Mentions {
		if user.ID == e.Session.State.User.ID {
			isMentioned = true
			break
		}
	}
	if isMentioned {
		if utils.SliceContains(EventHandlers.Settings.PingMuteServers, message.GuildID) {
			_ = utils.Discord.MuteUser(e.Session, member)
		}
		if utils.SliceContains(EventHandlers.Settings.PingTimeoutServers, message.GuildID) {
			until := time.Now().Add(5 * time.Minute)
			_ = e.Session.GuildMemberTimeout(member.GuildID, member.User.ID, &until)
		}
		if utils.SliceContains(EventHandlers.Settings.PingBanServers, message.GuildID) {
			_ = e.Session.GuildBanCreateWithReason(member.GuildID, member.User.ID, "Ping Ban Event Reaction", 0)
		}
		if utils.SliceContains(EventHandlers.Settings.PingKickServers, message.GuildID) {
			_ = e.Session.GuildMemberDeleteWithReason(member.GuildID, member.User.ID, "Ping Kick Event Reaction")
		}
		if utils.SliceContains(EventHandlers.Settings.PingInsultServers, message.GuildID) {
			insult, _ := http.Get("https://evilinsult.com/generate_insult.php?lang=en")
			_, _ = e.Session.ChannelMessageSend(message.ChannelID, insult)
		}
		if utils.SliceContains(EventHandlers.Settings.PingPingServers, message.GuildID) {
			_, _ = e.Session.ChannelMessageSend(message.ChannelID, message.Author.Mention())
		}
	}
}

func (e *Events) checkAFK(message *discordgo.Message, cfg *config.Config) {
	if message.Author.ID == e.Session.State.User.ID ||
		!cfg.AFK.Enabled ||
		utils.SliceContains(cfg.AFK.IgnoredUsers, message.Author.ID) ||
		utils.SliceContains(cfg.AFK.IgnoredServers, message.GuildID) {
		return
	}
	isMentioned := false
	for _, user := range message.Mentions {
		if user.ID == e.Session.State.User.ID {
			isMentioned = true
			break
		}
	}
	if (isMentioned && cfg.AFK.OnPing) || (message.GuildID == "" && cfg.AFK.OnDM) {
		_ = e.Session.MessageAck(message.ChannelID, message.ID)
		if cfg.AFK.Delay > 0 {
			if cfg.AFK.Typing {
				utils.Discord.SleepWithTyping(e.Session, message.ChannelID, time.Second*time.Duration(cfg.AFK.Delay))
			} else {
				time.Sleep(time.Second * time.Duration(cfg.AFK.Delay))
			}
			time.Sleep(2 * time.Second)
		}
		if cfg.AFK.Reply {
			_, _ = e.Session.ChannelMessageSendReply(message.ChannelID, cfg.AFK.Message, message.Reference())
		} else {
			_, _ = e.Session.ChannelMessageSend(message.ChannelID, cfg.AFK.Message)
		}
	}
}

func (e *Events) checkProtections(message *discordgo.Message, cfg *config.Config) {
	if message.Author.ID == e.Session.State.User.ID {
		return
	}
	if cfg.Protections.AntiMassDM && message.GuildID == "" {
		logging.Warn("message.go:477: ", message)
		amountOfDMs := 0
		messages, _ := utils.Discord.GetChannelMessages(e.Session, message.ChannelID, 100)
		for _, message := range messages {
			if message.Author.ID != e.Session.State.User.ID {
				amountOfDMs++
			}
		}
		if amountOfDMs < 2 {
			e.Settings.newDMs = append(e.Settings.newDMs, message)
			dmsInLastMinute := 0
			for _, dm := range e.Settings.newDMs {
				if time.Since(dm.Timestamp) < time.Minute {
					dmsInLastMinute++
				}
			}
			if dmsInLastMinute > 10 {
				var guildIDs []string
				for _, g := range e.Session.State.Guilds {
					guildIDs = append(guildIDs, g.ID)
				}
				_ = e.Session.MessageAck(message.ChannelID, message.ID)
				_, _ = api.ChangeSettings(map[string]interface{}{"restricted_guilds": guildIDs}, true)
				logging.UI.MassDMBlocked()
				go func() {
					time.Sleep(5 * 60)
					_, _ = api.ChangeSettings(map[string]interface{}{"restricted_guilds": []any{}}, true)
				}()
			}
		}
	}
	if message.GuildID != "" && utils.SliceContains(cfg.Protections.ServerProtections.ProtectedServers, message.GuildID) {
		userPerms, _ := e.Session.UserChannelPermissions(message.Author.ID, message.ChannelID)
		if userPerms&discordgo.PermissionManageMessages != 0 || userPerms&discordgo.PermissionAdministrator != 0 {
			return
		}
		if cfg.Protections.ServerProtections.AntiInvite && inviteRegex.MatchString(message.Content) {
			err := e.Session.ChannelMessageDelete(message.ChannelID, message.ID)
			if err == nil {
				ctx := &dgc.Ctx{
					Session: e.Session,
					Event:   &discordgo.MessageCreate{Message: message},
				}
				embed.NewEmbed().
					SetDescription("Invites are not allowed in this server.").
					Send(ctx)
			}
		}
		if cfg.Protections.ServerProtections.AntiSpam {
			messages, _ := utils.Discord.GetChannelMessages(e.Session, message.ChannelID, 100)
			userMessagesInLastMinute := 0
			for _, m := range messages {
				if m.Author.ID == message.Author.ID && time.Since(m.Timestamp) < time.Second*10 {
					userMessagesInLastMinute++
				}
			}
			if userMessagesInLastMinute > 10 {
				err := e.Session.ChannelMessageDelete(message.ChannelID, message.ID)
				if err == nil {
					ctx := &dgc.Ctx{
						Session: e.Session,
						Event:   &discordgo.MessageCreate{Message: message},
					}
					embed.NewEmbed().
						SetDescription("Spamming is not allowed in this server.").
						Send(ctx)
				}
			}
		}
		if cfg.Protections.ServerProtections.AntiUpper && len(message.Content) > 10 {
			upperChars := 0
			for _, char := range message.Content {
				if unicode.IsUpper(char) {
					upperChars++
				}
			}
			upperPercentage := float64(upperChars) / float64(len(message.Content)) * 100
			if upperPercentage >= 70 {
				err := e.Session.ChannelMessageDelete(message.ChannelID, message.ID)
				if err == nil {
					ctx := &dgc.Ctx{
						Session: e.Session,
						Event:   &discordgo.MessageCreate{Message: message},
					}
					embed.NewEmbed().
						SetDescription("Uppercase spam is not allowed in this server.").
						Send(ctx)
				}
			}
		}
	}
}

func (e *Events) checkProtectionsDelete(message *discordgo.Message, cfg *config.Config) {
	if cfg.Protections.ServerProtections.AntiDelete &&
		message.GuildID != "" &&
		utils.SliceContains(cfg.Protections.ServerProtections.ProtectedServers, message.GuildID) {
		userPerms, _ := e.Session.UserChannelPermissions(message.Author.ID, message.ChannelID)
		if userPerms&discordgo.PermissionManageMessages != 0 || userPerms&discordgo.PermissionAdministrator != 0 {
			return
		}
		logging.UI.AntiDelete(e.Session, message)
	}
}

func (e *Events) checkNotifiers(message *discordgo.Message, cfg *config.Config) {
	if message.Author.ID == e.Session.State.User.ID {
		return
	}

	if message.GuildID == "" {
		if cfg.Events.DMNotifier {
			logging.UI.DMNotifier(message)
		}
	}

	if message.Author.PublicFlags&discordgo.UserFlagDiscordEmployee != 0 {
		logging.UI.DiscordStaffDetection(e.Session, message)
	}
}

func (e *Events) checkNotifiersEdit(message *discordgo.Message, cfg *config.Config) {
	if message.Author.ID == e.Session.State.User.ID {
		return
	}

	if message.GuildID == "" {
		if cfg.Events.DMEditNotifier {
			logging.UI.DMEditNotifier(message)
		}

	}
}

func (e *Events) checkNotifiersDelete(message *discordgo.Message, cfg *config.Config) {
	if message.Author.ID == e.Session.State.User.ID {
		return
	}

	if message.GuildID == "" {
		if cfg.Events.DMDeleteNotifier {
			logging.UI.DMDeleteNotifier(message)
		}
	} else {
		if cfg.Events.MessageDeleteNotifier {
			logging.UI.MessageDeleteNotifier(e.Session, message)
		}
	}

	isMentioned := false
	for _, user := range message.Mentions {
		if user.ID == e.Session.State.User.ID {
			isMentioned = true
			break
		}
	}
	if isMentioned {
		logging.UI.GhostPingDetection(e.Session, message)
	}
}

func (e *Events) ensureUserEntry(userID string) SBDScore {
	if _, ok := e.Settings.SelfbotDetectionScore[userID]; !ok {
		e.Settings.SelfbotDetectionScore[userID] = SBDScore{UserID: userID, Score: 0}
	}
	return e.Settings.SelfbotDetectionScore[userID]
}

func (e *Events) logSelfbotDetection(score SBDScore, message *discordgo.Message) {
	if score.Score >= 15 && !score.Logged {
		logging.UI.SelfbotDetection(e.Session, message, "")
		logging.Info("Selfbot detection triggered for user " + message.Author.String())
		score.Logged = true
	}
}

var sbDetectables = []SBDetectable{
	{
		String:  "[lightningbot]",
		Selfbot: "LightningBot",
		Score:   1,
	},
	{
		String:  "www.team-luna.org",
		Selfbot: "Luna",
		Score:   2,
	},
	{
		String:  "enabled protections » ",
		Selfbot: "Luna",
		Score:   2,
	},
	{
		String:  "nighty.one",
		Selfbot: "Nighty",
		Score:   2,
	},
	{
		String:  "<> is required",
		Selfbot: "",
		Score:   2,
	},
	{
		String:  "<> = required",
		Selfbot: "",
		Score:   2,
	},
	{
		String:  "> **",
		Selfbot: "",
		Score:   0.5,
	},
}

func (e *Events) checkSelfbotDetection(message *discordgo.Message, cfg *config.Config) {
	if !cfg.Events.SelfbotDetection || message.Author.Bot { // check for state user
		return
	}
	score := e.ensureUserEntry(message.Author.ID)
	for _, detectable := range sbDetectables {
		if strings.Contains(strings.ToLower(message.Content), detectable.String) {
			score.Score += detectable.Score
			score.Events = append(score.Events, SBDEvent{
				Event:   "detectable",
				Message: message,
			})
			score.Selfbot = utils.Misc.Or(detectable.Selfbot, score.Selfbot)
			e.logSelfbotDetection(score, message)
		}
	}
}

func (e *Events) checkSelfbotDetectionEdit(message *discordgo.Message, cfg *config.Config) {
	if !cfg.Events.SelfbotDetection || message.Author.Bot { // check for state user
		return
	}
	e.ensureUserEntry(message.Author.ID)
}

func (e *Events) checkSelfbotDetectionDelete(message *discordgo.Message, cfg *config.Config) {
	if !cfg.Events.SelfbotDetection || message.Author.Bot { // check for state user
		return
	}
	score := e.ensureUserEntry(message.Author.ID)
	timeDiff := time.Since(utils.Discord.SnowflakeToTime(message.ID)) - e.Settings.discordTimeDiff
	if timeDiff < time.Second*1 {
		score.Events = append(score.Events, SBDEvent{Event: "fast_delete", Message: message})
		score.Score += 0.5
		e.logSelfbotDetection(score, message)
	}
	e.Settings.SelfbotDetectionScore[message.Author.ID] = score
}
