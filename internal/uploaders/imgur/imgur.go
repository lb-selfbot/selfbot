package imgur

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path"
	"strings"

	"github.com/lb-selfbot/go-jester"
	"github.com/lb-selfbot/selfbot/internal/uploaders"
	"github.com/lb-selfbot/selfbot/internal/uploaders/keys"
)

type Imgur struct{}

func (Imgur) SupportsAllFiles() bool {
	return false
}

func (Imgur) Upload(file *bytes.Buffer, filename string) (*uploaders.UploadResult, error) {
	buf := new(bytes.Buffer)
	writer := multipart.NewWriter(buf)

	part, err := writer.CreateFormFile("image", filename)
	if err != nil {
		return nil, err
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return nil, err
	}

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", "https://api.imgur.com/3/image", buf)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Client-ID "+keys.ImgurClientID)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := jester.NewReader(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("invalid Imgur response: %s", data.MustMap())
	}

	if !data.Get("success").MustBool() || data.Get("status").MustInt() != 200 {
		return nil, fmt.Errorf("invalid Imgur response: %s", data.MustMap())
	}

	directUrl := data.Get("data", "link").MustString()
	if directUrl == "" || !strings.Contains(directUrl, ".") {
		return nil, fmt.Errorf("invalid Imgur response: %s", data.MustMap())
	}

	ext := path.Ext(directUrl)
	result := &uploaders.UploadResult{
		URL:       directUrl[:len(directUrl)-len(ext)],
		DirectURL: directUrl,
	}

	return result, nil
}
