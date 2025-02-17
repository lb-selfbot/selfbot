package eventreactions

import (
	"strconv"
	"time"
	"unicode"

	"github.com/LightningDev1/LB-Selfbot/internal/constants"
	"github.com/LightningDev1/LB-Selfbot/selfbot/api"
	"github.com/LightningDev1/LB-Selfbot/selfbot/config"
	"github.com/LightningDev1/LB-Selfbot/selfbot/embed"
	"github.com/LightningDev1/LB-Selfbot/selfbot/logging"

	"github.com/LightningDev1/LB-Selfbot/selfbot/utils"
	"github.com/LightningDev1/dgc"
	"github.com/LightningDev1/discordgo"
)

type friendRequest struct {
	*discordgo.Relationship

	Time time.Time
}

type Protections struct {
	NewDMs         []*discordgo.Message
	FriendRequests []*friendRequest
	LastBlock      time.Time
}

func NewProtections() *Protections {
	return &Protections{}
}

func (p *Protections) Type() EventReactionType {
	return PROTECTIONS
}

func (p *Protections) Description() string {
	return "Protections (default)"
}

func (p *Protections) Equals(other EventReaction) bool {
	_, ok := other.(*Protections)

	return ok
}

func (p *Protections) Data() map[string]any {
	return map[string]any{}
}

func (p *Protections) Handler(session *discordgo.Session, event any) {
	cfg, err := config.GetConfig()
	if err != nil {
		utils.Logging.Debug("Protections.Handler: Error getting config: " + err.Error())
		return
	}

	switch event := event.(type) {
	case *discordgo.MessageCreate:
		p.MessageCreateHandler(session, event, &cfg)
	case *discordgo.MessageDelete:
		p.MessageDeleteHandler(session, event, &cfg)
	case *discordgo.RelationshipAdd:
		p.RelationshipAddHandler(session, event, &cfg)
	default:
		utils.Logging.Debug("Protections.Handler: Event is not a protections event")
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

		logging.Logging.MassDMBlocked()

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

		ctx := &dgc.Ctx{
			Session: session,
			Event:   &discordgo.MessageCreate{Message: message.Message},
		}

		if cfg.Protections.ServerProtections.AntiInvite && constants.INVITE_REGEX.MatchString(message.Content) {
			err := session.ChannelMessageDelete(message.ChannelID, message.ID)
			if err != nil {
				return
			}

			embed.NewEmbed().
				SetDescription("Invites are not allowed in this server.").
				Send(ctx)
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
				err := session.ChannelMessageDelete(message.ChannelID, message.ID)
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

		logging.Logging.AntiDelete(session, message.Message)
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

		logging.Logging.FriendRequestSpamBlocked()

		go func() {
			time.Sleep(5 * time.Minute)

			_ = api.AllowFriendRequests(session, true)
		}()
	}()
}
