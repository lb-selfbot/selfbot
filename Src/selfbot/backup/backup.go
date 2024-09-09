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
	"github.com/LightningDev1/discordgo/protos"
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

// Server represents a server in a backup.
type Server struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Invite string `json:"invite"`
}

// Relationships represents a backup of friends and blocked users.
type Relationships struct {
	Friends []User `json:"friends"`
	Blocked []User `json:"blocked"`
}

// Servers represents a backup of the user's servers.
type Servers struct {
	Servers []Server `json:"servers"`
}

// Profile represents a backup of the user's profile.
type Profile struct {
	User    *discordgo.User `json:"user"`
	Profile map[string]any  `json:"profile"`
	Avatar  string          `json:"avatar"`
}

// Settings represents a backup of the user's settings.
type Settings struct {
	Settings *protos.PreloadedUserSettings `json:"settings"`
}

// GIFs represents a backup of the user's saved GIFs.
type GIFs struct {
	Settings *protos.FrecencyUserSettings `json:"settings"`
}

// CheckEligible checks if the account is eligible for backup/restore.
func CheckEligible(session *discordgo.Session) error {
	user, err := session.User("@me")
	if err != nil {
		return err
	}

	time.Sleep(time.Second)

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
func RestoreBackup(backupType, backupName string, data any) error {
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

// FindInvitableChannel finds an invitable channel in a guild.
func FindInvitableChannel(session *discordgo.Session, guild *discordgo.Guild) *discordgo.Channel {
	for _, channel := range guild.Channels {
		validChannel := channel.Type == discordgo.ChannelTypeGuildText ||
			channel.Type == discordgo.ChannelTypeGuildVoice

		if !validChannel {
			continue
		}

		perms, err := utils.Discord.UserChannelPermissions(session, guild.ID, channel.ID, session.State.User.ID)
		if err != nil {
			continue
		}

		if perms&discordgo.PermissionCreateInstantInvite == 0 ||
			perms&discordgo.PermissionViewChannel == 0 {
			continue
		}

		return channel
	}

	return nil
}
