package eventreactions

import (
	"slices"
	"strings"

	"github.com/LightningDev1/LB-Selfbot/selfbot/api"
	"github.com/LightningDev1/LB-Selfbot/selfbot/logging"

	"github.com/LightningDev1/LB-Selfbot/selfbot/utils"
	"github.com/LightningDev1/discordgo"
)

type MessageReaction struct {
	ReactionType EventReactionType
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

func NewMessageReaction(rType EventReactionType, g *discordgo.Guild, c *discordgo.Channel, u *discordgo.User) *MessageReaction {
	return &MessageReaction{ReactionType: rType, Guild: g, Channel: c, User: u}
}

func NewMessageReactionFromData(session *discordgo.Session, reactionType EventReactionType, data map[string]any) (*MessageReaction, error) {
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
		user, ok = utils.Discord.GetUser(session, userID)
		if !ok {
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

func (m *MessageReaction) Type() EventReactionType {
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
	case RESPONDER:
		return "Responder" + location
	case AI_RESPONDER:
		return "AI Responder" + location
	case COPY:
		return "Copy" + location
	case REPLY_COPY:
		return "Reply Copy" + location
	case MOCK:
		return "Mock" + location
	case PIN_SPAM:
		return "Pin Spam" + location
	case REACT_SPAM:
		return "React Spam" + location
	case DELETE_MESSAGES:
		return "Delete Messages" + location
	case KEYWORD_NOTIFIER:
		return "Keyword Notifier" + location
	default:
		utils.Logging.Debug("MessageReaction.Description: Unknown reaction type")
		return "MessageReaction reaction" + location
	}
}

func (m *MessageReaction) Equals(other EventReaction) bool {
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

	allowedForSelf := []EventReactionType{REACT_SPAM}

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
	case RESPONDER:
		m, err := session.ChannelMessageSend(message.ChannelID, m.Response)
		HandleMessage(m, err)
	case AI_RESPONDER:
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

		m, err := session.ChannelMessageSend(message.ChannelID, response)
		HandleMessage(m, err)
	case COPY:
		m, err := session.ChannelMessageSendComplex(message.ChannelID, &discordgo.MessageSend{
			Content:         message.Content,
			AllowedMentions: &discordgo.MessageAllowedMentions{}, // Disable mentions
		})
		HandleMessage(m, err)
	case REPLY_COPY:
		m, err := session.ChannelMessageSendComplex(message.ChannelID, &discordgo.MessageSend{
			Content:         message.Content,
			AllowedMentions: &discordgo.MessageAllowedMentions{RepliedUser: true},
			Reference: &discordgo.MessageReference{
				MessageID: message.ID,
				ChannelID: message.ChannelID,
				GuildID:   message.GuildID,
			},
		})
		HandleMessage(m, err)
	case MOCK:
		m, err := session.ChannelMessageSendComplex(message.ChannelID, &discordgo.MessageSend{
			Content:         utils.String.Mock(message.Content),
			AllowedMentions: &discordgo.MessageAllowedMentions{}, // Disable mentions
		})
		HandleMessage(m, err)
	case PIN_SPAM:
		err := session.ChannelMessagePin(message.ChannelID, message.ID)
		if err == nil {
			PinSpamMessages.Add(message.Message)
		}
	case REACT_SPAM:
		_ = session.MessageReactionAdd(message.ChannelID, message.ID, m.Reaction)
	case DELETE_MESSAGES:
		_ = session.ChannelMessageDelete(message.ChannelID, message.ID)
	case KEYWORD_NOTIFIER:
		if strings.Contains(strings.ToLower(message.Content), strings.ToLower(m.Keyword)) {
			logging.Logging.KeywordNotifier(session, message.Message, m.Keyword)
		}
	}
}
