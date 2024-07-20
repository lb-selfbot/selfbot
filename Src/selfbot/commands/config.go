package commands

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/LightningDev1/LB-Selfbot/selfbot/config"
	"github.com/LightningDev1/LB-Selfbot/selfbot/embed"
	"github.com/LightningDev1/LB-Selfbot/selfbot/eventreactions"
	"github.com/LightningDev1/LB-Selfbot/selfbot/logging"
	"github.com/LightningDev1/LB-Selfbot/selfbot/utils"
	"github.com/LightningDev1/dgc"
	"github.com/LightningDev1/discordgo"
)

type ConfigCommands struct{}

func (c ConfigCommands) RegisterCommands(router *dgc.Router) {
	router.StartCategory("Config", "Change the bot's config")

	router.RegisterCmd(&dgc.Command{
		Name:        "commandprefix",
		Description: "Get or set the command prefix",
		Usage:       "[p]commandprefix <get/set> [value]",
		Flags:       []string{CmdFlagDangerous},
		Handler:     c.commandprefix,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "commandmode",
		Description: "Get or set the command mode",
		Usage:       "[p]commandmode <get/set/values> [value]",
		Flags:       []string{CmdFlagDangerous},
		Handler:     c.commandmode,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "title",
		Description: "Get or set the title",
		Usage:       "[p]title <get/set> [value]",
		Flags:       []string{CmdFlagDangerous},
		Handler:     c.title,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "footer",
		Description: "Get or set the footer",
		Usage:       "[p]footer <get/set> [value]",
		Flags:       []string{CmdFlagDangerous},
		Handler:     c.footer,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "delresponses",
		Description: "Get, set or toggle if the bot should delete responses to commands",
		Usage:       "[p]delresponses <get/set/toggle> [on/off]",
		Flags:       []string{CmdFlagDangerous},
		Handler:     c.delresponses,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "delresponsesdelay",
		Description: "Get or set the delay before deleting responses",
		Usage:       "[p]delresponsesdelay <get/set> [value]",
		Flags:       []string{CmdFlagDangerous},
		Handler:     c.delresponsesdelay,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "delcommands",
		Description: "Get, set or toggle if the bot should delete commands",
		Usage:       "[p]delcommands <get/set/toggle> [on/off]",
		Flags:       []string{CmdFlagDangerous},
		Handler:     c.delcommands,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "clientspoof",
		Description: "Get or set the bot's client spoof",
		Usage:       "[p]clientspoof <get/set/values> [value]",
		Flags:       []string{CmdFlagDangerous},
		Handler:     c.clientspoof,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "timezone",
		Description: "Get or set the timezone",
		Usage:       "[p]timezone <get/set> [value]",
		Flags:       []string{CmdFlagDangerous},
		Handler:     c.timezone,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "embedcolor",
		Description: "Get or set the Embed Color (for URL Embeds)",
		Usage:       "[p]embedcolor <get/set> [hex color]",
		Flags:       []string{CmdFlagDangerous},
		Handler:     c.embedcolor,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "embedimage",
		Description: "Get or set the Embed Image (for URL Embeds)",
		Usage:       "[p]embedimage <get/set> [url]",
		Flags:       []string{CmdFlagDangerous},
		Handler:     c.embedimage,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "embedurl",
		Description: "Get or set the Embed URL (for URL Embeds)",
		Usage:       "[p]embedurl <get/set> [url]",
		Flags:       []string{CmdFlagDangerous},
		Handler:     c.embedurl,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "aliases",
		Description: "Get, add, remove or clear command aliases",
		Usage:       "[p]aliases <get/add/rm/clear> [alias] [command]",
		Flags:       []string{CmdFlagDangerous},
		Handler:     c.aliases,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "notifications",
		Description: "Get, set, toggle, or get an example of toast notifications",
		Usage:       "[p]notifications <get/set/toggle/example> [on/off]",
		Flags:       []string{CmdFlagDangerous},
		Handler:     c.notifications,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "nitrosniper",
		Description: "Get, set or toggle the Nitro Sniper",
		Usage:       "[p]nitrosniper <get/set/toggle> [on/off]",
		Flags:       []string{CmdFlagDangerous},
		Handler:     c.nitrosniper,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "giveawaysniper",
		Description: "Get, set or toggle the Giveaway Sniper",
		Usage:       "[p]giveawaysniper <get/set/toggle> [on/off]",
		Flags:       []string{CmdFlagDangerous},
		Handler:     c.giveawaysniper,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "giveawaysnipersleep",
		Description: "Get or set the Giveaway Sniper sleep delay",
		Usage:       "[p]giveawaysnipersleep <get/set> [value]",
		Flags:       []string{CmdFlagDangerous},
		Handler:     c.giveawaysnipersleep,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "customgiveawaybots",
		Description: "Get, add, remove or clear custom giveaway bots",
		Usage:       "[p]customgiveawaybots <get/add/rm/clear> [bot]",
		Flags:       []string{CmdFlagDangerous},
		Handler:     c.customgiveawaybots,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "sniperblacklist",
		Description: "Get, add, remove or clear the sniper blacklist",
		Usage:       "[p]sniperblacklist <get/add/rm/clear> [value]",
		Flags:       []string{CmdFlagDangerous},
		Handler:     c.sniperblacklist,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "sniperserverblacklist",
		Description: "Get, add, remove or clear the sniper server blacklist",
		Usage:       "[p]sniperserverblacklist <get/add/rm/clear> [server id]",
		Flags:       []string{CmdFlagDangerous},
		Handler:     c.sniperserverblacklist,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "multitokensniper",
		Description: "Get, set or toggle the Multi Token Sniper",
		Usage:       "[p]multitokensniper <get/set/toggle> [on/off]",
		Flags:       []string{CmdFlagDangerous},
		Handler:     c.multitokensniper,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "snipertokens",
		Description: "Get, add, remove or clear Multi Token Sniper tokens",
		Usage:       "[p]snipertokens <get/add/rm/clear> [value]",
		Flags:       []string{CmdFlagDangerous},
		Handler:     c.snipertokens,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "privnotesniper",
		Description: "Get, set or toggle the Privnote Sniper",
		Usage:       "[p]privnotesniper <get/set/toggle> [on/off]",
		Flags:       []string{CmdFlagDangerous},
		Handler:     c.privnotesniper,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "poketwosniper",
		Description: "Get, set or toggle the PokéTwo Sniper",
		Usage:       "[p]poketwosniper <get/set/toggle> [on/off]",
		Flags:       []string{CmdFlagDangerous},
		Handler:     c.poketwosniper,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "poketwosleep",
		Description: "Get or set the PokéTwo Sniper catch delay",
		Usage:       "[p]poketwosleep <get/set> [value]",
		Flags:       []string{CmdFlagDangerous},
		Handler:     c.poketwosleep,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "poketworandomsleep",
		Description: "Get, set or toggle the PokéTwo Sniper random sleep",
		Usage:       "[p]poketworandomsleep <get/set/toggle> [on/off]",
		Flags:       []string{CmdFlagDangerous},
		Handler:     c.poketworandomsleep,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "poketwotimeout",
		Description: "Get or set the PokéTwo Sniper timeout",
		Usage:       "[p]poketwotimeout <get/set> [value]",
		Flags:       []string{CmdFlagDangerous},
		Handler:     c.poketwotimeout,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "poketwotyping",
		Description: "Get, set or toggle the PokéTwo typing",
		Usage:       "[p]poketwotyping <get/set/toggle> [on/off]",
		Flags:       []string{CmdFlagDangerous},
		Handler:     c.poketwotyping,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "poketwocatchrate",
		Description: "Get or set the PokéTwo catch rate",
		Usage:       "[p]poketwocatchrate <get/set> [percentage]",
		Flags:       []string{CmdFlagDangerous},
		Handler:     c.poketwocatchrate,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "poketwofakeactivity",
		Description: "Get, set or toggle the PokéTwo fake activity",
		Usage:       "[p]poketwofakeactivity <get/set/toggle> [on/off]",
		Flags:       []string{CmdFlagDangerous},
		Handler:     c.poketwofakeactivity,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "poketwochannels",
		Description: "Get, add, remove or clear PokéTwo channels",
		Usage:       "[p]poketwochannels <get/add/rm/clear> [channel]",
		Flags:       []string{CmdFlagDangerous},
		Handler:     c.poketwochannels,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "poketwolanguage",
		Description: "Get or set the PokéTwo Sniper catch language",
		Usage:       "[p]poketwolanguage <get/set/values> [language]",
		Flags:       []string{CmdFlagDangerous},
		Handler:     c.poketwolanguage,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "setupwebhooks",
		Description: "Automatically set up all webhooks",
		Usage:       "[p]setupwebhooks [server id]",
		Flags:       []string{CmdFlagDangerous},
		Handler:     c.setupwebhooks,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "wall",
		Description: "Set all webhooks at once",
		Usage:       "[p]wall [webhook url]",
		Flags:       []string{CmdFlagDangerous},
		Handler:     c.wall,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "wnitrosniper",
		Description: "Get or set the Nitro Sniper webhook",
		Usage:       "[p]wnitrosniper <get/set> [webhook url]",
		Flags:       []string{CmdFlagDangerous},
		Handler:     c.wnitrosniper,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "wgiveawaysniper",
		Description: "Get or set the Giveaway Sniper webhook",
		Usage:       "[p]wgiveawaysniper <get/set> [webhook url]",
		Flags:       []string{CmdFlagDangerous},
		Handler:     c.wgiveawaysniper,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "wselfbotdetection",
		Description: "Get or set the Selfbot Detection webhook",
		Usage:       "[p]wselfbotdetection <get/set> [webhook url]",
		Flags:       []string{CmdFlagDangerous},
		Handler:     c.wselfbotdetection,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "wstaffdetection",
		Description: "Get or set the Staff Detection webhook",
		Usage:       "[p]wstaffdetection <get/set> [webhook url]",
		Flags:       []string{CmdFlagDangerous},
		Handler:     c.wstaffdetection,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "wnotifiers",
		Description: "Get or set the Notifiers webhook",
		Usage:       "[p]wnotifiers <get/set> [webhook url]",
		Flags:       []string{CmdFlagDangerous},
		Handler:     c.wnotifiers,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "wpingdetection",
		Description: "Get or set the Ping Detection webhook",
		Usage:       "[p]wpingdetection <get/set> [webhook url]",
		Flags:       []string{CmdFlagDangerous},
		Handler:     c.wpingdetection,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "wmisc",
		Description: "Get or set the Miscellaneous webhook",
		Usage:       "[p]wmisc <get/set> [webhook url]",
		Flags:       []string{CmdFlagDangerous},
		Handler:     c.wmisc,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "wsessions",
		Description: "Get or set the Sessions webhook",
		Usage:       "[p]wsessions <get/set> [webhook url]",
		Flags:       []string{CmdFlagDangerous},
		Handler:     c.wsessions,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "wpingself",
		Description: "Get, set or toggle Ping Self for webhooks",
		Usage:       "[p]wpingself <get/set/values> [value]",
		Flags:       []string{CmdFlagDangerous},
		Handler:     c.wpingself,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "wpingeveryone",
		Description: "Get, set or toggle Ping Everyone for webhooks",
		Usage:       "[p]wpingeveryone <get/set/values> [value]",
		Flags:       []string{CmdFlagDangerous},
		Handler:     c.wpingeveryone,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "bannotifier",
		Description: "Get, set or toggle the Ban Notifier",
		Usage:       "[p]bannotifier <get/set/toggle> [on/off]",
		Flags:       []string{CmdFlagDangerous},
		Handler:     c.bannotifier,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "staffdetection",
		Description: "Get, set or toggle the Staff Detection",
		Usage:       "[p]staffdetection <get/set/toggle> [on/off]",
		Flags:       []string{CmdFlagDangerous},
		Handler:     c.staffdetection,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "pingdetection",
		Description: "Get, set or toggle the Ping Detection",
		Usage:       "[p]pingdetection <get/set/toggle> [on/off]",
		Flags:       []string{CmdFlagDangerous},
		Handler:     c.pingdetection,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "ticketnotifier",
		Description: "Get, set or toggle the Ticket Notifier",
		Usage:       "[p]ticketnotifier <get/set/toggle> [on/off]",
		Flags:       []string{CmdFlagDangerous},
		Handler:     c.ticketnotifier,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "sessionconnected",
		Description: "Get, set or toggle the Session Connected Notifier",
		Usage:       "[p]sessionconnected <get/set/toggle> [on/off]",
		Flags:       []string{CmdFlagDangerous},
		Handler:     c.sessionconnected,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "sessiondisconnected",
		Description: "Get, set or toggle the Session Disconnected Notifier",
		Usage:       "[p]sessiondisconnected <get/set/toggle> [on/off]",
		Flags:       []string{CmdFlagDangerous},
		Handler:     c.sessiondisconnected,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "selfbotdetection",
		Description: "Get, set or toggle the Selfbot Detection",
		Usage:       "[p]selfbotdetection <get/set/toggle> [on/off]",
		Flags:       []string{CmdFlagDangerous},
		Handler:     c.selfbotdetection,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "dmtypingnotifier",
		Description: "Get, set or toggle the DM Typing Notifier",
		Usage:       "[p]dmtypingnotifier <get/set/toggle> [on/off]",
		Flags:       []string{CmdFlagDangerous},
		Handler:     c.dmtypingnotifier,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "dmnotifier",
		Description: "Get, set or toggle the DM Notifier",
		Usage:       "[p]dmnotifier <get/set/toggle> [on/off]",
		Flags:       []string{CmdFlagDangerous},
		Handler:     c.dmnotifier,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "dmeditnotifier",
		Description: "Get, set or toggle the DM Edit Notifier",
		Usage:       "[p]dmeditnotifier <get/set/toggle> [on/off]",
		Flags:       []string{CmdFlagDangerous},
		Handler:     c.dmeditnotifier,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "dmdeletenotifier",
		Description: "Get, set or toggle the DM Delete Notifier",
		Usage:       "[p]dmdeletenotifier <get/set/toggle> [on/off]",
		Flags:       []string{CmdFlagDangerous},
		Handler:     c.dmdeletenotifier,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "messagedeletenotifier",
		Description: "Get, set or toggle the Message Delete Notifier",
		Usage:       "[p]messagedeletenotifier <get/set/toggle> [on/off]",
		Flags:       []string{CmdFlagDangerous},
		Handler:     c.messagedeletenotifier,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "uitheme",
		Description: "Get or set the UI theme",
		Usage:       "[p]uitheme <get/set/values> [value]",
		Flags:       []string{CmdFlagDangerous},
		Handler:     c.uitheme,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "hidetosystray",
		Description: "Get, set or toggle Hide UI To Systray",
		Usage:       "[p]hidetosystray <get/set/toggle> [on/off]",
		Flags:       []string{CmdFlagDangerous},
		Handler:     c.hidetosystray,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "runuionstartup",
		Description: "Get, set or toggle Run UI on Startup",
		Usage:       "[p]runuionstartup <get/set/toggle> [on/off]",
		Flags:       []string{CmdFlagDangerous},
		Handler:     c.runuionstartup,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "runbotonstartup",
		Description: "Get, set or toggle Run Selfbot on Startup",
		Usage:       "[p]runbotonstartup <get/set/toggle> [on/off]",
		Flags:       []string{CmdFlagDangerous},
		Handler:     c.runbotonstartup,
	})
}

