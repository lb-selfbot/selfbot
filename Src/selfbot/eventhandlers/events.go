package eventhandlers

import (
	"github.com/LightningDev1/LB-Selfbot/selfbot/snipers"
	"github.com/LightningDev1/dgc"
	"github.com/LightningDev1/discordgo"
	"github.com/bitly/go-simplejson"
)

type EventHandlerSettings struct {
	AlreadyCheckedNitroCodes []string
	SentNitroCodes           map[string][]snipers.NitroCode
	CountryCode              string
	PaymentSourceID          string
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
	e.Settings.SentNitroCodes = map[string][]snipers.NitroCode{}

	e.Router.RegisterMiddleware(e.OnCommand)

	e.Session.AddHandler(e.OnReady)

	e.Session.AddHandler(e.OnMessage)
	e.Session.AddHandler(e.OnMessageEdit)
	e.Session.AddHandler(e.OnMessageDelete)

	// Get data for nitro sniper
	result, err := e.Session.Request("GET", "https://discord.com/api/v9/users/@me/billing/payment-sources", nil)
	if err == nil {
		resultJson, err := simplejson.NewJson(result)
		if err != nil {
			resultJson = simplejson.New()
		}

		paymentSourcesAmount := len(resultJson.MustArray())

		if paymentSourcesAmount != 0 {
			paymentSource := getDefaultPaymentSource(paymentSourcesAmount, resultJson)
			e.Settings.PaymentSourceID = paymentSource.Get("id").MustString()
			e.Settings.CountryCode = paymentSource.Get("country").MustString()
		}
	}
}

func getDefaultPaymentSource(amount int, json *simplejson.Json) *simplejson.Json {
	// Try to get the default payment source
	for i := 0; i < amount; i++ {
		if !json.GetIndex(i).Get("default").MustBool() {
			continue
		}

		return json.GetIndex(i)
	}

	return json.GetIndex(0)
}

var EventHandlers = Events{}
