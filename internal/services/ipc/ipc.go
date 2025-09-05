package ipc

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/goccy/go-json"

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

	zlibWriter *zlib.Writer

	mu sync.Mutex

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
		_ = i.write(&IPCMessage{Op: OpCodeRestart})
	})

	i.eventBroker.Subscribe(eventbroker.EventLogError, func(payload any) {
		text, ok := payload.(string)
		if !ok {
			return
		}

		_ = i.write(&IPCMessage{Op: OpCodeError, Data: &IPCData{Text: text}})
	})

	i.eventBroker.Subscribe(eventbroker.EventSetConsole, func(payload any) {
		console, ok := payload.(consoledefs.Set)
		if !ok {
			return
		}

		if len(console.HTML) < htmlCutOff {
			_ = i.write(&IPCMessage{Op: OpCodeConsoleSet, Data: &IPCData{Text: console.HTML}})

			i.sentRawConsole = false
		} else {
			rawConsole := constants.ANSI_REGEX.ReplaceAllString(console.ANSI, "")
			_ = i.write(&IPCMessage{Op: OpCodeConsoleSet, Data: &IPCData{Text: rawConsole, IsRaw: true}})

			i.sentRawConsole = true
		}
	})

	i.eventBroker.Subscribe(eventbroker.EventAppendConsole, func(payload any) {
		console, ok := payload.(consoledefs.Append)
		if !ok {
			return
		}

		if len(console.HTML) < htmlCutOff {
			_ = i.write(&IPCMessage{Op: OpCodeConsoleAppend, Data: &IPCData{Text: console.HTML}})
		} else if i.sentRawConsole {
			rawConsole := constants.ANSI_REGEX.ReplaceAllString(console.ANSI, "")
			_ = i.write(&IPCMessage{Op: OpCodeConsoleAppend, Data: &IPCData{Text: rawConsole, IsRaw: true}})
		} else {
			rawConsole := constants.ANSI_REGEX.ReplaceAllString(i.consoleSvc.GetANSIConsole(), "")
			_ = i.write(&IPCMessage{Op: OpCodeConsoleAppend, Data: &IPCData{Text: rawConsole, IsRaw: true}})

			i.sentRawConsole = true
		}
	})

	i.eventBroker.Subscribe(eventbroker.EventClearConsole, func(any) {
		_ = i.write(&IPCMessage{Op: OpCodeConsoleClear})

		i.sentRawConsole = false
	})
}

func (i *ipcService) Connect() (err error) {
	i.mu.Lock()
	defer i.mu.Unlock()

	i.zlibWriter = zlib.NewWriter(os.Stdout)

	i.readyChan = make(chan struct{}, 1)

	go i.listen()

	return i.write(&IPCMessage{Op: OpCodePing})
}

func (i *ipcService) GetReadyChan() <-chan struct{} {
	return i.readyChan
}

func (i *ipcService) listen() {
	buf := make([]byte, 4096)
	var zlibBuf []byte

	for {
		n, err := os.Stdin.Read(buf)
		if err != nil {
			if err != io.EOF {
				fmt.Fprintf(os.Stderr, "Error reading from stdin: %v\n", err)
			}
			os.Exit(1)
		}

		zlibBuf = append(zlibBuf, buf[:n]...)

		for len(zlibBuf) > 0 {
			reader, err := zlib.NewReader(bytes.NewReader(zlibBuf))
			if err != nil {
				break
			}

			var decompressed bytes.Buffer
			_, err = io.Copy(&decompressed, reader)
			reader.Close()

			if err != nil {
				break
			}

			var msg *IPCMessage
			err = json.Unmarshal(decompressed.Bytes(), &msg)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Invalid JSON: %v\n", err)
				zlibBuf = nil
				break
			}

			err = i.handleMessage(msg)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error handling message: %v\n", err)
				os.Exit(1)
			}

			zlibBuf = nil
			break
		}
	}
}

func (i *ipcService) write(message *IPCMessage) error {
	if i.zlibWriter == nil {
		return nil
	}

	if err := json.NewEncoder(i.zlibWriter).Encode(message); err != nil {
		return err
	}

	return i.zlibWriter.Flush()
}
