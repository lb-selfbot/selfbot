package slashcommands

import (
	"crypto/md5"
	"fmt"
	"time"

	"github.com/LightningDev1/dgc"
	"github.com/LightningDev1/discordgo"
)

var commandCache = make(map[string]*discordgo.ApplicationCommand)

// Job represents a job to execute a slash command at a certain interval.
// The jobs are stored in a JSON file and are loaded on startup.
// The jobs are then executed by the job handler.
type Job struct {
	// The ID of the application command
	CommandID string `json:"command_id"`

	// The ID of the channel
	ChannelID string `json:"channel_id"`

	// The ID of the guild
	GuildID string `json:"guild_id"`

	// The raw arguments to execute the command with
	RawArguments string `json:"raw_arguments"`

	// The interval to execute the command at
	Interval time.Duration `json:"interval"`

	//////////////////////////////////////////////////////////////////
	// The following fields are initialized when the job is started //
	//////////////////////////////////////////////////////////////////

	// The ID of the job (calculated from the MD5 hash)
	ID string

	// The application command to execute
	Command *discordgo.ApplicationCommand

	// The arguments to execute the command with
	Arguments *dgc.Arguments
}

func (j *Job) CalculateID() string {
	stringRepresentation := j.CommandID + j.ChannelID + j.RawArguments + j.Interval.String()

	hash := md5.Sum([]byte(stringRepresentation))

	return fmt.Sprintf("%x", hash)
}

func (j *Job) GetCommand(session *discordgo.Session) (*discordgo.ApplicationCommand, error) {
	if command, ok := commandCache[j.ID]; ok {
		return command, nil
	}

	commands, _, err := session.ChannelApplicationCommandsSearch(&discordgo.ApplicationCommandsSearchParams{
		ChannelID:  j.ChannelID,
		Type:       1,
		Limit:      1,
		CommandIDs: []string{j.CommandID},
	})
	if err != nil {
		return nil, err
	}

	if len(commands) == 0 {
		return nil, fmt.Errorf("no commands found")
	}

	commandCache[j.ID] = commands[0]

	return commands[0], nil
}

func (j *Job) GetArguments() *dgc.Arguments {
	return dgc.ParseArguments(j.RawArguments)
}

func (j *Job) Load(session *discordgo.Session) error {
	j.ID = j.CalculateID()
	j.Arguments = j.GetArguments()

	command, err := j.GetCommand(session)
	if err != nil {
		return err
	}

	j.Command = command

	return nil
}

func (j *Job) ExecuteInteraction(session *discordgo.Session) error {
	return session.Interact(&discordgo.InteractData{
		Type:          int(discordgo.InteractionApplicationCommand),
		Data:          j.Command,
		ChannelID:     j.ChannelID,
		GuildID:       j.GuildID,
		ApplicationID: j.Command.ApplicationID,
	})
}
