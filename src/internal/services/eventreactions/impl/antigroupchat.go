package impl

import (
	"fmt"

	"github.com/LightningDev1/LB-Selfbot/internal/services/interfaces"
	"github.com/LightningDev1/LB-Selfbot/internal/utils"
	"github.com/LightningDev1/discordgo"
)

type AntiGroupChat struct {
	Silent bool
}

func NewAntiGroupChat(silent bool) *AntiGroupChat {
	return &AntiGroupChat{Silent: silent}
}

func NewAntiGroupChatFromData(_ *discordgo.Session, data map[string]any) (*AntiGroupChat, error) {
	// error doesn't matter here because it will default to false
	silent, _ := data["Silent"].(bool)

	return NewAntiGroupChat(silent), nil
}

func (a *AntiGroupChat) Type() interfaces.EventReactionType {
	return interfaces.EventReactionAntiGroupChat
}

func (a *AntiGroupChat) Description() string {
	return fmt.Sprintf("Anti Group Chat (silent: %s)", utils.Misc.FormatBoolYesNo(a.Silent))
}

func (a *AntiGroupChat) Equals(other interfaces.EventReaction) bool {
	// We shouldn't check silent because the value does not matter
	// There should only be one AntiGroupChat reaction at once
	_, ok := other.(*AntiGroupChat)
	return ok
}

func (a *AntiGroupChat) Data() map[string]any {
	return map[string]any{
		"Silent": a.Silent,
	}
}

func (a *AntiGroupChat) Handler(session *discordgo.Session, event any) {
	channel, ok := event.(*discordgo.ChannelCreate)
	if !ok {
		utils.Logging.Debug("AntiGroupChat.Handler: Event is not a channel create event")
		return
	}

	if channel.Type != discordgo.ChannelTypeGroupDM {
		return
	}

	_, _ = session.ChannelDelete(channel.ID, a.Silent)
}
