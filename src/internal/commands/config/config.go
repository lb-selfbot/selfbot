package config

import (
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/LightningDev1/LB-Selfbot/internal/config"
	"github.com/LightningDev1/LB-Selfbot/internal/router"
	"github.com/LightningDev1/LB-Selfbot/internal/services"

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
		Description: "Change bot configuration",
	}, []*router.Command{
		{
			Name:        "commandprefix",
			Description: "Command prefix commands",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			SubCommands: []*router.Command{
				{
					Name:        "get",
					Description: "Get the current command prefix",
					Handler:     c.getCommandPrefix,
				},
				{
					Name:        "set",
					Description: "Set the command prefix",
					Arguments:   "<prefix>",
					Handler:     c.setCommandPrefix,
				},
			},
		},
		{
			Name:        "commandmode",
			Description: "Command mode commands",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			SubCommands: []*router.Command{
				{
					Name:        "get",
					Description: "Get the current command mode",
					Handler:     c.getCommandMode,
				},
				{
					Name:        "set",
					Description: "Set the command mode",
					Arguments:   "<mode>",
					Handler:     c.setCommandMode,
				},
				{
					Name:        "values",
					Description: "List available command modes",
					Handler:     c.listCommandModes,
				},
			},
		},
		{
			Name:        "title",
			Description: "Embed title commands",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			SubCommands: []*router.Command{
				{
					Name:        "get",
					Description: "Get the current embed title",
					Handler:     c.getTitle,
				},
				{
					Name:        "set",
					Description: "Set the embed title",
					Arguments:   "<title>",
					Handler:     c.setTitle,
				},
			},
		},
		{
			Name:        "footer",
			Description: "Embed footer commands",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			SubCommands: []*router.Command{
				{
					Name:        "get",
					Description: "Get the current embed footer",
					Handler:     c.getFooter,
				},
				{
					Name:        "set",
					Description: "Set the embed footer",
					Arguments:   "<footer>",
					Handler:     c.setFooter,
				},
			},
		},
		{
			Name:        "delresponses",
			Description: "Delete responses commands",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			SubCommands: []*router.Command{
				{
					Name:        "get",
					Description: "Get delete responses status",
					Handler:     c.getDelResponses,
				},
				{
					Name:        "set",
					Description: "Set delete responses status",
					Arguments:   "<on/off>",
					Handler:     c.setDelResponses,
				},
				{
					Name:        "toggle",
					Description: "Toggle delete responses",
					Handler:     c.toggleDelResponses,
				},
			},
		},
		{
			Name:        "delresponsesdelay",
			Description: "Delete responses delay commands",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			SubCommands: []*router.Command{
				{
					Name:        "get",
					Description: "Get delete responses delay",
					Handler:     c.getDelResponsesDelay,
				},
				{
					Name:        "set",
					Description: "Set delete responses delay",
					Arguments:   "<seconds>",
					Handler:     c.setDelResponsesDelay,
				},
			},
		},
		{
			Name:        "delcommands",
			Description: "Delete commands commands",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			SubCommands: []*router.Command{
				{
					Name:        "get",
					Description: "Get delete commands status",
					Handler:     c.getDelCommands,
				},
				{
					Name:        "set",
					Description: "Set delete commands status",
					Arguments:   "<on/off>",
					Handler:     c.setDelCommands,
				},
				{
					Name:        "toggle",
					Description: "Toggle delete commands",
					Handler:     c.toggleDelCommands,
				},
			},
		},
		{
			Name:        "clientspoof",
			Description: "Client spoof commands",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			SubCommands: []*router.Command{
				{
					Name:        "get",
					Description: "Get client spoof status",
					Handler:     c.getClientSpoof,
				},
				{
					Name:        "set",
					Description: "Set client spoof",
					Arguments:   "<on/off>",
					Handler:     c.setClientSpoof,
				},
				{
					Name:        "toggle",
					Description: "Toggle client spoof",
					Handler:     c.toggleClientSpoof,
				},
			},
		},
		{
			Name:        "timezone",
			Description: "Timezone commands",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			SubCommands: []*router.Command{
				{
					Name:        "get",
					Description: "Get current timezone",
					Handler:     c.getTimezone,
				},
				{
					Name:        "set",
					Description: "Set timezone",
					Arguments:   "<timezone>",
					Handler:     c.setTimezone,
				},
			},
		},
		{
			Name:        "embedcolor",
			Description: "Embed color commands",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			SubCommands: []*router.Command{
				{
					Name:        "get",
					Description: "Get current embed color",
					Handler:     c.getEmbedColor,
				},
				{
					Name:        "set",
					Description: "Set embed color",
					Arguments:   "<color>",
					Handler:     c.setEmbedColor,
				},
			},
		},
		{
			Name:        "embedimage",
			Description: "Embed image commands",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			SubCommands: []*router.Command{
				{
					Name:        "get",
					Description: "Get current embed image",
					Handler:     c.getEmbedImage,
				},
				{
					Name:        "set",
					Description: "Set embed image",
					Arguments:   "<url>",
					Handler:     c.setEmbedImage,
				},
			},
		},
		{
			Name:        "embedurl",
			Description: "Embed URL commands",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			SubCommands: []*router.Command{
				{
					Name:        "get",
					Description: "Get current embed URL",
					Handler:     c.getEmbedURL,
				},
				{
					Name:        "set",
					Description: "Set embed URL",
					Arguments:   "<url>",
					Handler:     c.setEmbedURL,
				},
			},
		},
		{
			Name:        "forcesubscribedservers",
			Description: "Force subscribed servers commands",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			SubCommands: []*router.Command{
				{
					Name:        "get",
					Description: "Get force subscribed servers status",
					Handler:     c.getForceSubscribedServers,
				},
				{
					Name:        "set",
					Description: "Set force subscribed servers",
					Arguments:   "<on/off>",
					Handler:     c.setForceSubscribedServers,
				},
				{
					Name:        "toggle",
					Description: "Toggle force subscribed servers",
					Handler:     c.toggleForceSubscribedServers,
				},
			},
		},
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
		{
			Name:        "notifications",
			Description: "Notification settings commands",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			SubCommands: []*router.Command{
				{
					Name:        "get",
					Description: "Get notifications status",
					Handler:     c.getNotifications,
				},
				{
					Name:        "set",
					Description: "Set notifications status",
					Arguments:   "<on/off>",
					Handler:     c.setNotifications,
				},
				{
					Name:        "toggle",
					Description: "Toggle notifications",
					Handler:     c.toggleNotifications,
				},
			},
		},
		{
			Name:        "nitrosniper",
			Description: "Nitro sniper commands",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			SubCommands: []*router.Command{
				{
					Name:        "get",
					Description: "Get nitro sniper status",
					Handler:     c.getNitroSniper,
				},
				{
					Name:        "set",
					Description: "Set nitro sniper status",
					Arguments:   "<on/off>",
					Handler:     c.setNitroSniper,
				},
				{
					Name:        "toggle",
					Description: "Toggle nitro sniper",
					Handler:     c.toggleNitroSniper,
				},
			},
		},
		{
			Name:        "giveawaysniper",
			Description: "Giveaway sniper commands",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			SubCommands: []*router.Command{
				{
					Name:        "get",
					Description: "Get giveaway sniper status",
					Handler:     c.getGiveawaySniper,
				},
				{
					Name:        "set",
					Description: "Set giveaway sniper status",
					Arguments:   "<on/off>",
					Handler:     c.setGiveawaySniper,
				},
				{
					Name:        "toggle",
					Description: "Toggle giveaway sniper",
					Handler:     c.toggleGiveawaySniper,
				},
			},
		},
		{
			Name:        "giveawaysnipersleep",
			Description: "Giveaway sniper sleep commands",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			SubCommands: []*router.Command{
				{
					Name:        "get",
					Description: "Get giveaway sniper sleep duration",
					Handler:     c.getGiveawaySnipeDelay,
				},
				{
					Name:        "set",
					Description: "Set giveaway sniper sleep duration",
					Arguments:   "<seconds>",
					Handler:     c.setGiveawaySnipeDelay,
				},
			},
		},
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
		{
			Name:        "privnotesniper",
			Description: "Private note sniper commands",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			SubCommands: []*router.Command{
				{
					Name:        "get",
					Description: "Get private note sniper status",
					Handler:     c.getPrivnoteSniper,
				},
				{
					Name:        "set",
					Description: "Set private note sniper status",
					Arguments:   "<on/off>",
					Handler:     c.setPrivnoteSniper,
				},
				{
					Name:        "toggle",
					Description: "Toggle private note sniper",
					Handler:     c.togglePrivnoteSniper,
				},
			},
		},
	})
}

