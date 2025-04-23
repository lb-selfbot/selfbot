package ipc

import (
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/lb-selfbot/selfbot/internal/constants"
	"github.com/lb-selfbot/selfbot/internal/services/config"
	"github.com/lb-selfbot/selfbot/internal/services/console"
	"github.com/lb-selfbot/selfbot/internal/services/console/consoledefs"
	"github.com/lb-selfbot/selfbot/internal/services/eventbroker"
)

const (
	defaultIPCAddr = "127.0.0.1:13249"
	htmlCutOff     = 1e8
)

// Service handles the IPC with the UI process.
type Service interface {
	// Connect establishes a connection with the IPC client.
	Connect() error

	// GetReadyChan returns a channel that signals when the connection is ready.
	GetReadyChan() <-chan struct{}
}

type ipcService struct {
	eventBroker eventbroker.Service
	configSvc   config.Service
	consoleSvc  console.Service

	upgrader *websocket.Upgrader
	conn     *websocket.Conn
	mu       sync.Mutex

	readyChan chan struct{}

	sentRawConsole bool
}

// NewService creates a new IPCService instance.
func NewService(eventBroker eventbroker.Service, configSvc config.Service, consoleSvc console.Service) Service {
	svc := &ipcService{eventBroker: eventBroker, configSvc: configSvc, consoleSvc: consoleSvc}
	svc.subscribe()
	return svc
}

func (i *ipcService) subscribe() {
	i.eventBroker.Subscribe(eventbroker.EventRestart, func(any) {
		_ = i.write(&WSMessage{Op: OpCodeRestart})
	})

	i.eventBroker.Subscribe(eventbroker.EventLogError, func(payload any) {
		text, ok := payload.(string)
		if !ok {
			return
		}

		_ = i.write(&WSMessage{Op: OpCodeError, Data: WSData{Text: text}})
	})

	i.eventBroker.Subscribe(eventbroker.EventSetConsole, func(payload any) {
		console, ok := payload.(consoledefs.Set)
		if !ok {
			return
		}

		if len(console.HTML) < htmlCutOff {
			_ = i.write(&WSMessage{Op: OpCodeConsoleSet, Data: WSData{Text: console.HTML}})

			i.sentRawConsole = false
		} else {
			rawConsole := constants.ANSI_REGEX.ReplaceAllString(console.ANSI, "")
			_ = i.write(&WSMessage{Op: OpCodeConsoleSet, Data: WSData{Text: rawConsole, IsRaw: true}})

			i.sentRawConsole = true
		}
	})

	i.eventBroker.Subscribe(eventbroker.EventAppendConsole, func(payload any) {
		console, ok := payload.(consoledefs.Append)
		if !ok {
			return
		}

		if len(console.HTML) < htmlCutOff {
			_ = i.write(&WSMessage{Op: OpCodeConsoleAppend, Data: WSData{Text: console.HTML}})
		} else if i.sentRawConsole {
			rawConsole := constants.ANSI_REGEX.ReplaceAllString(console.ANSI, "")
			_ = i.write(&WSMessage{Op: OpCodeConsoleAppend, Data: WSData{Text: rawConsole, IsRaw: true}})
		} else {
			rawConsole := constants.ANSI_REGEX.ReplaceAllString(i.consoleSvc.GetANSIConsole(), "")
			_ = i.write(&WSMessage{Op: OpCodeConsoleAppend, Data: WSData{Text: rawConsole, IsRaw: true}})

			i.sentRawConsole = true
		}
	})

	i.eventBroker.Subscribe(eventbroker.EventClearConsole, func(any) {
		_ = i.write(&WSMessage{Op: OpCodeConsoleClear})

		i.sentRawConsole = false
	})
}

func (i *ipcService) Connect() error {
	i.mu.Lock()
	defer i.mu.Unlock()

	i.upgrader = &websocket.Upgrader{}
	i.readyChan = make(chan struct{}, 1)

	mux := http.NewServeMux()
	mux.HandleFunc("/", i.handleRequest)
	go func() {
		log.Fatal(http.ListenAndServe(defaultIPCAddr, mux))
	}()

	return nil
}

func (i *ipcService) GetReadyChan() <-chan struct{} {
	return i.readyChan
}

func (i *ipcService) handleRequest(w http.ResponseWriter, r *http.Request) {
	if i.conn != nil {
		return
	}

	// Upgrade the HTTP connection to a WebSocket connection.
	var err error
	i.conn, err = i.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer i.conn.Close()

	i.readyChan <- struct{}{}

	if err := i.conn.WriteJSON(&WSMessage{Op: OpCodeConnected}); err != nil {
		return
	}

	for {
		message := &WsMessageReceive{}

		err = i.conn.ReadJSON(message)
		if err != nil {
			break
		}

		if err = i.handleMessage(message); err != nil {
			break
		}
	}

	i.conn = nil
	close(i.readyChan)
	i.readyChan = make(chan struct{}, 1)

	if !constants.DEBUG {
		os.Exit(0)
	}
}

func (i *ipcService) write(message *WSMessage) error {
	if i.conn == nil {
		return nil
	}

	return i.conn.WriteJSON(message)
}
