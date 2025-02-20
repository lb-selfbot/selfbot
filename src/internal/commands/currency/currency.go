package currency

import (
	"github.com/LightningDev1/LB-Selfbot/internal/services"
	"github.com/LightningDev1/dgc"
)

type CurrencyCommands struct {
	services *services.Services
}

func NewCurrencyCommands(services *services.Services) *CurrencyCommands {
	return &CurrencyCommands{services: services}
}

func (c *CurrencyCommands) Register() {
	router := c.services.Discord.GetRouter()
	router.StartCategory("Currency", "Currency and crypto commands")

	router.RegisterCmd(&dgc.Command{
		Name:        "convert",
		Description: "Convert between currencies",
		Usage:       "[p]convert <amount> <from> <to>",
		Handler:     c.convert,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "price",
		Description: "Get current price of a cryptocurrency",
		Usage:       "[p]price <symbol>",
		Handler:     c.price,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "crypto",
		Description: "Get detailed crypto information",
		Usage:       "[p]crypto <symbol>",
		Handler:     c.crypto,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "cryptomarket",
		Description: "Get crypto market overview",
		Usage:       "[p]cryptomarket",
		Handler:     c.cryptomarket,
	})
}

func (c *CurrencyCommands) convert(ctx *dgc.Ctx) {
	// ... Existing implementation ...
}

// ... Keep all other existing command implementations ...
