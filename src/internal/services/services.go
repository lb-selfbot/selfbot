package services

import (
	"github.com/LightningDev1/LB-Selfbot/internal/services/auth"
	"github.com/LightningDev1/LB-Selfbot/internal/services/discord"
	"github.com/LightningDev1/LB-Selfbot/internal/services/interfaces"
	"github.com/LightningDev1/LB-Selfbot/internal/services/ipc"
	"github.com/LightningDev1/LB-Selfbot/internal/services/ui"
)

type Services struct {
	Auth    interfaces.AuthService
	UI      interfaces.UIService
	IPC     interfaces.IPCService
	Discord interfaces.DiscordService
}

func NewServices() *Services {
	services := &Services{}
	services.register()
	return services
}

func (s *Services) register() {
	// Create services
	s.Auth = &auth.AuthServiceImpl{}
	s.UI = &ui.UIServiceImpl{}
	s.IPC = &ipc.IPCServiceImpl{}
	s.Discord = &discord.DiscordServiceImpl{}

	// Inject dependencies
	s.Auth.Inject(s.UI, s.IPC)
	s.UI.Inject(s.Auth, s.IPC)
	s.Discord.Inject(s.Auth, s.UI)
}
