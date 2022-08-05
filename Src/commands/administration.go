package commands

import (
	"encoding/json"
	"fmt"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/LightningDev1/LB-Selfbot/api"
	"github.com/LightningDev1/LB-Selfbot/embed"
	"github.com/LightningDev1/LB-Selfbot/http"
	"github.com/LightningDev1/LB-Selfbot/utils"
	"github.com/LightningDev1/dgc"
	"github.com/LightningDev1/discordgo"
)

var (
	verificationLevels = map[int]string{
		0: "None",
		1: "Low",
		2: "Medium",
		3: "High",
		4: "Highest",
	}

	boostLevels = map[int]string{
		0: "None",
		1: "Level 1",
		2: "Level 2",
		3: "Level 3",
	}
)

type AdministrationCommands struct{}

func (m AdministrationCommands) RegisterCommands(router *dgc.Router) {
	router.StartCategory("Administration", "Administration commands")
	router.RegisterCmd(&dgc.Command{
		Name:        "channelinfo",
		Description: "Get info about a channel",
		Usage:       "[p]channelinfo <#channel>",
		Handler:     m.channelinfo,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "textchannel",
		Description: "Create a text channel",
		Usage:       "[p]textchannel <name>",
		Handler:     m.textchannel,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "voicechannel",
		Description: "Create a voice channel",
		Usage:       "[p]voicechannel <name>",
		Handler:     m.voicechannel,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "stagechannel",
		Description: "Create a stage channel",
		Usage:       "[p]stagechannel <name>",
		Handler:     m.stagechannel,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "newschannel",
		Description: "Create a news channel",
		Usage:       "[p]newschannel <name>",
		Handler:     m.newschannel,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "thread",
		Description: "Create a thread",
		Usage:       "[p]thread <name>",
		Handler:     m.thread,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "category",
		Description: "Create a category",
		Usage:       "[p]category <name>",
		Handler:     m.category,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "deletecategory",
		Description: "Delete a category",
		Usage:       "[p]deletecategory <category id>",
		Handler:     m.deletecategory,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "renamechannel",
		Description: "Rename a channel",
		Usage:       "[p]renamechannel <channel> <name>",
		Handler:     m.renamechannel,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "deletechannel",
		Description: "Delete a channel",
		Usage:       "[p]deletechannel <channel>",
		Handler:     m.deletechannel,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "slowmode",
		Description: "Change slowmode amount",
		Usage:       "[p]slowmode <channel> <seconds>",
		Handler:     m.slowmode,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "removeslowmode",
		Description: "Remove clowmode",
		Usage:       "[p]removeslowmode <channel>",
		Handler:     m.removeslowmode,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "lockchannel",
		Description: "Lock a channel",
		Usage:       "[p]lockchannel <channel>",
		Handler:     m.lockchannel,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "unlockchannel",
		Description: "Unlock a channel",
		Usage:       "[p]unlockchannel <channel>",
		Handler:     m.unlockchannel,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "purge",
		Description: "Delete messages in a channel",
		Usage:       "[p]purge <amount> [channel]",
		Handler:     m.purge,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "nuke",
		Description: "Nuke a channel",
		Usage:       "[p]nuke [channel]",
		Handler:     m.nuke,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "webhookinfo",
		Description: "Get info about a webhook",
		Usage:       "[p]webhookinfo <webhook url>",
		Handler:     m.webhookinfo,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "webhooksend",
		Description: "Send a webhook message",
		Usage:       "[p]webhooksend <webhook url> <\"message\">",
		Handler:     m.webhooksend,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "webhookdelete",
		Description: "Delete a webhook",
		Usage:       "[p]webhookdelete <webhook url>",
		Handler:     m.webhookdelete,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "userinfo",
		Description: "Get info about a user",
		Usage:       "[p]userinfo <user>",
		Handler:     m.userinfo,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "userage",
		Description: "See the age of a users account",
		Usage:       "[p]userage <user>",
		Handler:     m.userage,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "mutualfriends",
		Description: "Get mutual friends with user",
		Usage:       "[p]mutualfriends <user>",
		Handler:     m.mutualfriends,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "mutualservers",
		Description: "Get mutual servers with user",
		Usage:       "[p]mutualservers <user>",
		Handler:     m.mutualservers,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "mobile",
		Description: "See if a user is on mobile",
		Usage:       "[p]mobile <user>",
		Handler:     m.mobile,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "web",
		Description: "See if a user is on web",
		Usage:       "[p]web <user>",
		Handler:     m.web,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "desktop",
		Description: "See if a user is on desktop",
		Usage:       "[p]desktop <user>",
		Handler:     m.desktop,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "avatar",
		Description: "Get a users avatar",
		Usage:       "[p]avatar <user>",
		Handler:     m.avatar,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "banner",
		Description: "Get a users banner",
		Usage:       "[p]banner <user>",
		Handler:     m.banner,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "report",
		Description: "Report a message",
		Usage:       "[p]report <message link> <reason>",
		Handler:     m.report,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "mute",
		Description: "Mute a user",
		Usage:       "[p]mute <user>",
		Handler:     m.mute,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "unmute",
		Description: "Unmute a user",
		Usage:       "[p]unmute <user>",
		Handler:     m.unmute,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "timeout",
		Description: "Timeout a user",
		Usage:       "[p]timeout <user> <seconds>",
		Handler:     m.timeout,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "removetimeout",
		Description: "Remove a timeout",
		Usage:       "[p]removetimeout <user>",
		Handler:     m.removetimeout,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "kick",
		Description: "Kick a user",
		Usage:       "[p]kick <user> [reason]",
		Handler:     m.kick,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "ban",
		Description: "Ban a user",
		Usage:       "[p]ban <user> [reason]",
		Handler:     m.ban,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "hackban",
		Description: "Ban a user that's not in the server",
		Usage:       "[p]hackban <user id> [reason]",
		Handler:     m.ban,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "softban",
		Description: "Ban and instantly unban a user",
		Usage:       "[p]softban <user> [reason]",
		Handler:     m.softban,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "unban",
		Description: "Ban and instantly unban a user",
		Usage:       "[p]unban <user id>",
		Handler:     m.unban,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "banlist",
		Description: "Get a list of banned users",
		Usage:       "[p]banlist [server id]",
		Handler:     m.banlist,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "copybans",
		Description: "Copy bans from one server to another",
		Usage:       "[p]copybans <source server id> [dest server id]",
		Handler:     m.copybans,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "nick",
		Description: "Change your nickname",
		Usage:       "[p]nick <text>",
		Handler:     m.nick,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "nickuser",
		Description: "Change a users nickname",
		Usage:       "[p]nickuser <user> <text>",
		Handler:     m.nickuser,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "clearnick",
		Description: "Clear a users nickname",
		Usage:       "[p]clearnick [user]",
		Handler:     m.clearnick,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "serverinfo",
		Description: "Get info about a server",
		Usage:       "[p]serverinfo [server id]",
		Handler:     m.serverinfo,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "serverage",
		Description: "Get the age of a server",
		Usage:       "[p]serverage [server id]",
		Handler:     m.serverage,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "serverowner",
		Description: "Get a servers owner",
		Usage:       "[p]serverowner [server id]",
		Handler:     m.serverowner,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "servericon",
		Description: "Get a servers icon",
		Usage:       "[p]servericon [server id]",
		Handler:     m.servericon,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "serverbanner",
		Description: "Get a servers banner",
		Usage:       "[p]serverbanner [server id]",
		Handler:     m.serverbanner,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "serversplash",
		Description: "Get a servers invite splash",
		Usage:       "[p]serversplash [server id]",
		Handler:     m.serversplash,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "serverdsplash",
		Description: "Get a servers discovery splash",
		Usage:       "[p]serverdsplash [server id]",
		Handler:     m.serverdsplash,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "serverboosters",
		Description: "Get a list of server boosters",
		Usage:       "[p]serverboosters [server id]",
		Handler:     m.serverboosters,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "serverboostlevel",
		Description: "Get a servers boost level",
		Usage:       "[p]serverboostlevel [server id]",
		Handler:     m.serverboostlevel,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "setservername",
		Description: "Set a servers name",
		Usage:       "[p]setservername <name>",
		Handler:     m.setservername,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "setservericon",
		Description: "Set a servers icon",
		Usage:       "[p]setservericon <url>",
		Handler:     m.setservericon,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "setserverbanner",
		Description: "Set a servers banner",
		Usage:       "[p]setserverbanner <url>",
		Handler:     m.setserverbanner,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "lockserver",
		Description: "Lock a server",
		Usage:       "[p]lockserver [server id]",
		Handler:     m.lockserver,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "unlockserver",
		Description: "Unlock a server",
		Usage:       "[p]unlockserver [server id]",
		Handler:     m.unlockserver,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "channels",
		Description: "See all channels in a server",
		Usage:       "[p]channels [server id]",
		Handler:     m.channels,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "bots",
		Description: "See all bots in a server",
		Usage:       "[p]bots [server id]",
		Handler:     m.bots,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "admins",
		Description: "See all admins in a server",
		Usage:       "[p]admins [server id]",
		Handler:     m.admins,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "roles",
		Description: "See all roles in a server",
		Usage:       "[p]roles [server id]",
		Handler:     m.roles,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "invite",
		Description: "Create an invite",
		Usage:       "[p]invite [channel] [max age] [max uses]",
		Handler:     m.invite,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "friendinvite",
		Description: "Create a friend invite",
		Usage:       "[p]friendinvite",
		Handler:     m.friendinvite,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "inviteinfo",
		Description: "Get info about an invite",
		Usage:       "[p]inviteinfo <invite>",
		Handler:     m.inviteinfo,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "deleteinvite",
		Description: "Delete an invite",
		Usage:       "[p]deleteinvite <invite>",
		Handler:     m.deleteinvite,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "serverinvites",
		Description: "Get all invites in a server",
		Usage:       "[p]serverinvites [server id]",
		Handler:     m.serverinvites,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "channelinvites",
		Description: "Get all invites in a channel",
		Usage:       "[p]channelinvites [channel]",
		Handler:     m.channelinvites,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "userinvites",
		Description: "Get all invites created by a user",
		Usage:       "[p]userinvites <user> [server id]",
		Handler:     m.userinvites,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "roleinfo",
		Description: "Get info about a role",
		Usage:       "[p]roleinfo <role>",
		Handler:     m.roleinfo,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "giverole",
		Description: "Give a role to a user",
		Usage:       "[p]giverole <user> <role>",
		Handler:     m.giverole,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "removerole",
		Description: "Remove a role from a user",
		Usage:       "[p]removerole <user> <role>",
		Handler:     m.removerole,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "createrole",
		Description: "Create a role",
		Usage:       "[p]createrole <name>",
		Handler:     m.createrole,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "deleterole",
		Description: "Delete a role",
		Usage:       "[p]deleterole <role>",
		Handler:     m.deleterole,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "renamerole",
		Description: "Rename a role",
		Usage:       "[p]renamerole <role> <name>",
		Handler:     m.renamerole,
	})
}

