package trolling

import (
	"github.com/LightningDev1/LB-Selfbot/internal/services"
	"github.com/LightningDev1/dgc"
)

type TrollingCommands struct {
	services *services.Services
}

func NewTrollingCommands(services *services.Services) *TrollingCommands {
	return &TrollingCommands{services: services}
}

func (c *TrollingCommands) Register() {
	router := c.services.Discord.GetRouter()
	router.StartCategory("Trolling", "Trolling and joke commands")

	router.RegisterCmd(&dgc.Command{
		Name:        "spam",
		Description: "Spam a message",
		Usage:       "[p]spam <amount> <message>",
		Flags:       []string{"dangerous"},
		Handler:     c.spam,
	})

	// ... Register remaining trolling commands ...
}

// Keep original command handlers
func (c *TrollingCommands) spam(ctx *dgc.Ctx) {
	// ... Existing implementation ...
}
