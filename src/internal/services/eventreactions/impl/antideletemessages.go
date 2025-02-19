package impl

import (
	"fmt"
	"strings"

	"github.com/LightningDev1/LB-Selfbot/internal/services/interfaces"
	"github.com/LightningDev1/LB-Selfbot/selfbot/config"
	"github.com/LightningDev1/LB-Selfbot/selfbot/utils"
	"github.com/LightningDev1/discordgo"
)

type AntiDeleteMessages struct {
	Channel *discordgo.Channel
}

func NewAntiDeleteMessages(channel *discordgo.Channel) *AntiDeleteMessages {
	return &AntiDeleteMessages{Channel: channel}
}

func NewAntiDeleteMessagesFromData(session *discordgo.Session, data map[string]any) (*AntiDeleteMessages, error) {
	channelID, ok := data["channel"].(string)
	if ok {
		channel, err := utils.Discord.GetChannel(session, channelID)
		if err != nil {
			channel = &discordgo.Channel{ID: channelID}
		}

		return NewAntiDeleteMessages(channel), nil
	} else {
		utils.Logging.Error("AntiDeleteMessages.NewAntiDeleteMessagesFromData: Invalid antideletemessages data")
		return nil, fmt.Errorf("invalid antideletemessages data")
	}
}

func (a *AntiDeleteMessages) Type() interfaces.EventReactionType {
	return interfaces.EventReactionAntiDeleteMessages
}

func (a *AntiDeleteMessages) Description() string {
	if a.Channel == nil {
		return "Anti Delete Messages"
	}

	return fmt.Sprintf("Anti Delete Messages - %s", a.Channel.Name)
}

func (a *AntiDeleteMessages) Equals(other interfaces.EventReaction) bool {
	otherAntiDeleteMessages, ok := other.(*AntiDeleteMessages)
	if !ok {
		return false
	}

	return a.Channel != nil && otherAntiDeleteMessages.Channel != nil &&
		a.Channel.ID == otherAntiDeleteMessages.Channel.ID
}

func (a *AntiDeleteMessages) Data() map[string]any {
	return map[string]any{"channel": a.Channel.ID}
}

func (a *AntiDeleteMessages) Handler(session *discordgo.Session, event any) {
	message, ok := event.(*discordgo.MessageDelete)
	if !ok {
		utils.Logging.Debug("AntiDeleteMessages.Handler: Event is not a message delete event")
		return
	}

	if message.ChannelID != a.Channel.ID || message.Author.ID != session.State.User.ID {
		return
	}

	cfg, err := config.GetConfig()
	if err != nil {
		utils.Logging.Debug("AntiDeleteMessages.Handler: Error getting config: " + err.Error())
		return
	}

	// Don't resend messages with the command prefix.
	if strings.HasPrefix(message.Content, cfg.Selfbot.CommandPrefix) {
		return
	}

	// Message is not handled because we want messages sent by the handler to also be resent
	// when they are deleted.
	_, err = session.ChannelMessageSendComplex(message.ChannelID, &discordgo.MessageSend{
		Content:  message.Content,
		Activity: message.Activity,
		TTS:      message.TTS,
	})
	if err != nil {
		utils.Logging.Debug("AntiDeleteMessages.Handler: Error sending message: " + err.Error())
		return
	}
}
