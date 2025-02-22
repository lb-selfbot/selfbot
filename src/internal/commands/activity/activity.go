package activity

import (
	"encoding/json"
	"fmt"
	"net/url"
	"reflect"
	"slices"
	"strings"
	"time"

	"github.com/LightningDev1/LB-Selfbot/internal/api"
	"github.com/LightningDev1/LB-Selfbot/internal/config"
	"github.com/LightningDev1/LB-Selfbot/internal/http"
	"github.com/LightningDev1/LB-Selfbot/internal/router"
	"github.com/LightningDev1/LB-Selfbot/internal/services"
	"github.com/LightningDev1/LB-Selfbot/internal/utils"
	"github.com/LightningDev1/LB-Selfbot/selfbot/sessions"
	"github.com/LightningDev1/discordgo"
	"github.com/LightningDev1/discordgo/protos"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type ActivityCommands struct{ services *services.Services }

func NewActivityCommands(s *services.Services) *ActivityCommands {
	return &ActivityCommands{services: s}
}
func (c *ActivityCommands) Register() {
	cmdRouter := c.services.Discord.GetRouter()

	cmdRouter.RegisterCategory(&router.Category{
		Name:        "Activity",
		Description: "Customize your Discord activity",
	}, []*router.Command{
		{
			Name:        "status",
			Description: "Set your Discord status",
			Arguments:   "status <online/dnd/idle/invisible>",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.status,
		},
		{
			Name:        "customstatus",
			Description: "Change your custom status",
			Arguments:   "customstatus [text]",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.customstatus,
		},
		{
			Name:        "clearactivity",
			Description: "Clear activity",
			Arguments:   "clearactivity",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.clearactivity,
		},
		{
			Name:        "copyactivity",
			Description: "Copy a user's activity",
			Arguments:   "copyactivity <user> [save to config: on/off]",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.copyactivity,
		},
		{
			Name:        "getactivityjson",
			Description: "Get a users activity in JSON format",
			Arguments:   "getactivityjson <user>",
			Handler:     c.getactivityjson,
		},
		{
			Name:        "setactivityjson",
			Description: "Set your activity using raw JSON",
			Arguments:   "setactivityjson <save to config: on/off> <json>",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.setactivityjson,
		},
		{
			Name:        "cactivity",
			Description: "Get, set or toggle Custom Activity",
			Arguments:   "cactivity <get/set/toggle> [on/off]",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.cactivity,
		},
		{
			Name:        "capauseforother",
			Description: "Get, set or toggle Pause CA for other Activity",
			Arguments:   "capauseforother <get/set/toggle> [on/off]",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.capauseforother,
		},
		{
			Name:        "caconnectiontype",
			Description: "Get or set the Custom Activity connection type",
			Arguments:   "caconnectiontype <get/set/values> [value]",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.caconnectiontype,
		},
		{
			Name:        "catype",
			Description: "Get or set the Custom Activity type",
			Arguments:   "catype <get/set/values> [value]",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.catype,
		},
		{
			Name:        "caplatform",
			Description: "Get or set the Custom Activity platform",
			Arguments:   "caplatform <get/set/values> [value]",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.caplatform,
		},
		{
			Name:        "caclientid",
			Description: "Get or set the Custom Activity Client ID",
			Arguments:   "caclientid <get/set> [value]",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.caclientid,
		},
		{
			Name:        "caname",
			Description: "Get or set the Custom Activity Name",
			Arguments:   "caname <get/set> [value]",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.caname,
		},
		{
			Name:        "castate",
			Description: "Get or set the Custom Activity State",
			Arguments:   "castate <get/set> [value]",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.castate,
		},
		{
			Name:        "cadetails",
			Description: "Get or set the Custom Activity Details",
			Arguments:   "cadetails <get/set> [value]",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.cadetails,
		},
		{
			Name:        "cainvitable",
			Description: "Get, set or toggle if RPC is invitable",
			Arguments:   "cainvitable <get/set/toggle> [on/off]",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.cainvitable,
		},
		{
			Name:        "catime",
			Description: "Get, set or toggle if Custom Activity shows time",
			Arguments:   "catime <get/set/toggle> [on/off]",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.catime,
		},
		{
			Name:        "castreamingurl",
			Description: "Get or set the URL for the Streaming activity",
			Arguments:   "castreamingurl <get/set> [value]",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.castreamingurl,
		},
		{
			Name:        "cacreateimagekey",
			Description: "Create an image key for Custom Activity",
			Arguments:   "cacreateimagekey <image url>",
			Handler:     c.cacreateimagekey,
		},
		{
			Name:        "cacreatespotifykey",
			Description: "Create an Spotify image key for Custom Activity",
			Arguments:   "cacreatespotifykey <song name>",
			Handler:     c.cacreatespotifykey,
		},
		{
			Name:        "calargeimage",
			Description: "Get or set the Custom Activity large image key",
			Arguments:   "calargeimage <get/set> [value]",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.calargeimage,
		},
		{
			Name:        "calargetext",
			Description: "Get or set the Custom Activity large image text",
			Arguments:   "calargetext <get/set> [value]",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.calargetext,
		},
		{
			Name:        "casmallimage",
			Description: "Get or set the Custom Activity small image key",
			Arguments:   "casmallimage <get/set> [value]",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.casmallimage,
		},
		{
			Name:        "casmalltext",
			Description: "Get or set the Custom Activity small image text",
			Arguments:   "casmalltext <get/set> [value]",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.casmalltext,
		},
		{
			Name:        "cabutton1",
			Description: "Get or set the label and URL of the first button",
			Arguments:   "cabutton1 <get/set> [\"label\"] [url]",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.cabutton1,
		},
		{
			Name:        "cabutton2",
			Description: "Get or set the label and URL of the second button",
			Arguments:   "cabutton2 <get/set> [\"label\"] [url]",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.cabutton2,
		},
		{
			Name:        "capresets",
			Description: "Manage Custom Activity presets",
			Arguments:   "capresets <get/load/next/save/rm/list/clear> [name]",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.capresets,
		},
		{
			Name:        "activitypresets",
			Description: "Set or view built-in Custom Activity presets",
			Arguments:   "activitypresets <set/values> [name] [text]",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.activitypresets,
		},
	})
}

func (a *ActivityCommands) setStatus(session *discordgo.Session, status string) error {
	sessions.LastStatusChange = time.Now()
	sessions.Status = status

	session.State.UserSettings.Status.Status = wrapperspb.String(status)
	err := session.UserSettingsUpdate(session.State.UserSettings, true)
	if err != nil {
		return err
	}

	return session.UpdateStatusComplex(discordgo.UpdateStatusData{
		AFK:    true,
		Status: status,
	})
}

func (a *ActivityCommands) setCustomStatus(session *discordgo.Session, status string) error {
	var customStatus *protos.PreloadedUserSettings_CustomStatus
	var activity *discordgo.Activity

	if status != "" {
		customStatus = &protos.PreloadedUserSettings_CustomStatus{
			Text: status,
		}

		activity = &discordgo.Activity{
			Type:  discordgo.ActivityTypeCustom,
			Name:  "Custom Status",
			State: status,
		}
	}

	session.State.UserSettings.Status.CustomStatus = customStatus
	err := session.UserSettingsUpdate(session.State.UserSettings, true)
	if err != nil {
		return err
	}

	return session.UpdateStatusComplex(discordgo.UpdateStatusData{
		Activities: []*discordgo.Activity{activity},
		AFK:        true,
		Status:     sessions.Status,
	})
}

func (a *ActivityCommands) status(ctx *router.Ctx) {
	status := strings.ToLower(ctx.Arguments.Get(0).Raw())

	switch status {
	case "online", "dnd", "idle", "invisible":
		err := a.setStatus(ctx.Session, status)
		if err != nil {
			_ = ctx.RespondText("An error occurred while changing your status: " + err.Error())
			return
		}

		c.services.Embed.NewEmbed().
			SetDescription(fmt.Sprintf("Successfully set status to `%s`!", status)).
			Send(ctx)
	default:
		_ = ctx.RespondText("Invalid status argument. Use `online`, `dnd`, `idle`, or `invisible`.")
	}
}

func (a *ActivityCommands) customstatus(ctx *router.Ctx) {
	err := a.setCustomStatus(ctx.Session, ctx.Arguments.Raw())
	if err != nil {
		_ = ctx.RespondText("An error occurred while changing your custom status: " + err.Error())
		return
	}

	c.services.Embed.NewEmbed().
		SetDescription(fmt.Sprintf("Successfully set custom status to `%s`!", ctx.Arguments.Raw())).
		Send(ctx)
}

func (a *ActivityCommands) clearactivity(ctx *router.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config: " + err.Error())
		return
	}

	cfg.Activity.Enabled = false
	err = cfg.Save()
	if err != nil {
		_ = ctx.RespondText("Error saving config: " + err.Error())
		return
	}

	err = ctx.Session.UpdateStatusComplex(discordgo.UpdateStatusData{
		Game:       nil,
		Activities: []*discordgo.Activity{},
		AFK:        true,
		Status:     sessions.Status,
	})
	if err != nil {
		_ = ctx.RespondText("An error occurred while clearing your activity: " + err.Error())
		return
	}

	c.services.Embed.NewEmbed().
		SetDescription("Successfully cleared activity!").
		Send(ctx)
}

