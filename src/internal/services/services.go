package services

import (
	"github.com/LightningDev1/LB-Selfbot/internal/services/auth"
	"github.com/LightningDev1/LB-Selfbot/internal/services/backup"
	"github.com/LightningDev1/LB-Selfbot/internal/services/discord"
	"github.com/LightningDev1/LB-Selfbot/internal/services/embed"
	"github.com/LightningDev1/LB-Selfbot/internal/services/eventreactions"
	"github.com/LightningDev1/LB-Selfbot/internal/services/extensions"
	"github.com/LightningDev1/LB-Selfbot/internal/services/interfaces"
	"github.com/LightningDev1/LB-Selfbot/internal/services/ipc"
	"github.com/LightningDev1/LB-Selfbot/internal/services/multitokensniper"
	"github.com/LightningDev1/LB-Selfbot/internal/services/o11y"
	"github.com/LightningDev1/LB-Selfbot/internal/services/reminder"
	"github.com/LightningDev1/LB-Selfbot/internal/services/rpc"
	"github.com/LightningDev1/LB-Selfbot/internal/services/scjobs"
	"github.com/LightningDev1/LB-Selfbot/internal/services/sessions"
	"github.com/LightningDev1/LB-Selfbot/internal/services/spotify"
	"github.com/LightningDev1/LB-Selfbot/internal/services/ui"
	"github.com/LightningDev1/LB-Selfbot/internal/services/variables"
)

type Services struct {
	Auth             interfaces.AuthService
	UI               interfaces.UIService
	IPC              interfaces.IPCService
	O11y             interfaces.O11yService
	Discord          interfaces.DiscordService
	Embed            interfaces.EmbedService
	EventReactions   interfaces.EventReactionsService
	Extensions       interfaces.ExtensionsService
	Backup           interfaces.BackupService
	MultiTokenSniper interfaces.MultiTokenSniperService
	Reminder         interfaces.ReminderService
	Sessions         interfaces.SessionsService
	RPC              interfaces.RPCService
	Spotify          interfaces.SpotifyService
	Variables        interfaces.VariablesService
	SCJobs           interfaces.SCJobsService
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
	s.O11y = &o11y.O11yServiceImpl{}
	s.Discord = &discord.DiscordServiceImpl{}
	s.Embed = &embed.EmbedServiceImpl{}
	s.EventReactions = eventreactions.NewEventReactionsServiceImpl()
	s.Extensions = &extensions.ExtensionsServiceImpl{}
	s.Backup = &backup.BackupServiceImpl{}
	s.MultiTokenSniper = &multitokensniper.MultiTokenSniperServiceImpl{}
	s.Reminder = &reminder.ReminderServiceImpl{}
	s.Sessions = &sessions.SessionsServiceImpl{}
	s.RPC = &rpc.RPCServiceImpl{}
	s.Spotify = &spotify.SpotifyServiceImpl{}
	s.Variables = &variables.VariablesServiceImpl{}
	s.SCJobs = &scjobs.SCJobsServiceImpl{}

	// Inject dependencies
	s.Auth.Inject(s.UI, s.IPC)
	s.UI.Inject(s.Auth, s.IPC, s.Discord, s.Embed)
	s.IPC.Inject(s.Auth)
	s.O11y.Inject(s.UI)
	s.Discord.Inject(s.Auth, s.UI, s.O11y, s.EventReactions, s.Extensions)
	s.Embed.Inject(s.Variables)
	s.EventReactions.Inject(s.Auth, s.UI, s.Discord, s.Embed, s.Sessions, s.SCJobs, s.RPC, s.Variables)
	s.Extensions.Inject(s.UI, s.Discord, s.Embed)
	s.Backup.Inject(s.Discord)
	s.MultiTokenSniper.Inject(s.UI, s.O11y, s.Discord, s.EventReactions)
	s.Reminder.Inject(s.UI)
	s.Sessions.Inject(s.UI, s.Discord)
	s.RPC.Inject(s.UI, s.Discord, s.Sessions, s.Variables)
	s.Spotify.Inject(s.Discord)
	s.Variables.Inject(s.Discord, s.EventReactions, s.Extensions, s.Spotify)
	s.SCJobs.Inject(s.UI, s.Discord)
}
