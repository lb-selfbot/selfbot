package router

import (
	"fmt"

	"github.com/lb-selfbot/discordgo"
	"github.com/lb-selfbot/selfbot/internal/services/config"
)

// Ctx represents the context for command execution.
type Ctx struct {
	Router    Service
	configSvc config.Service

	Command   *Command
	Arguments *Arguments

	Session *discordgo.Session
	Message *discordgo.Message
}

func (ctx *Ctx) RespondText(format string, e ...any) error {
	_, err := ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, fmt.Sprintf(format, e...))
	return err
}
