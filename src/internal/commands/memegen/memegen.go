package memegen

import (
	"fmt"
	"math/rand"

	"github.com/LightningDev1/LB-Selfbot/internal/api"
	"github.com/LightningDev1/LB-Selfbot/internal/http"
	"github.com/LightningDev1/LB-Selfbot/internal/router"
	"github.com/LightningDev1/LB-Selfbot/internal/services"
	"github.com/LightningDev1/LB-Selfbot/internal/utils"
)

type MemeGenCommands struct{ services *services.Services }

func NewMemeGenCommands(s *services.Services) *MemeGenCommands {
	return &MemeGenCommands{services: s}
}
func (c *MemeGenCommands) Register() {
	cmdRouter := c.services.Discord.GetRouter()

	cmdRouter.RegisterCategory(&router.Category{
		Name:        "Memegen",
		Description: "Generate memes using imgflip",
	}, []*router.Command{
		{
			Name:        "randommeme",
			Description: "Send a random meme with custom text",
			Arguments:   "<\"top\"> <\"bottom\">",
			Handler:     m.randommeme,
		},
		{
			Name:        "custommeme",
			Description: "Send a meme based on a template",
			Arguments:   "<template id> <\"top\"> <\"bottom\">",
			Handler:     m.custommeme,
		},
		{
			Name:        meme[0],
			Description: fmt.Sprintf("Send a %s meme with custom text", meme[2]),
			Arguments:   fmt.Sprintf("%s <\"top\"> <\"bottom\">", meme[0]),
			Handler:     m.generateFunc(meme[1]),
		},
	})
}

func (MemeGenCommands) randommeme(ctx *router.Ctx) {
	template, err := http.GetJson("https://api.imgflip.com/get_memes")
	if err != nil {
		_ = ctx.RespondText("An error has occurred.")
		return
	}

	memes := template.Get("data").Get("memes")
	memesLen := len(memes.MustArray())

	// Get a random meme from the list
	meme := memes.GetIndex(rand.Intn(memesLen))
	templateID := meme.Get("id").MustString()

	image, err := api.CreateMeme(templateID, utils.Arguments["<\"top\">"], utils.Arguments["<\"bottom\">"])
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
		return
	}

	c.services.Embed.NewEmbed().
		AddField("Used Meme", meme.Get("name").MustString()).
		SetImage(image).
		Send(ctx)
}

func (MemeGenCommands) custommeme(ctx *router.Ctx) {
	image, err := api.CreateMeme(utils.Arguments["<template id>"], utils.Arguments["<\"top\">"], utils.Arguments["<\"bottom\">"])
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
		return
	}

	_ = ctx.RespondText(image)
}

func (MemeGenCommands) generateFunc(templateID string) func(ctx *router.Ctx) {
	// This makes sure that the function is not overwritten by the next command.
	// If you put the function in the RegisterCommands function, the element
	// variable will always be the last element in the command array.

	return func(ctx *router.Ctx) {
		image, err := api.CreateMeme(templateID, utils.Arguments["<\"top\">"], utils.Arguments["<\"bottom\">"])
		if err != nil {
			_ = ctx.RespondText("An error has occurred: " + err.Error())
			return
		}

		_ = ctx.RespondText(image)
	}
}
