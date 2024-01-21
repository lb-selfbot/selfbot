package commands

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/LightningDev1/LB-Selfbot/selfbot/api"
	"github.com/LightningDev1/LB-Selfbot/selfbot/config"
	"github.com/LightningDev1/LB-Selfbot/selfbot/embed"
	"github.com/LightningDev1/LB-Selfbot/selfbot/sessions"
	"github.com/LightningDev1/LB-Selfbot/selfbot/utils"
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

func (a *ActivityCommands) RegisterCommands(router *dgc.Router, session *discordgo.Session) {
	a.Session = session
	router.StartCategory("Activity", "Set your Discord activity")
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
		Usage:       "[p]customstatus [text]",
		Handler:     a.customstatus,
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
		Name:        "discordactivity",
		Description: "Set a Discord activity",
		Usage:       "[p]discordactivity <text>",
		Handler:     a.discordactivity,
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
	router.RegisterCmd(&dgc.Command{
		Name:        "getactivityjson",
		Description: "Get a users activity in JSON format",
		Usage:       "[p]getactivityjson <user>",
		Handler:     a.getactivityjson,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "setactivityjson",
		Description: "Set your activity using raw JSON",
		Usage:       "[p]setactivityjson <json>",
		Handler:     a.setactivityjson,
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
				a.setActivity(a.Session, "status", statuses[current%len(statuses)])
			case "activitytype":
				a.setActivity(a.Session, activityTypes[current%len(activityTypes)], a.CycleText[0])
			case "timestatus":
				a.setActivity(a.Session, "customstatus", "Current Time: "+time.Now().Format("15:04:05"))
			default:
				a.setActivity(a.Session, a.CycleMode, a.CycleText[current%len(a.CycleText)])
			}
		}
	}()
}

func (a *ActivityCommands) setActivity(session *discordgo.Session, activityType, text string) {
	cfg, err := config.GetConfig()
	if err != nil {
		utils.Logging.Error("Error getting config: " + err.Error())
		return
	}

	switch activityType {
	case "status":
		sessions.LastStatusChange = time.Now()
		sessions.Status = text
		
		_ = api.ChangeSettings(session, "/settings", map[string]any{"status": text})
		_ = a.Session.UpdateStatusComplex(discordgo.UpdateStatusData{
			AFK:    true,
			Status: text,
		})
	case "customstatus":
		data := map[string]any{"custom_status": map[string]any{"text": text}}
		if text == "" {
			data = map[string]any{"custom_status": nil}
		}

		_ = api.ChangeSettings(session, "/settings", data)
	case "listening":
		_ = a.Session.UpdateStatusComplex(discordgo.UpdateStatusData{
			Activities: []*discordgo.Activity{{
				Name: text,
				Type: discordgo.ActivityTypeListening,
			}},
			AFK:    true,
			Status: sessions.Status,
		})
	case "watching":
		_ = a.Session.UpdateStatusComplex(discordgo.UpdateStatusData{
			Activities: []*discordgo.Activity{{
				Name: text,
				Type: discordgo.ActivityTypeWatching,
			}},
			AFK:    true,
			Status: sessions.Status,
		})
	case "playing":
		_ = a.Session.UpdateStatusComplex(discordgo.UpdateStatusData{
			Activities: []*discordgo.Activity{{
				Name: text,
				Type: discordgo.ActivityTypeGame,
			}},
			AFK:    true,
			Status: sessions.Status,
		})
	case "streaming":
		_ = a.Session.UpdateStatusComplex(discordgo.UpdateStatusData{
			Activities: []*discordgo.Activity{{
				Name:    text,
				Details: text,
				URL:     utils.Misc.Or(a.StreamingURL, cfg.Activity.StreamingURL),
				Type:    discordgo.ActivityTypeStreaming,
				Assets:  &discordgo.Assets{},
			}},
			AFK:    true,
			Status: sessions.Status,
		})
	case "competing":
		_ = a.Session.UpdateStatusComplex(discordgo.UpdateStatusData{
			Activities: []*discordgo.Activity{{
				Name: text,
				Type: discordgo.ActivityTypeCompeting,
			}},
			AFK:    true,
			Status: sessions.Status,
		})
	}
}