func (c *ConfigCommands) getCommandPrefix(ctx *router.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		ctx.RespondText("Error getting config")
		return
	}

	c.services.Embed.NewEmbed().
		AddField("Command Prefix", cfg.Selfbot.CommandPrefix).
		Send(ctx)
}

func (c *ConfigCommands) setCommandPrefix(ctx *router.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		ctx.RespondText("Error getting config")
		return
	}

	prefix := ctx.Arguments.Raw()

	cfg.Selfbot.CommandPrefix = prefix
	_ = cfg.Save()

	c.services.Embed.NewEmbed().
		SetDescription("Command prefix set to `%s`", prefix).
		Send(ctx)
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

func (c *ConfigCommands) getNotifications(ctx *router.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		ctx.RespondText("Error getting config")
		return
	}

	c.services.Embed.NewEmbed().
		AddField("Notifications", utils.Misc.FormatBoolOnOff(cfg.Notifications.Enabled)).
		Send(ctx)
}

func (c *ConfigCommands) setNotifications(ctx *router.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		ctx.RespondText("Error getting config")
		return
	}

	toggle, err := ctx.Arguments.Get(0).Bool()
	if err != nil {
		ctx.RespondText("Invalid value. Use `on` or `off`.")
		return
	}

	cfg.Notifications.Enabled = toggle
	_ = cfg.Save()

	c.services.Embed.NewEmbed().
		SetDescription("Turned " + utils.Misc.FormatBoolOnOff(toggle) + " notifications").
		Send(ctx)
}

