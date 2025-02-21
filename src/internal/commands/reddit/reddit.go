package reddit

import (
	"github.com/LightningDev1/LB-Selfbot/internal/api"
	"github.com/LightningDev1/LB-Selfbot/internal/router"
	"github.com/LightningDev1/LB-Selfbot/internal/services"
	"github.com/LightningDev1/LB-Selfbot/internal/utils"
	"github.com/LightningDev1/dgc"
)

func (r RedditCommands) RegisterCommands(router *dgc.Router) {
	router.StartCategory("Reddit", "Reddit commands")
	router.RegisterCmd(&dgc.Command{
		Name:        "subreddit",
		Description: "Get a random post from a subreddit",
		Arguments:   "<subreddit>",
		Handler:     r.subreddit,
	})
	subreddits := []string{
		"mildlyinteresting",
		"thatsinsane",
		"funny",
		"gaming",
		"aww",
		"pics",
		"wellthatsucks",
		"crappydesign",
		"blackmagicfuckery",
		"oddlysatisfying",
		"interestingasfuck",
		"dataisbeautiful",
		"makemesuffer",
		"memes",
		"dankmemes",
		"nextfuckinglevel",
		"foodporn",
		"programmerhumor",
		"suspiciouslyspecific",
		"mademesmile",
		"whitepeopletwitter",
		"blackpeopletwitter",
		"nonpoliticaltwitter",
		"minecraft",
		"damnthatsinteresting",
		"mildlyinfuriating",
		"extremelyinfuriating",
		"whatswrongwithyourdog",
		"humansbeingbros",
		"humansaremetal",
		"animalsbeingbros",
		"animalsbeingderps",
		"tihi",
		"tumblr",
		"wholesomememes",
		"atbge",
		"europe",
		"funnyanimals",
		"holup",
		"dndmemes",
		"gifs",
		"technicallythetruth",
		"meirl",
		"oldschoolcool",
		"pcmasterrace",
	}
	for _, subreddit := range subreddits {
		router.RegisterCmd(&dgc.Command{
			Name:        subreddit,
			Description: "Get a random post from r/" + subreddit,
			Usage:       "[p]" + subreddit,
			Handler:     r.generateFunc(subreddit),
		})
	}
}

func (RedditCommands) generateFunc(subreddit string) func(ctx *router.Ctx) {
	// This makes sure that the function is not overwritten by the next command.
	// If you put the function in the RegisterCommands function, the element
	// variable will always be the last element in the command array.
	return func(ctx *router.Ctx) {
		post, success := api.GetImageFromSubreddit(subreddit)
		if !success || post == nil {
			_ = ctx.RespondText("That subreddit is invalid.")
			return
		}

		c.services.Embed.NewEmbed().
			AddField("Subreddit", subreddit).
			AddField("Title", post.Title).
			AddField("URL", post.ImageURL).
			Send(ctx)
	}
}

func (RedditCommands) subreddit(ctx *router.Ctx) {
	subreddit := utils.Arguments["<subreddit>"]
	post, success := api.GetImageFromSubreddit(subreddit)
	if !success || post == nil {
		_ = ctx.RespondText("That subreddit is invalid.")
		return
	}
	c.services.Embed.NewEmbed().
		AddField("Subreddit", subreddit).
		AddField("Title", post.Title).
		AddField("URL", post.ImageURL).
		Send(ctx)
}

type RedditCommands struct{ services *services.Services }

func NewRedditCommands(s *services.Services) *RedditCommands {
	return &RedditCommands{services: s}
}
func (c *RedditCommands) Register() {
	cmdRouter := c.services.Discord.GetRouter()

}
