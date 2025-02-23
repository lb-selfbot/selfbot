package config

import (
	"fmt"
	"slices"
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
			ValueType:    configCommandValueTypeInt,
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
		{
			Name:        "aliases",
			Description: "Command aliases commands",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			SubCommands: []*router.Command{
				{
					Name:        "get",
					Description: "Get a command alias",
					Arguments:   "<alias>",
					Handler:     c.getAlias,
				},
				{
					Name:        "set",
					Description: "Set a command alias",
					Arguments:   "<alias> <command>",
					Handler:     c.setAlias,
				},
				{
					Name:        "rm",
					Description: "Remove a command alias",
					Arguments:   "<alias>",
					Handler:     c.removeAlias,
				},
				{
					Name:        "list",
					Description: "List all command aliases",
					Handler:     c.listAliases,
				},
				{
					Name:        "clear",
					Description: "Clear all command aliases",
					Handler:     c.clearAliases,
				},
			},
		},
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
			ValueType:    configCommandValueTypeInt,
			GetDesc:      "Get giveaway sniper sleep duration",
			GetFunc:      func(cfg *config.Config) any { return cfg.Snipers.GiveawaySniper.Sleep },
			SetDesc:      "Set giveaway sniper sleep duration",
			SetFunc:      func(cfg *config.Config, value any) { cfg.Snipers.GiveawaySniper.Sleep = value.(int) },
		}),
		{
			Name:        "customgiveawaybots",
			Description: "Custom giveaway bots commands",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			SubCommands: []*router.Command{
				{
					Name:        "add",
					Description: "Add a custom giveaway bot",
					Arguments:   "<bot id>",
					Handler:     c.addCustomGiveawayBot,
				},
				{
					Name:        "rm",
					Description: "Remove a custom giveaway bot",
					Arguments:   "<bot id>",
					Handler:     c.removeCustomGiveawayBot,
				},
				{
					Name:        "list",
					Description: "List custom giveaway bots",
					Handler:     c.listCustomGiveawayBots,
				},
				{
					Name:        "clear",
					Description: "Clear all custom giveaway bots",
					Handler:     c.clearCustomGiveawayBots,
				},
			},
		},
		{
			Name:        "sniperblacklist",
			Description: "Sniper blacklist commands",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			SubCommands: []*router.Command{
				{
					Name:        "add",
					Description: "Add a user to sniper blacklist",
					Arguments:   "<user>",
					Handler:     c.addSniperBlacklist,
				},
				{
					Name:        "rm",
					Description: "Remove a user from sniper blacklist",
					Arguments:   "<user>",
					Handler:     c.removeSniperBlacklist,
				},
				{
					Name:        "list",
					Description: "List sniper blacklist",
					Handler:     c.listSniperBlacklist,
				},
				{
					Name:        "clear",
					Description: "Clear sniper blacklist",
					Handler:     c.clearSniperBlacklist,
				},
			},
		},
		{
			Name:        "sniperserverblacklist",
			Description: "Sniper server blacklist commands",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			SubCommands: []*router.Command{
				{
					Name:        "add",
					Description: "Add a server to sniper blacklist",
					Arguments:   "<server>",
					Handler:     c.addSniperServerBlacklist,
				},
				{
					Name:        "rm",
					Description: "Remove a server from sniper blacklist",
					Arguments:   "<server>",
					Handler:     c.removeSniperServerBlacklist,
				},
				{
					Name:        "list",
					Description: "List sniper server blacklist",
					Handler:     c.listSniperServerBlacklist,
				},
				{
					Name:        "clear",
					Description: "Clear sniper server blacklist",
					Handler:     c.clearSniperServerBlacklist,
				},
			},
		},
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
	})
}

type configCommandValueType int

const (
	configCommandValueTypeString configCommandValueType = iota
	configCommandValueTypeInt
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

				switch data.ValueType {
				case configCommandValueTypeString:
					value = ctx.Arguments.Raw()
				case configCommandValueTypeInt:
					value, _ = ctx.Arguments.Get(0).Int()
				case configCommandValueTypeBool:
					value, _ = ctx.Arguments.Get(0).Bool()
				case configCommandValueTypeDuration:
					value, _ = ctx.Arguments.Get(0).Duration()
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

func (c *ConfigCommands) getAlias(ctx *router.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		ctx.RespondText("Error getting config")
		return
	}

	alias := ctx.Arguments.Raw()
	if command, ok := cfg.Selfbot.CommandAliases[alias]; ok {
		c.services.Embed.NewEmbed().
			AddField("Alias", fmt.Sprintf("%s -> %s", alias, command)).
			Send(ctx)
	} else {
		ctx.RespondText("Alias not found")
	}
}

func (c *ConfigCommands) setAlias(ctx *router.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		ctx.RespondText("Error getting config")
		return
	}

	alias := ctx.Arguments.Get(0).Raw()
	command := ctx.Arguments.RawFrom(1)

	if cfg.Selfbot.CommandAliases == nil {
		cfg.Selfbot.CommandAliases = make(map[string]string)
	}

	cfg.Selfbot.CommandAliases[alias] = command
	_ = cfg.Save()

	c.services.Embed.NewEmbed().
		SetDescription("Set alias `%s` -> `%s`", alias, command).
		Send(ctx)
}

