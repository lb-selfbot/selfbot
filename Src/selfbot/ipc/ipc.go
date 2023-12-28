package ipc

import (
	"log"
	"net/http"
	"os"

	"github.com/LightningDev1/LB-Selfbot/selfbot/config"
	"github.com/LightningDev1/LB-Selfbot/selfbot/constants"
	"github.com/LightningDev1/LB-Selfbot/selfbot/utils"
	"github.com/gorilla/websocket"
)

var htmlCutoff int = 1e8

type JsonEventData struct {
	Text       string         `json:"text,omitempty"`
	IsRaw      bool           `json:"is_raw,omitempty"`
	MainFolder string         `json:"mainFolder,omitempty"`
	Config     *config.Config `json:"config,omitempty"`
}

type JsonEvent struct {
	Op   int           `json:"op"`
	Data JsonEventData `json:"data"`
}

var (
	OP_Connected     = 98
	OP_Hello         = 99
	OP_Exit          = 100
	OP_Restart       = 101
	OP_Notification  = 102
	OP_InvalidToken  = 103
	OP_Error         = 104
	OP_ConfigUpdate  = 105
	OP_NewToken      = 106
	OP_ConsoleAppend = 107
	OP_ConsoleClear  = 108
	OP_ConsoleSet    = 109
)

var HTMLConsole string
var RawConsole string

var ConfigUpdateCallbacks = []func(*config.Config){}

type Ipc struct {
	Upgrader       *websocket.Upgrader
	Conn           *websocket.Conn
	Connected      bool
	Connect        chan struct{}
	SentRawConsole bool
}

func (i *Ipc) Init() {
	i.Upgrader = &websocket.Upgrader{}
	i.Connect = make(chan struct{})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		c, err := i.Upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer c.Close()
		i.Connected = true
		close(i.Connect)

		i.Conn = c
		i.OnConnect()
		for {
			message := JsonEvent{}
			err := c.ReadJSON(&message)
			if err != nil {
				break
			}
			i.OnMessage(&message)
		}
		i.Connected = false
		i.Conn = nil
		i.Connect = make(chan struct{})
		if !constants.DEBUG {
			os.Exit(0)
		}
	})

	go func() {
		log.Fatal(http.ListenAndServe("127.0.0.1:13249", nil))
	}()
}

func (i *Ipc) Write(data interface{}) {
	if i.Conn != nil {
		i.Conn.WriteJSON(data)
	}
}

func (i *Ipc) OnConnect() {
	i.Write(JsonEvent{
		Op:   OP_Connected,
		Data: JsonEventData{},
	})
}

func (i *Ipc) OnMessage(jsonEvent *JsonEvent) {
	switch jsonEvent.Op {
	case OP_Hello:
		break
	case OP_Exit:
		os.Exit(0)
	case OP_ConfigUpdate:
		config.SaveConfigFromUI(jsonEvent.Data.Config)

		for _, callback := range ConfigUpdateCallbacks {
			go callback(jsonEvent.Data.Config)
		}
	default:
		utils.Logging.Warn("Unknown IPC message:", jsonEvent)
	}
}

func (i *Ipc) ConsoleSet(ansi, html string) {
	HTMLConsole = html
	newRawText := constants.ANSI_REGEX.ReplaceAllString(ansi, "") + "\n"
	RawConsole = newRawText

	if len(HTMLConsole) < htmlCutoff {
		i.Write(JsonEvent{
			Op: OP_ConsoleSet,
			Data: JsonEventData{
				Text: HTMLConsole,
			},
		})
		return
	} else {
		i.Write(JsonEvent{
			Op: OP_ConsoleSet,
			Data: JsonEventData{
				Text:  RawConsole,
				IsRaw: true,
			},
		})
		i.SentRawConsole = true
	}
}

func (i *Ipc) ConsoleAppend(ansi, html string) {
	HTMLConsole += html + "\n"
	newRawText := constants.ANSI_REGEX.ReplaceAllString(ansi, "") + "\n"
	RawConsole += newRawText

	if len(HTMLConsole) < htmlCutoff {
		i.Write(JsonEvent{
			Op: OP_ConsoleAppend,
			Data: JsonEventData{
				Text: html,
			},
		})
		return
	}

	if i.SentRawConsole {
		i.Write(JsonEvent{
			Op: OP_ConsoleAppend,
			Data: JsonEventData{
				Text:  newRawText,
				IsRaw: true,
			},
		})
	} else {
		i.Write(JsonEvent{
			Op: OP_ConsoleSet,
			Data: JsonEventData{
				Text:  RawConsole,
				IsRaw: true,
			},
		})
		i.SentRawConsole = true
	}
}

func (i *Ipc) ConsoleClear() {
	i.Write(JsonEvent{
		Op:   OP_ConsoleClear,
		Data: JsonEventData{},
	})
	HTMLConsole = ""
	RawConsole = ""
	i.SentRawConsole = false
}

func (i *Ipc) Restart() {
	i.Write(JsonEvent{
		Op:   OP_Restart,
		Data: JsonEventData{},
	})
	os.Exit(0)
}

func (i *Ipc) WaitForConnection() {
	if i.Connected {
		return
	}
	<-i.Connect
}

var IPC = Ipc{}

func SendNotification(text string) {
	IPC.Write(JsonEvent{
		Op: OP_Notification,
		Data: JsonEventData{
			Text: text,
		},
	})
}

func ConfigUpdate(cfg *config.Config) {
	for _, callback := range ConfigUpdateCallbacks {
		go callback(cfg)
	}

	IPC.Write(JsonEvent{
		Op: OP_ConfigUpdate,
		Data: JsonEventData{
			Config: cfg,
		},
	})
}
