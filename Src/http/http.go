package http

import (
	_json "encoding/json"

	"github.com/LightningDev1/LB-Selfbot/json"
	"github.com/bitly/go-simplejson"
	"github.com/valyala/fasthttp"
)

func GetBytesOld(url string, headersList ...map[string]string) ([]byte, int, error) {
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
		return []byte{}, -1, err
	}
	return resp.Body(), resp.StatusCode(), nil
}

func GetJsonOld(url string, headersList ...map[string]string) (*simplejson.Json, error) {
	bytes, _, err := GetBytesOld(url, headersList...)
	if err != nil {
		return nil, err
	}
	return json.SimpleUnmarshal(string(bytes))
}

func PostOld(url string, data string, headersList ...map[string]string) (string, error) {
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)

	req.SetRequestURI(url)
	req.Header.SetMethod("POST")
	req.SetBody([]byte(data))
	for _, headers := range headersList {
		for k, v := range headers {
			req.Header.Add(k, v)
		}
	}

	err := fasthttp.Do(req, resp)
	if err != nil {
		return "", err
	}
	return string(resp.Body()), nil
}

func PostJsonOld(url string, data string, headersList ...map[string]string) (*simplejson.Json, error) {
	headersList = append(headersList, map[string]string{"Content-Type": "application/json"})
	result, err := PostOld(url, data, headersList...)
	if err != nil {
		return nil, err
	}
	return json.SimpleUnmarshal(result)
}

func PatchOld(url string, data string, headersList ...map[string]string) (string, error) {
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)

	req.SetRequestURI(url)
	req.Header.SetMethod("PATCH")
	req.SetBody([]byte(data))
	for _, headers := range headersList {
		for k, v := range headers {
			req.Header.Add(k, v)
		}
	}

	err := fasthttp.Do(req, resp)
	if err != nil {
		return "", err
	}
	return string(resp.Body()), nil
}

func PatchJsonOld(url string, data string, headersList ...map[string]string) (*simplejson.Json, error) {
	headersList = append(headersList, map[string]string{"Content-Type": "application/json"})
	result, err := PatchOld(url, data, headersList...)
	if err != nil {
		return nil, err
	}
	return json.SimpleUnmarshal(result)
}

func PutOld(url string, headersList ...map[string]string) (string, error) {
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)

	req.SetRequestURI(url)
	req.Header.SetMethod("PUT")
	for _, headers := range headersList {
		for k, v := range headers {
			req.Header.Add(k, v)
		}
	}

	err := fasthttp.Do(req, resp)
	if err != nil {
		return "", err
	}
	return string(resp.Body()), nil
}

func DeleteOld(url string, headersList ...map[string]string) (string, error) {
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)

	req.SetRequestURI(url)
	req.Header.SetMethod("DELETE")
	for _, headers := range headersList {
		for k, v := range headers {
			req.Header.Add(k, v)
		}
	}

	err := fasthttp.Do(req, resp)
	if err != nil {
		return "", err
	}
	return string(resp.Body()), nil
}

func GetDiscordHeadersDEPRECATED(token string) map[string]string {
	return map[string]string{
		"content-type":  "application/json",
		"authorization": token,
	}
}

type HttpResponse struct {
	Success       bool
	Error         error
	Body          string
	BodyBytes     []byte
	StatusCode    int
	Headers       string
	InnerResponse *fasthttp.Response
}

func (r *HttpResponse) Json() *simplejson.Json {
	jsonData, err := simplejson.NewJson(r.BodyBytes)
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
		Headers:       "",
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

	return &HttpResponse{
		Success:       true,
		Error:         nil,
		Body:          string(resp.Body()),
		BodyBytes:     resp.Body(),
		StatusCode:    resp.StatusCode(),
		Headers:       resp.Header.String(),
		InnerResponse: resp,
	}
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

	return &HttpResponse{
		Success:       true,
		Error:         nil,
		Body:          string(resp.Body()),
		BodyBytes:     resp.Body(),
		StatusCode:    resp.StatusCode(),
		Headers:       resp.Header.String(),
		InnerResponse: resp,
	}
}

func Post(url string, data []byte, headersList ...map[string]string) *HttpResponse {
	return CustomMethodWithData(url, "POST", data, headersList...)
}

func Put(url string, data []byte, headersList ...map[string]string) *HttpResponse {
	return CustomMethodWithData(url, "PUT", data, headersList...)
}

func PostJson(url string, data map[string]any, headersList ...map[string]string) *HttpResponse {
	dataJson, err := _json.Marshal(data)
	if err != nil {
		return errorResponse(err)
	}

	headersList = append(headersList, map[string]string{"Content-Type": "application/json"})

	return CustomMethodWithData(url, "POST", dataJson, headersList...)
}
