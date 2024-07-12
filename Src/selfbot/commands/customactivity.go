package commands

import (
	"fmt"
	"reflect"
	"slices"
	"strings"

	"github.com/LightningDev1/LB-Selfbot/selfbot/api"
	"github.com/LightningDev1/LB-Selfbot/selfbot/config"
	"github.com/LightningDev1/LB-Selfbot/selfbot/embed"
	"github.com/LightningDev1/LB-Selfbot/selfbot/utils"
	"github.com/LightningDev1/dgc"
)

type CustomActivityCommands struct{}

func (c CustomActivityCommands) RegisterCommands(router *dgc.Router) {
	router.StartCategory("CustomActivity", "Commands to set a custom Discord activity")

	router.RegisterCmd(&dgc.Command{
		Name:        "cactivity",
		Description: "Get, set or toggle Custom Activity",
		Usage:       "[p]cactivity <get/set/toggle> [on/off]",
		Flags:       []string{CmdFlagDangerous},
		Handler:     c.cactivity,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "caconnectiontype",
		Description: "Get or set the Custom Activity connection type",
		Usage:       "[p]caconnectiontype <get/set/values> [value]",
		Flags:       []string{CmdFlagDangerous},
		Handler:     c.caconnectiontype,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "catype",
		Description: "Get or set the Custom Activity type",
		Usage:       "[p]catype <get/set/values> [value]",
		Flags:       []string{CmdFlagDangerous},
		Handler:     c.catype,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "caplatform",
		Description: "Get or set the Custom Activity platform",
		Usage:       "[p]caplatform <get/set/values> [value]",
		Flags:       []string{CmdFlagDangerous},
		Handler:     c.caplatform,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "caclientid",
		Description: "Get or set the Custom Activity Client ID",
		Usage:       "[p]caclientid <get/set> [value]",
		Flags:       []string{CmdFlagDangerous},
		Handler:     c.caclientid,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "caname",
		Description: "Get or set the Custom Activity Name",
		Usage:       "[p]caname <get/set> [value]",
		Flags:       []string{CmdFlagDangerous},
		Handler:     c.caname,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "castate",
		Description: "Get or set the Custom Activity State",
		Usage:       "[p]castate <get/set> [value]",
		Flags:       []string{CmdFlagDangerous},
		Handler:     c.castate,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "cadetails",
		Description: "Get or set the Custom Activity Details",
		Usage:       "[p]cadetails <get/set> [value]",
		Flags:       []string{CmdFlagDangerous},
		Handler:     c.cadetails,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "cainvitable",
		Description: "Get, set or toggle if RPC is invitable",
		Usage:       "[p]cainvitable <get/set/toggle> [on/off]",
		Flags:       []string{CmdFlagDangerous},
		Handler:     c.cainvitable,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "catime",
		Description: "Get, set or toggle if Custom Activity shows time",
		Usage:       "[p]catime <get/set/toggle> [on/off]",
		Flags:       []string{CmdFlagDangerous},
		Handler:     c.catime,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "castreamingurl",
		Description: "Get or set the URL for the Streaming activity",
		Usage:       "[p]castreamingurl <get/set> [value]",
		Flags:       []string{CmdFlagDangerous},
		Handler:     c.castreamingurl,
	})

	// TODO: Improve
	router.RegisterCmd(&dgc.Command{
		Name:        "cacreateimagekey",
		Description: "Create an image key for a Custom Activity",
		Usage:       "[p]cacreateimagekey <image url>",
		Handler:     c.cacreateimagekey,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "calargeimage",
		Description: "Get or set the Custom Activity large image key",
		Usage:       "[p]calargeimage <get/set> [value]",
		Flags:       []string{CmdFlagDangerous},
		Handler:     c.calargeimage,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "calargetext",
		Description: "Get or set the Custom Activity large image text",
		Usage:       "[p]calargetext <get/set> [value]",
		Flags:       []string{CmdFlagDangerous},
		Handler:     c.calargetext,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "casmallimage",
		Description: "Get or set the Custom Activity small image key",
		Usage:       "[p]casmallimage <get/set> [value]",
		Flags:       []string{CmdFlagDangerous},
		Handler:     c.casmallimage,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "casmalltext",
		Description: "Get or set the Custom Activity small image text",
		Usage:       "[p]casmalltext <get/set> [value]",
		Flags:       []string{CmdFlagDangerous},
		Handler:     c.casmalltext,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "cabutton1",
		Description: "Get or set the label and URL of the first button",
		Usage:       "[p]cabutton1 <get/set> [label] [url]",
		Flags:       []string{CmdFlagDangerous},
		Handler:     c.cabutton1,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "cabutton2",
		Description: "Get or set the label and URL of the second button",
		Usage:       "[p]cabutton2 <get/set> [label] [url]",
		Flags:       []string{CmdFlagDangerous},
		Handler:     c.cabutton2,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "capresets",
		Description: "Manage Custom Activity presets",
		Usage:       "[p]capresets <get/load/next/save/rm/list/clear> [name]",
		Flags:       []string{CmdFlagDangerous},
		Handler:     c.capresets,
	})
}

