package interfaces

import "context"

type MultiTokenSniperService interface {
	Inject(uiSvc UIService, o11ySvc O11yService, discordSvc DiscordService, eventReactionsSvc EventReactionsService)

	Init(context context.Context) error
}
