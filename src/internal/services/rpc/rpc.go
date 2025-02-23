package rpc

import (
	"fmt"
	"strings"
	"time"

	"github.com/LightningDev1/LB-Selfbot/internal/config"
	"github.com/LightningDev1/LB-Selfbot/internal/constants"
	"github.com/LightningDev1/LB-Selfbot/internal/services/interfaces"
	"github.com/LightningDev1/LB-Selfbot/internal/utils"
	"github.com/LightningDev1/discordgo"
	"github.com/hugolgst/rich-go/client"
)

type RPCServiceImpl struct {
	uiSvc        interfaces.UIService
	discordSvc   interfaces.DiscordService
	sessionsSvc  interfaces.SessionsService
	variablesSvc interfaces.VariablesService

	startTime time.Time
}

func (r *RPCServiceImpl) Inject(
	uiSvc interfaces.UIService,
	discordSvc interfaces.DiscordService,
	sessionsSvc interfaces.SessionsService,
	variablesSvc interfaces.VariablesService,
) {
	r.uiSvc = uiSvc
	r.discordSvc = discordSvc
	r.sessionsSvc = sessionsSvc
	r.variablesSvc = variablesSvc
}

func (r *RPCServiceImpl) Start() {
	go r.start()
}

func (r *RPCServiceImpl) start() {
	r.startTime = time.Now()

	timer := time.NewTimer(time.Second * 10)
	defer timer.Stop()

	for range timer.C {
		cfg, err := config.GetConfig()
		if err != nil {
			r.uiSvc.LogError("Failed to get config: ", err)
			continue
		}

		err = r.setRPC(&cfg)
		if err != nil {
			r.uiSvc.LogError("Failed to set RPC: ", err)
			continue
		}
	}
}

func (r *RPCServiceImpl) setRPC(cfg *config.Config) error {
	if !cfg.Activity.Enabled {
		client.Logout()
		return nil
	}

	if cfg.Activity.PauseForOther && r.checkPauseForOther() {
		return nil
	}

	switch cfg.Activity.ActivityType {
	case "RPC":
		if constants.HOSTING {
			r.uiSvc.LogError("RPC is not supported when hosting the bot. Please use Custom Activity instead.")
			cfg.Activity.ActivityType = "CustomActivity"
			return cfg.Save()
		}

		err := client.Login(cfg.Activity.ClientID)
		if err != nil {
			return fmt.Errorf("failed to login to RPC: %w", err)
		}

		err = client.SetActivity(client.Activity{
			Details:    r.variablesSvc.Parse(cfg.Activity.Details),
			State:      r.variablesSvc.Parse(cfg.Activity.State),
			LargeImage: cfg.Activity.LargeImage,
			LargeText:  r.variablesSvc.Parse(cfg.Activity.LargeImageText),
			SmallImage: cfg.Activity.SmallImage,
			SmallText:  r.variablesSvc.Parse(cfg.Activity.SmallImageText),
			Buttons:    r.getButtons(cfg),
			Timestamps: r.getTimestamps(cfg.Activity.Time),
			Party:      r.getParty(cfg),
			Secrets:    r.getSecrets(cfg),
		})
		if err != nil {
			return fmt.Errorf("failed to set RPC activity: %w", err)
		}
	case "CustomActivity":
		client.Logout()
		session := r.discordSvc.GetSession()

		buttons := []string{}
		buttonUrls := []string{}

		if cfg.Activity.Button1Label != "" {
			label := r.variablesSvc.Parse(cfg.Activity.Button1Label)

			buttons = append(buttons, utils.String.Limit(label, 31))
			buttonUrls = append(buttonUrls, r.variablesSvc.Parse(cfg.Activity.Button1URL))
		}

		if cfg.Activity.Button2Label != "" {
			label := r.variablesSvc.Parse(cfg.Activity.Button2Label)

			buttons = append(buttons, utils.String.Limit(label, 31))
			buttonUrls = append(buttonUrls, r.variablesSvc.Parse(cfg.Activity.Button2URL))
		}

		appID := "898932653049393172"
		if cfg.Activity.ClientID != "" {
			appID = cfg.Activity.ClientID
		}

		assets := &discordgo.Assets{
			LargeImageID: r.variablesSvc.Parse(cfg.Activity.LargeImage),
			LargeText:    utils.String.Limit(r.variablesSvc.Parse(cfg.Activity.LargeImageText), 128),
			SmallImageID: r.variablesSvc.Parse(cfg.Activity.SmallImage),
			SmallText:    utils.String.Limit(r.variablesSvc.Parse(cfg.Activity.SmallImageText), 128),
		}

		platform, platformID := r.getPlatform(cfg.Activity.Type, cfg.Activity.Platform)

		var activity *discordgo.Activity

		switch cfg.Activity.Type {
		case "Spotify":
			activity = &discordgo.Activity{
				Type:    discordgo.ActivityType(2),
				Name:    "Spotify",
				State:   utils.String.Limit(r.variablesSvc.Parse(cfg.Activity.State), 128),
				Details: utils.String.Limit(r.variablesSvc.Parse(cfg.Activity.Details), 128),
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
				Name:       utils.String.Limit(r.variablesSvc.Parse(cfg.Activity.Name), 128),
				State:      utils.String.Limit(r.variablesSvc.Parse(cfg.Activity.State), 128),
				Details:    utils.String.Limit(r.variablesSvc.Parse(cfg.Activity.Details), 128),
				Timestamps: r.getDGTimestamps(cfg.Activity.Time),
				Type:       r.getActivityType(cfg.Activity.Type),
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
			Status:     r.sessionsSvc.GetStatus(),
		})
	}

	return nil
}