func (a *ActivityCommands) copyactivity(ctx *router.Ctx) {
	user, err := ctx.Arguments.Get(0).User()
	if err != nil {
		_ = ctx.RespondText("Could not find user")
		return
	}

	rawSaveToConfig := strings.ToLower(utils.Arguments["[save to config: on/off]"])
	saveToConfig := rawSaveToConfig == "on"

	var guilds []string

	if user.ID == ctx.Session.State.User.ID {
		for _, guild := range ctx.Session.State.Guilds {
			guilds = append(guilds, guild.ID)
		}
	} else {
		profile, err := ctx.Session.UserProfile(user.ID)
		if err != nil {
			_ = ctx.RespondText("Could not find user")
			return
		}

		for _, guild := range profile.MutualGuilds {
			guilds = append(guilds, guild.ID)
		}
	}

	for _, guildID := range guilds {
		guild, err := utils.Discord.GetGuild(ctx.Session, guildID)
		if err != nil {
			continue
		}

		member, ok := utils.Discord.GetMember(ctx.Session, guild.ID, user.ID)
		if !ok || member.Presence == nil {
			continue
		}

		ctx.Session.UpdateStatusComplex(discordgo.UpdateStatusData{
			Activities: member.Presence.Activities,
			AFK:        false,
			Status:     sessions.Status,
		})

		for _, activity := range member.Presence.Activities {
			if !saveToConfig {
				break
			}

			if activity.Type == discordgo.ActivityTypeCustom {
				continue
			}

			cfg, err := config.GetConfig()
			if err != nil {
				utils.Logging.Error("Error getting config: " + err.Error())
				return
			}

			newActivity := config.ActivityFromDiscord(activity)
			newActivity.Enabled = cfg.Activity.Enabled
			newActivity.Button1URL = cfg.Activity.Button1URL
			newActivity.Button2URL = cfg.Activity.Button2URL

			cfg.Activity = newActivity

			err = cfg.Save()
			if err != nil {
				utils.Logging.Error("Error saving config: " + err.Error())
				return
			}
		}

		c.services.Embed.NewEmbed().
			SetDescription("Successfully copied user's activity!").
			Send(ctx)

		return
	}

	c.services.Embed.NewEmbed().
		SetDescription("No activity found").
		Send(ctx)
}

