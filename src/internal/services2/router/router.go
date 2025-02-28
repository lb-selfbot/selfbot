package router

import (
	"github.com/LightningDev1/LB-Selfbot/internal/services2/config"
	"github.com/LightningDev1/discordgo"
)

// RouterService is the interface for the router service.
type RouterService interface {
	// GetHandler returns the message create handler for the command router.
	GetHandler() func(*discordgo.Session, *discordgo.MessageCreate)

	// GetCategories returns all registered categories.
	GetCategories() []*Category
	// GetCategory returns the category with the given name.
	GetCategory(name string) *Category
	// RegisterCategory registers a category with the given commands.
	RegisterCategory(category *Category, commands []*Command)

	// GetCommands returns all registered commands.
	GetCommands() []*Command
	// GetCommand returns the command with the given name.
	GetCommand(name string) *Command

	// RegisterMiddleware registers a middleware.
	RegisterMiddleware(middleware Middleware)

	// InvokeRemoteCmd invokes a command in a specific channel.
	InvokeRemoteCmd(
		session *discordgo.Session,
		channel *discordgo.Channel,
		author *discordgo.User,
		cmd *Command,
		args string,
	)
}

// ExecutionHandler represents a handler for executing commands.
type ExecutionHandler func(*Ctx)

// Middleware represents a command middleware.
type Middleware func(following ExecutionHandler) ExecutionHandler

// Implementation check
var _ RouterService = (*routerService)(nil)

type routerService struct {
	cfgService config.ConfigService

	categories []*Category
	commands   []*Command
	middleware []Middleware
}

// NewRouterService creates a new router service.
func NewRouterService(cfg config.ConfigService) *routerService {
	return &routerService{
		cfgService: cfg,
	}
}

func (r *routerService) getCommandPrefix() string {
	return r.cfgService.GetConfig().Selfbot.CommandPrefix
}

func (r *routerService) GetHandler() func(*discordgo.Session, *discordgo.MessageCreate) {
	return r.handler
}

func (r *routerService) RegisterCategory(category *Category, commands []*Command) {
	for _, command := range commands {
		command.Category = category
	}

	r.categories = append(r.categories, category)
	r.commands = append(r.commands, commands...)
}

func (r *routerService) GetCategories() []*Category {
	return r.categories
}

func (r *routerService) GetCategory(name string) *Category {
	for _, category := range r.categories {
		if category.Name == name {
			return category
		}
	}

	return nil
}

func (r *routerService) RegisterCmd(command *Command) {
	r.commands = append(r.commands, command)
}

func (r *routerService) GetCommands() []*Command {
	return r.commands
}

func (r *routerService) GetCommand(name string) *Command {
	for _, command := range r.commands {
		if command.Name == name {
			return command
		}
	}

	return nil
}

func (r *routerService) RegisterMiddleware(middleware Middleware) {
	r.middleware = append(r.middleware, middleware)
}

func (r *routerService) InvokeRemoteCmd(
	session *discordgo.Session,
	channel *discordgo.Channel,
	author *discordgo.User,
	cmd *Command,
	args string,
) {
	// todo
}
