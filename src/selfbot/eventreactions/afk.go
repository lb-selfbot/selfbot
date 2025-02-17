package eventreactions

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/LightningDev1/LB-Selfbot/selfbot/api"
	"github.com/LightningDev1/LB-Selfbot/selfbot/config"

	"github.com/LightningDev1/LB-Selfbot/selfbot/utils"
	"github.com/LightningDev1/LB-Selfbot/selfbot/variables"
	"github.com/LightningDev1/discordgo"
)

type AFK struct{}

func NewAFK() *AFK {
	return &AFK{}
}

func (a *AFK) Type() EventReactionType {
	return AFK_RESPONDER
}

func (a *AFK) Description() string {
	return "AFK Responder (default)"
}

func (a *AFK) Equals(other EventReaction) bool {
	_, ok := other.(*AFK)

	return ok
}

func (a *AFK) HandleLastActivity(session *discordgo.Session, m *discordgo.Message) {
	if m.Author.ID == session.State.User.ID && !IsEventReactionMessage(m) && time.Since(LastAFKMessage) > 5*time.Second {
		LastActivity = time.Now()
	}
}

func (a *AFK) IsActive(cfg *config.Config) bool {
	if !cfg.AFK.MustBeInactive {
		return false
	}

	inactiveTimeout := time.Duration(cfg.AFK.InactiveTimeout) * time.Second

	return time.Since(LastActivity) < inactiveTimeout
}

func (a *AFK) Data() map[string]any {
	return map[string]any{}
}

func (a *AFK) GetAFKMessage(m *discordgo.Message, message string) string {
	message = strings.ReplaceAll(message, "$USER", m.Author.Username)
	message = strings.ReplaceAll(message, "$MENTION", m.Author.Mention())
	message = strings.ReplaceAll(message, "$MSGLINK", utils.Discord.GetMessageLink(m))
	message = strings.ReplaceAll(message, "$MSGCONTENT", m.ContentWithMentionsReplaced())
	message = strings.ReplaceAll(message, "$AFKSINCE", fmt.Sprintf("<t:%d:R>", LastActivity.Unix()))

	message = variables.Parse(message)

	return message
}

func (a *AFK) Handler(session *discordgo.Session, event any) {
	message, ok := event.(*discordgo.MessageCreate)
	if !ok {
		utils.Logging.Debug("AFK.Handler: Event is not a message create event")
		return
	}

	a.HandleLastActivity(session, message.Message)

	cfg, err := config.GetConfig()
	if err != nil {
		utils.Logging.Debug("AFK.Handler: Error getting config: " + err.Error())
		return
	}

	activatedOnPing := cfg.AFK.OnPing && utils.Discord.IsUserMentioned(message.Message, session.State.User.ID)
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

	if a.IsActive(&cfg) {
		return
	}

	if cfg.AFK.Delay > 0 {
		delay := cfg.AFK.Delay

		// Random delay
		if cfg.AFK.RandomDelay {
			delay = rand.Intn(delay)
		}

		if cfg.AFK.Typing {
			utils.Discord.SleepWithTyping(session, message.ChannelID, time.Second*time.Duration(delay))
		} else {
			time.Sleep(time.Second * time.Duration(delay))
		}

	}

	time.Sleep(2 * time.Second)

	if a.IsActive(&cfg) || time.Since(LastAFKMessage) < 60*time.Second {
		return
	}

	LastAFKMessage = time.Now()

	afkMessage := a.GetAFKMessage(message.Message, cfg.AFK.Message)
	if cfg.AFK.AI {
		instructions := utils.Misc.CreateAISystemAFK(afkMessage, session, message.Message)

		message.Content = strings.ReplaceAll(message.Content, session.State.User.Mention(), message.Author.Mention())

		afkMessage, err = api.GPT(instructions, message.Content)
		if err != nil {
			utils.Logging.Error("AFK.Handler: Error getting AI response: " + err.Error())
			return
		}
	}

	logging.Logging.AFKMessageSent(session, message.Message, afkMessage)

	if cfg.AFK.Reply {
		m, err := session.ChannelMessageSendReply(message.ChannelID, afkMessage, message.Reference())
		if err == nil {
			err = session.MessageUnack(message.ChannelID, message.ID, 1)
		}
		HandleMessage(m, err)
	} else {
		m, err := session.ChannelMessageSend(message.ChannelID, afkMessage)
		if err == nil {
			err = session.MessageUnack(message.ChannelID, message.ID, 1)
		}
		HandleMessage(m, err)
	}
}
