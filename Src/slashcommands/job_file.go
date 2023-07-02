package slashcommands

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"

	"github.com/LightningDev1/LB-Selfbot/utils"
	"github.com/LightningDev1/discordgo"
)

var jobsMutex = &sync.Mutex{}

var Jobs []*Job

type jobFileFormat struct {
	Jobs []*Job `json:"jobs"`
}

func ensureFileExists() error {
	path := filepath.Join(utils.IO.GetMainFolder(), "SlashCommands", "jobs.json")
	if !utils.IO.Exists(path) {
		return utils.IO.WriteBytes(path, []byte("{}"))
	}

	return nil
}

func LoadJobs(session *discordgo.Session) error {
	jobsMutex.Lock()
	defer jobsMutex.Unlock()

	err := ensureFileExists()
	if err != nil {
		return err
	}

	path := filepath.Join(utils.IO.GetMainFolder(), "SlashCommands", "jobs.json")
	jobsFile, err := os.Open(path)
	if err != nil {
		return err
	}

	defer jobsFile.Close()

	var jobs jobFileFormat

	err = json.NewDecoder(jobsFile).Decode(&jobs)
	if err != nil {
		return err
	}

	Jobs = jobs.Jobs

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

	err := ensureFileExists()
	if err != nil {
		return err
	}

	jobs := jobFileFormat{Jobs: Jobs}

	jobBytes, err := json.MarshalIndent(jobs, "", "  ")
	if err != nil {
		return err
	}

	path := filepath.Join(utils.IO.GetMainFolder(), "SlashCommands", "jobs.json")
	return utils.IO.WriteBytes(path, jobBytes)
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
