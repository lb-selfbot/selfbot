package commands

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go/token"
	"go/types"
	"image/png"
	"io"
	"math"
	nethttp "net/http"
	"net/url"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/LightningDev1/LB-Selfbot/selfbot/config"
	"github.com/LightningDev1/LB-Selfbot/selfbot/eventhandlers/eventreactions"
	"github.com/LightningDev1/LB-Selfbot/selfbot/ipc"
	"github.com/LightningDev1/LB-Selfbot/selfbot/logging"
	"github.com/LightningDev1/LB-Selfbot/selfbot/uploaders/destinations"
	"github.com/LightningDev1/screenshot"
	"github.com/bitly/go-simplejson"
	"github.com/kkdai/youtube/v2"
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

	// Discord utilities

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
		Usage:       "[p]aboutme [text]",
		Handler:     u.aboutme,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "displayname",
		Description: "Change your display name",
		Usage:       "[p]displayname <text>",
		Handler:     u.displayname,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "pronouns",
		Description: "Change your pronouns",
		Usage:       "[p]pronouns <text>",
		Handler:     u.pronouns,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "nickscan",
		Description: "Search for servers with a nickname",
		Usage:       "[p]nickscan",
		Handler:     u.nickscan,
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
		Name:        "sendinterval",
		Description: "Send a message every x seconds (flags: DELETE, INSTANT_DEL)",
		Usage:       "[p]sendinterval <interval> <amount> <\"message\"> [flags]",
		Handler:     u.sendinterval,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "optionpoll",
		Description: "Create a poll with options",
		Usage:       "[p]optionpoll <\"question\"> <\"option1\"> <\"option2\">",
		Handler:     u.optionpoll,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "checktoken",
		Description: "Check if a Discord token is valid",
		Usage:       "[p]checktoken <token>",
		Handler:     u.checktoken,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "discordstatus",
		Description: "Check the Discord status",
		Usage:       "[p]discordstatus",
		Handler:     u.discordstatus,
	})

	// LightningBot Utilities

	router.RegisterCmd(&dgc.Command{
		Name:        "website",
		Description: "Open the LightningBot website",
		Usage:       "[p]website",
		Flags:       []string{"HOSTING_DANGEROUS"},
		Handler:     u.website,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "dashboard",
		Description: "Open the LightningBot dashboard",
		Usage:       "[p]dashboard",
		Flags:       []string{"HOSTING_DANGEROUS"},
		Handler:     u.dashboard,
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
		Name:        "history",
		Description: "See used command history",
		Usage:       "[p]history",
		Handler:     u.history,
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

	// System utilities

	router.RegisterCmd(&dgc.Command{
		Name:        "systeminfo",
		Description: "See info about your system",
		Usage:       "[p]systeminfo",
		Handler:     u.systeminfo,
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
		Name:        "open",
		Description: "Open a URL in the browser",
		Usage:       "[p]open <url>",
		Flags:       []string{"DANGEROUS", "HOSTING_DANGEROUS"},
		Handler:     u.open,
	})

	// Uploading utilities

	router.RegisterCmd(&dgc.Command{
		Name:        "screenshot",
		Description: "Take a screenshot",
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

	// Paste utilities

	router.RegisterCmd(&dgc.Command{
		Name:        "paste",
		Description: "Create a paste",
		Usage:       "[p]pastebin <text>",
		Handler:     u.pastebin,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "pastefile",
		Description: "Create a paste from a file",
		Usage:       "[p]pastefile <filepath>",
		Handler:     u.pastefile,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "tinyurl",
		Description: "Shorten a URL with TinyURL",
		Usage:       "[p]tinyurl <url>",
		Handler:     u.tinyurl,
	})

	// Obfuscation utilities

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

	// ReviewDB utilities

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

	// Tag utilities

	router.RegisterCmd(&dgc.Command{
		Name:        "tag",
		Description: "Send a tag's content",
		Usage:       "[p]tag <tag>",
		Handler:     u.tag,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "taglist",
		Description: "List all tags",
		Usage:       "[p]taglist",
		Handler:     u.taglist,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "createtag",
		Description: "Add a tag",
		Usage:       "[p]createtag <tag> <content>",
		Handler:     u.createtag,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "deletetag",
		Description: "Delete a tag",
		Usage:       "[p]deletetag <tag>",
		Handler:     u.deletetag,
	})

	// Creditcard commands

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

	// Network utilities

	router.RegisterCmd(&dgc.Command{
		Name:        "ping",
		Description: "Check the latency of the bot",
		Usage:       "[p]ping",
		Handler:     u.ping,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "pingweb",
		Description: "Check the latency to a website",
		Usage:       "[p]pingweb <host>",
		Handler:     u.pingweb,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "ipinfo",
		Description: "Get info about an IP Address",
		Usage:       "[p]ipinfo <ip>",
		Handler:     u.ipinfo,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "portscan",
		Description: "Scan a host for open ports",
		Usage:       "[p]portscan <host>",
		Handler:     u.portscan,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "resolve",
		Description: "Get a domain's IP address",
		Usage:       "[p]resolve <domain>",
		Handler:     u.resolve,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "reversedns",
		Description: "Find all domains hosted on an IP address",
		Usage:       "[p]reversedns <ip>",
		Handler:     u.reversedns,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "reversemac",
		Description: "Find the vender of a MAC address",
		Usage:       "[p]reversemac <mac>",
		Handler:     u.reversemac,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "headers",
		Description: "Get URL response headers",
		Usage:       "[p]headers <url>",
		Handler:     u.headers,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "interface",
		Description: "See your public/local ip and MAC address",
		Usage:       "[p]interface",
		Flags:       []string{"DANGEROUS", "HOSTING_DANGEROUS"},
		Handler:     u._interface,
	})

	// Misc utilities

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
		Name:        "define",
		Description: "Define a word",
		Usage:       "[p]define <word>",
		Handler:     u.define,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "urbandictionary",
		Description: "Search someting on the Urban Dictionary",
		Usage:       "[p]urbandictionary <word>",
		Handler:     u.urbandictionary,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "wikipedia",
		Description: "Search someting on Wikipedia",
		Usage:       "[p]wikipedia <query>",
		Handler:     u.wikipedia,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "weather",
		Description: "Check the weather in a city",
		Usage:       "[p]weather <city>",
		Handler:     u.weather,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "lvbypass",
		Description: "Bypass a Linkvertise link",
		Usage:       "[p]lvbypass <url>",
		Handler:     u.lvbypass,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "flightinfo",
		Description: "Get info about a flight",
		Usage:       "[p]flightinfo <flight number>",
		Handler:     u.flightinfo,
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
		Name:        "revimage",
		Description: "Reverse image search",
		Usage:       "[p]revimage <image url>",
		Handler:     u.revimage,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "revavatar",
		Description: "Reverse image search a users avatar",
		Usage:       "[p]revavatar <user>",
		Handler:     u.revavatar,
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
		Name:        "calculate",
		Description: "Calculate something",
		Usage:       "[p]calculate <expression>",
		Handler:     u.calculate,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "googlesearch",
		Description: "Search someting on Google",
		Usage:       "[p]googlesearch <query>",
		Handler:     u.googlesearch,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "youtubesearch",
		Description: "Search someting on Youtube",
		Usage:       "[p]youtubesearch <query>",
		Handler:     u.youtubesearch,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "youtubedl",
		Description: "Download a YouTube video to the desktop",
		Usage:       "[p]youtubedl <query>",
		Handler:     u.youtubedl,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "amazon",
		Description: "Search products on Amazon",
		Usage:       "[p]amazon <query>",
		Handler:     u.amazon,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "steam",
		Description: "Search games on Steam",
		Usage:       "[p]steam <query>",
		Handler:     u.steam,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "github",
		Description: "Search repos on GitHub",
		Usage:       "[p]github <query>",
		Handler:     u.github,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "stackoverflow",
		Description: "Search answers on Stack Overflow",
		Usage:       "[p]stackoverflow <query>",
		Handler:     u.stackoverflow,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "searchcountry",
		Description: "Get info about a country",
		Usage:       "[p]searchcountry <country>",
		Handler:     u.searchcountry,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "randomcountry",
		Description: "Get info about a random country",
		Usage:       "[p]randomcountry <country>",
		Handler:     u.randomcountry,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "loremipsum",
		Description: "Generate a random placeholder text",
		Usage:       "[p]loremipsum [paragraphs]",
		Handler:     u.loremipsum,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "baconipsum",
		Description: "Generate a meaty placeholder text",
		Usage:       "[p]baconipsum [paragraphs]",
		Handler:     u.baconipsum,
	})
}

// Discord utilities

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
	var bio any = nil
	if arg := ctx.Arguments.Raw(); arg != "" {
		bio = arg
	}

	_ = api.ChangeSettings(ctx.Session, "", map[string]any{"bio": bio})
}

