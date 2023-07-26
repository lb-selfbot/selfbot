package commands

import (
	"fmt"
	"math"
	"strings"

	"github.com/LightningDev1/LB-Selfbot/config"
	"github.com/LightningDev1/LB-Selfbot/constants"
	"github.com/LightningDev1/LB-Selfbot/embed"
	"github.com/LightningDev1/LB-Selfbot/eventhandlers"
	"github.com/LightningDev1/LB-Selfbot/utils"
	"github.com/LightningDev1/dgc"
	"github.com/LightningDev1/discordgo"
)

type HelpCommand struct{}

type HelpMessage struct {
	Page     int
	Category *dgc.Category
	Command  *dgc.Command
}

// RegisterHelpCommand registers the help command to the router
func (h *HelpCommand) RegisterHelpCommand(router *dgc.Router, session *discordgo.Session) {
	// Register the default help command
	router.StopCategory()
	router.RegisterCmd(&dgc.Command{
		Name:    "help",
		Handler: generalHelpCommand,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "search",
		Usage:       "[p]search <query>",
		Description: "Search for a command",
		Handler:     h.search,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "repeat",
		Usage:       "[p]repeat",
		Description: "Repeat the last used command",
		Handler:     h.repeat,
	})

	for _, category := range router.Categories {
		name := strings.ReplaceAll(strings.ToLower(category.Name), " ", "")

		router.RegisterCmd(&dgc.Command{
			Name:        name,
			Usage:       "[p]" + name,
			Description: fmt.Sprintf("Get help for the %s category", category.Name),
			Handler:     generateCategoryHelpFunc(category, name),
		})
	}
}

func generateCategoryHelpFunc(category *dgc.Category, name string) dgc.ExecutionHandler {
	return func(ctx *dgc.Ctx) {
		page, _ := ctx.Arguments.Get(0).AsInt()

		renderCategoryHelp(ctx.Router, category, name, page).Send(ctx)
	}
}

// generalHelpCommand handles the general help command
func generalHelpCommand(ctx *dgc.Ctx) {
	// Check if the user provided an argument
	if ctx.Arguments.Amount() > 0 {
		specificHelpCommand(ctx)
		return
	}

	// Send the general help embed
	renderDefaultGeneralHelpEmbed(ctx.Router).Send(ctx)
}

func getNameFromArgument(arg string) string {
	arg = strings.ReplaceAll(arg, "eventreactions", "EventReactions")
	arg = strings.ReplaceAll(arg, "eventreaction", "EventReactions")

	arg = strings.ReplaceAll(arg, "nsfw", "NSFW")

	arg = strings.ReplaceAll(arg, "customactivity", "CustomActivity")

	arg = strings.ReplaceAll(arg, "slashcommands", "SlashCommands")
	arg = strings.ReplaceAll(arg, "slashcommand", "SlashCommands")

	arg = strings.ReplaceAll(arg, "ai", "AI")

	if arg == "admin" {
		arg = "Admininstration"
	}

	return arg
}

