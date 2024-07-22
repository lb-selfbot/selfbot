package backup

import (
	"encoding/json"
	"errors"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/LightningDev1/LB-Selfbot/selfbot/utils"
	"github.com/LightningDev1/discordgo"
)

var (
	// ErrAccountNotEligible is returned when the account is not eligible for backup/restore.
	ErrAccountNotEligible = errors.New("you must have a verified email and phone number and your account must be at least 1 year old to use this command")
)

// User represents a user in a backup.
type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

// Relationships represents a backup of friends and blocked users.
type Relationships struct {
	Friends []User `json:"friends"`
	Blocked []User `json:"blocked"`
}

// CheckEligible checks if the account is eligible for backup/restore.
func CheckEligible(session *discordgo.Session) error {
	user, err := session.User("@me")
	if err != nil {
		return err
	}

	creation := utils.Discord.SnowflakeToTime(user.ID)

	if !user.Verified || user.Phone == "" || time.Since(creation) < time.Hour*24*365 {
		return ErrAccountNotEligible
	}

	return nil
}

// SaveBackup saves a backup to a file.
func SaveBackup(backupType, backupName string, data any) error {
	if backupName != "" {
		backupName = "-" + backupName
	}

	fileName := utils.IO.SanitizeFilename(backupType + backupName + ".json")
	filePath := filepath.FromSlash(path.Join(utils.IO.GetMainFolder(), "Backups", fileName))

	err := os.MkdirAll(path.Dir(filePath), os.ModePerm)
	if err != nil {
		return err
	}

	rawData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(filePath, rawData, os.ModePerm)
	if err != nil {
		return err
	}

	return utils.System.OpenURL(filePath)
}

// RestoreBackup restores a backup from a file.
func RestoreBackup(backupType, backupName string, data any) (error) {
	if backupName != "" {
		backupName = "-" + backupName
	}

	fileName := utils.IO.SanitizeFilename(backupType + backupName + ".json")
	filePath := filepath.FromSlash(path.Join(utils.IO.GetMainFolder(), "Backups", fileName))

	rawData, err := os.ReadFile(filePath)
	if err != nil {
		return nil
	}

	return json.Unmarshal(rawData, &data)
}
