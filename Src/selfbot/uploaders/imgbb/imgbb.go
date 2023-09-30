package imgbb

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/LightningDev1/LB-Selfbot/selfbot/uploaders"
	"github.com/LightningDev1/LB-Selfbot/selfbot/uploaders/keys"
	"github.com/bitly/go-simplejson"
)

type ImgBB struct{}

func (ImgBB) SupportsAllFiles() bool {
	return false
}

func (i ImgBB) Upload(file *bytes.Buffer, filename string) (*uploaders.UploadResult, error) {
	var buf = new(bytes.Buffer)
	writer := multipart.NewWriter(buf)

	part, err := writer.CreateFormFile("source", filename)
	if err != nil {
		return nil, err
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return nil, err
	}

	_ = writer.WriteField("type", "file")
	_ = writer.WriteField("action", "upload")
	_ = writer.WriteField("timestamp", fmt.Sprint((time.Now().Second())))
	_ = writer.WriteField("auth_token", keys.ImgBBAuthToken)

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", "https://imgbb.com/json", buf)
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

	if resp.StatusCode != 200 || json.Get("status_code").MustInt() != 200 {
		return nil, fmt.Errorf("invalid ImgBB response: %s", json.MustMap())
	}

	result := &uploaders.UploadResult{
		URL:       json.GetPath("image", "url_viewer").MustString(),
		DirectURL: json.GetPath("image", "url").MustString(),
	}

	return result, nil
}