// Channel Control

func (AdministrationCommands) channelinfo(ctx *dgc.Ctx) {
	channel, success := utils.Discord.TryGetChannelFromMention(ctx, ctx.Arguments.Get(0))
	if !success {
		_ = ctx.RespondText("Could not find channel")
		return
	}

	embed.NewEmbed().
		AddField("Name", channel.Name).
		AddField("ID", channel.ID).
		AddField("Topic", channel.Topic).
		AddField("Category ID", channel.ParentID).
		AddField("Is NSFW", utils.Misc.FormatBool(channel.NSFW)).
		AddField("Created At", utils.Discord.SnowflakeToTime(channel.ID).Format(time.RFC1123)).
		AddField("Slowmode Delay", fmt.Sprint(channel.RateLimitPerUser)+" seconds").
		AddField("Position", fmt.Sprint(channel.Position)).
		AddField("Last Message ID", channel.LastMessageID).
		Send(ctx)
}

func (AdministrationCommands) textchannel(ctx *dgc.Ctx) {
	if !utils.Discord.CheckIfInServer(ctx) {
		return
	}
	channel, err := ctx.Session.GuildChannelCreate(ctx.Event.GuildID, ctx.Arguments.Raw(), discordgo.ChannelTypeGuildText)
	if err != nil {
		_ = ctx.RespondText("Error creating channel: " + err.Error())
		return
	}
	embed.NewEmbed().
		SetDescription("Successfully created text channel!").
		AddField("Channel", channel.Mention()).
		Send(ctx)
}

func (AdministrationCommands) voicechannel(ctx *dgc.Ctx) {
	if !utils.Discord.CheckIfInServer(ctx) {
		return
	}
	channel, err := ctx.Session.GuildChannelCreate(ctx.Event.GuildID, ctx.Arguments.Raw(), discordgo.ChannelTypeGuildVoice)
	if err != nil {
		_ = ctx.RespondText("Error creating channel: " + err.Error())
		return
	}
	embed.NewEmbed().
		SetDescription("Successfully created voice channel!").
		AddField("Channel", channel.Mention()).
		Send(ctx)
}

func (AdministrationCommands) stagechannel(ctx *dgc.Ctx) {
	if !utils.Discord.CheckIfInServer(ctx) {
		return
	}
	channel, err := ctx.Session.GuildChannelCreate(ctx.Event.GuildID, ctx.Arguments.Raw(), discordgo.ChannelTypeGuildStageVoice)
	if err != nil {
		_ = ctx.RespondText("Error creating channel: " + err.Error())
		return
	}
	embed.NewEmbed().
		SetDescription("Successfully created stage channel!").
		AddField("Channel", channel.Mention()).
		Send(ctx)
}

func (AdministrationCommands) newschannel(ctx *dgc.Ctx) {
	if !utils.Discord.CheckIfInServer(ctx) {
		return
	}
	channel, err := ctx.Session.GuildChannelCreate(ctx.Event.GuildID, ctx.Arguments.Raw(), discordgo.ChannelTypeGuildNews)
	if err != nil {
		_ = ctx.RespondText("Error creating channel: " + err.Error())
		return
	}
	embed.NewEmbed().
		SetDescription("Successfully created news channel!").
		AddField("Channel", channel.Mention()).
		Send(ctx)
}

func (AdministrationCommands) thread(ctx *dgc.Ctx) {
	if !utils.Discord.CheckIfInServer(ctx) {
		return
	}
	channel, err := ctx.Session.ThreadStart(ctx.Event.ChannelID, ctx.Arguments.Raw(), discordgo.ChannelTypeGuildPublicThread, 0)
	if err != nil {
		_ = ctx.RespondText("Error creating thread: " + err.Error())
		return
	}
	embed.NewEmbed().
		SetDescription("Successfully created thread!").
		AddField("Thread", channel.Mention()).
		Send(ctx)
}

func (AdministrationCommands) category(ctx *dgc.Ctx) {
	if !utils.Discord.CheckIfInServer(ctx) {
		return
	}
	channel, err := ctx.Session.GuildChannelCreate(ctx.Event.GuildID, ctx.Arguments.Raw(), discordgo.ChannelTypeGuildCategory)
	if err != nil {
		_ = ctx.RespondText("Error creating category: " + err.Error())
		return
	}
	embed.NewEmbed().
		SetDescription("Successfully created category!").
		AddField("Category", channel.Mention()).
		Send(ctx)
}

func (AdministrationCommands) deletecategory(ctx *dgc.Ctx) {
	channel, success := utils.Discord.TryGetChannelFromMention(ctx, ctx.Arguments.Get(0))
	if !success {
		_ = ctx.RespondText("Could not find category")
		return
	}
	_, err := ctx.Session.ChannelDelete(channel.ID)
	if err != nil {
		_ = ctx.RespondText("Error deleting category: " + err.Error())
		return
	}

	embed.NewEmbed().
		SetDescription("Successfully deleted category!").
		Send(ctx)
}

func (AdministrationCommands) renamechannel(ctx *dgc.Ctx) {
	channel, success := utils.Discord.TryGetChannelFromMention(ctx, ctx.Arguments.Get(0))
	if !success {
		_ = ctx.RespondText("Could not find channel")
		return
	}

	_, err := ctx.Session.ChannelEdit(channel.ID, utils.String.RemoveWords(ctx.Arguments.Raw(), 1))
	if err != nil {
		_ = ctx.RespondText("Error renaming channel: " + err.Error())
		return
	}

	embed.NewEmbed().
		SetDescription("Successfully renamed channel!").
		AddField("Channel", channel.Mention()).
		Send(ctx)
}

func (AdministrationCommands) deletechannel(ctx *dgc.Ctx) {
	channel, success := utils.Discord.TryGetChannelFromMention(ctx, ctx.Arguments.Get(0))
	if !success {
		_ = ctx.RespondText("Could not find channel")
		return
	}

	_, err := ctx.Session.ChannelDelete(channel.ID)
	if err != nil {
		_ = ctx.RespondText("Error deleting channel: " + err.Error())
		return
	}

	embed.NewEmbed().
		SetDescription("Successfully deleted channel!").
		Send(ctx)
}

