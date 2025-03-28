package eventbroker

import (
	"sync"
)

// EventType represents the type of event that can be published and subscribed to.
type EventType int

const (
	EventUnknown EventType = iota

	// Internal events
	EventCheckNonce
	EventSendConfigUpdateSync
	EventRestart
	EventCaptureError
	EventSessionReady

	// EventReactions events
	EventHandleReactionMessage
	EventHandlePinSpamMessage

	// UI Events
	EventSetConsole
	EventAppendConsole
	EventClearConsole
	EventShowNotification
	EventLogError
	EventLogInfo
	EventLogFeature
)

// EventHandler is a function that handles an event. It takes a payload of any type.
type EventHandler func(payload any)

// Service is an event broker.
type Service interface {
	// Subscribe subscribes to an event type with a handler.
	Subscribe(eventType EventType, handler EventHandler)

	// Publish publishes an event type with a payload.
	Publish(eventType EventType, payload any)
}

type eventBroker struct {
	subscribers map[EventType][]EventHandler
	mu          sync.RWMutex
}

func NewService() Service {
	return &eventBroker{
		subscribers: make(map[EventType][]EventHandler),
	}
}

func (b *eventBroker) Subscribe(eventType EventType, handler EventHandler) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.subscribers[eventType] = append(b.subscribers[eventType], handler)
}

func (b *eventBroker) Publish(eventType EventType, payload any) {
	b.mu.RLock()
	handlers := make([]EventHandler, len(b.subscribers[eventType]))
	copy(handlers, b.subscribers[eventType])
	b.mu.RUnlock()

	for _, handler := range handlers {
		go handler(payload)
	}
}
