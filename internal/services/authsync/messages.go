package authsync

import "github.com/lb-selfbot/selfbot/internal/services/config"

// WSMessage represents a websocket message with an operation code and data payload
type WSMessage struct {
	Op   WSOpCode `json:"op"`
	Data WSData   `json:"data"`
}

// WSData represents the data payload of a websocket message
type WSData struct {
	Username     string         `json:"username,omitempty"`
	Password     string         `json:"password,omitempty"`
	ComputerName string         `json:"computerName,omitempty"`
	Time         int64          `json:"time,omitempty"`
	ConsoleText  string         `json:"consoleText,omitempty"`
	Nonce        string         `json:"nonce,omitempty"`
	Config       *config.Config `json:"config,omitempty"`
	Text         string         `json:"text,omitempty"`
}

// WSOpCode represents the type of websocket operation being performed
type WSOpCode string

const (
	OpCodeHello              WSOpCode = "HELLO"
	OpCodeHelloAck           WSOpCode = "HELLO_ACK"
	OpCodePing               WSOpCode = "PING"
	OpCodePong               WSOpCode = "PONG"
	OpCodeExit               WSOpCode = "EXIT"
	OpCodeRestart            WSOpCode = "RESTART"
	OpCodeWebsiteConnect     WSOpCode = "WEBSITE_CONNECT"
	OpCodeConfigUpdate       WSOpCode = "CONFIG_UPDATE"
	OpCodeConsoleAppend      WSOpCode = "CONSOLE_APPEND"
	OpCodeConsoleClear       WSOpCode = "CONSOLE_CLEAR"
	OpCodeConsoleSet         WSOpCode = "CONSOLE_SET"
	OpCodeFreeTrialStart     WSOpCode = "FREE_TRIAL_START"
	OpCodeAnotherSBConnected WSOpCode = "ANOTHER_SELFBOT_CONNECTED"
)