func (AdministrationCommands) slowmode(ctx *dgc.Ctx) {
	channel, success := utils.Discord.TryGetChannelFromMention(ctx, ctx.Arguments.Get(0))
	if !success {
		_ = ctx.RespondText("Could not find channel")
		return
	}

	seconds, err := strconv.Atoi(utils.Arguments["<seconds>"])
	if err != nil {
		_ = ctx.RespondText("Argument <seconds> must be a number.")
		return
	}

	_, err = ctx.Session.ChannelEditComplex(channel.ID, &discordgo.ChannelEdit{
		RateLimitPerUser: &seconds,
	})
	if err != nil {
		_ = ctx.RespondText("Error setting slowmode: " + err.Error())
		return
	}

	embed.NewEmbed().
		SetDescription(fmt.Sprintf("Successfully set slowmode to %s seconds!", utils.Arguments["<seconds>"])).
		AddField("Channel", channel.Mention()).
		Send(ctx)
}

func (AdministrationCommands) removeslowmode(ctx *dgc.Ctx) {
	channel, success := utils.Discord.TryGetChannelFromMention(ctx, ctx.Arguments.Get(0))
	if !success {
		_ = ctx.RespondText("Could not find channel")
		return
	}

	_, err := ctx.Session.ChannelEditComplex(channel.ID, &discordgo.ChannelEdit{
		RateLimitPerUser: nil,
	})
	if err != nil {
		_ = ctx.RespondText("Error removing slowmode: " + err.Error())
		return
	}

	embed.NewEmbed().
		SetDescription("Successfully removed slowmode!").
		AddField("Channel", channel.Mention()).
		Send(ctx)
}

func (AdministrationCommands) lockchannel(ctx *dgc.Ctx) {
	channel, success := utils.Discord.TryGetChannelFromMention(ctx, ctx.Arguments.Get(0))
	if !success {
		_ = ctx.RespondText("Could not find channel")
		return
	}

	_, err := ctx.Session.ChannelEditComplex(channel.ID, &discordgo.ChannelEdit{
		PermissionOverwrites: []*discordgo.PermissionOverwrite{
			{
				ID:   channel.GuildID,
				Type: discordgo.PermissionOverwriteTypeRole,
				Deny: discordgo.PermissionSendMessages,
			},
		},
	})
	if err != nil {
		_ = ctx.RespondText("Error locking channel: " + err.Error())
		return
	}

	embed.NewEmbed().
		SetDescription("Successfully locked channel!").
		AddField("Channel", channel.Mention()).
		Send(ctx)
}

func (AdministrationCommands) unlockchannel(ctx *dgc.Ctx) {
	channel, success := utils.Discord.TryGetChannelFromMention(ctx, ctx.Arguments.Get(0))
	if !success {
		_ = ctx.RespondText("Could not find channel")
		return
	}

	_, err := ctx.Session.ChannelEditComplex(channel.ID, &discordgo.ChannelEdit{
		PermissionOverwrites: []*discordgo.PermissionOverwrite{
			{
				ID:   channel.GuildID,
				Type: discordgo.PermissionOverwriteTypeRole,
				Deny: 0,
			},
		},
	})
	if err != nil {
		_ = ctx.RespondText("Error unlocking channel: " + err.Error())
		return
	}

	embed.NewEmbed().
		SetDescription("Successfully unlocked channel!").
		AddField("Channel", channel.Mention()).
		Send(ctx)
}

func (AdministrationCommands) purge(ctx *dgc.Ctx) {
	amount, err := strconv.Atoi(utils.Arguments["<amount>"])
	if err != nil {
		_ = ctx.RespondText("Argument <amount> must be a number.")
		return
	}

	channel, success := utils.Discord.TryGetChannelFromMention(ctx, ctx.Arguments.Get(1))
	channelID := ctx.Event.ChannelID
	if success {
		channelID = channel.ID
	}

	messages, err := utils.Discord.GetChannelMessages(ctx, channelID, amount)
	if err != nil {
		_ = ctx.RespondText("Error getting messages: " + err.Error())
		return
	}

	for _, message := range messages {
		err = ctx.Session.ChannelMessageDelete(channelID, message.ID)
		if err != nil {
			_ = ctx.RespondText("Error deleting messages: " + err.Error())
			return
		}
	}
}

func (AdministrationCommands) nuke(ctx *dgc.Ctx) {
	channel, success := utils.Discord.TryGetChannelFromMention(ctx, ctx.Arguments.Get(0))
	if !success {
		channel, _ = ctx.Session.State.Channel(ctx.Event.ChannelID)
		if channel == nil {
			_ = ctx.RespondText("Could not find channel")
			return
		}
	}

	_, err := ctx.Session.ChannelDelete(channel.ID)
	if err != nil {
		_ = ctx.RespondText("Error deleting channel: " + err.Error())
		return
	}

	channel, err = ctx.Session.GuildChannelCreateComplex(ctx.Event.GuildID, discordgo.GuildChannelCreateData{
		Name:                 channel.Name,
		Type:                 channel.Type,
		Topic:                channel.Topic,
		Bitrate:              channel.Bitrate,
		UserLimit:            channel.UserLimit,
		RateLimitPerUser:     channel.RateLimitPerUser,
		Position:             channel.Position,
		PermissionOverwrites: channel.PermissionOverwrites,
		ParentID:             channel.ParentID,
		NSFW:                 channel.NSFW,
	})
	if err != nil {
		_ = ctx.RespondText("Error creating channel: " + err.Error())
		return
	}
	_, _ = ctx.Session.ChannelEditComplex(channel.ID, &discordgo.ChannelEdit{
		Position: channel.Position,
	})

	ctx.Event.ChannelID = channel.ID

	embed.NewEmbed().
		SetDescription("Successfully nuked channel!").
		Send(ctx)
}

// Webhook Control

func (AdministrationCommands) webhookinfo(ctx *dgc.Ctx) {
	webhookUrl, err := url.Parse(utils.Arguments["<webhook url>"])
	if err != nil {
		_ = ctx.RespondText("Invalid webhook URL!")
		return
	}
	if !strings.Contains(webhookUrl.Host, "discord.com") {
		_ = ctx.RespondText("Invalid webhook URL!")
		return
	}
	webhookInfo, statusCode, err := http.GetJsonStatusCode(webhookUrl.String())
	if err != nil {
		_ = ctx.RespondText("Invalid webhook URL!")
		return
	}
	if statusCode != 200 {
		_ = ctx.RespondText("Invalid webhook URL!")
		return
	}
	embed.NewEmbed().
		AddField("Type", map[int]string{1: "Incoming", 2: "Channel Follower", 3: "Application"}[webhookInfo.Get("type").MustInt()]).
		AddField("ID", webhookInfo.Get("id").MustString()).
		AddField("Name", webhookInfo.Get("name").MustString()).
		AddField("Avatar", webhookInfo.Get("avatar").MustString()).
		AddField("Channel ID", webhookInfo.Get("channel_id").MustString()).
		AddField("Server ID", webhookInfo.Get("guild_id").MustString()).
		AddField("Application ID", webhookInfo.Get("application_id").MustString()).
		AddField("Token", webhookInfo.Get("token").MustString()).
		Send(ctx)
}

func (AdministrationCommands) webhooksend(ctx *dgc.Ctx) {
	webhookUrl, err := url.Parse(utils.Arguments["<webhook url>"])
	if err != nil {
		_ = ctx.RespondText("Invalid webhook URL!")
		return
	}
	if !strings.Contains(webhookUrl.Host, "discord.com") {
		_ = ctx.RespondText("Invalid webhook URL!")
		return
	}
	data := struct {
		Content string `json:"content"`
	}{utils.Arguments["<\"message\">"]}
	jsonData, _ := json.Marshal(data)
	_, _ = http.Post(webhookUrl.String(), string(jsonData))
}

func (AdministrationCommands) webhookdelete(ctx *dgc.Ctx) {
	webhookUrl, err := url.Parse(utils.Arguments["<webhook url>"])
	if err != nil {
		_ = ctx.RespondText("Invalid webhook URL!")
		return
	}
	if !strings.Contains(webhookUrl.Host, "discord.com") {
		_ = ctx.RespondText("Invalid webhook URL!")
		return
	}
	_, _ = http.Delete(webhookUrl.String())
	embed.NewEmbed().
		SetDescription("Successfully deleted webhook.").
		Send(ctx)
}

// Member Control

func (AdministrationCommands) userinfo(ctx *dgc.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		_ = ctx.RespondText("Could not find user")
		return
	}
	member, _ := ctx.Session.GuildMember(ctx.Event.GuildID, user.ID)

	embed := embed.NewEmbed().
		SetDescription(fmt.Sprintf("User Info (%s - %s)", user.String(), user.ID)).
		AddField("Avatar", user.AvatarURL("")).
		AddField("Banner", user.BannerURL("")).
		AddField("Banner Color", fmt.Sprintf("#%x", user.AccentColor)).
		AddField("Is Bot", strconv.FormatBool(user.Bot)).
		AddField("Created At", utils.Discord.SnowflakeToTime(user.ID).Format(time.RFC1123)).
		AddField("Flags", fmt.Sprint(user.Flags))
	if member != nil {
		embed.AddField("Nickname", member.Nick)
		embed.AddField("Joined Server At", member.JoinedAt.Format(time.RFC1123))
		embed.AddField("Role Amount", fmt.Sprint(len(member.Roles)))
		embed.AddField("Roles", utils.Discord.GetRoleNames(ctx, member))
	}
	embed.Send(ctx)
}

