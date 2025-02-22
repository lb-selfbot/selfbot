package slashcommands

import (
	"fmt"
	"time"

	"github.com/LightningDev1/LB-Selfbot/internal/router"
	"github.com/LightningDev1/LB-Selfbot/internal/services"
	"github.com/LightningDev1/LB-Selfbot/internal/services/interfaces"
	"github.com/LightningDev1/LB-Selfbot/internal/services/scjobs"
	"github.com/LightningDev1/LB-Selfbot/internal/utils"
	"github.com/LightningDev1/discordgo"
)

type SlashCommands struct {
	services *services.Services
}

func NewSlashCommands(services *services.Services) *SlashCommands {
	return &SlashCommands{services: services}
}

func (c *SlashCommands) Register() {
	cmdRouter := c.services.Discord.GetRouter()

	cmdRouter.RegisterCategory(&router.Category{
		Name:        "SlashCommands",
		Description: "Slash command utilities",
	}, []*router.Command{
		{
			Name:        "searchslash",
			Description: "Search a slash command",
			Arguments:   "<query> [server id]",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.searchslash,
		},
		{
			Name:        "executeslash",
			Description: "Run a slash command",
			Arguments:   "<#channel> <@bot> <command> [options...]",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.executeslash,
		},
		{
			Name:        "slashcommandjobhelp",
			Description: "More information about slash command jobs",
			Arguments:   "[p]slashcommandjobhelp",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.slashcommandjobhelp,
		},
		{
			Name:        "listscjobs",
			Description: "List all slash command jobs",
			Arguments:   "[p]listscjobs",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.listscjobs,
		},
		{
			Name:        "getscjob",
			Description: "Get info about a slash command job",
			Arguments:   "<job name>",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.getscjob,
		},
		{
			Name:        "createscjob",
			Description: "Create a slash command job",
			Arguments:   "<job name> <#channel> <@bot> <interval> <command> [options...]",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.createscjob,
		},
		{
			Name:        "deletescjob",
			Description: "Delete a slash command job",
			Arguments:   "<job name>",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.deletescjob,
		},
		{
			Name:        "setscjobcommand",
			Description: "Set a job's command",
			Arguments:   "<job name> <@bot> <command> [options...]",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.setscjobcommand,
		},
		{
			Name:        "setscjobchannel",
			Description: "Set a job's channel",
			Arguments:   "<job name> <#channel>",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.setscjobchannel,
		},
		{
			Name:        "setscjobinterval",
			Description: "Set a job's interval",
			Arguments:   "<job name> <interval>",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.setscjobinterval,
		},
		{
			Name:        "setscjobpaused",
			Description: "Pause or unpause a job",
			Arguments:   "<job name> <on/off>",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.setscjobpaused,
		},
	})
}

