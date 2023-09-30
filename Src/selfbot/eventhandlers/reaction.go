package eventhandlers

import (
	"github.com/LightningDev1/discordgo"
)

func (e *Events) OnMessageReactionAdd(_ *discordgo.Session, reactionAdd *discordgo.MessageReactionAdd) {
	if reactionAdd.Member == nil {
		return
	}

	for _, event := range e.Settings.RemoveReactions {
		if event.GuildID != reactionAdd.Member.GuildID {
			continue
		}

		if event.UserID != reactionAdd.Member.User.ID {
			continue
		}

		_ = e.Session.MessageReactionRemove(
			reactionAdd.ChannelID,
			reactionAdd.MessageID,
			reactionAdd.Emoji.APIName(),
			reactionAdd.Member.User.ID,
		)
	}
}