func (AdministrationCommands) userage(ctx *dgc.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		_ = ctx.RespondText("Could not find user")
		return
	}
	age := time.Since(utils.Discord.SnowflakeToTime(user.ID))
	embed.NewEmbed().
		SetDescription(fmt.Sprintf("User Age (%s - %s)", user.String(), user.ID)).
		AddField("Age", utils.Misc.FormatDuration(age)).
		Send(ctx)
}

func (AdministrationCommands) mutualfriends(ctx *dgc.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		_ = ctx.RespondText("Could not find user")
		return
	}
	mutualFriends, err := ctx.Session.UserMutualFriends(user.ID)
	if err != nil {
		_ = ctx.RespondText("Could not find user")
		return
	}
	mutualFriendsString := ""
	for _, friend := range mutualFriends {
		mutualFriendsString += fmt.Sprintf("%s - %s\n", friend.String(), friend.ID)
	}
	embed.NewEmbed().
		SetDescription(fmt.Sprintf("Mutual Friends (%s - %s)", user.String(), user.ID)).
		AddField("Mutual Friends", mutualFriendsString).
		Send(ctx)
}

func (AdministrationCommands) mutualservers(ctx *dgc.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		_ = ctx.RespondText("Could not find user")
		return
	}
	mutualServers, err := ctx.Session.UserProfile(user.ID)
	if err != nil {
		_ = ctx.RespondText("Could not find user")
		return
	}
	mutualServersString := ""
	for _, guild_ := range mutualServers.MutualGuilds {
		guild, err := utils.Discord.GetGuild(ctx.Session, guild_.ID)
		if err != nil {
			continue
		}
		mutualServersString += fmt.Sprintf("%s - %s\n", guild.Name, guild.ID)
	}
	embed.NewEmbed().
		SetDescription(fmt.Sprintf("Mutual Servers (%s - %s)", user.String(), user.ID)).
		AddField("Mutual Servers", mutualServersString).
		Send(ctx)
}

func (AdministrationCommands) mobile(ctx *dgc.Ctx) {
	member, success := utils.Discord.TryGetMemberFromMention(ctx, ctx.Arguments.Get(0))
	if !success {
		_ = ctx.RespondText("Could not find user")
		return
	}
	embed.NewEmbed().
		SetDescription(fmt.Sprintf("%s - %s", member.User.String(), member.User.ID)).
		AddField("Is On Mobile", utils.Misc.FormatBool(member.Presence.ClientStatus["mobile"] != "")).
		Send(ctx)
}

func (AdministrationCommands) web(ctx *dgc.Ctx) {
	member, success := utils.Discord.TryGetMemberFromMention(ctx, ctx.Arguments.Get(0))
	if !success {
		_ = ctx.RespondText("Could not find user")
		return
	}
	embed.NewEmbed().
		SetDescription(fmt.Sprintf("%s - %s", member.User.String(), member.User.ID)).
		AddField("Is On Web", utils.Misc.FormatBool(member.Presence.ClientStatus["web"] != "")).
		Send(ctx)
}

func (AdministrationCommands) desktop(ctx *dgc.Ctx) {
	member, success := utils.Discord.TryGetMemberFromMention(ctx, ctx.Arguments.Get(0))
	if !success {
		_ = ctx.RespondText("Could not find user")
		return
	}
	embed.NewEmbed().
		SetDescription(fmt.Sprintf("%s - %s", member.User.String(), member.User.ID)).
		AddField("Is On Desktop", utils.Misc.FormatBool(member.Presence.ClientStatus["desktop"] != "")).
		Send(ctx)
}

func (AdministrationCommands) avatar(ctx *dgc.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		_ = ctx.RespondText("Could not find user")
		return
	}
	if user.Avatar == "" {
		_ = ctx.RespondText("User has no avatar")
		return
	}
	_ = ctx.RespondText(user.AvatarURL(""))
}

func (AdministrationCommands) banner(ctx *dgc.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		_ = ctx.RespondText("Could not find user")
		return
	}
	if user.Banner == "" {
		_ = ctx.RespondText("User has no banner")
		return
	}
	_ = ctx.RespondText(user.BannerURL(""))
}

func (AdministrationCommands) report(ctx *dgc.Ctx) {
	link := utils.Arguments["<message link>"]
	reason := utils.Arguments["<reason>"]

	// Verify message link
	url, err := url.Parse(link)
	if err != nil {
		_ = ctx.RespondText("Invalid message link.")
		return
	}
	urlPathParts := strings.Split(url.Path, "/")
	if len(urlPathParts) < 4 {
		_ = ctx.RespondText("Invalid message link.")
		return
	}
	guildID := urlPathParts[2]
	channelID := urlPathParts[3]
	messageID := urlPathParts[4]

	var reasonInt int
	switch reason {
	case "illegalcontent":
		reasonInt = 0
	case "harassment":
		reasonInt = 1
	case "spam":
	case "phishing":
		reasonInt = 2
	case "selfharm":
		reasonInt = 3
	case "nsfw":
		reasonInt = 4
	default:
		_ = ctx.RespondText("Argument `<reason>` must be one of: `illegalcontent`, `harassment`, `spam`, `phishing`, `selfharm`, `nsfw`")
		return
	}

	_, _ = api.ReportMessage(guildID, channelID, messageID, reasonInt)

	embed.NewEmbed().
		SetDescription("Successfully reported message!").
		Send(ctx)
}

func (AdministrationCommands) mute(ctx *dgc.Ctx) {
	if !utils.Discord.CheckIfInServer(ctx) {
		return
	}
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		_ = ctx.RespondText("Could not find user")
		return
	}

	member, _ := ctx.Session.GuildMember(ctx.Event.GuildID, user.ID)
	roles, err := ctx.Session.GuildRoles(ctx.Event.GuildID)
	if member == nil || err != nil {
		_ = ctx.RespondText("Error muting user: " + err.Error())
		return
	}

	for _, role := range roles {
		if strings.ToLower(role.Name) != "muted" {
			continue
		}
		if utils.SliceContains(member.Roles, role.ID) {
			embed.NewEmbed().
				SetDescription("User is already muted!").
				Send(ctx)
			return
		}
		err = ctx.Session.GuildMemberEdit(ctx.Event.GuildID, user.ID, append(member.Roles, role.ID))
		if err != nil {
			_ = ctx.RespondText("Error muting user: " + err.Error())
			return
		}
		embed.NewEmbed().
			SetDescription("Successfully muted user!").
			Send(ctx)
		return
	}

	channels, err := ctx.Session.GuildChannels(ctx.Event.GuildID)
	if err != nil {
		_ = ctx.RespondText("Error muting user: " + err.Error())
		return
	}

	role, err := ctx.Session.GuildRoleCreate(ctx.Event.GuildID, "Muted", 0, false, 0, false)
	if err != nil {
		_ = ctx.RespondText("Error muting user: " + err.Error())
		return
	}

	for _, channel := range channels {
		if channel.Type == 1 {
			err = ctx.Session.ChannelPermissionSet(channel.ID, role.ID, discordgo.PermissionOverwriteTypeRole, 0, discordgo.PermissionSendMessages)
			if err != nil {
				_ = ctx.RespondText("Error muting user: " + err.Error())
				return
			}
		}
	}

	err = ctx.Session.GuildMemberEdit(ctx.Event.GuildID, user.ID, append(member.Roles, role.ID))
	if err != nil {
		_ = ctx.RespondText("Error muting user: " + err.Error())
		return
	}

	embed.NewEmbed().
		SetDescription("Successfully muted user!").
		Send(ctx)
}

func (AdministrationCommands) unmute(ctx *dgc.Ctx) {
	if !utils.Discord.CheckIfInServer(ctx) {
		return
	}

	member, success := utils.Discord.TryGetMemberFromMention(ctx, ctx.Arguments.Get(0))
	if !success {
		_ = ctx.RespondText("Could not find user")
		return
	}

	roles, err := ctx.Session.GuildRoles(ctx.Event.GuildID)
	if member == nil || err != nil {
		_ = ctx.RespondText("Error unmuting user: " + err.Error())
		return
	}

	for _, role := range roles {
		if strings.ToLower(role.Name) != "muted" {
			continue
		}
		if !utils.SliceContains(member.Roles, role.ID) {
			embed.NewEmbed().
				SetDescription("User is already unmuted!").
				Send(ctx)
			return
		}
		for i, roleID := range member.Roles {
			if roleID == role.ID {
				member.Roles = append(member.Roles[:i], member.Roles[i+1:]...)
				break
			}
		}
		err = ctx.Session.GuildMemberEdit(ctx.Event.GuildID, member.User.ID, member.Roles)
		if err != nil {
			_ = ctx.RespondText("Error unmuting user: " + err.Error())
			return
		}
		embed.NewEmbed().
			SetDescription("Successfully unmuted user!").
			Send(ctx)
		return
	}
	_ = ctx.RespondText("Error: No muted role found!")
}

