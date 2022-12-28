package eventhandlers

import (
	"github.com/LightningDev1/LB-Selfbot/logging"
	"github.com/LightningDev1/discordgo"
)

func (e *Events) OnGuildBanAdd(_ *discordgo.Session, ban *discordgo.GuildBanAdd) {
	if ban.User == nil {
		return
	}

	if ban.User.ID != e.Session.State.User.ID {
		return
	}

	logging.Logging.BanNotifier(e.Session, ban.GuildID)
}
