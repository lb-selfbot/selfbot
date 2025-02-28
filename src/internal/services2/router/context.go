package router

import (
	"fmt"

	"github.com/LightningDev1/discordgo"
)

// Ctx represents the context for command execution.
type Ctx struct {
	Router    Router
	Command   *Command
	Arguments *Arguments

	Session *discordgo.Session
	Message *discordgo.Message
}

func (c *Ctx) RespondText(format string, e ...any) error {
	_, err := c.Session.ChannelMessageSend(c.Message.ChannelID, fmt.Sprintf(format, e...))
	return err
}
