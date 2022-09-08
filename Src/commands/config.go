package commands

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	//"github.com/electricbubble/go-toast"

	"github.com/LightningDev1/LB-Selfbot/config"
	"github.com/LightningDev1/LB-Selfbot/ipc"
	"github.com/LightningDev1/LB-Selfbot/utils"
	"github.com/LightningDev1/dgc"
	"github.com/LightningDev1/discordgo"
)

type ConfigCommands struct{}

func (c ConfigCommands) RegisterCommands(router *dgc.Router) {
	router.StartCategory("Config", "Change your config")
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
		Usage:       "[p]gsnipersleep <seconds>",
		Handler:     c.gsnipersleep,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "addgsniperbot",
		Description: "Add a bot to the custom bot list",
		Usage:       "[p]addgsniperbot <user>",
		Handler:     c.addgsniperbot,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "removegsniperbot",
		Description: "Remove a bot from the custom bot list",
		Usage:       "[p]removegsniperbot <user>",
		Handler:     c.removegsniperbot,
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
	router.RegisterCmd(&dgc.Command{
		Name:        "staffdetection",
		Description: "Turn staff detection on or off",
		Usage:       "[p]staffdetection <on/off>",
		Handler:     c.staffdetection,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "pingdetection",
		Description: "Turn ping detection on or off",
		Usage:       "[p]pingdetection <on/off>",
		Handler:     c.pingdetection,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "ticketnotifier",
		Description: "Turn ticket notifier on or off",
		Usage:       "[p]ticketnotifier <on/off>",
		Handler:     c.ticketnotifier,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "selfbotdetection",
		Description: "Turn selfbot detection on or off",
		Usage:       "[p]selfbotdetection <on/off>",
		Handler:     c.selfbotdetection,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "dmtypingnotifier",
		Description: "Turn DM typing detection on or off",
		Usage:       "[p]dmtypingnotifier <on/off>",
		Handler:     c.dmtypingnotifier,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "dmnotifier",
		Description: "Turn DM message detection on or off",
		Usage:       "[p]dmnotifier <on/off>",
		Handler:     c.dmnotifier,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "dmeditnotifier",
		Description: "Turn DM edit detection on or off",
		Usage:       "[p]dmeditnotifier <on/off>",
		Handler:     c.dmeditnotifier,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "dmdeletenotifier",
		Description: "Turn DM delete detection on or off",
		Usage:       "[p]dmdeletenotifier <on/off>",
		Handler:     c.dmdeletenotifier,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "messagedeletenotifier",
		Description: "Turn message delete detection on or off",
		Usage:       "[p]messagedeletenotifier <on/off>",
		Handler:     c.messagedeletenotifier,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "uitheme",
		Description: "Change the UI theme",
		Usage:       "[p]uitheme <light/dark>",
		Handler:     c.uitheme,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "hidetosystray",
		Description: "Turn hide to system tray on or off",
		Usage:       "[p]hidetosystray <on/off>",
		Handler:     c.hidetosystray,
	})
}

func (ConfigCommands) commandprefix(*dgc.Ctx) {
	cfg, _ := config.GetConfig()
	cfg.Selfbot.CommandPrefix = utils.Arguments["<prefix>"]
	_ = cfg.Save()
}

func (ConfigCommands) commandmode(ctx *dgc.Ctx) {
	cfg, _ := config.GetConfig()
	mode := utils.String.Capitalize(utils.Arguments["<codeblock/text/urlembed>"])
	if mode != "Codeblock" && mode != "Text" && mode != "Urlembed" {
		_ = ctx.RespondText("Invalid command mode.")
		return
	}
	cfg.Selfbot.CommandMode = strings.ReplaceAll(mode, "Urlembed", "URLEmbed")
	_ = cfg.Save()
}

func (ConfigCommands) title(ctx *dgc.Ctx) {
	cfg, _ := config.GetConfig()
	cfg.Selfbot.Title = ctx.Arguments.Raw()
	_ = cfg.Save()
}

func (ConfigCommands) footer(ctx *dgc.Ctx) {
	cfg, _ := config.GetConfig()
	cfg.Selfbot.Footer = ctx.Arguments.Raw()
	_ = cfg.Save()
}

func (ConfigCommands) delresponses(ctx *dgc.Ctx) {
	cfg, _ := config.GetConfig()
	toggle := strings.ToLower(utils.Arguments["<on/off>"])
	if toggle != "on" && toggle != "off" {
		_ = ctx.RespondText("Invalid value. Use `on` or `off`.")
		return
	}
	cfg.Selfbot.DeleteResponses = toggle == "on"
	_ = cfg.Save()
}

