package interfaces

import "github.com/LightningDev1/discordgo"

type BackupService interface {
	Inject(discordSvc DiscordService)

	CheckEligibility() error
	FindInvitableChannel(guild *discordgo.Guild) (*discordgo.Channel, error)

	SaveBackup(bType, bName string, data any) error
	LoadBackup(bType, bName string, data any) error
}
