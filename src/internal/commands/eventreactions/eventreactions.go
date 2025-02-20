package eventreactions

import (
	"github.com/LightningDev1/LB-Selfbot/internal/services"
	"github.com/LightningDev1/dgc"
)

type EventReactionsCommands struct {
	services *services.Services
}

func NewEventReactionsCommands(services *services.Services) *EventReactionsCommands {
	return &EventReactionsCommands{services: services}
}

func (c *EventReactionsCommands) Register() {
	router := c.services.Discord.GetRouter()
	router.StartCategory("EventReactions", "Event reaction commands")

	router.RegisterCmd(&dgc.Command{
		Name:        "autoreply",
		Description: "Configure auto-reply settings",
		Usage:       "[p]autoreply <add/remove/list/clear>",
		Handler:     c.autoreply,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "autoreplyglobal",
		Description: "Configure global auto-reply settings",
		Usage:       "[p]autoreplyglobal <get/set/toggle>",
		Handler:     c.autoreplyglobal,
	})

	// ... Register remaining event reaction commands ...
}

func (c *EventReactionsCommands) autoreply(ctx *dgc.Ctx) {
	// ... Existing implementation ...
}
