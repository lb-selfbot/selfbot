package commands

import (
	"fmt"

	"github.com/LightningDev1/LB-Selfbot/embed"
	"github.com/LightningDev1/LB-Selfbot/slashcommands"
	"github.com/LightningDev1/LB-Selfbot/utils"
	"github.com/LightningDev1/dgc"
	"github.com/LightningDev1/discordgo"
)

type SlashCommands struct{}

func (s SlashCommands) RegisterCommands(router *dgc.Router) {
	router.StartCategory("SlashCommands", "Slash command utilities")

	router.RegisterCmd(&dgc.Command{
		Name:        "searchslash",
		Description: "Search a slash command",
		Usage:       "[p]searchslash <#channel> <query>",
		Handler:     s.searchslash,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "executeslash",
		Description: "Run a slash command",
		Usage:       "[p]executeslash <#channel> <command> [options...]",
		Handler:     s.executeslash,
	})
}

func (s SlashCommands) searchslash(ctx *dgc.Ctx) {
	channel, ok := utils.Discord.TryGetChannelFromMention(ctx, ctx.Arguments.Get(0))
	if !ok {
		ctx.RespondText("Invalid channel")
		return
	}

	query := utils.Arguments["<query>"]

	commands, _, err := ctx.Session.ChannelApplicationCommandsSearch(&discordgo.ApplicationCommandsSearchParams{
		ChannelID:           channel.ID,
		Type:                1,
		Query:               query,
		Limit:               20,
		IncludeApplications: false,
	})
	if err != nil {
		ctx.RespondText("Error getting commands: " + err.Error())
		return
	}

	if len(commands) == 0 {
		ctx.RespondText("No commands found")
		return
	}

	desc := fmt.Sprintf("Found %d slash command(s) in #%s", len(commands), channel.Name) + "\n\n"

	for _, cmd := range commands {
		botName := "Unknown"

		bot, ok := utils.Discord.GetUser(ctx.Session, cmd.ApplicationID)
		if ok {
			botName = bot.Username
		}

		desc += fmt.Sprintf("`/%s` - %s - %s\n", cmd.Name, botName, cmd.Description)
	}

	embed.NewEmbed().
		SetTitle("Slash Commands: " + query).
		SetDescription(desc).
		Send(ctx)
}

func (s SlashCommands) executeslash(ctx *dgc.Ctx) {
	channel, ok := utils.Discord.TryGetChannelFromMention(ctx, ctx.Arguments.Get(0))
	if !ok {
		ctx.RespondText("Invalid channel")
		return
	}

	commands, _, err := ctx.Session.ChannelApplicationCommandsSearch(&discordgo.ApplicationCommandsSearchParams{
		ChannelID:           channel.ID,
		Type:                1,
		Query:               utils.Arguments["<command>"],
		Limit:               7,
		IncludeApplications: false,
	})
	if err != nil {
		ctx.RespondText("Error getting commands: " + err.Error())
		return
	}

	if len(commands) == 0 {
		ctx.RespondText("No commands found")
		return
	}

	for _, cmd := range commands {
		opts, err := slashcommands.ArgumentsToCommandOptions(cmd.Options, ctx.Arguments, 2)
		if err != nil {
			ctx.RespondText("Error parsing options: " + err.Error())
			return
		}

		cmd.Options = opts

		err = ctx.Session.Interact(&discordgo.InteractData{
			Type:          int(discordgo.InteractionApplicationCommand),
			Data:          cmd,
			ChannelID:     channel.ID,
			GuildID:       channel.GuildID,
			ApplicationID: cmd.ApplicationID,
		})
		if err != nil {
			ctx.RespondText("Error executing command: " + err.Error())
			return
		}
	}
}
