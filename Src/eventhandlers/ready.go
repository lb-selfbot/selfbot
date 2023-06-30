package eventhandlers

import (
	"github.com/LightningDev1/LB-Selfbot/logging"
	"github.com/LightningDev1/LB-Selfbot/sessions"
	"github.com/LightningDev1/LB-Selfbot/slashcommands"
	"github.com/LightningDev1/discordgo"
)

func (e *Events) OnReady(session *discordgo.Session, event *discordgo.Ready) {
	// Start the session handler
	sessions.RegisterSessionHandler(session)

	// Start slash command jobs
	go slashcommands.StartJobs(session)

	logging.Logging.OnReady(e.Router, event)
}
