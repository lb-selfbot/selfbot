package ipc

import "github.com/lb-selfbot/selfbot/internal/services/config"

// IPCMessage represents an IPC message with an operation code and data payload
type IPCMessage struct {
	Op   OpCode   `json:"op"`
	Data *IPCData `json:"data,omitempty"`
}

// IPCData represents the data payload of an IPC message
type IPCData struct {
	Text   string         `json:"text,omitempty"`
	IsRaw  bool           `json:"is_raw,omitempty"`
	Config *config.Config `json:"config,omitempty"`
}

// OpCode represents the type of operation being performed
type OpCode int

const (
	OpCodePing          OpCode = 0
	OpCodeExit          OpCode = 1
	OpCodeRestart       OpCode = 2
	OpCodeError         OpCode = 3
	OpCodeConfigUpdate  OpCode = 4
	OpCodeConsoleAppend OpCode = 5
	OpCodeConsoleClear  OpCode = 6
	OpCodeConsoleSet    OpCode = 7
)
