package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/LightningDev1/LB-Selfbot/internal/config"
	"github.com/LightningDev1/LB-Selfbot/internal/constants"
	"github.com/LightningDev1/LB-Selfbot/internal/encryption"
	"github.com/LightningDev1/LB-Selfbot/internal/logging"
	"github.com/LightningDev1/LB-Selfbot/internal/random"
	"github.com/LightningDev1/LB-Selfbot/internal/services/interfaces"
	"github.com/LightningDev1/LB-Selfbot/internal/setup"
	"github.com/LightningDev1/LB-Selfbot/internal/system"
	"github.com/gorilla/websocket"
)

const (
	defaultWebsocketURL = "wss://api.lightning-bot.com/api/v3/ws/selfbot"
	maxErrors           = 20
	errorWindow         = 15 * time.Minute
	pingInterval        = 60 * time.Second
	reconnectDelay      = 30 * time.Second
	handshakeDelay      = 100 * time.Millisecond
	nonceLength         = 64
)

type AuthServiceImpl struct {
	uiSvc  interfaces.UIService
	ipcSvc interfaces.IPCService

	conn *websocket.Conn
	mu   sync.RWMutex
	enc  *encryption.AESCBC

	readyChan chan struct{}

	creds interfaces.AuthCredentials
	nonce string

	receivedNonce string

	errors     int
	errorStart time.Time
}

func (a *AuthServiceImpl) Inject(uiSvc interfaces.UIService, ipcSvc interfaces.IPCService) {
	a.uiSvc = uiSvc
	a.ipcSvc = ipcSvc
}

func (a *AuthServiceImpl) Connect(context context.Context, creds interfaces.AuthCredentials) error {
	if err := a.setupEncryption(); err != nil {
		a.logError(err, -1)
		return err
	}

	a.readyChan = make(chan struct{})
	a.creds = creds
	a.nonce = random.String(nonceLength)

	err := a.connectWS(context)
	if err != nil {
		a.logError(err, -1)
	}

	return err
}

func (a *AuthServiceImpl) OnReady() <-chan struct{} {
	return a.readyChan
}

func (a *AuthServiceImpl) CheckNonce() {
	isValid := constants.HOSTING || a.receivedNonce == "" || a.receivedNonce == a.nonce
	if isValid {
		return
	}

	logging.Error("An AuthWS error has occurred. Error code: " + random.String(8))
	os.Exit(1)
}

func (a *AuthServiceImpl) setupEncryption() error {
	encrypted := []int{62, 46, 10, 38, 2, 61, 98, 56, 6, 96, 6, 55, 60, 1, 17, 50, 7, 97, 67, 11, 51, 28, 15, 123, 46, 48, 8, 50, 42, 43, 36, 16, 99, 61, 28, 34, 11, 65, 100, 4, 93, 22, 59, 68, 17, 51, 7, 49, 4, 60, 51, 15, 42, 16, 19, 56, 4, 91, 21, 10, 32, 103, 102, 6, 43, 97, 34, 31, 13, 32, 47, 42, 51, 8, 57, 14, 21, 33, 57, 78, 110, 13, 0, 55, 78, 41, 122, 86, 97, 125, 90, 22, 126, 7, 68, 47, 71, 28, 24, 106}
	keyInt := []int{89, 85, 112, 112, 105, 98, 113, 73, 102, 110, 75, 83, 117, 66, 71, 119, 111, 69, 117, 79, 121, 112, 71, 97, 105, 104, 90, 71, 89, 74, 66, 117, 81, 75, 83, 117, 83, 83, 85, 84, 67, 88, 108, 83, 90, 71, 74, 111, 118, 117}
	decrypted := ""
	for i := len(encrypted) - 1; i > -1; i-- {
		decrypted += string(rune(encrypted[i] ^ keyInt[(len(encrypted)-1-i)%len(keyInt)]))
	}

	encryptedKey := []byte{121, 106, 120, 89, 65, 65, 81, 69, 81, 88, 88, 52, 99, 79, 52, 76, 102, 49, 49, 83, 73, 120, 101, 105, 114, 104, 53, 119, 54, 73, 102, 81, 53, 109, 102, 55, 109, 43, 72, 83, 90, 76, 85, 61, 58, 113, 113, 69, 43, 83, 65, 54, 114, 77, 51, 86, 83, 66, 49, 113, 112, 103, 65, 78, 68, 78, 81, 61, 61, 58, 87, 107, 74, 106, 75, 115, 82, 70, 85, 43, 84, 84, 117, 90, 67, 67, 97, 120, 110, 72, 79, 86, 110, 110, 52, 83, 68, 77, 108, 115, 115, 112, 57, 74, 82, 66, 81, 70, 106, 69, 47, 88, 49, 98, 111, 104, 49, 73, 71, 97, 118, 122, 97, 54, 66, 101, 110, 80, 79, 99, 110, 56, 88, 53, 48, 122, 52, 50, 117, 54, 84, 55, 112, 73, 116, 67, 121, 78, 56, 71, 52, 67, 76, 87, 57, 83, 88, 88, 65, 68, 68, 82, 80, 82, 115, 113, 86, 109, 111, 116, 84, 117, 54, 98, 120, 51, 53, 79, 87, 110, 67, 88, 113, 117, 90, 102, 70, 78, 53, 84, 76, 66, 121, 115, 78, 115, 72, 101}

	keyEncryption := &encryption.AESCBC{Key: decrypted}
	key, err := keyEncryption.Decrypt(encryptedKey)
	if err != nil {
		return fmt.Errorf("failed to decrypt key: %w", err)
	}

	a.enc = &encryption.AESCBC{Key: string(key)}
	return nil
}

