package commands

import "github.com/LightningDev1/dgc"

type Commands struct {
	TextCommands      TextCommands
	EncodingCommands  EncodingCommands
	CodeblockCommands CodeblockCommands
	CryptoCommands    CryptoCommands
	FunCommands       FunCommands
	TrollingCommands  TrollingCommands
	ImageCommands     ImageCommands
	NetworkCommands   NetworkCommands
	SpotifyCommands   SpotifyCommands
	GameCommands      GameCommands
	ConfigCommands    ConfigCommands
	RPCCommands       RPCCommands
	NSFWCommands      NSFWCommands
	RedditCommands    RedditCommands
}

func (c Commands) RegisterCommands(router *dgc.Router) {
	c.TextCommands.RegisterCommands(router)
	c.EncodingCommands.RegisterCommands(router)
	c.CodeblockCommands.RegisterCommands(router)
	c.CryptoCommands.RegisterCommands(router)
	c.FunCommands.RegisterCommands(router)
	c.TrollingCommands.RegisterCommands(router)
	c.ImageCommands.RegisterCommands(router)
	c.NetworkCommands.RegisterCommands(router)
	c.SpotifyCommands.RegisterCommands(router)
	c.GameCommands.RegisterCommands(router)
	c.ConfigCommands.RegisterCommands(router)
	c.RPCCommands.RegisterCommands(router)
	c.NSFWCommands.RegisterCommands(router)
	c.RedditCommands.RegisterCommands(router)
}
