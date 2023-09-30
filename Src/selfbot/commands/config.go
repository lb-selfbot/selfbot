package commands

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/LightningDev1/LB-Selfbot/selfbot/config"
	"github.com/LightningDev1/LB-Selfbot/selfbot/embed"
	"github.com/LightningDev1/LB-Selfbot/selfbot/logging"
	"github.com/LightningDev1/LB-Selfbot/selfbot/utils"
	"github.com/LightningDev1/dgc"
	"github.com/LightningDev1/discordgo"
)

type ConfigCommands struct{}

func (c ConfigCommands) RegisterCommands(router *dgc.Router) {
	router.StartCategory("Config", "Change your config")

	router.RegisterCmd(&dgc.Command{
		Name:        "getcommandprefix",
		Description: "Get your current command prefix",
		Usage:       "[p]getcommandprefix",
		Flags:       []string{"DANGEROUS"},
		Handler:     c.getcommandprefix,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "commandprefix",
		Description: "Set your command prefix",
		Usage:       "[p]commandprefix <prefix>",
		Flags:       []string{"DANGEROUS"},
		Handler:     c.commandprefix,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "getcommandmode",
		Description: "Get your current command mode",
		Usage:       "[p]getcommandmode",
		Flags:       []string{"DANGEROUS"},
		Handler:     c.getcommandmode,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "commandmode",
		Description: "Set your command mode",
		Usage:       "[p]commandmode <codeblock/text/webhookembed/indent>",
		Flags:       []string{"DANGEROUS"},
		Handler:     c.commandmode,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "gettitle",
		Description: "Get your current title",
		Usage:       "[p]gettitle",
		Flags:       []string{"DANGEROUS"},
		Handler:     c.gettitle,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "title",
		Description: "Set your title",
		Usage:       "[p]title [title]",
		Flags:       []string{"DANGEROUS"},
		Handler:     c.title,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "getfooter",
		Description: "Get your current footer",
		Usage:       "[p]getfooter",
		Flags:       []string{"DANGEROUS"},
		Handler:     c.getfooter,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "footer",
		Description: "Set your footer",
		Usage:       "[p]footer [footer]",
		Flags:       []string{"DANGEROUS"},
		Handler:     c.footer,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "getdelresponses",
		Description: "Get your current delete responses setting",
		Usage:       "[p]getdelresponses",
		Flags:       []string{"DANGEROUS"},
		Handler:     c.getdelresponses,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "delresponses",
		Description: "Set if the bot should delete responses to commands",
		Usage:       "[p]delresponses <on/off>",
		Flags:       []string{"DANGEROUS"},
		Handler:     c.delresponses,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "getdelresponsesdelay",
		Description: "Get your current delete responses delay",
		Usage:       "[p]getdelresponsesdelay",
		Flags:       []string{"DANGEROUS"},
		Handler:     c.getdelresponsesdelay,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "delresponsesdelay",
		Description: "Set the delay before deleting responses",
		Usage:       "[p]delresponsesdelay <seconds>",
		Flags:       []string{"DANGEROUS"},
		Handler:     c.delresponsesdelay,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "getdelcommands",
		Description: "Get your current delete commands setting",
		Usage:       "[p]getdelcommands",
		Flags:       []string{"DANGEROUS"},
		Handler:     c.getdelcommands,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "delcommands",
		Description: "Set if the bot should delete commands",
		Usage:       "[p]delcommands <on/off>",
		Flags:       []string{"DANGEROUS"},
		Handler:     c.delcommands,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "getmobilefriendly",
		Description: "Get your current mobile friendly setting",
		Usage:       "[p]getmobilefriendly",
		Flags:       []string{"DANGEROUS"},
		Handler:     c.getmobilefriendly,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "mobilefriendly",
		Description: "Make the embeds mobile friendly",
		Usage:       "[p]mobilefriendly <on/off>",
		Flags:       []string{"DANGEROUS"},
		Handler:     c.mobilefriendly,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "getclientspoof",
		Description: "Get your current client spoof setting",
		Usage:       "[p]getclientspoof",
		Flags:       []string{"DANGEROUS"},
		Handler:     c.getclientspoof,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "clientspoof",
		Description: "Set your client spoof",
		Usage:       "[p]clientspoof <desktop/web/mobile>",
		Flags:       []string{"DANGEROUS"},
		Handler:     c.clientspoof,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "getnotifications",
		Description: "Get your current notifications setting",
		Usage:       "[p]getnotifications",
		Flags:       []string{"DANGEROUS"},
		Handler:     c.getnotifications,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "notifications",
		Description: "Turn toast notifications on or off",
		Usage:       "[p]notifications <on/off>",
		Flags:       []string{"DANGEROUS"},
		Handler:     c.notifications,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "examplenotification",
		Description: "Show an example notification",
		Usage:       "[p]examplenotification",
		Flags:       []string{"DANGEROUS"},
		Handler:     c.examplenotification,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "getnitrosniper",
		Description: "Get your current nitro sniper setting",
		Usage:       "[p]getnitrosniper",
		Flags:       []string{"DANGEROUS"},
		Handler:     c.getnitrosniper,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "nitrosniper",
		Description: "Turn the nitro sniper on or off",
		Usage:       "[p]nitrosniper <on/off>",
		Flags:       []string{"DANGEROUS"},
		Handler:     c.nitrosniper,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "getgiveawaysniper",
		Description: "Get your current giveaway sniper setting",
		Usage:       "[p]getgiveawaysniper",
		Flags:       []string{"DANGEROUS"},
		Handler:     c.getgiveawaysniper,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "giveawaysniper",
		Description: "Turn the giveaway sniper on or off",
		Usage:       "[p]giveawaysniper <on/off>",
		Flags:       []string{"DANGEROUS"},
		Handler:     c.giveawaysniper,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "getgsnipersleep",
		Description: "Get your current giveaway sniper sleep",
		Usage:       "[p]getgsnipersleep",
		Flags:       []string{"DANGEROUS"},
		Handler:     c.getgsnipersleep,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "gsnipersleep",
		Description: "Set the delay before the bot enters a giveaway",
		Usage:       "[p]gsnipersleep <seconds>",
		Flags:       []string{"DANGEROUS"},
		Handler:     c.gsnipersleep,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "getcustombots",
		Description: "Get a list of current custom bots",
		Usage:       "[p]getcustombots",
		Flags:       []string{"DANGEROUS"},
		Handler:     c.getcustombots,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "addcustombot",
		Description: "Add a custom bot for the giveaway sniper",
		Usage:       "[p]addcustombot <user>",
		Flags:       []string{"DANGEROUS"},
		Handler:     c.addcustombot,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "removecustombot",
		Description: "Remove a custom bot",
		Usage:       "[p]removecustombot <user>",
		Flags:       []string{"DANGEROUS"},
		Handler:     c.removecustombot,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "getsniperblacklist",
		Description: "Get the list of sniper blacklist entries",
		Usage:       "[p]getsniperblacklist",
		Flags:       []string{"DANGEROUS"},
		Handler:     c.getsniperblacklist,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "addsniperblacklist",
		Description: "Add to the sniper blacklist",
		Usage:       "[p]addsniperblacklist <text>",
		Flags:       []string{"DANGEROUS"},
		Handler:     c.addsniperblacklist,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "rmsniperblacklist",
		Description: "Remove from the sniper blacklist",
		Usage:       "[p]rmsniperblacklist <text>",
		Flags:       []string{"DANGEROUS"},
		Handler:     c.rmsniperblacklist,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "getsniperserverbl",
		Description: "Get the list of sniper server blacklist entries",
		Usage:       "[p]getsniperserverbl",
		Flags:       []string{"DANGEROUS"},
		Handler:     c.getsniperserverbl,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "addsniperserverbl",
		Description: "Add a server to the sniper blacklist",
		Usage:       "[p]addsniperserverbl <server id>",
		Flags:       []string{"DANGEROUS"},
		Handler:     c.addsniperserverbl,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "rmsniperserverbl",
		Description: "Remove a server from the sniper blacklist",
		Usage:       "[p]rmsniperserverbl <server id>",
		Flags:       []string{"DANGEROUS"},
		Handler:     c.rmsniperserverbl,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "getgsniperbots",
		Description: "Get the list of custom giveaway sniper bots",
		Usage:       "[p]getgsniperbots",
		Flags:       []string{"DANGEROUS"},
		Handler:     c.getgsniperbots,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "addgsniperbot",
		Description: "Add a bot to the custom bot list",
		Usage:       "[p]addgsniperbot <user>",
		Flags:       []string{"DANGEROUS"},
		Handler:     c.addgsniperbot,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "removegsniperbot",
		Description: "Remove a bot from the custom bot list",
		Usage:       "[p]removegsniperbot <user>",
		Flags:       []string{"DANGEROUS"},
		Handler:     c.removegsniperbot,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "getmultitokensniper",
		Description: "Get your current multi token sniper setting",
		Usage:       "[p]getmultitokensniper",
		Flags:       []string{"DANGEROUS"},
		Handler:     c.getmultitokensniper,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "multitokensniper",
		Description: "Turn the multi token sniper on or off",
		Usage:       "[p]multitokensniper <on/off>",
		Flags:       []string{"DANGEROUS"},
		Handler:     c.multitokensniper,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "getprivnotesniper",
		Description: "Get your current Privnote sniper setting",
		Usage:       "[p]getprivnotesniper",
		Flags:       []string{"DANGEROUS"},
		Handler:     c.getprivnotesniper,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "privnotesniper",
		Description: "Turn the Privnote sniper on or off",
		Usage:       "[p]privnotesniper <on/off>",
		Flags:       []string{"DANGEROUS"},
		Handler:     c.privnotesniper,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "getsnipertokens",
		Description: "Get list of multi token sniper tokens",
		Usage:       "[p]getsnipertokens",
		Flags:       []string{"DANGEROUS"},
		Handler:     c.getsnipertokens,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "addsnipertoken",
		Description: "Add a token to the multi token sniper",
		Usage:       "[p]addsnipertoken <token>",
		Flags:       []string{"DANGEROUS"},
		Handler:     c.addsnipertoken,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "rmsnipertoken",
		Description: "Remove a token from the multi token sniper",
		Usage:       "[p]rmsnipertoken <token>",
		Flags:       []string{"DANGEROUS"},
		Handler:     c.rmsnipertoken,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "setupwebhooks",
		Description: "Automatically set up all webhooks",
		Usage:       "[p]setupwebhooks [server id]",
		Flags:       []string{"DANGEROUS"},
		Handler:     c.setupwebhooks,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "wall",
		Description: "Set all webhooks at once",
		Usage:       "[p]wall [webhook url]",
		Flags:       []string{"DANGEROUS"},
		Handler:     c.wall,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "getwnitro",
		Description: "Get the current nitro sniper webhook",
		Usage:       "[p]getwnitro",
		Flags:       []string{"DANGEROUS"},
		Handler:     c.getwnitro,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "wnitro",
		Description: "Set the nitro sniper webhook",
		Usage:       "[p]wnitro [webhook url]",
		Flags:       []string{"DANGEROUS"},
		Handler:     c.wnitro,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "getwgiveaway",
		Description: "Get the current giveaway sniper webhook",
		Usage:       "[p]getwgiveaway",
		Flags:       []string{"DANGEROUS"},
		Handler:     c.getwgiveaway,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "wgiveaway",
		Description: "Set the giveaway sniper webhook",
		Usage:       "[p]wgiveaway [webhook url]",
		Flags:       []string{"DANGEROUS"},
		Handler:     c.wgiveaway,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "getwselfbotdetection",
		Description: "Get the current selfbot detection webhook",
		Usage:       "[p]getwselfbotdetection",
		Flags:       []string{"DANGEROUS"},
		Handler:     c.getwselfbotdetection,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "wselfbotdetection",
		Description: "Set the selfbot detection webhook",
		Usage:       "[p]wselfbotdetection [webhook url]",
		Flags:       []string{"DANGEROUS"},
		Handler:     c.wselfbotdetection,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "getwstaffdetection",
		Description: "Get the current staff detection webhook",
		Usage:       "[p]getwstaffdetection",
		Flags:       []string{"DANGEROUS"},
		Handler:     c.getwstaffdetection,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "wstaffdetection",
		Description: "Set the staff detection webhook",
		Usage:       "[p]wstaffdetection [webhook url]",
		Flags:       []string{"DANGEROUS"},
		Handler:     c.wstaffdetection,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "getwnotifiers",
		Description: "Get the current notifiers webhook",
		Usage:       "[p]getwnotifiers",
		Flags:       []string{"DANGEROUS"},
		Handler:     c.getwnotifiers,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "wnotifiers",
		Description: "Set the notifiers webhook",
		Usage:       "[p]wnotifiers [webhook url]",
		Flags:       []string{"DANGEROUS"},
		Handler:     c.wnotifiers,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "getwpingdetection",
		Description: "Get the current ping detection webhook",
		Usage:       "[p]getwpingdetection",
		Flags:       []string{"DANGEROUS"},
		Handler:     c.getwpingdetection,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "wpingdetection",
		Description: "Set the ping detection webhook",
		Usage:       "[p]wpingdetection [webhook url]",
		Flags:       []string{"DANGEROUS"},
		Handler:     c.wpingdetection,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "getwmisc",
		Description: "Get the current misc webhook",
		Usage:       "[p]getwmisc",
		Flags:       []string{"DANGEROUS"},
		Handler:     c.getwmisc,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "wmisc",
		Description: "Set the misc webhook",
		Usage:       "[p]wmisc [webhook url]",
		Flags:       []string{"DANGEROUS"},
		Handler:     c.wmisc,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "getwsessions",
		Description: "Get the current sessions webhook",
		Usage:       "[p]getwsessions",
		Flags:       []string{"DANGEROUS"},
		Handler:     c.getwsessions,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "wsessions",
		Description: "Set the sessions webhook",
		Usage:       "[p]wsessions [webhook url]",
		Flags:       []string{"DANGEROUS"},
		Handler:     c.wsessions,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "getbannotifier",
		Description: "Get the current ban notifier value",
		Usage:       "[p]getbannotifier",
		Flags:       []string{"DANGEROUS"},
		Handler:     c.getbannotifier,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "bannotifier",
		Description: "Turn ban notifier on or off",
		Usage:       "[p]bannotifier <on/off>",
		Flags:       []string{"DANGEROUS"},
		Handler:     c.bannotifier,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "getstaffdetection",
		Description: "Get the current staff detection value",
		Usage:       "[p]getstaffdetection",
		Flags:       []string{"DANGEROUS"},
		Handler:     c.getstaffdetection,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "staffdetection",
		Description: "Turn staff detection on or off",
		Usage:       "[p]staffdetection <on/off>",
		Flags:       []string{"DANGEROUS"},
		Handler:     c.staffdetection,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "getpingdetection",
		Description: "Get the current ping detection value",
		Usage:       "[p]getpingdetection",
		Flags:       []string{"DANGEROUS"},
		Handler:     c.getpingdetection,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "pingdetection",
		Description: "Turn ping detection on or off",
		Usage:       "[p]pingdetection <on/off>",
		Flags:       []string{"DANGEROUS"},
		Handler:     c.pingdetection,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "getticketnotifier",
		Description: "Get the current ticket notifier value",
		Usage:       "[p]getticketnotifier",
		Flags:       []string{"DANGEROUS"},
		Handler:     c.getticketnotifier,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "ticketnotifier",
		Description: "Turn ticket notifier on or off",
		Usage:       "[p]ticketnotifier <on/off>",
		Flags:       []string{"DANGEROUS"},
		Handler:     c.ticketnotifier,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "getsessionconnect",
		Description: "Get the current session connected value",
		Usage:       "[p]getsessionconnect",
		Flags:       []string{"DANGEROUS"},
		Handler:     c.getsessionconnect,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "sessionconnect",
		Description: "Turn session connected notifier on or off",
		Usage:       "[p]sessionconnect <on/off>",
		Flags:       []string{"DANGEROUS"},
		Handler:     c.sessionconnect,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "getsessiondisconnect",
		Description: "Get the current session disconnected value",
		Usage:       "[p]getsessiondisconnect",
		Flags:       []string{"DANGEROUS"},
		Handler:     c.getsessiondisconnect,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "sessiondisconnect",
		Description: "Turn session disconnected notifier on or off",
		Usage:       "[p]sessiondisconnect <on/off>",
		Flags:       []string{"DANGEROUS"},
		Handler:     c.sessiondisconnect,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "getselfbotdetection",
		Description: "Get the current selfbot detection value",
		Usage:       "[p]getselfbotdetection",
		Flags:       []string{"DANGEROUS"},
		Handler:     c.getselfbotdetection,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "selfbotdetection",
		Description: "Turn selfbot detection on or off",
		Usage:       "[p]selfbotdetection <on/off>",
		Flags:       []string{"DANGEROUS"},
		Handler:     c.selfbotdetection,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "getdmtypingnotifier",
		Description: "Get the current DM typing notifier value",
		Usage:       "[p]getdmtypingnotifier",
		Flags:       []string{"DANGEROUS"},
		Handler:     c.getdmtypingnotifier,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "dmtypingnotifier",
		Description: "Turn the DM typing notifier on or off",
		Usage:       "[p]dmtypingnotifier <on/off>",
		Flags:       []string{"DANGEROUS"},
		Handler:     c.dmtypingnotifier,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "getdmnotifier",
		Description: "Get the current DM notifier value",
		Usage:       "[p]getdmnotifier",
		Flags:       []string{"DANGEROUS"},
		Handler:     c.getdmnotifier,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "dmnotifier",
		Description: "Turn the DM notifier on or off",
		Usage:       "[p]dmnotifier <on/off>",
		Flags:       []string{"DANGEROUS"},
		Handler:     c.dmnotifier,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "getdmeditnotifier",
		Description: "Get the current DM edit notifier value",
		Usage:       "[p]getdmeditnotifier",
		Flags:       []string{"DANGEROUS"},
		Handler:     c.getdmeditnotifier,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "dmeditnotifier",
		Description: "Turn the DM edit notifier on or off",
		Usage:       "[p]dmeditnotifier <on/off>",
		Flags:       []string{"DANGEROUS"},
		Handler:     c.dmeditnotifier,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "getdmdeletenotifier",
		Description: "Get the current DM delete notifier value",
		Usage:       "[p]getdmdeletenotifier",
		Flags:       []string{"DANGEROUS"},
		Handler:     c.getdmdeletenotifier,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "dmdeletenotifier",
		Description: "Turn the DM delete notifier on or off",
		Usage:       "[p]dmdeletenotifier <on/off>",
		Flags:       []string{"DANGEROUS"},
		Handler:     c.dmdeletenotifier,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "getmsgdeletenotifier",
		Description: "Get the current message delete notifier value",
		Usage:       "[p]getmsgdeletenotifier",
		Flags:       []string{"DANGEROUS"},
		Handler:     c.getmsgdeletenotifier,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "messagedeletenotifier",
		Description: "Turn the message delete notifier on or off",
		Usage:       "[p]messagedeletenotifier <on/off>",
		Flags:       []string{"DANGEROUS"},
		Handler:     c.messagedeletenotifier,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "getuitheme",
		Description: "Get the current UI theme",
		Usage:       "[p]getuitheme",
		Flags:       []string{"DANGEROUS"},
		Handler:     c.getuitheme,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "uitheme",
		Description: "Change the UI theme",
		Usage:       "[p]uitheme <light/dark>",
		Flags:       []string{"DANGEROUS"},
		Handler:     c.uitheme,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "gethidetosystray",
		Description: "Get the current hide to system tray value",
		Usage:       "[p]gethidetosystray",
		Flags:       []string{"DANGEROUS"},
		Handler:     c.gethidetosystray,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "hidetosystray",
		Description: "Turn hide to system tray on or off",
		Usage:       "[p]hidetosystray <on/off>",
		Flags:       []string{"DANGEROUS"},
		Handler:     c.hidetosystray,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "getrunonstartup",
		Description: "Get the current run on startup value",
		Usage:       "[p]getrunonstartup",
		Flags:       []string{"DANGEROUS"},
		Handler:     c.getrunonstartup,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "runonstartup",
		Description: "Turn run on startup on or off",
		Usage:       "[p]runonstartup <on/off>",
		Flags:       []string{"DANGEROUS"},
		Handler:     c.runonstartup,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "getrunsbonstartup",
		Description: "Get the current run selfbot on startup value",
		Usage:       "[p]getrunsbonstartup",
		Flags:       []string{"DANGEROUS"},
		Handler:     c.getrunsbonstartup,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "runsbonstartup",
		Description: "Turn run selfbot on startup on or off",
		Usage:       "[p]runsbonstartup <on/off>",
		Flags:       []string{"DANGEROUS"},
		Handler:     c.runsbonstartup,
	})
}

