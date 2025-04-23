package ipc

import "github.com/lb-selfbot/selfbot/internal/services/config"

// WSMessage represents a websocket message with an operation code and data payload
type WSMessage struct {
	Op   WSOpCode `json:"op"`
	Data WSData   `json:"data"`
}

type WsMessageReceive struct {
	Op   WSOpCode       `json:"op"`
	Data *config.Config `json:"data"`
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
	OpCodeConnected     WSOpCode = 98
	OpCodeHello         WSOpCode = 99
	OpCodeExit          WSOpCode = 100
	OpCodeRestart       WSOpCode = 101
	OpCodeNotification  WSOpCode = 102
	OpCodeInvalidToken  WSOpCode = 103
	OpCodeError         WSOpCode = 104
	OpCodeConfigUpdate  WSOpCode = 105
	OpCodeNewToken      WSOpCode = 106
	OpCodeConsoleAppend WSOpCode = 107
	OpCodeConsoleClear  WSOpCode = 108
	OpCodeConsoleSet    WSOpCode = 109
)
