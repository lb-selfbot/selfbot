package eventreactions

import (
	"slices"
	"sync"
	"time"

	"github.com/lb-selfbot/discordgo"
	"github.com/lb-selfbot/selfbot/internal/services/config"
	"github.com/lb-selfbot/selfbot/internal/services/eventbroker"
	"github.com/lb-selfbot/selfbot/internal/services/eventreactions/types"
	"github.com/lb-selfbot/selfbot/internal/services/router"
	"github.com/lb-selfbot/selfbot/internal/utils"
	"github.com/lb-selfbot/selfbot/internal/utils2"
)

// Service handles event reactions.
type Service interface {
	// Register registers the event handlers.
	Register()

	// GetAllReactions returns all event reactions.
	GetAllReactions() []types.EventReaction

	// GetReactionsByType returns event reactions by type.
	GetReactionsByType(t types.EventReactionType) []types.EventReaction

	// AddReaction adds a new event reaction.
	AddReaction(r types.EventReaction)

	// RemoveReaction removes an event reaction.
	RemoveReaction(r types.EventReaction)

	// RemoveAllReactions removes all event reactions.
	RemoveAllReactions()

	// RemoveReactionsByType removes event reactions by type.
	RemoveReactionsByType(t types.EventReactionType)

	// IsReactionMessage checks if a message is a reaction message.
	IsReactionMessage(message *discordgo.Message) bool

	// HandleReactionMessage handles a reaction message.
	HandleReactionMessage(message *discordgo.Message, err error)

	// LoadCmdAliases loads command aliases.
	LoadCmdAliases() error

	// GetPinSpamMessages returns pin spam messages.
	GetPinSpamMessages() []*discordgo.Message

	// AddPinSpamMessage adds a pin spam message.
	AddPinSpamMessage(message *discordgo.Message)

	// ClearPinSpamMessages removes pin spam messages.
	ClearPinSpamMessages()

	// GetEditedMessages returns edited messages by channel ID.
	GetEditedMessages(channelID string) []*discordgo.MessageUpdate

	// GetDeletedMessages returns deleted messages by channel ID.
	GetDeletedMessages(channelID string) []*discordgo.Message

	// GetUsedCommands returns used commands.
	GetUsedCommands() []*router.Ctx

	// AddUsedCommand adds a used command.
	AddUsedCommand(ctx *router.Ctx)

	// GetLastActivity returns the last activity time.
	GetLastActivity() time.Time

	// SetLastActivity sets the last activity time.
	SetLastActivity()

	// GetPaymentSrcID returns the payment source ID.
	GetPaymentSrcID() string

	// GetCountryCode returns the country code.
	GetCountryCode() string
}

type eventReactionsService struct {
	eventBroker eventbroker.Service
	configSvc   config.Service
	routerSvc   router.RouterService

	reactions   []types.EventReaction
	reactionsMu sync.Mutex

	pinSpamMsgs   []*discordgo.Message
	pinSpamMsgsMu sync.Mutex

	eventReactionsMsgs   []*discordgo.Message
	eventReactionsMsgsMu sync.Mutex

	editedMsgs   map[string][]*discordgo.MessageUpdate
	editedMsgsMu sync.Mutex

	deletedMsgs   map[string][]*discordgo.Message
	deletedMsgsMu sync.Mutex

	usedCmds   []*router.Ctx
	usedCmdsMu sync.Mutex

	lastActivity time.Time
	lastAFKMsg   time.Time

	paymentSrcID string
	countryCode  string
}

func NewService(eventBroker eventbroker.Service, configService config.Service, routerSvc router.RouterService) *eventReactionsService {
	svc := &eventReactionsService{
		eventBroker:        eventBroker,
		configSvc:          configService,
		routerSvc:          routerSvc,
		reactions:          make([]types.EventReaction, 0),
		pinSpamMsgs:        make([]*discordgo.Message, 0),
		eventReactionsMsgs: make([]*discordgo.Message, 0),
		editedMsgs:         make(map[string][]*discordgo.MessageUpdate),
		deletedMsgs:        make(map[string][]*discordgo.Message),
		usedCmds:           make([]*router.Ctx, 0),
		lastActivity:       time.Now(),
		lastAFKMsg:         time.Now(),
	}
	svc.subscribe()
	return svc
}

func (e *eventReactionsService) subscribe() {
	e.eventBroker.Subscribe(eventbroker.EventSessionReady, func(payload any) {
		session, ok := payload.(*discordgo.Session)
		if !ok {
			return
		}

		e.register(session)
	})

	e.eventBroker.Subscribe(eventbroker.EventHandleReactionMessage, func(payload any) {
		msg, ok := payload.(*discordgo.Message)
		if !ok {
			return
		}

		e.HandleReactionMessage(msg)
	})

	e.eventBroker.Subscribe(eventbroker.EventHandlePinSpamMessage, func(payload any) {
		msg, ok := payload.(*discordgo.Message)
		if !ok {
			return
		}

		e.AddPinSpamMessage(msg)
	})
}

func (e *eventReactionsService) register(session *discordgo.Session) {
	session.AddHandler(func(_ *discordgo.Session, _ any) {
		e.eventBroker.Publish(eventbroker.EventCheckNonce, nil)
	})

	session.AddHandler(e.onReady)
	session.AddHandler(e.onMessageCreate)
	session.AddHandler(e.onMessageUpdate)
	session.AddHandler(e.onMessageDelete)
	session.AddHandler(e.onMessageReactionAdd)
	session.AddHandler(e.onTypingStart)
	session.AddHandler(e.onGuildBanAdd)
	session.AddHandler(e.onGuildMemberAdd)
	session.AddHandler(e.onGuildMemberUpdate)
	session.AddHandler(e.onRelationshipAdd)
	session.AddHandler(e.onVoiceStateUpdate)
	session.AddHandler(e.onChannelCreate)
	session.AddHandler(e.onChannelRecipientRemove)
}

