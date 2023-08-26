package destinations

import (
	"bytes"

	"github.com/LightningDev1/LB-Selfbot/uploaders"
	"github.com/LightningDev1/LB-Selfbot/uploaders/anonfiles"
	"github.com/LightningDev1/LB-Selfbot/uploaders/catbox"
	"github.com/LightningDev1/LB-Selfbot/uploaders/imgbb"
	"github.com/LightningDev1/LB-Selfbot/uploaders/imgur"
	"github.com/LightningDev1/LB-Selfbot/uploaders/suleu"
	"github.com/LightningDev1/LB-Selfbot/uploaders/sxcu"
	"github.com/LightningDev1/LB-Selfbot/uploaders/vgyme"
)

type Destination interface {
	SupportsAllFiles() bool
	Upload(file *bytes.Buffer, filename string) (*uploaders.UploadResult, error)
}

var Destinations = map[string]Destination{
	"imgur":     imgur.Imgur{},
	"vgyme":     vgyme.Vgyme{},
	"imgbb":     imgbb.ImgBB{},
	"suleu":     suleu.SulEU{},
	"anonfiles": anonfiles.Anonfiles{},
	"sxcu":      sxcu.SXCU{},
	"catbox":    catbox.Catbox{},
}
