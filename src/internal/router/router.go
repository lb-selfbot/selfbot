package router

import (
	"regexp"
	"strings"

	"github.com/LightningDev1/discordgo"
)

// regexSplitting represents the regex to split the arguments at.
var regexSplitting = regexp.MustCompile(`\s+`)

// Middleware represents a command middleware.
type Middleware func(following ExecutionHandler) ExecutionHandler

// RouterConfig represents the configuration for the router.
type RouterConfig struct {
	PrefixFunc        func() []string
	IgnorePrefixCase  bool
	IgnoreCmdCase     bool
	IsUserAllowedFunc func(ctx *Ctx) bool
}

// Router represents a command router.
type Router interface {
	GetConfig() *RouterConfig

	RegisterHandler(session *discordgo.Session)

	RegisterCategory(category *Category, commands []*Command)
	GetCategories() []*Category
	GetCategory(name string) *Category

	RegisterCmd(command *Command)
	GetCmds() []*Command
	GetCmd(name string) *Command

	RegisterMiddleware(middleware Middleware)
}

func NewRouter(cfg *RouterConfig) Router {
	return &router{cfg: cfg}
}

type router struct {
	cfg *RouterConfig

	categories  []*Category
	commands    []*Command
	middlewares []Middleware
}

func (r *router) GetConfig() *RouterConfig {
	return r.cfg
}

func (r *router) RegisterHandler(session *discordgo.Session) {
	session.AddHandler(r.handler)
}

func (r *router) RegisterCategory(category *Category, commands []*Command) {
	for _, command := range commands {
		command.Category = category
	}

	r.categories = append(r.categories, category)
	r.commands = append(r.commands, commands...)
}

func (r *router) GetCategories() []*Category {
	return r.categories
}

func (r *router) GetCategory(name string) *Category {
	for _, category := range r.categories {
		if category.Name == name {
			return category
		}
	}

	return nil
}

func (r *router) RegisterCmd(command *Command) {
	r.commands = append(r.commands, command)
}

func (r *router) GetCmds() []*Command {
	return r.commands
}

func (r *router) GetCmd(name string) *Command {
	for _, command := range r.commands {
		if command.Name == name {
			return command
		}
	}

	return nil
}

func (r *router) RegisterMiddleware(middleware Middleware) {
	r.middlewares = append(r.middlewares, middleware)
}

// handler provides the discordgo handler for the router.
func (r *router) handler(session *discordgo.Session, message *discordgo.MessageCreate) {
	// Ignore messages from bots
	if message.Author == nil || message.Author.Bot {
		return
	}

	prefixes := r.cfg.PrefixFunc()
	hasPrefix, content := stringHasPrefix(message.Content, prefixes, r.cfg.IgnorePrefixCase)
	if !hasPrefix {
		return
	}

	content = strings.TrimSpace(content)
	if content == "" {
		return
	}

	parts := regexSplitting.Split(content, -1)
	if len(parts) == 0 || parts[0] == "" {
		return
	}

	for _, cmd := range r.commands {
		if !stringArrayContains(cmd.GetIdentifiers(), parts[0], r.cfg.IgnoreCmdCase) {
			continue
		}

		content = strings.TrimSpace(strings.Join(parts[1:], " "))

		ctx := &Ctx{
			Router:  r,
			Command: cmd,
			Session: session,
			Message: message.Message,
		}

		ctx.Arguments = ParseArguments(ctx, content)

		if !r.cfg.IsUserAllowedFunc(ctx) {
			return
		}

		r.triggerCommand(ctx)
		return
	}
}

func (r *router) triggerCommand(ctx *Ctx) {
	if r.triggerSubCommand(ctx) {
		return
	}

	nextHandler := ctx.Command.Handler
	for _, middleware := range r.middlewares {
		nextHandler = middleware(nextHandler)
	}

	nextHandler(ctx)
}

func (r *router) triggerSubCommand(ctx *Ctx) bool {
	if ctx.Arguments.Amount() == 0 {
		return false
	}

	arg := ctx.Arguments.Get(0).Raw()
	subCmd := ctx.Command.GetSubCmd(arg, r.cfg.IgnoreCmdCase)
	if subCmd == nil {
		return false
	}

	subCtx := &Ctx{
		Router:  r,
		Command: subCmd,
		Session: ctx.Session,
		Message: ctx.Message,
	}

	subCtx.Arguments = ParseArguments(subCtx, "")
	if ctx.Arguments.Amount() > 1 {
		rawSplit := regexSplitting.Split(ctx.Arguments.Raw(), -1)[1:]
		subCtx.Arguments = ParseArguments(subCtx, strings.Join(rawSplit, " "))
	}

	r.triggerCommand(subCtx)

	return true
}
