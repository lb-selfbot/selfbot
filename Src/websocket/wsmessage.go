package websocket

import "github.com/LightningDev1/LB-Selfbot/config"

type WsData struct {
	Username     string         `json:"username,omitempty"`
	Password     string         `json:"password,omitempty"`
	ComputerName string         `json:"computer_name,omitempty"`
	Time         int64          `json:"time,omitempty"`
	ConsoleText  string         `json:"console_text,omitempty"`
	Text         string         `json:"text,omitempty"`
	Config       *config.Config `json:"config,omitempty"`
	Nonce        string         `json:"nonce,omitempty"`
}

type WsMessage struct {
	Op   string `json:"op,omitempty"`
	Data WsData `json:"data,omitempty"`
}