func (ConfigCommands) getcommandprefix(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	embed.NewEmbed().
		AddField("Command Prefix", cfg.Selfbot.CommandPrefix).
		Send(ctx)
}

func (ConfigCommands) commandprefix(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	cfg.Selfbot.CommandPrefix = utils.Arguments["<prefix>"]
	_ = cfg.Save()

	embed.NewEmbed().
		SetDescription("Command prefix set to " + utils.Arguments["<prefix>"]).
		Send(ctx)
}

func (ConfigCommands) getcommandmode(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	embed.NewEmbed().
		AddField("Command Mode", cfg.Selfbot.CommandMode).
		Send(ctx)
}

func (ConfigCommands) commandmode(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	mode := utils.String.Capitalize(utils.Arguments["<codeblock/text/webhookembed/indent>"])
	if !utils.SliceContains([]string{"Codeblock", "Text", "Webhookembed", "Indent"}, mode) {
		_ = ctx.RespondText("Invalid command mode.")
		return
	}

	cfg.Selfbot.CommandMode = strings.ReplaceAll(mode, "Webhookembed", "WebhookEmbed")
	_ = cfg.Save()

	embed.NewEmbed().
		SetDescription("Command mode set to " + cfg.Selfbot.CommandMode).
		Send(ctx)
}

