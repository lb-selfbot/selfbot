package websocket

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/LightningDev1/LB-Selfbot/constants"
	"github.com/LightningDev1/LB-Selfbot/encryption"
	"github.com/LightningDev1/LB-Selfbot/ipc"
	"github.com/rgamba/evtwebsocket"
)

var Conn *evtwebsocket.Conn
var Encryption *encryption.AESCBC
var Connected bool
var Connect chan struct{} = make(chan struct{})

type WsData struct {
	Username     string `json:"username,omitempty"`
	Password     string `json:"password,omitempty"`
	ComputerName string `json:"computer_name,omitempty"`
	Time         int64  `json:"time,omitempty"`
	ConsoleText  string `json:"console_text,omitempty"`
	Text         string `json:"text,omitempty"`
}

type WsMessage struct {
	Op   string `json:"op,omitempty"`
	Data WsData `json:"data,omitempty"`
}

func Init(username, password string) error {
	// XOR encryption
	encrypted := []int{62, 46, 10, 38, 2, 61, 98, 56, 6, 96, 6, 55, 60, 1, 17, 50, 7, 97, 67, 11, 51, 28, 15, 123, 46, 48, 8, 50, 42, 43, 36, 16, 99, 61, 28, 34, 11, 65, 100, 4, 93, 22, 59, 68, 17, 51, 7, 49, 4, 60, 51, 15, 42, 16, 19, 56, 4, 91, 21, 10, 32, 103, 102, 6, 43, 97, 34, 31, 13, 32, 47, 42, 51, 8, 57, 14, 21, 33, 57, 78, 110, 13, 0, 55, 78, 41, 122, 86, 97, 125, 90, 22, 126, 7, 68, 47, 71, 28, 24, 106}
	keyInt := []int{89, 85, 112, 112, 105, 98, 113, 73, 102, 110, 75, 83, 117, 66, 71, 119, 111, 69, 117, 79, 121, 112, 71, 97, 105, 104, 90, 71, 89, 74, 66, 117, 81, 75, 83, 117, 83, 83, 85, 84, 67, 88, 108, 83, 90, 71, 74, 111, 118, 117}
	decrypted := ""
	for i := len(encrypted) - 1; i > -1; i-- {
		decrypted += string(rune(encrypted[i] ^ keyInt[(len(encrypted)-1-i)%len(keyInt)]))
	}
	encryptedKey := []byte{121, 106, 120, 89, 65, 65, 81, 69, 81, 88, 88, 52, 99, 79, 52, 76, 102, 49, 49, 83, 73, 120, 101, 105, 114, 104, 53, 119, 54, 73, 102, 81, 53, 109, 102, 55, 109, 43, 72, 83, 90, 76, 85, 61, 58, 113, 113, 69, 43, 83, 65, 54, 114, 77, 51, 86, 83, 66, 49, 113, 112, 103, 65, 78, 68, 78, 81, 61, 61, 58, 87, 107, 74, 106, 75, 115, 82, 70, 85, 43, 84, 84, 117, 90, 67, 67, 97, 120, 110, 72, 79, 86, 110, 110, 52, 83, 68, 77, 108, 115, 115, 112, 57, 74, 82, 66, 81, 70, 106, 69, 47, 88, 49, 98, 111, 104, 49, 73, 71, 97, 118, 122, 97, 54, 66, 101, 110, 80, 79, 99, 110, 56, 88, 53, 48, 122, 52, 50, 117, 54, 84, 55, 112, 73, 116, 67, 121, 78, 56, 71, 52, 67, 76, 87, 57, 83, 88, 88, 65, 68, 68, 82, 80, 82, 115, 113, 86, 109, 111, 116, 84, 117, 54, 98, 120, 51, 53, 79, 87, 110, 67, 88, 113, 117, 90, 102, 70, 78, 53, 84, 76, 66, 121, 115, 78, 115, 72, 101}
	keyEncryption := &encryption.AESCBC{Key: decrypted}
	key, _ := keyEncryption.Decrypt(encryptedKey)

	Encryption = &encryption.AESCBC{Key: string(key)}
	pingMsg, _ := json.Marshal(WsMessage{Op: "PING"})
	pingMsgEncrypted, _ := Encryption.Encrypt(pingMsg)
	Conn = &evtwebsocket.Conn{
		OnConnected: func(*evtwebsocket.Conn) {
			hostname, _ := os.Hostname()
			Send(WsMessage{
				Op: "HELLO",
				Data: WsData{
					Username:     username,
					Password:     password,
					ComputerName: hostname,
					Time:         time.Now().Unix(),
					ConsoleText:  ipc.Console,
				},
			})
		},
		OnMessage: func(msg []byte, _ *evtwebsocket.Conn) {
			decryptedMsg, _ := Encryption.Decrypt(msg)
			wsMessage := WsMessage{}
			_ = json.Unmarshal(decryptedMsg, &wsMessage)
			switch wsMessage.Op {
			case "HELLO_ACK":
				Connected = true
				Connect <- struct{}{}
			case "PONG":
				break
			case "EXIT":
				os.Exit(0)
			case "RESTART":
				ipc.IPC.Restart()
			case "WEBSITE_CONNECT":
				Send(WsMessage{
					Op: "CONSOLE_SET",
					Data: WsData{
						Text: ipc.Console,
					},
				})
			default:
				fmt.Println(wsMessage)
			}
		},
		OnError: func(err error) {
			//logging.Error(err)
			if !constants.Debug {
				os.Exit(0)
			}
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
	_ = Conn.Send(evtwebsocket.Msg{Body: encrypted})
}

func WaitForConnection() {
	if Connected {
		return
	}
	<-Connect
}
