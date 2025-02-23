package help

import (
	"fmt"
	"math"
	"strings"

	"github.com/LightningDev1/LB-Selfbot/internal/config"
	"github.com/LightningDev1/LB-Selfbot/internal/constants"
	"github.com/LightningDev1/LB-Selfbot/internal/embed"
	"github.com/LightningDev1/LB-Selfbot/internal/router"
	"github.com/LightningDev1/LB-Selfbot/internal/services"
	"github.com/LightningDev1/LB-Selfbot/internal/utils"
	"github.com/LightningDev1/LB-Selfbot/internal/utils2"
)

type HelpCommands struct{ services *services.Services }

func NewHelpCommands(s *services.Services) *HelpCommands {
	return &HelpCommands{services: s}
}

func (c *HelpCommands) Register() {
	cmdRouter := c.services.Discord.GetRouter()

	cmds := []router.Command{
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

	for _, category := range cmdRouter.GetCategories() {
		name := strings.ReplaceAll(strings.ToLower(category.Name), " ", "")

		cmds = append(cmds, router.Command{
			Name:        name,
			Description: fmt.Sprintf("Get help for the %s category", category.Name),
			Handler:     c.generateCategoryHelpFunc(category, name),
		})
	}

}

type HelpMessage struct {
	Page     int
	Category *router.Category
	Command  *router.Command
}

func (c *HelpCommands) generateCategoryHelpFunc(category *router.Category, name string) router.ExecutionHandler {
	return func(ctx *router.Ctx) {
		page, _ := ctx.Arguments.Get(0).Int()

		c.renderCategoryHelp(ctx.Router, category, name, page).Send(ctx)
	}
}

// generalHelpCommand handles the general help command
func (c *HelpCommands) generalHelpCommand(ctx *router.Ctx) {
	// Check if the user provided an argument
	if ctx.Arguments.Amount() > 0 {
		c.specificHelpCommand(ctx)
		return
	}

	// Send the general help embed
	c.renderDefaultGeneralHelpEmbed(ctx.Router).Send(ctx)
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
func (c *HelpCommands) specificHelpCommand(ctx *router.Ctx) {
	oldName := strings.ToLower(ctx.Arguments.Get(0).Raw())
	name := c.getNameFromArgument(oldName)
	page, _ := ctx.Arguments.Get(1).Int()
	prefix := ctx.Router.GetConfig().PrefixFunc()[0]

	if category := ctx.Router.GetCategory(utils.String.Capitalize(name)); category != nil {
		c.renderCategoryHelp(ctx.Router, category, oldName, page).Send(ctx)
		return
	}

	if command := ctx.Router.GetCmd(name); command != nil {
		c.renderCommandHelp(ctx.Router, command).Send(ctx)
		return
	}

	// get close commands and categories
	var closeCommands []string
	for _, command := range ctx.Router.GetCmds() {
		if len(closeCommands) >= 5 {
			break
		}
		distance := utils.Misc.GetDistance(command.Name, name)
		if distance < 2 || strings.Contains(command.Name, name) {
			closeCommands = append(closeCommands, prefix+command.Name+" "+command.Arguments+" - "+command.Description)
		}
	}

	var closeCategories []string
	for _, category := range ctx.Router.GetCategories() {
		if len(closeCategories) >= 5 {
			break
		}
		distance := utils.Misc.GetDistance(category.Name, utils.String.Capitalize(name))
		if distance < 2 || strings.Contains(category.Name, utils.String.Capitalize(name)) {
			closeCategories = append(closeCategories, category.Name+" - "+category.Description)
		}
	}

	description := fmt.Sprintf("`%s` is not a command or category!", name)
	if len(closeCommands) > 0 {
		description += "\n\nDid you mean one of these commands?\n\n" + strings.Join(closeCommands, "\n")
	}
	if len(closeCategories) > 0 {
		description += "\n\nDid you mean one of these categories?\n\n" + strings.Join(closeCategories, "\n")
	}

	c.services.Embed.NewEmbed().SetDescription(description).Send(ctx)
}

func (c *HelpCommands) renderCommandHelp(cmdRouter router.Router, command *router.Command) *embed.Embed {
	return c.services.Embed.NewEmbed().
		SetDescription("Command Help: "+command.Name).
		AddField("Command Name", command.Name).
		AddField("Usage", "<> is required, [] is optional\n\n"+cmdRouter.GetConfig().PrefixFunc()[0]+command.Name+" "+command.Arguments).
		AddField("Description", command.Description).
		AddField("Category", command.Category.Name)
}

func (c *HelpCommands) renderCategoryHelp(cmdRouter router.Router, category *router.Category, rawArgument string, page int) *embed.Embed {
	helpPages := map[int]string{}
	currentPage := 0

	cmds := utils2.Filter(cmdRouter.GetCmds(), func(cmd *router.Command) bool {
		return cmd.Category == category
	})

	for _, command := range cmds {
		helpText := fmt.Sprintf("%s%s %s - %s\n", cmdRouter.GetConfig().PrefixFunc()[0], command.Name, command.Arguments, command.Description)
		if len(helpPages[currentPage])+len(helpText) > 1650 {
			currentPage++
		}
		helpPages[currentPage] += helpText
	}

	page = int(math.Max(0, float64(page-1)))

	emb := c.services.Embed.NewEmbed().
		SetDescription("<> is required, [] is optional\n\nAmount of %s commands: %d", category.Name, len(cmds)).
		SetSubtext(fmt.Sprintf("Page: %d/%d", page+1, len(helpPages)))

	cfg, _ := config.GetConfig()

	if page < len(helpPages) {
		if cfg.Selfbot.CommandMode == "WebhookEmbed" {
			emb.Description += "\n\n**Commands**\n" + helpPages[page]
		} else {
			emb.AddField("Commands", helpPages[page], true, true)
		}
	} else {
		return c.services.Embed.NewEmbed().SetDescription("That page does not exist! %s has %d pages.", category.Name, len(helpPages))
	}
	if len(helpPages) > 1 {
		if cfg.Selfbot.CommandMode == "WebhookEmbed" {
			emb.SetSubtext(fmt.Sprintf("Page: %d/%d\n\nUse %shelp %s <page> to go to the next page", page+1, len(helpPages), cmdRouter.GetConfig().PrefixFunc()[0], rawArgument))
		} else {
			emb.SetSubtext(fmt.Sprintf("Page: %d/%d\n\nUse \x1b[34;1m%shelp %s <page>\x1b[0m to go to the next page", page+1, len(helpPages), cmdRouter.GetConfig().PrefixFunc()[0], rawArgument))
		}
	}
	return emb
}

// renderDefaultGeneralHelpEmbed renders the general help embed on the given page
func (c *HelpCommands) renderDefaultGeneralHelpEmbed(cmdRouter router.Router) *embed.Embed {
	cfg, _ := config.GetConfig()
	categoriesString := ""

	for _, category := range cmdRouter.GetCategories() {
		// Get the commands in the category, will be empty otherwise
		category = cmdRouter.GetCategory(category.Name)
		categoriesString += fmt.Sprintf("%s - %s\n", category.Name+strings.Repeat(" ", 14-len(category.Name)), category.Description)
	}

	prefix := cmdRouter.GetConfig().PrefixFunc()[0]

	emb := c.services.Embed.NewEmbed().
		SetDescription("LightningBot %s Help\nAmount of Commands: %d, <> is required, [] is optional\n", constants.VERSION, len(cmdRouter.GetCmds()))

	if cfg.Selfbot.CommandMode == "WebhookEmbed" {
		emb.Description += "\n\n**Categories**\n" + categoriesString
		emb.SetSubtext(strings.ReplaceAll(
			"Use [p]help <category name> or [p]help <command name> to get more info.\nUse [p]search <query> to search commands.\nUse [p]repeat to repeat the last command.",
			"[p]",
			prefix,
		))
	} else {
		emb.AddField("Categories", categoriesString, true, true)
		emb.SetSubtext(strings.ReplaceAll(
			"Use \x1b[34;1m[p]help <category name>\x1b[0m or \x1b[34;1m[p]help <command name>\x1b[0m to get more info.\nUse \x1b[34;1m[p]search <query>\x1b[0m to search commands.\nUse \x1b[34;1m[p]repeat\x1b[0m to repeat the last command.",
			"[p]",
			prefix,
		))
	}

	return emb
}

func (c *HelpCommands) search(ctx *router.Ctx) {
	query := ctx.Arguments.Get(0).Raw()
	prefix := ctx.Router.GetConfig().PrefixFunc()[0]

	var closeCommands []string
	for _, command := range ctx.Router.GetCmds() {
		if len(closeCommands) >= 10 {
			break
		}
		distance := utils.Misc.GetDistance(command.Name, query)
		if distance < 2 || strings.Contains(command.Name, query) {
			closeCommands = append(closeCommands, prefix+command.Name+" "+command.Arguments+" - "+command.Description)
		}
	}

	if len(closeCommands) == 0 {
		c.services.Embed.NewEmbed().
			SetDescription("No commands found for query `%s`", query).
			Send(ctx)
		return
	}

	c.services.Embed.NewEmbed().
		SetDescription("Found %d commands for query `%s`", len(closeCommands), query).
		AddField("Commands", strings.Join(closeCommands, "\n"), true, true).
		Send(ctx)
}

func (c *HelpCommands) repeat(ctx *router.Ctx) {
	usedCommands := c.services.EventReactions.GetUsedCommands()

	if len(usedCommands) < 2 {
		c.services.Embed.NewEmbed().
			SetDescription("No commands have been used yet").
			Send(ctx)
		return
	}

	lastCommand := usedCommands[len(usedCommands)-2]

	channel, err := ctx.Session.State.Channel(ctx.Message.ChannelID)
	if err != nil {
		c.services.Embed.NewEmbed().
			SetDescription("Failed to get channel: %s", err.Error()).
			Send(ctx)
		return
	}

	ctx.Router.InvokeRemoteCmd(ctx.Session, channel, ctx.Message.Author, lastCommand.Command, lastCommand.Arguments.Raw())
}