func (c *SlashCommands) searchslash(ctx *router.Ctx) {
	query := ctx.Arguments.Get(0).Raw()

	guild, err := ctx.Arguments.Get(1).Guild(router.ArgumentFlagFallbackDefault)
	if err != nil {
		ctx.RespondText("Please provide a valid server ID or use the command in a server.")
		return
	}

	commands, applications, err := ctx.Session.GuildApplicationCommandIndex(guild.ID)
	if err != nil {
		ctx.RespondText("An error occurred while getting commands: " + err.Error())
		return
	}

	commands = utils.Filter(commands, func(command *discordgo.ApplicationCommand) bool {
		return command.Name == query || utils.Misc.GetDistance(command.Name, query) < 2
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

	c.services.Embed.NewEmbed().
		SetTitle("Slash Commands: " + query).
		SetDescription(desc).
		Send(ctx)
}

func (c *SlashCommands) executeslash(ctx *router.Ctx) {
	channel, err := ctx.Arguments.Get(0).Channel()
	if err != nil {
		ctx.RespondText("Could not find channel. Make sure to mention it or use the ID.")
		return
	}

	if channel.GuildID == "" {
		ctx.RespondText("This command can only be used in servers.")
		return
	}

	bot, err := ctx.Arguments.Get(1).User()
	if err != nil {
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
		return command.Name == ctx.Arguments.Get(2).Raw() && command.ApplicationID == appID
	})

	if len(commands) == 0 {
		ctx.RespondText("No commands have been found matching the query.")
		return
	}

	cmd := commands[0]

	opts, err := scjobs.ArgumentsToCommandOptions(cmd.Options, ctx.Arguments, 3)
	if err != nil {
		ctx.RespondText("An error occurred while parsing your arguments to command options: " + err.Error())
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

func (c *SlashCommands) slashcommandjobhelp(ctx *router.Ctx) {
	c.services.Embed.NewEmbed().
		SetDescription(
			"Slash command jobs are jobs that run slash commands on a set interval. " +
				"To get more information and examples, read the documentation at " +
				"https://docs.lightning-bot.com/commands/slash-commands",
		).
		Send(ctx)
}

func (c *SlashCommands) listscjobs(ctx *router.Ctx) {
	jobs := c.services.SCJobs.GetJobs()

	suffix := map[bool]string{
		true:  "",
		false: "s",
	}[len(jobs) == 1]

	emb := c.services.Embed.NewEmbed().
		SetTitle("Slash Command Jobs").
		SetDescription("Found %d job%s", len(jobs), suffix)

	for _, job := range jobs {
		isRunning := job.IsRunning()

		cmd := job.CommandObject
		channel := job.ChannelObject
		guild := job.GuildObject

		args := job.RawArguments
		if len(args) > 0 {
			args = " " + args
		}

		paused := utils.Misc.FormatBoolYesNo(job.Paused)
		running := utils.Misc.FormatBoolYesNo(isRunning)

		text := fmt.Sprintf("Command: `/%s%s` (%s)\n", cmd.Name, args, cmd.ID)
		text += fmt.Sprintf("Channel: #%s (%s)\n", channel.Name, channel.ID)
		text += fmt.Sprintf("Server: %s (%s)\n", guild.Name, guild.ID)
		text += fmt.Sprintf("Interval: %s\n", utils.Misc.FormatDuration(job.Interval))
		text += fmt.Sprintf("Paused: %s | Is Running: %s\n", paused, running)

		emb.AddField("Name: "+job.Name, text)
	}

	emb.Send(ctx)
}

func (c *SlashCommands) getscjob(ctx *router.Ctx) {
	name := ctx.Arguments.Get(0).Raw()

	job := c.services.SCJobs.GetJobByName(name)
	if job == nil {
		ctx.RespondText("Could not find a job with name `" + name + "`")
		return
	}

	isRunning := job.IsRunning()
	timeUntilNext := time.Until(job.GetNextExecution())

	cmd := job.CommandObject
	channel := job.ChannelObject
	guild := job.GuildObject

	args := job.RawArguments
	if len(args) > 0 {
		args = " " + args
	}

	c.services.Embed.NewEmbed().
		SetTitle("Slash command job information").
		AddField("Name", job.Name).
		AddField("Command", fmt.Sprintf("`/%s%s` (%s)", cmd.Name, args, cmd.ID)).
		AddField("Channel", fmt.Sprintf("#%s (%s)", channel.Name, channel.ID)).
		AddField("Server", fmt.Sprintf("%s (%s)", guild.Name, guild.ID)).
		AddField("Interval", utils.Misc.FormatDuration(job.Interval)).
		AddField("Paused", utils.Misc.FormatBoolYesNo(job.Paused)).
		AddField("Job handler is running", utils.Misc.FormatBoolYesNo(isRunning)).
		AddField("Time until next interaction", utils.Misc.FormatDuration(timeUntilNext)).
		Send(ctx)
}

func (c *SlashCommands) createscjob(ctx *router.Ctx) {
	name := ctx.Arguments.Get(0).Raw()

	for _, job := range c.services.SCJobs.GetJobs() {
		if job.Name != name {
			continue
		}

		ctx.RespondText("A job with the name `" + name + "` already exists.")
		return
	}

	channel, err := ctx.Arguments.Get(1).Channel()
	if err != nil {
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

	bot, err := ctx.Arguments.Get(2).User(router.ArgumentFlagFallbackDefault)
	if err != nil {
		ctx.RespondText("Could not find bot. Make sure to mention it or use the ID.")
		return
	}

	interval, err := ctx.Arguments.Get(3).Duration()
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
		return command.Name == ctx.Arguments.Get(4).Raw() && command.ApplicationID == appID
	})

	if len(commands) == 0 {
		ctx.RespondText("No commands have been found matching the query.")
		return
	}

	command := commands[0]

	_, err = scjobs.ArgumentsToCommandOptions(command.Options, ctx.Arguments, 5)
	if err != nil {
		ctx.RespondText("An error occurred while parsing your arguments to command options: " + err.Error())
		return
	}

	job := &interfaces.SCJob{
		Name: name,
		CommandObject: interfaces.DiscordObject{
			Name: command.Name,
			ID:   command.ID,
		},
		ApplicationObject: interfaces.DiscordObject{
			Name: bot.Username,
			ID:   appID,
		},
		ChannelObject: interfaces.DiscordObject{
			Name: channel.Name,
			ID:   channel.ID,
		},
		GuildObject: interfaces.DiscordObject{
			Name: guild.Name,
			ID:   guild.ID,
		},
		RawArguments: ctx.Arguments.RawFrom(5),
		IntervalSec:  int(interval.Seconds()),

		// Set command so the job handler won't have to fetch it again
		Command: command,
	}

	err = c.services.SCJobs.AddJob(job)
	if err != nil {
		ctx.RespondText("An error occurred while adding the job: " + err.Error())
		return
	}

	c.services.Embed.NewEmbed().
		SetDescription("Successfully created slash command job!").
		AddField("Job Name", job.Name).
		Send(ctx)
}

func (c *SlashCommands) deletescjob(ctx *router.Ctx) {
	name := ctx.Arguments.Get(0).Raw()

	job := c.services.SCJobs.GetJobByName(name)
	if job == nil {
		ctx.RespondText("Could not find a job with name `" + name + "`")
		return
	}

	err := c.services.SCJobs.RemoveJob(job)
	if err != nil {
		ctx.RespondText("An error occurred while removing the job: " + err.Error())
		return
	}

	c.services.Embed.NewEmbed().
		SetDescription("Successfully deleted slash command job!").
		AddField("Job Name", job.Name).
		Send(ctx)
}

func (c *SlashCommands) setscjobcommand(ctx *router.Ctx) {
	name := ctx.Arguments.Get(0).Raw()

	job := c.services.SCJobs.GetJobByName(name)
	if job == nil {
		ctx.RespondText("Could not find a job with name `" + name + "`")
		return
	}

	bot, err := ctx.Arguments.Get(1).User(router.ArgumentFlagFallbackDefault)
	if err != nil {
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
		return command.Name == ctx.Arguments.Get(2).Raw() && command.ApplicationID == appID
	})

	if len(commands) == 0 {
		ctx.RespondText("No commands have been found matching the query.")
		return
	}

	command := commands[0]

	_, err = scjobs.ArgumentsToCommandOptions(command.Options, ctx.Arguments, 4)
	if err != nil {
		ctx.RespondText("An error occurred while parsing your arguments to command options: " + err.Error())
		return
	}

	job.CommandObject = interfaces.DiscordObject{
		Name: command.Name,
		ID:   command.ID,
	}

	job.ApplicationObject = interfaces.DiscordObject{
		Name: bot.Username,
		ID:   command.ApplicationID,
	}

	job.RawArguments = ctx.Arguments.RawFrom(3)

	job.Command = command

	err = c.services.SCJobs.SaveJobs()
	if err != nil {
		ctx.RespondText("An error occurred while saving the jobs: " + err.Error())
		return
	}

	c.services.Embed.NewEmbed().
		SetDescription("Successfully updated slash command job").
		AddField("Job Name", job.Name).
		AddField("Command", fmt.Sprintf("`/%s%s` (%s)", command.Name, job.RawArguments, command.ID)).
		Send(ctx)
}

func (c *SlashCommands) setscjobchannel(ctx *router.Ctx) {
	name := ctx.Arguments.Get(0).Raw()

	job := c.services.SCJobs.GetJobByName(name)
	if job == nil {
		ctx.RespondText("Could not find a job with name `" + name + "`")
		return
	}

	channel, err := ctx.Arguments.Get(1).Channel()
	if err != nil {
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

	job.ChannelObject = interfaces.DiscordObject{
		Name: channel.Name,
		ID:   channel.ID,
	}

	job.GuildObject = interfaces.DiscordObject{
		Name: guild.Name,
		ID:   guild.ID,
	}

	err = c.services.SCJobs.SaveJobs()
	if err != nil {
		ctx.RespondText("An error occurred while saving the jobs: " + err.Error())
		return
	}

	c.services.Embed.NewEmbed().
		SetDescription("Successfully updated slash command job").
		AddField("Job Name", job.Name).
		AddField("Channel", channel.Mention()).
		Send(ctx)
}

func (c *SlashCommands) setscjobinterval(ctx *router.Ctx) {
	name := ctx.Arguments.Get(0).Raw()

	job := c.services.SCJobs.GetJobByName(name)
	if job == nil {
		ctx.RespondText("Could not find a job with name `" + name + "`")
		return
	}

	interval, err := ctx.Arguments.Get(1).Duration()
	if err != nil {
		ctx.RespondText("An error occurred while parsing the interval. Make sure to use a valid duration (i.e. 30s, 5m, or 5h).")
		return
	}

	job.IntervalSec = int(interval.Seconds())

	err = c.services.SCJobs.SaveJobs()
	if err != nil {
		ctx.RespondText("An error occurred while saving the jobs: " + err.Error())
		return
	}

	c.services.Embed.NewEmbed().
		SetDescription("Successfully updated slash command job").
		AddField("Job Name", job.Name).
		AddField("Interval", utils.Misc.FormatDuration(job.Interval)).
		Send(ctx)
}

func (c *SlashCommands) setscjobpaused(ctx *router.Ctx) {
	name := ctx.Arguments.Get(0).Raw()

	job := c.services.SCJobs.GetJobByName(name)
	if job == nil {
		ctx.RespondText("Could not find a job with name `" + name + "`")
		return
	}

	var err error
	job.Paused, err = ctx.Arguments.Get(1).Bool()
	if err != nil {
		ctx.RespondText("Invalid value. Use `on` or `off`.")
		return
	}

	err = c.services.SCJobs.SaveJobs()
	if err != nil {
		ctx.RespondText("An error occurred while saving the jobs: " + err.Error())
		return
	}

	c.services.Embed.NewEmbed().
		SetDescription("Successfully updated slash command job").
		AddField("Job Name", job.Name).
		AddField("Paused", utils.Misc.FormatBoolOnOff(job.Paused)).
		Send(ctx)
}
