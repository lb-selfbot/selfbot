package http

import (
	"bytes"
	"encoding/json"
	"io"
	"strings"

	"github.com/LightningDev1/discordgo"
	"github.com/bitly/go-simplejson"
	fhttp "github.com/bogdanfinn/fhttp"
	tls_client "github.com/bogdanfinn/tls-client"
	"github.com/bogdanfinn/tls-client/profiles"
)

var tlsClientOptions = []tls_client.HttpClientOption{
	tls_client.WithTimeoutSeconds(120),
	tls_client.WithClientProfile(profiles.Chrome_124),
	tls_client.WithCookieJar(tls_client.NewCookieJar()),
	tls_client.WithRandomTLSExtensionOrder(),
	tls_client.WithNotFollowRedirects(),
}

type HttpResponse struct {
	*fhttp.Response

	Error error
	Body  []byte
}

func (h *HttpResponse) GetHeaders() (headers map[string]string) {
	headers = make(map[string]string)

	for key, value := range h.Header {
		headers[key] = strings.Join(value, ", ")
	}

	return
}

func (h *HttpResponse) Json() *simplejson.Json {
	jsonData, err := simplejson.NewJson(h.Body)
	if err != nil {
		return simplejson.New()
	}

	return jsonData
}

func errorResponse(err error) *HttpResponse {
	return &HttpResponse{Error: err}
}

func Do(req *fhttp.Request) *HttpResponse {
	client, err := tls_client.NewHttpClient(tls_client.NewNoopLogger(), tlsClientOptions...)
	if err != nil {
		return errorResponse(err)
	}

	if req.Header.Get("User-Agent") == "" {
		req.Header.Set("User-Agent", discordgo.UserAgentWeb)
	}

	resp, err := client.Do(req)
	if err != nil {
		return errorResponse(err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return errorResponse(err)
	}

	return &HttpResponse{
		Response: resp,
		Error:    nil,
		Body:     body,
	}
}

func Get(url string, headersList ...map[string]string) *HttpResponse {
	req, err := fhttp.NewRequest("GET", url, nil)
	if err != nil {
		return errorResponse(err)
	}

	for _, headers := range headersList {
		for k, v := range headers {
			req.Header.Add(k, v)
		}
	}

	return Do(req)
}

func GetJson(url string, headersList ...map[string]string) (*simplejson.Json, error) {
	httpResponse := Get(url, headersList...)
	return httpResponse.Json(), httpResponse.Error
}

func CustomMethodWithData(url, method string, data []byte, headersList ...map[string]string) *HttpResponse {
	req, err := fhttp.NewRequest(method, url, bytes.NewReader(data))
	if err != nil {
		return errorResponse(err)
	}

	for _, headers := range headersList {
		for k, v := range headers {
			req.Header.Add(k, v)
		}
	}

	return Do(req)
}

func Post(url string, data any, headersList ...map[string]string) *HttpResponse {
	var rawData []byte

	switch data := data.(type) {
	case []byte:
		rawData = data
	default:
		var err error
		rawData, err = json.Marshal(data)
		if err != nil {
			return errorResponse(err)
		}

		headersList = append(headersList, map[string]string{
			"Content-Type": "application/json",
		})
	}

	return CustomMethodWithData(url, "POST", rawData, headersList...)
}

func PostJson(url string, data any, headersList ...map[string]string) (*simplejson.Json, error) {
	httpResponse := Post(url, data, headersList...)
	return httpResponse.Json(), httpResponse.Error
}

func Put(url string, data []byte, headersList ...map[string]string) *HttpResponse {
	return CustomMethodWithData(url, "PUT", data, headersList...)
}

func PutJson(url string, data any, headersList ...map[string]string) *HttpResponse {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return errorResponse(err)
	}

	headersList = append(headersList, map[string]string{
		"Content-Type": "application/json",
	})

	return CustomMethodWithData(url, "PUT", jsonData, headersList...)
}
