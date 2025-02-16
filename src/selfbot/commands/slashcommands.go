package commands

import (
	"fmt"
	"strings"
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
		Usage:       "[p]searchslash <query> [server id]",
		Flags:       []string{CmdFlagDangerous},
		Handler:     s.searchslash,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "executeslash",
		Description: "Run a slash command",
		Usage:       "[p]executeslash <#channel> <@bot> <command> [options...]",
		Flags:       []string{CmdFlagDangerous},
		Handler:     s.executeslash,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "slashcommandjobhelp",
		Description: "More information about slash command jobs",
		Usage:       "[p]slashcommandjobhelp",
		Flags:       []string{CmdFlagDangerous},
		Handler:     s.slashcommandjobhelp,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "listscjobs",
		Description: "List all slash command jobs",
		Usage:       "[p]listscjobs",
		Flags:       []string{CmdFlagDangerous},
		Handler:     s.listscjobs,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "getscjob",
		Description: "Get info about a slash command job",
		Usage:       "[p]getscjob <job name>",
		Flags:       []string{CmdFlagDangerous},
		Handler:     s.getscjob,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "createscjob",
		Description: "Create a slash command job",
		Usage:       "[p]createscjob <job name> <#channel> <@bot> <interval> <command> [options...]",
		Flags:       []string{CmdFlagDangerous},
		Handler:     s.createscjob,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "deletescjob",
		Description: "Delete a slash command job",
		Usage:       "[p]deletescjob <job name>",
		Flags:       []string{CmdFlagDangerous},
		Handler:     s.deletescjob,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "setscjobcommand",
		Description: "Set a job's command",
		Usage:       "[p]setscjobcommand <job name> <@bot> <command> [options...]",
		Flags:       []string{CmdFlagDangerous},
		Handler:     s.setscjobcommand,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "setscjobchannel",
		Description: "Set a job's channel",
		Usage:       "[p]setscjobcommand <job name> <#channel>",
		Flags:       []string{CmdFlagDangerous},
		Handler:     s.setscjobchannel,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "setscjobinterval",
		Description: "Set a job's interval",
		Usage:       "[p]setscjobinterval <job name> <interval>",
		Flags:       []string{CmdFlagDangerous},
		Handler:     s.setscjobinterval,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "setscjobpaused",
		Description: "Pause or unpause a job",
		Usage:       "[p]setscjobpaused <job name> <on/off>",
		Flags:       []string{CmdFlagDangerous},
		Handler:     s.setscjobpaused,
	})
}

func (s SlashCommands) searchslash(ctx *dgc.Ctx) {
	serverID := utils.Arguments["[server id]"]
	if serverID == "" {
		serverID = ctx.Event.GuildID
	}
	if serverID == "" {
		ctx.RespondText("Please provide a server ID or use the command in a server.")
		return
	}

	guild, err := utils.Discord.GetGuild(ctx.Session, serverID)
	if err != nil {
		ctx.RespondText("Error: Could not find server.")
		return
	}

	commands, applications, err := ctx.Session.GuildApplicationCommandIndex(serverID)
	if err != nil {
		ctx.RespondText("An error occurred while getting commands: " + err.Error())
		return
	}

	commands = utils.Filter(commands, func(command *discordgo.ApplicationCommand) bool {
		return utils.Misc.GetDistance(command.Name, utils.Arguments["<query>"]) < 2 ||
			command.Name == utils.Arguments["<query>"]
	})

	if len(commands) == 0 {
		ctx.RespondText("No commands have been found matching the query.")
		return
	}

	desc := fmt.Sprintf("Found %d slash command(s) in %s", len(commands), guild.Name) + "\n\n"

	for _, cmd := range commands {
		botName := "Unknown Bot"

		for _, app := range applications {
			if app.ID == cmd.ApplicationID {
				botName = app.Name
			}
		}

		desc += fmt.Sprintf("`/%s` - %s - %s\n", cmd.Name, botName, cmd.Description)
	}

	embed.NewEmbed().
		SetTitle("Slash Commands: " + utils.Arguments["<query>"]).
		SetDescription(desc).
		Send(ctx)
}

