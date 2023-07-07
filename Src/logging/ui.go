package logging

import (
	"fmt"
	"path/filepath"

	"github.com/LightningDev1/LB-Selfbot/config"
	"github.com/LightningDev1/LB-Selfbot/ipc"
	"github.com/LightningDev1/LB-Selfbot/utils"
	"github.com/LightningDev1/LB-Selfbot/websocket"
	"github.com/LightningDev1/beeep"
)

type UILogging struct{}

func (UILogging) Notification(message string) {
	cfg, err := config.GetConfig()
	if err != nil {
		Console.Error("Error getting config:", err)
		return
	}

	if !cfg.Notifications.Enabled {
		return
	}

	go beeep.Notify(
		"LightningBot",
		message,
		"LightningBot",
		filepath.Join(utils.IO.GetMainFolder(), "Logo.png"),
	)
}

func (UILogging) Error(errText string) {
	utils.Logging.Error(errText)

	ipc.IPC.Write(ipc.JsonEvent{
		Op: ipc.OP_Error,
		Data: ipc.JsonEventData{
			Text: errText,
		},
	})

	UI.Append(Console.Colored(fmt.Sprintf("[ERROR] %s", errText), 255, 100, 100, true))
}

func (UILogging) Append(text string) {
	text = "\n" + utils.Logging.AnsiToHtml(text)

	ipc.IPC.ConsoleAppend(text)
	websocket.Websocket.Send(websocket.WsMessage{
		Op: "CONSOLE_APPEND",
		Data: websocket.WsData{
			Text: text,
		},
	})
}

func (UILogging) Clear() {
	ipc.IPC.ConsoleClear()
	websocket.Websocket.Send(websocket.WsMessage{
		Op:   "CONSOLE_CLEAR",
		Data: websocket.WsData{},
	})
}

func (UILogging) Set(text string) {
	text = utils.Logging.AnsiToHtml(text)

	ipc.IPC.ConsoleSet(text)
	websocket.Websocket.Send(websocket.WsMessage{
		Op: "CONSOLE_SET",
		Data: websocket.WsData{
			Text: text,
		},
	})
}

var UI = UILogging{}
