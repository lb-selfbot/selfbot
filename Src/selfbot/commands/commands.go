package commands

import (
	"github.com/LightningDev1/dgc"
	"github.com/LightningDev1/discordgo"
)

var (
	CmdFlagDangerous        = "D"
	CmdFlagHostingDangerous = "H"
)

type CommandRegistrar struct {
	TextCommands           TextCommands
	EncodingCommands       EncodingCommands
	CurrencyCommands       CurrencyCommands
	FunCommands            FunCommands
	TrollingCommands       TrollingCommands
	EventReactionsCommands EventReactionsCommands
	ImageCommands          ImageCommands
	MemeGenCommands        MemeGenCommands
	AICommands             AICommands
	AdministrationCommands AdministrationCommands
	UtilityCommands        UtilityCommands
	SlashCommands          SlashCommands
	DumpingCommands        DumpingCommands
	BackupCommands         BackupCommands
	SpotifyCommands        SpotifyCommands
	GameCommands           GameCommands
	ActivityCommands       ActivityCommands
	AnimationCommands      AnimationCommands
	SessionsCommands       SessionsCommands
	NSFWCommands           NSFWCommands
	RedditCommands         RedditCommands
	SharingCommands        SharingCommands
	ProtectionsCommands    ProtectionsCommands
	ConfigCommands         ConfigCommands
	HelpCommand            HelpCommand
}

func (c CommandRegistrar) RegisterCommands(session *discordgo.Session, router *dgc.Router, registerHelp bool) {
	c.TextCommands.RegisterCommands(router)
	c.EncodingCommands.RegisterCommands(router)
	c.CurrencyCommands.RegisterCommands(router)
	c.FunCommands.RegisterCommands(router)
	c.TrollingCommands.RegisterCommands(router)
	c.EventReactionsCommands.RegisterCommands(router)
	c.ImageCommands.RegisterCommands(router)
	c.MemeGenCommands.RegisterCommands(router)
	c.AICommands.RegisterCommands(router)
	c.AdministrationCommands.RegisterCommands(router)
	c.UtilityCommands.RegisterCommands(router)
	c.SlashCommands.RegisterCommands(router)
	c.DumpingCommands.RegisterCommands(router)
	// c.BackupCommands.RegisterCommands(router)
	c.SpotifyCommands.RegisterCommands(router)
	c.GameCommands.RegisterCommands(router)
	c.ActivityCommands.RegisterCommands(router, session)
	c.AnimationCommands.RegisterCommands(router)
	c.SessionsCommands.RegisterCommands(router, session)
	c.NSFWCommands.RegisterCommands(router)
	c.RedditCommands.RegisterCommands(router)
	c.SharingCommands.RegisterCommands(router)
	c.ProtectionsCommands.RegisterCommands(router)
	c.ConfigCommands.RegisterCommands(router)
	if registerHelp {
		c.HelpCommand.RegisterHelpCommand(router, session)
	}
}

var Commands = CommandRegistrar{}
