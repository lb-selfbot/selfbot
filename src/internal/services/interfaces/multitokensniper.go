package interfaces

import "context"

type MultiTokenSniperService interface {
	Inject(uiSvc UIService, o11ySvc O11yService, discordSvc DiscordService)

	Init(context context.Context) error
}
