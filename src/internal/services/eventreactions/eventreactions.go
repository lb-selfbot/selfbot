package eventreactions

import (
	"sync"
	"time"

	"github.com/LightningDev1/LB-Selfbot/internal/services/interfaces"
	"github.com/LightningDev1/LB-Selfbot/selfbot/utils"
	"github.com/LightningDev1/dgc"
	"github.com/LightningDev1/discordgo"
)

type EventReactionsServiceImpl struct {
	authSvc     interfaces.AuthService
	uiSvc       interfaces.UIService
	discordSvc  interfaces.DiscordService
	sessionsSvc interfaces.SessionsService
	scJobsSvc   interfaces.SCJobsService
	rpcSvc      interfaces.RPCService

	reactions   []interfaces.EventReaction
	reactionsMu sync.Mutex

	pinSpamMsgs   []*discordgo.Message
	pinSpamMsgsMu sync.Mutex

	eventReactionsMsgs   []*discordgo.Message
	eventReactionsMsgsMu sync.Mutex

	editedMsgs   map[string][]*discordgo.MessageUpdate
	editedMsgsMu sync.Mutex

	deletedMsgs   map[string][]*discordgo.Message
	deletedMsgsMu sync.Mutex

	usedCmds   []*dgc.Ctx
	usedCmdsMu sync.Mutex

	lastActivity time.Time
	lastAFKMsg   time.Time

	paymentSrcID string
	countryCode  string

	timeDiff time.Duration
}

func NewEventReactionsServiceImpl() *EventReactionsServiceImpl {
	return &EventReactionsServiceImpl{
		reactions:          make([]interfaces.EventReaction, 0),
		pinSpamMsgs:        make([]*discordgo.Message, 0),
		eventReactionsMsgs: make([]*discordgo.Message, 0),
		editedMsgs:         make(map[string][]*discordgo.MessageUpdate),
		deletedMsgs:        make(map[string][]*discordgo.Message),
		usedCmds:           make([]*dgc.Ctx, 0),
		lastActivity:       time.Now(),
		lastAFKMsg:         time.Now(),
	}
}

func (e *EventReactionsServiceImpl) Inject(
	authSvc interfaces.AuthService,
	uiSvc interfaces.UIService,
	discordSvc interfaces.DiscordService,
	sessionsSvc interfaces.SessionsService,
	scJobsSvc interfaces.SCJobsService,
	rpcSvc interfaces.RPCService,
) {
	e.authSvc = authSvc
	e.uiSvc = uiSvc
	e.discordSvc = discordSvc
	e.sessionsSvc = sessionsSvc
	e.scJobsSvc = scJobsSvc
	e.rpcSvc = rpcSvc
}

func (e *EventReactionsServiceImpl) Register() {
	session := e.discordSvc.GetSession()

	session.AddHandler(func(_ *discordgo.Session, _ any) { e.authSvc.CheckNonce() })

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

func (e *EventReactionsServiceImpl) handleEvent(s *discordgo.Session, event any) {
	e.reactionsMu.Lock()
	defer e.reactionsMu.Unlock()

	for _, r := range e.reactions {
		go r.Handler(s, event)
	}
}

func (e *EventReactionsServiceImpl) addEditedMessage(message *discordgo.MessageUpdate) {
	e.editedMsgsMu.Lock()
	defer e.editedMsgsMu.Unlock()

	messages, ok := e.editedMsgs[message.ChannelID]
	if !ok {
		messages = []*discordgo.MessageUpdate{}
	}
	e.editedMsgs[message.ChannelID] = append(messages, message)
}

func (e *EventReactionsServiceImpl) addDeletedMessage(message *discordgo.Message) {
	e.deletedMsgsMu.Lock()
	defer e.deletedMsgsMu.Unlock()

	messages, ok := e.deletedMsgs[message.ChannelID]
	if !ok {
		messages = []*discordgo.Message{}
	}
	e.deletedMsgs[message.ChannelID] = append(messages, message)
}

func (e *EventReactionsServiceImpl) GetAllReactions() []interfaces.EventReaction {
	return e.reactions
}

func (e *EventReactionsServiceImpl) GetReactionsByType(t interfaces.EventReactionType) []interfaces.EventReaction {
	return utils.Filter(e.reactions, func(r interfaces.EventReaction) bool {
		return r.Type() == t
	})
}

func (e *EventReactionsServiceImpl) RemoveAllReactions() {
	e.reactionsMu.Lock()
	defer e.reactionsMu.Unlock()

	e.reactions = e.constructDefaultReactions()
}

func (e *EventReactionsServiceImpl) RemoveReactionsByType(t interfaces.EventReactionType) {
	e.reactionsMu.Lock()
	defer e.reactionsMu.Unlock()

	e.reactions = utils.Filter(e.reactions, func(r interfaces.EventReaction) bool {
		return r.Type() != t
	})
}

func (e *EventReactionsServiceImpl) IsReactionMessage(message *discordgo.Message) bool {
	for _, m := range e.eventReactionsMsgs {
		if m.ID == message.ID {
			return true
		}
	}

	return false
}

func (e *EventReactionsServiceImpl) HandleReactionMessage(message *discordgo.Message, err error) {
	if err != nil {
		return
	}

	e.eventReactionsMsgsMu.Lock()
	defer e.eventReactionsMsgsMu.Unlock()

	e.eventReactionsMsgs = append(e.eventReactionsMsgs, message)
}

func (e *EventReactionsServiceImpl) GetPinSpamMessages() []*discordgo.Message {
	e.pinSpamMsgsMu.Lock()
	defer e.pinSpamMsgsMu.Unlock()

	return e.pinSpamMsgs
}

func (e *EventReactionsServiceImpl) AddPinSpamMessage(message *discordgo.Message) {
	e.pinSpamMsgsMu.Lock()
	defer e.pinSpamMsgsMu.Unlock()

	e.pinSpamMsgs = append(e.pinSpamMsgs, message)
}

func (e *EventReactionsServiceImpl) GetEditedMessages(channelID string) []*discordgo.MessageUpdate {
	e.editedMsgsMu.Lock()
	defer e.editedMsgsMu.Unlock()

	return e.editedMsgs[channelID]
}

func (e *EventReactionsServiceImpl) GetDeletedMessages(channelID string) []*discordgo.Message {
	e.deletedMsgsMu.Lock()
	defer e.deletedMsgsMu.Unlock()

	return e.deletedMsgs[channelID]
}

func (e *EventReactionsServiceImpl) GetUsedCommands() []*dgc.Ctx {
	return e.usedCmds
}

func (e *EventReactionsServiceImpl) AddUsedCommand(ctx *dgc.Ctx) {
	e.usedCmdsMu.Lock()
	defer e.usedCmdsMu.Unlock()

	e.usedCmds = append(e.usedCmds, ctx)
}

func (e *EventReactionsServiceImpl) GetLastActivity() time.Time {
	return e.lastActivity
}

func (e *EventReactionsServiceImpl) SetLastActivity() {
	e.lastActivity = time.Now()
}

func (e *EventReactionsServiceImpl) GetPaymentSrcID() string {
	return e.paymentSrcID
}

func (e *EventReactionsServiceImpl) GetCountryCode() string {
	return e.countryCode
}

func (e *EventReactionsServiceImpl) GetTimeDifference() time.Duration {
	return e.timeDiff
}
