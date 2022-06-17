package commands

import (
	"github.com/LightningDev/LB-Selfbot/api"
	"github.com/LightningDev/LB-Selfbot/embed"
	"github.com/LightningDev/LB-Selfbot/utils"
	"github.com/LightningDev1/dgc"
)

type RedditCommands struct{}

func (t RedditCommands) RegisterCommands(router *dgc.Router) {
	router.RegisterCmd(&dgc.Command{
		Name:        "subreddit",
		Description: "Get a random post from a subreddit",
		Usage:       "[p]subreddit <subreddit>",
		Handler:     t.subreddit,
	})
	redditCommands := []string{
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
		"dankmemes",
		"nextfuckinglevel",
		"foodporn",
	}
	for _, subreddit := range redditCommands {
		router.RegisterCmd(&dgc.Command{
			Name:        subreddit,
			Description: "Get a random post from r/" + subreddit,
			Usage:       "[p]" + subreddit,
			Handler: func(ctx *dgc.Ctx) {
				post, success := api.GetImageFromSubreddit(subreddit)
				if !success {
					ctx.RespondText("That subreddit is invalid.")
					return
				}
				embed.NewEmbed().
					AddField("Subreddit", subreddit).
					AddField("Title", post.Title).
					AddField("URL", post.ImageURL).
					Send(ctx)
			},
		})
	}
}

func (t RedditCommands) subreddit(ctx *dgc.Ctx) {
	subreddit := utils.Arguments["<subreddit>"]
	post, success := api.GetImageFromSubreddit(subreddit)
	if !success {
		ctx.RespondText("That subreddit is invalid.")
		return
	}
	embed.NewEmbed().
		AddField("Subreddit", subreddit).
		AddField("Title", post.Title).
		AddField("URL", post.ImageURL).
		Send(ctx)
}
