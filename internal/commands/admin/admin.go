package admin

import (
	"fmt"
	"net/url"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/lb-selfbot/discordgo"
	"github.com/lb-selfbot/selfbot/internal/api"
	"github.com/lb-selfbot/selfbot/internal/http"
	"github.com/lb-selfbot/selfbot/internal/services/router"

	"github.com/lb-selfbot/selfbot/internal/fn"
	"github.com/lb-selfbot/selfbot/internal/services"
	"github.com/lb-selfbot/selfbot/internal/utils"
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

type AdminCommands struct{ services *services.Services }

func NewAdminCommands(s *services.Services) *AdminCommands {
	return &AdminCommands{services: s}
}

func (c *AdminCommands) Register() {
	c.services.Router.RegisterCategory(&router.Category{
		Name:        "Admin",
		Description: "Administration commands",
	}, []*router.Command{
		{
			Name:        "channelinfo",
			Description: "Get info about a channel",
			Arguments:   "<#channel>",
			Handler:     c.channelinfo,
		},
		{
			Name:        "textchannel",
			Description: "Create a text channel",
			Arguments:   "<name>",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous, router.CmdFlagGuildOnly},
			Handler:     c.textchannel,
		},
		{
			Name:        "voicechannel",
			Description: "Create a voice channel",
			Arguments:   "<name>",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous, router.CmdFlagGuildOnly},
			Handler:     c.voicechannel,
		},
		{
			Name:        "stagechannel",
			Description: "Create a stage channel",
			Arguments:   "<name>",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous, router.CmdFlagGuildOnly},
			Handler:     c.stagechannel,
		},
		{
			Name:        "newschannel",
			Description: "Create a news channel",
			Arguments:   "<name>",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous, router.CmdFlagGuildOnly},
			Handler:     c.newschannel,
		},
		{
			Name:        "thread",
			Description: "Create a thread",
			Arguments:   "<name>",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous, router.CmdFlagGuildOnly},
			Handler:     c.thread,
		},
		{
			Name:        "category",
			Description: "Create a category",
			Arguments:   "<name>",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous, router.CmdFlagGuildOnly},
			Handler:     c.category,
		},
		{
			Name:        "deletecategory",
			Description: "Delete a category",
			Arguments:   "<category id>",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.deletecategory,
		},
		{
			Name:        "renamechannel",
			Description: "Rename a channel",
			Arguments:   "<#channel> <name>",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.renamechannel,
		},
		{
			Name:        "deletechannel",
			Description: "Delete a channel",
			Arguments:   "<#channel>",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.deletechannel,
		},
		{
			Name:        "slowmode",
			Description: "Change slowmode amount",
			Arguments:   "<#channel> <duration>",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.slowmode,
		},
		{
			Name:        "removeslowmode",
			Description: "Remove slowmode",
			Arguments:   "<#channel>",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.removeslowmode,
		},
		{
			Name:        "lockchannel",
			Description: "Lock a channel",
			Arguments:   "<#channel>",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.lockchannel,
		},
		{
			Name:        "unlockchannel",
			Description: "Unlock a channel",
			Arguments:   "<#channel>",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.unlockchannel,
		},
		{
			Name:        "purge",
			Description: "Delete messages in a channel",
			Arguments:   "<amount> [delay] [channel]",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.purge,
		},
		{
			Name:        "purgeself",
			Description: "Delete your own messages in a channel",
			Arguments:   "<amount> [delay] [channel]",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.purgeself,
		},
		{
			Name:        "purgeuser",
			Description: "Delete a users messages in a channel",
			Arguments:   "<amount> <user> [delay] [channel]",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.purgeuser,
		},
		{
			Name:        "purgerange",
			Description: "Purge a range of messages in a channel",
			Arguments:   "<start message> <end message> [delay] [channel]",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.purgerange,
		},
		{
			Name:        "nuke",
			Description: "Nuke a channel",
			Arguments:   "[channel]",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.nuke,
		},
		{
			Name:        "ackchannel",
			Description: "Ack a channel",
			Arguments:   "[channel]",
			Handler:     c.ackchannel,
		},
		{
			Name:        "ackmessage",
			Description: "Ack a message",
			Arguments:   "<#channel> <message id>",
			Handler:     c.ackmessage,
		},
		{
			Name:        "listwebhooks",
			Description: "List all webhooks in a channel",
			Arguments:   "[channel]",
			Handler:     c.listwebhooks,
		},
		{
			Name:        "createwebhook",
			Description: "List all webhooks in a channel",
			Arguments:   "[channel] [name]",
			Handler:     c.createwebhook,
		},
		{
			Name:        "webhookinfo",
			Description: "Get info about a webhook",
			Arguments:   "<webhook url>",
			Handler:     c.webhookinfo,
		},
		{
			Name:        "webhooksend",
			Description: "Send a webhook message",
			Arguments:   "<webhook url> <\"message\">",
			Handler:     c.webhooksend,
		},
		{
			Name:        "webhookdelete",
			Description: "Delete a webhook",
			Arguments:   "<webhook url>",
			Handler:     c.webhookdelete,
		},
		{
			Name:        "userinfo",
			Description: "Get info about a user",
			Arguments:   "[user]",
			Handler:     c.userinfo,
		},
		{
			Name:        "userage",
			Description: "See the age of a users account",
			Arguments:   "<user>",
			Handler:     c.userage,
		},
		{
			Name:        "mutualfriends",
			Description: "Get mutual friends with user",
			Arguments:   "<user>",
			Handler:     c.mutualfriends,
		},
		{
			Name:        "mutualservers",
			Description: "Get mutual servers with user",
			Arguments:   "<user>",
			Handler:     c.mutualservers,
		},
		{
			Name:        "platform",
			Description: "See the platforms a user is online on",
			Arguments:   "[user]",
			Handler:     c.platform,
		},
		{
			Name:        "avatar",
			Description: "Get a users avatar",
			Arguments:   "<user>",
			Handler:     c.avatar,
		},
		{
			Name:        "banner",
			Description: "Get a users banner",
			Arguments:   "<user>",
			Handler:     c.banner,
		},
		{
			Name:        "mute",
			Description: "Mute a user",
			Arguments:   "<user>",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous, router.CmdFlagGuildOnly},
			Handler:     c.mute,
		},
		{
			Name:        "unmute",
			Description: "Unmute a user",
			Arguments:   "<user>",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous, router.CmdFlagGuildOnly},
			Handler:     c.unmute,
		},
		{
			Name:        "timeout",
			Description: "Timeout a user",
			Arguments:   "<user> <duration>",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous, router.CmdFlagGuildOnly},
			Handler:     c.timeout,
		},
		{
			Name:        "removetimeout",
			Description: "Remove a timeout",
			Arguments:   "<user>",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous, router.CmdFlagGuildOnly},
			Handler:     c.removetimeout,
		},
		{
			Name:        "kick",
			Description: "Kick a user",
			Arguments:   "<user> [reason]",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous, router.CmdFlagGuildOnly},
			Handler:     c.kick,
		},
		{
			Name:        "ban",
			Description: "Ban a user",
			Arguments:   "<user> [reason]",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous, router.CmdFlagGuildOnly},
			Handler:     c.ban,
		},
		{
			Name:        "softban",
			Description: "Ban and instantly unban a user",
			Arguments:   "<user> [reason]",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous, router.CmdFlagGuildOnly},
			Handler:     c.softban,
		},
		{
			Name:        "unban",
			Description: "Unban a user",
			Arguments:   "<user id>",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous, router.CmdFlagGuildOnly},
			Handler:     c.unban,
		},
		{
			Name:        "banlist",
			Description: "Get a list of banned users",
			Arguments:   "[server id]",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.banlist,
		},
		{
			Name:        "copybans",
			Description: "Copy bans from one server to another",
			Arguments:   "<source server id> [dest server id]",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.copybans,
		},
		{
			Name:        "nick",
			Description: "Change your nickname",
			Arguments:   "<text>",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous, router.CmdFlagGuildOnly},
			Handler:     c.nick,
		},
		{
			Name:        "nickuser",
			Description: "Change a users nickname",
			Arguments:   "<user> <text>",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous, router.CmdFlagGuildOnly},
			Handler:     c.nickuser,
		},
		{
			Name:        "clearnick",
			Description: "Clear a users nickname",
			Arguments:   "[user]",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous, router.CmdFlagGuildOnly},
			Handler:     c.clearnick,
		},
		{
			Name:        "serverinfo",
			Description: "Get info about a server",
			Arguments:   "[server id]",
			Handler:     c.serverinfo,
		},
		{
			Name:        "serverage",
			Description: "Get the age of a server",
			Arguments:   "[server id]",
			Handler:     c.serverage,
		},
		{
			Name:        "serverowner",
			Description: "Get a servers owner",
			Arguments:   "[server id]",
			Handler:     c.serverowner,
		},
		{
			Name:        "servericon",
			Description: "Get a servers icon",
			Arguments:   "[server id]",
			Handler:     c.servericon,
		},
		{
			Name:        "serverbanner",
			Description: "Get a servers banner",
			Arguments:   "[server id]",
			Handler:     c.serverbanner,
		},
		{
			Name:        "serversplash",
			Description: "Get a servers invite splash",
			Arguments:   "[server id]",
			Handler:     c.serversplash,
		},
		{
			Name:        "serverdsplash",
			Description: "Get a servers discovery splash",
			Arguments:   "[server id]",
			Handler:     c.serverdsplash,
		},
		{
			Name:        "serverboosters",
			Description: "Get a list of server boosters",
			Arguments:   "[server id]",
			Handler:     c.serverboosters,
		},
		{
			Name:        "serverboostlevel",
			Description: "Get a servers boost level",
			Arguments:   "[server id]",
			Handler:     c.serverboostlevel,
		},
		{
			Name:        "setservername",
			Description: "Set a servers name",
			Arguments:   "<name>",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous, router.CmdFlagGuildOnly},
			Handler:     c.setservername,
		},
		{
			Name:        "setservericon",
			Description: "Set a servers icon",
			Arguments:   "<url>",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous, router.CmdFlagGuildOnly},
			Handler:     c.setservericon,
		},
		{
			Name:        "setserverbanner",
			Description: "Set a servers banner",
			Arguments:   "<url>",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous, router.CmdFlagGuildOnly},
			Handler:     c.setserverbanner,
		},
		{
			Name:        "lockserver",
			Description: "Lock a server",
			Arguments:   "[server id]",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.lockserver,
		},
		{
			Name:        "unlockserver",
			Description: "Unlock a server",
			Arguments:   "[server id]",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.unlockserver,
		},
		{
			Name:        "channels",
			Description: "See all channels in a server",
			Arguments:   "[server id]",
			Handler:     c.channels,
		},
		{
			Name:        "bots",
			Description: "See all bots in a server",
			Arguments:   "[server id]",
			Handler:     c.bots,
		},
		{
			Name:        "admins",
			Description: "See all admins in a server",
			Arguments:   "[server id]",
			Handler:     c.admins,
		},
		{
			Name:        "roles",
			Description: "See all roles in a server",
			Arguments:   "[server id]",
			Handler:     c.roles,
		},
		{
			Name:        "leaveserver",
			Description: "Leave a server",
			Arguments:   "[server id]",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.leaveserver,
		},
		{
			Name:        "deleteserver",
			Description: "Delete a server",
			Arguments:   "[server id]",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.deleteserver,
		},
		{
			Name:        "invite",
			Description: "Create an invite",
			Arguments:   "[channel] [max age] [max uses]",
			Handler:     c.invite,
		},
		{
			Name:        "friendinvite",
			Description: "Create a friend invite",
			Arguments:   "friendinvite",
			Handler:     c.friendinvite,
		},
		{
			Name:        "inviteinfo",
			Description: "Get info about an invite",
			Arguments:   "<invite>",
			Handler:     c.inviteinfo,
		},
		{
			Name:        "deleteinvite",
			Description: "Delete an invite",
			Arguments:   "<invite>",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.deleteinvite,
		},
		{
			Name:        "serverinvites",
			Description: "Get all invites in a server",
			Arguments:   "[server id]",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.serverinvites,
		},
		{
			Name:        "channelinvites",
			Description: "Get all invites in a channel",
			Arguments:   "[channel]",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.channelinvites,
		},
		{
			Name:        "userinvites",
			Description: "Get all invites created by a user",
			Arguments:   "<user> [server id]",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous, router.CmdFlagGuildOnly},
			Handler:     c.userinvites,
		},
		{
			Name:        "roleinfo",
			Description: "Get info about a role",
			Arguments:   "<role>",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous, router.CmdFlagGuildOnly},
			Handler:     c.roleinfo,
		},
		{
			Name:        "giverole",
			Description: "Give a role to a user",
			Arguments:   "<user> <role>",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous, router.CmdFlagGuildOnly},
			Handler:     c.giverole,
		},
		{
			Name:        "removerole",
			Description: "Remove a role from a user",
			Arguments:   "<user> <role>",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous, router.CmdFlagGuildOnly},
			Handler:     c.removerole,
		},
		{
			Name:        "createrole",
			Description: "Create a role",
			Arguments:   "<name>",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous, router.CmdFlagGuildOnly},
			Handler:     c.createrole,
		},
		{
			Name:        "deleterole",
			Description: "Delete a role",
			Arguments:   "<role>",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous, router.CmdFlagGuildOnly},
			Handler:     c.deleterole,
		},
		{
			Name:        "renamerole",
			Description: "Rename a role",
			Arguments:   "<role> <name>",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous, router.CmdFlagGuildOnly},
			Handler:     c.renamerole,
		},
	})
}

