package commands

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/LightningDev1/LB-Selfbot/config"
	"github.com/LightningDev1/LB-Selfbot/ipc"
	"github.com/LightningDev1/LB-Selfbot/utils"
	"github.com/LightningDev1/dgc"
	"github.com/LightningDev1/discordgo"
)

type ConfigCommands struct{}

func (c ConfigCommands) RegisterCommands(router *dgc.Router) {
	router.RegisterCmd(&dgc.Command{
		Name:        "commandprefix",
		Description: "Set your command prefix",
		Usage:       "[p]commandprefix <prefix>",
		Handler:     c.commandprefix,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "commandmode",
		Description: "Set your command mode",
		Usage:       "[p]commandmode <codeblock/text/urlembed>",
		Handler:     c.commandmode,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "title",
		Description: "Set your title",
		Usage:       "[p]title <title>",
		Handler:     c.title,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "footer",
		Description: "Set your footer",
		Usage:       "[p]footer <footer>",
		Handler:     c.footer,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "delresponses",
		Description: "Set if the bot should delete responses to commands",
		Usage:       "[p]delresponses <on/off>",
		Handler:     c.delresponses,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "delresponsesdelay",
		Description: "Set the delay before deleting responses",
		Usage:       "[p]delresponsesdelay <seconds>",
		Handler:     c.delresponsesdelay,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "delcommands",
		Description: "Set if the bot should delete commands",
		Usage:       "[p]delcommands <on/off>",
		Handler:     c.delcommands,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "notifications",
		Description: "Turn toast notifications on or off",
		Usage:       "[p]notifications <on/off>",
		Handler:     c.notifications,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "notificationtype",
		Description: "Change the notification type",
		Usage:       "[p]notificationtype <default/custom>",
		Handler:     c.notificationtype,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "examplenotification",
		Description: "Show an example notification",
		Usage:       "[p]examplenotification",
		Handler:     c.examplenotification,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "nitrosniper",
		Description: "Turn the nitro sniper on or off",
		Usage:       "[p]nitrosniper <on/off>",
		Handler:     c.nitrosniper,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "giveawaysniper",
		Description: "Turn the giveaway sniper on or off",
		Usage:       "[p]giveawaysniper <on/off>",
		Handler:     c.giveawaysniper,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "gsnipersleep",
		Description: "Set the delay before the bot enters a giveaway",
		Usage:       "[p]gsnipersleep <on/off>",
		Handler:     c.gsnipersleep,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "setupwebhooks",
		Description: "Automatically set up all webhooks",
		Usage:       "[p]setupwebhooks [server id]",
		Handler:     c.setupwebhooks,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "wall",
		Description: "Set all webhooks at once",
		Usage:       "[p]wall <webhook url>",
		Handler:     c.wall,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "wnitro",
		Description: "Set the nitro sniper webhook",
		Usage:       "[p]wnitro <webhook url>",
		Handler:     c.wnitro,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "wgiveaway",
		Description: "Set the giveaway sniper webhook",
		Usage:       "[p]wgiveaway <webhook url>",
		Handler:     c.wgiveaway,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "wprivnote",
		Description: "Set the privnote sniper webhook",
		Usage:       "[p]wprivnote <webhook url>",
		Handler:     c.wprivnote,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "wselfbotdetection",
		Description: "Set the selfbot detection webhook",
		Usage:       "[p]wselfbotdetection <webhook url>",
		Handler:     c.wselfbotdetection,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "wstaffdetection",
		Description: "Set the staff detection webhook",
		Usage:       "[p]wstaffdetection <webhook url>",
		Handler:     c.wstaffdetection,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "wnotifiers",
		Description: "Set the notifiers webhook",
		Usage:       "[p]wnotifiers <webhook url>",
		Handler:     c.wnotifiers,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "wpingdetection",
		Description: "Set the ping detection webhook",
		Usage:       "[p]wpingdetection <webhook url>",
		Handler:     c.wpingdetection,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "wmisc",
		Description: "Set the misc webhook",
		Usage:       "[p]wmisc <webhook url>",
		Handler:     c.wmisc,
	})
}

func (ConfigCommands) commandprefix(ctx *dgc.Ctx) {
	cfg, _ := config.GetConfig()
	cfg.Selfbot.CommandPrefix = utils.Arguments["<prefix>"]
	cfg.Save()
}

