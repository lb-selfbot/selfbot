package game

import (
	"github.com/LightningDev1/LB-Selfbot/internal/services"
	"github.com/LightningDev1/dgc"
)

type GameCommands struct {
	services *services.Services
}

func NewGameCommands(services *services.Services) *GameCommands {
	return &GameCommands{services: services}
}

func (c *GameCommands) Register() {
	router := c.services.Discord.GetRouter()
	router.StartCategory("Game", "Game and activity commands")

	router.RegisterCmd(&dgc.Command{
		Name:        "tictactoe",
		Description: "Play Tic-tac-toe",
		Usage:       "[p]tictactoe <@user>",
		Handler:     c.tictactoe,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "rps",
		Description: "Play Rock, Paper, Scissors",
		Usage:       "[p]rps <rock/paper/scissors>",
		Handler:     c.rps,
	})

	// ... Register remaining game commands ...
}

func (c *GameCommands) tictactoe(ctx *dgc.Ctx) {
	// ... Existing implementation ...
}
