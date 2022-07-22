package commands

import (
	"time"

	"github.com/LightningDev1/LB-Selfbot/api"
	"github.com/LightningDev1/LB-Selfbot/utils"
	"github.com/LightningDev1/dgc"
	"github.com/LightningDev1/discordgo"
)

type ActivityCommands struct {
	CycleOn      bool
	CycleMode    string
	CycleText    []string
	StreamingURL string
	Session      *discordgo.Session
}

func (a *ActivityCommands) RegisterCommands(router *dgc.Router) {
	router.RegisterCmd(&dgc.Command{
		Name:        "online",
		Description: "Set your status to online",
		Usage:       "[p]online",
		Handler:     a.online,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "dnd",
		Description: "Set your status to do not disturb",
		Usage:       "[p]dnd",
		Handler:     a.dnd,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "idle",
		Description: "Set your status to idle",
		Usage:       "[p]idle",
		Handler:     a.idle,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "invisible",
		Description: "Set your status to invisible",
		Usage:       "[p]invisible",
		Handler:     a.invisible,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "customstatus",
		Description: "Change your custom status",
		Usage:       "[p]customstatus <text>",
		Handler:     a.customstatus,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "clearcustomstatus",
		Description: "Clear your custom status",
		Usage:       "[p]clearcustomstatus",
		Handler:     a.clearcustomstatus,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "aboutme",
		Description: "Change your About Me text",
		Usage:       "[p]aboutme <text>",
		Handler:     a.aboutme,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "clearaboutme",
		Description: "Clear your About Me text",
		Usage:       "[p]clearaboutme",
		Handler:     a.clearaboutme,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "listening",
		Description: "Set your activity to listening",
		Usage:       "[p]listening <text>",
		Handler:     a.listening,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "watching",
		Description: "Set your activity to watching",
		Usage:       "[p]watching <text>",
		Handler:     a.watching,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "playing",
		Description: "Set your activity to playing",
		Usage:       "[p]playing <text>",
		Handler:     a.playing,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "streaming",
		Description: "Set your activity to streaming",
		Usage:       "[p]streaming <url> <text>",
		Handler:     a.streaming,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "competing",
		Description: "Set your activity to competing",
		Usage:       "[p]competing <text>",
		Handler:     a.competing,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "clearactivity",
		Description: "Clear activity",
		Usage:       "[p]clearactivity",
		Handler:     a.clearactivity,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "rickrollactivity",
		Description: "Set a rick roll activity",
		Usage:       "[p]rickrollactivity",
		Handler:     a.rickrollactivity,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "youtubeactivity",
		Description: "Set a YouTube activity",
		Usage:       "[p]youtubeactivity <text>",
		Handler:     a.youtubeactivity,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "twitchactivity",
		Description: "Set a Twitch activity",
		Usage:       "[p]twitchactivity <text>",
		Handler:     a.twitchactivity,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "cyclestatus",
		Description: "Cycle between statuses",
		Usage:       "[p]cyclestatus",
		Handler:     a.cyclestatus,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "cyclecstatus",
		Description: "Cycle between custom statuses",
		Usage:       "[p]cyclecstatus <\"text1\"> <\"text2\"> [\"text3\"] [\"text4\"] [\"text5\"]",
		Handler:     a.cyclecstatus,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "cyclelistening",
		Description: "Cycle between listening activities",
		Usage:       "[p]cyclelistening <\"text1\"> <\"text2\"> [\"text3\"] [\"text4\"] [\"text5\"]",
		Handler:     a.cyclelistening,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "cyclewatching",
		Description: "Cycle between watching activities",
		Usage:       "[p]cyclewatching <\"text1\"> <\"text2\"> [\"text3\"] [\"text4\"] [\"text5\"]",
		Handler:     a.cyclewatching,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "cycleplaying",
		Description: "Cycle between playing activities",
		Usage:       "[p]cycleplaying <\"text1\"> <\"text2\"> [\"text3\"] [\"text4\"] [\"text5\"]",
		Handler:     a.cycleplaying,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "cyclestreaming",
		Description: "Cycle between streaming activities",
		Usage:       "[p]cyclestreaming <url> <\"text1\"> <\"text2\"> [\"text3\"] [\"text4\"] [\"text5\"]",
		Handler:     a.cyclestreaming,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "cyclecompeting",
		Description: "Cycle between competing activities",
		Usage:       "[p]cyclecompeting <\"text1\"> <\"text2\"> [\"text3\"] [\"text4\"] [\"text5\"]",
		Handler:     a.cyclecompeting,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "cycleactivity",
		Description: "Cycle between activity types",
		Usage:       "[p]cycleactivity <stream url> <text>",
		Handler:     a.cycleactivity,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "timestatus",
		Description: "Set current time as custom status",
		Usage:       "[p]timestatus",
		Handler:     a.timestatus,
	})

	statuses := []string{"online", "dnd", "idle", "invisible"}
	activityTypes := []string{"listening", "watching", "playing", "streaming", "competing"}

	go func() {
		current := 0
		for {
			time.Sleep(time.Second * 10)
			if !a.CycleOn {
				continue
			}
			current += 1
			if current >= 100 {
				current = 0
			}
			switch a.CycleMode {
			case "status":
				a.setActivity("status", statuses[current%len(statuses)])
			case "activitytype":
				a.setActivity(activityTypes[current%len(activityTypes)], a.CycleText[0])
			case "timestatus":
				a.setActivity("customstatus", "Current Time: "+time.Now().Format("15:04:05"))
			default:
				a.setActivity(a.CycleMode, a.CycleText[current%len(a.CycleText)])
			}
		}
	}()
}

