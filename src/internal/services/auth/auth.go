package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/LightningDev1/LB-Selfbot/internal/random"
	"github.com/LightningDev1/LB-Selfbot/internal/services/interfaces"
	"github.com/LightningDev1/LB-Selfbot/selfbot/constants"
	"github.com/LightningDev1/LB-Selfbot/selfbot/encryption"
	"github.com/gorilla/websocket"
)

const (
	defaultWebsocketURL = "wss://api.lightning-bot.com/api/v3/ws/selfbot"
	maxErrors           = 20
	errorWindowMinutes  = 15
	pingInterval        = 60 * time.Second
	reconnectDelay      = 30 * time.Second
	handshakeDelay      = 100 * time.Millisecond
	nonceLength         = 64
)

type AuthServiceImpl struct {
	uiSvc interfaces.UIService

	conn *websocket.Conn
	mu   sync.RWMutex
	enc  *encryption.AESCBC

	readyChan chan struct{}

	creds interfaces.AuthCredentials
	nonce string

	receivedNonce string
}

func (a *AuthServiceImpl) Inject(uiSvc interfaces.UIService) {
	a.uiSvc = uiSvc
}

func (a *AuthServiceImpl) Connect(context context.Context, creds interfaces.AuthCredentials) error {
	if err := a.setupEncryption(); err != nil {
		return err
	}

	a.readyChan = make(chan struct{}, 1)
	a.creds = creds
	a.nonce = random.String(nonceLength)

	return a.connectWS(context)
}

func (a *AuthServiceImpl) OnReady() <-chan struct{} {
	return a.readyChan
}

func (a *AuthServiceImpl) CheckNonce() {
	isValid := constants.HOSTING || a.receivedNonce == "" || a.receivedNonce == a.nonce
	if isValid {
		return
	}

	// TODO: log security error
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

	// TODO: can this be removevd?
	time.Sleep(handshakeDelay)

	// Get hostname
	hostname := "unknown"
	if constants.HOSTING {
		hostname = "LB Hosting"
	} else {
		hostname, _ = os.Hostname()
	}

	// Send HELLO
	err = a.write(&WSMessage{
		Op: OP_HELLO,
		Data: WSData{
			Username:     a.creds.Username,
			Password:     a.creds.Password,
			ComputerName: hostname,
			Time:         time.Now().Unix(),
			ConsoleText:  "", // TODO: console text IPC
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
			// TODO: error handling
			// a.OnError(err)
			continue
		}

		if err := a.onMessage(msg); err != nil {
			// TODO: error handling
			// a.OnError(err)
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
			if err := a.write(&WSMessage{Op: OP_PING}); err != nil {
				// TODO: error handling
				// a.OnError(err)
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
