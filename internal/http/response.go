package http

import (
	"io"

	fhttp "github.com/bogdanfinn/fhttp"
	"github.com/lb-selfbot/go-jester"
)

type Response struct {
	*fhttp.Response

	jsonData *jester.Data
}

func NewResponse(resp *fhttp.Response) *Response {
	return &Response{Response: resp}
}

func (r *Response) IsSuccess() bool {
	return r.StatusCode >= 200 && r.StatusCode < 300
}

func (r *Response) GetBody() ([]byte, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (r *Response) MustGetBody() []byte {
	body, err := r.GetBody()
	if err != nil {
		return nil
	}

	return body
}

func (r *Response) Json() *jester.Data {
	if r.jsonData != nil {
		return r.jsonData
	}

	data, err := jester.NewReader(r.Body)
	if err != nil {
		return jester.NewEmpty()
	}

	r.jsonData = data
	return data
}