func (ConfigCommands) commandmode(ctx *dgc.Ctx) {
	cfg, _ := config.GetConfig()
	mode := utils.Capitalize(utils.Arguments["<codeblock/text/urlembed>"])
	if mode != "Codeblock" && mode != "Text" && mode != "Urlembed" {
		ctx.RespondText("Invalid command mode.")
		return
	}
	cfg.Selfbot.CommandMode = strings.ReplaceAll(mode, "Urlembed", "URLEmbed")
	cfg.Save()
}

func (ConfigCommands) title(ctx *dgc.Ctx) {
	cfg, _ := config.GetConfig()
	cfg.Selfbot.Title = ctx.Arguments.Raw()
	cfg.Save()
}

func (ConfigCommands) footer(ctx *dgc.Ctx) {
	cfg, _ := config.GetConfig()
	cfg.Selfbot.Footer = ctx.Arguments.Raw()
	cfg.Save()
}

func (ConfigCommands) delresponses(ctx *dgc.Ctx) {
	cfg, _ := config.GetConfig()
	toggle := strings.ToLower(utils.Arguments["<on/off>"])
	if toggle != "on" && toggle != "off" {
		ctx.RespondText("Invalid value. Use `on` or `off`.")
		return
	}
	cfg.Selfbot.DeleteResponses = toggle == "on"
	cfg.Save()
}

func (ConfigCommands) delresponsesdelay(ctx *dgc.Ctx) {
	cfg, _ := config.GetConfig()
	delay, err := strconv.Atoi(utils.Arguments["<seconds>"])
	if err != nil {
		ctx.RespondText("Argument <seconds> must be a number.")
		return
	}
	cfg.Selfbot.DeleteResponsesDelay = delay
	cfg.Save()
}

func (ConfigCommands) delcommands(ctx *dgc.Ctx) {
	cfg, _ := config.GetConfig()
	toggle := strings.ToLower(utils.Arguments["<on/off>"])
	if toggle != "on" && toggle != "off" {
		ctx.RespondText("Invalid value. Use `on` or `off`.")
		return
	}
	cfg.Selfbot.DeleteCommands = toggle == "on"
	cfg.Save()
}

func (ConfigCommands) notifications(ctx *dgc.Ctx) {
	cfg, _ := config.GetConfig()
	toggle := strings.ToLower(utils.Arguments["<on/off>"])
	if toggle != "on" && toggle != "off" {
		ctx.RespondText("Invalid value. Use `on` or `off`.")
		return
	}
	cfg.Notifications.Enabled = toggle == "on"
	cfg.Save()
}

func (ConfigCommands) notificationtype(ctx *dgc.Ctx) {
	cfg, _ := config.GetConfig()
	notificationType := utils.Capitalize(utils.Arguments["<default/custom>"])
	if notificationType != "Default" && notificationType != "Custom" {
		ctx.RespondText("Invalid command mode.")
		return
	}
	cfg.Notifications.Type = notificationType
	cfg.Save()
}

func (ConfigCommands) examplenotification(ctx *dgc.Ctx) {
	ipc.SendNotification("This is an example notification.")
}

func (ConfigCommands) nitrosniper(ctx *dgc.Ctx) {
	cfg, _ := config.GetConfig()
	toggle := strings.ToLower(utils.Arguments["<on/off>"])
	if toggle != "on" && toggle != "off" {
		ctx.RespondText("Invalid value. Use `on` or `off`.")
		return
	}
	cfg.Snipers.NitroSniper = toggle == "on"
	cfg.Save()
}

func (ConfigCommands) giveawaysniper(ctx *dgc.Ctx) {
	cfg, _ := config.GetConfig()
	toggle := strings.ToLower(utils.Arguments["<on/off>"])
	if toggle != "on" && toggle != "off" {
		ctx.RespondText("Invalid value. Use `on` or `off`.")
		return
	}
	cfg.Snipers.GiveawaySniper.Sniper = toggle == "on"
	cfg.Save()
}

func (ConfigCommands) gsnipersleep(ctx *dgc.Ctx) {
	cfg, _ := config.GetConfig()
	delay, err := strconv.Atoi(utils.Arguments["<seconds>"])
	if err != nil {
		ctx.RespondText("Argument <seconds> must be a number.")
		return
	}
	cfg.Snipers.GiveawaySniper.Sleep = delay
	cfg.Save()
}

