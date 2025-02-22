package admin

import (
	"fmt"
	"net/url"
	"slices"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/LightningDev1/LB-Selfbot/internal/api"
	"github.com/LightningDev1/LB-Selfbot/internal/http"
	"github.com/LightningDev1/LB-Selfbot/internal/router"
	"github.com/LightningDev1/LB-Selfbot/internal/services"
	"github.com/LightningDev1/LB-Selfbot/internal/utils"
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

type AdminCommands struct{ services *services.Services }

func NewAdminCommands(s *services.Services) *AdminCommands {
	return &AdminCommands{services: s}
}

func (c *AdminCommands) Register() {
	cmdRouter := c.services.Discord.GetRouter()

	cmdRouter.RegisterCategory(&router.Category{
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
			Arguments:   "<#channel> <seconds>",
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
			Arguments:   "<user> <seconds>",
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
			Description: "Ban and instantly unban a user",
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
		_ = ctx.RespondText("Could not find channel")
		return
	}

	c.services.Embed.NewEmbed().
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

func (c *AdminCommands) textchannel(ctx *router.Ctx) {
	channel, err := ctx.Session.GuildChannelCreate(ctx.Message.GuildID, ctx.Arguments.Raw(), discordgo.ChannelTypeGuildText)
	if err != nil {
		_ = ctx.RespondText("Error creating channel: " + err.Error())
		return
	}
	c.services.Embed.NewEmbed().
		SetDescription("Successfully created text channel!").
		AddField("Channel", channel.Mention()).
		Send(ctx)
}

func (c *AdminCommands) voicechannel(ctx *router.Ctx) {
	channel, err := ctx.Session.GuildChannelCreate(ctx.Message.GuildID, ctx.Arguments.Raw(), discordgo.ChannelTypeGuildVoice)
	if err != nil {
		_ = ctx.RespondText("Error creating channel: " + err.Error())
		return
	}
	c.services.Embed.NewEmbed().
		SetDescription("Successfully created voice channel!").
		AddField("Channel", channel.Mention()).
		Send(ctx)
}

func (c *AdminCommands) stagechannel(ctx *router.Ctx) {
	channel, err := ctx.Session.GuildChannelCreate(ctx.Message.GuildID, ctx.Arguments.Raw(), discordgo.ChannelTypeGuildStageVoice)
	if err != nil {
		_ = ctx.RespondText("Error creating channel: " + err.Error())
		return
	}
	c.services.Embed.NewEmbed().
		SetDescription("Successfully created stage channel!").
		AddField("Channel", channel.Mention()).
		Send(ctx)
}

func (c *AdminCommands) newschannel(ctx *router.Ctx) {
	channel, err := ctx.Session.GuildChannelCreate(ctx.Message.GuildID, ctx.Arguments.Raw(), discordgo.ChannelTypeGuildNews)
	if err != nil {
		_ = ctx.RespondText("Error creating channel: " + err.Error())
		return
	}
	c.services.Embed.NewEmbed().
		SetDescription("Successfully created news channel!").
		AddField("Channel", channel.Mention()).
		Send(ctx)
}

func (c *AdminCommands) thread(ctx *router.Ctx) {
	channel, err := ctx.Session.ThreadStart(ctx.Message.ChannelID, ctx.Arguments.Raw(), discordgo.ChannelTypeGuildPublicThread, 0)
	if err != nil {
		_ = ctx.RespondText("Error creating thread: " + err.Error())
		return
	}
	c.services.Embed.NewEmbed().
		SetDescription("Successfully created thread!").
		AddField("Thread", channel.Mention()).
		Send(ctx)
}

func (c *AdminCommands) category(ctx *router.Ctx) {
	channel, err := ctx.Session.GuildChannelCreate(ctx.Message.GuildID, ctx.Arguments.Raw(), discordgo.ChannelTypeGuildCategory)
	if err != nil {
		_ = ctx.RespondText("Error creating category: " + err.Error())
		return
	}
	c.services.Embed.NewEmbed().
		SetDescription("Successfully created category!").
		AddField("Category", channel.Mention()).
		Send(ctx)
}

func (c *AdminCommands) deletecategory(ctx *router.Ctx) {
	channel, err := ctx.Arguments.Get(0).Channel()
	if err != nil {
		_ = ctx.RespondText("Could not find category")
		return
	}
	_, err = ctx.Session.ChannelDelete(channel.ID)
	if err != nil {
		_ = ctx.RespondText("Error deleting category: " + err.Error())
		return
	}

	c.services.Embed.NewEmbed().
		SetDescription("Successfully deleted category!").
		Send(ctx)
}

func (c *AdminCommands) renamechannel(ctx *router.Ctx) {
	channel, err := ctx.Arguments.Get(0).Channel()
	if err != nil {
		_ = ctx.RespondText("Could not find channel")
		return
	}

	_, err = ctx.Session.ChannelEdit(channel.ID, &discordgo.ChannelEdit{
		Name: ctx.Arguments.RawFrom(1),
	})
	if err != nil {
		_ = ctx.RespondText("Error renaming channel: " + err.Error())
		return
	}

	c.services.Embed.NewEmbed().
		SetDescription("Successfully renamed channel!").
		AddField("Channel", channel.Mention()).
		Send(ctx)
}

func (c *AdminCommands) deletechannel(ctx *router.Ctx) {
	channel, err := ctx.Arguments.Get(0).Channel()
	if err != nil {
		_ = ctx.RespondText("Could not find channel")
		return
	}

	_, err = ctx.Session.ChannelDelete(channel.ID)
	if err != nil {
		_ = ctx.RespondText("Error deleting channel: " + err.Error())
		return
	}

	c.services.Embed.NewEmbed().
		SetDescription("Successfully deleted channel!").
		Send(ctx)
}

func (c *AdminCommands) slowmode(ctx *router.Ctx) {
	channel, err := ctx.Arguments.Get(0).Channel()
	if err != nil {
		_ = ctx.RespondText("Could not find channel")
		return
	}

	seconds, err := strconv.Atoi(utils.Arguments["<seconds>"])
	if err != nil {
		_ = ctx.RespondText("Argument <seconds> must be a number.")
		return
	}

	_, err = ctx.Session.ChannelEditComplex(channel.ID, &discordgo.ChannelEdit{
		Position:         channel.Position,
		RateLimitPerUser: &seconds,
	})
	if err != nil {
		_ = ctx.RespondText("Error setting slowmode: " + err.Error())
		return
	}

	c.services.Embed.NewEmbed().
		SetDescription(fmt.Sprintf("Successfully set slowmode to %s seconds!", utils.Arguments["<seconds>"])).
		AddField("Channel", channel.Mention()).
		Send(ctx)
}

func (c *AdminCommands) removeslowmode(ctx *router.Ctx) {
	channel, err := ctx.Arguments.Get(0).Channel()
	if err != nil {
		_ = ctx.RespondText("Could not find channel")
		return
	}

	_, err = ctx.Session.ChannelEditComplex(channel.ID, &discordgo.ChannelEdit{
		Position:         channel.Position,
		RateLimitPerUser: nil,
	})
	if err != nil {
		_ = ctx.RespondText("Error removing slowmode: " + err.Error())
		return
	}

	c.services.Embed.NewEmbed().
		SetDescription("Successfully removed slowmode!").
		AddField("Channel", channel.Mention()).
		Send(ctx)
}

func (c *AdminCommands) lockchannel(ctx *router.Ctx) {
	channel, err := ctx.Arguments.Get(0).Channel()
	if err != nil {
		_ = ctx.RespondText("Could not find channel")
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
		_ = ctx.RespondText("Error locking channel: " + err.Error())
		return
	}

	c.services.Embed.NewEmbed().
		SetDescription("Successfully locked channel!").
		AddField("Channel", channel.Mention()).
		Send(ctx)
}

func (c *AdminCommands) unlockchannel(ctx *router.Ctx) {
	channel, err := ctx.Arguments.Get(0).Channel()
	if err != nil {
		_ = ctx.RespondText("Could not find channel")
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
		_ = ctx.RespondText("Error unlocking channel: " + err.Error())
		return
	}

	c.services.Embed.NewEmbed().
		SetDescription("Successfully unlocked channel!").
		AddField("Channel", channel.Mention()).
		Send(ctx)
}

func (c *AdminCommands) purge(ctx *router.Ctx) {
	amount, err := strconv.Atoi(utils.Arguments["<amount>"])
	if err != nil {
		_ = ctx.RespondText("Argument <amount> must be a number.")
		return
	}

	delay, err := strconv.Atoi(utils.Arguments["[delay]"])
	if err != nil {
		delay = 1
	}

	channel, err := ctx.Arguments.Get(2).Channel(router.ArgumentFlagFallbackDefault)
	if err != nil {
		_ = ctx.RespondText("Could not find channel")
		return
	}

	messages, err := utils.Discord.GetChannelMessages(ctx.Session, channel.ID, amount)
	if err != nil {
		_ = ctx.RespondText("Error getting messages: " + err.Error())
		return
	}

	for _, message := range messages {
		err = ctx.Session.ChannelMessageDelete(channel.ID, message.ID)
		if err != nil {
			_ = ctx.RespondText("Error deleting messages: " + err.Error())
			return
		}

		time.Sleep(time.Duration(delay) * time.Second)
	}
}

func (c *AdminCommands) purgeself(ctx *router.Ctx) {
	amount, err := strconv.Atoi(utils.Arguments["<amount>"])
	if err != nil {
		_ = ctx.RespondText("Argument <amount> must be a number.")
		return
	}

	delaySec, err := strconv.Atoi(utils.Arguments["[delay]"])
	if err != nil {
		delaySec = 1
	}

	channel, err := ctx.Arguments.Get(2).Channel(router.ArgumentFlagFallbackDefault)
	if err != nil {
		_ = ctx.RespondText("Could not find channel")
		return
	}

	userID := ctx.Message.Author.ID
	delay := time.Duration(delaySec) * time.Second

	// Allow the command to be deleted before the messages are deleted
	time.Sleep(1 * time.Second)

	deleted, err := utils.Discord.DeleteMessages(ctx.Session, channel.ID, userID, amount, delay)
	if err != nil {
		_ = ctx.RespondText("Error deleting messages: " + err.Error())
		return
	}

	c.services.Embed.NewEmbed().
		SetDescription(fmt.Sprintf("Successfully deleted %d messages!", deleted)).
		Send(ctx)
}

func (c *AdminCommands) purgeuser(ctx *router.Ctx) {
	amount, err := strconv.Atoi(utils.Arguments["<amount>"])
	if err != nil {
		_ = ctx.RespondText("Argument <amount> must be a number.")
		return
	}

	user, err := ctx.Arguments.Get(1).User()
	if err != nil {
		_ = ctx.RespondText("Could not find user")
		return
	}

	delaySec, err := strconv.Atoi(utils.Arguments["[delay]"])
	if err != nil {
		delaySec = 1
	}

	channel, err := ctx.Arguments.Get(2).Channel(router.ArgumentFlagFallbackDefault)
	if err != nil {
		_ = ctx.RespondText("Could not find channel")
		return
	}

	delay := time.Duration(delaySec) * time.Second

	// Allow the command to be deleted before the messages are deleted
	time.Sleep(1 * time.Second)

	deleted, err := utils.Discord.DeleteMessages(ctx.Session, channel.ID, user.ID, amount, delay)
	if err != nil {
		_ = ctx.RespondText("Error deleting messages: " + err.Error())
		return
	}

	c.services.Embed.NewEmbed().
		SetDescription(fmt.Sprintf("Successfully deleted %d messages!", deleted)).
		Send(ctx)
}

func (c *AdminCommands) purgerange(ctx *router.Ctx) {
	channel, success := ctx.Arguments.Get(2).Channel()
	channelID := ctx.Message.ChannelID
	if success {
		channelID = channel.ID
	}

	start, end := utils.Arguments["<start message>"], utils.Arguments["<end message>"]

	delay, err := strconv.Atoi(utils.Arguments["[delay]"])
	if err != nil {
		delay = 1
	}

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

		time.Sleep(time.Duration(delay) * time.Second)
	}
}

func (c *AdminCommands) nuke(ctx *router.Ctx) {
	channel, err := ctx.Arguments.Get(0).Channel(router.ArgumentFlagFallbackDefault)
	if err != nil {
		_ = ctx.RespondText("Could not find channel")
		return
	}

	_, err := ctx.Session.ChannelDelete(channel.ID)
	if err != nil {
		_ = ctx.RespondText("Error deleting channel: " + err.Error())
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
		_ = ctx.RespondText("Error creating channel: " + err.Error())
		return
	}
	_, _ = ctx.Session.ChannelEditComplex(channel.ID, &discordgo.ChannelEdit{
		Position: channel.Position,
	})

	ctx.Message.ChannelID = channel.ID

	c.services.Embed.NewEmbed().
		SetDescription("Successfully nuked channel!").
		Send(ctx)
}

func (c *AdminCommands) ackchannel(ctx *router.Ctx) {
	channel, err := ctx.Arguments.Get(0).Channel(router.ArgumentFlagFallbackDefault)
	if err != nil {
		_ = ctx.RespondText("Could not find channel")
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
		_ = ctx.RespondText("Could not find channel")
		return
	}
	_ = ctx.Session.MessageAck(channel.ID, ctx.Arguments.Get(1).Raw())
}

// Webhook Control

func (c *AdminCommands) listwebhooks(ctx *router.Ctx) {
	channel, success := utils.Discord.TryGetChannelFromMentionWithDefault(ctx, ctx.Arguments.Get(0))
	if !success {
		_ = ctx.RespondText("Could not find channel")
		return
	}

	webhooks, err := ctx.Session.ChannelWebhooks(channel.ID)
	if err != nil {
		_ = ctx.RespondText("Error getting webhooks: " + err.Error())
		return
	}

	if len(webhooks) == 0 {
		_ = ctx.RespondText("No webhooks found in channel")
		return
	}

	var fmtWebhooks []string

	for _, webhook := range webhooks {
		fmtWebhooks = append(
			fmtWebhooks,
			webhook.Name+" - "+utils.Discord.WebhookURL(webhook.ID, webhook.Token),
		)
	}

	c.services.Embed.NewEmbed().
		SetDescription("Webhooks in "+channel.Mention()).
		AddField("Webhooks", strings.Join(fmtWebhooks, "\n")).
		Send(ctx)
}

func (c *AdminCommands) createwebhook(ctx *router.Ctx) {
	channel, success := utils.Discord.TryGetChannelFromMentionWithDefault(ctx, ctx.Arguments.Get(0))
	if !success {
		_ = ctx.RespondText("Could not find channel")
		return
	}

	name := utils.Arguments["[name]"]
	if name == "" {
		name = "LightningBot Webhook"
	}

	webhook, err := ctx.Session.WebhookCreate(channel.ID, name, "")
	if err != nil {
		_ = ctx.RespondText("Error creating webhook: " + err.Error())
		return
	}

	c.services.Embed.NewEmbed().
		SetDescription("Successfully created webhook!").
		AddField("Webhook", webhook.Name+" - "+utils.Discord.WebhookURL(webhook.ID, webhook.Token)).
		Send(ctx)
}

func (c *AdminCommands) webhookinfo(ctx *router.Ctx) {
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

	c.services.Embed.NewEmbed().
		AddField("Type", map[int]string{1: "Incoming", 2: "Channel Follower", 3: "Application"}[webhookInfo.Get("type").MustInt()]).
		AddField("ID", webhookInfo.Get("id").MustString()).
		AddField("Name", webhookInfo.Get("name").MustString()).
		AddField("Avatar", webhookInfo.Get("avatar").MustString()).
		AddField("Channel ID", webhookInfo.Get("channel_id").MustString()).
		AddField("Server ID", webhookInfo.Get("guild_id").MustString()).
		AddField("Application ID", webhookInfo.Get("application_id").MustString()).
		AddField("Token", webhookInfo.Get("token").MustString()).
		AddField("URL", webhookUrl.String()).
		Send(ctx)
}

func (c *AdminCommands) webhooksend(ctx *router.Ctx) {
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

	httpResponse := http.Post(webhookUrl.String(), data)
	if httpResponse.Error != nil || httpResponse.StatusCode != 204 {
		_ = ctx.RespondText("Error sending message to webhook: " + string(httpResponse.Body))
		return
	}

	c.services.Embed.NewEmbed().
		SetDescription("Successfully sent message to webhook.").
		Send(ctx)
}

func (c *AdminCommands) webhookdelete(ctx *router.Ctx) {
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

	c.services.Embed.NewEmbed().
		SetDescription("Successfully deleted webhook.").
		Send(ctx)
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

	emb := c.services.Embed.NewEmbed().
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

func (c *AdminCommands) userage(ctx *router.Ctx) {
	user, err := ctx.Arguments.Get(0).User()
	if err != nil {
		_ = ctx.RespondText("Could not find user")
		return
	}
	age := time.Since(utils.Discord.SnowflakeToTime(user.ID))
	c.services.Embed.NewEmbed().
		SetDescription(fmt.Sprintf("User Age (%s - %s)", user.String(), user.ID)).
		AddField("Age", utils.Misc.FormatDuration(age)).
		Send(ctx)
}

func (c *AdminCommands) mutualfriends(ctx *router.Ctx) {
	user, err := ctx.Arguments.Get(0).User()
	if err != nil {
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
	c.services.Embed.NewEmbed().
		SetDescription(fmt.Sprintf("Mutual Friends (%s - %s)", user.String(), user.ID)).
		AddField("Mutual Friends", mutualFriendsString).
		Send(ctx)
}

func (c *AdminCommands) mutualservers(ctx *router.Ctx) {
	user, err := ctx.Arguments.Get(0).User()
	if err != nil {
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
	c.services.Embed.NewEmbed().
		SetDescription(fmt.Sprintf("Mutual Servers (%s - %s)", user.String(), user.ID)).
		AddField("Mutual Servers", mutualServersString).
		Send(ctx)
}

func (c *AdminCommands) platform(ctx *router.Ctx) {
	userIDArg := ctx.Arguments.Get(0)
	userID := utils.Misc.Or(userIDArg.AsUserMentionID(), userIDArg.Raw())
	if userID == "" {
		userID = ctx.Message.Author.ID
	}

	member, ok := utils.Discord.GetMember(ctx.Session, ctx.Message.GuildID, userID, true)
	if !ok {
		_ = ctx.RespondText("Could not find user!")
		return
	}

	if member.Presence == nil {
		_ = ctx.RespondText("Users presence is unknown!")
		return
	}

	c.services.Embed.NewEmbed().
		AddField("User", member.User.Username).
		AddField("Status", utils.Discord.FormatStatusEmoji(member.Presence.Status)).
		AddField("Platform", utils.Discord.FormatClientStatus(member.Presence.ClientStatus)).
		Send(ctx)
}

func (c *AdminCommands) avatar(ctx *router.Ctx) {
	user, err := ctx.Arguments.Get(0).User()
	if err != nil {
		_ = ctx.RespondText("Could not find user")
		return
	}
	if user.Avatar == "" {
		_ = ctx.RespondText("User has no avatar")
		return
	}
	_ = ctx.RespondText(user.AvatarURL(""))
}

func (c *AdminCommands) banner(ctx *router.Ctx) {
	user, err := ctx.Arguments.Get(0).User()
	if err != nil {
		_ = ctx.RespondText("Could not find user")
		return
	}
	if user.Banner == "" {
		_ = ctx.RespondText("User has no banner")
		return
	}
	_ = ctx.RespondText(user.BannerURL(""))
}

func (c *AdminCommands) mute(ctx *router.Ctx) {
	member, success := utils.Discord.TryGetMemberFromMention(ctx, ctx.Arguments.Get(0))
	if !success {
		_ = ctx.RespondText("Could not find user")
		return
	}

	if err := utils.Discord.MuteUser(ctx.Session, member); err != nil {
		_ = ctx.RespondText(utils.String.Capitalize(err.Error()))
		return
	}

	c.services.Embed.NewEmbed().
		SetDescription("Successfully muted user!").
		Send(ctx)
}

func (c *AdminCommands) unmute(ctx *router.Ctx) {
	member, success := utils.Discord.TryGetMemberFromMention(ctx, ctx.Arguments.Get(0))
	if !success {
		_ = ctx.RespondText("Could not find user")
		return
	}

	roles, err := ctx.Session.GuildRoles(ctx.Message.GuildID)
	if member == nil || err != nil {
		_ = ctx.RespondText("Error unmuting user: " + err.Error())
		return
	}

	for _, role := range roles {
		if strings.ToLower(role.Name) != "muted" {
			continue
		}

		if !slices.Contains(member.Roles, role.ID) {
			c.services.Embed.NewEmbed().
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

		_, err = ctx.Session.GuildMemberEdit(ctx.Message.GuildID, member.User.ID, &discordgo.GuildMemberParams{
			Roles: &member.Roles,
		})
		if err != nil {
			_ = ctx.RespondText("Error unmuting user: " + err.Error())
			return
		}

		c.services.Embed.NewEmbed().
			SetDescription("Successfully unmuted user!").
			Send(ctx)

		return
	}
	_ = ctx.RespondText("Error: No muted role found!")
}

func (c *AdminCommands) timeout(ctx *router.Ctx) {
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

	err = ctx.Session.GuildMemberTimeout(ctx.Message.GuildID, member.User.ID, &until)
	if err != nil {
		_ = ctx.RespondText("Error timing out user: " + err.Error())
		return
	}

	c.services.Embed.NewEmbed().
		SetDescription(fmt.Sprintf("Successfully timed out user for %d seconds!", seconds)).
		Send(ctx)
}

func (c *AdminCommands) removetimeout(ctx *router.Ctx) {
	member, success := utils.Discord.TryGetMemberFromMention(ctx, ctx.Arguments.Get(0))
	if !success {
		_ = ctx.RespondText("Could not find user")
		return
	}

	err := ctx.Session.GuildMemberTimeout(ctx.Message.GuildID, member.User.ID, nil)
	if err != nil {
		_ = ctx.RespondText("Error removing timeout: " + err.Error())
		return
	}

	c.services.Embed.NewEmbed().
		SetDescription("Successfully removed timeout!").
		Send(ctx)
}

func (c *AdminCommands) kick(ctx *router.Ctx) {
	member, success := utils.Discord.TryGetMemberFromMention(ctx, ctx.Arguments.Get(0))
	if !success {
		_ = ctx.RespondText("Could not find user")
		return
	}

	err := ctx.Session.GuildMemberDeleteWithReason(ctx.Message.GuildID, member.User.ID, utils.Arguments["[reason]"])
	if err != nil {
		_ = ctx.RespondText("Error kicking user: " + err.Error())
		return
	}

	c.services.Embed.NewEmbed().
		SetDescription("Successfully kicked user!").
		Send(ctx)
}

func (c *AdminCommands) ban(ctx *router.Ctx) {
	member, success := utils.Discord.TryGetMemberFromMention(ctx, ctx.Arguments.Get(0))
	memberID := ctx.Arguments.Get(0).Raw()
	if success {
		memberID = member.User.ID
	}

	err := ctx.Session.GuildBanCreateWithReason(ctx.Message.GuildID, memberID, utils.Arguments["[reason]"], 0)
	if err != nil {
		_ = ctx.RespondText("Error banning user: " + err.Error())
		return
	}

	c.services.Embed.NewEmbed().
		SetDescription("Successfully banned user!").
		Send(ctx)
}

func (c *AdminCommands) softban(ctx *router.Ctx) {
	member, success := utils.Discord.TryGetMemberFromMention(ctx, ctx.Arguments.Get(0))
	if !success {
		_ = ctx.RespondText("Could not find user")
		return
	}

	err := ctx.Session.GuildBanCreateWithReason(ctx.Message.GuildID, member.User.ID, utils.Arguments["[reason]"], 0)
	if err != nil {
		_ = ctx.RespondText("Error softbanning user: " + err.Error())
		return
	}
	err = ctx.Session.GuildBanDelete(ctx.Message.GuildID, member.User.ID)
	if err != nil {
		_ = ctx.RespondText("Error softbanning user: " + err.Error())
		return
	}

	c.services.Embed.NewEmbed().
		SetDescription("Successfully softbanned user!").
		Send(ctx)
}

func (c *AdminCommands) unban(ctx *router.Ctx) {
	err := ctx.Session.GuildBanDelete(ctx.Message.GuildID, utils.Arguments["<user id>"])
	if err != nil {
		_ = ctx.RespondText("Error unbanning user: " + err.Error())
		return
	}

	c.services.Embed.NewEmbed().
		SetDescription("Successfully unbanned user!").
		Send(ctx)
}

func (c *AdminCommands) banlist(ctx *router.Ctx) {
	guildID := utils.Arguments["[server id]"]
	if guildID == "" {
		guildID = ctx.Message.GuildID
	}
	if guildID == "" {
		_ = ctx.RespondText("Error: No server ID specified!")
		return
	}
	banList, err := ctx.Session.GuildBans(ctx.Message.GuildID, 1000, "", "")
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
	c.services.Embed.NewEmbed().
		AddField("Banned Amount", fmt.Sprint(len(banList))).
		AddField("Ban List", banListString).
		Send(ctx)
}

func (c *AdminCommands) copybans(ctx *router.Ctx) {
	sourceServerID := utils.Arguments["<source server id>"]
	destinationServerID := utils.Arguments["[dest server id]"]
	if destinationServerID == "" {
		destinationServerID = ctx.Message.GuildID
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

func (c *AdminCommands) nick(ctx *router.Ctx) {
	err := ctx.Session.GuildMemberNickname(ctx.Message.GuildID, "@me", ctx.Arguments.Raw())
	if err != nil {
		_ = ctx.RespondText("Error changing nickname: " + err.Error())
		return
	}
	c.services.Embed.NewEmbed().
		SetDescription("Successfully changed nickname!").
		Send(ctx)
}

func (c *AdminCommands) nickuser(ctx *router.Ctx) {
	member, success := utils.Discord.TryGetMemberFromMention(ctx, ctx.Arguments.Get(0))
	if !success {
		_ = ctx.RespondText("Could not find user")
		return
	}

	err := ctx.Session.GuildMemberNickname(ctx.Message.GuildID, strings.ReplaceAll(member.User.ID, ctx.Message.Author.ID, "@me"), ctx.Arguments.RawFrom(1))
	if err != nil {
		_ = ctx.RespondText("Error changing nickname: " + err.Error())
		return
	}
	c.services.Embed.NewEmbed().
		SetDescription("Successfully changed nickname!").
		Send(ctx)
}

func (c *AdminCommands) clearnick(ctx *router.Ctx) {
	userID := "@me"
	member, success := utils.Discord.TryGetMemberFromMention(ctx, ctx.Arguments.Get(0))
	if success {
		userID = member.User.ID
	}
	err := ctx.Session.GuildMemberNickname(ctx.Message.GuildID, userID, "")
	if err != nil {
		_ = ctx.RespondText("Error changing nickname: " + err.Error())
		return
	}
	c.services.Embed.NewEmbed().
		SetDescription("Successfully changed nickname!").
		Send(ctx)
}

// Server Control

func (c *AdminCommands) serverinfo(ctx *router.Ctx) {
	guild, err := utils.Discord.GetGuild(ctx.Session, utils.Misc.Or(utils.Arguments["[server id]"], ctx.Message.GuildID))
	if err != nil {
		_ = ctx.RespondText("Error getting server: " + err.Error())
		return
	}

	owner, ok := utils.Discord.GetMember(ctx.Session, guild.ID, guild.OwnerID)
	ownerMention := "Error"
	if ok {
		ownerMention = owner.User.String()
	}

	c.services.Embed.NewEmbed().
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

func (c *AdminCommands) serverage(ctx *router.Ctx) {
	guild, err := utils.Discord.GetGuild(ctx.Session, utils.Misc.Or(utils.Arguments["[server id]"], ctx.Message.GuildID))
	if err != nil {
		_ = ctx.RespondText("Error getting server: " + err.Error())
		return
	}
	age := time.Since(utils.Discord.SnowflakeToTime(guild.ID))
	c.services.Embed.NewEmbed().
		SetDescription(fmt.Sprintf("Server Age (%s - %s)", guild.Name, guild.ID)).
		AddField("Age", utils.Misc.FormatDuration(age)).
		Send(ctx)
}

func (c *AdminCommands) serverowner(ctx *router.Ctx) {
	guild, err := utils.Discord.GetGuild(ctx.Session, utils.Misc.Or(utils.Arguments["[server id]"], ctx.Message.GuildID))
	if err != nil {
		_ = ctx.RespondText("Error getting server: " + err.Error())
		return
	}
	owner, ok := utils.Discord.GetMember(ctx.Session, guild.ID, guild.OwnerID)
	if !ok {
		_ = ctx.RespondText("Error getting server owner")
		return
	}
	c.services.Embed.NewEmbed().
		SetDescription(guild.Name).
		AddField("Owner", owner.Mention()).
		Send(ctx)
}

func (c *AdminCommands) servericon(ctx *router.Ctx) {
	guild, err := utils.Discord.GetGuild(ctx.Session, utils.Misc.Or(utils.Arguments["[server id]"], ctx.Message.GuildID))
	if err != nil {
		_ = ctx.RespondText("Error getting server: " + err.Error())
		return
	}

	emb := c.services.Embed.NewEmbed().
		SetDescription(fmt.Sprintf("Server Icon - %s", guild.Name)).
		SetImage(guild.IconURL())
	if guild.Icon == "" {
		emb.AddField("Error", "Server has no icon")
	}
	emb.Send(ctx)
}

func (c *AdminCommands) serverbanner(ctx *router.Ctx) {
	guild, err := utils.Discord.GetGuild(ctx.Session, utils.Misc.Or(utils.Arguments["[server id]"], ctx.Message.GuildID))
	if err != nil {
		_ = ctx.RespondText("Error getting server: " + err.Error())
		return
	}
	if guild.Banner == "" {
		c.services.Embed.NewEmbed().AddField("Error", "Server has no banner!").Send(ctx)
	}
	c.services.Embed.NewEmbed().
		SetDescription(fmt.Sprintf("Server Banner - %s", guild.Name)).
		SetImage(guild.BannerURL()).
		Send(ctx)
}

func (c *AdminCommands) serversplash(ctx *router.Ctx) {
	guild, err := utils.Discord.GetGuild(ctx.Session, utils.Misc.Or(utils.Arguments["[server id]"], ctx.Message.GuildID))
	if err != nil {
		_ = ctx.RespondText("Error getting server: " + err.Error())
		return
	}
	if guild.Splash == "" {
		_ = ctx.RespondText("Server has no invite splash!")
		return
	}
	c.services.Embed.NewEmbed().
		SetDescription(fmt.Sprintf("Server Invite Splash - %s", guild.Name)).
		SetImage(discordgo.EndpointGuildSplash(guild.ID, guild.Splash)).
		Send(ctx)
}

func (c *AdminCommands) serverdsplash(ctx *router.Ctx) {
	guild, err := utils.Discord.GetGuild(ctx.Session, utils.Misc.Or(utils.Arguments["[server id]"], ctx.Message.GuildID))
	if err != nil {
		_ = ctx.RespondText("Error getting server: " + err.Error())
		return
	}
	if guild.DiscoverySplash == "" {
		_ = ctx.RespondText("Server has no discovery splash!")
		return
	}
	c.services.Embed.NewEmbed().
		SetDescription(fmt.Sprintf("Server Discovery Splash - %s", guild.Name)).
		SetImage(discordgo.EndpointGuildSplash(guild.ID, guild.DiscoverySplash)).
		Send(ctx)
}

func (c *AdminCommands) serverboosters(ctx *router.Ctx) {
	guild, err := utils.Discord.GetGuild(ctx.Session, utils.Misc.Or(utils.Arguments["[server id]"], ctx.Message.GuildID))
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

	c.services.Embed.NewEmbed().
		SetDescription(fmt.Sprintf("Server Boosters - %s", guild.Name)).
		AddField("Boosters", boostersString).
		Send(ctx)
}

func (c *AdminCommands) serverboostlevel(ctx *router.Ctx) {
	guild, err := utils.Discord.GetGuild(ctx.Session, utils.Misc.Or(utils.Arguments["[server id]"], ctx.Message.GuildID))
	if err != nil {
		_ = ctx.RespondText("Error getting server: " + err.Error())
		return
	}

	c.services.Embed.NewEmbed().
		SetDescription(fmt.Sprintf("Server Boost Level - %s", guild.Name)).
		AddField("Boost Level", boostLevels[int(guild.PremiumTier)]).
		Send(ctx)
}

func (c *AdminCommands) setservername(ctx *router.Ctx) {
	guild, err := ctx.Session.GuildEdit(ctx.Message.GuildID, &discordgo.GuildParams{
		Name: ctx.Arguments.Raw(),
	})
	if err != nil {
		_ = ctx.RespondText("Error changing server name: " + err.Error())
		return
	}
	c.services.Embed.NewEmbed().
		SetDescription(fmt.Sprintf("Successfully changed server name to %s!", guild.Name)).
		Send(ctx)
}

func (c *AdminCommands) setservericon(ctx *router.Ctx) {
	_, err := ctx.Session.GuildEdit(ctx.Message.GuildID, &discordgo.GuildParams{
		Icon: utils.Arguments["<url>"],
	})
	if err != nil {
		_ = ctx.RespondText("Error changing server icon: " + err.Error())
		return
	}
	c.services.Embed.NewEmbed().
		SetDescription("Successfully changed server icon!").
		Send(ctx)
}

func (c *AdminCommands) setserverbanner(ctx *router.Ctx) {
	_, err := ctx.Session.GuildEdit(ctx.Message.GuildID, &discordgo.GuildParams{
		Banner: utils.Arguments["<url>"],
	})
	if err != nil {
		_ = ctx.RespondText("Error changing server icon: " + err.Error())
		return
	}
	c.services.Embed.NewEmbed().
		SetDescription("Successfully changed server icon!").
		Send(ctx)
}

func (c *AdminCommands) lockserver(ctx *router.Ctx) {
	guild, err := utils.Discord.GetGuild(ctx.Session, utils.Misc.Or(utils.Arguments["[server id]"], ctx.Message.GuildID))
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
			Position:             channel.Position,
			PermissionOverwrites: overwrites,
		})

		if err != nil {
			_ = ctx.RespondText("Error locking server: " + err.Error())
			return
		}
	}

	c.services.Embed.NewEmbed().
		SetDescription("Successfully locked server!").
		Send(ctx)
}

func (c *AdminCommands) unlockserver(ctx *router.Ctx) {
	guild, err := utils.Discord.GetGuild(ctx.Session, utils.Misc.Or(utils.Arguments["[server id]"], ctx.Message.GuildID))
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
			Position:             channel.Position,
			PermissionOverwrites: overwrites,
		})

		if err != nil {
			_ = ctx.RespondText("Error unlocking server: " + err.Error())
			return
		}
	}

	c.services.Embed.NewEmbed().
		SetDescription("Successfully unlocked server!").
		Send(ctx)
}

func (c *AdminCommands) channels(ctx *router.Ctx) {
	guild, err := utils.Discord.GetGuild(ctx.Session, utils.Misc.Or(utils.Arguments["[server id]"], ctx.Message.GuildID))
	if err != nil {
		_ = ctx.RespondText("Error getting server: " + err.Error())
		return
	}

	categories, categoryChannels, looseChannels := utils.Discord.MapCategories(guild.Channels)
	emb := c.services.Embed.NewEmbed().SetDescription(fmt.Sprintf("Server Channels - %s", guild.Name))

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

func (c *AdminCommands) bots(ctx *router.Ctx) {
	guild, err := utils.Discord.GetGuild(ctx.Session, utils.Misc.Or(utils.Arguments["[server id]"], ctx.Message.GuildID))
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
	c.services.Embed.NewEmbed().
		SetDescription(fmt.Sprintf("Server Bots - %s", guild.Name)).
		AddField("Server Bots", botString).
		Send(ctx)
}

func (c *AdminCommands) admins(ctx *router.Ctx) {
	guild, err := utils.Discord.GetGuild(ctx.Session, utils.Misc.Or(utils.Arguments["[server id]"], ctx.Message.GuildID))
	if err != nil {
		_ = ctx.RespondText("Error getting server: " + err.Error())
		return
	}
	if len(guild.Channels) == 0 {
		_ = ctx.RespondText("Server has no channels!")
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
	c.services.Embed.NewEmbed().
		SetDescription(fmt.Sprintf("Server Admins - %s", guild.Name)).
		AddField("Server Admins", adminString).
		Send(ctx)
}

func (c *AdminCommands) roles(ctx *router.Ctx) {
	guild, err := utils.Discord.GetGuild(ctx.Session, utils.Misc.Or(utils.Arguments["[server id]"], ctx.Message.GuildID))
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
	c.services.Embed.NewEmbed().
		SetDescription(fmt.Sprintf("Server Roles - %s", guild.Name)).
		AddField("Server Roles", roleString).
		Send(ctx)
}

func (c *AdminCommands) leaveserver(ctx *router.Ctx) {
	guild, err := utils.Discord.GetGuild(ctx.Session, utils.Misc.Or(utils.Arguments["[server id]"], ctx.Message.GuildID))
	if err != nil {
		_ = ctx.RespondText("Error getting server: " + err.Error())
		return
	}
	_ = ctx.Session.GuildLeave(guild.ID)
}

func (c *AdminCommands) deleteserver(ctx *router.Ctx) {
	guild, err := utils.Discord.GetGuild(ctx.Session, utils.Misc.Or(utils.Arguments["[server id]"], ctx.Message.GuildID))
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

func (c *AdminCommands) invite(ctx *router.Ctx) {
	channel, err := ctx.Arguments.Get(0).Channel(router.ArgumentFlagFallbackDefault)
	if err != nil {
		_ = ctx.RespondText("Could not find channel")
		return
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

	invite, err := ctx.Session.ChannelInviteCreate(channel.ID, ageInt, usesInt, false)
	if err != nil {
		_ = ctx.RespondText("Error creating invite: " + err.Error())
		return
	}

	c.services.Embed.NewEmbed().
		SetDescription(fmt.Sprintf("Created invite: https://discord.gg/%s", invite.Code)).
		Send(ctx)
}

func (c *AdminCommands) friendinvite(ctx *router.Ctx) {
	invite, err := ctx.Session.CreateFriendInvite("@me")
	if err != nil {
		_ = ctx.RespondText("Error creating friend invite: " + err.Error())
		return
	}

	c.services.Embed.NewEmbed().
		SetDescription("Successfully created friend invite!").
		AddField("Invite", fmt.Sprintf("https://discord.gg/%s", invite.Code)).
		Send(ctx)
}

func (c *AdminCommands) inviteinfo(ctx *router.Ctx) {
	parsedInvite := utils.Discord.ParseInvite(utils.Arguments["<invite>"])
	invite, err := ctx.Session.InviteWithCounts(parsedInvite)
	if err != nil {
		_ = ctx.RespondText("Error getting invite: " + err.Error())
		return
	}

	emb := c.services.Embed.NewEmbed().SetDescription(fmt.Sprintf("Invite Info - %s", invite.Code))
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

func (c *AdminCommands) deleteinvite(ctx *router.Ctx) {
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

	c.services.Embed.NewEmbed().
		SetDescription(fmt.Sprintf("Deleted invite: %s", invite.Code)).
		Send(ctx)
}

func (c *AdminCommands) serverinvites(ctx *router.Ctx) {
	invites, err := ctx.Session.GuildInvites(utils.Misc.Or(utils.Arguments["[server id]"], ctx.Message.GuildID))
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
	c.services.Embed.NewEmbed().
		SetDescription("Invite List").
		AddField("Invites", invitesString).
		Send(ctx)
}

func (c *AdminCommands) channelinvites(ctx *router.Ctx) {
	channel, err := ctx.Arguments.Get(0).Channel(router.ArgumentFlagFallbackDefault)
	if err != nil {
		_ = ctx.RespondText("Could not find channel")
		return
	}
	invites, err := ctx.Session.ChannelInvites(channel.ID)
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
	c.services.Embed.NewEmbed().
		SetDescription("Invite List").
		AddField("Invites", invitesString).
		Send(ctx)
}

func (c *AdminCommands) userinvites(ctx *router.Ctx) {
	member, success := utils.Discord.TryGetMemberFromMention(ctx, ctx.Arguments.Get(0))
	if !success {
		_ = ctx.RespondText("Could not find user")
		return
	}

	invites, err := ctx.Session.GuildInvites(utils.Misc.Or(utils.Arguments["[server id]"], ctx.Message.GuildID))
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
	c.services.Embed.NewEmbed().
		SetDescription(fmt.Sprintf("Invite List (%s)", member.User.String())).
		AddField("Invites", invitesString).
		Send(ctx)
}

func (c *AdminCommands) roleinfo(ctx *router.Ctx) {
	role, success := utils.Discord.TryGetRoleFromMention(ctx, ctx.Arguments.Get(0))
	if !success {
		_ = ctx.RespondText("Could not find role")
		return
	}
	// Get all members in guild and check if they have the role
	guild, err := utils.Discord.GetGuild(ctx.Session, ctx.Message.GuildID)
	if err != nil {
		_ = ctx.RespondText("Error getting server: " + err.Error())
		return
	}
	memberCount := 0
	for _, member := range guild.Members {
		if slices.Contains(member.Roles, role.ID) {
			memberCount++
		}
	}
	createdAt := utils.Discord.SnowflakeToTime(role.ID)
	perms, _ := api.CreatePaste(utils.Discord.FormatPermissions(role.Name, role.ID, role.Permissions))
	c.services.Embed.NewEmbed().
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

func (c *AdminCommands) giverole(ctx *router.Ctx) {
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

	err := ctx.Session.GuildMemberRoleAdd(ctx.Message.GuildID, member.User.ID, role.ID)
	if err != nil {
		_ = ctx.RespondText("Error giving role: " + err.Error())
		return
	}
	c.services.Embed.NewEmbed().
		SetDescription(fmt.Sprintf("Successfully gave role (%s) to (%s)", role.Name, member.User.String())).
		Send(ctx)
}

func (c *AdminCommands) removerole(ctx *router.Ctx) {
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

	err := ctx.Session.GuildMemberRoleRemove(ctx.Message.GuildID, member.User.ID, role.ID)
	if err != nil {
		_ = ctx.RespondText("Error removing role: " + err.Error())
		return
	}
	c.services.Embed.NewEmbed().
		SetDescription(fmt.Sprintf("Successfully removed role (%s) from (%s)", role.Name, member.User.String())).
		Send(ctx)
}

func (c *AdminCommands) createrole(ctx *router.Ctx) {
	role, err := ctx.Session.GuildRoleCreate(ctx.Message.GuildID, &discordgo.RoleParams{
		Name: ctx.Arguments.Raw(),
	})
	if err != nil {
		_ = ctx.RespondText("Error creating role: " + err.Error())
		return
	}

	c.services.Embed.NewEmbed().
		SetDescription(fmt.Sprintf("Successfully created role (%s)", role.Name)).
		Send(ctx)
}

func (c *AdminCommands) deleterole(ctx *router.Ctx) {
	role, success := utils.Discord.TryGetRoleFromMention(ctx, ctx.Arguments.Get(0))
	if !success {
		_ = ctx.RespondText("Could not find role")
		return
	}
	err := ctx.Session.GuildRoleDelete(ctx.Message.GuildID, role.ID)
	if err != nil {
		_ = ctx.RespondText("Error deleting role: " + err.Error())
		return
	}
	c.services.Embed.NewEmbed().
		SetDescription(fmt.Sprintf("Successfully deleted role (%s)", role.Name)).
		Send(ctx)
}

func (c *AdminCommands) renamerole(ctx *router.Ctx) {
	role, success := utils.Discord.TryGetRoleFromMention(ctx, ctx.Arguments.Get(0))
	if !success {
		_ = ctx.RespondText("Could not find role")
		return
	}

	role, err := ctx.Session.GuildRoleEdit(ctx.Message.GuildID, role.ID, &discordgo.RoleParams{
		Name: ctx.Arguments.Raw(),
	})
	if err != nil {
		_ = ctx.RespondText("Error renaming role: " + err.Error())
		return
	}

	c.services.Embed.NewEmbed().
		SetDescription(fmt.Sprintf("Successfully renamed role (%s)", role.Name)).
		Send(ctx)
}
