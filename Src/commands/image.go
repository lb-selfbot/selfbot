package commands

import (
	"fmt"
	"net/url"

	"github.com/LightningDev1/LB-Selfbot/embed"
	"github.com/LightningDev1/LB-Selfbot/http"
	"github.com/LightningDev1/dgc"
)

type ImageCommands struct{}

func (i ImageCommands) RegisterCommands(router *dgc.Router) {
	router.RegisterCmd(&dgc.Command{
		Name:        "tweet",
		Description: "Tweet a message",
		Usage:       "[p]tweet <text>",
		Handler:     i.tweet,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "meme",
		Description: "Send a random meme",
		Usage:       "[p]meme",
		Handler:     i.meme,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "ytcomment",
		Description: "Make a fake YouTube comment",
		Usage:       "[p]ytcomment <text>",
		Handler:     i.ytcomment,
	})
}

func (ImageCommands) tweet(ctx *dgc.Ctx) {
	jsonResult, err := http.GetJson(
		fmt.Sprintf("https://nekobot.xyz/api/imagegen?type=tweet&text=%s&username=%s", url.QueryEscape(ctx.Arguments.Raw()), url.QueryEscape(ctx.Event.Author.Username)),
	)
	if err != nil {
		ctx.RespondText("An error has occurred: " + err.Error())
		return
	}
	ctx.RespondText(jsonResult.Get("message").MustString())
}

func (ImageCommands) meme(ctx *dgc.Ctx) {
	jsonResult, err := http.GetJson("https://meme-api.herokuapp.com/gimme")
	if err != nil {
		ctx.RespondText("An error has occurred: " + err.Error())
		return
	}
	embed.NewEmbed().
		SetTitle(jsonResult.Get("title").MustString()).
		SetImage(jsonResult.Get("url").MustString()).
		AddField("URL", fmt.Sprintf("`%s`", jsonResult.Get("postLink").MustString())).
		SetFooter(fmt.Sprintf("Subreddit: %s | Upvotes: %d | Posted by: u/%s", jsonResult.Get("subreddit").MustString(), jsonResult.Get("ups").MustInt(), jsonResult.Get("author").MustString())).
		Send(ctx)
}

func (ImageCommands) ytcomment(ctx *dgc.Ctx) {
	ctx.RespondText(fmt.Sprintf(
		"https://some-random-api.ml/canvas/youtube-comment?avatar=%s&comment=%s&username=%s",
		url.QueryEscape(ctx.Event.Author.AvatarURL("")),
		url.QueryEscape(ctx.Arguments.Raw()),
		url.QueryEscape(ctx.Event.Author.Username),
	))
}
