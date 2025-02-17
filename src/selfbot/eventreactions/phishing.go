package eventreactions

import (
	"net/url"
	"slices"
	"strings"

	"github.com/LightningDev1/LB-Selfbot/selfbot/config"
	"github.com/LightningDev1/LB-Selfbot/selfbot/http"
	"github.com/LightningDev1/LB-Selfbot/selfbot/logging"

	"github.com/LightningDev1/LB-Selfbot/selfbot/utils"
	"github.com/LightningDev1/discordgo"
)

var PhishingDomains []string
var SuspiciousDomains []string

type Phishing struct{}

func NewPhishing() *Phishing {
	if len(PhishingDomains) == 0 || len(SuspiciousDomains) == 0 {

		httpResult := http.Get("https://raw.githubusercontent.com/nikolaischunk/discord-phishing-links/main/domain-list.json")
		PhishingDomains = httpResult.Json().Get("domains").MustStringArray()

		httpResult = http.Get("https://raw.githubusercontent.com/nikolaischunk/discord-phishing-links/main/suspicious-list.json")
		SuspiciousDomains = httpResult.Json().Get("domains").MustStringArray()
	}

	return &Phishing{}
}

func (p *Phishing) Type() EventReactionType {
	return ANTI_PHISHING
}

func (p *Phishing) Description() string {
	return "Anti Phishing (default)"
}

func (p *Phishing) Equals(other EventReaction) bool {
	_, ok := other.(*Phishing)

	return ok
}

func (p *Phishing) Data() map[string]any {
	return map[string]any{}
}

func (p *Phishing) Handler(session *discordgo.Session, event any) {
	message, ok := event.(*discordgo.MessageCreate)
	if !ok {
		utils.Logging.Debug("Phishing.Handler: Event is not a message create event")
		return
	}

	cfg, err := config.GetConfig()
	if err != nil {
		utils.Logging.Debug("Phishing.Handler: Error getting config: " + err.Error())
		return
	}

	if !cfg.Protections.AntiPhishing {
		return
	}

	for _, part := range strings.Split(message.Content, " ") {
		link, err := url.Parse(part)
		if err != nil {
			continue
		}

		hostname := link.Hostname()

		isPhishing := slices.Contains(PhishingDomains, hostname)
		isSuspicious := slices.Contains(SuspiciousDomains, hostname)

		if isPhishing || isSuspicious {
			_ = session.MessageAck(message.ChannelID, message.ID)
			logging.Logging.PhishingDetected(session, message.Message, hostname)
		}
	}
}
