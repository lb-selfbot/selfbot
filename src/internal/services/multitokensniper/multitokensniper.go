package multitokensniper

import (
	"context"
	"fmt"
	"time"

	"github.com/LightningDev1/LB-Selfbot/internal/constants"
	"github.com/LightningDev1/LB-Selfbot/internal/services/interfaces"
	"github.com/LightningDev1/LB-Selfbot/selfbot/config"
	"github.com/LightningDev1/LB-Selfbot/selfbot/errorhandler"
	"github.com/LightningDev1/LB-Selfbot/selfbot/eventreactions"
	"github.com/LightningDev1/discordgo"
)

type SniperState int

const (
	StateDisconnected SniperState = iota
	StateConnecting
	StateConnected

	maxHostingTokens = 5
)

type SniperInstance struct {
	State   SniperState
	Token   string
	Guilds  int
	Session *discordgo.Session
}

type MultiTokenSniperServiceImpl struct {
	uiSvc      interfaces.UIService
	discordSvc interfaces.DiscordService

	instances []*SniperInstance
}

func (m *MultiTokenSniperServiceImpl) Inject(uiSvc interfaces.UIService, discordSvc interfaces.DiscordService) {
	m.discordSvc = discordSvc
}

func (m *MultiTokenSniperServiceImpl) Init(context context.Context) error {
	cfg, err := config.GetConfig()
	if err != nil {
		return err
	}

	if !cfg.Snipers.MultiTokenSniper {
		return nil
	}

	tokens := cfg.Snipers.Tokens
	if constants.HOSTING && len(tokens) > maxHostingTokens {
		tokens = tokens[:maxHostingTokens]
	}
	m.instances = make([]*SniperInstance, 0, len(tokens))

	for _, token := range tokens {
		if m.isTokenUsed(token) {
			continue
		}

		instance := &SniperInstance{
			State: StateDisconnected,
			Token: token,
		}

		m.instances = append(m.instances, instance)
		m.startInstance(context, instance)
	}

	if len(m.instances) > 0 {
		guilds := len(m.discordSvc.GetSession().State.Guilds)

		for _, instance := range m.instances {
			guilds += instance.Guilds
		}

		// m.uiSvc.MTSConnected(len(m.instances), guilds)
	}

	return nil
}

func (m *MultiTokenSniperServiceImpl) isTokenUsed(token string) bool {
	for _, instance := range m.instances {
		if instance.Token == token {
			return true
		}
	}

	return token == m.discordSvc.GetSession().Identify.Token
}

func (m *MultiTokenSniperServiceImpl) startInstance(context context.Context, instance *SniperInstance) error {
	instance.State = StateConnecting

	instance.Session, _ = discordgo.New(instance.Token)
	instance.Session.SetIdentify(discordgo.IdentifyDiscordClient)
	instance.Session.ErrorHandler = errorhandler.ErrorHandler
	instance.Session.StateEnabled = false

	var snipers *eventreactions.Snipers
	readyChan := make(chan struct{})

	instance.Session.AddHandler(func(s *discordgo.Session, e *discordgo.Ready) {
		instance.State = StateConnected
		instance.Guilds = len(e.Guilds)

		username := fmt.Sprintf("%s (%s)", e.User.DisplayName, e.User.Username)
		snipers = eventreactions.NewSnipers(m.discordSvc.GetSession(), username)

		close(readyChan)
	})

	instance.Session.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Message == nil || m.Author == nil || snipers == nil {
			return
		}

		snipers.Handler(s, m)
	})

	instance.Session.AddHandler(func(s *discordgo.Session, m *discordgo.MessageUpdate) {
		if m.Message == nil || m.Author == nil {
			return
		}

		if m.BeforeUpdate == nil {
			m.BeforeUpdate = m.Message
		}

		snipers.Handler(s, m)
	})

	err := instance.Session.Open()
	if err != nil {
		instance.State = StateDisconnected
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
