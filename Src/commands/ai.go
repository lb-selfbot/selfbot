package commands

import (
	"sort"
	"strings"

	"github.com/LightningDev1/LB-Selfbot/api"
	"github.com/LightningDev1/LB-Selfbot/embed"
	"github.com/LightningDev1/LB-Selfbot/utils"
	"github.com/LightningDev1/dgc"
)

type AICommands struct{}

func (a *AICommands) RegisterCommands(router *dgc.Router) {
	router.StartCategory("AI", "Interact with AIs")

	router.RegisterCmd(&dgc.Command{
		Name:        "chatgpt",
		Description: "Talk to OpenAI GPT-3.5",
		Usage:       "[p]chatgpt <\"system\"> <\"message\">",
		Handler:     a.chatgpt,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "clydeai",
		Description: "Talk to Discord's Clyde AI",
		Usage:       "[p]clydeai <message>",
		Handler:     a.clydeai,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "imaginestyles",
		Description: "List of available styles for [p]imagine",
		Usage:       "[p]imaginestyles",
		Handler:     a.imaginestyles,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "imaginemodels",
		Description: "List of available models for [p]imagine",
		Usage:       "[p]imaginemodels",
		Handler:     a.imaginemodels,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "imagine",
		Description: "Generate an image from text",
		Usage:       "[p]imagine <\"prompt\"> [\"negative prompt\"] [\"model\"] [\"style\"]",
		Handler:     a.imagine,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "ocrlink",
		Description: "Extract text from an image URL",
		Usage:       "[p]ocrlink <link>",
		Handler:     a.ocrlink,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "ocrfile",
		Description: "Extract text from an image file",
		Usage:       "[p]ocrfile [attachment]",
		Handler:     a.ocrfile,
	})
}

func (AICommands) chatgpt(ctx *dgc.Ctx) {
	response := api.GPT3(utils.Arguments["<\"system\">"], utils.Arguments["<\"message\">"])

	embed.NewEmbed().
		SetTitle("ChatGPT (GPT-3.5)").
		AddField("System", utils.Arguments["<\"system\">"]).
		AddField("Message", utils.Arguments["<\"message\">"]).
		AddField("Response", response).
		Send(ctx)
}

func (AICommands) clydeai(ctx *dgc.Ctx) {
	response, err := api.ClydeAI(ctx.Arguments.Raw())
	if err != nil {
		_ = ctx.RespondText("Error getting response: " + err.Error())
		return
	}

	response = strings.ReplaceAll(response, "1101230589195079822", ctx.Session.State.User.ID)

	embed.NewEmbed().
		SetTitle("Clyde AI").
		AddField("Message", ctx.Arguments.Raw()).
		AddField("Response", response).
		Send(ctx)
}

func (AICommands) imaginestyles(ctx *dgc.Ctx) {
	stylesString := []string{}

	for k := range api.DefaultStyles {
		stylesString = append(stylesString, k)
	}

	sort.Strings(stylesString)

	embed.NewEmbed().
		SetTitle("Image Styles").
		SetDescription(strings.Join(stylesString, "\n")).
		Send(ctx)
}

func (AICommands) imaginemodels(ctx *dgc.Ctx) {
	models := api.GetImageModels()

	modelsString := []string{}

	for k := range models {
		modelsString = append(modelsString, k)
	}

	sort.Strings(modelsString)

	embed.NewEmbed().
		SetTitle("Image Models").
		SetDescription(strings.Join(modelsString, "\n")).
		Send(ctx)
}

func (AICommands) imagine(ctx *dgc.Ctx) {
	prompt := utils.Arguments["<\"prompt\">"]

	style := strings.ToLower(utils.Arguments["[\"style\"]"])
	if styleValue, ok := api.DefaultStyles[style]; ok {
		prompt += ", " + styleValue
	}

	var model string

	if modelValue, ok := api.DefaultModels[model]; ok {
		model = modelValue
	}

	image, err := api.GenerateImage(prompt, utils.Arguments["[\"negative prompt\"]"], model)
	if err != nil {
		_ = ctx.RespondText("Error getting response: " + err.Error())
		return
	}

	_, _ = ctx.Session.ChannelFileSend(ctx.Event.ChannelID, "image.png", image)
}

func (AICommands) ocrlink(ctx *dgc.Ctx) {
	text, err := api.OCR(ctx.Arguments.Raw())
	if err != nil {
		_ = ctx.RespondText("Error getting response: " + err.Error())
		return
	}

	embed.NewEmbed().
		SetTitle("OCR").
		AddField("URL", ctx.Arguments.Raw()).
		AddField("Extracted Text", text).
		Send(ctx)
}

func (AICommands) ocrfile(ctx *dgc.Ctx) {
	attachments := ctx.Event.Message.Attachments

	if len(attachments) == 0 {
		_ = ctx.RespondText("Error: you must provide an image as an attachment.")
		return
	}

	text, err := api.OCR(attachments[0].URL)
	if err != nil {
		_ = ctx.RespondText("Error getting response: " + err.Error())
		return
	}

	embed.NewEmbed().
		SetTitle("OCR").
		AddField("URL", attachments[0].URL).
		AddField("Extracted Text", text).
		Send(ctx)
}