func (ConfigCommands) gettitle(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	embed.NewEmbed().
		AddField("Title", cfg.Selfbot.Title).
		Send(ctx)
}

func (ConfigCommands) title(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	cfg.Selfbot.Title = ctx.Arguments.Raw()
	_ = cfg.Save()

	embed.NewEmbed().
		SetDescription("Title set to " + cfg.Selfbot.CommandMode).
		Send(ctx)
}

func (ConfigCommands) getfooter(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	embed.NewEmbed().
		AddField("Footer", cfg.Selfbot.Footer).
		Send(ctx)
}

func (ConfigCommands) footer(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	cfg.Selfbot.Footer = ctx.Arguments.Raw()
	_ = cfg.Save()

	embed.NewEmbed().
		SetDescription("Footer set to " + cfg.Selfbot.Footer).
		Send(ctx)
}

func (ConfigCommands) getdelresponses(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	embed.NewEmbed().
		AddField("Delete Responses", utils.Misc.FormatBoolOnOff(cfg.Selfbot.DeleteResponses)).
		Send(ctx)
}

func (ConfigCommands) delresponses(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	toggle := strings.ToLower(utils.Arguments["<on/off>"])
	if toggle != "on" && toggle != "off" {
		_ = ctx.RespondText("Invalid value. Use `on` or `off`.")
		return
	}

	cfg.Selfbot.DeleteResponses = toggle == "on"
	_ = cfg.Save()

	embed.NewEmbed().
		SetDescription("Turned " + toggle + " delete responses").
		Send(ctx)
}

