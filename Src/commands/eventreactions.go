package commands

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/LightningDev1/LB-Selfbot/config"
	"github.com/LightningDev1/LB-Selfbot/embed"
	"github.com/LightningDev1/LB-Selfbot/eventhandlers"
	"github.com/LightningDev1/LB-Selfbot/utils"
	"github.com/LightningDev1/dgc"
	"github.com/LightningDev1/discordgo"
)

type EventReactionsCommands struct{}

func (e *EventReactionsCommands) RegisterCommands(router *dgc.Router) {
	router.StartCategory("EventReactions", "Event Reaction commands")
	router.RegisterCmd(&dgc.Command{
		Name:        "ignoreuser",
		Description: "Ignore messages from a user",
		Usage:       "[p]ignoreuser <user>",
		Handler:     e.ignoreuser,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "unignoreuser",
		Description: "Unignore a user",
		Usage:       "[p]unignoreuser <user>",
		Handler:     e.unignoreuser,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "userignorelist",
		Description: "See all ignored users",
		Usage:       "[p]userignorelist",
		Handler:     e.userignorelist,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "clearuserignorelist",
		Description: "Clear the list of ignored users",
		Usage:       "[p]clearuserignorelist",
		Handler:     e.clearuserignorelist,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "ignoreserver",
		Description: "Ignore messages from a server",
		Usage:       "[p]ignoreserver [server id]",
		Handler:     e.ignoreserver,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "unignoreserver",
		Description: "Unignore a server",
		Usage:       "[p]unignoreserver [server id]",
		Handler:     e.unignoreserver,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "serverignorelist",
		Description: "See all ignored servers",
		Usage:       "[p]serverignorelist",
		Handler:     e.serverignorelist,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "clearserverignorelist",
		Description: "Clear the list of ignored servers",
		Usage:       "[p]clearserverignorelist",
		Handler:     e.clearserverignorelist,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "pingmute",
		Description: "Mute a user when they ping you",
		Usage:       "[p]pingmute <on/off>",
		Handler:     e.pingmute,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "pingtimeout",
		Description: "Time a user out when they ping you",
		Usage:       "[p]pingtimeout <on/off>",
		Handler:     e.pingtimeout,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "pingban",
		Description: "Ban a user when they ping you",
		Usage:       "[p]pingban <on/off>",
		Handler:     e.pingban,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "pingkick",
		Description: "Kick a user when they ping you",
		Usage:       "[p]pingkick <on/off>",
		Handler:     e.pingkick,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "pinginsult",
		Description: "Insult a user when they ping you",
		Usage:       "[p]pinginsult <on/off>",
		Handler:     e.pinginsult,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "pingping",
		Description: "Ping a user back when they ping you",
		Usage:       "[p]pingping <on/off>",
		Handler:     e.pingping,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "afkenabled",
		Description: "Turn AFK mode on or off",
		Usage:       "[p]afkenabled <on/off>",
		Handler:     e.afkenabled,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "afkonping",
		Description: "Set if an AFK message should be sent on ping",
		Usage:       "[p]afkonping <on/off>",
		Handler:     e.afkonping,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "afkondm",
		Description: "Set if an AFK message should be sent on DM",
		Usage:       "[p]afkondm <on/off>",
		Handler:     e.afkondm,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "afkmessage",
		Description: "Set the AFK message",
		Usage:       "[p]afkmessage <text>",
		Handler:     e.afkmessage,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "afkreply",
		Description: "Set if the AFK message should be a reply",
		Usage:       "[p]afkreply <on/off>",
		Handler:     e.afkreply,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "afktyping",
		Description: "Set if the bot should show typing before sending message",
		Usage:       "[p]afktyping <on/off>",
		Handler:     e.afktyping,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "afkdelay",
		Description: "Set delay before sending AFK message",
		Usage:       "[p]afkdelay <seconds>",
		Handler:     e.afkdelay,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "afktriggertype",
		Description: "Set the AFK trigger type",
		Usage:       "[p]afktriggertype <type>",
		Handler:     e.afktriggertype,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "afkaddserver",
		Description: "Add a server to the AFK trigger list",
		Usage:       "[p]afkaddserver <server id>",
		Handler:     e.afkaddserver,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "afkremoveserver",
		Description: "Remove a server from the AFK trigger list",
		Usage:       "[p]afkremoveserver <server id>",
		Handler:     e.afkremoveserver,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "afkadduser",
		Description: "Add a user to the AFK trigger list",
		Usage:       "[p]afkadduser <<user>",
		Handler:     e.afkadduser,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "afkremoveuser",
		Description: "Remove a user from the AFK trigger list",
		Usage:       "[p]afkremoveuser <user>",
		Handler:     e.afkremoveuser,
	})
}

