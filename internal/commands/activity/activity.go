package activity

import (
	"fmt"
	"reflect"
	"slices"
	"strings"

	"github.com/goccy/go-json"

	"github.com/lb-selfbot/selfbot/internal/api"
	"github.com/lb-selfbot/selfbot/internal/logging"
	"github.com/lb-selfbot/selfbot/internal/services/config"
	"github.com/lb-selfbot/selfbot/internal/services/router"

	"github.com/lb-selfbot/discordgo"
	"github.com/lb-selfbot/discordgo/protos"
	"github.com/lb-selfbot/selfbot/internal/http"

	"github.com/lb-selfbot/selfbot/internal/services"
	"github.com/lb-selfbot/selfbot/internal/utils"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type ActivityCommands struct{ services *services.Services }

func NewActivityCommands(s *services.Services) *ActivityCommands {
	return &ActivityCommands{services: s}
}

func (c *ActivityCommands) Register() {
	c.services.Router.RegisterCategory(&router.Category{
		Name:        "Activity",
		Description: "Customize your Discord activity",
	}, []*router.Command{
		{
			Name:        "status",
			Description: "Set your Discord status",
			Arguments:   "<online/dnd/idle/invisible>",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.status,
		},
		{
			Name:        "customstatus",
			Description: "Change your custom status",
			Arguments:   "[text]",
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
			Arguments:   "<user> [save to config: on/off]",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.copyactivity,
		},
		{
			Name:        "getactivityjson",
			Description: "Get a users activity in JSON format",
			Arguments:   "<user>",
			Handler:     c.getactivityjson,
		},
		{
			Name:        "setactivityjson",
			Description: "Set your activity using raw JSON",
			Arguments:   "<save to config: on/off> <json>",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.setactivityjson,
		},
		{
			Name:        "cactivity",
			Description: "Custom Activity commands",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			SubCommands: []*router.Command{
				{
					Name:        "get",
					Description: "Get Custom Activity status",
					Handler:     c.getCAEnabled,
				},
				{
					Name:        "set",
					Description: "Set Custom Activity status",
					Arguments:   "<on/off>",
					Handler:     c.setCAEnabled,
				},
				{
					Name:        "toggle",
					Description: "Toggle Custom Activity status",
					Handler:     c.toggleCAEnabled,
				},
			},
		},
		{
			Name:        "capauseforother",
			Description: "Pause Custom Activity for other Activity commands",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			SubCommands: []*router.Command{
				{
					Name:        "get",
					Description: "Get Pause CA status",
					Handler:     c.getCAPauseForOther,
				},
				{
					Name:        "set",
					Description: "Set Pause CA status",
					Arguments:   "<on/off>",
					Handler:     c.setCAPauseForOther,
				},
				{
					Name:        "toggle",
					Description: "Toggle Pause CA status",
					Handler:     c.toggleCAPauseForOther,
				},
			},
		},
		{
			Name:        "caconnectiontype",
			Description: "Custom Activity Connection Type commands",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			SubCommands: []*router.Command{
				{
					Name:        "get",
					Description: "Get the Custom Activity Connection Type",
					Handler:     c.getCAConnectionType,
				},
				{
					Name:        "set",
					Description: "Set the Custom Activity Connection Type",
					Arguments:   "<type>",
					Handler:     c.setCAConnectionType,
				},
				{
					Name:        "values",
					Description: "List available connection types",
					Handler:     c.listCAConnectionTypes,
				},
			},
		},
		{
			Name:        "catype",
			Description: "Custom Activity Type commands",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			SubCommands: []*router.Command{
				{
					Name:        "get",
					Description: "Get the Custom Activity Type",
					Handler:     c.getCAType,
				},
				{
					Name:        "set",
					Description: "Set the Custom Activity Type",
					Arguments:   "<type>",
					Handler:     c.setCAType,
				},
				{
					Name:        "values",
					Description: "List available types",
					Handler:     c.listCATypes,
				},
			},
		},
		{
			Name:        "caplatform",
			Description: "Custom Activity Platform commands",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			SubCommands: []*router.Command{
				{
					Name:        "get",
					Description: "Get the Custom Activity Platform",
					Handler:     c.getCAPlatform,
				},
				{
					Name:        "set",
					Description: "Set the Custom Activity Platform",
					Arguments:   "<platform>",
					Handler:     c.setCAPlatform,
				},
				{
					Name:        "values",
					Description: "List available platforms",
					Handler:     c.listCAPlatforms,
				},
			},
		},
		{
			Name:        "cainvitable",
			Description: "Custom Activity Invitable commands",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			SubCommands: []*router.Command{
				{
					Name:        "get",
					Description: "Get Custom Activity Invitable status",
					Handler:     c.getCAInvitable,
				},
				{
					Name:        "set",
					Description: "Set Custom Activity Invitable status",
					Arguments:   "<on/off>",
					Handler:     c.setCAInvitable,
				},
				{
					Name:        "toggle",
					Description: "Toggle Custom Activity Invitable status",
					Handler:     c.toggleCAInvitable,
				},
			},
		},
		{
			Name:        "catime",
			Description: "Custom Activity Time commands",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			SubCommands: []*router.Command{
				{
					Name:        "get",
					Description: "Get Custom Activity Time status",
					Handler:     c.getCATime,
				},
				{
					Name:        "set",
					Description: "Set Custom Activity Time status",
					Arguments:   "<on/off>",
					Handler:     c.setCATime,
				},
				{
					Name:        "toggle",
					Description: "Toggle Custom Activity Time status",
					Handler:     c.toggleCATime,
				},
			},
		},
		{
			Name:        "caname",
			Description: "Custom Activity Name commands",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			SubCommands: []*router.Command{
				{
					Name:        "get",
					Description: "Get the Custom Activity Name",
					Handler:     c.getCAName,
				},
				{
					Name:        "set",
					Description: "Set the Custom Activity Name",
					Arguments:   "<name>",
					Handler:     c.setCAName,
				},
			},
		},
		{
			Name:        "castate",
			Description: "Custom Activity State commands",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			SubCommands: []*router.Command{
				{
					Name:        "get",
					Description: "Get the Custom Activity State",
					Handler:     c.getCAState,
				},
				{
					Name:        "set",
					Description: "Set the Custom Activity State",
					Arguments:   "<state>",
					Handler:     c.setCAState,
				},
			},
		},
		{
			Name:        "cadetails",
			Description: "Custom Activity Details commands",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			SubCommands: []*router.Command{
				{
					Name:        "get",
					Description: "Get the Custom Activity Details",
					Handler:     c.getCADetails,
				},
				{
					Name:        "set",
					Description: "Set the Custom Activity Details",
					Arguments:   "<details>",
					Handler:     c.setCADetails,
				},
			},
		},
		{
			Name:        "castreamingurl",
			Description: "Custom Activity Streaming URL commands",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			SubCommands: []*router.Command{
				{
					Name:        "get",
					Description: "Get the Custom Activity Streaming URL",
					Handler:     c.getCAStreamingURL,
				},
				{
					Name:        "set",
					Description: "Set the Custom Activity Streaming URL",
					Arguments:   "<url>",
					Handler:     c.setCAStreamingURL,
				},
			},
		},
		{
			Name:        "calargeimage",
			Description: "Custom Activity Large Image commands",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			SubCommands: []*router.Command{
				{
					Name:        "get",
					Description: "Get the Custom Activity Large Image",
					Handler:     c.getCALargeImage,
				},
				{
					Name:        "set",
					Description: "Set the Custom Activity Large Image",
					Arguments:   "<image>",
					Handler:     c.setCALargeImage,
				},
			},
		},
		{
			Name:        "calargetext",
			Description: "Custom Activity Large Text commands",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			SubCommands: []*router.Command{
				{
					Name:        "get",
					Description: "Get the Custom Activity Large Text",
					Handler:     c.getCALargeText,
				},
				{
					Name:        "set",
					Description: "Set the Custom Activity Large Text",
					Arguments:   "<text>",
					Handler:     c.setCALargeText,
				},
			},
		},
		{
			Name:        "casmallimage",
			Description: "Custom Activity Small Image commands",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			SubCommands: []*router.Command{
				{
					Name:        "get",
					Description: "Get the Custom Activity Small Image",
					Handler:     c.getCASmallImage,
				},
				{
					Name:        "set",
					Description: "Set the Custom Activity Small Image",
					Arguments:   "<image>",
					Handler:     c.setCASmallImage,
				},
			},
		},
		{
			Name:        "casmalltext",
			Description: "Custom Activity Small Text commands",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			SubCommands: []*router.Command{
				{
					Name:        "get",
					Description: "Get the Custom Activity Small Text",
					Handler:     c.getCASmallText,
				},
				{
					Name:        "set",
					Description: "Set the Custom Activity Small Text",
					Arguments:   "<text>",
					Handler:     c.setCASmallText,
				},
			},
		},
		{
			Name:        "cainvitable",
			Description: "Custom Activity Invitable commands",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			SubCommands: []*router.Command{
				{
					Name:        "get",
					Description: "Get Custom Activity Invitable status",
					Handler:     c.getCAInvitable,
				},
				{
					Name:        "set",
					Description: "Set Custom Activity Invitable status",
					Arguments:   "<on/off>",
					Handler:     c.setCAInvitable,
				},
				{
					Name:        "toggle",
					Description: "Toggle Custom Activity Invitable status",
					Handler:     c.toggleCAInvitable,
				},
			},
		},
		{
			Name:        "catime",
			Description: "Custom Activity Time commands",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			SubCommands: []*router.Command{
				{
					Name:        "get",
					Description: "Get Custom Activity Time status",
					Handler:     c.getCATime,
				},
				{
					Name:        "set",
					Description: "Set Custom Activity Time status",
					Arguments:   "<on/off>",
					Handler:     c.setCATime,
				},
				{
					Name:        "toggle",
					Description: "Toggle Custom Activity Time status",
					Handler:     c.toggleCATime,
				},
			},
		},
		{
			Name:        "cacreateimagekey",
			Description: "Create an image key for Custom Activity",
			Arguments:   "<image url>",
			Handler:     c.cacreateimagekey,
		},
		{
			Name:        "cacreatespotifykey",
			Description: "Create an Spotify image key for Custom Activity",
			Arguments:   "<song name>",
			Handler:     c.cacreatespotifykey,
		},
		{
			Name:        "cabutton1",
			Description: "Custom Activity Button 1 commands",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			SubCommands: []*router.Command{
				{
					Name:        "get",
					Description: "Get Button 1 settings",
					Handler:     c.getCAButton1,
				},
				{
					Name:        "set",
					Description: "Set Button 1 settings",
					Arguments:   "<\"label\"> <url>",
					Handler:     c.setCAButton1,
				},
			},
		},
		{
			Name:        "cabutton2",
			Description: "Custom Activity Button 2 commands",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			SubCommands: []*router.Command{
				{
					Name:        "get",
					Description: "Get Button 2 settings",
					Handler:     c.getCAButton2,
				},
				{
					Name:        "set",
					Description: "Set Button 2 settings",
					Arguments:   "<\"label\"> <url>",
					Handler:     c.setCAButton2,
				},
			},
		},
		{
			Name:        "capresets",
			Description: "Custom Activity Presets commands",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			SubCommands: []*router.Command{
				{
					Name:        "get",
					Description: "Get current preset",
					Handler:     c.getCAPreset,
				},
				{
					Name:        "load",
					Description: "Load a preset",
					Arguments:   "<name>",
					Handler:     c.loadCAPreset,
				},
				{
					Name:        "next",
					Description: "Load next preset",
					Handler:     c.nextCAPreset,
				},
				{
					Name:        "save",
					Description: "Save current activity as preset",
					Arguments:   "<name>",
					Handler:     c.saveCAPreset,
				},
				{
					Name:        "rm",
					Description: "Delete a preset",
					Arguments:   "<name>",
					Handler:     c.removeCAPreset,
				},
				{
					Name:        "list",
					Description: "List all presets",
					Handler:     c.listCAPresets,
				},
				{
					Name:        "clear",
					Description: "Clear all presets",
					Handler:     c.clearCAPresets,
				},
			},
		},
		{
			Name:        "activitypresets",
			Description: "Built-in Activity Presets commands",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			SubCommands: []*router.Command{
				{
					Name:        "set",
					Description: "Set a built-in preset",
					Arguments:   "<name> [text]",
					Handler:     c.setActivityPreset,
				},
				{
					Name:        "values",
					Description: "List available presets",
					Handler:     c.listActivityPresets,
				},
			},
		},
	})
}

func (c *ActivityCommands) setStatus(session *discordgo.Session, status string) error {
	session.State.UserSettings.Status.Status = wrapperspb.String(status)
	err := session.UserSettingsUpdate(session.State.UserSettings, true)
	if err != nil {
		return err
	}

	err = session.UpdateStatusComplex(discordgo.UpdateStatusData{
		AFK:    true,
		Status: status,
	})
	if err != nil {
		return err
	}

	c.services.Sessions.SetStatus(status)

	return nil
}

func (c *ActivityCommands) setCustomStatus(session *discordgo.Session, status string) error {
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
		Status:     c.services.Sessions.GetStatus(),
	})
}

