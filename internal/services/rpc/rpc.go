package rpc

import (
	"fmt"
	"strings"
	"time"

	"github.com/hugolgst/rich-go/client"
	"github.com/lb-selfbot/discordgo"
	"github.com/lb-selfbot/selfbot/internal/constants"
	"github.com/lb-selfbot/selfbot/internal/services/config"
	"github.com/lb-selfbot/selfbot/internal/services/eventbroker"
	"github.com/lb-selfbot/selfbot/internal/services/sessions"
	"github.com/lb-selfbot/selfbot/internal/services/variables"
	"github.com/lb-selfbot/selfbot/internal/utils"
)

// Service handles Discord Rich Presence.
type Service interface{}

type rpcService struct {
	eventBroker  eventbroker.Service
	configSvc    config.Service
	sessionsSvc  sessions.Service
	variablesSvc variables.Service

	startTime time.Time
}

// NewService creates a new instance of RPCService.
func NewService(eventBroker eventbroker.Service, configSvc config.Service, sessionsSvc sessions.Service, variablesSvc variables.Service) Service {
	svc := &rpcService{eventBroker: eventBroker, configSvc: configSvc, sessionsSvc: sessionsSvc, variablesSvc: variablesSvc}
	svc.subscribe()
	return svc
}

func (r *rpcService) subscribe() {
	r.eventBroker.Subscribe(eventbroker.EventSessionReady, func(payload any) {
		session, ok := payload.(*discordgo.Session)
		if !ok {
			return
		}

		r.startTime = time.Now()
		// NOTE: this is a blocking function, but we are already in a goroutine.
		r.start(session)
	})
}

func (r *rpcService) start(session *discordgo.Session) {
	timer := time.NewTimer(time.Second * 10)
	defer timer.Stop()

	for range timer.C {
		cfg := r.configSvc.GetConfig()
		if err := r.setRPC(session, cfg); err != nil {
			r.eventBroker.Publish(eventbroker.EventLogError, "Failed to set RPC: "+err.Error())
			continue
		}
	}
}

func (r *rpcService) setRPC(session *discordgo.Session, cfg *config.Config) error {
	if !cfg.Activity.Enabled {
		client.Logout()
		return nil
	}

	if cfg.Activity.PauseForOther && r.checkPauseForOther(session) {
		return nil
	}

	switch cfg.Activity.ActivityType {
	case "RPC":
		if constants.HOSTING {
			r.eventBroker.Publish(eventbroker.EventLogError, "RPC is not supported when hosting the bot. Please use Custom Activity instead.")
			cfg.Activity.ActivityType = "CustomActivity"
			return r.configSvc.SaveConfig(cfg)
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

func (r *rpcService) checkPauseForOther(session *discordgo.Session) bool {
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

func (r *rpcService) getActivityType(activityType string) discordgo.ActivityType {
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

func (r *rpcService) getPlatform(activityType, platform string) (string, string) {
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

func (r *rpcService) getButtons(cfg *config.Config) []*client.Button {
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

func (r *rpcService) getTimestamps(timeEnabled bool) *client.Timestamps {
	if !timeEnabled {
		return nil
	}

	return &client.Timestamps{Start: &r.startTime}
}

func (r *rpcService) getDGTimestamps(timeEnabled bool) *discordgo.TimeStamps {
	if !timeEnabled {
		return nil
	}

	return &discordgo.TimeStamps{StartTimestamp: r.startTime.UnixMilli()}
}

func (r *rpcService) getParty(cfg *config.Config) *client.Party {
	if !cfg.Activity.Invitable {
		return nil
	}

	return &client.Party{
		ID:         "party_id",
		Players:    1,
		MaxPlayers: 69,
	}
}

func (r *rpcService) getSecrets(cfg *config.Config) *client.Secrets {
	if !cfg.Activity.Invitable {
		return nil
	}

	return &client.Secrets{
		Join:     "join_secret",
		Spectate: "spectate_secret",
		Match:    "match_secret",
	}
}