func (UtilityCommands) displayname(ctx *dgc.Ctx) {
	var globalName any = nil
	if arg := ctx.Arguments.Raw(); arg != "" {
		globalName = arg
	}

	_ = api.ChangeSettings(ctx.Session, "", map[string]any{"global_name": globalName})
}

func (UtilityCommands) pronouns(ctx *dgc.Ctx) {
	_ = api.ChangeSettings(ctx.Session, "/profile", map[string]any{"pronouns": ctx.Arguments.Raw()})
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

func (UtilityCommands) checktoken(ctx *dgc.Ctx) {
	token := utils.Arguments["<token>"]

	response, err := api.GetTokenData(ctx.Session, token)
	if err != nil {
		if !bytes.Contains(response, []byte("401: Unauthorized")) {
			_ = ctx.RespondText("An error has occurred: " + err.Error())
			return
		}

		embed.NewEmbed().
			AddField("Status", "Invalid Token").
			Send(ctx)
		return
	}

	tokenData, err := simplejson.NewJson(response)
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
		return
	}

	embed.NewEmbed().
		AddField("Status", "Valid Token").
		AddField("Username", fmt.Sprintf("%s#%s", tokenData.Get("username").MustString(), tokenData.Get("discriminator").MustString())).
		AddField("ID", tokenData.Get("id").MustString()).
		AddField("Public Flags", fmt.Sprintf("%.f", tokenData.Get("public_flags").MustFloat64())).
		AddField("Flags", fmt.Sprintf("%.f", tokenData.Get("flags").MustFloat64())).
		AddField("Purchased Flags", fmt.Sprintf("%.f", tokenData.Get("purchased_flags").MustFloat64())).
		AddField("Banner", tokenData.Get("banner").MustString()).
		AddField("Banner Color", tokenData.Get("banner_color").MustString()).
		AddField("Bio", tokenData.Get("bio").MustString()).
		AddField("Locale", tokenData.Get("locale").MustString()).
		AddField("NSFW Allowed", utils.Misc.FormatBoolYesNo(tokenData.Get("nsfw_allowed").MustBool())).
		AddField("2FA Enabled", utils.Misc.FormatBoolYesNo(tokenData.Get("mfa_enabled").MustBool())).
		AddField("Email", tokenData.Get("email").MustString()).
		AddField("Verified", utils.Misc.FormatBoolYesNo(tokenData.Get("verified").MustBool())).
		AddField("Phone", tokenData.Get("phone").MustString()).
		Send(ctx)
}

