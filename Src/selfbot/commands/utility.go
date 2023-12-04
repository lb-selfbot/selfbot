package commands

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go/token"
	"go/types"
	"image/png"
	"math"
	nethttp "net/http"
	"net/url"
	"os"
	"os/exec"
	"path"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/LightningDev1/LB-Selfbot/selfbot/eventhandlers/eventreactions"
	"github.com/LightningDev1/LB-Selfbot/selfbot/ipc"
	"github.com/LightningDev1/LB-Selfbot/selfbot/logging"
	"github.com/LightningDev1/LB-Selfbot/selfbot/uploaders/destinations"
	"github.com/LightningDev1/screenshot"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"

	"github.com/LightningDev1/LB-Selfbot/selfbot/api"
	"github.com/LightningDev1/LB-Selfbot/selfbot/constants"
	"github.com/LightningDev1/LB-Selfbot/selfbot/embed"
	"github.com/LightningDev1/LB-Selfbot/selfbot/http"
	"github.com/LightningDev1/LB-Selfbot/selfbot/utils"
	"github.com/LightningDev1/dgc"
	"github.com/LightningDev1/discordgo"
)

type UtilityCommands struct{}

func (u UtilityCommands) RegisterCommands(router *dgc.Router) {
	router.StartCategory("Utility", "Utility commands")
	router.RegisterCmd(&dgc.Command{
		Name:        "copyserver",
		Description: "Make an exact copy of a server",
		Usage:       "[p]copyserver [server id] [delay]",
		Handler:     u.copyserver,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "copychannel",
		Description: "Make an exact copy of a channel",
		Usage:       "[p]copychannel <source channel> <dest channel> [limit]",
		Handler:     u.copychannel,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "leaveservers",
		Description: "Leave all of your servers",
		Usage:       "[p]leaveservers",
		Handler:     u.leaveservers,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "leavegroups",
		Description: "Leave all of your groups",
		Usage:       "[p]leavegroups [silent: on/off]",
		Handler:     u.leavegroups,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "unfriendall",
		Description: "Remove all of your friends",
		Usage:       "[p]unfriendall",
		Handler:     u.unfriendall,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "tts",
		Description: "Send a Text-To-Speech message",
		Usage:       "[p]tts <language> <text>",
		Handler:     u.tts,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "translate",
		Description: "Translate text to a language",
		Usage:       "[p]translate <language> <text>",
		Handler:     u.translate,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "detectlang",
		Description: "Detect a texts language",
		Usage:       "[p]detectlang <text>",
		Handler:     u.detectlang,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "english",
		Description: "Translate text to English",
		Usage:       "[p]english <text>",
		Handler:     u.english,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "adminservers",
		Description: "See all servers you have administator in",
		Usage:       "[p]adminservers",
		Handler:     u.adminservers,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "admingroups",
		Description: "See all groupds you have owner in",
		Usage:       "[p]admingroups",
		Handler:     u.admingroups,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "joinvc",
		Description: "Join a voice channel",
		Usage:       "[p]joinvc <channel id>",
		Handler:     u.joinvc,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "optionpoll",
		Description: "Create a poll with options",
		Usage:       "[p]optionpoll <\"question\"> <\"option1\"> <\"option2\">",
		Handler:     u.optionpoll,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "systeminfo",
		Description: "See info about your system",
		Usage:       "[p]systeminfo",
		Handler:     u.systeminfo,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "findserver",
		Description: "Find a server by ID",
		Usage:       "[p]findserver <server id>",
		Handler:     u.findserver,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "finduser",
		Description: "Find a user by ID",
		Usage:       "[p]finduser <server id>",
		Handler:     u.finduser,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "passgen",
		Description: "Generate a password",
		Usage:       "[p]passgen <character amount>",
		Handler:     u.passgen,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "tokengen",
		Description: "Generate a Discord token",
		Usage:       "[p]tokengen",
		Handler:     u.tokengen,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "nickscan",
		Description: "Search for servers with a nickname",
		Usage:       "[p]nickscan",
		Handler:     u.nickscan,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "shutdownpc",
		Description: "Shut your PC down",
		Usage:       "[p]shutdownpc",
		Flags:       []string{"DANGEROUS", "HOSTING_DANGEROUS"},
		Handler:     u.shutdownpc,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "restartpc",
		Description: "Restart your PC",
		Usage:       "[p]restartpc",
		Flags:       []string{"DANGEROUS", "HOSTING_DANGEROUS"},
		Handler:     u.restartpc,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "sleeppc",
		Description: "Put your PC to sleep",
		Usage:       "[p]sleeppc",
		Flags:       []string{"DANGEROUS", "HOSTING_DANGEROUS"},
		Handler:     u.sleeppc,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "lockpc",
		Description: "Lock your PC",
		Usage:       "[p]lockpc",
		Flags:       []string{"DANGEROUS", "HOSTING_DANGEROUS"},
		Handler:     u.lockpc,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "snipe",
		Description: "See the last deleted message in a channel",
		Usage:       "[p]snipe [channel]",
		Handler:     u.snipe,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "editsnipe",
		Description: "See the last edited message in a channel",
		Usage:       "[p]editsnipe [channel]",
		Handler:     u.editsnipe,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "history",
		Description: "See used command history",
		Usage:       "[p]history",
		Handler:     u.history,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "website",
		Description: "Open the LightningBot website",
		Usage:       "[p]website",
		Flags:       []string{"HOSTING_DANGEROUS"},
		Handler:     u.website,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "dashboard",
		Description: "Opens dashboard in browser",
		Usage:       "[p]dashboard",
		Flags:       []string{"HOSTING_DANGEROUS"},
		Handler:     u.dashboard,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "restartbot",
		Description: "Restart LightningBot",
		Usage:       "[p]restartbot",
		Flags:       []string{"DANGEROUS"},
		Handler:     u.restartbot,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "stopbot",
		Description: "Stop LightningBot",
		Usage:       "[p]stopbot",
		Flags:       []string{"DANGEROUS"},
		Handler:     u.stopbot,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "changelog",
		Description: "See the latest changelog",
		Usage:       "[p]changelog",
		Handler:     u.changelog,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "version",
		Description: "See the current LightningBot version",
		Usage:       "[p]version",
		Handler:     u.version,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "uptime",
		Description: "See how long the bot has been running for",
		Usage:       "[p]uptime",
		Handler:     u.uptime,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "cls",
		Description: "Clear UI console",
		Usage:       "[p]cls",
		Flags:       []string{"DANGEROUS"},
		Handler:     u.cls,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "openfolder",
		Description: "Open the LightningBot folder",
		Usage:       "[p]openfolder",
		Flags:       []string{"HOSTING_DANGEROUS"},
		Handler:     u.openfolder,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "mastercard",
		Description: "Generate a mastercard credit card",
		Usage:       "[p]mastercard [name]",
		Handler:     u.mastercard,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "visa16",
		Description: "Generate a visa credit card (16 characters)",
		Usage:       "[p]visa16 [name]",
		Handler:     u.visa16,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "visa13",
		Description: "Generate a visa credit card (13 characters)",
		Usage:       "[p]visa13 [name]",
		Handler:     u.visa13,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "amex",
		Description: "Generate an amex credit card",
		Usage:       "[p]amex [name]",
		Handler:     u.amex,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "discover",
		Description: "Generate a discover credit card",
		Usage:       "[p]discover [name]",
		Handler:     u.discover,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "diners",
		Description: "Generate a diners credit card",
		Usage:       "[p]diners [name]",
		Handler:     u.diners,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "jcb",
		Description: "Generate a JCB credit card",
		Usage:       "[p]jcb [name]",
		Handler:     u.jcb,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "revavatarsearch",
		Description: "Reverse image search a users avatar",
		Usage:       "[p]revavatarsearch <user>",
		Handler:     u.revavatarsearch,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "calc",
		Description: "Open the calculator",
		Usage:       "[p]calc",
		Flags:       []string{"DANGEROUS", "HOSTING_DANGEROUS"},
		Handler:     u.calc,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "notepad",
		Description: "Open notepad",
		Usage:       "[p]notepad",
		Flags:       []string{"DANGEROUS", "HOSTING_DANGEROUS"},
		Handler:     u.notepad,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "firstmessage",
		Description: "See the first message in the channel",
		Usage:       "[p]firstmessage [channel]",
		Handler:     u.firstmessage,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "hypesquad",
		Description: "Change your Hypesquad House",
		Usage:       "[p]hypesquad <bravery/brilliance/balance>",
		Handler:     u.hypesquad,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "leavehypesquad",
		Description: "Leave Hypesquad",
		Usage:       "[p]leavehypesquad",
		Handler:     u.leavehypesquad,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "theme",
		Description: "Change your Discord theme",
		Usage:       "[p]theme <light/dark>",
		Handler:     u.theme,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "compact",
		Description: "Toggle Discord compact mode",
		Usage:       "[p]compact <on/off>",
		Handler:     u.compact,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "devmode",
		Description: "Toggle Discord developer mode",
		Usage:       "[p]devmode <on/off>",
		Handler:     u.devmode,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "sessionid",
		Description: "See LightningBot's Discord session ID",
		Usage:       "[p]sessionid",
		Handler:     u.sessionid,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "aboutme",
		Description: "Change your About Me text",
		Usage:       "[p]aboutme <text>",
		Handler:     u.aboutme,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "clearaboutme",
		Description: "Clear your About Me text",
		Usage:       "[p]clearaboutme",
		Handler:     u.clearaboutme,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "displayname",
		Description: "Change your display name",
		Usage:       "[p]displayname <text>",
		Handler:     u.displayname,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "cleardisplayname",
		Description: "Clear your display name",
		Usage:       "[p]cleardisplayname",
		Handler:     u.cleardisplayname,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "pronouns",
		Description: "Change your pronouns",
		Usage:       "[p]pronouns <text>",
		Handler:     u.pronouns,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "clearpronouns",
		Description: "Clear your pronouns",
		Usage:       "[p]clearpronouns",
		Handler:     u.clearpronouns,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "timezones",
		Description: "Get a list of timezones for a continent",
		Usage:       "[p]timezones <continent>",
		Handler:     u.timezones,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "timezone",
		Description: "Get information about a timezone",
		Usage:       "[p]timezone <timezone>",
		Handler:     u.timezone,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "getmessagejson",
		Description: "Get a message's JSON data",
		Usage:       "[p]getmessagejson <#channel> <message id>",
		Handler:     u.getmessagejson,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "getembedjson",
		Description: "Get an embed's JSON data",
		Usage:       "[p]getembedjson <#channel> <message id>",
		Handler:     u.getembedjson,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "calculate",
		Description: "Calculate something",
		Usage:       "[p]calculate <expression>",
		Handler:     u.calculate,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "sendinterval",
		Description: "Send a message every x seconds (flags: DELETE, INSTANT_DEL)",
		Usage:       "[p]sendinterval <interval> <amount> <\"message\"> [flags]",
		Handler:     u.sendinterval,
	})

	// Uploading
	router.RegisterCmd(&dgc.Command{
		Name:        "screenshot",
		Description: "Make a screenshot",
		Usage:       "[p]screenshot [delay]",
		Flags:       []string{"DANGEROUS", "HOSTING_DANGEROUS"},
		Handler:     u.screenshot,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "listdestinations",
		Description: "List all upload destinations",
		Usage:       "[p]listdestinations",
		Handler:     u.listdestinations,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "uploadss",
		Description: "Upload a screenshot",
		Usage:       "[p]uploadss <destination> [delay]",
		Flags:       []string{"DANGEROUS", "HOSTING_DANGEROUS"},
		Handler:     u.uploadss,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "uploadfile",
		Description: "Upload a file",
		Usage:       "[p]uploadfile <destination> [attachment]",
		Handler:     u.uploadfile,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "uploadlink",
		Description: "Upload a link",
		Usage:       "[p]uploadlink <destination> <link>",
		Handler:     u.uploadlink,
	})

	// Obfuscation
	router.RegisterCmd(&dgc.Command{
		Name:        "obfuscatepython",
		Description: "Obfuscate a Python script",
		Usage:       "[p]obfuscatepython [script/attachment]",
		Handler:     u.obfuscateFunc("python", "py"),
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "obfuscatelua",
		Description: "Obfuscate a Lua script",
		Usage:       "[p]obfuscatelua [script/attachment]",
		Handler:     u.obfuscateFunc("lua", "lua"),
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "obfuscatejs",
		Description: "Obfuscate a JS script",
		Usage:       "[p]obfuscatejs [script/attachment]",
		Handler:     u.obfuscateFunc("javascript", "js"),
	})

	// ReviewDB
	router.RegisterCmd(&dgc.Command{
		Name:        "reviewdbtoken",
		Description: "Get your ReviewDB token",
		Usage:       "[p]reviewdbtoken",
		Handler:     u.reviewdbtoken,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "getreviews",
		Description: "Get a users ReviewDB reviews",
		Usage:       "[p]getreviews <user>",
		Handler:     u.getreviews,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "addreview",
		Description: "Add a ReviewDB review to a user",
		Usage:       "[p]addreview <user> <review>",
		Handler:     u.addreview,
	})
}

