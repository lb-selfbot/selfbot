package slashcommand

import (
	"errors"
	"fmt"
	"slices"
	"sync"
	"time"

	"github.com/goccy/go-json"

	"github.com/lb-selfbot/discordgo"
	"github.com/lb-selfbot/selfbot/internal/fn"
	"github.com/lb-selfbot/selfbot/internal/logging"
	"github.com/lb-selfbot/selfbot/internal/services/config"
	"github.com/lb-selfbot/selfbot/internal/services/eventbroker"
	"github.com/lb-selfbot/selfbot/internal/services/router"
	"github.com/lb-selfbot/selfbot/internal/services/ui/featurelogs"
)

var (
	// ErrInvalidSubCommand is an error for an invalid subcommand.
	ErrInvalidSubCommand = errors.New("invalid subcommand")

	// ErrCommandNotFound is an error for a command not found.
	ErrCommandNotFound = errors.New("command not found")
)

// Service handles slash command interactions and jobs.
type Service interface {
	// GetJobs returns all the jobs.
	GetJobs() []*SCJob

	// SaveJobs saves the jobs to the config.
	SaveJobs() error

	// GetJobByName returns a job by its name.
	GetJobByName(name string) *SCJob

	// AddJob adds a new job.
	AddJob(job *SCJob) error

	// RemoveJob removes a job.
	RemoveJob(job *SCJob) error
}

type slashCommandService struct {
	eventBroker eventbroker.Service
	configSvc   config.Service

	jobs   []*SCJob
	jobsMu sync.Mutex

	cmdCache map[string]*discordgo.ApplicationCommand
	cmdMu    sync.Mutex
}

// NewService creates a new slash command service.
func NewService(eventBroker eventbroker.Service, configSvc config.Service) Service {
	svc := &slashCommandService{eventBroker: eventBroker, configSvc: configSvc}
	svc.subscribe()
	return svc
}

func (s *slashCommandService) subscribe() {
	s.eventBroker.Subscribe(eventbroker.EventSessionReady, func(payload any) {
		session, ok := payload.(*discordgo.Session)
		if !ok {
			logging.Error("Failed to cast payload to *discordgo.Session")
			return
		}

		go s.start(session)
	})
}

func (s *slashCommandService) start(session *discordgo.Session) {
	isFirstIteration := true

	for range time.Tick(time.Second * 5) {
		err := s.LoadJobs(session)
		if err != nil {
			s.eventBroker.Publish(eventbroker.EventLogError, "Failed to load SC jobs: "+err.Error())
			continue
		}

		// First start all job handlers
		for _, job := range s.jobs {
			if job.IsRunning() {
				continue
			}

			go s.runJob(session, job)
		}

		if isFirstIteration {
			s.eventBroker.Publish(eventbroker.EventLogFeature, &featurelogs.SCJobs{
				NumJobs: len(s.jobs),
			})

			isFirstIteration = false
		}
	}
}

func (s *slashCommandService) runJob(session *discordgo.Session, job *SCJob) {
	// Check if the job is already running
	if !job.Lock.CompareAndSwap(0, 1) {
		return
	}
	defer job.Lock.Store(0)

	for slices.Contains(s.jobs, job) {
		// If the job is paused, wait until it is unpaused
		if job.Paused {
			time.Sleep(time.Second * 5)
			continue
		}

		time.Sleep(time.Until(job.GetNextExecution()))

		// Check if the job was removed while waiting
		if !slices.Contains(s.jobs, job) {
			return
		}

		// Check if the job was paused while waiting
		if job.Paused {
			continue
		}

		err := s.interactJob(session, job)
		if err != nil {
			logging.Error("Failed to interact with job:", err.Error())
			continue
		}

		job.LastExecution = time.Now()
		err = s.SaveJobs()
		if err != nil {
			logging.Error("Failed to save jobs:", err.Error())
			continue
		}
	}
}

func (s *slashCommandService) interactJob(session *discordgo.Session, j *SCJob) error {
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

func (s *slashCommandService) LoadJobs(session *discordgo.Session) error {
	s.jobsMu.Lock()
	defer s.jobsMu.Unlock()

	cfg := s.configSvc.GetConfig()

	err := json.Unmarshal(cfg.SlashCommandJobs, &s.jobs)
	if err != nil {
		logging.Error("Failed to unmarshal jobs:", err.Error())

		// Reset the jobs when the unmarshal has failed
		cfg.SlashCommandJobs = []byte("[]")
		_ = s.configSvc.SaveConfig(cfg)

		return err
	}

	s.jobs = fn.Filter(s.jobs, func(job *SCJob) bool {
		job.Arguments = router.ParseArguments(&router.Ctx{Session: session}, job.RawArguments)

		if job.Command == nil {
			job.Command, err = s.getCommand(session, job)
			if err != nil {
				return false
			}
		}

		job.Options, err = ArgumentsToCommandOptions(job.Command.Options, job.Arguments, 0)
		if err != nil {
			return false
		}

		job.Interval = time.Duration(job.IntervalSec) * time.Second

		return true
	})

	return nil
}

func (s *slashCommandService) getCommand(session *discordgo.Session, job *SCJob) (*discordgo.ApplicationCommand, error) {
	s.cmdMu.Lock()
	defer s.cmdMu.Unlock()

	cmdID := job.GetCommandID()

	if cmd, ok := s.cmdCache[cmdID]; ok {
		return cmd, nil
	}

	commands, _, err := session.GuildApplicationCommandIndex(job.GuildObject.ID)
	if err != nil {
		return nil, err
	}

	commands = fn.Filter(commands, func(command *discordgo.ApplicationCommand) bool {
		return command.Name == job.Command.Name &&
			command.ApplicationID == job.ApplicationObject.ID
	})

	if len(commands) == 0 {
		return nil, ErrCommandNotFound
	}

	cmd := commands[0]
	s.cmdCache[cmdID] = cmd
	return cmd, nil
}

func (s *slashCommandService) GetJobs() []*SCJob {
	return s.jobs
}

func (s *slashCommandService) SaveJobs() error {
	cfg := s.configSvc.GetConfig()

	jobs, err := json.MarshalIndent(s.jobs, "", "  ")
	if err != nil {
		return err
	}

	cfg.SlashCommandJobs = jobs

	return s.configSvc.SaveConfig(cfg)
}

func (s *slashCommandService) GetJobByName(name string) *SCJob {
	for _, job := range s.jobs {
		if job.Name == name {
			return job
		}
	}

	return nil
}

func (s *slashCommandService) AddJob(job *SCJob) error {
	s.jobsMu.Lock()
	s.jobs = append(s.jobs, job)
	s.jobsMu.Unlock()

	return s.SaveJobs()
}

func (s *slashCommandService) RemoveJob(job *SCJob) error {
	s.jobsMu.Lock()
	s.jobs = fn.Filter(s.jobs, func(j *SCJob) bool {
		return j.Name != job.Name
	})
	s.jobsMu.Unlock()

	return s.SaveJobs()
}
