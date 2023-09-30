package anonfiles

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/LightningDev1/LB-Selfbot/selfbot/uploaders"
	"github.com/bitly/go-simplejson"
)

type Anonfiles struct{}

func (Anonfiles) SupportsAllFiles() bool {
	return true
}

func (Anonfiles) Upload(file *bytes.Buffer, filename string) (*uploaders.UploadResult, error) {
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

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", "https://api.anonfiles.com/upload", buf)
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
		return nil, fmt.Errorf("invalid Anonfiles response: %s", json.MustMap())
	}

	if !json.Get("status").MustBool() {
		return nil, fmt.Errorf("invalid Anonfiles response: %s", json.MustMap())
	}

	result := &uploaders.UploadResult{
		URL:       json.GetPath("data", "file", "url", "full").MustString(),
		DirectURL: json.GetPath("data", "file", "url", "short").MustString(),
	}

	return result, nil
}
