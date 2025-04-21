package backup

import (
	"errors"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/goccy/go-json"

	"github.com/lb-selfbot/discordgo"
	"github.com/lb-selfbot/discordgo/protos"
	"github.com/lb-selfbot/selfbot/internal/fs"
	"github.com/lb-selfbot/selfbot/internal/system"
	"github.com/lb-selfbot/selfbot/internal/utils"
)

var (
	// ErrAccountNotEligible is returned when the account is not eligible for backup/restore.
	ErrAccountNotEligible = errors.New("you must have a verified email and phone number and your account must be at least 1 year old to use this command")

	// ErrNoInvitableChannel is returned when no invitable channel is found in a guild.
	ErrNoInvitableChannel = errors.New("no invitable channel found in the guild")
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

// Service handles Discord account backup and restore operations.
type Service interface {
	// CheckEligibility checks if the account is eligible for backup/restore.
	CheckEligibility(session *discordgo.Session) error

	// FindInvitableChannel finds an invitable channel in a guild.
	FindInvitableChannel(session *discordgo.Session, guild *discordgo.Guild) (*discordgo.Channel, error)

	// SaveBackup saves a backup to a file.
	SaveBackup(bType, bName string, data any) error

	// LoadBackup restores a backup from a file.
	LoadBackup(bType, bName string, data any) error
}

type backupService struct{}

// NewService creates a new BackupService instance.
func NewService() Service {
	return &backupService{}
}

// CheckEligibility checks if the account is eligible for backup/restore.
func (b *backupService) CheckEligibility(session *discordgo.Session) error {
	user, err := session.User("@me")
	if err != nil {
		return err
	}

	time.Sleep(time.Second)

	creation := discordgo.MustSnowflakeTimestamp(user.ID)

	if !user.Verified || user.Phone == "" || time.Since(creation) < time.Hour*24*365 {
		return ErrAccountNotEligible
	}

	return nil
}

// FindInvitableChannel finds an invitable channel in a guild.
func (b *backupService) FindInvitableChannel(session *discordgo.Session, guild *discordgo.Guild) (*discordgo.Channel, error) {
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

		return channel, nil
	}

	return nil, ErrNoInvitableChannel
}

// SaveBackup saves a backup to a file.
func (b *backupService) SaveBackup(bType, bName string, data any) error {
	if bName != "" {
		bName = "-" + bName
	}

	fileName := fs.SanitizeFilename(bType + bName + ".json")
	filePath := filepath.FromSlash(path.Join(fs.GetMainFolder(), "Backups", fileName))

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

	return system.OpenURL(filePath)
}

// LoadBackup restores a backup from a file.
func (b *backupService) LoadBackup(bType, bName string, data any) error {
	if bName != "" {
		bName = "-" + bName
	}

	fileName := fs.SanitizeFilename(bType + bName + ".json")
	filePath := filepath.FromSlash(path.Join(fs.GetMainFolder(), "Backups", fileName))

	rawData, err := os.ReadFile(filePath)
	if err != nil {
		return nil
	}

	return json.Unmarshal(rawData, &data)
}
