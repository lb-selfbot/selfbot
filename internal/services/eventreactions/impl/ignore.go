package impl

import (
	"fmt"

	"github.com/lb-selfbot/discordgo"
	"github.com/lb-selfbot/selfbot/internal/services/eventreactions/types"
	"github.com/lb-selfbot/selfbot/internal/utils"
)

type Ignore struct {
	User  *discordgo.User
	Guild *discordgo.Guild
}

func NewIgnoreUser(user *discordgo.User) *Ignore {
	return &Ignore{User: user}
}

func NewIgnoreGuild(guild *discordgo.Guild) *Ignore {
	return &Ignore{Guild: guild}
}

func NewIgnoreFromData(session *discordgo.Session, data map[string]any) (*Ignore, error) {
	userID, ok := data["user"].(string)
	if ok {
		user, err := utils.Discord.GetUser(session, userID)
		if err != nil {
			user = &discordgo.User{ID: userID}
		}

		return NewIgnoreUser(user), nil
	}

	guildID, ok := data["guild"].(string)
	if ok {
		guild, err := utils.Discord.GetGuild(session, guildID)
		if err != nil {
			guild = &discordgo.Guild{ID: guildID}
		}

		return NewIgnoreGuild(guild), nil
	}

	utils.Logging.Error("Ignore.NewIgnoreFromData: Invalid ignore data")

	return nil, fmt.Errorf("invalid reaction data")
}

func (i *Ignore) Type() types.EventReactionType {
	if i.User != nil {
		return types.IgnoreUser
	} else if i.Guild != nil {
		return types.IgnoreGuild
	} else {
		utils.Logging.Debug("Ignore.Type: Invalid ignore type")
		return -1
	}
}

func (i *Ignore) Description() string {
	if i.User != nil {
		return "Ignore User - " + i.User.Username
	} else if i.Guild != nil {
		return "Ignore Server - " + i.Guild.Name
	} else {
		utils.Logging.Debug("Ignore.Description: Invalid ignore type")
		return "Ignore"
	}
}

func (i *Ignore) Equals(other types.EventReaction) bool {
	otherIgnore, ok := other.(*Ignore)
	if !ok {
		utils.Logging.Debug("Ignore.Equals: Other is not an Ignore")
		return false
	}

	return i.User == otherIgnore.User && i.Guild == otherIgnore.Guild
}

func (i *Ignore) Data() map[string]any {
	data := make(map[string]any)

	if i.User != nil {
		data["user"] = i.User.ID
	}
	if i.Guild != nil {
		data["guild"] = i.Guild.ID
	}

	return data
}

func (i *Ignore) Handler(session *discordgo.Session, event any) {
	message, ok := event.(*discordgo.MessageCreate)
	if !ok {
		utils.Logging.Debug("Ignore.Handler: Event is not a message create event")
		return
	}

	if message.Author.ID == session.State.User.ID {
		return
	}

	if i.Guild != nil && i.Guild.ID == message.GuildID {
		_ = session.GuildAck(i.Guild.ID)
	} else if i.User != nil && i.User.ID == message.Author.ID {
		_ = session.MessageAck(message.ChannelID, message.ID)
	}
}
