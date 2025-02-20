package sharing

import (
	"github.com/LightningDev1/LB-Selfbot/internal/services"
	"github.com/LightningDev1/dgc"
)

type SharingCommands struct {
	services *services.Services
}

func NewSharingCommands(services *services.Services) *SharingCommands {
	return &SharingCommands{services: services}
}

func (c *SharingCommands) Register() {
	router := c.services.Discord.GetRouter()
	router.StartCategory("Sharing", "Command sharing settings")

	router.RegisterCmd(&dgc.Command{
		Name:        "share",
		Description: "Share commands with users",
		Usage:       "[p]share <add/remove/list/clear> [user]",
		Handler:     c.share,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "sharecategory",
		Description: "Share command categories",
		Usage:       "[p]sharecategory <add/remove/list/clear> [category]",
		Handler:     c.sharecategory,
	})

	// ... Register remaining sharing commands ...
}

func (c *SharingCommands) share(ctx *dgc.Ctx) {
	// ... Existing implementation ...
}
