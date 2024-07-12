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

func (r CustomActivityCommands) RegisterCommands(router *dgc.Router) {
	router.StartCategory("CustomActivity", "Commands to set a custom Discord activity")
	router.RegisterCmd(&dgc.Command{
		Name:        "caenabled",
		Description: "Turn Custom Activity on or off",
		Usage:       "[p]caenabled <on/off>",
		Handler:     r.caenabled,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "caconnectiontype",
		Description: "Set Custom Activity connection type",
		Usage:       "[p]caconnectiontype <rpc/custom>",
		Handler:     r.caconnectiontype,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "catype",
		Description: "Set custom activity type",
		Usage:       "[p]cactivitytype <playing/listening/watching/streaming/competing/spotify>",
		Handler:     r.caconnectiontype,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "caplaying",
		Description: "Set Custom Activity to playing",
		Usage:       "[p]caplaying",
		Handler:     r.caplaying,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "calistening",
		Description: "Set Custom Activity to listening",
		Usage:       "[p]calistening",
		Handler:     r.calistening,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "cawatching",
		Description: "Set Custom Activity to watching",
		Usage:       "[p]cawatching",
		Handler:     r.cawatching,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "castreaming",
		Description: "Set Custom Activity to streaming",
		Usage:       "[p]castreaming",
		Handler:     r.castreaming,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "cacompeting",
		Description: "Set Custom Activity to competing",
		Usage:       "[p]cacompeting",
		Handler:     r.cacompeting,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "caclientid",
		Description: "Set the RPC client id",
		Usage:       "[p]caclientid <client id>",
		Handler:     r.caclientid,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "caname",
		Description: "Set the Custom Activity name",
		Usage:       "[p]caname [text]",
		Handler:     r.caname,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "castate",
		Description: "Set the Custom Activity state",
		Usage:       "[p]castate [text]",
		Handler:     r.castate,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "cadetails",
		Description: "Set the Custom Activity details",
		Usage:       "[p]cadetails [text]",
		Handler:     r.cadetails,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "cainvitable",
		Description: "Set if Custom Activity is invitable",
		Usage:       "[p]cainvitable <on/off>",
		Handler:     r.cainvitable,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "catime",
		Description: "Turn Custom Activity time on or off",
		Usage:       "[p]catime <on/off>",
		Handler:     r.catime,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "castreamingurl",
		Description: "Set the URL for the streaming activity",
		Usage:       "[p]castreamingurl [url]",
		Handler:     r.castreamingurl,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "cacreateimagekey",
		Description: "Create an image key for a Custom Activity",
		Usage:       "[p]cacreateimagekey <image url>",
		Handler:     r.cacreateimagekey,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "calargeimage",
		Description: "Set the Custom Activity large image key",
		Usage:       "[p]calargeimage [image]",
		Handler:     r.calargeimage,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "calargetext",
		Description: "Set the Custom Activity large image text",
		Usage:       "[p]calargetext [text]",
		Handler:     r.calargetext,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "casmallimage",
		Description: "Set the Custom Activity small image",
		Usage:       "[p]casmallimage [image]",
		Handler:     r.casmallimage,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "casmalltext",
		Description: "Set the Custom Activity small image text",
		Usage:       "[p]casmalltext [text]",
		Handler:     r.casmalltext,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "cabutton1text",
		Description: "Set the RPC first button text",
		Usage:       "[p]cabutton1text [text]",
		Handler:     r.cabutton1text,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "cabutton1url",
		Description: "Set the RPC first button url",
		Usage:       "[p]cabutton1url [url]",
		Handler:     r.cabutton1url,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "cabutton2text",
		Description: "Set the RPC second button text",
		Usage:       "[p]cabutton2text [text]",
		Handler:     r.cabutton2text,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "cabutton2url",
		Description: "Set the RPC second button url",
		Usage:       "[p]cabutton2url [url]",
		Handler:     r.cabutton2url,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "capresets",
		Description: "List all Custom Activity presets",
		Usage:       "[p]capresets",
		Handler:     r.capresets,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "clearcapresets",
		Description: "Clear all Custom Activity presets",
		Usage:       "[p]clearcapresets",
		Handler:     r.clearcapresets,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "savecapreset",
		Description: "Save current activity as a preset",
		Usage:       "[p]savecapreset <name>",
		Handler:     r.savecapreset,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "loadcapreset",
		Description: "Load a Custom Activity preset",
		Usage:       "[p]loadcapreset <name>",
		Handler:     r.loadcapreset,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "deletecapreset",
		Description: "Delete a Custom Activity preset",
		Usage:       "[p]deletecapreset <name>",
		Handler:     r.deletecapreset,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "currentcapreset",
		Description: "Get the name of the current Custom Activity preset",
		Usage:       "[p]currentcapreset",
		Handler:     r.currentcapreset,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "nextcapreset",
		Description: "Load the next Custom Activity preset",
		Usage:       "[p]nextcapreset",
		Handler:     r.nextcapreset,
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
