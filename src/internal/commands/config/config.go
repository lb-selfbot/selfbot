package config

import (
	"fmt"

	"github.com/LightningDev1/LB-Selfbot/internal/config"
	"github.com/LightningDev1/LB-Selfbot/internal/services"
	"github.com/LightningDev1/dgc"
)

type ConfigCommands struct {
	services *services.Services
}

func NewConfigCommands(services *services.Services) *ConfigCommands {
	return &ConfigCommands{services: services}
}

func (c *ConfigCommands) Register() {
	router := c.services.Discord.GetRouter()
	router.StartCategory("Config", "Configuration commands")

	router.RegisterCmd(&dgc.Command{
		Name:        "prefix",
		Description: "Change the command prefix",
		Usage:       "[p]prefix <new prefix>",
		Handler:     c.prefix,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "deletecommands",
		Description: "Toggle command deletion",
		Usage:       "[p]deletecommands <on/off>",
		Handler:     c.deletecommands,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "embedcolor",
		Description: "Change the embed color",
		Usage:       "[p]embedcolor <hex color>",
		Handler:     c.embedcolor,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "saveconfig",
		Description: "Save the current config",
		Usage:       "[p]saveconfig",
		Handler:     c.saveconfig,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "reloadconfig",
		Description: "Reload the config from disk",
		Usage:       "[p]reloadconfig",
		Handler:     c.reloadconfig,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "resetconfig",
		Description: "Reset the config to defaults",
		Usage:       "[p]resetconfig",
		Handler:     c.resetconfig,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "backupconfig",
		Description: "Backup the current config",
		Usage:       "[p]backupconfig",
		Handler:     c.backupconfig,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "showconfig",
		Description: "Show the current config",
		Usage:       "[p]showconfig",
		Handler:     c.showconfig,
	})
}

func (c *ConfigCommands) prefix(ctx *dgc.Ctx) {
	if len(ctx.Arguments) < 1 {
		ctx.RespondText("Please provide a new prefix")
		return
	}

	newPrefix := ctx.Arguments[0]
	cfg, err := config.GetConfig()
	if err != nil {
		ctx.RespondText("Error getting config: " + err.Error())
		return
	}

	cfg.Prefix = newPrefix
	if err := cfg.Save(); err != nil {
		ctx.RespondText("Error saving config: " + err.Error())
		return
	}

	ctx.RespondEmbed(fmt.Sprintf("Successfully changed prefix to `%s`", newPrefix))
}

// Add remaining command handlers...