func (c *ActivityCommands) status(ctx *router.Ctx) {
	status := strings.ToLower(ctx.Arguments.Get(0).Raw())

	switch status {
	case "online", "dnd", "idle", "invisible":
		err := c.setStatus(ctx.Session, status)
		if err != nil {
			ctx.RespondText("An error occurred while changing your status: %v", err)
			return
		}

		c.services.Embed.New(ctx).
			SetDescription("Successfully set status to `%s`!", status).
			Send()
	default:
		ctx.RespondText("Invalid status argument. Use `online`, `dnd`, `idle`, or `invisible`.")
	}
}

func (c *ActivityCommands) customstatus(ctx *router.Ctx) {
	err := c.setCustomStatus(ctx.Session, ctx.Arguments.Raw())
	if err != nil {
		ctx.RespondText("An error occurred while changing your custom status: %v", err)
		return
	}

	c.services.Embed.New(ctx).
		SetDescription("Successfully set custom status to `%s`!", ctx.Arguments.Raw()).
		Send()
}

func (c *ActivityCommands) clearactivity(ctx *router.Ctx) {
	cfg := c.services.Config.GetConfig()
	cfg.Activity.Enabled = false
	if err := c.services.Config.SaveConfig(cfg); err != nil {
		ctx.RespondText("Error saving config: %v", err)
		return
	}

	err := ctx.Session.UpdateStatusComplex(discordgo.UpdateStatusData{
		Game:       nil,
		Activities: []*discordgo.Activity{},
		AFK:        true,
		Status:     c.services.Sessions.GetStatus(),
	})
	if err != nil {
		ctx.RespondText("An error occurred while clearing your activity: %v", err)
		return
	}

	c.services.Embed.New(ctx).
		SetDescription("Successfully cleared activity!").
		Send()
}