func (a *ActivityCommands) online(ctx *dgc.Ctx) {
	a.setActivity(ctx.Session, "status", "online")
}

func (a *ActivityCommands) dnd(ctx *dgc.Ctx) {
	a.setActivity(ctx.Session, "status", "dnd")
}

func (a *ActivityCommands) idle(ctx *dgc.Ctx) {
	a.setActivity(ctx.Session, "status", "idle")
}

func (a *ActivityCommands) invisible(ctx *dgc.Ctx) {
	a.setActivity(ctx.Session, "status", "invisible")
}

func (a *ActivityCommands) customstatus(ctx *dgc.Ctx) {
	a.setActivity(ctx.Session, "customstatus", ctx.Arguments.Raw())
}

func (a *ActivityCommands) listening(ctx *dgc.Ctx) {
	a.setActivity(ctx.Session, "listening", ctx.Arguments.Raw())
}

func (a *ActivityCommands) watching(ctx *dgc.Ctx) {
	a.setActivity(ctx.Session, "watching", ctx.Arguments.Raw())
}

func (a *ActivityCommands) playing(ctx *dgc.Ctx) {
	a.setActivity(ctx.Session, "playing", ctx.Arguments.Raw())
}

func (a *ActivityCommands) streaming(ctx *dgc.Ctx) {
	a.StreamingURL = utils.Arguments["<url>"]
	a.setActivity(ctx.Session, "streaming", utils.String.RemoveWords(ctx.Arguments.Raw(), 1))
}

func (a *ActivityCommands) competing(ctx *dgc.Ctx) {
	a.setActivity(ctx.Session, "competing", ctx.Arguments.Raw())
}

func (a *ActivityCommands) clearactivity(ctx *dgc.Ctx) {
	a.CycleOn = false
	a.CycleMode = ""
	a.StreamingURL = ""
	a.CycleText = []string{}

	_ = ctx.Session.UpdateStatusComplex(discordgo.UpdateStatusData{
		Game:       nil,
		Activities: []*discordgo.Activity{},
		AFK:        true,
		Status:     sessions.Status,
	})
}

func (ActivityCommands) rickrollactivity(ctx *dgc.Ctx) {
	_ = ctx.Session.UpdateStatusComplex(discordgo.UpdateStatusData{
		Activities: []*discordgo.Activity{{
			Name:    "Never Gonna Give You Up",
			Details: "Never Gonna Give You Up",
			URL:     "https://www.youtube.com/watch?v=dQw4w9WgXcQ",
			Type:    discordgo.ActivityTypeListening,
			Assets: &discordgo.Assets{
				LargeImageID: "spotify:ab67616d0000b2735755e164993798e0c9ef7d7a",
				LargeText:    "Rick Astley",
			},
		}},
		AFK:    true,
		Status: sessions.Status,
	})
}

func (ActivityCommands) youtubeactivity(ctx *dgc.Ctx) {
	_ = ctx.Session.UpdateStatusComplex(discordgo.UpdateStatusData{
		Activities: []*discordgo.Activity{{
			Name:    "YouTube",
			Details: ctx.Arguments.Raw(),
			Type:    discordgo.ActivityTypeWatching,
			Flags:   0,
			Assets: &discordgo.Assets{
				LargeImageID: "spotify:ab67706c0000bebb7d1c495e63620a99f3cbcbf5",
				LargeText:    "YouTube",
			},
		}},
		AFK:    true,
		Status: sessions.Status,
	})
}

func (ActivityCommands) twitchactivity(ctx *dgc.Ctx) {
	_ = ctx.Session.UpdateStatusComplex(discordgo.UpdateStatusData{
		Activities: []*discordgo.Activity{{
			Name:    "Twitch",
			Details: ctx.Arguments.Raw(),
			Type:    discordgo.ActivityTypeWatching,
			Flags:   0,
			Assets: &discordgo.Assets{
				LargeImageID: "spotify:ab67706c0000bebb4081b043d43093915cb73114",
				LargeText:    "Twitch",
			},
		}},
		AFK:    true,
		Status: sessions.Status,
	})
}

