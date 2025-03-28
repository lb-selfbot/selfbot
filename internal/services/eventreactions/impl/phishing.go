package impl

import (
	"net/url"
	"slices"
	"strings"

	"github.com/lb-selfbot/selfbot/internal/http"
	"github.com/lb-selfbot/selfbot/internal/logging"
	"github.com/lb-selfbot/selfbot/internal/services/config"
	"github.com/lb-selfbot/selfbot/internal/services/eventbroker"
	"github.com/lb-selfbot/selfbot/internal/services/eventreactions/types"
	"github.com/lb-selfbot/selfbot/internal/services/ui/featurelogs"

	"github.com/lb-selfbot/discordgo"
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
	phishing := &Phishing{
		eventBroker: eventBroker,
		configSvc:   configSvc,
	}

	resp, err := http.DefaultClient().Get(phishingDomainsURL)
	if err != nil {
		logging.Error("Phishing.NewPhishing: Failed to fetch phishing domains:", err)
		return nil
	}
	phishing.phishingDomains = resp.Json().Get("domains").MustStringSlice()

	resp, err = http.DefaultClient().Get(suspiciousDomainsURL)
	if err != nil {
		logging.Error("Phishing.NewPhishing: Failed to fetch suspicious domains:", err)
		return nil
	}
	phishing.suspiciousDomains = resp.Json().Get("domains").MustStringSlice()

	return phishing
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