func (c *ActivityCommands) copyactivity(ctx *router.Ctx) {
	user, err := ctx.Arguments.Get(0).User()
	if err != nil {
		ctx.RespondText("Could not find user")
		return
	}

	saveToConfig, _ := ctx.Arguments.Get(1).Bool()

	var guilds []string

	if user.ID == ctx.Session.State.User.ID {
		for _, guild := range ctx.Session.State.Guilds {
			guilds = append(guilds, guild.ID)
		}
	} else {
		profile, err := ctx.Session.UserProfile(user.ID)
		if err != nil {
			ctx.RespondText("Could not find user")
			return
		}

		for _, guild := range profile.MutualGuilds {
			guilds = append(guilds, guild.ID)
		}
	}

	for _, guildID := range guilds {
		guild, err := ctx.Session.GuildWithFallback(guildID)
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
			Status:     c.services.Sessions.GetStatus(),
		})

		for _, activity := range member.Presence.Activities {
			if !saveToConfig {
				break
			}

			if activity.Type == discordgo.ActivityTypeCustom {
				continue
			}

			cfg := c.services.Config.GetConfig()

			newActivity := config.ActivityFromDiscord(activity)
			newActivity.Enabled = cfg.Activity.Enabled
			newActivity.Button1URL = cfg.Activity.Button1URL
			newActivity.Button2URL = cfg.Activity.Button2URL

			cfg.Activity = newActivity

			if err := c.services.Config.SaveConfig(cfg); err != nil {
				logging.Error("Error saving config: " + err.Error())
				return
			}
		}

		c.services.Embed.New(ctx).
			SetDescription("Successfully copied user's activity!").
			Send()

		return
	}

	c.services.Embed.New(ctx).
		SetDescription("No activity found").
		Send()
}

