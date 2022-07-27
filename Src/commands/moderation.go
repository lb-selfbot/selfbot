package commands

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/LightningDev1/LB-Selfbot/api"
	"github.com/LightningDev1/LB-Selfbot/embed"
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

type ModerationCommands struct{}

func (m ModerationCommands) RegisterCommands(router *dgc.Router) {
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
		Usage:       "[p]purge <amount>",
		Handler:     m.purge,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "userinfo",
		Description: "Get info about a user",
		Usage:       "[p]userinfo <user>",
		Handler:     m.userinfo,
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
		Name:        "serverinfo",
		Description: "Get info about a server",
		Usage:       "[p]serverinfo [server id]",
		Handler:     m.serverinfo,
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
}

// Channel Control

func (ModerationCommands) channelinfo(ctx *dgc.Ctx) {
	channel, success := utils.Discord.TryGetChannelFromMention(ctx, ctx.Arguments.Get(0))
	if !success {
		ctx.RespondText("Could not find channel")
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

func (ModerationCommands) textchannel(ctx *dgc.Ctx) {
	if !utils.Discord.CheckIfInServer(ctx) {
		return
	}
	channel, err := ctx.Session.GuildChannelCreate(ctx.Event.GuildID, ctx.Arguments.Raw(), discordgo.ChannelTypeGuildText)
	if err != nil {
		ctx.RespondText("Error creating channel: " + err.Error())
		return
	}
	embed.NewEmbed().
		SetDescription("Successfully created text channel!").
		AddField("Channel", channel.Mention()).
		Send(ctx)
}

func (ModerationCommands) voicechannel(ctx *dgc.Ctx) {
	if !utils.Discord.CheckIfInServer(ctx) {
		return
	}
	channel, err := ctx.Session.GuildChannelCreate(ctx.Event.GuildID, ctx.Arguments.Raw(), discordgo.ChannelTypeGuildVoice)
	if err != nil {
		ctx.RespondText("Error creating channel: " + err.Error())
		return
	}
	embed.NewEmbed().
		SetDescription("Successfully created voice channel!").
		AddField("Channel", channel.Mention()).
		Send(ctx)
}

func (ModerationCommands) stagechannel(ctx *dgc.Ctx) {
	if !utils.Discord.CheckIfInServer(ctx) {
		return
	}
	channel, err := ctx.Session.GuildChannelCreate(ctx.Event.GuildID, ctx.Arguments.Raw(), discordgo.ChannelTypeGuildStageVoice)
	if err != nil {
		ctx.RespondText("Error creating channel: " + err.Error())
		return
	}
	embed.NewEmbed().
		SetDescription("Successfully created stage channel!").
		AddField("Channel", channel.Mention()).
		Send(ctx)
}

func (ModerationCommands) newschannel(ctx *dgc.Ctx) {
	if !utils.Discord.CheckIfInServer(ctx) {
		return
	}
	channel, err := ctx.Session.GuildChannelCreate(ctx.Event.GuildID, ctx.Arguments.Raw(), discordgo.ChannelTypeGuildNews)
	if err != nil {
		ctx.RespondText("Error creating channel: " + err.Error())
		return
	}
	embed.NewEmbed().
		SetDescription("Successfully created stage channel!").
		AddField("Channel", channel.Mention()).
		Send(ctx)
}

func (ModerationCommands) thread(ctx *dgc.Ctx) {
	if !utils.Discord.CheckIfInServer(ctx) {
		return
	}
	ctx.Session.ThreadStart(ctx.Event.ChannelID, ctx.Arguments.Raw(), discordgo.ChannelTypeGuildPublicThread, 0)
	channel, err := ctx.Session.ThreadStart(ctx.Event.ChannelID, ctx.Arguments.Raw(), discordgo.ChannelTypeGuildPublicThread, 0)
	if err != nil {
		ctx.RespondText("Error creating channel: " + err.Error())
		return
	}
	embed.NewEmbed().
		SetDescription("Successfully created thread!").
		AddField("Thread", channel.Mention()).
		Send(ctx)
}

func (ModerationCommands) renamechannel(ctx *dgc.Ctx) {
	channel, success := utils.Discord.TryGetChannelFromMention(ctx, ctx.Arguments.Get(0))
	if !success {
		ctx.RespondText("Could not find channel")
		return
	}

	_, err := ctx.Session.ChannelEdit(channel.ID, ctx.Arguments.Raw())
	if err != nil {
		ctx.RespondText("Error renaming channel: " + err.Error())
		return
	}

	embed.NewEmbed().
		SetDescription("Successfully renamed channel!").
		AddField("Channel", channel.Mention()).
		Send(ctx)
}

func (ModerationCommands) deletechannel(ctx *dgc.Ctx) {
	channel, success := utils.Discord.TryGetChannelFromMention(ctx, ctx.Arguments.Get(0))
	if !success {
		ctx.RespondText("Could not find channel")
		return
	}

	_, err := ctx.Session.ChannelDelete(channel.ID)
	if err != nil {
		ctx.RespondText("Error deleting channel: " + err.Error())
		return
	}

	embed.NewEmbed().
		SetDescription("Successfully deleted channel!").
		Send(ctx)
}

func (ModerationCommands) slowmode(ctx *dgc.Ctx) {
	channel, success := utils.Discord.TryGetChannelFromMention(ctx, ctx.Arguments.Get(0))
	if !success {
		ctx.RespondText("Could not find channel")
		return
	}

	seconds, err := strconv.Atoi(utils.Arguments["<seconds>"])
	if err != nil {
		ctx.RespondText("Argument <seconds> must be a number.")
		return
	}

	_, err = ctx.Session.ChannelEditComplex(channel.ID, &discordgo.ChannelEdit{
		RateLimitPerUser: &seconds,
	})
	if err != nil {
		ctx.RespondText("Error setting slowmode: " + err.Error())
		return
	}

	embed.NewEmbed().
		SetDescription(fmt.Sprintf("Successfully set slowmode to %s seconds!", utils.Arguments["<seconds>"])).
		AddField("Channel", channel.Mention()).
		Send(ctx)
}

func (ModerationCommands) removeslowmode(ctx *dgc.Ctx) {
	channel, success := utils.Discord.TryGetChannelFromMention(ctx, ctx.Arguments.Get(0))
	if !success {
		ctx.RespondText("Could not find channel")
		return
	}

	_, err := ctx.Session.ChannelEditComplex(channel.ID, &discordgo.ChannelEdit{
		RateLimitPerUser: nil,
	})
	if err != nil {
		ctx.RespondText("Error removing slowmode: " + err.Error())
		return
	}

	embed.NewEmbed().
		SetDescription("Successfully removed slowmode!").
		AddField("Channel", channel.Mention()).
		Send(ctx)
}

func (ModerationCommands) lockchannel(ctx *dgc.Ctx) {
	channel, success := utils.Discord.TryGetChannelFromMention(ctx, ctx.Arguments.Get(0))
	if !success {
		ctx.RespondText("Could not find channel")
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
		ctx.RespondText("Error locking channel: " + err.Error())
		return
	}

	embed.NewEmbed().
		SetDescription("Successfully locked channel!").
		AddField("Channel", channel.Mention()).
		Send(ctx)
}

func (ModerationCommands) unlockchannel(ctx *dgc.Ctx) {
	channel, success := utils.Discord.TryGetChannelFromMention(ctx, ctx.Arguments.Get(0))
	if !success {
		ctx.RespondText("Could not find channel")
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
		ctx.RespondText("Error unlocking channel: " + err.Error())
		return
	}

	embed.NewEmbed().
		SetDescription("Successfully unlocked channel!").
		AddField("Channel", channel.Mention()).
		Send(ctx)
}

func (ModerationCommands) lockserver(ctx *dgc.Ctx) {
	guild, err := utils.Discord.GetGuild(ctx.Session, utils.Misc.Or(utils.Arguments["[server id]"], ctx.Event.GuildID))
	if err != nil {
		ctx.RespondText("Error getting server: " + err.Error())
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
				ctx.RespondText("Error locking server: " + err.Error())
				return
			}
		}
	}

	embed.NewEmbed().
		SetDescription("Successfully locked server!").
		Send(ctx)
}

