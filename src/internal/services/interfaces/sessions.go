package interfaces

import (
	"time"

	"github.com/LightningDev1/discordgo"
)

type SessionsService interface {
	Inject(uiSvc UIService, discordService DiscordService)

	Register()

	GetGatewaySessions() []*discordgo.GatewaySession
	GetLoginSessions() []*discordgo.LoginSession
	SetLastStatusChange(time time.Time)
	GetStatus() string
}