func (c *ConfigCommands) toggleNotifications(ctx *router.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		ctx.RespondText("Error getting config")
		return
	}

	cfg.Notifications.Enabled = !cfg.Notifications.Enabled
	_ = cfg.Save()

	c.services.Embed.NewEmbed().
		SetDescription("Turned " + utils.Misc.FormatBoolOnOff(cfg.Notifications.Enabled) + " notifications").
		Send(ctx)
}

func (c *ConfigCommands) getNitroSniper(ctx *router.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		ctx.RespondText("Error getting config")
		return
	}

	c.services.Embed.NewEmbed().
		AddField("Nitro Sniper", utils.Misc.FormatBoolOnOff(cfg.Snipers.NitroSniper)).
		Send(ctx)
}

func (c *ConfigCommands) setNitroSniper(ctx *router.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		ctx.RespondText("Error getting config")
		return
	}

	toggle, err := ctx.Arguments.Get(0).Bool()
	if err != nil {
		ctx.RespondText("Invalid value. Use `on` or `off`.")
		return
	}

	cfg.Snipers.NitroSniper = toggle
	_ = cfg.Save()

	c.services.Embed.NewEmbed().
		SetDescription("Turned " + utils.Misc.FormatBoolOnOff(toggle) + " nitro sniper").
		Send(ctx)
}

func (c *ConfigCommands) toggleNitroSniper(ctx *router.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		ctx.RespondText("Error getting config")
		return
	}

	cfg.Snipers.NitroSniper = !cfg.Snipers.NitroSniper
	_ = cfg.Save()

	c.services.Embed.NewEmbed().
		SetDescription("Turned " + utils.Misc.FormatBoolOnOff(cfg.Snipers.NitroSniper) + " nitro sniper").
		Send(ctx)
}

