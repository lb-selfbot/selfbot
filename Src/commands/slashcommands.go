package commands

import (
	"github.com/LightningDev1/dgc"
)

type SlashCommands struct{}

func (s SlashCommands) RegisterCommands(router *dgc.Router) {
	router.StartCategory("SlashCommands", "Slash command utilities")
}
