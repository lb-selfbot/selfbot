package eventreactions

import (
	"encoding/json"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/LightningDev1/LB-Selfbot/internal/data"
	"github.com/LightningDev1/LB-Selfbot/selfbot/config"
	"github.com/LightningDev1/LB-Selfbot/selfbot/ipc"
	"github.com/LightningDev1/LB-Selfbot/selfbot/logging"

	"github.com/LightningDev1/LB-Selfbot/selfbot/reminder"
	"github.com/LightningDev1/LB-Selfbot/selfbot/sessions"
	"github.com/LightningDev1/LB-Selfbot/selfbot/slashcommands"
	"github.com/LightningDev1/LB-Selfbot/selfbot/utils"
	"github.com/LightningDev1/LB-Selfbot/selfbot/websocket"
	"github.com/LightningDev1/dgc"
	"github.com/LightningDev1/discordgo"
	"github.com/bitly/go-simplejson"
)

// Event reaction types

type EventReactionType int

const (
	// Ignore reactions
	IGNORE_USER  EventReactionType = 0
	IGNORE_GUILD EventReactionType = 1

	// Ping reactions
	PING_PING        EventReactionType = 2
	PING_KICK        EventReactionType = 3
	PING_BAN         EventReactionType = 4
	PING_MUTE        EventReactionType = 5
	PING_TIMEOUT     EventReactionType = 6
	PING_INSULT      EventReactionType = 7
	PING_REACT       EventReactionType = 8
	PING_RESPONSE    EventReactionType = 9
	PING_AI_RESPONSE EventReactionType = 10

	// Mirror reactions
	MIRROR_GUILD   EventReactionType = 11
	MIRROR_CHANNEL EventReactionType = 12

	// Responder reactions
	RESPONDER    EventReactionType = 13
	AI_RESPONDER EventReactionType = 14

	// Message reactions
	COPY             EventReactionType = 15
	REPLY_COPY       EventReactionType = 16
	MOCK             EventReactionType = 17
	PIN_SPAM         EventReactionType = 18
	REACT_SPAM       EventReactionType = 19
	DELETE_MESSAGES  EventReactionType = 20
	KEYWORD_NOTIFIER EventReactionType = 21

	// Snipers
	SNIPERS EventReactionType = 22

	// AFK
	AFK_RESPONDER EventReactionType = 23

	// Phishing
	ANTI_PHISHING EventReactionType = 24

	// Protections
	PROTECTIONS EventReactionType = 25

	// Notifiers
	NOTIFIERS EventReactionType = 26

	// Selfbot detection
	SELFBOT_DETECTION EventReactionType = 27

	// Remove reactions
	REMOVE_REACTIONS EventReactionType = 28

	// Force disconnect
	FORCE_DISCONNECT EventReactionType = 29

	// Force nick
	FORCE_NICK EventReactionType = 30

	// Force kick
	FORCE_KICK EventReactionType = 31

	// No leave
	NO_LEAVE EventReactionType = 32

	// Anti group chat
	ANTI_GROUP_CHAT EventReactionType = 33

	// Anti delete messages
	ANTI_DELETE_MESSAGES EventReactionType = 34
)

func defaultEventReactions() []EventReaction {
	return []EventReaction{
		NewSnipers(nil, ""),
		NewAFK(),
		NewPhishing(),
		NewNotifiers(),
		NewProtections(),
		NewSelfbotDetection(),
	}
}

// ReactionsMutex is a mutex for the list of active event reactions
var ReactionsMutex sync.Mutex

// ActiveReactions is a list of currently active event reactions
var ActiveReactions = []EventReaction{}

// PinSpamMessages is a list of messages that are pinned by the pin spam reaction
var PinSpamMessages = utils.NewList[*discordgo.Message]()

// EventReactionMessages is a list of messages that are created by event reactions
var EventReactionMessages = utils.NewList[*discordgo.Message]()

// EditedMessages is a map of messages that are edited by users in a channel
var EditedMessages = sync.Map{}

