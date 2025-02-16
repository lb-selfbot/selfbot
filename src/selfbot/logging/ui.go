package logging

import (
	"fmt"
	"path/filepath"
	"time"

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

func (UILogging) Error(text string) {
	utils.Logging.Error(text)

	ipc.IPC.Write(ipc.JsonEvent{
		Op:   ipc.OP_Error,
		Data: ipc.JsonEventData{Text: text},
	})

	time := time.Now().Local().Format("2006-01-02 15:04:05")
	text = fmt.Sprintf("[%s] [ERROR] %s", time, text)

	UI.Append(Console.Colored(text, 255, 100, 100, true))
}

func (UILogging) Info(text string) {
	utils.Logging.Info(text)

	time := time.Now().Local().Format("2006-01-02 15:04:05")
	text = fmt.Sprintf("[%s] [INFO] %s", time, text)

	UI.Append(Console.Colored(text, 193, 239, 255, true))
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
