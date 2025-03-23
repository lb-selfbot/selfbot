package catbox

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/lb-selfbot/selfbot/internal/uploaders"
)

type Catbox struct{}

func (Catbox) SupportsAllFiles() bool {
	return false
}

func (c Catbox) Upload(file *bytes.Buffer, filename string) (*uploaders.UploadResult, error) {
	buf := new(bytes.Buffer)
	writer := multipart.NewWriter(buf)

	part, err := writer.CreateFormFile("fileToUpload", filename)
	if err != nil {
		return nil, err
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return nil, err
	}

	_ = writer.WriteField("reqtype", "fileupload")

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", "https://catbox.moe/user/api.php", buf)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	result := &uploaders.UploadResult{
		URL:       string(body),
		DirectURL: string(body),
	}

	return result, nil
}
