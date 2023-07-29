package sxcu

import (
	"bytes"
	_json "encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/LightningDev1/LB-Selfbot/uploaders"
	"github.com/LightningDev1/LB-Selfbot/utils"
	"github.com/bitly/go-simplejson"
)

var ErrSXCUFileDoesNotExist = fmt.Errorf("sxcu file does not exist! please put your sxcu file in LightningBot/uploader.sxcu")

type sxcuData struct {
	RequestURL   string `json:"RequestURL"`
	FileFormName string `json:"FileFormName"`
	URL          string `json:"URL"`

	Arguments map[string]string
}

func (s *sxcuData) UnmarshalJSON(data []byte) error {
	type Alias sxcuData
	aux := &struct {
		Arguments any `json:"Arguments"`
		*Alias
	}{Alias: (*Alias)(s)}

	if err := _json.Unmarshal(data, &aux); err != nil {
		return err
	}

	if m, ok := aux.Arguments.(map[string]any); ok {
		s.Arguments = make(map[string]string)

		for k, v := range m {
			s.Arguments[k] = fmt.Sprint(v)
		}
	}

	return nil
}

type SXCU struct{}

func (SXCU) SupportsAllFiles() bool {
	return false
}

func (SXCU) getSXCU() (data *sxcuData, err error) {
	filePath := path.Join(utils.IO.GetMainFolder(), "uploader.sxcu")

	if !utils.IO.Exists(filePath) {
		return nil, ErrSXCUFileDoesNotExist
	}

	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	err = _json.NewDecoder(f).Decode(&data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (s SXCU) Upload(file *bytes.Buffer, filename string) (*uploaders.UploadResult, error) {
	data, err := s.getSXCU()
	if err != nil {
		return nil, err
	}

	var buf = new(bytes.Buffer)
	writer := multipart.NewWriter(buf)

	part, err := writer.CreateFormFile(data.FileFormName, filename)
	if err != nil {
		return nil, err
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return nil, err
	}

	for k, v := range data.Arguments {
		_ = writer.WriteField(k, v)
	}

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", data.RequestURL, buf)
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
		return nil, fmt.Errorf("invalid response: %s", json.MustMap())
	}

	data.URL = strings.ReplaceAll(data.URL, "$json:", "")
	data.URL = strings.ReplaceAll(data.URL, "$", "")

	link := json.GetPath(strings.Split(data.URL, ".")...).MustString()

	if link == "" {
		return nil, fmt.Errorf("invalid response: %s", json.MustMap())
	}

	result := &uploaders.UploadResult{
		URL:       link,
		DirectURL: link,
	}

	return result, nil
}