func getChannelType(channelType discordgo.ChannelType) discordgo.ChannelType {
	if channelType == discordgo.ChannelTypeGuildVoice || channelType == discordgo.ChannelTypeGuildStageVoice {
		return discordgo.ChannelTypeGuildVoice
	}
	return discordgo.ChannelTypeGuildText
}

func (UtilityCommands) copyserver(ctx *dgc.Ctx) {
	duration, _ := time.ParseDuration(utils.Arguments["[delay]"])

	guild, err := utils.Discord.GetGuild(ctx.Session, utils.Misc.Or(utils.Arguments["[server id]"], ctx.Event.GuildID))
	if err != nil {
		_ = ctx.RespondText("Error getting server: " + err.Error())
		return
	}

	icon := http.Get(guild.IconURL()).Body

	newGuild, err := api.CreateGuild(ctx.Session, fmt.Sprintf("%s - Copy", guild.Name), icon, "") //ctx.Session.GuildCreate(fmt.Sprintf("%s - Copy", guild.Name))
	if err != nil {
		_ = ctx.RespondText("Error copying server: " + err.Error())
		return
	}
	_, _ = ctx.Session.GuildEdit(newGuild.ID, &discordgo.GuildParams{
		VerificationLevel:           &guild.VerificationLevel,
		DefaultMessageNotifications: int(guild.DefaultMessageNotifications),
		AfkChannelID:                guild.AfkChannelID,
		AfkTimeout:                  guild.AfkTimeout,
		Icon:                        icon,
	})

	IDs := map[string]string{guild.ID: newGuild.ID}

	for _, role := range guild.Roles {
		if role.ID == guild.ID || role.Managed {
			continue
		}

		newRole, err := ctx.Session.GuildRoleCreate(newGuild.ID, &discordgo.RoleParams{
			Name:        role.Name,
			Color:       &role.Color,
			Hoist:       &role.Hoist,
			Permissions: &role.Permissions,
			Mentionable: &role.Mentionable,
		})
		if err != nil {
			_ = ctx.RespondText("Error copying server: " + err.Error())
			return
		}

		IDs[role.ID] = newRole.ID

		time.Sleep(duration)
	}

	newChannels, err := ctx.Session.GuildChannels(newGuild.ID)
	if err != nil {
		_ = ctx.RespondText("Error copying server: " + err.Error())
		return
	}

	for _, channel := range newChannels {
		_, err := ctx.Session.ChannelDelete(channel.ID)
		if err != nil {
			_ = ctx.RespondText("Error copying server: " + err.Error())
			return
		}

		time.Sleep(duration)
	}

	channels, err := ctx.Session.GuildChannels(guild.ID)
	if err != nil {
		_ = ctx.RespondText("Error copying server: " + err.Error())
		return
	}

	categories, categoryChannels, looseChannels := utils.Discord.MapCategories(channels)

	for _, channel := range looseChannels {
		overwrites := utils.Discord.MapPermissionOverwrites(channel.PermissionOverwrites, IDs)

		_, err = ctx.Session.GuildChannelCreateComplex(newGuild.ID, discordgo.GuildChannelCreateData{
			Name:                 channel.Name,
			Type:                 getChannelType(channel.Type),
			Topic:                channel.Topic,
			Bitrate:              int(math.Min(float64(channel.Bitrate), 96000)),
			UserLimit:            int(math.Min(float64(channel.UserLimit), 99)),
			RateLimitPerUser:     channel.RateLimitPerUser,
			Position:             channel.Position,
			PermissionOverwrites: overwrites,
			NSFW:                 channel.NSFW,
		})
		if err != nil {
			_ = ctx.RespondText("Error copying server: " + err.Error())
			return
		}

		time.Sleep(duration)
	}

	for _, category := range categories {
		overwrites := utils.Discord.MapPermissionOverwrites(category.PermissionOverwrites, IDs)

		newCategory, err := ctx.Session.GuildChannelCreateComplex(newGuild.ID, discordgo.GuildChannelCreateData{
			Name:                 category.Name,
			Type:                 discordgo.ChannelTypeGuildCategory,
			Position:             category.Position,
			PermissionOverwrites: overwrites,
		})
		if err != nil {
			_ = ctx.RespondText("Error copying server: " + err.Error())
			return
		}

		for _, channel := range categoryChannels[category.ID] {
			overwrites := utils.Discord.MapPermissionOverwrites(channel.PermissionOverwrites, IDs)

			_, err = ctx.Session.GuildChannelCreateComplex(newGuild.ID, discordgo.GuildChannelCreateData{
				Name:                 channel.Name,
				Type:                 getChannelType(channel.Type),
				Topic:                channel.Topic,
				Bitrate:              int(math.Min(float64(channel.Bitrate), 96000)),
				UserLimit:            channel.UserLimit,
				RateLimitPerUser:     channel.RateLimitPerUser,
				Position:             channel.Position,
				PermissionOverwrites: overwrites,
				ParentID:             newCategory.ID,
				NSFW:                 channel.NSFW,
			})
			if err != nil {
				_ = ctx.RespondText("Error copying server: " + err.Error())
				return
			}

			time.Sleep(duration)
		}
	}
	embed.NewEmbed().
		SetDescription("Server copied successfully!").
		Send(ctx)
}