func (c *ActivityCommands) getactivityjson(ctx *router.Ctx) {
	user, err := ctx.Arguments.Get(0).User()
	if err != nil {
		ctx.RespondText("Could not find user")
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
			ctx.RespondText("Could not find user")
			return
		}

		for _, guild := range profile.MutualGuilds {
			guilds = append(guilds, guild.ID)
		}
	}

	for _, guildID := range guilds {
		guild, err := ctx.Session.GuildWithFallback(guildID)
		if err != nil {
			continue
		}

		member, ok := utils.Discord.GetMember(ctx.Session, guild.ID, user.ID)
		if !ok || member.Presence == nil {
			continue
		}

		emb := c.services.Embed.New(ctx)

		for index, activity := range member.Presence.Activities {
			jsonActivity, _ := json.MarshalIndent(activity, "", "    ")
			emb.AddField(fmt.Sprintf("Activity %d", index+1), string(jsonActivity))
		}

		emb.Send()

		return
	}

	c.services.Embed.New(ctx).
		SetDescription("No activity found").
		Send()
}

func (c *ActivityCommands) setactivityjson(ctx *router.Ctx) {
	rawJson := ctx.Arguments.RawFrom(1)

	var activity *discordgo.Activity
	err := json.Unmarshal([]byte(rawJson), &activity)
	if err != nil {
		ctx.RespondText("Invalid activity JSON! %v", err)
		return
	}

	saveToConfig, _ := ctx.Arguments.Get(0).Bool()

	if saveToConfig {
		cfg := c.services.Config.GetConfig()

		newActivity := config.ActivityFromDiscord(activity)
		newActivity.Enabled = cfg.Activity.Enabled
		newActivity.Button1URL = cfg.Activity.Button1URL
		newActivity.Button2URL = cfg.Activity.Button2URL

		cfg.Activity = newActivity

		if err := c.services.Config.SaveConfig(cfg); err != nil {
			logging.Error("Error saving config: " + err.Error())
			return
		}
	}

	_ = ctx.Session.UpdateStatusComplex(discordgo.UpdateStatusData{
		Activities: []*discordgo.Activity{activity},
		AFK:        true,
		Status:     c.services.Sessions.GetStatus(),
	})
}

func (c *ActivityCommands) getCAEnabled(ctx *router.Ctx) {
	cfg := c.services.Config.GetConfig()

	c.services.Embed.New(ctx).
		AddField("Custom Activity", utils.Misc.FormatBoolOnOff(cfg.Activity.Enabled)).
		Send()
}

func (c *ActivityCommands) setCAEnabled(ctx *router.Ctx) {
	cfg := c.services.Config.GetConfig()

	toggle, err := ctx.Arguments.Get(0).Bool()
	if err != nil {
		ctx.RespondText("Invalid value. Use `on` or `off`.")
		return
	}

	cfg.Activity.Enabled = toggle
	_ = c.services.Config.SaveConfig(cfg)

	c.services.Embed.New(ctx).
		SetDescription("%s", "Turned "+utils.Misc.FormatBoolOnOff(toggle)+" Custom Activity").
		Send()
}

func (c *ActivityCommands) toggleCAEnabled(ctx *router.Ctx) {
	cfg := c.services.Config.GetConfig()

	cfg.Activity.Enabled = !cfg.Activity.Enabled
	_ = c.services.Config.SaveConfig(cfg)

	c.services.Embed.New(ctx).
		SetDescription("%s", "Turned "+utils.Misc.FormatBoolOnOff(cfg.Activity.Enabled)+" Custom Activity").
		Send()
}

func (c *ActivityCommands) getCAPauseForOther(ctx *router.Ctx) {
	cfg := c.services.Config.GetConfig()

	c.services.Embed.New(ctx).
		AddField("Pause Custom Activity for other Activity", utils.Misc.FormatBoolOnOff(cfg.Activity.PauseForOther)).
		Send()
}

func (c *ActivityCommands) setCAPauseForOther(ctx *router.Ctx) {
	cfg := c.services.Config.GetConfig()

	toggle, err := ctx.Arguments.Get(0).Bool()
	if err != nil {
		ctx.RespondText("Invalid value. Use `on` or `off`.")
		return
	}

	cfg.Activity.PauseForOther = toggle
	_ = c.services.Config.SaveConfig(cfg)

	c.services.Embed.New(ctx).
		SetDescription("%s", "Turned "+utils.Misc.FormatBoolOnOff(toggle)+" Pause Custom Activity for other Activity").
		Send()
}

func (c *ActivityCommands) toggleCAPauseForOther(ctx *router.Ctx) {
	cfg := c.services.Config.GetConfig()

	cfg.Activity.PauseForOther = !cfg.Activity.PauseForOther
	_ = c.services.Config.SaveConfig(cfg)

	c.services.Embed.New(ctx).
		SetDescription("%s", "Turned "+utils.Misc.FormatBoolOnOff(cfg.Activity.PauseForOther)+" Pause Custom Activity for other Activity").
		Send()
}

func (c *ActivityCommands) getCAConnectionType(ctx *router.Ctx) {
	cfg := c.services.Config.GetConfig()

	c.services.Embed.New(ctx).
		AddField("Custom Activity Connection Type", cfg.Activity.ActivityType).
		Send()
}