// DeletedMessages is a map of messages that are deleted by users in a channel
var DeletedMessages = sync.Map{}

// UsedCommands is a list of commands that have been used
var UsedCommands = utils.NewList[*dgc.Ctx]()

// TimeDifference is the time difference between the local time and the Discord time
var TimeDifference time.Duration = 0

// LastActivity is the last time the user was active (used for inactive timeout)
var LastActivity = time.Now()

// LastAFKMessage is the last time an AFK message was sent
var LastAFKMessage = time.Now()

// EventReaction is an interface for event reactions
type EventReaction interface {
	Type() EventReactionType
	Description() string
	Equals(other EventReaction) bool
	Data() map[string]any
	Handler(session *discordgo.Session, event any)
}

// Persistence
func StartLoading(session *discordgo.Session, router *dgc.Router) {
	err := Load(session, true)
	if err != nil {
		utils.Logging.Error("Failed to load event reactions:", err.Error())
	}

	for {
		time.Sleep(time.Second * 5)

		err = Load(session, false)
		if err != nil {
			utils.Logging.Error("Failed to load event reactions:", err.Error())
		}

		// Load command aliases
		err = LoadAliases(router)
		if err != nil {
			utils.Logging.Error("Failed to load command aliases:", err.Error())
		}
	}
}
func Load(session *discordgo.Session, first bool) error {
	ReactionsMutex.Lock()
	defer ReactionsMutex.Unlock()

	defer func() {
		if r := recover(); r != nil {
			cfg, err := config.GetConfig()
			if err != nil {
				return
			}

			cfg.EventReactions = []byte("[]")

			err = cfg.Save()
			if err != nil {
				return
			}

			utils.Logging.Error("An error occurred while loading event reactions. Event reactions have been reset.")
		}
	}()

	cfg, err := config.GetConfig()
	if err != nil {
		return err
	}

	var reactions []struct {
		Type EventReactionType `json:"Type"`
		Data map[string]any    `json:"Data"`
	}
	_ = json.Unmarshal(cfg.EventReactions, &reactions)

	ActiveReactions = defaultEventReactions()
	oldLength := len(ActiveReactions)

	for _, reaction := range reactions {
		switch reaction.Type {
		case IGNORE_USER, IGNORE_GUILD:
			ignore, err := NewIgnoreFromData(session, reaction.Data)
			if err == nil {
				ActiveReactions = append(ActiveReactions, ignore)
			}
		case PING_PING, PING_KICK, PING_BAN, PING_MUTE, PING_TIMEOUT, PING_INSULT, PING_REACT, PING_RESPONSE, PING_AI_RESPONSE:
			ping, err := NewPingFromData(session, reaction.Type, reaction.Data)
			if err == nil {
				ActiveReactions = append(ActiveReactions, ping)
			}
		case MIRROR_GUILD, MIRROR_CHANNEL:
			mirror, err := NewMirrorFromData(session, reaction.Data)
			if err == nil {
				ActiveReactions = append(ActiveReactions, mirror)
			}
		case RESPONDER, AI_RESPONDER, COPY, REPLY_COPY, MOCK, PIN_SPAM, REACT_SPAM, DELETE_MESSAGES, KEYWORD_NOTIFIER:
			responder, err := NewMessageReactionFromData(session, reaction.Type, reaction.Data)
			if err == nil {
				ActiveReactions = append(ActiveReactions, responder)
			}
		case REMOVE_REACTIONS:
			remove, err := NewRemoveReactionsFromData(session, reaction.Data)
			if err == nil {
				ActiveReactions = append(ActiveReactions, remove)
			}
		case FORCE_DISCONNECT:
			forceDisconnect, err := NewForceDisconnectFromData(session, reaction.Data)
			if err == nil {
				ActiveReactions = append(ActiveReactions, forceDisconnect)
			}
		case FORCE_NICK:
			forceNick, err := NewForceNickFromData(session, reaction.Data)
			if err == nil {
				ActiveReactions = append(ActiveReactions, forceNick)
			}
		case FORCE_KICK:
			forceKick, err := NewForceKickFromData(session, reaction.Data)
			if err == nil {
				ActiveReactions = append(ActiveReactions, forceKick)
			}
		case NO_LEAVE:
			noLeave, err := NewNoLeaveFromData(session, reaction.Data)
			if err == nil {
				ActiveReactions = append(ActiveReactions, noLeave)
			}
		case ANTI_GROUP_CHAT:
			antiGroupChat, err := NewAntiGroupChatFromData(session, reaction.Data)
			if err == nil {
				ActiveReactions = append(ActiveReactions, antiGroupChat)
			}
		case ANTI_DELETE_MESSAGES:
			antiDeleteMessages, err := NewAntiDeleteMessagesFromData(session, reaction.Data)
			if err == nil {
				ActiveReactions = append(ActiveReactions, antiDeleteMessages)
			}
		}
	}

	new := len(ActiveReactions) - oldLength
	if first && new > 0 {
		logging.Logging.EventReactions(new)
	}

	return nil
}

