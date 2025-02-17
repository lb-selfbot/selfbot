package ipc

import "github.com/LightningDev1/LB-Selfbot/selfbot/config"

// WSMessage represents a websocket message with an operation code and data payload
type WSMessage struct {
	Op   WSOpCode `json:"op"`
	Data WSData   `json:"data"`
}

// WSData represents the data payload of a websocket message
type WSData struct {
	Text       string         `json:"text,omitempty"`
	IsRaw      bool           `json:"is_raw,omitempty"`
	MainFolder string         `json:"mainFolder,omitempty"`
	Config     *config.Config `json:"config,omitempty"`
}

// WSOpCode represents the type of websocket operation being performed
type WSOpCode int

const (
	OP_CONNECTED      WSOpCode = 98
	OP_HELLO          WSOpCode = 99
	OP_EXIT           WSOpCode = 100
	OP_RESTART        WSOpCode = 101
	OP_NOTIFICATION   WSOpCode = 102
	OP_INVALID_TOKEN  WSOpCode = 103
	OP_ERROR          WSOpCode = 104
	OP_CONFIG_UPDATE  WSOpCode = 105
	OP_NEW_TOKEN      WSOpCode = 106
	OP_CONSOLE_APPEND WSOpCode = 107
	OP_CONSOLE_CLEAR  WSOpCode = 108
	OP_CONSOLE_SET    WSOpCode = 109
)