func (AdministrationCommands) timeout(ctx *dgc.Ctx) {
	if !utils.Discord.CheckIfInServer(ctx) {
		return
	}

	member, success := utils.Discord.TryGetMemberFromMention(ctx, ctx.Arguments.Get(0))
	if !success {
		_ = ctx.RespondText("Could not find user")
		return
	}

	seconds, err := strconv.Atoi(utils.Arguments["<seconds>"])
	if err != nil {
		_ = ctx.RespondText("Argument <seconds> must be a number.")
		return
	}
	until := time.Now().Add(time.Duration(seconds) * time.Second)

	err = ctx.Session.GuildMemberTimeout(ctx.Event.GuildID, member.User.ID, &until)
	if err != nil {
		_ = ctx.RespondText("Error timing out user: " + err.Error())
		return
	}

	embed.NewEmbed().
		SetDescription(fmt.Sprintf("Successfully timed out user for %d seconds!", seconds)).
		Send(ctx)
}

func (AdministrationCommands) removetimeout(ctx *dgc.Ctx) {
	if !utils.Discord.CheckIfInServer(ctx) {
		return
	}

	member, success := utils.Discord.TryGetMemberFromMention(ctx, ctx.Arguments.Get(0))
	if !success {
		_ = ctx.RespondText("Could not find user")
		return
	}

	err := ctx.Session.GuildMemberTimeout(ctx.Event.GuildID, member.User.ID, nil)
	if err != nil {
		_ = ctx.RespondText("Error removing timeout: " + err.Error())
		return
	}

	embed.NewEmbed().
		SetDescription("Successfully removed timeout!").
		Send(ctx)
}

func (AdministrationCommands) kick(ctx *dgc.Ctx) {
	if !utils.Discord.CheckIfInServer(ctx) {
		return
	}

	member, success := utils.Discord.TryGetMemberFromMention(ctx, ctx.Arguments.Get(0))
	if !success {
		_ = ctx.RespondText("Could not find user")
		return
	}

	err := ctx.Session.GuildMemberDeleteWithReason(ctx.Event.GuildID, member.User.ID, utils.Arguments["[reason]"])
	if err != nil {
		_ = ctx.RespondText("Error kicking user: " + err.Error())
		return
	}

	embed.NewEmbed().
		SetDescription("Successfully kicked user!").
		Send(ctx)
}

func (AdministrationCommands) ban(ctx *dgc.Ctx) {
	if !utils.Discord.CheckIfInServer(ctx) {
		return
	}

	member, success := utils.Discord.TryGetMemberFromMention(ctx, ctx.Arguments.Get(0))
	memberID := ctx.Arguments.Get(0).Raw()
	if success {
		memberID = member.User.ID
	}

	err := ctx.Session.GuildBanCreateWithReason(ctx.Event.GuildID, memberID, utils.Arguments["[reason]"], 0)
	if err != nil {
		_ = ctx.RespondText("Error banning user: " + err.Error())
		return
	}

	embed.NewEmbed().
		SetDescription("Successfully banned user!").
		Send(ctx)
}

func (AdministrationCommands) softban(ctx *dgc.Ctx) {
	if !utils.Discord.CheckIfInServer(ctx) {
		return
	}

	member, success := utils.Discord.TryGetMemberFromMention(ctx, ctx.Arguments.Get(0))
	if !success {
		_ = ctx.RespondText("Could not find user")
		return
	}

	err := ctx.Session.GuildBanCreateWithReason(ctx.Event.GuildID, member.User.ID, utils.Arguments["[reason]"], 0)
	if err != nil {
		_ = ctx.RespondText("Error softbanning user: " + err.Error())
		return
	}
	err = ctx.Session.GuildBanDelete(ctx.Event.GuildID, member.User.ID)
	if err != nil {
		_ = ctx.RespondText("Error softbanning user: " + err.Error())
		return
	}

	embed.NewEmbed().
		SetDescription("Successfully softbanned user!").
		Send(ctx)
}

func (AdministrationCommands) unban(ctx *dgc.Ctx) {
	if !utils.Discord.CheckIfInServer(ctx) {
		return
	}

	err := ctx.Session.GuildBanDelete(ctx.Event.GuildID, utils.Arguments["<user id>"])
	if err != nil {
		_ = ctx.RespondText("Error unbanning user: " + err.Error())
		return
	}

	embed.NewEmbed().
		SetDescription("Successfully softbanned user!").
		Send(ctx)
}

func (AdministrationCommands) banlist(ctx *dgc.Ctx) {
	guildID := utils.Arguments["[server id]"]
	if guildID == "" {
		guildID = ctx.Event.GuildID
	}
	if guildID == "" {
		_ = ctx.RespondText("Error: No server ID specified!")
		return
	}
	banList, err := ctx.Session.GuildBans(ctx.Event.GuildID, 1000, "", "")
	if err != nil {
		_ = ctx.RespondText("Error getting ban list: " + err.Error())
		return
	}
	banListString := ""
	for i, ban := range banList {
		// Max 40
		if i >= 40 {
			break
		}
		banListString += fmt.Sprintf("%s - %s\n", ban.User.Username, ban.User.ID)
	}
	embed.NewEmbed().
		AddField("Banned Amount", fmt.Sprint(len(banList))).
		AddField("Ban List", banListString).
		Send(ctx)
}

func (AdministrationCommands) copybans(ctx *dgc.Ctx) {
	sourceServerID := utils.Arguments["<source server id>"]
	destinationServerID := utils.Arguments["[dest server id]"]
	if destinationServerID == "" {
		destinationServerID = ctx.Event.GuildID
	}
	if destinationServerID == "" {
		_ = ctx.RespondText("Error: No server ID specified!")
		return
	}
	banList, err := ctx.Session.GuildBans(sourceServerID, 1000, "", "")
	if err != nil {
		_ = ctx.RespondText("Error copying bans: " + err.Error())
		return
	}

	for _, ban := range banList {
		err := ctx.Session.GuildBanCreateWithReason(destinationServerID, ban.User.ID, ban.Reason, 0)
		if err != nil {
			_ = ctx.RespondText("Error copying bans: " + err.Error())
			return
		}
	}
}

// Nickname Control

func (AdministrationCommands) nick(ctx *dgc.Ctx) {
	if !utils.Discord.CheckIfInServer(ctx) {
		return
	}
	err := ctx.Session.GuildMemberNickname(ctx.Event.GuildID, "@me", ctx.Arguments.Raw())
	if err != nil {
		_ = ctx.RespondText("Error changing nickname: " + err.Error())
		return
	}
	embed.NewEmbed().
		SetDescription("Successfully changed nickname!").
		Send(ctx)
}

func (AdministrationCommands) nickuser(ctx *dgc.Ctx) {
	if !utils.Discord.CheckIfInServer(ctx) {
		return
	}

	member, success := utils.Discord.TryGetMemberFromMention(ctx, ctx.Arguments.Get(0))
	if !success {
		_ = ctx.RespondText("Could not find user")
		return
	}

	err := ctx.Session.GuildMemberNickname(ctx.Event.GuildID, strings.ReplaceAll(member.User.ID, ctx.Event.Author.ID, "@me"), utils.String.RemoveWords(ctx.Arguments.Raw(), 1))
	if err != nil {
		_ = ctx.RespondText("Error changing nickname: " + err.Error())
		return
	}
	embed.NewEmbed().
		SetDescription("Successfully changed nickname!").
		Send(ctx)
}

func (AdministrationCommands) clearnick(ctx *dgc.Ctx) {
	if !utils.Discord.CheckIfInServer(ctx) {
		return
	}

	userID := "@me"
	member, success := utils.Discord.TryGetMemberFromMention(ctx, ctx.Arguments.Get(0))
	if success {
		userID = member.User.ID
	}
	err := ctx.Session.GuildMemberNickname(ctx.Event.GuildID, userID, "")
	if err != nil {
		_ = ctx.RespondText("Error changing nickname: " + err.Error())
		return
	}
	embed.NewEmbed().
		SetDescription("Successfully changed nickname!").
		Send(ctx)
}

// Server Control