func (UtilityCommands) discordstatus(ctx *dgc.Ctx) {
	httpResponse := http.Get("https://srhpyqt94yxb.statuspage.io/api/v2/summary.json")
	if httpResponse.Error != nil {
		_ = ctx.RespondText("An error has occurred.")
		return
	}

	components := httpResponse.Json().Get("components")

	emb := embed.NewEmbed()

	for i := 0; i < len(components.MustArray()); i++ {
		component := components.GetIndex(i)
		emb.AddField(
			component.Get("name").MustString(),
			component.Get("status").MustString(),
		)
	}

	emb.Send(ctx)
}

// LightningBot Utilities

func (UtilityCommands) website(*dgc.Ctx) {
	_ = utils.System.OpenURL("https://lightning-bot.com/")
}

func (UtilityCommands) dashboard(*dgc.Ctx) {
	_ = utils.System.OpenURL("https://lightning-bot.com/dashboard")
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

func (UtilityCommands) cls(*dgc.Ctx) {
	logging.UI.Clear()
}

func (UtilityCommands) openfolder(*dgc.Ctx) {
	_ = utils.System.OpenURL(utils.IO.GetMainFolder())
}

// System Utilities

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

func (UtilityCommands) open(ctx *dgc.Ctx) {
	_ = utils.System.OpenURL(ctx.Arguments.Raw())
}

// Uploading utilities

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

// Paste utilities

func (UtilityCommands) pastebin(ctx *dgc.Ctx) {
	paste, err := api.CreatePaste(ctx.Arguments.Raw())
	if err != nil {
		_ = ctx.RespondText("An error has occurred.")
		return
	}
	embed.NewEmbed().
		AddField("Paste URL", paste).
		Send(ctx)
}

func (UtilityCommands) pastefile(ctx *dgc.Ctx) {
	filePath := utils.Arguments["<filepath>"]
	if content, success := utils.IO.Read(filePath); success {
		paste, err := api.CreatePaste(content)
		if err != nil {
			_ = ctx.RespondText("An error has occurred.")
			return
		}
		embed.NewEmbed().
			AddField("Paste URL", paste).
			Send(ctx)
	} else {
		_ = ctx.RespondText("That file does not exist.")
	}
}

func (UtilityCommands) tinyurl(ctx *dgc.Ctx) {
	sourceUrl := ctx.Arguments.Raw()
	httpResponse := http.Get(fmt.Sprintf("https://tinyurl.com/api-create.php?url=%s", sourceUrl))
	if httpResponse.Error != nil {
		_ = ctx.RespondText("An error has occurred.")
		return
	}
	embed.NewEmbed().
		AddField("URL", sourceUrl).
		AddField("TinyURL", httpResponse.Body).
		Send(ctx)
}

// Obfuscation utilities

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

// ReviewDB utilities

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

// Tag utilities

func (UtilityCommands) tag(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config: " + err.Error())
		return
	}

	content, ok := cfg.Tags[ctx.Arguments.Get(0).Raw()]
	if !ok {
		_ = ctx.RespondText("Error: tag does not exist")
		return
	}

	_ = ctx.RespondText(content)
}

