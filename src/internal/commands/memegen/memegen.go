package memegen

import (
	"github.com/LightningDev1/LB-Selfbot/internal/services"
	"github.com/LightningDev1/dgc"
)

type MemeGenCommands struct {
	services *services.Services
}

func NewMemeGenCommands(services *services.Services) *MemeGenCommands {
	return &MemeGenCommands{services: services}
}

func (c *MemeGenCommands) Register() {
	router := c.services.Discord.GetRouter()
	router.StartCategory("MemeGen", "Meme generation commands")

	router.RegisterCmd(&dgc.Command{
		Name:        "meme",
		Description: "Generate a meme from a template",
		Usage:       "[p]meme <template> <text>",
		Handler:     c.meme,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "memetemplates",
		Description: "List available meme templates",
		Usage:       "[p]memetemplates",
		Handler:     c.memetemplates,
	})

	// ... Register remaining memegen commands ...
}

func (c *MemeGenCommands) meme(ctx *dgc.Ctx) {
	// ... Existing implementation ...
}