func (ConfigCommands) commandprefix(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	switch utils.Arguments["<get/set>"] {
	case "get":
		embed.NewEmbed().
			AddField("Command Prefix", cfg.Selfbot.CommandPrefix).
			Send(ctx)
	case "set":
		value := utils.Arguments["[value]"]
		if value == "" {
			_ = ctx.RespondText("You must provide a value to set!")
			return
		}

		cfg.Selfbot.CommandPrefix = value
		_ = cfg.Save()

		embed.NewEmbed().
			SetDescription("Command prefix set to " + value).
			Send(ctx)
	default:
		_ = ctx.RespondText("Invalid <get/set> argument. Use `get` or `set`.")
	}
}

func (ConfigCommands) commandmode(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	switch utils.Arguments["<get/set/values>"] {
	case "get":
		embed.NewEmbed().
			AddField("Command Mode", cfg.Selfbot.CommandMode).
			Send(ctx)
	case "set":
		value := strings.ToLower(utils.Arguments["[value]"])

		switch value {
		case "codeblock":
			cfg.Selfbot.CommandMode = "Codeblock"
		case "text":
			cfg.Selfbot.CommandMode = "Text"
		case "indent":
			cfg.Selfbot.CommandMode = "Indent"
		case "webhookembed":
			cfg.Selfbot.CommandMode = "WebhookEmbed"
		case "urlembed":
			cfg.Selfbot.CommandMode = "URLEmbed"
		default:
			_ = ctx.RespondText("Invalid command mode. Use `[p]commandmode values` to get a list of available command modes.")
			return
		}

		_ = cfg.Save()

		embed.NewEmbed().
			SetDescription("Command mode set to " + cfg.Selfbot.CommandMode).
			Send(ctx)
	case "values":
		embed.NewEmbed().
			SetDescription("Available command modes: Codeblock, Text, Indent, WebhookEmbed, URLEmbed").
			Send(ctx)
	default:
		_ = ctx.RespondText("Invalid <get/set/values> argument. Use `get` or `set` or `values`.")
	}
}

func (ConfigCommands) title(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	switch utils.Arguments["<get/set>"] {
	case "get":
		embed.NewEmbed().
			AddField("Title", cfg.Selfbot.Title).
			Send(ctx)
	case "set":
		value := utils.String.RemoveWords(ctx.Arguments.Raw(), 1)

		cfg.Selfbot.Title = value
		_ = cfg.Save()

		embed.NewEmbed().
			SetDescription("Title set to " + value).
			Send(ctx)
	default:
		_ = ctx.RespondText("Invalid <get/set> argument. Use `get` or `set`.")
	}
}

func (ConfigCommands) footer(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	switch utils.Arguments["<get/set>"] {
	case "get":
		embed.NewEmbed().
			AddField("Footer", cfg.Selfbot.Footer).
			Send(ctx)
	case "set":
		value := utils.String.RemoveWords(ctx.Arguments.Raw(), 1)

		cfg.Selfbot.Footer = value
		_ = cfg.Save()

		embed.NewEmbed().
			SetDescription("Footer set to " + value).
			Send(ctx)
	default:
		_ = ctx.RespondText("Invalid <get/set> argument. Use `get` or `set`.")
	}
}

func (ConfigCommands) delresponses(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	switch utils.Arguments["<get/set/toggle>"] {
	case "get":
		embed.NewEmbed().
			AddField("Delete Responses", utils.Misc.FormatBoolOnOff(cfg.Selfbot.DeleteResponses)).
			Send(ctx)
	case "set":
		toggle := strings.ToLower(utils.Arguments["[on/off]"])
		if toggle != "on" && toggle != "off" {
			_ = ctx.RespondText("Invalid value. Use `on` or `off`.")
			return
		}

		cfg.Selfbot.DeleteResponses = toggle == "on"
		_ = cfg.Save()

		embed.NewEmbed().
			SetDescription("Turned " + toggle + " delete responses").
			Send(ctx)
	case "toggle":
		cfg.Selfbot.DeleteResponses = !cfg.Selfbot.DeleteResponses
		_ = cfg.Save()

		embed.NewEmbed().
			SetDescription("Turned " + utils.Misc.FormatBoolOnOff(cfg.Selfbot.DeleteResponses) + " delete responses").
			Send(ctx)
	default:
		_ = ctx.RespondText("Invalid <get/set/toggle> argument. Use `get` or `set` or `toggle`.")
	}
}

func (ConfigCommands) delresponsesdelay(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	switch utils.Arguments["<get/set>"] {
	case "get":
		formatted := utils.Misc.FormatDuration(time.Duration(cfg.Selfbot.DeleteResponsesDelay) * time.Second)

		embed.NewEmbed().
			AddField("Delete Responses Delay", formatted).
			Send(ctx)
	case "set":
		delay, err := utils.Misc.ParseDuration(utils.Arguments["[value]"])
		if err != nil {
			_ = ctx.RespondText("Invalid duration: " + err.Error())
			return
		}

		cfg.Selfbot.DeleteResponsesDelay = int(delay.Seconds())
		_ = cfg.Save()

		formatted := utils.Misc.FormatDuration(time.Duration(cfg.Selfbot.DeleteResponsesDelay) * time.Second)

		embed.NewEmbed().
			SetDescription(fmt.Sprintf("Delete responses delay set to %s.", formatted)).
			Send(ctx)
	default:
		_ = ctx.RespondText("Invalid <get/set> argument. Use `get` or `set`.")
	}
}

func (ConfigCommands) delcommands(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	switch utils.Arguments["<get/set/toggle>"] {
	case "get":
		embed.NewEmbed().
			AddField("Delete Commands", utils.Misc.FormatBoolOnOff(cfg.Selfbot.DeleteCommands)).
			Send(ctx)
	case "set":
		toggle := strings.ToLower(utils.Arguments["[on/off]"])
		if toggle != "on" && toggle != "off" {
			_ = ctx.RespondText("Invalid value. Use `on` or `off`.")
			return
		}

		cfg.Selfbot.DeleteCommands = toggle == "on"
		_ = cfg.Save()

		embed.NewEmbed().
			SetDescription("Turned " + toggle + " delete commands").
			Send(ctx)
	case "toggle":
		cfg.Selfbot.DeleteCommands = !cfg.Selfbot.DeleteCommands
		_ = cfg.Save()

		embed.NewEmbed().
			SetDescription("Turned " + utils.Misc.FormatBoolOnOff(cfg.Selfbot.DeleteCommands) + " delete commands").
			Send(ctx)
	default:
		_ = ctx.RespondText("Invalid <get/set/toggle> argument. Use `get` or `set` or `toggle`.")
	}
}