func (UtilityCommands) taglist(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config: " + err.Error())
		return
	}

	var tagList []string

	for tag := range cfg.Tags {
		tagList = append(tagList, fmt.Sprintf("%s:\n%s\n", tag, cfg.Tags[tag]))
	}

	embed.NewEmbed().
		AddField("Tags", strings.Join(tagList, "\n")).
		Send(ctx)
}

func (UtilityCommands) createtag(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config: " + err.Error())
		return
	}

	if cfg.Tags == nil {
		cfg.Tags = map[string]string{}
	}

	cfg.Tags[ctx.Arguments.Get(0).Raw()] = strings.Join(strings.Split(ctx.Event.Content, " ")[2:], " ")

	err = cfg.Save()
	if err != nil {
		_ = ctx.RespondText("Error saving config: " + err.Error())
		return
	}

	embed.NewEmbed().
		SetDescription("Successfully created tag").
		Send(ctx)
}

func (UtilityCommands) deletetag(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config: " + err.Error())
		return
	}

	delete(cfg.Tags, ctx.Arguments.Get(0).Raw())

	err = cfg.Save()
	if err != nil {
		_ = ctx.RespondText("Error saving config: " + err.Error())
		return
	}

	embed.NewEmbed().
		SetDescription("Successfully deleted tag").
		Send(ctx)
}

// Creditcard commands

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

// Network utilities

func (UtilityCommands) ping(ctx *dgc.Ctx) {
	emb := embed.NewEmbed()
	time, err := api.Ping("https://lightning-bot.com")
	if err != nil {
		emb.AddField("Ping to LightningBot", "Error")
	} else {
		emb.AddField("Ping to LightningBot", fmt.Sprintf("%dms", time.Milliseconds()))
	}
	time, err = api.Ping("https://google.com")
	if err != nil {
		emb.AddField("Ping to Google", "Error")
	} else {
		emb.AddField("Ping to Google", fmt.Sprintf("%dms", time.Milliseconds()))
	}
	emb.
		AddField("Ping to Discord", fmt.Sprintf("%dms", ctx.Session.HeartbeatLatency().Milliseconds())).
		Send(ctx)
}

func (UtilityCommands) pingweb(ctx *dgc.Ctx) {
	host := utils.Arguments["<host>"]
	time, err := api.Ping(host)
	emb := embed.NewEmbed()
	if err != nil {
		emb.AddField("Ping to "+host, "Error")
	} else {
		emb.AddField("Ping to "+host, fmt.Sprintf("%dms", time.Milliseconds()))
	}
	emb.Send(ctx)
}

func (UtilityCommands) ipinfo(ctx *dgc.Ctx) {
	url := fmt.Sprintf(
		"https://ipinfo.io/%s/json",
		url.QueryEscape(utils.Arguments["<ip>"]),
	)

	httpResponse := http.Get(url)
	if httpResponse.Error != nil {
		_ = ctx.RespondText("An error has occurred.")
		return
	}

	ipinfo := httpResponse.Json()

	if ipinfo.Get("status").MustInt(0) == 404 {
		_ = ctx.RespondText("No results found.")
		return
	}

	location := fmt.Sprintf(
		"%s, %s, %s (postal code %s) (coordinates %s)",
		ipinfo.Get("city").MustString(), ipinfo.Get("region").MustString(),
		ipinfo.Get("country").MustString(), ipinfo.Get("postal").MustString(),
		ipinfo.Get("loc").MustString(),
	)

	embed.NewEmbed().
		AddField("IP", utils.Arguments["<ip>"]).
		AddField("Hostname", ipinfo.Get("hostname").MustString()).
		AddField("Timezone", ipinfo.Get("timezone").MustString()).
		AddField("Location", location).
		Send(ctx)
}

func (UtilityCommands) portscan(ctx *dgc.Ctx) {
	host := utils.Arguments["<host>"]
	ports := map[string]int{
		"FTP":            21,
		"SSH":            22,
		"Telnet":         23,
		"SMTP":           25,
		"DNS":            53,
		"HTTP":           80,
		"POP3":           110,
		"RPC":            135,
		"IMAP":           143,
		"HTTPS":          443,
		"SMB":            445,
		"MS-SQL":         1433,
		"MySQL":          3306,
		"Remote Desktop": 3389,
		"VNC":            5900,
		"HTTP Alt":       8080,
	}
	emb := embed.NewEmbed()
	for name, port := range ports {
		if utils.Misc.IsPortOpen(host, port) {
			emb.AddField(fmt.Sprintf("%s (%d)", name, port), "Open")
		}
	}
	emb.Send(ctx)
}

func (UtilityCommands) resolve(ctx *dgc.Ctx) {
	domain := utils.Arguments["<domain>"]

	httpResponse := http.Get(fmt.Sprintf("https://json.geoiplookup.io/%s", domain))
	if httpResponse.Error != nil {
		_ = ctx.RespondText("An error has occurred.")
		return
	}

	response := httpResponse.Json()

	if response.Get("ip").MustString() == "" {
		embed.NewEmbed().SetDescription("Invalid domain.").Send(ctx)
		return
	}

	embed.NewEmbed().
		AddField("Domain", domain).
		AddField("IP", response.Get("ip").MustString()).
		Send(ctx)
}

