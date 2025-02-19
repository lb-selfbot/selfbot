package impl

import (
	"strings"

	"github.com/LightningDev1/LB-Selfbot/internal/services/interfaces"
	"github.com/LightningDev1/LB-Selfbot/selfbot/config"
	"github.com/LightningDev1/LB-Selfbot/selfbot/logging"

	"github.com/LightningDev1/LB-Selfbot/selfbot/utils"
	"github.com/LightningDev1/discordgo"
)

type Notifiers struct{}

func NewNotifiers() *Notifiers {
	return &Notifiers{}
}

func (n *Notifiers) Type() interfaces.EventReactionType {
	return interfaces.EventReactionNotifiers
}

func (n *Notifiers) Description() string {
	return "Notifiers (default)"
}

func (n *Notifiers) Equals(other interfaces.EventReaction) bool {
	_, ok := other.(*Notifiers)

	return ok
}

func (n *Notifiers) Data() map[string]any {
	return map[string]any{}
}

func (n *Notifiers) Handler(session *discordgo.Session, event any) {
	cfg, err := config.GetConfig()
	if err != nil {
		utils.Logging.Debug("Notifiers.Handler: Error getting config: " + err.Error())
		return
	}

	switch event := event.(type) {
	case *discordgo.MessageCreate:
		n.MessageCreateHandler(session, event, &cfg)
	case *discordgo.MessageUpdate:
		n.MessageUpdateHandler(session, event, &cfg)
	case *discordgo.MessageDelete:
		n.MessageDeleteHandler(session, event, &cfg)
	case *discordgo.GuildBanAdd:
		n.GuildBanHandler(session, event, &cfg)
	case *discordgo.ChannelCreate:
		n.ChannelCreateHandler(session, event, &cfg)
	case *discordgo.TypingStart:
		n.TypingStartHandler(session, event, &cfg)
	default:
		utils.Logging.Debug("Notifiers.Handler: Event is not a notifiers event")
	}
}

// MessageCreateHandler: DMNotifier, DiscordStaffDetection
func (n *Notifiers) MessageCreateHandler(session *discordgo.Session, message *discordgo.MessageCreate, cfg *config.Config) {
	if message.Author.ID == session.State.User.ID {
		return
	}

	isDM := utils.Discord.IsDM(session, message.Message)

	// DM Notifier
	if isDM && cfg.Events.DMNotifier {
		logging.Logging.DMNotifier(message.Message)
	}

	// Discord Staff Detection
	authorIsEmployee := message.Author.PublicFlags&discordgo.UserFlagDiscordEmployee != 0

	if authorIsEmployee && cfg.Events.DiscordStaffDetection {
		logging.Logging.DiscordStaffDetection(session, message.Message)
	}
}

// MessageUpdateHandler: DMEditNotifier
func (n *Notifiers) MessageUpdateHandler(session *discordgo.Session, message *discordgo.MessageUpdate, cfg *config.Config) {
	if message.Author.ID == session.State.User.ID {
		return
	}

	isDM := utils.Discord.IsDM(session, message.Message)
	embedded := message.Embeds != nil && len(message.Embeds) > 0

	//  DM Edit Notifier
	if isDM && !embedded && cfg.Events.DMEditNotifier {
		logging.Logging.DMEditNotifier(message.Message, message.BeforeUpdate.Content, message.Content)
	}
}

// MessageDeleteHandler: DMDeleteNotifier, MessageDeleteNotifier, GhostPingDetection
func (n *Notifiers) MessageDeleteHandler(session *discordgo.Session, message *discordgo.MessageDelete, cfg *config.Config) {
	if message.Author.ID == session.State.User.ID {
		return
	}

	isDM := utils.Discord.IsDM(session, message.Message)
	dmNotifier := cfg.Events.DMDeleteNotifier && isDM
	messageNotifier := cfg.Events.MessageDeleteNotifier && !isDM

	// DM Delete Notifier, Message Delete Notifier
	func() {
		if !dmNotifier && !messageNotifier {
			return
		}

		var attachments []string

		for _, attachment := range message.Attachments {
			if attachment.Size > 100*1024*1024 {
				// > 100MB: only show the URL
				attachments = append(attachments, attachment.Filename+": "+attachment.URL+" "+attachment.ProxyURL)
				continue
			}

			filepath, err := utils.Discord.SaveAttachment(attachment)
			if err != nil {
				utils.Logging.Error("Notifiers.DeleteHandler: Error saving attachment: " + err.Error())
				attachments = append(attachments, attachment.Filename+": "+attachment.URL+" "+attachment.ProxyURL)
			}

			attachments = append(attachments, attachment.Filename+": "+filepath)
		}

		if dmNotifier {
			logging.Logging.DMDeleteNotifier(message.Message, attachments)
		} else {
			logging.Logging.MessageDeleteNotifier(session, message.Message, attachments)
		}
	}()

	// Ghost Ping Detection
	mentionsUser := strings.Contains(message.Content, session.State.User.Mention())

	if mentionsUser && cfg.Events.GhostPingDetection {
		logging.Logging.GhostPingDetection(session, message.Message)
	}
}

// GuildBanHandler: GuildBanNotifier
func (n *Notifiers) GuildBanHandler(session *discordgo.Session, ban *discordgo.GuildBanAdd, cfg *config.Config) {
	if ban.User == nil {
		return
	}

	if ban.User.ID != session.State.User.ID {
		return
	}

	// Ban Notifier
	if cfg.Events.BanNotifier {
		logging.Logging.BanNotifier(session, ban.GuildID)
	}
}

// ChannelCreateHandler: TicketNotifier
func (n *Notifiers) ChannelCreateHandler(session *discordgo.Session, channel *discordgo.ChannelCreate, cfg *config.Config) {
	// Ticket Notifier
	if channel.Type == discordgo.ChannelTypeGuildText && strings.Contains(channel.Name, "ticket") {
		perms, err := session.State.UserChannelPermissions(channel.ID, session.State.User.ID)
		if err != nil {
			return
		}

		canSendMessages := perms&discordgo.PermissionSendMessages != 0
		canViewChannel := perms&discordgo.PermissionViewChannel != 0

		if canSendMessages && canViewChannel {
			logging.Logging.TicketNotifier(session, channel.Channel)
		}
	}
}

// TypingStartHandler: DMTypingNotifier
func (n *Notifiers) TypingStartHandler(session *discordgo.Session, event *discordgo.TypingStart, cfg *config.Config) {
	isDM := event.GuildID == ""

	// DM Typing Notifier
	if cfg.Events.DMTypingNotifier && isDM {
		username := utils.Discord.GetUserName(session, event.UserID)
		if username == "Unknown" {
			return
		}

		logging.Logging.DMTypingNotifier(username)
	}
}
