package interfaces

type IPCService interface {
	Serve() error
	OnConnect() <-chan struct{}

	ClearConsole() error
	SetConsole(ansi, html string) error
	AppendConsole(ansi, html string) error

	SendError(text string) error
	SendRestart() error
}
