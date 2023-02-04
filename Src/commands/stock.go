package commands

import "github.com/LightningDev1/dgc"

type StockCommands struct{}

func (s StockCommands) RegisterCommands(router *dgc.Router) {
	router.StartCategory("Stock", "Stock market commands")
	router.RegisterCmd(&dgc.Command{
		Name:        "stock",
		Description: "Get stock information",
		Usage:       "[p]stock <symbol>",
		Handler:     s.stock,
	})
}

func (StockCommands) stock(ctx *dgc.Ctx) {
	// TODO
}