func (ActivityCommands) discordactivity(ctx *dgc.Ctx) {
	_ = ctx.Session.UpdateStatusComplex(discordgo.UpdateStatusData{
		Activities: []*discordgo.Activity{{
			Name:    "Discord",
			Details: "Discord",
			State:   ctx.Arguments.Raw(),
			Type:    discordgo.ActivityTypeWatching,
			Flags:   0,
			Assets: &discordgo.Assets{
				LargeImageID: "spotify:ab67616d0000b273db20400815ca2a9481bd255e",
			},
		}},
		AFK:    true,
		Status: sessions.Status,
	})
}

func (a *ActivityCommands) cyclestatus(ctx *dgc.Ctx) {
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
	a.CycleOn = true
	a.CycleMode = "customstatus"
	a.setCycleText()
}

func (a *ActivityCommands) cyclelistening(ctx *dgc.Ctx) {
	a.CycleOn = true
	a.CycleMode = "listening"
	a.setCycleText()
}

func (a *ActivityCommands) cyclewatching(ctx *dgc.Ctx) {
	a.CycleOn = true
	a.CycleMode = "watching"
	a.setCycleText()
}

func (a *ActivityCommands) cycleplaying(ctx *dgc.Ctx) {
	a.CycleOn = true
	a.CycleMode = "playing"
	a.setCycleText()
}

func (a *ActivityCommands) cyclestreaming(ctx *dgc.Ctx) {
	a.CycleOn = true
	a.CycleMode = "streaming"
	a.StreamingURL = utils.Arguments["<url>"]
	a.setCycleText()
}

func (a *ActivityCommands) cyclecompeting(ctx *dgc.Ctx) {
	a.CycleOn = true
	a.CycleMode = "competing"
	a.setCycleText()
}

func (a *ActivityCommands) cycleactivity(ctx *dgc.Ctx) {
	a.CycleOn = true
	a.CycleMode = "activitytype"
	a.StreamingURL = utils.Arguments["<stream url>"]
	a.CycleText = []string{utils.String.RemoveWords(ctx.Arguments.Raw(), 1)}
}

func (a *ActivityCommands) timestatus(ctx *dgc.Ctx) {
	a.CycleOn = true
	a.CycleMode = "timestatus"
}

func (a *ActivityCommands) getactivityjson(ctx *dgc.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		_ = ctx.RespondText("Could not find user")
		return
	}
	mutualServers, err := ctx.Session.UserProfile(user.ID)
	if err != nil {
		_ = ctx.RespondText("Could not find user")
		return
	}
	if len(mutualServers.MutualGuilds) < 1 {
		_ = ctx.RespondText("User is not in any mutual servers")
		return
	}
	emb := embed.NewEmbed()
	for _, guild_ := range mutualServers.MutualGuilds {
		guild, err := utils.Discord.GetGuild(ctx.Session, guild_.ID)
		if err != nil {
			continue
		}
		member, ok := utils.Discord.GetMember(ctx.Session, guild.ID, user.ID)
		if !ok {
			continue
		}
		if member.Presence == nil {
			continue
		}
		for index, activity := range member.Presence.Activities {
			jsonActivity, _ := json.MarshalIndent(activity, "", "    ")
			emb.AddField(fmt.Sprintf("Activity %d", index+1), string(jsonActivity))
		}
		break
	}
	emb.Send(ctx)
}

func (a *ActivityCommands) setactivityjson(ctx *dgc.Ctx) {
	rawJson := ctx.Arguments.Raw()
	var activity discordgo.Activity
	err := json.Unmarshal([]byte(rawJson), &activity)
	if err != nil {
		_ = ctx.RespondText("Invalid activity JSON! " + err.Error())
		return
	}
	_ = ctx.Session.UpdateStatusComplex(discordgo.UpdateStatusData{
		Activities: []*discordgo.Activity{&activity},
		AFK:        true,
		Status:     sessions.Status,
	})
}