func (ModerationCommands) unlockserver(ctx *dgc.Ctx) {
	guild, err := utils.Discord.GetGuild(ctx.Session, utils.Misc.Or(utils.Arguments["[server id]"], ctx.Event.GuildID))
	if err != nil {
		ctx.RespondText("Error getting server: " + err.Error())
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
				ctx.RespondText("Error locking server: " + err.Error())
				return
			}
		}
	}

	embed.NewEmbed().
		SetDescription("Successfully locked server!").
		Send(ctx)
}

func (ModerationCommands) category(ctx *dgc.Ctx) {
	if !utils.Discord.CheckIfInServer(ctx) {
		return
	}
	channel, err := ctx.Session.GuildChannelCreate(ctx.Event.GuildID, ctx.Arguments.Raw(), discordgo.ChannelTypeGuildCategory)
	if err != nil {
		ctx.RespondText("Error creating category: " + err.Error())
		return
	}
	embed.NewEmbed().
		SetDescription("Successfully created category!").
		AddField("Category", channel.Mention()).
		Send(ctx)
}

func (ModerationCommands) deletecategory(ctx *dgc.Ctx) {
	channel, success := utils.Discord.TryGetChannelFromMention(ctx, ctx.Arguments.Get(0))
	if !success {
		ctx.RespondText("Could not find category")
		return
	}
	_, err := ctx.Session.ChannelDelete(channel.ID)
	if err != nil {
		ctx.RespondText("Error deleting category: " + err.Error())
		return
	}

	embed.NewEmbed().
		SetDescription("Successfully deleted category!").
		Send(ctx)
}