func (a *ActivityCommands) getactivityjson(ctx *router.Ctx) {
	user, err := ctx.Arguments.Get(0).User()
	if err != nil {
		_ = ctx.RespondText("Could not find user")
		return
	}

	var guilds []string

	if user.ID == ctx.Session.State.User.ID {
		for _, guild := range ctx.Session.State.Guilds {
			guilds = append(guilds, guild.ID)
		}
	} else {
		profile, err := ctx.Session.UserProfile(user.ID)
		if err != nil {
			_ = ctx.RespondText("Could not find user")
			return
		}

		for _, guild := range profile.MutualGuilds {
			guilds = append(guilds, guild.ID)
		}
	}

	for _, guildID := range guilds {
		guild, err := utils.Discord.GetGuild(ctx.Session, guildID)
		if err != nil {
			continue
		}

		member, ok := utils.Discord.GetMember(ctx.Session, guild.ID, user.ID)
		if !ok || member.Presence == nil {
			continue
		}

		emb := c.services.Embed.NewEmbed()

		for index, activity := range member.Presence.Activities {
			jsonActivity, _ := json.MarshalIndent(activity, "", "    ")
			emb.AddField(fmt.Sprintf("Activity %d", index+1), string(jsonActivity))
		}

		emb.Send(ctx)

		return
	}

	c.services.Embed.NewEmbed().
		SetDescription("No activity found").
		Send(ctx)
}

func (a *ActivityCommands) setactivityjson(ctx *router.Ctx) {
	rawJson := ctx.Arguments.RawFrom(1)

	var activity *discordgo.Activity
	err := json.Unmarshal([]byte(rawJson), &activity)
	if err != nil {
		_ = ctx.RespondText("Invalid activity JSON! " + err.Error())
		return
	}

	rawSaveToConfig := strings.ToLower(utils.Arguments["<save to config: on/off>"])
	saveToConfig := rawSaveToConfig == "on"

	if saveToConfig {
		cfg, err := config.GetConfig()
		if err != nil {
			utils.Logging.Error("Error getting config: " + err.Error())
			return
		}

		newActivity := config.ActivityFromDiscord(activity)
		newActivity.Enabled = cfg.Activity.Enabled
		newActivity.Button1URL = cfg.Activity.Button1URL
		newActivity.Button2URL = cfg.Activity.Button2URL

		cfg.Activity = newActivity

		err = cfg.Save()
		if err != nil {
			utils.Logging.Error("Error saving config: " + err.Error())
			return
		}
	}

	_ = ctx.Session.UpdateStatusComplex(discordgo.UpdateStatusData{
		Activities: []*discordgo.Activity{activity},
		AFK:        true,
		Status:     sessions.Status,
	})
}

