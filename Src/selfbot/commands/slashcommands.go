package commands

import (
	"fmt"
	"time"

	"github.com/LightningDev1/LB-Selfbot/selfbot/embed"
	"github.com/LightningDev1/LB-Selfbot/selfbot/slashcommands"
	"github.com/LightningDev1/LB-Selfbot/selfbot/utils"
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
		Flags:       []string{"DANGEROUS"},
		Handler:     s.searchslash,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "executeslash",
		Description: "Run a slash command",
		Usage:       "[p]executeslash <#channel> <@bot> <command> [options...]",
		Flags:       []string{"DANGEROUS"},
		Handler:     s.executeslash,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "slashcommandjobhelp",
		Description: "More information about slash command jobs",
		Usage:       "[p]slashcommandjobhelp",
		Flags:       []string{"DANGEROUS"},
		Handler:     s.slashcommandjobhelp,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "listscjobs",
		Description: "List all slash command jobs",
		Usage:       "[p]listscjobs",
		Flags:       []string{"DANGEROUS"},
		Handler:     s.listscjobs,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "getscjob",
		Description: "Get info about a slash command job",
		Usage:       "[p]getscjob <job id>",
		Flags:       []string{"DANGEROUS"},
		Handler:     s.getscjob,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "createscjob",
		Description: "Create a slash command job",
		Usage:       "[p]createscjob <#channel> <@bot> <interval> <command> [options...]",
		Flags:       []string{"DANGEROUS"},
		Handler:     s.createscjob,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "deletescjob",
		Description: "Delete a slash command job",
		Usage:       "[p]deletescjob <job id>",
		Flags:       []string{"DANGEROUS"},
		Handler:     s.deletescjob,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "setscjobcommand",
		Description: "Set a job's command",
		Usage:       "[p]setscjobcommand <job id> <@bot> <command> [options...]",
		Flags:       []string{"DANGEROUS"},
		Handler:     s.setscjobcommand,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "setscjobchannel",
		Description: "Set a job's channel",
		Usage:       "[p]setscjobcommand <job id> <#channel>",
		Flags:       []string{"DANGEROUS"},
		Handler:     s.setscjobchannel,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "setscjobinterval",
		Description: "Set a job's interval (in seconds)",
		Usage:       "[p]setscjobinterval <job id> <interval>",
		Flags:       []string{"DANGEROUS"},
		Handler:     s.setscjobinterval,
	})
}

func (s SlashCommands) searchslash(ctx *dgc.Ctx) {
	channel, ok := utils.Discord.TryGetChannelFromMention(ctx, ctx.Arguments.Get(0))
	if !ok {
		ctx.RespondText("Could not find channel. Make sure to mention it or use the ID.")
		return
	}

	query := utils.Arguments["<query>"]

	commands, _, err := ctx.Session.ChannelApplicationCommandsSearch(&discordgo.ApplicationCommandsSearchParams{
		ChannelID:           channel.ID,
		Type:                1,
		Query:               query,
		Limit:               20,
		IncludeApplications: true,
	})
	if err != nil {
		ctx.RespondText("An error occurred while getting commands from Discord: " + err.Error())
		return
	}

	if len(commands) == 0 {
		ctx.RespondText("No commands have been found matching the query.")
		return
	}

	desc := fmt.Sprintf("Found %d slash command(s) in #%s", len(commands), channel.Name) + "\n\n"

	for _, cmd := range commands {
		botName := "Unknown"

		bot, ok := utils.Discord.GetUser(ctx.Session, cmd.BotUserID)
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
		ctx.RespondText("Could not find channel. Make sure to mention it or use the ID.")
		return
	}

	bot, ok := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(1))
	if !ok {
		ctx.RespondText("Could not find bot. Make sure to mention it or use the ID.")
		return
	}

	commands, _, err := ctx.Session.ChannelApplicationCommandsSearch(&discordgo.ApplicationCommandsSearchParams{
		ChannelID:           channel.ID,
		Type:                1,
		Query:               utils.Arguments["<command>"],
		Limit:               20,
		IncludeApplications: true,
	})
	if err != nil {
		ctx.RespondText("An error occurred while getting commands from Discord: " + err.Error())
		return
	}

	rightCommands := utils.Filter(commands, func(command *discordgo.ApplicationCommand) bool {
		return command.Name == utils.Arguments["<command>"] && command.BotUserID == bot.ID
	})

	if len(rightCommands) == 0 {
		ctx.RespondText("No commands have been found matching the query.")
		return
	}

	cmd := rightCommands[0]

	opts, err := slashcommands.ArgumentsToCommandOptions(cmd.Options, ctx.Arguments, 3)
	if err != nil {
		ctx.RespondText("An error occured while parsing your arguments to command options: " + err.Error())
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
		ctx.RespondText("An error occurred while executing the command: " + err.Error())
		return
	}
}

