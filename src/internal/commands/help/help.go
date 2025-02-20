package help

import (
	"fmt"
	"strings"

	"github.com/LightningDev1/LB-Selfbot/internal/services"
	"github.com/LightningDev1/dgc"
)

type HelpCommands struct {
	services *services.Services
}

func NewHelpCommands(services *services.Services) *HelpCommands {
	return &HelpCommands{services: services}
}

func (c *HelpCommands) Register() {
	router := c.services.Discord.GetRouter()
	router.StartCategory("Help", "Help and documentation commands")

	router.RegisterCmd(&dgc.Command{
		Name:        "help",
		Description: "Show help for commands",
		Usage:       "[p]help [command/category]",
		Handler:     c.help,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "commands",
		Description: "List all commands",
		Usage:       "[p]commands [category]",
		Handler:     c.commands,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "categories",
		Description: "List all command categories",
		Usage:       "[p]categories",
		Handler:     c.categories,
	})
}

func (c *HelpCommands) help(ctx *dgc.Ctx) {
	if len(ctx.Arguments) == 0 {
		// Show general help
		c.showGeneralHelp(ctx)
		return
	}

	// Show specific command/category help
	query := strings.ToLower(ctx.Arguments[0])
	if cmd := c.services.Commands.GetCommand(query); cmd != nil {
		c.showCommandHelp(ctx, cmd)
		return
	}

	if cat := c.services.Commands.GetCategory(query); cat != nil {
		c.showCategoryHelp(ctx, cat)
		return
	}

	ctx.RespondText(fmt.Sprintf("No command or category found for `%s`", query))
}

func (c *HelpCommands) commands(ctx *dgc.Ctx) {
	// ... Existing implementation ...
}

func (c *HelpCommands) categories(ctx *dgc.Ctx) {
	// ... Existing implementation ...
}

// Add remaining command handlers and helper methods...