func (EventReactionsCommands) ignoreuser(ctx *dgc.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		_ = ctx.RespondText("Could not find user.")
		return
	}
	for _, ignoredUser := range eventhandlers.EventHandlers.Settings.IgnoredUsers {
		if user.ID == ignoredUser.ID {
			_ = ctx.RespondText("User is already ignored.")
			return
		}
	}
	eventhandlers.EventHandlers.Settings.IgnoredUsers = append(eventhandlers.EventHandlers.Settings.IgnoredUsers, user)
}

func (EventReactionsCommands) unignoreuser(ctx *dgc.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		_ = ctx.RespondText("Could not find user.")
		return
	}
	isIgnored := false
	for _, ignoredUser := range eventhandlers.EventHandlers.Settings.IgnoredUsers {
		if user.ID == ignoredUser.ID {
			isIgnored = true
		}
	}
	if !isIgnored {
		_ = ctx.RespondText("User is already not ignored.")
		return
	}
	eventhandlers.EventHandlers.Settings.IgnoredUsers = utils.RemoveSliceItem(
		eventhandlers.EventHandlers.Settings.IgnoredUsers,
		user,
	)
}

func (EventReactionsCommands) userignorelist(ctx *dgc.Ctx) {
	if len(eventhandlers.EventHandlers.Settings.IgnoredUsers) == 0 {
		_ = ctx.RespondText("No users are ignored.")
		return
	}
	ignoredUserString := ""
	for _, user := range eventhandlers.EventHandlers.Settings.IgnoredUsers {
		ignoredUserString += fmt.Sprintf("%s - %s\n", user.String(), user.ID)
	}
	embed.NewEmbed().
		AddField("Ignored Users", ignoredUserString).
		Send(ctx)
}

func (EventReactionsCommands) clearuserignorelist(*dgc.Ctx) {
	eventhandlers.EventHandlers.Settings.IgnoredUsers = []*discordgo.User{}
}

func (EventReactionsCommands) ignoreserver(ctx *dgc.Ctx) {
	guild, err := utils.Discord.GetGuild(ctx.Session, utils.Misc.Or(utils.Arguments["[server id]"], ctx.Event.GuildID))
	if err != nil {
		_ = ctx.RespondText("Error getting server: " + err.Error())
		return
	}
	for _, ignoredGuild := range eventhandlers.EventHandlers.Settings.IgnoredGuilds {
		if guild.ID == ignoredGuild.ID {
			_ = ctx.RespondText("Server is already ignored.")
			return
		}
	}
	eventhandlers.EventHandlers.Settings.IgnoredGuilds = append(eventhandlers.EventHandlers.Settings.IgnoredGuilds, guild)
}

func (EventReactionsCommands) unignoreserver(ctx *dgc.Ctx) {
	guild, err := utils.Discord.GetGuild(ctx.Session, utils.Misc.Or(utils.Arguments["[server id]"], ctx.Event.GuildID))
	if err != nil {
		_ = ctx.RespondText("Error getting server: " + err.Error())
		return
	}
	isIgnored := false
	for _, ignoredGuild := range eventhandlers.EventHandlers.Settings.IgnoredGuilds {
		if guild.ID == ignoredGuild.ID {
			isIgnored = true
		}
	}
	if !isIgnored {
		_ = ctx.RespondText("Server is already not ignored.")
		return
	}
	eventhandlers.EventHandlers.Settings.IgnoredGuilds = utils.RemoveSliceItem(
		eventhandlers.EventHandlers.Settings.IgnoredGuilds,
		guild,
	)

}

func (EventReactionsCommands) serverignorelist(ctx *dgc.Ctx) {
	if len(eventhandlers.EventHandlers.Settings.IgnoredGuilds) == 0 {
		_ = ctx.RespondText("No servers are ignored.")
		return
	}
	ignoredServersString := ""
	for _, guild := range eventhandlers.EventHandlers.Settings.IgnoredGuilds {
		ignoredServersString += fmt.Sprintf("%s - %s\n", guild.Name, guild.ID)
	}
	embed.NewEmbed().
		AddField("Ignored Servers", ignoredServersString).
		Send(ctx)
}