func (ConfigCommands) clientspoof(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	switch utils.Arguments["<get/set/values>"] {
	case "get":
		embed.NewEmbed().
			AddField("Client", cfg.Selfbot.Client).
			Send(ctx)
	case "set":
		value := strings.ToLower(utils.Arguments["[value]"])

		switch value {
		case "desktop":
			cfg.Selfbot.Client = "Desktop"
		case "web":
			cfg.Selfbot.Client = "Web"
		case "mobile":
			cfg.Selfbot.Client = "Mobile"
		case "console":
			cfg.Selfbot.Client = "Console"
		default:
			_ = ctx.RespondText("Invalid client. Use `[p]clientspoof values` to get a list of available clients.")
			return
		}

		_ = cfg.Save()

		embed.NewEmbed().
			SetDescription("Client Spoof set to " + cfg.Selfbot.Client).
			Send(ctx)
	case "values":
		embed.NewEmbed().
			SetDescription("Available client spoof values: Desktop, Web, Mobile, Console").
			Send(ctx)
	default:
		_ = ctx.RespondText("Invalid <get/set/values> argument. Use `get` or `set` or `values`.")
	}
}

func (ConfigCommands) timezone(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	switch utils.Arguments["<get/set>"] {
	case "get":
		location, err := time.LoadLocation(cfg.Selfbot.Timezone)
		if err != nil {
			location = time.UTC
		}

		currTime := time.Now().In(location).Format("15:04:05")

		embed.NewEmbed().
			AddField("Timezone", cfg.Selfbot.Timezone).
			AddField("Current Time", currTime).
			Send(ctx)
	case "set":
		timezone := utils.Arguments["[value]"]
		if _, err := time.LoadLocation(timezone); err != nil {
			_ = ctx.RespondText("Invalid timezone: " + err.Error())
			return
		}

		cfg.Selfbot.Timezone = timezone
		_ = cfg.Save()

		embed.NewEmbed().
			SetDescription("Timezone set to " + cfg.Selfbot.Timezone).
			Send(ctx)
	default:
		_ = ctx.RespondText("Invalid <get/set> argument. Use `get` or `set`.")
	}
}

func (ConfigCommands) embedcolor(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	switch utils.Arguments["<get/set>"] {
	case "get":
		embed.NewEmbed().
			AddField("Embed Color", cfg.Selfbot.EmbedColor).
			Send(ctx)
	case "set":
		color := strings.TrimPrefix(utils.Arguments["[hex color]"], "#")
		if len(color) != 6 {
			_ = ctx.RespondText("Invalid hex color. Use a 6-digit hex color.")
			return
		}

		_, _, _, err := utils.Misc.HexToRGB(color)
		if err != nil {
			_ = ctx.RespondText("Invalid hex color: " + err.Error() + ". Use a 6-digit hex color.")
			return
		}

		cfg.Selfbot.EmbedColor = color
		_ = cfg.Save()

		embed.NewEmbed().
			SetDescription("Embed Color set to " + cfg.Selfbot.EmbedColor).
			Send(ctx)
	default:
		_ = ctx.RespondText("Invalid <get/set> argument. Use `get` or `set`.")
	}
}

func (ConfigCommands) embedimage(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	switch utils.Arguments["<get/set>"] {
	case "get":
		embed.NewEmbed().
			AddField("Embed Image", cfg.Selfbot.EmbedImage).
			Send(ctx)
	case "set":
		url := utils.Arguments["[url]"]

		cfg.Selfbot.EmbedImage = url
		_ = cfg.Save()

		embed.NewEmbed().
			SetDescription("Embed Image set to " + cfg.Selfbot.EmbedImage).
			Send(ctx)
	default:
		_ = ctx.RespondText("Invalid <get/set> argument. Use `get` or `set`.")
	}
}

func (ConfigCommands) embedurl(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	switch utils.Arguments["<get/set>"] {
	case "get":
		embed.NewEmbed().
			AddField("Embed URL", cfg.Selfbot.EmbedURL).
			Send(ctx)
	case "set":
		url := utils.Arguments["[url]"]

		cfg.Selfbot.EmbedURL = url
		_ = cfg.Save()

		embed.NewEmbed().
			SetDescription("Embed URL set to " + cfg.Selfbot.EmbedURL).
			Send(ctx)
	default:
		_ = ctx.RespondText("Invalid <get/set> argument. Use `get` or `set`.")
	}
}

func (ConfigCommands) aliases(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config.")
		return
	}

	switch utils.Arguments["<get/add/rm/clear>"] {
	case "get":
		prefix := cfg.Selfbot.CommandPrefix
		var aliases []string

		for alias, command := range cfg.Selfbot.CommandAliases {
			aliases = append(aliases, fmt.Sprintf("%s%s - %s%s", prefix, alias, prefix, command))
		}

		if len(aliases) == 0 {
			aliases = append(aliases, "No aliases set.")
		}

		embed.NewEmbed().
			AddField("Command Aliases", strings.Join(aliases, "\n")).
			Send(ctx)
	case "add":
		alias, cmd := utils.Arguments["[alias]"], utils.Arguments["[command]"]

		if alias == "" || cmd == "" {
			_ = ctx.RespondText("You must provide an alias and a command.")
			return
		}

		if _, ok := cfg.Selfbot.CommandAliases[alias]; ok {
			_ = ctx.RespondText("Alias already exists.")
			return
		}

		if ctx.Router.GetCmd(cmd) == nil {
			_ = ctx.RespondText("Command does not exist.")
			return
		}

		cfg.Selfbot.CommandAliases[alias] = cmd
		_ = cfg.Save()

		_ = eventreactions.LoadAliases(ctx.Router)

		embed.NewEmbed().
			SetDescription(fmt.Sprintf("Added alias `%s` for command `%s`.", alias, cmd)).
			Send(ctx)
	case "rm":
		alias := utils.Arguments["[alias]"]

		cmd, ok := cfg.Selfbot.CommandAliases[alias]
		if !ok {
			_ = ctx.RespondText("Alias doesn't exist.")
			return
		}

		delete(cfg.Selfbot.CommandAliases, alias)
		_ = cfg.Save()

		_ = eventreactions.LoadAliases(ctx.Router)

		embed.NewEmbed().
			SetDescription(fmt.Sprintf("Removed alias `%s` for command `%s`.", alias, cmd)).
			Send(ctx)
	case "clear":
		cfg.Selfbot.CommandAliases = make(map[string]string)
		_ = cfg.Save()

		_ = eventreactions.LoadAliases(ctx.Router)

		embed.NewEmbed().
			SetDescription("Cleared all aliases.").
			Send(ctx)
	default:
		_ = ctx.RespondText("Invalid <get/add/rm/clear> argument. Use `get`, `add`, `rm`, or `clear`.")
	}
}

func (ConfigCommands) notifications(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	switch utils.Arguments["<get/set/toggle/example>"] {
	case "get":
		embed.NewEmbed().
			AddField("Notifications", utils.Misc.FormatBoolOnOff(cfg.Notifications.Enabled)).
			Send(ctx)
	case "set":
		toggle := strings.ToLower(utils.Arguments["[on/off]"])
		if toggle != "on" && toggle != "off" {
			_ = ctx.RespondText("Invalid value. Use `on` or `off`.")
			return
		}

		cfg.Notifications.Enabled = toggle == "on"
		_ = cfg.Save()

		embed.NewEmbed().
			SetDescription("Turned " + toggle + " notifications").
			Send(ctx)
	case "toggle":
		cfg.Notifications.Enabled = !cfg.Notifications.Enabled
		_ = cfg.Save()

		embed.NewEmbed().
			SetDescription("Turned " + utils.Misc.FormatBoolOnOff(cfg.Notifications.Enabled) + " notifications").
			Send(ctx)
	case "example":
		logging.UI.Notification("This is an example notification.")

		embed.NewEmbed().
			SetDescription("Sent an example notification.").
			Send(ctx)

	default:
		_ = ctx.RespondText("Invalid <get/set/toggle/example> argument. Use `get` or `set` or `toggle` or `example`.")
	}
}

func (ConfigCommands) nitrosniper(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	switch utils.Arguments["<get/set/toggle>"] {
	case "get":
		embed.NewEmbed().
			AddField("Nitro Sniper", utils.Misc.FormatBoolOnOff(cfg.Snipers.NitroSniper)).
			Send(ctx)
	case "set":
		toggle := strings.ToLower(utils.Arguments["[on/off]"])
		if toggle != "on" && toggle != "off" {
			_ = ctx.RespondText("Invalid value. Use `on` or `off`.")
			return
		}

		cfg.Snipers.NitroSniper = toggle == "on"
		_ = cfg.Save()

		embed.NewEmbed().
			SetDescription("Turned " + toggle + " nitro sniper").
			Send(ctx)
	case "toggle":
		cfg.Snipers.NitroSniper = !cfg.Snipers.NitroSniper
		_ = cfg.Save()

		embed.NewEmbed().
			SetDescription("Turned " + utils.Misc.FormatBoolOnOff(cfg.Snipers.NitroSniper) + " nitro sniper").
			Send(ctx)
	default:
		_ = ctx.RespondText("Invalid <get/set/toggle> argument. Use `get` or `set` or `toggle`.")
	}
}

func (ConfigCommands) giveawaysniper(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	switch utils.Arguments["<get/set/toggle>"] {
	case "get":
		embed.NewEmbed().
			AddField("Giveaway Sniper", utils.Misc.FormatBoolOnOff(cfg.Snipers.GiveawaySniper.Enabled)).
			Send(ctx)
	case "set":
		toggle := strings.ToLower(utils.Arguments["[on/off]"])
		if toggle != "on" && toggle != "off" {
			_ = ctx.RespondText("Invalid value. Use `on` or `off`.")
			return
		}

		cfg.Snipers.GiveawaySniper.Enabled = toggle == "on"
		_ = cfg.Save()

		embed.NewEmbed().
			SetDescription("Turned " + toggle + " giveaway sniper").
			Send(ctx)
	case "toggle":
		cfg.Snipers.GiveawaySniper.Enabled = !cfg.Snipers.GiveawaySniper.Enabled
		_ = cfg.Save()

		embed.NewEmbed().
			SetDescription("Turned " + utils.Misc.FormatBoolOnOff(cfg.Snipers.GiveawaySniper.Enabled) + " giveaway sniper").
			Send(ctx)
	default:
		_ = ctx.RespondText("Invalid <get/set/toggle> argument. Use `get` or `set` or `toggle`.")
	}
}

