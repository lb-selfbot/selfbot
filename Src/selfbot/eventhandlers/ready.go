package eventhandlers

import (
	"github.com/LightningDev1/LB-Selfbot/selfbot/eventhandlers/eventreactions"
	"github.com/LightningDev1/LB-Selfbot/selfbot/logging"
	"github.com/LightningDev1/LB-Selfbot/selfbot/sessions"
	"github.com/LightningDev1/LB-Selfbot/selfbot/slashcommands"
	"github.com/LightningDev1/discordgo"
)

func (e *Events) OnReady(session *discordgo.Session, event *discordgo.Ready) {
	// Start the session handler
	sessions.RegisterSessionHandler(session)

	// Start slash command jobs
	go slashcommands.StartJobs(session)
	go eventreactions.StartLoading(session, e.Router)

	logging.Logging.OnReady(e.Router, event)
}
