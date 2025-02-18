package interfaces

import "github.com/LightningDev1/discordgo"

type O11yService interface {
	Inject(uiSvc UIService)

	Init() error

	ReportError(err any)
	DiscordErrorHandler(s *discordgo.Session, err any)
}
