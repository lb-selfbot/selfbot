package impl

import (
	"slices"
	"strings"

	"github.com/LightningDev1/LB-Selfbot/internal/api"
	"github.com/LightningDev1/LB-Selfbot/internal/services/interfaces"

	"github.com/LightningDev1/LB-Selfbot/internal/utils"
	"github.com/LightningDev1/discordgo"
)

type MessageReaction struct {
	uiSvc             interfaces.UIService
	eventReactionsSvc interfaces.EventReactionsService

	ReactionType interfaces.EventReactionType
	Guild        *discordgo.Guild
	Channel      *discordgo.Channel
	User         *discordgo.User
	Global       bool

	// Responders
	System   string
	Response string

	// React User
	Reaction string

	// Keyword Notifier
	Keyword string
}

func NewMessageReaction(
	uiSvc interfaces.UIService,
	eventReactionsSvc interfaces.EventReactionsService,
	reactionType interfaces.EventReactionType,
	guild *discordgo.Guild,
	channel *discordgo.Channel,
	user *discordgo.User,
) *MessageReaction {
	return &MessageReaction{
		uiSvc:             uiSvc,
		eventReactionsSvc: eventReactionsSvc,
		ReactionType:      reactionType,
		Guild:             guild,
		Channel:           channel,
		User:              user,
	}
}

func NewMessageReactionFromData(session *discordgo.Session, reactionType interfaces.EventReactionType, data map[string]any) (*MessageReaction, error) {
	var guild *discordgo.Guild
	var channel *discordgo.Channel
	var user *discordgo.User
	var err error

	guildID, ok := data["Guild"].(string)
	if ok {
		guild, err = utils.Discord.GetGuild(session, guildID)
		if err != nil {
			guild = &discordgo.Guild{ID: guildID}
		}
	}

	channelID, ok := data["Channel"].(string)
	if ok {
		channel, err = utils.Discord.GetChannel(session, channelID)
		if err != nil {
			channel = &discordgo.Channel{ID: channelID}
		}
	}

	userID, ok := data["User"].(string)
	if ok {
		user, err = utils.Discord.GetUser(session, userID)
		if err != nil {
			user = &discordgo.User{ID: userID}
		}
	}

	global, ok := data["Global"].(bool)
	if !ok {
		global = false
	}

	system, ok := data["System"].(string)
	if !ok {
		system = ""
	}

	response, ok := data["Response"].(string)
	if !ok {
		response = ""
	}

	reaction, ok := data["Reaction"].(string)
	if !ok {
		reaction = ""
	}

	keyword, ok := data["Keyword"].(string)
	if !ok {
		keyword = ""
	}

	return &MessageReaction{
		ReactionType: reactionType,
		Guild:        guild,
		Channel:      channel,
		User:         user,
		Global:       global,
		System:       system,
		Response:     response,
		Reaction:     reaction,
		Keyword:      keyword,
	}, nil
}

func (m *MessageReaction) Type() interfaces.EventReactionType {
	return m.ReactionType
}

func (m *MessageReaction) Description() string {
	var locations []string

	if m.User != nil {
		locations = append(locations, m.User.Username)
	}
	if m.Channel != nil {
		locations = append(locations, m.Channel.Name)
	}
	if m.Guild != nil {
		locations = append(locations, m.Guild.Name)
	}
	if m.Global {
		locations = append(locations, "Global")
	}
	if m.Reaction != "" {
		locations = append(locations, m.Reaction)
	}
	if m.Keyword != "" {
		locations = append(locations, m.Keyword)
	}

	location := " (" + strings.Join(locations, " - ") + ")"

	switch m.ReactionType {
	case interfaces.EventReactionResponder:
		return "Responder" + location
	case interfaces.EventReactionAIResponder:
		return "AI Responder" + location
	case interfaces.EventReactionCopy:
		return "Copy" + location
	case interfaces.EventReactionReplyCopy:
		return "Reply Copy" + location
	case interfaces.EventReactionMock:
		return "Mock" + location
	case interfaces.EventReactionPinSpam:
		return "Pin Spam" + location
	case interfaces.EventReactionReactSpam:
		return "React Spam" + location
	case interfaces.EventReactionDeleteMessages:
		return "Delete Messages" + location
	case interfaces.EventReactionKeywordNotifier:
		return "Keyword Notifier" + location
	default:
		utils.Logging.Debug("MessageReaction.Description: Unknown reaction type")
		return "MessageReaction reaction" + location
	}
}

