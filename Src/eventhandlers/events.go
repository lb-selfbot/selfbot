package eventhandlers

import (
	"time"

	"github.com/LightningDev1/LB-Selfbot/http"
	"github.com/LightningDev1/LB-Selfbot/snipers"
	"github.com/LightningDev1/dgc"
	"github.com/LightningDev1/discordgo"
	"github.com/bitly/go-simplejson"
)

type FriendRequest struct {
	Time time.Time
	*discordgo.Relationship
}

type EventEntry struct {
	UserID    string
	GuildID   string
	ChannelID string
	Argument  string
}

type SBDEvent struct {
	Event   string
	Message *discordgo.Message
}

type SBDScore struct {
	UserID  string
	Score   float64
	Logged  bool
	Selfbot string
	Events  []SBDEvent
}

type SBDetectable struct {
	String  string
	Selfbot string
	Score   float64
}

type EventHandlerSettings struct {
	IgnoredUsers             []*discordgo.User
	IgnoredGuilds            []*discordgo.Guild
	UsedCommands             []*dgc.Ctx
	PinSpamMessages          []*discordgo.Message
	DeletedMessages          map[string][]*discordgo.Message
	EditedMessages           map[string][]*discordgo.Message
	AlreadyCheckedNitroCodes []string
	SentNitroCodes           map[string][]snipers.NitroCode
	SelfbotDetectionScore    map[string]SBDScore

	PingMuteServers    []string
	PingTimeoutServers []string
	PingBanServers     []string
	PingKickServers    []string
	PingInsultServers  []string
	PingPingServers    []string

	ForceNick       []EventEntry
	ForceKick       []EventEntry
	NoLeave         []EventEntry
	Copy            []EventEntry
	Mock            []EventEntry
	ReplyCopy       []EventEntry
	PinSpam         []EventEntry
	DeleteMessages  []EventEntry
	ReactUser       []EventEntry
	RemoveReactions []EventEntry
	ForceDisconnect []EventEntry

	CountryCode     string
	PaymentSourceID string

	phisingDomains    []string
	suspiciousDomains []string
	friendRequests    []*FriendRequest
	newDMs            []*discordgo.Message
	discordTimeDiff   time.Duration
}

type Events struct {
	Session   *discordgo.Session
	Router    *dgc.Router
	Settings  *EventHandlerSettings
	StartTime time.Time
}

func (e *Events) RegisterHandlers(session *discordgo.Session, router *dgc.Router) {
	e.Session = session
	e.Router = router
	e.Settings = &EventHandlerSettings{}
	e.StartTime = time.Now()
	e.Settings.DeletedMessages = map[string][]*discordgo.Message{}
	e.Settings.EditedMessages = map[string][]*discordgo.Message{}
	e.Settings.SentNitroCodes = map[string][]snipers.NitroCode{}
	e.Settings.SelfbotDetectionScore = map[string]SBDScore{}

	e.Router.RegisterMiddleware(e.OnCommand)

	e.Session.AddHandler(e.OnReady)

	e.Session.AddHandler(e.OnMessage)
	e.Session.AddHandler(e.OnMessageEdit)
	e.Session.AddHandler(e.OnMessageDelete)

	e.Session.AddHandler(e.OnRelationshipAdd)

	e.Session.AddHandler(e.OnGuildBanAdd)

	e.Session.AddHandler(e.OnTypingStart)

	phisingDomains, err := http.GetJsonOld("https://raw.githubusercontent.com/nikolaischunk/discord-phishing-links/main/domain-list.json")
	if err == nil {
		e.Settings.phisingDomains = phisingDomains.Get("domains").MustStringArray()
	}

	suspiciousDomains, err := http.GetJsonOld("https://raw.githubusercontent.com/nikolaischunk/discord-phishing-links/main/suspicious-list.json")
	if err == nil {
		e.Settings.suspiciousDomains = suspiciousDomains.Get("domains").MustStringArray()
	}

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
