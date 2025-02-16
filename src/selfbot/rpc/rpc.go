package rpc

import (
	"strings"
	"time"

	"github.com/LightningDev1/LB-Selfbot/selfbot/config"
	"github.com/LightningDev1/LB-Selfbot/selfbot/constants"
	"github.com/LightningDev1/LB-Selfbot/selfbot/logging"
	"github.com/LightningDev1/LB-Selfbot/selfbot/sessions"
	"github.com/LightningDev1/LB-Selfbot/selfbot/utils"
	"github.com/LightningDev1/LB-Selfbot/selfbot/variables"
	"github.com/LightningDev1/discordgo"
	"github.com/hugolgst/rich-go/client"
)

func getActivityType(activityType string) discordgo.ActivityType {
	switch activityType {
	case "Playing":
		return discordgo.ActivityTypeGame
	case "Streaming":
		return discordgo.ActivityTypeStreaming
	case "Listening":
		return discordgo.ActivityTypeListening
	case "Watching":
		return discordgo.ActivityTypeWatching
	case "Competing":
		return discordgo.ActivityTypeCompeting
	default:
		return discordgo.ActivityTypeGame
	}
}

func getPlatform(activityType, platform string) (string, string) {
	if activityType != "Playing" {
		return "", ""
	}

	switch platform {
	case "PS4":
		return "ps4", "f01b97df8a090da"
	case "PS5":
		return "ps5", "f01b97df8a090da"
	case "Xbox":
		return "xbox", "f01b97df8a090da"
	case "Samsung":
		return "samsung", "f01b97df8a090da"
	default:
		return "", ""
	}
}

func getButtons(cfg *config.Config) []*client.Button {
	var buttons []*client.Button

	if cfg.Activity.Button1Label != "" {
		buttons = append(buttons, &client.Button{
			Label: cfg.Activity.Button1Label,
			Url:   cfg.Activity.Button1URL,
		})
	}

	if cfg.Activity.Button2Label != "" {
		buttons = append(buttons, &client.Button{
			Label: cfg.Activity.Button2Label,
			Url:   cfg.Activity.Button2URL,
		})
	}

	return buttons
}

func getTimestamps(timeEnabled bool, startTime time.Time) *client.Timestamps {
	if timeEnabled {
		return &client.Timestamps{Start: &startTime}
	}
	return nil
}

func getDGTimestamps(timeEnabled bool, startTime time.Time) *discordgo.TimeStamps {
	if timeEnabled {
		return &discordgo.TimeStamps{StartTimestamp: startTime.UnixMilli()}
	}
	return nil
}

func getParty(cfg *config.Config) *client.Party {
	if cfg.Activity.Invitable {
		return &client.Party{
			ID:         "party_id",
			Players:    1,
			MaxPlayers: 69,
		}
	}
	return nil
}
func getSecrets(cfg *config.Config) *client.Secrets {
	if cfg.Activity.Invitable {
		return &client.Secrets{
			Join:     "join_secret",
			Spectate: "spectate_secret",
			Match:    "match_secret",
		}
	}
	return nil
}

var errorLogged = false