func (CustomActivityCommands) caenabled(ctx *dgc.Ctx) {
	toggle := strings.ToLower(utils.Arguments["<on/off>"])
	if toggle != "on" && toggle != "off" {
		_ = ctx.RespondText("Invalid value. Use `on` or `off`.")
		return
	}
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Failed to get config.")
		return
	}

	cfg.Activity.Enabled = toggle == "on"
	_ = cfg.Save()

	embed.NewEmbed().
		SetDescription(fmt.Sprintf("Turned Custom Activity %s.", toggle)).
		Send(ctx)
}

func (CustomActivityCommands) caconnectiontype(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Failed to get config.")
		return
	}

	toggle := strings.ToLower(utils.Arguments["<rpc/custom>"])
	if toggle == "rpc" {
		cfg.Activity.ActivityType = "RPC"
	} else if toggle == "custom" {
		cfg.Activity.ActivityType = "CustomActivity"
	} else {
		_ = ctx.RespondText("Invalid value. Use `rpc` or `custom`.")
		return
	}

	_ = cfg.Save()

	embed.NewEmbed().
		SetDescription(fmt.Sprintf("Set activity type to %s.", toggle)).
		Send(ctx)
}

func (CustomActivityCommands) caplaying(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Failed to get config.")
		return
	}

	cfg.Activity.Type = "Playing"
	_ = cfg.Save()

	embed.NewEmbed().
		SetDescription("Set custom activity type to Playing.").
		Send(ctx)
}

func (CustomActivityCommands) calistening(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Failed to get config.")
		return
	}

	cfg.Activity.Type = "Listening"
	_ = cfg.Save()

	embed.NewEmbed().
		SetDescription("Set custom activity type to Listening.").
		Send(ctx)
}

func (CustomActivityCommands) cawatching(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Failed to get config.")
		return
	}

	cfg.Activity.Type = "Watching"
	_ = cfg.Save()

	embed.NewEmbed().
		SetDescription("Set custom activity type to Watching.").
		Send(ctx)
}

func (CustomActivityCommands) castreaming(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Failed to get config.")
		return
	}

	cfg.Activity.Type = "Streaming"
	_ = cfg.Save()

	embed.NewEmbed().
		SetDescription("Set custom activity type to Streaming.").
		Send(ctx)
}

func (CustomActivityCommands) cacompeting(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Failed to get config.")
		return
	}

	cfg.Activity.Type = "Competing"
	_ = cfg.Save()

	embed.NewEmbed().
		SetDescription("Set custom activity type to Competing.").
		Send(ctx)
}

func (CustomActivityCommands) caclientid(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Failed to get config.")
		return
	}

	cfg.Activity.ClientID = utils.Arguments["<client id>"]
	_ = cfg.Save()

	embed.NewEmbed().
		SetDescription(fmt.Sprintf("Set activity client id to `%s`.", cfg.Activity.ClientID)).
		Send(ctx)
}

func (CustomActivityCommands) caname(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Failed to get config.")
		return
	}

	cfg.Activity.Name = ctx.Arguments.Raw()
	_ = cfg.Save()

	embed.NewEmbed().
		SetDescription(fmt.Sprintf("Set activity name to `%s`.", cfg.Activity.Name)).
		Send(ctx)
}

func (CustomActivityCommands) castate(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Failed to get config.")
		return
	}

	cfg.Activity.State = ctx.Arguments.Raw()
	_ = cfg.Save()

	embed.NewEmbed().
		SetDescription(fmt.Sprintf("Set activity state to `%s`.", cfg.Activity.State)).
		Send(ctx)
}

