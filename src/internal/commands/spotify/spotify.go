package spotify

import (
	"github.com/LightningDev1/LB-Selfbot/internal/services"
	"github.com/LightningDev1/dgc"
)

type SpotifyCommands struct {
	services *services.Services
}

func NewSpotifyCommands(services *services.Services) *SpotifyCommands {
	return &SpotifyCommands{services: services}
}

func (c *SpotifyCommands) Register() {
	router := c.services.Discord.GetRouter()
	router.StartCategory("Spotify", "Spotify integration commands")

	router.RegisterCmd(&dgc.Command{
		Name:        "spotify",
		Description: "Show current Spotify track",
		Usage:       "[p]spotify",
		Handler:     c.spotify,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "spotifyadd",
		Description: "Add current track to playlist",
		Usage:       "[p]spotifyadd",
		Handler:     c.spotifyadd,
	})

	// ... Register remaining spotify commands ...
}

func (c *SpotifyCommands) spotify(ctx *dgc.Ctx) {
	// ... Existing implementation ...
}
