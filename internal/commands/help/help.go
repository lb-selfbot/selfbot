package help

import (
	"fmt"
	"math"
	"strings"

	"github.com/agnivade/levenshtein"
	"github.com/lb-selfbot/selfbot/internal/constants"
	"github.com/lb-selfbot/selfbot/internal/embed"

	"github.com/lb-selfbot/selfbot/internal/fn"
	"github.com/lb-selfbot/selfbot/internal/services"
	"github.com/lb-selfbot/selfbot/internal/services/router"
	"github.com/lb-selfbot/selfbot/internal/utils"
)

type HelpCommands struct{ services *services.Services }

func NewHelpCommands(s *services.Services) *HelpCommands {
	return &HelpCommands{services: s}
}

func (c *HelpCommands) Register() {
	cmds := []*router.Command{
		{
			Name:    "help",
			Handler: c.generalHelpCommand,
		},
		{
			Name:        "search",
			Description: "Search for a command",
			Arguments:   "<query>",
			Handler:     c.search,
		},
		{
			Name:        "repeat",
			Description: "Repeat the last used command",
			Handler:     c.repeat,
		},
	}

	for _, category := range c.services.Router.GetCategories() {
		name := strings.ReplaceAll(strings.ToLower(category.Name), " ", "")

		cmds = append(cmds, &router.Command{
			Name:        name,
			Description: fmt.Sprintf("Get help for the %s category", category.Name),
			Handler:     c.generateCategoryHelpFunc(category, name),
		})
	}

	c.services.Router.RegisterCategory(&router.Category{
		Name:        "Help",
		Description: "Get help for commands",
	}, cmds)
}

type HelpMessage struct {
	Page     int
	Category *router.Category
	Command  *router.Command
}

func (c *HelpCommands) generateCategoryHelpFunc(category *router.Category, name string) router.ExecutionHandler {
	return func(ctx *router.Ctx) error {
		page, _ := ctx.Arguments.Get(0).Int()

		c.renderCategoryHelp(ctx, category, name, page).Send()
		return nil
	}
}

// generalHelpCommand handles the general help command
func (c *HelpCommands) generalHelpCommand(ctx *router.Ctx) error {
	// Check if the user provided an argument
	if ctx.Arguments.Amount() > 0 {
		c.specificHelpCommand(ctx)
		return nil
	}

	// Send the general help embed
	c.renderDefaultGeneralHelpEmbed(ctx).Send()
	return nil
}

func (c *HelpCommands) getNameFromArgument(arg string) string {
	arg = strings.ReplaceAll(arg, "eventreactions", "EventReactions")
	arg = strings.ReplaceAll(arg, "eventreaction", "EventReactions")

	arg = strings.ReplaceAll(arg, "nsfw", "NSFW")

	arg = strings.ReplaceAll(arg, "customactivity", "CustomActivity")

	arg = strings.ReplaceAll(arg, "slashcommands", "SlashCommands")
	arg = strings.ReplaceAll(arg, "slashcommand", "SlashCommands")

	if arg == "ai" {
		arg = "AI"
	}

	if arg == "admin" {
		arg = "Admininstration"
	}

	return arg
}