func (c *ConfigCommands) removeAlias(ctx *router.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		ctx.RespondText("Error getting config")
		return
	}

	alias := ctx.Arguments.Raw()

	if _, ok := cfg.Selfbot.CommandAliases[alias]; !ok {
		ctx.RespondText("Alias not found")
		return
	}

	delete(cfg.Selfbot.CommandAliases, alias)
	_ = cfg.Save()

	c.services.Embed.NewEmbed().
		SetDescription("Removed alias `%s`", alias).
		Send(ctx)
}

func (c *ConfigCommands) listAliases(ctx *router.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		ctx.RespondText("Error getting config")
		return
	}

	if len(cfg.Selfbot.CommandAliases) == 0 {
		ctx.RespondText("No aliases set")
		return
	}

	var aliases []string
	for k, v := range cfg.Selfbot.CommandAliases {
		aliases = append(aliases, fmt.Sprintf("%s -> %s", k, v))
	}

	c.services.Embed.NewEmbed().
		AddField("Aliases", strings.Join(aliases, "\n")).
		Send(ctx)
}

func (c *ConfigCommands) clearAliases(ctx *router.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		ctx.RespondText("Error getting config")
		return
	}

	cfg.Selfbot.CommandAliases = make(map[string]string)
	_ = cfg.Save()

	c.services.Embed.NewEmbed().
		SetDescription("Cleared all aliases").
		Send(ctx)
}

func (c *ConfigCommands) addCustomGiveawayBot(ctx *router.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		ctx.RespondText("Error getting config")
		return
	}

	botID := ctx.Arguments.Get(0).Raw()
	if botID == "" {
		ctx.RespondText("You must provide a bot ID!")
		return
	}

	if cfg.Snipers.GiveawaySniper.CustomBots == nil {
		cfg.Snipers.GiveawaySniper.CustomBots = []string{}
	}

	if slices.Contains(cfg.Snipers.GiveawaySniper.CustomBots, botID) {
		ctx.RespondText("Bot already in list!")
		return
	}

	cfg.Snipers.GiveawaySniper.CustomBots = append(cfg.Snipers.GiveawaySniper.CustomBots, botID)
	_ = cfg.Save()

	c.services.Embed.NewEmbed().
		SetDescription("Added bot to custom giveaway bots").
		Send(ctx)
}

func (c *ConfigCommands) removeCustomGiveawayBot(ctx *router.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		ctx.RespondText("Error getting config")
		return
	}

	botID := ctx.Arguments.Get(0).Raw()
	if botID == "" {
		ctx.RespondText("You must provide a bot ID!")
		return
	}

	index := slices.Index(cfg.Snipers.GiveawaySniper.CustomBots, botID)
	if index == -1 {
		ctx.RespondText("Bot not in list!")
		return
	}

	cfg.Snipers.GiveawaySniper.CustomBots = slices.Delete(cfg.Snipers.GiveawaySniper.CustomBots, index, index+1)
	_ = cfg.Save()

	c.services.Embed.NewEmbed().
		SetDescription("Removed bot from custom giveaway bots").
		Send(ctx)
}

func (c *ConfigCommands) listCustomGiveawayBots(ctx *router.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		ctx.RespondText("Error getting config")
		return
	}

	if len(cfg.Snipers.GiveawaySniper.CustomBots) == 0 {
		ctx.RespondText("No custom giveaway bots")
		return
	}

	c.services.Embed.NewEmbed().
		AddField("Custom Giveaway Bots", strings.Join(cfg.Snipers.GiveawaySniper.CustomBots, "\n")).
		Send(ctx)
}

func (c *ConfigCommands) clearCustomGiveawayBots(ctx *router.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		ctx.RespondText("Error getting config")
		return
	}

	cfg.Snipers.GiveawaySniper.CustomBots = []string{}
	_ = cfg.Save()

	c.services.Embed.NewEmbed().
		SetDescription("Cleared custom giveaway bots").
		Send(ctx)
}

