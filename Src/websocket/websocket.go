package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/gorilla/websocket"

	"github.com/LightningDev1/beeep"

	"github.com/LightningDev1/LB-Selfbot/config"
	"github.com/LightningDev1/LB-Selfbot/encryption"
	"github.com/LightningDev1/LB-Selfbot/ipc"
	"github.com/LightningDev1/LB-Selfbot/utils"
)

type WebsocketConnection struct {
	Conn          *websocket.Conn
	RequestHeader http.Header
	WriteMutex    *sync.Mutex

	Context context.Context

	Encryption *encryption.AESCBC

	Connected    bool
	Ready        bool
	ReadyChannel chan struct{}

	// Sliding window for errors
	Errors      int
	MaxErrors   int
	ErrorStart  time.Time
	ErrorWindow time.Duration

	LogError func(err error, closeCode int)

	Username string
	Password string
	Hostname string
	URL      string
}

func (w *WebsocketConnection) InitEncryption() {
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

	w.Encryption = &encryption.AESCBC{Key: string(key)}
}

func (w *WebsocketConnection) Init(username, password string) error {
	w.InitEncryption()

	w.RequestHeader = http.Header{"Origin": []string{"http://localhost/"}}
	w.WriteMutex = &sync.Mutex{}

	w.Connected = false
	w.Ready = false
	w.ReadyChannel = make(chan struct{}, 128)

	w.MaxErrors = 5
	w.ErrorStart = time.Now()
	w.ErrorWindow = 30 * time.Second

	w.Username = username
	w.Password = password
	w.Hostname, _ = os.Hostname()
	w.URL = "wss://api.lightning-bot.com/api/v3/ws/selfbot"

	go w.Listen()
	go w.Ping()

	return w.Connect()
}

func (w *WebsocketConnection) Connect() (err error) {
	w.Conn, _, err = websocket.DefaultDialer.Dial(w.URL, w.RequestHeader)
	if err != nil {
		return err
	}

	w.Conn.SetCloseHandler(w.CloseHandler)

	w.Connected = true

	// Wait a bit before sending HELLO to allow
	// the listener to start listening
	time.Sleep(100 * time.Millisecond)

	err = w.SendHello()
	if err != nil {
		return err
	}

	return nil
}

func (w *WebsocketConnection) SendHello() error {
	return w.Send(WsMessage{
		Op: "HELLO",
		Data: WsData{
			Username:     w.Username,
			Password:     w.Password,
			ComputerName: w.Hostname,
			Time:         time.Now().Unix(),
			ConsoleText:  ipc.Console,
		},
	})
}

func (w *WebsocketConnection) SendConfigUpdate() error {
	cfg, _ := config.GetConfig()

	password := cfg.Account.LightningBot.Password
	token := cfg.Account.Discord.Token

	cfg.Account.LightningBot.Password = "redacted"
	cfg.Account.Discord.Token = "redacted"
	cfg.Snipers.Tokens = []string{}

	err := w.Send(WsMessage{
		Op:   "CONFIG_UPDATE",
		Data: WsData{Config: &cfg},
	})

	cfg.Account.LightningBot.Password = password
	cfg.Account.Discord.Token = token

	return err
}

func (w *WebsocketConnection) Send(msg WsMessage) error {
	if w.Encryption == nil || w.Conn == nil {
		return fmt.Errorf("websocket connection is not initialized")
	}

	jsonMsg, _ := json.Marshal(msg)
	encrypted, _ := w.Encryption.Encrypt(jsonMsg)

	return w.Conn.WriteMessage(websocket.TextMessage, encrypted)
}

func (w *WebsocketConnection) Read() (msg WsMessage, err error) {
	msgType, data, err := w.Conn.ReadMessage()
	if err != nil {
		return msg, err
	}

	if msgType != websocket.TextMessage {
		return msg, fmt.Errorf("received non-text message")
	}

	decrypted, err := w.Encryption.Decrypt(data)
	if err != nil {
		return msg, err
	}

	err = json.Unmarshal(decrypted, &msg)
	return msg, err
}

func (w *WebsocketConnection) Listen() {
	for {
		w.WaitUntilConnected()

		if !w.Connected {
			continue
		}

		msg, err := w.Read()
		if err != nil {
			w.OnError(err)
			continue
		}

		utils.Logging.Debug("received message:", msg)

		go w.OnMessage(msg)
	}
}

func (w *WebsocketConnection) Ping() {
	for {
		w.WaitUntilConnected()

		time.Sleep(60 * time.Second)

		err := w.Send(WsMessage{Op: "PING"})

		if err != nil {
			w.OnError(err)
		}
	}
}

func (w *WebsocketConnection) WaitUntilConnected() {
	// We can't use a channel here because we need this to work in multiple goroutines
	for !w.Connected {
		time.Sleep(20 * time.Millisecond)
	}
}

func (w *WebsocketConnection) WaitUntilReady() {
	if w.Ready {
		return
	}

	<-w.ReadyChannel
}

func (w *WebsocketConnection) OnMessage(msg WsMessage) {
	switch msg.Op {
	case "HELLO_ACK":
		w.Ready = true
		w.ReadyChannel <- struct{}{}

		w.SendConfigUpdate()
	case "FREE_TRIAL_START":
		go beeep.Notify(
			"LightningBot",
			"Your free trial has started! You will be able to use LightningBot for free for the next 7 days.",
			"LightningBot",
			filepath.Join(utils.IO.GetMainFolder(), "Logo.png"),
		)
	case "EXIT":
		os.Exit(0)
	case "RESTART":
		ipc.IPC.Restart()
	case "WEBSITE_CONNECT":
		_ = w.Send(WsMessage{
			Op:   "CONSOLE_SET",
			Data: WsData{Text: ipc.Console},
		})

		w.SendConfigUpdate()
	case "CONFIG_UPDATE":
		if msg.Data.Config == nil {
			break
		}

		cfg, _ := config.GetConfig()

		msg.Data.Config.Account = cfg.Account

		_ = msg.Data.Config.Save()

		ipc.ConfigUpdate(msg.Data.Config)
	case "PONG":
		break
	case "ANOTHER_SELFBOT_CONNECTED":
		// handled by error handler
		break
	default:
		utils.Logging.Warn("Unknown op", msg.Op)
	}
}

func (w *WebsocketConnection) CloseHandler(code int, text string) error {
	w.Connected = false

	utils.Logging.Error("websocket closed", code, text)

	return nil
}

func (w *WebsocketConnection) OnError(err error) {
	w.Connected = false

	utils.Logging.Error("websocket error:", err)

	closeError, ok := err.(*websocket.CloseError)
	if ok && closeError.Code > 3000 {
		w.LogError(err, closeError.Code)
		return
	}

	w.Errors++

	if time.Since(w.ErrorStart) > w.ErrorWindow {
		w.ErrorStart = time.Now()
		w.Errors = 1
	}

	if w.Errors > w.MaxErrors {
		w.LogError(err, -1)
		os.Exit(1)
	}

	time.Sleep(2 * time.Second)

	utils.Logging.Info("reconnecting websocket...")

	err = w.Connect()
	if err != nil {
		w.OnError(err)
	}
}

var Websocket = &WebsocketConnection{}
