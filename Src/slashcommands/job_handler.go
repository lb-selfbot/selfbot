package slashcommands

import (
	"fmt"
	"sync"
	"time"

	"github.com/LightningDev1/LB-Selfbot/utils"
	"github.com/LightningDev1/discordgo"
)

var jobHandlers = make(map[string]*JobHandler)

type JobHandler struct {
	Job   *Job
	Mutex *sync.Mutex
	Stop  chan bool
}

func (j *JobHandler) ExecuteJobHandler(session *discordgo.Session) error {
	locked := j.Mutex.TryLock()
	if !locked {
		return fmt.Errorf("job is already running")
	}

	defer j.Mutex.Unlock()

	for {
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

func (j *JobHandler) StopJobHandler() {
	j.Stop <- true
}

func GetHandlerForJob(job *Job) *JobHandler {
	if handler, ok := jobHandlers[job.ID]; ok {
		return handler
	}

	handler := &JobHandler{
		Job:   job,
		Mutex: &sync.Mutex{},
		Stop:  make(chan bool),
	}

	jobHandlers[job.ID] = handler

	return handler
}

func StartJobs(session *discordgo.Session) {
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

			go handler.ExecuteJobHandler(session)
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

		time.Sleep(time.Second * 5)
	}
}