func (*ActivityCommands) cactivity(ctx *router.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	switch utils.Arguments["<get/set/toggle>"] {
	case "get":
		c.services.Embed.NewEmbed().
			AddField("Custom Activity", utils.Misc.FormatBoolOnOff(cfg.Activity.Enabled)).
			Send(ctx)
	case "set":
		toggle := strings.ToLower(utils.Arguments["[on/off]"])
		if toggle != "on" && toggle != "off" {
			_ = ctx.RespondText("Invalid value. Use `on` or `off`.")
			return
		}

		cfg.Activity.Enabled = toggle == "on"
		_ = cfg.Save()

		c.services.Embed.NewEmbed().
			SetDescription("Turned " + toggle + " Custom Activity").
			Send(ctx)
	case "toggle":
		cfg.Activity.Enabled = !cfg.Activity.Enabled
		_ = cfg.Save()

		c.services.Embed.NewEmbed().
			SetDescription("Turned " + utils.Misc.FormatBoolOnOff(cfg.Activity.Enabled) + " Custom Activity").
			Send(ctx)
	default:
		_ = ctx.RespondText("Invalid <get/set/toggle> argument. Use `get`, `set` or `toggle`.")
	}
}

func (*ActivityCommands) capauseforother(ctx *router.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	switch utils.Arguments["<get/set/toggle>"] {
	case "get":
		c.services.Embed.NewEmbed().
			AddField("Pause Custom Activity for other Activity", utils.Misc.FormatBoolOnOff(cfg.Activity.PauseForOther)).
			Send(ctx)
	case "set":
		toggle := strings.ToLower(utils.Arguments["[on/off]"])
		if toggle != "on" && toggle != "off" {
			_ = ctx.RespondText("Invalid value. Use `on` or `off`.")
			return
		}

		cfg.Activity.PauseForOther = toggle == "on"
		_ = cfg.Save()

		c.services.Embed.NewEmbed().
			SetDescription("Turned " + toggle + " Pause Custom Activity for other Activity").
			Send(ctx)
	case "toggle":
		cfg.Activity.PauseForOther = !cfg.Activity.PauseForOther
		_ = cfg.Save()

		c.services.Embed.NewEmbed().
			SetDescription("Turned " + utils.Misc.FormatBoolOnOff(cfg.Activity.PauseForOther) + " Pause Custom Activity for other Activity").
			Send(ctx)
	default:
		_ = ctx.RespondText("Invalid <get/set/toggle> argument. Use `get`, `set` or `toggle`.")
	}
}

func (*ActivityCommands) caconnectiontype(ctx *router.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	switch utils.Arguments["<get/set/values>"] {
	case "get":
		c.services.Embed.NewEmbed().
			AddField("Custom Activity Connection Type", cfg.Activity.ActivityType).
			Send(ctx)
	case "set":
		value := strings.ToLower(utils.Arguments["[value]"])

		switch value {
		case "rpc":
			cfg.Activity.ActivityType = "RPC"
		case "ca", "customactivity":
			cfg.Activity.ActivityType = "CustomActivity"
		default:
			_ = ctx.RespondText("Invalid connection type. Use `[p]caconnectiontype values` to get a list of available values.")
			return
		}

		_ = cfg.Save()

		c.services.Embed.NewEmbed().
			SetDescription("Custom Activity Connection Type set to " + cfg.Activity.ActivityType).
			Send(ctx)
	case "values":
		c.services.Embed.NewEmbed().
			SetDescription("Available connection types: `RPC`, `CA`/`CustomActivity`").
			Send(ctx)
	default:
		_ = ctx.RespondText("Invalid <get/set/values> argument. Use `get`, `set` or `values`.")
	}
}

func (*ActivityCommands) catype(ctx *router.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	switch utils.Arguments["<get/set/values>"] {
	case "get":
		c.services.Embed.NewEmbed().
			AddField("Custom Activity Type", cfg.Activity.Type).
			Send(ctx)
	case "set":
		value := strings.ToLower(utils.Arguments["[value]"])

		switch value {
		case "playing":
			cfg.Activity.Type = "Playing"
		case "listening":
			cfg.Activity.Type = "Listening"
		case "watching":
			cfg.Activity.Type = "Watching"
		case "streaming":
			cfg.Activity.Type = "Streaming"
		case "competing":
			cfg.Activity.Type = "Competing"
		case "spotify":
			cfg.Activity.Type = "Spotify"
		default:
			_ = ctx.RespondText("Invalid connection type. Use `[p]catype values` to get a list of available values.")
			return
		}

		_ = cfg.Save()

		c.services.Embed.NewEmbed().
			SetDescription("Custom Activity Connection Type set to " + cfg.Activity.Type).
			Send(ctx)
	case "values":
		c.services.Embed.NewEmbed().
			SetDescription("Available connection types: `Playing`, `Listening`, `Watching`, `Streaming`, `Competing`, `Spotify`.").
			Send(ctx)
	default:
		_ = ctx.RespondText("Invalid <get/set/values> argument. Use `get`, `set` or `values`.")
	}
}

