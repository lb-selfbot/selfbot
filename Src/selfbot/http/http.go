package http

import (
	_json "encoding/json"
	"strings"

	"github.com/bitly/go-simplejson"
	"github.com/valyala/fasthttp"
)

type HttpResponse struct {
	Success       bool
	Error         error
	Body          string
	BodyBytes     []byte
	StatusCode    int
	Headers       map[string]string
	HeaderString  string
	InnerResponse *fasthttp.Response
}

func (h *HttpResponse) GetHeaders() {
	headers := make(map[string]string)

	h.InnerResponse.Header.VisitAll(func(key, value []byte) {
		keyString := strings.ToLower(string(key))

		headers[keyString] = string(value)
	})

	h.Headers = headers
	h.HeaderString = h.InnerResponse.Header.String()
}

func (h *HttpResponse) Json() *simplejson.Json {
	jsonData, err := simplejson.NewJson(h.BodyBytes)
	if err != nil {
		return simplejson.New()
	}
	return jsonData
}

func errorResponse(err error) *HttpResponse {
	return &HttpResponse{
		Success:       false,
		Error:         err,
		Body:          "",
		BodyBytes:     nil,
		StatusCode:    0,
		InnerResponse: nil,
	}
}

func Get(url string, headersList ...map[string]string) *HttpResponse {
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)

	req.SetRequestURI(url)

	for _, headers := range headersList {
		for k, v := range headers {
			req.Header.Add(k, v)
		}
	}

	err := fasthttp.Do(req, resp)
	if err != nil {
		return errorResponse(err)
	}

	response := &HttpResponse{
		Success:       true,
		Error:         nil,
		Body:          string(resp.Body()),
		BodyBytes:     resp.Body(),
		StatusCode:    resp.StatusCode(),
		InnerResponse: resp,
	}

	response.GetHeaders()

	return response
}

func GetJson(url string, headersList ...map[string]string) (*simplejson.Json, error) {
	httpResponse := Get(url, headersList...)
	if !httpResponse.Success {
		return nil, httpResponse.Error
	}

	return httpResponse.Json(), nil
}

func CustomMethodWithData(url, method string, data []byte, headersList ...map[string]string) *HttpResponse {
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)

	req.SetRequestURI(url)
	req.Header.SetMethod(method)
	req.SetBody(data)

	for _, headers := range headersList {
		for k, v := range headers {
			req.Header.Add(k, v)
		}
	}

	err := fasthttp.Do(req, resp)
	if err != nil {
		return errorResponse(err)
	}

	response := &HttpResponse{
		Success:       true,
		Error:         nil,
		Body:          string(resp.Body()),
		BodyBytes:     resp.Body(),
		StatusCode:    resp.StatusCode(),
		InnerResponse: resp,
	}

	response.GetHeaders()

	return response
}

func Post(url string, data []byte, headersList ...map[string]string) *HttpResponse {
	return CustomMethodWithData(url, "POST", data, headersList...)
}

func Put(url string, data []byte, headersList ...map[string]string) *HttpResponse {
	return CustomMethodWithData(url, "PUT", data, headersList...)
}

func PostJson(url string, data any, headersList ...map[string]string) *HttpResponse {
	dataJson, err := _json.Marshal(data)
	if err != nil {
		return errorResponse(err)
	}

	headersList = append(headersList, map[string]string{"Content-Type": "application/json"})

	return CustomMethodWithData(url, "POST", dataJson, headersList...)
}
