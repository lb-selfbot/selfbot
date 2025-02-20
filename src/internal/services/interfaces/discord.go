package interfaces

import (
	"github.com/LightningDev1/LB-Selfbot/internal/router"
	"github.com/LightningDev1/discordgo"
)

type DiscordService interface {
	Inject(authSvc AuthService, uiSvc UIService, o11ySvc O11yService, eventReactionsSvc EventReactionsService, extensionsSvc ExtensionsService)

	Init(token string, verbose bool) error

	GetSession() *discordgo.Session
	GetRouter() router.Router

	Open() error
}
