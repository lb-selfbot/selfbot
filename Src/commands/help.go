package commands

import (
	"fmt"
	"github.com/LightningDev1/LB-Selfbot/constants"
	"github.com/LightningDev1/LB-Selfbot/embed"
	"github.com/LightningDev1/dgc"
	"strings"
	"time"

	"github.com/LightningDev1/discordgo"
)

type HelpCommand struct{}

type HelpMessage struct {
	Page     int
	Category *dgc.Category
	Command  *dgc.Command
}

// RegisterHelpCommand registers the help command to the router
func (HelpCommand) RegisterHelpCommand(router *dgc.Router, session *discordgo.Session) {
	// Initialize the help messages storage
	router.InitializeStorage("helpMessages")

	// Initialize the reaction add listener
	session.AddHandler(func(session *discordgo.Session, event *discordgo.MessageReactionRemove) {
		// Define useful variables
		channelID := event.ChannelID
		messageID := event.MessageID
		userID := event.UserID

		// Check whether the reaction was added by the current user
		if event.UserID != session.State.User.ID && router.SelfBot {
			return
		}

		// Check whether the message is a help message
		rawPage, ok := router.Storage["helpMessages"].Get(channelID + ":" + messageID + ":" + userID)
		if !ok {
			return
		}
		page := rawPage.(int)
		if page <= 0 {
			return
		}

		// Check which reaction was added
		//reactionName := event.Emoji.Name
		//switch reactionName {
		//case "⬅️":
		//	// Update the help message
		//	embed, newPage := renderDefaultGeneralHelpEmbed(router, page-1)
		//	page = newPage
		//	//session.ChannelMessageEditEmbed(channelID, messageID, embed)
		//
		//	// Remove the reaction
		//	session.MessageReactionRemove(channelID, messageID, reactionName, userID)
		//	break
		//case "❌":
		//	// Delete the help message
		//	session.ChannelMessageDelete(channelID, messageID)
		//	break
		//case "➡️":
		//	// Update the help message
		//	//embed, newPage := renderDefaultGeneralHelpEmbed(router, page+1)
		//	//page = newPage
		//	//session.ChannelMessageEditEmbed(channelID, messageID, embed)
		//
		//	// Remove the reaction
		//	session.MessageReactionRemove(channelID, messageID, reactionName, userID)
		//	break
		//}

		// Update the stores page
		router.Storage["dgc_helpMessages"].Set(channelID+":"+messageID+":"+event.UserID, page)
	})

	// Register the default help command
	router.StopCategory()
	router.RegisterCmd(&dgc.Command{
		Name:    "help",
		Handler: generalHelpCommand,
	})
}

// generalHelpCommand handles the general help command
func generalHelpCommand(ctx *dgc.Ctx) {
	// Check if the user provided an argument
	if ctx.Arguments.Amount() > 0 {
		specificHelpCommand(ctx)
		return
	}

	// Define useful variables
	channelID := ctx.Event.ChannelID

	// Send the general help embed
	message := renderDefaultGeneralHelpEmbed(ctx.Router).Send(ctx)

	// Add the reactions to the message
	//_ = session.MessageReactionAdd(channelID, message.ID, "❌")
	//_ = session.MessageReactionAdd(channelID, message.ID, "➡️")

	// Define the message as a help message
	ctx.Router.Storage["helpMessages"].Set(channelID+":"+message.ID+":"+ctx.Event.Author.ID, 1)
}

// specificHelpCommand handles the specific help command
func specificHelpCommand(ctx *dgc.Ctx) {
	// Define the command names
	commandNames := strings.Split(ctx.Arguments.Raw(), " ")

	// Define the command
	var command *dgc.Command
	for index, commandName := range commandNames {
		if index == 0 {
			command = ctx.Router.GetCmd(commandName)
			continue
		}
		if command == nil {
			break
		}
		command = command.GetSubCmd(commandName)
	}

	// Send the help embed
	ctx.Session.ChannelMessageSendEmbed(ctx.Event.ChannelID, renderDefaultSpecificHelpEmbed(ctx, command))
}

// renderDefaultGeneralHelpEmbed renders the general help embed on the given page
func renderDefaultGeneralHelpEmbed(router *dgc.Router) *embed.Embed {
	// Define useful variables

	categoriesString := "ansi\n\x1b[1;31m"
	for _, category := range router.Categories {
		// Get the commands in the category, will be empty otherwise
		category = router.GetCategory(category.Name)
		categoriesString += fmt.Sprintf("%s - %s\n", category.Name+strings.Repeat(" ", 14-len(category.Name)), category.Description)
	}
	return embed.NewEmbed().
		SetDescription(fmt.Sprintf("LightningBot %s Help\nAmount of Commands: %d, <> is required, [] is optional\n", constants.VersionString, len(router.Commands))).
		AddField("Categories", categoriesString, true, true).
		SetSubtext(fmt.Sprintf("Use \x1b[34;1m%shelp <category name>\x1b[0m or \u001B[34;1m%shelp <command name>\u001B[0m to get more info.", router.GetPrefixes()[0], router.GetPrefixes()[0]))
}

// renderDefaultSpecificHelpEmbed renders the specific help embed of the given command
func renderDefaultSpecificHelpEmbed(ctx *dgc.Ctx, command *dgc.Command) *discordgo.MessageEmbed {
	// Define useful variables
	prefix := ctx.Router.Prefixes[0]

	// Check if the command is invalid
	if command == nil {
		return &discordgo.MessageEmbed{
			Type:      "rich",
			Title:     "Error",
			Timestamp: time.Now().Format(time.RFC3339),
			Color:     0xff0000,
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:   "Message",
					Value:  "```The given command doesn't exist. Type `" + prefix + "help` for a list of available commands.```",
					Inline: false,
				},
			},
		}
	}

	// Define the sub commands string
	subCommands := "No sub commands"
	if len(command.SubCommands) > 0 {
		subCommandNames := make([]string, len(command.SubCommands))
		for index, subCommand := range command.SubCommands {
			subCommandNames[index] = subCommand.Name
		}
		subCommands = "`" + strings.Join(subCommandNames, "`, `") + "`"
	}

	// Define the aliases string
	aliases := "No aliases"
	if len(command.Aliases) > 0 {
		aliases = "`" + strings.Join(command.Aliases, "`, `") + "`"
	}

	// Return the embed
	return &discordgo.MessageEmbed{
		Type:        "rich",
		Title:       "Command Information",
		Description: "Displaying the information for the `" + command.Name + "` command.",
		Timestamp:   time.Now().Format(time.RFC3339),
		Color:       0xffff00,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Name",
				Value:  "`" + command.Name + "`",
				Inline: false,
			},
			{
				Name:   "Sub Commands",
				Value:  subCommands,
				Inline: false,
			},
			{
				Name:   "Aliases",
				Value:  aliases,
				Inline: false,
			},
			{
				Name:   "Description",
				Value:  "```" + command.Description + "```",
				Inline: false,
			},
			{
				Name:   "Usage",
				Value:  "```" + prefix + command.Usage + "```",
				Inline: false,
			},
			{
				Name:   "Example",
				Value:  "```" + prefix + command.Example + "```",
				Inline: false,
			},
		},
	}
}