func (UtilityCommands) copychannel(ctx *dgc.Ctx) {
	sourceChannel, success := utils.Discord.TryGetChannelFromMention(ctx, ctx.Arguments.Get(0))
	if !success {
		_ = ctx.RespondText("Could not find channel")
		return
	}

	destChannel, success := utils.Discord.TryGetChannelFromMention(ctx, ctx.Arguments.Get(1))
	if !success {
		_ = ctx.RespondText("Could not find channel")
		return
	}

	limit, err := ctx.Arguments.Get(2).AsInt()
	if err != nil {
		limit = 100
	}

	if limit < 1 {
		_ = ctx.RespondText("Limit must be greater than 0")
		return
	}

	messages, err := utils.Discord.GetChannelMessages(ctx.Session, sourceChannel.ID, limit)
	if err != nil {
		_ = ctx.RespondText("Error getting source messages: " + err.Error())
		return
	}

	sort.Slice(messages, func(i, j int) bool {
		timeI := utils.Discord.SnowflakeToTime(messages[i].ID)
		timeJ := utils.Discord.SnowflakeToTime(messages[j].ID)

		return timeI.Before(timeJ)
	})

	_ = ctx.RespondText(fmt.Sprintf(
		"Copying %d messages...\n\nEstimated time required: %s",
		len(messages),
		time.Duration(len(messages)*2)*time.Second,
	))

	webhook, err := ctx.Session.WebhookCreate(destChannel.ID, "LightningBot - Copy Channel", "")
	if err != nil {
		_ = ctx.RespondText("Error creating destination webhook: " + err.Error())
		return
	}

	requests := 0
	rateLimit := 30
	start := time.Now()

	for _, message := range messages {
		attachmentList := []string{}

		for _, attachment := range message.Attachments {
			attachmentList = append(attachmentList, attachment.URL)
		}

		attachments := "\n" + strings.Join(attachmentList, "\n")

		maxLength := 2000 - len(attachments)
		length := int(math.Min(float64(len(message.Content)), float64(maxLength)))

		content := message.Content[:length] + "\n" + attachments

		_, err := ctx.Session.WebhookExecute(webhook.ID, webhook.Token, false, &discordgo.WebhookParams{
			Content:    content,
			Username:   message.Author.Username,
			AvatarURL:  message.Author.AvatarURL(""),
			TTS:        message.TTS,
			Components: message.Components,
			Embeds:     message.Embeds,
		})
		if err != nil {
			_ = ctx.RespondText("Error sending message: " + err.Error())
		}

		time.Sleep(750 * time.Millisecond)

		requests++

		if requests >= rateLimit && time.Since(start).Seconds() < 60 {
			time.Sleep(time.Duration(60-time.Since(start).Seconds()) * time.Second)
			requests = 0
			start = time.Now()
		}
	}
}

func (UtilityCommands) leaveservers(ctx *dgc.Ctx) {
	for _, guild := range ctx.Session.State.Guilds {
		err := ctx.Session.GuildLeave(guild.ID)
		if err != nil {
			_ = ctx.RespondText("Error leaving server: " + err.Error())
		}

		time.Sleep(2 * time.Second)
	}
}

func (UtilityCommands) leavegroups(ctx *dgc.Ctx) {
	rawSilentToggle := strings.ToLower(utils.Arguments["[silent leave: on/off]"])
	silentToggle := rawSilentToggle == "on"

	for _, channel := range ctx.Session.State.PrivateChannels {
		if channel.Type == discordgo.ChannelTypeGroupDM {
			_, err := ctx.Session.ChannelDelete(channel.ID, silentToggle)
			if err != nil {
				_ = ctx.RespondText("Error leaving group: " + err.Error())
			}

			time.Sleep(2 * time.Second)
		}
	}
}