func (ConfigCommands) giveawaysnipersleep(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	switch utils.Arguments["<get/set>"] {
	case "get":
		formatted := utils.Misc.FormatDuration(time.Duration(cfg.Snipers.GiveawaySniper.Sleep) * time.Second)

		embed.NewEmbed().
			AddField("Giveaway Sniper Sleep", formatted).
			Send(ctx)
	case "set":
		delay, err := utils.Misc.ParseDuration(utils.Arguments["[value]"])
		if err != nil {
			_ = ctx.RespondText("Invalid duration: " + err.Error())
			return
		}

		cfg.Snipers.GiveawaySniper.Sleep = int(delay.Seconds())
		_ = cfg.Save()

		formatted := utils.Misc.FormatDuration(time.Duration(cfg.Snipers.GiveawaySniper.Sleep) * time.Second)

		embed.NewEmbed().
			SetDescription(fmt.Sprintf("Giveaway sniper sleep set to %s.", formatted)).
			Send(ctx)
	default:
		_ = ctx.RespondText("Invalid <get/set> argument. Use `get` or `set`.")
	}
}

func (ConfigCommands) customgiveawaybots(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config.")
		return
	}

	switch utils.Arguments["<get/add/rm/clear>"] {
	case "get":
		botString := ""
		for _, userID := range cfg.Snipers.GiveawaySniper.CustomBots {
			user, ok := utils.Discord.GetUser(ctx.Session, userID)
			if !ok {
				botString += fmt.Sprintf("%s - %s\n", "Unknown", userID)
				continue
			}

			botString += fmt.Sprintf("%s - %s\n", user.Username, userID)
		}

		if botString == "" {
			botString = "No custom bots set."
		}

		embed.NewEmbed().
			AddField("Custom Giveaway Bots", botString).
			Send(ctx)
	case "add":
		bot, ok := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(1))
		if !ok {
			_ = ctx.RespondText("Could not find bot.")
			return
		}

		if utils.SliceContains(cfg.Snipers.GiveawaySniper.CustomBots, bot.ID) {
			_ = ctx.RespondText("Bot is already in the custom bots list.")
			return
		}

		cfg.Snipers.GiveawaySniper.CustomBots = append(cfg.Snipers.GiveawaySniper.CustomBots, bot.ID)

		_ = cfg.Save()

		embed.NewEmbed().
			SetDescription(fmt.Sprintf("Added %s to the custom giveaway bots list.", bot.Mention())).
			Send(ctx)
	case "rm":
		bot, ok := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(1))
		if !ok {
			_ = ctx.RespondText("Could not find bot.")
			return
		}

		if !utils.SliceContains(cfg.Snipers.GiveawaySniper.CustomBots, bot.ID) {
			_ = ctx.RespondText("Bot already isn't in the custom bots list.")
			return
		}

		cfg.Snipers.GiveawaySniper.CustomBots = utils.RemoveSliceItem(cfg.Snipers.GiveawaySniper.CustomBots, bot.ID)

		_ = cfg.Save()

		embed.NewEmbed().
			SetDescription(fmt.Sprintf("Removed %s from the custom giveaway bots list.", bot.Mention())).
			Send(ctx)
	case "clear":
		cfg.Snipers.GiveawaySniper.CustomBots = []string{}
		_ = cfg.Save()

		embed.NewEmbed().
			SetDescription("Cleared all custom giveaway bots.").
			Send(ctx)
	default:
		_ = ctx.RespondText("Invalid <get/add/rm/clear> argument. Use `get`, `add`, `rm`, or `clear`.")
	}
}

func (ConfigCommands) sniperblacklist(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config.")
		return
	}

	switch utils.Arguments["<get/add/rm/clear>"] {
	case "get":
		embed.NewEmbed().
			AddField("Giveaway Sniper Blacklist", strings.Join(cfg.Snipers.Blacklist, ", ")).
			Send(ctx)
	case "add":
		entry := utils.String.RemoveWords(ctx.Arguments.Raw(), 1)

		if utils.SliceContains(cfg.Snipers.Blacklist, entry) {
			_ = ctx.RespondText("Value is already in the blacklist!")
			return
		}

		cfg.Snipers.Blacklist = append(cfg.Snipers.Blacklist, entry)
		_ = cfg.Save()

		embed.NewEmbed().
			SetDescription("Added " + entry + " to the sniper blacklist.").
			Send(ctx)
	case "rm":
		entry := utils.String.RemoveWords(ctx.Arguments.Raw(), 1)

		if !utils.SliceContains(cfg.Snipers.Blacklist, entry) {
			_ = ctx.RespondText("Value already isn't in the blacklist!")
			return
		}

		cfg.Snipers.Blacklist = utils.RemoveSliceItem(cfg.Snipers.Blacklist, entry)
		_ = cfg.Save()

		embed.NewEmbed().
			SetDescription("Removed " + entry + " from the sniper blacklist.").
			Send(ctx)
	case "clear":
		cfg.Snipers.Blacklist = []string{}
		_ = cfg.Save()

		embed.NewEmbed().
			SetDescription("Cleared all sniper blacklist entries.").
			Send(ctx)
	default:
		_ = ctx.RespondText("Invalid <get/add/rm/clear> argument. Use `get`, `add`, `rm`, or `clear`.")
	}
}

func (ConfigCommands) sniperserverblacklist(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config.")
		return
	}

	switch utils.Arguments["<get/add/rm/clear>"] {
	case "get":
		var servers []string

		for _, serverID := range cfg.Snipers.ServerBlacklist {
			guild, err := utils.Discord.GetGuild(ctx.Session, serverID)
			if err != nil {
				servers = append(servers, serverID)
				continue
			}

			servers = append(servers, guild.Name+" ("+guild.ID+")")
		}

		embed.NewEmbed().
			AddField("Sniper Server Blacklist", strings.Join(servers, ", ")).
			Send(ctx)
	case "add":
		guild, err := utils.Discord.GetGuild(ctx.Session, utils.Misc.Or(utils.Arguments["[server id]"], ctx.Event.GuildID))
		if err != nil {
			_ = ctx.RespondText("Could not find server.")
			return
		}

		if utils.SliceContains(cfg.Snipers.ServerBlacklist, guild.ID) {
			_ = ctx.RespondText(fmt.Sprintf("Server %s is already in the blacklist!", guild.Name))
			return
		}

		cfg.Snipers.ServerBlacklist = append(cfg.Snipers.ServerBlacklist, guild.ID)
		_ = cfg.Save()

		embed.NewEmbed().
			SetDescription(fmt.Sprintf("Added %s to the sniper server blacklist.", guild.Name)).
			Send(ctx)
	case "rm":
		guild, err := utils.Discord.GetGuild(ctx.Session, utils.Misc.Or(utils.Arguments["[server id]"], ctx.Event.GuildID))
		if err != nil {
			_ = ctx.RespondText("Could not find server.")
			return
		}

		if !utils.SliceContains(cfg.Snipers.ServerBlacklist, guild.ID) {
			_ = ctx.RespondText(fmt.Sprintf("Server %s already isn't in the blacklist!", guild.Name))
			return
		}

		cfg.Snipers.ServerBlacklist = utils.RemoveSliceItem(cfg.Snipers.ServerBlacklist, guild.ID)
		_ = cfg.Save()

		embed.NewEmbed().
			SetDescription(fmt.Sprintf("Removed %s from the sniper server blacklist.", guild.Name)).
			Send(ctx)
	case "clear":
		cfg.Snipers.ServerBlacklist = []string{}
		_ = cfg.Save()

		embed.NewEmbed().
			SetDescription("Cleared all sniper server blacklist entries.").
			Send(ctx)
	default:
		_ = ctx.RespondText("Invalid <get/add/rm/clear> argument. Use `get`, `add`, `rm`, or `clear`.")
	}
}

func (ConfigCommands) privnotesniper(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	switch utils.Arguments["<get/set/toggle>"] {
	case "get":
		embed.NewEmbed().
			AddField("Privnote Sniper", utils.Misc.FormatBoolOnOff(cfg.Snipers.PrivnoteSniper)).
			Send(ctx)
	case "set":
		toggle := strings.ToLower(utils.Arguments["[on/off]"])
		if toggle != "on" && toggle != "off" {
			_ = ctx.RespondText("Invalid value. Use `on` or `off`.")
			return
		}

		cfg.Snipers.PrivnoteSniper = toggle == "on"
		_ = cfg.Save()

		embed.NewEmbed().
			SetDescription("Turned " + toggle + " Privnote sniper").
			Send(ctx)
	case "toggle":
		cfg.Snipers.PrivnoteSniper = !cfg.Snipers.PrivnoteSniper
		_ = cfg.Save()

		embed.NewEmbed().
			SetDescription("Turned " + utils.Misc.FormatBoolOnOff(cfg.Snipers.PrivnoteSniper) + " Privnote sniper").
			Send(ctx)
	default:
		_ = ctx.RespondText("Invalid <get/set/toggle> argument. Use `get` or `set` or `toggle`.")
	}
}

func (ConfigCommands) multitokensniper(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	switch utils.Arguments["<get/set/toggle>"] {
	case "get":
		embed.NewEmbed().
			AddField("Multi Token Sniper", utils.Misc.FormatBoolOnOff(cfg.Snipers.MultiTokenSniper)).
			Send(ctx)
	case "set":
		toggle := strings.ToLower(utils.Arguments["[on/off]"])
		if toggle != "on" && toggle != "off" {
			_ = ctx.RespondText("Invalid value. Use `on` or `off`.")
			return
		}

		cfg.Snipers.MultiTokenSniper = toggle == "on"
		_ = cfg.Save()

		embed.NewEmbed().
			SetDescription("Turned " + toggle + " multi token sniper").
			Send(ctx)
	case "toggle":
		cfg.Snipers.MultiTokenSniper = !cfg.Snipers.MultiTokenSniper
		_ = cfg.Save()

		embed.NewEmbed().
			SetDescription("Turned " + utils.Misc.FormatBoolOnOff(cfg.Snipers.MultiTokenSniper) + " multi token sniper").
			Send(ctx)
	default:
		_ = ctx.RespondText("Invalid <get/set/toggle> argument. Use `get` or `set` or `toggle`.")
	}
}

