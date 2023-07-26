package commands

import (
	"github.com/LightningDev1/dgc"
	"github.com/LightningDev1/discordgo"
)

type CommandRegistrar struct {
	TextCommands           TextCommands
	EncodingCommands       EncodingCommands
	CodeblockCommands      CodeblockCommands
	CryptoCommands         CryptoCommands
	FunCommands            FunCommands
	TrollingCommands       TrollingCommands
	EventReactionsCommands EventReactionsCommands
	ImageCommands          ImageCommands
	MemeGenCommands        MemeGenCommands
	NetworkCommands        NetworkCommands
	AICommands             AICommands
	StockCommands          StockCommands
	AdministrationCommands AdministrationCommands
	UtilityCommands        UtilityCommands
	SlashCommands          SlashCommands
	ConversionCommands     ConversionCommands
	DumpingCommands        DumpingCommands
	BackupCommands         BackupCommands
	SpotifyCommands        SpotifyCommands
	GameCommands           GameCommands
	ActivityCommands       ActivityCommands
	AnimationCommands      AnimationCommands
	CustomActivityCommands CustomActivityCommands
	SessionsCommands       SessionsCommands
	NSFWCommands           NSFWCommands
	RedditCommands         RedditCommands
	SharingCommands        SharingCommands
	ProtectionsCommands    ProtectionsCommands
	ConfigCommands         ConfigCommands
	HelpCommand            HelpCommand
}

func (c CommandRegistrar) RegisterCommands(session *discordgo.Session, router *dgc.Router) {
	c.TextCommands.RegisterCommands(router)
	c.EncodingCommands.RegisterCommands(router)
	c.CodeblockCommands.RegisterCommands(router)
	c.CryptoCommands.RegisterCommands(router)
	c.FunCommands.RegisterCommands(router)
	c.TrollingCommands.RegisterCommands(router)
	c.EventReactionsCommands.RegisterCommands(router)
	c.ImageCommands.RegisterCommands(router)
	c.MemeGenCommands.RegisterCommands(router)
	c.NetworkCommands.RegisterCommands(router)
	c.AICommands.RegisterCommands(router)
	c.StockCommands.RegisterCommands(router)
	c.AdministrationCommands.RegisterCommands(router)
	c.UtilityCommands.RegisterCommands(router)
	c.SlashCommands.RegisterCommands(router)
	c.ConversionCommands.RegisterCommands(router)
	c.DumpingCommands.RegisterCommands(router)
	c.BackupCommands.RegisterCommands(router)
	c.SpotifyCommands.RegisterCommands(router)
	c.GameCommands.RegisterCommands(router)
	c.ActivityCommands.RegisterCommands(router, session)
	c.AnimationCommands.RegisterCommands(router)
	c.CustomActivityCommands.RegisterCommands(router)
	c.SessionsCommands.RegisterCommands(router, session)
	c.NSFWCommands.RegisterCommands(router)
	c.RedditCommands.RegisterCommands(router)
	c.SharingCommands.RegisterCommands(router)
	c.ProtectionsCommands.RegisterCommands(router)
	c.ConfigCommands.RegisterCommands(router)
	c.HelpCommand.RegisterHelpCommand(router, session)
}

var Commands = CommandRegistrar{}
