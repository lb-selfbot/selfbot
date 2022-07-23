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
}

// Channel Control

func (ModerationCommands) channelinfo(ctx *dgc.Ctx) {
	channel, success := utils.TryGetChannelFromMention(ctx, ctx.Arguments.Get(0))
	if !success {
		ctx.RespondText("Could not find channel")
		return
	}

	embed.NewEmbed().
		AddField("Name", channel.Name).
		AddField("ID", channel.ID).
		AddField("Topic", channel.Topic).
		AddField("Category ID", channel.ParentID).
		AddField("Is NSFW", utils.FormatBool(channel.NSFW)).
		AddField("Created At", utils.SnowflakeToTime(channel.ID).Format(time.RFC1123)).
		AddField("Slowmode Delay", fmt.Sprint(channel.RateLimitPerUser)+" seconds").
		AddField("Position", fmt.Sprint(channel.Position)).
		AddField("Last Message ID", channel.LastMessageID).
		Send(ctx)
}

func (ModerationCommands) textchannel(ctx *dgc.Ctx) {
	if !utils.CheckIfInServer(ctx) {
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
	if !utils.CheckIfInServer(ctx) {
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
	if !utils.CheckIfInServer(ctx) {
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
	if !utils.CheckIfInServer(ctx) {
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
	if !utils.CheckIfInServer(ctx) {
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
	channel, success := utils.TryGetChannelFromMention(ctx, ctx.Arguments.Get(0))
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
	channel, success := utils.TryGetChannelFromMention(ctx, ctx.Arguments.Get(0))
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
	channel, success := utils.TryGetChannelFromMention(ctx, ctx.Arguments.Get(0))
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
	channel, success := utils.TryGetChannelFromMention(ctx, ctx.Arguments.Get(0))
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
	channel, success := utils.TryGetChannelFromMention(ctx, ctx.Arguments.Get(0))
	if !success {
		ctx.RespondText("Could not find channel")
		return
	}

	_, err := ctx.Session.ChannelEditComplex(channel.ID, &discordgo.ChannelEdit{
		PermissionOverwrites: []*discordgo.PermissionOverwrite{
			{
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
	channel, success := utils.TryGetChannelFromMention(ctx, ctx.Arguments.Get(0))
	if !success {
		ctx.RespondText("Could not find channel")
		return
	}
	guild, err := ctx.Session.Guild(ctx.Event.GuildID)
	if err != nil {
		ctx.RespondText("Error locking channel: " + err.Error())
	}
	// TODO: @everyone
	fmt.Println(guild.Roles[0])

	_, err = ctx.Session.ChannelEditComplex(channel.ID, &discordgo.ChannelEdit{
		PermissionOverwrites: []*discordgo.PermissionOverwrite{
			{
				ID:    guild.Roles[0].ID,
				Type:  discordgo.PermissionOverwriteTypeRole,
				Allow: discordgo.PermissionSendMessages,
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

func (ModerationCommands) category(ctx *dgc.Ctx) {
	if !utils.CheckIfInServer(ctx) {
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
	channel, success := utils.TryGetChannelFromMention(ctx, ctx.Arguments.Get(0))
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

	messages, err := utils.GetChannelMessages(ctx, amount)
	if err != nil {
		ctx.RespondText("Error getting messages: " + err.Error())
		return
	}

	for _, message := range messages {
		ctx.Session.ChannelMessageDelete(ctx.Event.ChannelID, message.ID)
	}
}

func (ModerationCommands) nuke(ctx *dgc.Ctx) {
	channel, success := utils.TryGetChannelFromMention(ctx, ctx.Arguments.Get(0))
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
	user, success := utils.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
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
		AddField("Created At", utils.SnowflakeToTime(user.ID).Format(time.RFC1123)).
		AddField("Flags", fmt.Sprint(user.Flags))
	if member != nil {
		embed.AddField("Nickname", member.Nick)
		embed.AddField("Joined Server At", member.JoinedAt.Format(time.RFC1123))
		embed.AddField("Role Amount", fmt.Sprint(len(member.Roles)))
		embed.AddField("Roles", utils.GetRoleNames(ctx, member))
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
	if !utils.CheckIfInServer(ctx) {
		return
	}
	user, success := utils.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
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

	role, err := ctx.Session.GuildRoleCreate(ctx.Event.GuildID)
	if err != nil {
		ctx.RespondText("Error muting user: " + err.Error())
		return
	}
	role, err = ctx.Session.GuildRoleEdit(ctx.Event.GuildID, role.ID, "Muted", role.Color, role.Hoist, role.Permissions, false)
	if err != nil {
		ctx.RespondText("Error muting user: " + err.Error())
		return
	}

	for _, channel := range channels {
		if channel.Type == discordgo.ChannelTypeGuildText {
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
	if !utils.CheckIfInServer(ctx) {
		return
	}

	member, success := utils.TryGetMemberFromMention(ctx, ctx.Arguments.Get(0))
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
	if !utils.CheckIfInServer(ctx) {
		return
	}

	member, success := utils.TryGetMemberFromMention(ctx, ctx.Arguments.Get(0))
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
	if !utils.CheckIfInServer(ctx) {
		return
	}

	member, success := utils.TryGetMemberFromMention(ctx, ctx.Arguments.Get(0))
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
	if !utils.CheckIfInServer(ctx) {
		return
	}

	member, success := utils.TryGetMemberFromMention(ctx, ctx.Arguments.Get(0))
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

func (ModerationCommands) ban(ctx *dgc.Ctx) {
	if !utils.CheckIfInServer(ctx) {
		return
	}

	member, success := utils.TryGetMemberFromMention(ctx, ctx.Arguments.Get(0))
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
	if !utils.CheckIfInServer(ctx) {
		return
	}

	member, success := utils.TryGetMemberFromMention(ctx, ctx.Arguments.Get(0))
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
	if !utils.CheckIfInServer(ctx) {
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