func setRPC(session *discordgo.Session, cfg *config.Config, startTime time.Time) error {
	if !cfg.Activity.Enabled {
		client.Logout()
		return nil
	}

	if cfg.Activity.PauseForOther {
		var hasOther bool
		var selfHasActivity bool

		for _, gSession := range sessions.GatewaySessions {
			if gSession.SessionID == session.State.SessionID {
				selfHasActivity = len(gSession.Activities) > 0
				continue
			}

			for _, activity := range gSession.Activities {
				if activity.Type == discordgo.ActivityTypeCustom {
					continue
				}

				hasOther = true
				break
			}
		}

		if hasOther {
			client.Logout()

			if selfHasActivity {
				return session.UpdateStatusComplex(discordgo.UpdateStatusData{
					Activities: []*discordgo.Activity{},
					AFK:        true,
					Status:     sessions.Status,
				})
			}

			return nil
		}
	}

	switch cfg.Activity.ActivityType {
	case "RPC":
		if constants.HOSTING {
			if !errorLogged {
				logging.UI.Error("RPC is not supported when hosting the bot. Please use Custom Activity instead.")
				errorLogged = true
			}

			return nil
		}

		err := client.Login(cfg.Activity.ClientID)
		if err != nil {
			if !errorLogged {
				logging.UI.Error("RPC Login failed: " + err.Error())
				errorLogged = true
			}

			return err
		}

		err = client.SetActivity(client.Activity{
			Details:    variables.Parse(cfg.Activity.Details),
			State:      variables.Parse(cfg.Activity.State),
			LargeImage: cfg.Activity.LargeImage,
			LargeText:  variables.Parse(cfg.Activity.LargeImageText),
			SmallImage: cfg.Activity.SmallImage,
			SmallText:  variables.Parse(cfg.Activity.SmallImageText),
			Buttons:    getButtons(cfg),
			Timestamps: getTimestamps(cfg.Activity.Time, startTime),
			Party:      getParty(cfg),
			Secrets:    getSecrets(cfg),
		})
		if err != nil {
			logging.UI.Error("RPC Set Activity failed: " + err.Error())
			return err
		}

		return nil
	case "CustomActivity":
		client.Logout()

		buttons := []string{}
		buttonUrls := []string{}

		if cfg.Activity.Button1Label != "" {
			label := variables.Parse(cfg.Activity.Button1Label)

			buttons = append(buttons, utils.String.Limit(label, 31))
			buttonUrls = append(buttonUrls, variables.Parse(cfg.Activity.Button1URL))
		}

		if cfg.Activity.Button2Label != "" {
			label := variables.Parse(cfg.Activity.Button2Label)

			buttons = append(buttons, utils.String.Limit(label, 31))
			buttonUrls = append(buttonUrls, variables.Parse(cfg.Activity.Button2URL))
		}
		appID := "898932653049393172"

		if cfg.Activity.ClientID != "" {
			appID = cfg.Activity.ClientID
		}

		assets := &discordgo.Assets{
			LargeImageID: variables.Parse(cfg.Activity.LargeImage),
			LargeText:    utils.String.Limit(variables.Parse(cfg.Activity.LargeImageText), 128),
			SmallImageID: variables.Parse(cfg.Activity.SmallImage),
			SmallText:    utils.String.Limit(variables.Parse(cfg.Activity.SmallImageText), 128),
		}

		platform, platformID := getPlatform(cfg.Activity.Type, cfg.Activity.Platform)

		var activity *discordgo.Activity

		switch cfg.Activity.Type {
		case "Spotify":
			activity = &discordgo.Activity{
				Type:    discordgo.ActivityType(2),
				Name:    "Spotify",
				State:   utils.String.Limit(variables.Parse(cfg.Activity.State), 128),
				Details: utils.String.Limit(variables.Parse(cfg.Activity.Details), 128),
				Assets:  assets,
				Timestamps: &discordgo.TimeStamps{
					StartTimestamp: time.Now().Unix() - 4140000,
					EndTimestamp:   time.Now().Unix() - 3880000,
				},
				Party: &discordgo.Party{
					ID: "spotify:" + session.State.User.ID,
				},
				Flags: 48,
			}
		default:
			activity = &discordgo.Activity{
				Name:       utils.String.Limit(variables.Parse(cfg.Activity.Name), 128),
				State:      utils.String.Limit(variables.Parse(cfg.Activity.State), 128),
				Details:    utils.String.Limit(variables.Parse(cfg.Activity.Details), 128),
				Timestamps: getDGTimestamps(cfg.Activity.Time, startTime),
				Type:       getActivityType(cfg.Activity.Type),
				Assets:     assets,
				Buttons:    buttons,
				Metadata: &discordgo.ActivityMetadata{
					ButtonURLs: buttonUrls,
				},
				ApplicationID: &appID,
				Platform:      platform,
				ID:            platformID,
			}
		}

		if cfg.Activity.Type == "Streaming" {
			activity.URL = cfg.Activity.StreamingURL

			if !strings.HasPrefix(activity.URL, "https://") {
				activity.URL = "https://twitch.tv/ninja"
			}
		}

		return session.UpdateStatusComplex(discordgo.UpdateStatusData{
			Activities: []*discordgo.Activity{activity},
			AFK:        true,
			Status:     sessions.Status,
		})
	default:
		return nil
	}
}

func StartRPC(session *discordgo.Session) {
	go func() {
		startTime := time.Now()

		for {
			time.Sleep(time.Second * 10)

			cfg, err := config.GetConfig()
			if err != nil {
				logging.UI.Error("Error getting config: " + err.Error())
				return
			}

			_ = setRPC(session, &cfg, startTime)
		}
	}()
}