// specificHelpCommand handles the specific help command
func specificHelpCommand(ctx *dgc.Ctx) {
	oldName := strings.ToLower(ctx.Arguments.Get(0).Raw())
	name := getNameFromArgument(oldName)
	page, _ := ctx.Arguments.Get(1).AsInt()
	prefix := ctx.Router.GetPrefixes()[0]

	if category := ctx.Router.GetCategory(utils.String.Capitalize(name)); category != nil {
		renderCategoryHelp(ctx.Router, category, oldName, page).Send(ctx)
		return
	}

	if command := ctx.Router.GetCmd(name); command != nil {
		renderCommandHelp(ctx.Router, command).Send(ctx)
		return
	}

	// get close commands and categories
	var closeCommands []string
	for _, command := range ctx.Router.Commands {
		if len(closeCommands) >= 5 {
			break
		}
		distance := utils.Misc.GetDistance(command.Name, name)
		if distance < 2 || strings.Contains(command.Name, name) {
			closeCommands = append(closeCommands, strings.ReplaceAll(command.Usage, "[p]", prefix)+" - "+command.Description)
		}
	}

	var closeCategories []string
	for _, category := range ctx.Router.Categories {
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

	embed.NewEmbed().SetDescription(description).Send(ctx)
}

func renderCommandHelp(router *dgc.Router, command *dgc.Command) *embed.Embed {
	categoryName := ""
	// Try to get the category the command belongs in
	for _, category := range router.Categories {
		category = router.GetCategory(category.Name)
		for _, cmd := range category.Commands {
			if cmd.Name != command.Name {
				continue
			}

			categoryName = category.Name
			break
		}
	}
	return embed.NewEmbed().
		SetDescription("Command Help: "+command.Name).
		AddField("Command Name", command.Name).
		AddField("Usage", "<> is required, [] is optional\n\n"+strings.ReplaceAll(command.Usage, "[p]", router.GetPrefixes()[0])).
		AddField("Description", command.Description).
		AddField("Category", categoryName)
}

func renderCategoryHelp(router *dgc.Router, category *dgc.Category, rawArgument string, page int) *embed.Embed {
	helpPages := map[int]string{}
	currentPage := 0
	for _, command := range category.Commands {
		splitUsage := strings.SplitN(command.Usage, " ", 2)
		usage := ""
		if len(splitUsage) > 1 {
			usage = " " + splitUsage[1]
		}
		helpText := fmt.Sprintf("%s%s%s - %s\n", router.GetPrefixes()[0], command.Name, usage, command.Description)
		if len(helpPages[currentPage])+len(helpText) > 1650 {
			currentPage++
		}
		helpPages[currentPage] += helpText
	}
	page = int(math.Max(0, float64(page-1)))

	emb := embed.NewEmbed().
		SetDescription(fmt.Sprintf("<> is required, [] is optional\n\nAmount of %s commands: %d", category.Name, len(category.Commands))).
		SetSubtext(fmt.Sprintf("Page: %d/%d", page+1, len(helpPages)))

	cfg, _ := config.GetConfig()

	if page < len(helpPages) {
		if cfg.Selfbot.CommandMode == "WebhookEmbed" {
			emb.Description += "\n\n**Commands**\n" + helpPages[page]
		} else {
			emb.AddField("Commands", helpPages[page], true, true)
		}
	} else {
		return embed.NewEmbed().SetDescription(fmt.Sprintf("That page does not exist! %s has %d pages.", category.Name, len(helpPages)))
	}
	if len(helpPages) > 1 {
		if cfg.Selfbot.CommandMode == "WebhookEmbed" {
			emb.SetSubtext(fmt.Sprintf("Page: %d/%d\n\nUse %shelp %s <page> to go to the next page", page+1, len(helpPages), router.GetPrefixes()[0], rawArgument))
		} else {
			emb.SetSubtext(fmt.Sprintf("Page: %d/%d\n\nUse \x1b[34;1m%shelp %s <page>\x1b[0m to go to the next page", page+1, len(helpPages), router.GetPrefixes()[0], rawArgument))
		}
	}
	return emb
}

// renderDefaultGeneralHelpEmbed renders the general help embed on the given page
func renderDefaultGeneralHelpEmbed(router *dgc.Router) *embed.Embed {
	cfg, _ := config.GetConfig()
	categoriesString := ""

	for _, category := range router.Categories {
		// Get the commands in the category, will be empty otherwise
		category = router.GetCategory(category.Name)
		categoriesString += fmt.Sprintf("%s - %s\n", category.Name+strings.Repeat(" ", 15-len(category.Name)), category.Description)
	}

	prefix := router.GetPrefixes()[0]

	emb := embed.NewEmbed().
		SetDescription(fmt.Sprintf("LightningBot %s Help\nAmount of Commands: %d, <> is required, [] is optional\n", constants.VERSION, len(router.Commands)))

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

func (HelpCommand) search(ctx *dgc.Ctx) {
	query := utils.Arguments["<query>"]
	prefix := ctx.Router.GetPrefixes()[0]
	var closeCommands []string
	for _, command := range ctx.Router.Commands {
		if len(closeCommands) >= 10 {
			break
		}
		distance := utils.Misc.GetDistance(command.Name, query)
		if distance < 2 || strings.Contains(command.Name, query) {
			closeCommands = append(closeCommands, strings.ReplaceAll(command.Usage, "[p]", prefix)+" - "+command.Description)
		}
	}
	if len(closeCommands) == 0 {
		embed.NewEmbed().
			SetDescription(fmt.Sprintf("No commands found for query `%s`", query)).
			Send(ctx)
		return
	}
	embed.NewEmbed().
		SetDescription(fmt.Sprintf("Found %d commands for query `%s`", len(closeCommands), query)).
		AddField("Commands", strings.Join(closeCommands, "\n"), true, true).
		Send(ctx)
}

func (HelpCommand) repeat(ctx *dgc.Ctx) {
	usedCommands := eventhandlers.EventHandlers.Settings.UsedCommands
	if len(eventhandlers.EventHandlers.Settings.UsedCommands) < 2 {
		embed.NewEmbed().
			SetDescription("No commands have been used yet").
			Send(ctx)
		return
	}
	lastCommand := usedCommands[len(usedCommands)-2]
	lastCommand.Event = ctx.Event

	// Restore arguments
	arguments, _ := utils.Discord.GetArguments(lastCommand)
	utils.Arguments = arguments

	// Invoke the last command
	lastCommand.Command.Handler(lastCommand)
}