func (ConfigCommands) poketwosniper(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	switch utils.Arguments["<get/set/toggle>"] {
	case "get":
		embed.NewEmbed().
			AddField("PokéTwo Sniper", utils.Misc.FormatBoolOnOff(cfg.Snipers.PokeTwoSniper.Enabled)).
			Send(ctx)
	case "set":
		toggle := strings.ToLower(utils.Arguments["[on/off]"])
		if toggle != "on" && toggle != "off" {
			_ = ctx.RespondText("Invalid value. Use `on` or `off`.")
			return
		}

		cfg.Snipers.PokeTwoSniper.Enabled = toggle == "on"
		_ = cfg.Save()

		embed.NewEmbed().
			SetDescription("Turned " + toggle + " PokéTwo sniper").
			Send(ctx)
	case "toggle":
		cfg.Snipers.PokeTwoSniper.Enabled = !cfg.Snipers.PokeTwoSniper.Enabled
		_ = cfg.Save()

		embed.NewEmbed().
			SetDescription("Turned " + utils.Misc.FormatBoolOnOff(cfg.Snipers.PokeTwoSniper.Enabled) + " PokéTwo sniper").
			Send(ctx)
	default:
		_ = ctx.RespondText("Invalid <get/set/toggle> argument. Use `get` or `set` or `toggle`.")
	}
}

func (ConfigCommands) poketwosleep(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	switch utils.Arguments["<get/set>"] {
	case "get":
		formatted := utils.Misc.FormatDuration(time.Duration(cfg.Snipers.PokeTwoSniper.Sleep) * time.Second)

		embed.NewEmbed().
			AddField("PokéTwo Sniper Sleep", formatted).
			Send(ctx)
	case "set":
		delay, err := utils.Misc.ParseDuration(utils.Arguments["[value]"])
		if err != nil {
			_ = ctx.RespondText("Invalid duration: " + err.Error())
			return
		}

		cfg.Snipers.PokeTwoSniper.Sleep = int(delay.Seconds())
		_ = cfg.Save()

		formatted := utils.Misc.FormatDuration(time.Duration(cfg.Snipers.PokeTwoSniper.Sleep) * time.Second)

		embed.NewEmbed().
			SetDescription(fmt.Sprintf("PokéTwo sniper sleep set to %s.", formatted)).
			Send(ctx)
	default:
		_ = ctx.RespondText("Invalid <get/set> argument. Use `get` or `set`.")
	}
}

func (ConfigCommands) poketworandomsleep(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	switch utils.Arguments["<get/set/toggle>"] {
	case "get":
		embed.NewEmbed().
			AddField("PokéTwo Random Sleep", utils.Misc.FormatBoolOnOff(cfg.Snipers.PokeTwoSniper.RandomSleep)).
			Send(ctx)
	case "set":
		toggle := strings.ToLower(utils.Arguments["[on/off]"])
		if toggle != "on" && toggle != "off" {
			_ = ctx.RespondText("Invalid value. Use `on` or `off`.")
			return
		}

		cfg.Snipers.PokeTwoSniper.RandomSleep = toggle == "on"
		_ = cfg.Save()

		embed.NewEmbed().
			SetDescription("Turned " + toggle + " PokéTwo random sleep").
			Send(ctx)
	case "toggle":
		cfg.Snipers.PokeTwoSniper.RandomSleep = !cfg.Snipers.PokeTwoSniper.RandomSleep
		_ = cfg.Save()

		embed.NewEmbed().
			SetDescription("Turned " + utils.Misc.FormatBoolOnOff(cfg.Snipers.PokeTwoSniper.RandomSleep) + " PokéTwo random sleep").
			Send(ctx)
	default:
		_ = ctx.RespondText("Invalid <get/set/toggle> argument. Use `get` or `set` or `toggle`.")
	}
}

func (ConfigCommands) poketwotimeout(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	switch utils.Arguments["<get/set>"] {
	case "get":
		formatted := utils.Misc.FormatDuration(time.Duration(cfg.Snipers.PokeTwoSniper.Timeout) * time.Second)

		embed.NewEmbed().
			AddField("PokéTwo Sniper Timeout", formatted).
			Send(ctx)
	case "set":
		delay, err := utils.Misc.ParseDuration(utils.Arguments["[value]"])
		if err != nil {
			_ = ctx.RespondText("Invalid duration: " + err.Error())
			return
		}

		cfg.Snipers.PokeTwoSniper.Timeout = int(delay.Seconds())
		_ = cfg.Save()

		formatted := utils.Misc.FormatDuration(time.Duration(cfg.Snipers.PokeTwoSniper.Timeout) * time.Second)

		embed.NewEmbed().
			SetDescription(fmt.Sprintf("PokéTwo sniper tiumeout set to %s.", formatted)).
			Send(ctx)
	default:
		_ = ctx.RespondText("Invalid <get/set> argument. Use `get` or `set`.")
	}
}

func (ConfigCommands) poketwotyping(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	switch utils.Arguments["<get/set/toggle>"] {
	case "get":
		embed.NewEmbed().
			AddField("PokéTwo Typing", utils.Misc.FormatBoolOnOff(cfg.Snipers.PokeTwoSniper.Typing)).
			Send(ctx)
	case "set":
		toggle := strings.ToLower(utils.Arguments["[on/off]"])
		if toggle != "on" && toggle != "off" {
			_ = ctx.RespondText("Invalid value. Use `on` or `off`.")
			return
		}

		cfg.Snipers.PokeTwoSniper.Typing = toggle == "on"
		_ = cfg.Save()

		embed.NewEmbed().
			SetDescription("Turned " + toggle + " PokéTwo typing").
			Send(ctx)
	case "toggle":
		cfg.Snipers.PokeTwoSniper.Typing = !cfg.Snipers.PokeTwoSniper.Typing
		_ = cfg.Save()

		embed.NewEmbed().
			SetDescription("Turned " + utils.Misc.FormatBoolOnOff(cfg.Snipers.PokeTwoSniper.Typing) + " PokéTwo typing").
			Send(ctx)
	default:
		_ = ctx.RespondText("Invalid <get/set/toggle> argument. Use `get` or `set` or `toggle`.")
	}
}

func (ConfigCommands) poketwocatchrate(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	switch utils.Arguments["<get/set>"] {
	case "get":
		embed.NewEmbed().
			AddField("PokéTwo Catch Rate", fmt.Sprintf("%d%%", cfg.Snipers.PokeTwoSniper.CatchRate)).
			Send(ctx)
	case "set":
		catchRate, err := strconv.Atoi(strings.ReplaceAll(utils.Arguments["[percentage]"], "%", ""))
		if err != nil || catchRate < 1 || catchRate > 100 {
			_ = ctx.RespondText("Argument [percentage] must be a percentage in the format `1-100%`.")
			return
		}

		cfg.Snipers.PokeTwoSniper.CatchRate = catchRate
		_ = cfg.Save()

		embed.NewEmbed().
			SetDescription(fmt.Sprintf("PokéTwo catch rate set to %d%%", cfg.Snipers.PokeTwoSniper.CatchRate)).
			Send(ctx)
	default:
		_ = ctx.RespondText("Invalid <get/set> argument. Use `get` or `set`.")
	}
}

func (ConfigCommands) poketwofakeactivity(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	switch utils.Arguments["<get/set/toggle>"] {
	case "get":
		embed.NewEmbed().
			AddField("PokéTwo Fake Activity", utils.Misc.FormatBoolOnOff(cfg.Snipers.PokeTwoSniper.FakeActivity)).
			Send(ctx)
	case "set":
		toggle := strings.ToLower(utils.Arguments["[on/off]"])
		if toggle != "on" && toggle != "off" {
			_ = ctx.RespondText("Invalid value. Use `on` or `off`.")
			return
		}

		cfg.Snipers.PokeTwoSniper.FakeActivity = toggle == "on"
		_ = cfg.Save()

		embed.NewEmbed().
			SetDescription("Turned " + toggle + " PokéTwo fake activity").
			Send(ctx)
	case "toggle":
		cfg.Snipers.PokeTwoSniper.FakeActivity = !cfg.Snipers.PokeTwoSniper.FakeActivity
		_ = cfg.Save()

		embed.NewEmbed().
			SetDescription("Turned " + utils.Misc.FormatBoolOnOff(cfg.Snipers.PokeTwoSniper.FakeActivity) + " PokéTwo fake activity").
			Send(ctx)
	default:
		_ = ctx.RespondText("Invalid <get/set/toggle> argument. Use `get` or `set` or `toggle`.")
	}
}

func (ConfigCommands) poketwochannels(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	switch utils.Arguments["<get/add/rm>"] {
	case "get":
		var channels []string

		for _, channelID := range cfg.Snipers.PokeTwoSniper.Channels {
			channel, err := utils.Discord.GetChannel(ctx.Session, channelID)
			if err == nil {
				guildName := "Unknown server"

				guild, err := utils.Discord.GetGuild(ctx.Session, channel.GuildID)
				if err == nil {
					guildName = guild.Name
				}

				channels = append(channels, fmt.Sprintf("%s - %s (%s)", channel.Mention(), channel.ID, guildName))
			} else {
				channels = append(channels, channelID)
			}
		}

		embed.NewEmbed().
			AddField("PokéTwo Sniper Channels", strings.Join(channels, "\n")).
			Send(ctx)
	case "add":
		channel, success := utils.Discord.TryGetChannelFromMention(ctx, ctx.Arguments.Get(1))
		if !success {
			_ = ctx.RespondText("Could not find channel")
			return
		}

		if utils.SliceContains(cfg.Snipers.PokeTwoSniper.Channels, channel.ID) {
			_ = ctx.RespondText("Channel is already in the catch list!")
			return
		}

		cfg.Snipers.PokeTwoSniper.Channels = append(cfg.Snipers.PokeTwoSniper.Channels, channel.ID)
		_ = cfg.Save()

		embed.NewEmbed().
			SetDescription(fmt.Sprintf("Added %s to the PokéTwo sniper catch list.", channel.Mention())).
			Send(ctx)
	case "rm":
		channel, success := utils.Discord.TryGetChannelFromMention(ctx, ctx.Arguments.Get(1))
		if !success {
			_ = ctx.RespondText("Could not find channel")
			return
		}

		if !utils.SliceContains(cfg.Snipers.PokeTwoSniper.Channels, channel.ID) {
			_ = ctx.RespondText("Channel is not in the catch list!")
			return
		}

		cfg.Snipers.PokeTwoSniper.Channels = utils.RemoveSliceItem(cfg.Snipers.PokeTwoSniper.Channels, channel.ID)
		_ = cfg.Save()

		embed.NewEmbed().
			SetDescription(fmt.Sprintf("Removed %s from the PokéTwo sniper catch list.", channel.Mention())).
			Send(ctx)
	case "clear":
		cfg.Snipers.PokeTwoSniper.Channels = []string{}
		_ = cfg.Save()

		embed.NewEmbed().
			SetDescription("Cleared all PokéTwo sniper channels.").
			Send(ctx)
	default:
		_ = ctx.RespondText("Invalid <get/add/rm/clear> argument. Use `get`, `add`, `rm`, or `clear`.")
	}
}