func (EventReactionsCommands) clearserverignorelist(*dgc.Ctx) {
	eventhandlers.EventHandlers.Settings.IgnoredGuilds = []*discordgo.Guild{}
}

func (EventReactionsCommands) pingmute(ctx *dgc.Ctx) {
	if !utils.Discord.CheckIfInServer(ctx) {
		return
	}
	toggleString := strings.ToLower(utils.Arguments["<on/off>"])
	if toggleString != "on" && toggleString != "off" {
		_ = ctx.RespondText("Invalid value. Use `on` or `off`.")
		return
	}
	alreadyOn := utils.SliceContains(eventhandlers.EventHandlers.Settings.PingMuteServers, ctx.Event.GuildID)
	if toggleString == "on" && alreadyOn {
		_ = ctx.RespondText("Ping mute is already on in this server.")
		return
	} else if toggleString == "off" && !alreadyOn {
		_ = ctx.RespondText("Ping mute is already off in this server.")
		return
	}
	if toggleString == "on" {
		eventhandlers.EventHandlers.Settings.PingMuteServers = append(eventhandlers.EventHandlers.Settings.PingMuteServers, ctx.Event.GuildID)
	} else {
		eventhandlers.EventHandlers.Settings.PingMuteServers = utils.RemoveSliceItem(
			eventhandlers.EventHandlers.Settings.PingMuteServers,
			ctx.Event.GuildID,
		)
	}
}

func (EventReactionsCommands) pingtimeout(ctx *dgc.Ctx) {
	if !utils.Discord.CheckIfInServer(ctx) {
		return
	}
	toggleString := strings.ToLower(utils.Arguments["<on/off>"])
	if toggleString != "on" && toggleString != "off" {
		_ = ctx.RespondText("Invalid value. Use `on` or `off`.")
		return
	}
	alreadyOn := utils.SliceContains(eventhandlers.EventHandlers.Settings.PingTimeoutServers, ctx.Event.GuildID)
	if toggleString == "on" && alreadyOn {
		_ = ctx.RespondText("Ping timeout is already on in this server.")
		return
	} else if toggleString == "off" && !alreadyOn {
		_ = ctx.RespondText("Ping timeout is already off in this server.")
		return
	}
	if toggleString == "on" {
		eventhandlers.EventHandlers.Settings.PingTimeoutServers = append(eventhandlers.EventHandlers.Settings.PingTimeoutServers, ctx.Event.GuildID)
	} else {
		eventhandlers.EventHandlers.Settings.PingTimeoutServers = utils.RemoveSliceItem(
			eventhandlers.EventHandlers.Settings.PingTimeoutServers,
			ctx.Event.GuildID,
		)
	}
}

func (EventReactionsCommands) pingban(ctx *dgc.Ctx) {
	if !utils.Discord.CheckIfInServer(ctx) {
		return
	}
	toggleString := strings.ToLower(utils.Arguments["<on/off>"])
	if toggleString != "on" && toggleString != "off" {
		_ = ctx.RespondText("Invalid value. Use `on` or `off`.")
		return
	}
	alreadyOn := utils.SliceContains(eventhandlers.EventHandlers.Settings.PingBanServers, ctx.Event.GuildID)
	if toggleString == "on" && alreadyOn {
		_ = ctx.RespondText("Ping ban is already on in this server.")
		return
	} else if toggleString == "off" && !alreadyOn {
		_ = ctx.RespondText("Ping ban is already off in this server.")
		return
	}
	if toggleString == "on" {
		eventhandlers.EventHandlers.Settings.PingBanServers = append(eventhandlers.EventHandlers.Settings.PingBanServers, ctx.Event.GuildID)
	} else {
		eventhandlers.EventHandlers.Settings.PingBanServers = utils.RemoveSliceItem(
			eventhandlers.EventHandlers.Settings.PingBanServers,
			ctx.Event.GuildID,
		)
	}
}

