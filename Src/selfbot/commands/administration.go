package commands

import (
	"fmt"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/LightningDev1/LB-Selfbot/selfbot/api"
	"github.com/LightningDev1/LB-Selfbot/selfbot/embed"
	"github.com/LightningDev1/LB-Selfbot/selfbot/http"
	"github.com/LightningDev1/LB-Selfbot/selfbot/utils"
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

func (m *AdministrationCommands) RegisterCommands(router *dgc.Router) {
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
		Flags:       []string{"DANGEROUS"},
		Handler:     m.textchannel,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "voicechannel",
		Description: "Create a voice channel",
		Usage:       "[p]voicechannel <name>",
		Flags:       []string{"DANGEROUS"},
		Handler:     m.voicechannel,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "stagechannel",
		Description: "Create a stage channel",
		Usage:       "[p]stagechannel <name>",
		Flags:       []string{"DANGEROUS"},
		Handler:     m.stagechannel,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "newschannel",
		Description: "Create a news channel",
		Usage:       "[p]newschannel <name>",
		Flags:       []string{"DANGEROUS"},
		Handler:     m.newschannel,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "thread",
		Description: "Create a thread",
		Usage:       "[p]thread <name>",
		Flags:       []string{"DANGEROUS"},
		Handler:     m.thread,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "category",
		Description: "Create a category",
		Usage:       "[p]category <name>",
		Flags:       []string{"DANGEROUS"},
		Handler:     m.category,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "deletecategory",
		Description: "Delete a category",
		Usage:       "[p]deletecategory <category id>",
		Flags:       []string{"DANGEROUS"},
		Handler:     m.deletecategory,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "renamechannel",
		Description: "Rename a channel",
		Usage:       "[p]renamechannel <channel> <name>",
		Flags:       []string{"DANGEROUS"},
		Handler:     m.renamechannel,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "deletechannel",
		Description: "Delete a channel",
		Usage:       "[p]deletechannel <channel>",
		Flags:       []string{"DANGEROUS"},
		Handler:     m.deletechannel,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "slowmode",
		Description: "Change slowmode amount",
		Usage:       "[p]slowmode <channel> <seconds>",
		Flags:       []string{"DANGEROUS"},
		Handler:     m.slowmode,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "removeslowmode",
		Description: "Remove slowmode",
		Usage:       "[p]removeslowmode <channel>",
		Flags:       []string{"DANGEROUS"},
		Handler:     m.removeslowmode,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "lockchannel",
		Description: "Lock a channel",
		Usage:       "[p]lockchannel <channel>",
		Flags:       []string{"DANGEROUS"},
		Handler:     m.lockchannel,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "unlockchannel",
		Description: "Unlock a channel",
		Usage:       "[p]unlockchannel <channel>",
		Flags:       []string{"DANGEROUS"},
		Handler:     m.unlockchannel,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "purge",
		Description: "Delete messages in a channel",
		Usage:       "[p]purge <amount> [channel]",
		Flags:       []string{"DANGEROUS"},
		Handler:     m.purge,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "purgeself",
		Description: "Delete your own messages in a channel",
		Usage:       "[p]purgeself <amount> [channel]",
		Flags:       []string{"DANGEROUS"},
		Handler:     m.purgeself,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "purgeuser",
		Description: "Delete a users messages in a channel",
		Usage:       "[p]purgeuser <amount> <user> [channel]",
		Flags:       []string{"DANGEROUS"},
		Handler:     m.purgeuser,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "purgerange",
		Description: "Purge a range of messages in a channel",
		Usage:       "[p]purgerange <start message> <end message> [channel]",
		Flags:       []string{"DANGEROUS"},
		Handler:     m.purgerange,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "nuke",
		Description: "Nuke a channel",
		Usage:       "[p]nuke [channel]",
		Flags:       []string{"DANGEROUS"},
		Handler:     m.nuke,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "ackchannel",
		Description: "Ack a channel",
		Usage:       "[p]ackchannel [channel]",
		Handler:     m.ackchannel,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "ackmessage",
		Description: "Ack a message",
		Usage:       "[p]ackmessage <channel> <message id>",
		Handler:     m.ackmessage,
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
		Usage:       "[p]userinfo [user]",
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
		Name:        "platform",
		Description: "See the platforms a user is online on",
		Usage:       "[p]platform [user]",
		Handler:     m.platform,
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
		Flags:       []string{"DANGEROUS"},
		Handler:     m.report,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "mute",
		Description: "Mute a user",
		Usage:       "[p]mute <user>",
		Flags:       []string{"DANGEROUS"},
		Handler:     m.mute,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "unmute",
		Description: "Unmute a user",
		Usage:       "[p]unmute <user>",
		Flags:       []string{"DANGEROUS"},
		Handler:     m.unmute,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "timeout",
		Description: "Timeout a user",
		Usage:       "[p]timeout <user> <seconds>",
		Flags:       []string{"DANGEROUS"},
		Handler:     m.timeout,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "removetimeout",
		Description: "Remove a timeout",
		Usage:       "[p]removetimeout <user>",
		Flags:       []string{"DANGEROUS"},
		Handler:     m.removetimeout,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "kick",
		Description: "Kick a user",
		Usage:       "[p]kick <user> [reason]",
		Flags:       []string{"DANGEROUS"},
		Handler:     m.kick,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "ban",
		Description: "Ban a user",
		Usage:       "[p]ban <user> [reason]",
		Flags:       []string{"DANGEROUS"},
		Handler:     m.ban,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "hackban",
		Description: "Ban a user that's not in the server",
		Usage:       "[p]hackban <user id> [reason]",
		Flags:       []string{"DANGEROUS"},
		Handler:     m.ban,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "softban",
		Description: "Ban and instantly unban a user",
		Usage:       "[p]softban <user> [reason]",
		Flags:       []string{"DANGEROUS"},
		Handler:     m.softban,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "unban",
		Description: "Ban and instantly unban a user",
		Usage:       "[p]unban <user id>",
		Flags:       []string{"DANGEROUS"},
		Handler:     m.unban,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "banlist",
		Description: "Get a list of banned users",
		Usage:       "[p]banlist [server id]",
		Flags:       []string{"DANGEROUS"},
		Handler:     m.banlist,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "copybans",
		Description: "Copy bans from one server to another",
		Usage:       "[p]copybans <source server id> [dest server id]",
		Flags:       []string{"DANGEROUS"},
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
		Flags:       []string{"DANGEROUS"},
		Handler:     m.nickuser,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "clearnick",
		Description: "Clear a users nickname",
		Usage:       "[p]clearnick [user]",
		Flags:       []string{"DANGEROUS"},
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
		Flags:       []string{"DANGEROUS"},
		Handler:     m.setservername,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "setservericon",
		Description: "Set a servers icon",
		Usage:       "[p]setservericon <url>",
		Flags:       []string{"DANGEROUS"},
		Handler:     m.setservericon,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "setserverbanner",
		Description: "Set a servers banner",
		Usage:       "[p]setserverbanner <url>",
		Flags:       []string{"DANGEROUS"},
		Handler:     m.setserverbanner,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "lockserver",
		Description: "Lock a server",
		Usage:       "[p]lockserver [server id]",
		Flags:       []string{"DANGEROUS"},
		Handler:     m.lockserver,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "unlockserver",
		Description: "Unlock a server",
		Usage:       "[p]unlockserver [server id]",
		Flags:       []string{"DANGEROUS"},
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
		Name:        "leaveserver",
		Description: "Leave a server",
		Usage:       "[p]leaveserver [server id]",
		Flags:       []string{"DANGEROUS"},
		Handler:     m.leaveserver,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "deleteserver",
		Description: "Delete a server",
		Usage:       "[p]deleteserver [server id]",
		Flags:       []string{"DANGEROUS"},
		Handler:     m.deleteserver,
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
		Flags:       []string{"DANGEROUS"},
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
		Flags:       []string{"DANGEROUS"},
		Handler:     m.giverole,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "removerole",
		Description: "Remove a role from a user",
		Usage:       "[p]removerole <user> <role>",
		Flags:       []string{"DANGEROUS"},
		Handler:     m.removerole,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "createrole",
		Description: "Create a role",
		Usage:       "[p]createrole <name>",
		Flags:       []string{"DANGEROUS"},
		Handler:     m.createrole,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "deleterole",
		Description: "Delete a role",
		Usage:       "[p]deleterole <role>",
		Flags:       []string{"DANGEROUS"},
		Handler:     m.deleterole,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "renamerole",
		Description: "Rename a role",
		Usage:       "[p]renamerole <role> <name>",
		Flags:       []string{"DANGEROUS"},
		Handler:     m.renamerole,
	})
}

