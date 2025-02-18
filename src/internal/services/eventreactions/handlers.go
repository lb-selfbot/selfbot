package eventreactions

import (
	"encoding/json"
	"runtime"
	"strings"
	"time"

	"github.com/LightningDev1/LB-Selfbot/internal/services/ui"
	"github.com/LightningDev1/LB-Selfbot/selfbot/utils"
	"github.com/LightningDev1/discordgo"
)

func (e *EventReactionsServiceImpl) onReady(s *discordgo.Session, r *discordgo.Ready) {
	e.sessionsSvc.Register()
	e.scJobsSvc.Start()
	e.rpcSvc.Start()

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
			_ = e.uiSvc.ClearConsole()
			_ = e.uiSvc.ShowBanner(ui.BannerUI)
		}
	}()
}

func (e *EventReactionsServiceImpl) onMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Message == nil || m.Author == nil {
		return
	}

	_ = s.State.MessageAdd(m.Message)

	e.handleEvent(s, m)

	e.timeDiff = time.Since(utils.Discord.SnowflakeToTime(m.ID))
}

func (e *EventReactionsServiceImpl) onMessageUpdate(s *discordgo.Session, m *discordgo.MessageUpdate) {
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

	if !e.isEventReactionMessage(m.Message) {
		e.handleEvent(s, m)
	}

	if m.Author.ID != s.State.User.ID || !strings.HasPrefix(m.Message.Content, e.configSvc.GetPrefix()) {
		e.addEditedMessage(m)
	}
}

func (e *EventReactionsServiceImpl) onMessageDelete(s *discordgo.Session, m *discordgo.MessageDelete) {
	if m.BeforeDelete == nil {
		return
	}

	m.Message = m.BeforeDelete

	if m.Message == nil || m.Author == nil {
		return
	}

	if !e.isEventReactionMessage(m.Message) {
		e.handleEvent(s, m)
	}

	if m.Author.ID != s.State.User.ID || !strings.HasPrefix(m.Message.Content, e.configSvc.GetPrefix()) {
		e.addDeletedMessage(m.BeforeDelete)
	}
}

func (e *EventReactionsServiceImpl) onMessageReactionAdd(s *discordgo.Session, m *discordgo.MessageReactionAdd) {
	msg := &discordgo.Message{ID: m.MessageID}

	if !e.isEventReactionMessage(msg) {
		e.handleEvent(s, m)
	}
}

func (e *EventReactionsServiceImpl) onTypingStart(s *discordgo.Session, t *discordgo.TypingStart) {
	e.handleEvent(s, t)
}

func (e *EventReactionsServiceImpl) onGuildBanAdd(s *discordgo.Session, g *discordgo.GuildBanAdd) {
	e.handleEvent(s, g)
}

func (e *EventReactionsServiceImpl) onGuildMemberAdd(s *discordgo.Session, g *discordgo.GuildMemberAdd) {
	e.handleEvent(s, g)
}

func (e *EventReactionsServiceImpl) onGuildMemberUpdate(s *discordgo.Session, g *discordgo.GuildMemberUpdate) {
	e.handleEvent(s, g)
}

func (e *EventReactionsServiceImpl) onRelationshipAdd(s *discordgo.Session, r *discordgo.RelationshipAdd) {
	e.handleEvent(s, r)
}

func (e *EventReactionsServiceImpl) onVoiceStateUpdate(s *discordgo.Session, r *discordgo.VoiceStateUpdate) {
	e.handleEvent(s, r)
}

func (e *EventReactionsServiceImpl) onChannelCreate(s *discordgo.Session, r *discordgo.ChannelCreate) {
	e.handleEvent(s, r)
}

func (e *EventReactionsServiceImpl) onChannelRecipientRemove(s *discordgo.Session, r *discordgo.ChannelRecipientRemove) {
	e.handleEvent(s, r)
}
