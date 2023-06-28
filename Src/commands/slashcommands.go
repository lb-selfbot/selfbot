package commands

import (
	"github.com/LightningDev1/LB-Selfbot/utils"
	"github.com/LightningDev1/dgc"
)

type SlashCommands struct{}

func (s SlashCommands) RegisterCommands(router *dgc.Router) {
	router.StartCategory("SlashCommands", "Slash command utilities")

	router.RegisterCmd(&dgc.Command{
		Name:        "executeslash",
		Description: "Run a slash command",
		Usage:       "[p]executeslash <#channel> <command> [options...]",
		Handler:     s.executeslash,
	})
}

func (s SlashCommands) executeslash(ctx *dgc.Ctx) {
	channel, ok := utils.Discord.TryGetChannelFromMention(ctx, ctx.Arguments.Get(0))
	if !ok {
		ctx.RespondText("Invalid channel")
		return
	}

	command := utils.Arguments["<command>]"]

	_ = channel
	_ = command
}
