package logging

import (
	"fmt"
	"path/filepath"

	"github.com/LightningDev1/LB-Selfbot/selfbot/config"
	"github.com/LightningDev1/LB-Selfbot/selfbot/ipc"
	"github.com/LightningDev1/LB-Selfbot/selfbot/utils"
	"github.com/LightningDev1/LB-Selfbot/selfbot/websocket"
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

	go func() {
		_ = beeep.Notify(
			"LightningBot",
			message,
			"LightningBot",
			filepath.Join(utils.IO.GetMainFolder(), "Logo.png"),
		)
	}()
}

func (UILogging) Error(errText string) {
	utils.Logging.Error(errText)

	ipc.IPC.Write(ipc.JsonEvent{
		Op:   ipc.OP_Error,
		Data: ipc.JsonEventData{Text: errText},
	})

	UI.Append(Console.Colored(fmt.Sprintf("[ERROR] %s", errText), 255, 100, 100, true))
}

func (UILogging) Append(ansi string) {
	html := "\n" + utils.Logging.AnsiToHtml(ansi)

	ipc.IPC.ConsoleAppend(ansi, html)
	_ = websocket.Websocket.Send(websocket.WsMessage{
		Op:   "CONSOLE_APPEND",
		Data: websocket.WsData{Text: html},
	})
}

func (UILogging) Clear() {
	ipc.IPC.ConsoleClear()
	_ = websocket.Websocket.Send(websocket.WsMessage{
		Op:   "CONSOLE_CLEAR",
		Data: websocket.WsData{},
	})
}

func (UILogging) Set(ansi string) {
	html := utils.Logging.AnsiToHtml(ansi)

	ipc.IPC.ConsoleSet(ansi, html)
	_ = websocket.Websocket.Send(websocket.WsMessage{
		Op:   "CONSOLE_SET",
		Data: websocket.WsData{Text: html},
	})
}

var UI = UILogging{}