func (c *ConfigCommands) getGiveawaySniper(ctx *router.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		ctx.RespondText("Error getting config")
		return
	}

	c.services.Embed.NewEmbed().
		AddField("Giveaway Sniper", utils.Misc.FormatBoolOnOff(cfg.Snipers.GiveawaySniper.Enabled)).
		Send(ctx)
}

func (c *ConfigCommands) setGiveawaySniper(ctx *router.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		ctx.RespondText("Error getting config")
		return
	}

	toggle, err := ctx.Arguments.Get(0).Bool()
	if err != nil {
		ctx.RespondText("Invalid value. Use `on` or `off`.")
		return
	}

	cfg.Snipers.GiveawaySniper.Enabled = toggle
	_ = cfg.Save()

	c.services.Embed.NewEmbed().
		SetDescription("Turned " + utils.Misc.FormatBoolOnOff(toggle) + " giveaway sniper").
		Send(ctx)
}

func (c *ConfigCommands) toggleGiveawaySniper(ctx *router.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		ctx.RespondText("Error getting config")
		return
	}

	cfg.Snipers.GiveawaySniper.Enabled = !cfg.Snipers.GiveawaySniper.Enabled
	_ = cfg.Save()

	c.services.Embed.NewEmbed().
		SetDescription("Turned " + utils.Misc.FormatBoolOnOff(cfg.Snipers.GiveawaySniper.Enabled) + " giveaway sniper").
		Send(ctx)
}

func (c *ConfigCommands) getGiveawaySnipeDelay(ctx *router.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		ctx.RespondText("Error getting config")
		return
	}

	c.services.Embed.NewEmbed().
		AddField("Giveaway Sniper Sleep", fmt.Sprintf("%d seconds", cfg.Snipers.GiveawaySniper.Sleep)).
		Send(ctx)
}

func (c *ConfigCommands) setGiveawaySnipeDelay(ctx *router.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		ctx.RespondText("Error getting config")
		return
	}

	delay, err := ctx.Arguments.Get(0).Int()
	if err != nil {
		ctx.RespondText("Invalid delay value. Must be a number in seconds.")
		return
	}

	cfg.Snipers.GiveawaySniper.Sleep = delay
	_ = cfg.Save()

	c.services.Embed.NewEmbed().
		SetDescription("Set giveaway sniper sleep to %d seconds", delay).
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

func (c *ConfigCommands) getPrivnoteSniper(ctx *router.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		ctx.RespondText("Error getting config")
		return
	}

	c.services.Embed.NewEmbed().
		AddField("Private Note Sniper", utils.Misc.FormatBoolOnOff(cfg.Snipers.PrivnoteSniper)).
		Send(ctx)
}

func (c *ConfigCommands) setPrivnoteSniper(ctx *router.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		ctx.RespondText("Error getting config")
		return
	}

	toggle, err := ctx.Arguments.Get(0).Bool()
	if err != nil {
		ctx.RespondText("Invalid value. Use `on` or `off`.")
		return
	}

	cfg.Snipers.PrivnoteSniper = toggle
	_ = cfg.Save()

	c.services.Embed.NewEmbed().
		SetDescription("Turned " + utils.Misc.FormatBoolOnOff(toggle) + " private note sniper").
		Send(ctx)
}

func (c *ConfigCommands) togglePrivnoteSniper(ctx *router.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		ctx.RespondText("Error getting config")
		return
	}

	cfg.Snipers.PrivnoteSniper = !cfg.Snipers.PrivnoteSniper
	_ = cfg.Save()

	c.services.Embed.NewEmbed().
		SetDescription("Turned " + utils.Misc.FormatBoolOnOff(cfg.Snipers.PrivnoteSniper) + " private note sniper").
		Send(ctx)
}

