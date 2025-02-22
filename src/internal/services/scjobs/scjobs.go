package scjobs

import (
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"sync"
	"time"

	"github.com/LightningDev1/LB-Selfbot/internal/config"
	"github.com/LightningDev1/LB-Selfbot/internal/router"
	"github.com/LightningDev1/LB-Selfbot/internal/services/interfaces"
	"github.com/LightningDev1/LB-Selfbot/internal/utils"
	"github.com/LightningDev1/discordgo"
)

var (
	// ErrInvalidSubCommand is an error for an invalid subcommand.
	ErrInvalidSubCommand = errors.New("invalid subcommand")

	// ErrCommandNotFound is an error for a command not found.
	ErrCommandNotFound = errors.New("command not found")
)

type SCJobsServiceImpl struct {
	uiSvc      interfaces.UIService
	discordSvc interfaces.DiscordService

	jobs   []*interfaces.SCJob
	jobsMu sync.Mutex

	cmdCache map[string]*discordgo.ApplicationCommand
	cmdMu    sync.Mutex
}

func (s *SCJobsServiceImpl) Inject(uiSvc interfaces.UIService, discordSvc interfaces.DiscordService) {
	s.uiSvc = uiSvc
	s.discordSvc = discordSvc
}

func (s *SCJobsServiceImpl) Start() {
	go s.start()
}

func (s *SCJobsServiceImpl) start() {
	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()

	isFirstIteration := true
	for range ticker.C {
		err := s.LoadJobs()
		if err != nil {
			s.uiSvc.LogError("Failed to load SC jobs:", err.Error())
			continue
		}

		// First start all job handlers
		for _, job := range s.jobs {
			if job.IsRunning() {
				continue
			}

			go s.runJob(job)
		}

		if isFirstIteration {
			s.uiSvc.LogSCJobs(len(s.jobs))
			isFirstIteration = false
		}
	}
}

func (s *SCJobsServiceImpl) runJob(job *interfaces.SCJob) {
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

		err := s.interactJob(job)
		if err != nil {
			utils.Logging.Error("Failed to interact with job:", err.Error())
			continue
		}

		job.LastExecution = time.Now()
		err = s.SaveJobs()
		if err != nil {
			utils.Logging.Error("Failed to save jobs:", err.Error())
			continue
		}
	}
}

func (s *SCJobsServiceImpl) interactJob(j *interfaces.SCJob) error {
	if j.Command == nil {
		return fmt.Errorf("command is nil")
	}

	session := s.discordSvc.GetSession()

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

func (s *SCJobsServiceImpl) LoadJobs() error {
	s.jobsMu.Lock()
	defer s.jobsMu.Unlock()

	cfg, err := config.GetConfig()
	if err != nil {
		return err
	}

	err = json.Unmarshal(cfg.SlashCommandJobs, &s.jobs)
	if err != nil {
		utils.Logging.Error("Failed to unmarshal jobs:", err.Error())

		// Reset the jobs when the unmarshal has failed
		cfg.SlashCommandJobs = []byte("[]")
		_ = cfg.Save()

		return err
	}

	session := s.discordSvc.GetSession()

	s.jobs = utils.Filter(s.jobs, func(job *interfaces.SCJob) bool {
		job.Arguments = router.ParseArguments(&router.Ctx{Session: session}, job.RawArguments)

		if job.Command == nil {
			job.Command, err = s.getCommand(job)
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

func (s *SCJobsServiceImpl) getCommand(job *interfaces.SCJob) (*discordgo.ApplicationCommand, error) {
	s.cmdMu.Lock()
	defer s.cmdMu.Unlock()

	session := s.discordSvc.GetSession()
	cmdID := job.GetCommandID()

	if cmd, ok := s.cmdCache[cmdID]; ok {
		return cmd, nil
	}

	commands, _, err := session.GuildApplicationCommandIndex(job.GuildObject.ID)
	if err != nil {
		return nil, err
	}

	commands = utils.Filter(commands, func(command *discordgo.ApplicationCommand) bool {
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

func (s *SCJobsServiceImpl) SaveJobs() error {
	cfg, err := config.GetConfig()
	if err != nil {
		return err
	}

	jobs, err := json.MarshalIndent(s.jobs, "", "  ")
	if err != nil {
		return err
	}

	cfg.SlashCommandJobs = jobs

	return cfg.Save()
}

func (s *SCJobsServiceImpl) GetJobs() []*interfaces.SCJob {
	return s.jobs
}

func (s *SCJobsServiceImpl) GetJobByName(name string) *interfaces.SCJob {
	for _, job := range s.jobs {
		if job.Name == name {
			return job
		}
	}

	return nil
}

func (s *SCJobsServiceImpl) AddJob(job *interfaces.SCJob) error {
	s.jobsMu.Lock()
	s.jobs = append(s.jobs, job)
	s.jobsMu.Unlock()

	return s.SaveJobs()
}

func (s *SCJobsServiceImpl) RemoveJob(job *interfaces.SCJob) error {
	s.jobsMu.Lock()
	s.jobs = utils.Filter(s.jobs, func(j *interfaces.SCJob) bool {
		return j.Name != job.Name
	})
	s.jobsMu.Unlock()

	return s.SaveJobs()
}
