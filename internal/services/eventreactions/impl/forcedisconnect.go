package impl

import (
	"github.com/lb-selfbot/discordgo"
	"github.com/lb-selfbot/selfbot/internal/logging"
	"github.com/lb-selfbot/selfbot/internal/services/eventreactions/types"
)

type ForceDisconnect struct {
	Guild   *discordgo.Guild
	Channel *discordgo.Channel
	User    *discordgo.User
}

func NewForceDisconnect(g *discordgo.Guild, c *discordgo.Channel, u *discordgo.User) *ForceDisconnect {
	return &ForceDisconnect{Guild: g, Channel: c, User: u}
}

func NewForceDisconnectFromData(session *discordgo.Session, data map[string]any) (*ForceDisconnect, error) {
	var guild *discordgo.Guild
	var channel *discordgo.Channel
	var user *discordgo.User
	var err error

	guildID, ok := data["guild"].(string)
	if ok {
		guild, err = session.GuildWithFallback(guildID)
		if err != nil {
			guild = &discordgo.Guild{ID: guildID}
		}
	}

	channelID, ok := data["channel"].(string)
	if ok {
		channel, err = session.ChannelWithFallback(channelID)
		if err != nil {
			channel = &discordgo.Channel{ID: channelID}
		}
	}

	userID, ok := data["user"].(string)
	if ok {
		user, err = session.UserWithFallback(userID)
		if err != nil {
			user = &discordgo.User{ID: userID}
		}
	}

	return NewForceDisconnect(guild, channel, user), nil
}

func (r *ForceDisconnect) Type() types.EventReactionType {
	return types.ForceDisconnect
}

func (r *ForceDisconnect) Description() string {
	location := r.User.Username

	if r.Channel != nil {
		location = location + " - " + r.Channel.Name
	}
	if r.Guild != nil {
		location = location + " - " + r.Guild.Name
	}

	return "Force Disconnect (" + location + ")"
}

func (r *ForceDisconnect) Equals(other types.EventReaction) bool {
	otherForceDisconnect, ok := other.(*ForceDisconnect)
	if !ok {
		logging.Debug("ForceDisconnect.Equals: Other is not a ForceDisconnect")
		return false
	}

	return r.Guild == otherForceDisconnect.Guild &&
		r.Channel == otherForceDisconnect.Channel &&
		r.User == otherForceDisconnect.User
}

func (r *ForceDisconnect) Data() map[string]any {
	data := make(map[string]any)

	if r.Guild != nil {
		data["guild"] = r.Guild.ID
	}
	if r.Channel != nil {
		data["channel"] = r.Channel.ID
	}
	if r.User != nil {
		data["user"] = r.User.ID
	}

	return data
}

func (r *ForceDisconnect) Handler(session *discordgo.Session, event any) {
	vcUpdate, ok := event.(*discordgo.VoiceStateUpdate)
	if !ok {
		return
	}

	if vcUpdate.UserID == session.State.User.ID {
		return
	}

	if r.Guild != nil && r.Guild.ID != vcUpdate.GuildID {
		return
	} else if r.Channel != nil && r.Channel.ID != vcUpdate.ChannelID {
		return
	} else if r.User != nil && r.User.ID != vcUpdate.UserID {
		return
	}

	_ = session.GuildMemberMove(vcUpdate.GuildID, vcUpdate.UserID, nil)
}