func (ConfigCommands) setupwebhooks(ctx *dgc.Ctx) {
	serverID := ctx.Event.GuildID
	if utils.Arguments["[server id]"] != "" {
		serverID = utils.Arguments["[server id]"]
	}
	channels := map[string]string{
		"nitrosniper":      "NitroSniper",
		"giveawaysniper":   "GiveawaySniper",
		"privnotesniper":   "PrivnoteSniper",
		"selfbotdetection": "SelfbotDetection",
		"staffdetection":   "StaffDetection",
		"notifiers":        "Notifiers",
		"pingdetection":    "PingDetection",
		"misc":             "Misc",
	}
	cfg, _ := config.GetConfig()
	category, err := ctx.Session.GuildChannelCreate(serverID, "LightningBot", discordgo.ChannelTypeGuildCategory)
	if err != nil {
		ctx.RespondText("Failed to create channel: " + err.Error())
		return
	}
	for channelName, configName := range channels {
		channel, err := ctx.Session.GuildChannelCreateComplex(serverID, discordgo.GuildChannelCreateData{
			Type:     discordgo.ChannelTypeGuildText,
			ParentID: category.ID,
			Name:     channelName,
		})
		if err != nil {
			ctx.RespondText("Failed to create channel: " + err.Error())
			return
		}
		webhook, err := ctx.Session.WebhookCreate(channel.ID, channelName, "")
		if err != nil {
			ctx.RespondText("Failed to create webhook: " + err.Error())
			return
		}
		f := reflect.Indirect(reflect.ValueOf(cfg.Webhooks)).FieldByName(configName)
		f.SetString(fmt.Sprintf("https://discord.com/api/webhooks/%s/%s", webhook.ID, webhook.Token))
	}
	cfg.Save()
}

func (ConfigCommands) wall(ctx *dgc.Ctx) {
	cfg, _ := config.GetConfig()
	cfg.Webhooks.NitroSniper = utils.Arguments["<webhook url>"]
	cfg.Webhooks.GiveawaySniper = utils.Arguments["<webhook url>"]
	cfg.Webhooks.PrivnoteSniper = utils.Arguments["<webhook url>"]
	cfg.Webhooks.SelfbotDetection = utils.Arguments["<webhook url>"]
	cfg.Webhooks.StaffDetection = utils.Arguments["<webhook url>"]
	cfg.Webhooks.Notifiers = utils.Arguments["<webhook url>"]
	cfg.Webhooks.PingDetection = utils.Arguments["<webhook url>"]
	cfg.Webhooks.Misc = utils.Arguments["<webhook url>"]
	cfg.Save()
}

func (ConfigCommands) wnitro(ctx *dgc.Ctx) {
	cfg, _ := config.GetConfig()
	cfg.Webhooks.NitroSniper = utils.Arguments["<webhook url>"]
	cfg.Save()
}

func (ConfigCommands) wgiveaway(ctx *dgc.Ctx) {
	cfg, _ := config.GetConfig()
	cfg.Webhooks.GiveawaySniper = utils.Arguments["<webhook url>"]
	cfg.Save()
}

func (ConfigCommands) wprivnote(ctx *dgc.Ctx) {
	cfg, _ := config.GetConfig()
	cfg.Webhooks.PrivnoteSniper = utils.Arguments["<webhook url>"]
	cfg.Save()
}

func (ConfigCommands) wselfbotdetection(ctx *dgc.Ctx) {
	cfg, _ := config.GetConfig()
	cfg.Webhooks.SelfbotDetection = utils.Arguments["<webhook url>"]
	cfg.Save()
}

func (ConfigCommands) wstaffdetection(ctx *dgc.Ctx) {
	cfg, _ := config.GetConfig()
	cfg.Webhooks.StaffDetection = utils.Arguments["<webhook url>"]
	cfg.Save()
}

func (ConfigCommands) wnotifiers(ctx *dgc.Ctx) {
	cfg, _ := config.GetConfig()
	cfg.Webhooks.Notifiers = utils.Arguments["<webhook url>"]
	cfg.Save()
}

func (ConfigCommands) wpingdetection(ctx *dgc.Ctx) {
	cfg, _ := config.GetConfig()
	cfg.Webhooks.PingDetection = utils.Arguments["<webhook url>"]
	cfg.Save()
}

func (ConfigCommands) wmisc(ctx *dgc.Ctx) {
	cfg, _ := config.GetConfig()
	cfg.Webhooks.Misc = utils.Arguments["<webhook url>"]
	cfg.Save()
}
