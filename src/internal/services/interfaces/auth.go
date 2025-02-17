package interfaces

import "context"

type AuthCredentials struct {
	Username string
	Password string
}

type AuthService interface {
	Inject(uiSvc UIService, ipcSvc IPCService)

	Connect(context context.Context, creds AuthCredentials) error
	OnReady() <-chan struct{}
	CheckNonce()

	SendConfigUpdate() error
	AppendConsole(html string) error
	ClearConsole() error
	SetConsole(html string) error
}
