package slashcommands

import (
	"fmt"
	"time"

	"github.com/LightningDev1/LB-Selfbot/selfbot/utils"
	"github.com/LightningDev1/dgc"
	"github.com/LightningDev1/discordgo"
)

// commandCache is a cache for application commands.
var commandCache = make(map[string]*discordgo.ApplicationCommand)

// DiscordObject represents a Discord object with a name and ID.
type DiscordObject struct {
	Name string `json:"Name"`
	ID   string `json:"ID"`
}

// Job represents a job to execute a slash command at a certain interval.
// The jobs are stored in a JSON file and are loaded on startup.
// The jobs are then executed by the job handler.
type Job struct {
	// The name of the job
	Name string `json:"Name"`

	// The Discord object of the application command
	CommandObject DiscordObject `json:"Command"`

	// The Discord object of the application
	ApplicationObject DiscordObject `json:"Application"`

	// The Discord object of the channel
	ChannelObject DiscordObject `json:"Channel"`

	// The Discord object of the guild
	GuildObject DiscordObject `json:"Guild"`

	// The raw arguments to execute the command with
	RawArguments string `json:"RawArguments"`

	// The last execution time of the job
	LastExecution time.Time `json:"LastExecution"`

	// The interval to execute the command at
	IntervalSec int `json:"Interval"`

	// The paused state of the job
	Paused bool `json:"Paused"`

	// Runtime fields

	// The fetched application command
	Command *discordgo.ApplicationCommand `json:"-"`

	// The command options
	Options []*discordgo.ApplicationCommandOption `json:"-"`

	// The arguments to execute the command with
	Arguments *dgc.Arguments `json:"-"`

	// The interval to execute the command at
	Interval time.Duration `json:"-"`
}

// NextExecution returns the next execution time of the job.
func (j *Job) NextExecution() time.Time {
	if j.LastExecution.IsZero() {
		return time.Now()
	}

	return j.LastExecution.Add(j.Interval)
}

// GetCommand gets the application command of the job.
func (j *Job) GetCommand(session *discordgo.Session) (*discordgo.ApplicationCommand, error) {
	cmdID := CommandIdentifier(j)
	if command, ok := commandCache[cmdID]; ok {
		return command, nil
	}

	commands, _, err := session.GuildApplicationCommandIndex(j.GuildObject.ID)
	if err != nil {
		return nil, err
	}

	commands = utils.Filter(commands, func(command *discordgo.ApplicationCommand) bool {
		return command.Name == j.CommandObject.Name && command.ApplicationID == j.ApplicationObject.ID
	})

	if len(commands) == 0 {
		return nil, fmt.Errorf("no commands found")
	}

	commandCache[cmdID] = commands[0]

	return commands[0], nil
}

// Load loads the job from the configuration.
func (j *Job) Load(session *discordgo.Session) (err error) {
	j.Arguments = dgc.ParseArguments(j.RawArguments)

	if j.Command == nil {
		j.Command, err = j.GetCommand(session)
		if err != nil {
			return err
		}
	}

	j.Options, err = ArgumentsToCommandOptions(j.Command.Options, j.Arguments, 0)
	if err != nil {
		return err
	}

	j.Interval = time.Duration(j.IntervalSec) * time.Second

	return nil
}

// Interact executes the job's slash command.
func (j *Job) Interact(session *discordgo.Session) error {
	if j.Command == nil {
		return fmt.Errorf("command is nil")
	}

	return session.Interact(&discordgo.InteractData{
		Type: int(discordgo.InteractionApplicationCommand),
		Data: map[string]any{
			"version":             j.Command.Version,
			"id":                  j.Command.ID,
			"name":                j.Command.Name,
			"type":                j.Command.Type,
			"options":             j.Options,
			"application_command": j.Command,
			"attachments":         []any{},
		},
		ChannelID:     j.ChannelObject.ID,
		GuildID:       j.GuildObject.ID,
		ApplicationID: j.Command.ApplicationID,
	})
}
