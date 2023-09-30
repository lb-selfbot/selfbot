package eventhandlers

import (
	"github.com/LightningDev1/LB-Selfbot/selfbot/config"
	"github.com/LightningDev1/LB-Selfbot/selfbot/logging"
	"github.com/LightningDev1/LB-Selfbot/selfbot/utils"
	"github.com/LightningDev1/discordgo"
)

func (e *Events) OnTypingStart(_ *discordgo.Session, event *discordgo.TypingStart) {
	cfg, err := config.GetConfig()
	if err != nil {
		return
	}

	if !cfg.Events.DMTypingNotifier || event.GuildID != "" {
		return
	}

	userName := utils.Discord.GetUserName(e.Session, event.UserID)
	if userName == "Unknown" {
		return
	}

	logging.Logging.DMTypingNotifier(userName)
}