func (*ActivityCommands) caplatform(ctx *router.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	switch utils.Arguments["<get/set/values>"] {
	case "get":
		c.services.Embed.NewEmbed().
			AddField("Custom Activity Platform", cfg.Activity.Platform).
			Send(ctx)
	case "set":
		value := strings.ToLower(utils.Arguments["[value]"])

		switch value {
		case "default":
			cfg.Activity.Platform = "Default"
		case "ps4":
			cfg.Activity.Platform = "PS4"
		case "ps5":
			cfg.Activity.Platform = "PS5"
		case "xbox":
			cfg.Activity.Platform = "Xbox"
		case "samsung":
			cfg.Activity.Platform = "Samsung"
		default:
			_ = ctx.RespondText("Invalid platform. Use `[p]caplatform values` to get a list of available values.")
			return
		}

		_ = cfg.Save()

		c.services.Embed.NewEmbed().
			SetDescription("Custom Activity Platform set to " + cfg.Activity.Platform).
			Send(ctx)
	case "values":
		c.services.Embed.NewEmbed().
			SetDescription("Available connection types: `Default`, `PS4`, `PS5`, `Xbox`, `Samsung`.").
			Send(ctx)
	default:
		_ = ctx.RespondText("Invalid <get/set/values> argument. Use `get`, `set` or `values`.")
	}
}

func (*ActivityCommands) caclientid(ctx *router.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	switch utils.Arguments["<get/set>"] {
	case "get":
		c.services.Embed.NewEmbed().
			AddField("Custom Activity Client ID", cfg.Activity.ClientID).
			Send(ctx)
	case "set":
		value := utils.Arguments["[value]"]
		if value == "" {
			_ = ctx.RespondText("You must provide a value to set!")
			return
		}

		cfg.Activity.ClientID = value
		_ = cfg.Save()

		c.services.Embed.NewEmbed().
			SetDescription("Custom Activity Client ID set to " + value).
			Send(ctx)
	default:
		_ = ctx.RespondText("Invalid <get/set> argument. Use `get` or `set`.")
	}
}

func (*ActivityCommands) caname(ctx *router.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	switch utils.Arguments["<get/set>"] {
	case "get":
		c.services.Embed.NewEmbed().
			AddField("Custom Activity Name", cfg.Activity.Name).
			Send(ctx)
	case "set":
		value := ctx.Arguments.RawFrom(1)

		cfg.Activity.Name = value
		_ = cfg.Save()

		c.services.Embed.NewEmbed().
			SetDescription("Custom Activity Name set to " + value).
			Send(ctx)
	default:
		_ = ctx.RespondText("Invalid <get/set> argument. Use `get` or `set`.")
	}
}

func (*ActivityCommands) castate(ctx *router.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	switch utils.Arguments["<get/set>"] {
	case "get":
		c.services.Embed.NewEmbed().
			AddField("Custom Activity State", cfg.Activity.State).
			Send(ctx)
	case "set":
		value := ctx.Arguments.RawFrom(1)

		cfg.Activity.State = value
		_ = cfg.Save()

		c.services.Embed.NewEmbed().
			SetDescription("Custom Activity State set to " + value).
			Send(ctx)
	default:
		_ = ctx.RespondText("Invalid <get/set> argument. Use `get` or `set`.")
	}
}

func (*ActivityCommands) cadetails(ctx *router.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	switch utils.Arguments["<get/set>"] {
	case "get":
		c.services.Embed.NewEmbed().
			AddField("Custom Activity Details", cfg.Activity.Details).
			Send(ctx)
	case "set":
		value := ctx.Arguments.RawFrom(1)

		cfg.Activity.Details = value
		_ = cfg.Save()

		c.services.Embed.NewEmbed().
			SetDescription("Custom Activity Details set to " + value).
			Send(ctx)
	default:
		_ = ctx.RespondText("Invalid <get/set> argument. Use `get` or `set`.")
	}
}

func (*ActivityCommands) cainvitable(ctx *router.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	switch utils.Arguments["<get/set/toggle>"] {
	case "get":
		c.services.Embed.NewEmbed().
			AddField("Custom Activity Invitable", utils.Misc.FormatBoolOnOff(cfg.Activity.Invitable)).
			Send(ctx)
	case "set":
		toggle := strings.ToLower(utils.Arguments["[on/off]"])
		if toggle != "on" && toggle != "off" {
			_ = ctx.RespondText("Invalid value. Use `on` or `off`.")
			return
		}

		cfg.Activity.Invitable = toggle == "on"
		_ = cfg.Save()

		c.services.Embed.NewEmbed().
			SetDescription("Turned " + toggle + " Custom Activity Invitable").
			Send(ctx)
	case "toggle":
		cfg.Activity.Invitable = !cfg.Activity.Invitable
		_ = cfg.Save()

		c.services.Embed.NewEmbed().
			SetDescription("Turned " + utils.Misc.FormatBoolOnOff(cfg.Activity.Invitable) + " Custom Activity Invitable").
			Send(ctx)
	default:
		_ = ctx.RespondText("Invalid <get/set/toggle> argument. Use `get`, `set` or `toggle`.")
	}
}

