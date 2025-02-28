package router

import (
	"strings"

	"github.com/LightningDev1/LB-Selfbot/internal/utils2"
	"github.com/LightningDev1/discordgo"
)

func (r *routerService) handler(session *discordgo.Session, message *discordgo.MessageCreate) {
	// Ignore messages from bots
	if message.Author == nil || message.Author.Bot {
		return
	}

	prefixes := r.getCommandPrefix()
	hasPrefix, content := stringHasPrefix(message.Content, prefixes, true)
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

	handler := ctx.Command.Handler
	if handler == nil {
		subCmdsReadable := utils2.Map(ctx.Command.SubCommands, func(subCmd *Command) string {
			return subCmd.Name
		})

		ctx.RespondText("You must provide a subcommand for [p]%s! The available subcommands are: %s", ctx.Command.Name, strings.Join(subCmdsReadable, ", "))
		return
	}

	for _, middleware := range r.middlewares {
		handler = middleware(handler)
	}

	handler(ctx)
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