func (ConfigCommands) delresponsesdelay(ctx *dgc.Ctx) {
	cfg, _ := config.GetConfig()
	delay, err := strconv.Atoi(utils.Arguments["<seconds>"])
	if err != nil {
		_ = ctx.RespondText("Argument <seconds> must be a number.")
		return
	}
	cfg.Selfbot.DeleteResponsesDelay = delay
	_ = cfg.Save()
}

func (ConfigCommands) delcommands(ctx *dgc.Ctx) {
	cfg, _ := config.GetConfig()
	toggle := strings.ToLower(utils.Arguments["<on/off>"])
	if toggle != "on" && toggle != "off" {
		_ = ctx.RespondText("Invalid value. Use `on` or `off`.")
		return
	}
	cfg.Selfbot.DeleteCommands = toggle == "on"
	_ = cfg.Save()
}

func (ConfigCommands) notifications(ctx *dgc.Ctx) {
	cfg, _ := config.GetConfig()
	toggle := strings.ToLower(utils.Arguments["<on/off>"])
	if toggle != "on" && toggle != "off" {
		_ = ctx.RespondText("Invalid value. Use `on` or `off`.")
		return
	}
	cfg.Notifications.Enabled = toggle == "on"
	_ = cfg.Save()
}

func (ConfigCommands) notificationtype(ctx *dgc.Ctx) {
	cfg, _ := config.GetConfig()
	notificationType := utils.String.Capitalize(utils.Arguments["<default/custom>"])
	if notificationType != "Default" && notificationType != "Custom" {
		_ = ctx.RespondText("Invalid command mode.")
		return
	}
	cfg.Notifications.Type = notificationType
	_ = cfg.Save()
}

func (ConfigCommands) examplenotification(*dgc.Ctx) {
	/*_ = toast.Push("test message",
		toast.WithTitle("app title"),
		toast.WithMessage("test message 2"),
		toast.WithAudio(toast.LoopingAlarm4),
		toast.WithProtocolAction("Open Maps", "bingmaps:?q=sushi"),
	)*/
	ipc.SendNotification("This is an example notification.")
}

func (ConfigCommands) nitrosniper(ctx *dgc.Ctx) {
	cfg, _ := config.GetConfig()
	toggle := strings.ToLower(utils.Arguments["<on/off>"])
	if toggle != "on" && toggle != "off" {
		_ = ctx.RespondText("Invalid value. Use `on` or `off`.")
		return
	}
	cfg.Snipers.NitroSniper = toggle == "on"
	_ = cfg.Save()
}

func (ConfigCommands) giveawaysniper(ctx *dgc.Ctx) {
	cfg, _ := config.GetConfig()
	toggle := strings.ToLower(utils.Arguments["<on/off>"])
	if toggle != "on" && toggle != "off" {
		_ = ctx.RespondText("Invalid value. Use `on` or `off`.")
		return
	}
	cfg.Snipers.GiveawaySniper.Sniper = toggle == "on"
	_ = cfg.Save()
}

func (ConfigCommands) gsnipersleep(ctx *dgc.Ctx) {
	cfg, _ := config.GetConfig()
	delay, err := strconv.Atoi(utils.Arguments["<seconds>"])
	if err != nil {
		_ = ctx.RespondText("Argument <seconds> must be a number.")
		return
	}
	cfg.Snipers.GiveawaySniper.Sleep = delay
	_ = cfg.Save()
}

func (ConfigCommands) addgsniperbot(ctx *dgc.Ctx) {
	cfg, _ := config.GetConfig()
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		_ = ctx.RespondText("Could not find bot.")
		return
	}
	if utils.SliceContains(cfg.Snipers.GiveawaySniper.CustomBots, user.ID) {
		_ = ctx.RespondText(fmt.Sprintf("%s is already in the custom bot list.", user.String()))
		return
	}
	cfg.Snipers.GiveawaySniper.CustomBots = append(cfg.Snipers.GiveawaySniper.CustomBots, user.ID)
	_ = cfg.Save()
}

