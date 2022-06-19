package commands

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/LightningDev1/LB-Selfbot/config"
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
}

func (c ConfigCommands) commandprefix(ctx *dgc.Ctx) {
	cfg, _ := config.GetConfig()
	cfg.Selfbot.CommandPrefix = utils.Arguments["<prefix>"]
	cfg.Save()
}

func (c ConfigCommands) commandmode(ctx *dgc.Ctx) {
	cfg, _ := config.GetConfig()
	mode := utils.Capitalize(utils.Arguments["<codeblock/text/urlembed>"])
	if mode != "Codeblock" && mode != "Text" && mode != "Urlembed" {
		ctx.RespondText("Invalid command mode.")
		return
	}
	cfg.Selfbot.CommandMode = strings.ReplaceAll(mode, "Urlembed", "URLEmbed")
	cfg.Save()
}

func (c ConfigCommands) title(ctx *dgc.Ctx) {
	cfg, _ := config.GetConfig()
	cfg.Selfbot.Title = ctx.Arguments.Raw()
	cfg.Save()
}

func (c ConfigCommands) footer(ctx *dgc.Ctx) {
	cfg, _ := config.GetConfig()
	cfg.Selfbot.Footer = ctx.Arguments.Raw()
	cfg.Save()
}

func (c ConfigCommands) delresponses(ctx *dgc.Ctx) {
	cfg, _ := config.GetConfig()
	toggle := strings.ToLower(utils.Arguments["<on/off>"])
	if toggle != "on" && toggle != "off" {
		ctx.RespondText("Invalid value. Use `on` or `off`.")
		return
	}
	cfg.Selfbot.DeleteResponses = toggle == "on"
	cfg.Save()
}

func (c ConfigCommands) delresponsesdelay(ctx *dgc.Ctx) {
	cfg, _ := config.GetConfig()
	delay, err := strconv.Atoi(utils.Arguments["<seconds>"])
	if err != nil {
		ctx.RespondText("Argument <seconds> must be a number.")
		return
	}
	cfg.Selfbot.DeleteResponsesDelay = delay
	cfg.Save()
}

func (c ConfigCommands) delcommands(ctx *dgc.Ctx) {
	cfg, _ := config.GetConfig()
	toggle := strings.ToLower(utils.Arguments["<on/off>"])
	if toggle != "on" && toggle != "off" {
		ctx.RespondText("Invalid value. Use `on` or `off`.")
		return
	}
	cfg.Selfbot.DeleteCommands = toggle == "on"
	cfg.Save()
}

func (c ConfigCommands) notifications(ctx *dgc.Ctx) {
	cfg, _ := config.GetConfig()
	toggle := strings.ToLower(utils.Arguments["<on/off>"])
	if toggle != "on" && toggle != "off" {
		ctx.RespondText("Invalid value. Use `on` or `off`.")
		return
	}
	cfg.Notifications.Enabled = toggle == "on"
	cfg.Save()
}

func (c ConfigCommands) notificationtype(ctx *dgc.Ctx) {
	cfg, _ := config.GetConfig()
	notificationType := utils.Capitalize(utils.Arguments["<default/custom>"])
	if notificationType != "Default" && notificationType != "Custom" {
		ctx.RespondText("Invalid command mode.")
		return
	}
	cfg.Notifications.Type = notificationType
	cfg.Save()
}

func (c ConfigCommands) examplenotification(ctx *dgc.Ctx) {
	utils.SendNotification("This is an example notification.")
}

func (c ConfigCommands) nitrosniper(ctx *dgc.Ctx) {
	cfg, _ := config.GetConfig()
	toggle := strings.ToLower(utils.Arguments["<on/off>"])
	if toggle != "on" && toggle != "off" {
		ctx.RespondText("Invalid value. Use `on` or `off`.")
		return
	}
	cfg.Snipers.NitroSniper = toggle == "on"
	cfg.Save()
}

func (c ConfigCommands) giveawaysniper(ctx *dgc.Ctx) {
	cfg, _ := config.GetConfig()
	toggle := strings.ToLower(utils.Arguments["<on/off>"])
	if toggle != "on" && toggle != "off" {
		ctx.RespondText("Invalid value. Use `on` or `off`.")
		return
	}
	cfg.Snipers.GiveawaySniper.Sniper = toggle == "on"
	cfg.Save()
}

func (c ConfigCommands) gsnipersleep(ctx *dgc.Ctx) {
	cfg, _ := config.GetConfig()
	delay, err := strconv.Atoi(utils.Arguments["<seconds>"])
	if err != nil {
		ctx.RespondText("Argument <seconds> must be a number.")
		return
	}
	cfg.Snipers.GiveawaySniper.Sleep = delay
	cfg.Save()
}

func (c ConfigCommands) setupwebhooks(ctx *dgc.Ctx) {
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

func (c ConfigCommands) wall(ctx *dgc.Ctx) {
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
