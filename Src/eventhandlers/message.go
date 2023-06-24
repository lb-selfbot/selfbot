package eventhandlers

import (
	"net/url"
	"runtime"
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
	if m.Message == nil || m.Author == nil {
		utils.Logging.Debug("Message: Author or Message is nil: ", m.Author, m.Message)
		return
	}

	if m.Type != discordgo.MessageTypeDefault && m.Type != discordgo.MessageTypeReply && m.Type != discordgo.MessageTypeChatInputCommand {
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
	e.checkPingReactions(m.Message)
	e.checkEventReactions(m.Message)
	e.checkAFK(m.Message, &cfg)
	e.checkProtections(m.Message, &cfg)
	e.checkNotifiers(m.Message, &cfg)
	e.checkSelfbotDetection(m.Message, &cfg)
}

func (e *Events) OnMessageEdit(_ *discordgo.Session, m *discordgo.MessageUpdate) {
	message, err := e.Session.State.Message(m.Message.ChannelID, m.Message.ID)
	if err != nil {
		m.Message = message
	}

	if m.Message == nil {
		utils.Logging.Debug("MessageEdit: Message is nil: ", m.Message)
		return
	}

	if m.Author == nil {
		utils.Logging.Debug("MessageEdit: Author is nil: ", m.Author)
		return
	}

	if m.BeforeUpdate == nil {
		m.BeforeUpdate = message
	}

	if m.Type != discordgo.MessageTypeDefault && m.Type != discordgo.MessageTypeReply && m.Type != discordgo.MessageTypeChatInputCommand {
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
	
	snipers.CheckGiveawaySniperEdit(e.Session, m.Message, &cfg, "")

	e.checkIgnore(m.Message)
	e.checkPhising(m.Message, &cfg)
	e.checkNotifiersEdit(m.Message, &cfg)
	e.checkSelfbotDetectionEdit(m.Message, &cfg)

	if strings.HasPrefix(m.Message.Content, cfg.Selfbot.CommandPrefix) {
		return
	}

	messages, ok := e.Settings.EditedMessages[m.ChannelID]
	if ok {
		e.Settings.EditedMessages[m.ChannelID] = append(messages, m)
	} else {
		e.Settings.EditedMessages[m.ChannelID] = []*discordgo.MessageUpdate{m}
	}
}

func (e *Events) OnMessageDelete(_ *discordgo.Session, m *discordgo.MessageDelete) {
	if m.BeforeDelete != nil {
		m.Message = m.BeforeDelete
	}

	if m.Message == nil {
		utils.Logging.Debug("MessageDelete: Message is nil: ", m.Message)
		return
	}

	if m.Author == nil {
		utils.Logging.Debug("MessageDelete: Author is nil: ", m.Author)
		return
	}

	if m.Type != discordgo.MessageTypeDefault && m.Type != discordgo.MessageTypeReply {
		return
	}

	cfg, err := config.GetConfig()
	if err != nil {
		return
	}

	e.checkProtectionsDelete(m.Message, &cfg)
	e.checkNotifiersDelete(m.Message, &cfg)
	e.checkSelfbotDetectionDelete(m.Message, &cfg)

	if strings.HasPrefix(m.Message.Content, cfg.Selfbot.CommandPrefix) {
		return
	}

	messages, ok := e.Settings.DeletedMessages[m.ChannelID]
	if ok {
		e.Settings.DeletedMessages[m.ChannelID] = append(messages, m.Message)
	} else {
		e.Settings.DeletedMessages[m.ChannelID] = []*discordgo.Message{m.Message}
	}
}

func (e *Events) checkIgnore(message *discordgo.Message) {
	for _, user := range e.Settings.IgnoredUsers {
		if message.Author.ID != user.ID {
			continue
		}

		_ = e.Session.MessageAck(message.ChannelID, message.ID)
		break
	}

	for _, guild := range e.Settings.IgnoredGuilds {
		if message.GuildID != guild.ID {
			continue
		}

		_ = e.Session.GuildAck(message.GuildID)
		break
	}
}

func (e *Events) checkPhising(message *discordgo.Message, cfg *config.Config) {
	if !cfg.Protections.AntiPhising {
		return
	}

	for _, part := range strings.Split(message.Content, " ") {
		link, err := url.Parse(part)
		if err != nil {
			continue
		}

		hostname := link.Hostname()

		isPhising := utils.SliceContains(e.Settings.phisingDomains, hostname)
		isSuspicious := utils.SliceContains(e.Settings.suspiciousDomains, hostname)

		if isPhising || isSuspicious {
			_ = e.Session.MessageAck(message.ChannelID, message.ID)
			logging.Logging.PhisingDetected(e.Session, message, hostname)
		}
	}
}

func (e *Events) checkPingReactions(message *discordgo.Message) {
	if message.Author.ID == e.Session.State.User.ID || message.GuildID == "" {
		return
	}

	mentionsUser := strings.Contains(message.Content, e.Session.State.User.Mention())

	if !mentionsUser {
		return
	}

	member, err := utils.Discord.GetMember(e.Session, message.GuildID, message.Author.ID)
	if err != nil {
		return
	}

	if utils.Discord.IsGuildInList(e.Session, message.GuildID, EventHandlers.Settings.PingMuteServers) {
		_ = utils.Discord.MuteUser(e.Session, member)
	}

	if utils.Discord.IsGuildInList(e.Session, message.GuildID, EventHandlers.Settings.PingTimeoutServers) {
		until := time.Now().Add(5 * time.Minute)
		_ = e.Session.GuildMemberTimeout(member.GuildID, member.User.ID, &until)
	}

	if utils.Discord.IsGuildInList(e.Session, message.GuildID, EventHandlers.Settings.PingBanServers) {
		_ = e.Session.GuildBanCreateWithReason(member.GuildID, member.User.ID, "Ping Ban Event Reaction", 0)
	}

	if utils.Discord.IsGuildInList(e.Session, message.GuildID, EventHandlers.Settings.PingKickServers) {
		_ = e.Session.GuildMemberDeleteWithReason(member.GuildID, member.User.ID, "Ping Kick Event Reaction")
	}

	if utils.Discord.IsGuildInList(e.Session, message.GuildID, EventHandlers.Settings.PingInsultServers) {
		httpResponse := http.Get("https://evilinsult.com/generate_insult.php?lang=en")
		_, _ = e.Session.ChannelMessageSend(message.ChannelID, httpResponse.Body)
	}

	if utils.Discord.IsGuildInList(e.Session, message.GuildID, EventHandlers.Settings.PingPingServers) {
		_, _ = e.Session.ChannelMessageSend(message.ChannelID, message.Author.Mention())
	}
}

func (e *Events) checkEventReactions(message *discordgo.Message) {
	for _, event := range e.Settings.Copy {
		if event.GuildID != message.GuildID && event.ChannelID != message.ChannelID {
			continue
		}

		if event.UserID != message.Author.ID {
			continue
		}

		e.Session.ChannelMessageSendComplex(message.ChannelID, &discordgo.MessageSend{
			Content:         message.Content,
			AllowedMentions: &discordgo.MessageAllowedMentions{},
		})
	}

	for _, event := range e.Settings.Mock {
		if event.GuildID != message.GuildID && event.ChannelID != message.ChannelID {
			continue
		}

		if event.UserID != message.Author.ID {
			continue
		}

		e.Session.ChannelMessageSendComplex(message.ChannelID, &discordgo.MessageSend{
			Content:         utils.String.Mock(message.Content),
			AllowedMentions: &discordgo.MessageAllowedMentions{},
		})
	}

	for _, event := range e.Settings.ReplyCopy {
		if event.GuildID != message.GuildID && event.ChannelID != message.ChannelID {
			continue
		}

		if event.UserID != message.Author.ID {
			continue
		}

		e.Session.ChannelMessageSendComplex(message.ChannelID, &discordgo.MessageSend{
			Content:         message.Content,
			AllowedMentions: &discordgo.MessageAllowedMentions{RepliedUser: true},
			Reference: &discordgo.MessageReference{
				MessageID: message.ID,
				ChannelID: message.ChannelID,
				GuildID:   message.GuildID,
			},
		})
	}

	for _, event := range e.Settings.PinSpam {
		if event.GuildID != message.GuildID && event.ChannelID != message.ChannelID {
			continue
		}

		if event.UserID != message.Author.ID {
			continue
		}

		err := e.Session.ChannelMessagePin(message.ChannelID, message.ID)
		if err == nil {
			e.Settings.PinSpamMessages = append(e.Settings.PinSpamMessages, message)
		}
	}

	for _, event := range e.Settings.DeleteMessages {
		if event.GuildID != message.GuildID && event.ChannelID != message.ChannelID {
			continue
		}

		if event.UserID != message.Author.ID {
			continue
		}

		_ = e.Session.ChannelMessageDelete(message.ChannelID, message.ID)
	}

	for _, event := range e.Settings.ReactUser {
		if event.GuildID != message.GuildID && event.ChannelID != message.ChannelID {
			continue
		}

		if event.UserID != message.Author.ID {
			continue
		}

		_ = e.Session.MessageReactionAdd(message.ChannelID, message.ID, event.Argument)
	}

}

func (e *Events) checkAFK(message *discordgo.Message, cfg *config.Config) {
	if !cfg.AFK.Enabled || message.Author.ID == e.Session.State.User.ID {
		return
	}

	if !cfg.AFK.TriggerList.IsTriggered(e.Session, message.Author.ID, message.GuildID) {
		return
	}

	mentionsUser := strings.Contains(message.Content, e.Session.State.User.Mention())

	activatedOnPing := cfg.AFK.OnPing && mentionsUser
	activatedOnDM := cfg.AFK.OnDM && message.GuildID == ""

	if !activatedOnPing && !activatedOnDM {
		return
	}

	if runtime.GOOS == "windows" && cfg.AFK.MustBeInactive {
		idleTimeout := time.Duration(cfg.AFK.InactiveTimeout) * time.Second

		idleTime, err := utils.System.GetIdleTime()
		if err != nil || idleTime < idleTimeout {
			return
		}
	}

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

func (e *Events) checkProtections(message *discordgo.Message, cfg *config.Config) {
	if message.Author.ID == e.Session.State.User.ID {
		return
	}

	if cfg.Protections.AntiMassDM && message.GuildID == "" {
		messages, _ := utils.Discord.GetChannelMessages(e.Session, message.ChannelID, 100)

		directMessages := utils.Filter(messages, func(message *discordgo.Message) bool {
			return message.Author.ID != e.Session.State.User.ID
		})

		// There are more than 2 messages from this user,
		// so this is not a new DM
		if len(directMessages) > 2 {
			return
		}

		e.Settings.newDMs = append(e.Settings.newDMs, message)

		dmsInLastMinute := utils.Filter(e.Settings.newDMs, func(message *discordgo.Message) bool {
			return time.Since(message.Timestamp) < time.Minute
		})

		if len(dmsInLastMinute) < 10 {
			return
		}

		_ = e.Session.MessageAck(message.ChannelID, message.ID)

		var guildIDs []string
		for _, g := range e.Session.State.Guilds {
			guildIDs = append(guildIDs, g.ID)
		}

		_ = api.ChangeSettings(e.Session, map[string]any{"restricted_guilds": guildIDs}, true)

		logging.Logging.MassDMBlocked()

		go func() {
			// Enable DMs after 5 minutes
			time.Sleep(5 * 60)

			_ = api.ChangeSettings(e.Session, map[string]any{"restricted_guilds": []any{}}, true)
		}()
	}

	isGuildProtected := utils.Discord.IsGuildInList(
		e.Session,
		message.GuildID,
		cfg.Protections.ServerProtections.ProtectedServers,
	)

	if isGuildProtected {
		userPerms, _ := e.Session.UserChannelPermissions(message.Author.ID, message.ChannelID)

		userHasManageMessages := userPerms&discordgo.PermissionManageMessages != 0
		userIsAdmin := userPerms&discordgo.PermissionAdministrator != 0

		if userHasManageMessages || userIsAdmin {
			return
		}

		ctx := &dgc.Ctx{
			Session: e.Session,
			Event:   &discordgo.MessageCreate{Message: message},
		}

		if cfg.Protections.ServerProtections.AntiInvite && constants.INVITE_REGEX.MatchString(message.Content) {
			err := e.Session.ChannelMessageDelete(message.ChannelID, message.ID)
			if err != nil {
				return
			}

			embed.NewEmbed().
				SetDescription("Invites are not allowed in this server.").
				Send(ctx)
		}

		if cfg.Protections.ServerProtections.AntiSpam {
			messages, _ := utils.Discord.GetChannelMessages(e.Session, message.ChannelID, 100)

			authorMessagesInLastMinute := utils.Filter(messages, func(m *discordgo.Message) bool {
				return m.Author.ID == message.Author.ID && time.Since(message.Timestamp) < time.Minute
			})

			if len(authorMessagesInLastMinute) > 10 {
				err := e.Session.ChannelMessageDelete(message.ChannelID, message.ID)
				if err != nil {
					return
				}

				embed.NewEmbed().
					SetDescription("Spamming is not allowed in this server.").
					Send(ctx)
			}
		}

		if cfg.Protections.ServerProtections.AntiUpper && len(message.Content) > 10 {
			contentRunes := []rune(message.Content)

			upperRunes := utils.Filter(contentRunes, func(char rune) bool {
				return unicode.IsUpper(char)
			})

			upperPercentage := float64(len(upperRunes)) / float64(len(message.Content)) * 100

			if upperPercentage >= 70 {
				err := e.Session.ChannelMessageDelete(message.ChannelID, message.ID)
				if err != nil {
					return
				}

				embed.NewEmbed().
					SetDescription("Uppercase spam is not allowed in this server.").
					Send(ctx)
			}
		}
	}
}

func (e *Events) checkProtectionsDelete(message *discordgo.Message, cfg *config.Config) {
	isGuildProtected := utils.Discord.IsGuildInList(
		e.Session,
		message.GuildID,
		cfg.Protections.ServerProtections.ProtectedServers,
	)

	if isGuildProtected && cfg.Protections.ServerProtections.AntiDelete {
		userPerms, _ := e.Session.UserChannelPermissions(message.Author.ID, message.ChannelID)

		userHasManageMessages := userPerms&discordgo.PermissionManageMessages != 0
		userIsAdmin := userPerms&discordgo.PermissionAdministrator != 0

		if userHasManageMessages || userIsAdmin {
			return
		}

		logging.Logging.AntiDelete(e.Session, message)
	}
}

func (e *Events) checkNotifiers(message *discordgo.Message, cfg *config.Config) {
	if message.Author.ID == e.Session.State.User.ID {
		return
	}

	if message.GuildID == "" && cfg.Events.DMNotifier {
		logging.Logging.DMNotifier(message)
	}

	authorIsEmployee := message.Author.PublicFlags&discordgo.UserFlagDiscordEmployee != 0

	if authorIsEmployee && cfg.Events.DiscordStaffDetection {
		logging.Logging.DiscordStaffDetection(e.Session, message)
	}
}

func (e *Events) checkNotifiersEdit(message *discordgo.Message, cfg *config.Config) {
	if message.Author.ID == e.Session.State.User.ID {
		return
	}

	if message.GuildID == "" && cfg.Events.DMEditNotifier {
		logging.Logging.DMEditNotifier(message)
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

	mentionsUser := strings.Contains(message.Content, e.Session.State.User.Mention())

	if mentionsUser && cfg.Events.GhostPingDetection {
		logging.Logging.GhostPingDetection(e.Session, message)
	}
}

func (e *Events) ensureUserEntry(userID string) *SBDScore {
	rawScore, _ := e.Settings.SelfbotDetectionScore.LoadOrStore(userID, SBDScore{UserID: userID, Score: 0})
	score := rawScore.(SBDScore)
	return &score
}

func (e *Events) logSelfbotDetection(score *SBDScore, message *discordgo.Message) {
	if score.Score < 15 || score.Logged {
		return
	}

	selfbot := utils.Misc.Or(score.Selfbot, "Unknown")
	logging.Logging.SelfbotDetection(e.Session, message, selfbot)
	score.Logged = true
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

	e.Settings.SelfbotDetectionScore.Store(message.Author.ID, *score)
}

func (e *Events) checkSelfbotDetectionEdit(message *discordgo.Message, cfg *config.Config) {
	if !cfg.Events.SelfbotDetection || message.Author.Bot || message.Author.ID == e.Session.State.User.ID {
		return
	}

	e.ensureUserEntry(message.Author.ID)
}

func (e *Events) checkSelfbotDetectionDelete(message *discordgo.Message, cfg *config.Config) {
	if !cfg.Events.SelfbotDetection || message.Author.Bot || message.Author.ID == e.Session.State.User.ID {
		return
	}

	score := e.ensureUserEntry(message.Author.ID)
	timeDiff := time.Since(utils.Discord.SnowflakeToTime(message.ID)) - e.Settings.discordTimeDiff
	if timeDiff < time.Second*1 {
		score.Events = append(score.Events, SBDEvent{Event: "fast_delete", Message: message})
		score.Score += 0.5
		e.logSelfbotDetection(score, message)
	}

	e.Settings.SelfbotDetectionScore.Store(message.Author.ID, *score)
}