func (s SlashCommands) executeslash(ctx *dgc.Ctx) {
	channel, ok := utils.Discord.TryGetChannelFromMention(ctx, ctx.Arguments.Get(0))
	if !ok {
		ctx.RespondText("Could not find channel. Make sure to mention it or use the ID.")
		return
	}

	if channel.GuildID == "" {
		ctx.RespondText("This command can only be used in servers.")
		return
	}

	bot, ok := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(1))
	if !ok {
		ctx.RespondText("Could not find bot. Make sure to mention it or use the ID.")
		return
	}

	commands, applications, err := ctx.Session.GuildApplicationCommandIndex(channel.GuildID)
	if err != nil {
		ctx.RespondText("An error occurred while getting commands: " + err.Error())
		return
	}

	appID := bot.ID
	for _, app := range applications {
		if app.BotID == bot.ID {
			appID = app.ID
		}
	}

	commands = utils.Filter(commands, func(command *discordgo.ApplicationCommand) bool {
		return command.Name == utils.Arguments["<command>"] && command.ApplicationID == appID
	})

	if len(commands) == 0 {
		ctx.RespondText("No commands have been found matching the query.")
		return
	}

	cmd := commands[0]

	opts, err := slashcommands.ArgumentsToCommandOptions(cmd.Options, ctx.Arguments, 3)
	if err != nil {
		ctx.RespondText("An error occured while parsing your arguments to command options: " + err.Error())
		return
	}

	err = ctx.Session.Interact(&discordgo.InteractData{
		Type: int(discordgo.InteractionApplicationCommand),
		Data: map[string]any{
			"version":             cmd.Version,
			"id":                  cmd.ID,
			"name":                cmd.Name,
			"type":                cmd.Type,
			"options":             opts,
			"application_command": cmd,
			"attachments":         []any{},
		},
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
		jobID := slashcommands.JobIdentifier(job)
		handler := slashcommands.GetJobHandler(job.Name, jobID)

		cmd := job.CommandObject
		channel := job.ChannelObject
		guild := job.GuildObject

		args := job.RawArguments
		if len(args) > 0 {
			args = " " + args
		}

		paused := utils.Misc.FormatBoolYesNo(job.Paused)
		running := utils.Misc.FormatBoolYesNo(handler.IsRunning())

		text := fmt.Sprintf("Command: `/%s%s` (%s)\n", cmd.Name, args, cmd.ID)
		text += fmt.Sprintf("Channel: #%s (%s)\n", channel.Name, channel.ID)
		text += fmt.Sprintf("Server: %s (%s)\n", guild.Name, guild.ID)
		text += fmt.Sprintf("Interval: %s\n", utils.Misc.FormatDuration(job.Interval))
		text += fmt.Sprintf("Paused: %s | Is Running: %s\n", paused, running)

		embed.AddField("Name: "+job.Name, text)
	}

	embed.Send(ctx)
}

func (s SlashCommands) getscjob(ctx *dgc.Ctx) {
	name := utils.Arguments["<job name>"]

	job := slashcommands.GetJobByName(name)
	if job == nil {
		ctx.RespondText("Could not find a job with name `" + name + "`")
		return
	}

	id := slashcommands.JobIdentifier(job)
	handler := slashcommands.GetJobHandler(job.Name, id)

	timeUntilNext := time.Until(job.NextExecution())

	cmd := job.CommandObject
	channel := job.ChannelObject
	guild := job.GuildObject

	args := job.RawArguments
	if len(args) > 0 {
		args = " " + args
	}

	embed.NewEmbed().
		SetTitle("Slash command job information").
		AddField("Name", job.Name).
		AddField("Command", fmt.Sprintf("`/%s%s` (%s)", cmd.Name, args, cmd.ID)).
		AddField("Channel", fmt.Sprintf("#%s (%s)", channel.Name, channel.ID)).
		AddField("Server", fmt.Sprintf("%s (%s)", guild.Name, guild.ID)).
		AddField("Interval", utils.Misc.FormatDuration(job.Interval)).
		AddField("Paused", utils.Misc.FormatBoolYesNo(job.Paused)).
		AddField("Job handler is running", utils.Misc.FormatBoolYesNo(handler.IsRunning())).
		AddField("Time until next interaction", utils.Misc.FormatDuration(timeUntilNext)).
		Send(ctx)
}