func (a *ActivityCommands) setActivity(activityType, text string) {
	switch activityType {
	case "status":
		api.ChangeSettings(map[string]interface{}{"status": text}, true)
	case "customstatus":
		api.ChangeSettings(map[string]interface{}{"custom_status": map[string]string{"text": text}}, true)
	case "listening":
		a.Session.UpdateStatusComplex(discordgo.UpdateStatusData{
			Game: &discordgo.Activity{
				Name:  text,
				Type:  discordgo.ActivityTypeListening,
				Flags: 0,
			},
		})
	case "watching":
		a.Session.UpdateStatusComplex(discordgo.UpdateStatusData{
			Game: &discordgo.Activity{
				Name:  text,
				Type:  discordgo.ActivityTypeWatching,
				Flags: 0,
			},
		})
	case "playing":
		a.Session.UpdateStatusComplex(discordgo.UpdateStatusData{
			Game: &discordgo.Activity{
				Name:  text,
				Type:  discordgo.ActivityTypeGame,
				Flags: 0,
			},
		})
	case "streaming":
		a.Session.UpdateStatusComplex(discordgo.UpdateStatusData{
			Activities: []*discordgo.Activity{},
			Game: &discordgo.Activity{
				Name:    text,
				Details: text,
				URL:     a.StreamingURL,
				Type:    discordgo.ActivityTypeStreaming,
				Assets:  &discordgo.Assets{},
			},
		})
	case "competing":
		a.Session.UpdateStatusComplex(discordgo.UpdateStatusData{
			Game: &discordgo.Activity{
				Name:  text,
				Type:  discordgo.ActivityTypeCompeting,
				Flags: 0,
			},
		})
	}
}

func (a *ActivityCommands) online(ctx *dgc.Ctx) {
	a.Session = ctx.Session
	a.setActivity("status", "online")
}

func (a *ActivityCommands) dnd(ctx *dgc.Ctx) {
	a.Session = ctx.Session
	a.setActivity("status", "dnd")
}

func (a *ActivityCommands) idle(ctx *dgc.Ctx) {
	a.Session = ctx.Session
	a.setActivity("status", "idle")
}

func (a *ActivityCommands) invisible(ctx *dgc.Ctx) {
	a.Session = ctx.Session
	a.setActivity("status", "invisible")
}

func (a *ActivityCommands) customstatus(ctx *dgc.Ctx) {
	a.Session = ctx.Session
	a.setActivity("customstatus", ctx.Arguments.Raw())
}

func (a *ActivityCommands) clearcustomstatus(ctx *dgc.Ctx) {
	a.Session = ctx.Session
	a.setActivity("customstatus", "")
}

func (a *ActivityCommands) aboutme(ctx *dgc.Ctx) {
	api.ChangeSettings(map[string]interface{}{"bio": ctx.Arguments.Raw()}, false)
}

func (a *ActivityCommands) clearaboutme(ctx *dgc.Ctx) {
	api.ChangeSettings(map[string]interface{}{"bio": ""}, false)
}

func (a *ActivityCommands) listening(ctx *dgc.Ctx) {
	a.Session = ctx.Session
	a.setActivity("listening", ctx.Arguments.Raw())
}

func (a *ActivityCommands) watching(ctx *dgc.Ctx) {
	a.Session = ctx.Session
	a.setActivity("watching", ctx.Arguments.Raw())
}

func (a *ActivityCommands) playing(ctx *dgc.Ctx) {
	a.Session = ctx.Session
	a.setActivity("playing", ctx.Arguments.Raw())
}

func (a *ActivityCommands) streaming(ctx *dgc.Ctx) {
	a.Session = ctx.Session
	a.StreamingURL = utils.Arguments["<url>"]
	a.setActivity("streaming", utils.RemoveWords(ctx.Arguments.Raw(), 1))
}

func (a *ActivityCommands) competing(ctx *dgc.Ctx) {
	a.Session = ctx.Session
	a.setActivity("competing", ctx.Arguments.Raw())
}

