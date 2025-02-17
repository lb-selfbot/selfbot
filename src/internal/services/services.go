package services

import (
	"github.com/LightningDev1/LB-Selfbot/internal/services/auth"
	"github.com/LightningDev1/LB-Selfbot/internal/services/backup"
	"github.com/LightningDev1/LB-Selfbot/internal/services/discord"
	"github.com/LightningDev1/LB-Selfbot/internal/services/extensions"
	"github.com/LightningDev1/LB-Selfbot/internal/services/interfaces"
	"github.com/LightningDev1/LB-Selfbot/internal/services/ipc"
	"github.com/LightningDev1/LB-Selfbot/internal/services/multitokensniper"
	"github.com/LightningDev1/LB-Selfbot/internal/services/reminder"
	"github.com/LightningDev1/LB-Selfbot/internal/services/rpc"
	"github.com/LightningDev1/LB-Selfbot/internal/services/sessions"
	"github.com/LightningDev1/LB-Selfbot/internal/services/spotify"
	"github.com/LightningDev1/LB-Selfbot/internal/services/ui"
	"github.com/LightningDev1/LB-Selfbot/internal/services/variables"
)

type Services struct {
	Auth             interfaces.AuthService
	UI               interfaces.UIService
	IPC              interfaces.IPCService
	Discord          interfaces.DiscordService
	Extensions       interfaces.ExtensionsService
	Backup           interfaces.BackupService
	MultiTokenSniper interfaces.MultiTokenSniperService
	Reminder         interfaces.ReminderService
	Sessions         interfaces.SessionsService
	RPC              interfaces.RPCService
	Spotify          interfaces.SpotifyService
	Variables        interfaces.VariablesService
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
	s.Extensions = &extensions.ExtensionsServiceImpl{}
	s.Backup = &backup.BackupServiceImpl{}
	s.MultiTokenSniper = &multitokensniper.MultiTokenSniperServiceImpl{}
	s.Reminder = &reminder.ReminderServiceImpl{}
	s.Sessions = &sessions.SessionsServiceImpl{}
	s.RPC = &rpc.RPCServiceImpl{}
	s.Spotify = &spotify.SpotifyServiceImpl{}
	s.Variables = &variables.VariablesServiceImpl{}

	// Inject dependencies
	s.Auth.Inject(s.UI, s.IPC)
	s.UI.Inject(s.Auth, s.IPC)
	s.Discord.Inject(s.Auth, s.UI, s.Extensions)
	s.Extensions.Inject(s.UI, s.Discord)
	s.Backup.Inject(s.Discord)
	s.MultiTokenSniper.Inject(s.UI, s.Discord)
	s.Reminder.Inject(s.UI)
	s.Sessions.Inject(s.UI, s.Discord)
	s.RPC.Inject(s.UI, s.Discord, s.Sessions)
	s.Spotify.Inject(s.Discord)
	s.Variables.Inject(s.Discord, s.Extensions, s.Spotify)
}
