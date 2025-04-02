package impl

import (
	"github.com/lb-selfbot/discordgo"
	"github.com/lb-selfbot/selfbot/internal/logging"
	"github.com/lb-selfbot/selfbot/internal/services/eventreactions/types"
)

type RemoveReactions struct {
	Guild   *discordgo.Guild
	Channel *discordgo.Channel
	User    *discordgo.User
}

func NewRemoveReactions(g *discordgo.Guild, c *discordgo.Channel, u *discordgo.User) *RemoveReactions {
	return &RemoveReactions{Guild: g, Channel: c, User: u}
}

func NewRemoveReactionsFromData(session *discordgo.Session, data map[string]any) (*RemoveReactions, error) {
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

	return NewRemoveReactions(guild, channel, user), nil
}

func (r *RemoveReactions) Type() types.EventReactionType {
	return types.RemoveReactions
}

func (r *RemoveReactions) Description() string {
	location := r.User.Username

	if r.Channel != nil {
		location = location + " - " + r.Channel.Name
	}
	if r.Guild != nil {
		location = location + " - " + r.Guild.Name
	}

	return "Remove Reactions (" + location + ")"
}

func (r *RemoveReactions) Equals(other types.EventReaction) bool {
	otherRemoveReactions, ok := other.(*RemoveReactions)
	if !ok {
		logging.Debug("RemoveReactions.Equals: Other is not a RemoveReactions")
		return false
	}
	return r.Guild == otherRemoveReactions.Guild &&
		r.Channel == otherRemoveReactions.Channel &&
		r.User == otherRemoveReactions.User
}

func (r *RemoveReactions) Data() map[string]any {
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

func (r *RemoveReactions) Handler(session *discordgo.Session, event any) {
	reactionAdd, ok := event.(*discordgo.MessageReactionAdd)
	if !ok {
		return
	}

	if reactionAdd.UserID == session.State.User.ID {
		return
	}

	if r.Guild != nil && r.Guild.ID != reactionAdd.GuildID {
		return
	} else if r.Channel != nil && r.Channel.ID != reactionAdd.ChannelID {
		return
	} else if r.User != nil && r.User.ID != reactionAdd.UserID {
		return
	}

	_ = session.MessageReactionRemove(
		reactionAdd.ChannelID,
		reactionAdd.MessageID,
		reactionAdd.Emoji.APIName(),
		reactionAdd.UserID,
	)
}
