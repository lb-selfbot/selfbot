package commands

import (
	"github.com/LightningDev1/dgc"
)

type BackupCommands struct{}

func (b BackupCommands) RegisterCommands(router *dgc.Router) {
	router.StartCategory("Backup", "Create and restore backups of your account")
}