func (UtilityCommands) reversedns(ctx *dgc.Ctx) {
	url := fmt.Sprintf(
		"https://api.domaintools.com/v1/%s/host-domains/",
		url.QueryEscape(utils.Arguments["<ip>"]),
	)

	httpResponse := http.Get(url)
	if httpResponse.Error != nil {
		_ = ctx.RespondText("An error has occurred.")
		return
	}

	response := httpResponse.Json()

	if response.Get("error").MustArray() != nil {
		embed.NewEmbed().SetDescription("Invalid IP Address.").Send(ctx)
		return
	}

	domains := response.GetPath("response", "ip_addresses", "domain_names").MustStringArray()

	embed.NewEmbed().
		AddField("IP Address", utils.Arguments["<ip>"]).
		AddField("Domains", strings.Join(domains, ", ")).
		Send(ctx)
}

func (UtilityCommands) reversemac(ctx *dgc.Ctx) {
	mac := utils.Arguments["<mac>"]
	httpResponse := http.Get(fmt.Sprintf("https://api.macvendors.com/%s", mac))
	if httpResponse.Error != nil {
		_ = ctx.RespondText("An error has occurred.")
		return
	}
	if strings.Contains(httpResponse.Body, "Not Found") {
		embed.NewEmbed().SetDescription("Invalid MAC Address.").Send(ctx)
		return
	}
	embed.NewEmbed().
		AddField("MAC Address", mac).
		AddField("Vendor", httpResponse.Body).
		Send(ctx)
}

func (UtilityCommands) headers(ctx *dgc.Ctx) {
	httpResponse := http.Get(utils.Arguments["<url>"])
	if httpResponse.Error != nil {
		_ = ctx.RespondText("An error has occurred.")
		return
	}

	var headers string
	for key, value := range httpResponse.Headers {
		headers += fmt.Sprintf("%s: %s\n", key, value)
	}

	embed.NewEmbed().
		AddField("URL", utils.Arguments["<url>"]).
		AddField("Headers", headers).
		Send(ctx)
}

func (UtilityCommands) _interface(ctx *dgc.Ctx) {
	embed.NewEmbed().
		AddField("Public IP", api.GetPublicIP()).
		AddField("Local IP", api.GetLocalIP()).
		AddField("MAC Address", api.GetMacAddress()).
		Send(ctx)
}

// Misc utilities

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

func (UtilityCommands) define(ctx *dgc.Ctx) {
	url := fmt.Sprintf(
		"https://api.dictionaryapi.dev/api/v2/entries/en/%s",
		url.QueryEscape(ctx.Arguments.Raw()),
	)

	httpResponse := http.Get(url)
	if httpResponse.Error != nil {
		_ = ctx.RespondText("An error has occurred.")
		return
	}

	list := httpResponse.Json()

	if len(list.MustArray()) == 0 {
		_ = ctx.RespondText("No results found.")
		return
	}

	result := list.GetIndex(0)

	meanings := utils.Map(result.Get("meanings").MustArray(), func(obj any) string {
		meaning := obj.(map[string]any)
		definitions := meaning["definitions"].([]any)

		if len(definitions) == 0 {
			return ""
		}

		definition := definitions[0].(map[string]any)["definition"].(string)

		return fmt.Sprintf("**%s**: %s", meaning["partOfSpeech"].(string), definition)
	})

	meanings = utils.Filter(meanings, func(obj string) bool {
		return obj != ""
	})

	embed.NewEmbed().
		AddField("Word", ctx.Arguments.Raw()).
		AddField("Meanings", strings.Join(meanings, "\n")).
		Send(ctx)
}

func (UtilityCommands) urbandictionary(ctx *dgc.Ctx) {
	url := fmt.Sprintf(
		"https://api.urbandictionary.com/v0/define?term=%s",
		url.QueryEscape(ctx.Arguments.Raw()),
	)

	httpResponse := http.Get(url)
	if httpResponse.Error != nil {
		_ = ctx.RespondText("An error has occurred.")
		return
	}

	list := httpResponse.Json().Get("list")

	if len(list.MustArray()) == 0 {
		_ = ctx.RespondText("No results found.")
		return
	}

	result := list.GetIndex(0)

	embed.NewEmbed().
		AddField("Word", ctx.Arguments.Raw()).
		AddField("Definition", result.Get("definition").MustString()).
		AddField("Permalink", result.Get("permalink").MustString()).
		AddField("Thumbs Up", fmt.Sprint(result.Get("thumbs_up").MustFloat64())).
		AddField("Thumbs Down", fmt.Sprint(result.Get("thumbs_down").MustFloat64())).
		AddField("Author", result.Get("author").MustString()).
		Send(ctx)
}