func (ConfigCommands) poketwolanguage(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	switch utils.Arguments["<get/set/values>"] {
	case "get":
		embed.NewEmbed().
			AddField("PokéTwo Catch Language", cfg.Snipers.PokeTwoSniper.CatchLanguage).
			Send(ctx)
	case "set":
		language := strings.ToLower(utils.Arguments["[language]"])

		switch language {
		case "en", "english":
			cfg.Snipers.PokeTwoSniper.CatchLanguage = "English"
		case "ja", "japanese":
			cfg.Snipers.PokeTwoSniper.CatchLanguage = "Japanese"
		case "ja_t", "japanese_t":
			cfg.Snipers.PokeTwoSniper.CatchLanguage = "JapaneseTransliterated"
		case "fr", "french":
			cfg.Snipers.PokeTwoSniper.CatchLanguage = "French"
		case "de", "german":
			cfg.Snipers.PokeTwoSniper.CatchLanguage = "German"
		default:
			_ = ctx.RespondText("Invalid language. `[p]poketwolanguages values` to view all available values.")
			return
		}

		_ = cfg.Save()

		embed.NewEmbed().
			SetDescription(fmt.Sprintf("PokéTwo catch language set to %s.", cfg.Snipers.PokeTwoSniper.CatchLanguage)).
			Send(ctx)
	case "values":
		embed.NewEmbed().
			SetDescription("Valid languages: `english`/`en`, `japanese`/`ja`, `japanese_t`/`ja_t` (transliterated), `french`/`fr`, `german`/`de`.").
			Send(ctx)
	default:
		_ = ctx.RespondText("Invalid <get/set/values> argument. Use `get`, `set`, or `values`.")
	}
}

func (ConfigCommands) snipertokens(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	switch utils.Arguments["<get/add/rm/clear>"] {
	case "get":
		embed.NewEmbed().
			AddField("Sniper Tokens", strings.Join(cfg.Snipers.Tokens, ", ")).
			Send(ctx)
	case "add":
		entries := strings.Split(ctx.Arguments.Raw(), ",")

		var added int

		for _, entry := range entries {
			entry = strings.TrimSpace(entry)

			if utils.SliceContains(cfg.Snipers.Tokens, entry) {
				continue
			}

			cfg.Snipers.Tokens = append(cfg.Snipers.Tokens, entry)
			added++
		}
		_ = cfg.Save()

		embed.NewEmbed().
			SetDescription(fmt.Sprintf("Added %d token(s) to the multi token sniper.", added)).
			Send(ctx)
	case "rm":
		entries := strings.Split(ctx.Arguments.Raw(), ",")

		var added int

		for _, entry := range entries {
			entry = strings.TrimSpace(entry)

			if !utils.SliceContains(cfg.Snipers.Tokens, entry) {
				continue
			}

			cfg.Snipers.Tokens = append(cfg.Snipers.Tokens, entry)
			added++
		}
		_ = cfg.Save()

		embed.NewEmbed().
			SetDescription(fmt.Sprintf("Removed %d token(s) from the multi token sniper.", added)).
			Send(ctx)
	case "clear":
		cfg.Snipers.Tokens = []string{}
		_ = cfg.Save()

		embed.NewEmbed().
			SetDescription("Cleared all sniper tokens.").
			Send(ctx)
	default:
		_ = ctx.RespondText("Invalid <get/add/rm/clear> argument. Use `get`, `add`, `rm`, or `clear`.")
	}
}

func createWebhookChannel(ctx *dgc.Ctx, serverID, categoryID, name string) (string, error) {
	channel, err := ctx.Session.GuildChannelCreateComplex(serverID, discordgo.GuildChannelCreateData{
		Type:     discordgo.ChannelTypeGuildText,
		ParentID: categoryID,
		Name:     name,
	})
	if err != nil {
		return "", err
	}

	webhook, err := ctx.Session.WebhookCreate(channel.ID, name, "")
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("https://discord.com/api/webhooks/%s/%s", webhook.ID, webhook.Token), nil
}

func (ConfigCommands) setupwebhooks(ctx *dgc.Ctx) {
	serverID := ctx.Event.GuildID
	if utils.Arguments["[server id]"] != "" {
		serverID = utils.Arguments["[server id]"]
	}

	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	category, err := ctx.Session.GuildChannelCreate(serverID, "LightningBot", discordgo.ChannelTypeGuildCategory)
	if err != nil {
		_ = ctx.RespondText("Failed to create category: " + err.Error())
		return
	}

	cfg.Webhooks.NitroSniper, _ = createWebhookChannel(ctx, serverID, category.ID, "nitrosniper")
	cfg.Webhooks.GiveawaySniper, _ = createWebhookChannel(ctx, serverID, category.ID, "giveawaysniper")
	cfg.Webhooks.SelfbotDetection, _ = createWebhookChannel(ctx, serverID, category.ID, "selfbotdetection")
	cfg.Webhooks.StaffDetection, _ = createWebhookChannel(ctx, serverID, category.ID, "staffdetection")
	cfg.Webhooks.Notifiers, _ = createWebhookChannel(ctx, serverID, category.ID, "notifiers")
	cfg.Webhooks.PingDetection, _ = createWebhookChannel(ctx, serverID, category.ID, "pingdetection")
	cfg.Webhooks.Misc, _ = createWebhookChannel(ctx, serverID, category.ID, "misc")
	cfg.Webhooks.Sessions, _ = createWebhookChannel(ctx, serverID, category.ID, "sessions")

	_ = cfg.Save()

	embed.NewEmbed().
		SetDescription("Webhooks setup successfully.").
		Send(ctx)
}

func (ConfigCommands) wall(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	webhookURL := utils.Arguments["[webhook url]"]

	cfg.Webhooks.NitroSniper = webhookURL
	cfg.Webhooks.GiveawaySniper = webhookURL
	cfg.Webhooks.SelfbotDetection = webhookURL
	cfg.Webhooks.StaffDetection = webhookURL
	cfg.Webhooks.Notifiers = webhookURL
	cfg.Webhooks.PingDetection = webhookURL
	cfg.Webhooks.Misc = webhookURL
	cfg.Webhooks.Sessions = webhookURL

	_ = cfg.Save()

	embed.NewEmbed().
		SetDescription("Set all webhooks to " + webhookURL).
		Send(ctx)
}

func (ConfigCommands) wnitrosniper(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	switch utils.Arguments["<get/set>"] {
	case "get":
		embed.NewEmbed().
			AddField("Nitro Sniper Webhook", cfg.Webhooks.NitroSniper).
			Send(ctx)
	case "set":
		cfg.Webhooks.NitroSniper = utils.Arguments["[webhook url]"]
		_ = cfg.Save()

		embed.NewEmbed().
			SetDescription("Set nitro sniper webhook to " + utils.Arguments["[webhook url]"]).
			Send(ctx)
	default:
		_ = ctx.RespondText("Invalid <get/set> argument. Use `get` or `set`.")
	}
}

func (ConfigCommands) wgiveawaysniper(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	switch utils.Arguments["<get/set>"] {
	case "get":
		embed.NewEmbed().
			AddField("Giveaway Sniper Webhook", cfg.Webhooks.GiveawaySniper).
			Send(ctx)
	case "set":
		cfg.Webhooks.GiveawaySniper = utils.Arguments["[webhook url]"]
		_ = cfg.Save()

		embed.NewEmbed().
			SetDescription("Set giveaway sniper webhook to " + utils.Arguments["[webhook url]"]).
			Send(ctx)
	default:
		_ = ctx.RespondText("Invalid <get/set> argument. Use `get` or `set`.")
	}
}

func (ConfigCommands) wselfbotdetection(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	switch utils.Arguments["<get/set>"] {
	case "get":
		embed.NewEmbed().
			AddField("Selfbot Detection Webhook", cfg.Webhooks.SelfbotDetection).
			Send(ctx)
	case "set":
		cfg.Webhooks.SelfbotDetection = utils.Arguments["[webhook url]"]
		_ = cfg.Save()

		embed.NewEmbed().
			SetDescription("Set selfbot detection webhook to " + utils.Arguments["[webhook url]"]).
			Send(ctx)
	default:
		_ = ctx.RespondText("Invalid <get/set> argument. Use `get` or `set`.")
	}
}

func (ConfigCommands) wstaffdetection(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	switch utils.Arguments["<get/set>"] {
	case "get":
		embed.NewEmbed().
			AddField("Staff Detection Webhook", cfg.Webhooks.StaffDetection).
			Send(ctx)
	case "set":
		cfg.Webhooks.StaffDetection = utils.Arguments["[webhook url]"]
		_ = cfg.Save()

		embed.NewEmbed().
			SetDescription("Set staff detection webhook to " + utils.Arguments["[webhook url]"]).
			Send(ctx)
	default:
		_ = ctx.RespondText("Invalid <get/set> argument. Use `get` or `set`.")
	}
}

func (ConfigCommands) wnotifiers(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	switch utils.Arguments["<get/set>"] {
	case "get":
		embed.NewEmbed().
			AddField("Notifiers Webhook", cfg.Webhooks.Notifiers).
			Send(ctx)
	case "set":
		cfg.Webhooks.Notifiers = utils.Arguments["[webhook url]"]
		_ = cfg.Save()

		embed.NewEmbed().
			SetDescription("Set notifiers webhook to " + utils.Arguments["[webhook url]"]).
			Send(ctx)
	default:
		_ = ctx.RespondText("Invalid <get/set> argument. Use `get` or `set`.")
	}
}

func (ConfigCommands) wpingdetection(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	switch utils.Arguments["<get/set>"] {
	case "get":
		embed.NewEmbed().
			AddField("Ping Detection Webhook", cfg.Webhooks.PingDetection).
			Send(ctx)
	case "set":
		cfg.Webhooks.PingDetection = utils.Arguments["[webhook url]"]
		_ = cfg.Save()

		embed.NewEmbed().
			SetDescription("Set ping detection webhook to " + utils.Arguments["[webhook url]"]).
			Send(ctx)
	default:
		_ = ctx.RespondText("Invalid <get/set> argument. Use `get` or `set`.")
	}
}

func (ConfigCommands) wmisc(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	switch utils.Arguments["<get/set>"] {
	case "get":
		embed.NewEmbed().
			AddField("Misc Webhook", cfg.Webhooks.Misc).
			Send(ctx)
	case "set":
		cfg.Webhooks.Misc = utils.Arguments["[webhook url]"]
		_ = cfg.Save()

		embed.NewEmbed().
			SetDescription("Set misc webhook to " + utils.Arguments["[webhook url]"]).
			Send(ctx)
	default:
		_ = ctx.RespondText("Invalid <get/set> argument. Use `get` or `set`.")
	}
}

