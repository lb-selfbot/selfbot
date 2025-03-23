package mts

import (
	"context"
	"fmt"
	"time"

	"github.com/lb-selfbot/discordgo"
	"github.com/lb-selfbot/selfbot/internal/constants"
	"github.com/lb-selfbot/selfbot/internal/services/config"
	"github.com/lb-selfbot/selfbot/internal/services/eventbroker"
	"github.com/lb-selfbot/selfbot/internal/services/eventreactions"
	"github.com/lb-selfbot/selfbot/internal/services/eventreactions/impl"
	"github.com/lb-selfbot/selfbot/internal/services/ui/featurelogs"
)

type sniperState int

const (
	stateDisconnected sniperState = iota
	stateConnecting
	stateConnected

	maxHostingTokens = 5
)

// Service handles the Multi-Token Sniper.
type Service interface{}

type sniperInstance struct {
	state   sniperState
	token   string
	guilds  int
	session *discordgo.Session
}

type mtsService struct {
	eventBroker       eventbroker.Service
	configSvc         config.Service
	eventReactionsSvc eventreactions.Service

	instances []*sniperInstance
}

// NewService creates a new instance of MTSService.
func NewService(eventBroker eventbroker.Service, configSvc config.Service, eventReactionsSvc eventreactions.Service) Service {
	svc := &mtsService{
		eventBroker:       eventBroker,
		configSvc:         configSvc,
		eventReactionsSvc: eventReactionsSvc,
	}
	svc.subscribe()
	return svc
}

func (m *mtsService) subscribe() {
	m.eventBroker.Subscribe(eventbroker.EventSessionReady, func(payload any) {
		session, ok := payload.(*discordgo.Session)
		if !ok {
			return
		}

		m.init(session)
	})
}

func (m *mtsService) init(session *discordgo.Session) {
	cfg := m.configSvc.GetConfig()
	if !cfg.Snipers.MultiTokenSniper {
		return
	}

	tokens := cfg.Snipers.Tokens
	if constants.HOSTING && len(tokens) > maxHostingTokens {
		tokens = tokens[:maxHostingTokens]
	}

	m.instances = make([]*sniperInstance, 0, len(tokens))

	for _, token := range tokens {
		if m.isTokenUsed(session, token) {
			continue
		}

		instance := &sniperInstance{
			state: stateDisconnected,
			token: token,
		}

		m.instances = append(m.instances, instance)
		m.startInstance(context.Background(), session, instance)
	}

	if len(m.instances) > 0 {
		guilds := len(session.State.Guilds)

		for _, instance := range m.instances {
			guilds += instance.guilds
		}

		m.eventBroker.Publish(eventbroker.EventLogFeature, &featurelogs.MultiTokenSniper{
			NumTokens: len(m.instances),
			NumGuilds: guilds,
		})
	}
}

func (m *mtsService) isTokenUsed(session *discordgo.Session, token string) bool {
	if session.Identify.Token == token {
		return true
	}

	for _, instance := range m.instances {
		if instance.token == token {
			return true
		}
	}

	return false
}

func (m *mtsService) startInstance(context context.Context, mainSession *discordgo.Session, instance *sniperInstance) error {
	instance.state = stateConnecting

	instance.session, _ = discordgo.New(instance.token)
	instance.session.SetIdentify(discordgo.IdentifyDiscordClient)
	instance.session.ErrorHandler = func(session *discordgo.Session, err any) {
		m.eventBroker.Publish(eventbroker.EventCaptureError, err)
	}
	instance.session.StateEnabled = false

	var snipers *impl.Snipers
	readyChan := make(chan struct{})

	instance.session.AddHandler(func(s *discordgo.Session, e *discordgo.Ready) {
		instance.state = stateConnected
		instance.guilds = len(e.Guilds)

		snipers = impl.NewSnipers(
			m.eventBroker,
			m.configSvc,
			mainSession,
			m.eventReactionsSvc.GetCountryCode,
			m.eventReactionsSvc.GetPaymentSrcID,
			fmt.Sprintf("%s (%s)", e.User.DisplayName, e.User.Username),
		)

		close(readyChan)
	})

	instance.session.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Message == nil || m.Author == nil || snipers == nil {
			return
		}

		snipers.Handler(s, m)
	})

	instance.session.AddHandler(func(s *discordgo.Session, m *discordgo.MessageUpdate) {
		if m.Message == nil || m.Author == nil {
			return
		}

		if m.BeforeUpdate == nil {
			m.BeforeUpdate = m.Message
		}

		snipers.Handler(s, m)
	})

	err := instance.session.Open()
	if err != nil {
		instance.state = stateDisconnected
		return err
	}

	select {
	case <-readyChan:
		return nil
	case <-time.After(time.Second * 10):
		return fmt.Errorf("timed out waiting for ready event")
	case <-context.Done():
		return context.Err()
	}
}