func (ConfigCommands) getdelresponsesdelay(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	embed.NewEmbed().
		AddField("Delete Responses Delay", fmt.Sprintf("%d seconds", cfg.Selfbot.DeleteResponsesDelay)).
		Send(ctx)
}

func (ConfigCommands) delresponsesdelay(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	delay, err := strconv.Atoi(utils.Arguments["<seconds>"])
	if err != nil {
		_ = ctx.RespondText("Argument <seconds> must be a number.")
		return
	}

	cfg.Selfbot.DeleteResponsesDelay = delay
	_ = cfg.Save()

	embed.NewEmbed().
		SetDescription(fmt.Sprintf("Delete responses delay set to %d seconds", cfg.Selfbot.DeleteResponsesDelay)).
		Send(ctx)
}

func (ConfigCommands) getdelcommands(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	embed.NewEmbed().
		AddField("Delete Commands", utils.Misc.FormatBoolOnOff(cfg.Selfbot.DeleteCommands)).
		Send(ctx)
}

func (ConfigCommands) delcommands(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	toggle := strings.ToLower(utils.Arguments["<on/off>"])
	if toggle != "on" && toggle != "off" {
		_ = ctx.RespondText("Invalid value. Use `on` or `off`.")
		return
	}

	cfg.Selfbot.DeleteCommands = toggle == "on"
	_ = cfg.Save()

	embed.NewEmbed().
		SetDescription("Turned " + toggle + " delete commands").
		Send(ctx)
}

func (ConfigCommands) getmobilefriendly(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	embed.NewEmbed().
		AddField("Mobile Friendly", utils.Misc.FormatBoolOnOff(cfg.Selfbot.MobileFriendly)).
		Send(ctx)
}

func (ConfigCommands) mobilefriendly(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	toggle := strings.ToLower(utils.Arguments["<on/off>"])
	if toggle != "on" && toggle != "off" {
		_ = ctx.RespondText("Invalid value. Use `on` or `off`.")
		return
	}

	cfg.Selfbot.MobileFriendly = toggle == "on"
	_ = cfg.Save()

	embed.NewEmbed().
		SetDescription("Turned " + toggle + " mobile friendly").
		Send(ctx)
}

func (ConfigCommands) getclientspoof(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	embed.NewEmbed().
		AddField("Client", cfg.Selfbot.Client).
		Send(ctx)
}

func (ConfigCommands) clientspoof(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	client := utils.String.Capitalize(utils.Arguments["<desktop/web/mobile>"])
	if !utils.SliceContains([]string{"Desktop", "Web", "Mobile"}, client) {
		_ = ctx.RespondText("Invalid client.")
		return
	}

	cfg.Selfbot.CommandMode = client
	_ = cfg.Save()

	embed.NewEmbed().
		SetDescription("Client set to " + cfg.Selfbot.CommandMode).
		Send(ctx)
}

func (ConfigCommands) getnotifications(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	embed.NewEmbed().
		AddField("Notifications", utils.Misc.FormatBoolOnOff(cfg.Notifications.Enabled)).
		Send(ctx)
}