func (UtilityCommands) unfriendall(ctx *dgc.Ctx) {
	rels, err := ctx.Session.UserRelationships()
	if err != nil {
		_ = ctx.RespondText("Error getting relationships: " + err.Error())
		return
	}

	for _, rel := range rels {
		err = ctx.Session.UserRelationshipDelete(rel.User.ID)
		if err != nil {
			_ = ctx.RespondText("Error removing friend: " + err.Error())
			return
		}

		time.Sleep(2 * time.Second)
	}
}

func (UtilityCommands) tts(ctx *dgc.Ctx) {
	httpResponse := api.TTS(utils.String.RemoveWords(ctx.Arguments.Raw(), 1), utils.Arguments["<language>"])
	if httpResponse.Error != nil {
		_ = ctx.RespondText("An error has occurred.")
		return
	}

	if strings.Contains(httpResponse.Body, "Bad Request") {
		_ = ctx.RespondText("Invalid language. Please use an ISO 639-1 language code. See https://en.wikipedia.org/wiki/List_of_ISO_639-1_codes for more information.")
		return
	}

	_, _ = ctx.Session.ChannelFileSend(
		ctx.Event.ChannelID,
		"tts.mp3",
		bytes.NewReader(httpResponse.BodyBytes),
	)
}

func (UtilityCommands) translate(ctx *dgc.Ctx) {
	translateResult := api.Translate(utils.String.RemoveWords(ctx.Arguments.Raw(), 1), utils.Arguments["<language>"])
	if translateResult.Error != nil {
		_ = ctx.RespondText("An error has occurred: " + translateResult.Error.Error())
		return
	}
	embed.NewEmbed().
		AddField("Source Language", utils.String.Capitalize(constants.LANGUAGES[translateResult.SourceLanguage])).
		AddField("Destination Language", utils.String.Capitalize(constants.LANGUAGES[translateResult.DestinationLanguage])).
		AddField("Original Text", translateResult.Original).
		AddField("Translated Text", translateResult.Result).
		Send(ctx)
}

func (UtilityCommands) detectlang(ctx *dgc.Ctx) {
	translateResult := api.Translate(ctx.Arguments.Raw(), "en")
	if translateResult.Error != nil {
		_ = ctx.RespondText("An error has occurred: " + translateResult.Error.Error())
		return
	}
	embed.NewEmbed().
		AddField("Text", translateResult.Original).
		AddField("Detected Language", utils.String.Capitalize(constants.LANGUAGES[translateResult.SourceLanguage])).
		AddField("Confidence", fmt.Sprintf("%.2f", translateResult.Confidence)).
		Send(ctx)
}

func (UtilityCommands) english(ctx *dgc.Ctx) {
	translateResult := api.Translate(ctx.Arguments.Raw(), "en")
	if translateResult.Error != nil {
		_ = ctx.RespondText("An error has occurred: " + translateResult.Error.Error())
		return
	}
	embed.NewEmbed().
		AddField("Source Language", utils.String.Capitalize(constants.LANGUAGES[translateResult.SourceLanguage])).
		AddField("Original Text", translateResult.Original).
		AddField("Translated Text", translateResult.Result).
		Send(ctx)
}

func (UtilityCommands) adminservers(ctx *dgc.Ctx) {
	userId, _ := strconv.Atoi(ctx.Session.State.User.ID)

	adminServerString := ""
	for _, guild := range ctx.Session.State.Guilds {
		if len(guild.Channels) < 1 {
			continue
		}

		_, err := ctx.Session.State.Member(guild.ID, ctx.Session.State.User.ID)
		if err != nil {
			params := discordgo.NewQueryGuildMembersParams(guild.ID)
			params.UserIDs = []int{userId}

			_, _ = ctx.Session.QueryGuildMembers(params)
		}

		permissions, err := ctx.Session.State.UserChannelPermissions(ctx.Session.State.User.ID, guild.Channels[0].ID)
		if err != nil {
			continue
		}

		if permissions&discordgo.PermissionAdministrator != 0 {
			adminServerString += guild.Name + "\n"
		}
	}

	embed.NewEmbed().
		AddField("Admin Servers", adminServerString).
		Send(ctx)
}

func (UtilityCommands) admingroups(ctx *dgc.Ctx) {
	adminGroupsString := ""
	for _, channel := range ctx.Session.State.PrivateChannels {
		if channel.Type == discordgo.ChannelTypeGroupDM {
			if channel.OwnerID == ctx.Session.State.User.ID {
				adminGroupsString += channel.Name + "\n"
			}
		}
	}
	embed.NewEmbed().
		AddField("Admin Groups", adminGroupsString).
		Send(ctx)
}

func (UtilityCommands) joinvc(ctx *dgc.Ctx) {
	channel, success := utils.Discord.TryGetChannelFromMention(ctx, ctx.Arguments.Get(0))
	if !success {
		_ = ctx.RespondText("Could not find channel")
		return
	}

	_, err := ctx.Session.ChannelVoiceJoin(channel.GuildID, channel.ID, false, false)
	if err != nil {
		_ = ctx.RespondText("Could not join channel")
		return
	}

	_ = ctx.RespondText("Joined channel " + channel.Name)
}

func (UtilityCommands) optionpoll(ctx *dgc.Ctx) {
	messages := embed.NewEmbed().
		SetDescription("Poll: "+utils.Arguments["<\"question\">"]).
		AddField("👍", utils.Arguments["<\"option1\">"]).
		AddField("👎", utils.Arguments["<\"option2\">"]).
		Send(ctx)

	if len(messages) != 0 {
		_ = ctx.Session.MessageReactionAdd(ctx.Event.ChannelID, messages[0].ID, "👍")
		_ = ctx.Session.MessageReactionAdd(ctx.Event.ChannelID, messages[0].ID, "👎")
	}
}

func (UtilityCommands) systeminfo(ctx *dgc.Ctx) {
	path := "/"
	if runtime.GOOS == "windows" {
		path = "C:"
	}
	hostStat, _ := host.Info()
	cpuStat, _ := cpu.Info()
	cpuInfo := cpuStat[0]
	vmStat, _ := mem.VirtualMemory()
	diskStat, _ := disk.Usage(path)
	embed.NewEmbed().
		SetDescription("System Information").
		AddField("OS", utils.String.Capitalize(hostStat.OS)).
		AddField("Platform", fmt.Sprintf("%s - %s", hostStat.Platform, hostStat.PlatformVersion)).
		AddField("Kernel Arch and Version", fmt.Sprintf("%s - %s", hostStat.KernelArch, hostStat.KernelVersion)).
		AddField("CPU Info", fmt.Sprintf("Model: %s\nFrequency: %.f MHz\nCores: %d\nVendor: %s", cpuInfo.ModelName, cpuInfo.Mhz, cpuInfo.Cores, cpuInfo.VendorID)).
		AddField("Memory Info", fmt.Sprintf("Total: %d MB\nFree: %d MB\nUsed: %d MB\nUsed Percentage: %.2f%%", vmStat.Total/1024/1024, vmStat.Free/1024/1024, vmStat.Used/1024/1024, vmStat.UsedPercent)).
		AddField("Disk Info", fmt.Sprintf("Total: %d MB\nFree: %d MB\nUsed: %d MB\nUsed Percentage: %.2f%%", diskStat.Total/1024/1024, diskStat.Free/1024/1024, diskStat.Used/1024/1024, diskStat.UsedPercent)).
		Send(ctx)
}