func (c *ActivityCommands) setCAConnectionType(ctx *router.Ctx) {
	cfg := c.services.Config.GetConfig()

	value := strings.ToLower(ctx.Arguments.Raw())

	switch value {
	case "rpc":
		cfg.Activity.ActivityType = "RPC"
	case "ca", "customactivity":
		cfg.Activity.ActivityType = "CustomActivity"
	default:
		ctx.RespondText("Invalid connection type. Use `[p]caconnectiontype values` to get a list of available values.")
		return
	}

	_ = c.services.Config.SaveConfig(cfg)

	c.services.Embed.New(ctx).
		SetDescription("%s", "Custom Activity Connection Type set to "+cfg.Activity.ActivityType).
		Send()
}

func (c *ActivityCommands) listCAConnectionTypes(ctx *router.Ctx) {
	c.services.Embed.New(ctx).
		SetDescription("Available connection types: `RPC`, `CA`/`CustomActivity`").
		Send()
}

func (c *ActivityCommands) getCAType(ctx *router.Ctx) {
	cfg := c.services.Config.GetConfig()

	c.services.Embed.New(ctx).
		AddField("Custom Activity Type", cfg.Activity.Type).
		Send()
}

func (c *ActivityCommands) setCAType(ctx *router.Ctx) {
	cfg := c.services.Config.GetConfig()

	value := strings.ToLower(ctx.Arguments.Raw())

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
		ctx.RespondText("Invalid connection type. Use `[p]catype values` to get a list of available values.")
		return
	}

	_ = c.services.Config.SaveConfig(cfg)

	c.services.Embed.New(ctx).
		SetDescription("%s", "Custom Activity Type set to "+cfg.Activity.Type).
		Send()
}

func (c *ActivityCommands) listCATypes(ctx *router.Ctx) {
	c.services.Embed.New(ctx).
		SetDescription("Available connection types: `Playing`, `Listening`, `Watching`, `Streaming`, `Competing`, `Spotify`.").
		Send()
}

func (c *ActivityCommands) getCAPlatform(ctx *router.Ctx) {
	cfg := c.services.Config.GetConfig()

	c.services.Embed.New(ctx).
		AddField("Custom Activity Platform", cfg.Activity.Platform).
		Send()
}

func (c *ActivityCommands) setCAPlatform(ctx *router.Ctx) {
	cfg := c.services.Config.GetConfig()

	value := strings.ToLower(ctx.Arguments.Raw())

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
		ctx.RespondText("Invalid platform. Use `[p]caplatform values` to get a list of available values.")
		return
	}

	_ = c.services.Config.SaveConfig(cfg)

	c.services.Embed.New(ctx).
		SetDescription("%s", "Custom Activity Platform set to "+cfg.Activity.Platform).
		Send()
}

func (c *ActivityCommands) listCAPlatforms(ctx *router.Ctx) {
	c.services.Embed.New(ctx).
		SetDescription("Available connection types: `Default`, `PS4`, `PS5`, `Xbox`, `Samsung`.").
		Send()
}

func (c *ActivityCommands) getCAInvitable(ctx *router.Ctx) {
	cfg := c.services.Config.GetConfig()

	c.services.Embed.New(ctx).
		AddField("Custom Activity Invitable", utils.Misc.FormatBoolOnOff(cfg.Activity.Invitable)).
		Send()
}

func (c *ActivityCommands) setCAInvitable(ctx *router.Ctx) {
	cfg := c.services.Config.GetConfig()

	toggle, err := ctx.Arguments.Get(0).Bool()
	if err != nil {
		ctx.RespondText("Invalid value. Use `on` or `off`.")
		return
	}

	cfg.Activity.Invitable = toggle
	_ = c.services.Config.SaveConfig(cfg)

	c.services.Embed.New(ctx).
		SetDescription("%s", "Turned "+utils.Misc.FormatBoolOnOff(toggle)+" Custom Activity Invitable").
		Send()
}

func (c *ActivityCommands) toggleCAInvitable(ctx *router.Ctx) {
	cfg := c.services.Config.GetConfig()

	cfg.Activity.Invitable = !cfg.Activity.Invitable
	_ = c.services.Config.SaveConfig(cfg)

	c.services.Embed.New(ctx).
		SetDescription("%s", "Turned "+utils.Misc.FormatBoolOnOff(cfg.Activity.Invitable)+" Custom Activity Invitable").
		Send()
}

func (c *ActivityCommands) getCATime(ctx *router.Ctx) {
	cfg := c.services.Config.GetConfig()

	c.services.Embed.New(ctx).
		AddField("Custom Activity Time", utils.Misc.FormatBoolOnOff(cfg.Activity.Time)).
		Send()
}

func (c *ActivityCommands) setCATime(ctx *router.Ctx) {
	cfg := c.services.Config.GetConfig()

	toggle, err := ctx.Arguments.Get(0).Bool()
	if err != nil {
		ctx.RespondText("Invalid value. Use `on` or `off`.")
		return
	}

	cfg.Activity.Time = toggle
	_ = c.services.Config.SaveConfig(cfg)

	c.services.Embed.New(ctx).
		SetDescription("%s", "Turned "+utils.Misc.FormatBoolOnOff(toggle)+" Custom Activity Time").
		Send()
}

func (c *ActivityCommands) toggleCATime(ctx *router.Ctx) {
	cfg := c.services.Config.GetConfig()

	cfg.Activity.Time = !cfg.Activity.Time
	_ = c.services.Config.SaveConfig(cfg)

	c.services.Embed.New(ctx).
		SetDescription("%s", "Turned "+utils.Misc.FormatBoolOnOff(cfg.Activity.Time)+" Custom Activity Time").
		Send()
}

