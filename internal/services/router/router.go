package router

import (
	"sync"

	"github.com/lb-selfbot/discordgo"
	"github.com/lb-selfbot/selfbot/internal/services/config"
	"github.com/lb-selfbot/selfbot/internal/services/eventbroker"
)

// Service is the interface for the router service.
type Service interface {
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

	// GetNumUsedCommands returns the number of used commands.
	GetNumUsedCommands() int

	// GetUsedCommands returns the used commands.
	GetUsedCommands() []*Ctx

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
type ExecutionHandler func(ctx *Ctx) error

// Middleware represents a command middleware.
type Middleware func(following ExecutionHandler) ExecutionHandler

type routerService struct {
	eventBroker eventbroker.Service
	configSvc   config.Service

	categories []*Category
	commands   []*Command
	middleware []Middleware

	usedCmdsMu    sync.Mutex
	usedCmds      []*Ctx
	usedCmdsCount int
}

// NewService creates a new router service.
func NewService(eventBroker eventbroker.Service, configSvc config.Service) *routerService {
	svc := &routerService{eventBroker: eventBroker, configSvc: configSvc}
	svc.subscribe()
	return svc
}

func (r *routerService) subscribe() {
	r.eventBroker.Subscribe(eventbroker.EventSessionReady, func(payload any) {
		session, ok := payload.(*discordgo.Session)
		if !ok {
			return
		}

		session.AddHandler(r.handler)
	})
}

func (r *routerService) getCommandPrefix() string {
	return r.configSvc.GetConfig().Selfbot.CommandPrefix
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

func (r *routerService) GetNumUsedCommands() int {
	return r.usedCmdsCount
}

func (r *routerService) GetUsedCommands() []*Ctx {
	return r.usedCmds
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
	ctx := &Ctx{
		Router:    r,
		configSvc: r.configSvc,
		Command:   cmd,
		Session:   session,
		Message: &discordgo.Message{
			ChannelID: channel.ID,
			Content:   r.getCommandPrefix() + cmd.Name + " " + args,
		},
	}

	ctx.Arguments = ParseArguments(ctx, args)

	r.triggerCommand(ctx)
}

func (r *routerService) handleCmdError(ctx *Ctx, err error) {
	ctx.RespondText("An error occurred while executing the %s command: %s", ctx.Command.Name, err.Error())
}