func (*ActivityCommands) catime(ctx *router.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	switch utils.Arguments["<get/set/toggle>"] {
	case "get":
		c.services.Embed.NewEmbed().
			AddField("Custom Activity Time", utils.Misc.FormatBoolOnOff(cfg.Activity.Time)).
			Send(ctx)
	case "set":
		toggle := strings.ToLower(utils.Arguments["[on/off]"])
		if toggle != "on" && toggle != "off" {
			_ = ctx.RespondText("Invalid value. Use `on` or `off`.")
			return
		}

		cfg.Activity.Time = toggle == "on"
		_ = cfg.Save()

		c.services.Embed.NewEmbed().
			SetDescription("Turned " + toggle + " Custom Activity Time").
			Send(ctx)
	case "toggle":
		cfg.Activity.Time = !cfg.Activity.Time
		_ = cfg.Save()

		c.services.Embed.NewEmbed().
			SetDescription("Turned " + utils.Misc.FormatBoolOnOff(cfg.Activity.Time) + " Custom Activity Time").
			Send(ctx)
	default:
		_ = ctx.RespondText("Invalid <get/set/toggle> argument. Use `get`, `set` or `toggle`.")
	}
}

func (*ActivityCommands) castreamingurl(ctx *router.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	switch utils.Arguments["<get/set>"] {
	case "get":
		c.services.Embed.NewEmbed().
			AddField("Custom Activity Streaming URL", cfg.Activity.StreamingURL).
			Send(ctx)
	case "set":
		value := utils.Arguments["[value]"]

		cfg.Activity.StreamingURL = value
		_ = cfg.Save()

		c.services.Embed.NewEmbed().
			SetDescription("Custom Activity Streaming URL set to " + value).
			Send(ctx)
	default:
		_ = ctx.RespondText("Invalid <get/set> argument. Use `get` or `set`.")
	}
}

func (*ActivityCommands) calargeimage(ctx *router.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	switch utils.Arguments["<get/set>"] {
	case "get":
		c.services.Embed.NewEmbed().
			AddField("Custom Activity Large Image", cfg.Activity.LargeImage).
			Send(ctx)
	case "set":
		value := utils.Arguments["[value]"]

		cfg.Activity.LargeImage = value
		_ = cfg.Save()

		c.services.Embed.NewEmbed().
			SetDescription("Custom Activity Large Image set to " + value).
			Send(ctx)
	default:
		_ = ctx.RespondText("Invalid <get/set> argument. Use `get` or `set`.")
	}
}

func (*ActivityCommands) calargetext(ctx *router.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	switch utils.Arguments["<get/set>"] {
	case "get":
		c.services.Embed.NewEmbed().
			AddField("Custom Activity Large Image Text", cfg.Activity.LargeImageText).
			Send(ctx)
	case "set":
		value := ctx.Arguments.RawFrom(1)

		cfg.Activity.LargeImageText = value
		_ = cfg.Save()

		c.services.Embed.NewEmbed().
			SetDescription("Custom Activity Large Image Text set to " + value).
			Send(ctx)
	default:
		_ = ctx.RespondText("Invalid <get/set> argument. Use `get` or `set`.")
	}
}

func (*ActivityCommands) casmallimage(ctx *router.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	switch utils.Arguments["<get/set>"] {
	case "get":
		c.services.Embed.NewEmbed().
			AddField("Custom Activity Small Image", cfg.Activity.SmallImage).
			Send(ctx)
	case "set":
		value := utils.Arguments["[value]"]

		cfg.Activity.SmallImage = value
		_ = cfg.Save()

		c.services.Embed.NewEmbed().
			SetDescription("Custom Activity Small Image set to " + value).
			Send(ctx)
	default:
		_ = ctx.RespondText("Invalid <get/set> argument. Use `get` or `set`.")
	}
}

func (*ActivityCommands) casmalltext(ctx *router.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	switch utils.Arguments["<get/set>"] {
	case "get":
		c.services.Embed.NewEmbed().
			AddField("Custom Activity Small Image Text", cfg.Activity.SmallImageText).
			Send(ctx)
	case "set":
		value := ctx.Arguments.RawFrom(1)

		cfg.Activity.SmallImageText = value
		_ = cfg.Save()

		c.services.Embed.NewEmbed().
			SetDescription("Custom Activity Small Image Text set to " + value).
			Send(ctx)
	default:
		_ = ctx.RespondText("Invalid <get/set> argument. Use `get` or `set`.")
	}
}

func (*ActivityCommands) cabutton1(ctx *router.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	switch utils.Arguments["<get/set>"] {
	case "get":
		c.services.Embed.NewEmbed().
			AddField("CA Button 1 Label", cfg.Activity.Button1Label).
			AddField("CA Button 1 URL", cfg.Activity.Button1URL).
			Send(ctx)
	case "set":
		label, url := utils.Arguments["[\"label\"]"], utils.Arguments["[url]"]

		cfg.Activity.Button1Label = label
		cfg.Activity.Button1URL = url
		_ = cfg.Save()

		c.services.Embed.NewEmbed().
			SetDescription(fmt.Sprintf("CA Button 1 set to %s (%s)", label, url)).
			Send(ctx)
	default:
		_ = ctx.RespondText("Invalid <get/set> argument. Use `get` or `set`.")
	}
}

