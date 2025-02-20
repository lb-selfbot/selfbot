package image

import (
	"github.com/LightningDev1/LB-Selfbot/internal/services"
	"github.com/LightningDev1/dgc"
)

type ImageCommands struct {
	services *services.Services
}

func NewImageCommands(services *services.Services) *ImageCommands {
	return &ImageCommands{services: services}
}

func (c *ImageCommands) Register() {
	router := c.services.Discord.GetRouter()
	router.StartCategory("Image", "Image manipulation commands")

	router.RegisterCmd(&dgc.Command{
		Name:        "upload",
		Description: "Upload an image to specified host",
		Usage:       "[p]upload <url/attachment> [host]",
		Handler:     c.upload,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "qrcode",
		Description: "Generate a QR code",
		Usage:       "[p]qrcode <text>",
		Handler:     c.qrcode,
	})

	// ... Register remaining image commands ...
}

func (c *ImageCommands) upload(ctx *dgc.Ctx) {
	// ... Existing implementation ...
}