// Channel Control

func (c *AdminCommands) channelinfo(ctx *router.Ctx) {
	channel, err := ctx.Arguments.Get(0).Channel()
	if err != nil {
		ctx.RespondText("Could not find channel")
		return
	}

	c.services.Embed.New(ctx).
		AddField("Name", channel.Name).
		AddField("ID", channel.ID).
		AddField("Topic", channel.Topic).
		AddField("Category ID", channel.ParentID).
		AddField("Is NSFW", utils.Misc.FormatBoolYesNo(channel.NSFW)).
		AddField("Created At", utils.Discord.SnowflakeToTime(channel.ID).Format(time.RFC1123)).
		AddField("Slowmode Delay", fmt.Sprint(channel.RateLimitPerUser)+" seconds").
		AddField("Position", fmt.Sprint(channel.Position)).
		AddField("Last Message ID", channel.LastMessageID).
		Send()
}

func (c *AdminCommands) textchannel(ctx *router.Ctx) {
	channel, err := ctx.Session.GuildChannelCreate(ctx.Message.GuildID, ctx.Arguments.Raw(), discordgo.ChannelTypeGuildText)
	if err != nil {
		ctx.RespondText("Error creating channel: %v", err)
		return
	}
	c.services.Embed.New(ctx).
		SetDescription("Successfully created text channel!").
		AddField("Channel", channel.Mention()).
		Send()
}

func (c *AdminCommands) voicechannel(ctx *router.Ctx) {
	channel, err := ctx.Session.GuildChannelCreate(ctx.Message.GuildID, ctx.Arguments.Raw(), discordgo.ChannelTypeGuildVoice)
	if err != nil {
		ctx.RespondText("Error creating channel: %v", err)
		return
	}
	c.services.Embed.New(ctx).
		SetDescription("Successfully created voice channel!").
		AddField("Channel", channel.Mention()).
		Send()
}

func (c *AdminCommands) stagechannel(ctx *router.Ctx) {
	channel, err := ctx.Session.GuildChannelCreate(ctx.Message.GuildID, ctx.Arguments.Raw(), discordgo.ChannelTypeGuildStageVoice)
	if err != nil {
		ctx.RespondText("Error creating channel: %v", err)
		return
	}
	c.services.Embed.New(ctx).
		SetDescription("Successfully created stage channel!").
		AddField("Channel", channel.Mention()).
		Send()
}

func (c *AdminCommands) newschannel(ctx *router.Ctx) {
	channel, err := ctx.Session.GuildChannelCreate(ctx.Message.GuildID, ctx.Arguments.Raw(), discordgo.ChannelTypeGuildNews)
	if err != nil {
		ctx.RespondText("Error creating channel: %v", err)
		return
	}
	c.services.Embed.New(ctx).
		SetDescription("Successfully created news channel!").
		AddField("Channel", channel.Mention()).
		Send()
}

