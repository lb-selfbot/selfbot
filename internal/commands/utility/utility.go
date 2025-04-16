package utility

import (
	"bytes"
	"errors"
	"fmt"
	"go/token"
	"go/types"
	"image/png"
	"math"
	"net"
	nethttp "net/http"
	"net/url"
	"os"
	"os/exec"
	"path"
	"regexp"
	"runtime"
	"slices"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/goccy/go-json"

	"github.com/lb-selfbot/selfbot/internal/conversion"
	"github.com/lb-selfbot/selfbot/internal/logging"
	"github.com/lb-selfbot/selfbot/internal/network"
	"github.com/lb-selfbot/selfbot/internal/random"
	"github.com/lb-selfbot/selfbot/internal/services/config"
	"github.com/lb-selfbot/selfbot/internal/services/eventbroker"
	"github.com/lb-selfbot/selfbot/internal/services/router"

	"github.com/lb-selfbot/selfbot/internal/fn"
	"github.com/lb-selfbot/selfbot/internal/services"

	"github.com/LightningDev1/go-foreground"
	"github.com/kbinani/screenshot"
	"github.com/lb-selfbot/discordgo/protos"
	"github.com/lb-selfbot/selfbot/internal/uploaders/destinations"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/host"
	"github.com/shirou/gopsutil/v4/mem"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/lb-selfbot/discordgo"
	"github.com/lb-selfbot/selfbot/internal/api"
	"github.com/lb-selfbot/selfbot/internal/constants"
	"github.com/lb-selfbot/selfbot/internal/http"
	"github.com/lb-selfbot/selfbot/internal/system"
	"github.com/lb-selfbot/selfbot/internal/utils"
)

type UtilityCommands struct{ services *services.Services }

func NewUtilityCommands(s *services.Services) *UtilityCommands {
	return &UtilityCommands{services: s}
}

func (c *UtilityCommands) Register() {
	c.services.Router.RegisterCategory(&router.Category{
		Name:        "Utility",
		Description: "Utility commands",
	}, []*router.Command{
		// Discord utilities
		{
			Name:        "copyserver",
			Description: "Make an exact copy of a server",
			Arguments:   "[server id] [delay]",
			Handler:     c.copyserver,
		},
		{
			Name:        "copychannel",
			Description: "Make an exact copy of a channel",
			Arguments:   "<source channel> <dest channel> [limit]",
			Handler:     c.copychannel,
		},
		{
			Name:        "leavegroups",
			Description: "Leave all of your groups",
			Arguments:   "[silent: on/off]",
			Handler:     c.leavegroups,
		},
		{
			Name:        "adminservers",
			Description: "See all servers you have administator in",
			Handler:     c.adminservers,
		},
		{
			Name:        "admingroups",
			Description: "See all groups you are the owner of",
			Handler:     c.admingroups,
		},
		{
			Name:        "joinvc",
			Description: "Join a voice channel",
			Arguments:   "<channel id>",
			Handler:     c.joinvc,
		},
		{
			Name:        "snipe",
			Description: "See the last deleted message in a channel",
			Arguments:   "[channel]",
			Handler:     c.snipe,
		},
		{
			Name:        "editsnipe",
			Description: "See the last edited message in a channel",
			Arguments:   "[channel]",
			Handler:     c.editsnipe,
		},
		{
			Name:        "firstmessage",
			Description: "See the first message in the channel",
			Arguments:   "[channel]",
			Handler:     c.firstmessage,
		},
		{
			Name:        "hypesquad",
			Description: "Change your Hypesquad House",
			Arguments:   "<bravery/brilliance/balance>",
			Handler:     c.hypesquad,
		},
		{
			Name:        "leavehypesquad",
			Description: "Leave Hypesquad",
			Handler:     c.leavehypesquad,
		},
		{
			Name:        "theme",
			Description: "Change your Discord theme",
			Arguments:   "<dark/light/darker/midnight>",
			Handler:     c.theme,
		},
		{
			Name:        "compact",
			Description: "Toggle Discord compact mode",
			Arguments:   "<on/off>",
			Handler:     c.compact,
		},
		{
			Name:        "devmode",
			Description: "Toggle Discord developer mode",
			Arguments:   "<on/off>",
			Handler:     c.devmode,
		},
		{
			Name:        "sessionid",
			Description: "See LightningBot's Discord session ID",
			Handler:     c.sessionid,
		},
		{
			Name:        "aboutme",
			Description: "Change your About Me text",
			Arguments:   "[text]",
			Handler:     c.aboutme,
		},
		{
			Name:        "displayname",
			Description: "Change your display name",
			Arguments:   "<text>",
			Handler:     c.displayname,
		},
		{
			Name:        "pronouns",
			Description: "Change your pronouns",
			Arguments:   "<text>",
			Handler:     c.pronouns,
		},
		{
			Name:        "nickscan",
			Description: "Search for servers with a nickname",
			Handler:     c.nickscan,
		},
		{
			Name:        "getmessagejson",
			Description: "Get a message's JSON data",
			Arguments:   "<#channel> <message id>",
			Handler:     c.getmessagejson,
		},
		{
			Name:        "getembedjson",
			Description: "Get an embed's JSON data",
			Arguments:   "<#channel> <message id>",
			Handler:     c.getembedjson,
		},
		{
			Name:        "sendmessagejson",
			Description: "Send a message using raw JSON data",
			Arguments:   "<#channel> <json>",
			Handler:     c.sendmessagejson,
		},
		{
			Name:        "sendinterval",
			Description: "Send a message every x seconds (flags: DELETE, INSTANT_DEL)",
			Arguments:   "<interval> <amount> <\"message\"> [flags]",
			Handler:     c.sendinterval,
		},
		{
			Name:        "optionpoll",
			Description: "Create a poll with options",
			Arguments:   "<\"question\"> <\"option1\"> <\"option2\">",
			Handler:     c.optionpoll,
		},
		{
			Name:        "discordpoll",
			Description: "Create a native Discord poll (options/emojis split by comma)",
			Arguments:   "<duration> <\"question\"> <\"options\"> [\"emojis\"]",
			Handler:     c.discordpoll,
		},
		{
			Name:        "checktoken",
			Description: "Check if a Discord token is valid",
			Arguments:   "<token>",
			Handler:     c.checktoken,
		},
		{
			Name:        "discordstatus",
			Description: "Check the Discord status",
			Handler:     c.discordstatus,
		},

		// LightningBot Utilities
		{
			Name:        "website",
			Description: "Open the LightningBot website",
			Flags:       []router.CommandFlag{router.CmdFlagHostingDangerous},
			Handler:     c.website,
		},
		{
			Name:        "dashboard",
			Description: "Open the LightningBot dashboard",
			Flags:       []router.CommandFlag{router.CmdFlagHostingDangerous},
			Handler:     c.dashboard,
		},
		{
			Name:        "changelog",
			Description: "See the latest changelog",
			Handler:     c.changelog,
		},
		{
			Name:        "version",
			Description: "See the current LightningBot version",
			Handler:     c.version,
		},
		{
			Name:        "uptime",
			Description: "See how long the bot has been running for",
			Handler:     c.uptime,
		},
		{
			Name:        "history",
			Description: "See used command history",
			Handler:     c.history,
		},
		{
			Name:        "restartbot",
			Description: "Restart LightningBot",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.restartbot,
		},
		{
			Name:        "stopbot",
			Description: "Stop LightningBot",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.stopbot,
		},
		{
			Name:        "cls",
			Description: "Clear UI console",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.cls,
		},
		{
			Name:        "openfolder",
			Description: "Open the LightningBot folder",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous, router.CmdFlagHostingDangerous},
			Handler:     c.openfolder,
		},
		{
			Name:        "remote",
			Description: "Run a command remotely",
			Arguments:   "<#channel> <command> [args...]",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.remote,
		},

		// System utilities
		{
			Name:        "systeminfo",
			Description: "See info about your system",
			Handler:     c.systeminfo,
		},
		{
			Name:        "shutdownpc",
			Description: "Shut your PC down",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous, router.CmdFlagHostingDangerous},
			Handler:     c.shutdownpc,
		},
		{
			Name:        "restartpc",
			Description: "Restart your PC",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous, router.CmdFlagHostingDangerous},
			Handler:     c.restartpc,
		},
		{
			Name:        "sleeppc",
			Description: "Put your PC to sleep",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous, router.CmdFlagHostingDangerous},
			Handler:     c.sleeppc,
		},
		{
			Name:        "hibernatepc",
			Description: "Hibernate your PC",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous, router.CmdFlagHostingDangerous},
			Handler:     c.hibernatepc,
		},
		{
			Name:        "lockpc",
			Description: "Lock your PC",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous, router.CmdFlagHostingDangerous},
			Handler:     c.lockpc,
		},
		{
			Name:        "activewindow",
			Description: "View the active window on your PC",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous, router.CmdFlagHostingDangerous},
			Handler:     c.activewindow,
		},
		{
			Name:        "calc",
			Description: "Open the calculator",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous, router.CmdFlagHostingDangerous},
			Handler:     c.calc,
		},
		{
			Name:        "notepad",
			Description: "Open notepad",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous, router.CmdFlagHostingDangerous},
			Handler:     c.notepad,
		},
		{
			Name:        "open",
			Description: "Open a URL in the browser",
			Arguments:   "<url>",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous, router.CmdFlagHostingDangerous},
			Handler:     c.open,
		},

		// Uploading utilities
		{
			Name:        "screenshot",
			Description: "Take a screenshot",
			Arguments:   "[delay]",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous, router.CmdFlagHostingDangerous},
			Handler:     c.screenshot,
		},
		{
			Name:        "listdestinations",
			Description: "List all upload destinations",
			Handler:     c.listdestinations,
		},
		{
			Name:        "uploadss",
			Description: "Upload a screenshot",
			Arguments:   "<destination> [delay]",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous, router.CmdFlagHostingDangerous},
			Handler:     c.uploadss,
		},
		{
			Name:        "uploadfile",
			Description: "Upload a file",
			Arguments:   "<destination> <attachment>",
			Handler:     c.uploadfile,
		},
		{
			Name:        "uploadlink",
			Description: "Upload a link",
			Arguments:   "<destination> <link>",
			Handler:     c.uploadlink,
		},

		// Paste utilities
		{
			Name:        "paste",
			Description: "Create a paste",
			Arguments:   "<text>",
			Handler:     c.pastebin,
		},
		{
			Name:        "pastefile",
			Description: "Create a paste from a file",
			Arguments:   "<filepath>",
			Handler:     c.pastefile,
		},
		{
			Name:        "tinyurl",
			Description: "Shorten a URL with TinyURL",
			Arguments:   "<url>",
			Handler:     c.tinyurl,
		},

		// Obfuscation utilities
		{
			Name:        "obfuscatepython",
			Description: "Obfuscate a Python script",
			Arguments:   "[script/attachment]",
			Handler:     c.obfuscateFunc("python", "py"),
		},
		{
			Name:        "obfuscatelua",
			Description: "Obfuscate a Lua script",
			Arguments:   "[script/attachment]",
			Handler:     c.obfuscateFunc("lua", "lua"),
		},
		{
			Name:        "obfuscatejs",
			Description: "Obfuscate a JS script",
			Arguments:   "[script/attachment]",
			Handler:     c.obfuscateFunc("javascript", "js"),
		},

		// ReviewDB utilities
		{
			Name:        "reviewdbtoken",
			Description: "Get your ReviewDB token",
			Handler:     c.reviewdbtoken,
		},
		{
			Name:        "getreviews",
			Description: "Get a users ReviewDB reviews",
			Arguments:   "<user>",
			Handler:     c.getreviews,
		},
		{
			Name:        "addreview",
			Description: "Add a ReviewDB review to a user",
			Arguments:   "<user> <review>",
			Handler:     c.addreview,
		},

		// Tag utilities
		{
			Name:        "tags",
			Description: "Manage tags",
			Arguments:   "<get/add/rm/clear/send> [name] [content]",
			SubCommands: []*router.Command{
				{
					Name:        "get",
					Description: "Get all tags",
					Handler:     c.getTags,
				},
				{
					Name:        "add",
					Description: "Add a tag",
					Arguments:   "<name> <content>",
					Handler:     c.addTag,
				},
				{
					Name:        "rm",
					Description: "Remove a tag",
					Arguments:   "<name>",
					Handler:     c.removeTag,
				},
				{
					Name:        "clear",
					Description: "Clear all tags",
					Handler:     c.clearTags,
				},
			},
		},

		// Creditcard commands
		{
			Name:        "mastercard",
			Description: "Generate a mastercard credit card",
			Arguments:   "[name]",
			Handler:     c.mastercard,
		},
		{
			Name:        "visa16",
			Description: "Generate a visa credit card (16 characters)",
			Arguments:   "[name]",
			Handler:     c.visa16,
		},
		{
			Name:        "visa13",
			Description: "Generate a visa credit card (13 characters)",
			Arguments:   "[name]",
			Handler:     c.visa13,
		},
		{
			Name:        "amex",
			Description: "Generate an amex credit card",
			Arguments:   "[name]",
			Handler:     c.amex,
		},
		{
			Name:        "discover",
			Description: "Generate a discover credit card",
			Arguments:   "[name]",
			Handler:     c.discover,
		},
		{
			Name:        "diners",
			Description: "Generate a diners credit card",
			Arguments:   "[name]",
			Handler:     c.diners,
		},
		{
			Name:        "jcb",
			Description: "Generate a JCB credit card",
			Arguments:   "[name]",
			Handler:     c.jcb,
		},
		// Network utilities

		{
			Name:        "ping",
			Description: "Check the latency of the bot",
			Handler:     c.ping,
		},
		{
			Name:        "pingweb",
			Description: "Check the latency to a website",
			Arguments:   "<host>",
			Handler:     c.pingweb,
		},
		{
			Name:        "ipinfo",
			Description: "Get info about an IP Address",
			Arguments:   "<ip>",
			Handler:     c.ipinfo,
		},
		{
			Name:        "portscan",
			Description: "Scan a host for open ports",
			Arguments:   "<host>",
			Handler:     c.portscan,
		},
		{
			Name:        "resolve",
			Description: "Get a domain's IP address",
			Arguments:   "<domain>",
			Handler:     c.resolve,
		},
		{
			Name:        "reversedns",
			Description: "Find all domains hosted on an IP address",
			Arguments:   "<ip>",
			Handler:     c.reversedns,
		},
		{
			Name:        "reversemac",
			Description: "Find the vender of a MAC address",
			Arguments:   "<mac>",
			Handler:     c.reversemac,
		},
		{
			Name:        "headers",
			Description: "Get URL response headers",
			Arguments:   "<url>",
			Handler:     c.headers,
		},
		{
			Name:        "interface",
			Description: "See your public/local ip and MAC address",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous, router.CmdFlagHostingDangerous},
			Handler:     c._interface,
		},

		// PokeTwo utilities
		{
			Name:        "snipepokemon",
			Description: "Snipe a PokéTwo Pokemon",
			Arguments:   "[#channel]",
			Handler:     c.snipepokemon,
		},
		{
			Name:        "detectpokemon",
			Description: "Detect a PokéTwo Pokemon",
			Arguments:   "[#channel]",
			Handler:     c.detectpokemon,
		},

		// Reminder utilities
		{
			Name:        "reminders",
			Description: "Manage your reminders",
			Arguments:   "<get/add/rm/clear> [name] [duration] [message]",
			SubCommands: []*router.Command{
				{
					Name:        "get",
					Description: "Get all reminders",
					Handler:     c.getReminders,
				},
				{
					Name:        "add",
					Description: "Add a reminder",
					Arguments:   "<name> <duration> <message>",
					Handler:     c.addReminder,
				},
				{
					Name:        "rm",
					Description: "Remove a reminder",
					Arguments:   "<name>",
					Handler:     c.removeReminder,
				},
				{
					Name:        "clear",
					Description: "Clear all reminders",
					Handler:     c.clearReminders,
				},
			},
		},

		// Conversion utilities
		{
			Name:        "convertvelocity",
			Description: "Convert velocity from one unit to another",
			Arguments:   "<value> <source unit> <destination unit>",
			Handler:     c.convertvelocity,
		},
		{
			Name:        "convertarea",
			Description: "Convert area from one unit to another",
			Arguments:   "<value> <source unit> <destination unit>",
			Handler:     c.convertarea,
		},
		{
			Name:        "convertlength",
			Description: "Convert length from one unit to another",
			Arguments:   "<value> <source unit> <destination unit>",
			Handler:     c.convertlength,
		},
		{
			Name:        "converttemperature",
			Description: "Convert temperature from one unit to another",
			Arguments:   "<value> <source unit> <destination unit>",
			Handler:     c.converttemperature,
		},
		{
			Name:        "converttime",
			Description: "Convert time from one unit to another",
			Arguments:   "<value> <source unit> <destination unit>",
			Handler:     c.converttime,
		},
		{
			Name:        "convertvolume",
			Description: "Convert volume from one unit to another",
			Arguments:   "<value> <source unit> <destination unit>",
			Handler:     c.convertvolume,
		},
		{
			Name:        "convertmass",
			Description: "Convert mass from one unit to another",
			Arguments:   "<value> <source unit> <destination unit>",
			Handler:     c.convertmass,
		},

		// Misc utilities
		{
			Name:        "tts",
			Description: "Send a Text-To-Speech message",
			Arguments:   "<language> <text>",
			Handler:     c.tts,
		},
		{
			Name:        "translate",
			Description: "Translate text to a language",
			Arguments:   "<language> <text>",
			Handler:     c.translate,
		},
		{
			Name:        "detectlang",
			Description: "Detect a texts language",
			Arguments:   "<text>",
			Handler:     c.detectlang,
		},
		{
			Name:        "english",
			Description: "Translate text to English",
			Arguments:   "<text>",
			Handler:     c.english,
		},
		{
			Name:        "define",
			Description: "Define a word",
			Arguments:   "<word>",
			Handler:     c.define,
		},
		{
			Name:        "urbandictionary",
			Description: "Search someting on the Urban Dictionary",
			Arguments:   "<word>",
			Handler:     c.urbandictionary,
		},
		{
			Name:        "wikipedia",
			Description: "Search someting on Wikipedia",
			Arguments:   "<query>",
			Handler:     c.wikipedia,
		},
		{
			Name:        "weather",
			Description: "Check the weather in a city",
			Arguments:   "<city>",
			Handler:     c.weather,
		},
		{
			Name:        "lvbypass",
			Description: "Bypass a Linkvertise link",
			Arguments:   "<url>",
			Handler:     c.lvbypass,
		},
		{
			Name:        "flightinfo",
			Description: "Get info about a flight",
			Arguments:   "<flight number>",
			Handler:     c.flightinfo,
		},
		{
			Name:        "passgen",
			Description: "Generate a password",
			Arguments:   "<character amount>",
			Handler:     c.passgen,
		},
		{
			Name:        "tokengen",
			Description: "Generate a Discord token",
			Handler:     c.tokengen,
		},
		{
			Name:        "revimage",
			Description: "Reverse image search",
			Arguments:   "<image url>",
			Handler:     c.revimage,
		},
		{
			Name:        "revavatar",
			Description: "Reverse image search a users avatar",
			Arguments:   "<user>",
			Handler:     c.revavatar,
		},
		{
			Name:        "timezones",
			Description: "Get a list of timezones for a continent",
			Arguments:   "<continent>",
			Handler:     c.timezones,
		},
		{
			Name:        "timezone",
			Description: "Get information about a timezone",
			Arguments:   "<timezone>",
			Handler:     c.timezoneinfo,
		},
		{
			Name:        "calculate",
			Description: "Calculate something",
			Arguments:   "<expression>",
			Handler:     c.calculate,
		},
		{
			Name:        "googlesearch",
			Description: "Search someting on Google",
			Arguments:   "<query>",
			Handler:     c.googlesearch,
		},
		{
			Name:        "youtubesearch",
			Description: "Search someting on Youtube",
			Arguments:   "<query>",
			Handler:     c.youtubesearch,
		},
		{
			Name:        "amazon",
			Description: "Search products on Amazon",
			Arguments:   "<query>",
			Handler:     c.amazon,
		},
		{
			Name:        "steam",
			Description: "Search games on Steam",
			Arguments:   "<query>",
			Handler:     c.steam,
		},
		{
			Name:        "github",
			Description: "Search repos on GitHub",
			Arguments:   "<query>",
			Handler:     c.github,
		},
		{
			Name:        "stackoverflow",
			Description: "Search answers on Stack Overflow",
			Arguments:   "<query>",
			Handler:     c.stackoverflow,
		},
		{
			Name:        "searchcountry",
			Description: "Get info about a country",
			Arguments:   "<country>",
			Handler:     c.searchcountry,
		},
		{
			Name:        "randomcountry",
			Description: "Get info about a random country",
			Arguments:   "<country>",
			Handler:     c.randomcountry,
		},
		{
			Name:        "loremipsum",
			Description: "Generate a random placeholder text",
			Arguments:   "[paragraphs]",
			Handler:     c.loremipsum,
		},
		{
			Name:        "baconipsum",
			Description: "Generate a meaty placeholder text",
			Arguments:   "[paragraphs]",
			Handler:     c.baconipsum,
		},
		{
			Name:        "previewcolor",
			Description: "Preview a color",
			Arguments:   "<rgb/hex>",
			Handler:     c.previewcolor,
		},
		{
			Name:        "generatekey",
			Description: "Generate a random license key (X or x for letters, 9 for numbers)",
			Arguments:   "<format> [amount]",
			Handler:     c.generatekey,
		},
		{
			Name:        "timestamp",
			Description: "Create a Discord timestamp",
			Arguments:   "<time>",
			Handler:     c.timestampCommand(""),
		},
		{
			Name:        "reltimestamp",
			Description: "Create a relative Discord timestamp",
			Arguments:   "<time>",
			Handler:     c.timestampCommand(":R"),
		},
	})
}

