package authsync

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/goccy/go-json"

	"github.com/gorilla/websocket"
	"github.com/lb-selfbot/selfbot/internal/constants"
	"github.com/lb-selfbot/selfbot/internal/encryption"
	"github.com/lb-selfbot/selfbot/internal/logging"
	"github.com/lb-selfbot/selfbot/internal/random"
	"github.com/lb-selfbot/selfbot/internal/services/config"
	"github.com/lb-selfbot/selfbot/internal/services/console"
	"github.com/lb-selfbot/selfbot/internal/services/eventbroker"
	"github.com/lb-selfbot/selfbot/internal/setup"
	"github.com/lb-selfbot/selfbot/internal/system"
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

// Credentials represents the credentials for authentication.
type Credentials struct {
	Username string
	Password string
}

// Service handles the authentication and synchronization process.
type Service interface {
	// Connect establishes a connection using the provided credentials.
	Connect(Credentials) error

	// GetReadyChan returns a channel that signals when the connection is ready.
	GetReadyChan() <-chan struct{}

	// CheckNonce verifies the nonce for authentication.
	CheckNonce()

	// SendSync sends synchronization data to the server.
	SendSync(SyncType, any) error
}

type authSyncService struct {
	eventBroker eventbroker.Service
	configSvc   config.Service
	consoleSvc  console.Service

	conn *websocket.Conn
	mu   sync.RWMutex
	enc  *encryption.AESCBC

	readyChan chan struct{}

	creds Credentials
	nonce string

	receivedNonce string

	// Data for sliding window error handling
	errors     int
	errorStart time.Time
}

// NewService creates a new instance of AuthSyncService.
func NewService(eventBroker eventbroker.Service, configSvc config.Service, consoleSvc console.Service) Service {
	svc := &authSyncService{eventBroker: eventBroker, configSvc: configSvc, consoleSvc: consoleSvc}
	svc.subscribe()
	return svc
}

func (a *authSyncService) subscribe() {
	a.eventBroker.Subscribe(eventbroker.EventCheckNonce, func(any) { a.CheckNonce() })

	a.eventBroker.Subscribe(eventbroker.EventSendConfigUpdateSync, func(payload any) {
		cfg, ok := payload.(*config.Config)
		if !ok {
			return
		}

		if err := a.SendSync(SyncTypeConfigUpdate, cfg); err != nil {
			logging.Error("Failed to send config update sync:", err)
		}
	})

	a.eventBroker.Subscribe(eventbroker.EventSetConsole, func(payload any) {
		console, ok := payload.(string)
		if !ok {
			return
		}

		if err := a.SendSync(SyncTypeConsoleSet, console); err != nil {
			logging.Error("Failed to send console set sync:", err)
		}
	})

	a.eventBroker.Subscribe(eventbroker.EventAppendConsole, func(payload any) {
		html, ok := payload.(string)
		if !ok {
			return
		}

		if err := a.SendSync(SyncTypeConsoleAppend, html); err != nil {
			logging.Error("Failed to send console append sync:", err)
		}
	})

	a.eventBroker.Subscribe(eventbroker.EventClearConsole, func(any) {
		if err := a.SendSync(SyncTypeConsoleClear, nil); err != nil {
			logging.Error("Failed to send console clear sync:", err)
		}
	})
}

func (a *authSyncService) Connect(creds Credentials) error {
	// Setup encryption
	if err := a.setupEncryption(); err != nil {
		return err
	}

	// Setup connection data
	a.readyChan = make(chan struct{}, 1)
	a.creds = creds
	a.nonce = random.String(nonceLength)

	return a.connectWS(context.Background())
}

func (a *authSyncService) GetReadyChan() <-chan struct{} {
	return a.readyChan
}

func (a *authSyncService) CheckNonce() {
	// Return if nonce is valid
	if constants.HOSTING || a.nonce == "" || a.nonce == a.receivedNonce {
		return
	}

	fmt.Println("An AuthSync error has occurred. Error code: " + random.String(8))
	os.Exit(1)
}

func (a *authSyncService) SendSync(syncType SyncType, data any) error {
	var msg *WSMessage

	switch syncType {
	case SyncTypeConfigUpdate:
		cfg, ok := data.(*config.Config)
		if !ok {
			return fmt.Errorf("invalid data type for config update")
		}

		// Remove sensitive data
		cfgCopy := *cfg
		cfgCopy.Account.LightningBot.Password = ""
		cfgCopy.Account.Discord.Token = ""

		msg = &WSMessage{
			Op:   OpCodeConfigUpdate,
			Data: WSData{Config: &cfgCopy},
		}

	case SyncTypeConsoleAppend:
		html, ok := data.(string)
		if !ok {
			return fmt.Errorf("invalid data type for console append")
		}

		msg = &WSMessage{
			Op:   OpCodeConsoleAppend,
			Data: WSData{Text: html},
		}

	case SyncTypeConsoleClear:
		msg = &WSMessage{
			Op: OpCodeConsoleClear,
		}

	case SyncTypeConsoleSet:
		html, ok := data.(string)
		if !ok {
			return fmt.Errorf("invalid data type for console set")
		}
		msg = &WSMessage{
			Op:   OpCodeConsoleSet,
			Data: WSData{Text: html},
		}

	default:
		return fmt.Errorf("unsupported sync type: %d", syncType)
	}

	if err := a.write(msg); err != nil {
		return fmt.Errorf("failed to send sync message: %w", err)
	}

	return nil
}

