package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/LightningDev1/LB-Selfbot/internal/config"
	"github.com/LightningDev1/LB-Selfbot/internal/router"
	"github.com/LightningDev1/LB-Selfbot/internal/services"
	"github.com/LightningDev1/LB-Selfbot/internal/utils2"

	"github.com/LightningDev1/LB-Selfbot/internal/utils"
)

type ConfigCommands struct{ services *services.Services }

func NewConfigCommands(s *services.Services) *ConfigCommands {
	return &ConfigCommands{services: s}
}

func (c *ConfigCommands) Register() {
	cmdRouter := c.services.Discord.GetRouter()

	cmdRouter.RegisterCategory(&router.Category{
		Name:        "Config",
		Description: "Change the bot's settings",
	}, []*router.Command{
		c.createConfigActionCommand(configCommandData{
			Name:         "commandprefix",
			ReadableName: "Command Prefix",
			ValueType:    configCommandValueTypeString,
			GetDesc:      "Get the current command prefix",
			GetFunc:      func(cfg *config.Config) any { return cfg.Selfbot.CommandPrefix },
			SetDesc:      "Set the command prefix",
			SetFunc:      func(cfg *config.Config, value any) { cfg.Selfbot.CommandPrefix = value.(string) },
		}),
		c.createConfigActionCommand(configCommandData{
			Name:         "commandmode",
			ReadableName: "Command Mode",
			ValueType:    configCommandValueTypeString,
			GetDesc:      "Get the current command mode",
			GetFunc:      func(cfg *config.Config) any { return cfg.Selfbot.CommandMode },
			SetDesc:      "Set the command mode",
			SetFunc:      func(cfg *config.Config, value any) { cfg.Selfbot.CommandMode = value.(string) },
			ValuesDesc:   "List available command modes",
			ValuesFunc:   func() []any { return []any{"Codeblock", "Text", "Indent", "WebhookEmbed", "URLEmbed"} },
		}),
		c.createConfigActionCommand(configCommandData{
			Name:         "title",
			ReadableName: "Embed Title",
			ValueType:    configCommandValueTypeString,
			GetDesc:      "Get the current embed title",
			GetFunc:      func(cfg *config.Config) any { return cfg.Selfbot.Title },
			SetDesc:      "Set the embed title",
			SetFunc:      func(cfg *config.Config, value any) { cfg.Selfbot.Title = value.(string) },
		}),

		c.createConfigActionCommand(configCommandData{
			Name:         "footer",
			ReadableName: "Embed Footer",
			ValueType:    configCommandValueTypeString,
			GetDesc:      "Get the current embed footer",
			GetFunc:      func(cfg *config.Config) any { return cfg.Selfbot.Footer },
			SetDesc:      "Set the embed footer",
			SetFunc:      func(cfg *config.Config, value any) { cfg.Selfbot.Footer = value.(string) },
		}),
		c.createConfigActionCommand(configCommandData{
			Name:         "delresponses",
			ReadableName: "Delete Responses",
			ValueType:    configCommandValueTypeBool,
			GetDesc:      "Get delete responses status",
			GetFunc:      func(cfg *config.Config) any { return cfg.Selfbot.DeleteResponses },
			SetDesc:      "Set delete responses status",
			SetFunc:      func(cfg *config.Config, value any) { cfg.Selfbot.DeleteResponses = value.(bool) },
			ToggleDesc:   "Toggle delete responses",
			ToggleFunc:   func(cfg *config.Config) { cfg.Selfbot.DeleteResponses = !cfg.Selfbot.DeleteResponses },
		}),
		c.createConfigActionCommand(configCommandData{
			Name:         "delresponsesdelay",
			ReadableName: "Delete Responses Delay",
			ValueType:    configCommandValueTypeDuration,
			GetDesc:      "Get delete responses delay",
			GetFunc:      func(cfg *config.Config) any { return cfg.Selfbot.DeleteResponsesDelay },
			SetDesc:      "Set delete responses delay",
			SetFunc:      func(cfg *config.Config, value any) { cfg.Selfbot.DeleteResponsesDelay = value.(int) },
		}),
		c.createConfigActionCommand(configCommandData{
			Name:         "delcommands",
			ReadableName: "Delete Commands",
			ValueType:    configCommandValueTypeBool,
			GetDesc:      "Get delete commands status",
			GetFunc:      func(cfg *config.Config) any { return cfg.Selfbot.DeleteCommands },
			SetDesc:      "Set delete commands status",
			SetFunc:      func(cfg *config.Config, value any) { cfg.Selfbot.DeleteCommands = value.(bool) },
			ToggleDesc:   "Toggle delete commands",
			ToggleFunc:   func(cfg *config.Config) { cfg.Selfbot.DeleteCommands = !cfg.Selfbot.DeleteCommands },
		}),
		c.createConfigActionCommand(configCommandData{
			Name:         "clientspoof",
			ReadableName: "Client Spoof",
			ValueType:    configCommandValueTypeString,
			GetDesc:      "Get current client spoof setting",
			GetFunc:      func(cfg *config.Config) any { return cfg.Selfbot.Client },
			SetDesc:      "Set client spoof setting",
			SetFunc:      func(cfg *config.Config, value any) { cfg.Selfbot.Client = utils.String.Capitalize(value.(string)) },
			ValuesDesc:   "List available clients",
			ValuesFunc:   func() []any { return []any{"Mobile", "Web", "Desktop", "Console"} },
		}),
		c.createConfigActionCommand(configCommandData{
			Name:         "timezone",
			ReadableName: "Timezone",
			ValueType:    configCommandValueTypeString,
			GetDesc:      "Get current timezone",
			GetFunc:      func(cfg *config.Config) any { return cfg.Selfbot.Timezone },
			SetDesc:      "Set timezone",
			SetFunc: func(cfg *config.Config, value any) {
				timezone := value.(string)
				if _, err := time.LoadLocation(timezone); err == nil {
					cfg.Selfbot.Timezone = timezone
				}
			},
		}),
		c.createConfigActionCommand(configCommandData{
			Name:         "embedcolor",
			ReadableName: "Embed Color",
			ValueType:    configCommandValueTypeString,
			GetDesc:      "Get current embed color",
			GetFunc:      func(cfg *config.Config) any { return cfg.Selfbot.EmbedColor },
			SetDesc:      "Set embed color",
			SetFunc:      func(cfg *config.Config, value any) { cfg.Selfbot.EmbedColor = value.(string) },
		}),
		c.createConfigActionCommand(configCommandData{
			Name:         "embedimage",
			ReadableName: "Embed Image",
			ValueType:    configCommandValueTypeString,
			GetDesc:      "Get current embed image",
			GetFunc:      func(cfg *config.Config) any { return cfg.Selfbot.EmbedImage },
			SetDesc:      "Set embed image",
			SetFunc:      func(cfg *config.Config, value any) { cfg.Selfbot.EmbedImage = value.(string) },
		}),
		c.createConfigActionCommand(configCommandData{
			Name:         "embedurl",
			ReadableName: "Embed URL",
			ValueType:    configCommandValueTypeString,
			GetDesc:      "Get current embed URL",
			GetFunc:      func(cfg *config.Config) any { return cfg.Selfbot.EmbedURL },
			SetDesc:      "Set embed URL",
			SetFunc:      func(cfg *config.Config, value any) { cfg.Selfbot.EmbedURL = value.(string) },
		}),
		// c.createConfigActionCommand(configCommandData{
		// 	Name:         "forcesubscribedservers",
		// 	ReadableName: "Force Subscribed Servers",
		// 	ValueType:    configCommandValueTypeBool,
		// 	GetDesc:      "Get force subscribed servers status",
		// 	GetFunc:      func(cfg *config.Config) any { return len(cfg.Selfbot.ForceSubscribedServers) > 0 },
		// 	SetDesc:      "Set force subscribed servers status",
		// 	SetFunc: func(cfg *config.Config, value any) {
		// 		if value.(bool) {
		// 			guildIDs := make([]string, 0)
		// 			for _, guild := range c.services.Discord.GetState().Guilds {
		// 				guildIDs = append(guildIDs, guild.ID)
		// 			}
		// 			cfg.Selfbot.ForceSubscribedServers = guildIDs
		// 		} else {
		// 			cfg.Selfbot.ForceSubscribedServers = []string{}
		// 		}
		// 	},
		// 	ToggleDesc: "Toggle force subscribed servers",
		// 	ToggleFunc: func(cfg *config.Config) {
		// 		if len(cfg.Selfbot.ForceSubscribedServers) > 0 {
		// 			cfg.Selfbot.ForceSubscribedServers = []string{}
		// 		} else {
		// 			guildIDs := make([]string, 0)
		// 			for _, guild := range c.services.Discord.GetState().Guilds {
		// 				guildIDs = append(guildIDs, guild.ID)
		// 			}
		// 			cfg.Selfbot.ForceSubscribedServers = guildIDs
		// 		}
		// 	},
		// }),
		// {
		// 	Name:        "aliases",
		// 	Description: "Command aliases commands",
		// 	Flags:       []router.CommandFlag{router.CmdFlagDangerous},
		// 	SubCommands: []*router.Command{
		// 		{
		// 			Name:        "get",
		// 			Description: "Get a command alias",
		// 			Arguments:   "<alias>",
		// 			Handler:     c.getAlias,
		// 		},
		// 		{
		// 			Name:        "set",
		// 			Description: "Set a command alias",
		// 			Arguments:   "<alias> <command>",
		// 			Handler:     c.setAlias,
		// 		},
		// 		{
		// 			Name:        "rm",
		// 			Description: "Remove a command alias",
		// 			Arguments:   "<alias>",
		// 			Handler:     c.removeAlias,
		// 		},
		// 		{
		// 			Name:        "list",
		// 			Description: "List all command aliases",
		// 			Handler:     c.listAliases,
		// 		},
		// 		{
		// 			Name:        "clear",
		// 			Description: "Clear all command aliases",
		// 			Handler:     c.clearAliases,
		// 		},
		// 	},
		// },
		c.createConfigActionCommand(configCommandData{
			Name:         "notifications",
			ReadableName: "Notifications",
			ValueType:    configCommandValueTypeBool,
			GetDesc:      "Get notifications status",
			GetFunc:      func(cfg *config.Config) any { return cfg.Notifications.Enabled },
			SetDesc:      "Set notifications status",
			SetFunc:      func(cfg *config.Config, value any) { cfg.Notifications.Enabled = value.(bool) },
			ToggleDesc:   "Toggle notifications",
			ToggleFunc:   func(cfg *config.Config) { cfg.Notifications.Enabled = !cfg.Notifications.Enabled },
			ExampleDesc:  "Show an example of a notification",
			ExampleFunc: func() {
				c.services.UI.ShowNotification("This is an example notification.")
			},
		}),
		c.createConfigActionCommand(configCommandData{
			Name:         "nitrosniper",
			ReadableName: "Nitro Sniper",
			ValueType:    configCommandValueTypeBool,
			GetDesc:      "Get nitro sniper status",
			GetFunc:      func(cfg *config.Config) any { return cfg.Snipers.NitroSniper },
			SetDesc:      "Set nitro sniper status",
			SetFunc:      func(cfg *config.Config, value any) { cfg.Snipers.NitroSniper = value.(bool) },
			ToggleDesc:   "Toggle nitro sniper",
			ToggleFunc:   func(cfg *config.Config) { cfg.Snipers.NitroSniper = !cfg.Snipers.NitroSniper },
		}),
		c.createConfigActionCommand(configCommandData{
			Name:         "giveawaysniper",
			ReadableName: "Giveaway Sniper",
			ValueType:    configCommandValueTypeBool,
			GetDesc:      "Get giveaway sniper status",
			GetFunc:      func(cfg *config.Config) any { return cfg.Snipers.GiveawaySniper.Enabled },
			SetDesc:      "Set giveaway sniper status",
			SetFunc:      func(cfg *config.Config, value any) { cfg.Snipers.GiveawaySniper.Enabled = value.(bool) },
			ToggleDesc:   "Toggle giveaway sniper",
			ToggleFunc:   func(cfg *config.Config) { cfg.Snipers.GiveawaySniper.Enabled = !cfg.Snipers.GiveawaySniper.Enabled },
		}),
		c.createConfigActionCommand(configCommandData{
			Name:         "giveawaysnipersleep",
			ReadableName: "Giveaway Sniper Sleep",
			ValueType:    configCommandValueTypeDuration,
			GetDesc:      "Get giveaway sniper sleep duration",
			GetFunc:      func(cfg *config.Config) any { return cfg.Snipers.GiveawaySniper.Sleep },
			SetDesc:      "Set giveaway sniper sleep duration",
			SetFunc:      func(cfg *config.Config, value any) { cfg.Snipers.GiveawaySniper.Sleep = value.(int) },
		}),
		// {
		// 	Name:        "customgiveawaybots",
		// 	Description: "Custom giveaway bots commands",
		// 	Flags:       []router.CommandFlag{router.CmdFlagDangerous},
		// 	SubCommands: []*router.Command{
		// 		{
		// 			Name:        "add",
		// 			Description: "Add a custom giveaway bot",
		// 			Arguments:   "<bot id>",
		// 			Handler:     c.addCustomGiveawayBot,
		// 		},
		// 		{
		// 			Name:        "rm",
		// 			Description: "Remove a custom giveaway bot",
		// 			Arguments:   "<bot id>",
		// 			Handler:     c.removeCustomGiveawayBot,
		// 		},
		// 		{
		// 			Name:        "list",
		// 			Description: "List custom giveaway bots",
		// 			Handler:     c.listCustomGiveawayBots,
		// 		},
		// 		{
		// 			Name:        "clear",
		// 			Description: "Clear all custom giveaway bots",
		// 			Handler:     c.clearCustomGiveawayBots,
		// 		},
		// 	},
		// },
		// {
		// 	Name:        "sniperblacklist",
		// 	Description: "Sniper blacklist commands",
		// 	Flags:       []router.CommandFlag{router.CmdFlagDangerous},
		// 	SubCommands: []*router.Command{
		// 		{
		// 			Name:        "add",
		// 			Description: "Add a user to sniper blacklist",
		// 			Arguments:   "<user>",
		// 			Handler:     c.addSniperBlacklist,
		// 		},
		// 		{
		// 			Name:        "rm",
		// 			Description: "Remove a user from sniper blacklist",
		// 			Arguments:   "<user>",
		// 			Handler:     c.removeSniperBlacklist,
		// 		},
		// 		{
		// 			Name:        "list",
		// 			Description: "List sniper blacklist",
		// 			Handler:     c.listSniperBlacklist,
		// 		},
		// 		{
		// 			Name:        "clear",
		// 			Description: "Clear sniper blacklist",
		// 			Handler:     c.clearSniperBlacklist,
		// 		},
		// 	},
		// },
		// {
		// 	Name:        "sniperserverblacklist",
		// 	Description: "Sniper server blacklist commands",
		// 	Flags:       []router.CommandFlag{router.CmdFlagDangerous},
		// 	SubCommands: []*router.Command{
		// 		{
		// 			Name:        "add",
		// 			Description: "Add a server to sniper blacklist",
		// 			Arguments:   "<server>",
		// 			Handler:     c.addSniperServerBlacklist,
		// 		},
		// 		{
		// 			Name:        "rm",
		// 			Description: "Remove a server from sniper blacklist",
		// 			Arguments:   "<server>",
		// 			Handler:     c.removeSniperServerBlacklist,
		// 		},
		// 		{
		// 			Name:        "list",
		// 			Description: "List sniper server blacklist",
		// 			Handler:     c.listSniperServerBlacklist,
		// 		},
		// 		{
		// 			Name:        "clear",
		// 			Description: "Clear sniper server blacklist",
		// 			Handler:     c.clearSniperServerBlacklist,
		// 		},
		// 	},
		// },
		c.createConfigActionCommand(configCommandData{
			Name:         "multitokensniper",
			ReadableName: "Multi-Token Sniper",
			ValueType:    configCommandValueTypeBool,
			GetDesc:      "Get multi-token sniper status",
			GetFunc:      func(cfg *config.Config) any { return cfg.Snipers.MultiTokenSniper },
			SetDesc:      "Set multi-token sniper status",
			SetFunc:      func(cfg *config.Config, value any) { cfg.Snipers.MultiTokenSniper = value.(bool) },
			ToggleDesc:   "Toggle multi-token sniper",
			ToggleFunc:   func(cfg *config.Config) { cfg.Snipers.MultiTokenSniper = !cfg.Snipers.MultiTokenSniper },
		}),
		// snipertokens
		c.createConfigActionCommand(configCommandData{
			Name:         "privnotesniper",
			ReadableName: "Private Note Sniper",
			ValueType:    configCommandValueTypeBool,
			GetDesc:      "Get private note sniper status",
			GetFunc:      func(cfg *config.Config) any { return cfg.Snipers.PrivnoteSniper },
			SetDesc:      "Set private note sniper status",
			SetFunc:      func(cfg *config.Config, value any) { cfg.Snipers.PrivnoteSniper = value.(bool) },
			ToggleDesc:   "Toggle private note sniper",
			ToggleFunc:   func(cfg *config.Config) { cfg.Snipers.PrivnoteSniper = !cfg.Snipers.PrivnoteSniper },
		}),
		c.createConfigActionCommand(configCommandData{
			Name:         "poketwosniper",
			ReadableName: "Poketwo Sniper",
			ValueType:    configCommandValueTypeBool,
			GetDesc:      "Get poketwo sniper status",
			GetFunc:      func(cfg *config.Config) any { return cfg.Snipers.PokeTwoSniper.Enabled },
			SetDesc:      "Set poketwo sniper status",
			SetFunc:      func(cfg *config.Config, value any) { cfg.Snipers.PokeTwoSniper.Enabled = value.(bool) },
			ToggleDesc:   "Toggle poketwo sniper",
			ToggleFunc:   func(cfg *config.Config) { cfg.Snipers.PokeTwoSniper.Enabled = !cfg.Snipers.PokeTwoSniper.Enabled },
		}),
		c.createConfigActionCommand(configCommandData{
			Name:         "poketwosleep",
			ReadableName: "Poketwo Sniper Sleep",
			ValueType:    configCommandValueTypeDuration,
			GetDesc:      "Get poketwo sniper sleep duration",
			GetFunc:      func(cfg *config.Config) any { return cfg.Snipers.PokeTwoSniper.Sleep },
			SetDesc:      "Set poketwo sniper sleep duration",
			SetFunc:      func(cfg *config.Config, value any) { cfg.Snipers.PokeTwoSniper.Sleep = value.(int) },
		}),
		c.createConfigActionCommand(configCommandData{
			Name:         "poketworandomsleep",
			ReadableName: "Poketwo Sniper Random Sleep",
			ValueType:    configCommandValueTypeBool,
			GetDesc:      "Get poketwo sniper random sleep status",
			GetFunc:      func(cfg *config.Config) any { return cfg.Snipers.PokeTwoSniper.RandomSleep },
			SetDesc:      "Set poketwo sniper random sleep status",
			SetFunc:      func(cfg *config.Config, value any) { cfg.Snipers.PokeTwoSniper.RandomSleep = value.(bool) },
			ToggleDesc:   "Toggle poketwo sniper random sleep",
			ToggleFunc: func(cfg *config.Config) {
				cfg.Snipers.PokeTwoSniper.RandomSleep = !cfg.Snipers.PokeTwoSniper.RandomSleep
			},
		}),
		c.createConfigActionCommand(configCommandData{
			Name:         "poketwotimeout",
			ReadableName: "Poketwo Sniper Timeout",
			ValueType:    configCommandValueTypeDuration,
			GetDesc:      "Get poketwo sniper timeout duration",
			GetFunc:      func(cfg *config.Config) any { return cfg.Snipers.PokeTwoSniper.Timeout },
			SetDesc:      "Set poketwo sniper timeout duration",
			SetFunc:      func(cfg *config.Config, value any) { cfg.Snipers.PokeTwoSniper.Timeout = value.(int) },
		}),
		c.createConfigActionCommand(configCommandData{
			Name:         "poketwotyping",
			ReadableName: "Poketwo Sniper Typing",
			ValueType:    configCommandValueTypeBool,
			GetDesc:      "Get poketwo sniper typing status",
			GetFunc:      func(cfg *config.Config) any { return cfg.Snipers.PokeTwoSniper.Typing },
			SetDesc:      "Set poketwo sniper typing status",
			SetFunc:      func(cfg *config.Config, value any) { cfg.Snipers.PokeTwoSniper.Typing = value.(bool) },
			ToggleDesc:   "Toggle poketwo sniper typing",
			ToggleFunc:   func(cfg *config.Config) { cfg.Snipers.PokeTwoSniper.Typing = !cfg.Snipers.PokeTwoSniper.Typing },
		}),
		c.createConfigActionCommand(configCommandData{
			Name:         "poketwocatchrate",
			ReadableName: "Poketwo Sniper Catch Rate",
			ValueType:    configCommandValueTypePercentage,
			GetDesc:      "Get poketwo sniper catch rate",
			GetFunc:      func(cfg *config.Config) any { return cfg.Snipers.PokeTwoSniper.CatchRate },
			SetDesc:      "Set poketwo sniper catch rate",
			SetFunc:      func(cfg *config.Config, value any) { cfg.Snipers.PokeTwoSniper.CatchRate = value.(int) },
		}),
		c.createConfigActionCommand(configCommandData{
			Name:         "poketwofakeactivity",
			ReadableName: "Poketwo Sniper Fake Activity",
			ValueType:    configCommandValueTypeBool,
			GetDesc:      "Get poketwo sniper fake activity status",
			GetFunc:      func(cfg *config.Config) any { return cfg.Snipers.PokeTwoSniper.FakeActivity },
			SetDesc:      "Set poketwo sniper fake activity status",
			SetFunc:      func(cfg *config.Config, value any) { cfg.Snipers.PokeTwoSniper.FakeActivity = value.(bool) },
			ToggleDesc:   "Toggle poketwo sniper fake activity",
			ToggleFunc: func(cfg *config.Config) {
				cfg.Snipers.PokeTwoSniper.FakeActivity = !cfg.Snipers.PokeTwoSniper.FakeActivity
			},
		}),
		// poketwochannels
		c.createConfigActionCommand(configCommandData{
			Name:         "poketwolanguage",
			ReadableName: "Poketwo Sniper Language",
			ValueType:    configCommandValueTypeString,
			GetDesc:      "Get poketwo sniper language",
			GetFunc:      func(cfg *config.Config) any { return cfg.Snipers.PokeTwoSniper.CatchLanguage },
			SetDesc:      "Set poketwo sniper language",
			SetFunc:      func(cfg *config.Config, value any) { cfg.Snipers.PokeTwoSniper.CatchLanguage = value.(string) },
			ValuesDesc:   "List available languages",
			ValuesFunc:   func() []any { return []any{"en", "ja", "ja_t", "fr", "de"} },
		}),
		// setupwebhooks
		// wall
		c.createConfigActionCommand(configCommandData{
			Name:         "wnitro",
			ReadableName: "Nitro Sniper Webhook",
			ValueType:    configCommandValueTypeString,
			GetDesc:      "Get nitro sniper webhook",
			GetFunc:      func(cfg *config.Config) any { return cfg.Webhooks.NitroSniper },
			SetDesc:      "Set nitro sniper webhook",
			SetFunc:      func(cfg *config.Config, value any) { cfg.Webhooks.NitroSniper = value.(string) },
		}),
		c.createConfigActionCommand(configCommandData{
			Name:         "wgiveawaysniper",
			ReadableName: "Giveaway Sniper Webhook",
			ValueType:    configCommandValueTypeString,
			GetDesc:      "Get giveaway sniper webhook",
			GetFunc:      func(cfg *config.Config) any { return cfg.Webhooks.GiveawaySniper },
			SetDesc:      "Set giveaway sniper webhook",
			SetFunc:      func(cfg *config.Config, value any) { cfg.Webhooks.GiveawaySniper = value.(string) },
		}),
		c.createConfigActionCommand(configCommandData{
			Name:         "wselfbotdetection",
			ReadableName: "Selfbot Detection Webhook",
			ValueType:    configCommandValueTypeString,
			GetDesc:      "Get selfbot detection webhook",
			GetFunc:      func(cfg *config.Config) any { return cfg.Webhooks.SelfbotDetection },
			SetDesc:      "Set selfbot detection webhook",
			SetFunc:      func(cfg *config.Config, value any) { cfg.Webhooks.SelfbotDetection = value.(string) },
		}),
		c.createConfigActionCommand(configCommandData{
			Name:         "wstaffdetection",
			ReadableName: "Staff Detection Webhook",
			ValueType:    configCommandValueTypeString,
			GetDesc:      "Get staff detection webhook",
			GetFunc:      func(cfg *config.Config) any { return cfg.Webhooks.StaffDetection },
			SetDesc:      "Set staff detection webhook",
			SetFunc:      func(cfg *config.Config, value any) { cfg.Webhooks.StaffDetection = value.(string) },
		}),
		c.createConfigActionCommand(configCommandData{
			Name:         "wnotifiers",
			ReadableName: "Notifiers Webhook",
			ValueType:    configCommandValueTypeString,
			GetDesc:      "Get notifiers webhook",
			GetFunc:      func(cfg *config.Config) any { return cfg.Webhooks.Notifiers },
			SetDesc:      "Set notifiers webhook",
			SetFunc:      func(cfg *config.Config, value any) { cfg.Webhooks.Notifiers = value.(string) },
		}),
		c.createConfigActionCommand(configCommandData{
			Name:         "wpingdetection",
			ReadableName: "Ping Detection Webhook",
			ValueType:    configCommandValueTypeString,
			GetDesc:      "Get ping detection webhook",
			GetFunc:      func(cfg *config.Config) any { return cfg.Webhooks.PingDetection },
			SetDesc:      "Set ping detection webhook",
			SetFunc:      func(cfg *config.Config, value any) { cfg.Webhooks.PingDetection = value.(string) },
		}),
		c.createConfigActionCommand(configCommandData{
			Name:         "wmisc",
			ReadableName: "Misc Webhook",
			ValueType:    configCommandValueTypeString,
			GetDesc:      "Get misc webhook",
			GetFunc:      func(cfg *config.Config) any { return cfg.Webhooks.Misc },
			SetDesc:      "Set misc webhook",
			SetFunc:      func(cfg *config.Config, value any) { cfg.Webhooks.Misc = value.(string) },
		}),
		c.createConfigActionCommand(configCommandData{
			Name:         "wsessions",
			ReadableName: "Sessions Webhook",
			ValueType:    configCommandValueTypeString,
			GetDesc:      "Get sessions webhook",
			GetFunc:      func(cfg *config.Config) any { return cfg.Webhooks.Sessions },
			SetDesc:      "Set sessions webhook",
			SetFunc:      func(cfg *config.Config, value any) { cfg.Webhooks.Sessions = value.(string) },
		}),
		c.createConfigActionCommand(configCommandData{
			Name:         "wpingself",
			ReadableName: "Ping Self Webhook",
			ValueType:    configCommandValueTypeString,
			GetDesc:      "Get ping self webhook status",
			GetFunc:      func(cfg *config.Config) any { return cfg.Webhooks.PingSelfMode },
			SetDesc:      "Set ping self webhook status",
			SetFunc:      func(cfg *config.Config, value any) { cfg.Webhooks.PingSelfMode = value.(string) },
			ValuesDesc:   "List available ping self webhook modes",
			ValuesFunc:   func() []any { return []any{"Disabled, Always, Important"} },
		}),
		c.createConfigActionCommand(configCommandData{
			Name:         "wpingeveryone",
			ReadableName: "Ping Everyone Webhook",
			ValueType:    configCommandValueTypeString,
			GetDesc:      "Get ping everyone webhook status",
			GetFunc:      func(cfg *config.Config) any { return cfg.Webhooks.PingEveryoneMode },
			SetDesc:      "Set ping everyone webhook status",
			SetFunc:      func(cfg *config.Config, value any) { cfg.Webhooks.PingEveryoneMode = value.(string) },
			ValuesDesc:   "List available ping everyone webhook modes",
			ValuesFunc:   func() []any { return []any{"Disabled, Always, Important"} },
		}),
		c.createConfigActionCommand(configCommandData{
			Name:         "bannotifier",
			ReadableName: "Ban Notifier",
			ValueType:    configCommandValueTypeBool,
			GetDesc:      "Get ban notifier status",
			GetFunc:      func(cfg *config.Config) any { return cfg.Events.BanNotifier },
			SetDesc:      "Set ban notifier status",
			SetFunc:      func(cfg *config.Config, value any) { cfg.Events.BanNotifier = value.(bool) },
			ToggleDesc:   "Toggle ban notifier",
			ToggleFunc:   func(cfg *config.Config) { cfg.Events.BanNotifier = !cfg.Events.BanNotifier },
		}),
		c.createConfigActionCommand(configCommandData{
			Name:         "staffdetection",
			ReadableName: "Staff Detection",
			ValueType:    configCommandValueTypeBool,
			GetDesc:      "Get staff detection status",
			GetFunc:      func(cfg *config.Config) any { return cfg.Events.DiscordStaffDetection },
			SetDesc:      "Set staff detection status",
			SetFunc:      func(cfg *config.Config, value any) { cfg.Events.DiscordStaffDetection = value.(bool) },
			ToggleDesc:   "Toggle staff detection",
			ToggleFunc:   func(cfg *config.Config) { cfg.Events.DiscordStaffDetection = !cfg.Events.DiscordStaffDetection },
		}),
		c.createConfigActionCommand(configCommandData{
			Name:         "pingdetection",
			ReadableName: "Ping Detection",
			ValueType:    configCommandValueTypeBool,
			GetDesc:      "Get ping detection status",
			GetFunc:      func(cfg *config.Config) any { return cfg.Events.GhostPingDetection },
			SetDesc:      "Set ping detection status",
			SetFunc:      func(cfg *config.Config, value any) { cfg.Events.GhostPingDetection = value.(bool) },
			ToggleDesc:   "Toggle ping detection",
			ToggleFunc:   func(cfg *config.Config) { cfg.Events.GhostPingDetection = !cfg.Events.GhostPingDetection },
		}),
		c.createConfigActionCommand(configCommandData{
			Name:         "ticketnotifier",
			ReadableName: "Ticket Notifier",
			ValueType:    configCommandValueTypeBool,
			GetDesc:      "Get ticket notifier status",
			GetFunc:      func(cfg *config.Config) any { return cfg.Events.TicketNotifier },
			SetDesc:      "Set ticket notifier status",
			SetFunc:      func(cfg *config.Config, value any) { cfg.Events.TicketNotifier = value.(bool) },
			ToggleDesc:   "Toggle ticket notifier",
			ToggleFunc:   func(cfg *config.Config) { cfg.Events.TicketNotifier = !cfg.Events.TicketNotifier },
		}),
		c.createConfigActionCommand(configCommandData{
			Name:         "sessionconnected",
			ReadableName: "Session Connected",
			ValueType:    configCommandValueTypeBool,
			GetDesc:      "Get session connected status",
			GetFunc:      func(cfg *config.Config) any { return cfg.Events.SessionConnected },
			SetDesc:      "Set session connected status",
			SetFunc:      func(cfg *config.Config, value any) { cfg.Events.SessionConnected = value.(bool) },
			ToggleDesc:   "Toggle session connected",
			ToggleFunc:   func(cfg *config.Config) { cfg.Events.SessionConnected = !cfg.Events.SessionConnected },
		}),
		c.createConfigActionCommand(configCommandData{
			Name:         "sessiondisconnected",
			ReadableName: "Session Disconnected",
			ValueType:    configCommandValueTypeBool,
			GetDesc:      "Get session disconnected status",
			GetFunc:      func(cfg *config.Config) any { return cfg.Events.SessionDisconnected },
			SetDesc:      "Set session disconnected status",
			SetFunc:      func(cfg *config.Config, value any) { cfg.Events.SessionDisconnected = value.(bool) },
			ToggleDesc:   "Toggle session disconnected",
			ToggleFunc:   func(cfg *config.Config) { cfg.Events.SessionDisconnected = !cfg.Events.SessionDisconnected },
		}),
		c.createConfigActionCommand(configCommandData{
			Name:         "selfbotdetection",
			ReadableName: "Selfbot Detection",
			ValueType:    configCommandValueTypeBool,
			GetDesc:      "Get selfbot detection status",
			GetFunc:      func(cfg *config.Config) any { return cfg.Events.SelfbotDetection },
			SetDesc:      "Set selfbot detection status",
			SetFunc:      func(cfg *config.Config, value any) { cfg.Events.SelfbotDetection = value.(bool) },
			ToggleDesc:   "Toggle selfbot detection",
			ToggleFunc:   func(cfg *config.Config) { cfg.Events.SelfbotDetection = !cfg.Events.SelfbotDetection },
		}),
		c.createConfigActionCommand(configCommandData{
			Name:         "dmtypingnotifier",
			ReadableName: "DM Typing Notifier",
			ValueType:    configCommandValueTypeBool,
			GetDesc:      "Get DM typing notifier status",
			GetFunc:      func(cfg *config.Config) any { return cfg.Events.DMTypingNotifier },
			SetDesc:      "Set DM typing notifier status",
			SetFunc:      func(cfg *config.Config, value any) { cfg.Events.DMTypingNotifier = value.(bool) },
			ToggleDesc:   "Toggle DM typing notifier",
			ToggleFunc:   func(cfg *config.Config) { cfg.Events.DMTypingNotifier = !cfg.Events.DMTypingNotifier },
		}),
		c.createConfigActionCommand(configCommandData{
			Name:         "dmnotifier",
			ReadableName: "DM Notifier",
			ValueType:    configCommandValueTypeBool,
			GetDesc:      "Get DM notifier status",
			GetFunc:      func(cfg *config.Config) any { return cfg.Events.DMNotifier },
			SetDesc:      "Set DM notifier status",
			SetFunc:      func(cfg *config.Config, value any) { cfg.Events.DMNotifier = value.(bool) },
			ToggleDesc:   "Toggle DM notifier",
			ToggleFunc:   func(cfg *config.Config) { cfg.Events.DMNotifier = !cfg.Events.DMNotifier },
		}),
		c.createConfigActionCommand(configCommandData{
			Name:         "dmeditnotifier",
			ReadableName: "DM Edit Notifier",
			ValueType:    configCommandValueTypeBool,
			GetDesc:      "Get DM edit notifier status",
			GetFunc:      func(cfg *config.Config) any { return cfg.Events.DMEditNotifier },
			SetDesc:      "Set DM edit notifier status",
			SetFunc:      func(cfg *config.Config, value any) { cfg.Events.DMEditNotifier = value.(bool) },
			ToggleDesc:   "Toggle DM edit notifier",
			ToggleFunc:   func(cfg *config.Config) { cfg.Events.DMEditNotifier = !cfg.Events.DMEditNotifier },
		}),
		c.createConfigActionCommand(configCommandData{
			Name:         "dmdeletenotifier",
			ReadableName: "DM Delete Notifier",
			ValueType:    configCommandValueTypeBool,
			GetDesc:      "Get DM delete notifier status",
			GetFunc:      func(cfg *config.Config) any { return cfg.Events.DMDeleteNotifier },
			SetDesc:      "Set DM delete notifier status",
			SetFunc:      func(cfg *config.Config, value any) { cfg.Events.DMDeleteNotifier = value.(bool) },
			ToggleDesc:   "Toggle DM delete notifier",
			ToggleFunc:   func(cfg *config.Config) { cfg.Events.DMDeleteNotifier = !cfg.Events.DMDeleteNotifier },
		}),
		c.createConfigActionCommand(configCommandData{
			Name:         "messagedeletenotifier",
			ReadableName: "Message Delete Notifier",
			ValueType:    configCommandValueTypeBool,
			GetDesc:      "Get message delete notifier status",
			GetFunc:      func(cfg *config.Config) any { return cfg.Events.MessageDeleteNotifier },
			SetDesc:      "Set message delete notifier status",
			SetFunc:      func(cfg *config.Config, value any) { cfg.Events.MessageDeleteNotifier = value.(bool) },
			ToggleDesc:   "Toggle message delete notifier",
			ToggleFunc:   func(cfg *config.Config) { cfg.Events.MessageDeleteNotifier = !cfg.Events.MessageDeleteNotifier },
		}),
		c.createConfigActionCommand(configCommandData{
			Name:         "uitheme",
			ReadableName: "UI Theme",
			ValueType:    configCommandValueTypeString,
			GetDesc:      "Get UI theme",
			GetFunc:      func(cfg *config.Config) any { return cfg.UI.Theme },
			SetDesc:      "Set UI theme",
			SetFunc:      func(cfg *config.Config, value any) { cfg.UI.Theme = value.(string) },
			ValuesDesc:   "List available UI themes",
			ValuesFunc:   func() []any { return []any{"Dark", "Light"} },
		}),
		c.createConfigActionCommand(configCommandData{
			Name:         "hidetosystray",
			ReadableName: "Hide to System Tray",
			ValueType:    configCommandValueTypeBool,
			GetDesc:      "Get hide to system tray status",
			GetFunc:      func(cfg *config.Config) any { return cfg.UI.HideToSystemTray },
			SetDesc:      "Set hide to system tray status",
			SetFunc:      func(cfg *config.Config, value any) { cfg.UI.HideToSystemTray = value.(bool) },
			ToggleDesc:   "Toggle hide to system tray",
			ToggleFunc:   func(cfg *config.Config) { cfg.UI.HideToSystemTray = !cfg.UI.HideToSystemTray },
		}),
		c.createConfigActionCommand(configCommandData{
			Name:         "runuionstartup",
			ReadableName: "Run UI on Startup",
			ValueType:    configCommandValueTypeBool,
			GetDesc:      "Get run UI on startup status",
			GetFunc:      func(cfg *config.Config) any { return cfg.UI.RunOnStartup },
			SetDesc:      "Set run UI on startup status",
			SetFunc:      func(cfg *config.Config, value any) { cfg.UI.RunOnStartup = value.(bool) },
			ToggleDesc:   "Toggle run UI on startup",
			ToggleFunc:   func(cfg *config.Config) { cfg.UI.RunOnStartup = !cfg.UI.RunOnStartup },
		}),
		c.createConfigActionCommand(configCommandData{
			Name:         "runbotonstartup",
			ReadableName: "Run Bot on Startup",
			ValueType:    configCommandValueTypeBool,
			GetDesc:      "Get run bot on startup status",
			GetFunc:      func(cfg *config.Config) any { return cfg.UI.RunOnStartup },
			SetDesc:      "Set run bot on startup status",
			SetFunc:      func(cfg *config.Config, value any) { cfg.UI.RunOnStartup = value.(bool) },
			ToggleDesc:   "Toggle run bot on startup",
			ToggleFunc:   func(cfg *config.Config) { cfg.UI.RunOnStartup = !cfg.UI.RunOnStartup },
		}),
	})
}