func LoadAliases(router *dgc.Router) error {
	cfg, err := config.GetConfig()
	if err != nil {
		return err
	}

	aliases := make(map[string][]string)

	for alias, cmd := range cfg.Selfbot.CommandAliases {
		if aliases[cmd] == nil {
			aliases[cmd] = []string{}
		}

		aliases[cmd] = append(aliases[cmd], alias)
	}

	for _, cmd := range cfg.Selfbot.CommandAliases {
		if cmd := router.GetCmd(cmd); cmd != nil {
			cmd.Aliases = aliases[cmd.Name]
		}
	}

	return nil
}

func Save() error {
	var reactions []any

	for _, reaction := range ActiveReactions {
		if strings.Contains(reaction.Description(), "default") {
			continue
		}

		reactions = append(reactions, map[string]any{
			"Type": reaction.Type(),
			"Data": reaction.Data(),
		})
	}

	cfg, err := config.GetConfig()
	if err != nil {
		return err
	}

	reactionsBytes, err := json.MarshalIndent(reactions, "", "  ")
	if err != nil {
		return err
	}

	cfg.EventReactions = reactionsBytes

	return cfg.Save()
}

// Exists checks if an event reaction already exists
func Exists(reaction EventReaction) bool {
	ReactionsMutex.Lock()
	defer ReactionsMutex.Unlock()

	for _, r := range ActiveReactions {
		if r.Equals(reaction) {
			return true
		}
	}

	return false
}

// Add adds an event reaction to the list of active event reactions
func Add(reaction EventReaction) {
	if Exists(reaction) {
		return
	}

	ReactionsMutex.Lock()
	defer ReactionsMutex.Unlock()

	ActiveReactions = append(ActiveReactions, reaction)

	_ = Save()
}

// GetAll gets all event reactions of a certain type
func GetAll(reactionType EventReactionType) (reactions []EventReaction) {
	ReactionsMutex.Lock()
	defer ReactionsMutex.Unlock()

	for _, r := range ActiveReactions {
		if r.Type() == reactionType {
			reactions = append(reactions, r)
		}
	}

	return
}

// Remove removes an event reaction from the list of active event reactions
func Remove(reaction EventReaction) {
	ReactionsMutex.Lock()
	defer ReactionsMutex.Unlock()

	for i, r := range ActiveReactions {
		if r.Equals(reaction) {
			ActiveReactions = append(ActiveReactions[:i], ActiveReactions[i+1:]...)
		}
	}

	_ = Save()
}

// RemoveAll removes all event reactions of a certain type
func RemoveAll(reactionType EventReactionType) {
	ReactionsMutex.Lock()
	defer ReactionsMutex.Unlock()

	for i := len(ActiveReactions) - 1; i >= 0; i-- {
		if ActiveReactions[i].Type() == reactionType {
			ActiveReactions = append(ActiveReactions[:i], ActiveReactions[i+1:]...)
		}
	}

	_ = Save()
}

