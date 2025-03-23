package reddit

import (
	"github.com/lb-selfbot/selfbot/internal/api"
	"github.com/lb-selfbot/selfbot/internal/services/router"

	"github.com/lb-selfbot/selfbot/internal/services"
)

type RedditCommands struct{ services *services.Services }

func NewRedditCommands(s *services.Services) *RedditCommands {
	return &RedditCommands{services: s}
}

func (c *RedditCommands) Register() {
	cmdRouter := c.services.Discord.GetRouter()

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

	commands := make([]*router.Command, 0, len(subreddits)+1)

	commands = append(commands, &router.Command{
		Name:        "subreddit",
		Description: "Get a random post from a subreddit",
		Arguments:   "<name>",
		Handler:     c.subreddit,
	})

	for _, subreddit := range subreddits {
		commands = append(commands, &router.Command{
			Name:        subreddit,
			Description: "Get a random post from r/" + subreddit,
			Handler:     c.generateFunc(subreddit),
		})
	}

	cmdRouter.RegisterCategory(&router.Category{
		Name:        "Reddit",
		Description: "Reddit commands",
	}, commands)
}

func (c *RedditCommands) generateFunc(subreddit string) func(ctx *router.Ctx) {
	// This makes sure that the function is not overwritten by the next command.
	// If you put the function in the RegisterCommands function, the element
	// variable will always be the last element in the command array.
	return func(ctx *router.Ctx) {
		post, success := api.GetImageFromSubreddit(subreddit)
		if !success || post == nil {
			ctx.RespondText("That subreddit is invalid.")
			return
		}

		c.services.Embed.NewEmbed().
			AddField("Subreddit", subreddit).
			AddField("Title", post.Title).
			AddField("URL", post.ImageURL).
			Send(ctx)
	}
}

func (c *RedditCommands) subreddit(ctx *router.Ctx) {
	subreddit := ctx.Arguments.Get(0).Raw()

	post, success := api.GetImageFromSubreddit(subreddit)
	if !success || post == nil {
		ctx.RespondText("That subreddit is invalid.")
		return
	}

	c.services.Embed.NewEmbed().
		AddField("Subreddit", subreddit).
		AddField("Title", post.Title).
		AddField("URL", post.ImageURL).
		Send(ctx)
}
