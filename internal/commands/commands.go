package commands

import (
	"github.com/LightningDev1/LB-Selfbot/internal/commands/activity"
	"github.com/LightningDev1/LB-Selfbot/internal/commands/admin"
	"github.com/LightningDev1/LB-Selfbot/internal/commands/ai"
	"github.com/LightningDev1/LB-Selfbot/internal/commands/backup"
	"github.com/LightningDev1/LB-Selfbot/internal/commands/config"
	"github.com/LightningDev1/LB-Selfbot/internal/commands/currency"
	"github.com/LightningDev1/LB-Selfbot/internal/commands/dumping"
	"github.com/LightningDev1/LB-Selfbot/internal/commands/eventreactions"
	"github.com/LightningDev1/LB-Selfbot/internal/commands/fun"
	"github.com/LightningDev1/LB-Selfbot/internal/commands/game"
	"github.com/LightningDev1/LB-Selfbot/internal/commands/help"
	"github.com/LightningDev1/LB-Selfbot/internal/commands/image"
	"github.com/LightningDev1/LB-Selfbot/internal/commands/memegen"
	"github.com/LightningDev1/LB-Selfbot/internal/commands/nsfw"
	"github.com/LightningDev1/LB-Selfbot/internal/commands/protections"
	"github.com/LightningDev1/LB-Selfbot/internal/commands/reddit"
	"github.com/LightningDev1/LB-Selfbot/internal/commands/sessions"
	"github.com/LightningDev1/LB-Selfbot/internal/commands/sharing"
	"github.com/LightningDev1/LB-Selfbot/internal/commands/slashcommands"
	"github.com/LightningDev1/LB-Selfbot/internal/commands/spotify"
	"github.com/LightningDev1/LB-Selfbot/internal/commands/text"
	"github.com/LightningDev1/LB-Selfbot/internal/commands/trolling"
	"github.com/LightningDev1/LB-Selfbot/internal/commands/utility"
	"github.com/LightningDev1/LB-Selfbot/internal/services"
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