func (ConfigCommands) removegsniperbot(ctx *dgc.Ctx) {
	cfg, _ := config.GetConfig()
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		_ = ctx.RespondText("Could not find bot.")
		return
	}
	if !utils.SliceContains(cfg.Snipers.GiveawaySniper.CustomBots, user.ID) {
		_ = ctx.RespondText(fmt.Sprintf("%s is not in the custom bot list.", user.String()))
		return
	}
	cfg.Snipers.GiveawaySniper.CustomBots = utils.RemoveSliceItem(cfg.Snipers.GiveawaySniper.CustomBots, user.ID)
	_ = cfg.Save()
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
		_ = ctx.RespondText("Failed to create channel: " + err.Error())
		return
	}
	for channelName, configName := range channels {
		channel, err := ctx.Session.GuildChannelCreateComplex(serverID, discordgo.GuildChannelCreateData{
			Type:     discordgo.ChannelTypeGuildText,
			ParentID: category.ID,
			Name:     channelName,
		})
		if err != nil {
			_ = ctx.RespondText("Failed to create channel: " + err.Error())
			return
		}
		webhook, err := ctx.Session.WebhookCreate(channel.ID, channelName, "")
		if err != nil {
			_ = ctx.RespondText("Failed to create webhook: " + err.Error())
			return
		}
		f := reflect.Indirect(reflect.ValueOf(cfg.Webhooks)).FieldByName(configName)
		f.SetString(fmt.Sprintf("https://discord.com/api/webhooks/%s/%s", webhook.ID, webhook.Token))
	}
	_ = cfg.Save()
}

func (ConfigCommands) wall(*dgc.Ctx) {
	cfg, _ := config.GetConfig()
	cfg.Webhooks.NitroSniper = utils.Arguments["<webhook url>"]
	cfg.Webhooks.GiveawaySniper = utils.Arguments["<webhook url>"]
	cfg.Webhooks.PrivnoteSniper = utils.Arguments["<webhook url>"]
	cfg.Webhooks.SelfbotDetection = utils.Arguments["<webhook url>"]
	cfg.Webhooks.StaffDetection = utils.Arguments["<webhook url>"]
	cfg.Webhooks.Notifiers = utils.Arguments["<webhook url>"]
	cfg.Webhooks.PingDetection = utils.Arguments["<webhook url>"]
	cfg.Webhooks.Misc = utils.Arguments["<webhook url>"]
	_ = cfg.Save()
}

func (ConfigCommands) wnitro(*dgc.Ctx) {
	cfg, _ := config.GetConfig()
	cfg.Webhooks.NitroSniper = utils.Arguments["<webhook url>"]
	_ = cfg.Save()
}

func (ConfigCommands) wgiveaway(*dgc.Ctx) {
	cfg, _ := config.GetConfig()
	cfg.Webhooks.GiveawaySniper = utils.Arguments["<webhook url>"]
	_ = cfg.Save()
}

func (ConfigCommands) wprivnote(*dgc.Ctx) {
	cfg, _ := config.GetConfig()
	cfg.Webhooks.PrivnoteSniper = utils.Arguments["<webhook url>"]
	_ = cfg.Save()
}

func (ConfigCommands) wselfbotdetection(*dgc.Ctx) {
	cfg, _ := config.GetConfig()
	cfg.Webhooks.SelfbotDetection = utils.Arguments["<webhook url>"]
	_ = cfg.Save()
}

func (ConfigCommands) wstaffdetection(*dgc.Ctx) {
	cfg, _ := config.GetConfig()
	cfg.Webhooks.StaffDetection = utils.Arguments["<webhook url>"]
	_ = cfg.Save()
}

func (ConfigCommands) wnotifiers(*dgc.Ctx) {
	cfg, _ := config.GetConfig()
	cfg.Webhooks.Notifiers = utils.Arguments["<webhook url>"]
	_ = cfg.Save()
}

func (ConfigCommands) wpingdetection(*dgc.Ctx) {
	cfg, _ := config.GetConfig()
	cfg.Webhooks.PingDetection = utils.Arguments["<webhook url>"]
	_ = cfg.Save()
}

func (ConfigCommands) wmisc(*dgc.Ctx) {
	cfg, _ := config.GetConfig()
	cfg.Webhooks.Misc = utils.Arguments["<webhook url>"]
	_ = cfg.Save()
}

func (ConfigCommands) staffdetection(ctx *dgc.Ctx) {
	toggle := strings.ToLower(utils.Arguments["<on/off>"])
	if toggle != "on" && toggle != "off" {
		_ = ctx.RespondText("Invalid value. Use `on` or `off`.")
		return
	}
	cfg, _ := config.GetConfig()
	cfg.Events.DiscordStaffDetection = toggle == "on"
	_ = cfg.Save()
}

