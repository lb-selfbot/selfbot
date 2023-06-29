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

func LoadJobs(session *discordgo.Session) error {
	jobsMutex.Lock()
	defer jobsMutex.Unlock()

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
		job.Load(session)
	}

	return nil
}

func SaveJobs() error {
	jobsMutex.Lock()
	defer jobsMutex.Unlock()

	jobs := jobFileFormat{Jobs: Jobs}

	jobBytes, err := json.MarshalIndent(jobs, "", "  ")
	if err != nil {
		return err
	}

	path := filepath.Join(utils.IO.GetMainFolder(), "SlashCommands", "jobs.json")
	return utils.IO.WriteBytes(path, jobBytes)
}