func (UtilityCommands) wikipedia(ctx *dgc.Ctx) {
	query := ctx.Arguments.Raw()
	article, err := api.WikipediaSearch(query)
	if err != nil {
		_ = ctx.RespondText("An error has occurred.")
		return
	}
	embed.NewEmbed().
		AddField("Query", query).
		AddField("Wikipedia Summary", article[:1500]).
		Send(ctx)
}

func (UtilityCommands) weather(ctx *dgc.Ctx) {
	url := fmt.Sprintf(
		"https://api.openweathermap.org/data/2.5/weather?appid=533d00f29739edba54c5fb2b1a8262b3&q=%s",
		url.QueryEscape(ctx.Arguments.Raw()),
	)

	httpResponse := http.Get(url)
	if httpResponse.Error != nil {
		_ = ctx.RespondText("An error has occurred: " + httpResponse.Error.Error())
		return
	}

	weather := httpResponse.Json()

	if fmt.Sprint(weather.Get("cod").MustInt()) == "404" {
		_ = ctx.RespondText("No results found.")
		return
	}

	embed.NewEmbed().
		SetDescription(fmt.Sprintf("Weather in %s, %s", weather.Get("name").MustString(), weather.GetPath("sys", "country").MustString())).
		AddField("Description", utils.String.Capitalize(weather.Get("weather").GetIndex(0).Get("description").MustString())).
		AddField("Temperature", fmt.Sprintf("%.2f °C", math.Round(weather.GetPath("main", "temp").MustFloat64()-273.15))).
		AddField("Feels Like", fmt.Sprintf("%.2f °C", math.Round(weather.GetPath("main", "feels_like").MustFloat64()-273.15))).
		AddField("Min. Temperature", fmt.Sprintf("%.2f °C", math.Round(weather.GetPath("main", "temp_min").MustFloat64()-273.15))).
		AddField("Max. Temperature", fmt.Sprintf("%.2f °C", math.Round(weather.GetPath("main", "temp_max").MustFloat64()-273.15))).
		AddField("Pressure", fmt.Sprintf("%.2f hPa", weather.GetPath("main", "pressure").MustFloat64())).
		AddField("Humidity", fmt.Sprintf("%.2f %%", weather.GetPath("main", "humidity").MustFloat64())).
		AddField("Visibility", fmt.Sprintf("%.2f m", weather.Get("visibility").MustFloat64())).
		AddField("Wind Speed", fmt.Sprintf("%.2f Km/h", math.Round(weather.GetPath("wind", "speed").MustFloat64()*3.6))).
		AddField("Wind Gust", fmt.Sprintf("%.2f Km/h", math.Round(weather.GetPath("wind", "gust").MustFloat64()*3.6))).
		AddField("Wind Direction", fmt.Sprintf("%.2f°", weather.GetPath("wind", "deg").MustFloat64())).
		Send(ctx)
}

func (UtilityCommands) lvbypass(ctx *dgc.Ctx) {
	resp, err := http.GetJson("https://bypass.pm/bypass2?url=" + ctx.Arguments.Raw())
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
		return
	}

	destination := resp.Get("destination").MustString()

	embed.NewEmbed().
		AddField("URL", ctx.Arguments.Raw()).
		AddField("Bypassed URL", utils.Misc.Or(destination, "An error has occurred")).
		Send(ctx)
}

