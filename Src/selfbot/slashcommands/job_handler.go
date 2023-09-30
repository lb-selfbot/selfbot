package slashcommands

import (
	"fmt"
	"sync/atomic"
	"time"

	"github.com/LightningDev1/LB-Selfbot/selfbot/logging"
	"github.com/LightningDev1/LB-Selfbot/selfbot/utils"
	"github.com/LightningDev1/discordgo"
)

var jobHandlers = make(map[string]*JobHandler)
var errAlreadyRunning = fmt.Errorf("job is already running")

type JobHandler struct {
	Job             *Job
	Locker          atomic.Uint32
	NextInteraction time.Time
	Stop            chan bool
}

func (j *JobHandler) ExecuteJobHandler(session *discordgo.Session) error {
	if !j.Locker.CompareAndSwap(0, 1) {
		return errAlreadyRunning
	}
	defer j.Locker.Store(0)

	for {
		j.NextInteraction = time.Now().Add(j.Job.Interval)

		select {
		case <-j.Stop:
			return nil
		case <-time.After(j.Job.Interval):
			err := j.Job.ExecuteInteraction(session)
			if err != nil {
				return err
			}
		}
	}
}

func (j *JobHandler) TimeUntilNextInteraction() time.Duration {
	return time.Until(j.NextInteraction)
}

func (j *JobHandler) IsJobRunning() bool {
	// Check if the job is locked
	return j.Locker.Load() == 1
}

func (j *JobHandler) StopJobHandler() {
	j.Stop <- true
}

func GetHandlerForJob(job *Job) *JobHandler {
	if handler, ok := jobHandlers[job.ID]; ok {
		return handler
	}

	handler := &JobHandler{
		Job:    job,
		Locker: atomic.Uint32{},
		Stop:   make(chan bool),
	}

	jobHandlers[job.ID] = handler

	return handler
}

func StartJobs(session *discordgo.Session) {
	hasLogged := false

	for {
		err := LoadJobs(session)
		if err != nil {
			utils.Logging.Error("Failed to load jobs: " + err.Error())
			time.Sleep(time.Second * 5)
			continue
		}

		// Start the job handlers
		for _, job := range Jobs {
			handler := GetHandlerForJob(job)

			if handler.IsJobRunning() {
				continue
			}

			go func() {
				ret := handler.ExecuteJobHandler(session)
				if ret != nil {
					utils.Logging.Error("Error executing job handler: " + ret.Error())
				}
			}()
		}

		// Remove handlers for jobs that no longer exist
		for _, handler := range jobHandlers {
			jobs := utils.Filter(Jobs, func(job *Job) bool {
				return job.ID == handler.Job.ID
			})

			jobExists := len(jobs) > 0

			if !jobExists {
				handler.StopJobHandler()

				delete(jobHandlers, handler.Job.ID)
			}
		}

		if !hasLogged {
			if len(Jobs) > 0 {
				logging.Logging.SlashCommandJobs(len(Jobs))
			}

			hasLogged = true
		}

		time.Sleep(time.Second * 5)
	}
}