func (s SlashCommands) createscjob(ctx *dgc.Ctx) {
	name := utils.Arguments["<job name>"]

	var found bool
	for _, job := range slashcommands.Jobs {
		if job.Name == name {
			found = true
			break
		}
	}

	if found {
		ctx.RespondText("A job with the name `" + name + "` already exists.")
		return
	}

	channel, ok := utils.Discord.TryGetChannelFromMention(ctx, ctx.Arguments.Get(1))
	if !ok {
		ctx.RespondText("Could not find channel. Make sure to mention it or use the ID.")
		return
	}

	if channel.GuildID == "" {
		ctx.RespondText("This command can only be used in servers.")
		return
	}

	guild, err := utils.Discord.GetGuild(ctx.Session, channel.GuildID)
	if err != nil {
		ctx.RespondText("Could not find server: " + err.Error())
		return
	}

	bot, ok := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(2))
	if !ok {
		ctx.RespondText("Could not find bot. Make sure to mention it or use the ID.")
		return
	}

	interval, err := utils.Misc.ParseDuration(utils.Arguments["<interval>"])
	if err != nil {
		ctx.RespondText("An error occurred while parsing the interval. Make sure to use a valid duration (i.e. 30s, 5m, or 5h).")
		return
	}

	commands, applications, err := ctx.Session.GuildApplicationCommandIndex(channel.GuildID)
	if err != nil {
		ctx.RespondText("An error occurred while getting commands: " + err.Error())
		return
	}

	appID := bot.ID
	for _, app := range applications {
		if app.BotID == bot.ID {
			appID = app.ID
		}
	}

	commands = utils.Filter(commands, func(command *discordgo.ApplicationCommand) bool {
		return command.Name == utils.Arguments["<command>"] && command.ApplicationID == appID
	})

	if len(commands) == 0 {
		ctx.RespondText("No commands have been found matching the query.")
		return
	}

	command := commands[0]

	_, err = slashcommands.ArgumentsToCommandOptions(command.Options, ctx.Arguments, 5)
	if err != nil {
		ctx.RespondText("An error occured while parsing your arguments to command options: " + err.Error())
		return
	}

	job := &slashcommands.Job{
		Name: name,
		CommandObject: slashcommands.DiscordObject{
			Name: command.Name,
			ID:   command.ID,
		},
		ApplicationObject: slashcommands.DiscordObject{
			Name: bot.Username,
			ID:   appID,
		},
		ChannelObject: slashcommands.DiscordObject{
			Name: channel.Name,
			ID:   channel.ID,
		},
		GuildObject: slashcommands.DiscordObject{
			Name: guild.Name,
			ID:   guild.ID,
		},
		RawArguments: utils.String.RemoveWords(ctx.Arguments.Raw(), 5),
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
		AddField("Job Name", job.Name).
		Send(ctx)
}

func (s SlashCommands) deletescjob(ctx *dgc.Ctx) {
	name := utils.Arguments["<job name>"]

	job := slashcommands.GetJobByName(name)
	if job == nil {
		ctx.RespondText("Could not find a job with name `" + name + "`")
		return
	}

	err := slashcommands.RemoveJob(job)
	if err != nil {
		ctx.RespondText("An error occurred while removing the job: " + err.Error())
		return
	}

	embed.NewEmbed().
		SetDescription("Successfully deleted slash command job!").
		AddField("Job Name", job.Name).
		Send(ctx)
}

