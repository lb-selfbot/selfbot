package destinations

import (
	"bytes"

	"github.com/lb-selfbot/selfbot/internal/uploaders"
	"github.com/lb-selfbot/selfbot/internal/uploaders/catbox"
	"github.com/lb-selfbot/selfbot/internal/uploaders/imgbb"
	"github.com/lb-selfbot/selfbot/internal/uploaders/imgur"
	"github.com/lb-selfbot/selfbot/internal/uploaders/suleu"
	"github.com/lb-selfbot/selfbot/internal/uploaders/sxcu"
	"github.com/lb-selfbot/selfbot/internal/uploaders/vgyme"
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