func (s SlashCommands) slashcommandjobhelp(ctx *dgc.Ctx) {
	embed.NewEmbed().
		SetDescription(
			"Slash command jobs are jobs that run slash commands on a set interval. " +
				"To get more information and examples, read the documentation at " +
				"https://docs.lightning-bot.com/commands/slash-commands",
		).
		Send(ctx)
}

func (s SlashCommands) listscjobs(ctx *dgc.Ctx) {
	embed := embed.NewEmbed().
		SetTitle("Slash Command Jobs").
		SetDescription(fmt.Sprintf("Found %d job", len(slashcommands.Jobs)))

	if len(slashcommands.Jobs) != 1 {
		embed.SetDescription(embed.Description + "s")
	}

	for _, job := range slashcommands.Jobs {
		cmd := job.CommandObject
		channel := job.ChannelObject
		guild := job.GuildObject

		args := job.RawArguments
		if len(args) > 0 {
			args = " " + args
		}

		text := fmt.Sprintf("Command: `/%s%s` (%s)\n", cmd.Name, args, cmd.ID)
		text += fmt.Sprintf("Channel: #%s (%s)\n", channel.Name, channel.ID)
		text += fmt.Sprintf("Server: %s (%s)\n", guild.Name, guild.ID)
		text += fmt.Sprintf("Interval: %s\n", utils.Misc.FormatDuration(job.Interval))

		embed.AddField("ID: "+job.ID, text)
	}

	embed.Send(ctx)
}

func (s SlashCommands) getscjob(ctx *dgc.Ctx) {
	id := utils.Arguments["<job id>"]

	job := slashcommands.GetJobByID(id)
	if job == nil {
		ctx.RespondText("Could not find a job with ID `" + id + "`")
		return
	}

	jobHandler := slashcommands.GetHandlerForJob(job)
	timeUntilNext := jobHandler.TimeUntilNextInteraction()

	cmd := job.CommandObject
	channel := job.ChannelObject
	guild := job.GuildObject

	args := job.RawArguments
	if len(args) > 0 {
		args = " " + args
	}

	embed.NewEmbed().
		SetTitle("Slash command job information").
		AddField("Command", fmt.Sprintf("`/%s%s` (%s)", cmd.Name, args, cmd.ID)).
		AddField("Channel", fmt.Sprintf("#%s (%s)", channel.Name, channel.ID)).
		AddField("Server", fmt.Sprintf("%s (%s)", guild.Name, guild.ID)).
		AddField("Interval", utils.Misc.FormatDuration(job.Interval)).
		AddField("Job handler is running", utils.Misc.FormatBoolYesNo(jobHandler.IsJobRunning())).
		AddField("Time until next interaction", utils.Misc.FormatDuration(timeUntilNext)).
		Send(ctx)
}