func (ModerationCommands) purge(ctx *dgc.Ctx) {
	amount, err := strconv.Atoi(utils.Arguments["<amount>"])
	if err != nil {
		ctx.RespondText("Argument <amount> must be a number.")
		return
	}

	messages, err := utils.Discord.GetChannelMessages(ctx, amount)
	if err != nil {
		ctx.RespondText("Error getting messages: " + err.Error())
		return
	}

	for _, message := range messages {
		ctx.Session.ChannelMessageDelete(ctx.Event.ChannelID, message.ID)
	}
}

func (ModerationCommands) nuke(ctx *dgc.Ctx) {
	channel, success := utils.Discord.TryGetChannelFromMention(ctx, ctx.Arguments.Get(0))
	if !success {
		ctx.RespondText("Could not find channel")
		return
	}

	_, err := ctx.Session.ChannelDelete(channel.ID)
	if err != nil {
		ctx.RespondText("Error deleting channel: " + err.Error())
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
		ctx.RespondText("Error creating channel: " + err.Error())
		return
	}

	ctx.Event.ChannelID = channel.ID

	embed.NewEmbed().
		SetDescription("Successfully nuked channel!").
		Send(ctx)
}

// Member Control

func (ModerationCommands) userinfo(ctx *dgc.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		ctx.RespondText("Could not find user")
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

func (ModerationCommands) report(ctx *dgc.Ctx) {
	link := utils.Arguments["<message link>"]
	reason := utils.Arguments["<reason>"]

	// Verify message link
	url, err := url.Parse(link)
	if err != nil {
		ctx.RespondText("Invalid message link.")
		return
	}
	urlPathParts := strings.Split(url.Path, "/")
	if len(urlPathParts) < 4 {
		ctx.RespondText("Invalid message link.")
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
		ctx.RespondText("Argument `<reason>` must be one of: `illegalcontent`, `harassment`, `spam`, `phishing`, `selfharm`, `nsfw`")
		return
	}

	api.ReportMessage(guildID, channelID, messageID, reasonInt)

	embed.NewEmbed().
		SetDescription("Successfully reported message!").
		Send(ctx)
}

func (ModerationCommands) mute(ctx *dgc.Ctx) {
	if !utils.Discord.CheckIfInServer(ctx) {
		return
	}
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		ctx.RespondText("Could not find user")
		return
	}

	member, _ := ctx.Session.GuildMember(ctx.Event.GuildID, user.ID)
	roles, err := ctx.Session.GuildRoles(ctx.Event.GuildID)
	if member == nil || err != nil {
		ctx.RespondText("Error muting user: " + err.Error())
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
			ctx.RespondText("Error muting user: " + err.Error())
			return
		}
		embed.NewEmbed().
			SetDescription("Successfully muted user!").
			Send(ctx)
		return
	}

	channels, err := ctx.Session.GuildChannels(ctx.Event.GuildID)
	if err != nil {
		ctx.RespondText("Error muting user: " + err.Error())
		return
	}

	role, err := ctx.Session.GuildRoleCreate(ctx.Event.GuildID, "Muted", 0, false, 0, false)
	if err != nil {
		ctx.RespondText("Error muting user: " + err.Error())
		return
	}

	for _, channel := range channels {
		if channel.Type == 1 {
			err = ctx.Session.ChannelPermissionSet(channel.ID, role.ID, discordgo.PermissionOverwriteTypeRole, 0, discordgo.PermissionSendMessages)
			if err != nil {
				ctx.RespondText("Error muting user: " + err.Error())
				return
			}
		}
	}

	err = ctx.Session.GuildMemberEdit(ctx.Event.GuildID, user.ID, append(member.Roles, role.ID))
	if err != nil {
		ctx.RespondText("Error muting user: " + err.Error())
		return
	}

	embed.NewEmbed().
		SetDescription("Successfully muted user!").
		Send(ctx)
}