// Discord utilities

func getChannelType(channelType discordgo.ChannelType) discordgo.ChannelType {
	if channelType == discordgo.ChannelTypeGuildVoice || channelType == discordgo.ChannelTypeGuildStageVoice {
		return discordgo.ChannelTypeGuildVoice
	}
	return discordgo.ChannelTypeGuildText
}

func (c *UtilityCommands) copyserver(ctx *router.Ctx) {
	duration, err := ctx.Arguments.Get(1).Duration()
	if err != nil || duration < time.Second*10 {
		duration = time.Second * 10
	}

	guild, err := ctx.Arguments.Get(0).Guild()
	if err != nil {
		ctx.RespondText("Error getting server: %v", err)
		return
	}

	newGuild, err := ctx.Session.GuildCreate(fmt.Sprintf("[COPY] %.90s", guild.Name))
	if err != nil {
		ctx.RespondText("Error copying server: %v", err)
		return
	}

	time.Sleep(duration)

	resp, _ := http.DefaultClient().Get(guild.IconURL())
	icon := utils.Misc.ImageToBase64(resp.MustGetBody())

	_, _ = ctx.Session.GuildEdit(newGuild.ID, &discordgo.GuildParams{
		VerificationLevel:           &guild.VerificationLevel,
		DefaultMessageNotifications: int(guild.DefaultMessageNotifications),
		AfkChannelID:                guild.AfkChannelID,
		AfkTimeout:                  guild.AfkTimeout,
		Icon:                        icon,
	})

	time.Sleep(duration)

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
			ctx.RespondText("Error copying server: %v", err)
			return
		}

		IDs[role.ID] = newRole.ID

		time.Sleep(duration)
	}

	newChannels, err := ctx.Session.GuildChannels(newGuild.ID)
	if err != nil {
		ctx.RespondText("Error copying server: %v", err)
		return
	}

	time.Sleep(duration)

	for _, channel := range newChannels {
		_, err := ctx.Session.ChannelDelete(channel.ID)
		if err != nil {
			ctx.RespondText("Error copying server: %v", err)
			return
		}

		time.Sleep(duration)
	}

	channels, err := ctx.Session.GuildChannels(guild.ID)
	if err != nil {
		ctx.RespondText("Error copying server: %v", err)
		return
	}

	time.Sleep(duration)

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
			ctx.RespondText("Error copying server: %v", err)
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
			ctx.RespondText("Error copying server: %v", err)
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
				ctx.RespondText("Error copying server: %v", err)
				return
			}

			time.Sleep(duration)
		}
	}
	c.services.Embed.New(ctx).
		SetDescription("Server copied successfully!").
		Send()
}