func (*ActivityCommands) cabutton2(ctx *router.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	switch utils.Arguments["<get/set>"] {
	case "get":
		c.services.Embed.NewEmbed().
			AddField("CA Button 2 Label", cfg.Activity.Button2Label).
			AddField("CA Button 2 URL", cfg.Activity.Button2URL).
			Send(ctx)
	case "set":
		label, url := utils.Arguments["[\"label\"]"], utils.Arguments["[url]"]

		cfg.Activity.Button2Label = label
		cfg.Activity.Button2URL = url
		_ = cfg.Save()

		c.services.Embed.NewEmbed().
			SetDescription(fmt.Sprintf("CA Button 2 set to %s (%s)", label, url)).
			Send(ctx)
	default:
		_ = ctx.RespondText("Invalid <get/set> argument. Use `get` or `set`.")
	}
}

func (*ActivityCommands) capresets(ctx *router.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	switch utils.Arguments["<get/load/next/save/rm/list/clear>"] {
	case "get":
		var currentPreset string

		for name, preset := range cfg.ActivityPresets {
			if reflect.DeepEqual(preset, cfg.Activity) {
				currentPreset = name
				break
			}
		}

		if currentPreset == "" {
			c.services.Embed.NewEmbed().
				SetDescription("You must load a preset before using the get current preset command!").
				Send(ctx)
			return
		}

		c.services.Embed.NewEmbed().
			SetDescription(fmt.Sprintf("Current Custom Activity preset is `%s`.", currentPreset)).
			Send(ctx)
	case "load":
		name := utils.Arguments["[name]"]

		preset, ok := cfg.ActivityPresets[name]
		if !ok {
			c.services.Embed.NewEmbed().SetDescription(fmt.Sprintf("CA Preset `%s` not found.", name)).Send(ctx)
			return
		}

		cfg.Activity = preset
		_ = cfg.Save()

		c.services.Embed.NewEmbed().
			SetDescription(fmt.Sprintf("Loaded Custom Activity preset `%s`.", name)).
			Send(ctx)
	case "next":
		var currentPreset string

		for name, preset := range cfg.ActivityPresets {
			if reflect.DeepEqual(preset, cfg.Activity) {
				currentPreset = name
				break
			}
		}

		// Sort presets alphabetically
		var keys []string
		for k := range cfg.ActivityPresets {
			keys = append(keys, k)
		}

		slices.Sort(keys)

		// Get next preset
		nextPreset := 0
		for i, k := range keys {
			if k == currentPreset {
				nextPreset = (i + 1) % len(keys)
				break
			}
		}

		nextPresetName := keys[nextPreset]
		cfg.Activity = cfg.ActivityPresets[nextPresetName]
		_ = cfg.Save()

		c.services.Embed.NewEmbed().
			SetDescription(fmt.Sprintf("Loaded next CA preset `%s`.", nextPresetName)).
			Send(ctx)
	case "save":
		name := utils.Arguments["[name]"]
		if name == "" {
			_ = ctx.RespondText("You must provide a name to save the preset as!")
			return
		}

		if cfg.ActivityPresets == nil {
			cfg.ActivityPresets = make(map[string]config.Activity)
		}

		cfg.ActivityPresets[name] = cfg.Activity
		_ = cfg.Save()

		c.services.Embed.NewEmbed().
			SetDescription(fmt.Sprintf("Saved current activity as preset `%s`.", name)).
			Send(ctx)
	case "rm":
		name := utils.Arguments["[name]"]

		if _, ok := cfg.ActivityPresets[name]; !ok {
			c.services.Embed.NewEmbed().SetDescription(fmt.Sprintf("Preset `%s` not found.", name)).Send(ctx)
			return
		}

		delete(cfg.ActivityPresets, name)
		_ = cfg.Save()

		c.services.Embed.NewEmbed().
			SetDescription(fmt.Sprintf("Deleted activity preset `%s`.", name)).
			Send(ctx)
	case "list":
		var presets []string

		for k, v := range cfg.ActivityPresets {
			presets = append(presets, fmt.Sprintf("%s: %s %s", k, v.Type, v.Name))
		}

		if len(presets) == 0 {
			c.services.Embed.NewEmbed().SetDescription("No Custom Activity presets found.").Send(ctx)
			return
		}

		c.services.Embed.NewEmbed().
			AddField("Custom Activity Presets", strings.Join(presets, "\n")).
			Send(ctx)
	case "clear":
		cfg.ActivityPresets = make(map[string]config.Activity)
		_ = cfg.Save()

		c.services.Embed.NewEmbed().
			SetDescription("Cleared all activity presets.").
			Send(ctx)
	default:
		_ = ctx.RespondText("Invalid <get/load/next/save/rm/list/clear> argument. Use `get`, `load`, `next`, `save`, `rm`, `list` or `clear`.")
	}
}