// Clear clears all event reactions
func Clear() {
	ReactionsMutex.Lock()
	defer ReactionsMutex.Unlock()

	ActiveReactions = defaultEventReactions()

	_ = Save()
}

// HandleMessage adds a message to the list of event reaction messages
func HandleMessage(message *discordgo.Message, err error) {
	if err != nil {
		return
	}

	EventReactionMessages.Mutex.Lock()
	defer EventReactionMessages.Mutex.Unlock()

	EventReactionMessages.Add(message)
}

// IsEventReactionMessage checks if a message is an event reaction message
func IsEventReactionMessage(message *discordgo.Message) bool {
	EventReactionMessages.Mutex.Lock()
	defer EventReactionMessages.Mutex.Unlock()

	for _, m := range EventReactionMessages.Items {
		if m.ID == message.ID {
			return true
		}
	}

	return false
}

func AddEditedMessage(message *discordgo.MessageUpdate) {
	messages, _ := EditedMessages.LoadOrStore(message.ChannelID, []*discordgo.MessageUpdate{})

	messages = append(messages.([]*discordgo.MessageUpdate), message)

	EditedMessages.Store(message.ChannelID, messages)
}

func GetEditedMessages(channelID string) []*discordgo.MessageUpdate {
	messages, _ := EditedMessages.LoadOrStore(channelID, []*discordgo.MessageUpdate{})

	return messages.([]*discordgo.MessageUpdate)
}

func AddDeletedMessage(message *discordgo.Message) {
	messages, _ := DeletedMessages.LoadOrStore(message.ChannelID, []*discordgo.Message{})

	messages = append(messages.([]*discordgo.Message), message)

	DeletedMessages.Store(message.ChannelID, messages)
}

func GetDeletedMessages(channelID string) []*discordgo.Message {
	messages, _ := DeletedMessages.LoadOrStore(channelID, []*discordgo.Message{})

	return messages.([]*discordgo.Message)
}

func HandleEvent(session *discordgo.Session, event any) {
	ReactionsMutex.Lock()
	defer ReactionsMutex.Unlock()

	for _, r := range ActiveReactions {
		go r.Handler(session, event)
	}
}

