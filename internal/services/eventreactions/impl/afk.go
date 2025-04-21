package impl

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/lb-selfbot/selfbot/internal/api"
	"github.com/lb-selfbot/selfbot/internal/logging"

	"github.com/lb-selfbot/selfbot/internal/services/config"
	"github.com/lb-selfbot/selfbot/internal/services/eventbroker"
	"github.com/lb-selfbot/selfbot/internal/services/eventreactions/types"
	"github.com/lb-selfbot/selfbot/internal/services/ui/featurelogs"
	"github.com/lb-selfbot/selfbot/internal/services/variables"

	"github.com/lb-selfbot/discordgo"
	"github.com/lb-selfbot/selfbot/internal/utils"
)

type AFK struct {
	eventBroker  eventbroker.Service
	configSvc    config.Service
	variablesSvc variables.Service

	isReactionMessage func(*discordgo.Message) bool

	lastAFKMessage time.Time
	lastActivity   time.Time
}

func NewAFK(eventBroker eventbroker.Service, configSvc config.Service, variablesSvc variables.Service, isReactionMessage func(*discordgo.Message) bool) *AFK {
	return &AFK{
		eventBroker:       eventBroker,
		configSvc:         configSvc,
		variablesSvc:      variablesSvc,
		isReactionMessage: isReactionMessage,
	}
}

func (a *AFK) Type() types.EventReactionType {
	return types.AFKResponder
}

func (a *AFK) Description() string {
	return "AFK Responder (default)"
}

func (a *AFK) Equals(other types.EventReaction) bool {
	_, ok := other.(*AFK)

	return ok
}

func (a *AFK) HandleLastActivity(session *discordgo.Session, m *discordgo.Message) {
	if m.Author.ID != session.State.User.ID ||
		a.isReactionMessage(m) ||
		time.Since(a.lastAFKMessage) < 5*time.Second {
		return
	}

	a.lastActivity = time.Now()
}

func (a *AFK) IsActive() bool {
	cfg := a.configSvc.GetConfig()

	if !cfg.AFK.MustBeInactive {
		return false
	}

	inactiveTimeout := time.Duration(cfg.AFK.InactiveTimeout) * time.Second
	return time.Since(a.lastActivity) < inactiveTimeout
}

func (a *AFK) Data() map[string]any {
	return map[string]any{}
}

func (a *AFK) GetAFKMessage(m *discordgo.Message, message string) string {
	message = strings.ReplaceAll(message, "$USER", m.Author.Username)
	message = strings.ReplaceAll(message, "$MENTION", m.Author.Mention())
	message = strings.ReplaceAll(message, "$MSGLINK", m.GetLink())
	message = strings.ReplaceAll(message, "$MSGCONTENT", m.ContentWithMentionsReplaced())
	message = strings.ReplaceAll(message, "$AFKSINCE",
		fmt.Sprintf("<t:%d:R>", a.lastActivity.Unix()))

	message = a.variablesSvc.Parse(message)

	return message
}

func (a *AFK) Handler(session *discordgo.Session, event any) {
	message, ok := event.(*discordgo.MessageCreate)
	if !ok {
		return
	}

	a.HandleLastActivity(session, message.Message)

	cfg := a.configSvc.GetConfig()

	activatedOnPing := cfg.AFK.OnPing && message.Message.IsUserMentioned(session.State.User.ID)
	activatedOnDM := cfg.AFK.OnDM && utils.Discord.IsDM(session, message.Message)

	if !(activatedOnPing || activatedOnDM) {
		return
	}

	if message.Author.Bot {
		return
	}

	// Check if AFK is enabled
	if !cfg.AFK.Enabled || message.Author.ID == session.State.User.ID {
		return
	}

	// Check if AFK is triggered here
	if !cfg.AFK.TriggerList.IsTriggered(session, message.Author.ID, message.GuildID) {
		return
	}

	if a.IsActive() {
		return
	}

	if cfg.AFK.Delay > 0 {
		delay := cfg.AFK.Delay

		// Random delay
		if cfg.AFK.RandomDelay {
			delay = rand.Intn(delay)
		}

		if cfg.AFK.Typing {
			time.Sleep(5 * time.Second)

			duration := time.Duration(delay) * time.Second
			duration = max(duration-time.Second*5, 0)

			session.ChannelTypingDuration(message.ChannelID, time.Second*time.Duration(delay))
		} else {
			time.Sleep(time.Second * time.Duration(delay))
		}

	}

	time.Sleep(2 * time.Second)

	if a.IsActive() || time.Since(a.lastAFKMessage) < 60*time.Second {
		return
	}

	a.lastAFKMessage = time.Now()

	afkMessage := a.GetAFKMessage(message.Message, cfg.AFK.Message)
	if cfg.AFK.AI {
		instructions := utils.Misc.CreateAISystemAFK(afkMessage, session, message.Message)

		message.Content = strings.ReplaceAll(message.Content, session.State.User.Mention(), message.Author.Mention())

		creds := api.LBCreds{
			Username: cfg.Account.LightningBot.Username,
			Password: cfg.Account.LightningBot.Password,
		}

		var err error
		afkMessage, err = api.LBGPT(creds, instructions, message.Content)
		if err != nil {
			logging.Error("AFK.Handler: Error getting AI response: " + err.Error())
			return
		}
	}

	a.eventBroker.Publish(eventbroker.EventLogFeature, &featurelogs.AFKMessage{
		Message:  message.Message,
		Response: afkMessage,
	})

	if cfg.AFK.Reply {
		m, err := session.ChannelMessageSendReply(message.ChannelID, afkMessage, message.Reference())
		if err == nil {
			_ = session.MessageUnack(message.ChannelID, message.ID, 1)
			return
		}

		a.eventBroker.Publish(eventbroker.EventHandleReactionMessage, m)
	} else {
		m, err := session.ChannelMessageSend(message.ChannelID, afkMessage)
		if err == nil {
			_ = session.MessageUnack(message.ChannelID, message.ID, 1)
			return
		}

		a.eventBroker.Publish(eventbroker.EventHandleReactionMessage, m)
	}
}