func (ConfigCommands) wsessions(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	switch utils.Arguments["<get/set>"] {
	case "get":
		embed.NewEmbed().
			AddField("Sessions Webhook", cfg.Webhooks.Sessions).
			Send(ctx)
	case "set":
		cfg.Webhooks.Sessions = utils.Arguments["[webhook url]"]
		_ = cfg.Save()

		embed.NewEmbed().
			SetDescription("Set sessions webhook to " + utils.Arguments["[webhook url]"]).
			Send(ctx)
	default:
		_ = ctx.RespondText("Invalid <get/set> argument. Use `get` or `set`.")
	}
}

func (ConfigCommands) wpingself(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	switch utils.Arguments["<get/set/values>"] {
	case "get":
		embed.NewEmbed().
			AddField("Webhook Ping Self", cfg.Webhooks.PingSelfMode).
			Send(ctx)
	case "set":
		value := strings.ToLower(utils.Arguments["[value]"])

		switch value {
		case "disabled":
			cfg.Webhooks.PingSelfMode = "Disabled"
		case "always":
			cfg.Webhooks.PingSelfMode = "Always"
		case "important":
			cfg.Webhooks.PingSelfMode = "Important"
		default:
			_ = ctx.RespondText("Invalid Webhook Ping Self mode. Use `[p]wpingself values` to get a list of available modes.")
			return
		}

		_ = cfg.Save()

		embed.NewEmbed().
			SetDescription("Webhook Ping Self set to " + cfg.Webhooks.PingSelfMode).
			Send(ctx)
	case "values":
		embed.NewEmbed().
			SetDescription("Available Webhook Ping Self modes: Disabled, Always, Important").
			Send(ctx)
	default:
		_ = ctx.RespondText("Invalid <get/set/values> argument. Use `get`, `set`, or `values`.")
	}
}

func (ConfigCommands) wpingeveryone(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	switch utils.Arguments["<get/set/values>"] {
	case "get":
		embed.NewEmbed().
			AddField("Webhook Ping Everyone", cfg.Webhooks.PingEveryoneMode).
			Send(ctx)
	case "set":
		value := strings.ToLower(utils.Arguments["[value]"])

		switch value {
		case "disabled":
			cfg.Webhooks.PingEveryoneMode = "Disabled"
		case "always":
			cfg.Webhooks.PingEveryoneMode = "Always"
		case "important":
			cfg.Webhooks.PingEveryoneMode = "Important"
		default:
			_ = ctx.RespondText("Invalid Webhook Ping Everyone mode. Use `[p]wpingeveryone values` to get a list of available modes.")
			return
		}

		_ = cfg.Save()

		embed.NewEmbed().
			SetDescription("Webhook Ping Everyone set to " + cfg.Webhooks.PingEveryoneMode).
			Send(ctx)
	case "values":
		embed.NewEmbed().
			SetDescription("Available Webhook Ping Everyone modes: Disabled, Always, Important").
			Send(ctx)
	default:
		_ = ctx.RespondText("Invalid <get/set/values> argument. Use `get`, `set`, or `values`.")
	}
}

func (ConfigCommands) bannotifier(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	switch utils.Arguments["<get/set/toggle>"] {
	case "get":
		embed.NewEmbed().
			AddField("Ban Notifier", utils.Misc.FormatBoolOnOff(cfg.Events.BanNotifier)).
			Send(ctx)
	case "set":
		toggle := strings.ToLower(utils.Arguments["[on/off]"])
		if toggle != "on" && toggle != "off" {
			_ = ctx.RespondText("Invalid value. Use `on` or `off`.")
			return
		}

		cfg.Events.BanNotifier = toggle == "on"
		_ = cfg.Save()

		embed.NewEmbed().
			SetDescription("Turned " + toggle + " ban notifier").
			Send(ctx)
	case "toggle":
		cfg.Events.BanNotifier = !cfg.Events.BanNotifier
		_ = cfg.Save()

		embed.NewEmbed().
			SetDescription("Turned " + utils.Misc.FormatBoolOnOff(cfg.Events.BanNotifier) + " ban notifier").
			Send(ctx)
	default:
		_ = ctx.RespondText("Invalid <get/set/toggle> argument. Use `get` or `set` or `toggle`.")
	}
}

func (ConfigCommands) staffdetection(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	switch utils.Arguments["<get/set/toggle>"] {
	case "get":
		embed.NewEmbed().
			AddField("Staff Detection", utils.Misc.FormatBoolOnOff(cfg.Events.DiscordStaffDetection)).
			Send(ctx)
	case "set":
		toggle := strings.ToLower(utils.Arguments["[on/off]"])
		if toggle != "on" && toggle != "off" {
			_ = ctx.RespondText("Invalid value. Use `on` or `off`.")
			return
		}

		cfg.Events.DiscordStaffDetection = toggle == "on"
		_ = cfg.Save()

		embed.NewEmbed().
			SetDescription("Turned " + toggle + " staff detection").
			Send(ctx)
	case "toggle":
		cfg.Events.DiscordStaffDetection = !cfg.Events.DiscordStaffDetection
		_ = cfg.Save()

		embed.NewEmbed().
			SetDescription("Turned " + utils.Misc.FormatBoolOnOff(cfg.Events.DiscordStaffDetection) + " staff detection").
			Send(ctx)
	default:
		_ = ctx.RespondText("Invalid <get/set/toggle> argument. Use `get` or `set` or `toggle`.")
	}
}

func (ConfigCommands) pingdetection(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	switch utils.Arguments["<get/set/toggle>"] {
	case "get":
		embed.NewEmbed().
			AddField("Ping Detection", utils.Misc.FormatBoolOnOff(cfg.Events.GhostPingDetection)).
			Send(ctx)
	case "set":
		toggle := strings.ToLower(utils.Arguments["[on/off]"])
		if toggle != "on" && toggle != "off" {
			_ = ctx.RespondText("Invalid value. Use `on` or `off`.")
			return
		}

		cfg.Events.GhostPingDetection = toggle == "on"
		_ = cfg.Save()

		embed.NewEmbed().
			SetDescription("Turned " + toggle + " ping detection").
			Send(ctx)
	case "toggle":
		cfg.Events.GhostPingDetection = !cfg.Events.GhostPingDetection
		_ = cfg.Save()

		embed.NewEmbed().
			SetDescription("Turned " + utils.Misc.FormatBoolOnOff(cfg.Events.GhostPingDetection) + " ping detection").
			Send(ctx)
	default:
		_ = ctx.RespondText("Invalid <get/set/toggle> argument. Use `get` or `set` or `toggle`.")
	}
}

func (ConfigCommands) ticketnotifier(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	switch utils.Arguments["<get/set/toggle>"] {
	case "get":
		embed.NewEmbed().
			AddField("Ticket Notifier", utils.Misc.FormatBoolOnOff(cfg.Events.TicketNotifier)).
			Send(ctx)
	case "set":
		toggle := strings.ToLower(utils.Arguments["[on/off]"])
		if toggle != "on" && toggle != "off" {
			_ = ctx.RespondText("Invalid value. Use `on` or `off`.")
			return
		}

		cfg.Events.TicketNotifier = toggle == "on"
		_ = cfg.Save()

		embed.NewEmbed().
			SetDescription("Turned " + toggle + " ticket notifier").
			Send(ctx)
	case "toggle":
		cfg.Events.TicketNotifier = !cfg.Events.TicketNotifier
		_ = cfg.Save()

		embed.NewEmbed().
			SetDescription("Turned " + utils.Misc.FormatBoolOnOff(cfg.Events.TicketNotifier) + " ticket notifier").
			Send(ctx)
	default:
		_ = ctx.RespondText("Invalid <get/set/toggle> argument. Use `get` or `set` or `toggle`.")
	}
}

func (ConfigCommands) sessionconnected(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	switch utils.Arguments["<get/set/toggle>"] {
	case "get":
		embed.NewEmbed().
			AddField("Session Connected Notifier", utils.Misc.FormatBoolOnOff(cfg.Events.SessionConnected)).
			Send(ctx)
	case "set":
		toggle := strings.ToLower(utils.Arguments["[on/off]"])
		if toggle != "on" && toggle != "off" {
			_ = ctx.RespondText("Invalid value. Use `on` or `off`.")
			return
		}

		cfg.Events.SessionConnected = toggle == "on"
		_ = cfg.Save()

		embed.NewEmbed().
			SetDescription("Turned " + toggle + " session connected notifier").
			Send(ctx)
	case "toggle":
		cfg.Events.SessionConnected = !cfg.Events.SessionConnected
		_ = cfg.Save()

		embed.NewEmbed().
			SetDescription("Turned " + utils.Misc.FormatBoolOnOff(cfg.Events.SessionConnected) + " session connected notifier").
			Send(ctx)
	default:
		_ = ctx.RespondText("Invalid <get/set/toggle> argument. Use `get` or `set` or `toggle`.")
	}
}

func (ConfigCommands) sessiondisconnected(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	switch utils.Arguments["<get/set/toggle>"] {
	case "get":
		embed.NewEmbed().
			AddField("Session Disconnected Notifier", utils.Misc.FormatBoolOnOff(cfg.Events.SessionDisconnected)).
			Send(ctx)
	case "set":
		toggle := strings.ToLower(utils.Arguments["[on/off]"])
		if toggle != "on" && toggle != "off" {
			_ = ctx.RespondText("Invalid value. Use `on` or `off`.")
			return
		}

		cfg.Events.SessionDisconnected = toggle == "on"
		_ = cfg.Save()

		embed.NewEmbed().
			SetDescription("Turned " + toggle + " session disconnected notifier").
			Send(ctx)
	case "toggle":
		cfg.Events.SessionDisconnected = !cfg.Events.SessionDisconnected
		_ = cfg.Save()

		embed.NewEmbed().
			SetDescription("Turned " + utils.Misc.FormatBoolOnOff(cfg.Events.SessionDisconnected) + " session disconnected notifier").
			Send(ctx)
	default:
		_ = ctx.RespondText("Invalid <get/set/toggle> argument. Use `get` or `set` or `toggle`.")
	}
}

func (ConfigCommands) selfbotdetection(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	switch utils.Arguments["<get/set/toggle>"] {
	case "get":
		embed.NewEmbed().
			AddField("Selfbot Detection", utils.Misc.FormatBoolOnOff(cfg.Events.SelfbotDetection)).
			Send(ctx)
	case "set":
		toggle := strings.ToLower(utils.Arguments["[on/off]"])
		if toggle != "on" && toggle != "off" {
			_ = ctx.RespondText("Invalid value. Use `on` or `off`.")
			return
		}

		cfg.Events.SelfbotDetection = toggle == "on"
		_ = cfg.Save()

		embed.NewEmbed().
			SetDescription("Turned " + toggle + " selfbot detection").
			Send(ctx)
	case "toggle":
		cfg.Events.SelfbotDetection = !cfg.Events.SelfbotDetection
		_ = cfg.Save()

		embed.NewEmbed().
			SetDescription("Turned " + utils.Misc.FormatBoolOnOff(cfg.Events.SelfbotDetection) + " selfbot detection").
			Send(ctx)
	default:
		_ = ctx.RespondText("Invalid <get/set/toggle> argument. Use `get` or `set` or `toggle`.")
	}
}