func (r *RPCServiceImpl) checkPauseForOther() bool {
	session := r.discordSvc.GetSession()
	var hasOther, selfHasActivity bool

	for _, gSession := range r.sessionsSvc.GetGatewaySessions() {
		if gSession.SessionID == session.State.SessionID {
			selfHasActivity = len(gSession.Activities) > 0
			continue
		}

		for _, activity := range gSession.Activities {
			if activity.Type != discordgo.ActivityTypeCustom {
				hasOther = true
				break
			}
		}
	}

	if !hasOther {
		return false
	}

	client.Logout()

	if selfHasActivity {
		_ = session.UpdateStatusComplex(discordgo.UpdateStatusData{
			Activities: []*discordgo.Activity{},
			AFK:        true,
			Status:     r.sessionsSvc.GetStatus(),
		})
	}

	return true
}

func (r *RPCServiceImpl) getActivityType(activityType string) discordgo.ActivityType {
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

func (r *RPCServiceImpl) getPlatform(activityType, platform string) (string, string) {
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

func (r *RPCServiceImpl) getButtons(cfg *config.Config) []*client.Button {
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

func (r *RPCServiceImpl) getTimestamps(timeEnabled bool) *client.Timestamps {
	if !timeEnabled {
		return nil
	}

	return &client.Timestamps{Start: &r.startTime}
}

func (r *RPCServiceImpl) getDGTimestamps(timeEnabled bool) *discordgo.TimeStamps {
	if !timeEnabled {
		return nil
	}

	return &discordgo.TimeStamps{StartTimestamp: r.startTime.UnixMilli()}
}

func (r *RPCServiceImpl) getParty(cfg *config.Config) *client.Party {
	if !cfg.Activity.Invitable {
		return nil
	}

	return &client.Party{
		ID:         "party_id",
		Players:    1,
		MaxPlayers: 69,
	}
}

func (r *RPCServiceImpl) getSecrets(cfg *config.Config) *client.Secrets {
	if !cfg.Activity.Invitable {
		return nil
	}

	return &client.Secrets{
		Join:     "join_secret",
		Spectate: "spectate_secret",
		Match:    "match_secret",
	}
}
