package eventhandlers

import (
	"net/url"
	"strings"
	"time"
	"unicode"

	"github.com/LightningDev1/LB-Selfbot/api"
	"github.com/LightningDev1/LB-Selfbot/config"
	"github.com/LightningDev1/LB-Selfbot/constants"
	"github.com/LightningDev1/LB-Selfbot/embed"
	"github.com/LightningDev1/LB-Selfbot/http"
	"github.com/LightningDev1/LB-Selfbot/logging"
	"github.com/LightningDev1/LB-Selfbot/snipers"
	"github.com/LightningDev1/LB-Selfbot/utils"
	"github.com/LightningDev1/dgc"
	"github.com/LightningDev1/discordgo"
)

func (e *Events) OnMessage(_ *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author == nil || m.Message == nil {
		utils.Logging.Warn("Message: Author or Message is nil: ", m.Author, m.Message)
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
		e.Settings.CountryCode,
		e.Settings.PaymentSourceID,
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
		utils.Logging.Warn("MessageEdit: Author or Message is nil: ", m.Author, m.Message)
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
		e.Settings.CountryCode,
		e.Settings.PaymentSourceID,
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
				logging.Logging.PhisingDetected(e.Session, message, hostname)
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
			httpResponse := http.Get("https://evilinsult.com/generate_insult.php?lang=en")
			_, _ = e.Session.ChannelMessageSend(message.ChannelID, httpResponse.Body)
		}
		if utils.SliceContains(EventHandlers.Settings.PingPingServers, message.GuildID) {
			_, _ = e.Session.ChannelMessageSend(message.ChannelID, message.Author.Mention())
		}
	}
}

func (e *Events) checkAFK(message *discordgo.Message, cfg *config.Config) {
	guild, err := e.Session.State.Guild(message.GuildID)
	guildName := ""
	if err == nil {
		guildName = guild.Name
	}
	if message.Author.ID == e.Session.State.User.ID ||
		!cfg.AFK.Enabled ||
		utils.SliceContains(cfg.AFK.IgnoredUsers, message.Author.ID) ||
		utils.SliceContains(cfg.AFK.IgnoredUsers, message.Author.String()) ||
		utils.SliceContains(cfg.AFK.IgnoredServers, message.GuildID) ||
		utils.SliceContains(cfg.AFK.IgnoredServers, guildName) {
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
				_ = api.ChangeSettings(e.Session, map[string]interface{}{"restricted_guilds": guildIDs}, true)
				logging.Logging.MassDMBlocked()
				go func() {
					time.Sleep(5 * 60)
					_ = api.ChangeSettings(e.Session, map[string]interface{}{"restricted_guilds": []any{}}, true)
				}()
			}
		}
	}

	guild, err := e.Session.State.Guild(message.GuildID)
	guildName := ""
	if err == nil {
		guildName = guild.Name
	}

	if message.GuildID != "" && (utils.SliceContains(cfg.Protections.ServerProtections.ProtectedServers, message.GuildID) ||
		utils.SliceContains(cfg.Protections.ServerProtections.ProtectedServers, guildName)) {
		userPerms, _ := e.Session.UserChannelPermissions(message.Author.ID, message.ChannelID)
		if userPerms&discordgo.PermissionManageMessages != 0 || userPerms&discordgo.PermissionAdministrator != 0 {
			return
		}
		if cfg.Protections.ServerProtections.AntiInvite && constants.INVITE_REGEX.MatchString(message.Content) {
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
		logging.Logging.AntiDelete(e.Session, message)
	}
}

func (e *Events) checkNotifiers(message *discordgo.Message, cfg *config.Config) {
	if message.Author.ID == e.Session.State.User.ID {
		return
	}

	if message.GuildID == "" {
		if cfg.Events.DMNotifier {
			logging.Logging.DMNotifier(message)
		}
	}

	if cfg.Events.DiscordStaffDetection && message.Author.PublicFlags&discordgo.UserFlagDiscordEmployee != 0 {
		logging.Logging.DiscordStaffDetection(e.Session, message)
	}
}

func (e *Events) checkNotifiersEdit(message *discordgo.Message, cfg *config.Config) {
	if message.Author.ID == e.Session.State.User.ID {
		return
	}

	if message.GuildID == "" {
		if cfg.Events.DMEditNotifier {
			logging.Logging.DMEditNotifier(message)
		}

	}
}

func (e *Events) checkNotifiersDelete(message *discordgo.Message, cfg *config.Config) {
	if message.Author.ID == e.Session.State.User.ID {
		return
	}

	if message.GuildID == "" {
		if cfg.Events.DMDeleteNotifier {
			logging.Logging.DMDeleteNotifier(message)
		}
	} else {
		if cfg.Events.MessageDeleteNotifier {
			logging.Logging.MessageDeleteNotifier(e.Session, message)
		}
	}

	isMentioned := false
	for _, user := range message.Mentions {
		if user.ID == e.Session.State.User.ID {
			isMentioned = true
			break
		}
	}
	if isMentioned && cfg.Events.GhostPingDetection {
		logging.Logging.GhostPingDetection(e.Session, message)
	}
}

func (e *Events) ensureUserEntry(userID string) *SBDScore {
	if _, ok := e.Settings.SelfbotDetectionScore[userID]; !ok {
		e.Settings.SelfbotDetectionScore[userID] = SBDScore{UserID: userID, Score: 0}
	}
	a := e.Settings.SelfbotDetectionScore[userID]
	return &a
}

func (e *Events) logSelfbotDetection(score *SBDScore, message *discordgo.Message) {
	if score.Score >= 15 && !score.Logged {
		selfbot := utils.Misc.Or(score.Selfbot, "Unknown")
		logging.Logging.SelfbotDetection(e.Session, message, selfbot)
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
		String:  "[lightningbot free]",
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
		Score:   1,
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
	if !cfg.Events.SelfbotDetection || message.Author.Bot || message.Author.ID == e.Session.State.User.ID {
		return
	}

	score := e.ensureUserEntry(message.Author.ID)
	for _, detectable := range sbDetectables {
		if !strings.Contains(strings.ToLower(message.Content), detectable.String) {
			continue
		}

		score.Score += detectable.Score
		score.Events = append(score.Events, SBDEvent{
			Event:   "detectable",
			Message: message,
		})
		score.Selfbot = utils.Misc.Or(detectable.Selfbot, score.Selfbot)
		e.logSelfbotDetection(score, message)
	}
	e.Settings.SelfbotDetectionScore[message.Author.ID] = *score
}

func (e *Events) checkSelfbotDetectionEdit(message *discordgo.Message, cfg *config.Config) {
	if !cfg.Events.SelfbotDetection || message.Author.Bot || message.Author.ID == e.Session.State.User.ID {
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
	e.Settings.SelfbotDetectionScore[message.Author.ID] = *score
}