func (*ActivityCommands) cacreateimagekey(ctx *router.Ctx) {
	url := utils.Arguments["<image url>"]

	if strings.Contains(url, "discordapp.com") {
		_ = ctx.RespondText("Your image seems to be hosted on Discord. Please use a different link to create an image key, as this might not work. You can try uploading it to a free image hoster like Imgur or ImgBB (but make sure to use the direct link to the image if you do. The URL should end with .png, .jpg or .gif).")
	}

	imageKey := api.GenImageKey(url)

	c.services.Embed.NewEmbed().
		AddField("Image URL", utils.Arguments["<image url>"]).
		AddField("Image Key", imageKey).
		Send(ctx)
}

func (*ActivityCommands) cacreatespotifykey(ctx *router.Ctx) {
	bearerToken := utils.Spotify.GetBearerToken(ctx.Session.Identify.Token)

	if bearerToken == "" {
		_ = ctx.RespondText("Could not find a Spotify connection. Please connect your Spotify account to your Discord.")
		return
	}

	httpResponse := http.Get(
		fmt.Sprintf("https://api.spotify.com/v1/search?q=%s&type=track", url.QueryEscape(ctx.Arguments.Raw())),
		map[string]string{"Authorization": bearerToken},
	)
	json := httpResponse.Json()

	song := json.GetPath("tracks", "items").GetIndex(0)

	imageUrl := song.Get("album").Get("images").GetIndex(0).Get("url").MustString()
	if httpResponse.Error != nil || imageUrl == "" {
		_ = ctx.RespondText("A Spotify error occurred: Song not found.")
		return
	}

	splitUrl := strings.Split(imageUrl, "/")

	c.services.Embed.NewEmbed().
		AddField("Song", song.Get("name").MustString()).
		AddField("Image Key", "spotify:"+splitUrl[len(splitUrl)-1]).
		Send(ctx)
}

func (ActivityCommands) activitypresets(ctx *router.Ctx) {
	switch strings.ToLower(utils.Arguments["<set/values>"]) {
	case "set":
		name := strings.ToLower(utils.Arguments["[name]"])
		text := ctx.Arguments.RawFrom(2)

		cfg, err := config.GetConfig()
		if err != nil {
			_ = ctx.RespondText("Error getting config: " + err.Error())
			return
		}

		activity := config.Activity{
			Enabled:      true,
			ActivityType: "CustomActivity",
			Platform:     "Default",
			ClientID:     "898932653049393172",
			Time:         true,
		}

		switch name {
		case "youtube":
			activity.Name = "YouTube"
			activity.Details = text
			activity.Type = "Watching"
			activity.LargeImage = "spotify:ab67706c0000bebb7d1c495e63620a99f3cbcbf5"
			activity.LargeImageText = "YouTube"
		case "twitch":
			activity.Name = "Twitch"
			activity.Details = text
			activity.Type = "Watching"
			activity.LargeImage = "spotify:ab67706c0000bebb4081b043d43093915cb73114"
			activity.LargeImageText = "Twitch"
		case "discord":
			activity.Name = "Discord"
			activity.Details = "Discord"
			activity.State = text
			activity.Type = "Watching"
			activity.LargeImage = "spotify:ab67616d0000b273db20400815ca2a9481bd255e"
		case "rickroll":
			activity.Name = "Never Gonna Give You Up"
			activity.Details = "Never Gonna Give You Up"
			activity.State = text
			activity.Type = "Listening"
			activity.LargeImage = "spotify:ab67616d0000b2735755e164993798e0c9ef7d7a"
			activity.LargeImageText = "Rick Astley"
		case "gta6":
			activity.Name = "Grand Theft Auto VI"
			activity.Details = "Grand Theft Auto VI"
			activity.State = text
			activity.Type = "Playing"
			activity.Platform = "PS5"
			activity.LargeImage = "mp:external/mLiXAArsDKgwP6Wb7mvDrBn4mFbJUZLxndTK5D2dL8s/https/i.vgy.me/JhGfvq.jpg"
		default:
			_ = ctx.RespondText("Invalid preset name. Use `youtube`, `twitch`, `discord`, `rickroll`, or `gta6`.")
			return
		}

		cfg.Activity = activity
		_ = cfg.Save()

		c.services.Embed.NewEmbed().
			SetDescription(fmt.Sprintf("Successfully set activity preset `%s`!", name)).
			Send(ctx)
	case "values":
		_ = ctx.RespondText("Available activity presets: `youtube`, `twitch`, `discord`, `rickroll`, `gta6`.")
	default:
		_ = ctx.RespondText("Invalid argument. Use `set` or `values`.")
	}
}
