package activity

import (
	"fmt"
	"strings"

	"github.com/LightningDev1/LB-Selfbot/internal/services"
	"github.com/LightningDev1/dgc"
)

type ActivityCommands struct {
	services *services.Services
}

func NewActivityCommands(services *services.Services) *ActivityCommands {
	return &ActivityCommands{services: services}
}

func (c *ActivityCommands) Register() {
	router := c.services.Discord.GetRouter()
	router.StartCategory("Activity", "Customize your Discord activity")

	router.RegisterCmd(&dgc.Command{
		Name:        "status",
		Description: "Set your Discord status",
		Usage:       "[p]status <online/dnd/idle/invisible>",
		Flags:       []string{CmdFlagDangerous},
		Handler:     c.status,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "customstatus",
		Description: "Change your custom status",
		Usage:       "[p]customstatus [text]",
		Flags:       []string{CmdFlagDangerous},
		Handler:     c.customstatus,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "clearactivity",
		Description: "Clear activity",
		Usage:       "[p]clearactivity",
		Flags:       []string{CmdFlagDangerous},
		Handler:     c.clearactivity,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "copyactivity",
		Description: "Copy a user's activity",
		Usage:       "[p]copyactivity <user> [save to config: on/off]",
		Flags:       []string{CmdFlagDangerous},
		Handler:     c.copyactivity,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "getactivityjson",
		Description: "Get a users activity in JSON format",
		Usage:       "[p]getactivityjson <user>",
		Handler:     c.getactivityjson,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "setactivityjson",
		Description: "Set your activity using raw JSON",
		Usage:       "[p]setactivityjson <save to config: on/off> <json>",
		Flags:       []string{CmdFlagDangerous},
		Handler:     c.setactivityjson,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "cactivity",
		Description: "Get, set or toggle Custom Activity",
		Usage:       "[p]cactivity <get/set/toggle> [on/off]",
		Flags:       []string{CmdFlagDangerous},
		Handler:     c.cactivity,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "capauseforother",
		Description: "Get, set or toggle Pause CA for other Activity",
		Usage:       "[p]capauseforother <get/set/toggle> [on/off]",
		Flags:       []string{CmdFlagDangerous},
		Handler:     c.capauseforother,
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

	router.RegisterCmd(&dgc.Command{
		Name:        "cacreateimagekey",
		Description: "Create an image key for Custom Activity",
		Usage:       "[p]cacreateimagekey <image url>",
		Handler:     c.cacreateimagekey,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "cacreatespotifykey",
		Description: "Create an Spotify image key for Custom Activity",
		Usage:       "[p]cacreatespotifykey <song name>",
		Handler:     c.cacreatespotifykey,
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
		Usage:       "[p]cabutton1 <get/set> [\"label\"] [url]",
		Flags:       []string{CmdFlagDangerous},
		Handler:     c.cabutton1,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "cabutton2",
		Description: "Get or set the label and URL of the second button",
		Usage:       "[p]cabutton2 <get/set> [\"label\"] [url]",
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

	router.RegisterCmd(&dgc.Command{
		Name:        "activitypresets",
		Description: "Set or view built-in Custom Activity presets",
		Usage:       "[p]activitypresets <set/values> [name] [text]",
		Flags:       []string{CmdFlagDangerous},
		Handler:     c.activitypresets,
	})
}

// Keep original command handlers and helper functions
func (a *ActivityCommands) status(ctx *dgc.Ctx) {
	status := strings.ToLower(ctx.Arguments[0])

	switch status {
	case "online", "dnd", "idle", "invisible":
		err := a.setStatus(ctx.Session, status)
		if err != nil {
			ctx.RespondText("An error occurred while changing your status: " + err.Error())
			return
		}

		ctx.RespondEmbed(fmt.Sprintf("Successfully set status to `%s`!", status))
	default:
		ctx.RespondText("Invalid status argument. Use `online`, `dnd`, `idle`, or `invisible`.")
	}
}

func (a *ActivityCommands) customstatus(ctx *dgc.Ctx) {
	// ... Existing implementation ...
}

// ... Keep all other existing command implementations as is ...
