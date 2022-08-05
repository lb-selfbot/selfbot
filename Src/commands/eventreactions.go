package commands

import (
	"fmt"
	"github.com/LightningDev1/LB-Selfbot/embed"
	"github.com/LightningDev1/LB-Selfbot/eventhandlers"
	"github.com/LightningDev1/LB-Selfbot/utils"
	"github.com/LightningDev1/dgc"
	"github.com/LightningDev1/discordgo"
)

type EventReactionsCommands struct{}

func (e *EventReactionsCommands) RegisterCommands(router *dgc.Router) {
	router.StartCategory("Event Reaction", "Event Reaction commands")
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