func (AdministrationCommands) serverinfo(ctx *dgc.Ctx) {
	if !utils.Discord.CheckIfInServer(ctx) {
		return
	}

	guild, err := utils.Discord.GetGuild(ctx.Session, utils.Misc.Or(utils.Arguments["[server id]"], ctx.Event.GuildID))
	if err != nil {
		_ = ctx.RespondText("Error getting server: " + err.Error())
		return
	}

	owner, err := utils.Discord.GetMember(ctx.Session, guild.ID, guild.OwnerID)
	ownerMention := "Error"
	if err == nil {
		ownerMention = owner.User.String()
	}

	embed.NewEmbed().
		SetDescription(fmt.Sprintf("Server Info - %s", guild.Name)).
		AddField("ID", guild.ID).
		AddField("Owner", ownerMention).
		AddField("Created At", utils.Discord.SnowflakeToTime(guild.ID).Format(time.RFC1123)).
		AddField("Verification Level", verificationLevels[int(guild.VerificationLevel)]).
		AddField("Boost Level", boostLevels[int(guild.PremiumTier)]).
		AddField("Amount of boosts", fmt.Sprint(guild.PremiumSubscriptionCount)).
		AddField("Member Count", fmt.Sprint(guild.MemberCount)).
		AddField("Channel Count", fmt.Sprint(len(guild.Channels))).
		AddField("Role Count", fmt.Sprint(len(guild.Roles))).
		AddField("Joined At", guild.JoinedAt.Format(time.RFC1123)).
		SetImage(guild.IconURL()).
		Send(ctx)
}

func (AdministrationCommands) serverage(ctx *dgc.Ctx) {
	guild, err := utils.Discord.GetGuild(ctx.Session, utils.Misc.Or(utils.Arguments["[server id]"], ctx.Event.GuildID))
	if err != nil {
		_ = ctx.RespondText("Error getting server: " + err.Error())
		return
	}
	age := time.Since(utils.Discord.SnowflakeToTime(guild.ID))
	embed.NewEmbed().
		SetDescription(fmt.Sprintf("Server Age (%s - %s)", guild.Name, guild.ID)).
		AddField("Age", utils.Misc.FormatDuration(age)).
		Send(ctx)
}

func (AdministrationCommands) serverowner(ctx *dgc.Ctx) {
	if !utils.Discord.CheckIfInServer(ctx) {
		return
	}

	guild, err := utils.Discord.GetGuild(ctx.Session, utils.Misc.Or(utils.Arguments["[server id]"], ctx.Event.GuildID))
	if err != nil {
		_ = ctx.RespondText("Error getting server: " + err.Error())
		return
	}
	owner, err := utils.Discord.GetMember(ctx.Session, guild.ID, guild.OwnerID)
	if err != nil {
		_ = ctx.RespondText("Error getting server owner: " + err.Error())
		return
	}
	embed.NewEmbed().
		SetDescription(guild.Name).
		AddField("Owner", owner.Mention()).
		Send(ctx)
}

func (AdministrationCommands) servericon(ctx *dgc.Ctx) {
	if !utils.Discord.CheckIfInServer(ctx) {
		return
	}

	guild, err := utils.Discord.GetGuild(ctx.Session, utils.Misc.Or(utils.Arguments["[server id]"], ctx.Event.GuildID))
	if err != nil {
		_ = ctx.RespondText("Error getting server: " + err.Error())
		return
	}

	embed := embed.NewEmbed().
		SetDescription(fmt.Sprintf("Server Icon - %s", guild.Name)).
		SetImage(guild.IconURL())
	if guild.Icon == "" {
		embed.AddField("Error", "Server has no icon")
	}
	embed.Send(ctx)
}

func (AdministrationCommands) serverbanner(ctx *dgc.Ctx) {
	if !utils.Discord.CheckIfInServer(ctx) {
		return
	}
	guild, err := utils.Discord.GetGuild(ctx.Session, utils.Misc.Or(utils.Arguments["[server id]"], ctx.Event.GuildID))
	if err != nil {
		_ = ctx.RespondText("Error getting server: " + err.Error())
		return
	}
	if guild.Banner == "" {
		embed.NewEmbed().AddField("Error", "Server has no banner!").Send(ctx)
	}
	embed.NewEmbed().
		SetDescription(fmt.Sprintf("Server Banner - %s", guild.Name)).
		SetImage(guild.BannerURL()).
		Send(ctx)
}

func (AdministrationCommands) serversplash(ctx *dgc.Ctx) {
	if !utils.Discord.CheckIfInServer(ctx) {
		return
	}
	guild, err := utils.Discord.GetGuild(ctx.Session, utils.Misc.Or(utils.Arguments["[server id]"], ctx.Event.GuildID))
	if err != nil {
		_ = ctx.RespondText("Error getting server: " + err.Error())
		return
	}
	if guild.Splash == "" {
		_ = ctx.RespondText("Server has no invite splash!")
		return
	}
	embed.NewEmbed().
		SetDescription(fmt.Sprintf("Server Invite Splash - %s", guild.Name)).
		SetImage(discordgo.EndpointGuildSplash(guild.ID, guild.Splash)).
		Send(ctx)
}

func (AdministrationCommands) serverdsplash(ctx *dgc.Ctx) {
	if !utils.Discord.CheckIfInServer(ctx) {
		return
	}
	guild, err := utils.Discord.GetGuild(ctx.Session, utils.Misc.Or(utils.Arguments["[server id]"], ctx.Event.GuildID))
	if err != nil {
		_ = ctx.RespondText("Error getting server: " + err.Error())
		return
	}
	if guild.DiscoverySplash == "" {
		_ = ctx.RespondText("Server has no discovery splash!")
		return
	}
	embed.NewEmbed().
		SetDescription(fmt.Sprintf("Server Discovery Splash - %s", guild.Name)).
		SetImage(discordgo.EndpointGuildSplash(guild.ID, guild.DiscoverySplash)).
		Send(ctx)
}

func (AdministrationCommands) serverboosters(ctx *dgc.Ctx) {
	if !utils.Discord.CheckIfInServer(ctx) {
		return
	}

	guild, err := utils.Discord.GetGuild(ctx.Session, utils.Misc.Or(utils.Arguments["[server id]"], ctx.Event.GuildID))
	if err != nil {
		_ = ctx.RespondText("Error getting server: " + err.Error())
		return
	}

	boosters := utils.Discord.GetBoosters(guild)
	boostersString := ""

	for _, booster := range boosters {
		boostersString += fmt.Sprintf("%s - %s\n", booster.User.String(), booster.User.ID)
	}

	if boostersString == "" {
		boostersString = "No boosters"
	}

	embed.NewEmbed().
		SetDescription(fmt.Sprintf("Server Boosters - %s", guild.Name)).
		AddField("Boosters", boostersString).
		Send(ctx)
}

func (AdministrationCommands) serverboostlevel(ctx *dgc.Ctx) {
	if !utils.Discord.CheckIfInServer(ctx) {
		return
	}

	guild, err := utils.Discord.GetGuild(ctx.Session, utils.Misc.Or(utils.Arguments["[server id]"], ctx.Event.GuildID))
	if err != nil {
		_ = ctx.RespondText("Error getting server: " + err.Error())
		return
	}

	embed.NewEmbed().
		SetDescription(fmt.Sprintf("Server Boost Level - %s", guild.Name)).
		AddField("Boost Level", boostLevels[int(guild.PremiumTier)]).
		Send(ctx)
}

func (AdministrationCommands) setservername(ctx *dgc.Ctx) {
	if !utils.Discord.CheckIfInServer(ctx) {
		return
	}
	guild, err := ctx.Session.GuildEdit(ctx.Event.GuildID, discordgo.GuildParams{
		Name: ctx.Arguments.Raw(),
	})
	if err != nil {
		_ = ctx.RespondText("Error changing server name: " + err.Error())
		return
	}
	embed.NewEmbed().
		SetDescription(fmt.Sprintf("Successfully changed server name to %s!", guild.Name)).
		Send(ctx)
}

func (AdministrationCommands) setservericon(ctx *dgc.Ctx) {
	if !utils.Discord.CheckIfInServer(ctx) {
		return
	}
	_, err := ctx.Session.GuildEdit(ctx.Event.GuildID, discordgo.GuildParams{
		Icon: utils.Arguments["<url>"],
	})
	if err != nil {
		_ = ctx.RespondText("Error changing server icon: " + err.Error())
		return
	}
	embed.NewEmbed().
		SetDescription("Successfully changed server icon!").
		Send(ctx)
}

func (AdministrationCommands) setserverbanner(ctx *dgc.Ctx) {
	if !utils.Discord.CheckIfInServer(ctx) {
		return
	}
	_, err := ctx.Session.GuildEdit(ctx.Event.GuildID, discordgo.GuildParams{
		Banner: utils.Arguments["<url>"],
	})
	if err != nil {
		_ = ctx.RespondText("Error changing server icon: " + err.Error())
		return
	}
	embed.NewEmbed().
		SetDescription("Successfully changed server icon!").
		Send(ctx)
}

