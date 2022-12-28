package eventhandlers

import (
	"strings"

	"github.com/LightningDev1/LB-Selfbot/logging"
	"github.com/LightningDev1/discordgo"
)

func (e *Events) OnChannelCreate(_ *discordgo.Session, channel *discordgo.ChannelCreate) {
	if channel.Type != discordgo.ChannelTypeGuildText {
		return
	}

	if !strings.Contains(channel.Name, "ticket") {
		return
	}

	perms, err := e.Session.State.UserChannelPermissions(channel.ID, e.Session.State.User.ID)
	if err != nil {
		return
	}

	canSendMessages := perms&discordgo.PermissionSendMessages != 0
	canViewChannel := perms&discordgo.PermissionViewChannel != 0

	if canSendMessages && canViewChannel {
		logging.Logging.TicketNotifier(e.Session, channel.Channel)
	}
}
