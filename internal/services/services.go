package services

import (
	"github.com/lb-selfbot/selfbot/internal/services/authsync"
	"github.com/lb-selfbot/selfbot/internal/services/backup"
	"github.com/lb-selfbot/selfbot/internal/services/config"
	"github.com/lb-selfbot/selfbot/internal/services/console"
	"github.com/lb-selfbot/selfbot/internal/services/discord"
	"github.com/lb-selfbot/selfbot/internal/services/embed"
	"github.com/lb-selfbot/selfbot/internal/services/eventbroker"
	"github.com/lb-selfbot/selfbot/internal/services/eventreactions"
	"github.com/lb-selfbot/selfbot/internal/services/extensions"
	"github.com/lb-selfbot/selfbot/internal/services/ipc"
	"github.com/lb-selfbot/selfbot/internal/services/mts"
	"github.com/lb-selfbot/selfbot/internal/services/o11y"
	"github.com/lb-selfbot/selfbot/internal/services/reminder"
	"github.com/lb-selfbot/selfbot/internal/services/router"
	"github.com/lb-selfbot/selfbot/internal/services/rpc"
	"github.com/lb-selfbot/selfbot/internal/services/sessions"
	"github.com/lb-selfbot/selfbot/internal/services/slashcommand"
	"github.com/lb-selfbot/selfbot/internal/services/spotify"
	"github.com/lb-selfbot/selfbot/internal/services/ui"
	"github.com/lb-selfbot/selfbot/internal/services/variables"
)

// Services is the container for all service instances.
type Services struct {
	EventBroker eventbroker.Service
	Config      config.Service
	Backup      backup.Service
	Spotify     spotify.Service

	Console      console.Service
	Router       router.Service
	Sessions     sessions.Service
	Extensions   extensions.Service
	Discord      discord.Service
	O11y         o11y.Service
	IPC          ipc.Service
	Reminder     reminder.Service
	SlashCommand slashcommand.Service
	AuthSync     authsync.Service

	Variables      variables.Service
	Embed          embed.Service
	UI             ui.Service
	EventReactions eventreactions.Service
	MTS            mts.Service
	RPC            rpc.Service
}

// NewServices initializes and wires all services.
func NewServices() *Services {
	eventBroker := eventbroker.NewService()
	cfg := config.NewService()
	backup := backup.NewService()
	spotify := spotify.NewService()

	console := console.NewService(eventBroker)
	router := router.NewService(eventBroker, cfg)
	sessions := sessions.NewService(eventBroker)
	extensions := extensions.NewService(eventBroker)
	discord := discord.NewService(eventBroker, cfg, console)
	o11y := o11y.NewService(eventBroker, cfg)
	ipc := ipc.NewService(eventBroker, cfg, console)
	reminder := reminder.NewService(eventBroker, cfg)
	slashCommand := slashcommand.NewService(eventBroker, cfg)
	authSync := authsync.NewService(eventBroker, cfg, console)

	variables := variables.NewService(eventBroker, router, extensions, spotify)
	embed := embed.NewService(cfg, variables)
	ui := ui.NewService(eventBroker, console, cfg, discord)
	eventReactions := eventreactions.NewService(eventBroker, cfg, router, variables, discord)
	mts := mts.NewService(eventBroker, cfg, eventReactions)
	rpc := rpc.NewService(eventBroker, cfg, sessions, variables)

	return &Services{
		EventBroker: eventBroker,
		Config:      cfg,
		Backup:      backup,
		Spotify:     spotify,

		Console:      console,
		Router:       router,
		Sessions:     sessions,
		Extensions:   extensions,
		Discord:      discord,
		O11y:         o11y,
		IPC:          ipc,
		Reminder:     reminder,
		SlashCommand: slashCommand,
		AuthSync:     authSync,

		Variables:      variables,
		Embed:          embed,
		UI:             ui,
		EventReactions: eventReactions,
		MTS:            mts,
		RPC:            rpc,
	}
}