func (AdministrationCommands) lockserver(ctx *dgc.Ctx) {
	guild, err := utils.Discord.GetGuild(ctx.Session, utils.Misc.Or(utils.Arguments["[server id]"], ctx.Event.GuildID))
	if err != nil {
		_ = ctx.RespondText("Error getting server: " + err.Error())
		return
	}

	for _, channel := range guild.Channels {
		if channel.Type == discordgo.ChannelTypeGuildText {
			_, err := ctx.Session.ChannelEditComplex(channel.ID, &discordgo.ChannelEdit{
				PermissionOverwrites: []*discordgo.PermissionOverwrite{
					{
						ID:   guild.ID,
						Type: discordgo.PermissionOverwriteTypeRole,
						Deny: discordgo.PermissionSendMessages,
					},
				},
			})
			if err != nil {
				_ = ctx.RespondText("Error locking server: " + err.Error())
				return
			}
		}
	}

	embed.NewEmbed().
		SetDescription("Successfully locked server!").
		Send(ctx)
}

func (AdministrationCommands) unlockserver(ctx *dgc.Ctx) {
	guild, err := utils.Discord.GetGuild(ctx.Session, utils.Misc.Or(utils.Arguments["[server id]"], ctx.Event.GuildID))
	if err != nil {
		_ = ctx.RespondText("Error getting server: " + err.Error())
		return
	}

	for _, channel := range guild.Channels {
		if channel.Type == discordgo.ChannelTypeGuildText {
			_, err := ctx.Session.ChannelEditComplex(channel.ID, &discordgo.ChannelEdit{
				PermissionOverwrites: []*discordgo.PermissionOverwrite{
					{
						ID:   guild.ID,
						Type: discordgo.PermissionOverwriteTypeRole,
						Deny: 0,
					},
				},
			})
			if err != nil {
				_ = ctx.RespondText("Error locking server: " + err.Error())
				return
			}
		}
	}

	embed.NewEmbed().
		SetDescription("Successfully locked server!").
		Send(ctx)
}

func (AdministrationCommands) channels(ctx *dgc.Ctx) {
	guild, err := utils.Discord.GetGuild(ctx.Session, utils.Misc.Or(utils.Arguments["[server id]"], ctx.Event.GuildID))
	if err != nil {
		_ = ctx.RespondText("Error getting server: " + err.Error())
		return
	}

	categories, categoryChannels := utils.Discord.MapCategories(guild.Channels)
	embed := embed.NewEmbed().SetDescription(fmt.Sprintf("Server Channels - %s", guild.Name))
	for _, category := range categories {
		channelsString := ""
		for _, channel := range categoryChannels[category.ID] {
			channelsString += channel.Mention() + "\n"
		}
		embed.AddField("Category: "+category.Name, channelsString)
	}
	embed.Send(ctx)
}

func (AdministrationCommands) bots(ctx *dgc.Ctx) {
	guild, err := utils.Discord.GetGuild(ctx.Session, utils.Misc.Or(utils.Arguments["[server id]"], ctx.Event.GuildID))
	if err != nil {
		_ = ctx.RespondText("Error getting server: " + err.Error())
		return
	}
	fmt.Println(len(guild.Members))
	botString := ""
	for _, member := range guild.Members {
		if member.User.Bot {
			botString += member.User.Mention() + "\n"
		}
	}
	botString = utils.Misc.Or(botString, "No bots found!")
	embed.NewEmbed().
		SetDescription(fmt.Sprintf("Server Bots - %s", guild.Name)).
		AddField("Server Bots", botString).
		Send(ctx)
}

func (AdministrationCommands) admins(ctx *dgc.Ctx) {
	guild, err := utils.Discord.GetGuild(ctx.Session, utils.Misc.Or(utils.Arguments["[server id]"], ctx.Event.GuildID))
	if err != nil {
		_ = ctx.RespondText("Error getting server: " + err.Error())
		return
	}
	if len(guild.Channels) == 0 {
		_ = ctx.RespondText("Error getting server: " + err.Error())
		return
	}
	adminString := ""
	for _, member := range guild.Members {
		permissions, _ := ctx.Session.UserChannelPermissions(member.User.ID, guild.Channels[0].ID)
		if permissions&discordgo.PermissionAdministrator != 0 && !member.User.Bot {
			adminString += fmt.Sprintf("%s - %s\n", member.User.String(), member.User.ID)
		}
	}
	adminString = utils.Misc.Or(adminString, "No admins found!")
	embed.NewEmbed().
		SetDescription(fmt.Sprintf("Server Admins - %s", guild.Name)).
		AddField("Server Admins", adminString).
		Send(ctx)
}

func (AdministrationCommands) roles(ctx *dgc.Ctx) {
	guild, err := utils.Discord.GetGuild(ctx.Session, utils.Misc.Or(utils.Arguments["[server id]"], ctx.Event.GuildID))
	if err != nil {
		_ = ctx.RespondText("Error getting server: " + err.Error())
		return
	}
	sort.Slice(guild.Roles, func(i, j int) bool {
		return guild.Roles[i].Position < guild.Roles[j].Position
	})
	roleString := ""
	for _, role := range guild.Roles {
		if role.ID == guild.ID {
			// Ignore @everyone role
			continue
		}
		roleString += fmt.Sprintf("%s - %s\n", role.Name, role.ID)
	}
	roleString = utils.Misc.Or(roleString, "No roles found!")
	embed.NewEmbed().
		SetDescription(fmt.Sprintf("Server Roles - %s", guild.Name)).
		AddField("Server Roles", roleString).
		Send(ctx)
}

// Invite Control

func (AdministrationCommands) invite(ctx *dgc.Ctx) {
	channel, success := utils.Discord.TryGetChannelFromMention(ctx, ctx.Arguments.Get(0))
	channelID := ctx.Event.ChannelID
	if !success {
		channelID = channel.ID
	}

	ageInt, err := strconv.Atoi(utils.Arguments["[max age]"])
	if err != nil && utils.Arguments["[max age]"] != "" {
		_ = ctx.RespondText("Argument [max age] must be a number.")
		return
	}

	usesInt, err := strconv.Atoi(utils.Arguments["[max uses]"])
	if err != nil && utils.Arguments["[max uses]"] != "" {
		_ = ctx.RespondText("Argument [max uses] must be a number.")
		return
	}

	invite, err := ctx.Session.ChannelInviteCreate(channelID, discordgo.Invite{
		MaxAge:    ageInt,
		MaxUses:   usesInt,
		Temporary: false,
		Unique:    true,
	})
	if err != nil {
		_ = ctx.RespondText("Error creating invite: " + err.Error())
		return
	}

	embed.NewEmbed().
		SetDescription(fmt.Sprintf("Created invite: https://discord.gg/%s", invite.Code)).
		Send(ctx)
}

func (AdministrationCommands) friendinvite(ctx *dgc.Ctx) {
	invite, err := ctx.Session.CreateFriendInvite("@me")
	if err != nil {
		_ = ctx.RespondText("Error creating friend invite: " + err.Error())
		return
	}

	embed.NewEmbed().
		SetDescription("Successfully created friend invite!").
		AddField("Invite", fmt.Sprintf("https://discord.gg/%s", invite.Code)).
		Send(ctx)
}

func (AdministrationCommands) inviteinfo(ctx *dgc.Ctx) {
	parsedInvite := utils.Discord.ParseInvite(utils.Arguments["<invite>"])
	invite, err := ctx.Session.InviteWithCounts(parsedInvite)
	if err != nil {
		_ = ctx.RespondText("Error getting invite: " + err.Error())
		return
	}

	embed := embed.NewEmbed().SetDescription(fmt.Sprintf("Invite Info - %s", invite.Code))

	if invite.Guild != nil {
		embed.AddField("Server", invite.Guild.Name)
		if invite.Guild.Icon != "" {
			embed.SetImage(invite.Guild.IconURL())
		}
	}
	if invite.Channel != nil {
		embed.AddField("Channel", invite.Channel.Name)
	}
	if invite.Inviter != nil {
		embed.AddField("Inviter", invite.Inviter.String())
	}
	if invite.ExpiresAt != nil {
		embed.AddField("Expires At", invite.ExpiresAt.Format(time.RFC1123))
	}
	embed.AddField("Max Age", fmt.Sprint(invite.MaxAge)).
		AddField("Uses", fmt.Sprint(invite.Uses)).
		AddField("Max Uses", fmt.Sprint(invite.MaxUses)).
		AddField("Member Count", fmt.Sprint(invite.ApproximateMemberCount)).
		AddField("Online Member Count", fmt.Sprint(invite.ApproximatePresenceCount)).
		Send(ctx)
}

func (AdministrationCommands) deleteinvite(ctx *dgc.Ctx) {
	parsedInvite := utils.Discord.ParseInvite(utils.Arguments["<invite>"])
	invite, err := ctx.Session.InviteWithCounts(parsedInvite)
	if err != nil {
		_ = ctx.RespondText("Error getting invite: " + err.Error())
		return
	}

	invite, err = ctx.Session.InviteDelete(invite.Code)
	if err != nil {
		_ = ctx.RespondText("Error deleting invite: " + err.Error())
		return
	}

	embed.NewEmbed().
		SetDescription(fmt.Sprintf("Deleted invite: %s", invite.Code)).
		Send(ctx)
}

