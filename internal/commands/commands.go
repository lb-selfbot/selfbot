package commands

import (
	"github.com/lb-selfbot/selfbot/internal/commands/activity"
	"github.com/lb-selfbot/selfbot/internal/commands/admin"
	"github.com/lb-selfbot/selfbot/internal/commands/ai"
	"github.com/lb-selfbot/selfbot/internal/commands/backup"
	"github.com/lb-selfbot/selfbot/internal/commands/config"
	"github.com/lb-selfbot/selfbot/internal/commands/currency"
	"github.com/lb-selfbot/selfbot/internal/commands/dumping"
	"github.com/lb-selfbot/selfbot/internal/commands/eventreactions"
	"github.com/lb-selfbot/selfbot/internal/commands/fun"
	"github.com/lb-selfbot/selfbot/internal/commands/game"
	"github.com/lb-selfbot/selfbot/internal/commands/help"
	"github.com/lb-selfbot/selfbot/internal/commands/image"
	"github.com/lb-selfbot/selfbot/internal/commands/memegen"
	"github.com/lb-selfbot/selfbot/internal/commands/nsfw"
	"github.com/lb-selfbot/selfbot/internal/commands/protections"
	"github.com/lb-selfbot/selfbot/internal/commands/reddit"
	"github.com/lb-selfbot/selfbot/internal/commands/sessions"
	"github.com/lb-selfbot/selfbot/internal/commands/sharing"
	"github.com/lb-selfbot/selfbot/internal/commands/slashcommands"
	"github.com/lb-selfbot/selfbot/internal/commands/spotify"
	"github.com/lb-selfbot/selfbot/internal/commands/text"
	"github.com/lb-selfbot/selfbot/internal/commands/trolling"
	"github.com/lb-selfbot/selfbot/internal/commands/utility"
	"github.com/lb-selfbot/selfbot/internal/services"
)

type CmdFlag int

const (
	CmdFlagDangerous CmdFlag = iota
	CmdFlagHostingDangerous
)

type CmdRegistry interface {
	Register(services *services.Services)
}

type cmdRegistry struct{}

func NewCmdRegistry() CmdRegistry {
	return &cmdRegistry{}
}

func (s *cmdRegistry) Register(services *services.Services) {
	// Register all command categories
	activity.NewActivityCommands(services).Register()
	admin.NewAdminCommands(services).Register()
	ai.NewAICommands(services).Register()
	backup.NewBackupCommands(services).Register()
	config.NewConfigCommands(services).Register()
	currency.NewCurrencyCommands(services).Register()
	dumping.NewDumpingCommands(services).Register()
	eventreactions.NewEventReactionsCommands(services).Register()
	fun.NewFunCommands(services).Register()
	game.NewGameCommands(services).Register()
	image.NewImageCommands(services).Register()
	memegen.NewMemeGenCommands(services).Register()
	nsfw.NewNSFWCommands(services).Register()
	protections.NewProtectionsCommands(services).Register()
	reddit.NewRedditCommands(services).Register()
	sessions.NewSessionsCommands(services).Register()
	sharing.NewSharingCommands(services).Register()
	slashcommands.NewSlashCommands(services).Register()
	spotify.NewSpotifyCommands(services).Register()
	text.NewTextCommands(services).Register()
	trolling.NewTrollingCommands(services).Register()
	utility.NewUtilityCommands(services).Register()

	// Register help commands last so they can access all commands
	help.NewHelpCommands(services).Register()
}
