package vgyme

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

type Vgyme struct{}

func (Vgyme) SupportsAllFiles() bool {
	return false
}

func (Vgyme) Upload(file *bytes.Buffer, filename string) (*uploaders.UploadResult, error) {
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

	_ = writer.WriteField("userkey", keys.VgymeUserKey)

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", "https://vgy.me/upload", buf)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

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
		return nil, fmt.Errorf("invalid vgy.me response: %s", data.MustMap())
	}

	if data.Get("error").MustBool() {
		return nil, fmt.Errorf("invalid vgy.me response: %s", data.MustMap())
	}

	result := &uploaders.UploadResult{
		URL:       data.Get("url").MustString(),
		DirectURL: data.Get("image").MustString(),
	}

	return result, nil
}