func (m *AdministrationCommands) serverinvites(ctx *dgc.Ctx) {
	invites, err := ctx.Session.GuildInvites(utils.Misc.Or(utils.Arguments["[server id]"], ctx.Event.GuildID))
	if err != nil {
		_ = ctx.RespondText("Error getting invites: " + err.Error())
		return
	}
	invitesString := ""
	for _, invite := range invites {
		username := "Unknown"
		if invite.Inviter != nil {
			username = invite.Inviter.String()
		}
		invitesString += fmt.Sprintf(
			"Link: https://discord.gg/%s - Expires At: %s - Created By: %s\n",
			invite.Code,
			utils.Discord.GetTimestamp(invite.ExpiresAt, "Never"),
			username,
		)
	}
	embed.NewEmbed().
		SetDescription("Invite List").
		AddField("Invites", invitesString).
		Send(ctx)
}

func (m *AdministrationCommands) channelinvites(ctx *dgc.Ctx) {
	channel, success := utils.Discord.TryGetChannelFromMention(ctx, ctx.Arguments.Get(0))
	channelID := ctx.Event.ChannelID
	if success {
		channelID = channel.ID
	}
	invites, err := ctx.Session.ChannelInvites(channelID)
	if err != nil {
		_ = ctx.RespondText("Error getting invites: " + err.Error())
		return
	}
	invitesString := ""
	for _, invite := range invites {
		username := "Unknown"
		if invite.Inviter != nil {
			username = invite.Inviter.String()
		}
		invitesString += fmt.Sprintf(
			"Link: https://discord.gg/%s - Expires At: %s - Created By: %s\n",
			invite.Code,
			utils.Discord.GetTimestamp(invite.ExpiresAt, "Never"),
			username,
		)
	}
	embed.NewEmbed().
		SetDescription("Invite List").
		AddField("Invites", invitesString).
		Send(ctx)
}

func (m *AdministrationCommands) userinvites(ctx *dgc.Ctx) {
	if !utils.Discord.CheckIfInServer(ctx) {
		return
	}

	member, success := utils.Discord.TryGetMemberFromMention(ctx, ctx.Arguments.Get(0))
	if !success {
		_ = ctx.RespondText("Could not find user")
		return
	}

	invites, err := ctx.Session.GuildInvites(utils.Misc.Or(utils.Arguments["[server id]"], ctx.Event.GuildID))
	if err != nil {
		_ = ctx.RespondText("Error getting invites: " + err.Error())
		return
	}
	invitesString := ""
	for _, invite := range invites {
		if invite.Inviter == nil {
			continue
		}
		if invite.Inviter.ID != member.User.ID {
			continue
		}
		invitesString += fmt.Sprintf(
			"Link: https://discord.gg/%s - Expires At: %s\n",
			invite.Code,
			utils.Discord.GetTimestamp(invite.ExpiresAt, "Never"),
		)
	}
	embed.NewEmbed().
		SetDescription(fmt.Sprintf("Invite List (%s)", member.User.String())).
		AddField("Invites", invitesString).
		Send(ctx)
}

func (AdministrationCommands) roleinfo(ctx *dgc.Ctx) {
	if !utils.Discord.CheckIfInServer(ctx) {
		return
	}
	role, success := utils.Discord.TryGetRoleFromMention(ctx, ctx.Arguments.Get(0))
	if !success {
		_ = ctx.RespondText("Could not find role")
		return
	}
	// Get all members in guild and check if they have the role
	guild, err := utils.Discord.GetGuild(ctx.Session, ctx.Event.GuildID)
	if err != nil {
		_ = ctx.RespondText("Error getting server: " + err.Error())
		return
	}
	memberCount := 0
	for _, member := range guild.Members {
		if utils.SliceContains(member.Roles, role.ID) {
			memberCount++
		}
	}
	createdAt := utils.Discord.SnowflakeToTime(role.ID)
	perms, _ := api.CreatePaste(utils.Discord.FormatPermissions(role.Name, role.ID, role.Permissions))
	embed.NewEmbed().
		SetDescription(fmt.Sprintf("Role Info (%s)", role.Name)).
		AddField("ID", role.ID).
		AddField("Color", fmt.Sprintf("#%x", role.Color)).
		AddField("Created at", utils.Discord.GetTimestamp(&createdAt, "Error")).
		AddField("Position", fmt.Sprint(role.Position)).
		AddField("Is Hoist", utils.Misc.FormatBool(role.Hoist)).
		AddField("Is Mentionable", utils.Misc.FormatBool(role.Mentionable)).
		AddField("Members With Role", fmt.Sprint(memberCount)).
		AddField("Permissions", perms).
		Send(ctx)
}

func (AdministrationCommands) giverole(ctx *dgc.Ctx) {
	if !utils.Discord.CheckIfInServer(ctx) {
		return
	}
	member, success := utils.Discord.TryGetMemberFromMention(ctx, ctx.Arguments.Get(0))
	if !success {
		_ = ctx.RespondText("Could not find user")
		return
	}
	role, success := utils.Discord.TryGetRoleFromMention(ctx, ctx.Arguments.Get(1))
	if !success {
		_ = ctx.RespondText("Could not find role")
		return
	}

	err := ctx.Session.GuildMemberRoleAdd(ctx.Event.GuildID, member.User.ID, role.ID)
	if err != nil {
		_ = ctx.RespondText("Error giving role: " + err.Error())
		return
	}
	embed.NewEmbed().
		SetDescription(fmt.Sprintf("Successfully gave role (%s) to (%s)", role.Name, member.User.String())).
		Send(ctx)
}

func (AdministrationCommands) removerole(ctx *dgc.Ctx) {
	if !utils.Discord.CheckIfInServer(ctx) {
		return
	}
	member, success := utils.Discord.TryGetMemberFromMention(ctx, ctx.Arguments.Get(0))
	if !success {
		_ = ctx.RespondText("Could not find user")
		return
	}
	role, success := utils.Discord.TryGetRoleFromMention(ctx, ctx.Arguments.Get(1))
	if !success {
		_ = ctx.RespondText("Could not find role")
		return
	}

	err := ctx.Session.GuildMemberRoleRemove(ctx.Event.GuildID, member.User.ID, role.ID)
	if err != nil {
		_ = ctx.RespondText("Error removing role: " + err.Error())
		return
	}
	embed.NewEmbed().
		SetDescription(fmt.Sprintf("Successfully removed role (%s) from (%s)", role.Name, member.User.String())).
		Send(ctx)
}

func (AdministrationCommands) createrole(ctx *dgc.Ctx) {
	if !utils.Discord.CheckIfInServer(ctx) {
		return
	}
	role, err := ctx.Session.GuildRoleCreate(ctx.Event.GuildID, ctx.Arguments.Raw(), 0, false, 0, false)
	if err != nil {
		_ = ctx.RespondText("Error creating role: " + err.Error())
		return
	}
	embed.NewEmbed().
		SetDescription(fmt.Sprintf("Successfully created role (%s)", role.Name)).
		Send(ctx)
}

func (AdministrationCommands) deleterole(ctx *dgc.Ctx) {
	if !utils.Discord.CheckIfInServer(ctx) {
		return
	}
	role, success := utils.Discord.TryGetRoleFromMention(ctx, ctx.Arguments.Get(0))
	if !success {
		_ = ctx.RespondText("Could not find role")
		return
	}
	err := ctx.Session.GuildRoleDelete(ctx.Event.GuildID, role.ID)
	if err != nil {
		_ = ctx.RespondText("Error deleting role: " + err.Error())
		return
	}
	embed.NewEmbed().
		SetDescription(fmt.Sprintf("Successfully deleted role (%s)", role.Name)).
		Send(ctx)
}

func (AdministrationCommands) renamerole(ctx *dgc.Ctx) {
	if !utils.Discord.CheckIfInServer(ctx) {
		return
	}
	role, success := utils.Discord.TryGetRoleFromMention(ctx, ctx.Arguments.Get(0))
	if !success {
		_ = ctx.RespondText("Could not find role")
		return
	}
	role, err := ctx.Session.GuildRoleEdit(ctx.Event.GuildID, role.ID, utils.String.RemoveWords(ctx.Arguments.Raw(), 1), role.Color, role.Hoist, role.Permissions, role.Mentionable)
	if err != nil {
		_ = ctx.RespondText("Error renaming role: " + err.Error())
		return
	}
	embed.NewEmbed().
		SetDescription(fmt.Sprintf("Successfully renamed role (%s)", role.Name)).
		Send(ctx)
}