func (ConfigCommands) notifications(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	toggle := strings.ToLower(utils.Arguments["<on/off>"])
	if toggle != "on" && toggle != "off" {
		_ = ctx.RespondText("Invalid value. Use `on` or `off`.")
		return
	}

	cfg.Notifications.Enabled = toggle == "on"
	_ = cfg.Save()

	embed.NewEmbed().
		SetDescription("Turned " + toggle + " notifications").
		Send(ctx)
}

func (ConfigCommands) examplenotification(ctx *dgc.Ctx) {
	logging.UI.Notification("This is an example notification.")
}

func (ConfigCommands) getnitrosniper(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	embed.NewEmbed().
		AddField("Nitro Sniper", utils.Misc.FormatBoolOnOff(cfg.Snipers.NitroSniper)).
		Send(ctx)
}

func (ConfigCommands) nitrosniper(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	toggle := strings.ToLower(utils.Arguments["<on/off>"])
	if toggle != "on" && toggle != "off" {
		_ = ctx.RespondText("Invalid value. Use `on` or `off`.")
		return
	}

	cfg.Snipers.NitroSniper = toggle == "on"
	_ = cfg.Save()

	embed.NewEmbed().
		SetDescription("Turned " + toggle + " nitro sniper").
		Send(ctx)
}

func (ConfigCommands) getgiveawaysniper(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	embed.NewEmbed().
		AddField("Giveaway Sniper", utils.Misc.FormatBoolOnOff(cfg.Snipers.GiveawaySniper.Enabled)).
		Send(ctx)
}

func (ConfigCommands) giveawaysniper(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	toggle := strings.ToLower(utils.Arguments["<on/off>"])
	if toggle != "on" && toggle != "off" {
		_ = ctx.RespondText("Invalid value. Use `on` or `off`.")
		return
	}

	cfg.Snipers.GiveawaySniper.Enabled = toggle == "on"
	_ = cfg.Save()

	embed.NewEmbed().
		SetDescription("Turned " + toggle + " giveaway sniper").
		Send(ctx)
}

func (ConfigCommands) getgsnipersleep(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	embed.NewEmbed().
		AddField("Giveaway Sniper Sleep", fmt.Sprintf("%d seconds", cfg.Snipers.GiveawaySniper.Sleep)).
		Send(ctx)
}

func (ConfigCommands) gsnipersleep(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	delay, err := strconv.Atoi(utils.Arguments["<seconds>"])
	if err != nil {
		_ = ctx.RespondText("Argument <seconds> must be a number.")
		return
	}

	cfg.Snipers.GiveawaySniper.Sleep = delay
	_ = cfg.Save()

	embed.NewEmbed().
		SetDescription(fmt.Sprintf("Giveaway sniper sleep set to %d seconds", cfg.Snipers.GiveawaySniper.Sleep)).
		Send(ctx)
}

func (ConfigCommands) getcustombots(ctx *dgc.Ctx) {
	cfg, _ := config.GetConfig()

	botString := ""
	for _, userID := range cfg.Snipers.GiveawaySniper.CustomBots {
		user, ok := utils.Discord.GetUser(ctx.Session, userID)
		if !ok {
			botString += fmt.Sprintf("%s - %s\n", "Unknown", userID)
			continue
		}

		botString += fmt.Sprintf("%s - %s\n", user.Username, userID)
	}

	embed.NewEmbed().
		AddField("Custom Bots", botString).
		Send(ctx)
}

func (ConfigCommands) addcustombot(ctx *dgc.Ctx) {
	user, ok := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !ok {
		_ = ctx.RespondText("Could not find bot.")
		return
	}

	cfg, _ := config.GetConfig()

	if utils.SliceContains(cfg.Snipers.GiveawaySniper.CustomBots, user.ID) {
		_ = ctx.RespondText("Bot is already in the custom bots list.")
		return
	}

	cfg.Snipers.GiveawaySniper.CustomBots = append(cfg.Snipers.GiveawaySniper.CustomBots, user.ID)

	_ = cfg.Save()
}

func (ConfigCommands) removecustombot(ctx *dgc.Ctx) {
	user, ok := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !ok {
		_ = ctx.RespondText("Could not find bot.")
		return
	}

	cfg, _ := config.GetConfig()

	if !utils.SliceContains(cfg.Snipers.GiveawaySniper.CustomBots, user.ID) {
		_ = ctx.RespondText("Bot is already not in the custom bots list.")
		return
	}

	cfg.Snipers.GiveawaySniper.CustomBots = utils.RemoveSliceItem(cfg.Snipers.GiveawaySniper.CustomBots, user.ID)

	_ = cfg.Save()
}

func (ConfigCommands) getsniperblacklist(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	embed.NewEmbed().
		AddField("Giveaway Sniper Blacklist", strings.Join(cfg.Snipers.Blacklist, ", ")).
		Send(ctx)
}

func (ConfigCommands) addsniperblacklist(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	entry := utils.String.RemoveWords(ctx.Arguments.Raw(), 1)

	if utils.SliceContains(cfg.Snipers.Blacklist, entry) {
		_ = ctx.RespondText("That is already in the blacklist!")
		return
	}

	cfg.Snipers.Blacklist = append(cfg.Snipers.Blacklist, entry)
	_ = cfg.Save()

	embed.NewEmbed().
		SetDescription("Added " + entry + " to the sniper blacklist").
		Send(ctx)
}

func (ConfigCommands) rmsniperblacklist(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	entry := utils.String.RemoveWords(ctx.Arguments.Raw(), 1)

	if !utils.SliceContains(cfg.Snipers.Blacklist, entry) {
		_ = ctx.RespondText("That is already not in the blacklist!")
		return
	}

	cfg.Snipers.Blacklist = utils.RemoveSliceItem(cfg.Snipers.Blacklist, entry)
	_ = cfg.Save()

	embed.NewEmbed().
		SetDescription("Removed " + entry + " from the sniper blacklist").
		Send(ctx)
}

func (ConfigCommands) getsniperserverbl(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	embed.NewEmbed().
		AddField("Giveaway Sniper Server Blacklist", strings.Join(cfg.Snipers.ServerBlacklist, ", ")).
		Send(ctx)
}

func (ConfigCommands) addsniperserverbl(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	guild, err := utils.Discord.GetGuild(ctx.Session, utils.Misc.Or(utils.Arguments["<server id>"], ctx.Event.GuildID))
	if err != nil {
		_ = ctx.RespondText("Could not find server.")
		return
	}

	if utils.SliceContains(cfg.Snipers.ServerBlacklist, guild.ID) {
		_ = ctx.RespondText("Server is already in the blacklist!")
		return
	}

	cfg.Snipers.ServerBlacklist = append(cfg.Snipers.ServerBlacklist, guild.ID)
	_ = cfg.Save()

	embed.NewEmbed().
		SetDescription(fmt.Sprintf("Added %s to the server blacklist.", guild.Name)).
		Send(ctx)
}