func (c *ConfigCommands) getCommandMode(ctx *router.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		ctx.RespondText("Error getting config")
		return
	}

	c.services.Embed.NewEmbed().
		AddField("Command Mode", cfg.Selfbot.CommandMode).
		Send(ctx)
}

func (c *ConfigCommands) setCommandMode(ctx *router.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		ctx.RespondText("Error getting config")
		return
	}

	mode := ctx.Arguments.Raw()
	if mode != "prefix" && mode != "mention" {
		ctx.RespondText("Invalid mode. Use `prefix` or `mention`.")
		return
	}

	cfg.Selfbot.CommandMode = mode
	_ = cfg.Save()

	c.services.Embed.NewEmbed().
		SetDescription("Command mode set to `%s`", mode).
		Send(ctx)
}

func (c *ConfigCommands) listCommandModes(ctx *router.Ctx) {
	c.services.Embed.NewEmbed().
		SetDescription("Available command modes: `prefix`, `mention`").
		Send(ctx)
}

func (c *ConfigCommands) getTitle(ctx *router.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		ctx.RespondText("Error getting config")
		return
	}

	c.services.Embed.NewEmbed().
		AddField("Title", cfg.Selfbot.Title).
		Send(ctx)
}

func (c *ConfigCommands) setTitle(ctx *router.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		ctx.RespondText("Error getting config")
		return
	}

	title := ctx.Arguments.Raw()
	cfg.Selfbot.Title = title
	_ = cfg.Save()

	c.services.Embed.NewEmbed().
		SetDescription("Title set to `%s`", title).
		Send(ctx)
}

func (c *ConfigCommands) getFooter(ctx *router.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		ctx.RespondText("Error getting config")
		return
	}

	c.services.Embed.NewEmbed().
		AddField("Footer", cfg.Selfbot.Footer).
		Send(ctx)
}

func (c *ConfigCommands) setFooter(ctx *router.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		ctx.RespondText("Error getting config")
		return
	}

	footer := ctx.Arguments.Raw()
	cfg.Selfbot.Footer = footer
	_ = cfg.Save()

	c.services.Embed.NewEmbed().
		SetDescription("Footer set to `%s`", footer).
		Send(ctx)
}

func (c *ConfigCommands) getDelResponses(ctx *router.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		ctx.RespondText("Error getting config")
		return
	}

	c.services.Embed.NewEmbed().
		AddField("Delete Responses", utils.Misc.FormatBoolOnOff(cfg.Selfbot.DeleteResponses)).
		Send(ctx)
}

func (c *ConfigCommands) setDelResponses(ctx *router.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		ctx.RespondText("Error getting config")
		return
	}

	toggle, err := ctx.Arguments.Get(0).Bool()
	if err != nil {
		ctx.RespondText("Invalid value. Use `on` or `off`.")
		return
	}

	cfg.Selfbot.DeleteResponses = toggle
	_ = cfg.Save()

	c.services.Embed.NewEmbed().
		SetDescription("Turned " + utils.Misc.FormatBoolOnOff(toggle) + " delete responses").
		Send(ctx)
}

func (c *ConfigCommands) toggleDelResponses(ctx *router.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		ctx.RespondText("Error getting config")
		return
	}

	cfg.Selfbot.DeleteResponses = !cfg.Selfbot.DeleteResponses
	_ = cfg.Save()

	c.services.Embed.NewEmbed().
		SetDescription("Turned " + utils.Misc.FormatBoolOnOff(cfg.Selfbot.DeleteResponses) + " delete responses").
		Send(ctx)
}

func (c *ConfigCommands) getDelResponsesDelay(ctx *router.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		ctx.RespondText("Error getting config")
		return
	}

	c.services.Embed.NewEmbed().
		AddField("Delete Responses Delay", fmt.Sprintf("%d seconds", cfg.Selfbot.DeleteResponsesDelay)).
		Send(ctx)
}