func (c *UtilityCommands) copychannel(ctx *router.Ctx) {
	sourceChannel, err := ctx.Arguments.Get(0).Channel()
	if err != nil {
		ctx.RespondText("Could not find channel")
		return
	}

	destChannel, err := ctx.Arguments.Get(1).Channel()
	if err != nil {
		ctx.RespondText("Could not find channel")
		return
	}

	limit, _ := ctx.Arguments.Get(2).Int()
	limit = fn.Clamp(limit, 1, 100)

	if limit < 1 {
		ctx.RespondText("Limit must be greater than 0")
		return
	}

	messages, err := utils.Discord.GetChannelMessages(ctx.Session, sourceChannel.ID, limit)
	if err != nil {
		ctx.RespondText("Error getting source messages: %v", err)
		return
	}

	slices.SortFunc(messages, func(i, j *discordgo.Message) int {
		timeI := utils.Discord.SnowflakeToTime(i.ID)
		timeJ := utils.Discord.SnowflakeToTime(j.ID)

		if timeI.Before(timeJ) {
			return -1
		} else if timeI.After(timeJ) {
			return 1
		}
		return 0
	})

	ctx.RespondText(
		"Copying %d messages...\n\nEstimated time required: %s",
		len(messages),
		time.Duration(len(messages)*2)*time.Second,
	)

	webhook, err := ctx.Session.WebhookCreate(destChannel.ID, "LightningBot - Copy Channel", "")
	if err != nil {
		ctx.RespondText("Error creating destination webhook: %v", err)
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
			ctx.RespondText("Error sending message: %v", err)
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

func (c *UtilityCommands) leavegroups(ctx *router.Ctx) {
	isSilent, _ := ctx.Arguments.Get(0).Bool()

	for _, channel := range ctx.Session.State.PrivateChannels {
		if channel.Type == discordgo.ChannelTypeGroupDM {
			_, err := ctx.Session.ChannelDelete(channel.ID, isSilent)
			if err != nil {
				ctx.RespondText("Error leaving group: %v", err)
			}

			time.Sleep(2 * time.Second)
		}
	}
}

func (c *UtilityCommands) adminservers(ctx *router.Ctx) {
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

	c.services.Embed.New(ctx).
		AddField("Admin Servers", adminServerString).
		Send()
}

func (c *UtilityCommands) admingroups(ctx *router.Ctx) {
	adminGroupsString := ""
	for _, channel := range ctx.Session.State.PrivateChannels {
		if channel.Type == discordgo.ChannelTypeGroupDM {
			if channel.OwnerID == ctx.Session.State.User.ID {
				adminGroupsString += channel.Name + "\n"
			}
		}
	}
	c.services.Embed.New(ctx).
		AddField("Admin Groups", adminGroupsString).
		Send()
}

func (c *UtilityCommands) joinvc(ctx *router.Ctx) {
	channel, err := ctx.Arguments.Get(0).Channel()
	if err != nil {
		ctx.RespondText("Could not find channel")
		return
	}

	_, err = ctx.Session.ChannelVoiceJoin(channel.GuildID, channel.ID, false, false)
	if err != nil {
		ctx.RespondText("Could not join channel")
		return
	}

	ctx.RespondText("%s", "Joined channel "+channel.Name)
}

func (c *UtilityCommands) snipe(ctx *router.Ctx) {
	channel, err := ctx.Arguments.Get(0).Channel(router.ArgumentFlagFallbackDefault)
	if err != nil {
		ctx.RespondText("Could not find channel")
		return
	}

	messages := c.services.EventReactions.GetDeletedMessages(channel.ID)
	if len(messages) == 0 {
		c.services.Embed.New(ctx).
			SetDescription("No deleted message was found in this channel!").
			Send()
		return
	}

	formatted := ""

	for i, message := range messages {
		if i > 5 {
			break
		}

		var attachments []string

		for _, attachment := range message.Attachments {
			attachments = append(attachments, utils.Misc.Or(attachment.URL, attachment.ProxyURL))
		}

		var fmtAttachments string
		if len(attachments) > 0 {
			fmtAttachments = "\n" + strings.Join(attachments, " , ")
		}

		formatted += fmt.Sprintf(
			"[%s] [%s]: %s%s\n\n",
			utils.Discord.SnowflakeToTime(message.ID).Format(time.RFC1123),
			message.Author.String()+" | "+message.Author.ID,
			message.Content,
			fmtAttachments,
		)
	}

	c.services.Embed.New(ctx).
		AddField("Deleted Messages", formatted).
		Send()
}

func (c *UtilityCommands) editsnipe(ctx *router.Ctx) {
	channel, err := ctx.Arguments.Get(0).Channel(router.ArgumentFlagFallbackDefault)
	if err != nil {
		ctx.RespondText("Could not find channel")
		return
	}

	messages := c.services.EventReactions.GetEditedMessages(channel.ID)
	if len(messages) == 0 {
		c.services.Embed.New(ctx).
			SetDescription("No deleted message was found in this channel!").
			Send()
		return
	}

	message := messages[len(messages)-1]
	emb := c.services.Embed.New(ctx)
	if message.Author != nil {
		emb.AddField("Author", message.Author.String()+" | "+message.Author.ID)
	}
	emb.AddField("Sent At", utils.Discord.SnowflakeToTime(message.ID).Format(time.RFC1123))
	emb.AddField("Old Content", message.BeforeUpdate.Content)
	emb.AddField("New Content", message.Content)
	emb.Send()
}

func (c *UtilityCommands) firstmessage(ctx *router.Ctx) {
	channel, err := ctx.Arguments.Get(0).Channel(router.ArgumentFlagFallbackDefault)
	if err != nil {
		ctx.RespondText("Could not find channel")
		return
	}
	messages, err := ctx.Session.ChannelMessages(channel.ID, 1, "", "1", "")
	if err != nil {
		ctx.RespondText("Error getting messages")
		return
	}
	if len(messages) != 1 {
		ctx.RespondText("Error getting messages")
		return
	}
	c.services.Embed.New(ctx).
		SetDescription("First Message").
		AddField("Author", messages[0].Author.String()).
		AddField("Message", messages[0].Content).
		AddField("Jump URL", utils.Discord.GetMessageLink(messages[0])).
		Send()
}

func (c *UtilityCommands) hypesquad(ctx *router.Ctx) {
	houseName := strings.ToLower(ctx.Arguments.Get(0).Raw())
	house, ok := map[string]int{
		"bravery":    1,
		"brilliance": 2,
		"balance":    3,
	}[houseName]
	if !ok {
		ctx.RespondText("Missing required argument: <bravery/brilliance/balance>")
		return
	}

	if err := c.services.Discord.SetHypesquad(house); err != nil {
		ctx.RespondText("Error setting hypesquad: %v", err)
	}

	c.services.Embed.New(ctx).
		SetDescription("Hypesquad set successfully!").
		Send()
}

func (c *UtilityCommands) leavehypesquad(ctx *router.Ctx) {
	if err := c.services.Discord.LeaveHypesquad(); err != nil {
		ctx.RespondText("Error leaving hypesquad: %v", err)
	}

	c.services.Embed.New(ctx).
		SetDescription("Left hypesquad successfully!").
		Send()
}

func (c *UtilityCommands) theme(ctx *router.Ctx) {
	themeName := strings.ToLower(ctx.Arguments.Get(0).Raw())
	theme, ok := map[string]protos.PreloadedUserSettings_Theme{
		"dark":     1,
		"light":    2,
		"darker":   3,
		"midnight": 4,
	}[themeName]
	if !ok {
		ctx.RespondText("Missing required argument: <dark/light/darker/midnight>")
		return
	}

	ctx.Session.State.UserSettings.Appearance.Theme = theme
	err := ctx.Session.UserSettingsUpdate(ctx.Session.State.UserSettings, true)
	if err != nil {
		ctx.RespondText("Error setting theme: %v", err)
		return
	}

	c.services.Embed.New(ctx).
		SetDescription("Theme set successfully!").
		Send()
}

func (c *UtilityCommands) compact(ctx *router.Ctx) {
	toggle, err := ctx.Arguments.Get(0).Bool()
	if err != nil {
		ctx.RespondText("Error parsing argument: %v", err)
		return
	}

	ctx.Session.State.UserSettings.TextAndImages.MessageDisplayCompact = wrapperspb.Bool(toggle)
	err = ctx.Session.UserSettingsUpdate(ctx.Session.State.UserSettings, true)
	if err != nil {
		ctx.RespondText("Error setting theme: %v", err)
		return
	}

	c.services.Embed.New(ctx).
		SetDescription("Compact mode set successfully!").
		Send()
}

func (c *UtilityCommands) devmode(ctx *router.Ctx) {
	toggle, err := ctx.Arguments.Get(0).Bool()
	if err != nil {
		ctx.RespondText("Error parsing argument: %v", err)
		return
	}

	ctx.Session.State.UserSettings.Appearance.DeveloperMode = toggle
	err = ctx.Session.UserSettingsUpdate(ctx.Session.State.UserSettings, true)
	if err != nil {
		ctx.RespondText("Error setting theme: %v", err)
		return
	}

	c.services.Embed.New(ctx).
		SetDescription("Developer mode set successfully!").
		Send()
}

func (c *UtilityCommands) sessionid(ctx *router.Ctx) {
	c.services.Embed.New(ctx).
		AddField("Session ID", ctx.Session.State.SessionID).
		Send()
}

func (c *UtilityCommands) aboutme(ctx *router.Ctx) {
	bio := any(ctx.Arguments.Raw())
	if bio == "" {
		bio = nil
	}

	err := c.services.Discord.UpdateSettings("/profile", map[string]any{"bio": bio})
	if err != nil {
		ctx.RespondText("Error setting bio: %v", err)
		return
	}

	c.services.Embed.New(ctx).
		SetDescription("Bio set successfully!").
		Send()
}

func (c *UtilityCommands) displayname(ctx *router.Ctx) {
	globalName := any(ctx.Arguments.Raw())
	if globalName == "" {
		globalName = nil
	}

	err := c.services.Discord.UpdateSettings("", map[string]any{"global_name": globalName})
	if err != nil {
		ctx.RespondText("Error setting display name: %v", err)
		return
	}

	c.services.Embed.New(ctx).
		SetDescription("Display name set successfully!").
		Send()
}

func (c *UtilityCommands) pronouns(ctx *router.Ctx) {
	err := c.services.Discord.UpdateSettings("/profile", map[string]any{"pronouns": ctx.Arguments.Raw()})
	if err != nil {
		ctx.RespondText("Error setting pronouns: %v", err)
		return
	}

	c.services.Embed.New(ctx).
		SetDescription("Pronouns set successfully!").
		Send()
}

func (c *UtilityCommands) nickscan(ctx *router.Ctx) {
	nickText := ""
	for _, guild := range ctx.Session.State.Guilds {
		me, err := ctx.Session.State.Member(guild.ID, ctx.Message.Author.ID)
		if err != nil {
			me, err = ctx.Session.GuildMember(guild.ID, ctx.Message.Author.ID)
			if err != nil {
				continue
			}
		}
		if me.Nick != "" {
			nickText += fmt.Sprintf("%s: %s\n", guild.Name, me.Nick)
		}
	}

	c.services.Embed.New(ctx).
		AddField("Nicknames", nickText).
		Send()
}

func (c *UtilityCommands) getmessagejson(ctx *router.Ctx) {
	channel, err := ctx.Arguments.Get(0).Channel()
	if err != nil {
		ctx.RespondText("Could not find channel")
		return
	}

	messages, err := ctx.Session.ChannelMessages(channel.ID, 1, "", "", ctx.Arguments.Get(1).Raw())
	if err != nil {
		ctx.RespondText("Could not get message: %v", err)
		return
	}

	if len(messages) == 0 {
		ctx.RespondText("Could not find message")
		return
	}

	message := messages[0]

	jsonData, err := json.MarshalIndent(message, "", "  ")
	if err != nil {
		ctx.RespondText("Error marshaling JSON")
		return
	}

	c.services.Embed.New(ctx).
		SetTitle("Message JSON").
		SetDescription("%s", string(jsonData)).
		Send()
}

func (c *UtilityCommands) getembedjson(ctx *router.Ctx) {
	channel, err := ctx.Arguments.Get(0).Channel()
	if err != nil {
		ctx.RespondText("Could not find channel")
		return
	}

	messages, err := ctx.Session.ChannelMessages(channel.ID, 1, "", "", ctx.Arguments.Get(1).Raw())
	if err != nil {
		ctx.RespondText("Could not get message: %v", err)
		return
	}

	if len(messages) == 0 {
		ctx.RespondText("Could not find message")
		return
	}

	message := messages[0]

	jsonData, err := json.MarshalIndent(message.Embeds, "", "  ")
	if err != nil {
		ctx.RespondText("Error marshaling JSON")
		return
	}

	c.services.Embed.New(ctx).
		SetTitle("Message JSON").
		SetDescription("%s", string(jsonData)).
		Send()
}

func (c *UtilityCommands) sendmessagejson(ctx *router.Ctx) {
	channel, err := ctx.Arguments.Get(0).Channel()
	if err != nil {
		ctx.RespondText("Could not find channel")
		return
	}

	rawData := ctx.Arguments.RawFrom(1)

	var data *discordgo.MessageSend
	err = json.Unmarshal([]byte(rawData), &data)
	if err != nil {
		ctx.RespondText("Error unmarshaling JSON: %v", err)
		return
	}

	_, err = ctx.Session.ChannelMessageSendComplex(channel.ID, data)
	if err != nil {
		ctx.RespondText("Error sending message: %v", err)
		return
	}
}

func (c *UtilityCommands) sendinterval(ctx *router.Ctx) {
	interval, err := ctx.Arguments.Get(0).Duration()
	if err != nil {
		ctx.RespondText("Error parsing interval: %v", err)
		return
	}

	amount, err := ctx.Arguments.Get(1).Int()
	if err != nil {
		ctx.RespondText("Error parsing amount: %v", err)
		return
	}

	if amount < 1 {
		ctx.RespondText("Error parsing amount: amount must be at least 1")
		return
	}

	message, flags := ctx.Arguments.Get(2).Raw(), ctx.Arguments.Get(3).Raw()
	instantDelete := strings.Contains(flags, "INSTANT_DEL")
	delete := strings.Contains(flags, "DELETE")

	for range amount {
		msg, err := ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, message)
		if err != nil {
			return
		}

		if instantDelete {
			_ = ctx.Session.ChannelMessageDelete(ctx.Message.ChannelID, msg.ID)
		}

		time.Sleep(interval)

		if delete {
			_ = ctx.Session.ChannelMessageDelete(ctx.Message.ChannelID, ctx.Message.ID)
		}
	}
}

func (c *UtilityCommands) optionpoll(ctx *router.Ctx) {
	messages := c.services.Embed.New(ctx).
		SetDescription("%s", "Poll: "+ctx.Arguments.Get(0).Raw()).
		AddField("👍", ctx.Arguments.Get(0).Raw()).
		AddField("👎", ctx.Arguments.Get(0).Raw()).
		Send()

	if len(messages) > 0 {
		_ = ctx.Session.MessageReactionAdd(ctx.Message.ChannelID, messages[0].ID, "👍")
		_ = ctx.Session.MessageReactionAdd(ctx.Message.ChannelID, messages[0].ID, "👎")
	}
}

func (c *UtilityCommands) discordpoll(ctx *router.Ctx) {
	duration, err := ctx.Arguments.Get(0).Duration()
	if err != nil {
		ctx.RespondText("Error parsing duration: %v", err)
		return
	}

	duration = fn.Clamp(duration, 1*time.Hour, 168*time.Hour)

	question := ctx.Arguments.Get(1).Raw()
	options := strings.Split(ctx.Arguments.Get(2).Raw(), ",")
	emojis := strings.Split(ctx.Arguments.Get(3).Raw(), ",")

	var pollAnswers []discordgo.PollAnswer

	for i, option := range options {
		if option == "" {
			continue
		}

		answer := discordgo.PollAnswer{
			Media: discordgo.PollMedia{Text: option},
		}

		if i < len(emojis) && emojis[i] != "" {
			emojiID := ""
			emojiName := strings.TrimSpace(emojis[i])

			emojiIDs := regexp.MustCompile(`<a?:.*:(\d+)>`).FindStringSubmatch(emojiName)
			if len(emojiIDs) > 0 {
				emojiID = emojiIDs[1]
				emojiName = ""
			}

			answer.Media.Emoji = &discordgo.PollEmoji{
				ID:   emojiID,
				Name: emojiName,
			}
		}

		pollAnswers = append(pollAnswers, answer)
	}

	if len(pollAnswers) == 0 {
		ctx.RespondText("Error parsing options: At least one option must be provided.")
		return
	}

	_, err = ctx.Session.ChannelMessageSendComplex(ctx.Message.ChannelID, &discordgo.MessageSend{
		Poll: &discordgo.Poll{
			Question:         discordgo.PollQuestion{Text: question},
			Answers:          pollAnswers,
			AllowMultiselect: false,
			Duration:         int(duration.Hours()),
			LayoutType:       1,
		},
	})
	if err != nil {
		ctx.RespondText("Error sending poll: %v", err)
	}
}

func (c *UtilityCommands) checktoken(ctx *router.Ctx) {
	token := ctx.Arguments.Get(0).Raw()

	tokenData, err := c.services.Discord.GetTokenInfo(token)
	if err != nil {
		ctx.RespondText("An error has occurred: %v", err)
		return
	}

	c.services.Embed.New(ctx).
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
		Send()
}

func (c *UtilityCommands) discordstatus(ctx *router.Ctx) {
	resp, err := http.DefaultClient().Get("https://srhpyqt94yxb.statuspage.io/api/v2/summary.json")
	if err != nil {
		ctx.RespondText("An error has occurred.")
		return
	}

	components := resp.Json().Get("components")

	emb := c.services.Embed.New(ctx)

	for component := range components.Iterator() {
		emb.AddField(
			component.Get("name").MustString(),
			component.Get("status").MustString(),
		)
	}

	emb.Send()
}

// LightningBot Utilities

func (c *UtilityCommands) website(*router.Ctx) {
	_ = system.OpenURL("https://lightning-bot.com/")
}

func (c *UtilityCommands) dashboard(*router.Ctx) {
	_ = system.OpenURL("https://lightning-bot.com/dashboard")
}

func (c *UtilityCommands) changelog(ctx *router.Ctx) {
	resp, err := http.DefaultClient().Post(
		"https://api.LBlightning-bot.com/api/v3/selfbot/check-version",
		map[string]string{"branch": "selfbot", "version": ""},
	)
	if err != nil {
		ctx.RespondText("An error has occurred: %v", err)
		return
	}

	c.services.Embed.New(ctx).
		AddField("Latest Version", resp.Json().Get("latest", "version").MustString()).
		AddField("Changelog", resp.Json().Get("latest", "changelog").MustString()).
		Send()
}

func (c *UtilityCommands) version(ctx *router.Ctx) {
	resp, err := http.DefaultClient().Post(
		"https://api.LBlightning-bot.com/api/v3/selfbot/check-version",
		map[string]string{"branch": "selfbot", "version": ""},
	)
	if err != nil {
		ctx.RespondText("An error has occurred: %v", err)
		return
	}

	c.services.Embed.New(ctx).
		AddField("Current Version", constants.VERSION[1:]).
		AddField("Latest Version", resp.Json().Get("latest", "version").MustString()).
		Send()
}

func (c *UtilityCommands) uptime(ctx *router.Ctx) {
	c.services.Embed.New(ctx).
		AddField("Uptime", utils.Misc.FormatDuration(time.Since(constants.START_TIME))).
		Send()
}

func (c *UtilityCommands) history(ctx *router.Ctx) {
	usedCmds := fn.Limit(c.services.EventReactions.GetUsedCommands(), 10)

	if len(usedCmds) == 0 {
		c.services.Embed.New(ctx).
			SetDescription("No commands have been used yet.").
			Send()
		return
	}

	usedCmdsReadable := fn.MapIndexed(usedCmds, func(i int, cmd *router.Ctx) string {
		return fmt.Sprintf("%d. %s", i+1, cmd.Message.Content)
	})

	c.services.Embed.New(ctx).
		AddField("Used Commands", strings.Join(usedCmdsReadable, "\n")).
		Send()
}

func (c *UtilityCommands) restartbot(*router.Ctx) {
	if constants.CONSOLE {
		_ = system.RestartExecution()
	} else {
		c.services.EventBroker.Publish(eventbroker.EventRestart, nil)
	}
}

func (c *UtilityCommands) stopbot(*router.Ctx) {
	os.Exit(0)
}

func (c *UtilityCommands) cls(*router.Ctx) {
	c.services.UI.ClearConsole()
	c.services.UI.ShowBanner()
}

func (c *UtilityCommands) openfolder(*router.Ctx) {
	_ = system.OpenURL(utils.IO.GetMainFolder())
}

func (c *UtilityCommands) remote(ctx *router.Ctx) {
	channel, err := ctx.Arguments.Get(0).Channel(router.ArgumentFlagFallbackDefault)
	if err != nil {
		ctx.RespondText("Could not find channel")
		return
	}

	commandName := ctx.Arguments.Get(1).Raw()
	args := ctx.Arguments.RawFrom(2)

	cmd := ctx.Router.GetCommand(commandName)
	if cmd == nil {
		ctx.RespondText("Command not found")
		return
	}

	ctx.Router.InvokeRemoteCmd(ctx.Session, channel, ctx.Message.Author, cmd, args)
}

// System Utilities

func (c *UtilityCommands) systeminfo(ctx *router.Ctx) {
	path := "/"
	if runtime.GOOS == "windows" {
		path = "C:"
	}

	hostStat, err := host.Info()
	if err != nil {
		ctx.RespondText("Error getting host info: %v", err)
		return
	}

	cpuStat, err := cpu.Info()
	if err != nil {
		ctx.RespondText("Error getting cpu info: %v", err)
		return
	}

	if len(cpuStat) < 1 {
		ctx.RespondText("Error getting cpu info: No CPUs found")
		return
	}

	cpuInfo := cpuStat[0]

	vmStat, err := mem.VirtualMemory()
	if err != nil {
		ctx.RespondText("Error getting memory info: %v", err)
		return
	}

	diskStat, err := disk.Usage(path)
	if err != nil {
		ctx.RespondText("Error getting disk info: %v", err)
		return
	}

	c.services.Embed.New(ctx).
		SetDescription("System Information").
		AddField("OS", utils.String.Capitalize(hostStat.OS)).
		AddField("Platform", fmt.Sprintf("%s - %s", hostStat.Platform, hostStat.PlatformVersion)).
		AddField("Kernel Arch and Version", fmt.Sprintf("%s - %s", hostStat.KernelArch, hostStat.KernelVersion)).
		AddField("CPU Info", fmt.Sprintf("Model: %s\nFrequency: %.f MHz\nCores: %d\nVendor: %s", cpuInfo.ModelName, cpuInfo.Mhz, cpuInfo.Cores, cpuInfo.VendorID)).
		AddField("Memory Info", fmt.Sprintf("Total: %d MB\nFree: %d MB\nUsed: %d MB\nUsed Percentage: %.2f%%", vmStat.Total/1024/1024, vmStat.Free/1024/1024, vmStat.Used/1024/1024, vmStat.UsedPercent)).
		AddField("Disk Info", fmt.Sprintf("Total: %d MB\nFree: %d MB\nUsed: %d MB\nUsed Percentage: %.2f%%", diskStat.Total/1024/1024, diskStat.Free/1024/1024, diskStat.Used/1024/1024, diskStat.UsedPercent)).
		Send()
}

func (c *UtilityCommands) shutdownpc(ctx *router.Ctx) {
	_ = system.Shutdown()
}

func (c *UtilityCommands) restartpc(ctx *router.Ctx) {
	_ = system.Restart()
}

func (c *UtilityCommands) sleeppc(ctx *router.Ctx) {
	_ = system.Sleep()
}

func (c *UtilityCommands) hibernatepc(ctx *router.Ctx) {
	_ = system.Hibernate()
}

func (c *UtilityCommands) lockpc(ctx *router.Ctx) {
	_ = system.Lock()
}

func (c *UtilityCommands) activewindow(ctx *router.Ctx) {
	pid, _ := foreground.GetForegroundPID()
	title, _ := foreground.GetForegroundTitle()

	c.services.Embed.New(ctx).
		AddField("Active Window PID", fmt.Sprint(pid)).
		AddField("Active Window Title", title).
		Send()
}

func (c *UtilityCommands) calc(ctx *router.Ctx) {
	if runtime.GOOS != "windows" {
		ctx.RespondText("This command is only available on Windows")
	}
	_ = exec.Command("C:\\Windows\\System32\\calc.exe").Start()
}

func (c *UtilityCommands) notepad(ctx *router.Ctx) {
	if runtime.GOOS != "windows" {
		ctx.RespondText("This command is only available on Windows")
	}
	_ = exec.Command("C:\\Windows\\System32\\notepad.exe").Start()
}

func (c *UtilityCommands) open(ctx *router.Ctx) {
	_ = system.OpenURL(ctx.Arguments.Raw())
}

// Uploading utilities

func (c *UtilityCommands) screenshot(ctx *router.Ctx) {
	delay, _ := ctx.Arguments.Get(0).Duration()
	time.Sleep(delay)

	img, err := screenshot.CaptureDisplay(0)
	if err != nil {
		ctx.RespondText("Error capturing screenshot: %v", err)
		return
	}

	buf := new(bytes.Buffer)

	err = png.Encode(buf, img)
	if err != nil {
		ctx.RespondText("Error encoding image: %v", err)
		return
	}

	_, _ = ctx.Session.ChannelFileSend(ctx.Message.ChannelID, "screenshot.png", buf)
}

func (c *UtilityCommands) listdestinations(ctx *router.Ctx) {
	destinationsString := []string{}

	for name, dest := range destinations.Destinations {
		desc := fmt.Sprintf("`%s`  - Supports all file types: **%s**", name, utils.Misc.FormatBoolYesNo(dest.SupportsAllFiles()))

		destinationsString = append(destinationsString, desc)
	}

	sort.Strings(destinationsString)

	c.services.Embed.New(ctx).
		AddField("Destinations", strings.Join(destinationsString, "\n")).
		Send()
}

func (c *UtilityCommands) uploadss(ctx *router.Ctx) {
	destName := strings.ToLower(ctx.Arguments.Get(0).Raw())
	dest, ok := destinations.Destinations[destName]
	if !ok {
		ctx.RespondText("Error: invalid destination\n\nUse the `listdestinations` command to view all destinations")
		return
	}

	delay, _ := ctx.Arguments.Get(1).Duration()
	time.Sleep(delay)

	img, err := screenshot.CaptureDisplay(0)
	if err != nil {
		ctx.RespondText("Error capturing screenshot: %v", err)
		return
	}

	buf := new(bytes.Buffer)

	err = png.Encode(buf, img)
	if err != nil {
		ctx.RespondText("Error encoding image: %v", err)
		return
	}

	result, err := dest.Upload(buf, "screenshot.png")
	if err != nil {
		ctx.RespondText("Error uploading screenshot: %v", err)
		return
	}

	c.services.Embed.New(ctx).
		SetDescription("Screenshot").
		AddField("URL", result.URL).
		AddField("Direct URL", result.DirectURL).
		Send()
}

func (c *UtilityCommands) uploadfile(ctx *router.Ctx) {
	destName := strings.ToLower(ctx.Arguments.Get(0).Raw())
	dest, ok := destinations.Destinations[destName]
	if !ok {
		ctx.RespondText("Error: invalid destination\n\nUse the `listdestinations` command to view all destinations")
		return
	}

	attachments := ctx.Message.Attachments

	if len(attachments) == 0 {
		ctx.RespondText("Error: you must provide a file as an attachment.")
		return
	}

	resp, err := nethttp.Get(attachments[0].URL)
	if err != nil {
		ctx.RespondText("Error downloading file: %v", err)
		return
	}
	defer resp.Body.Close()

	buf := new(bytes.Buffer)

	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		ctx.RespondText("Error reading file body: %v", err)
		return
	}

	result, err := dest.Upload(buf, attachments[0].Filename)
	if err != nil {
		ctx.RespondText("Error uploading file: %v", err)
		return
	}

	c.services.Embed.New(ctx).
		SetDescription("Uploaded file: %s", attachments[0].Filename).
		AddField("URL", result.URL).
		AddField("Direct URL", result.DirectURL).
		Send()
}