func (ConfigCommands) rmsniperserverbl(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	guild, err := utils.Discord.GetGuild(ctx.Session, utils.Misc.Or(utils.Arguments["<server id>"], ctx.Event.GuildID))
	if err != nil {
		_ = ctx.RespondText("Could not find server.")
		return
	}

	if !utils.SliceContains(cfg.Snipers.ServerBlacklist, guild.ID) {
		_ = ctx.RespondText("Server is already not in the blacklist!")
		return
	}

	cfg.Snipers.ServerBlacklist = utils.RemoveSliceItem(cfg.Snipers.ServerBlacklist, guild.ID)
	_ = cfg.Save()

	embed.NewEmbed().
		SetDescription(fmt.Sprintf("Removed %s from the server blacklist.", guild.Name)).
		Send(ctx)
}

func (ConfigCommands) getgsniperbots(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	embed.NewEmbed().
		AddField("Giveaway Sniper Custom Bots", strings.Join(cfg.Snipers.GiveawaySniper.CustomBots, ", ")).
		Send(ctx)
}

func (ConfigCommands) addgsniperbot(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

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

	embed.NewEmbed().
		SetDescription(fmt.Sprintf("Added %s to the custom giveaway sniper bot list.", user.String())).
		Send(ctx)
}

func (ConfigCommands) removegsniperbot(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

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

	embed.NewEmbed().
		SetDescription(fmt.Sprintf("Removed %s from the custom giveaway sniper bot list.", user.String())).
		Send(ctx)
}

func (ConfigCommands) getprivnotesniper(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	embed.NewEmbed().
		AddField("Privnote Sniper", utils.Misc.FormatBoolOnOff(cfg.Snipers.PrivnoteSniper)).
		Send(ctx)
}

func (ConfigCommands) privnotesniper(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	toggle := strings.ToLower(utils.Arguments["<on/off>"])
	if toggle != "on" && toggle != "off" {
		_ = ctx.RespondText("Invalid value. Use `on` or `off`.")
		return
	}

	cfg.Snipers.PrivnoteSniper = toggle == "on"
	_ = cfg.Save()

	embed.NewEmbed().
		SetDescription("Turned " + toggle + " Privnote sniper").
		Send(ctx)
}

func (ConfigCommands) getmultitokensniper(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	embed.NewEmbed().
		AddField("Multi Token Sniper", utils.Misc.FormatBoolOnOff(cfg.Snipers.MultiTokenSniper)).
		Send(ctx)
}

func (ConfigCommands) multitokensniper(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	toggle := strings.ToLower(utils.Arguments["<on/off>"])
	if toggle != "on" && toggle != "off" {
		_ = ctx.RespondText("Invalid value. Use `on` or `off`.")
		return
	}

	cfg.Snipers.MultiTokenSniper = toggle == "on"
	_ = cfg.Save()

	embed.NewEmbed().
		SetDescription("Turned " + toggle + " multi token sniper").
		Send(ctx)
}

func (ConfigCommands) getsnipertokens(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	embed.NewEmbed().
		AddField("Sniper Tokens", strings.Join(cfg.Snipers.Tokens, ", ")).
		Send(ctx)
}

func (ConfigCommands) addsnipertoken(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	entry := utils.String.RemoveWords(ctx.Arguments.Raw(), 1)

	if utils.SliceContains(cfg.Snipers.Tokens, entry) {
		_ = ctx.RespondText("That token is already in the multi token sniper!")
		return
	}

	cfg.Snipers.Tokens = append(cfg.Snipers.Tokens, entry)
	_ = cfg.Save()

	embed.NewEmbed().
		SetDescription("Added token the multi token sniper").
		Send(ctx)
}

func (ConfigCommands) rmsnipertoken(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	entry := utils.String.RemoveWords(ctx.Arguments.Raw(), 1)

	if !utils.SliceContains(cfg.Snipers.Tokens, entry) {
		_ = ctx.RespondText("That token is already not in the multi token sniper!")
		return
	}

	cfg.Snipers.Tokens = utils.RemoveSliceItem(cfg.Snipers.Tokens, entry)
	_ = cfg.Save()

	embed.NewEmbed().
		SetDescription("Added token the multi token sniper").
		Send(ctx)
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
}

func (ConfigCommands) wall(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	cfg.Webhooks.NitroSniper = utils.Arguments["[webhook url]"]
	cfg.Webhooks.GiveawaySniper = utils.Arguments["[webhook url]"]
	cfg.Webhooks.SelfbotDetection = utils.Arguments["[webhook url]"]
	cfg.Webhooks.StaffDetection = utils.Arguments["[webhook url]"]
	cfg.Webhooks.Notifiers = utils.Arguments["[webhook url]"]
	cfg.Webhooks.PingDetection = utils.Arguments["[webhook url]"]
	cfg.Webhooks.Misc = utils.Arguments["[webhook url]"]
	cfg.Webhooks.Sessions = utils.Arguments["[webhook url]"]

	_ = cfg.Save()

	embed.NewEmbed().
		SetDescription("Set all webhooks to " + utils.Arguments["[webhook url]"]).
		Send(ctx)
}

func (ConfigCommands) getwnitro(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	embed.NewEmbed().
		AddField("Nitro Sniper Webhook", cfg.Webhooks.NitroSniper).
		Send(ctx)
}

func (ConfigCommands) wnitro(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	cfg.Webhooks.NitroSniper = utils.Arguments["[webhook url]"]
	_ = cfg.Save()

	embed.NewEmbed().
		SetDescription("Set nitro sniper webhook to " + utils.Arguments["[webhook url]"]).
		Send(ctx)
}

func (ConfigCommands) getwgiveaway(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	embed.NewEmbed().
		AddField("Giveaway Sniper Webhook", cfg.Webhooks.GiveawaySniper).
		Send(ctx)
}

func (ConfigCommands) wgiveaway(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	cfg.Webhooks.GiveawaySniper = utils.Arguments["[webhook url]"]
	_ = cfg.Save()

	embed.NewEmbed().
		SetDescription("Set giveaway sniper webhook to " + utils.Arguments["[webhook url]"]).
		Send(ctx)
}

func (ConfigCommands) getwselfbotdetection(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	embed.NewEmbed().
		AddField("Selfbot Detection Webhook", cfg.Webhooks.SelfbotDetection).
		Send(ctx)
}

func (ConfigCommands) wselfbotdetection(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	cfg.Webhooks.SelfbotDetection = utils.Arguments["[webhook url]"]
	_ = cfg.Save()

	embed.NewEmbed().
		SetDescription("Set selfbot detection webhook to " + utils.Arguments["[webhook url]"]).
		Send(ctx)
}

func (ConfigCommands) getwstaffdetection(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	embed.NewEmbed().
		AddField("Staff Detection Webhook", cfg.Webhooks.StaffDetection).
		Send(ctx)
}