func (c *ConfigCommands) setDelResponsesDelay(ctx *router.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		ctx.RespondText("Error getting config")
		return
	}

	delay, err := ctx.Arguments.Get(0).Int()
	if err != nil {
		ctx.RespondText("Invalid delay value. Must be a number in seconds.")
		return
	}

	cfg.Selfbot.DeleteResponsesDelay = delay
	_ = cfg.Save()

	c.services.Embed.NewEmbed().
		SetDescription("Set delete responses delay to %d seconds", delay).
		Send(ctx)
}

func (c *ConfigCommands) getDelCommands(ctx *router.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		ctx.RespondText("Error getting config")
		return
	}

	c.services.Embed.NewEmbed().
		AddField("Delete Commands", utils.Misc.FormatBoolOnOff(cfg.Selfbot.DeleteCommands)).
		Send(ctx)
}

func (c *ConfigCommands) setDelCommands(ctx *router.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		ctx.RespondText("Error getting config")
		return
	}

	toggle, err := ctx.Arguments.Get(0).Bool()
	if err != nil {
		ctx.RespondText("Invalid value. Use `on` or `off`.")
		return
	}

	cfg.Selfbot.DeleteCommands = toggle
	_ = cfg.Save()

	c.services.Embed.NewEmbed().
		SetDescription("Turned " + utils.Misc.FormatBoolOnOff(toggle) + " delete commands").
		Send(ctx)
}

func (c *ConfigCommands) toggleDelCommands(ctx *router.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		ctx.RespondText("Error getting config")
		return
	}

	cfg.Selfbot.DeleteCommands = !cfg.Selfbot.DeleteCommands
	_ = cfg.Save()

	c.services.Embed.NewEmbed().
		SetDescription("Turned " + utils.Misc.FormatBoolOnOff(cfg.Selfbot.DeleteCommands) + " delete commands").
		Send(ctx)
}

func (c *ConfigCommands) getClientSpoof(ctx *router.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		ctx.RespondText("Error getting config")
		return
	}

	c.services.Embed.NewEmbed().
		AddField("Client", cfg.Selfbot.Client).
		Send(ctx)
}

func (c *ConfigCommands) setClientSpoof(ctx *router.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		ctx.RespondText("Error getting config")
		return
	}

	client := strings.Title(strings.ToLower(ctx.Arguments.Raw()))
	switch client {
	case "Mobile", "Web", "Desktop", "Console":
		cfg.Selfbot.Client = client
	default:
		ctx.RespondText("Invalid client. Use `Mobile`, `Web`, `Desktop` or `Console`.")
		return
	}

	_ = cfg.Save()

	c.services.Embed.NewEmbed().
		SetDescription("Client set to `%s`", client).
		Send(ctx)
}

func (c *ConfigCommands) toggleClientSpoof(ctx *router.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		ctx.RespondText("Error getting config")
		return
	}

	switch cfg.Selfbot.Client {
	case "Mobile":
		cfg.Selfbot.Client = "Web"
	case "Web":
		cfg.Selfbot.Client = "Desktop"
	case "Desktop":
		cfg.Selfbot.Client = "Console"
	case "Console", "":
		cfg.Selfbot.Client = "Mobile"
	}

	_ = cfg.Save()

	c.services.Embed.NewEmbed().
		SetDescription("Client set to `%s`", cfg.Selfbot.Client).
		Send(ctx)
}

func (c *ConfigCommands) getTimezone(ctx *router.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		ctx.RespondText("Error getting config")
		return
	}

	c.services.Embed.NewEmbed().
		AddField("Timezone", cfg.Selfbot.Timezone).
		Send(ctx)
}

func (c *ConfigCommands) setTimezone(ctx *router.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		ctx.RespondText("Error getting config")
		return
	}

	timezone := ctx.Arguments.Raw()
	_, err = time.LoadLocation(timezone)
	if err != nil {
		ctx.RespondText("Invalid timezone. Use a valid IANA time zone name.")
		return
	}

	cfg.Selfbot.Timezone = timezone
	_ = cfg.Save()

	c.services.Embed.NewEmbed().
		SetDescription("Timezone set to `%s`", timezone).
		Send(ctx)
}

