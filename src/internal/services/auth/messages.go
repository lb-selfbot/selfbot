package auth

import "github.com/LightningDev1/LB-Selfbot/selfbot/config"

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
	OP_HELLO                     WSOpCode = "HELLO"
	OP_HELLO_ACK                 WSOpCode = "HELLO_ACK"
	OP_PING                      WSOpCode = "PING"
	OP_PONG                      WSOpCode = "PONG"
	OP_EXIT                      WSOpCode = "EXIT"
	OP_RESTART                   WSOpCode = "RESTART"
	OP_WEBSITE_CONNECT           WSOpCode = "WEBSITE_CONNECT"
	OP_CONFIG_UPDATE             WSOpCode = "CONFIG_UPDATE"
	OP_CONSOLE_APPEND            WSOpCode = "CONSOLE_APPEND"
	OP_CONSOLE_CLEAR             WSOpCode = "CONSOLE_CLEAR"
	OP_CONSOLE_SET               WSOpCode = "CONSOLE_SET"
	OP_FREE_TRIAL_START          WSOpCode = "FREE_TRIAL_START"
	OP_ANOTHER_SELFBOT_CONNECTED WSOpCode = "ANOTHER_SELFBOT_CONNECTED"
)