func (ConfigCommands) wstaffdetection(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	cfg.Webhooks.StaffDetection = utils.Arguments["[webhook url]"]
	_ = cfg.Save()

	embed.NewEmbed().
		SetDescription("Set staff detection webhook to " + utils.Arguments["[webhook url]"]).
		Send(ctx)
}

func (ConfigCommands) getwnotifiers(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	embed.NewEmbed().
		AddField("Notifiers Webhook", cfg.Webhooks.Notifiers).
		Send(ctx)
}

func (ConfigCommands) wnotifiers(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	cfg.Webhooks.Notifiers = utils.Arguments["[webhook url]"]
	_ = cfg.Save()

	embed.NewEmbed().
		SetDescription("Set notifiers webhook to " + utils.Arguments["[webhook url]"]).
		Send(ctx)
}

func (ConfigCommands) getwpingdetection(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	embed.NewEmbed().
		AddField("Ping Detection Webhook", cfg.Webhooks.PingDetection).
		Send(ctx)
}

func (ConfigCommands) wpingdetection(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	cfg.Webhooks.PingDetection = utils.Arguments["[webhook url]"]
	_ = cfg.Save()

	embed.NewEmbed().
		SetDescription("Set ping detection webhook to " + utils.Arguments["[webhook url]"]).
		Send(ctx)
}

func (ConfigCommands) getwmisc(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	embed.NewEmbed().
		AddField("Misc Webhook", cfg.Webhooks.Misc).
		Send(ctx)
}

func (ConfigCommands) wmisc(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	cfg.Webhooks.Misc = utils.Arguments["[webhook url]"]
	_ = cfg.Save()

	embed.NewEmbed().
		SetDescription("Set misc webhook to " + utils.Arguments["[webhook url]"]).
		Send(ctx)
}

func (ConfigCommands) getwsessions(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	embed.NewEmbed().
		AddField("Sessions Webhook", cfg.Webhooks.Sessions).
		Send(ctx)
}

func (ConfigCommands) wsessions(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	cfg.Webhooks.Sessions = utils.Arguments["[webhook url]"]
	_ = cfg.Save()

	embed.NewEmbed().
		SetDescription("Set sessions webhook to " + utils.Arguments["[webhook url]"]).
		Send(ctx)
}

func (ConfigCommands) getbannotifier(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	embed.NewEmbed().
		AddField("Ban Notifier", utils.Misc.FormatBoolOnOff(cfg.Events.BanNotifier)).
		Send(ctx)
}

func (ConfigCommands) bannotifier(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	toggle := strings.ToLower(utils.Arguments["<on/off>"])
	if toggle != "on" && toggle != "off" {
		_ = ctx.RespondText("Invalid value. Use `on` or `off`.")
		return
	}

	cfg.Events.BanNotifier = toggle == "on"
	_ = cfg.Save()

	embed.NewEmbed().
		SetDescription("Turned " + toggle + " ban notifier").
		Send(ctx)
}

func (ConfigCommands) getstaffdetection(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	embed.NewEmbed().
		AddField("Staff Detection", utils.Misc.FormatBoolOnOff(cfg.Events.DiscordStaffDetection)).
		Send(ctx)
}

func (ConfigCommands) staffdetection(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	toggle := strings.ToLower(utils.Arguments["<on/off>"])
	if toggle != "on" && toggle != "off" {
		_ = ctx.RespondText("Invalid value. Use `on` or `off`.")
		return
	}

	cfg.Events.DiscordStaffDetection = toggle == "on"
	_ = cfg.Save()

	embed.NewEmbed().
		SetDescription("Turned " + toggle + " staff detection").
		Send(ctx)
}

func (ConfigCommands) getpingdetection(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	embed.NewEmbed().
		AddField("Ping Detection", utils.Misc.FormatBoolOnOff(cfg.Events.GhostPingDetection)).
		Send(ctx)
}

func (ConfigCommands) pingdetection(ctx *dgc.Ctx) {
	toggle := strings.ToLower(utils.Arguments["<on/off>"])
	if toggle != "on" && toggle != "off" {
		_ = ctx.RespondText("Invalid value. Use `on` or `off`.")
		return
	}
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	cfg.Events.GhostPingDetection = toggle == "on"
	_ = cfg.Save()

	embed.NewEmbed().
		SetDescription("Turned " + toggle + " ping detection").
		Send(ctx)
}

func (ConfigCommands) getticketnotifier(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	embed.NewEmbed().
		AddField("Ticket Notifier", utils.Misc.FormatBoolOnOff(cfg.Events.TicketNotifier)).
		Send(ctx)
}

func (ConfigCommands) ticketnotifier(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	toggle := strings.ToLower(utils.Arguments["<on/off>"])
	if toggle != "on" && toggle != "off" {
		_ = ctx.RespondText("Invalid value. Use `on` or `off`.")
		return
	}

	cfg.Events.TicketNotifier = toggle == "on"
	_ = cfg.Save()

	embed.NewEmbed().
		SetDescription("Turned " + toggle + " ticket notifier").
		Send(ctx)
}

func (ConfigCommands) getsessionconnect(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	embed.NewEmbed().
		AddField("Session Connected Notifier", utils.Misc.FormatBoolOnOff(cfg.Events.SessionConnected)).
		Send(ctx)
}

func (ConfigCommands) sessionconnect(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	toggle := strings.ToLower(utils.Arguments["<on/off>"])
	if toggle != "on" && toggle != "off" {
		_ = ctx.RespondText("Invalid value. Use `on` or `off`.")
		return
	}

	cfg.Events.SessionConnected = toggle == "on"
	_ = cfg.Save()

	embed.NewEmbed().
		SetDescription("Turned " + toggle + " session connected notifier").
		Send(ctx)
}

func (ConfigCommands) getsessiondisconnect(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	embed.NewEmbed().
		AddField("Session Disconnected Notifier", utils.Misc.FormatBoolOnOff(cfg.Events.SessionDisconnected)).
		Send(ctx)
}

func (ConfigCommands) sessiondisconnect(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	toggle := strings.ToLower(utils.Arguments["<on/off>"])
	if toggle != "on" && toggle != "off" {
		_ = ctx.RespondText("Invalid value. Use `on` or `off`.")
		return
	}

	cfg.Events.SessionDisconnected = toggle == "on"
	_ = cfg.Save()

	embed.NewEmbed().
		SetDescription("Turned " + toggle + " session disconnected notifier").
		Send(ctx)
}

func (ConfigCommands) getselfbotdetection(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	embed.NewEmbed().
		AddField("Selfbot Detection", utils.Misc.FormatBoolOnOff(cfg.Events.SelfbotDetection)).
		Send(ctx)
}

