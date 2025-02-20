package backup

import (
	"github.com/LightningDev1/LB-Selfbot/internal/services"
	"github.com/LightningDev1/dgc"
)

type BackupCommands struct {
	services *services.Services
}

func NewBackupCommands(services *services.Services) *BackupCommands {
	return &BackupCommands{services: services}
}

func (c *BackupCommands) Register() {
	router := c.services.Discord.GetRouter()
	router.StartCategory("Backup", "Backup and restore commands")

	router.RegisterCmd(&dgc.Command{
		Name:        "backup",
		Description: "Create a backup of server/channel/user data",
		Usage:       "[p]backup <server/channel/user> <id>",
		Handler:     c.backup,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "restore",
		Description: "Restore a backup",
		Usage:       "[p]restore <backup id>",
		Handler:     c.restore,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "listbackups",
		Description: "List all backups",
		Usage:       "[p]listbackups",
		Handler:     c.listbackups,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "deletebackup",
		Description: "Delete a backup",
		Usage:       "[p]deletebackup <backup id>",
		Handler:     c.deletebackup,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "backupinfo",
		Description: "Get info about a backup",
		Usage:       "[p]backupinfo <backup id>",
		Handler:     c.backupinfo,
	})
}

func (c *BackupCommands) backup(ctx *dgc.Ctx) {
	// ... Existing implementation ...
}

// ... Keep all other existing command implementations ...