func (c *ActivityCommands) getCAName(ctx *router.Ctx) {
	cfg := c.services.Config.GetConfig()

	c.services.Embed.New(ctx).
		AddField("Custom Activity Name", cfg.Activity.Name).
		Send()
}

func (c *ActivityCommands) setCAName(ctx *router.Ctx) {
	cfg := c.services.Config.GetConfig()

	value := ctx.Arguments.Raw()

	cfg.Activity.Name = value
	_ = c.services.Config.SaveConfig(cfg)

	c.services.Embed.New(ctx).
		SetDescription("%s", "Custom Activity Name set to "+value).
		Send()
}

func (c *ActivityCommands) getCAState(ctx *router.Ctx) {
	cfg := c.services.Config.GetConfig()

	c.services.Embed.New(ctx).
		AddField("Custom Activity State", cfg.Activity.State).
		Send()
}

func (c *ActivityCommands) setCAState(ctx *router.Ctx) {
	cfg := c.services.Config.GetConfig()

	value := ctx.Arguments.Raw()

	cfg.Activity.State = value
	_ = c.services.Config.SaveConfig(cfg)

	c.services.Embed.New(ctx).
		SetDescription("%s", "Custom Activity State set to "+value).
		Send()
}

func (c *ActivityCommands) getCADetails(ctx *router.Ctx) {
	cfg := c.services.Config.GetConfig()

	c.services.Embed.New(ctx).
		AddField("Custom Activity Details", cfg.Activity.Details).
		Send()
}

func (c *ActivityCommands) setCADetails(ctx *router.Ctx) {
	cfg := c.services.Config.GetConfig()

	value := ctx.Arguments.Raw()

	cfg.Activity.Details = value
	_ = c.services.Config.SaveConfig(cfg)

	c.services.Embed.New(ctx).
		SetDescription("%s", "Custom Activity Details set to "+value).
		Send()
}

func (c *ActivityCommands) getCAStreamingURL(ctx *router.Ctx) {
	cfg := c.services.Config.GetConfig()

	c.services.Embed.New(ctx).
		AddField("Custom Activity Streaming URL", cfg.Activity.StreamingURL).
		Send()
}

func (c *ActivityCommands) setCAStreamingURL(ctx *router.Ctx) {
	cfg := c.services.Config.GetConfig()

	value := ctx.Arguments.Raw()

	cfg.Activity.StreamingURL = value
	_ = c.services.Config.SaveConfig(cfg)

	c.services.Embed.New(ctx).
		SetDescription("%s", "Custom Activity Streaming URL set to "+value).
		Send()
}

func (c *ActivityCommands) getCALargeImage(ctx *router.Ctx) {
	cfg := c.services.Config.GetConfig()

	c.services.Embed.New(ctx).
		AddField("Custom Activity Large Image", cfg.Activity.LargeImage).
		Send()
}

func (c *ActivityCommands) setCALargeImage(ctx *router.Ctx) {
	cfg := c.services.Config.GetConfig()

	value := ctx.Arguments.Raw()

	cfg.Activity.LargeImage = value
	_ = c.services.Config.SaveConfig(cfg)

	c.services.Embed.New(ctx).
		SetDescription("%s", "Custom Activity Large Image set to "+value).
		Send()
}

func (c *ActivityCommands) getCALargeText(ctx *router.Ctx) {
	cfg := c.services.Config.GetConfig()

	c.services.Embed.New(ctx).
		AddField("Custom Activity Large Image Text", cfg.Activity.LargeImageText).
		Send()
}

func (c *ActivityCommands) setCALargeText(ctx *router.Ctx) {
	cfg := c.services.Config.GetConfig()

	value := ctx.Arguments.Raw()

	cfg.Activity.LargeImageText = value
	_ = c.services.Config.SaveConfig(cfg)

	c.services.Embed.New(ctx).
		SetDescription("%s", "Custom Activity Large Image Text set to "+value).
		Send()
}

func (c *ActivityCommands) getCASmallImage(ctx *router.Ctx) {
	cfg := c.services.Config.GetConfig()

	c.services.Embed.New(ctx).
		AddField("Custom Activity Small Image", cfg.Activity.SmallImage).
		Send()
}

func (c *ActivityCommands) setCASmallImage(ctx *router.Ctx) {
	cfg := c.services.Config.GetConfig()

	value := ctx.Arguments.Raw()

	cfg.Activity.SmallImage = value
	_ = c.services.Config.SaveConfig(cfg)

	c.services.Embed.New(ctx).
		SetDescription("%s", "Custom Activity Small Image set to "+value).
		Send()
}

func (c *ActivityCommands) getCASmallText(ctx *router.Ctx) {
	cfg := c.services.Config.GetConfig()

	c.services.Embed.New(ctx).
		AddField("Custom Activity Small Image Text", cfg.Activity.SmallImageText).
		Send()
}

func (c *ActivityCommands) setCASmallText(ctx *router.Ctx) {
	cfg := c.services.Config.GetConfig()

	value := ctx.Arguments.Raw()

	cfg.Activity.SmallImageText = value
	_ = c.services.Config.SaveConfig(cfg)

	c.services.Embed.New(ctx).
		SetDescription("%s", "Custom Activity Small Image Text set to "+value).
		Send()
}

func (c *ActivityCommands) getCAButton1(ctx *router.Ctx) {
	cfg := c.services.Config.GetConfig()

	c.services.Embed.New(ctx).
		AddField("Custom Activity Button 1", fmt.Sprintf("%s (%s)", cfg.Activity.Button1Label, cfg.Activity.Button1URL)).
		Send()
}