func (ConfigCommands) pingdetection(ctx *dgc.Ctx) {
	toggle := strings.ToLower(utils.Arguments["<on/off>"])
	if toggle != "on" && toggle != "off" {
		_ = ctx.RespondText("Invalid value. Use `on` or `off`.")
		return
	}
	cfg, _ := config.GetConfig()
	cfg.Events.GhostPingDetection = toggle == "on"
	_ = cfg.Save()
}

func (ConfigCommands) ticketnotifier(ctx *dgc.Ctx) {
	toggle := strings.ToLower(utils.Arguments["<on/off>"])
	if toggle != "on" && toggle != "off" {
		_ = ctx.RespondText("Invalid value. Use `on` or `off`.")
		return
	}
	cfg, _ := config.GetConfig()
	cfg.Events.TicketNotifier = toggle == "on"
	_ = cfg.Save()
}

func (ConfigCommands) selfbotdetection(ctx *dgc.Ctx) {
	toggle := strings.ToLower(utils.Arguments["<on/off>"])
	if toggle != "on" && toggle != "off" {
		_ = ctx.RespondText("Invalid value. Use `on` or `off`.")
		return
	}
	cfg, _ := config.GetConfig()
	cfg.Events.SelfbotDetection = toggle == "on"
	_ = cfg.Save()
}

func (ConfigCommands) dmtypingnotifier(ctx *dgc.Ctx) {
	toggle := strings.ToLower(utils.Arguments["<on/off>"])
	if toggle != "on" && toggle != "off" {
		_ = ctx.RespondText("Invalid value. Use `on` or `off`.")
		return
	}
	cfg, _ := config.GetConfig()
	cfg.Events.DMTypingNotifier = toggle == "on"
	_ = cfg.Save()
}

func (ConfigCommands) dmnotifier(ctx *dgc.Ctx) {
	toggle := strings.ToLower(utils.Arguments["<on/off>"])
	if toggle != "on" && toggle != "off" {
		_ = ctx.RespondText("Invalid value. Use `on` or `off`.")
		return
	}
	cfg, _ := config.GetConfig()
	cfg.Events.DMNotifier = toggle == "on"
	_ = cfg.Save()
}

func (ConfigCommands) dmdeletenotifier(ctx *dgc.Ctx) {
	toggle := strings.ToLower(utils.Arguments["<on/off>"])
	if toggle != "on" && toggle != "off" {
		_ = ctx.RespondText("Invalid value. Use `on` or `off`.")
		return
	}
	cfg, _ := config.GetConfig()
	cfg.Events.DMDeleteNotifier = toggle == "on"
	_ = cfg.Save()
}

func (ConfigCommands) dmeditnotifier(ctx *dgc.Ctx) {
	toggle := strings.ToLower(utils.Arguments["<on/off>"])
	if toggle != "on" && toggle != "off" {
		_ = ctx.RespondText("Invalid value. Use `on` or `off`.")
		return
	}
	cfg, _ := config.GetConfig()
	cfg.Events.DMEditNotifier = toggle == "on"
	_ = cfg.Save()
}

func (ConfigCommands) messagedeletenotifier(ctx *dgc.Ctx) {
	toggle := strings.ToLower(utils.Arguments["<on/off>"])
	if toggle != "on" && toggle != "off" {
		_ = ctx.RespondText("Invalid value. Use `on` or `off`.")
		return
	}
	cfg, _ := config.GetConfig()
	cfg.Events.MessageDeleteNotifier = toggle == "on"
	_ = cfg.Save()
}

func (ConfigCommands) uitheme(ctx *dgc.Ctx) {
	theme := strings.ToLower(utils.Arguments["<light/dark>"])
	if theme != "light" && theme != "dark" {
		_ = ctx.RespondText("Invalid value. Use `light` or `dark`.")
		return
	}
	cfg, _ := config.GetConfig()
	cfg.UI.Theme = utils.String.Capitalize(theme)
	_ = cfg.Save()
}

func (ConfigCommands) hidetosystray(ctx *dgc.Ctx) {
	toggle := strings.ToLower(utils.Arguments["<on/off>"])
	if toggle != "on" && toggle != "off" {
		_ = ctx.RespondText("Invalid value. Use `on` or `off`.")
		return
	}
	cfg, _ := config.GetConfig()
	cfg.UI.HideToSystemTray = toggle == "on"
	_ = cfg.Save()
}
