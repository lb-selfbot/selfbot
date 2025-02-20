package interfaces

import (
	"github.com/LightningDev1/discordgo"
)

// CommandContext represents the context for a command execution
type CommandContext struct {
	Session   *discordgo.Session
	Arguments []string
	Message   *discordgo.Message
	Services  interface{} // Will be *services.Services
}

// Command represents a registered bot command
type Command struct {
	Name        string
	Description string
	Usage       string
	Handler     func(*CommandContext)
}

// CommandCategory represents a group of related commands
type CommandCategory struct {
	Name        string
	Description string
	Commands    []*Command
}

// CommandsService handles command registration and execution
type CommandsService interface {
	// Command registration
	RegisterCommand(name string, description string, usage string, handler func(*CommandContext))
	RegisterCategory(name string, description string) *CommandCategory

	// Command lookup
	GetCommand(name string) *Command
	GetCategory(name string) *CommandCategory

	// Command execution
	HandleCommand(ctx *CommandContext)
}
