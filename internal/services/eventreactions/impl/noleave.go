package impl

import (
	"github.com/lb-selfbot/discordgo"
	"github.com/lb-selfbot/selfbot/internal/logging"
	"github.com/lb-selfbot/selfbot/internal/services/eventreactions/types"
	"github.com/lb-selfbot/selfbot/internal/utils"
)

type NoLeave struct {
	Channel *discordgo.Channel
	User    *discordgo.User
}

func NewNoLeave(c *discordgo.Channel, u *discordgo.User) *NoLeave {
	return &NoLeave{Channel: c, User: u}
}

func NewNoLeaveFromData(session *discordgo.Session, data map[string]any) (*NoLeave, error) {
	var channel *discordgo.Channel
	var user *discordgo.User
	var err error

	channelID, ok := data["channel"].(string)
	if ok {
		channel, err = utils.Discord.GetChannel(session, channelID)
		if err != nil {
			channel = &discordgo.Channel{ID: channelID}
		}
	}

	userID, ok := data["user"].(string)
	if ok {
		user, err = utils.Discord.GetUser(session, userID)
		if err != nil {
			user = &discordgo.User{ID: userID}
		}
	}

	return NewNoLeave(channel, user), nil
}

func (n *NoLeave) Type() types.EventReactionType {
	return types.NoLeave
}

func (n *NoLeave) Description() string {
	location := n.User.Username

	if n.Channel != nil {
		location = location + " - " + n.Channel.Name
	}

	return "No Leave (" + location + ")"
}

func (n *NoLeave) Equals(other types.EventReaction) bool {
	otherNoLeave, ok := other.(*NoLeave)
	if !ok {
		logging.Debug("NoLeave.Equals: Other is not a NoLeave")
		return false
	}

	return n.Channel == otherNoLeave.Channel &&
		n.User == otherNoLeave.User
}

func (n *NoLeave) Data() map[string]any {
	data := make(map[string]any)

	if n.Channel != nil {
		data["channel"] = n.Channel.ID
	}
	if n.User != nil {
		data["user"] = n.User.ID
	}

	return data
}

func (r *NoLeave) Handler(session *discordgo.Session, event any) {
	recipientRemove, ok := event.(*discordgo.ChannelRecipientRemove)
	if !ok {
		logging.Debug("NoLeave.Handler: Event is not a channel recipient remove event")
		return
	}

	if recipientRemove.User == nil {
		return
	} else if recipientRemove.User.ID == session.State.User.ID {
		return
	}

	if r.Channel != nil && r.Channel.ID != recipientRemove.ChannelID {
		return
	} else if r.User != nil && r.User.ID != recipientRemove.User.ID {
		return
	}

	_ = session.ChannelRecipientAdd(recipientRemove.ChannelID, recipientRemove.User.ID)
}
