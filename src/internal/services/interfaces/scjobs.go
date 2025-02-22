package interfaces

import (
	"crypto/md5"
	"fmt"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/LightningDev1/LB-Selfbot/internal/router"
	"github.com/LightningDev1/discordgo"
)

// DiscordObject represents a Discord object with a name and ID.
type DiscordObject struct {
	Name string `json:"Name"`
	ID   string `json:"ID"`
}

// SCJob represents a job to execute a slash command at a certain interval.
// The jobs are stored in a JSON file and are loaded on startup.
// The jobs are then executed by the job handler.
type SCJob struct {
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
	Arguments *router.Arguments `json:"-"`

	// The interval to execute the command at
	Interval time.Duration `json:"-"`

	// Lock for the job
	Lock atomic.Uint32 `json:"-"`

	// Channel to stop the job
	StopChannel chan struct{} `json:"-"`
}

// IsRunning returns true if the job is running.
func (j *SCJob) IsRunning() bool {
	return j.Lock.Load() == 1
}

// Stop stops the job handler.
func (j *SCJob) Stop() {
	close(j.StopChannel)
}

// GetNextExecution returns the next execution time of the job.
func (j *SCJob) GetNextExecution() time.Time {
	if j.LastExecution.IsZero() {
		return time.Now()
	}

	return j.LastExecution.Add(j.Interval)
}

func (j *SCJob) GetCommandID() string {
	return j.GuildObject.ID + j.ApplicationObject.ID + j.CommandObject.ID
}

func (j *SCJob) GetJobID() string {
	id := []byte(j.GuildObject.ID +
		j.ApplicationObject.ID +
		j.ChannelObject.ID +
		j.CommandObject.ID +
		j.Name +
		j.RawArguments +
		strconv.Itoa(j.IntervalSec))

	return fmt.Sprintf("%x", md5.Sum(id))
}

type SCJobsService interface {
	Inject(uiSvc UIService, discordSvc DiscordService)

	Start()
	LoadJobs() error
	SaveJobs() error
	GetJobs() []*SCJob
	GetJobByName(name string) *SCJob
	AddJob(job *SCJob) error
	RemoveJob(job *SCJob) error
}
