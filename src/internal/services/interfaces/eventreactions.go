package interfaces

import (
	"time"

	"github.com/LightningDev1/LB-Selfbot/internal/router"
	"github.com/LightningDev1/discordgo"
)

type EventReactionType int

const (
	// Ignore reactions
	EventReactionIgnoreUser  EventReactionType = 0
	EventReactionIgnoreGuild EventReactionType = 1

	// Ping reactions
	EventReactionPingPing       EventReactionType = 2
	EventReactionPingKick       EventReactionType = 3
	EventReactionPingBan        EventReactionType = 4
	EventReactionPingMute       EventReactionType = 5
	EventReactionPingTimeout    EventReactionType = 6
	EventReactionPingInsult     EventReactionType = 7
	EventReactionPingReact      EventReactionType = 8
	EventReactionPingResponse   EventReactionType = 9
	EventReactionPingAIResponse EventReactionType = 10

	// Mirror reactions
	EventReactionMirrorGuild   EventReactionType = 11
	EventReactionMirrorChannel EventReactionType = 12

	// Responder reactions
	EventReactionResponder   EventReactionType = 13
	EventReactionAIResponder EventReactionType = 14

	// Message reactions
	EventReactionCopy            EventReactionType = 15
	EventReactionReplyCopy       EventReactionType = 16
	EventReactionMock            EventReactionType = 17
	EventReactionPinSpam         EventReactionType = 18
	EventReactionReactSpam       EventReactionType = 19
	EventReactionDeleteMessages  EventReactionType = 20
	EventReactionKeywordNotifier EventReactionType = 21

	// Snipers
	EventReactionSnipers EventReactionType = 22

	// AFK
	EventReactionAFKResponder EventReactionType = 23

	// Phishing
	EventReactionAntiPhishing EventReactionType = 24

	// Protections
	EventReactionProtections EventReactionType = 25

	// Notifiers
	EventReactionNotifiers EventReactionType = 26

	// Selfbot detection
	EventReactionSelfbotDetection EventReactionType = 27

	// Remove reactions
	EventReactionRemoveReactions EventReactionType = 28

	// Force disconnect
	EventReactionForceDisconnect EventReactionType = 29

	// Force nick
	EventReactionForceNick EventReactionType = 30

	// Force kick
	EventReactionForceKick EventReactionType = 31

	// No leave
	EventReactionNoLeave EventReactionType = 32

	// Anti group chat
	EventReactionAntiGroupChat EventReactionType = 33

	// Anti delete messages
	EventReactionAntiDeleteMessages EventReactionType = 34
)

var DefaultEventReactions = []EventReactionType{
	EventReactionSnipers,
	EventReactionAFKResponder,
	EventReactionAntiPhishing,
	EventReactionNotifiers,
	EventReactionProtections,
	EventReactionSelfbotDetection,
}

type EventReaction interface {
	Type() EventReactionType
	Description() string
	Equals(other EventReaction) bool
	Data() map[string]any
	Handler(session *discordgo.Session, event any)
}

type EventReactionsService interface {
	Inject(authSvc AuthService, uiSvc UIService, discordSvc DiscordService, sessionsSvc SessionsService, scJobsSvc SCJobsService, rpcSvc RPCService)

	Register()

	GetAllReactions() []EventReaction
	GetReactionsByType(t EventReactionType) []EventReaction
	AddReaction(r EventReaction)
	RemoveReaction(r EventReaction)
	RemoveAllReactions()
	RemoveReactionsByType(t EventReactionType)

	IsReactionMessage(message *discordgo.Message) bool
	HandleReactionMessage(message *discordgo.Message, err error)

	LoadCmdAliases() error

	GetPinSpamMessages() []*discordgo.Message
	AddPinSpamMessage(message *discordgo.Message)
	RemovePinSpamMessages()

	GetEditedMessages(channelID string) []*discordgo.MessageUpdate
	GetDeletedMessages(channelID string) []*discordgo.Message

	GetUsedCommands() []*router.Ctx
	AddUsedCommand(ctx *router.Ctx)

	GetLastActivity() time.Time
	SetLastActivity()

	GetPaymentSrcID() string
	GetCountryCode() string

	GetTimeDifference() time.Duration
}
