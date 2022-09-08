package commands

import (
	"github.com/LightningDev1/LB-Selfbot/api"
	"github.com/LightningDev1/LB-Selfbot/embed"
	"github.com/LightningDev1/LB-Selfbot/utils"
	"github.com/LightningDev1/dgc"
)

type RedditCommands struct{}

func (r RedditCommands) RegisterCommands(router *dgc.Router) {
	router.StartCategory("Reddit", "Reddit commands")
	router.RegisterCmd(&dgc.Command{
		Name:        "subreddit",
		Description: "Get a random post from a subreddit",
		Usage:       "[p]subreddit <subreddit>",
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

func (RedditCommands) generateFunc(subreddit string) func(ctx *dgc.Ctx) {
	// This makes sure that the function is not overwritten by the next command.
	// If you put the function in the RegisterCommands function, the element
	// variable will always be the last element in the command array.
	return func(ctx *dgc.Ctx) {
		post, success := api.GetImageFromSubreddit(subreddit)
		if !success {
			_ = ctx.RespondText("That subreddit is invalid.")
			return
		}
		embed.NewEmbed().
			AddField("Subreddit", subreddit).
			AddField("Title", post.Title).
			AddField("URL", post.ImageURL).
			Send(ctx)
	}
}

func (RedditCommands) subreddit(ctx *dgc.Ctx) {
	subreddit := utils.Arguments["<subreddit>"]
	post, success := api.GetImageFromSubreddit(subreddit)
	if !success {
		_ = ctx.RespondText("That subreddit is invalid.")
		return
	}
	embed.NewEmbed().
		AddField("Subreddit", subreddit).
		AddField("Title", post.Title).
		AddField("URL", post.ImageURL).
		Send(ctx)
}