func (c *ActivityCommands) setCAButton1(ctx *router.Ctx) {
	cfg := c.services.Config.GetConfig()

	label := ctx.Arguments.Get(0).Raw()
	url := ctx.Arguments.Get(1).Raw()

	cfg.Activity.Button1Label = label
	cfg.Activity.Button1URL = url
	_ = c.services.Config.SaveConfig(cfg)

	c.services.Embed.New(ctx).
		SetDescription("%s", fmt.Sprintf("CA Button 1 set to %s (%s)", label, url)).
		Send()
}

func (c *ActivityCommands) getCAButton2(ctx *router.Ctx) {
	cfg := c.services.Config.GetConfig()

	c.services.Embed.New(ctx).
		AddField("Custom Activity Button 2", fmt.Sprintf("%s (%s)", cfg.Activity.Button2Label, cfg.Activity.Button2URL)).
		Send()
}

func (c *ActivityCommands) setCAButton2(ctx *router.Ctx) {
	cfg := c.services.Config.GetConfig()

	label := ctx.Arguments.Get(0).Raw()
	url := ctx.Arguments.Get(1).Raw()

	cfg.Activity.Button2Label = label
	cfg.Activity.Button2URL = url
	_ = c.services.Config.SaveConfig(cfg)

	c.services.Embed.New(ctx).
		SetDescription("%s", fmt.Sprintf("CA Button 2 set to %s (%s)", label, url)).
		Send()
}

func (c *ActivityCommands) getCAPreset(ctx *router.Ctx) {
	cfg := c.services.Config.GetConfig()

	var currentPreset string

	for name, preset := range cfg.ActivityPresets {
		if reflect.DeepEqual(preset, cfg.Activity) {
			currentPreset = name
			break
		}
	}

	if currentPreset == "" {
		c.services.Embed.New(ctx).
			SetDescription("You must load a preset before using the get current preset command!").
			Send()
		return
	}

	c.services.Embed.New(ctx).
		SetDescription("Current Custom Activity preset is `%s`.", currentPreset).
		Send()
}

func (c *ActivityCommands) loadCAPreset(ctx *router.Ctx) {
	cfg := c.services.Config.GetConfig()

	name := ctx.Arguments.Raw()
	if name == "" {
		ctx.RespondText("You must provide a preset name!")
		return
	}

	preset, ok := cfg.ActivityPresets[name]
	if !ok {
		c.services.Embed.New(ctx).SetDescription("CA Preset `%s` not found.", name).Send()
		return
	}

	cfg.Activity = preset
	_ = c.services.Config.SaveConfig(cfg)

	c.services.Embed.New(ctx).
		SetDescription("Loaded Custom Activity preset `%s`.", name).
		Send()
}

func (c *ActivityCommands) nextCAPreset(ctx *router.Ctx) {
	cfg := c.services.Config.GetConfig()

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
	_ = c.services.Config.SaveConfig(cfg)

	c.services.Embed.New(ctx).
		SetDescription("Loaded next CA preset `%s`.", nextPresetName).
		Send()
}

func (c *ActivityCommands) saveCAPreset(ctx *router.Ctx) {
	cfg := c.services.Config.GetConfig()

	name := ctx.Arguments.Raw()
	if name == "" {
		ctx.RespondText("You must provide a name to save the preset as!")
		return
	}

	if cfg.ActivityPresets == nil {
		cfg.ActivityPresets = make(map[string]config.Activity)
	}

	cfg.ActivityPresets[name] = cfg.Activity
	_ = c.services.Config.SaveConfig(cfg)

	c.services.Embed.New(ctx).
		SetDescription("Saved current activity as preset `%s`.", name).
		Send()
}

func (c *ActivityCommands) removeCAPreset(ctx *router.Ctx) {
	cfg := c.services.Config.GetConfig()

	name := ctx.Arguments.Raw()

	if _, ok := cfg.ActivityPresets[name]; !ok {
		c.services.Embed.New(ctx).SetDescription("Preset `%s` not found.", name).Send()
		return
	}

	delete(cfg.ActivityPresets, name)
	_ = c.services.Config.SaveConfig(cfg)

	c.services.Embed.New(ctx).
		SetDescription("Deleted activity preset `%s`.", name).
		Send()
}

func (c *ActivityCommands) listCAPresets(ctx *router.Ctx) {
	cfg := c.services.Config.GetConfig()

	var presets []string

	for k, v := range cfg.ActivityPresets {
		presets = append(presets, fmt.Sprintf("%s: %s %s", k, v.Type, v.Name))
	}

	if len(presets) == 0 {
		c.services.Embed.New(ctx).SetDescription("No Custom Activity presets found.").Send()
		return
	}

	c.services.Embed.New(ctx).
		AddField("Custom Activity Presets", strings.Join(presets, "\n")).
		Send()
}

func (c *ActivityCommands) clearCAPresets(ctx *router.Ctx) {
	cfg := c.services.Config.GetConfig()

	cfg.ActivityPresets = make(map[string]config.Activity)
	_ = c.services.Config.SaveConfig(cfg)

	c.services.Embed.New(ctx).
		SetDescription("Cleared all activity presets.").
		Send()
}

