package commands

import (
	"bytes"
	"fmt"
	"math"
	"net/url"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/LightningDev1/LB-Selfbot/eventhandlers"
	"github.com/LightningDev1/LB-Selfbot/ipc"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"

	"github.com/LightningDev1/LB-Selfbot/api"
	"github.com/LightningDev1/LB-Selfbot/constants"
	"github.com/LightningDev1/LB-Selfbot/embed"
	"github.com/LightningDev1/LB-Selfbot/http"
	"github.com/LightningDev1/LB-Selfbot/utils"
	"github.com/LightningDev1/dgc"
	"github.com/LightningDev1/discordgo"
)

type UtilityCommands struct{}

func (u UtilityCommands) RegisterCommands(router *dgc.Router) {
	router.StartCategory("Utility", "Utility commands")
	router.RegisterCmd(&dgc.Command{
		Name:        "copyserver",
		Description: "Make an exact copy of a server",
		Usage:       "[p]copyserver [server id]",
		Handler:     u.copyserver,
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
		Usage:       "[p]leavegroups",
		Handler:     u.leavegroups,
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
		Name:        "shudownpc",
		Description: "Shut your PC down",
		Usage:       "[p]shudownpc",
		Handler:     u.shudownpc,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "restartpc",
		Description: "Restart your PC",
		Usage:       "[p]restartpc",
		Handler:     u.restartpc,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "snipe",
		Description: "See the last deleted message in this channel",
		Usage:       "[p]snipe",
		Handler:     u.snipe,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "editsnipe",
		Description: "See the last edited message in this channel",
		Usage:       "[p]editsnipe",
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
		Handler:     u.website,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "dashboard",
		Description: "Opens dashboard in browser",
		Usage:       "[p]dashboard",
		Handler:     u.dashboard,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "restartbot",
		Description: "Restart LightningBot",
		Usage:       "[p]restartbot",
		Handler:     u.restartbot,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "stopbot",
		Description: "Stop LightningBot",
		Usage:       "[p]stopbot",
		Handler:     u.stopbot,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "changelog",
		Description: "See the latest changelog",
		Usage:       "[p]changelog",
		Handler:     u.changelog,
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
		Handler:     u.cls,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "mastercard",
		Description: "Generate a mastercard credit card",
		Usage:       "[p]mastercard",
		Handler:     u.mastercard,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "visa16",
		Description: "Generate a visa credit card (16 characters)",
		Usage:       "[p]visa16",
		Handler:     u.visa16,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "visa13",
		Description: "Generate a visa credit card (13 characters)",
		Usage:       "[p]visa13",
		Handler:     u.visa13,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "amex",
		Description: "Generate an amex credit card",
		Usage:       "[p]amex",
		Handler:     u.amex,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "discover",
		Description: "Generate a discover credit card",
		Usage:       "[p]discover",
		Handler:     u.discover,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "diners",
		Description: "Generate a diners credit card",
		Usage:       "[p]diners",
		Handler:     u.diners,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "jcb",
		Description: "Generate a JCB credit card",
		Usage:       "[p]jcb",
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
		Handler:     u.calc,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "notepad",
		Description: "Open notepad",
		Usage:       "[p]notepad",
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
}

func getChannelType(channelType discordgo.ChannelType) discordgo.ChannelType {
	if channelType == discordgo.ChannelTypeGuildVoice || channelType == discordgo.ChannelTypeGuildStageVoice {
		return discordgo.ChannelTypeGuildVoice
	}
	return discordgo.ChannelTypeGuildText
}

func (UtilityCommands) copyserver(ctx *dgc.Ctx) {
	guild, err := utils.Discord.GetGuild(ctx.Session, utils.Misc.Or(utils.Arguments["[server id]"], ctx.Event.GuildID))
	if err != nil {
		_ = ctx.RespondText("Error getting server: " + err.Error())
		return
	}

	icon := http.Get(guild.IconURL()).Body

	newGuild, err := ctx.Session.GuildCreate(fmt.Sprintf("%s - Copy", guild.Name))
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

	for _, role := range guild.Roles {
		if role.ID == guild.ID || role.Managed {
			continue
		}
		
		_, err := ctx.Session.GuildRoleCreate(newGuild.ID, &discordgo.RoleParams{
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
	}

	channels, err := ctx.Session.GuildChannels(guild.ID)
	if err != nil {
		_ = ctx.RespondText("Error copying server: " + err.Error())
		return
	}

	categories, categoryChannels := utils.Discord.MapCategories(channels)

	for _, category := range categories {
		newCategory, err := ctx.Session.GuildChannelCreateComplex(newGuild.ID, discordgo.GuildChannelCreateData{
			Name:                 category.Name,
			Type:                 discordgo.ChannelTypeGuildCategory,
			Position:             category.Position,
			PermissionOverwrites: category.PermissionOverwrites,
		})
		if err != nil {
			_ = ctx.RespondText("Error copying server: " + err.Error())
			return
		}
		for _, channel := range categoryChannels[category.ID] {
			_, err = ctx.Session.GuildChannelCreateComplex(newGuild.ID, discordgo.GuildChannelCreateData{
				Name:                 channel.Name,
				Type:                 getChannelType(channel.Type),
				Topic:                channel.Topic,
				Bitrate:              int(math.Min(float64(channel.Bitrate), 96000)),
				UserLimit:            channel.UserLimit,
				RateLimitPerUser:     channel.RateLimitPerUser,
				Position:             channel.Position,
				PermissionOverwrites: channel.PermissionOverwrites,
				ParentID:             newCategory.ID,
				NSFW:                 channel.NSFW,
			})
			if err != nil {
				_ = ctx.RespondText("Error copying server: " + err.Error())
				return
			}
		}
	}
	embed.NewEmbed().
		SetDescription("Server copied successfully!").
		Send(ctx)
}

func (UtilityCommands) leaveservers(ctx *dgc.Ctx) {
	for _, guild := range ctx.Session.State.Guilds {
		_ = ctx.Session.GuildLeave(guild.ID)
	}
}

func (UtilityCommands) leavegroups(ctx *dgc.Ctx) {
	for _, channel := range ctx.Session.State.PrivateChannels {
		if channel.Type == discordgo.ChannelTypeGroupDM {
			_, _ = ctx.Session.ChannelDelete(channel.ID)
		}
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
	adminServerString := ""
	for _, guild := range ctx.Session.State.Guilds {
		if len(guild.Channels) < 1 {
			continue
		}
		member, _ := ctx.Session.GuildMember(guild.ID, ctx.Session.State.User.ID)
		_ = ctx.Session.State.MemberAdd(member)
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

func (UtilityCommands) optionpoll(ctx *dgc.Ctx) {
	message := embed.NewEmbed().
		SetDescription("Poll: "+utils.Arguments["<\"question\">"]).
		AddField("👍", utils.Arguments["<\"option1\">"]).
		AddField("👎", utils.Arguments["<\"option2\">"]).
		Send(ctx)[0]
	if message != nil {
		_ = ctx.Session.MessageReactionAdd(ctx.Event.ChannelID, message.ID, "👍")
		_ = ctx.Session.MessageReactionAdd(ctx.Event.ChannelID, message.ID, "👎")
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
	httpResponse := http.Get("https://some-random-api.ml/bottoken")
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

func (UtilityCommands) shudownpc(ctx *dgc.Ctx) {
	_ = utils.Process.ShutdownPC()
}

func (UtilityCommands) restartpc(ctx *dgc.Ctx) {
	_ = utils.Process.RestartPC()
}

func (UtilityCommands) snipe(ctx *dgc.Ctx) {
	messages, ok := eventhandlers.EventHandlers.Settings.DeletedMessages[ctx.Event.ChannelID]
	if !ok || len(messages) == 0 {
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
	messages, ok := eventhandlers.EventHandlers.Settings.EditedMessages[ctx.Event.ChannelID]
	if !ok || len(messages) == 0 {
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

func (UtilityCommands) history(ctx *dgc.Ctx) {
	usedCommands := ""
	index := 1
	for i := len(eventhandlers.EventHandlers.Settings.UsedCommands) - 1; i >= 0 && i <= 10; i-- {
		cmd := eventhandlers.EventHandlers.Settings.UsedCommands[i]
		usedCommands += fmt.Sprintf("%d. %s\n", index, cmd.Event.Content)
		index++
	}
	embed.NewEmbed().
		AddField("Used Commands", usedCommands).
		Send(ctx)
}

func (UtilityCommands) website(*dgc.Ctx) {
	_ = utils.Process.OpenURL("https://lightning-bot.com/")
}

func (UtilityCommands) dashboard(*dgc.Ctx) {
	_ = utils.Process.OpenURL("https://lightning-bot.com/dashboard")
}

func (UtilityCommands) restartbot(*dgc.Ctx) {
	ipc.IPC.Restart()
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

func (UtilityCommands) uptime(ctx *dgc.Ctx) {
	embed.NewEmbed().
		AddField("Uptime", utils.Misc.FormatDuration(time.Since(eventhandlers.EventHandlers.StartTime))).
		Send(ctx)
}

func (UtilityCommands) cls(*dgc.Ctx) {
	ipc.IPC.ConsoleClear()
}

func (UtilityCommands) mastercard(ctx *dgc.Ctx) {
	creditCard := api.GenerateCreditCard("mastercard")
	embed.NewEmbed().
		AddField("Card Type", "Mastercard").
		AddField("Card Number", creditCard.Number).
		AddField("CCV", creditCard.CCV).
		AddField("Expiry Date", creditCard.ExpiryDate).
		Send(ctx)
}

func (UtilityCommands) visa16(ctx *dgc.Ctx) {
	creditCard := api.GenerateCreditCard("visa16")
	embed.NewEmbed().
		AddField("Card Type", "Visa").
		AddField("Card Number", creditCard.Number).
		AddField("CCV", creditCard.CCV).
		AddField("Expiry Date", creditCard.ExpiryDate).
		Send(ctx)
}

func (UtilityCommands) visa13(ctx *dgc.Ctx) {
	creditCard := api.GenerateCreditCard("visa13")
	embed.NewEmbed().
		AddField("Card Type", "Visa").
		AddField("Card Number", creditCard.Number).
		AddField("CCV", creditCard.CCV).
		AddField("Expiry Date", creditCard.ExpiryDate).
		Send(ctx)
}

func (UtilityCommands) amex(ctx *dgc.Ctx) {
	creditCard := api.GenerateCreditCard("amex")
	embed.NewEmbed().
		AddField("Card Type", "American Express").
		AddField("Card Number", creditCard.Number).
		AddField("CCV", creditCard.CCV).
		AddField("Expiry Date", creditCard.ExpiryDate).
		Send(ctx)
}

func (UtilityCommands) discover(ctx *dgc.Ctx) {
	creditCard := api.GenerateCreditCard("discover")
	embed.NewEmbed().
		AddField("Card Type", "Discover").
		AddField("Card Number", creditCard.Number).
		AddField("CCV", creditCard.CCV).
		AddField("Expiry Date", creditCard.ExpiryDate).
		Send(ctx)
}

func (UtilityCommands) diners(ctx *dgc.Ctx) {
	creditCard := api.GenerateCreditCard("diners")
	embed.NewEmbed().
		AddField("Card Type", "Diners").
		AddField("Card Number", creditCard.Number).
		AddField("CCV", creditCard.CCV).
		AddField("Expiry Date", creditCard.ExpiryDate).
		Send(ctx)
}

func (UtilityCommands) jcb(ctx *dgc.Ctx) {
	creditCard := api.GenerateCreditCard("jcb")
	embed.NewEmbed().
		AddField("Card Type", "JCB").
		AddField("Card Number", creditCard.Number).
		AddField("CCV", creditCard.CCV).
		AddField("Expiry Date", creditCard.ExpiryDate).
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
		AddField("Jump URL", fmt.Sprintf("https://discord.com/channels/%s/%s/%s", channel.GuildID, channel.ID, messages[0].ID)).
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
	_ = api.ChangeSettings(ctx.Session, map[string]interface{}{"theme": theme}, true)
}

func (UtilityCommands) compact(ctx *dgc.Ctx) {
	toggle := strings.ToLower(utils.Arguments["<on/off>"])
	if toggle != "on" && toggle != "off" {
		_ = ctx.RespondText("Error: <on/off> must be either \"on\" or \"off\"")
		return
	}
	_ = api.ChangeSettings(ctx.Session, map[string]interface{}{"message_display_compact": toggle == "on"}, true)
}

func (UtilityCommands) devmode(ctx *dgc.Ctx) {
	toggle := strings.ToLower(utils.Arguments["<on/off>"])
	if toggle != "on" && toggle != "off" {
		_ = ctx.RespondText("Error: <on/off> must be either \"on\" or \"off\"")
		return
	}
	_ = api.ChangeSettings(ctx.Session, map[string]interface{}{"developer_mode": toggle == "on"}, true)
}

func (UtilityCommands) sessionid(ctx *dgc.Ctx) {
	embed.NewEmbed().
		AddField("Session ID", ctx.Session.State.SessionID).
		Send(ctx)
}
