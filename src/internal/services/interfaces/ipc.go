package interfaces

import "github.com/LightningDev1/LB-Selfbot/internal/config"

type IPCService interface {
	Serve() error
	OnConnect() <-chan struct{}

	GetHTMLConsole() string
	GetRawConsole() string

	ClearConsole() error
	SetConsole(ansi, html string) error
	AppendConsole(ansi, html string) error

	SendError(text string) error
	SendRestart() error
	SendConfigUpdate(cfg *config.Config) error
}