func (EventReactionsCommands) pingkick(ctx *dgc.Ctx) {
	if !utils.Discord.CheckIfInServer(ctx) {
		return
	}
	toggleString := strings.ToLower(utils.Arguments["<on/off>"])
	if toggleString != "on" && toggleString != "off" {
		_ = ctx.RespondText("Invalid value. Use `on` or `off`.")
		return
	}
	alreadyOn := utils.SliceContains(eventhandlers.EventHandlers.Settings.PingKickServers, ctx.Event.GuildID)
	if toggleString == "on" && alreadyOn {
		_ = ctx.RespondText("Ping kick is already on in this server.")
		return
	} else if toggleString == "off" && !alreadyOn {
		_ = ctx.RespondText("Ping kick is already off in this server.")
		return
	}
	if toggleString == "on" {
		eventhandlers.EventHandlers.Settings.PingKickServers = append(eventhandlers.EventHandlers.Settings.PingKickServers, ctx.Event.GuildID)
	} else {
		eventhandlers.EventHandlers.Settings.PingKickServers = utils.RemoveSliceItem(
			eventhandlers.EventHandlers.Settings.PingKickServers,
			ctx.Event.GuildID,
		)
	}
}

func (EventReactionsCommands) pinginsult(ctx *dgc.Ctx) {
	if !utils.Discord.CheckIfInServer(ctx) {
		return
	}
	toggleString := strings.ToLower(utils.Arguments["<on/off>"])
	if toggleString != "on" && toggleString != "off" {
		_ = ctx.RespondText("Invalid value. Use `on` or `off`.")
		return
	}
	alreadyOn := utils.SliceContains(eventhandlers.EventHandlers.Settings.PingInsultServers, ctx.Event.GuildID)
	if toggleString == "on" && alreadyOn {
		_ = ctx.RespondText("Ping insult is already on in this server.")
		return
	} else if toggleString == "off" && !alreadyOn {
		_ = ctx.RespondText("Ping insult is already off in this server.")
		return
	}
	if toggleString == "on" {
		eventhandlers.EventHandlers.Settings.PingInsultServers = append(eventhandlers.EventHandlers.Settings.PingInsultServers, ctx.Event.GuildID)
	} else {
		eventhandlers.EventHandlers.Settings.PingInsultServers = utils.RemoveSliceItem(
			eventhandlers.EventHandlers.Settings.PingInsultServers,
			ctx.Event.GuildID,
		)
	}
}

func (EventReactionsCommands) pingping(ctx *dgc.Ctx) {
	if !utils.Discord.CheckIfInServer(ctx) {
		return
	}
	toggleString := strings.ToLower(utils.Arguments["<on/off>"])
	if toggleString != "on" && toggleString != "off" {
		_ = ctx.RespondText("Invalid value. Use `on` or `off`.")
		return
	}
	alreadyOn := utils.SliceContains(eventhandlers.EventHandlers.Settings.PingPingServers, ctx.Event.GuildID)
	if toggleString == "on" && alreadyOn {
		_ = ctx.RespondText("Ping ping is already on in this server.")
		return
	} else if toggleString == "off" && !alreadyOn {
		_ = ctx.RespondText("Ping ping is already off in this server.")
		return
	}
	if toggleString == "on" {
		eventhandlers.EventHandlers.Settings.PingPingServers = append(eventhandlers.EventHandlers.Settings.PingPingServers, ctx.Event.GuildID)
	} else {
		eventhandlers.EventHandlers.Settings.PingPingServers = utils.RemoveSliceItem(
			eventhandlers.EventHandlers.Settings.PingPingServers,
			ctx.Event.GuildID,
		)
	}
}

func (EventReactionsCommands) afkenabled(ctx *dgc.Ctx) {
	toggle := strings.ToLower(utils.Arguments["<on/off>"])
	if toggle != "on" && toggle != "off" {
		_ = ctx.RespondText("Invalid value. Use `on` or `off`.")
		return
	}
	cfg, _ := config.GetConfig()
	cfg.AFK.Enabled = toggle == "on"
	_ = cfg.Save()
}

func (EventReactionsCommands) afkonping(ctx *dgc.Ctx) {
	toggle := strings.ToLower(utils.Arguments["<on/off>"])
	if toggle != "on" && toggle != "off" {
		_ = ctx.RespondText("Invalid value. Use `on` or `off`.")
		return
	}
	cfg, _ := config.GetConfig()
	cfg.AFK.OnPing = toggle == "on"
	_ = cfg.Save()
}

func (EventReactionsCommands) afkondm(ctx *dgc.Ctx) {
	toggle := strings.ToLower(utils.Arguments["<on/off>"])
	if toggle != "on" && toggle != "off" {
		_ = ctx.RespondText("Invalid value. Use `on` or `off`.")
		return
	}
	cfg, _ := config.GetConfig()
	cfg.AFK.OnDM = toggle == "on"
	_ = cfg.Save()
}