func (c *ConfigCommands) getEmbedColor(ctx *router.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		ctx.RespondText("Error getting config")
		return
	}

	c.services.Embed.NewEmbed().
		AddField("Embed Color", cfg.Selfbot.EmbedColor).
		Send(ctx)
}

func (c *ConfigCommands) setEmbedColor(ctx *router.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		ctx.RespondText("Error getting config")
		return
	}

	color := ctx.Arguments.Raw()
	cfg.Selfbot.EmbedColor = color
	_ = cfg.Save()

	c.services.Embed.NewEmbed().
		SetDescription("Embed color set to `%s`", color).
		Send(ctx)
}

func (c *ConfigCommands) getEmbedImage(ctx *router.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		ctx.RespondText("Error getting config")
		return
	}

	c.services.Embed.NewEmbed().
		AddField("Embed Image", cfg.Selfbot.EmbedImage).
		Send(ctx)
}

func (c *ConfigCommands) setEmbedImage(ctx *router.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		ctx.RespondText("Error getting config")
		return
	}

	url := ctx.Arguments.Raw()
	cfg.Selfbot.EmbedImage = url
	_ = cfg.Save()

	c.services.Embed.NewEmbed().
		SetDescription("Embed image set to `%s`", url).
		Send(ctx)
}

func (c *ConfigCommands) getEmbedURL(ctx *router.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		ctx.RespondText("Error getting config")
		return
	}

	c.services.Embed.NewEmbed().
		AddField("Embed URL", cfg.Selfbot.EmbedURL).
		Send(ctx)
}

func (c *ConfigCommands) setEmbedURL(ctx *router.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		ctx.RespondText("Error getting config")
		return
	}

	url := ctx.Arguments.Raw()
	cfg.Selfbot.EmbedURL = url
	_ = cfg.Save()

	c.services.Embed.NewEmbed().
		SetDescription("Embed URL set to `%s`", url).
		Send(ctx)
}

func (c *ConfigCommands) getForceSubscribedServers(ctx *router.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		ctx.RespondText("Error getting config")
		return
	}

	c.services.Embed.NewEmbed().
		AddField("Force Subscribed Servers", utils.Misc.FormatBoolOnOff(len(cfg.Selfbot.ForceSubscribedServers) > 0)).
		Send(ctx)
}

func (c *ConfigCommands) setForceSubscribedServers(ctx *router.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		ctx.RespondText("Error getting config")
		return
	}

	toggle, err := ctx.Arguments.Get(0).Bool()
	if err != nil {
		ctx.RespondText("Invalid value. Use `on` or `off`.")
		return
	}

	if toggle {
		guildIDs := make([]string, 0)
		for _, guild := range ctx.Session.State.Guilds {
			guildIDs = append(guildIDs, guild.ID)
		}
		cfg.Selfbot.ForceSubscribedServers = guildIDs
	} else {
		cfg.Selfbot.ForceSubscribedServers = []string{}
	}

	_ = cfg.Save()

	c.services.Embed.NewEmbed().
		SetDescription("Turned " + utils.Misc.FormatBoolOnOff(toggle) + " force subscribed servers").
		Send(ctx)
}

func (c *ConfigCommands) toggleForceSubscribedServers(ctx *router.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		ctx.RespondText("Error getting config")
		return
	}

	if len(cfg.Selfbot.ForceSubscribedServers) > 0 {
		cfg.Selfbot.ForceSubscribedServers = []string{}
	} else {
		guildIDs := make([]string, 0)
		for _, guild := range ctx.Session.State.Guilds {
			guildIDs = append(guildIDs, guild.ID)
		}
		cfg.Selfbot.ForceSubscribedServers = guildIDs
	}

	_ = cfg.Save()

	c.services.Embed.NewEmbed().
		SetDescription("Turned " + utils.Misc.FormatBoolOnOff(len(cfg.Selfbot.ForceSubscribedServers) > 0) + " force subscribed servers").
		Send(ctx)
}