func (c *AdminCommands) thread(ctx *router.Ctx) {
	channel, err := ctx.Session.ThreadStart(ctx.Message.ChannelID, ctx.Arguments.Raw(), discordgo.ChannelTypeGuildPublicThread, 0)
	if err != nil {
		ctx.RespondText("Error creating thread: %v", err)
		return
	}
	c.services.Embed.New(ctx).
		SetDescription("Successfully created thread!").
		AddField("Thread", channel.Mention()).
		Send()
}

func (c *AdminCommands) category(ctx *router.Ctx) {
	channel, err := ctx.Session.GuildChannelCreate(ctx.Message.GuildID, ctx.Arguments.Raw(), discordgo.ChannelTypeGuildCategory)
	if err != nil {
		ctx.RespondText("Error creating category: %v", err)
		return
	}
	c.services.Embed.New(ctx).
		SetDescription("Successfully created category!").
		AddField("Category", channel.Mention()).
		Send()
}

func (c *AdminCommands) deletecategory(ctx *router.Ctx) {
	channel, err := ctx.Arguments.Get(0).Channel()
	if err != nil {
		ctx.RespondText("Could not find category")
		return
	}
	_, err = ctx.Session.ChannelDelete(channel.ID)
	if err != nil {
		ctx.RespondText("Error deleting category: %v", err)
		return
	}

	c.services.Embed.New(ctx).
		SetDescription("Successfully deleted category!").
		Send()
}

func (c *AdminCommands) renamechannel(ctx *router.Ctx) {
	channel, err := ctx.Arguments.Get(0).Channel()
	if err != nil {
		ctx.RespondText("Could not find channel")
		return
	}

	_, err = ctx.Session.ChannelEdit(channel.ID, &discordgo.ChannelEdit{
		Name: ctx.Arguments.RawFrom(1),
	})
	if err != nil {
		ctx.RespondText("Error renaming channel: %v", err)
		return
	}

	c.services.Embed.New(ctx).
		SetDescription("Successfully renamed channel!").
		AddField("Channel", channel.Mention()).
		Send()
}

func (c *AdminCommands) deletechannel(ctx *router.Ctx) {
	channel, err := ctx.Arguments.Get(0).Channel()
	if err != nil {
		ctx.RespondText("Could not find channel")
		return
	}

	_, err = ctx.Session.ChannelDelete(channel.ID)
	if err != nil {
		ctx.RespondText("Error deleting channel: %v", err)
		return
	}

	c.services.Embed.New(ctx).
		SetDescription("Successfully deleted channel!").
		Send()
}

func (c *AdminCommands) slowmode(ctx *router.Ctx) {
	channel, err := ctx.Arguments.Get(0).Channel()
	if err != nil {
		ctx.RespondText("Could not find channel")
		return
	}

	duration, err := ctx.Arguments.Get(1).Duration()
	if err != nil {
		ctx.RespondText("Could not parse duration")
		return
	}

	durationSec := int(duration.Seconds())

	_, err = ctx.Session.ChannelEditComplex(channel.ID, &discordgo.ChannelEdit{
		Position:         channel.Position,
		RateLimitPerUser: &durationSec,
	})
	if err != nil {
		ctx.RespondText("Error setting slowmode: %v", err)
		return
	}

	c.services.Embed.New(ctx).
		SetDescription("Successfully set slowmode to %s", duration.String()).
		AddField("Channel", channel.Mention()).
		Send()
}

func (c *AdminCommands) removeslowmode(ctx *router.Ctx) {
	channel, err := ctx.Arguments.Get(0).Channel()
	if err != nil {
		ctx.RespondText("Could not find channel")
		return
	}

	_, err = ctx.Session.ChannelEditComplex(channel.ID, &discordgo.ChannelEdit{
		Position:         channel.Position,
		RateLimitPerUser: nil,
	})
	if err != nil {
		ctx.RespondText("Error removing slowmode: %v", err)
		return
	}

	c.services.Embed.New(ctx).
		SetDescription("Successfully removed slowmode!").
		AddField("Channel", channel.Mention()).
		Send()
}

func (c *AdminCommands) lockchannel(ctx *router.Ctx) {
	channel, err := ctx.Arguments.Get(0).Channel()
	if err != nil {
		ctx.RespondText("Could not find channel")
		return
	}

	_, err = ctx.Session.ChannelEditComplex(channel.ID, &discordgo.ChannelEdit{
		Position: channel.Position,
		PermissionOverwrites: []*discordgo.PermissionOverwrite{
			{
				ID:   channel.GuildID,
				Type: discordgo.PermissionOverwriteTypeRole,
				Deny: discordgo.PermissionSendMessages,
			},
		},
	})
	if err != nil {
		ctx.RespondText("Error locking channel: %v", err)
		return
	}

	c.services.Embed.New(ctx).
		SetDescription("Successfully locked channel!").
		AddField("Channel", channel.Mention()).
		Send()
}

func (c *AdminCommands) unlockchannel(ctx *router.Ctx) {
	channel, err := ctx.Arguments.Get(0).Channel()
	if err != nil {
		ctx.RespondText("Could not find channel")
		return
	}

	_, err = ctx.Session.ChannelEditComplex(channel.ID, &discordgo.ChannelEdit{
		Position: channel.Position,
		PermissionOverwrites: []*discordgo.PermissionOverwrite{
			{
				ID:   channel.GuildID,
				Type: discordgo.PermissionOverwriteTypeRole,
				Deny: 0,
			},
		},
	})
	if err != nil {
		ctx.RespondText("Error unlocking channel: %v", err)
		return
	}

	c.services.Embed.New(ctx).
		SetDescription("Successfully unlocked channel!").
		AddField("Channel", channel.Mention()).
		Send()
}

func (c *AdminCommands) purge(ctx *router.Ctx) {
	amount, err := ctx.Arguments.Get(0).Int()
	amount = fn.Clamp(amount, 1, 10000)
	if err != nil {
		ctx.RespondText("Argument <amount> must be a number.")
		return
	}

	delay, err := ctx.Arguments.Get(1).Duration()
	delay = fn.Clamp(delay, 1*time.Second, 60*time.Second)
	if err != nil {
		ctx.RespondText("Could not parse delay")
		return
	}

	channel, err := ctx.Arguments.Get(2).Channel(router.ArgumentFlagFallbackDefault)
	if err != nil {
		ctx.RespondText("Could not find channel")
		return
	}

	messages, err := utils.Discord.GetChannelMessages(ctx.Session, channel.ID, amount)
	if err != nil {
		ctx.RespondText("Error getting messages: %v", err)
		return
	}

	for _, message := range messages {
		err = ctx.Session.ChannelMessageDelete(channel.ID, message.ID)
		if err != nil {
			ctx.RespondText("Error deleting messages: %v", err)
			return
		}

		time.Sleep(delay)
	}
}

func (c *AdminCommands) purgeself(ctx *router.Ctx) {
	amount, err := ctx.Arguments.Get(0).Int()
	amount = fn.Clamp(amount, 1, 10000)
	if err != nil {
		ctx.RespondText("Argument <amount> must be a number.")
		return
	}

	delay, err := ctx.Arguments.Get(1).Duration()
	delay = fn.Clamp(delay, 1*time.Second, 60*time.Second)
	if err != nil {
		ctx.RespondText("Could not parse delay")
		return
	}

	channel, err := ctx.Arguments.Get(2).Channel(router.ArgumentFlagFallbackDefault)
	if err != nil {
		ctx.RespondText("Could not find channel")
		return
	}

	userID := ctx.Message.Author.ID

	// Allow the command to be deleted before the messages are deleted
	time.Sleep(1 * time.Second)

	deleted, err := utils.Discord.DeleteMessages(ctx.Session, channel.ID, userID, amount, delay)
	if err != nil {
		ctx.RespondText("Error deleting messages: %v", err)
		return
	}

	c.services.Embed.New(ctx).
		SetDescription("Successfully deleted %d messages!", deleted).
		Send()
}