func (ConfigCommands) dmtypingnotifier(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	switch utils.Arguments["<get/set/toggle>"] {
	case "get":
		embed.NewEmbed().
			AddField("DM Typing Notifier", utils.Misc.FormatBoolOnOff(cfg.Events.DMTypingNotifier)).
			Send(ctx)
	case "set":
		toggle := strings.ToLower(utils.Arguments["[on/off]"])
		if toggle != "on" && toggle != "off" {
			_ = ctx.RespondText("Invalid value. Use `on` or `off`.")
			return
		}

		cfg.Events.DMTypingNotifier = toggle == "on"
		_ = cfg.Save()

		embed.NewEmbed().
			SetDescription("Turned " + toggle + " DM typing notifier").
			Send(ctx)
	case "toggle":
		cfg.Events.DMTypingNotifier = !cfg.Events.DMTypingNotifier
		_ = cfg.Save()

		embed.NewEmbed().
			SetDescription("Turned " + utils.Misc.FormatBoolOnOff(cfg.Events.DMTypingNotifier) + " DM typing notifier").
			Send(ctx)
	default:
		_ = ctx.RespondText("Invalid <get/set/toggle> argument. Use `get` or `set` or `toggle`.")
	}
}

func (ConfigCommands) dmnotifier(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	switch utils.Arguments["<get/set/toggle>"] {
	case "get":
		embed.NewEmbed().
			AddField("DM Notifier", utils.Misc.FormatBoolOnOff(cfg.Events.DMNotifier)).
			Send(ctx)
	case "set":
		toggle := strings.ToLower(utils.Arguments["[on/off]"])
		if toggle != "on" && toggle != "off" {
			_ = ctx.RespondText("Invalid value. Use `on` or `off`.")
			return
		}

		cfg.Events.DMNotifier = toggle == "on"
		_ = cfg.Save()

		embed.NewEmbed().
			SetDescription("Turned " + toggle + " DM notifier").
			Send(ctx)
	case "toggle":
		cfg.Events.DMNotifier = !cfg.Events.DMNotifier
		_ = cfg.Save()

		embed.NewEmbed().
			SetDescription("Turned " + utils.Misc.FormatBoolOnOff(cfg.Events.DMNotifier) + " DM notifier").
			Send(ctx)
	default:
		_ = ctx.RespondText("Invalid <get/set/toggle> argument. Use `get` or `set` or `toggle`.")
	}
}

func (ConfigCommands) dmdeletenotifier(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	switch utils.Arguments["<get/set/toggle>"] {
	case "get":
		embed.NewEmbed().
			AddField("DM Delete Notifier", utils.Misc.FormatBoolOnOff(cfg.Events.DMDeleteNotifier)).
			Send(ctx)
	case "set":
		toggle := strings.ToLower(utils.Arguments["[on/off]"])
		if toggle != "on" && toggle != "off" {
			_ = ctx.RespondText("Invalid value. Use `on` or `off`.")
			return
		}

		cfg.Events.DMDeleteNotifier = toggle == "on"
		_ = cfg.Save()

		embed.NewEmbed().
			SetDescription("Turned " + toggle + " DM delete notifier").
			Send(ctx)
	case "toggle":
		cfg.Events.DMDeleteNotifier = !cfg.Events.DMDeleteNotifier
		_ = cfg.Save()

		embed.NewEmbed().
			SetDescription("Turned " + utils.Misc.FormatBoolOnOff(cfg.Events.DMDeleteNotifier) + " DM delete notifier").
			Send(ctx)
	default:
		_ = ctx.RespondText("Invalid <get/set/toggle> argument. Use `get` or `set` or `toggle`.")
	}
}

func (ConfigCommands) dmeditnotifier(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	switch utils.Arguments["<get/set/toggle>"] {
	case "get":
		embed.NewEmbed().
			AddField("DM Edit Notifier", utils.Misc.FormatBoolOnOff(cfg.Events.DMEditNotifier)).
			Send(ctx)
	case "set":
		toggle := strings.ToLower(utils.Arguments["[on/off]"])
		if toggle != "on" && toggle != "off" {
			_ = ctx.RespondText("Invalid value. Use `on` or `off`.")
			return
		}

		cfg.Events.DMEditNotifier = toggle == "on"
		_ = cfg.Save()

		embed.NewEmbed().
			SetDescription("Turned " + toggle + " DM edit notifier").
			Send(ctx)
	case "toggle":
		cfg.Events.DMEditNotifier = !cfg.Events.DMEditNotifier
		_ = cfg.Save()

		embed.NewEmbed().
			SetDescription("Turned " + utils.Misc.FormatBoolOnOff(cfg.Events.DMEditNotifier) + " DM edit notifier").
			Send(ctx)
	default:
		_ = ctx.RespondText("Invalid <get/set/toggle> argument. Use `get` or `set` or `toggle`.")
	}
}

func (ConfigCommands) messagedeletenotifier(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	switch utils.Arguments["<get/set/toggle>"] {
	case "get":
		embed.NewEmbed().
			AddField("Message Delete Notifier", utils.Misc.FormatBoolOnOff(cfg.Events.MessageDeleteNotifier)).
			Send(ctx)
	case "set":
		toggle := strings.ToLower(utils.Arguments["[on/off]"])
		if toggle != "on" && toggle != "off" {
			_ = ctx.RespondText("Invalid value. Use `on` or `off`.")
			return
		}

		cfg.Events.MessageDeleteNotifier = toggle == "on"
		_ = cfg.Save()

		embed.NewEmbed().
			SetDescription("Turned " + toggle + " message delete notifier").
			Send(ctx)
	case "toggle":
		cfg.Events.MessageDeleteNotifier = !cfg.Events.MessageDeleteNotifier
		_ = cfg.Save()

		embed.NewEmbed().
			SetDescription("Turned " + utils.Misc.FormatBoolOnOff(cfg.Events.MessageDeleteNotifier) + " message delete notifier").
			Send(ctx)
	default:
		_ = ctx.RespondText("Invalid <get/set/toggle> argument. Use `get` or `set` or `toggle`.")
	}
}

func (ConfigCommands) uitheme(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	switch utils.Arguments["<get/set/values>"] {
	case "get":
		embed.NewEmbed().
			AddField("UI Theme", cfg.UI.Theme).
			Send(ctx)
	case "set":
		theme := strings.ToLower(utils.Arguments["[value]"])
		if theme != "light" && theme != "dark" {
			_ = ctx.RespondText("Invalid value. Use `light` or `dark`.")
			return
		}

		cfg.UI.Theme = utils.String.Capitalize(theme)
		_ = cfg.Save()

		embed.NewEmbed().
			SetDescription("Set UI theme to " + cfg.UI.Theme).
			Send(ctx)
	case "values":
		embed.NewEmbed().
			SetDescription("Valid themes: `light`, `dark`.").
			Send(ctx)
	default:
		_ = ctx.RespondText("Invalid <get/set/values> argument. Use `get`, `set`, or `values`.")
	}
}

func (ConfigCommands) hidetosystray(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	switch utils.Arguments["<get/set/toggle>"] {
	case "get":
		embed.NewEmbed().
			AddField("Hide to System Tray", utils.Misc.FormatBoolOnOff(cfg.UI.HideToSystemTray)).
			Send(ctx)
	case "set":
		toggle := strings.ToLower(utils.Arguments["[on/off]"])
		if toggle != "on" && toggle != "off" {
			_ = ctx.RespondText("Invalid value. Use `on` or `off`.")
			return
		}

		cfg.UI.HideToSystemTray = toggle == "on"
		_ = cfg.Save()

		embed.NewEmbed().
			SetDescription("Turned " + toggle + " hide to system tray").
			Send(ctx)
	case "toggle":
		cfg.UI.HideToSystemTray = !cfg.UI.HideToSystemTray
		_ = cfg.Save()

		embed.NewEmbed().
			SetDescription("Turned " + utils.Misc.FormatBoolOnOff(cfg.UI.HideToSystemTray) + " hide to system tray").
			Send(ctx)
	default:
		_ = ctx.RespondText("Invalid <get/set/toggle> argument. Use `get` or `set` or `toggle`.")
	}
}

func (ConfigCommands) runuionstartup(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	switch utils.Arguments["<get/set/toggle>"] {
	case "get":
		embed.NewEmbed().
			AddField("Run UI On Startup", utils.Misc.FormatBoolOnOff(cfg.UI.RunOnStartup)).
			Send(ctx)
	case "set":
		toggle := strings.ToLower(utils.Arguments["[on/off]"])
		if toggle != "on" && toggle != "off" {
			_ = ctx.RespondText("Invalid value. Use `on` or `off`.")
			return
		}

		cfg.UI.RunOnStartup = toggle == "on"
		_ = cfg.Save()

		embed.NewEmbed().
			SetDescription("Turned " + toggle + " run UI on startup").
			Send(ctx)
	case "toggle":
		cfg.UI.RunOnStartup = !cfg.UI.RunOnStartup
		_ = cfg.Save()

		embed.NewEmbed().
			SetDescription("Turned " + utils.Misc.FormatBoolOnOff(cfg.UI.RunOnStartup) + " run UI on startup").
			Send(ctx)
	default:
		_ = ctx.RespondText("Invalid <get/set/toggle> argument. Use `get` or `set` or `toggle`.")
	}
}

func (ConfigCommands) runbotonstartup(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	switch utils.Arguments["<get/set/toggle>"] {
	case "get":
		embed.NewEmbed().
			AddField("Run Selfbot On Startup", utils.Misc.FormatBoolOnOff(cfg.UI.RunSelfbotOnStartup)).
			Send(ctx)
	case "set":
		toggle := strings.ToLower(utils.Arguments["[on/off]"])
		if toggle != "on" && toggle != "off" {
			_ = ctx.RespondText("Invalid value. Use `on` or `off`.")
			return
		}

		cfg.UI.RunSelfbotOnStartup = toggle == "on"
		_ = cfg.Save()

		embed.NewEmbed().
			SetDescription("Turned " + toggle + " run selfbot on startup").
			Send(ctx)
	case "toggle":
		cfg.UI.RunSelfbotOnStartup = !cfg.UI.RunSelfbotOnStartup
		_ = cfg.Save()

		embed.NewEmbed().
			SetDescription("Turned " + utils.Misc.FormatBoolOnOff(cfg.UI.RunSelfbotOnStartup) + " run selfbot on startup").
			Send(ctx)
	default:
		_ = ctx.RespondText("Invalid <get/set/toggle> argument. Use `get` or `set` or `toggle`.")
	}
}
