package ai

import (
	"github.com/LightningDev1/LB-Selfbot/internal/services"
	"github.com/LightningDev1/dgc"
)

type AICommands struct {
	services *services.Services
}

func NewAICommands(services *services.Services) *AICommands {
	return &AICommands{services: services}
}

func (c *AICommands) Register() {
	router := c.services.Discord.GetRouter()
	router.StartCategory("AI", "AI commands")

	router.RegisterCmd(&dgc.Command{
		Name:        "chat",
		Description: "Chat with an AI",
		Usage:       "[p]chat <message>",
		Handler:     c.chat,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "imagine",
		Description: "Generate an image with AI",
		Usage:       "[p]imagine <prompt>",
		Handler:     c.imagine,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "aimodel",
		Description: "Set or view the current AI model",
		Usage:       "[p]aimodel <get/set> [model]",
		Handler:     c.aimodel,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "aipersonality",
		Description: "Set or view AI personality",
		Usage:       "[p]aipersonality <get/set> [personality]",
		Handler:     c.aipersonality,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "aicontext",
		Description: "Manage AI conversation context",
		Usage:       "[p]aicontext <clear/show>",
		Handler:     c.aicontext,
	})
}

func (c *AICommands) chat(ctx *dgc.Ctx) {
	// ... Existing implementation ...
}

// ... Keep all other existing command implementations ...
