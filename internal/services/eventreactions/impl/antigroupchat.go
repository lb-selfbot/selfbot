package impl

import (
	"fmt"

	"github.com/lb-selfbot/discordgo"
	"github.com/lb-selfbot/selfbot/internal/formatting"
	"github.com/lb-selfbot/selfbot/internal/services/eventreactions/types"
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

func (a *AntiGroupChat) Type() types.EventReactionType {
	return types.AntiGroupChat
}

func (a *AntiGroupChat) Description() string {
	return fmt.Sprintf("Anti Group Chat (silent: %s)", formatting.BoolYesNo(a.Silent))
}

func (a *AntiGroupChat) Equals(other types.EventReaction) bool {
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
		return
	}

	if channel.Type != discordgo.ChannelTypeGroupDM {
		return
	}

	_, _ = session.ChannelDelete(channel.ID, a.Silent)
}
