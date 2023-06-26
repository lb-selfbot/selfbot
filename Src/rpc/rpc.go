package rpc

import (
	"strings"
	"time"

	"github.com/LightningDev1/LB-Selfbot/config"
	"github.com/LightningDev1/LB-Selfbot/constants"
	"github.com/LightningDev1/LB-Selfbot/logging"
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
		return &client.Timestamps{
			Start: &startTime,
		}
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

func parseText(text string) string {
	text = strings.ReplaceAll(text, "$VERSION", constants.VERSION)

	return text
}

func setRPC(session *discordgo.Session, cfg *config.Config, startTime time.Time) error {
	if !cfg.Activity.Enabled {
		client.Logout()
		return nil
	}

	switch cfg.Activity.ActivityType {
	case "RPC":
		err := client.Login(cfg.Activity.ClientID)
		if err != nil {
			logging.UI.Error("RPC Login failed: " + err.Error())
			return err
		}

		err = client.SetActivity(client.Activity{
			Details:    parseText(cfg.Activity.Details),
			State:      parseText(cfg.Activity.State),
			LargeImage: cfg.Activity.LargeImage,
			LargeText:  parseText(cfg.Activity.LargeImageText),
			SmallImage: cfg.Activity.SmallImage,
			SmallText:  parseText(cfg.Activity.SmallImageText),
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
			buttons = append(buttons, cfg.Activity.Button1Label)
			buttonUrls = append(buttonUrls, cfg.Activity.Button1URL)
		}

		if cfg.Activity.Button2Label != "" {
			buttons = append(buttons, cfg.Activity.Button2Label)
			buttonUrls = append(buttonUrls, cfg.Activity.Button2URL)
		}

		activity := &discordgo.Activity{
			Name:    parseText(cfg.Activity.Name),
			State:   parseText(cfg.Activity.State),
			Details: parseText(cfg.Activity.Details),
			Type:    getActivityType(cfg.Activity.Type),
			Assets: &discordgo.Assets{
				LargeImageID: cfg.Activity.LargeImage,
				LargeText:    parseText(cfg.Activity.LargeImageText),
				SmallImageID: cfg.Activity.SmallImage,
				SmallText:    parseText(cfg.Activity.SmallImageText),
			},
			Buttons: buttons,
			Metadata: &discordgo.ActivityMetadata{
				ButtonURLs: buttonUrls,
			},
		}

		return session.UpdateStatusComplex(discordgo.UpdateStatusData{
			Activities: []*discordgo.Activity{},
			Game:       activity,
			AFK:        true,
			Status:     "unknown",
		})
	default:
		return nil
	}
}

func StartRPC(session *discordgo.Session) {
	go func() {
		startTime := time.Now()
		for {
			time.Sleep(time.Second * 5)

			cfg, err := config.GetConfig()
			if err != nil {
				logging.UI.Error("Error getting config: " + err.Error())
				return
			}

			_ = setRPC(session, &cfg, startTime)
		}
	}()
}