func (c *AdminCommands) purgeuser(ctx *router.Ctx) {
	amount, err := ctx.Arguments.Get(0).Int()
	amount = fn.Clamp(amount, 1, 10000)
	if err != nil {
		ctx.RespondText("Argument <amount> must be a number.")
		return
	}

	user, err := ctx.Arguments.Get(1).User()
	if err != nil {
		ctx.RespondText("Could not find user")
		return
	}

	delay, _ := ctx.Arguments.Get(2).Duration()
	delay = fn.Clamp(delay, 1*time.Second, 60*time.Second)

	channel, err := ctx.Arguments.Get(3).Channel(router.ArgumentFlagFallbackDefault)
	if err != nil {
		ctx.RespondText("Could not find channel")
		return
	}

	// Allow the command to be deleted before the messages are deleted
	time.Sleep(1 * time.Second)

	deleted, err := utils.Discord.DeleteMessages(ctx.Session, channel.ID, user.ID, amount, delay)
	if err != nil {
		ctx.RespondText("Error deleting messages: %v", err)
		return
	}

	c.services.Embed.New(ctx).
		SetDescription("Successfully deleted %d messages!", deleted).
		Send()
}

func (c *AdminCommands) purgerange(ctx *router.Ctx) {
	start, end := ctx.Arguments.Get(0).Raw(), ctx.Arguments.Get(1).Raw()

	delay, _ := ctx.Arguments.Get(2).Duration()
	delay = fn.Clamp(delay, 1*time.Second, 60*time.Second)

	channel, err := ctx.Arguments.Get(3).Channel(router.ArgumentFlagFallbackDefault)
	if err != nil {
		ctx.RespondText("Could not find channel")
		return
	}

	messages, err := utils.Discord.GetChannelMessagesRange(ctx.Session, channel.ID, start, end)
	if err != nil {
		ctx.RespondText("Error getting messages: %v", err)
		return
	}

	for _, message := range messages {
		if message.ID < start || message.ID > end {
			continue
		}

		err = ctx.Session.ChannelMessageDelete(channel.ID, message.ID)
		if err != nil {
			ctx.RespondText("Error deleting messages: %v", err)
			return
		}

		time.Sleep(delay)
	}
}

func (c *AdminCommands) nuke(ctx *router.Ctx) {
	channel, err := ctx.Arguments.Get(0).Channel(router.ArgumentFlagFallbackDefault)
	if err != nil {
		ctx.RespondText("Could not find channel")
		return
	}

	_, err = ctx.Session.ChannelDelete(channel.ID)
	if err != nil {
		ctx.RespondText("Error deleting channel: %v", err)
		return
	}

	channel, err = ctx.Session.GuildChannelCreateComplex(ctx.Message.GuildID, discordgo.GuildChannelCreateData{
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
		ctx.RespondText("Error creating channel: %v", err)
		return
	}
	_, _ = ctx.Session.ChannelEditComplex(channel.ID, &discordgo.ChannelEdit{
		Position: channel.Position,
	})

	ctx.Message.ChannelID = channel.ID

	c.services.Embed.New(ctx).
		SetDescription("Successfully nuked channel!").
		Send()
}

func (c *AdminCommands) ackchannel(ctx *router.Ctx) {
	channel, err := ctx.Arguments.Get(0).Channel(router.ArgumentFlagFallbackDefault)
	if err != nil {
		ctx.RespondText("Could not find channel")
		return
	}
	messages, _ := utils.Discord.GetChannelMessages(ctx.Session, channel.ID, 50)
	if len(messages) > 0 {
		_ = ctx.Session.MessageAck(channel.ID, messages[0].ID)
	}
}

func (c *AdminCommands) ackmessage(ctx *router.Ctx) {
	channel, err := ctx.Arguments.Get(0).Channel()
	if err != nil {
		ctx.RespondText("Could not find channel")
		return
	}
	_ = ctx.Session.MessageAck(channel.ID, ctx.Arguments.Get(1).Raw())
}

// Webhook Control

func (c *AdminCommands) listwebhooks(ctx *router.Ctx) {
	channel, err := ctx.Arguments.Get(0).Channel(router.ArgumentFlagFallbackDefault)
	if err != nil {
		ctx.RespondText("Could not find channel")
		return
	}

	webhooks, err := ctx.Session.ChannelWebhooks(channel.ID)
	if err != nil {
		ctx.RespondText("Error getting webhooks: %v", err)
		return
	}

	if len(webhooks) == 0 {
		ctx.RespondText("No webhooks found in channel")
		return
	}

	var fmtWebhooks []string

	for _, webhook := range webhooks {
		fmtWebhooks = append(
			fmtWebhooks,
			webhook.Name+" - "+discordgo.EndpointWebhookToken(webhook.ID, webhook.Token),
		)
	}

	c.services.Embed.New(ctx).
		SetDescription("%s", "Webhooks in "+channel.Mention()).
		AddField("Webhooks", strings.Join(fmtWebhooks, "\n")).
		Send()
}

func (c *AdminCommands) createwebhook(ctx *router.Ctx) {
	channel, err := ctx.Arguments.Get(0).Channel(router.ArgumentFlagFallbackDefault)
	if err != nil {
		ctx.RespondText("Could not find channel")
		return
	}

	name := ctx.Arguments.Get(1).Raw()
	if name == "" {
		name = "LightningBot Webhook"
	}

	webhook, err := ctx.Session.WebhookCreate(channel.ID, name, "")
	if err != nil {
		ctx.RespondText("Error creating webhook: %v", err)
		return
	}

	c.services.Embed.New(ctx).
		SetDescription("Successfully created webhook!").
		AddField("Webhook", webhook.Name+" - "+discordgo.EndpointWebhookToken(webhook.ID, webhook.Token)).
		Send()
}

func (c *AdminCommands) webhookinfo(ctx *router.Ctx) {
	webhookUrl, err := url.Parse(ctx.Arguments.Get(0).Raw())
	if err != nil {
		ctx.RespondText("Invalid webhook URL!")
		return
	}

	if !strings.Contains(webhookUrl.Host, "discord.com") {
		ctx.RespondText("Invalid webhook URL!")
		return
	}

	resp, err := http.DefaultClient().Get(webhookUrl.String())
	if err != nil || resp.StatusCode != 200 {
		ctx.RespondText("Invalid webhook URL!")
		return
	}

	webhookInfo := resp.Json()

	c.services.Embed.New(ctx).
		AddField("Type", map[int]string{1: "Incoming", 2: "Channel Follower", 3: "Application"}[webhookInfo.Get("type").MustInt()]).
		AddField("ID", webhookInfo.Get("id").MustString()).
		AddField("Name", webhookInfo.Get("name").MustString()).
		AddField("Avatar", webhookInfo.Get("avatar").MustString()).
		AddField("Channel ID", webhookInfo.Get("channel_id").MustString()).
		AddField("Server ID", webhookInfo.Get("guild_id").MustString()).
		AddField("Application ID", webhookInfo.Get("application_id").MustString()).
		AddField("Token", webhookInfo.Get("token").MustString()).
		AddField("URL", webhookUrl.String()).
		Send()
}

func (c *AdminCommands) webhooksend(ctx *router.Ctx) {
	webhookUrl, err := url.Parse(ctx.Arguments.Get(0).Raw())
	if err != nil {
		ctx.RespondText("Invalid webhook URL!")
		return
	}

	if !strings.Contains(webhookUrl.Host, "discord.com") {
		ctx.RespondText("Invalid webhook URL!")
		return
	}

	data := map[string]any{
		"content": ctx.Arguments.RawFrom(1),
	}

	resp, err := http.DefaultClient().Post(webhookUrl.String(), data)
	if err != nil || resp.StatusCode != 204 {
		ctx.RespondText("%s", "Error sending message to webhook: "+string(resp.MustGetBody()))
		return
	}

	c.services.Embed.New(ctx).
		SetDescription("Successfully sent message to webhook.").
		Send()
}

func (c *AdminCommands) webhookdelete(ctx *router.Ctx) {
	webhookUrl, err := url.Parse(ctx.Arguments.Get(0).Raw())
	if err != nil {
		ctx.RespondText("Invalid webhook URL!")
		return
	}

	if !strings.Contains(webhookUrl.Host, "discord.com") {
		ctx.RespondText("Invalid webhook URL!")
		return
	}

	_, _ = http.DefaultClient().Delete(webhookUrl.String(), nil)

	c.services.Embed.New(ctx).
		SetDescription("Successfully deleted webhook.").
		Send()
}

// Member Control

func (c *AdminCommands) userinfo(ctx *router.Ctx) {
	user, err := ctx.Arguments.Get(0).User(router.ArgumentFlagFallbackDefault)
	if err != nil {
		ctx.RespondText("Cannot find user")
		return
	}

	var member *discordgo.Member
	if ctx.Message.GuildID != "" {
		member, _ = utils.Discord.GetMember(ctx.Session, ctx.Message.GuildID, user.ID, true)
	}

	emb := c.services.Embed.New(ctx).
		SetDescription("User Info (%s - %s)", user.String(), user.ID).
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
			AddField("Roles", fmt.Sprintf("%d Roles:\n%s", len(member.Roles), utils.Discord.GetRoleNames(ctx.Session, member)))

		if p := member.Presence; p != nil {
			emb.AddField("Status", utils.Discord.FormatStatusEmoji(p.Status)).
				AddField("Platform", utils.Discord.FormatClientStatus(p.ClientStatus))

			if len(p.Activities) > 0 {
				emb.AddField("Activity", utils.Discord.FormatActivities(p.Activities))
			}
		}
	}

	emb.Send()
}

