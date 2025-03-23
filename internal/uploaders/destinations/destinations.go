package destinations

import (
	"bytes"

	"github.com/LightningDev1/LB-Selfbot/internal/uploaders"
	"github.com/LightningDev1/LB-Selfbot/internal/uploaders/catbox"
	"github.com/LightningDev1/LB-Selfbot/internal/uploaders/imgbb"
	"github.com/LightningDev1/LB-Selfbot/internal/uploaders/imgur"
	"github.com/LightningDev1/LB-Selfbot/internal/uploaders/suleu"
	"github.com/LightningDev1/LB-Selfbot/internal/uploaders/sxcu"
	"github.com/LightningDev1/LB-Selfbot/internal/uploaders/vgyme"
)

type Destination interface {
	SupportsAllFiles() bool
	Upload(file *bytes.Buffer, filename string) (*uploaders.UploadResult, error)
}

var Destinations = map[string]Destination{
	"imgur":  imgur.Imgur{},
	"vgyme":  vgyme.Vgyme{},
	"imgbb":  imgbb.ImgBB{},
	"suleu":  suleu.SulEU{},
	"sxcu":   sxcu.SXCU{},
	"catbox": catbox.Catbox{},
}
