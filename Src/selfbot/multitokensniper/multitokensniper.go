package multitokensniper

import (
	"time"

	"github.com/LightningDev1/LB-Selfbot/selfbot/config"
	"github.com/LightningDev1/LB-Selfbot/selfbot/errorhandler"
	"github.com/LightningDev1/LB-Selfbot/selfbot/eventhandlers"
	"github.com/LightningDev1/LB-Selfbot/selfbot/logging"
	"github.com/LightningDev1/LB-Selfbot/selfbot/snipers"
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

	session, err := discordgo.New(s.Token)
	if err != nil {
		return err
	}

	s.Session = session
	s.Session.SetIdentify(discordgo.IdentifyDiscordClient)
	s.Session.ErrorHandler = errorhandler.ErrorHandler
	s.Session.StateEnabled = false

	s.Session.AddHandler(func(_ *discordgo.Session, ready *discordgo.Ready) {
		s.State = StateConnected
		s.Guilds = len(ready.Guilds)
	})

	s.Session.AddHandler(func(_ *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author == nil || m.Message == nil {
			return
		}
		cfg, err := config.GetConfig()
		if err != nil {
			return
		}
		eventhandlers.EventHandlers.Settings.AlreadyCheckedNitroCodes, eventhandlers.EventHandlers.Settings.SentNitroCodes = snipers.CheckNitroSniper(
			eventhandlers.EventHandlers.Session,
			m.Message,
			cfg,
			eventhandlers.EventHandlers.Settings.CountryCode,
			eventhandlers.EventHandlers.Settings.PaymentSourceID,
			eventhandlers.EventHandlers.Settings.AlreadyCheckedNitroCodes,
			eventhandlers.EventHandlers.Settings.SentNitroCodes,
			s.Session.State.User.String(),
		)

		snipers.CheckPrivnoteSniper(s.Session, m.Message, &cfg, s.Session.State.User.String())

		snipers.CheckGiveawaySniper(s.Session, m.Message, &cfg, s.Session.State.User.String())
	})

	s.Session.AddHandler(func(_ *discordgo.Session, m *discordgo.MessageUpdate) {
		if m.Author == nil || m.Message == nil {
			return
		}
		cfg, err := config.GetConfig()
		if err != nil {
			return
		}

		eventhandlers.EventHandlers.Settings.AlreadyCheckedNitroCodes, eventhandlers.EventHandlers.Settings.SentNitroCodes = snipers.CheckNitroSniper(
			s.Session,
			m.Message,
			cfg,
			eventhandlers.EventHandlers.Settings.CountryCode,
			eventhandlers.EventHandlers.Settings.PaymentSourceID,
			eventhandlers.EventHandlers.Settings.AlreadyCheckedNitroCodes,
			eventhandlers.EventHandlers.Settings.SentNitroCodes,
			s.Session.State.User.String(),
		)

		snipers.CheckPrivnoteSniper(s.Session, m.Message, &cfg, s.Session.State.User.String())

		snipers.CheckGiveawaySniperEdit(s.Session, m.Message, &cfg, s.Session.State.User.String())
	})

	if err = s.Session.Open(); err == nil {
		s.State = StateConnecting
	}

	return nil
}

var SniperInstances = []*SniperInstance{}

func Init(cfg config.Config) {
	if !cfg.Snipers.MultiTokenSniper {
		return
	}

	for _, token := range cfg.Snipers.Tokens {
		if token == cfg.Account.Discord.Token {
			return
		}
		for _, instance := range SniperInstances {
			if instance.Token == token {
				return
			}
		}
		instance := &SniperInstance{
			State: StateDisconnected,
			Token: token,
		}
		instance.Start()
		SniperInstances = append(SniperInstances, instance)
	}

	if len(SniperInstances) > 0 {
		go func() {
			time.Sleep(time.Second * 15)

			guilds := 0

			for _, instance := range SniperInstances {
				guilds += instance.Guilds
			}

			guilds += len(eventhandlers.EventHandlers.Session.State.Guilds)

			logging.Logging.MultiTokenSniperConnected(len(SniperInstances), guilds)
		}()
	}
}