func (s SlashCommands) createscjob(ctx *dgc.Ctx) {
	channel, ok := utils.Discord.TryGetChannelFromMention(ctx, ctx.Arguments.Get(0))
	if !ok {
		ctx.RespondText("Could not find channel. Make sure to mention it or use the ID.")
		return
	}

	bot, ok := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(1))
	if !ok {
		ctx.RespondText("Could not find bot. Make sure to mention it or use the ID.")
		return
	}

	interval, err := time.ParseDuration(utils.Arguments["<interval>"])
	if err != nil {
		ctx.RespondText("An error occurred while parsing the interval. Make sure to use a valid duration (i.e. 30s, 5m, or 5h).")
		return
	}

	commandName := utils.Arguments["<command>"]

	commands, _, err := ctx.Session.ChannelApplicationCommandsSearch(&discordgo.ApplicationCommandsSearchParams{
		ChannelID:           channel.ID,
		Type:                1,
		Query:               commandName,
		Limit:               20,
		IncludeApplications: true,
	})
	if err != nil {
		ctx.RespondText("An error occurred while getting commands from Discord: " + err.Error())
		return
	}

	rightCommands := utils.Filter(commands, func(command *discordgo.ApplicationCommand) bool {
		return command.Name == commandName && command.BotUserID == bot.ID
	})

	if len(rightCommands) == 0 {
		ctx.RespondText("No commands have been found matching the query.")
		return
	}

	command := rightCommands[0]

	options, err := slashcommands.ArgumentsToCommandOptions(command.Options, ctx.Arguments, 4)
	if err != nil {
		ctx.RespondText("An error occured while parsing your arguments to command options: " + err.Error())
		return
	}

	command.Options = options

	job := &slashcommands.Job{
		CommandObject: slashcommands.DiscordObject{
			Name: command.Name,
			ID:   command.ID,
		},
		ApplicationObject: slashcommands.DiscordObject{
			Name: bot.Username,
			ID:   command.ApplicationID,
		},
		ChannelObject: slashcommands.DiscordObject{
			Name: channel.Name,
			ID:   channel.ID,
		},
		GuildObject: slashcommands.DiscordObject{
			Name: channel.GuildID,
			ID:   channel.GuildID,
		},
		RawArguments: utils.String.RemoveWords(ctx.Arguments.Raw(), 4),
		IntervalSec:  int(interval.Seconds()),

		// Set command so the job handler won't have to fetch it again
		Command: command,
	}

	err = job.Load(ctx.Session)
	if err != nil {
		ctx.RespondText("An error occurred while loading the job data: " + err.Error())
		return
	}

	err = slashcommands.AddJob(job)
	if err != nil {
		ctx.RespondText("An error occurred while adding the job: " + err.Error())
		return
	}

	embed.NewEmbed().
		SetDescription("Successfully created slash command job!").
		AddField("Job ID", job.ID).
		Send(ctx)
}

func (s SlashCommands) deletescjob(ctx *dgc.Ctx) {
	id := utils.Arguments["<job id>"]

	job := slashcommands.GetJobByID(id)
	if job == nil {
		ctx.RespondText("Could not find a job with ID `" + id + "`")
		return
	}

	err := slashcommands.RemoveJob(job)
	if err != nil {
		ctx.RespondText("An error occurred while removing the job: " + err.Error())
		return
	}

	embed.NewEmbed().
		SetDescription("Successfully deleted slash command job!").
		AddField("Job ID", job.ID).
		Send(ctx)
}

