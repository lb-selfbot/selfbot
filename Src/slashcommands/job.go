package slashcommands

import (
	"crypto/md5"
	"fmt"
	"time"

	"github.com/LightningDev1/LB-Selfbot/utils"
	"github.com/LightningDev1/dgc"
	"github.com/LightningDev1/discordgo"
)

var commandCache = make(map[string]*discordgo.ApplicationCommand)

type DiscordObject struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

// Job represents a job to execute a slash command at a certain interval.
// The jobs are stored in a JSON file and are loaded on startup.
// The jobs are then executed by the job handler.
type Job struct {
	// The name and ID of the application command
	CommandObject DiscordObject `json:"command"`

	// The name and ID of the application
	ApplicationObject DiscordObject `json:"application"`

	// The name and ID of the channel
	ChannelObject DiscordObject `json:"channel"`

	// The name and ID of the guild
	GuildObject DiscordObject `json:"guild"`

	// The raw arguments to execute the command with
	RawArguments string `json:"raw_arguments"`

	// The interval to execute the command at
	IntervalSec int `json:"interval"`

	//////////////////////////////////////////////////////////////////
	// The following fields are initialized when the job is started //
	//////////////////////////////////////////////////////////////////

	// The ID of the job (calculated from the MD5 hash)
	ID string `json:"-"`

	// The application command to execute
	Command *discordgo.ApplicationCommand `json:"-"`

	// The arguments to execute the command with
	Arguments *dgc.Arguments `json:"-"`

	// The interval to execute the command at
	Interval time.Duration `json:"-"`
}

func (j *Job) CalculateID() string {
	cID := j.CommandObject.ID
	aID := j.ApplicationObject.ID
	chID := j.ChannelObject.ID
	gID := j.GuildObject.ID
	args := j.RawArguments
	interval := fmt.Sprint(j.IntervalSec)

	byteRepresentation := []byte(cID + aID + chID + gID + args + interval)

	return fmt.Sprintf("%x", md5.Sum(byteRepresentation))
}

func (j *Job) GetCommand(session *discordgo.Session) (*discordgo.ApplicationCommand, error) {
	if command, ok := commandCache[j.ID]; ok {
		return command, nil
	}

	commands, _, err := session.ChannelApplicationCommandsSearch(&discordgo.ApplicationCommandsSearchParams{
		ChannelID:           j.ChannelObject.ID,
		Type:                1,
		Query:               j.CommandObject.Name,
		Limit:               20,
		IncludeApplications: false,
	})
	if err != nil {
		return nil, err
	}

	rightCommands := utils.Filter(commands, func(command *discordgo.ApplicationCommand) bool {
		return command.Name == j.CommandObject.Name && command.ApplicationID == j.ApplicationObject.ID
	})

	if len(rightCommands) == 0 {
		return nil, fmt.Errorf("no commands found")
	}

	commandCache[j.ID] = rightCommands[0]

	return rightCommands[0], nil
}

func (j *Job) GetArguments() *dgc.Arguments {
	return dgc.ParseArguments(j.RawArguments)
}

func (j *Job) Load(session *discordgo.Session) error {
	j.ID = j.CalculateID()
	j.Arguments = j.GetArguments()

	if j.Command == nil {
		command, err := j.GetCommand(session)
		if err != nil {
			return err
		}

		j.Command = command
	}

	options, err := ArgumentsToCommandOptions(j.Command.Options, j.Arguments, 0)
	if err != nil {
		return err
	}

	j.Command.Options = options

	j.Interval = time.Duration(j.IntervalSec) * time.Second

	return nil
}

func (j *Job) ExecuteInteraction(session *discordgo.Session) error {
	if j.Command == nil {
		return fmt.Errorf("command is nil")
	}

	return session.Interact(&discordgo.InteractData{
		Type:          int(discordgo.InteractionApplicationCommand),
		Data:          j.Command,
		ChannelID:     j.ChannelObject.ID,
		GuildID:       j.GuildObject.ID,
		ApplicationID: j.Command.ApplicationID,
	})
}
