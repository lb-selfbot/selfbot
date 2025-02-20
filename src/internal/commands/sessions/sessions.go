package sessions

import (
	"github.com/LightningDev1/LB-Selfbot/internal/services"
	"github.com/LightningDev1/dgc"
)

type SessionsCommands struct {
	services *services.Services
}

func NewSessionsCommands(services *services.Services) *SessionsCommands {
	return &SessionsCommands{services: services}
}

func (c *SessionsCommands) Register() {
	router := c.services.Discord.GetRouter()
	router.StartCategory("Sessions", "Session management commands")

	router.RegisterCmd(&dgc.Command{
		Name:        "session",
		Description: "Session management commands",
		Usage:       "[p]session <list/add/remove/clear>",
		Handler:     c.session,
	})

	// ... Register remaining session commands ...
}

func (c *SessionsCommands) session(ctx *dgc.Ctx) {
	// ... Existing implementation ...
}
