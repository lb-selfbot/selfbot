package admin

import (
	"fmt"
	"net/url"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/lb-selfbot/discordgo"
	"github.com/lb-selfbot/selfbot/internal/formatting"
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

func (c *AdminCommands) channelinfo(ctx *router.Ctx) error {
	channel, err := ctx.Arguments.Get(0).Channel()
	if err != nil {
		return fmt.Errorf("Could not find channel")
	}

	c.services.Embed.New(ctx).
		AddField("Name", channel.Name).
		AddField("ID", channel.ID).
		AddField("Topic", channel.Topic).
		AddField("Category ID", channel.ParentID).
		AddField("Is NSFW", formatting.BoolYesNo(channel.NSFW)).
		AddField("Created At", discordgo.MustSnowflakeTimestamp(channel.ID).Format(time.RFC1123)).
		AddField("Slowmode Delay", fmt.Sprint(channel.RateLimitPerUser)+" seconds").
		AddField("Position", fmt.Sprint(channel.Position)).
		AddField("Last Message ID", channel.LastMessageID).
		Send()
	return nil
}

func (c *AdminCommands) textchannel(ctx *router.Ctx) error {
	channel, err := ctx.Session.GuildChannelCreate(ctx.Message.GuildID, ctx.Arguments.Raw(), discordgo.ChannelTypeGuildText)
	if err != nil {
		return fmt.Errorf("Error creating channel: %v", err)
	}
	c.services.Embed.New(ctx).
		SetDescription("Successfully created text channel!").
		AddField("Channel", channel.Mention()).
		Send()
	return nil
}

func (c *AdminCommands) voicechannel(ctx *router.Ctx) error {
	channel, err := ctx.Session.GuildChannelCreate(ctx.Message.GuildID, ctx.Arguments.Raw(), discordgo.ChannelTypeGuildVoice)
	if err != nil {
		return fmt.Errorf("Error creating channel: %v", err)
	}
	c.services.Embed.New(ctx).
		SetDescription("Successfully created voice channel!").
		AddField("Channel", channel.Mention()).
		Send()
	return nil
}

func (c *AdminCommands) stagechannel(ctx *router.Ctx) error {
	channel, err := ctx.Session.GuildChannelCreate(ctx.Message.GuildID, ctx.Arguments.Raw(), discordgo.ChannelTypeGuildStageVoice)
	if err != nil {
		return fmt.Errorf("Error creating channel: %v", err)
	}
	c.services.Embed.New(ctx).
		SetDescription("Successfully created stage channel!").
		AddField("Channel", channel.Mention()).
		Send()
	return nil
}

func (c *AdminCommands) newschannel(ctx *router.Ctx) error {
	channel, err := ctx.Session.GuildChannelCreate(ctx.Message.GuildID, ctx.Arguments.Raw(), discordgo.ChannelTypeGuildNews)
	if err != nil {
		return fmt.Errorf("Error creating channel: %v", err)
	}
	c.services.Embed.New(ctx).
		SetDescription("Successfully created news channel!").
		AddField("Channel", channel.Mention()).
		Send()
	return nil
}

func (c *AdminCommands) thread(ctx *router.Ctx) error {
	channel, err := ctx.Session.ThreadStart(ctx.Message.ChannelID, ctx.Arguments.Raw(), discordgo.ChannelTypeGuildPublicThread, 0)
	if err != nil {
		return fmt.Errorf("Error creating thread: %v", err)
	}
	c.services.Embed.New(ctx).
		SetDescription("Successfully created thread!").
		AddField("Thread", channel.Mention()).
		Send()
	return nil
}

func (c *AdminCommands) category(ctx *router.Ctx) error {
	channel, err := ctx.Session.GuildChannelCreate(ctx.Message.GuildID, ctx.Arguments.Raw(), discordgo.ChannelTypeGuildCategory)
	if err != nil {
		return fmt.Errorf("Error creating category: %v", err)
	}
	c.services.Embed.New(ctx).
		SetDescription("Successfully created category!").
		AddField("Category", channel.Mention()).
		Send()
	return nil
}

func (c *AdminCommands) deletecategory(ctx *router.Ctx) error {
	channel, err := ctx.Arguments.Get(0).Channel()
	if err != nil {
		return fmt.Errorf("Could not find category")
	}
	_, err = ctx.Session.ChannelDelete(channel.ID)
	if err != nil {
		return fmt.Errorf("Error deleting category: %v", err)
	}

	c.services.Embed.Respond(ctx, "Successfully deleted category!")
	return nil
}

func (c *AdminCommands) renamechannel(ctx *router.Ctx) error {
	channel, err := ctx.Arguments.Get(0).Channel()
	if err != nil {
		return fmt.Errorf("Could not find channel")
	}

	_, err = ctx.Session.ChannelEdit(channel.ID, &discordgo.ChannelEdit{
		Name: ctx.Arguments.RawFrom(1),
	})
	if err != nil {
		return fmt.Errorf("Error renaming channel: %v", err)
	}

	c.services.Embed.New(ctx).
		SetDescription("Successfully renamed channel!").
		AddField("Channel", channel.Mention()).
		Send()
	return nil
}

func (c *AdminCommands) deletechannel(ctx *router.Ctx) error {
	channel, err := ctx.Arguments.Get(0).Channel()
	if err != nil {
		return fmt.Errorf("Could not find channel")
	}

	_, err = ctx.Session.ChannelDelete(channel.ID)
	if err != nil {
		return fmt.Errorf("Error deleting channel: %v", err)
	}

	c.services.Embed.Respond(ctx, "Successfully deleted channel!")
	return nil
}

func (c *AdminCommands) slowmode(ctx *router.Ctx) error {
	channel, err := ctx.Arguments.Get(0).Channel()
	if err != nil {
		return fmt.Errorf("Could not find channel")
	}

	duration, err := ctx.Arguments.Get(1).Duration()
	if err != nil {
		return fmt.Errorf("Could not parse duration")
	}

	durationSec := int(duration.Seconds())

	_, err = ctx.Session.ChannelEditComplex(channel.ID, &discordgo.ChannelEdit{
		Position:         channel.Position,
		RateLimitPerUser: &durationSec,
	})
	if err != nil {
		return fmt.Errorf("Error setting slowmode: %v", err)
	}

	c.services.Embed.New(ctx).
		SetDescription("Successfully set slowmode to %s", duration.String()).
		AddField("Channel", channel.Mention()).
		Send()
	return nil
}

func (c *AdminCommands) removeslowmode(ctx *router.Ctx) error {
	channel, err := ctx.Arguments.Get(0).Channel()
	if err != nil {
		return fmt.Errorf("Could not find channel")
	}

	_, err = ctx.Session.ChannelEditComplex(channel.ID, &discordgo.ChannelEdit{
		Position:         channel.Position,
		RateLimitPerUser: nil,
	})
	if err != nil {
		return fmt.Errorf("Error removing slowmode: %v", err)
	}

	c.services.Embed.New(ctx).
		SetDescription("Successfully removed slowmode!").
		AddField("Channel", channel.Mention()).
		Send()
	return nil
}

func (c *AdminCommands) lockchannel(ctx *router.Ctx) error {
	channel, err := ctx.Arguments.Get(0).Channel()
	if err != nil {
		return fmt.Errorf("Could not find channel")
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
		return fmt.Errorf("Error locking channel: %v", err)
	}

	c.services.Embed.New(ctx).
		SetDescription("Successfully locked channel!").
		AddField("Channel", channel.Mention()).
		Send()
	return nil
}

func (c *AdminCommands) unlockchannel(ctx *router.Ctx) error {
	channel, err := ctx.Arguments.Get(0).Channel()
	if err != nil {
		return fmt.Errorf("Could not find channel")
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
		return fmt.Errorf("Error unlocking channel: %v", err)
	}

	c.services.Embed.New(ctx).
		SetDescription("Successfully unlocked channel!").
		AddField("Channel", channel.Mention()).
		Send()
	return nil
}

func (c *AdminCommands) purge(ctx *router.Ctx) error {
	amount, err := ctx.Arguments.Get(0).Int()
	amount = fn.Clamp(amount, 1, 10000)
	if err != nil {
		return fmt.Errorf("Argument <amount> must be a number.")
	}

	delay, err := ctx.Arguments.Get(1).Duration()
	delay = fn.Clamp(delay, 1*time.Second, 60*time.Second)
	if err != nil {
		return fmt.Errorf("Could not parse delay")
	}

	channel, err := ctx.Arguments.Get(2).Channel(router.ArgumentFlagFallbackDefault)
	if err != nil {
		return fmt.Errorf("Could not find channel")
	}

	messages, err := utils.Discord.GetChannelMessages(ctx.Session, channel.ID, amount)
	if err != nil {
		return fmt.Errorf("Error getting messages: %v", err)
	}

	for _, message := range messages {
		err = ctx.Session.ChannelMessageDelete(channel.ID, message.ID)
		if err != nil {
			return fmt.Errorf("Error deleting messages: %v", err)
		}

		time.Sleep(delay)
	}

	return nil
}

func (c *AdminCommands) purgeself(ctx *router.Ctx) error {
	amount, err := ctx.Arguments.Get(0).Int()
	amount = fn.Clamp(amount, 1, 10000)
	if err != nil {
		return fmt.Errorf("Argument <amount> must be a number.")
	}

	delay, err := ctx.Arguments.Get(1).Duration()
	delay = fn.Clamp(delay, 1*time.Second, 60*time.Second)
	if err != nil {
		return fmt.Errorf("Could not parse delay")
	}

	channel, err := ctx.Arguments.Get(2).Channel(router.ArgumentFlagFallbackDefault)
	if err != nil {
		return fmt.Errorf("Could not find channel")
	}

	userID := ctx.Message.Author.ID

	// Allow the command to be deleted before the messages are deleted
	time.Sleep(1 * time.Second)

	deleted, err := utils.Discord.DeleteMessages(ctx.Session, channel.ID, userID, amount, delay)
	if err != nil {
		return fmt.Errorf("Error deleting messages: %v", err)
	}

	c.services.Embed.Respond(ctx, "Successfully deleted %d messages!", deleted)
	return nil
}

func (c *AdminCommands) purgeuser(ctx *router.Ctx) error {
	amount, err := ctx.Arguments.Get(0).Int()
	amount = fn.Clamp(amount, 1, 10000)
	if err != nil {
		return fmt.Errorf("Argument <amount> must be a number.")
	}

	user, err := ctx.Arguments.Get(1).User()
	if err != nil {
		return fmt.Errorf("Could not find user")
	}

	delay, _ := ctx.Arguments.Get(2).Duration()
	delay = fn.Clamp(delay, 1*time.Second, 60*time.Second)

	channel, err := ctx.Arguments.Get(3).Channel(router.ArgumentFlagFallbackDefault)
	if err != nil {
		return fmt.Errorf("Could not find channel")
	}

	// Allow the command to be deleted before the messages are deleted
	time.Sleep(1 * time.Second)

	deleted, err := utils.Discord.DeleteMessages(ctx.Session, channel.ID, user.ID, amount, delay)
	if err != nil {
		return fmt.Errorf("Error deleting messages: %v", err)
	}

	c.services.Embed.Respond(ctx, "Successfully deleted %d messages!", deleted)
	return nil
}

func (c *AdminCommands) purgerange(ctx *router.Ctx) error {
	start, end := ctx.Arguments.Get(0).Raw(), ctx.Arguments.Get(1).Raw()

	delay, _ := ctx.Arguments.Get(2).Duration()
	delay = fn.Clamp(delay, 1*time.Second, 60*time.Second)

	channel, err := ctx.Arguments.Get(3).Channel(router.ArgumentFlagFallbackDefault)
	if err != nil {
		return fmt.Errorf("Could not find channel")
	}

	messages, err := utils.Discord.GetChannelMessagesRange(ctx.Session, channel.ID, start, end)
	if err != nil {
		return fmt.Errorf("Error getting messages: %v", err)
	}

	for _, message := range messages {
		if message.ID < start || message.ID > end {
			continue
		}

		err = ctx.Session.ChannelMessageDelete(channel.ID, message.ID)
		if err != nil {
			return fmt.Errorf("Error deleting messages: %v", err)
		}

		time.Sleep(delay)
	}

	return nil
}

func (c *AdminCommands) nuke(ctx *router.Ctx) error {
	channel, err := ctx.Arguments.Get(0).Channel(router.ArgumentFlagFallbackDefault)
	if err != nil {
		return fmt.Errorf("Could not find channel")
	}

	_, err = ctx.Session.ChannelDelete(channel.ID)
	if err != nil {
		return fmt.Errorf("Error deleting channel: %v", err)
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
		return fmt.Errorf("Error creating channel: %v", err)
	}
	_, _ = ctx.Session.ChannelEditComplex(channel.ID, &discordgo.ChannelEdit{
		Position: channel.Position,
	})

	ctx.Message.ChannelID = channel.ID

	c.services.Embed.Respond(ctx, "Successfully nuked channel!")
	return nil
}

func (c *AdminCommands) ackchannel(ctx *router.Ctx) error {
	channel, err := ctx.Arguments.Get(0).Channel(router.ArgumentFlagFallbackDefault)
	if err != nil {
		return fmt.Errorf("Could not find channel")
	}
	messages, _ := utils.Discord.GetChannelMessages(ctx.Session, channel.ID, 50)
	if len(messages) > 0 {
		return ctx.Session.MessageAck(channel.ID, messages[0].ID)
	}

	return nil
}

func (c *AdminCommands) ackmessage(ctx *router.Ctx) error {
	channel, err := ctx.Arguments.Get(0).Channel()
	if err != nil {
		return fmt.Errorf("Could not find channel")
	}
	return ctx.Session.MessageAck(channel.ID, ctx.Arguments.Get(1).Raw())
}

// Webhook Control

func (c *AdminCommands) listwebhooks(ctx *router.Ctx) error {
	channel, err := ctx.Arguments.Get(0).Channel(router.ArgumentFlagFallbackDefault)
	if err != nil {
		return fmt.Errorf("Could not find channel")
	}

	webhooks, err := ctx.Session.ChannelWebhooks(channel.ID)
	if err != nil {
		return fmt.Errorf("Error getting webhooks: %v", err)
	}

	if len(webhooks) == 0 {
		return fmt.Errorf("No webhooks found in channel")
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
	return nil
}

func (c *AdminCommands) createwebhook(ctx *router.Ctx) error {
	channel, err := ctx.Arguments.Get(0).Channel(router.ArgumentFlagFallbackDefault)
	if err != nil {
		return fmt.Errorf("Could not find channel")
	}

	name := ctx.Arguments.Get(1).Raw()
	if name == "" {
		name = "LightningBot Webhook"
	}

	webhook, err := ctx.Session.WebhookCreate(channel.ID, name, "")
	if err != nil {
		return fmt.Errorf("Error creating webhook: %v", err)
	}

	c.services.Embed.New(ctx).
		SetDescription("Successfully created webhook!").
		AddField("Webhook", webhook.Name+" - "+discordgo.EndpointWebhookToken(webhook.ID, webhook.Token)).
		Send()
	return nil
}

func (c *AdminCommands) webhookinfo(ctx *router.Ctx) error {
	webhookUrl, err := url.Parse(ctx.Arguments.Get(0).Raw())
	if err != nil {
		return fmt.Errorf("Invalid webhook URL!")
	}

	if !strings.Contains(webhookUrl.Host, "discord.com") {
		return fmt.Errorf("Invalid webhook URL!")
	}

	resp, err := http.DefaultClient().Get(webhookUrl.String())
	if err != nil || resp.StatusCode != 200 {
		return fmt.Errorf("Invalid webhook URL!")
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
	return nil
}

func (c *AdminCommands) webhooksend(ctx *router.Ctx) error {
	webhookUrl, err := url.Parse(ctx.Arguments.Get(0).Raw())
	if err != nil {
		return fmt.Errorf("Invalid webhook URL!")
	}

	if !strings.Contains(webhookUrl.Host, "discord.com") {
		return fmt.Errorf("Invalid webhook URL!")
	}

	data := map[string]any{
		"content": ctx.Arguments.RawFrom(1),
	}

	resp, err := http.DefaultClient().Post(webhookUrl.String(), data)
	if err != nil || resp.StatusCode != 204 {
		return fmt.Errorf("Error sending message to webhook: %s", string(resp.MustGetBody()))
	}

	c.services.Embed.Respond(ctx, "Successfully sent message to webhook.")
	return nil
}

func (c *AdminCommands) webhookdelete(ctx *router.Ctx) error {
	webhookUrl, err := url.Parse(ctx.Arguments.Get(0).Raw())
	if err != nil {
		return fmt.Errorf("Invalid webhook URL!")
	}

	if !strings.Contains(webhookUrl.Host, "discord.com") {
		return fmt.Errorf("Invalid webhook URL!")
	}

	_, _ = http.DefaultClient().Delete(webhookUrl.String(), nil)

	c.services.Embed.Respond(ctx, "Successfully deleted webhook.")
	return nil
}

// Member Control

func (c *AdminCommands) userinfo(ctx *router.Ctx) error {
	user, err := ctx.Arguments.Get(0).User(router.ArgumentFlagFallbackDefault)
	if err != nil {
		return fmt.Errorf("Cannot find user")
	}

	var member *discordgo.Member
	if ctx.Message.GuildID != "" {
		member, _ = ctx.Session.QueryMember(ctx.Message.GuildID, user.ID, true)
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

	emb.AddField("Is Bot", formatting.BoolYesNo(user.Bot)).
		AddField("Created At", discordgo.MustSnowflakeTimestamp(user.ID).Format(time.RFC1123)).
		AddField("Premium Type", formatting.PremiumType(user.PremiumType))

	if flags := formatting.UserFlags(user.PublicFlags); flags != "" {
		emb.AddField("Flags", formatting.UserFlags(user.PublicFlags))
	}

	if member != nil {
		if member.Nick != "" {
			emb.AddField("Nickname", member.Nick)
		}

		roles := fn.Map(member.Roles, func(roleID string) string {
			role, err := ctx.Session.State.Role(ctx.Message.GuildID, roleID)
			if err != nil {
				return roleID
			}
			return fmt.Sprintf("%s - %s", role.Name, role.ID)
		})

		emb.AddField("Joined Server On", member.JoinedAt.Format(time.RFC1123)).
			AddField("Roles", fmt.Sprintf("%d Roles:\n%s", len(member.Roles), strings.Join(roles, "\n")))

		if p := member.Presence; p != nil {
			emb.AddField("Status", formatting.StatusEmoji(p.Status)).
				AddField("Platform", formatting.ClientStatus(p.ClientStatus))

			if len(p.Activities) > 0 {
				emb.AddField("Activity", formatting.Activities(p.Activities))
			}
		}
	}

	emb.Send()
	return nil
}

func (c *AdminCommands) userage(ctx *router.Ctx) error {
	user, err := ctx.Arguments.Get(0).User()
	if err != nil {
		return fmt.Errorf("Could not find user")
	}
	age := time.Since(discordgo.MustSnowflakeTimestamp(user.ID))
	c.services.Embed.New(ctx).
		SetDescription("User Age (%s - %s)", user.String(), user.ID).
		AddField("Age", formatting.DurationLong(age)).
		Send()
	return nil
}

func (c *AdminCommands) mutualfriends(ctx *router.Ctx) error {
	user, err := ctx.Arguments.Get(0).User()
	if err != nil {
		return fmt.Errorf("Could not find user")
	}
	mutualFriends, err := ctx.Session.UserMutualFriends(user.ID)
	if err != nil {
		return fmt.Errorf("Could not find user")
	}
	mutualFriendsString := ""
	for _, friend := range mutualFriends {
		mutualFriendsString += fmt.Sprintf("%s - %s\n", friend.String(), friend.ID)
	}
	c.services.Embed.New(ctx).
		SetDescription("Mutual Friends (%s - %s)", user.String(), user.ID).
		AddField("Mutual Friends", mutualFriendsString).
		Send()
	return nil
}

func (c *AdminCommands) mutualservers(ctx *router.Ctx) error {
	user, err := ctx.Arguments.Get(0).User()
	if err != nil {
		return fmt.Errorf("Could not find user")
	}
	mutualServers, err := ctx.Session.UserProfile(user.ID)
	if err != nil {
		return fmt.Errorf("Could not find user")
	}
	mutualServersString := ""
	for _, guild_ := range mutualServers.MutualGuilds {
		guild, err := ctx.Session.GuildWithFallback(guild_.ID)
		if err != nil {
			continue
		}
		mutualServersString += fmt.Sprintf("%s - %s\n", guild.Name, guild.ID)
	}
	c.services.Embed.New(ctx).
		SetDescription("Mutual Servers (%s - %s)", user.String(), user.ID).
		AddField("Mutual Servers", mutualServersString).
		Send()
	return nil
}

func (c *AdminCommands) platform(ctx *router.Ctx) error {
	member, err := ctx.Arguments.Get(0).Member()
	if err != nil {
		return fmt.Errorf("Could not find user")
	}

	if member.Presence == nil {
		return fmt.Errorf("Users presence is unknown!")
	}

	c.services.Embed.New(ctx).
		AddField("User", member.User.Username).
		AddField("Status", formatting.StatusEmoji(member.Presence.Status)).
		AddField("Platform", formatting.ClientStatus(member.Presence.ClientStatus)).
		Send()
	return nil
}

func (c *AdminCommands) avatar(ctx *router.Ctx) error {
	user, err := ctx.Arguments.Get(0).User()
	if err != nil {
		return fmt.Errorf("Could not find user")
	}
	if user.Avatar == "" {
		return fmt.Errorf("User has no avatar")
	}
	return ctx.RespondText("%s", user.AvatarURL(""))
}

func (c *AdminCommands) banner(ctx *router.Ctx) error {
	user, err := ctx.Arguments.Get(0).User()
	if err != nil {
		return fmt.Errorf("Could not find user")
	}
	if user.Banner == "" {
		return fmt.Errorf("User has no banner")
	}
	return ctx.RespondText("%s", user.BannerURL(""))
}

func (c *AdminCommands) timeout(ctx *router.Ctx) error {
	member, err := ctx.Arguments.Get(0).Member()
	if err != nil {
		return fmt.Errorf("Could not find user")
	}

	duration, err := ctx.Arguments.Get(1).Duration()
	duration = fn.Clamp(duration, 1*time.Second, 28*24*time.Hour)
	if err != nil {
		return fmt.Errorf("Could not parse duration")
	}

	until := time.Now().Add(duration)

	err = ctx.Session.GuildMemberTimeout(ctx.Message.GuildID, member.User.ID, &until)
	if err != nil {
		return fmt.Errorf("Error timing out user: %v", err)
	}

	c.services.Embed.Respond(ctx, "Successfully timed out user for %s!", duration.String())
	return nil
}

func (c *AdminCommands) removetimeout(ctx *router.Ctx) error {
	member, err := ctx.Arguments.Get(0).Member()
	if err != nil {
		return fmt.Errorf("Could not find user")
	}

	err = ctx.Session.GuildMemberTimeout(ctx.Message.GuildID, member.User.ID, nil)
	if err != nil {
		return fmt.Errorf("Error removing timeout: %v", err)
	}

	c.services.Embed.Respond(ctx, "Successfully removed timeout!")
	return nil
}

func (c *AdminCommands) kick(ctx *router.Ctx) error {
	member, err := ctx.Arguments.Get(0).Member()
	if err != nil {
		return fmt.Errorf("Could not find user")
	}

	err = ctx.Session.GuildMemberDeleteWithReason(ctx.Message.GuildID, member.User.ID, ctx.Arguments.RawFrom(1))
	if err != nil {
		return fmt.Errorf("Error kicking user: %v", err)
	}

	c.services.Embed.Respond(ctx, "Successfully kicked user!")
	return nil
}

func (c *AdminCommands) ban(ctx *router.Ctx) error {
	member, err := ctx.Arguments.Get(0).Member(router.ArgumentFlagFallbackDefault)
	if err != nil {
		return fmt.Errorf("Could not find user")
	}

	err = ctx.Session.GuildBanCreateWithReason(ctx.Message.GuildID, member.User.ID, ctx.Arguments.RawFrom(1), 0)
	if err != nil {
		return fmt.Errorf("Error banning user: %v", err)
	}

	c.services.Embed.Respond(ctx, "Successfully banned user!")
	return nil
}

func (c *AdminCommands) softban(ctx *router.Ctx) error {
	member, err := ctx.Arguments.Get(0).Member()
	if err != nil {
		return fmt.Errorf("Could not find user")
	}

	err = ctx.Session.GuildBanCreateWithReason(ctx.Message.GuildID, member.User.ID, ctx.Arguments.RawFrom(1), 0)
	if err != nil {
		return fmt.Errorf("Error softbanning user: %v", err)
	}
	err = ctx.Session.GuildBanDelete(ctx.Message.GuildID, member.User.ID)
	if err != nil {
		return fmt.Errorf("Error softbanning user: %v", err)
	}

	c.services.Embed.Respond(ctx, "Successfully softbanned user!")
	return nil
}

func (c *AdminCommands) unban(ctx *router.Ctx) error {
	err := ctx.Session.GuildBanDelete(ctx.Message.GuildID, ctx.Arguments.Get(0).UserMentionID())
	if err != nil {
		return fmt.Errorf("Error unbanning user: %v", err)
	}

	c.services.Embed.Respond(ctx, "Successfully unbanned user!")
	return nil
}

func (c *AdminCommands) banlist(ctx *router.Ctx) error {
	guild, err := ctx.Arguments.Get(0).Guild()
	if err != nil {
		return fmt.Errorf("Error getting server: %v", err)
	}

	bans, err := ctx.Session.GuildBans(guild.ID, 100, "", "")
	if err != nil {
		return fmt.Errorf("Error getting ban list: %v", err)
	}

	banReadable := fn.Map(bans, func(ban *discordgo.GuildBan) string {
		return fmt.Sprintf("%s - %s", ban.User.Username, ban.User.ID)
	})

	c.services.Embed.New(ctx).
		AddField("Banned Amount", strconv.Itoa(len(bans))).
		AddField("Ban List", strings.Join(banReadable, "\n")).
		Send()
	return nil
}

func (c *AdminCommands) copybans(ctx *router.Ctx) error {
	srcGuild, err := ctx.Arguments.Get(0).Guild()
	if err != nil {
		return fmt.Errorf("Error getting server: %v", err)
	}

	destGuild, err := ctx.Arguments.Get(1).Guild(router.ArgumentFlagFallbackDefault)
	if err != nil {
		return fmt.Errorf("Error getting server: %v", err)
	}

	bans, err := ctx.Session.GuildBans(srcGuild.ID, 1000, "", "")
	if err != nil {
		return fmt.Errorf("Error copying bans: %v", err)
	}

	for _, ban := range bans {
		err := ctx.Session.GuildBanCreateWithReason(destGuild.ID, ban.User.ID, ban.Reason, 0)
		if err != nil {
			return fmt.Errorf("Error copying bans: %v", err)
		}

		time.Sleep(2 * time.Second)
	}

	c.services.Embed.Respond(ctx, "Successfully copied bans from %s to %s!", srcGuild.Name, destGuild.Name)
	return nil
}

// Nickname Control

func (c *AdminCommands) nick(ctx *router.Ctx) error {
	err := ctx.Session.GuildMemberNickname(ctx.Message.GuildID, "@me", ctx.Arguments.Raw())
	if err != nil {
		return fmt.Errorf("Error changing nickname: %v", err)
	}
	c.services.Embed.Respond(ctx, "Successfully changed nickname!")
	return nil
}

func (c *AdminCommands) nickuser(ctx *router.Ctx) error {
	member, err := ctx.Arguments.Get(0).Member()
	if err != nil {
		return fmt.Errorf("Could not find user")
	}

	err = ctx.Session.GuildMemberNickname(ctx.Message.GuildID, strings.ReplaceAll(member.User.ID, ctx.Message.Author.ID, "@me"), ctx.Arguments.RawFrom(1))
	if err != nil {
		return fmt.Errorf("Error changing nickname: %v", err)
	}
	c.services.Embed.Respond(ctx, "Successfully changed nickname!")
	return nil
}

func (c *AdminCommands) clearnick(ctx *router.Ctx) error {
	member, err := ctx.Arguments.Get(0).Member()
	if err != nil {
		return fmt.Errorf("Could not find user")
	}

	err = ctx.Session.GuildMemberNickname(ctx.Message.GuildID, member.User.ID, "")
	if err != nil {
		return fmt.Errorf("Error changing nickname: %v", err)
	}
	c.services.Embed.Respond(ctx, "Successfully changed nickname!")
	return nil
}

// Server Control

func (c *AdminCommands) serverinfo(ctx *router.Ctx) error {
	guild, err := ctx.Arguments.Get(0).Guild(router.ArgumentFlagFallbackDefault)
	if err != nil {
		return fmt.Errorf("Error getting server: %v", err)
	}

	ownerMention := "Error"

	owner, err := ctx.Session.QueryMember(guild.ID, guild.OwnerID, false)
	if err == nil {
		ownerMention = owner.User.String()
	}

	c.services.Embed.New(ctx).
		SetDescription("Server Info - %s", guild.Name).
		AddField("ID", guild.ID).
		AddField("Owner", ownerMention).
		AddField("Created At", discordgo.MustSnowflakeTimestamp(guild.ID).Format(time.RFC1123)).
		AddField("Verification Level", verificationLevels[int(guild.VerificationLevel)]).
		AddField("Boost Level", boostLevels[int(guild.PremiumTier)]).
		AddField("Amount of boosts", fmt.Sprint(guild.PremiumSubscriptionCount)).
		AddField("Member Count", fmt.Sprint(guild.MemberCount)).
		AddField("Channel Count", fmt.Sprint(len(guild.Channels))).
		AddField("Role Count", fmt.Sprint(len(guild.Roles))).
		AddField("Joined At", guild.JoinedAt.Format(time.RFC1123)).
		SetImage(guild.IconURL()).
		Send()
	return nil
}

func (c *AdminCommands) serverage(ctx *router.Ctx) error {
	guild, err := ctx.Arguments.Get(0).Guild(router.ArgumentFlagFallbackDefault)
	if err != nil {
		return fmt.Errorf("Error getting server: %v", err)
	}
	age := time.Since(discordgo.MustSnowflakeTimestamp(guild.ID))
	c.services.Embed.New(ctx).
		SetDescription("Server Age (%s - %s)", guild.Name, guild.ID).
		AddField("Age", formatting.DurationLong(age)).
		Send()
	return nil
}

func (c *AdminCommands) serverowner(ctx *router.Ctx) error {
	guild, err := ctx.Arguments.Get(0).Guild(router.ArgumentFlagFallbackDefault)
	if err != nil {
		return fmt.Errorf("Error getting server: %v", err)
	}
	owner, err := ctx.Session.QueryMember(guild.ID, guild.OwnerID, false)
	if err != nil {
		return fmt.Errorf("Error getting server owner: %v", err)
	}
	c.services.Embed.New(ctx).
		SetDescription("%s", guild.Name).
		AddField("Owner", owner.Mention()).
		Send()
	return nil
}

func (c *AdminCommands) servericon(ctx *router.Ctx) error {
	guild, err := ctx.Arguments.Get(0).Guild(router.ArgumentFlagFallbackDefault)
	if err != nil {
		return fmt.Errorf("Error getting server: %v", err)
	}

	emb := c.services.Embed.New(ctx).
		SetDescription("Server Icon - %s", guild.Name).
		SetImage(guild.IconURL())
	if guild.Icon == "" {
		emb.AddField("Error", "Server has no icon")
	}
	emb.Send()
	return nil
}

func (c *AdminCommands) serverbanner(ctx *router.Ctx) error {
	guild, err := ctx.Arguments.Get(0).Guild(router.ArgumentFlagFallbackDefault)
	if err != nil {
		return fmt.Errorf("Error getting server: %v", err)
	}
	if guild.Banner == "" {
		c.services.Embed.New(ctx).AddField("Error", "Server has no banner!").Send()
	}
	c.services.Embed.New(ctx).
		SetDescription("Server Banner - %s", guild.Name).
		SetImage(guild.BannerURL()).
		Send()
	return nil
}

func (c *AdminCommands) serversplash(ctx *router.Ctx) error {
	guild, err := ctx.Arguments.Get(0).Guild(router.ArgumentFlagFallbackDefault)
	if err != nil {
		return fmt.Errorf("Error getting server: %v", err)
	}
	if guild.Splash == "" {
		return fmt.Errorf("Server has no invite splash!")
	}
	c.services.Embed.New(ctx).
		SetDescription("Server Invite Splash - %s", guild.Name).
		SetImage(discordgo.EndpointGuildSplash(guild.ID, guild.Splash)).
		Send()
	return nil
}

func (c *AdminCommands) serverdsplash(ctx *router.Ctx) error {
	guild, err := ctx.Arguments.Get(0).Guild(router.ArgumentFlagFallbackDefault)
	if err != nil {
		return fmt.Errorf("Error getting server: %v", err)
	}
	if guild.DiscoverySplash == "" {
		return fmt.Errorf("Server has no discovery splash!")
	}
	c.services.Embed.New(ctx).
		SetDescription("Server Discovery Splash - %s", guild.Name).
		SetImage(discordgo.EndpointGuildSplash(guild.ID, guild.DiscoverySplash)).
		Send()
	return nil
}

func (c *AdminCommands) serverboosters(ctx *router.Ctx) error {
	guild, err := ctx.Arguments.Get(0).Guild(router.ArgumentFlagFallbackDefault)
	if err != nil {
		return fmt.Errorf("Error getting server: %v", err)
	}

	boosters := fn.Filter(guild.Members, func(member *discordgo.Member) bool {
		return member.PremiumSince != nil
	})

	boostersFormatted := fn.Map(boosters, func(member *discordgo.Member) string {
		return fmt.Sprintf("%s - %s", member.User.String(), member.User.ID)
	})

	if len(boostersFormatted) == 0 {
		boostersFormatted = []string{"No boosters"}
	}

	c.services.Embed.New(ctx).
		SetDescription("Server Boosters - %s", guild.Name).
		AddField("Boosters", strings.Join(boostersFormatted, "\n")).
		Send()
	return nil
}

func (c *AdminCommands) serverboostlevel(ctx *router.Ctx) error {
	guild, err := ctx.Arguments.Get(0).Guild(router.ArgumentFlagFallbackDefault)
	if err != nil {
		return fmt.Errorf("Error getting server: %v", err)
	}

	c.services.Embed.New(ctx).
		SetDescription("Server Boost Level - %s", guild.Name).
		AddField("Boost Level", boostLevels[int(guild.PremiumTier)]).
		Send()
	return nil
}

func (c *AdminCommands) setservername(ctx *router.Ctx) error {
	guild, err := ctx.Session.GuildEdit(ctx.Message.GuildID, &discordgo.GuildParams{
		Name: ctx.Arguments.Raw(),
	})
	if err != nil {
		return fmt.Errorf("Error changing server name: %v", err)
	}
	c.services.Embed.Respond(ctx, "Successfully changed server name to %s!", guild.Name)
	return nil
}

func (c *AdminCommands) setservericon(ctx *router.Ctx) error {
	_, err := ctx.Session.GuildEdit(ctx.Message.GuildID, &discordgo.GuildParams{
		Icon: ctx.Arguments.Get(0).Raw(),
	})
	if err != nil {
		return fmt.Errorf("Error changing server icon: %v", err)
	}
	c.services.Embed.Respond(ctx, "Successfully changed server icon!")
	return nil
}

func (c *AdminCommands) setserverbanner(ctx *router.Ctx) error {
	_, err := ctx.Session.GuildEdit(ctx.Message.GuildID, &discordgo.GuildParams{
		Banner: ctx.Arguments.Get(0).Raw(),
	})
	if err != nil {
		return fmt.Errorf("Error changing server icon: %v", err)
	}
	c.services.Embed.Respond(ctx, "Successfully changed server icon!")
	return nil
}

func (c *AdminCommands) lockserver(ctx *router.Ctx) error {
	guild, err := ctx.Arguments.Get(0).Guild(router.ArgumentFlagFallbackDefault)
	if err != nil {
		return fmt.Errorf("Error getting server: %v", err)
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
			return fmt.Errorf("Error locking server: %v", err)
		}
	}

	c.services.Embed.Respond(ctx, "Successfully locked server!")
	return nil
}

func (c *AdminCommands) unlockserver(ctx *router.Ctx) error {
	guild, err := ctx.Arguments.Get(0).Guild(router.ArgumentFlagFallbackDefault)
	if err != nil {
		return fmt.Errorf("Error getting server: %v", err)
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
			return fmt.Errorf("Error unlocking server: %v", err)
		}
	}

	c.services.Embed.Respond(ctx, "Successfully unlocked server!")
	return nil
}

func (c *AdminCommands) channels(ctx *router.Ctx) error {
	guild, err := ctx.Arguments.Get(0).Guild(router.ArgumentFlagFallbackDefault)
	if err != nil {
		return fmt.Errorf("Error getting server: %v", err)
	}

	emb := c.services.Embed.New(ctx).SetDescription("Server Channels - %s", guild.Name)
	categoryData := discordgo.MapCategories(guild.Channels)

	for _, category := range categoryData {
		emb.AddField("Category: "+category.GetCategoryName(), strings.Join(fn.Map(category.Channels, func(channel *discordgo.Channel) string {
			return fmt.Sprintf("  %s (ID: %s)", channel.Name, channel.ID)
		}), "\n"))
	}

	emb.Send()
	return nil
}

func (c *AdminCommands) bots(ctx *router.Ctx) error {
	guild, err := ctx.Arguments.Get(0).Guild(router.ArgumentFlagFallbackDefault)
	if err != nil {
		return fmt.Errorf("Error getting server: %v", err)
	}

	bots := fn.Filter(guild.Members, func(member *discordgo.Member) bool { return member.User.Bot })
	if len(bots) == 0 {
		return fmt.Errorf("No bots found!")
	}

	botsReadable := fn.Map(bots, func(member *discordgo.Member) string {
		return fmt.Sprintf("%s - %s", member.User.String(), member.User.ID)
	})

	c.services.Embed.New(ctx).
		SetDescription("Server Bots - %s", guild.Name).
		AddField("Server Bots", strings.Join(botsReadable, "\n")).
		Send()
	return nil
}

func (c *AdminCommands) admins(ctx *router.Ctx) error {
	guild, err := ctx.Arguments.Get(0).Guild(router.ArgumentFlagFallbackDefault)
	if err != nil {
		return fmt.Errorf("Error getting server: %v", err)
	}

	if len(guild.Channels) == 0 {
		return fmt.Errorf("Server has no channels!")
	}

	admins := fn.Filter(guild.Members, func(member *discordgo.Member) bool {
		permissions, _ := ctx.Session.UserChannelPermissions(member.User.ID, guild.Channels[0].ID)
		return permissions&discordgo.PermissionAdministrator != 0 && !member.User.Bot
	})
	if len(admins) == 0 {
		return fmt.Errorf("No admins found!")
	}

	adminsReadable := fn.Map(admins, func(member *discordgo.Member) string {
		return fmt.Sprintf("%s - %s", member.User.String(), member.User.ID)
	})

	c.services.Embed.New(ctx).
		SetDescription("Server Admins - %s", guild.Name).
		AddField("Server Admins", strings.Join(adminsReadable, "\n")).
		Send()
	return nil
}

func (c *AdminCommands) roles(ctx *router.Ctx) error {
	guild, err := ctx.Arguments.Get(0).Guild(router.ArgumentFlagFallbackDefault)
	if err != nil {
		return fmt.Errorf("Error getting server: %v", err)
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
	return nil
}

func (c *AdminCommands) leaveserver(ctx *router.Ctx) error {
	guild, err := ctx.Arguments.Get(0).Guild(router.ArgumentFlagFallbackDefault)
	if err != nil {
		return fmt.Errorf("Error getting server: %v", err)
	}
	return ctx.Session.GuildLeave(guild.ID)
}

func (c *AdminCommands) deleteserver(ctx *router.Ctx) error {
	guild, err := ctx.Arguments.Get(0).Guild(router.ArgumentFlagFallbackDefault)
	if err != nil {
		return fmt.Errorf("Error getting server: %v", err)
	}

	_, err = ctx.Session.GuildDelete(guild.ID)
	if err != nil {
		return fmt.Errorf("Error deleting server: %v", err)
	}
	return nil
}

// Invite Control

func (c *AdminCommands) invite(ctx *router.Ctx) error {
	channel, err := ctx.Arguments.Get(0).Channel(router.ArgumentFlagFallbackDefault)
	if err != nil {
		return fmt.Errorf("Could not find channel")
	}

	maxAge, _ := ctx.Arguments.Get(1).Int()
	maxAge = fn.Clamp(maxAge, 0, 604800)

	maxUses, _ := ctx.Arguments.Get(2).Int()
	maxUses = fn.Clamp(maxUses, 0, 100)

	invite, err := ctx.Session.ChannelInviteCreate(channel.ID, maxAge, maxUses, false)
	if err != nil {
		return fmt.Errorf("Error creating invite: %v", err)
	}

	c.services.Embed.Respond(ctx, "Created invite: https://discord.gg/%s", invite.Code)
	return nil
}

func (c *AdminCommands) friendinvite(ctx *router.Ctx) error {
	invite, err := ctx.Session.CreateFriendInvite("@me")
	if err != nil {
		return fmt.Errorf("Error creating friend invite: %v", err)
	}

	c.services.Embed.New(ctx).
		SetDescription("Successfully created friend invite!").
		AddField("Invite", fmt.Sprintf("https://discord.gg/%s", invite.Code)).
		Send()
	return nil
}

func (c *AdminCommands) inviteinfo(ctx *router.Ctx) error {
	inviteCode, err := ctx.Arguments.Get(0).Invite()
	if err != nil {
		return fmt.Errorf("Could not find invite: %v", err)
	}

	invite, err := ctx.Session.InviteWithCounts(inviteCode)
	if err != nil {
		return fmt.Errorf("Error getting invite: %v", err)
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
	return nil
}

func (c *AdminCommands) deleteinvite(ctx *router.Ctx) error {
	inviteCode, err := ctx.Arguments.Get(0).Invite()
	if err != nil {
		return fmt.Errorf("Could not find invite: %v", err)
	}

	invite, err := ctx.Session.InviteWithCounts(inviteCode)
	if err != nil {
		return fmt.Errorf("Error getting invite: %v", err)
	}

	invite, err = ctx.Session.InviteDelete(invite.Code)
	if err != nil {
		return fmt.Errorf("Error deleting invite: %v", err)
	}

	c.services.Embed.Respond(ctx, "Deleted invite: %s", invite.Code)
	return nil
}

func (c *AdminCommands) serverinvites(ctx *router.Ctx) error {
	guild, err := ctx.Arguments.Get(0).Guild(router.ArgumentFlagFallbackDefault)
	if err != nil {
		return fmt.Errorf("Error getting server: %v", err)
	}

	invites, err := ctx.Session.GuildInvites(guild.ID)
	if err != nil {
		return fmt.Errorf("Error getting invites: %v", err)
	}

	invitesFormatted := fn.Map(invites, func(invite *discordgo.Invite) string {
		username := "Unknown"
		if invite.Inviter != nil {
			username = invite.Inviter.String()
		}

		expiresAt := "Never"
		if invite.ExpiresAt != nil {
			expiresAt = formatting.Timestamp(*invite.ExpiresAt)
		}

		return fmt.Sprintf(
			"Link: https://discord.gg/%s - Expires At: %s - Created By: %s",
			invite.Code,
			expiresAt,
			username,
		)
	})

	c.services.Embed.New(ctx).
		SetDescription("Invite List").
		AddField("Invites", strings.Join(invitesFormatted, "\n")).
		Send()
	return nil
}

func (c *AdminCommands) channelinvites(ctx *router.Ctx) error {
	channel, err := ctx.Arguments.Get(0).Channel(router.ArgumentFlagFallbackDefault)
	if err != nil {
		return fmt.Errorf("Could not find channel")
	}

	invites, err := ctx.Session.ChannelInvites(channel.ID)
	if err != nil {
		return fmt.Errorf("Error getting invites: %v", err)
	}

	invitesFormatted := fn.Map(invites, func(invite *discordgo.Invite) string {
		username := "Unknown"
		if invite.Inviter != nil {
			username = invite.Inviter.String()
		}

		expiresAt := "Never"
		if invite.ExpiresAt != nil {
			expiresAt = formatting.Timestamp(*invite.ExpiresAt)
		}

		return fmt.Sprintf(
			"Link: https://discord.gg/%s - Expires At: %s - Created By: %s",
			invite.Code,
			expiresAt,
			username,
		)
	})

	c.services.Embed.New(ctx).
		SetDescription("Invite List").
		AddField("Invites", strings.Join(invitesFormatted, "\n")).
		Send()
	return nil
}

func (c *AdminCommands) userinvites(ctx *router.Ctx) error {
	member, err := ctx.Arguments.Get(0).Member()
	if err != nil {
		return fmt.Errorf("Could not find user")
	}

	guild, err := ctx.Arguments.Get(1).Guild(router.ArgumentFlagFallbackDefault)
	if err != nil {
		return fmt.Errorf("Error getting server: %v", err)
	}

	invites, err := ctx.Session.GuildInvites(guild.ID)
	if err != nil {
		return fmt.Errorf("Error getting invites: %v", err)
	}

	invites = fn.FilterInPlace(invites, func(invite *discordgo.Invite) bool {
		return invite.Inviter != nil && invite.Inviter.ID == member.User.ID
	})

	invitesFormatted := fn.Map(invites, func(invite *discordgo.Invite) string {
		expiresAt := "Never"
		if invite.ExpiresAt != nil {
			expiresAt = formatting.Timestamp(*invite.ExpiresAt)
		}

		return fmt.Sprintf(
			"Link: https://discord.gg/%s - Expires At: %s",
			invite.Code,
			expiresAt,
		)
	})

	c.services.Embed.New(ctx).
		SetDescription("Invite List").
		AddField("Invites", strings.Join(invitesFormatted, "\n")).
		Send()
	return nil
}

func (c *AdminCommands) roleinfo(ctx *router.Ctx) error {
	role, err := ctx.Arguments.Get(0).Role()
	if err != nil {
		return fmt.Errorf("Could not find role")
	}
	// Get all members in guild and check if they have the role
	guild, err := ctx.Session.GuildWithFallback(ctx.Message.GuildID)
	if err != nil {
		return fmt.Errorf("Error getting server: %v", err)
	}
	memberCount := 0
	for _, member := range guild.Members {
		if slices.Contains(member.Roles, role.ID) {
			memberCount++
		}
	}
	createdAt := discordgo.MustSnowflakeTimestamp(role.ID)
	c.services.Embed.New(ctx).
		SetDescription("Role Info (%s)", role.Name).
		AddField("ID", role.ID).
		AddField("Color", fmt.Sprintf("#%x", role.Color)).
		AddField("Created at", formatting.Timestamp(createdAt)).
		AddField("Position", strconv.Itoa(role.Position)).
		AddField("Is Hoist", formatting.BoolYesNo(role.Hoist)).
		AddField("Is Mentionable", formatting.BoolYesNo(role.Mentionable)).
		AddField("Members With Role", strconv.Itoa(memberCount)).
		AddField("Permissions", formatting.Permissions(role.Permissions)).
		Send()
	return nil
}

func (c *AdminCommands) giverole(ctx *router.Ctx) error {
	member, err := ctx.Arguments.Get(0).Member()
	if err != nil {
		return fmt.Errorf("Could not find user")
	}
	role, err := ctx.Arguments.Get(1).Role()
	if err != nil {
		return fmt.Errorf("Could not find role")
	}

	err = ctx.Session.GuildMemberRoleAdd(ctx.Message.GuildID, member.User.ID, role.ID)
	if err != nil {
		return fmt.Errorf("Error giving role: %v", err)
	}
	c.services.Embed.Respond(ctx, "Successfully gave role (%s) to (%s)", role.Name, member.User.String())
	return nil
}

func (c *AdminCommands) removerole(ctx *router.Ctx) error {
	member, err := ctx.Arguments.Get(0).Member()
	if err != nil {
		return fmt.Errorf("Could not find user")
	}
	role, err := ctx.Arguments.Get(1).Role()
	if err != nil {
		return fmt.Errorf("Could not find role")
	}

	err = ctx.Session.GuildMemberRoleRemove(ctx.Message.GuildID, member.User.ID, role.ID)
	if err != nil {
		return fmt.Errorf("Error removing role: %v", err)
	}
	c.services.Embed.Respond(ctx, "Successfully removed role (%s) from (%s)", role.Name, member.User.String())
	return nil
}

func (c *AdminCommands) createrole(ctx *router.Ctx) error {
	role, err := ctx.Session.GuildRoleCreate(ctx.Message.GuildID, &discordgo.RoleParams{
		Name: ctx.Arguments.Raw(),
	})
	if err != nil {
		return fmt.Errorf("Error creating role: %v", err)
	}

	c.services.Embed.Respond(ctx, "Successfully created role (%s)", role.Name)
	return nil
}

func (c *AdminCommands) deleterole(ctx *router.Ctx) error {
	role, err := ctx.Arguments.Get(0).Role()
	if err != nil {
		return fmt.Errorf("Could not find role")
	}
	err = ctx.Session.GuildRoleDelete(ctx.Message.GuildID, role.ID)
	if err != nil {
		return fmt.Errorf("Error deleting role: %v", err)
	}
	c.services.Embed.Respond(ctx, "Successfully deleted role (%s)", role.Name)
	return nil
}

func (c *AdminCommands) renamerole(ctx *router.Ctx) error {
	role, err := ctx.Arguments.Get(0).Role()
	if err != nil {
		return fmt.Errorf("Could not find role")
	}

	role, err = ctx.Session.GuildRoleEdit(ctx.Message.GuildID, role.ID, &discordgo.RoleParams{
		Name: ctx.Arguments.Raw(),
	})
	if err != nil {
		return fmt.Errorf("Error renaming role: %v", err)
	}

	c.services.Embed.Respond(ctx, "Successfully renamed role (%s)", role.Name)
	return nil
}
