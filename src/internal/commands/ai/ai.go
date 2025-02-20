package ai

import (
	"sort"
	"strings"

	"github.com/LightningDev1/LB-Selfbot/internal/api"
	"github.com/LightningDev1/LB-Selfbot/internal/router"
	"github.com/LightningDev1/LB-Selfbot/internal/services"
	"github.com/LightningDev1/LB-Selfbot/internal/utils"
)

type AICommands struct{ services *services.Services }

func NewAICommands(s *services.Services) *AICommands {
	return &AICommands{services: s}
}
func (c *AICommands) Register() {
	cmdRouter := c.services.Discord.GetRouter()

	cmdRouter.RegisterCategory(&router.Category{
		Name:        "AI",
		Description: "Interact with AI's",
	}, []*router.Command{
		{
			Name:        "chatgpt",
			Description: "Talk to OpenAI GPT-3.5",
			Arguments:   "chatgpt <\"message\"> [\"system\"]",
			Handler:     c.chatgpt,
		},
		{
			Name:        "gpt4",
			Description: "Talk to OpenAI GPT-4",
			Arguments:   "gpt4 <\"message\"> [\"system\"]",
			Handler:     c.gpt4,
		},
		{
			Name:        "imaginestyles",
			Description: "List of available styles for [p]imagine",
			Arguments:   "imaginestyles",
			Handler:     c.imaginestyles,
		},
		{
			Name:        "imaginemodels",
			Description: "List of available models for [p]imagine",
			Arguments:   "imaginemodels",
			Handler:     c.imaginemodels,
		},
		{
			Name:        "imagine",
			Description: "Generate an image from text",
			Arguments:   "imagine <\"prompt\"> [\"negative prompt\"] [\"model\"] [\"style\"]",
			Handler:     c.imagine,
		},
		{
			Name:        "ocrlink",
			Description: "Extract text from an image URL",
			Arguments:   "ocrlink <link>",
			Handler:     c.ocrlink,
		},
		{
			Name:        "ocrfile",
			Description: "Extract text from an image file",
			Arguments:   "ocrfile [attachment]",
			Handler:     c.ocrfile,
		},
	})
}

func (AICommands) chatgpt(ctx *router.Ctx) {
	response, err := api.GPT(utils.Arguments["[\"system\"]"], utils.Arguments["<\"message\">"])
	if err != nil {
		_ = ctx.RespondText("Error getting response: " + err.Error())
		return
	}

	c.services.Embed.NewEmbed().
		SetTitle("ChatGPT (GPT-3.5)").
		AddField("System", utils.Arguments["[\"system\"]"]).
		AddField("Message", utils.Arguments["<\"message\">"]).
		AddField("Response", response).
		Send(ctx)
}

func (AICommands) gpt4(ctx *router.Ctx) {
	response, err := api.GPT(utils.Arguments["[\"system\"]"], utils.Arguments["<\"message\">"], "gpt-4")
	if err != nil {
		_ = ctx.RespondText("Error getting response: " + err.Error())
		return
	}

	c.services.Embed.NewEmbed().
		SetTitle("ChatGPT (GPT-4)").
		AddField("System", utils.Arguments["[\"system\"]"]).
		AddField("Message", utils.Arguments["<\"message\">"]).
		AddField("Response", response).
		Send(ctx)
}

func (AICommands) imaginestyles(ctx *router.Ctx) {
	stylesString := []string{}

	for k := range api.DefaultStyles {
		stylesString = append(stylesString, k)
	}

	sort.Strings(stylesString)

	c.services.Embed.NewEmbed().
		SetTitle("Image Styles").
		SetDescription(strings.Join(stylesString, "\n")).
		Send(ctx)
}

func (AICommands) imaginemodels(ctx *router.Ctx) {
	models := api.GetImageModels()

	modelsString := []string{}

	for k := range models {
		modelsString = append(modelsString, k)
	}

	sort.Strings(modelsString)

	c.services.Embed.NewEmbed().
		SetTitle("Image Models").
		SetDescription(strings.Join(modelsString, "\n")).
		Send(ctx)
}

func (AICommands) imagine(ctx *router.Ctx) {
	prompt := utils.Arguments["<\"prompt\">"]

	style := strings.ToLower(utils.Arguments["[\"style\"]"])
	if styleValue, ok := api.DefaultStyles[style]; ok {
		prompt += ", " + styleValue
	}

	var model string

	if modelValue, ok := api.DefaultModels[utils.Arguments["[\"model\"]"]]; ok {
		model = modelValue
	}

	image, err := api.GenerateImage(prompt, utils.Arguments["[\"negative prompt\"]"], model)
	if err != nil {
		_ = ctx.RespondText("Error getting response: " + err.Error())
		return
	}

	_, _ = ctx.Session.ChannelFileSend(ctx.Message.ChannelID, "image.png", image)
}

func (AICommands) ocrlink(ctx *router.Ctx) {
	text, err := api.OCR(ctx.Arguments.Raw())
	if err != nil {
		_ = ctx.RespondText("Error getting response: " + err.Error())
		return
	}

	c.services.Embed.NewEmbed().
		SetTitle("OCR").
		AddField("URL", ctx.Arguments.Raw()).
		AddField("Extracted Text", text).
		Send(ctx)
}

func (AICommands) ocrfile(ctx *router.Ctx) {
	attachments := ctx.Message.Message.Attachments

	if len(attachments) == 0 {
		_ = ctx.RespondText("Error: you must provide an image as an attachment.")
		return
	}

	text, err := api.OCR(attachments[0].URL)
	if err != nil {
		_ = ctx.RespondText("Error getting response: " + err.Error())
		return
	}

	c.services.Embed.NewEmbed().
		SetTitle("OCR").
		AddField("URL", attachments[0].URL).
		AddField("Extracted Text", text).
		Send(ctx)
}