func (ModerationCommands) unmute(ctx *dgc.Ctx) {
	if !utils.Discord.CheckIfInServer(ctx) {
		return
	}

	member, success := utils.Discord.TryGetMemberFromMention(ctx, ctx.Arguments.Get(0))
	if !success {
		ctx.RespondText("Could not find user")
		return
	}

	roles, err := ctx.Session.GuildRoles(ctx.Event.GuildID)
	if member == nil || err != nil {
		ctx.RespondText("Error unmuting user: " + err.Error())
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
			ctx.RespondText("Error unmuting user: " + err.Error())
			return
		}
		embed.NewEmbed().
			SetDescription("Successfully unmuted user!").
			Send(ctx)
		return
	}
	ctx.RespondText("Error: No muted role found!")
}

func (ModerationCommands) timeout(ctx *dgc.Ctx) {
	if !utils.Discord.CheckIfInServer(ctx) {
		return
	}

	member, success := utils.Discord.TryGetMemberFromMention(ctx, ctx.Arguments.Get(0))
	if !success {
		ctx.RespondText("Could not find user")
		return
	}

	seconds, err := strconv.Atoi(utils.Arguments["<seconds>"])
	if err != nil {
		ctx.RespondText("Argument <seconds> must be a number.")
		return
	}
	time := time.Now().Add(time.Duration(seconds) * time.Second)

	err = ctx.Session.GuildMemberTimeout(ctx.Event.GuildID, member.User.ID, &time)
	if err != nil {
		ctx.RespondText("Error timing out user: " + err.Error())
		return
	}

	embed.NewEmbed().
		SetDescription(fmt.Sprintf("Successfully timed out user for %d seconds!", seconds)).
		Send(ctx)
}

func (ModerationCommands) removetimeout(ctx *dgc.Ctx) {
	if !utils.Discord.CheckIfInServer(ctx) {
		return
	}

	member, success := utils.Discord.TryGetMemberFromMention(ctx, ctx.Arguments.Get(0))
	if !success {
		ctx.RespondText("Could not find user")
		return
	}

	err := ctx.Session.GuildMemberTimeout(ctx.Event.GuildID, member.User.ID, nil)
	if err != nil {
		ctx.RespondText("Error removing timeout: " + err.Error())
		return
	}

	embed.NewEmbed().
		SetDescription("Successfully removed timeout!").
		Send(ctx)
}

func (ModerationCommands) kick(ctx *dgc.Ctx) {
	if !utils.Discord.CheckIfInServer(ctx) {
		return
	}

	member, success := utils.Discord.TryGetMemberFromMention(ctx, ctx.Arguments.Get(0))
	if !success {
		ctx.RespondText("Could not find user")
		return
	}

	err := ctx.Session.GuildMemberDeleteWithReason(ctx.Event.GuildID, member.User.ID, utils.Arguments["[reason]"])
	if err != nil {
		ctx.RespondText("Error kicking user: " + err.Error())
		return
	}

	embed.NewEmbed().
		SetDescription("Successfully kicked user!").
		Send(ctx)
}

func (ModerationCommands) ban(ctx *dgc.Ctx) { // hackban: same as ban
	if !utils.Discord.CheckIfInServer(ctx) {
		return
	}

	member, success := utils.Discord.TryGetMemberFromMention(ctx, ctx.Arguments.Get(0))
	if !success {
		ctx.RespondText("Could not find user")
		return
	}

	err := ctx.Session.GuildBanCreateWithReason(ctx.Event.GuildID, member.User.ID, utils.Arguments["[reason]"], 0)
	if err != nil {
		ctx.RespondText("Error banning user: " + err.Error())
		return
	}

	embed.NewEmbed().
		SetDescription("Successfully banned user!").
		Send(ctx)
}

func (ModerationCommands) softban(ctx *dgc.Ctx) {
	if !utils.Discord.CheckIfInServer(ctx) {
		return
	}

	member, success := utils.Discord.TryGetMemberFromMention(ctx, ctx.Arguments.Get(0))
	if !success {
		ctx.RespondText("Could not find user")
		return
	}

	err := ctx.Session.GuildBanCreateWithReason(ctx.Event.GuildID, member.User.ID, utils.Arguments["[reason]"], 0)
	if err != nil {
		ctx.RespondText("Error softbanning user: " + err.Error())
		return
	}
	err = ctx.Session.GuildBanDelete(ctx.Event.GuildID, member.User.ID)
	if err != nil {
		ctx.RespondText("Error softbanning user: " + err.Error())
		return
	}

	embed.NewEmbed().
		SetDescription("Successfully softbanned user!").
		Send(ctx)
}