func (c *UtilityCommands) uploadlink(ctx *router.Ctx) {
	destName := strings.ToLower(ctx.Arguments.Get(0).Raw())
	dest, ok := destinations.Destinations[destName]
	if !ok {
		ctx.RespondText("Error: invalid destination\n\nUse the `listdestinations` command to view all destinations")
		return
	}

	link := ctx.Arguments.Get(1).Raw()
	parsedURL, err := url.ParseRequestURI(link)
	if err != nil {
		ctx.RespondText("Error: invalid link: %v", err)
		return
	}

	fileName := path.Base(parsedURL.Path)

	resp, err := nethttp.Get(link)
	if err != nil {
		ctx.RespondText("Error downloading file: %v", err)
		return
	}
	defer resp.Body.Close()

	buf := new(bytes.Buffer)

	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		ctx.RespondText("Error reading file body: %v", err)
		return
	}

	result, err := dest.Upload(buf, fileName)
	if err != nil {
		ctx.RespondText("Error uploading file: %v", err)
		return
	}

	c.services.Embed.New(ctx).
		SetDescription("Uploaded file: %s", fileName).
		AddField("URL", result.URL).
		AddField("Direct URL", result.DirectURL).
		Send()
}

// Paste utilities

func (c *UtilityCommands) pastebin(ctx *router.Ctx) {
	paste, err := api.LBCreatePaste(ctx.Arguments.Raw())
	if err != nil {
		ctx.RespondText("An error has occurred.")
		return
	}
	c.services.Embed.New(ctx).
		AddField("Paste URL", paste).
		Send()
}

