package multitokensniper

import (
	"fmt"
	"time"

	"github.com/LightningDev1/LB-Selfbot/selfbot/config"
	"github.com/LightningDev1/LB-Selfbot/selfbot/errorhandler"
	"github.com/LightningDev1/LB-Selfbot/selfbot/eventreactions"
	"github.com/LightningDev1/LB-Selfbot/selfbot/logging"
	"github.com/LightningDev1/LB-Selfbot/selfbot/constants"
	"github.com/LightningDev1/discordgo"
)

type SniperState int

const (
	StateDisconnected SniperState = iota
	StateConnecting
	StateConnected
)

type SniperInstance struct {
	State   SniperState
	Token   string
	Guilds  int
	Session *discordgo.Session
}

func (s *SniperInstance) Start() error {
	s.State = StateDisconnected

	var err error
	s.Session, err = discordgo.New(s.Token)
	if err != nil {
		return err
	}

	s.Session.SetIdentify(discordgo.IdentifyDiscordClient)
	s.Session.ErrorHandler = errorhandler.ErrorHandler
	s.Session.StateEnabled = false

	var snipers *eventreactions.Snipers

	s.Session.AddHandler(func(_ *discordgo.Session, ready *discordgo.Ready) {
		s.State = StateConnected
		s.Guilds = len(ready.Guilds)

		username := fmt.Sprintf("%s (%s)", ready.User.DisplayName, ready.User.Username)
		snipers = eventreactions.NewSnipers(username)
	})

	s.Session.AddHandler(func(_ *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Message == nil || m.Author == nil {
			return
		}

		snipers.Handler(s.Session, m)
	})

	s.Session.AddHandler(func(_ *discordgo.Session, m *discordgo.MessageUpdate) {
		if m.Message == nil || m.Author == nil {
			return
		}

		if m.BeforeUpdate == nil {
			m.BeforeUpdate = m.Message
		}

		snipers.Handler(s.Session, m)
	})

	if err = s.Session.Open(); err == nil {
		s.State = StateConnecting
	}

	return nil
}

var SniperInstances = []*SniperInstance{}

var maxHostingTokens = 5

func Init(mainSession *discordgo.Session, cfg config.Config) {
	if !cfg.Snipers.MultiTokenSniper {
		return
	}

	tokens := cfg.Snipers.Tokens
	if len(tokens) > maxHostingTokens && constants.HOSTING {
		tokens = tokens[:maxHostingTokens]
	}

	for _, token := range tokens {
		// Check if the token is already in use
		if token == mainSession.Identify.Token {
			return
		}

		for _, instance := range SniperInstances {
			if instance.Token == token {
				return
			}
		}

		// Create a new sniper instance
		instance := &SniperInstance{
			State: StateDisconnected,
			Token: token,
		}

		instance.Start()

		SniperInstances = append(SniperInstances, instance)
	}

	go func() {
		if len(SniperInstances) == 0 {
			return
		}

		time.Sleep(time.Second * 15)

		// Count the total amount of guilds
		guilds := len(mainSession.State.Guilds)

		for _, instance := range SniperInstances {
			guilds += instance.Guilds
		}

		logging.Logging.MultiTokenSniperConnected(len(SniperInstances), guilds)
	}()
}