func (a *ActivityCommands) clearactivity(ctx *dgc.Ctx) {
	a.CycleOn = false
	a.CycleMode = ""
	a.StreamingURL = ""
	a.CycleText = []string{}
	ctx.Session.UpdateStatusComplex(discordgo.UpdateStatusData{
		Game: &discordgo.Activity{
			Type: discordgo.ActivityType(4),
		},
		Activities: []*discordgo.Activity{},
	})
}

func (ActivityCommands) rickrollactivity(ctx *dgc.Ctx) {
	ctx.Session.UpdateStatusComplex(discordgo.UpdateStatusData{
		Activities: []*discordgo.Activity{},
		Game: &discordgo.Activity{
			Name:    "Never Gonna Give You Up",
			Details: "Never Gonna Give You Up",
			URL:     "https://www.youtube.com/watch?v=dQw4w9WgXcQ",
			Type:    discordgo.ActivityTypeListening,
			Assets: &discordgo.Assets{
				LargeImageID: "spotify:ab67616d0000b2735755e164993798e0c9ef7d7a",
				LargeText:    "Rick Astley",
			},
		},
	})
}

func (ActivityCommands) youtubeactivity(ctx *dgc.Ctx) {
	ctx.Session.UpdateStatusComplex(discordgo.UpdateStatusData{
		Game: &discordgo.Activity{
			Name:    "YouTube",
			Details: ctx.Arguments.Raw(),
			Type:    discordgo.ActivityTypeWatching,
			Flags:   0,
			Assets: &discordgo.Assets{
				LargeImageID: "spotify:ab67706c0000bebb7d1c495e63620a99f3cbcbf5",
				LargeText:    "YouTube",
			},
		},
	})
}

func (ActivityCommands) twitchactivity(ctx *dgc.Ctx) {
	ctx.Session.UpdateStatusComplex(discordgo.UpdateStatusData{
		Game: &discordgo.Activity{
			Name:    "Twitch",
			Details: ctx.Arguments.Raw(),
			Type:    discordgo.ActivityTypeWatching,
			Flags:   0,
			Assets: &discordgo.Assets{
				LargeImageID: "spotify:ab67706c0000bebb4081b043d43093915cb73114",
				LargeText:    "Twitch",
			},
		},
	})
}

func (a *ActivityCommands) cyclestatus(ctx *dgc.Ctx) {
	a.Session = ctx.Session
	a.CycleOn = true
	a.CycleMode = "status"
}

func (a *ActivityCommands) setCycleText() {
	a.CycleText = []string{utils.Arguments["<\"text1\">"], utils.Arguments["<\"text2\">"]}
	if utils.Arguments["[\"text3\"]"] != "" {
		a.CycleText = append(a.CycleText, utils.Arguments["[\"text3\"]"])
	}
	if utils.Arguments["[\"text4\"]"] != "" {
		a.CycleText = append(a.CycleText, utils.Arguments["[\"text4\"]"])
	}
	if utils.Arguments["[\"text5\"]"] != "" {
		a.CycleText = append(a.CycleText, utils.Arguments["[\"text5\"]"])
	}
}

func (a *ActivityCommands) cyclecstatus(ctx *dgc.Ctx) {
	a.Session = ctx.Session
	a.CycleOn = true
	a.CycleMode = "customstatus"
	a.setCycleText()
}

func (a *ActivityCommands) cyclelistening(ctx *dgc.Ctx) {
	a.Session = ctx.Session
	a.CycleOn = true
	a.CycleMode = "listening"
	a.setCycleText()
}

func (a *ActivityCommands) cyclewatching(ctx *dgc.Ctx) {
	a.Session = ctx.Session
	a.CycleOn = true
	a.CycleMode = "watching"
	a.setCycleText()
}

func (a *ActivityCommands) cycleplaying(ctx *dgc.Ctx) {
	a.Session = ctx.Session
	a.CycleOn = true
	a.CycleMode = "playing"
	a.setCycleText()
}

func (a *ActivityCommands) cyclestreaming(ctx *dgc.Ctx) {
	a.Session = ctx.Session
	a.CycleOn = true
	a.CycleMode = "streaming"
	a.StreamingURL = utils.Arguments["<url>"]
	a.setCycleText()
}

func (a *ActivityCommands) cyclecompeting(ctx *dgc.Ctx) {
	a.Session = ctx.Session
	a.CycleOn = true
	a.CycleMode = "competing"
	a.setCycleText()
}

func (a *ActivityCommands) cycleactivity(ctx *dgc.Ctx) {
	a.Session = ctx.Session
	a.CycleOn = true
	a.CycleMode = "activitytype"
	a.StreamingURL = utils.Arguments["<stream url>"]
	a.CycleText = []string{utils.RemoveWords(ctx.Arguments.Raw(), 1)}
}

func (a *ActivityCommands) timestatus(ctx *dgc.Ctx) {
	a.Session = ctx.Session
	a.CycleOn = true
	a.CycleMode = "timestatus"
}