func (a *AuthServiceImpl) connectWS(context context.Context) error {
	header := http.Header{"origin": []string{"http://localhost/"}}

	var err error
	a.conn, _, err = websocket.DefaultDialer.DialContext(context, defaultWebsocketURL, header)
	if err != nil {
		return err
	}

	// Get hostname
	hostname := "unknown"
	if constants.HOSTING {
		hostname = "LB Hosting"
	} else {
		hostname, _ = os.Hostname()
	}

	// Send HELLO
	err = a.write(&WSMessage{
		Op: OpCodeHello,
		Data: WSData{
			Username:     a.creds.Username,
			Password:     a.creds.Password,
			ComputerName: hostname,
			Time:         time.Now().Unix(),
			ConsoleText:  a.ipcSvc.GetHTMLConsole(),
			Nonce:        a.nonce,
		},
	})
	if err != nil {
		return fmt.Errorf("failed to send HELLO: %w", err)
	}

	go a.listen()
	go a.ping(context)

	return nil
}

func (a *AuthServiceImpl) listen() {
	for {
		msg, err := a.read()
		if err != nil {
			a.onError(err)
			continue
		}

		if err := a.onMessage(msg); err != nil {
			a.onError(err)
			continue
		}
	}
}

func (a *AuthServiceImpl) ping(context context.Context) {
	ticker := time.NewTicker(pingInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := a.write(&WSMessage{Op: OpCodePing}); err != nil {
				a.onError(err)
				continue
			}
		case <-context.Done():
			return
		}
	}
}

func (a *AuthServiceImpl) read() (*WSMessage, error) {
	if a.conn == nil || a.enc == nil {
		return nil, fmt.Errorf("connection not initialized")
	}

	msgType, data, err := a.conn.ReadMessage()
	if err != nil {
		return nil, fmt.Errorf("failed to read message: %w", err)
	}

	if msgType != websocket.TextMessage {
		return nil, fmt.Errorf("unexpected message type: %d", msgType)
	}

	decrypted, err := a.enc.Decrypt(data)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt message: %w", err)
	}

	msg := &WSMessage{}
	if err := json.Unmarshal(decrypted, msg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal message: %w", err)
	}

	return msg, nil
}

func (a *AuthServiceImpl) write(msg *WSMessage) error {
	if a.conn == nil || a.enc == nil {
		return fmt.Errorf("connection not initialized")
	}

	marshaled, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	encrypted, err := a.enc.Encrypt(marshaled)
	if err != nil {
		return fmt.Errorf("failed to encrypt message: %w", err)
	}

	a.mu.Lock()
	defer a.mu.Unlock()

	return a.conn.WriteMessage(websocket.TextMessage, encrypted)
}

func (a *AuthServiceImpl) onError(err error) {
	a.conn = nil

	logging.Error("AuthWS error: " + err.Error())

	if closeErr, ok := err.(*websocket.CloseError); ok && closeErr.Code > 3000 {
		a.logError(err, closeErr.Code)
		os.Exit(1)
	}

	a.errors++

	if time.Since(a.errorStart) > errorWindow {
		a.errors = 1
		a.errorStart = time.Now()
	}

	if a.errors > maxErrors && !constants.HOSTING {
		a.logError(err, -1)
		os.Exit(1)
	}

	time.Sleep(reconnectDelay)

	logging.Info("Reconnecting AuthWS...")

	if err := a.connectWS(context.Background()); err != nil {
		a.onError(err)
	}
}

func (a *AuthServiceImpl) logError(err error, code int) {
	switch code {
	case 3004, 3010:
		if !constants.CONSOLE {
			a.uiSvc.LogError("Your LightningBot credentials are invalid! Please re-enter them in the Config tab.")
			break
		}

		logging.Error("Your LightningBot credentials are invalid!")

		username, password := setup.GetLogin()

		cfg, err := config.GetConfig()
		if err != nil {
			logging.Error("Error loading config: " + err.Error())
			os.Exit(0)
		}

		cfg.Account.LightningBot.Username = username
		cfg.Account.LightningBot.Password = password

		err = cfg.Save()
		if err != nil {
			logging.Error("Error saving config: " + err.Error())
			os.Exit(1)
		}

		_ = system.RestartExecution()
	case 3011:
		a.uiSvc.LogError("Your free trial has expired! Please purchase a license key to continue using LightningBot.")
	case 3003:
		a.uiSvc.LogError("AuthWS Error 3003: " + err.Error())
	case 3006:
		a.uiSvc.LogError("Another selfbot has connected using your LightningBot account! This instance will now exit.")
	default:
		a.uiSvc.LogError("An unexpected error has occurred while connecting to the websocket: " + err.Error())
	}
}