func (UtilityCommands) flightinfo(ctx *dgc.Ctx) {
	uri := fmt.Sprintf(
		"https://api.flightradar24.com/common/v1/flight/list.json?query=%s&fetchBy=flight&page=1",
		url.QueryEscape(utils.Arguments["<flight number>"]),
	)

	httpResponse := http.Get(uri)
	if httpResponse.Error != nil {
		_ = ctx.RespondText("An error has occurred.")
		return
	}

	json := httpResponse.Json()

	if _, exists := json.CheckGet("errors"); exists {
		_ = ctx.RespondText("That flight does not exist or an error has occurred.")
		return
	}

	data := json.GetPath("result", "response", "data").GetIndex(0)

	departureTime := data.GetPath("time", "scheduled", "departure").MustInt64()
	arrivalTime := data.GetPath("time", "scheduled", "arrival").MustInt64()

	embed.NewEmbed().
		AddField("Flight Number", data.GetPath("identification", "number", "default").MustString()).
		AddField("Airline", data.GetPath("airline", "name").MustString()).
		AddField("Departure Airport", data.GetPath("airport", "origin", "name").MustString()).
		AddField("Departure Time", time.Unix(departureTime, 0).Format("2006-01-02 15:04:05")).
		AddField("Arrival Airport", data.GetPath("airport", "destination", "name").MustString()).
		AddField("Arrival Time", time.Unix(arrivalTime, 0).Format("2006-01-02 15:04:05")).
		AddField("Status", data.GetPath("status", "text").MustString()).
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

func (UtilityCommands) revimage(ctx *dgc.Ctx) {
	imageUrl := utils.Arguments["<image url>"]
	embed.NewEmbed().
		AddField("Reverse Image Search URL", fmt.Sprintf("https://www.google.com/searchbyimage?image_url=%s", imageUrl)).
		Send(ctx)
}

func (UtilityCommands) revavatar(ctx *dgc.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		_ = ctx.RespondText("Could not find user")
		return
	}
	embed.NewEmbed().
		AddField("Search URL", fmt.Sprintf("https://www.google.com/searchbyimage?&image_url=%s", url.QueryEscape(user.AvatarURL("")))).
		Send(ctx)
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

func (UtilityCommands) googlesearch(ctx *dgc.Ctx) {
	results, err := api.GoogleSearch(ctx.Arguments.Raw())
	if err != nil {
		_ = ctx.RespondText("An error has occurred.")
		return
	}
	emb := embed.NewEmbed()
	for i := 0; i < len(results.MustArray()); i++ {
		result := results.GetIndex(i)
		emb.AddField(result.Get("title").MustString(), result.Get("url").MustString())
	}
	emb.Send(ctx)
}

func (UtilityCommands) youtubesearch(ctx *dgc.Ctx) {
	results, err := api.GoogleSearch(ctx.Arguments.Raw() + " site:youtube.com")
	if err != nil {
		_ = ctx.RespondText("An error has occurred.")
		return
	}
	emb := embed.NewEmbed()
	for i := 0; i < len(results.MustArray()); i++ {
		result := results.GetIndex(i)
		resultUrl := result.Get("url").MustString()
		if strings.Contains(resultUrl, "playlist") {
			continue
		}
		emb.AddField(result.Get("title").MustString(), result.Get("url").MustString())
	}
	emb.Send(ctx)
}

func (UtilityCommands) youtubedl(ctx *dgc.Ctx) {
	results, err := api.GoogleSearch(ctx.Arguments.Raw() + " site:youtube.com")
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
		return
	}
	for i := 0; i < len(results.MustArray()); i++ {
		result := results.GetIndex(i)
		resultUrl := result.Get("url").MustString()

		if !strings.Contains(resultUrl, "/watch?v=") {
			continue
		}

		client := youtube.Client{}

		video, err := client.GetVideo(resultUrl)
		if err != nil {
			_ = ctx.RespondText("An error has occurred: " + err.Error())
			return
		}

		_ = ctx.RespondText(fmt.Sprintf("Downloading %s...", video.Title))

		formats := video.Formats.WithAudioChannels()

		stream, _, err := client.GetStream(video, &formats[0])
		if err != nil {
			_ = ctx.RespondText("An error has occurred: " + err.Error())
			return
		}

		downloadPath := filepath.Join(
			utils.IO.GetMainFolder(),
			"Downloads",
			utils.IO.SanitizeFilename(fmt.Sprintf("%s.mp4", video.Title)),
		)

		file, err := os.Create(downloadPath)
		if err != nil {
			_ = ctx.RespondText("An error has occurred: " + err.Error())
			return
		}
		defer file.Close()

		_, err = io.Copy(file, stream)
		if err != nil {
			_ = ctx.RespondText("An error has occurred: " + err.Error())
			_ = file.Close()
			return
		}

		_ = utils.System.OpenURL(downloadPath)
		_ = ctx.RespondText(fmt.Sprintf("Successfully downloaded %s", video.Title))

		break
	}
}

func (UtilityCommands) amazon(ctx *dgc.Ctx) {
	query := ctx.Arguments.Raw()
	embed.NewEmbed().
		AddField("Amazon Search", fmt.Sprintf("https://www.amazon.com/s?k=%s", url.QueryEscape(query))).
		Send(ctx)
}

func (UtilityCommands) steam(ctx *dgc.Ctx) {
	query := ctx.Arguments.Raw()
	embed.NewEmbed().
		AddField("Steam Search", fmt.Sprintf("https://store.steampowered.com/search/?term=%s", url.QueryEscape(query))).
		Send(ctx)
}

func (UtilityCommands) github(ctx *dgc.Ctx) {
	query := ctx.Arguments.Raw()
	embed.NewEmbed().
		AddField("GitHub Search", fmt.Sprintf("https://github.com/search?q=%s", url.QueryEscape(query))).
		Send(ctx)
}

func (UtilityCommands) stackoverflow(ctx *dgc.Ctx) {
	query := ctx.Arguments.Raw()
	embed.NewEmbed().
		AddField("Stack Overflow Search", fmt.Sprintf("https://stackoverflow.com/search?q=%s", url.QueryEscape(query))).
		Send(ctx)
}

