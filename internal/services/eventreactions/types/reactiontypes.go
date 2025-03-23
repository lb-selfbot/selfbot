package types

import "github.com/lb-selfbot/discordgo"

type EventReactionType int

const (
	// Ignore reactions
	IgnoreUser  EventReactionType = 0
	IgnoreGuild EventReactionType = 1

	// Ping reactions
	PingPing       EventReactionType = 2
	PingKick       EventReactionType = 3
	PingBan        EventReactionType = 4
	PingMute       EventReactionType = 5
	PingTimeout    EventReactionType = 6
	PingInsult     EventReactionType = 7
	PingReact      EventReactionType = 8
	PingResponse   EventReactionType = 9
	PingAIResponse EventReactionType = 10

	// Mirror reactions
	MirrorGuild   EventReactionType = 11
	MirrorChannel EventReactionType = 12

	// Responder reactions
	Responder   EventReactionType = 13
	AIResponder EventReactionType = 14

	// Message reactions
	Copy            EventReactionType = 15
	ReplyCopy       EventReactionType = 16
	Mock            EventReactionType = 17
	PinSpam         EventReactionType = 18
	ReactSpam       EventReactionType = 19
	DeleteMessages  EventReactionType = 20
	KeywordNotifier EventReactionType = 21

	// Snipers
	Snipers EventReactionType = 22

	// AFK
	AFKResponder EventReactionType = 23

	// Phishing
	AntiPhishing EventReactionType = 24

	// Protections
	Protections EventReactionType = 25

	// Notifiers
	Notifiers EventReactionType = 26

	// Selfbot detection
	SelfbotDetection EventReactionType = 27

	// Remove reactions
	RemoveReactions EventReactionType = 28

	// Force disconnect
	ForceDisconnect EventReactionType = 29

	// Force nick
	ForceNick EventReactionType = 30

	// Force kick
	ForceKick EventReactionType = 31

	// No leave
	NoLeave EventReactionType = 32

	// Anti group chat
	AntiGroupChat EventReactionType = 33

	// Anti delete messages
	AntiDeleteMessages EventReactionType = 34
)

var DefaultReactions = []EventReactionType{
	Snipers,
	AFKResponder,
	AntiPhishing,
	Notifiers,
	Protections,
	SelfbotDetection,
}

type EventReaction interface {
	Type() EventReactionType
	Description() string
	Equals(other EventReaction) bool
	Data() map[string]any
	Handler(session *discordgo.Session, event any)
}

type MessageHandler func(*discordgo.Message, error)