func (UtilityCommands) findserver(ctx *dgc.Ctx) {
	guild, err := utils.Discord.GetGuild(ctx.Session, utils.Misc.Or(utils.Arguments["<server id>"], ctx.Event.GuildID))
	if err != nil {
		_ = ctx.RespondText("Error getting server: " + err.Error())
		return
	}
	embed.NewEmbed().
		AddField("Server Name", guild.Name).
		AddField("Created At", utils.Discord.SnowflakeToTime(guild.ID).Format("2006-01-02 15:04:05")).
		Send(ctx)
}

func (UtilityCommands) finduser(ctx *dgc.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		_ = ctx.RespondText("Could not find user")
		return
	}
	embed.NewEmbed().
		AddField("User", user.String()).
		AddField("Created At", utils.Discord.SnowflakeToTime(user.ID).Format("2006-01-02 15:04:05")).
		Send(ctx)
}

func (UtilityCommands) passgen(ctx *dgc.Ctx) {
	charAmount, err := strconv.Atoi(utils.Arguments["<character amount>"])
	if err != nil {
		_ = ctx.RespondText("Argument <character amount> must be a number.")
		return
	}
	embed.NewEmbed().
		AddField("Password", utils.Random.String(charAmount, "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789#$%&")).
		Send(ctx)
}

func (UtilityCommands) tokengen(ctx *dgc.Ctx) {
	httpResponse := http.Get("https://some-random-api.com/bottoken")
	if httpResponse.Error != nil {
		_ = ctx.RespondText("Error: " + httpResponse.Error.Error())
		return
	}

	token := httpResponse.Json().Get("token").MustString()

	embed.NewEmbed().
		AddField("Token", token).
		Send(ctx)
}

func (UtilityCommands) nickscan(ctx *dgc.Ctx) {
	nickText := ""
	for _, guild := range ctx.Session.State.Guilds {
		me, err := ctx.Session.State.Member(guild.ID, ctx.Event.Author.ID)
		if err != nil {
			me, err = ctx.Session.GuildMember(guild.ID, ctx.Event.Author.ID)
			if err != nil {
				continue
			}
		}
		if me.Nick != "" {
			nickText += fmt.Sprintf("%s: %s\n", guild.Name, me.Nick)
		}
	}
	embed.NewEmbed().
		AddField("Nicknames", nickText).
		Send(ctx)
}

func (UtilityCommands) shutdownpc(ctx *dgc.Ctx) {
	_ = utils.System.ShutdownPC()
}

func (UtilityCommands) restartpc(ctx *dgc.Ctx) {
	_ = utils.System.RestartPC()
}

func (UtilityCommands) sleeppc(ctx *dgc.Ctx) {
	_ = utils.System.SleepPC()
}

func (UtilityCommands) lockpc(ctx *dgc.Ctx) {
	_ = utils.System.LockPC()
}

func (UtilityCommands) snipe(ctx *dgc.Ctx) {
	channel, success := utils.Discord.TryGetChannelFromMention(ctx, ctx.Arguments.Get(0))
	if !success {
		channel, _ = ctx.Session.State.Channel(ctx.Event.ChannelID)
		if channel == nil {
			_ = ctx.RespondText("Could not find channel")
			return
		}
	}

	messages := eventreactions.GetDeletedMessages(channel.ID)
	if len(messages) == 0 {
		embed.NewEmbed().
			SetDescription("No deleted message was found in this channel!").
			Send(ctx)
	} else {
		message := messages[len(messages)-1]
		emb := embed.NewEmbed()
		if message.Author != nil {
			emb.AddField("Author", message.Author.String()+" | "+message.Author.ID)
		}
		emb.AddField("Sent At", utils.Discord.SnowflakeToTime(message.ID).Format(time.RFC1123))
		emb.AddField("Content", message.Content)
		emb.Send(ctx)
	}
}

func (UtilityCommands) editsnipe(ctx *dgc.Ctx) {
	channel, success := utils.Discord.TryGetChannelFromMention(ctx, ctx.Arguments.Get(0))
	if !success {
		channel, _ = ctx.Session.State.Channel(ctx.Event.ChannelID)
		if channel == nil {
			_ = ctx.RespondText("Could not find channel")
			return
		}
	}

	messages := eventreactions.GetEditedMessages(channel.ID)
	if len(messages) == 0 {
		embed.NewEmbed().
			SetDescription("No deleted message was found in this channel!").
			Send(ctx)
	} else {
		message := messages[len(messages)-1]
		emb := embed.NewEmbed()
		if message.Author != nil {
			emb.AddField("Author", message.Author.String()+" | "+message.Author.ID)
		}
		emb.AddField("Sent At", utils.Discord.SnowflakeToTime(message.ID).Format(time.RFC1123))
		emb.AddField("Old Content", message.BeforeUpdate.Content)
		emb.AddField("New Content", message.Content)
		emb.Send(ctx)
	}
}

func (UtilityCommands) history(ctx *dgc.Ctx) {
	eventreactions.UsedCommands.Mutex.Lock()
	defer eventreactions.UsedCommands.Mutex.Unlock()

	usedCommands := ""
	index := 1
	for i := len(eventreactions.UsedCommands.Items) - 1; i >= 0 && i <= 10; i-- {
		cmd := eventreactions.UsedCommands.Items[i]
		usedCommands += fmt.Sprintf("%d. %s\n", index, cmd.Event.Content)
		index++
	}
	embed.NewEmbed().
		AddField("Used Commands", usedCommands).
		Send(ctx)
}

func (UtilityCommands) website(*dgc.Ctx) {
	_ = utils.System.OpenURL("https://lightning-bot.com/")
}

func (UtilityCommands) dashboard(*dgc.Ctx) {
	_ = utils.System.OpenURL("https://lightning-bot.com/dashboard")
}

func (UtilityCommands) restartbot(*dgc.Ctx) {
	if constants.CONSOLE {
		utils.System.Restart()
	} else {
		ipc.IPC.Restart()
	}
}

func (UtilityCommands) stopbot(*dgc.Ctx) {
	os.Exit(0)
}

func (UtilityCommands) changelog(ctx *dgc.Ctx) {
	httpResponse := http.PostJson(
		"https://api.lightning-bot.com/api/v3/selfbot/check-version",
		map[string]string{"branch": "selfbot", "version": ""},
	)
	if httpResponse.Error != nil {
		_ = ctx.RespondText("An error has occurred: " + httpResponse.Error.Error())
		return
	}

	json := httpResponse.Json()

	latestVersion := json.GetPath("latest", "version").MustString()
	changelog := json.GetPath("latest", "changelog").MustString()

	embed.NewEmbed().
		AddField("Latest Version", latestVersion).
		AddField("Changelog", changelog).
		Send(ctx)
}

func (UtilityCommands) version(ctx *dgc.Ctx) {
	httpResponse := http.PostJson(
		"https://api.lightning-bot.com/api/v3/selfbot/check-version",
		map[string]string{"branch": "selfbot", "version": ""},
	)
	if httpResponse.Error != nil {
		_ = ctx.RespondText("An error has occurred: " + httpResponse.Error.Error())
		return
	}

	json := httpResponse.Json()

	latestVersion := json.GetPath("latest", "version").MustString()

	embed.NewEmbed().
		AddField("Current Version", constants.VERSION[1:]).
		AddField("Latest Version", latestVersion).
		Send(ctx)
}

