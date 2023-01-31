package commands

import (
	"fmt"
	"math/rand"
	"net/url"

	"github.com/LightningDev1/LB-Selfbot/embed"
	"github.com/LightningDev1/LB-Selfbot/http"
	"github.com/LightningDev1/LB-Selfbot/utils"
	"github.com/LightningDev1/dgc"
)

type MemeGenCommands struct{}

func (m *MemeGenCommands) RegisterCommands(router *dgc.Router) {
	router.RegisterCmd(&dgc.Command{
		Name:        "randommeme",
		Description: "Send a random meme with custom text",
		Usage:       "[p]randommeme <\"top\"> <\"bottom\">",
		Handler:     m.randommeme,
	})
}

func (MemeGenCommands) randommeme(ctx *dgc.Ctx) {
	template, err := http.GetJson("https://api.imgflip.com/get_memes")
	if err != nil {
		_ = ctx.RespondText("An error has occurred.")
		return
	}

	memes := template.Get("data").Get("memes")
	memesLen := len(memes.MustArray())

	// Get a random meme from the list
	meme := memes.GetIndex(rand.Intn(memesLen))
	templateId := meme.Get("id").MustString()

	uri := "https://api.imgflip.com/caption_image?template_id=%s&username=LightningBotSB&password=LightningBotSB&text0=%s&text1=%s"
	uri = fmt.Sprintf(uri, templateId, url.QueryEscape(utils.Arguments["<\"top\">"]), url.QueryEscape(utils.Arguments["<\"bottom\">"]))

	image, err := http.GetJson(uri)
	if err != nil {
		_ = ctx.RespondText("An error has occurred.")
		return
	}

	embed.NewEmbed().
		AddField("Used Meme", meme.Get("name").MustString()).
		SetImage(image.GetPath("data", "url").MustString()).
		Send(ctx)
}