func (c *AdminCommands) userage(ctx *router.Ctx) {
	user, err := ctx.Arguments.Get(0).User()
	if err != nil {
		ctx.RespondText("Could not find user")
		return
	}
	age := time.Since(utils.Discord.SnowflakeToTime(user.ID))
	c.services.Embed.New(ctx).
		SetDescription("User Age (%s - %s)", user.String(), user.ID).
		AddField("Age", utils.Misc.FormatDuration(age)).
		Send()
}

func (c *AdminCommands) mutualfriends(ctx *router.Ctx) {
	user, err := ctx.Arguments.Get(0).User()
	if err != nil {
		ctx.RespondText("Could not find user")
		return
	}
	mutualFriends, err := ctx.Session.UserMutualFriends(user.ID)
	if err != nil {
		ctx.RespondText("Could not find user")
		return
	}
	mutualFriendsString := ""
	for _, friend := range mutualFriends {
		mutualFriendsString += fmt.Sprintf("%s - %s\n", friend.String(), friend.ID)
	}
	c.services.Embed.New(ctx).
		SetDescription("Mutual Friends (%s - %s)", user.String(), user.ID).
		AddField("Mutual Friends", mutualFriendsString).
		Send()
}

func (c *AdminCommands) mutualservers(ctx *router.Ctx) {
	user, err := ctx.Arguments.Get(0).User()
	if err != nil {
		ctx.RespondText("Could not find user")
		return
	}
	mutualServers, err := ctx.Session.UserProfile(user.ID)
	if err != nil {
		ctx.RespondText("Could not find user")
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
	c.services.Embed.New(ctx).
		SetDescription("Mutual Servers (%s - %s)", user.String(), user.ID).
		AddField("Mutual Servers", mutualServersString).
		Send()
}

func (c *AdminCommands) platform(ctx *router.Ctx) {
	member, err := ctx.Arguments.Get(0).Member()
	if err != nil {
		ctx.RespondText("Could not find user")
		return
	}

	if member.Presence == nil {
		ctx.RespondText("Users presence is unknown!")
		return
	}

	c.services.Embed.New(ctx).
		AddField("User", member.User.Username).
		AddField("Status", utils.Discord.FormatStatusEmoji(member.Presence.Status)).
		AddField("Platform", utils.Discord.FormatClientStatus(member.Presence.ClientStatus)).
		Send()
}

func (c *AdminCommands) avatar(ctx *router.Ctx) {
	user, err := ctx.Arguments.Get(0).User()
	if err != nil {
		ctx.RespondText("Could not find user")
		return
	}
	if user.Avatar == "" {
		ctx.RespondText("User has no avatar")
		return
	}
	ctx.RespondText("%s", user.AvatarURL(""))
}

func (c *AdminCommands) banner(ctx *router.Ctx) {
	user, err := ctx.Arguments.Get(0).User()
	if err != nil {
		ctx.RespondText("Could not find user")
		return
	}
	if user.Banner == "" {
		ctx.RespondText("User has no banner")
		return
	}
	ctx.RespondText("%s", user.BannerURL(""))
}

func (c *AdminCommands) mute(ctx *router.Ctx) {
	member, err := ctx.Arguments.Get(0).Member()
	if err != nil {
		ctx.RespondText("Could not find user")
		return
	}

	if err := utils.Discord.MuteUser(ctx.Session, member); err != nil {
		ctx.RespondText("%s", utils.String.Capitalize(err.Error()))
		return
	}

	c.services.Embed.New(ctx).
		SetDescription("Successfully muted user!").
		Send()
}

func (c *AdminCommands) unmute(ctx *router.Ctx) {
	member, err := ctx.Arguments.Get(0).Member()
	if err != nil {
		ctx.RespondText("Could not find user")
		return
	}

	roles, err := ctx.Session.GuildRoles(ctx.Message.GuildID)
	if member == nil || err != nil {
		ctx.RespondText("Error unmuting user: %v", err)
		return
	}

	for _, role := range roles {
		if strings.ToLower(role.Name) != "muted" {
			continue
		}

		if !slices.Contains(member.Roles, role.ID) {
			c.services.Embed.New(ctx).
				SetDescription("User is already unmuted!").
				Send()
			return
		}

		for i, roleID := range member.Roles {
			if roleID == role.ID {
				member.Roles = append(member.Roles[:i], member.Roles[i+1:]...)
				break
			}
		}

		_, err = ctx.Session.GuildMemberEdit(ctx.Message.GuildID, member.User.ID, &discordgo.GuildMemberParams{
			Roles: &member.Roles,
		})
		if err != nil {
			ctx.RespondText("Error unmuting user: %v", err)
			return
		}

		c.services.Embed.New(ctx).
			SetDescription("Successfully unmuted user!").
			Send()

		return
	}
	ctx.RespondText("Error: No muted role found!")
}

func (c *AdminCommands) timeout(ctx *router.Ctx) {
	member, err := ctx.Arguments.Get(0).Member()
	if err != nil {
		ctx.RespondText("Could not find user")
		return
	}

	duration, err := ctx.Arguments.Get(1).Duration()
	duration = fn.Clamp(duration, 1*time.Second, 28*24*time.Hour)
	if err != nil {
		ctx.RespondText("Could not parse duration")
		return
	}

	until := time.Now().Add(duration)

	err = ctx.Session.GuildMemberTimeout(ctx.Message.GuildID, member.User.ID, &until)
	if err != nil {
		ctx.RespondText("Error timing out user: %v", err)
		return
	}

	c.services.Embed.New(ctx).
		SetDescription("Successfully timed out user for %s!", duration.String()).
		Send()
}

func (c *AdminCommands) removetimeout(ctx *router.Ctx) {
	member, err := ctx.Arguments.Get(0).Member()
	if err != nil {
		ctx.RespondText("Could not find user")
		return
	}

	err = ctx.Session.GuildMemberTimeout(ctx.Message.GuildID, member.User.ID, nil)
	if err != nil {
		ctx.RespondText("Error removing timeout: %v", err)
		return
	}

	c.services.Embed.New(ctx).
		SetDescription("Successfully removed timeout!").
		Send()
}

func (c *AdminCommands) kick(ctx *router.Ctx) {
	member, err := ctx.Arguments.Get(0).Member()
	if err != nil {
		ctx.RespondText("Could not find user")
		return
	}

	err = ctx.Session.GuildMemberDeleteWithReason(ctx.Message.GuildID, member.User.ID, ctx.Arguments.RawFrom(1))
	if err != nil {
		ctx.RespondText("Error kicking user: %v", err)
		return
	}

	c.services.Embed.New(ctx).
		SetDescription("Successfully kicked user!").
		Send()
}

func (c *AdminCommands) ban(ctx *router.Ctx) {
	member, err := ctx.Arguments.Get(0).Member(router.ArgumentFlagFallbackDefault)
	if err != nil {
		ctx.RespondText("Could not find user")
		return
	}

	err = ctx.Session.GuildBanCreateWithReason(ctx.Message.GuildID, member.User.ID, ctx.Arguments.RawFrom(1), 0)
	if err != nil {
		ctx.RespondText("Error banning user: %v", err)
		return
	}

	c.services.Embed.New(ctx).
		SetDescription("Successfully banned user!").
		Send()
}

func (c *AdminCommands) softban(ctx *router.Ctx) {
	member, err := ctx.Arguments.Get(0).Member()
	if err != nil {
		ctx.RespondText("Could not find user")
		return
	}

	err = ctx.Session.GuildBanCreateWithReason(ctx.Message.GuildID, member.User.ID, ctx.Arguments.RawFrom(1), 0)
	if err != nil {
		ctx.RespondText("Error softbanning user: %v", err)
		return
	}
	err = ctx.Session.GuildBanDelete(ctx.Message.GuildID, member.User.ID)
	if err != nil {
		ctx.RespondText("Error softbanning user: %v", err)
		return
	}

	c.services.Embed.New(ctx).
		SetDescription("Successfully softbanned user!").
		Send()
}

func (c *AdminCommands) unban(ctx *router.Ctx) {
	err := ctx.Session.GuildBanDelete(ctx.Message.GuildID, ctx.Arguments.Get(0).UserMentionID())
	if err != nil {
		ctx.RespondText("Error unbanning user: %v", err)
		return
	}

	c.services.Embed.New(ctx).
		SetDescription("Successfully unbanned user!").
		Send()
}

func (c *AdminCommands) banlist(ctx *router.Ctx) {
	guild, err := ctx.Arguments.Get(0).Guild()
	if err != nil {
		ctx.RespondText("Error getting server: %v", err)
		return
	}

	bans, err := ctx.Session.GuildBans(guild.ID, 100, "", "")
	if err != nil {
		ctx.RespondText("Error getting ban list: %v", err)
		return
	}

	banReadable := fn.Map(bans, func(ban *discordgo.GuildBan) string {
		return fmt.Sprintf("%s - %s", ban.User.Username, ban.User.ID)
	})

	c.services.Embed.New(ctx).
		AddField("Banned Amount", strconv.Itoa(len(bans))).
		AddField("Ban List", strings.Join(banReadable, "\n")).
		Send()
}

func (c *AdminCommands) copybans(ctx *router.Ctx) {
	srcGuild, err := ctx.Arguments.Get(0).Guild()
	if err != nil {
		ctx.RespondText("Error getting server: %v", err)
		return
	}

	destGuild, err := ctx.Arguments.Get(1).Guild(router.ArgumentFlagFallbackDefault)
	if err != nil {
		ctx.RespondText("Error getting server: %v", err)
		return
	}

	bans, err := ctx.Session.GuildBans(srcGuild.ID, 1000, "", "")
	if err != nil {
		ctx.RespondText("Error copying bans: %v", err)
		return
	}

	for _, ban := range bans {
		err := ctx.Session.GuildBanCreateWithReason(destGuild.ID, ban.User.ID, ban.Reason, 0)
		if err != nil {
			ctx.RespondText("Error copying bans: %v", err)
			return
		}

		time.Sleep(2 * time.Second)
	}

	c.services.Embed.New(ctx).
		SetDescription("Successfully copied bans from %s to %s!", srcGuild.Name, destGuild.Name).
		Send()
}

// Nickname Control

func (c *AdminCommands) nick(ctx *router.Ctx) {
	err := ctx.Session.GuildMemberNickname(ctx.Message.GuildID, "@me", ctx.Arguments.Raw())
	if err != nil {
		ctx.RespondText("Error changing nickname: %v", err)
		return
	}
	c.services.Embed.New(ctx).
		SetDescription("Successfully changed nickname!").
		Send()
}

func (c *AdminCommands) nickuser(ctx *router.Ctx) {
	member, err := ctx.Arguments.Get(0).Member()
	if err != nil {
		ctx.RespondText("Could not find user")
		return
	}

	err = ctx.Session.GuildMemberNickname(ctx.Message.GuildID, strings.ReplaceAll(member.User.ID, ctx.Message.Author.ID, "@me"), ctx.Arguments.RawFrom(1))
	if err != nil {
		ctx.RespondText("Error changing nickname: %v", err)
		return
	}
	c.services.Embed.New(ctx).
		SetDescription("Successfully changed nickname!").
		Send()
}

func (c *AdminCommands) clearnick(ctx *router.Ctx) {
	member, err := ctx.Arguments.Get(0).Member()
	if err != nil {
		ctx.RespondText("Could not find user")
		return
	}

	err = ctx.Session.GuildMemberNickname(ctx.Message.GuildID, member.User.ID, "")
	if err != nil {
		ctx.RespondText("Error changing nickname: %v", err)
		return
	}
	c.services.Embed.New(ctx).
		SetDescription("Successfully changed nickname!").
		Send()
}

// Server Control

func (c *AdminCommands) serverinfo(ctx *router.Ctx) {
	guild, err := ctx.Arguments.Get(0).Guild(router.ArgumentFlagFallbackDefault)
	if err != nil {
		ctx.RespondText("Error getting server: %v", err)
		return
	}

	owner, ok := utils.Discord.GetMember(ctx.Session, guild.ID, guild.OwnerID)
	ownerMention := "Error"
	if ok {
		ownerMention = owner.User.String()
	}

	c.services.Embed.New(ctx).
		SetDescription("Server Info - %s", guild.Name).
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
		Send()
}

func (c *AdminCommands) serverage(ctx *router.Ctx) {
	guild, err := ctx.Arguments.Get(0).Guild(router.ArgumentFlagFallbackDefault)
	if err != nil {
		ctx.RespondText("Error getting server: %v", err)
		return
	}
	age := time.Since(utils.Discord.SnowflakeToTime(guild.ID))
	c.services.Embed.New(ctx).
		SetDescription("Server Age (%s - %s)", guild.Name, guild.ID).
		AddField("Age", utils.Misc.FormatDuration(age)).
		Send()
}

func (c *AdminCommands) serverowner(ctx *router.Ctx) {
	guild, err := ctx.Arguments.Get(0).Guild(router.ArgumentFlagFallbackDefault)
	if err != nil {
		ctx.RespondText("Error getting server: %v", err)
		return
	}
	owner, ok := utils.Discord.GetMember(ctx.Session, guild.ID, guild.OwnerID)
	if !ok {
		ctx.RespondText("Error getting server owner")
		return
	}
	c.services.Embed.New(ctx).
		SetDescription("%s", guild.Name).
		AddField("Owner", owner.Mention()).
		Send()
}

func (c *AdminCommands) servericon(ctx *router.Ctx) {
	guild, err := ctx.Arguments.Get(0).Guild(router.ArgumentFlagFallbackDefault)
	if err != nil {
		ctx.RespondText("Error getting server: %v", err)
		return
	}

	emb := c.services.Embed.New(ctx).
		SetDescription("Server Icon - %s", guild.Name).
		SetImage(guild.IconURL())
	if guild.Icon == "" {
		emb.AddField("Error", "Server has no icon")
	}
	emb.Send()
}

func (c *AdminCommands) serverbanner(ctx *router.Ctx) {
	guild, err := ctx.Arguments.Get(0).Guild(router.ArgumentFlagFallbackDefault)
	if err != nil {
		ctx.RespondText("Error getting server: %v", err)
		return
	}
	if guild.Banner == "" {
		c.services.Embed.New(ctx).AddField("Error", "Server has no banner!").Send()
	}
	c.services.Embed.New(ctx).
		SetDescription("Server Banner - %s", guild.Name).
		SetImage(guild.BannerURL()).
		Send()
}

func (c *AdminCommands) serversplash(ctx *router.Ctx) {
	guild, err := ctx.Arguments.Get(0).Guild(router.ArgumentFlagFallbackDefault)
	if err != nil {
		ctx.RespondText("Error getting server: %v", err)
		return
	}
	if guild.Splash == "" {
		ctx.RespondText("Server has no invite splash!")
		return
	}
	c.services.Embed.New(ctx).
		SetDescription("Server Invite Splash - %s", guild.Name).
		SetImage(discordgo.EndpointGuildSplash(guild.ID, guild.Splash)).
		Send()
}

func (c *AdminCommands) serverdsplash(ctx *router.Ctx) {
	guild, err := ctx.Arguments.Get(0).Guild(router.ArgumentFlagFallbackDefault)
	if err != nil {
		ctx.RespondText("Error getting server: %v", err)
		return
	}
	if guild.DiscoverySplash == "" {
		ctx.RespondText("Server has no discovery splash!")
		return
	}
	c.services.Embed.New(ctx).
		SetDescription("Server Discovery Splash - %s", guild.Name).
		SetImage(discordgo.EndpointGuildSplash(guild.ID, guild.DiscoverySplash)).
		Send()
}

func (c *AdminCommands) serverboosters(ctx *router.Ctx) {
	guild, err := ctx.Arguments.Get(0).Guild(router.ArgumentFlagFallbackDefault)
	if err != nil {
		ctx.RespondText("Error getting server: %v", err)
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

	c.services.Embed.New(ctx).
		SetDescription("Server Boosters - %s", guild.Name).
		AddField("Boosters", boostersString).
		Send()
}

func (c *AdminCommands) serverboostlevel(ctx *router.Ctx) {
	guild, err := ctx.Arguments.Get(0).Guild(router.ArgumentFlagFallbackDefault)
	if err != nil {
		ctx.RespondText("Error getting server: %v", err)
		return
	}

	c.services.Embed.New(ctx).
		SetDescription("Server Boost Level - %s", guild.Name).
		AddField("Boost Level", boostLevels[int(guild.PremiumTier)]).
		Send()
}

func (c *AdminCommands) setservername(ctx *router.Ctx) {
	guild, err := ctx.Session.GuildEdit(ctx.Message.GuildID, &discordgo.GuildParams{
		Name: ctx.Arguments.Raw(),
	})
	if err != nil {
		ctx.RespondText("Error changing server name: %v", err)
		return
	}
	c.services.Embed.New(ctx).
		SetDescription("Successfully changed server name to %s!", guild.Name).
		Send()
}

func (c *AdminCommands) setservericon(ctx *router.Ctx) {
	_, err := ctx.Session.GuildEdit(ctx.Message.GuildID, &discordgo.GuildParams{
		Icon: ctx.Arguments.Get(0).Raw(),
	})
	if err != nil {
		ctx.RespondText("Error changing server icon: %v", err)
		return
	}
	c.services.Embed.New(ctx).
		SetDescription("Successfully changed server icon!").
		Send()
}

func (c *AdminCommands) setserverbanner(ctx *router.Ctx) {
	_, err := ctx.Session.GuildEdit(ctx.Message.GuildID, &discordgo.GuildParams{
		Banner: ctx.Arguments.Get(0).Raw(),
	})
	if err != nil {
		ctx.RespondText("Error changing server icon: %v", err)
		return
	}
	c.services.Embed.New(ctx).
		SetDescription("Successfully changed server icon!").
		Send()
}

func (c *AdminCommands) lockserver(ctx *router.Ctx) {
	guild, err := ctx.Arguments.Get(0).Guild(router.ArgumentFlagFallbackDefault)
	if err != nil {
		ctx.RespondText("Error getting server: %v", err)
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
			Position:             channel.Position,
			PermissionOverwrites: overwrites,
		})
		if err != nil {
			ctx.RespondText("Error locking server: %v", err)
			return
		}
	}

	c.services.Embed.New(ctx).
		SetDescription("Successfully locked server!").
		Send()
}