func (c *ConfigCommands) addSniperBlacklist(ctx *router.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		ctx.RespondText("Error getting config")
		return
	}

	user, err := ctx.Arguments.Get(0).User()
	if err != nil {
		ctx.RespondText("Could not find user")
		return
	}

	if cfg.Snipers.Blacklist == nil {
		cfg.Snipers.Blacklist = []string{}
	}

	if slices.Contains(cfg.Snipers.Blacklist, user.ID) {
		ctx.RespondText("User already in blacklist!")
		return
	}

	cfg.Snipers.Blacklist = append(cfg.Snipers.Blacklist, user.ID)
	_ = cfg.Save()

	c.services.Embed.NewEmbed().
		SetDescription("Added user to sniper blacklist").
		Send(ctx)
}

func (c *ConfigCommands) removeSniperBlacklist(ctx *router.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		ctx.RespondText("Error getting config")
		return
	}

	user, err := ctx.Arguments.Get(0).User()
	if err != nil {
		ctx.RespondText("Could not find user")
		return
	}

	index := slices.Index(cfg.Snipers.Blacklist, user.ID)
	if index == -1 {
		ctx.RespondText("User not in blacklist!")
		return
	}

	cfg.Snipers.Blacklist = slices.Delete(cfg.Snipers.Blacklist, index, index+1)
	_ = cfg.Save()

	c.services.Embed.NewEmbed().
		SetDescription("Removed user from sniper blacklist").
		Send(ctx)
}

func (c *ConfigCommands) listSniperBlacklist(ctx *router.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		ctx.RespondText("Error getting config")
		return
	}

	if len(cfg.Snipers.Blacklist) == 0 {
		ctx.RespondText("No users in sniper blacklist")
		return
	}

	var users []string
	for _, id := range cfg.Snipers.Blacklist {
		users = append(users, fmt.Sprintf("<@%s>", id))
	}

	c.services.Embed.NewEmbed().
		AddField("Sniper Blacklist", strings.Join(users, "\n")).
		Send(ctx)
}

func (c *ConfigCommands) clearSniperBlacklist(ctx *router.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		ctx.RespondText("Error getting config")
		return
	}

	cfg.Snipers.Blacklist = []string{}
	_ = cfg.Save()

	c.services.Embed.NewEmbed().
		SetDescription("Cleared sniper blacklist").
		Send(ctx)
}

func (c *ConfigCommands) addSniperServerBlacklist(ctx *router.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		ctx.RespondText("Error getting config")
		return
	}

	guild, err := ctx.Arguments.Get(0).Guild()
	if err != nil {
		ctx.RespondText("Could not find server")
		return
	}

	if cfg.Snipers.ServerBlacklist == nil {
		cfg.Snipers.ServerBlacklist = []string{}
	}

	if slices.Contains(cfg.Snipers.ServerBlacklist, guild.ID) {
		ctx.RespondText("Server already in blacklist!")
		return
	}

	cfg.Snipers.ServerBlacklist = append(cfg.Snipers.ServerBlacklist, guild.ID)
	_ = cfg.Save()

	c.services.Embed.NewEmbed().
		SetDescription("Added server to sniper blacklist").
		Send(ctx)
}

func (c *ConfigCommands) removeSniperServerBlacklist(ctx *router.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		ctx.RespondText("Error getting config")
		return
	}

	guild, err := ctx.Arguments.Get(0).Guild()
	if err != nil {
		ctx.RespondText("Could not find server")
		return
	}

	index := slices.Index(cfg.Snipers.ServerBlacklist, guild.ID)
	if index == -1 {
		ctx.RespondText("Server not in blacklist!")
		return
	}

	cfg.Snipers.ServerBlacklist = slices.Delete(cfg.Snipers.ServerBlacklist, index, index+1)
	_ = cfg.Save()

	c.services.Embed.NewEmbed().
		SetDescription("Removed server from sniper blacklist").
		Send(ctx)
}

func (c *ConfigCommands) listSniperServerBlacklist(ctx *router.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		ctx.RespondText("Error getting config")
		return
	}

	if len(cfg.Snipers.ServerBlacklist) == 0 {
		ctx.RespondText("No servers in sniper blacklist")
		return
	}

	var servers []string
	for _, id := range cfg.Snipers.ServerBlacklist {
		guild, err := ctx.Session.State.Guild(id)
		if err != nil {
			servers = append(servers, id)
		} else {
			servers = append(servers, guild.Name)
		}
	}

	c.services.Embed.NewEmbed().
		AddField("Sniper Server Blacklist", strings.Join(servers, "\n")).
		Send(ctx)
}

func (c *ConfigCommands) clearSniperServerBlacklist(ctx *router.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		ctx.RespondText("Error getting config")
		return
	}

	cfg.Snipers.ServerBlacklist = []string{}
	_ = cfg.Save()

	c.services.Embed.NewEmbed().
		SetDescription("Cleared sniper server blacklist").
		Send(ctx)
}
