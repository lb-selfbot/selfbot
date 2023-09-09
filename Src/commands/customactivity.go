package commands

import (
	"strings"

	"github.com/LightningDev1/LB-Selfbot/api"
	"github.com/LightningDev1/LB-Selfbot/config"
	"github.com/LightningDev1/LB-Selfbot/embed"
	"github.com/LightningDev1/LB-Selfbot/utils"
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
		Name:        "cactivitytype",
		Description: "Set custom activity type",
		Usage:       "[p]cactivitytype <rpc/custom>",
		Handler:     r.cactivitytype,
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
		Description: "Set Custom Activity time on or off",
		Usage:       "[p]catime <on/off>",
		Handler:     r.catime,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "cacreateimagekey",
		Description: "Create an image key for Custom Activity",
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
}

func (CustomActivityCommands) caenabled(ctx *dgc.Ctx) {
	toggle := strings.ToLower(utils.Arguments["<on/off>"])
	if toggle != "on" && toggle != "off" {
		_ = ctx.RespondText("Invalid value. Use `on` or `off`.")
		return
	}
	cfg, _ := config.GetConfig()
	cfg.Activity.Enabled = toggle == "on"
	_ = cfg.Save()
}

func (CustomActivityCommands) cactivitytype(ctx *dgc.Ctx) {
	cfg, _ := config.GetConfig()

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
}

func (CustomActivityCommands) caplaying(*dgc.Ctx) {
	cfg, _ := config.GetConfig()
	cfg.Activity.Type = "Playing"
	_ = cfg.Save()
}

func (CustomActivityCommands) calistening(*dgc.Ctx) {
	cfg, _ := config.GetConfig()
	cfg.Activity.Type = "Listening"
	_ = cfg.Save()
}

func (CustomActivityCommands) cawatching(*dgc.Ctx) {
	cfg, _ := config.GetConfig()
	cfg.Activity.Type = "Watching"
	_ = cfg.Save()
}

func (CustomActivityCommands) castreaming(*dgc.Ctx) {
	cfg, _ := config.GetConfig()
	cfg.Activity.Type = "Streaming"
	_ = cfg.Save()
}

func (CustomActivityCommands) cacompeting(*dgc.Ctx) {
	cfg, _ := config.GetConfig()
	cfg.Activity.Type = "Competing"
	_ = cfg.Save()
}

func (CustomActivityCommands) caclientid(*dgc.Ctx) {
	cfg, _ := config.GetConfig()
	cfg.Activity.ClientID = utils.Arguments["<client id>"]
	_ = cfg.Save()
}

func (CustomActivityCommands) caname(ctx *dgc.Ctx) {
	cfg, _ := config.GetConfig()
	cfg.Activity.Name = ctx.Arguments.Raw()
	_ = cfg.Save()
}

func (CustomActivityCommands) castate(ctx *dgc.Ctx) {
	cfg, _ := config.GetConfig()
	cfg.Activity.State = ctx.Arguments.Raw()
	_ = cfg.Save()
}

func (CustomActivityCommands) cadetails(ctx *dgc.Ctx) {
	cfg, _ := config.GetConfig()
	cfg.Activity.Details = ctx.Arguments.Raw()
	_ = cfg.Save()
}

func (CustomActivityCommands) catime(ctx *dgc.Ctx) {
	toggle := strings.ToLower(utils.Arguments["<on/off>"])
	if toggle != "on" && toggle != "off" {
		_ = ctx.RespondText("Invalid value. Use `on` or `off`.")
		return
	}
	cfg, _ := config.GetConfig()
	cfg.Activity.Time = toggle == "on"
	_ = cfg.Save()
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
	toggle := strings.ToLower(utils.Arguments["<on/off>"])
	if toggle != "on" && toggle != "off" {
		_ = ctx.RespondText("Invalid value. Use `on` or `off`.")
		return
	}
	cfg, _ := config.GetConfig()
	cfg.Activity.Invitable = toggle == "on"
	_ = cfg.Save()
}

func (CustomActivityCommands) calargeimage(ctx *dgc.Ctx) {
	cfg, _ := config.GetConfig()
	cfg.Activity.LargeImage = ctx.Arguments.Raw()
	_ = cfg.Save()
}

func (CustomActivityCommands) calargetext(ctx *dgc.Ctx) {
	cfg, _ := config.GetConfig()
	cfg.Activity.LargeImageText = ctx.Arguments.Raw()
	_ = cfg.Save()
}

func (CustomActivityCommands) casmallimage(ctx *dgc.Ctx) {
	cfg, _ := config.GetConfig()
	cfg.Activity.SmallImage = ctx.Arguments.Raw()
	_ = cfg.Save()
}

func (CustomActivityCommands) casmalltext(ctx *dgc.Ctx) {
	cfg, _ := config.GetConfig()
	cfg.Activity.SmallImageText = ctx.Arguments.Raw()
	_ = cfg.Save()
}

func (CustomActivityCommands) cabutton1text(ctx *dgc.Ctx) {
	cfg, _ := config.GetConfig()
	cfg.Activity.Button1Label = ctx.Arguments.Raw()
	_ = cfg.Save()
}

func (CustomActivityCommands) cabutton1url(ctx *dgc.Ctx) {
	cfg, _ := config.GetConfig()
	cfg.Activity.Button1URL = ctx.Arguments.Raw()
	_ = cfg.Save()
}

func (CustomActivityCommands) cabutton2text(ctx *dgc.Ctx) {
	cfg, _ := config.GetConfig()
	cfg.Activity.Button2Label = ctx.Arguments.Raw()
	_ = cfg.Save()
}

func (CustomActivityCommands) cabutton2url(ctx *dgc.Ctx) {
	cfg, _ := config.GetConfig()
	cfg.Activity.Button2URL = ctx.Arguments.Raw()
	_ = cfg.Save()
}
