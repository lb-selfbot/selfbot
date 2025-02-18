package eventreactions

import (
	"sync"
	"time"

	"github.com/LightningDev1/LB-Selfbot/internal/services/interfaces"
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
	configSvc   interfaces.ConfigService

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
	configSvc interfaces.ConfigService,
) {
	e.authSvc = authSvc
	e.uiSvc = uiSvc
	e.discordSvc = discordSvc
	e.sessionsSvc = sessionsSvc
	e.scJobsSvc = scJobsSvc
	e.rpcSvc = rpcSvc
	e.configSvc = configSvc
}

func (e *EventReactionsServiceImpl) Register() {
	session, router := e.discordSvc.GetSession(), e.discordSvc.GetRouter()

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
		go r.Handler(event)
	}
}

func (e *EventReactionsServiceImpl) isEventReactionMessage(message *discordgo.Message) bool {
	e.eventReactionsMsgsMu.Lock()
	defer e.eventReactionsMsgsMu.Unlock()

	for _, m := range e.eventReactionsMsgs {
		if m.ID == message.ID {
			return true
		}
	}
	return false
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