func (c *AdminCommands) unlockserver(ctx *router.Ctx) {
	guild, err := ctx.Arguments.Get(0).Guild(router.ArgumentFlagFallbackDefault)
	if err != nil {
		ctx.RespondText("Error getting server: %v", err)
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
			Position:             channel.Position,
			PermissionOverwrites: overwrites,
		})
		if err != nil {
			ctx.RespondText("Error unlocking server: %v", err)
			return
		}
	}

	c.services.Embed.New(ctx).
		SetDescription("Successfully unlocked server!").
		Send()
}

func (c *AdminCommands) channels(ctx *router.Ctx) {
	guild, err := ctx.Arguments.Get(0).Guild(router.ArgumentFlagFallbackDefault)
	if err != nil {
		ctx.RespondText("Error getting server: %v", err)
		return
	}

	categories, categoryChannels, looseChannels := utils.Discord.MapCategories(guild.Channels)
	emb := c.services.Embed.New(ctx).SetDescription("Server Channels - %s", guild.Name)

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
	emb.Send()
}

func (c *AdminCommands) bots(ctx *router.Ctx) {
	guild, err := ctx.Arguments.Get(0).Guild(router.ArgumentFlagFallbackDefault)
	if err != nil {
		ctx.RespondText("Error getting server: %v", err)
		return
	}

	bots := fn.Filter(guild.Members, func(member *discordgo.Member) bool { return member.User.Bot })
	if len(bots) == 0 {
		ctx.RespondText("No bots found!")
		return
	}

	botsReadable := fn.Map(bots, func(member *discordgo.Member) string {
		return fmt.Sprintf("%s - %s", member.User.String(), member.User.ID)
	})

	c.services.Embed.New(ctx).
		SetDescription("Server Bots - %s", guild.Name).
		AddField("Server Bots", strings.Join(botsReadable, "\n")).
		Send()
}

