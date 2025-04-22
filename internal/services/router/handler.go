package router

import (
	"fmt"
	"regexp"
	"slices"
	"strings"

	"github.com/lb-selfbot/discordgo"
	"github.com/lb-selfbot/selfbot/internal/constants"
	"github.com/lb-selfbot/selfbot/internal/fn"
	"github.com/lb-selfbot/selfbot/internal/services/eventbroker"
	"github.com/lb-selfbot/selfbot/internal/services/ui/featurelogs"
)

// regexSplitting represents the regex to split the arguments at.
var regexSplitting = regexp.MustCompile(`\s+`)

func (r *routerService) isUserAllowed(ctx *Ctx) bool {
	if constants.HOSTING && slices.Contains(ctx.Command.Flags, CmdFlagHostingDangerous) {
		// d.uiSvc.LogDangerousHostingCmdBlocked(ctx.Command.Name)
		return false
	}

	if ctx.Message.Author.ID == ctx.Session.State.User.ID {
		return true
	}

	cfg := r.configSvc.GetConfig()

	if !slices.Contains(cfg.Sharing.SharedUsers, ctx.Message.Author.ID) {
		return false
	}

	if cfg.Sharing.BlockDangerousCommands && slices.Contains(ctx.Command.Flags, CmdFlagDangerous) {
		// d.uiSvc.LogDangerousCmdBlocked(ctx.Message.Author, ctx.Command.Name)
		return false
	}

	if cfg.Sharing.ShareAll || slices.Contains([]string{"help", "search", "repeat"}, ctx.Command.Name) {
		return true
	}

	categoryName := ""
	if ctx.Command.Category != nil {
		categoryName = ctx.Command.Category.Name
	}

	for _, category := range cfg.Sharing.SharedCategories {
		if strings.EqualFold(category, categoryName) {
			return true
		}
	}

	return slices.Contains(cfg.Sharing.SharedCommands, ctx.Command.Name)
}

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
		if !stringArrayContains(cmd.GetIdentifiers(), parts[0], true) {
			continue
		}

		content = strings.TrimSpace(strings.Join(parts[1:], " "))

		ctx := &Ctx{
			Router:    r,
			configSvc: r.configSvc,
			Command:   cmd,
			Session:   session,
			Message:   message.Message,
		}

		ctx.Arguments = ParseArguments(ctx, content)

		if !r.isUserAllowed(ctx) {
			return
		}

		r.eventBroker.Publish(eventbroker.EventLogFeature, &featurelogs.CmdUsed{Name: cmd.Name})

		r.triggerCommand(ctx)

		r.usedCmdsMu.Lock()
		r.usedCmdsCount++
		if cmd.Name != "repeat" {
			r.usedCmds = append(r.usedCmds, ctx)
		}
		r.usedCmdsMu.Unlock()

		return
	}
}

func (r *routerService) triggerCommand(ctx *Ctx) {
	if r.triggerSubCommand(ctx) {
		return
	}

	handler := ctx.Command.Handler
	if handler == nil {
		subCmdsReadable := fn.Map(ctx.Command.SubCommands, func(subCmd *Command) string {
			return subCmd.Name
		})
		r.handleCmdError(ctx, fmt.Errorf("You must provide a subcommand! The available subcommands are: %s", strings.Join(subCmdsReadable, ", ")))
		return
	}

	for _, middleware := range r.middleware {
		handler = middleware(handler)
	}

	err := handler(ctx)
	if err != nil {
		r.handleCmdError(ctx, err)
		return
	}
}

func (r *routerService) triggerSubCommand(ctx *Ctx) bool {
	if ctx.Arguments.Amount() == 0 {
		return false
	}

	arg := ctx.Arguments.Get(0).Raw()
	subCmd := ctx.Command.GetSubCmd(arg, true)
	if subCmd == nil {
		return false
	}

	subCtx := &Ctx{
		Router:    r,
		configSvc: r.configSvc,
		Command:   subCmd,
		Session:   ctx.Session,
		Message:   ctx.Message,
	}

	subCtx.Arguments = ParseArguments(subCtx, "")
	if ctx.Arguments.Amount() > 1 {
		rawSplit := regexSplitting.Split(ctx.Arguments.Raw(), -1)[1:]
		subCtx.Arguments = ParseArguments(subCtx, strings.Join(rawSplit, " "))
	}

	r.triggerCommand(subCtx)

	return true
}