func (c *ActivityCommands) setActivityPreset(ctx *router.Ctx) {
	cfg := c.services.Config.GetConfig()

	name := strings.ToLower(ctx.Arguments.Get(0).Raw())
	text := ctx.Arguments.RawFrom(1)

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
		ctx.RespondText("Invalid preset name. Use `youtube`, `twitch`, `discord`, `rickroll`, or `gta6`.")
		return
	}

	cfg.Activity = activity
	_ = c.services.Config.SaveConfig(cfg)

	c.services.Embed.New(ctx).
		SetDescription("Successfully set activity preset `%s`!", name).
		Send()
}

func (c *ActivityCommands) listActivityPresets(ctx *router.Ctx) {
	ctx.RespondText("Available activity presets: `youtube`, `twitch`, `discord`, `rickroll`, `gta6`.")
}

func (c *ActivityCommands) capresets(ctx *router.Ctx) {
	cfg := c.services.Config.GetConfig()

	action := ctx.Arguments.Get(0).Raw()
	name := ctx.Arguments.Get(1).Raw()

	switch action {
	case "get":
		var currentPreset string

		for name, preset := range cfg.ActivityPresets {
			if reflect.DeepEqual(preset, cfg.Activity) {
				currentPreset = name
				break
			}
		}

		if currentPreset == "" {
			c.services.Embed.New(ctx).
				SetDescription("You must load a preset before using the get current preset command!").
				Send()
			return
		}

		c.services.Embed.New(ctx).
			SetDescription("Current Custom Activity preset is `%s`.", currentPreset).
			Send()
	case "load":
		if name == "" {
			ctx.RespondText("You must provide a preset name!")
			return
		}

		preset, ok := cfg.ActivityPresets[name]
		if !ok {
			c.services.Embed.New(ctx).SetDescription("CA Preset `%s` not found.", name).Send()
			return
		}

		cfg.Activity = preset
		_ = c.services.Config.SaveConfig(cfg)

		c.services.Embed.New(ctx).
			SetDescription("Loaded Custom Activity preset `%s`.", name).
			Send()
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
		_ = c.services.Config.SaveConfig(cfg)

		c.services.Embed.New(ctx).
			SetDescription("Loaded next CA preset `%s`.", nextPresetName).
			Send()
	case "save":
		name := ctx.Arguments.Get(1).Raw()
		if name == "" {
			ctx.RespondText("You must provide a name to save the preset as!")
			return
		}

		if cfg.ActivityPresets == nil {
			cfg.ActivityPresets = make(map[string]config.Activity)
		}

		cfg.ActivityPresets[name] = cfg.Activity
		_ = c.services.Config.SaveConfig(cfg)

		c.services.Embed.New(ctx).
			SetDescription("Saved current activity as preset `%s`.", name).
			Send()
	case "rm":
		name := ctx.Arguments.Get(1).Raw()

		if _, ok := cfg.ActivityPresets[name]; !ok {
			c.services.Embed.New(ctx).SetDescription("Preset `%s` not found.", name).Send()
			return
		}

		delete(cfg.ActivityPresets, name)
		_ = c.services.Config.SaveConfig(cfg)

		c.services.Embed.New(ctx).
			SetDescription("Deleted activity preset `%s`.", name).
			Send()
	case "list":
		var presets []string

		for k, v := range cfg.ActivityPresets {
			presets = append(presets, fmt.Sprintf("%s: %s %s", k, v.Type, v.Name))
		}

		if len(presets) == 0 {
			c.services.Embed.New(ctx).SetDescription("No Custom Activity presets found.").Send()
			return
		}

		c.services.Embed.New(ctx).
			AddField("Custom Activity Presets", strings.Join(presets, "\n")).
			Send()
	case "clear":
		cfg.ActivityPresets = make(map[string]config.Activity)
		_ = c.services.Config.SaveConfig(cfg)

		c.services.Embed.New(ctx).
			SetDescription("Cleared all activity presets.").
			Send()
	default:
		ctx.RespondText("Invalid <get/load/next/save/rm/list/clear> argument. Use `get`, `load`, `next`, `save`, `rm`, `list` or `clear`.")
	}
}

func (c *ActivityCommands) cacreateimagekey(ctx *router.Ctx) {
	url := ctx.Arguments.Get(0).Raw()

	if strings.Contains(url, "discordapp.com") {
		ctx.RespondText("Your image seems to be hosted on Discord. Please use a different link to create an image key, as this might not work. You can try uploading it to a free image hoster like Imgur or ImgBB (but make sure to use the direct link to the image if you do. The URL should end with .png, .jpg or .gif).")
	}

	imageKey := api.LBGenImageKey(url)

	c.services.Embed.New(ctx).
		AddField("Image URL", url).
		AddField("Image Key", imageKey).
		Send()
}

func (c *ActivityCommands) cacreatespotifykey(ctx *router.Ctx) {
	bearerToken, err := c.services.Spotify.GetBearerToken(ctx.Session)
	if err != nil {
		ctx.RespondText("Could not find a Spotify connection. Please connect your Spotify account to your Discord.")
		return
	}

	resp, err := http.DefaultClient().Get(
		fmt.Sprintf("https://api.spotify.com/v1/search?q=%s&type=track", ctx.Arguments.Escaped()),
		http.WithHeader("Authorization", bearerToken),
	)
	json := resp.Json()

	song := json.Get("tracks", "items").Get(0)

	imageUrl := song.Get("album").Get("images").Get(0).Get("url").MustString()
	if err != nil || imageUrl == "" {
		ctx.RespondText("A Spotify error occurred: Song not found.")
		return
	}

	splitUrl := strings.Split(imageUrl, "/")

	c.services.Embed.New(ctx).
		AddField("Song", song.Get("name").MustString()).
		AddField("Image Key", "spotify:"+splitUrl[len(splitUrl)-1]).
		Send()
}
