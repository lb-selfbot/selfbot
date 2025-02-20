package nsfw

import (
	"github.com/LightningDev1/LB-Selfbot/internal/services"
)

type NSFWCommands struct {
	services *services.Services
}

func NewNSFWCommands(services *services.Services) *NSFWCommands {
	return &NSFWCommands{services: services}
}

func (c *NSFWCommands) Register() {
	router := c.services.Discord.GetRouter()
	router.StartCategory("NSFW", "NSFW commands (18+ only)")

	// Register NSFW commands...
}

// Command implementations...
