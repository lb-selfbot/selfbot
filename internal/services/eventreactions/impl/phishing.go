package impl

import (
	"net/url"
	"slices"
	"strings"

	"github.com/LightningDev1/LB-Selfbot/internal/http"
	"github.com/LightningDev1/LB-Selfbot/internal/services/config"
	"github.com/LightningDev1/LB-Selfbot/internal/services/eventbroker"
	"github.com/LightningDev1/LB-Selfbot/internal/services/eventreactions/types"
	"github.com/LightningDev1/LB-Selfbot/internal/services/ui/featurelogs"

	"github.com/LightningDev1/LB-Selfbot/internal/utils"
	"github.com/LightningDev1/discordgo"
)

const (
	phishingDomainsURL   = "https://raw.githubusercontent.com/nikolaischunk/discord-phishing-links/main/domain-list.json"
	suspiciousDomainsURL = "https://raw.githubusercontent.com/nikolaischunk/discord-phishing-links/main/suspicious-list.json"
)

type Phishing struct {
	eventBroker eventbroker.Service
	configSvc   config.Service

	phishingDomains   []string
	suspiciousDomains []string
}

func NewPhishing(eventBroker eventbroker.Service, configSvc config.Service) *Phishing {
	return &Phishing{
		eventBroker:       eventBroker,
		configSvc:         configSvc,
		phishingDomains:   http.Get(phishingDomainsURL).Json().Get("domains").MustStringArray(),
		suspiciousDomains: http.Get(suspiciousDomainsURL).Json().Get("domains").MustStringArray(),
	}
}

func (p *Phishing) Type() types.EventReactionType {
	return types.AntiPhishing
}

func (p *Phishing) Description() string {
	return "Anti Phishing (default)"
}

func (p *Phishing) Equals(other types.EventReaction) bool {
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

	cfg := p.configSvc.GetConfig()
	if !cfg.Protections.AntiPhishing {
		return
	}

	for _, part := range strings.Split(message.Content, " ") {
		link, err := url.Parse(part)
		if err != nil {
			continue
		}

		hostname := link.Hostname()

		isPhishing := slices.Contains(p.phishingDomains, hostname)
		isSuspicious := slices.Contains(p.suspiciousDomains, hostname)

		if isPhishing || isSuspicious {
			_ = session.MessageAck(message.ChannelID, message.ID)

			p.eventBroker.Publish(eventbroker.EventLogFeature, &featurelogs.PhishingDetection{
				Message: message.Message,
				Domain:  hostname,
			})
		}
	}
}
