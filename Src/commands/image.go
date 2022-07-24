package commands

import (
	"fmt"
	"net/url"

	"github.com/LightningDev1/LB-Selfbot/embed"
	"github.com/LightningDev1/LB-Selfbot/http"
	"github.com/LightningDev1/LB-Selfbot/utils"
	"github.com/LightningDev1/dgc"
)

type ImageCommands struct{}

func (i ImageCommands) RegisterCommands(router *dgc.Router) {
	router.RegisterCmd(&dgc.Command{
		Name:        "tweet",
		Description: "Tweet a message",
		Usage:       "[p]tweet <username> <text>",
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
		Usage:       "[p]ytcomment <user> <text>",
		Handler:     i.ytcomment,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "pokemon",
		Description: "Send a random Pokémon image",
		Usage:       "[p]pokemon",
		Handler:     i.pokemon,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "wasted",
		Description: "Add a wasted overlay on the users avatar",
		Usage:       "[p]wasted [user]",
		Handler:     i.wasted,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "triggered",
		Description: "Add a triggered overlay on the users avatar",
		Usage:       "[p]triggered [user]",
		Handler:     i.triggered,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "glass",
		Description: "Add a glass overlay on the users avatar",
		Usage:       "[p]glass [user]",
		Handler:     i.glass,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "greyscale",
		Description: "Add a greyscale overlay on the users avatar",
		Usage:       "[p]greyscale [user]",
		Handler:     i.greyscale,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "invert",
		Description: "Invert the users avatar",
		Usage:       "[p]invert [user]",
		Handler:     i.invert,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "invertgreyscale",
		Description: "Invert and greyscale the users avatar",
		Usage:       "[p]invertgreyscale [user]",
		Handler:     i.invertgreyscale,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "brightness",
		Description: "Make the users avatar brighter",
		Usage:       "[p]brightness [user]",
		Handler:     i.brightness,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "red",
		Description: "Make the users avatar red",
		Usage:       "[p]red [user]",
		Handler:     i.red,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "green",
		Description: "Make the users avatar green",
		Usage:       "[p]green [user]",
		Handler:     i.green,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "blue",
		Description: "Make the users avatar blue",
		Usage:       "[p]blue [user]",
		Handler:     i.blue,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "pixelate",
		Description: "Pixelate the users avatar",
		Usage:       "[p]pixelate [user]",
		Handler:     i.pixelate,
	})
}

func (ImageCommands) tweet(ctx *dgc.Ctx) {
	jsonResult, err := http.GetJson(fmt.Sprintf(
		"https://nekobot.xyz/api/imagegen?type=tweet&text=%s&username=%s",
		url.QueryEscape(utils.String.RemoveWords(ctx.Arguments.Raw(), 1)),
		url.QueryEscape(utils.Arguments["<username>"]),
	))
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
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Event.Author
	}
	ctx.RespondText(fmt.Sprintf(
		"https://some-random-api.ml/canvas/youtube-comment?avatar=%s&comment=%s&username=%s",
		url.QueryEscape(user.AvatarURL("")),
		url.QueryEscape(utils.String.RemoveWords(ctx.Arguments.Raw(), 1)),
		url.QueryEscape(user.Username),
	))
}

func (ImageCommands) pokemon(ctx *dgc.Ctx) {
	ctx.RespondText(fmt.Sprintf("https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/%d.png", utils.Random.Number(1, 898)))
}

func (ImageCommands) wasted(ctx *dgc.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Event.Author
	}
	ctx.RespondText(fmt.Sprintf("https://some-random-api.ml/canvas/wasted?avatar=%s", url.QueryEscape(user.AvatarURL(""))))
}

func (ImageCommands) triggered(ctx *dgc.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Event.Author
	}
	ctx.RespondText(fmt.Sprintf("https://some-random-api.ml/canvas/triggered?avatar=%s", url.QueryEscape(user.AvatarURL(""))))
}

func (ImageCommands) glass(ctx *dgc.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Event.Author
	}
	ctx.RespondText(fmt.Sprintf("https://some-random-api.ml/canvas/glass?avatar=%s", url.QueryEscape(user.AvatarURL(""))))
}

func (ImageCommands) greyscale(ctx *dgc.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Event.Author
	}
	ctx.RespondText(fmt.Sprintf("https://some-random-api.ml/canvas/greyscale?avatar=%s", url.QueryEscape(user.AvatarURL(""))))
}

func (ImageCommands) invert(ctx *dgc.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Event.Author
	}
	ctx.RespondText(fmt.Sprintf("https://some-random-api.ml/canvas/invert?avatar=%s", url.QueryEscape(user.AvatarURL(""))))
}

func (ImageCommands) invertgreyscale(ctx *dgc.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Event.Author
	}
	ctx.RespondText(fmt.Sprintf("https://some-random-api.ml/canvas/invertgreyscale?avatar=%s", url.QueryEscape(user.AvatarURL(""))))
}

func (ImageCommands) brightness(ctx *dgc.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Event.Author
	}
	ctx.RespondText(fmt.Sprintf("https://some-random-api.ml/canvas/brightness?avatar=%s", url.QueryEscape(user.AvatarURL(""))))
}

func (ImageCommands) red(ctx *dgc.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Event.Author
	}
	ctx.RespondText(fmt.Sprintf("https://some-random-api.ml/canvas/red?avatar=%s", url.QueryEscape(user.AvatarURL(""))))
}

func (ImageCommands) green(ctx *dgc.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Event.Author
	}
	ctx.RespondText(fmt.Sprintf("https://some-random-api.ml/canvas/green?avatar=%s", url.QueryEscape(user.AvatarURL(""))))
}

func (ImageCommands) blue(ctx *dgc.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Event.Author
	}
	ctx.RespondText(fmt.Sprintf("https://some-random-api.ml/canvas/blue?avatar=%s", url.QueryEscape(user.AvatarURL(""))))
}

func (ImageCommands) pixelate(ctx *dgc.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Event.Author
	}
	ctx.RespondText(fmt.Sprintf("https://some-random-api.ml/canvas/pixelate?avatar=%s", url.QueryEscape(user.AvatarURL(""))))
}
