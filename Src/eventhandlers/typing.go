package eventhandlers

import (
	"github.com/LightningDev1/LB-Selfbot/config"
	"github.com/LightningDev1/LB-Selfbot/logging"
	"github.com/LightningDev1/discordgo"
)

func (e *Events) OnTypingStart(_ *discordgo.Session, typingStart *discordgo.TypingStart) {
	cfg, err := config.GetConfig()
	if err != nil {
		return
	}
	if cfg.Events.DMTypingNotifier && typingStart.GuildID == "" {
		user, err := e.Session.State.GetUser(typingStart.UserID)
		if err == nil {
			logging.Logging.DMTypingNotifier(user.String())
		}
	}
}
