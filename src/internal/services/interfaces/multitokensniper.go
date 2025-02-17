package interfaces

import "context"

type MultiTokenSniperService interface {
	Inject(uiSvc UIService, discordSvc DiscordService)

	Init(context context.Context) error
}