func (a *authSyncService) setupEncryption() error {
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

func (a *authSyncService) connectWS(ctx context.Context) error {
	header := http.Header{"origin": []string{"http://localhost/"}}

	var err error
	a.conn, _, err = websocket.DefaultDialer.DialContext(ctx, defaultWebsocketURL, header)
	if err != nil {
		return err
	}

	var hostname string
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
			ConsoleText:  a.consoleSvc.GetHTMLConsole(),
			Nonce:        a.nonce,
		},
	})
	if err != nil {
		return fmt.Errorf("failed to send HELLO: %w", err)
	}

	go a.listen()
	go a.ping(ctx)

	return nil
}

func (a *authSyncService) read() (*WSMessage, error) {
	if a.conn == nil || a.enc == nil {
		return nil, fmt.Errorf("connection not initialized")
	}

	a.mu.RLock()
	defer a.mu.RUnlock()

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

func (a *authSyncService) write(msg *WSMessage) error {
	if a.enc == nil {
		return fmt.Errorf("encryption not initialized")
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	encrypted, err := a.enc.Encrypt(data)
	if err != nil {
		return fmt.Errorf("failed to encrypt message: %w", err)
	}

	a.mu.Lock()
	defer a.mu.Unlock()

	if a.conn == nil {
		return fmt.Errorf("cannot write message: connection is nil")
	}

	if err := a.conn.WriteMessage(websocket.TextMessage, encrypted); err != nil {
		return fmt.Errorf("failed to write message: %w", err)
	}

	return nil
}

func (a *authSyncService) listen() {
	for {
		msg, err := a.read()
		if err != nil {
			a.onError(err)
			return
		}

		if err := a.handleMessage(msg); err != nil {
			a.onError(err)
			return
		}
	}
}

func (a *authSyncService) ping(ctx context.Context) {
	ticker := time.NewTicker(pingInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := a.write(&WSMessage{Op: OpCodePing}); err != nil {
				a.onError(err)
				return
			}
		case <-ctx.Done():
			return
		}
	}
}

func (a *authSyncService) onError(err error) {
	if a.conn != nil {
		_ = a.conn.Close()
		a.conn = nil
	}

	// Check if the error was induced by the AuthSync server
	var closeErr *websocket.CloseError
	if errors.As(err, &closeErr) && closeErr.Code > 3000 {
		a.logError(err, closeErr.Code)
		os.Exit(1)
	}

	// Sliding window error handling
	if time.Since(a.errorStart) > errorWindow {
		a.errors = 0
		a.errorStart = time.Now()
	}

	a.errors++

	if a.errors > maxErrors && !constants.HOSTING {
		a.logError(err, -1)
		os.Exit(1)
	}

	// Reconnecting
	time.Sleep(reconnectDelay)

	logging.Info("Reconnecting AuthSync...")

	if err := a.connectWS(context.Background()); err != nil {
		a.onError(err)
	}
}

func (a *authSyncService) logError(err error, code int) {
	switch code {
	case 3004, 3010:
		if !constants.CONSOLE {
			a.eventBroker.Publish(eventbroker.EventLogError, "Your LightningBot credentials are invalid! Please re-enter them in the Config tab.")
			break
		}

		logging.Error("Your LightningBot credentials are invalid!")

		username, password := setup.GetLogin()

		cfg := a.configSvc.GetConfig()
		cfg.Account.LightningBot.Username = username
		cfg.Account.LightningBot.Password = password

		err = a.configSvc.SaveConfig(cfg)
		if err != nil {
			logging.Error("Error saving config: " + err.Error())
			os.Exit(1)
		}

		_ = system.RestartExecution()
	case 3011:
		a.eventBroker.Publish(eventbroker.EventLogError, "Your free trial has expired! Please purchase a license key to continue using LightningBot.")
	case 3003:
		a.eventBroker.Publish(eventbroker.EventLogError, "AuthWS Error 3003: "+err.Error())
	case 3006:
		a.eventBroker.Publish(eventbroker.EventLogError, "Another selfbot has connected using your LightningBot account! This instance will now exit.")
	default:
		a.eventBroker.Publish(eventbroker.EventLogError, "An unexpected error has occurred while connecting to the websocket: "+err.Error())
	}
}