func (s SlashCommands) setscjobcommand(ctx *dgc.Ctx) {
	id := utils.Arguments["<job id>"]

	job := slashcommands.GetJobByID(id)
	if job == nil {
		ctx.RespondText("Could not find a job with ID `" + id + "`")
		return
	}

	bot, ok := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(1))
	if !ok {
		ctx.RespondText("Could not find bot. Make sure to mention it or use the ID.")
		return
	}

	commandName := utils.Arguments["<command>"]

	commands, _, err := ctx.Session.ChannelApplicationCommandsSearch(&discordgo.ApplicationCommandsSearchParams{
		ChannelID:           job.ChannelObject.ID,
		Type:                1,
		Query:               commandName,
		Limit:               20,
		IncludeApplications: true,
	})
	if err != nil {
		ctx.RespondText("An error occurred while getting commands from Discord: " + err.Error())
		return
	}

	rightCommands := utils.Filter(commands, func(command *discordgo.ApplicationCommand) bool {
		return command.Name == commandName && command.BotUserID == bot.ID
	})

	if len(rightCommands) == 0 {
		ctx.RespondText("No commands have been found matching the query.")
		return
	}

	command := rightCommands[0]

	options, err := slashcommands.ArgumentsToCommandOptions(command.Options, ctx.Arguments, 4)
	if err != nil {
		ctx.RespondText("An error occured while parsing your arguments to command options: " + err.Error())
		return
	}

	command.Options = options

	job.CommandObject = slashcommands.DiscordObject{
		Name: command.Name,
		ID:   command.ID,
	}

	job.ApplicationObject = slashcommands.DiscordObject{
		Name: bot.Username,
		ID:   command.ApplicationID,
	}

	job.RawArguments = utils.String.RemoveWords(ctx.Arguments.Raw(), 3)

	job.Command = command

	err = job.Load(ctx.Session)
	if err != nil {
		ctx.RespondText("An error occurred while loading the job data: " + err.Error())
		return
	}

	err = slashcommands.SaveJobs()
	if err != nil {
		ctx.RespondText("An error occurred while saving the jobs: " + err.Error())
		return
	}

	embed.NewEmbed().
		SetDescription("Successfully updated slash command job").
		AddField("Job ID", job.ID).
		Send(ctx)
}

func (s SlashCommands) setscjobchannel(ctx *dgc.Ctx) {
	id := utils.Arguments["<job id>"]

	job := slashcommands.GetJobByID(id)
	if job == nil {
		ctx.RespondText("Could not find a job with ID `" + id + "`")
		return
	}

	channel, ok := utils.Discord.TryGetChannelFromMention(ctx, ctx.Arguments.Get(1))
	if !ok {
		ctx.RespondText("Could not find channel. Make sure to mention it or use the ID.")
		return
	}

	job.ChannelObject = slashcommands.DiscordObject{
		Name: channel.Name,
		ID:   channel.ID,
	}

	err := job.Load(ctx.Session)
	if err != nil {
		ctx.RespondText("An error occurred while loading the job data: " + err.Error())
		return
	}

	err = slashcommands.SaveJobs()
	if err != nil {
		ctx.RespondText("An error occurred while saving the jobs: " + err.Error())
		return
	}

	embed.NewEmbed().
		SetDescription("Successfully updated slash command job").
		AddField("Job ID", job.ID).
		Send(ctx)
}

func (s SlashCommands) setscjobinterval(ctx *dgc.Ctx) {
	id := utils.Arguments["<job id>"]

	job := slashcommands.GetJobByID(id)
	if job == nil {
		ctx.RespondText("Could not find a job with ID `" + id + "`")
		return
	}

	interval, err := time.ParseDuration(utils.Arguments["<interval>"])
	if err != nil {
		ctx.RespondText("An error occurred while parsing the interval. Make sure to use a valid duration (i.e. 30s, 5m, or 5h).")
		return
	}

	job.IntervalSec = int(interval.Seconds())

	err = job.Load(ctx.Session)
	if err != nil {
		ctx.RespondText("An error occurred while loading the job data: " + err.Error())
		return
	}

	err = slashcommands.SaveJobs()
	if err != nil {
		ctx.RespondText("An error occurred while saving the jobs: " + err.Error())
		return
	}

	embed.NewEmbed().
		SetDescription("Successfully updated slash command job").
		AddField("Job ID", job.ID).
		Send(ctx)
}
