package impl

import (
	"strconv"
	"time"
	"unicode"

	"github.com/lb-selfbot/selfbot/internal/api"
	"github.com/lb-selfbot/selfbot/internal/logging"

	"github.com/lb-selfbot/selfbot/internal/constants"
	"github.com/lb-selfbot/selfbot/internal/services/config"
	"github.com/lb-selfbot/selfbot/internal/services/eventbroker"
	"github.com/lb-selfbot/selfbot/internal/services/eventreactions/types"
	"github.com/lb-selfbot/selfbot/internal/services/ui/featurelogs"

	"github.com/lb-selfbot/discordgo"
	"github.com/lb-selfbot/selfbot/internal/utils"
)

type friendRequest struct {
	*discordgo.Relationship

	Time time.Time
}

type Protections struct {
	eventBroker eventbroker.Service
	configSvc   config.Service

	NewDMs         []*discordgo.Message
	FriendRequests []*friendRequest
	LastBlock      time.Time
}

func NewProtections(eventBroker eventbroker.Service, configSvc config.Service) *Protections {
	return &Protections{eventBroker: eventBroker, configSvc: configSvc}
}

func (p *Protections) Type() types.EventReactionType {
	return types.Protections
}

func (p *Protections) Description() string {
	return "Protections (default)"
}

func (p *Protections) Equals(other types.EventReaction) bool {
	_, ok := other.(*Protections)

	return ok
}

func (p *Protections) Data() map[string]any {
	return map[string]any{}
}

func (p *Protections) Handler(session *discordgo.Session, event any) {
	cfg := p.configSvc.GetConfig()

	switch event := event.(type) {
	case *discordgo.MessageCreate:
		p.MessageCreateHandler(session, event, cfg)
	case *discordgo.MessageDelete:
		p.MessageDeleteHandler(session, event, cfg)
	case *discordgo.RelationshipAdd:
		p.RelationshipAddHandler(session, event, cfg)
	default:
		logging.Debug("Protections.Handler: Event is not a protections event")
	}
}

