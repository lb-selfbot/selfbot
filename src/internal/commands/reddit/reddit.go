package reddit

import (
	"github.com/LightningDev1/LB-Selfbot/internal/services"
	"github.com/LightningDev1/dgc"
)

type RedditCommands struct {
	services *services.Services
}

func NewRedditCommands(services *services.Services) *RedditCommands {
	return &RedditCommands{services: services}
}

func (c *RedditCommands) Register() {
	router := c.services.Discord.GetRouter()
	router.StartCategory("Reddit", "Reddit integration commands")

	router.RegisterCmd(&dgc.Command{
		Name:        "reddit",
		Description: "Get posts from a subreddit",
		Usage:       "[p]reddit <subreddit> [sort]",
		Handler:     c.reddit,
	})

	// ... Register remaining reddit commands ...
}

func (c *RedditCommands) reddit(ctx *dgc.Ctx) {
	// ... Existing implementation ...
}