func RegisterHandlers(session *discordgo.Session, router *dgc.Router) {
	// Ready Event
	session.AddHandler(func(_ *discordgo.Session, ready *discordgo.Ready) {
		runtime.GC()

		// Start the session handler
		sessions.RegisterSessionHandler(session)

		// Start slash command jobs
		go slashcommands.ExecuteJobs(session)

		// Start reminders
		go reminder.StartReminder()

		// Start event reactions
		go StartLoading(session, router)

		// Start RPC
		// rpc.StartRPC(session)

		// Setup snipers
		go func() {
			resp, err := session.Request("GET", "https://discord.com/api/v9/users/@me/billing/payment-sources", nil)
			if err != nil {
				utils.Logging.Error("Error getting payment sources: ", err)
				return
			}

			json, err := simplejson.NewJson(resp)
			if err != nil {
				utils.Logging.Error("Error parsing payment sources: ", err)
				return
			}

			// No payment sources
			if len(json.MustArray()) == 0 {
				return
			}

			for i := 0; i < len(json.MustArray()); i++ {
				source := json.GetIndex(i)

				if i == 0 || source.Get("default").MustBool() {
					data.PaymentSourceID = source.Get("id").MustString()
					data.CountryCode = source.Get("country").MustString()
				}
			}
		}()

		// GC
		gcTicker := time.NewTicker(time.Second * 10)
		go func() {
			for range gcTicker.C {
				runtime.GC()
			}
		}()

		// Console clear
		clsTicker := time.NewTicker(time.Hour * 24)
		go func() {
			for range clsTicker.C {
				if len(ipc.HTMLConsole) < 1e5 {
					continue
				}

				logging.UI.Clear()
			}
		}()

		data.Token = session.Identify.Token

		logging.Logging.OnReady(router, ready)
	})

	// Message Events
	session.AddHandler(func(_ *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Message == nil || m.Author == nil {
			return
		}

		_ = session.State.MessageAdd(m.Message)

		HandleEvent(session, m)

		TimeDifference = time.Since(utils.Discord.SnowflakeToTime(m.ID))

		websocket.Websocket.CheckNonce()
	})

	session.AddHandler(func(_ *discordgo.Session, m *discordgo.MessageUpdate) {
		message, err := session.State.Message(m.Message.ChannelID, m.Message.ID)
		if err == nil {
			m.Message = message
		}

		if m.Message == nil || m.Author == nil {
			return
		}

		if m.BeforeUpdate == nil {
			m.BeforeUpdate = m.Message
		}

		_ = session.State.MessageAdd(m.Message)

		if !IsEventReactionMessage(m.Message) {
			HandleEvent(session, m)
		}

		// Add edited message
		cfg, err := config.GetConfig()
		if err != nil {
			return
		}

		if m.Author.ID != session.State.User.ID || !strings.HasPrefix(m.Message.Content, cfg.Selfbot.CommandPrefix) {
			AddEditedMessage(m)
		}

		websocket.Websocket.CheckNonce()
	})

	session.AddHandler(func(_ *discordgo.Session, m *discordgo.MessageDelete) {
		if m.BeforeDelete == nil {
			return
		}

		m.Message = m.BeforeDelete

		if m.Message == nil || m.Author == nil {
			return
		}

		if !IsEventReactionMessage(m.Message) {
			HandleEvent(session, m)
		}

		// Add deleted message
		cfg, err := config.GetConfig()
		if err != nil {
			return
		}

		if m.Author.ID != session.State.User.ID || !strings.HasPrefix(m.Message.Content, cfg.Selfbot.CommandPrefix) {
			AddDeletedMessage(m.BeforeDelete)
		}

		websocket.Websocket.CheckNonce()
	})

	session.AddHandler(func(_ *discordgo.Session, m *discordgo.MessageReactionAdd) {
		msg := &discordgo.Message{ID: m.MessageID}

		if !IsEventReactionMessage(msg) {
			HandleEvent(session, m)
		}

		websocket.Websocket.CheckNonce()
	})

	session.AddHandler(func(_ *discordgo.Session, t *discordgo.TypingStart) {
		HandleEvent(session, t)

		websocket.Websocket.CheckNonce()
	})

	// Guild Events
	session.AddHandler(func(_ *discordgo.Session, g *discordgo.GuildBanAdd) {
		HandleEvent(session, g)

		websocket.Websocket.CheckNonce()
	})

	session.AddHandler(func(_ *discordgo.Session, g *discordgo.GuildMemberAdd) {
		HandleEvent(session, g)

		websocket.Websocket.CheckNonce()
	})

	session.AddHandler(func(_ *discordgo.Session, g *discordgo.GuildMemberUpdate) {
		HandleEvent(session, g)

		websocket.Websocket.CheckNonce()
	})

	// User Events
	session.AddHandler(func(_ *discordgo.Session, r *discordgo.RelationshipAdd) {
		HandleEvent(session, r)

		websocket.Websocket.CheckNonce()
	})

	// Channel Events
	session.AddHandler(func(_ *discordgo.Session, r *discordgo.VoiceStateUpdate) {
		HandleEvent(session, r)

		websocket.Websocket.CheckNonce()
	})

	session.AddHandler(func(_ *discordgo.Session, r *discordgo.ChannelCreate) {
		HandleEvent(session, r)

		websocket.Websocket.CheckNonce()
	})

	session.AddHandler(func(_ *discordgo.Session, r *discordgo.ChannelRecipientRemove) {
		HandleEvent(session, r)

		websocket.Websocket.CheckNonce()
	})

}