func (c *UtilityCommands) pastefile(ctx *router.Ctx) {
	filePath := ctx.Arguments.Get(0).Raw()
	content, err := os.ReadFile(filePath)
	if err != nil {
		ctx.RespondText("An error has occurred.")
		return
	}

	paste, err := api.LBCreatePaste(string(content))
	if err != nil {
		ctx.RespondText("An error has occurred.")
		return
	}

	c.services.Embed.New(ctx).
		AddField("Paste URL", paste).
		Send()
}

func (c *UtilityCommands) tinyurl(ctx *router.Ctx) {
	sourceUrl := ctx.Arguments.Raw()
	resp, err := http.DefaultClient().Get(fmt.Sprintf("https://tinyurl.com/api-create.php?url=%s", sourceUrl))
	if err != nil {
		ctx.RespondText("An error has occurred.")
		return
	}
	c.services.Embed.New(ctx).
		AddField("URL", sourceUrl).
		AddField("TinyURL", string(resp.MustGetBody())).
		Send()
}

// Obfuscation utilities

func (c *UtilityCommands) obfuscateFunc(lang, ext string) func(ctx *router.Ctx) {
	return func(ctx *router.Ctx) {
		script := ctx.Arguments.Raw()

		attachments := ctx.Message.Attachments
		if script == "" && len(attachments) > 0 {
			resp, err := http.DefaultClient().Get(attachments[0].URL)
			if err != nil {
				ctx.RespondText("Error getting script: %v", err)
				return
			}

			script = string(resp.MustGetBody())
		}

		if script == "" {
			ctx.RespondText("Error: you must provide a script as an attachment or argument.")
			return
		}

		obf := api.LBObfuscate(lang, script)
		if strings.HasPrefix(obf, "An error") {
			ctx.RespondText("%s", obf)
			return
		}

		_, _ = ctx.Session.ChannelMessageSendComplex(ctx.Message.ChannelID, &discordgo.MessageSend{
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

func (c *UtilityCommands) reviewdbtoken(ctx *router.Ctx) {
	token, err := api.GetReviewDBToken(ctx.Session)
	if err != nil {
		c.services.Embed.New(ctx).
			SetDescription("Error getting token: %v", err).
			Send()
		return
	}

	c.services.Embed.New(ctx).
		AddField("Token", token).
		Send()
}

func (c *UtilityCommands) getreviews(ctx *router.Ctx) {
	user, err := ctx.Arguments.Get(0).User()
	if err != nil {
		c.services.Embed.New(ctx).
			SetDescription("Could not find user").
			Send()
		return
	}

	reviews, err := api.GetReviewDBReviews(user.ID)
	if err != nil {
		c.services.Embed.New(ctx).
			SetDescription("Error getting reviews: %v", err).
			Send()
		return
	}

	reviewString := ""

	for _, review := range reviews {
		reviewString += fmt.Sprintf("**[%s] %s:**\n%s\n\n", review.Timestamp.Format(time.RFC1123), review.Username, review.Comment)
	}

	c.services.Embed.New(ctx).
		AddField("Reviews", reviewString).
		Send()
}

func (c *UtilityCommands) addreview(ctx *router.Ctx) {
	user, err := ctx.Arguments.Get(0).User()
	if err != nil {
		c.services.Embed.New(ctx).
			SetDescription("Could not find user").
			Send()
		return
	}

	err = api.AddReviewDBReview(ctx.Session, user.ID, ctx.Arguments.RawFrom(1))
	if err != nil {
		c.services.Embed.New(ctx).
			SetDescription("Error adding review: %v", err).
			Send()
		return
	}

	c.services.Embed.New(ctx).
		SetDescription("Successfully added review to %s!", user.Username).
		Send()
}

// Tag utilities

func (c *UtilityCommands) getTags(ctx *router.Ctx) {
	cfg := c.services.Config.GetConfig()

	if len(cfg.Tags) == 0 {
		c.services.Embed.New(ctx).
			SetDescription("No tags found.").
			Send()
		return
	}

	var tagList []string

	for tag, value := range cfg.Tags {
		tagList = append(tagList, fmt.Sprintf("%s:\n```\n%s\n```\n", tag, value))
	}

	c.services.Embed.New(ctx).
		AddField("Tags", strings.Join(tagList, "\n")).
		Send()
}

func (c *UtilityCommands) addTag(ctx *router.Ctx) {
	cfg := c.services.Config.GetConfig()

	if cfg.Tags == nil {
		cfg.Tags = map[string]string{}
	}

	name := ctx.Arguments.Get(0).Raw()
	content := ctx.Arguments.RawFrom(1)

	cfg.Tags[name] = content

	if err := c.services.Config.SaveConfig(cfg); err != nil {
		ctx.RespondText("Error saving tags: %v", err)
		return
	}

	c.services.Embed.New(ctx).
		SetDescription("Successfully created tag.").
		Send()
}

func (c *UtilityCommands) removeTag(ctx *router.Ctx) {
	cfg := c.services.Config.GetConfig()

	name := ctx.Arguments.Get(0).Raw()

	delete(cfg.Tags, name)

	if err := c.services.Config.SaveConfig(cfg); err != nil {
		ctx.RespondText("Error saving tags: %v", err)
		return
	}

	c.services.Embed.New(ctx).
		SetDescription("Successfully deleted tag.").
		Send()
}

func (c *UtilityCommands) clearTags(ctx *router.Ctx) {
	cfg := c.services.Config.GetConfig()

	cfg.Tags = map[string]string{}

	if err := c.services.Config.SaveConfig(cfg); err != nil {
		ctx.RespondText("Error saving tags: %v", err)
		return
	}

	c.services.Embed.New(ctx).
		SetDescription("Successfully cleared tags.").
		Send()
}

// Creditcard commands

func (c *UtilityCommands) mastercard(ctx *router.Ctx) {
	creditCard, err := api.LBGenerateCC("mastercard", ctx.Arguments.Raw())
	if err != nil {
		ctx.RespondText("An error has occurred.")
		return
	}
	c.services.Embed.New(ctx).
		AddField("Card Type", "Mastercard").
		AddField("Card Number", creditCard.Number).
		AddField("CVV", creditCard.CVV).
		AddField("Expiry Date", creditCard.ExpiryDate).
		SetImage(creditCard.Visualisation).
		Send()
}

func (c *UtilityCommands) visa16(ctx *router.Ctx) {
	creditCard, err := api.LBGenerateCC("visa16", ctx.Arguments.Raw())
	if err != nil {
		ctx.RespondText("An error has occurred.")
		return
	}
	c.services.Embed.New(ctx).
		AddField("Card Type", "Visa").
		AddField("Card Number", creditCard.Number).
		AddField("CVV", creditCard.CVV).
		AddField("Expiry Date", creditCard.ExpiryDate).
		SetImage(creditCard.Visualisation).
		Send()
}

func (c *UtilityCommands) visa13(ctx *router.Ctx) {
	creditCard, err := api.LBGenerateCC("visa13", ctx.Arguments.Raw())
	if err != nil {
		ctx.RespondText("An error has occurred.")
		return
	}
	c.services.Embed.New(ctx).
		AddField("Card Type", "Visa").
		AddField("Card Number", creditCard.Number).
		AddField("CVV", creditCard.CVV).
		AddField("Expiry Date", creditCard.ExpiryDate).
		SetImage(creditCard.Visualisation).
		Send()
}

func (c *UtilityCommands) amex(ctx *router.Ctx) {
	creditCard, err := api.LBGenerateCC("amex", ctx.Arguments.Raw())
	if err != nil {
		ctx.RespondText("An error has occurred.")
		return
	}
	c.services.Embed.New(ctx).
		AddField("Card Type", "American Express").
		AddField("Card Number", creditCard.Number).
		AddField("CVV", creditCard.CVV).
		AddField("Expiry Date", creditCard.ExpiryDate).
		SetImage(creditCard.Visualisation).
		Send()
}

func (c *UtilityCommands) discover(ctx *router.Ctx) {
	creditCard, err := api.LBGenerateCC("discover", ctx.Arguments.Raw())
	if err != nil {
		ctx.RespondText("An error has occurred.")
		return
	}
	c.services.Embed.New(ctx).
		AddField("Card Type", "Discover").
		AddField("Card Number", creditCard.Number).
		AddField("CVV", creditCard.CVV).
		AddField("Expiry Date", creditCard.ExpiryDate).
		SetImage(creditCard.Visualisation).
		Send()
}

func (c *UtilityCommands) diners(ctx *router.Ctx) {
	creditCard, err := api.LBGenerateCC("diners", ctx.Arguments.Raw())
	if err != nil {
		ctx.RespondText("An error has occurred.")
		return
	}
	c.services.Embed.New(ctx).
		AddField("Card Type", "Diners").
		AddField("Card Number", creditCard.Number).
		AddField("CVV", creditCard.CVV).
		AddField("Expiry Date", creditCard.ExpiryDate).
		SetImage(creditCard.Visualisation).
		Send()
}

func (c *UtilityCommands) jcb(ctx *router.Ctx) {
	creditCard, err := api.LBGenerateCC("jcb", ctx.Arguments.Raw())
	if err != nil {
		ctx.RespondText("An error has occurred.")
		return
	}
	c.services.Embed.New(ctx).
		AddField("Card Type", "JCB").
		AddField("Card Number", creditCard.Number).
		AddField("CVV", creditCard.CVV).
		AddField("Expiry Date", creditCard.ExpiryDate).
		SetImage(creditCard.Visualisation).
		Send()
}

// Network utilities

func (c *UtilityCommands) ping(ctx *router.Ctx) {
	emb := c.services.Embed.New(ctx)
	time, err := network.GetPingLatency("https://lightning-bot.com")
	if err != nil {
		emb.AddField("Ping to LightningBot", "Error: "+err.Error())
	} else {
		emb.AddField("Ping to LightningBot", fmt.Sprintf("%dms", time.Milliseconds()))
	}
	time, err = network.GetPingLatency("https://google.com")
	if err != nil {
		emb.AddField("Ping to Google", "Error: "+err.Error())
	} else {
		emb.AddField("Ping to Google", fmt.Sprintf("%dms", time.Milliseconds()))
	}
	emb.
		AddField("Discord Socket Latency", fmt.Sprintf("%dms", ctx.Session.HeartbeatLatency().Milliseconds())).
		Send()
}

func (c *UtilityCommands) pingweb(ctx *router.Ctx) {
	host := ctx.Arguments.Get(0).Raw()
	emb := c.services.Embed.New(ctx)

	time, err := network.GetPingLatency(host)
	if err != nil {
		emb.AddField("Ping to "+host, "Error")
	} else {
		emb.AddField("Ping to "+host, fmt.Sprintf("%dms", time.Milliseconds()))
	}

	emb.Send()
}

func (c *UtilityCommands) ipinfo(ctx *router.Ctx) {
	ip := ctx.Arguments.Get(0).Raw()
	url := fmt.Sprintf("https://ipinfo.io/%s/json", url.QueryEscape(ip))

	resp, err := http.DefaultClient().Get(url)
	if err != nil {
		ctx.RespondText("An error has occurred.")
		return
	}

	ipinfo := resp.Json()

	if ipinfo.Get("status").MustInt(0) == 404 {
		ctx.RespondText("No results found.")
		return
	}

	location := fmt.Sprintf(
		"%s, %s, %s (postal code %s) (coordinates %s)",
		ipinfo.Get("city").MustString(), ipinfo.Get("region").MustString(),
		ipinfo.Get("country").MustString(), ipinfo.Get("postal").MustString(),
		ipinfo.Get("loc").MustString(),
	)

	c.services.Embed.New(ctx).
		AddField("IP", ip).
		AddField("Hostname", ipinfo.Get("hostname").MustString()).
		AddField("Timezone", ipinfo.Get("timezone").MustString()).
		AddField("Location", location).
		Send()
}

func (c *UtilityCommands) portscan(ctx *router.Ctx) {
	host := ctx.Arguments.Get(0).Raw()
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
	emb := c.services.Embed.New(ctx)
	for name, port := range ports {
		address := net.JoinHostPort(host, strconv.Itoa(port))
		conn, err := net.DialTimeout("tcp", address, 50*time.Millisecond)
		if err == nil {
			emb.AddField(fmt.Sprintf("%s (%d)", name, port), "Open")
			_ = conn.Close()
		}

	}
	emb.Send()
}

func (c *UtilityCommands) resolve(ctx *router.Ctx) {
	domain := ctx.Arguments.Get(0).Raw()

	resp, err := http.DefaultClient().Get(fmt.Sprintf("https://json.geoiplookup.io/%s", domain))
	if err != nil {
		ctx.RespondText("An error has occurred.")
		return
	}

	response := resp.Json()

	if response.Get("ip").MustString() == "" {
		c.services.Embed.New(ctx).SetDescription("Invalid domain.").Send()
		return
	}

	c.services.Embed.New(ctx).
		AddField("Domain", domain).
		AddField("IP", response.Get("ip").MustString()).
		Send()
}

func (c *UtilityCommands) reversedns(ctx *router.Ctx) {
	ip := ctx.Arguments.Get(0).Raw()
	url := fmt.Sprintf("https://api.LBdomaintools.com/v1/%s/host-domains/", ip)

	resp, err := http.DefaultClient().Get(url)
	if err != nil {
		ctx.RespondText("An error has occurred.")
		return
	}

	response := resp.Json()

	if response.Get("error").MustSlice() != nil {
		c.services.Embed.New(ctx).SetDescription("Invalid IP Address.").Send()
		return
	}

	domains := response.Get("response", "ip_addresses", "domain_names").MustStringSlice()

	c.services.Embed.New(ctx).
		AddField("IP Address", ip).
		AddField("Domains", strings.Join(domains, ", ")).
		Send()
}

func (c *UtilityCommands) reversemac(ctx *router.Ctx) {
	mac := ctx.Arguments.Get(0).Raw()
	resp, err := http.DefaultClient().Get(fmt.Sprintf("https://api.LBmacvendors.com/%s", mac))
	if err != nil {
		ctx.RespondText("An error has occurred.")
		return
	}
	vendor := string(resp.MustGetBody())
	if strings.Contains(vendor, "Not Found") {
		c.services.Embed.New(ctx).SetDescription("Invalid MAC Address.").Send()
		return
	}
	c.services.Embed.New(ctx).
		AddField("MAC Address", mac).
		AddField("Vendor", vendor).
		Send()
}

func (c *UtilityCommands) headers(ctx *router.Ctx) {
	url := ctx.Arguments.Get(0).Raw()
	resp, err := http.DefaultClient().Get(url)
	if err != nil {
		ctx.RespondText("An error has occurred.")
		return
	}

	var headers string
	for key, value := range resp.Header {
		headers += fmt.Sprintf("%s: %s\n", key, strings.Join(value, ", "))
	}

	c.services.Embed.New(ctx).
		AddField("URL", url).
		AddField("Headers", headers).
		Send()
}

func (c *UtilityCommands) _interface(ctx *router.Ctx) {
	c.services.Embed.New(ctx).
		AddField("Public IP", network.GetPublicIP()).
		AddField("Local IP", network.GetLocalIP()).
		AddField("MAC Address", network.GetMacAddress()).
		Send()
}

// Poketwo utilities

func (c *UtilityCommands) getPoketwoMessage(ctx *router.Ctx) (*discordgo.Message, error) {
	channel, err := ctx.Arguments.Get(0).Channel(router.ArgumentFlagFallbackDefault)
	if err != nil {
		ctx.RespondText("Error: invalid channel")
		return nil, errors.New("invalid channel")
	}

	messages, err := utils.Discord.GetChannelMessages(ctx.Session, channel.ID, 100)
	if err != nil {
		ctx.RespondText("Error getting messages: %v", err)
		return nil, err
	}

	for _, m := range messages {
		isPokeTwo := m.Author.ID == "716390085896962058" &&
			len(m.Embeds) == 1 &&
			strings.Contains(m.Embeds[0].Title, "has appeared!") &&
			m.Embeds[0].Image != nil

		if isPokeTwo {
			return m, nil
		}
	}

	ctx.RespondText("Error: no Poketwo message found in channel")
	return nil, errors.New("no Poketwo message found")
}

func (c *UtilityCommands) snipepokemon(ctx *router.Ctx) {
	message, err := c.getPoketwoMessage(ctx)
	if err != nil {
		return
	}

	emb := message.Embeds[0]

	// Download the Pokémon image
	resp, err := http.DefaultClient().Get(emb.Image.URL)
	if err != nil {
		ctx.RespondText("Error downloading image: %v", err)
		return
	}

	image := utils.Misc.ImageToBase64(resp.MustGetBody())

	pokemon, _, err := api.PokeTwo(image)
	if err != nil {
		logging.Debug("Snipers.CheckPokeTwoSniper: Error getting Pokémon: ", err)
		return
	}

	cfg := c.services.Config.GetConfig()

	// Catch the Pokémon
	label := pokemon.NameEnglish
	switch cfg.Snipers.PokeTwoSniper.CatchLanguage {
	case "Japanese":
		label = pokemon.NameJapanese
	case "JapaneseTransliterated":
		label = pokemon.NameJapaneseT
	case "German":
		label = pokemon.NameGerman
	case "French":
		label = pokemon.NameFrench
	}

	if label == "" {
		label = pokemon.NameEnglish
	}

	label = strings.ToLower(label)

	_, err = ctx.Session.ChannelMessageSend(message.ChannelID, "<@716390085896962058> c "+strings.ToLower(label))
	if err != nil {
		ctx.RespondText("Error catching Pokémon: %v", err)
		return
	}

	c.services.Embed.New(ctx).
		SetDescription("Successfully sniped Pokémon").
		AddField("Pokémon", label).
		Send()
}

func (c *UtilityCommands) detectpokemon(ctx *router.Ctx) {
	message, err := c.getPoketwoMessage(ctx)
	if err != nil {
		return
	}

	emb := message.Embeds[0]

	// Download the Pokémon image
	resp, err := http.DefaultClient().Get(emb.Image.URL)
	if err != nil {
		ctx.RespondText("Error downloading image: %v", err)
		return
	}

	image := utils.Misc.ImageToBase64(resp.MustGetBody())

	pokemon, _, err := api.PokeTwo(image)
	if err != nil {
		logging.Debug("Snipers.CheckPokeTwoSniper: Error getting Pokémon: ", err)
		return
	}

	names := fmt.Sprintf(
		"🇬🇧: %s\n🇯🇵: %s\n🇯🇵 (transliterated): %s\n🇩🇪: %s\n🇫🇷: %s",
		pokemon.NameEnglish,
		pokemon.NameJapanese,
		pokemon.NameJapaneseT,
		pokemon.NameGerman,
		pokemon.NameFrench,
	)

	info := fmt.Sprintf(
		"ID: %s\nRegion: %s\nType: %s\nDescription: %s\nEvent: %s\nHeight: %s\nWeight: %s\nHP: %s\nAttack: %s\nDefense: %s",
		pokemon.ID,
		pokemon.Region,
		pokemon.Type,
		pokemon.Description,
		pokemon.Event,
		pokemon.Height,
		pokemon.Weight,
		pokemon.HP,
		pokemon.Attack,
		pokemon.Defense,
	)

	c.services.Embed.New(ctx).
		AddField("Names", names).
		AddField("Info", info).
		Send()
}

// Reminder utilities

func (c *UtilityCommands) getReminders(ctx *router.Ctx) {
	cfg := c.services.Config.GetConfig()

	if len(cfg.Reminders) == 0 {
		ctx.RespondText("No reminders have been set!")
		return
	}

	var reminders []string
	for _, reminder := range cfg.Reminders {
		reminders = append(reminders, fmt.Sprintf("**ID: %d**: %s - %s", reminder.ID, reminder.Content, reminder.At.Format(time.RFC1123)))
	}

	c.services.Embed.New(ctx).
		AddField("Reminders", strings.Join(reminders, "\n")).
		Send()
}

func (c *UtilityCommands) addReminder(ctx *router.Ctx) {
	cfg := c.services.Config.GetConfig()

	name := ctx.Arguments.Get(0).Raw()

	for _, reminder := range cfg.Reminders {
		if reminder.Name != name {
			continue
		}

		ctx.RespondText("Reminder with that name already exists.")
		return
	}

	duration, err := ctx.Arguments.Get(2).Duration()
	if err != nil {
		ctx.RespondText("Invalid duration.")
		return
	}

	duration = fn.Clamp(duration, 5*time.Second, 24*time.Hour)

	message := ctx.Arguments.RawFrom(3)

	cfg.Reminders = append(cfg.Reminders, config.Reminder{
		Name:    ctx.Arguments.Get(1).Raw(),
		Content: message,
		At:      time.Now().Add(duration),
	})

	if err := c.services.Config.SaveConfig(cfg); err != nil {
		ctx.RespondText("Error saving config: %v", err)
		return
	}

	c.services.Embed.New(ctx).
		SetDescription("Successfully set reminder!").
		Send()
}

func (c *UtilityCommands) removeReminder(ctx *router.Ctx) {
	cfg := c.services.Config.GetConfig()

	name := ctx.Arguments.Get(0).Raw()

	cfg.Reminders = fn.Filter(cfg.Reminders, func(r config.Reminder) bool {
		return r.Name != name
	})

	if err := c.services.Config.SaveConfig(cfg); err != nil {
		ctx.RespondText("Error saving config: %v", err)
		return
	}

	c.services.Embed.New(ctx).
		SetDescription("Successfully deleted reminder!").
		Send()
}

func (c *UtilityCommands) clearReminders(ctx *router.Ctx) {
	cfg := c.services.Config.GetConfig()

	cfg.Reminders = []config.Reminder{}

	if err := c.services.Config.SaveConfig(cfg); err != nil {
		ctx.RespondText("Error saving config: %v", err)
		return
	}

	c.services.Embed.New(ctx).
		SetDescription("Successfully cleared reminders!").
		Send()
}

// Conversion utilities

func (c *UtilityCommands) checkIfUnitValid(ctx *router.Ctx, unit string, validUnits []string) bool {
	isValid := slices.Contains(validUnits, strings.ToLower(unit))
	if isValid {
		return true
	}

	c.services.Embed.New(ctx).
		AddField("Conversion error", "You have specified an invalid unit.").
		AddField("Valid units", strings.Join(validUnits, ", ")).
		Send()

	return false
}

func (c *UtilityCommands) convertvelocity(ctx *router.Ctx) {
	value, err := ctx.Arguments.Get(0).Float64()
	if err != nil {
		ctx.RespondText("Invalid value")
		return
	}

	sourceUnit := ctx.Arguments.Get(1).Raw()
	destinationUnit := ctx.Arguments.Get(2).Raw()

	if !c.checkIfUnitValid(ctx, sourceUnit, conversion.SpeedUnits) || !c.checkIfUnitValid(ctx, destinationUnit, conversion.SpeedUnits) {
		return
	}

	result, err := conversion.ConvertSpeed(sourceUnit, destinationUnit, value)
	if err != nil {
		ctx.RespondText("Error converting value: %v", err)
	}

	c.services.Embed.New(ctx).
		AddField("Value", strconv.FormatFloat(value, 'g', -1, 64)+" "+sourceUnit).
		AddField("Converted", strconv.FormatFloat(result, 'g', -1, 64)+" "+destinationUnit).
		Send()
}

func (c *UtilityCommands) convertarea(ctx *router.Ctx) {
	value, err := ctx.Arguments.Get(0).Float64()
	if err != nil {
		ctx.RespondText("Invalid value")
		return
	}

	sourceUnit := ctx.Arguments.Get(1).Raw()
	destinationUnit := ctx.Arguments.Get(2).Raw()

	if !c.checkIfUnitValid(ctx, sourceUnit, conversion.AreaUnits) || !c.checkIfUnitValid(ctx, destinationUnit, conversion.AreaUnits) {
		return
	}

	result, err := conversion.ConvertArea(sourceUnit, destinationUnit, value)
	if err != nil {
		ctx.RespondText("Error converting value: %v", err)
	}

	c.services.Embed.New(ctx).
		AddField("Value", strconv.FormatFloat(value, 'g', -1, 64)+" "+sourceUnit).
		AddField("Converted", strconv.FormatFloat(result, 'g', -1, 64)+" "+destinationUnit).
		Send()
}

func (c *UtilityCommands) convertlength(ctx *router.Ctx) {
	value, err := ctx.Arguments.Get(0).Float64()
	if err != nil {
		ctx.RespondText("Invalid value")
		return
	}

	sourceUnit := ctx.Arguments.Get(1).Raw()
	destinationUnit := ctx.Arguments.Get(2).Raw()

	if !c.checkIfUnitValid(ctx, sourceUnit, conversion.LengthUnits) || !c.checkIfUnitValid(ctx, destinationUnit, conversion.LengthUnits) {
		return
	}

	result, err := conversion.ConvertLength(sourceUnit, destinationUnit, value)
	if err != nil {
		ctx.RespondText("Error converting value: %v", err)
	}

	c.services.Embed.New(ctx).
		AddField("Value", strconv.FormatFloat(value, 'g', -1, 64)+" "+sourceUnit).
		AddField("Converted", strconv.FormatFloat(result, 'g', -1, 64)+" "+destinationUnit).
		Send()
}

func (c *UtilityCommands) converttemperature(ctx *router.Ctx) {
	value, err := ctx.Arguments.Get(0).Float64()
	if err != nil {
		ctx.RespondText("Invalid value")
		return
	}

	sourceUnit := ctx.Arguments.Get(1).Raw()
	destinationUnit := ctx.Arguments.Get(2).Raw()

	if !c.checkIfUnitValid(ctx, sourceUnit, conversion.TemperatureUnits) || !c.checkIfUnitValid(ctx, destinationUnit, conversion.TemperatureUnits) {
		return
	}

	result, err := conversion.ConvertTemperature(sourceUnit, destinationUnit, value)
	if err != nil {
		ctx.RespondText("Error converting value: %v", err)
	}

	c.services.Embed.New(ctx).
		AddField("Value", strconv.FormatFloat(value, 'g', -1, 64)+" "+sourceUnit).
		AddField("Converted", strconv.FormatFloat(result, 'g', -1, 64)+" "+destinationUnit).
		Send()
}

func (c *UtilityCommands) converttime(ctx *router.Ctx) {
	value, err := ctx.Arguments.Get(0).Float64()
	if err != nil {
		ctx.RespondText("Invalid value")
		return
	}

	sourceUnit := ctx.Arguments.Get(1).Raw()
	destinationUnit := ctx.Arguments.Get(2).Raw()

	if !c.checkIfUnitValid(ctx, sourceUnit, conversion.TimeUnits) || !c.checkIfUnitValid(ctx, destinationUnit, conversion.TimeUnits) {
		return
	}

	result, err := conversion.ConvertTime(sourceUnit, destinationUnit, value)
	if err != nil {
		ctx.RespondText("Error converting value: %v", err)
	}

	c.services.Embed.New(ctx).
		AddField("Value", strconv.FormatFloat(value, 'g', -1, 64)+" "+sourceUnit).
		AddField("Converted", strconv.FormatFloat(result, 'g', -1, 64)+" "+destinationUnit).
		Send()
}

func (c *UtilityCommands) convertvolume(ctx *router.Ctx) {
	value, err := ctx.Arguments.Get(0).Float64()
	if err != nil {
		ctx.RespondText("Invalid value")
		return
	}

	sourceUnit := ctx.Arguments.Get(1).Raw()
	destinationUnit := ctx.Arguments.Get(2).Raw()

	if !c.checkIfUnitValid(ctx, sourceUnit, conversion.VolumeUnits) || !c.checkIfUnitValid(ctx, destinationUnit, conversion.VolumeUnits) {
		return
	}

	result, err := conversion.ConvertVolume(sourceUnit, destinationUnit, value)
	if err != nil {
		ctx.RespondText("Error converting value: %v", err)
	}

	c.services.Embed.New(ctx).
		AddField("Value", strconv.FormatFloat(value, 'g', -1, 64)+" "+sourceUnit).
		AddField("Converted", strconv.FormatFloat(result, 'g', -1, 64)+" "+destinationUnit).
		Send()
}

func (c *UtilityCommands) convertmass(ctx *router.Ctx) {
	value, err := ctx.Arguments.Get(0).Float64()
	if err != nil {
		ctx.RespondText("Invalid value")
		return
	}

	sourceUnit := ctx.Arguments.Get(1).Raw()
	destinationUnit := ctx.Arguments.Get(2).Raw()

	if !c.checkIfUnitValid(ctx, sourceUnit, conversion.MassUnits) || !c.checkIfUnitValid(ctx, destinationUnit, conversion.MassUnits) {
		return
	}

	result, err := conversion.ConvertMass(sourceUnit, destinationUnit, value)
	if err != nil {
		ctx.RespondText("Error converting value: %v", err)
	}

	c.services.Embed.New(ctx).
		AddField("Value", strconv.FormatFloat(value, 'g', -1, 64)+" "+sourceUnit).
		AddField("Converted", strconv.FormatFloat(result, 'g', -1, 64)+" "+destinationUnit).
		Send()
}

// Misc utilities

func (c *UtilityCommands) tts(ctx *router.Ctx) {
	reader, err := api.TTS(ctx.Arguments.RawFrom(1), ctx.Arguments.Get(0).Raw())
	if err != nil {
		ctx.RespondText("An error has occurred: %v", err)
		return
	}

	_, _ = ctx.Session.ChannelFileSend(ctx.Message.ChannelID, "tts.mp3", reader)
}

func (c *UtilityCommands) translate(ctx *router.Ctx) {
	translateResult := api.Translate(ctx.Arguments.RawFrom(1), ctx.Arguments.Get(0).Raw())
	if translateResult.Error != nil {
		ctx.RespondText("An error has occurred: %v", translateResult.Error)
		return
	}
	c.services.Embed.New(ctx).
		AddField("Source Language", utils.String.Capitalize(constants.LANGUAGES[translateResult.SourceLanguage])).
		AddField("Destination Language", utils.String.Capitalize(constants.LANGUAGES[translateResult.DestinationLanguage])).
		AddField("Original Text", translateResult.Original).
		AddField("Translated Text", translateResult.Result).
		Send()
}

func (c *UtilityCommands) detectlang(ctx *router.Ctx) {
	translateResult := api.Translate(ctx.Arguments.Raw(), "en")
	if translateResult.Error != nil {
		ctx.RespondText("An error has occurred: %v", translateResult.Error)
		return
	}
	c.services.Embed.New(ctx).
		AddField("Text", translateResult.Original).
		AddField("Detected Language", utils.String.Capitalize(constants.LANGUAGES[translateResult.SourceLanguage])).
		AddField("Confidence", fmt.Sprintf("%.2f", translateResult.Confidence)).
		Send()
}

func (c *UtilityCommands) english(ctx *router.Ctx) {
	translateResult := api.Translate(ctx.Arguments.Raw(), "en")
	if translateResult.Error != nil {
		ctx.RespondText("An error has occurred: %v", translateResult.Error)
		return
	}
	c.services.Embed.New(ctx).
		AddField("Source Language", utils.String.Capitalize(constants.LANGUAGES[translateResult.SourceLanguage])).
		AddField("Original Text", translateResult.Original).
		AddField("Translated Text", translateResult.Result).
		Send()
}

func (c *UtilityCommands) define(ctx *router.Ctx) {
	url := fmt.Sprintf(
		"https://api.LBdictionaryapi.LBdev/api/v2/entries/en/%s",
		ctx.Arguments.Escaped(),
	)

	resp, err := http.DefaultClient().Get(url)
	if err != nil {
		ctx.RespondText("An error has occurred.")
		return
	}

	list := resp.Json()

	if list.Len() == 0 {
		ctx.RespondText("No results found.")
		return
	}

	result := list.Get(0)

	meanings := fn.Map(result.Get("meanings").MustSlice(), func(obj any) string {
		meaning := obj.(map[string]any)
		definitions := meaning["definitions"].([]any)

		if len(definitions) == 0 {
			return ""
		}

		definition := definitions[0].(map[string]any)["definition"].(string)

		return fmt.Sprintf("**%s**: %s", meaning["partOfSpeech"].(string), definition)
	})

	meanings = fn.Filter(meanings, func(obj string) bool {
		return obj != ""
	})

	c.services.Embed.New(ctx).
		AddField("Word", ctx.Arguments.Raw()).
		AddField("Meanings", strings.Join(meanings, "\n")).
		Send()
}

func (c *UtilityCommands) urbandictionary(ctx *router.Ctx) {
	url := fmt.Sprintf(
		"https://api.LBurbandictionary.com/v0/define?term=%s",
		ctx.Arguments.Escaped(),
	)

	resp, err := http.DefaultClient().Get(url)
	if err != nil {
		ctx.RespondText("An error has occurred.")
		return
	}

	list := resp.Json().Get("list")

	if list.Len() == 0 {
		ctx.RespondText("No results found.")
		return
	}

	result := list.Get(0)

	c.services.Embed.New(ctx).
		AddField("Word", ctx.Arguments.Raw()).
		AddField("Definition", result.Get("definition").MustString()).
		AddField("Permalink", result.Get("permalink").MustString()).
		AddField("Thumbs Up", fmt.Sprint(result.Get("thumbs_up").MustFloat64())).
		AddField("Thumbs Down", fmt.Sprint(result.Get("thumbs_down").MustFloat64())).
		AddField("Author", result.Get("author").MustString()).
		Send()
}

func (c *UtilityCommands) wikipedia(ctx *router.Ctx) {
	query := ctx.Arguments.Raw()
	article, err := api.WikipediaSearch(query)
	if err != nil {
		ctx.RespondText("An error has occurred.")
		return
	}

	if len(article) > 2000 {
		article = article[:2000]
	}

	c.services.Embed.New(ctx).
		AddField("Query", query).
		AddField("Wikipedia Summary", article).
		Send()
}

func (c *UtilityCommands) weather(ctx *router.Ctx) {
	weather, err := api.GetWeather(ctx.Arguments.Raw())
	if err != nil {
		ctx.RespondText("An error has occurred: %v", err)
		return
	}

	c.services.Embed.New(ctx).
		SetDescription("Weather in %s, %s", weather.Get("name").MustString(), weather.Get("sys", "country").MustString()).
		AddField("Description", utils.String.Capitalize(weather.Get("weather").Get(0).Get("description").MustString())).
		AddField("Temperature", fmt.Sprintf("%.2f °C", math.Round(weather.Get("main", "temp").MustFloat64()-273.15))).
		AddField("Feels Like", fmt.Sprintf("%.2f °C", math.Round(weather.Get("main", "feels_like").MustFloat64()-273.15))).
		AddField("Min. Temperature", fmt.Sprintf("%.2f °C", math.Round(weather.Get("main", "temp_min").MustFloat64()-273.15))).
		AddField("Max. Temperature", fmt.Sprintf("%.2f °C", math.Round(weather.Get("main", "temp_max").MustFloat64()-273.15))).
		AddField("Pressure", fmt.Sprintf("%.2f hPa", weather.Get("main", "pressure").MustFloat64())).
		AddField("Humidity", fmt.Sprintf("%.2f %%", weather.Get("main", "humidity").MustFloat64())).
		AddField("Visibility", fmt.Sprintf("%.2f m", weather.Get("visibility").MustFloat64())).
		AddField("Wind Speed", fmt.Sprintf("%.2f Km/h", math.Round(weather.Get("wind", "speed").MustFloat64()*3.6))).
		AddField("Wind Gust", fmt.Sprintf("%.2f Km/h", math.Round(weather.Get("wind", "gust").MustFloat64()*3.6))).
		AddField("Wind Direction", fmt.Sprintf("%.2f°", weather.Get("wind", "deg").MustFloat64())).
		Send()
}

func (c *UtilityCommands) lvbypass(ctx *router.Ctx) {
	resp, err := http.DefaultClient().Get("https://bypass.pm/bypass2?url=" + ctx.Arguments.Raw())
	if err != nil {
		ctx.RespondText("An error has occurred: %v", err)
		return
	}

	destination := resp.Json().Get("destination").MustString()

	c.services.Embed.New(ctx).
		AddField("URL", ctx.Arguments.Raw()).
		AddField("Bypassed URL", utils.Misc.Or(destination, "An error has occurred")).
		Send()
}

func (c *UtilityCommands) flightinfo(ctx *router.Ctx) {
	uri := fmt.Sprintf(
		"https://api.LBflightradar24.com/common/v1/flight/list.json?query=%s&fetchBy=flight&page=1",
		url.QueryEscape(ctx.Arguments.Get(0).Raw()),
	)

	resp, err := http.DefaultClient().Get(uri)
	if err != nil {
		ctx.RespondText("An error has occurred.")
		return
	}

	json := resp.Json()

	if errors := json.Get("errors").Interface(); errors != nil {
		ctx.RespondText("That flight does not exist or an error has occurred.")
		return
	}

	data := json.Get("result", "response", "data").Get(0)

	departureTime := data.Get("time", "scheduled", "departure").MustInt64()
	arrivalTime := data.Get("time", "scheduled", "arrival").MustInt64()

	c.services.Embed.New(ctx).
		AddField("Flight Number", data.Get("identification", "number", "default").MustString()).
		AddField("Airline", data.Get("airline", "name").MustString()).
		AddField("Departure Airport", data.Get("airport", "origin", "name").MustString()).
		AddField("Departure Time", time.Unix(departureTime, 0).Format("2006-01-02 15:04:05")).
		AddField("Arrival Airport", data.Get("airport", "destination", "name").MustString()).
		AddField("Arrival Time", time.Unix(arrivalTime, 0).Format("2006-01-02 15:04:05")).
		AddField("Status", data.Get("status", "text").MustString()).
		Send()
}

func (c *UtilityCommands) passgen(ctx *router.Ctx) {
	charAmount, err := ctx.Arguments.Get(0).Int()
	if err != nil {
		ctx.RespondText("Argument <character amount> must be a number.")
		return
	}

	charAmount = fn.Clamp(charAmount, 1, 1000)

	c.services.Embed.New(ctx).
		AddField("Password", random.String(charAmount, "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789#$%&")).
		Send()
}

func (c *UtilityCommands) tokengen(ctx *router.Ctx) {
	resp, err := http.DefaultClient().Get("https://some-random-api.LBcom/bottoken")
	if err != nil {
		ctx.RespondText("Error: %v", err)
		return
	}

	token := resp.Json().Get("token").MustString()

	c.services.Embed.New(ctx).
		AddField("Token", token).
		Send()
}

func (c *UtilityCommands) revimage(ctx *router.Ctx) {
	imageUrl := ctx.Arguments.Get(0).Raw()
	c.services.Embed.New(ctx).
		AddField("Reverse Image Search URL", fmt.Sprintf("https://www.google.com/searchbyimage?image_url=%s", imageUrl)).
		Send()
}

func (c *UtilityCommands) revavatar(ctx *router.Ctx) {
	user, err := ctx.Arguments.Get(0).User()
	if err != nil {
		ctx.RespondText("Could not find user")
		return
	}
	c.services.Embed.New(ctx).
		AddField("Search URL", fmt.Sprintf("https://www.google.com/searchbyimage?&image_url=%s", url.QueryEscape(user.AvatarURL("")))).
		Send()
}

func (c *UtilityCommands) timezones(ctx *router.Ctx) {
	continent := ctx.Arguments.Get(0).Raw()
	resp, err := http.DefaultClient().Get("http://worldtimeapi.LBorg/api/timezone/" + url.PathEscape(continent))
	if err != nil {
		ctx.RespondText("Error getting timezones: %v", err)
		return
	}

	json := resp.Json()

	if err := json.Get("error").MustString(); err != "" {
		ctx.RespondText("Error getting timezones: %v", err)
		return
	}

	c.services.Embed.New(ctx).
		AddField("Continent", utils.String.Capitalize(continent)).
		AddField("Timezones", strings.Join(json.MustStringSlice(), ", ")).
		Send()
}

func (c *UtilityCommands) timezoneinfo(ctx *router.Ctx) {
	resp, err := http.DefaultClient().Get("http://worldtimeapi.LBorg/api/timezone/" + ctx.Arguments.Get(0).Raw())
	if err != nil {
		ctx.RespondText("Error getting timezones: %v", err)
		return
	}

	json := resp.Json()

	if err := json.Get("error").MustString(); err != "" {
		ctx.RespondText("Error getting timezones: %v", err)
		return
	}

	c.services.Embed.New(ctx).
		AddField("Timezone", json.Get("timezone").MustString()).
		AddField("Abbreviation", json.Get("abbreviation").MustString()).
		AddField("UTC Offset", json.Get("utc_offset").MustString()).
		AddField("Current Time", json.Get("datetime").MustString()).
		Send()
}

func (c *UtilityCommands) calculate(ctx *router.Ctx) {
	expression := ctx.Arguments.Get(0).Raw()

	tv, err := types.Eval(token.NewFileSet(), nil, token.NoPos, expression)
	if err != nil {
		ctx.RespondText("Error evaluating expression: %v", err)
		return
	}

	if tv.Value == nil {
		ctx.RespondText("Error evaluating expression: nil value")
		return
	}

	c.services.Embed.New(ctx).
		AddField("Expression", ctx.Arguments.Get(0).Raw()).
		AddField("Solution", tv.Value.String()).
		Send()
}

func (c *UtilityCommands) getCreds() api.LBCreds {
	cfg := c.services.Config.GetConfig()
	return api.LBCreds{
		Username: cfg.Account.LightningBot.Username,
		Password: cfg.Account.LightningBot.Password,
	}
}

func (c *UtilityCommands) googlesearch(ctx *router.Ctx) {
	results, err := api.LBSearch(c.getCreds(), ctx.Arguments.Raw())
	if err != nil {
		ctx.RespondText("An error has occurred.")
		return
	}
	emb := c.services.Embed.New(ctx)
	for result := range results.Iterator() {
		emb.AddField(result.Get("title").MustString(), result.Get("url").MustString())
	}
	emb.Send()
}

func (c *UtilityCommands) youtubesearch(ctx *router.Ctx) {
	results, err := api.LBSearch(c.getCreds(), ctx.Arguments.Raw()+" site:youtube.com")
	if err != nil {
		ctx.RespondText("An error has occurred.")
		return
	}
	emb := c.services.Embed.New(ctx)
	for result := range results.Iterator() {
		resultUrl := result.Get("url").MustString()
		if strings.Contains(resultUrl, "playlist") {
			continue
		}
		emb.AddField(result.Get("title").MustString(), result.Get("url").MustString())
	}
	emb.Send()
}

func (c *UtilityCommands) amazon(ctx *router.Ctx) {
	query := ctx.Arguments.Raw()
	c.services.Embed.New(ctx).
		AddField("Amazon Search", fmt.Sprintf("https://www.amazon.com/s?k=%s", url.QueryEscape(query))).
		Send()
}

func (c *UtilityCommands) steam(ctx *router.Ctx) {
	query := ctx.Arguments.Raw()
	c.services.Embed.New(ctx).
		AddField("Steam Search", fmt.Sprintf("https://store.steampowered.com/search/?term=%s", url.QueryEscape(query))).
		Send()
}

func (c *UtilityCommands) github(ctx *router.Ctx) {
	query := ctx.Arguments.Raw()
	c.services.Embed.New(ctx).
		AddField("GitHub Search", fmt.Sprintf("https://github.com/search?q=%s", url.QueryEscape(query))).
		Send()
}

func (c *UtilityCommands) stackoverflow(ctx *router.Ctx) {
	query := ctx.Arguments.Raw()
	c.services.Embed.New(ctx).
		AddField("Stack Overflow Search", fmt.Sprintf("https://stackoverflow.com/search?q=%s", url.QueryEscape(query))).
		Send()
}

func (c *UtilityCommands) searchcountry(ctx *router.Ctx) {
	resp, err := http.DefaultClient().Get(fmt.Sprintf("https://restcountries.com/v3.1/name/%s?fullText=true", ctx.Arguments.Raw()))
	if err != nil {
		ctx.RespondText("An error has occurred.")
		return
	}

	countryList := resp.Json()

	country := countryList.Get(0)

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

	c.services.Embed.New(ctx).
		SetTitle("%s - %s", country.Get("name", "common").MustString(), country.Get("cca2").MustString()).
		AddField("Name", country.Get("name", "common").MustString()).
		AddField("Top-Level Domain", strings.Join(country.Get("tld").MustStringSlice(), ", ")).
		AddField("Country Code", country.Get("cca2").MustString()).
		AddField("Capital", strings.Join(country.Get("capital").MustStringSlice(), ", ")).
		AddField("Region", country.Get("region").MustString()).
		AddField("Subregion", country.Get("subregion").MustString()).
		AddField("Population", fmt.Sprint(country.Get("population").MustInt())).
		AddField("Latitude/Longitude", fmt.Sprintf("%f, %f", country.Get("latlng").Get(0).MustFloat64(), country.Get("latlng").Get(1).MustFloat64())).
		AddField("Timezone", strings.Join(country.Get("timezones").MustStringSlice(), ", ")).
		AddField("Currency", strings.Join(currencies, ", ")).
		AddField("Language", strings.Join(languages, ", ")).
		Send()
}

func (c *UtilityCommands) randomcountry(ctx *router.Ctx) {
	resp, err := http.DefaultClient().Get("https://restcountries.com/v3.1/all")
	if err != nil {
		ctx.RespondText("An error has occurred.")
		return
	}

	countryList := resp.Json()

	country := countryList.Get(random.Int(0, countryList.Len()))

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

	c.services.Embed.New(ctx).
		SetTitle("%s - %s", country.Get("name", "common").MustString(), country.Get("cca2").MustString()).
		AddField("Name", country.Get("name", "common").MustString()).
		AddField("Top-Level Domain", strings.Join(country.Get("tld").MustStringSlice(), ", ")).
		AddField("Country Code", country.Get("cca2").MustString()).
		AddField("Capital", strings.Join(country.Get("capital").MustStringSlice(), ", ")).
		AddField("Region", country.Get("region").MustString()).
		AddField("Subregion", country.Get("subregion").MustString()).
		AddField("Population", fmt.Sprint(country.Get("population").MustInt())).
		AddField("Latitude/Longitude", fmt.Sprintf("%f, %f", country.Get("latlng").Get(0).MustFloat64(), country.Get("latlng").Get(1).MustFloat64())).
		AddField("Timezone", strings.Join(country.Get("timezones").MustStringSlice(), ", ")).
		AddField("Currency", strings.Join(currencies, ", ")).
		AddField("Language", strings.Join(languages, ", ")).
		Send()
}

func (c *UtilityCommands) loremipsum(ctx *router.Ctx) {
	paragrahps, err := ctx.Arguments.Get(0).Int()
	if err != nil {
		paragrahps = 1
	}

	resp, err := http.DefaultClient().Get(fmt.Sprintf("https://loripsum.net/api/%d/short/plaintext", paragrahps))
	if err != nil {
		ctx.RespondText("An error has occurred.")
		return
	}

	c.services.Embed.New(ctx).
		SetDescription("%s", string(resp.MustGetBody())).
		Send()
}

func (c *UtilityCommands) baconipsum(ctx *router.Ctx) {
	paragrahps, err := ctx.Arguments.Get(0).Int()
	if err != nil {
		paragrahps = 1
	}

	resp, err := http.DefaultClient().Get(fmt.Sprintf("https://baconipsum.com/api/?type=all-meat&paras=%d&format=text", paragrahps))
	if err != nil {
		ctx.RespondText("An error has occurred.")
		return
	}

	c.services.Embed.New(ctx).
		SetDescription("%s", string(resp.MustGetBody())).
		Send()
}

func (c *UtilityCommands) previewcolor(ctx *router.Ctx) {
	color := ctx.Arguments.Get(0).Raw()
	color = strings.ReplaceAll(color, "#", "")

	var r, g, b int

	// Parse hex
	if len(color) == 6 {
		r, g, b, _ = utils.Misc.HexToRGB(color)
	} else {
		// Parse RGB
		rgb := strings.Split(color, ",")
		if len(rgb) != 3 {
			ctx.RespondText("Invalid RGB format. Use R,G,B (e.g. 255,255,255)")
			return
		}

		r, _ = strconv.Atoi(rgb[0])
		g, _ = strconv.Atoi(rgb[1])
		b, _ = strconv.Atoi(rgb[2])
	}

	hex := fmt.Sprintf("%02x%02x%02x", r, g, b)
	name := "Unknown"

	resp, err := http.DefaultClient().Get("https://www.thecolorapi.LBcom/id?hex=" + hex)
	if err == nil {
		name = resp.Json().Get("name", "value").MustString()
	}

	c.services.Embed.New(ctx).
		AddField("Color", name).
		AddField("RGB", fmt.Sprintf("%d, %d, %d", r, g, b)).
		AddField("Hex", "#"+hex).
		Send()

	preview, err := utils.Misc.PreviewColorImage(500, 100, r, g, b)
	if err == nil {
		_, _ = ctx.Session.ChannelFileSend(ctx.Message.ChannelID, "color.png", preview)
	}
}

func (c *UtilityCommands) generatekey(ctx *router.Ctx) {
	format := ctx.Arguments.Get(0).Raw()
	amount, _ := ctx.Arguments.Get(1).Int()
	amount = fn.Clamp(amount, 1, 100)

	var keys []string

	for range amount {
		var key string

		for _, char := range format {
			switch char {
			case 'x':
				key += random.String(1, constants.ALPHABET)
			case 'X':
				key += strings.ToUpper(random.String(1, constants.ALPHABET))
			case '9':
				key += strconv.Itoa(random.Int(0, 10))
			default:
				key += string(char)
			}
		}

		keys = append(keys, key)
	}

	c.services.Embed.New(ctx).
		AddField("License Keys", strings.Join(keys, "\n")).
		Send()
}

func (c *UtilityCommands) timestampCommand(flags string) func(ctx *router.Ctx) {
	return func(ctx *router.Ctx) {
		timestamp, err := api.LBGPT(c.getCreds(), utils.Misc.CreateAISystemTimestamp(), ctx.Arguments.Raw())
		if err != nil {
			ctx.RespondText("An error has occurred while parsing the timestamp.")
			return
		}

		if timestamp == "UNABLE" {
			ctx.RespondText("Please provide a valid date/time.")
			return
		}

		time, err := time.Parse(time.RFC3339, timestamp)
		if err != nil {
			ctx.RespondText("An error has occurred while parsing the timestamp.")
			return
		}

		ctx.RespondText("<t:%d%s>", time.Unix(), flags)
	}
}
