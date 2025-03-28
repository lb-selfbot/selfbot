package suleu

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/lb-selfbot/go-jester"
	"github.com/lb-selfbot/selfbot/internal/uploaders"
	"github.com/lb-selfbot/selfbot/internal/uploaders/keys"
)

type SulEU struct{}

func (SulEU) SupportsAllFiles() bool {
	return true
}

func (s SulEU) Upload(file *bytes.Buffer, filename string) (*uploaders.UploadResult, error) {
	buf := new(bytes.Buffer)
	writer := multipart.NewWriter(buf)

	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return nil, err
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return nil, err
	}

	_ = writer.WriteField("wizard", "true")
	_ = writer.WriteField("key", keys.SulEUToken)

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", "https://s-ul.eu/api/v1/upload", buf)
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

	json, err := jester.NewReader(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("invalid s-ul.eu response: %s", json.MustMap())
	}

	url := json.Get("url").MustString()
	if url == "" {
		return nil, fmt.Errorf("invalid s-ul.eu response: %s", json.MustMap())
	}
	result := &uploaders.UploadResult{
		URL:       url,
		DirectURL: url,
	}

	return result, nil
}