func (e *eventReactionsService) handleEvent(s *discordgo.Session, event any) {
	e.reactionsMu.Lock()
	defer e.reactionsMu.Unlock()

	for _, r := range e.reactions {
		go r.Handler(s, event)
	}
}

func (e *eventReactionsService) addEditedMessage(message *discordgo.MessageUpdate) {
	e.editedMsgsMu.Lock()
	defer e.editedMsgsMu.Unlock()

	messages, ok := e.editedMsgs[message.ChannelID]
	if !ok {
		messages = []*discordgo.MessageUpdate{}
	}
	e.editedMsgs[message.ChannelID] = append(messages, message)
}

func (e *eventReactionsService) addDeletedMessage(message *discordgo.Message) {
	e.deletedMsgsMu.Lock()
	defer e.deletedMsgsMu.Unlock()

	messages, ok := e.deletedMsgs[message.ChannelID]
	if !ok {
		messages = []*discordgo.Message{}
	}
	e.deletedMsgs[message.ChannelID] = append(messages, message)
}

func (e *eventReactionsService) GetAllReactions() []types.EventReaction {
	return e.reactions
}

func (e *eventReactionsService) GetReactionsByType(t types.EventReactionType) []types.EventReaction {
	return utils.Filter(e.reactions, func(r types.EventReaction) bool {
		return r.Type() == t
	})
}

func (e *eventReactionsService) AddReaction(r types.EventReaction) {
	e.reactionsMu.Lock()
	defer e.reactionsMu.Unlock()

	e.reactions = append(e.reactions, r)

	_ = e.saveEventReactions()
}

func (e *eventReactionsService) RemoveReaction(r types.EventReaction) {
	e.reactionsMu.Lock()
	defer e.reactionsMu.Unlock()

	e.reactions = utils.Filter(e.reactions, func(rr types.EventReaction) bool {
		return !rr.Equals(r)
	})

	_ = e.saveEventReactions()
}

func (e *eventReactionsService) RemoveAllReactions() {
	e.reactionsMu.Lock()
	defer e.reactionsMu.Unlock()

	e.reactions = utils2.Filter(e.reactions, func(r types.EventReaction) bool {
		return !slices.Contains(types.DefaultReactions, r.Type())
	})

	_ = e.saveEventReactions()
}

func (e *eventReactionsService) RemoveReactionsByType(t types.EventReactionType) {
	e.reactionsMu.Lock()
	defer e.reactionsMu.Unlock()

	e.reactions = utils.Filter(e.reactions, func(r types.EventReaction) bool {
		return r.Type() != t
	})

	_ = e.saveEventReactions()
}

func (e *eventReactionsService) IsReactionMessage(message *discordgo.Message) bool {
	for _, m := range e.eventReactionsMsgs {
		if m.ID == message.ID {
			return true
		}
	}

	return false
}

func (e *eventReactionsService) HandleReactionMessage(message *discordgo.Message) {
	if message == nil {
		return
	}

	e.eventReactionsMsgsMu.Lock()
	defer e.eventReactionsMsgsMu.Unlock()

	e.eventReactionsMsgs = append(e.eventReactionsMsgs, message)
}

func (e *eventReactionsService) GetPinSpamMessages() []*discordgo.Message {
	e.pinSpamMsgsMu.Lock()
	defer e.pinSpamMsgsMu.Unlock()

	return e.pinSpamMsgs
}

func (e *eventReactionsService) AddPinSpamMessage(message *discordgo.Message) {
	e.pinSpamMsgsMu.Lock()
	defer e.pinSpamMsgsMu.Unlock()

	e.pinSpamMsgs = append(e.pinSpamMsgs, message)
}

func (e *eventReactionsService) ClearPinSpamMessages() {
	e.pinSpamMsgsMu.Lock()
	defer e.pinSpamMsgsMu.Unlock()

	e.pinSpamMsgs = []*discordgo.Message{}
}

func (e *eventReactionsService) GetEditedMessages(channelID string) []*discordgo.MessageUpdate {
	e.editedMsgsMu.Lock()
	defer e.editedMsgsMu.Unlock()

	return e.editedMsgs[channelID]
}

func (e *eventReactionsService) GetDeletedMessages(channelID string) []*discordgo.Message {
	e.deletedMsgsMu.Lock()
	defer e.deletedMsgsMu.Unlock()

	return e.deletedMsgs[channelID]
}

func (e *eventReactionsService) GetUsedCommands() []*router.Ctx { return e.usedCmds }
func (e *eventReactionsService) AddUsedCommand(ctx *router.Ctx) {
	e.usedCmdsMu.Lock()
	defer e.usedCmdsMu.Unlock()

	e.usedCmds = append(e.usedCmds, ctx)
}

func (e *eventReactionsService) GetLastActivity() time.Time { return e.lastActivity }
func (e *eventReactionsService) SetLastActivity()           { e.lastActivity = time.Now() }

func (e *eventReactionsService) GetPaymentSrcID() string { return e.paymentSrcID }
func (e *eventReactionsService) GetCountryCode() string  { return e.countryCode }
