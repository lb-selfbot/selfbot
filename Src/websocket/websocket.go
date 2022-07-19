package websocket

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/LightningDev1/LB-Selfbot/encryption"
	"github.com/LightningDev1/LB-Selfbot/utils"
	"github.com/rgamba/evtwebsocket"
)

var Conn *evtwebsocket.Conn
var Encryption *encryption.Encryption

type WsData struct {
	Username     string `json:"username,omitempty"`
	Password     string `json:"password,omitempty"`
	ComputerName string `json:"computer_name,omitempty"`
	Time         int64  `json:"time,omitempty"`
}

type WsMessage struct {
	Op   string `json:"op,omitempty"`
	Data WsData `json:"data,omitempty"`
}

func Init() error {
	Encryption = &encryption.Encryption{
		Key: "?aZV:3EQ_;'_B%g#q*Eg2#RQU.jCTo`qE%a>6RFRn-KXW`'.GA`38Zf+%+4]7qJF2UBDha1OoFEoC9rP",
	}
	pingMsg, _ := json.Marshal(WsMessage{Op: "PING"})
	pingMsgEncrypted, _ := Encryption.Encrypt(pingMsg)
	Conn = &evtwebsocket.Conn{
		OnConnected: func(*evtwebsocket.Conn) {
			Auth("Lightning", "Fatos_1984")
			fmt.Println("Connected!")
		},
		OnMessage: func(msg []byte, _ *evtwebsocket.Conn) {
			decryptedMsg, _ := Encryption.Decrypt(msg)
			wsMessage := WsMessage{}
			json.Unmarshal(decryptedMsg, &wsMessage)
			switch wsMessage.Op {
			case "PONG":
				fmt.Println("Received pong")
			case "EXIT":
				os.Exit(0)
			case "RESTART":
				utils.Restart()
			default:
				fmt.Println(wsMessage)
			}
		},
		OnError: func(err error) {
			fmt.Printf("Error: %s\n", err.Error())
		},
		PingIntervalSecs: 5,
		PingMsg:          pingMsgEncrypted,
		Reconnect:        true,
	}
	return Conn.Dial("ws://127.0.0.1:8000/api/v3/ws/selfbot", "")
}

func Send(message WsMessage) {
	jsonMessage, _ := json.Marshal(message)
	encrypted, _ := Encryption.Encrypt(jsonMessage)
	Conn.Send(evtwebsocket.Msg{
		Body: encrypted,
	})
}

func Auth(username, password string) {
	hostname, _ := os.Hostname()
	Send(WsMessage{
		Op: "HELLO",
		Data: WsData{
			Username:     username,
			Password:     password,
			ComputerName: hostname,
			Time:         time.Now().Unix(),
		},
	})

}
