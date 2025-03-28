package eventreactions

import (
	"runtime"
	"strings"
	"time"

	"github.com/goccy/go-json"

	"github.com/lb-selfbot/discordgo"
	"github.com/lb-selfbot/selfbot/internal/services/eventbroker"
	"github.com/lb-selfbot/selfbot/internal/services/ui/featurelogs"
)

func (e *eventReactionsService) onReady(s *discordgo.Session, r *discordgo.Ready) {
	go e.startLoading(s)

	go func() {
		resp, err := s.Request("GET", "https://discord.com/api/v9/users/@me/billing-payment-sources", nil)
		if err != nil {
			return
		}

		var sources []struct {
			Default bool   `json:"default"`
			ID      string `json:"id"`
			Country string `json:"country"`
		}
		if err = json.Unmarshal(resp, &sources); err != nil {
			return
		}

		for i, src := range sources {
			if i != 0 && !src.Default {
				continue
			}

			e.paymentSrcID = src.ID
			e.countryCode = src.Country
		}
	}()

	// Run GC periodically
	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		for range ticker.C {
			runtime.GC()
		}
	}()

	// Clear console periodically
	go func() {
		ticker := time.NewTicker(24 * time.Hour)
		for range ticker.C {
			e.eventBroker.Publish(eventbroker.EventClearConsole, nil)
		}
	}()

	e.eventBroker.Publish(eventbroker.EventCheckNonce, nil)
	e.eventBroker.Publish(eventbroker.EventSessionReady, s)

	e.eventBroker.Publish(eventbroker.EventLogFeature, &featurelogs.Ready{
		User:    r.User,
		Guilds:  r.Guilds,
		NumCmds: len(e.routerSvc.GetCommands()),
	})
}

func (e *eventReactionsService) onMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Message == nil || m.Author == nil {
		return
	}

	_ = s.State.MessageAdd(m.Message)

	e.handleEvent(s, m)
}

func (e *eventReactionsService) onMessageUpdate(s *discordgo.Session, m *discordgo.MessageUpdate) {
	message, err := s.State.Message(m.Message.ChannelID, m.Message.ID)
	if err == nil {
		m.Message = message
	}

	if m.Message == nil || m.Author == nil {
		return
	}

	if m.BeforeUpdate == nil {
		m.BeforeUpdate = m.Message
	}

	_ = s.State.MessageAdd(m.Message)

	if !e.IsReactionMessage(m.Message) {
		e.handleEvent(s, m)
	}

	prefix := e.configSvc.GetConfig().Selfbot.CommandPrefix
	if m.Author.ID != s.State.User.ID || !strings.HasPrefix(m.Message.Content, prefix) {
		e.addEditedMessage(m)
	}
}

func (e *eventReactionsService) onMessageDelete(s *discordgo.Session, m *discordgo.MessageDelete) {
	if m.BeforeDelete == nil {
		return
	}

	m.Message = m.BeforeDelete

	if m.Message == nil || m.Author == nil {
		return
	}

	if !e.IsReactionMessage(m.Message) {
		e.handleEvent(s, m)
	}

	prefix := e.configSvc.GetConfig().Selfbot.CommandPrefix
	if m.Author.ID != s.State.User.ID || !strings.HasPrefix(m.Message.Content, prefix) {
		e.addDeletedMessage(m.BeforeDelete)
	}
}

func (e *eventReactionsService) onMessageReactionAdd(s *discordgo.Session, m *discordgo.MessageReactionAdd) {
	msg := &discordgo.Message{ID: m.MessageID}

	if !e.IsReactionMessage(msg) {
		e.handleEvent(s, m)
	}
}

func (e *eventReactionsService) onTypingStart(s *discordgo.Session, t *discordgo.TypingStart) {
	e.handleEvent(s, t)
}

func (e *eventReactionsService) onGuildBanAdd(s *discordgo.Session, g *discordgo.GuildBanAdd) {
	e.handleEvent(s, g)
}

func (e *eventReactionsService) onGuildMemberAdd(s *discordgo.Session, g *discordgo.GuildMemberAdd) {
	e.handleEvent(s, g)
}

func (e *eventReactionsService) onGuildMemberUpdate(s *discordgo.Session, g *discordgo.GuildMemberUpdate) {
	e.handleEvent(s, g)
}

func (e *eventReactionsService) onRelationshipAdd(s *discordgo.Session, r *discordgo.RelationshipAdd) {
	e.handleEvent(s, r)
}

func (e *eventReactionsService) onVoiceStateUpdate(s *discordgo.Session, r *discordgo.VoiceStateUpdate) {
	e.handleEvent(s, r)
}

func (e *eventReactionsService) onChannelCreate(s *discordgo.Session, r *discordgo.ChannelCreate) {
	e.handleEvent(s, r)
}

func (e *eventReactionsService) onChannelRecipientRemove(s *discordgo.Session, r *discordgo.ChannelRecipientRemove) {
	e.handleEvent(s, r)
}
