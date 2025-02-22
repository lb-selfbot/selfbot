package interfaces

import (
	"github.com/LightningDev1/discordgo"
)

type SessionsService interface {
	Inject(uiSvc UIService, discordService DiscordService)

	Register()

	GetGatewaySessions() []*discordgo.GatewaySession
	GetLoginSessions() []*discordgo.LoginSession
	GetStatus() string
	SetStatus(status string)
}