func (c *AdminCommands) admins(ctx *router.Ctx) {
	guild, err := ctx.Arguments.Get(0).Guild(router.ArgumentFlagFallbackDefault)
	if err != nil {
		ctx.RespondText("Error getting server: %v", err)
		return
	}

	if len(guild.Channels) == 0 {
		ctx.RespondText("Server has no channels!")
		return
	}

	admins := fn.Filter(guild.Members, func(member *discordgo.Member) bool {
		permissions, _ := ctx.Session.UserChannelPermissions(member.User.ID, guild.Channels[0].ID)
		return permissions&discordgo.PermissionAdministrator != 0 && !member.User.Bot
	})
	if len(admins) == 0 {
		ctx.RespondText("No admins found!")
		return
	}

	adminsReadable := fn.Map(admins, func(member *discordgo.Member) string {
		return fmt.Sprintf("%s - %s", member.User.String(), member.User.ID)
	})

	c.services.Embed.New(ctx).
		SetDescription("Server Admins - %s", guild.Name).
		AddField("Server Admins", strings.Join(adminsReadable, "\n")).
		Send()
}

func (c *AdminCommands) roles(ctx *router.Ctx) {
	guild, err := ctx.Arguments.Get(0).Guild(router.ArgumentFlagFallbackDefault)
	if err != nil {
		ctx.RespondText("Error getting server: %v", err)
		return
	}

	slices.SortFunc(guild.Roles, func(i, j *discordgo.Role) int {
		return i.Position - j.Position
	})

	rolesReadable := fn.Map(guild.Roles, func(role *discordgo.Role) string {
		return fmt.Sprintf("%s - %s", role.Name, role.ID)
	})

	c.services.Embed.New(ctx).
		SetDescription("Server Roles - %s", guild.Name).
		AddField("Server Roles", strings.Join(rolesReadable, "\n")).
		Send()
}

func (c *AdminCommands) leaveserver(ctx *router.Ctx) {
	guild, err := ctx.Arguments.Get(0).Guild(router.ArgumentFlagFallbackDefault)
	if err != nil {
		ctx.RespondText("Error getting server: %v", err)
		return
	}
	_ = ctx.Session.GuildLeave(guild.ID)
}

func (c *AdminCommands) deleteserver(ctx *router.Ctx) {
	guild, err := ctx.Arguments.Get(0).Guild(router.ArgumentFlagFallbackDefault)
	if err != nil {
		ctx.RespondText("Error getting server: %v", err)
		return
	}

	_, err = ctx.Session.GuildDelete(guild.ID)
	if err != nil {
		ctx.RespondText("Error deleting server: %v", err)
		return
	}
}

// Invite Control

