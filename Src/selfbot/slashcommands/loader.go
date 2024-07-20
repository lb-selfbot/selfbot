package slashcommands

import (
	"encoding/json"
	"sync"

	"github.com/LightningDev1/LB-Selfbot/selfbot/config"
	"github.com/LightningDev1/LB-Selfbot/selfbot/utils"
	"github.com/LightningDev1/discordgo"
)

var (
	// jobsMutex is a mutex to lock the jobs slice.
	jobsMutex = &sync.Mutex{}

	// Jobs is a slice of jobs that are loaded.
	Jobs []*Job
)

// LoadJobs loads the jobs from the configuration.
func LoadJobs(session *discordgo.Session) error {
	jobsMutex.Lock()
	defer jobsMutex.Unlock()

	cfg, err := config.GetConfig()
	if err != nil {
		return err
	}

	// Unmarshal the jobs from the configuration.
	err = json.Unmarshal(cfg.SlashCommandJobs, &Jobs)
	if err != nil {
		utils.Logging.Error("Failed to unmarshal jobs:", err.Error())

		// Reset the jobs if the unmarshal failed.
		cfg.SlashCommandJobs = []byte("[]")
		_ = cfg.Save()

		return err
	}

	// Load each job.
	var invalidJobNames []string

	for _, job := range Jobs {
		err = job.Load(session)
		if err != nil {
			utils.Logging.Error("Failed to load job:", err.Error())

			invalidJobNames = append(invalidJobNames, job.Name)
		}
	}

	// Remove invalid jobs.
	if len(invalidJobNames) > 0 {
		Jobs = utils.Filter(Jobs, func(job *Job) bool {
			return !utils.SliceContains(invalidJobNames, job.Name)
		})

		_ = SaveJobs()
	}

	return nil
}

// SaveJobs saves the jobs to the configuration.
func SaveJobs() error {
	jobsMutex.Lock()
	defer jobsMutex.Unlock()

	cfg, err := config.GetConfig()
	if err != nil {
		return err
	}

	// Marshal the jobs to the configuration.
	jobsBytes, err := json.MarshalIndent(Jobs, "", "  ")
	if err != nil {
		return err
	}

	cfg.SlashCommandJobs = jobsBytes

	return cfg.Save()
}

// GetJobByName returns a job by its name.
func GetJobByName(name string) *Job {
	for _, job := range Jobs {
		if job.Name == name {
			return job
		}
	}

	return nil
}

// AddJob adds a job to the jobs slice.
func AddJob(job *Job) error {
	Jobs = append(Jobs, job)

	return SaveJobs()
}

// RemoveJob removes a job from the jobs slice.
func RemoveJob(job *Job) error {
	for i, j := range Jobs {
		if j.Name == job.Name {
			Jobs = append(Jobs[:i], Jobs[i+1:]...)
			break
		}
	}

	return SaveJobs()
}
