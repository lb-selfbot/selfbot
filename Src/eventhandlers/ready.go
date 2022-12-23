package eventhandlers

import (
	"github.com/LightningDev1/LB-Selfbot/logging"
	"github.com/LightningDev1/discordgo"
)

func (e *Events) OnReady(_ *discordgo.Session, event *discordgo.Ready) {
	logging.Logging.OnReady(e.Router, event)
}