type configCommandValueType int

const (
	configCommandValueTypeString configCommandValueType = iota
	configCommandValueTypeInt
	configCommandValueTypePercentage
	configCommandValueTypeBool
	configCommandValueTypeDuration
)

type configCommandData struct {
	Name         string
	ReadableName string
	ValueType    configCommandValueType

	GetDesc string
	GetFunc func(cfg *config.Config) any

	SetDesc string
	SetFunc func(cfg *config.Config, value any)

	ToggleDesc string
	ToggleFunc func(cfg *config.Config)

	ValuesDesc string
	ValuesFunc func() []any

	ExampleDesc string
	ExampleFunc func()
}

// type configCommandDataList struct {
// 	Name         string
// 	ReadableName string

// 	GetDesc string
// 	GetFunc func(ctx *router.Ctx) ([]any, error)

// 	AddDesc string
// 	AddFunc func(ctx *router.Ctx, value any) error

// 	RmDesc string
// 	RmFunc func(ctx *router.Ctx, value any) error

// 	ClearDesc string
// 	ClearFunc func(ctx *router.Ctx) error
// }

func (c *ConfigCommands) createConfigActionCommand(data configCommandData) *router.Command {
	subCommands := []*router.Command{}

	if data.GetFunc != nil {
		subCommands = append(subCommands, &router.Command{
			Name:        "get",
			Description: data.GetDesc,
			Handler: func(ctx *router.Ctx) {
				cfg, err := config.GetConfig()
				if err != nil {
					ctx.RespondText("Error getting config")
					return
				}

				value := data.GetFunc(&cfg)

				c.services.Embed.NewEmbed().
					SetDescription("%s is currently set to: %v", data.ReadableName, value).
					Send(ctx)
			},
		})
	}

	if data.SetFunc != nil {
		subCommands = append(subCommands, &router.Command{
			Name:        "set",
			Description: data.SetDesc,
			Arguments:   "[value]",
			Handler: func(ctx *router.Ctx) {
				var value any
				var err error

				switch data.ValueType {
				case configCommandValueTypeString:
					value = ctx.Arguments.Raw()
				case configCommandValueTypeInt:
					value, err = ctx.Arguments.Get(0).Int()
				case configCommandValueTypePercentage:
					value, err = ctx.Arguments.Get(0).Int()
					value = utils2.Clamp(value.(int), 0, 100)
				case configCommandValueTypeBool:
					value, err = ctx.Arguments.Get(0).Bool()
				case configCommandValueTypeDuration:
					value, err = ctx.Arguments.Get(0).Duration()
					value = int(value.(time.Duration).Seconds())
				}

				if err != nil {
					ctx.RespondText("Error parsing value")
					return
				}

				cfg, err := config.GetConfig()
				if err != nil {
					ctx.RespondText("Error getting config")
					return
				}

				data.SetFunc(&cfg, value)

				if err := cfg.Save(); err != nil {
					ctx.RespondText("Error saving config")
					return
				}

				c.services.Embed.NewEmbed().
					SetDescription("%s set to %v", data.ReadableName, value).
					Send(ctx)
			},
		})
	}

	if data.ToggleFunc != nil {
		subCommands = append(subCommands, &router.Command{
			Name:        "toggle",
			Description: data.ToggleDesc,
			Handler: func(ctx *router.Ctx) {
				cfg, err := config.GetConfig()
				if err != nil {
					ctx.RespondText("Error getting config")
					return
				}

				data.ToggleFunc(&cfg)

				if err := cfg.Save(); err != nil {
					ctx.RespondText("Error saving config")
					return
				}

				c.services.Embed.NewEmbed().
					SetDescription("%s toggled", data.ReadableName).
					Send(ctx)
			},
		})
	}

	if data.ValuesFunc != nil {
		subCommands = append(subCommands, &router.Command{
			Name:        "values",
			Description: data.ValuesDesc,
			Handler: func(ctx *router.Ctx) {
				values := data.ValuesFunc()
				stringValues := utils2.Map(values, func(value any) string {
					return fmt.Sprint(value)
				})

				c.services.Embed.NewEmbed().
					AddField(data.ReadableName, strings.Join(stringValues, "\n")).
					Send(ctx)
			},
		})
	}

	return &router.Command{
		Name:        data.Name,
		Flags:       []router.CommandFlag{router.CmdFlagDangerous},
		SubCommands: subCommands,
	}
}