func (p *Protections) MessageCreateHandler(session *discordgo.Session, message *discordgo.MessageCreate, cfg *config.Config) {
	if message.Author.ID == session.State.User.ID {
		return
	}

	isDM := utils.Discord.IsDM(session, message.Message)

	if cfg.Protections.AntiMassDM && time.Since(p.LastBlock) > 5*time.Minute && isDM {
		messages, _ := utils.Discord.GetChannelMessages(session, message.ChannelID, 100)

		directMessages := utils.Filter(messages, func(message *discordgo.Message) bool {
			return message.Author != nil && message.Author.ID != session.State.User.ID
		})

		// There are more than 2 messages from this user, so this is not a new DM
		if len(directMessages) > 2 {
			return
		}

		p.NewDMs = append(p.NewDMs, message.Message)

		newDMsInLastMinute := utils.Filter(p.NewDMs, func(message *discordgo.Message) bool {
			return time.Since(message.Timestamp) < time.Minute
		})

		if len(newDMsInLastMinute) < 10 {
			return
		}

		_ = session.MessageAck(message.ChannelID, message.ID)

		var guildIDs []uint64
		for _, g := range session.State.Guilds {
			id, _ := strconv.Atoi(g.ID)

			guildIDs = append(guildIDs, uint64(id))
		}

		session.State.UserSettings.Privacy.RestrictedGuildIds = guildIDs
		err := session.UserSettingsUpdate(session.State.UserSettings, true)
		if err != nil {
			return
		}

		p.LastBlock = time.Now()

		p.eventBroker.Publish(eventbroker.EventLogFeature, &featurelogs.MassDMBlocked{})

		go func() {
			// Enable DMs after 5 minutes
			time.Sleep(5 * 60)

			session.State.UserSettings.Privacy.RestrictedGuildIds = []uint64{}
			err := session.UserSettingsUpdate(session.State.UserSettings, true)
			if err != nil {
				return
			}
		}()
	}

	isGuildProtected := utils.Discord.IsGuildInList(
		session,
		message.GuildID,
		cfg.Protections.ServerProtections.ProtectedServers,
	)

	if isGuildProtected {
		userPerms, _ := session.UserChannelPermissions(message.Author.ID, message.ChannelID)

		userHasManageMessages := userPerms&discordgo.PermissionManageMessages != 0
		userIsAdmin := userPerms&discordgo.PermissionAdministrator != 0

		if userHasManageMessages || userIsAdmin {
			return
		}

		if cfg.Protections.ServerProtections.AntiInvite && constants.INVITE_REGEX.MatchString(message.Content) {
			err := session.ChannelMessageDelete(message.ChannelID, message.ID)
			if err != nil {
				return
			}

			msg, err := session.ChannelMessageSend(message.ChannelID, "Invites are not allowed in this server.")
			if err != nil {
				return
			}

			p.eventBroker.Publish(eventbroker.EventHandleReactionMessage, msg)
		}

		if cfg.Protections.ServerProtections.AntiSpam {
			messages, _ := utils.Discord.GetChannelMessages(session, message.ChannelID, 100)

			authorMessagesInLast30Sec := utils.Filter(messages, func(m *discordgo.Message) bool {
				return m.Author.ID == message.Author.ID && time.Since(message.Timestamp) < time.Second*30
			})

			if len(authorMessagesInLast30Sec) > 10 {
				err := session.ChannelMessageDelete(message.ChannelID, message.ID)
				if err != nil {
					return
				}

				msg, err := session.ChannelMessageSend(message.ChannelID, "Spamming is not allowed in this server.")
				if err != nil {
					return
				}

				p.eventBroker.Publish(eventbroker.EventHandleReactionMessage, msg)
			}
		}

		if cfg.Protections.ServerProtections.AntiUpper && len(message.Content) > 10 {
			contentRunes := []rune(message.Content)

			upperRunes := utils.Filter(contentRunes, func(char rune) bool {
				return unicode.IsUpper(char)
			})

			upperPercentage := float64(len(upperRunes)) / float64(len(message.Content)) * 100

			if upperPercentage >= 70 {
				err := session.ChannelMessageDelete(message.ChannelID, message.ID)
				if err != nil {
					return
				}

				msg, err := session.ChannelMessageSend(message.ChannelID, "Uppercase spam is not allowed in this server.")
				if err != nil {
					return
				}

				p.eventBroker.Publish(eventbroker.EventHandleReactionMessage, msg)
			}
		}
	}
}

func (p *Protections) MessageDeleteHandler(session *discordgo.Session, message *discordgo.MessageDelete, cfg *config.Config) {
	if message.Author.ID == session.State.User.ID {
		return
	}

	isGuildProtected := utils.Discord.IsGuildInList(
		session,
		message.GuildID,
		cfg.Protections.ServerProtections.ProtectedServers,
	)

	if isGuildProtected && cfg.Protections.ServerProtections.AntiDelete {
		userPerms, _ := session.UserChannelPermissions(message.Author.ID, message.ChannelID)

		userHasManageMessages := userPerms&discordgo.PermissionManageMessages != 0
		userIsAdmin := userPerms&discordgo.PermissionAdministrator != 0

		if userHasManageMessages || userIsAdmin {
			return
		}

		p.eventBroker.Publish(eventbroker.EventLogFeature, &featurelogs.AntiDelete{
			Message: message.Message,
		})
	}
}

func (p *Protections) RelationshipAddHandler(session *discordgo.Session, relationship *discordgo.RelationshipAdd, cfg *config.Config) {
	// Anti friend request
	func() {
		if relationship.Type != 3 {
			return
		}

		p.FriendRequests = append(p.FriendRequests, &friendRequest{
			Relationship: relationship.Relationship,
			Time:         time.Now(),
		})

		if !cfg.Protections.AntiFriendRequest {
			return
		}

		friendRequestsInLastMinute := utils.Filter(p.FriendRequests, func(fr *friendRequest) bool {
			return time.Since(fr.Time) < time.Minute
		})

		if len(friendRequestsInLastMinute) < 10 {
			return
		}

		_ = session.UserRelationshipDelete(relationship.User.ID)
		_ = api.AllowFriendRequests(session, false)

		p.eventBroker.Publish(eventbroker.EventLogFeature, &featurelogs.FriendRequestSpamBlocked{})

		go func() {
			time.Sleep(5 * time.Minute)

			_ = api.AllowFriendRequests(session, true)
		}()
	}()
}