// specificHelpCommand handles the specific help command
func (c *HelpCommands) specificHelpCommand(ctx *router.Ctx) error {
	oldName := strings.ToLower(ctx.Arguments.Get(0).Raw())
	name := c.getNameFromArgument(oldName)
	page, _ := ctx.Arguments.Get(1).Int()
	cfg := c.services.Config.GetConfig()
	prefix := cfg.Selfbot.CommandPrefix

	if category := ctx.Router.GetCategory(utils.String.Capitalize(name)); category != nil {
		c.renderCategoryHelp(ctx, category, oldName, page).Send()
		return nil
	}

	if command := ctx.Router.GetCommand(name); command != nil {
		c.renderCommandHelp(ctx, command).Send()
		return nil
	}

	// get close commands and categories
	var closeCommands []string
	for _, command := range ctx.Router.GetCommands() {
		if len(closeCommands) >= 5 {
			break
		}
		distance := levenshtein.ComputeDistance(command.Name, name)
		if distance < 3 || strings.Contains(command.Name, name) {
			closeCommands = append(closeCommands, prefix+command.Name+" "+command.Arguments+" - "+command.Description)
		}
	}

	var closeCategories []string
	for _, category := range ctx.Router.GetCategories() {
		if len(closeCategories) >= 5 {
			break
		}
		distance := levenshtein.ComputeDistance(category.Name, utils.String.Capitalize(name))
		if distance < 3 || strings.Contains(category.Name, utils.String.Capitalize(name)) {
			closeCategories = append(closeCategories, category.Name+" - "+category.Description)
		}
	}

	var suggestions string

	if len(closeCommands) > 0 {
		suggestions += "\n\nDid you mean one of these commands?\n\n" + strings.Join(closeCommands, "\n")
	}
	if len(closeCategories) > 0 {
		suggestions += "\n\nDid you mean one of these categories?\n\n" + strings.Join(closeCategories, "\n")
	}

	c.services.Embed.Respond(ctx, "`%s` is not a command or category!\n%s", name, suggestions)
	return nil
}

func (c *HelpCommands) renderCommandHelp(ctx *router.Ctx, command *router.Command) *embed.Embed {
	cfg := c.services.Config.GetConfig()
	helpText := fmt.Sprintf("%s%s %s", cfg.Selfbot.CommandPrefix, command.Name, command.Arguments)
	if len(command.SubCommands) > 0 {
		helpText = ""

		for _, subCommand := range command.SubCommands {
			helpText += fmt.Sprintf("%s%s %s %s - %s\n", cfg.Selfbot.CommandPrefix, command.Name, subCommand.Name, subCommand.Arguments, subCommand.Description)
		}
	}

	helpText = strings.ReplaceAll(helpText, "  ", " ")
	return c.services.Embed.New(ctx).
		SetDescription("Command Help: %s", command.Name).
		AddField("Command Name", command.Name).
		AddField("Usage", "<> is required, [] is optional\n\n"+helpText).
		AddField("Description", command.Description, command.Description != "").
		AddField("Category", command.Category.Name)
}

func (c *HelpCommands) renderCategoryHelp(ctx *router.Ctx, category *router.Category, rawArgument string, page int) *embed.Embed {
	helpPages := map[int]string{}
	currentPage := 0

	cmds := fn.Filter(ctx.Router.GetCommands(), func(cmd *router.Command) bool {
		return cmd.Category == category
	})

	cfg := c.services.Config.GetConfig()

	for _, command := range cmds {
		helpText := fmt.Sprintf("%s%s %s - %s\n", cfg.Selfbot.CommandPrefix, command.Name, command.Arguments, command.Description)
		if len(command.SubCommands) > 0 {
			helpText = fmt.Sprintf("%s%s %s - Use `!help %s` for more info\n", cfg.Selfbot.CommandPrefix, command.Name, command.Arguments, command.Name)
			// helpText = ""

			// for _, subCommand := range command.SubCommands {
			// 	helpText += fmt.Sprintf("%s%s %s %s - %s\n", cfg.Selfbot.CommandPrefix, command.Name, subCommand.Name, subCommand.Arguments, subCommand.Description)
			// }
		}

		helpText = strings.ReplaceAll(helpText, "  ", " ")

		if len(helpPages[currentPage])+len(helpText) > 1800 {
			currentPage++
		}
		helpPages[currentPage] += helpText
	}

	page = int(math.Max(0, float64(page-1)))

	emb := c.services.Embed.New(ctx).
		SetDescription("<> is required, [] is optional\n\nAmount of %s commands: %d", category.Name, len(cmds)).
		SetSubtext("Page: %d/%d", page+1, len(helpPages))

	if page < len(helpPages) {
		if cfg.Selfbot.CommandMode == "WebhookEmbed" {
			emb.Description += "\n\n**Commands**\n" + helpPages[page]
		} else {
			emb.AddField("Commands", helpPages[page], true, true)
		}
	} else {
		return c.services.Embed.New(ctx).SetDescription("That page does not exist! %s has %d pages.", category.Name, len(helpPages))
	}
	if len(helpPages) > 1 {
		if cfg.Selfbot.CommandMode == "WebhookEmbed" {
			emb.SetSubtext("Page: %d/%d\n\nUse %shelp %s <page> to go to the next page", page+1, len(helpPages), cfg.Selfbot.CommandPrefix, rawArgument)
		} else {
			emb.SetSubtext("Page: %d/%d\n\nUse \x1b[34;1m%shelp %s <page>\x1b[0m to go to the next page", page+1, len(helpPages), cfg.Selfbot.CommandPrefix, rawArgument)
		}
	}
	return emb
}

