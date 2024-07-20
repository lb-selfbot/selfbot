package slashcommands

import (
	"fmt"
	"sync/atomic"
	"time"

	"github.com/LightningDev1/LB-Selfbot/selfbot/logging"
	"github.com/LightningDev1/LB-Selfbot/selfbot/utils"
	"github.com/LightningDev1/discordgo"
)

var (
	// jobHandlers is a map of job ID to *JobHandler.
	jobHandlers = make(map[string]*JobHandler)

	// errAlreadyRunning is an error that is returned when a job is already running.
	errAlreadyRunning = fmt.Errorf("job is already running")
)

// JobHandler represents a handler for a job.
type JobHandler struct {
	// The job to handle.
	JobName string

	// The ID of the job.
	JobID string

	// Locker is a lock for the job.
	Locker atomic.Uint32

	// StopChannel is a channel to stop the job.
	StopChannel chan bool
}

// GetJobHandler gets or creates a job handler.
func GetJobHandler(jobName, jobID string) *JobHandler {
	if handler, ok := jobHandlers[jobID]; ok {
		return handler
	}

	handler := &JobHandler{
		JobName:     jobName,
		JobID:       jobID,
		Locker:      atomic.Uint32{},
		StopChannel: make(chan bool),
	}

	jobHandlers[jobID] = handler

	return handler
}

// Execute executes the job handler.
func (j *JobHandler) Execute(session *discordgo.Session) error {
	if !j.Locker.CompareAndSwap(0, 1) {
		return errAlreadyRunning
	}
	defer j.Locker.Store(0)

	for {
		job := GetJobByName(j.JobName)
		if job == nil {
			return fmt.Errorf("job not found")
		}

		// Wait until the job is unpaused.
		for job.Paused {
			select {
			case <-j.StopChannel:
				return nil
			case <-time.After(time.Second * 5):
				continue
			}
		}

		untilNext := time.Until(job.NextExecution())

		select {
		case <-j.StopChannel:
			return nil
		case <-time.After(untilNext):
			if job.Paused {
				continue
			}
			
			// Execute the job's interaction.
			err := job.Interact(session)
			if err != nil {
				return err
			}

			// Update the last execution time.
			job.LastExecution = time.Now()
			err = SaveJobs()
			if err != nil {
				return err
			}
		}
	}
}

// IsRunning checks if the job handler is running.
func (j *JobHandler) IsRunning() bool {
	return j.Locker.Load() == 1
}

// Stop stops the job handler.
func (j *JobHandler) Stop() {
	j.StopChannel <- true
}

// ExecuteJobs executes all jobs.
func ExecuteJobs(session *discordgo.Session) {
	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()

	var logged bool

	for range ticker.C {
		err := LoadJobs(session)
		if err != nil {
			utils.Logging.Error("Failed to load jobs:", err.Error())
			continue
		}

		var jobIDs []string

		// Run all job handlers.
		for _, job := range Jobs {
			jobID := JobIdentifier(job)
			jobIDs = append(jobIDs, jobID)

			handler := GetJobHandler(job.Name, jobID)

			if handler.IsRunning() {
				continue
			}

			go func() {
				err := handler.Execute(session)
				if err != nil {
					utils.Logging.Error("Failed to execute job:", err.Error())
				}
			}()
		}

		// Remove job handlers for jobs that no longer exist.
		for _, handler := range jobHandlers {
			// If the handler's job ID is in the list of job IDs, continue.
			if utils.SliceContains(jobIDs, handler.JobID) {
				continue
			}

			// Delete the handler
			handler.Stop()
			delete(jobHandlers, handler.JobID)
		}

		if !logged {
			logging.Logging.SlashCommandJobs(len(Jobs))
			logged = true
		}
	}
}