func (m *MessageReaction) Equals(other interfaces.EventReaction) bool {
	otherMessageReaction, ok := other.(*MessageReaction)
	if !ok {
		utils.Logging.Debug("MessageReaction.Equals: Other is not a MessageReaction")
		return false
	}

	return m.ReactionType == otherMessageReaction.ReactionType &&
		m.Guild == otherMessageReaction.Guild &&
		m.Channel == otherMessageReaction.Channel &&
		m.Global == otherMessageReaction.Global &&
		m.User == otherMessageReaction.User &&
		m.System == otherMessageReaction.System &&
		m.Response == otherMessageReaction.Response &&
		m.Reaction == otherMessageReaction.Reaction &&
		m.Keyword == otherMessageReaction.Keyword
}

func (m *MessageReaction) Data() map[string]any {
	data := make(map[string]any)

	if m.Guild != nil {
		data["Guild"] = m.Guild.ID
	}
	if m.Channel != nil {
		data["Channel"] = m.Channel.ID
	}
	if m.User != nil {
		data["User"] = m.User.ID
	}
	data["Global"] = m.Global
	data["System"] = m.System
	data["Response"] = m.Response
	data["Reaction"] = m.Reaction
	data["Keyword"] = m.Keyword

	return data
}

func (m *MessageReaction) Handler(session *discordgo.Session, event any) {
	message, ok := event.(*discordgo.MessageCreate)
	if !ok {
		utils.Logging.Debug("MessageReaction.Handler: Event is not a message create event")
		return
	}

	allowedForSelf := []interfaces.EventReactionType{interfaces.EventReactionReactSpam}

	if message.Author.ID == session.State.User.ID && !slices.Contains(allowedForSelf, m.ReactionType) {
		return
	}

	triggered := m.Guild != nil && m.Guild.ID == message.GuildID ||
		m.Channel != nil && m.Channel.ID == message.ChannelID ||
		m.Global

	if !triggered || m.User != nil && m.User.ID != message.Author.ID {
		return
	}

	switch m.ReactionType {
	case interfaces.EventReactionResponder:
		msg, err := session.ChannelMessageSend(message.ChannelID, m.Response)
		m.eventReactionsSvc.HandleReactionMessage(msg, err)
	case interfaces.EventReactionAIResponder:
		messages, err := utils.Discord.GetChannelMessages(session, message.ChannelID, 20)
		if err != nil {
			utils.Logging.Error("Responder.Handler: Error getting channel messages: " + err.Error())
			return
		}

		gptMessages := utils.Misc.CreateAISystemPingMessages(m.System, session, messages)

		response, err := api.GPTCustom(gptMessages)
		if err != nil {
			utils.Logging.Error("Responder.Handler: Error getting AI response: " + response + err.Error())
			return
		}

		msg, err := session.ChannelMessageSend(message.ChannelID, response)
		m.eventReactionsSvc.HandleReactionMessage(msg, err)
	case interfaces.EventReactionCopy:
		msg, err := session.ChannelMessageSendComplex(message.ChannelID, &discordgo.MessageSend{
			Content:         message.Content,
			AllowedMentions: &discordgo.MessageAllowedMentions{}, // Disable mentions
		})
		m.eventReactionsSvc.HandleReactionMessage(msg, err)
	case interfaces.EventReactionReplyCopy:
		msg, err := session.ChannelMessageSendComplex(message.ChannelID, &discordgo.MessageSend{
			Content:         message.Content,
			AllowedMentions: &discordgo.MessageAllowedMentions{RepliedUser: true},
			Reference: &discordgo.MessageReference{
				MessageID: message.ID,
				ChannelID: message.ChannelID,
				GuildID:   message.GuildID,
			},
		})
		m.eventReactionsSvc.HandleReactionMessage(msg, err)
	case interfaces.EventReactionMock:
		msg, err := session.ChannelMessageSendComplex(message.ChannelID, &discordgo.MessageSend{
			Content:         utils.String.Mock(message.Content),
			AllowedMentions: &discordgo.MessageAllowedMentions{}, // Disable mentions
		})
		m.eventReactionsSvc.HandleReactionMessage(msg, err)
	case interfaces.EventReactionPinSpam:
		err := session.ChannelMessagePin(message.ChannelID, message.ID)
		if err == nil {
			m.eventReactionsSvc.AddPinSpamMessage(message.Message)
		}
	case interfaces.EventReactionReactSpam:
		_ = session.MessageReactionAdd(message.ChannelID, message.ID, m.Reaction)
	case interfaces.EventReactionDeleteMessages:
		_ = session.ChannelMessageDelete(message.ChannelID, message.ID)
	case interfaces.EventReactionKeywordNotifier:
		if strings.Contains(strings.ToLower(message.Content), strings.ToLower(m.Keyword)) {
			m.uiSvc.LogKeywordNotifier(message.Message, m.Keyword)
		}
	}
}