func (EventReactionsCommands) afkmessage(ctx *dgc.Ctx) {
	cfg, _ := config.GetConfig()
	cfg.AFK.Message = ctx.Arguments.Raw()
	_ = cfg.Save()
}

func (EventReactionsCommands) afkreply(ctx *dgc.Ctx) {
	toggle := strings.ToLower(utils.Arguments["<on/off>"])
	if toggle != "on" && toggle != "off" {
		_ = ctx.RespondText("Invalid value. Use `on` or `off`.")
		return
	}
	cfg, _ := config.GetConfig()
	cfg.AFK.Reply = toggle == "on"
	_ = cfg.Save()
}

func (EventReactionsCommands) afktyping(ctx *dgc.Ctx) {
	toggle := strings.ToLower(utils.Arguments["<on/off>"])
	if toggle != "on" && toggle != "off" {
		_ = ctx.RespondText("Invalid value. Use `on` or `off`.")
		return
	}
	cfg, _ := config.GetConfig()
	cfg.AFK.Typing = toggle == "on"
	_ = cfg.Save()
}

func (EventReactionsCommands) afkdelay(ctx *dgc.Ctx) {
	cfg, _ := config.GetConfig()
	delay, err := strconv.Atoi(utils.Arguments["<seconds>"])
	if err != nil {
		_ = ctx.RespondText("Argument <seconds> must be a number.")
		return
	}
	cfg.AFK.Delay = delay
	_ = cfg.Save()
}

func (EventReactionsCommands) afktriggertype(ctx *dgc.Ctx) {
	arg := strings.ToLower(utils.Arguments["<type>"])
	if !utils.SliceContains([]string{"global", "whitelist", "blacklist"}, arg) {
		_ = ctx.RespondText("Invalid value. Use `Global`, `Whitelist`, or `Blacklist`.")
		return
	}

	cfg, _ := config.GetConfig()
	cfg.AFK.TriggerList.Type = utils.String.Capitalize(arg)
	_ = cfg.Save()
}

func (EventReactionsCommands) afkaddserver(ctx *dgc.Ctx) {
	guild, err := utils.Discord.GetGuild(ctx.Session, utils.Misc.Or(utils.Arguments["<server id>"], ctx.Event.GuildID))
	if err != nil {
		_ = ctx.RespondText("Could not find server.")
		return
	}
	cfg, _ := config.GetConfig()
	if utils.SliceContains(cfg.AFK.TriggerList.Servers, guild.ID) {
		_ = ctx.RespondText("Server is already added.")
		return
	}
	cfg.AFK.TriggerList.Servers = append(cfg.AFK.TriggerList.Servers, guild.ID)
	_ = cfg.Save()
}

func (EventReactionsCommands) afkremoveserver(ctx *dgc.Ctx) {
	guild, err := utils.Discord.GetGuild(ctx.Session, utils.Misc.Or(utils.Arguments["<server id>"], ctx.Event.GuildID))
	if err != nil {
		_ = ctx.RespondText("Could not find server.")
		return
	}
	cfg, _ := config.GetConfig()
	if !utils.SliceContains(cfg.AFK.TriggerList.Servers, guild.ID) {
		_ = ctx.RespondText("Server is already removed.")
		return
	}
	cfg.AFK.TriggerList.Servers = utils.RemoveSliceItem(cfg.AFK.TriggerList.Servers, guild.ID)
	_ = cfg.Save()
}

func (EventReactionsCommands) afkadduser(ctx *dgc.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		_ = ctx.RespondText("Could not find user.")
		return
	}
	cfg, _ := config.GetConfig()
	if utils.SliceContains(cfg.AFK.TriggerList.Users, user.ID) {
		_ = ctx.RespondText("User is already added.")
		return
	}
	cfg.AFK.TriggerList.Users = append(cfg.AFK.TriggerList.Users, user.ID)
	_ = cfg.Save()
}

func (EventReactionsCommands) afkremoveuser(ctx *dgc.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		_ = ctx.RespondText("Could not find user.")
		return
	}
	cfg, _ := config.GetConfig()
	if !utils.SliceContains(cfg.AFK.TriggerList.Users, user.ID) {
		_ = ctx.RespondText("User is already removed.")
		return
	}
	cfg.AFK.TriggerList.Users = utils.RemoveSliceItem(cfg.AFK.TriggerList.Users, user.ID)
	_ = cfg.Save()
}
