package fun

import (
	"github.com/LightningDev1/LB-Selfbot/internal/services"
	"github.com/LightningDev1/dgc"
)

type FunCommands struct {
	services *services.Services
}

func NewFunCommands(services *services.Services) *FunCommands {
	return &FunCommands{services: services}
}

func (c *FunCommands) Register() {
	router := c.services.Discord.GetRouter()
	router.StartCategory("Fun", "Fun and entertainment commands")

	router.RegisterCmd(&dgc.Command{
		Name:        "8ball",
		Description: "Ask the magic 8ball a question",
		Usage:       "[p]8ball <question>",
		Handler:     c.eightball,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "coinflip",
		Description: "Flip a coin",
		Usage:       "[p]coinflip",
		Handler:     c.coinflip,
	})

	// ... Register remaining fun commands ...
}

// Keep original command handlers
func (c *FunCommands) eightball(ctx *dgc.Ctx) {
	// ... Existing implementation ...
}

func (c *FunCommands) coinflip(ctx *dgc.Ctx) {
	// ... Existing implementation ...
}