func (UtilityCommands) uptime(ctx *dgc.Ctx) {
	embed.NewEmbed().
		AddField("Uptime", utils.Misc.FormatDuration(time.Since(constants.START_TIME))).
		Send(ctx)
}

func (UtilityCommands) cls(*dgc.Ctx) {
	logging.UI.Clear()
}

func (UtilityCommands) openfolder(*dgc.Ctx) {
	_ = utils.System.OpenURL(utils.IO.GetMainFolder())
}

func (UtilityCommands) mastercard(ctx *dgc.Ctx) {
	creditCard := api.GenerateCreditCard("mastercard", ctx.Arguments.Raw())
	embed.NewEmbed().
		AddField("Card Type", "Mastercard").
		AddField("Card Number", creditCard.Number).
		AddField("CVV", creditCard.CVV).
		AddField("Expiry Date", creditCard.ExpiryDate).
		SetImage(creditCard.Visualisation).
		Send(ctx)
}

func (UtilityCommands) visa16(ctx *dgc.Ctx) {
	creditCard := api.GenerateCreditCard("visa16", ctx.Arguments.Raw())
	embed.NewEmbed().
		AddField("Card Type", "Visa").
		AddField("Card Number", creditCard.Number).
		AddField("CVV", creditCard.CVV).
		AddField("Expiry Date", creditCard.ExpiryDate).
		SetImage(creditCard.Visualisation).
		Send(ctx)
}

func (UtilityCommands) visa13(ctx *dgc.Ctx) {
	creditCard := api.GenerateCreditCard("visa13", ctx.Arguments.Raw())
	embed.NewEmbed().
		AddField("Card Type", "Visa").
		AddField("Card Number", creditCard.Number).
		AddField("CVV", creditCard.CVV).
		AddField("Expiry Date", creditCard.ExpiryDate).
		SetImage(creditCard.Visualisation).
		Send(ctx)
}

func (UtilityCommands) amex(ctx *dgc.Ctx) {
	creditCard := api.GenerateCreditCard("amex", ctx.Arguments.Raw())
	embed.NewEmbed().
		AddField("Card Type", "American Express").
		AddField("Card Number", creditCard.Number).
		AddField("CVV", creditCard.CVV).
		AddField("Expiry Date", creditCard.ExpiryDate).
		SetImage(creditCard.Visualisation).
		Send(ctx)
}

func (UtilityCommands) discover(ctx *dgc.Ctx) {
	creditCard := api.GenerateCreditCard("discover", ctx.Arguments.Raw())
	embed.NewEmbed().
		AddField("Card Type", "Discover").
		AddField("Card Number", creditCard.Number).
		AddField("CVV", creditCard.CVV).
		AddField("Expiry Date", creditCard.ExpiryDate).
		SetImage(creditCard.Visualisation).
		Send(ctx)
}

func (UtilityCommands) diners(ctx *dgc.Ctx) {
	creditCard := api.GenerateCreditCard("diners", ctx.Arguments.Raw())
	embed.NewEmbed().
		AddField("Card Type", "Diners").
		AddField("Card Number", creditCard.Number).
		AddField("CVV", creditCard.CVV).
		AddField("Expiry Date", creditCard.ExpiryDate).
		SetImage(creditCard.Visualisation).
		Send(ctx)
}

func (UtilityCommands) jcb(ctx *dgc.Ctx) {
	creditCard := api.GenerateCreditCard("jcb", ctx.Arguments.Raw())
	embed.NewEmbed().
		AddField("Card Type", "JCB").
		AddField("Card Number", creditCard.Number).
		AddField("CVV", creditCard.CVV).
		AddField("Expiry Date", creditCard.ExpiryDate).
		SetImage(creditCard.Visualisation).
		Send(ctx)
}

func (UtilityCommands) revavatarsearch(ctx *dgc.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		_ = ctx.RespondText("Could not find user")
		return
	}
	embed.NewEmbed().
		AddField("Search URL", fmt.Sprintf("https://www.google.com/searchbyimage?&image_url=%s", url.QueryEscape(user.AvatarURL("")))).
		Send(ctx)
}

func (UtilityCommands) calc(ctx *dgc.Ctx) {
	if runtime.GOOS != "windows" {
		_ = ctx.RespondText("This command is only available on Windows")
	}
	_ = exec.Command("C:\\Windows\\System32\\calc.exe").Start()
}

func (UtilityCommands) notepad(ctx *dgc.Ctx) {
	if runtime.GOOS != "windows" {
		_ = ctx.RespondText("This command is only available on Windows")
	}
	_ = exec.Command("C:\\Windows\\System32\\notepad.exe").Start()
}

func (UtilityCommands) firstmessage(ctx *dgc.Ctx) {
	channel, success := utils.Discord.TryGetChannelFromMention(ctx, ctx.Arguments.Get(0))
	if !success {
		channel, _ = ctx.Session.State.Channel(ctx.Event.ChannelID)
		if channel == nil {
			_ = ctx.RespondText("Could not find channel")
			return
		}
	}
	messages, err := ctx.Session.ChannelMessages(channel.ID, 1, "", "1", "")
	if err != nil {
		_ = ctx.RespondText("Error getting messages")
		return
	}
	if len(messages) != 1 {
		_ = ctx.RespondText("Error getting messages")
		return
	}
	embed.NewEmbed().
		SetDescription("First Message").
		AddField("Author", messages[0].Author.String()).
		AddField("Message", messages[0].Content).
		AddField("Jump URL", utils.Discord.GetMessageLink(messages[0])).
		Send(ctx)
}

func (UtilityCommands) hypesquad(ctx *dgc.Ctx) {
	house := utils.Arguments["<bravery/brilliance/balance>"]
	house = strings.ReplaceAll(house, "bravery", "1")
	house = strings.ReplaceAll(house, "brilliance", "2")
	house = strings.ReplaceAll(house, "balance", "3")
	if house != "1" && house != "2" && house != "3" {
		_ = ctx.RespondText("Missing required argument: <bravery/brilliance/balance>")
		return
	}
	_ = api.SetHypesquad(ctx.Session, house)
}

func (UtilityCommands) leavehypesquad(ctx *dgc.Ctx) {
	_ = api.LeaveHypesquad(ctx.Session)
}

func (UtilityCommands) theme(ctx *dgc.Ctx) {
	theme := strings.ToLower(utils.Arguments["<light/dark>"])
	if theme != "light" && theme != "dark" {
		_ = ctx.RespondText("Error: <light/dark> must be either \"light\" or \"dark\"")
		return
	}
	_ = api.ChangeSettings(ctx.Session, "/settings", map[string]any{"theme": theme})
}

func (UtilityCommands) compact(ctx *dgc.Ctx) {
	toggle := strings.ToLower(utils.Arguments["<on/off>"])
	if toggle != "on" && toggle != "off" {
		_ = ctx.RespondText("Error: <on/off> must be either \"on\" or \"off\"")
		return
	}
	_ = api.ChangeSettings(ctx.Session, "/settings", map[string]any{"message_display_compact": toggle == "on"})
}