func (CustomActivityCommands) cadetails(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Failed to get config.")
		return
	}

	cfg.Activity.Details = ctx.Arguments.Raw()
	_ = cfg.Save()

	embed.NewEmbed().
		SetDescription(fmt.Sprintf("Set activity details to `%s`.", cfg.Activity.Details)).
		Send(ctx)
}

func (CustomActivityCommands) catime(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Failed to get config.")
		return
	}

	toggle := strings.ToLower(utils.Arguments["<on/off>"])
	if toggle != "on" && toggle != "off" {
		_ = ctx.RespondText("Invalid value. Use `on` or `off`.")
		return
	}

	cfg.Activity.Time = toggle == "on"
	_ = cfg.Save()

	embed.NewEmbed().
		SetDescription(fmt.Sprintf("Turned activity time %s.", toggle)).
		Send(ctx)
}

func (CustomActivityCommands) castreamingurl(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Failed to get config.")
		return
	}

	cfg.Activity.StreamingURL = ctx.Arguments.Raw()
	_ = cfg.Save()

	embed.NewEmbed().
		SetDescription(fmt.Sprintf("Set streaming URL to `%s`.", cfg.Activity.StreamingURL)).
		Send(ctx)
}

func (CustomActivityCommands) cacreateimagekey(ctx *dgc.Ctx) {
	url := utils.Arguments["<image url>"]

	if strings.Contains(url, "discordapp.com") {
		_ = ctx.RespondText("Your image seems to be hosted on Discord. Please use a different link to create an image key, as this might not work. You can try uploading it to a free image hoster like Imgur or ImgBB (but make sure to use the direct link to the image if you do. The URL should end with .png, .jpg or .gif).")
	}

	imageKey := api.GenImageKey(url)

	embed.NewEmbed().
		AddField("Image URL", utils.Arguments["<image url>"]).
		AddField("Image Key", imageKey).
		Send(ctx)
}

func (CustomActivityCommands) cainvitable(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Failed to get config.")
		return
	}

	toggle := strings.ToLower(utils.Arguments["<on/off>"])
	if toggle != "on" && toggle != "off" {
		_ = ctx.RespondText("Invalid value. Use `on` or `off`.")
		return
	}

	cfg.Activity.Invitable = toggle == "on"
	_ = cfg.Save()

	embed.NewEmbed().
		SetDescription(fmt.Sprintf("Turned activity invitable %s.", toggle)).
		Send(ctx)
}

func (CustomActivityCommands) calargeimage(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Failed to get config.")
		return
	}

	cfg.Activity.LargeImage = ctx.Arguments.Raw()
	_ = cfg.Save()

	embed.NewEmbed().
		SetDescription(fmt.Sprintf("Set large image key to `%s`.", cfg.Activity.LargeImage)).
		Send(ctx)
}

func (CustomActivityCommands) calargetext(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Failed to get config.")
		return
	}

	cfg.Activity.LargeImageText = ctx.Arguments.Raw()
	_ = cfg.Save()

	embed.NewEmbed().
		SetDescription(fmt.Sprintf("Set large image text to `%s`.", cfg.Activity.LargeImageText)).
		Send(ctx)
}

func (CustomActivityCommands) casmallimage(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Failed to get config.")
		return
	}

	cfg.Activity.SmallImage = ctx.Arguments.Raw()
	_ = cfg.Save()

	embed.NewEmbed().
		SetDescription(fmt.Sprintf("Set small image key to `%s`.", cfg.Activity.SmallImage)).
		Send(ctx)
}

func (CustomActivityCommands) casmalltext(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Failed to get config.")
		return
	}

	cfg.Activity.SmallImageText = ctx.Arguments.Raw()
	_ = cfg.Save()

	embed.NewEmbed().
		SetDescription(fmt.Sprintf("Set small image text to `%s`.", cfg.Activity.SmallImageText)).
		Send(ctx)
}

func (CustomActivityCommands) cabutton1text(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Failed to get config.")
		return
	}

	cfg.Activity.Button1Label = ctx.Arguments.Raw()
	_ = cfg.Save()

	embed.NewEmbed().
		SetDescription(fmt.Sprintf("Set button 1 text to `%s`.", cfg.Activity.Button1Label)).
		Send(ctx)
}

func (CustomActivityCommands) cabutton1url(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Failed to get config.")
		return
	}

	cfg.Activity.Button1URL = ctx.Arguments.Raw()
	_ = cfg.Save()

	embed.NewEmbed().
		SetDescription(fmt.Sprintf("Set button 1 url to `%s`.", cfg.Activity.Button1URL)).
		Send(ctx)
}

