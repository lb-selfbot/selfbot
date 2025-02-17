package ipc

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/LightningDev1/LB-Selfbot/internal/constants"
	"github.com/gorilla/websocket"
)

const (
	defaultIPCAddr = "127.0.0.1:13249"
	htmlCutOff     = 1e8
)

type IPCServiceImpl struct {
	upgrader *websocket.Upgrader
	conn     *websocket.Conn
	mu       sync.RWMutex

	connectChan chan struct{}

	htmlConsole    string
	rawConsole     string
	sentRawConsole bool
}

func (i *IPCServiceImpl) Serve() error {
	i.upgrader = &websocket.Upgrader{}
	i.connectChan = make(chan struct{}, 1)

	mux := http.NewServeMux()
	mux.HandleFunc("/", i.handleRequest)
	go log.Fatal(http.ListenAndServe(defaultIPCAddr, mux))

	return nil
}

func (i *IPCServiceImpl) OnConnect() <-chan struct{} {
	return i.connectChan
}

func (i *IPCServiceImpl) GetHTMLConsole() string { return i.htmlConsole }
func (i *IPCServiceImpl) GetRawConsole() string  { return i.rawConsole }

func (i *IPCServiceImpl) handleRequest(w http.ResponseWriter, r *http.Request) {
	var err error
	i.conn, err = i.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer i.conn.Close()

	i.connectChan <- struct{}{}

	for {
		message := &WSMessage{}

		i.mu.RLock()
		err = i.conn.ReadJSON(message)
		i.mu.RUnlock()

		if err != nil {
			break
		}

		if err = i.onMessage(message); err != nil {
			break
		}
	}

	i.conn = nil
	close(i.connectChan)
	i.connectChan = make(chan struct{}, 1)

	if !constants.DEBUG {
		os.Exit(0)
	}
}

func (i *IPCServiceImpl) write(data *WSMessage) error {
	if i.conn == nil {
		return fmt.Errorf("connection not initialized")
	}

	i.mu.Lock()
	defer i.mu.Unlock()

	err := i.conn.WriteJSON(data)
	if err != nil {
		return fmt.Errorf("failed to write message: %w", err)
	}

	return nil
}