func (UtilityCommands) devmode(ctx *dgc.Ctx) {
	toggle := strings.ToLower(utils.Arguments["<on/off>"])
	if toggle != "on" && toggle != "off" {
		_ = ctx.RespondText("Error: <on/off> must be either \"on\" or \"off\"")
		return
	}
	_ = api.ChangeSettings(ctx.Session, "/settings", map[string]any{"developer_mode": toggle == "on"})
}

func (UtilityCommands) sessionid(ctx *dgc.Ctx) {
	embed.NewEmbed().
		AddField("Session ID", ctx.Session.State.SessionID).
		Send(ctx)
}

func (UtilityCommands) aboutme(ctx *dgc.Ctx) {
	_ = api.ChangeSettings(ctx.Session, "", map[string]any{"bio": ctx.Arguments.Raw()})
}

func (UtilityCommands) clearaboutme(ctx *dgc.Ctx) {
	_ = api.ChangeSettings(ctx.Session, "", map[string]any{"bio": nil})
}

func (UtilityCommands) displayname(ctx *dgc.Ctx) {
	_ = api.ChangeSettings(ctx.Session, "", map[string]any{"global_name": ctx.Arguments.Raw()})
}

func (UtilityCommands) cleardisplayname(ctx *dgc.Ctx) {
	_ = api.ChangeSettings(ctx.Session, "", map[string]any{"global_name": nil})
}

func (UtilityCommands) pronouns(ctx *dgc.Ctx) {
	_ = api.ChangeSettings(ctx.Session, "/profile", map[string]any{"pronouns": ctx.Arguments.Raw()})
}

func (UtilityCommands) clearpronouns(ctx *dgc.Ctx) {
	_ = api.ChangeSettings(ctx.Session, "/profile", map[string]any{"pronouns": ""})
}

func (UtilityCommands) timezones(ctx *dgc.Ctx) {
	httpResponse := http.Get("http://worldtimeapi.org/api/timezone/" + url.PathEscape(utils.Arguments["<continent>"]))
	if httpResponse.Error != nil {
		_ = ctx.RespondText("Error getting timezones: " + httpResponse.Error.Error())
		return
	}

	json := httpResponse.Json()

	if err := json.Get("error").MustString(); err != "" {
		_ = ctx.RespondText("Error getting timezones: " + err)
		return
	}

	embed.NewEmbed().
		AddField("Continent", utils.String.Capitalize(utils.Arguments["<continent>"])).
		AddField("Timezones", strings.Join(json.MustStringArray(), ", ")).
		Send(ctx)
}

func (UtilityCommands) timezone(ctx *dgc.Ctx) {
	httpResponse := http.Get("http://worldtimeapi.org/api/timezone/" + utils.Arguments["<timezone>"])
	if httpResponse.Error != nil {
		_ = ctx.RespondText("Error getting timezones: " + httpResponse.Error.Error())
		return
	}

	json := httpResponse.Json()

	if err := json.Get("error").MustString(); err != "" {
		_ = ctx.RespondText("Error getting timezones: " + err)
		return
	}

	embed.NewEmbed().
		AddField("Timezone", json.Get("timezone").MustString()).
		AddField("Abbreviation", json.Get("abbreviation").MustString()).
		AddField("UTC Offset", json.Get("utc_offset").MustString()).
		AddField("Current Time", json.Get("datetime").MustString()).
		Send(ctx)
}

func (UtilityCommands) getmessagejson(ctx *dgc.Ctx) {
	channel, success := utils.Discord.TryGetChannelFromMention(ctx, ctx.Arguments.Get(0))
	if !success {
		_ = ctx.RespondText("Could not find channel")
		return
	}

	messages, err := ctx.Session.ChannelMessages(channel.ID, 1, "", "", ctx.Arguments.Get(1).Raw())
	if err != nil {
		_ = ctx.RespondText("Could not get message: " + err.Error())
		return
	}

	if len(messages) == 0 {
		_ = ctx.RespondText("Could not find message")
		return
	}

	message := messages[0]

	jsonData, err := json.MarshalIndent(message, "", "  ")
	if err != nil {
		_ = ctx.RespondText("Error marshaling JSON")
		return
	}

	embed.NewEmbed().
		SetTitle("Message JSON").
		SetDescription(string(jsonData)).
		Send(ctx)
}

func (UtilityCommands) getembedjson(ctx *dgc.Ctx) {
	channel, success := utils.Discord.TryGetChannelFromMention(ctx, ctx.Arguments.Get(0))
	if !success {
		_ = ctx.RespondText("Could not find channel")
		return
	}

	messages, err := ctx.Session.ChannelMessages(channel.ID, 1, "", "", ctx.Arguments.Get(1).Raw())
	if err != nil {
		_ = ctx.RespondText("Could not get message: " + err.Error())
		return
	}

	if len(messages) == 0 {
		_ = ctx.RespondText("Could not find message")
		return
	}

	message := messages[0]

	jsonData, err := json.MarshalIndent(message.Embeds, "", "  ")
	if err != nil {
		_ = ctx.RespondText("Error marshaling JSON")
		return
	}

	embed.NewEmbed().
		SetTitle("Message JSON").
		SetDescription(string(jsonData)).
		Send(ctx)
}

func (UtilityCommands) calculate(ctx *dgc.Ctx) {
	tv, err := types.Eval(token.NewFileSet(), nil, token.NoPos, utils.Arguments["<expression>"])
	if err != nil {
		_ = ctx.RespondText("Error evaluating expression: " + err.Error())
		return
	}

	if tv.Value == nil {
		_ = ctx.RespondText("Error evaluating expression: nil value")
		return
	}

	embed.NewEmbed().
		AddField("Expression", utils.Arguments["<expression>"]).
		AddField("Solution", tv.Value.String()).
		Send(ctx)
}

func (UtilityCommands) sendinterval(ctx *dgc.Ctx) {
	interval, err := strconv.Atoi(utils.Arguments["<interval>"])
	if err != nil {
		_ = ctx.RespondText("Error parsing interval: " + err.Error())
		return
	}

	if interval < 1 {
		_ = ctx.RespondText("Error parsing interval: interval must be at least 1 second")
		return
	}

	amount, err := strconv.Atoi(utils.Arguments["<amount>"])
	if err != nil {
		_ = ctx.RespondText("Error parsing amount: " + err.Error())
		return
	}

	if amount < 1 {
		_ = ctx.RespondText("Error parsing amount: amount must be at least 1")
		return
	}

	message := utils.Arguments["<\"message\">"]
	flags := utils.Arguments["[flags]"]

	for i := 0; i < amount; i++ {
		msg, err := ctx.Session.ChannelMessageSend(ctx.Event.ChannelID, message)
		if err != nil {
			return
		}

		if strings.Contains(flags, "INSTANT_DEL") {
			_ = ctx.Session.ChannelMessageDelete(ctx.Event.ChannelID, msg.ID)
		} else if strings.Contains(flags, "DELETE") {
			go func() {
				time.Sleep(time.Duration(interval) * time.Second)
				_ = ctx.Session.ChannelMessageDelete(ctx.Event.ChannelID, msg.ID)
			}()
		}

		time.Sleep(time.Duration(interval) * time.Second)
	}
}

