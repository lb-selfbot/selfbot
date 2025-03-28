package sxcu

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"strings"

	_json "github.com/goccy/go-json"

	"github.com/lb-selfbot/go-jester"
	"github.com/lb-selfbot/selfbot/internal/uploaders"
	"github.com/lb-selfbot/selfbot/internal/utils"
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

	f, err := os.Open(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, ErrSXCUFileDoesNotExist
		}

		return nil, err
	}

	err = _json.NewDecoder(f).Decode(&data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (s SXCU) Upload(file *bytes.Buffer, filename string) (*uploaders.UploadResult, error) {
	sxcuData, err := s.getSXCU()
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	writer := multipart.NewWriter(buf)

	part, err := writer.CreateFormFile(sxcuData.FileFormName, filename)
	if err != nil {
		return nil, err
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return nil, err
	}

	for k, v := range sxcuData.Arguments {
		_ = writer.WriteField(k, v)
	}

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", sxcuData.RequestURL, buf)
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
		return nil, fmt.Errorf("invalid response: %s", data.MustMap())
	}

	url := strings.ReplaceAll(sxcuData.URL, "$json:", "")
	url = strings.ReplaceAll(url, "$", "")
	path := strings.Split(url, ".")

	link := data.Get(toIfaceSlice(path)...).MustString()

	if link == "" {
		return nil, fmt.Errorf("invalid response: %s", data.MustMap())
	}

	result := &uploaders.UploadResult{
		URL:       link,
		DirectURL: link,
	}

	return result, nil
}

func toIfaceSlice(arr ...any) []any {
	return arr
}
