package ai

import (
	"github.com/lb-selfbot/selfbot/internal/api"
	"github.com/lb-selfbot/selfbot/internal/services/router"

	"github.com/lb-selfbot/selfbot/internal/services"
)

type AICommands struct{ services *services.Services }

func NewAICommands(s *services.Services) *AICommands {
	return &AICommands{services: s}
}

func (c *AICommands) Register() {
	c.services.Router.RegisterCategory(&router.Category{
		Name:        "AI",
		Description: "Interact with AI's",
	}, []*router.Command{
		{
			Name:        "deepseek",
			Description: "Interact with Deepseek-V3",
			Arguments:   "<\"message\"> [\"system\"]",
			Handler:     c.deepseek,
		},
		{
			Name:        "ocrlink",
			Description: "Extract text from an image URL",
			Arguments:   "<link>",
			Handler:     c.ocrlink,
		},
		{
			Name:        "ocrfile",
			Description: "Extract text from an image file",
			Arguments:   "[attachment]",
			Handler:     c.ocrfile,
		},
	})
}

func (c *AICommands) getCreds() api.LBCreds {
	cfg := c.services.Config.GetConfig()
	return api.LBCreds{
		Username: cfg.Account.LightningBot.Username,
		Password: cfg.Account.LightningBot.Password,
	}
}

func (c *AICommands) deepseek(ctx *router.Ctx) {
	message, system := ctx.Arguments.Get(0).Raw(), ctx.Arguments.Get(1).Raw()

	response, err := api.LBGPT(c.getCreds(), system, message)
	if err != nil {
		ctx.RespondText("Error getting response: %v", err)
		return
	}

	c.services.Embed.New(ctx).
		SetTitle("Deepseek V3").
		AddField("System", system).
		AddField("Message", message).
		AddField("Response", response).
		Send()
}

func (c *AICommands) ocrlink(ctx *router.Ctx) {
	text, err := api.LBOCR(c.getCreds(), ctx.Arguments.Raw())
	if err != nil {
		ctx.RespondText("Error getting response: %v", err)
		return
	}

	c.services.Embed.New(ctx).
		SetTitle("OCR").
		AddField("URL", ctx.Arguments.Raw()).
		AddField("Extracted Text", text).
		Send()
}

func (c *AICommands) ocrfile(ctx *router.Ctx) {
	attachments := ctx.Message.Attachments

	if len(attachments) == 0 {
		ctx.RespondText("Error: you must provide an image as an attachment.")
		return
	}

	text, err := api.LBOCR(c.getCreds(), attachments[0].URL)
	if err != nil {
		ctx.RespondText("Error getting response: %v", err)
		return
	}

	c.services.Embed.New(ctx).
		SetTitle("OCR").
		AddField("URL", attachments[0].URL).
		AddField("Extracted Text", text).
		Send()
}