func (ConfigCommands) selfbotdetection(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	toggle := strings.ToLower(utils.Arguments["<on/off>"])
	if toggle != "on" && toggle != "off" {
		_ = ctx.RespondText("Invalid value. Use `on` or `off`.")
		return
	}

	cfg.Events.SelfbotDetection = toggle == "on"
	_ = cfg.Save()

	embed.NewEmbed().
		SetDescription("Turned " + toggle + " selfbot detection").
		Send(ctx)
}

func (ConfigCommands) getdmtypingnotifier(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	embed.NewEmbed().
		AddField("DM Typing Notifier", utils.Misc.FormatBoolOnOff(cfg.Events.DMTypingNotifier)).
		Send(ctx)
}

func (ConfigCommands) dmtypingnotifier(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	toggle := strings.ToLower(utils.Arguments["<on/off>"])
	if toggle != "on" && toggle != "off" {
		_ = ctx.RespondText("Invalid value. Use `on` or `off`.")
		return
	}

	cfg.Events.DMTypingNotifier = toggle == "on"
	_ = cfg.Save()

	embed.NewEmbed().
		SetDescription("Turned " + toggle + " DM typing notifier").
		Send(ctx)
}

func (ConfigCommands) getdmnotifier(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	embed.NewEmbed().
		AddField("DM Notifier", utils.Misc.FormatBoolOnOff(cfg.Events.DMNotifier)).
		Send(ctx)
}

func (ConfigCommands) dmnotifier(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	toggle := strings.ToLower(utils.Arguments["<on/off>"])
	if toggle != "on" && toggle != "off" {
		_ = ctx.RespondText("Invalid value. Use `on` or `off`.")
		return
	}

	cfg.Events.DMNotifier = toggle == "on"
	_ = cfg.Save()

	embed.NewEmbed().
		SetDescription("Turned " + toggle + " DM notifier").
		Send(ctx)
}

func (ConfigCommands) getdmdeletenotifier(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	embed.NewEmbed().
		AddField("DM Delete Notifier", utils.Misc.FormatBoolOnOff(cfg.Events.DMDeleteNotifier)).
		Send(ctx)
}

func (ConfigCommands) dmdeletenotifier(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	toggle := strings.ToLower(utils.Arguments["<on/off>"])
	if toggle != "on" && toggle != "off" {
		_ = ctx.RespondText("Invalid value. Use `on` or `off`.")
		return
	}

	cfg.Events.DMDeleteNotifier = toggle == "on"
	_ = cfg.Save()

	embed.NewEmbed().
		SetDescription("Turned " + toggle + " DM delete notifier").
		Send(ctx)
}

func (ConfigCommands) getdmeditnotifier(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	embed.NewEmbed().
		AddField("DM Edit Notifier", utils.Misc.FormatBoolOnOff(cfg.Events.DMEditNotifier)).
		Send(ctx)
}

func (ConfigCommands) dmeditnotifier(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	toggle := strings.ToLower(utils.Arguments["<on/off>"])
	if toggle != "on" && toggle != "off" {
		_ = ctx.RespondText("Invalid value. Use `on` or `off`.")
		return
	}

	cfg.Events.DMEditNotifier = toggle == "on"
	_ = cfg.Save()

	embed.NewEmbed().
		SetDescription("Turned " + toggle + " DM edit notifier").
		Send(ctx)
}

func (ConfigCommands) getmsgdeletenotifier(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	embed.NewEmbed().
		AddField("Message Delete Notifier", utils.Misc.FormatBoolOnOff(cfg.Events.MessageDeleteNotifier)).
		Send(ctx)
}

func (ConfigCommands) messagedeletenotifier(ctx *dgc.Ctx) {
	toggle := strings.ToLower(utils.Arguments["<on/off>"])
	if toggle != "on" && toggle != "off" {
		_ = ctx.RespondText("Invalid value. Use `on` or `off`.")
		return
	}
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	cfg.Events.MessageDeleteNotifier = toggle == "on"
	_ = cfg.Save()

	embed.NewEmbed().
		SetDescription("Turned " + toggle + " message delete notifiernotifier").
		Send(ctx)
}

func (ConfigCommands) getuitheme(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	embed.NewEmbed().
		AddField("UI Theme", cfg.UI.Theme).
		Send(ctx)
}

func (ConfigCommands) uitheme(ctx *dgc.Ctx) {
	theme := strings.ToLower(utils.Arguments["<light/dark>"])
	if theme != "light" && theme != "dark" {
		_ = ctx.RespondText("Invalid value. Use `light` or `dark`.")
		return
	}
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	cfg.UI.Theme = utils.String.Capitalize(theme)
	_ = cfg.Save()

	embed.NewEmbed().
		SetDescription("Set UI theme to " + cfg.UI.Theme).
		Send(ctx)
}

func (ConfigCommands) gethidetosystray(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	embed.NewEmbed().
		AddField("Hide to System Tray", utils.Misc.FormatBoolOnOff(cfg.UI.HideToSystemTray)).
		Send(ctx)
}

func (ConfigCommands) hidetosystray(ctx *dgc.Ctx) {
	toggle := strings.ToLower(utils.Arguments["<on/off>"])
	if toggle != "on" && toggle != "off" {
		_ = ctx.RespondText("Invalid value. Use `on` or `off`.")
		return
	}
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	cfg.UI.HideToSystemTray = toggle == "on"
	_ = cfg.Save()

	embed.NewEmbed().
		SetDescription("Turned " + toggle + " hide to system tray").
		Send(ctx)
}

func (ConfigCommands) getrunonstartup(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	embed.NewEmbed().
		AddField("Run On Startup", utils.Misc.FormatBoolOnOff(cfg.UI.RunOnStartup)).
		Send(ctx)
}

func (ConfigCommands) runonstartup(ctx *dgc.Ctx) {
	toggle := strings.ToLower(utils.Arguments["<on/off>"])
	if toggle != "on" && toggle != "off" {
		_ = ctx.RespondText("Invalid value. Use `on` or `off`.")
		return
	}
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	cfg.UI.RunOnStartup = toggle == "on"
	_ = cfg.Save()

	embed.NewEmbed().
		SetDescription("Turned " + toggle + " run on startup").
		Send(ctx)
}

func (ConfigCommands) getrunsbonstartup(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	embed.NewEmbed().
		AddField("Run Selfbot On Startup", utils.Misc.FormatBoolOnOff(cfg.UI.RunSelfbotOnStartup)).
		Send(ctx)
}

func (ConfigCommands) runsbonstartup(ctx *dgc.Ctx) {
	toggle := strings.ToLower(utils.Arguments["<on/off>"])
	if toggle != "on" && toggle != "off" {
		_ = ctx.RespondText("Invalid value. Use `on` or `off`.")
		return
	}
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	cfg.UI.RunSelfbotOnStartup = toggle == "on"
	_ = cfg.Save()

	embed.NewEmbed().
		SetDescription("Turned " + toggle + " run selfbot on startup").
		Send(ctx)
}
