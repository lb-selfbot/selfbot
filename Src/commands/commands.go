package commands

import (
	"github.com/LightningDev1/dgc"
	"github.com/LightningDev1/discordgo"
)

type Commands struct {
	TextCommands       TextCommands
	EncodingCommands   EncodingCommands
	CodeblockCommands  CodeblockCommands
	CryptoCommands     CryptoCommands
	FunCommands        FunCommands
	TrollingCommands   TrollingCommands
	ImageCommands      ImageCommands
	NetworkCommands    NetworkCommands
	ModerationCommands ModerationCommands
	UtilityCommands    UtilityCommands
	SpotifyCommands    SpotifyCommands
	GameCommands       GameCommands
	ActivityCommands   ActivityCommands
	AnimationCommands  AnimationCommands
	ConfigCommands     ConfigCommands
	RPCCommands        RPCCommands
	SessionsCommands   SessionsCommands
	NSFWCommands       NSFWCommands
	RedditCommands     RedditCommands
}

func (c Commands) RegisterCommands(router *dgc.Router, session *discordgo.Session) {
	c.TextCommands.RegisterCommands(router)
	c.EncodingCommands.RegisterCommands(router)
	c.CodeblockCommands.RegisterCommands(router)
	c.CryptoCommands.RegisterCommands(router)
	c.FunCommands.RegisterCommands(router)
	c.TrollingCommands.RegisterCommands(router)
	c.ImageCommands.RegisterCommands(router)
	c.NetworkCommands.RegisterCommands(router)
	c.ModerationCommands.RegisterCommands(router)
	c.UtilityCommands.RegisterCommands(router)
	c.SpotifyCommands.RegisterCommands(router)
	c.GameCommands.RegisterCommands(router)
	c.ActivityCommands.RegisterCommands(router)
	c.AnimationCommands.RegisterCommands(router)
	c.ConfigCommands.RegisterCommands(router)
	c.RPCCommands.RegisterCommands(router)
	c.SessionsCommands.RegisterCommands(router, session)
	c.NSFWCommands.RegisterCommands(router)
	c.RedditCommands.RegisterCommands(router)
}