func (c *AdminCommands) invite(ctx *router.Ctx) {
	channel, err := ctx.Arguments.Get(0).Channel(router.ArgumentFlagFallbackDefault)
	if err != nil {
		ctx.RespondText("Could not find channel")
		return
	}

	maxAge, _ := ctx.Arguments.Get(1).Int()
	maxAge = fn.Clamp(maxAge, 0, 604800)

	maxUses, _ := ctx.Arguments.Get(2).Int()
	maxUses = fn.Clamp(maxUses, 0, 100)

	invite, err := ctx.Session.ChannelInviteCreate(channel.ID, maxAge, maxUses, false)
	if err != nil {
		ctx.RespondText("Error creating invite: %v", err)
		return
	}

	c.services.Embed.New(ctx).
		SetDescription("Created invite: https://discord.gg/%s", invite.Code).
		Send()
}

func (c *AdminCommands) friendinvite(ctx *router.Ctx) {
	invite, err := ctx.Session.CreateFriendInvite("@me")
	if err != nil {
		ctx.RespondText("Error creating friend invite: %v", err)
		return
	}

	c.services.Embed.New(ctx).
		SetDescription("Successfully created friend invite!").
		AddField("Invite", fmt.Sprintf("https://discord.gg/%s", invite.Code)).
		Send()
}

func (c *AdminCommands) inviteinfo(ctx *router.Ctx) {
	parsedInvite := utils.Discord.ParseInvite(ctx.Arguments.Get(0).Raw())
	invite, err := ctx.Session.InviteWithCounts(parsedInvite)
	if err != nil {
		ctx.RespondText("Error getting invite: %v", err)
		return
	}

	emb := c.services.Embed.New(ctx).SetDescription("Invite Info - %s", invite.Code)
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
		Send()
}

func (c *AdminCommands) deleteinvite(ctx *router.Ctx) {
	parsedInvite := utils.Discord.ParseInvite(ctx.Arguments.Get(0).Raw())
	invite, err := ctx.Session.InviteWithCounts(parsedInvite)
	if err != nil {
		ctx.RespondText("Error getting invite: %v", err)
		return
	}

	invite, err = ctx.Session.InviteDelete(invite.Code)
	if err != nil {
		ctx.RespondText("Error deleting invite: %v", err)
		return
	}

	c.services.Embed.New(ctx).
		SetDescription("Deleted invite: %s", invite.Code).
		Send()
}

func (c *AdminCommands) serverinvites(ctx *router.Ctx) {
	guild, err := ctx.Arguments.Get(0).Guild(router.ArgumentFlagFallbackDefault)
	if err != nil {
		ctx.RespondText("Error getting server: %v", err)
		return
	}
	invites, err := ctx.Session.GuildInvites(guild.ID)
	if err != nil {
		ctx.RespondText("Error getting invites: %v", err)
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
	c.services.Embed.New(ctx).
		SetDescription("Invite List").
		AddField("Invites", invitesString).
		Send()
}

func (c *AdminCommands) channelinvites(ctx *router.Ctx) {
	channel, err := ctx.Arguments.Get(0).Channel(router.ArgumentFlagFallbackDefault)
	if err != nil {
		ctx.RespondText("Could not find channel")
		return
	}
	invites, err := ctx.Session.ChannelInvites(channel.ID)
	if err != nil {
		ctx.RespondText("Error getting invites: %v", err)
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
	c.services.Embed.New(ctx).
		SetDescription("Invite List").
		AddField("Invites", invitesString).
		Send()
}

func (c *AdminCommands) userinvites(ctx *router.Ctx) {
	member, err := ctx.Arguments.Get(0).Member()
	if err != nil {
		ctx.RespondText("Could not find user")
		return
	}

	guild, err := ctx.Arguments.Get(1).Guild(router.ArgumentFlagFallbackDefault)
	if err != nil {
		ctx.RespondText("Error getting server: %v", err)
		return
	}

	invites, err := ctx.Session.GuildInvites(guild.ID)
	if err != nil {
		ctx.RespondText("Error getting invites: %v", err)
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
	c.services.Embed.New(ctx).
		SetDescription("Invite List (%s)", member.User.String()).
		AddField("Invites", invitesString).
		Send()
}

func (c *AdminCommands) roleinfo(ctx *router.Ctx) {
	role, err := ctx.Arguments.Get(0).Role()
	if err != nil {
		ctx.RespondText("Could not find role")
		return
	}
	// Get all members in guild and check if they have the role
	guild, err := utils.Discord.GetGuild(ctx.Session, ctx.Message.GuildID)
	if err != nil {
		ctx.RespondText("Error getting server: %v", err)
		return
	}
	memberCount := 0
	for _, member := range guild.Members {
		if slices.Contains(member.Roles, role.ID) {
			memberCount++
		}
	}
	createdAt := utils.Discord.SnowflakeToTime(role.ID)
	perms, _ := api.LBCreatePaste(utils.Discord.FormatPermissions(role.Name, role.ID, role.Permissions))
	c.services.Embed.New(ctx).
		SetDescription("Role Info (%s)", role.Name).
		AddField("ID", role.ID).
		AddField("Color", fmt.Sprintf("#%x", role.Color)).
		AddField("Created at", utils.Discord.GetTimestamp(&createdAt, "Error")).
		AddField("Position", fmt.Sprint(role.Position)).
		AddField("Is Hoist", utils.Misc.FormatBoolYesNo(role.Hoist)).
		AddField("Is Mentionable", utils.Misc.FormatBoolYesNo(role.Mentionable)).
		AddField("Members With Role", fmt.Sprint(memberCount)).
		AddField("Permissions", perms).
		Send()
}

func (c *AdminCommands) giverole(ctx *router.Ctx) {
	member, err := ctx.Arguments.Get(0).Member()
	if err != nil {
		ctx.RespondText("Could not find user")
		return
	}
	role, err := ctx.Arguments.Get(1).Role()
	if err != nil {
		ctx.RespondText("Could not find role")
		return
	}

	err = ctx.Session.GuildMemberRoleAdd(ctx.Message.GuildID, member.User.ID, role.ID)
	if err != nil {
		ctx.RespondText("Error giving role: %v", err)
		return
	}
	c.services.Embed.New(ctx).
		SetDescription("Successfully gave role (%s) to (%s)", role.Name, member.User.String()).
		Send()
}

func (c *AdminCommands) removerole(ctx *router.Ctx) {
	member, err := ctx.Arguments.Get(0).Member()
	if err != nil {
		ctx.RespondText("Could not find user")
		return
	}
	role, err := ctx.Arguments.Get(1).Role()
	if err != nil {
		ctx.RespondText("Could not find role")
		return
	}

	err = ctx.Session.GuildMemberRoleRemove(ctx.Message.GuildID, member.User.ID, role.ID)
	if err != nil {
		ctx.RespondText("Error removing role: %v", err)
		return
	}
	c.services.Embed.New(ctx).
		SetDescription("Successfully removed role (%s) from (%s)", role.Name, member.User.String()).
		Send()
}

func (c *AdminCommands) createrole(ctx *router.Ctx) {
	role, err := ctx.Session.GuildRoleCreate(ctx.Message.GuildID, &discordgo.RoleParams{
		Name: ctx.Arguments.Raw(),
	})
	if err != nil {
		ctx.RespondText("Error creating role: %v", err)
		return
	}

	c.services.Embed.New(ctx).
		SetDescription("Successfully created role (%s)", role.Name).
		Send()
}

func (c *AdminCommands) deleterole(ctx *router.Ctx) {
	role, err := ctx.Arguments.Get(0).Role()
	if err != nil {
		ctx.RespondText("Could not find role")
		return
	}
	err = ctx.Session.GuildRoleDelete(ctx.Message.GuildID, role.ID)
	if err != nil {
		ctx.RespondText("Error deleting role: %v", err)
		return
	}
	c.services.Embed.New(ctx).
		SetDescription("Successfully deleted role (%s)", role.Name).
		Send()
}

func (c *AdminCommands) renamerole(ctx *router.Ctx) {
	role, err := ctx.Arguments.Get(0).Role()
	if err != nil {
		ctx.RespondText("Could not find role")
		return
	}

	role, err = ctx.Session.GuildRoleEdit(ctx.Message.GuildID, role.ID, &discordgo.RoleParams{
		Name: ctx.Arguments.Raw(),
	})
	if err != nil {
		ctx.RespondText("Error renaming role: %v", err)
		return
	}

	c.services.Embed.New(ctx).
		SetDescription("Successfully renamed role (%s)", role.Name).
		Send()
}