func (UtilityCommands) searchcountry(ctx *dgc.Ctx) {
	httpResponse := http.Get(fmt.Sprintf("https://restcountries.com/v3.1/name/%s?fullText=true", ctx.Arguments.Raw()))
	if httpResponse.Error != nil {
		_ = ctx.RespondText("An error has occurred.")
		return
	}

	countryList := httpResponse.Json()

	country := countryList.GetIndex(0)

	rawCurrencies := country.Get("currencies").Interface().(map[string]any)
	var currencies []string
	for shortName, ifa := range rawCurrencies {
		currencies = append(currencies, fmt.Sprintf("%s (%s)", shortName, ifa.(map[string]any)["name"].(string)))
	}

	rawLanguages := country.Get("languages").Interface().(map[string]any)
	var languages []string
	for _, language := range rawLanguages {
		languages = append(languages, language.(string))
	}

	embed.NewEmbed().
		SetAuthor(fmt.Sprintf("%s - %s", country.GetPath("name", "common").MustString(), country.Get("cca2").MustString())).
		AddField("Name", country.GetPath("name", "common").MustString()).
		AddField("Top-Level Domain", strings.Join(country.Get("tld").MustStringArray(), ", ")).
		AddField("Country Code", country.Get("cca2").MustString()).
		AddField("Capital", strings.Join(country.Get("capital").MustStringArray(), ", ")).
		AddField("Region", country.Get("region").MustString()).
		AddField("Subregion", country.Get("subregion").MustString()).
		AddField("Population", fmt.Sprint(country.Get("population").MustInt())).
		AddField("Latitude/Longitude", fmt.Sprintf("%f, %f", country.Get("latlng").GetIndex(0).MustFloat64(), country.Get("latlng").GetIndex(1).MustFloat64())).
		AddField("Timezone", strings.Join(country.Get("timezones").MustStringArray(), ", ")).
		AddField("Currency", strings.Join(currencies, ", ")).
		AddField("Language", strings.Join(languages, ", ")).
		Send(ctx)
}

func (UtilityCommands) randomcountry(ctx *dgc.Ctx) {
	httpResponse := http.Get("https://restcountries.com/v3.1/all")
	if httpResponse.Error != nil {
		_ = ctx.RespondText("An error has occurred.")
		return
	}

	countryList := httpResponse.Json()

	country := countryList.GetIndex(utils.Random.Number(0, len(countryList.MustArray())))

	rawCurrencies := country.Get("currencies").Interface().(map[string]any)
	var currencies []string
	for shortName, ifa := range rawCurrencies {
		currencies = append(currencies, fmt.Sprintf("%s (%s)", shortName, ifa.(map[string]any)["name"].(string)))
	}

	rawLanguages := country.Get("languages").Interface().(map[string]any)
	var languages []string
	for _, language := range rawLanguages {
		languages = append(languages, language.(string))
	}

	embed.NewEmbed().
		SetAuthor(fmt.Sprintf("%s - %s", country.GetPath("name", "common").MustString(), country.Get("cca2").MustString())).
		AddField("Name", country.GetPath("name", "common").MustString()).
		AddField("Top-Level Domain", strings.Join(country.Get("tld").MustStringArray(), ", ")).
		AddField("Country Code", country.Get("cca2").MustString()).
		AddField("Capital", strings.Join(country.Get("capital").MustStringArray(), ", ")).
		AddField("Region", country.Get("region").MustString()).
		AddField("Subregion", country.Get("subregion").MustString()).
		AddField("Population", fmt.Sprint(country.Get("population").MustInt())).
		AddField("Latitude/Longitude", fmt.Sprintf("%f, %f", country.Get("latlng").GetIndex(0).MustFloat64(), country.Get("latlng").GetIndex(1).MustFloat64())).
		AddField("Timezone", strings.Join(country.Get("timezones").MustStringArray(), ", ")).
		AddField("Currency", strings.Join(currencies, ", ")).
		AddField("Language", strings.Join(languages, ", ")).
		Send(ctx)
}

func (UtilityCommands) loremipsum(ctx *dgc.Ctx) {
	paragrahps, err := ctx.Arguments.Get(0).AsInt()
	if err != nil {
		paragrahps = 1
	}

	httpResponse := http.Get(fmt.Sprintf("https://loripsum.net/api/%d/short/plaintext", paragrahps))
	if httpResponse.Error != nil {
		_ = ctx.RespondText("An error has occurred.")
		return
	}

	embed.NewEmbed().
		SetDescription(httpResponse.Body).
		Send(ctx)
}

func (UtilityCommands) baconipsum(ctx *dgc.Ctx) {
	paragrahps, err := ctx.Arguments.Get(0).AsInt()
	if err != nil {
		paragrahps = 1
	}

	httpResponse := http.Get(fmt.Sprintf("https://baconipsum.com/api/?type=all-meat&paras=%d&format=text", paragrahps))
	if httpResponse.Error != nil {
		_ = ctx.RespondText("An error has occurred.")
		return
	}

	embed.NewEmbed().
		SetDescription(httpResponse.Body).
		Send(ctx)
}
