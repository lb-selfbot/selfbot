package ipc

import (
	"encoding/json"
	"net"
	"os"
)

var (
	Ipc = IPC{}
)

type JsonEventData struct {
	Text string `json:"text,omitempty"`
}

type JsonEvent struct {
	Op   int           `json:"op"`
	Data JsonEventData `json:"data"`
}

type OpCodes struct {
	EXIT         int
	CLS          int
	NOTIFICATION int
	InvalidToken int
	ERROR        int
	ConfigEdit   int
	CONSOLE      int
	RESTART      int
	NewToken     int
}

type IPC struct {
	readPipe  net.Conn
	writeFile *os.File
	OnRead    func(string)
	OpCodes   OpCodes
}

func SendNotification(text string) {
	jsonBytes, _ := json.Marshal(JsonEvent{
		Op: Ipc.OpCodes.NOTIFICATION,
		Data: JsonEventData{
			Text: text,
		},
	})
	Ipc.Write(string(jsonBytes))
}