// renderDefaultGeneralHelpEmbed renders the general help embed on the given page
func (c *HelpCommands) renderDefaultGeneralHelpEmbed(ctx *router.Ctx) *embed.Embed {
	cfg := c.services.Config.GetConfig()
	categoriesString := ""

	for _, category := range ctx.Router.GetCategories() {
		// Get the commands in the category, will be empty otherwise
		category = ctx.Router.GetCategory(category.Name)
		categoriesString += fmt.Sprintf("%s - %s\n", category.Name+strings.Repeat(" ", 14-len(category.Name)), category.Description)
	}

	prefix := cfg.Selfbot.CommandPrefix

	emb := c.services.Embed.New(ctx).
		SetDescription("LightningBot %s Help\nAmount of Commands: %d, <> is required, [] is optional\n", constants.VERSION, len(ctx.Router.GetCommands()))

	if cfg.Selfbot.CommandMode == "WebhookEmbed" {
		emb.Description += "\n\n**Categories**\n" + categoriesString
		emb.SetSubtext("Use the %[1]shelp <category name> or %[1]shelp <command name> to get more info.\nUse %[1]ssearch <query> to search commands.\nUse %[1]srepeat to repeat the last command.", prefix)
	} else {
		emb.AddField("Categories", categoriesString, true, true)
		emb.SetSubtext("Use \x1b[34;1m%[1]shelp <category name>\x1b[0m or \x1b[34;1m%[1]shelp <command name>\x1b[0m to get more info.\nUse \x1b[34;1m%[1]ssearch <query>\x1b[0m to search commands.\nUse \x1b[34;1m%[1]srepeat\x1b[0m to repeat the last command.", prefix)
	}

	return emb
}

func (c *HelpCommands) search(ctx *router.Ctx) error {
	query := ctx.Arguments.Get(0).Raw()

	cfg := c.services.Config.GetConfig()

	var closeCommands []string
	for _, command := range ctx.Router.GetCommands() {
		if len(closeCommands) >= 10 {
			break
		}
		distance := levenshtein.ComputeDistance(command.Name, query)
		if distance < 2 || strings.Contains(command.Name, query) {
			closeCommands = append(closeCommands, cfg.Selfbot.CommandPrefix+command.Name+" "+command.Arguments+" - "+command.Description)
		}
	}

	if len(closeCommands) == 0 {
		c.services.Embed.New(ctx).
			SetDescription("No commands found for query `%s`", query).
			Send()
		return nil
	}

	c.services.Embed.New(ctx).
		SetDescription("Found %d command%s for query `%s`", len(closeCommands), utils.String.Pluralize(len(closeCommands)), query).
		AddField("Commands", strings.Join(closeCommands, "\n"), true, true).
		Send()
	return nil
}

func (c *HelpCommands) repeat(ctx *router.Ctx) error {
	usedCommands := c.services.EventReactions.GetUsedCommands()

	if len(usedCommands) < 1 {
		c.services.Embed.New(ctx).
			SetDescription("No commands have been used yet").
			Send()
		return nil
	}

	lastCommand := usedCommands[len(usedCommands)-1]

	channel, err := ctx.Session.State.Channel(ctx.Message.ChannelID)
	if err != nil {
		c.services.Embed.New(ctx).
			SetDescription("Failed to get channel: %s", err.Error()).
			Send()
		return nil
	}

	ctx.Router.InvokeRemoteCmd(ctx.Session, channel, ctx.Message.Author, lastCommand.Command, lastCommand.Arguments.Raw())
	return nil
}
