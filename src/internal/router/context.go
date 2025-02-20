package router

import "github.com/LightningDev1/discordgo"

// Ctx represents the context for command execution.
type Ctx struct {
	Router    Router
	Command   *Command
	Arguments *Arguments

	Session *discordgo.Session
	Message *discordgo.Message
}

// ExecutionHandler represents a handler for executing commands.
type ExecutionHandler func(*Ctx)
