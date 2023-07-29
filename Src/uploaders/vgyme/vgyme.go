package vgyme

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/LightningDev1/LB-Selfbot/uploaders"
	"github.com/LightningDev1/LB-Selfbot/uploaders/keys"
	"github.com/bitly/go-simplejson"
)

type Vgyme struct{}

func (Vgyme) SupportsAllFiles() bool {
	return false
}

func (Vgyme) Upload(file *bytes.Buffer, filename string) (*uploaders.UploadResult, error) {
	var buf = new(bytes.Buffer)
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

	json, err := simplejson.NewFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("invalid vgy.me response: %s", json.MustMap())
	}

	if json.Get("error").MustBool() {
		return nil, fmt.Errorf("invalid vgy.me response: %s", json.MustMap())
	}

	result := &uploaders.UploadResult{
		URL:       json.Get("url").MustString(),
		DirectURL: json.Get("image").MustString(),
	}

	return result, nil
}
