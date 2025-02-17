package ipc

import (
	"github.com/LightningDev1/LB-Selfbot/internal/constants"
	"github.com/LightningDev1/LB-Selfbot/selfbot/config"
)

func (i *IPCServiceImpl) ClearConsole() error {
	i.htmlConsole = ""
	i.rawConsole = ""
	i.sentRawConsole = false

	return i.write(&WSMessage{Op: OP_CONSOLE_CLEAR})
}

func (i *IPCServiceImpl) SetConsole(ansi, html string) error {
	i.htmlConsole = html
	i.rawConsole = constants.ANSI_REGEX.ReplaceAllString(ansi, "") + "\n"

	if len(i.htmlConsole) < htmlCutOff {
		return i.write(&WSMessage{
			Op:   OP_CONSOLE_SET,
			Data: WSData{Text: i.htmlConsole},
		})
	}

	i.sentRawConsole = true
	return i.write(&WSMessage{
		Op: OP_CONSOLE_SET,
		Data: WSData{
			Text:  i.rawConsole,
			IsRaw: true,
		},
	})
}

func (i *IPCServiceImpl) AppendConsole(ansi, html string) error {
	i.htmlConsole += html + "\n"
	i.rawConsole += constants.ANSI_REGEX.ReplaceAllString(ansi, "") + "\n"

	if len(i.htmlConsole) < htmlCutOff {
		return i.write(&WSMessage{
			Op:   OP_CONSOLE_APPEND,
			Data: WSData{Text: html},
		})
	}

	if i.sentRawConsole {
		return i.write(&WSMessage{
			Op: OP_CONSOLE_APPEND,
			Data: WSData{
				Text:  i.rawConsole,
				IsRaw: true,
			},
		})
	}

	i.sentRawConsole = true
	return i.write(&WSMessage{
		Op: OP_CONSOLE_SET,
		Data: WSData{
			Text:  i.rawConsole,
			IsRaw: true,
		},
	})
}

func (i *IPCServiceImpl) SendError(text string) error {
	return i.write(&WSMessage{
		Op:   OP_ERROR,
		Data: WSData{Text: text},
	})
}

func (i *IPCServiceImpl) SendRestart() error {
	return i.write(&WSMessage{Op: OP_RESTART})
}

func (i *IPCServiceImpl) SendConfigUpdate(cfg *config.Config) error {
	// TODO: callbacks

	return i.write(&WSMessage{
		Op:   OP_CONFIG_UPDATE,
		Data: WSData{Config: cfg},
	})
}
