package destinations

import (
	"bytes"

	"github.com/LightningDev1/LB-Selfbot/selfbot/uploaders"
	"github.com/LightningDev1/LB-Selfbot/selfbot/uploaders/anonfiles"
	"github.com/LightningDev1/LB-Selfbot/selfbot/uploaders/catbox"
	"github.com/LightningDev1/LB-Selfbot/selfbot/uploaders/imgbb"
	"github.com/LightningDev1/LB-Selfbot/selfbot/uploaders/imgur"
	"github.com/LightningDev1/LB-Selfbot/selfbot/uploaders/suleu"
	"github.com/LightningDev1/LB-Selfbot/selfbot/uploaders/sxcu"
	"github.com/LightningDev1/LB-Selfbot/selfbot/uploaders/vgyme"
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
