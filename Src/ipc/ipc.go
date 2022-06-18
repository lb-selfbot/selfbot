package ipc

import (
	"net"
	"os"
)

var (
	Ipc IPC = IPC{}
)

type JsonEventData struct {
	Text string `json:"text,omitempty"`
}

type JsonEvent struct {
	Op   int           `json:"op"`
	Data JsonEventData `json:"data"`
}

type OpCodes struct {
	EXIT          int
	CLS           int
	NOTIFICATION  int
	INVALID_TOKEN int
	ERROR         int
	CONFIG_EDIT   int
	CONSOLE       int
	RESTART       int
	NEW_TOKEN     int
}

type IPC struct {
	readPipe  net.Conn
	writeFile *os.File
	OnRead    func(string)
	OpCodes   OpCodes
}
