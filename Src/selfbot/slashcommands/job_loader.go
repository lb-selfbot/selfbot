package slashcommands

import (
	"encoding/json"
	"sync"

	"github.com/LightningDev1/LB-Selfbot/selfbot/config"
	"github.com/LightningDev1/LB-Selfbot/selfbot/utils"
	"github.com/LightningDev1/discordgo"
)

var jobsMutex = &sync.Mutex{}
var Jobs []*Job

func LoadJobs(session *discordgo.Session) error {
	jobsMutex.Lock()
	defer jobsMutex.Unlock()

	cfg, err := config.GetConfig()
	if err != nil {
		return err
	}

	_ = json.Unmarshal(cfg.SlashCommandJobs, &Jobs)

	for _, job := range Jobs {
		err = job.Load(session)
		if err != nil {
			utils.Logging.Error("Failed to load job:", err.Error())
		}
	}

	return nil
}

func SaveJobs() error {
	jobsMutex.Lock()
	defer jobsMutex.Unlock()

	cfg, err := config.GetConfig()
	if err != nil {
		return err
	}

	jobsBytes, err := json.MarshalIndent(Jobs, "", "  ")
	if err != nil {
		return err
	}

	cfg.SlashCommandJobs = jobsBytes

	return cfg.Save()
}

func GetJobByID(id string) *Job {
	for _, job := range Jobs {
		if job.ID == id {
			return job
		}
	}

	return nil
}

func AddJob(job *Job) error {
	Jobs = append(Jobs, job)

	return SaveJobs()
}

func RemoveJob(job *Job) error {
	for i, j := range Jobs {
		if j.ID == job.ID {
			Jobs = append(Jobs[:i], Jobs[i+1:]...)
			break
		}
	}

	return SaveJobs()
}
