package utility

import (
	"github.com/LightningDev1/LB-Selfbot/internal/services"
	"github.com/LightningDev1/dgc"
)

type UtilityCommands struct {
	services *services.Services
}

func NewUtilityCommands(services *services.Services) *UtilityCommands {
	return &UtilityCommands{services: services}
}

func (c *UtilityCommands) Register() {
	router := c.services.Discord.GetRouter()
	router.StartCategory("Utility", "Utility and helper commands")

	router.RegisterCmd(&dgc.Command{
		Name:        "ping",
		Description: "Check bot latency",
		Usage:       "[p]ping",
		Handler:     c.ping,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "uptime",
		Description: "Show bot uptime",
		Usage:       "[p]uptime",
		Handler:     c.uptime,
	})

	// ... Register remaining utility commands ...
}

func (c *UtilityCommands) ping(ctx *dgc.Ctx) {
	// ... Existing implementation ...
}
