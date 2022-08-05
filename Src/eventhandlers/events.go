package eventhandlers

import (
	"github.com/LightningDev1/LB-Selfbot/http"
	"github.com/LightningDev1/dgc"
	"github.com/LightningDev1/discordgo"
)

type EventHandlerSettings struct {
	IgnoredUsers  []*discordgo.User
	IgnoredGuilds []*discordgo.Guild

	alreadyCheckedNitroCodes []string
	phisingDomains           []string
	suspiciousDomains        []string
}

type Events struct {
	Session  *discordgo.Session
	Router   *dgc.Router
	Settings *EventHandlerSettings
}

func (e *Events) RegisterHandlers(session *discordgo.Session, router *dgc.Router) {
	e.Session = session
	e.Router = router
	e.Settings = &EventHandlerSettings{}

	e.Router.RegisterMiddleware(e.OnCommand)
	e.Session.AddHandler(e.OnReady)
	e.Session.AddHandler(e.OnMessage)
	e.Session.AddHandler(e.OnMessageEdit)

	phisingDomains, err := http.GetJson("https://raw.githubusercontent.com/nikolaischunk/discord-phishing-links/main/domain-list.json")
	if err == nil {
		e.Settings.phisingDomains = phisingDomains.Get("domains").MustStringArray()
	}

	suspiciousDomains, err := http.GetJson("https://raw.githubusercontent.com/nikolaischunk/discord-phishing-links/main/suspicious-list.json")
	if err == nil {
		e.Settings.suspiciousDomains = suspiciousDomains.Get("domains").MustStringArray()
	}
}

var EventHandlers = Events{}