// Channel Control

func (*AdministrationCommands) channelinfo(ctx *dgc.Ctx) {
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
		AddField("Is NSFW", utils.Misc.FormatBoolYesNo(channel.NSFW)).
		AddField("Created At", utils.Discord.SnowflakeToTime(channel.ID).Format(time.RFC1123)).
		AddField("Slowmode Delay", fmt.Sprint(channel.RateLimitPerUser)+" seconds").
		AddField("Position", fmt.Sprint(channel.Position)).
		AddField("Last Message ID", channel.LastMessageID).
		Send(ctx)
}

func (*AdministrationCommands) textchannel(ctx *dgc.Ctx) {
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

func (*AdministrationCommands) voicechannel(ctx *dgc.Ctx) {
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

func (*AdministrationCommands) stagechannel(ctx *dgc.Ctx) {
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

func (*AdministrationCommands) newschannel(ctx *dgc.Ctx) {
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

func (*AdministrationCommands) thread(ctx *dgc.Ctx) {
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

func (*AdministrationCommands) category(ctx *dgc.Ctx) {
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

func (*AdministrationCommands) deletecategory(ctx *dgc.Ctx) {
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

func (*AdministrationCommands) renamechannel(ctx *dgc.Ctx) {
	channel, success := utils.Discord.TryGetChannelFromMention(ctx, ctx.Arguments.Get(0))
	if !success {
		_ = ctx.RespondText("Could not find channel")
		return
	}

	_, err := ctx.Session.ChannelEdit(channel.ID, &discordgo.ChannelEdit{
		Name: utils.String.RemoveWords(ctx.Arguments.Raw(), 1),
	})
	if err != nil {
		_ = ctx.RespondText("Error renaming channel: " + err.Error())
		return
	}

	embed.NewEmbed().
		SetDescription("Successfully renamed channel!").
		AddField("Channel", channel.Mention()).
		Send(ctx)
}

func (*AdministrationCommands) deletechannel(ctx *dgc.Ctx) {
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

func (*AdministrationCommands) slowmode(ctx *dgc.Ctx) {
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

func (*AdministrationCommands) removeslowmode(ctx *dgc.Ctx) {
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

func (*AdministrationCommands) lockchannel(ctx *dgc.Ctx) {
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

func (*AdministrationCommands) unlockchannel(ctx *dgc.Ctx) {
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

func (*AdministrationCommands) purge(ctx *dgc.Ctx) {
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

	messages, err := utils.Discord.GetChannelMessages(ctx.Session, channelID, amount)
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

func (*AdministrationCommands) purgeself(ctx *dgc.Ctx) {
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

	messages, err := utils.Discord.GetChannelMessages(ctx.Session, channelID, amount)
	if err != nil {
		_ = ctx.RespondText("Error getting messages: " + err.Error())
		return
	}

	for _, message := range messages {
		if message.Author.ID != ctx.Session.State.User.ID {
			continue
		}
		err = ctx.Session.ChannelMessageDelete(channelID, message.ID)
		if err != nil {
			_ = ctx.RespondText("Error deleting messages: " + err.Error())
			return
		}
	}
}

func (*AdministrationCommands) purgeuser(ctx *dgc.Ctx) {
	amount, err := strconv.Atoi(utils.Arguments["<amount>"])
	if err != nil {
		_ = ctx.RespondText("Argument <amount> must be a number.")
		return
	}

	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(1))
	if !success {
		_ = ctx.RespondText("Could not find user")
		return
	}

	channel, success := utils.Discord.TryGetChannelFromMention(ctx, ctx.Arguments.Get(2))
	channelID := ctx.Event.ChannelID
	if success {
		channelID = channel.ID
	}

	messages, err := utils.Discord.GetChannelMessages(ctx.Session, channelID, amount)
	if err != nil {
		_ = ctx.RespondText("Error getting messages: " + err.Error())
		return
	}

	for _, message := range messages {
		if message.Author.ID != user.ID {
			continue
		}
		err = ctx.Session.ChannelMessageDelete(channelID, message.ID)
		if err != nil {
			_ = ctx.RespondText("Error deleting messages: " + err.Error())
			return
		}
	}
}

func (*AdministrationCommands) purgerange(ctx *dgc.Ctx) {
	channel, success := utils.Discord.TryGetChannelFromMention(ctx, ctx.Arguments.Get(2))
	channelID := ctx.Event.ChannelID
	if success {
		channelID = channel.ID
	}

	start, end := utils.Arguments["<start message>"], utils.Arguments["<end message>"]

	messages, err := utils.Discord.GetChannelMessagesRange(ctx.Session, channelID, utils.Arguments["<start message>"], utils.Arguments["<end message>"])
	if err != nil {
		_ = ctx.RespondText("Error getting messages: " + err.Error())
		return
	}

	for _, message := range messages {
		if message.ID < start || message.ID > end {
			continue
		}

		err = ctx.Session.ChannelMessageDelete(channelID, message.ID)
		if err != nil {
			_ = ctx.RespondText("Error deleting messages: " + err.Error())
			return
		}
	}
}

func (*AdministrationCommands) nuke(ctx *dgc.Ctx) {
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

func (*AdministrationCommands) ackchannel(ctx *dgc.Ctx) {
	channel, success := utils.Discord.TryGetChannelFromMention(ctx, ctx.Arguments.Get(0))
	if !success {
		channel, _ = ctx.Session.State.Channel(ctx.Event.ChannelID)
		if channel == nil {
			_ = ctx.RespondText("Could not find channel")
			return
		}
	}
	messages, _ := utils.Discord.GetChannelMessages(ctx.Session, channel.ID, 50)
	if len(messages) > 0 {
		_ = ctx.Session.MessageAck(channel.ID, messages[0].ID)
	}
}

func (*AdministrationCommands) ackmessage(ctx *dgc.Ctx) {
	channel, success := utils.Discord.TryGetChannelFromMention(ctx, ctx.Arguments.Get(0))
	if !success {
		_ = ctx.RespondText("Could not find channel")
		return
	}
	_ = ctx.Session.MessageAck(channel.ID, ctx.Arguments.Get(1).Raw())
}

// Webhook Control

func (*AdministrationCommands) webhookinfo(ctx *dgc.Ctx) {
	webhookUrl, err := url.Parse(utils.Arguments["<webhook url>"])
	if err != nil {
		_ = ctx.RespondText("Invalid webhook URL!")
		return
	}

	if !strings.Contains(webhookUrl.Host, "discord.com") {
		_ = ctx.RespondText("Invalid webhook URL!")
		return
	}

	httpResponse := http.Get(webhookUrl.String())
	if httpResponse.Error != nil || httpResponse.StatusCode != 200 {
		_ = ctx.RespondText("Invalid webhook URL!")
		return
	}

	webhookInfo := httpResponse.Json()

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

func (*AdministrationCommands) webhooksend(ctx *dgc.Ctx) {
	webhookUrl, err := url.Parse(utils.Arguments["<webhook url>"])
	if err != nil {
		_ = ctx.RespondText("Invalid webhook URL!")
		return
	}

	if !strings.Contains(webhookUrl.Host, "discord.com") {
		_ = ctx.RespondText("Invalid webhook URL!")
		return
	}

	data := map[string]any{
		"content": utils.Arguments["<\"message\">"],
	}

	_ = http.PostJson(webhookUrl.String(), data)
}

func (*AdministrationCommands) webhookdelete(ctx *dgc.Ctx) {
	webhookUrl, err := url.Parse(utils.Arguments["<webhook url>"])
	if err != nil {
		_ = ctx.RespondText("Invalid webhook URL!")
		return
	}

	if !strings.Contains(webhookUrl.Host, "discord.com") {
		_ = ctx.RespondText("Invalid webhook URL!")
		return
	}

	_ = http.CustomMethodWithData(webhookUrl.String(), "DELETE", nil)

	embed.NewEmbed().
		SetDescription("Successfully deleted webhook.").
		Send(ctx)
}

// Member Control

func (*AdministrationCommands) userinfo(ctx *dgc.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Event.Author
	}

	var member *discordgo.Member
	if ctx.Event.GuildID != "" {
		member, _ = utils.Discord.GetMember(ctx.Session, ctx.Event.GuildID, user.ID, true)
	}

	emb := embed.NewEmbed().
		SetDescription(fmt.Sprintf("User Info (%s - %s)", user.String(), user.ID)).
		AddField("Avatar", user.AvatarURL(""))

	if user.BannerURL("") != "" {
		emb.AddField("Banner", user.BannerURL(""))
	}

	if user.AccentColor != 0 {
		emb.AddField("Banner Color", fmt.Sprintf("#%x", user.AccentColor))
	}

	emb.AddField("Is Bot", utils.Misc.FormatBoolYesNo(user.Bot)).
		AddField("Created At", utils.Discord.SnowflakeToTime(user.ID).Format(time.RFC1123)).
		AddField("Premium Type", utils.Discord.FormatPremiumType(user.PremiumType))

	if flags := utils.Discord.FormatUserFlags(user.PublicFlags); flags != "" {
		emb.AddField("Flags", utils.Discord.FormatUserFlags(user.PublicFlags))
	}

	if member != nil {
		if member.Nick != "" {
			emb.AddField("Nickname", member.Nick)
		}

		emb.AddField("Joined Server On", member.JoinedAt.Format(time.RFC1123)).
			AddField("Roles", fmt.Sprintf("%d Roles:\n%s", len(member.Roles), utils.Discord.GetRoleNames(ctx, member)))

		if p := member.Presence; p != nil {
			emb.AddField("Status", utils.Discord.FormatStatusEmoji(p.Status)).
				AddField("Platform", utils.Discord.FormatClientStatus(p.ClientStatus))

			if len(p.Activities) > 0 {
				emb.AddField("Activity", utils.Discord.FormatActivities(p.Activities))
			}
		}
	}

	emb.Send(ctx)
}

func (*AdministrationCommands) userage(ctx *dgc.Ctx) {
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

func (*AdministrationCommands) mutualfriends(ctx *dgc.Ctx) {
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

func (*AdministrationCommands) mutualservers(ctx *dgc.Ctx) {
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

func (*AdministrationCommands) platform(ctx *dgc.Ctx) {
	userIDArg := ctx.Arguments.Get(0)
	userID := utils.Misc.Or(userIDArg.AsUserMentionID(), userIDArg.Raw())
	if userID == "" {
		userID = ctx.Event.Author.ID
	}

	member, ok := utils.Discord.GetMember(ctx.Session, ctx.Event.GuildID, userID, true)
	if !ok {
		_ = ctx.RespondText("Could not find user!")
		return
	}

	if member.Presence == nil {
		_ = ctx.RespondText("Users presence is unknown!")
		return
	}

	embed.NewEmbed().
		AddField("User", member.User.Username).
		AddField("Status", utils.Discord.FormatStatusEmoji(member.Presence.Status)).
		AddField("Platform", utils.Discord.FormatClientStatus(member.Presence.ClientStatus)).
		Send(ctx)
}

func (*AdministrationCommands) avatar(ctx *dgc.Ctx) {
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

func (*AdministrationCommands) banner(ctx *dgc.Ctx) {
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

func (*AdministrationCommands) report(ctx *dgc.Ctx) {
	link := utils.Arguments["<message link>"]
	reason := utils.Arguments["<reason>"]

	// Verify message link
	messageUrl, err := url.Parse(link)
	if err != nil {
		_ = ctx.RespondText("Invalid message link.")
		return
	}
	urlPathParts := strings.Split(messageUrl.Path, "/")
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

	_ = api.ReportMessage(ctx.Session, guildID, channelID, messageID, reasonInt)

	embed.NewEmbed().
		SetDescription("Successfully reported message!").
		Send(ctx)
}

func (*AdministrationCommands) mute(ctx *dgc.Ctx) {
	if !utils.Discord.CheckIfInServer(ctx) {
		return
	}
	member, success := utils.Discord.TryGetMemberFromMention(ctx, ctx.Arguments.Get(0))
	if !success {
		_ = ctx.RespondText("Could not find user")
		return
	}

	if err := utils.Discord.MuteUser(ctx.Session, member); err != nil {
		_ = ctx.RespondText(utils.String.Capitalize(err.Error()))
		return
	}

	embed.NewEmbed().
		SetDescription("Successfully muted user!").
		Send(ctx)
}

func (*AdministrationCommands) unmute(ctx *dgc.Ctx) {
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

		_, err = ctx.Session.GuildMemberEdit(ctx.Event.GuildID, member.User.ID, &discordgo.GuildMemberParams{
			Roles: &member.Roles,
		})
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

func (*AdministrationCommands) timeout(ctx *dgc.Ctx) {
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

func (*AdministrationCommands) removetimeout(ctx *dgc.Ctx) {
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

func (*AdministrationCommands) kick(ctx *dgc.Ctx) {
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

func (*AdministrationCommands) ban(ctx *dgc.Ctx) {
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

func (*AdministrationCommands) softban(ctx *dgc.Ctx) {
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

func (*AdministrationCommands) unban(ctx *dgc.Ctx) {
	if !utils.Discord.CheckIfInServer(ctx) {
		return
	}

	err := ctx.Session.GuildBanDelete(ctx.Event.GuildID, utils.Arguments["<user id>"])
	if err != nil {
		_ = ctx.RespondText("Error unbanning user: " + err.Error())
		return
	}

	embed.NewEmbed().
		SetDescription("Successfully unbanned user!").
		Send(ctx)
}

func (*AdministrationCommands) banlist(ctx *dgc.Ctx) {
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

func (*AdministrationCommands) copybans(ctx *dgc.Ctx) {
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

func (*AdministrationCommands) nick(ctx *dgc.Ctx) {
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

func (*AdministrationCommands) nickuser(ctx *dgc.Ctx) {
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

func (*AdministrationCommands) clearnick(ctx *dgc.Ctx) {
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

func (*AdministrationCommands) serverinfo(ctx *dgc.Ctx) {
	guild, err := utils.Discord.GetGuild(ctx.Session, utils.Misc.Or(utils.Arguments["[server id]"], ctx.Event.GuildID))
	if err != nil {
		_ = ctx.RespondText("Error getting server: " + err.Error())
		return
	}

	owner, ok := utils.Discord.GetMember(ctx.Session, guild.ID, guild.OwnerID)
	ownerMention := "Error"
	if ok {
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

func (*AdministrationCommands) serverage(ctx *dgc.Ctx) {
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

func (*AdministrationCommands) serverowner(ctx *dgc.Ctx) {
	guild, err := utils.Discord.GetGuild(ctx.Session, utils.Misc.Or(utils.Arguments["[server id]"], ctx.Event.GuildID))
	if err != nil {
		_ = ctx.RespondText("Error getting server: " + err.Error())
		return
	}
	owner, ok := utils.Discord.GetMember(ctx.Session, guild.ID, guild.OwnerID)
	if !ok {
		_ = ctx.RespondText("Error getting server owner")
		return
	}
	embed.NewEmbed().
		SetDescription(guild.Name).
		AddField("Owner", owner.Mention()).
		Send(ctx)
}

func (*AdministrationCommands) servericon(ctx *dgc.Ctx) {
	guild, err := utils.Discord.GetGuild(ctx.Session, utils.Misc.Or(utils.Arguments["[server id]"], ctx.Event.GuildID))
	if err != nil {
		_ = ctx.RespondText("Error getting server: " + err.Error())
		return
	}

	emb := embed.NewEmbed().
		SetDescription(fmt.Sprintf("Server Icon - %s", guild.Name)).
		SetImage(guild.IconURL())
	if guild.Icon == "" {
		emb.AddField("Error", "Server has no icon")
	}
	emb.Send(ctx)
}

func (*AdministrationCommands) serverbanner(ctx *dgc.Ctx) {
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

func (*AdministrationCommands) serversplash(ctx *dgc.Ctx) {
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

func (*AdministrationCommands) serverdsplash(ctx *dgc.Ctx) {
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

func (*AdministrationCommands) serverboosters(ctx *dgc.Ctx) {
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

func (*AdministrationCommands) serverboostlevel(ctx *dgc.Ctx) {
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

func (*AdministrationCommands) setservername(ctx *dgc.Ctx) {
	if !utils.Discord.CheckIfInServer(ctx) {
		return
	}
	guild, err := ctx.Session.GuildEdit(ctx.Event.GuildID, &discordgo.GuildParams{
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

func (*AdministrationCommands) setservericon(ctx *dgc.Ctx) {
	if !utils.Discord.CheckIfInServer(ctx) {
		return
	}
	_, err := ctx.Session.GuildEdit(ctx.Event.GuildID, &discordgo.GuildParams{
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

func (*AdministrationCommands) setserverbanner(ctx *dgc.Ctx) {
	if !utils.Discord.CheckIfInServer(ctx) {
		return
	}
	_, err := ctx.Session.GuildEdit(ctx.Event.GuildID, &discordgo.GuildParams{
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

func (*AdministrationCommands) lockserver(ctx *dgc.Ctx) {
	guild, err := utils.Discord.GetGuild(ctx.Session, utils.Misc.Or(utils.Arguments["[server id]"], ctx.Event.GuildID))
	if err != nil {
		_ = ctx.RespondText("Error getting server: " + err.Error())
		return
	}

	overwrites := []*discordgo.PermissionOverwrite{
		{
			ID:   guild.ID,
			Type: discordgo.PermissionOverwriteTypeRole,
			Deny: discordgo.PermissionSendMessages,
		},
	}

	for _, channel := range guild.Channels {
		if channel.Type != discordgo.ChannelTypeGuildText {
			continue
		}

		_, err = ctx.Session.ChannelEditComplex(channel.ID, &discordgo.ChannelEdit{
			PermissionOverwrites: overwrites,
		})

		if err != nil {
			_ = ctx.RespondText("Error locking server: " + err.Error())
			return
		}
	}

	embed.NewEmbed().
		SetDescription("Successfully locked server!").
		Send(ctx)
}

func (*AdministrationCommands) unlockserver(ctx *dgc.Ctx) {
	guild, err := utils.Discord.GetGuild(ctx.Session, utils.Misc.Or(utils.Arguments["[server id]"], ctx.Event.GuildID))
	if err != nil {
		_ = ctx.RespondText("Error getting server: " + err.Error())
		return
	}

	overwrites := []*discordgo.PermissionOverwrite{
		{
			ID:   guild.ID,
			Type: discordgo.PermissionOverwriteTypeRole,
			Deny: 0,
		},
	}

	for _, channel := range guild.Channels {
		if channel.Type != discordgo.ChannelTypeGuildText {
			continue
		}

		_, err = ctx.Session.ChannelEditComplex(channel.ID, &discordgo.ChannelEdit{
			PermissionOverwrites: overwrites,
		})

		if err != nil {
			_ = ctx.RespondText("Error locking server: " + err.Error())
			return
		}
	}

	embed.NewEmbed().
		SetDescription("Successfully locked server!").
		Send(ctx)
}

func (*AdministrationCommands) channels(ctx *dgc.Ctx) {
	guild, err := utils.Discord.GetGuild(ctx.Session, utils.Misc.Or(utils.Arguments["[server id]"], ctx.Event.GuildID))
	if err != nil {
		_ = ctx.RespondText("Error getting server: " + err.Error())
		return
	}

	categories, categoryChannels, looseChannels := utils.Discord.MapCategories(guild.Channels)
	emb := embed.NewEmbed().SetDescription(fmt.Sprintf("Server Channels - %s", guild.Name))

	looseString := ""
	for _, channel := range looseChannels {
		looseString += channel.Mention() + "\n"
	}
	emb.AddField("Channels (No Category)", looseString)

	for _, category := range categories {
		channelsString := ""
		for _, channel := range categoryChannels[category.ID] {
			channelsString += channel.Mention() + "\n"
		}
		emb.AddField("Category: "+category.Name, channelsString)
	}
	emb.Send(ctx)
}

func (*AdministrationCommands) bots(ctx *dgc.Ctx) {
	guild, err := utils.Discord.GetGuild(ctx.Session, utils.Misc.Or(utils.Arguments["[server id]"], ctx.Event.GuildID))
	if err != nil {
		_ = ctx.RespondText("Error getting server: " + err.Error())
		return
	}
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

func (*AdministrationCommands) admins(ctx *dgc.Ctx) {
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

func (*AdministrationCommands) roles(ctx *dgc.Ctx) {
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

func (*AdministrationCommands) leaveserver(ctx *dgc.Ctx) {
	guild, err := utils.Discord.GetGuild(ctx.Session, utils.Misc.Or(utils.Arguments["[server id]"], ctx.Event.GuildID))
	if err != nil {
		_ = ctx.RespondText("Error getting server: " + err.Error())
		return
	}
	_ = ctx.Session.GuildLeave(guild.ID)
}

func (*AdministrationCommands) deleteserver(ctx *dgc.Ctx) {
	guild, err := utils.Discord.GetGuild(ctx.Session, utils.Misc.Or(utils.Arguments["[server id]"], ctx.Event.GuildID))
	if err != nil {
		_ = ctx.RespondText("Error getting server: " + err.Error())
		return
	}

	_, err = ctx.Session.GuildDelete(guild.ID)
	if err != nil {
		_ = ctx.RespondText("Error deleting server: " + err.Error())
		return
	}
}

// Invite Control

func (*AdministrationCommands) invite(ctx *dgc.Ctx) {
	channel, success := utils.Discord.TryGetChannelFromMention(ctx, ctx.Arguments.Get(0))
	channelID := ctx.Event.ChannelID
	if success {
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

func (*AdministrationCommands) friendinvite(ctx *dgc.Ctx) {
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

func (*AdministrationCommands) inviteinfo(ctx *dgc.Ctx) {
	parsedInvite := utils.Discord.ParseInvite(utils.Arguments["<invite>"])
	invite, err := ctx.Session.InviteWithCounts(parsedInvite)
	if err != nil {
		_ = ctx.RespondText("Error getting invite: " + err.Error())
		return
	}

	emb := embed.NewEmbed().SetDescription(fmt.Sprintf("Invite Info - %s", invite.Code))
	if invite.Guild != nil {
		emb.AddField("Server", invite.Guild.Name)
		if invite.Guild.Icon != "" {
			emb.SetImage(invite.Guild.IconURL())
		}
	}
	if invite.Channel != nil {
		emb.AddField("Channel", invite.Channel.Name)
	}
	if invite.Inviter != nil {
		emb.AddField("Inviter", invite.Inviter.String())
	}
	if invite.ExpiresAt != nil {
		emb.AddField("Expires At", invite.ExpiresAt.Format(time.RFC1123))
	}
	emb.AddField("Max Age", fmt.Sprint(invite.MaxAge)).
		AddField("Uses", fmt.Sprint(invite.Uses)).
		AddField("Max Uses", fmt.Sprint(invite.MaxUses)).
		AddField("Member Count", fmt.Sprint(invite.ApproximateMemberCount)).
		AddField("Online Member Count", fmt.Sprint(invite.ApproximatePresenceCount)).
		Send(ctx)
}

func (*AdministrationCommands) deleteinvite(ctx *dgc.Ctx) {
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

func (*AdministrationCommands) roleinfo(ctx *dgc.Ctx) {
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
		AddField("Is Hoist", utils.Misc.FormatBoolYesNo(role.Hoist)).
		AddField("Is Mentionable", utils.Misc.FormatBoolYesNo(role.Mentionable)).
		AddField("Members With Role", fmt.Sprint(memberCount)).
		AddField("Permissions", perms).
		Send(ctx)
}

func (*AdministrationCommands) giverole(ctx *dgc.Ctx) {
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

func (*AdministrationCommands) removerole(ctx *dgc.Ctx) {
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

func (*AdministrationCommands) createrole(ctx *dgc.Ctx) {
	if !utils.Discord.CheckIfInServer(ctx) {
		return
	}

	role, err := ctx.Session.GuildRoleCreate(ctx.Event.GuildID, &discordgo.RoleParams{
		Name: ctx.Arguments.Raw(),
	})
	if err != nil {
		_ = ctx.RespondText("Error creating role: " + err.Error())
		return
	}

	embed.NewEmbed().
		SetDescription(fmt.Sprintf("Successfully created role (%s)", role.Name)).
		Send(ctx)
}

func (*AdministrationCommands) deleterole(ctx *dgc.Ctx) {
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

func (*AdministrationCommands) renamerole(ctx *dgc.Ctx) {
	if !utils.Discord.CheckIfInServer(ctx) {
		return
	}

	role, success := utils.Discord.TryGetRoleFromMention(ctx, ctx.Arguments.Get(0))
	if !success {
		_ = ctx.RespondText("Could not find role")
		return
	}

	role, err := ctx.Session.GuildRoleEdit(ctx.Event.GuildID, role.ID, &discordgo.RoleParams{
		Name: ctx.Arguments.Raw(),
	})
	if err != nil {
		_ = ctx.RespondText("Error renaming role: " + err.Error())
		return
	}

	embed.NewEmbed().
		SetDescription(fmt.Sprintf("Successfully renamed role (%s)", role.Name)).
		Send(ctx)
}