func (s SlashCommands) setscjobcommand(ctx *dgc.Ctx) {
	name := utils.Arguments["<job name>"]

	job := slashcommands.GetJobByName(name)
	if job == nil {
		ctx.RespondText("Could not find a job with name `" + name + "`")
		return
	}

	bot, ok := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(1))
	if !ok {
		ctx.RespondText("Could not find bot. Make sure to mention it or use the ID.")
		return
	}

	commands, applications, err := ctx.Session.GuildApplicationCommandIndex(job.GuildObject.ID)
	if err != nil {
		ctx.RespondText("An error occurred while getting commands: " + err.Error())
		return
	}

	appID := bot.ID
	for _, app := range applications {
		if app.BotID == bot.ID {
			appID = app.ID
		}
	}

	commands = utils.Filter(commands, func(command *discordgo.ApplicationCommand) bool {
		return command.Name == utils.Arguments["<command>"] && command.ApplicationID == appID
	})

	if len(commands) == 0 {
		ctx.RespondText("No commands have been found matching the query.")
		return
	}

	command := commands[0]

	_, err = slashcommands.ArgumentsToCommandOptions(command.Options, ctx.Arguments, 4)
	if err != nil {
		ctx.RespondText("An error occured while parsing your arguments to command options: " + err.Error())
		return
	}

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
		AddField("Job Name", job.Name).
		AddField("Command", fmt.Sprintf("`/%s%s` (%s)", command.Name, job.RawArguments, command.ID)).
		Send(ctx)
}

func (s SlashCommands) setscjobchannel(ctx *dgc.Ctx) {
	name := utils.Arguments["<job name>"]

	job := slashcommands.GetJobByName(name)
	if job == nil {
		ctx.RespondText("Could not find a job with name `" + name + "`")
		return
	}

	channel, ok := utils.Discord.TryGetChannelFromMention(ctx, ctx.Arguments.Get(1))
	if !ok {
		ctx.RespondText("Could not find channel. Make sure to mention it or use the ID.")
		return
	}

	if channel.GuildID == "" {
		ctx.RespondText("This command can only be used in servers.")
		return
	}

	guild, err := utils.Discord.GetGuild(ctx.Session, channel.GuildID)
	if err != nil {
		ctx.RespondText("Could not find server: " + err.Error())
		return
	}

	job.ChannelObject = slashcommands.DiscordObject{
		Name: channel.Name,
		ID:   channel.ID,
	}

	job.GuildObject = slashcommands.DiscordObject{
		Name: guild.Name,
		ID:   guild.ID,
	}

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
		AddField("Job Name", job.Name).
		AddField("Channel", channel.Mention()).
		Send(ctx)
}

func (s SlashCommands) setscjobinterval(ctx *dgc.Ctx) {
	name := utils.Arguments["<job name>"]

	job := slashcommands.GetJobByName(name)
	if job == nil {
		ctx.RespondText("Could not find a job with name `" + name + "`")
		return
	}

	interval, err := utils.Misc.ParseDuration(utils.Arguments["<interval>"])
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
		AddField("Job Name", job.Name).
		AddField("Interval", utils.Misc.FormatDuration(job.Interval)).
		Send(ctx)
}

func (s SlashCommands) setscjobpaused(ctx *dgc.Ctx) {
	name := utils.Arguments["<job name>"]

	job := slashcommands.GetJobByName(name)
	if job == nil {
		ctx.RespondText("Could not find a job with name `" + name + "`")
		return
	}

	toggle := strings.ToLower(utils.Arguments["<on/off>"])
	if toggle != "on" && toggle != "off" {
		_ = ctx.RespondText("Invalid value. Use `on` or `off`.")
		return
	}

	job.Paused = toggle == "on"

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
		AddField("Job Name", job.Name).
		AddField("Paused", toggle).
		Send(ctx)
}
