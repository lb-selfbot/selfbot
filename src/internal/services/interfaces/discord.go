package interfaces

import (
	"github.com/LightningDev1/dgc"
	"github.com/LightningDev1/discordgo"
)

type DiscordService interface {
	Inject(authSvc AuthService, uiSvc UIService, extensionsSvc ExtensionsService)

	Init(token string, verbose bool) error

	GetSession() *discordgo.Session
	GetRouter() *dgc.Router

	Open() error
}