func (UtilityCommands) screenshot(ctx *dgc.Ctx) {
	duration, _ := time.ParseDuration(utils.Arguments["[delay]"])
	time.Sleep(duration)

	img, err := screenshot.CaptureScreen()
	if err != nil {
		_ = ctx.RespondText("Error capturing screenshot: " + err.Error())
		return
	}

	buf := new(bytes.Buffer)

	err = png.Encode(buf, img)
	if err != nil {
		_ = ctx.RespondText("Error encoding image: " + err.Error())
		return
	}

	_, _ = ctx.Session.ChannelFileSend(ctx.Event.ChannelID, "screenshot.png", buf)
}

func (UtilityCommands) listdestinations(ctx *dgc.Ctx) {
	destinationsString := []string{}

	for name, dest := range destinations.Destinations {
		desc := fmt.Sprintf("`%s`  - Supports all file types: **%s**", name, utils.Misc.FormatBoolYesNo(dest.SupportsAllFiles()))

		destinationsString = append(destinationsString, desc)
	}

	sort.Strings(destinationsString)

	embed.NewEmbed().
		AddField("Destinations", strings.Join(destinationsString, "\n")).
		Send(ctx)
}

func (UtilityCommands) uploadss(ctx *dgc.Ctx) {
	dest, ok := destinations.Destinations[utils.Arguments["<destination>"]]
	if !ok {
		_ = ctx.RespondText("Error: invalid destination\n\nUse the `listdestinations` command to view all destinations")
		return
	}

	duration, _ := time.ParseDuration(utils.Arguments["[delay]"])
	time.Sleep(duration)

	img, err := screenshot.CaptureScreen()
	if err != nil {
		_ = ctx.RespondText("Error capturing screenshot: " + err.Error())
		return
	}

	buf := new(bytes.Buffer)

	err = png.Encode(buf, img)
	if err != nil {
		_ = ctx.RespondText("Error encoding image: " + err.Error())
		return
	}

	result, err := dest.Upload(buf, "screenshot.png")
	if err != nil {
		_ = ctx.RespondText("Error uploading screenshot: " + err.Error())
		return
	}

	embed.NewEmbed().
		SetDescription("Screenshot").
		AddField("URL", result.URL).
		AddField("Direct URL", result.DirectURL).
		Send(ctx)
}

func (UtilityCommands) uploadfile(ctx *dgc.Ctx) {
	dest, ok := destinations.Destinations[utils.Arguments["<destination>"]]
	if !ok {
		_ = ctx.RespondText("Error: invalid destination\n\nUse the `listdestinations` command to view all destinations")
		return
	}

	attachments := ctx.Event.Message.Attachments

	if len(attachments) == 0 {
		_ = ctx.RespondText("Error: you must provide a file as an attachment.")
		return
	}

	resp, err := nethttp.Get(attachments[0].URL)
	if err != nil {
		_ = ctx.RespondText("Error downloading file: " + err.Error())
		return
	}
	defer resp.Body.Close()

	buf := new(bytes.Buffer)

	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		_ = ctx.RespondText("Error reading file body: " + err.Error())
		return
	}

	result, err := dest.Upload(buf, attachments[0].Filename)
	if err != nil {
		_ = ctx.RespondText("Error uploading file: " + err.Error())
		return
	}

	embed.NewEmbed().
		SetDescription("Uploaded file: "+attachments[0].Filename).
		AddField("URL", result.URL).
		AddField("Direct URL", result.DirectURL).
		Send(ctx)
}

func (UtilityCommands) uploadlink(ctx *dgc.Ctx) {
	dest, ok := destinations.Destinations[utils.Arguments["<destination>"]]
	if !ok {
		_ = ctx.RespondText("Error: invalid destination\n\nUse the `listdestinations` command to view all destinations")
		return
	}

	link := utils.Arguments["<link>"]
	parsedURL, err := url.ParseRequestURI(link)
	if err != nil {
		_ = ctx.RespondText("Error: invalid link: " + err.Error())
		return
	}

	fileName := path.Base(parsedURL.Path)

	resp, err := nethttp.Get(link)
	if err != nil {
		_ = ctx.RespondText("Error downloading file: " + err.Error())
		return
	}
	defer resp.Body.Close()

	buf := new(bytes.Buffer)

	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		_ = ctx.RespondText("Error reading file body: " + err.Error())
		return
	}

	result, err := dest.Upload(buf, fileName)
	if err != nil {
		_ = ctx.RespondText("Error uploading file: " + err.Error())
		return
	}

	embed.NewEmbed().
		SetDescription("Uploaded file: "+fileName).
		AddField("URL", result.URL).
		AddField("Direct URL", result.DirectURL).
		Send(ctx)
}

func (UtilityCommands) obfuscateFunc(lang, ext string) func(ctx *dgc.Ctx) {
	return func(ctx *dgc.Ctx) {
		script := utils.Arguments["[script/attachment]"]

		if script == "" && len(ctx.Event.Message.Attachments) > 0 {
			httpResponse := http.Get(ctx.Event.Message.Attachments[0].URL)
			if httpResponse.Error != nil {
				_ = ctx.RespondText("Error getting script: " + httpResponse.Error.Error())
				return
			}

			script = httpResponse.Body
		}

		if script == "" {
			_ = ctx.RespondText("Error: you must provide a script as an attachment or argument.")
			return
		}

		obf := api.Obfuscate(lang, script)
		if strings.HasPrefix(obf, "An error") {
			_ = ctx.RespondText(obf)
			return
		}

		_, _ = ctx.Session.ChannelMessageSendComplex(ctx.Event.ChannelID, &discordgo.MessageSend{
			Content: "Obfuscated Script",
			Files: []*discordgo.File{
				{
					Name:        "script." + ext,
					ContentType: "text/plain",
					Reader:      strings.NewReader(obf),
				},
			},
		})
	}
}

// ReviewDB

func (UtilityCommands) reviewdbtoken(ctx *dgc.Ctx) {
	token, err := api.GetReviewDBToken(ctx.Session)
	if err != nil {
		embed.NewEmbed().
			SetDescription("Error getting token: " + err.Error()).
			Send(ctx)
		return
	}

	embed.NewEmbed().
		AddField("Token", token).
		Send(ctx)
}

func (UtilityCommands) getreviews(ctx *dgc.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		embed.NewEmbed().
			SetDescription("Could not find user").
			Send(ctx)
		return
	}

	reviews, err := api.GetReviewDBReviews(user.ID)
	if err != nil {
		embed.NewEmbed().
			SetDescription("Error getting reviews: " + err.Error()).
			Send(ctx)
		return
	}

	reviewString := ""

	for _, review := range reviews {
		reviewString += fmt.Sprintf("**[%s] %s:**\n%s\n\n", review.Timestamp.Format(time.RFC1123), review.Username, review.Comment)
	}

	embed.NewEmbed().
		AddField("Reviews", reviewString).
		Send(ctx)
}

func (UtilityCommands) addreview(ctx *dgc.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		embed.NewEmbed().
			SetDescription("Could not find user").
			Send(ctx)
		return
	}

	err := api.AddReviewDBReview(ctx.Session, user.ID, utils.String.RemoveWords(ctx.Arguments.Raw(), 1))
	if err != nil {
		embed.NewEmbed().
			SetDescription("Error adding review: " + err.Error()).
			Send(ctx)
		return
	}

	embed.NewEmbed().
		SetDescription("Successfully added review to " + user.Username + "!").
		Send(ctx)
}