func (CustomActivityCommands) cabutton2text(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Failed to get config.")
		return
	}

	cfg.Activity.Button2Label = ctx.Arguments.Raw()
	_ = cfg.Save()

	embed.NewEmbed().
		SetDescription(fmt.Sprintf("Set button 2 text to `%s`.", cfg.Activity.Button2Label)).
		Send(ctx)
}

func (CustomActivityCommands) cabutton2url(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Failed to get config.")
		return
	}

	cfg.Activity.Button2URL = ctx.Arguments.Raw()
	_ = cfg.Save()

	embed.NewEmbed().
		SetDescription(fmt.Sprintf("Set button 2 url to `%s`.", cfg.Activity.Button2URL)).
		Send(ctx)
}

func (CustomActivityCommands) capresets(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Failed to get config.")
		return
	}

	var presets []string

	for k, v := range cfg.ActivityPresets {
		presets = append(presets, fmt.Sprintf("%s: %s %s", k, v.Type, v.Name))
	}

	if len(presets) == 0 {
		embed.NewEmbed().SetDescription("No Custom Activity presets found.").Send(ctx)
		return
	}

	embed.NewEmbed().
		SetDescription("Custom Activity presets\n\n" + strings.Join(presets, "\n")).
		Send(ctx)
}

func (CustomActivityCommands) clearcapresets(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Failed to get config.")
		return
	}

	cfg.ActivityPresets = make(map[string]config.Activity)
	_ = cfg.Save()

	embed.NewEmbed().
		SetDescription("Cleared all activity presets.").
		Send(ctx)
}

func (CustomActivityCommands) savecapreset(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Failed to get config.")
		return
	}

	name := utils.Arguments["<name>"]

	if cfg.ActivityPresets == nil {
		cfg.ActivityPresets = make(map[string]config.Activity)
	}

	cfg.ActivityPresets[name] = cfg.Activity
	_ = cfg.Save()

	embed.NewEmbed().
		SetDescription(fmt.Sprintf("Saved current activity as preset `%s`.", name)).
		Send(ctx)
}

func (CustomActivityCommands) loadcapreset(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Failed to get config.")
		return
	}

	name := utils.Arguments["<name>"]

	preset, ok := cfg.ActivityPresets[name]
	if !ok {
		embed.NewEmbed().SetDescription(fmt.Sprintf("Preset `%s` not found.", name)).Send(ctx)
		return
	}

	cfg.Activity = preset
	_ = cfg.Save()

	embed.NewEmbed().
		SetDescription(fmt.Sprintf("Loaded activity preset `%s`.", name)).
		Send(ctx)
}

func (CustomActivityCommands) deletecapreset(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Failed to get config.")
		return
	}

	name := utils.Arguments["<name>"]

	if _, ok := cfg.ActivityPresets[name]; !ok {
		embed.NewEmbed().SetDescription(fmt.Sprintf("Preset `%s` not found.", name)).Send(ctx)
		return
	}

	delete(cfg.ActivityPresets, name)
	_ = cfg.Save()

	embed.NewEmbed().
		SetDescription(fmt.Sprintf("Deleted activity preset `%s`.", name)).
		Send(ctx)
}

func (CustomActivityCommands) currentcapreset(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Failed to get config.")
		return
	}

	// Get current preset
	var currentPreset string

	for name, preset := range cfg.ActivityPresets {
		if reflect.DeepEqual(preset, cfg.Activity) {
			currentPreset = name
			break
		}
	}

	if currentPreset == "" {
		embed.NewEmbed().
			SetDescription("You must load a preset before using the current preset command!").
			Send(ctx)
		return
	}

	embed.NewEmbed().
		SetDescription(fmt.Sprintf("Current activity preset is `%s`.", currentPreset)).
		Send(ctx)
}

func (CustomActivityCommands) nextcapreset(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Failed to get config.")
		return
	}

	// Get current preset
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

	err = cfg.Save()
	if err != nil {
		embed.NewEmbed().
			SetDescription("Failed to load next activity preset.").
			Send(ctx)
		return
	}

	embed.NewEmbed().
		SetDescription(fmt.Sprintf("Loaded next activity preset `%s`.", nextPresetName)).
		Send(ctx)
}