func (ModerationCommands) unban(ctx *dgc.Ctx) {
	if !utils.Discord.CheckIfInServer(ctx) {
		return
	}

	err := ctx.Session.GuildBanDelete(ctx.Event.GuildID, utils.Arguments["<user id>"])
	if err != nil {
		ctx.RespondText("Error unbanning user: " + err.Error())
		return
	}

	embed.NewEmbed().
		SetDescription("Successfully softbanned user!").
		Send(ctx)
}

func (ModerationCommands) banlist(ctx *dgc.Ctx) {
	guildID := utils.Arguments["[server id]"]
	if guildID == "" {
		guildID = ctx.Event.GuildID
	}
	if guildID == "" {
		ctx.RespondText("Error: No server ID specified!")
		return
	}
	banList, err := ctx.Session.GuildBans(ctx.Event.GuildID, 1000, "", "")
	if err != nil {
		ctx.RespondText("Error getting ban list: " + err.Error())
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

func (ModerationCommands) copybans(ctx *dgc.Ctx) {
	sourceServerID := utils.Arguments["<source server id>"]
	destinationServerID := utils.Arguments["<destination server id>"]
	if destinationServerID == "" {
		destinationServerID = ctx.Event.GuildID
	}
	banList, err := ctx.Session.GuildBans(sourceServerID, 1000, "", "")
	if err != nil {
		ctx.RespondText("Error copying bans: " + err.Error())
		return
	}

	for _, ban := range banList {
		err := ctx.Session.GuildBanCreateWithReason(destinationServerID, ban.User.ID, ban.Reason, 0)
		if err != nil {
			ctx.RespondText("Error copying bans: " + err.Error())
			return
		}
	}
}

// Nickname Control

func (ModerationCommands) nick(ctx *dgc.Ctx) {
	if !utils.Discord.CheckIfInServer(ctx) {
		return
	}
	err := ctx.Session.GuildMemberNickname(ctx.Event.GuildID, ctx.Event.Author.ID, utils.Arguments["<text>"])
	if err != nil {
		ctx.RespondText("Error changing nickname: " + err.Error())
		return
	}
	embed.NewEmbed().
		SetDescription("Successfully changed nickname!").
		Send(ctx)
}

func (ModerationCommands) nickuser(ctx *dgc.Ctx) {
	if !utils.Discord.CheckIfInServer(ctx) {
		return
	}

	member, success := utils.Discord.TryGetMemberFromMention(ctx, ctx.Arguments.Get(0))
	if !success {
		ctx.RespondText("Could not find user")
		return
	}

	err := ctx.Session.GuildMemberNickname(ctx.Event.GuildID, member.User.ID, utils.Arguments["<text>"])
	if err != nil {
		ctx.RespondText("Error changing nickname: " + err.Error())
		return
	}
	embed.NewEmbed().
		SetDescription("Successfully changed nickname!").
		Send(ctx)
}

func (ModerationCommands) clearnick(ctx *dgc.Ctx) {
	if !utils.Discord.CheckIfInServer(ctx) {
		return
	}

	member, success := utils.Discord.TryGetMemberFromMention(ctx, ctx.Arguments.Get(0))
	if !success {
		member = ctx.Event.Member
	}

	err := ctx.Session.GuildMemberNickname(ctx.Event.GuildID, member.User.ID, "")
	if err != nil {
		ctx.RespondText("Error changing nickname: " + err.Error())
		return
	}
	embed.NewEmbed().
		SetDescription("Successfully changed nickname!").
		Send(ctx)
}

// Server Control

func (ModerationCommands) serverowner(ctx *dgc.Ctx) {
	if !utils.Discord.CheckIfInServer(ctx) {
		return
	}

	guild, err := utils.Discord.GetGuild(ctx.Session, utils.Misc.Or(utils.Arguments["[server id]"], ctx.Event.GuildID))
	if err != nil {
		ctx.RespondText("Error getting server: " + err.Error())
		return
	}
	owner, err := utils.Discord.GetMember(ctx.Session, guild.ID, guild.OwnerID)
	if err != nil {
		ctx.RespondText("Error getting server owner: " + err.Error())
		return
	}
	embed.NewEmbed().
		SetDescription(guild.Name).
		AddField("Owner", owner.Mention()).
		Send(ctx)
}

func (ModerationCommands) servericon(ctx *dgc.Ctx) {
	if !utils.Discord.CheckIfInServer(ctx) {
		return
	}

	guild, err := utils.Discord.GetGuild(ctx.Session, utils.Misc.Or(utils.Arguments["[server id]"], ctx.Event.GuildID))
	if err != nil {
		ctx.RespondText("Error getting server: " + err.Error())
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

func (ModerationCommands) serverbanner(ctx *dgc.Ctx) {
	if !utils.Discord.CheckIfInServer(ctx) {
		return
	}

	guild, err := utils.Discord.GetGuild(ctx.Session, utils.Misc.Or(utils.Arguments["[server id]"], ctx.Event.GuildID))
	if err != nil {
		ctx.RespondText("Error getting server: " + err.Error())
		return
	}

	embed := embed.NewEmbed().
		SetDescription(fmt.Sprintf("Server Banner - %s", guild.Name)).
		SetImage(guild.BannerURL())
	if guild.Banner == "" {
		embed.AddField("Error", "Server has no banner!")
	}
	embed.Send(ctx)
}

func (ModerationCommands) serverboosters(ctx *dgc.Ctx) {
	if !utils.Discord.CheckIfInServer(ctx) {
		return
	}

	guild, err := utils.Discord.GetGuild(ctx.Session, utils.Misc.Or(utils.Arguments["[server id]"], ctx.Event.GuildID))
	if err != nil {
		ctx.RespondText("Error getting server: " + err.Error())
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

func (ModerationCommands) serverboostlevel(ctx *dgc.Ctx) {
	if !utils.Discord.CheckIfInServer(ctx) {
		return
	}

	guild, err := utils.Discord.GetGuild(ctx.Session, utils.Misc.Or(utils.Arguments["[server id]"], ctx.Event.GuildID))
	if err != nil {
		ctx.RespondText("Error getting server: " + err.Error())
		return
	}

	embed.NewEmbed().
		SetDescription(fmt.Sprintf("Server Boost Level - %s", guild.Name)).
		AddField("Boost Level", boostLevels[int(guild.PremiumTier)]).
		Send(ctx)
}

func (ModerationCommands) setservername(ctx *dgc.Ctx) {
	if !utils.Discord.CheckIfInServer(ctx) {
		return
	}
	guild, err := ctx.Session.GuildEdit(ctx.Event.GuildID, discordgo.GuildParams{
		Name: utils.Arguments["<text>"],
	})
	if err != nil {
		ctx.RespondText("Error changing server name: " + err.Error())
		return
	}
	embed.NewEmbed().
		SetDescription(fmt.Sprintf("Successfully changed server name to %s!", guild.Name)).
		Send(ctx)
}

func (ModerationCommands) setservericon(ctx *dgc.Ctx) {
	if !utils.Discord.CheckIfInServer(ctx) {
		return
	}
	_, err := ctx.Session.GuildEdit(ctx.Event.GuildID, discordgo.GuildParams{
		Icon: utils.Arguments["<url>"],
	})
	if err != nil {
		ctx.RespondText("Error changing server icon: " + err.Error())
		return
	}
	embed.NewEmbed().
		SetDescription("Successfully changed server icon!").
		Send(ctx)
}

func (ModerationCommands) setserverbanner(ctx *dgc.Ctx) {
	if !utils.Discord.CheckIfInServer(ctx) {
		return
	}
	_, err := ctx.Session.GuildEdit(ctx.Event.GuildID, discordgo.GuildParams{
		Banner: utils.Arguments["<url>"],
	})
	if err != nil {
		ctx.RespondText("Error changing server icon: " + err.Error())
		return
	}
	embed.NewEmbed().
		SetDescription("Successfully changed server icon!").
		Send(ctx)
}

func (ModerationCommands) serverinfo(ctx *dgc.Ctx) {
	if !utils.Discord.CheckIfInServer(ctx) {
		return
	}

	guild, err := utils.Discord.GetGuild(ctx.Session, utils.Misc.Or(utils.Arguments["[server id]"], ctx.Event.GuildID))
	if err != nil {
		ctx.RespondText("Error getting server: " + err.Error())
		return
	}

	owner, err := utils.Discord.GetMember(ctx.Session, guild.ID, guild.OwnerID)
	ownerMention := "Error"
	if err == nil {
		ownerMention = owner.User.String()
	}

	embed := embed.NewEmbed().
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
		AddField("Joined At", guild.JoinedAt.Format(time.RFC1123))
	if guild.Icon != "" {
		embed.SetImage(guild.IconURL())
	}
	embed.Send(ctx)
}

// Invite Control

func (ModerationCommands) invite(ctx *dgc.Ctx) {
	channel, success := utils.Discord.TryGetChannelFromMention(ctx, ctx.Arguments.Get(0))
	if !success {
		ctx.RespondText("Could not find channel")
		return
	}

	ageInt, err := strconv.Atoi(utils.Arguments["[maxage]"])
	if err != nil && utils.Arguments["[maxage]"] != "" {
		ctx.RespondText("Argument [maxage] must be a number.")
		return
	}

	usesInt, err := strconv.Atoi(utils.Arguments["[maxuses]"])
	if err != nil && utils.Arguments["[maxuses]"] != "" {
		ctx.RespondText("Argument [maxuses] must be a number.")
		return
	}

	invite, err := ctx.Session.ChannelInviteCreate(channel.ID, discordgo.Invite{
		MaxAge:    ageInt,
		MaxUses:   usesInt,
		Temporary: false,
		Unique:    true,
	})
	if err != nil {
		ctx.RespondText("Error creating invite: " + err.Error())
		return
	}

	embed.NewEmbed().
		SetDescription(fmt.Sprintf("Created invite: https://discord.gg/%s", invite.Code)).
		AddField("Max Age", fmt.Sprint(invite.MaxAge)).
		AddField("Max Uses", fmt.Sprint(invite.MaxUses)).
		AddField("Channel", invite.Channel.Name).
		Send(ctx)
}

func (ModerationCommands) friendinvite(ctx *dgc.Ctx) {
	invite, err := ctx.Session.CreateFriendInvite("@me")
	if err != nil {
		ctx.RespondText("Error creating friend invite: " + err.Error())
		return
	}

	embed.NewEmbed().
		SetDescription("Successfully created friend invite!").
		AddField("Invite", fmt.Sprintf("https://discord.gg/%s", invite.Code)).
		Send(ctx)
}

func (ModerationCommands) inviteinfo(ctx *dgc.Ctx) {
	parsedInvite := utils.Discord.ParseInvite(utils.Arguments["<invite>"])
	invite, err := ctx.Session.InviteWithCounts(parsedInvite)
	if err != nil {
		ctx.RespondText("Error getting invite: " + err.Error())
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

func (ModerationCommands) deleteinvite(ctx *dgc.Ctx) {
	parsedInvite := utils.Discord.ParseInvite(utils.Arguments["<invite>"])
	invite, err := ctx.Session.InviteWithCounts(parsedInvite)
	if err != nil {
		ctx.RespondText("Error getting invite: " + err.Error())
		return
	}

	invite, err = ctx.Session.InviteDelete(invite.Code)
	if err != nil {
		ctx.RespondText("Error deleting invite: " + err.Error())
		return
	}

	embed.NewEmbed().
		SetDescription(fmt.Sprintf("Deleted invite: %s", invite.Code)).
		Send(ctx)
}

func (m *ModerationCommands) serverinvites(ctx *dgc.Ctx) {
	invites, err := ctx.Session.GuildInvites(utils.Misc.Or(utils.Arguments["[server id]"], ctx.Event.GuildID))
	if err != nil {
		ctx.RespondText("Error getting invites: " + err.Error())
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

func (m *ModerationCommands) channelinvites(ctx *dgc.Ctx) {
	channel, success := utils.Discord.TryGetChannelFromMention(ctx, ctx.Arguments.Get(0))
	channelID := ctx.Event.ChannelID
	if success {
		channelID = channel.ID
	}
	invites, err := ctx.Session.ChannelInvites(channelID)
	if err != nil {
		ctx.RespondText("Error getting invites: " + err.Error())
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

func (m *ModerationCommands) userinvites(ctx *dgc.Ctx) {
	if !utils.Discord.CheckIfInServer(ctx) {
		return
	}

	member, success := utils.Discord.TryGetMemberFromMention(ctx, ctx.Arguments.Get(0))
	if !success {
		ctx.RespondText("Could not find user")
		return
	}

	invites, err := ctx.Session.GuildInvites(utils.Misc.Or(utils.Arguments["[server id]"], ctx.Event.GuildID))
	if err != nil {
		ctx.RespondText("Error getting invites: " + err.Error())
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

func (ModerationCommands) roleinfo(ctx *dgc.Ctx) {
	role, success := utils.Discord.TryGetRoleFromMention(ctx, ctx.Arguments.Get(0))
	if !success {
		ctx.RespondText("Could not find role")
		return
	}
	// Get all members in guild and check if they have the role
	members, _ := ctx.Session.LazyRequestMembers(ctx.Event.GuildID, 1000, true)
	memberCount := 0
	for _, member := range members {
		if utils.SliceContains(member.Roles, role.ID) {
			memberCount++
		}
	}
	fmt.Println("Lazy Request", len(members))
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

func (ModerationCommands) giverole(ctx *dgc.Ctx) {
	if !utils.Discord.CheckIfInServer(ctx) {
		return
	}
	member, success := utils.Discord.TryGetMemberFromMention(ctx, ctx.Arguments.Get(0))
	if !success {
		ctx.RespondText("Could not find user")
		return
	}
	role, success := utils.Discord.TryGetRoleFromMention(ctx, ctx.Arguments.Get(1))
	if !success {
		ctx.RespondText("Could not find role")
		return
	}

	err := ctx.Session.GuildMemberRoleAdd(ctx.Event.GuildID, member.User.ID, role.ID)
	if err != nil {
		ctx.RespondText("Error giving role: " + err.Error())
		return
	}
	embed.NewEmbed().
		SetDescription(fmt.Sprintf("Successfully gave role (%s) to (%s)", role.Name, member.User.String())).
		Send(ctx)
}

func (ModerationCommands) removerole(ctx *dgc.Ctx) {
	if !utils.Discord.CheckIfInServer(ctx) {
		return
	}
	member, success := utils.Discord.TryGetMemberFromMention(ctx, ctx.Arguments.Get(0))
	if !success {
		ctx.RespondText("Could not find user")
		return
	}
	role, success := utils.Discord.TryGetRoleFromMention(ctx, ctx.Arguments.Get(1))
	if !success {
		ctx.RespondText("Could not find role")
		return
	}

	err := ctx.Session.GuildMemberRoleRemove(ctx.Event.GuildID, member.User.ID, role.ID)
	if err != nil {
		ctx.RespondText("Error removing role: " + err.Error())
		return
	}
	embed.NewEmbed().
		SetDescription(fmt.Sprintf("Successfully removed role (%s) from (%s)", role.Name, member.User.String())).
		Send(ctx)
}

func (ModerationCommands) createrole(ctx *dgc.Ctx) {
	if !utils.Discord.CheckIfInServer(ctx) {
		return
	}
	role, err := ctx.Session.GuildRoleCreate(ctx.Event.GuildID, ctx.Arguments.Raw(), 0, false, 0, false)
	if err != nil {
		ctx.RespondText("Error creating role: " + err.Error())
		return
	}
	embed.NewEmbed().
		SetDescription(fmt.Sprintf("Successfully created role (%s)", role.Name)).
		Send(ctx)
}

func (ModerationCommands) deleterole(ctx *dgc.Ctx) {
	if !utils.Discord.CheckIfInServer(ctx) {
		return
	}
	role, success := utils.Discord.TryGetRoleFromMention(ctx, ctx.Arguments.Get(0))
	if !success {
		ctx.RespondText("Could not find role")
		return
	}
	err := ctx.Session.GuildRoleDelete(ctx.Event.GuildID, role.ID)
	if err != nil {
		ctx.RespondText("Error deleting role: " + err.Error())
		return
	}
	embed.NewEmbed().
		SetDescription(fmt.Sprintf("Successfully deleted role (%s)", role.Name)).
		Send(ctx)
}

func (ModerationCommands) renamerole(ctx *dgc.Ctx) {
	if !utils.Discord.CheckIfInServer(ctx) {
		return
	}
	role, success := utils.Discord.TryGetRoleFromMention(ctx, ctx.Arguments.Get(0))
	if !success {
		ctx.RespondText("Could not find role")
		return
	}
	role, err := ctx.Session.GuildRoleEdit(ctx.Event.GuildID, role.ID, utils.String.RemoveWords(ctx.Arguments.Raw(), 1), role.Color, role.Hoist, role.Permissions, role.Mentionable)
	if err != nil {
		ctx.RespondText("Error renaming role: " + err.Error())
		return
	}
	embed.NewEmbed().
		SetDescription(fmt.Sprintf("Successfully renamed role (%s)", role.Name)).
		Send(ctx)
}
