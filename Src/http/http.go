package http

import (
	"github.com/LightningDev1/LB-Selfbot/json"
	"github.com/bitly/go-simplejson"
	"github.com/valyala/fasthttp"
)

func GetBytes(url string, headersList ...map[string]string) ([]byte, int, error) {
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

func Get(url string, headersList ...map[string]string) (string, error) {
	bytes, _, err := GetBytes(url, headersList...)
	return string(bytes), err
}

func GetJson(url string, headersList ...map[string]string) (*simplejson.Json, error) {
	bytes, _, err := GetBytes(url, headersList...)
	if err != nil {
		return nil, err
	}
	return json.SimpleUnmarshal(string(bytes))
}

func GetJsonStatusCode(url string, headersList ...map[string]string) (*simplejson.Json, int, error) {
	bytes, statusCode, err := GetBytes(url, headersList...)
	if err != nil {
		return nil, -1, err
	}
	jsonData, err := json.SimpleUnmarshal(string(bytes))
	return jsonData, statusCode, err
}

func GetHeadersString(url string) (string, error) {
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)

	req.SetRequestURI(url)
	req.Header.SetMethod("HEAD")
	err := fasthttp.Do(req, resp)
	if err != nil {
		return "", err
	}
	return string(resp.Header.Header()), nil

}

func Post(url string, data string, headersList ...map[string]string) (string, error) {
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

func PostJson(url string, data string, headersList ...map[string]string) (*simplejson.Json, error) {
	headersList = append(headersList, map[string]string{"Content-Type": "application/json"})
	result, err := Post(url, data, headersList...)
	if err != nil {
		return nil, err
	}
	return json.SimpleUnmarshal(result)
}

func Patch(url string, data string, headersList ...map[string]string) (string, error) {
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

func PatchJson(url string, data string, headersList ...map[string]string) (*simplejson.Json, error) {
	headersList = append(headersList, map[string]string{"Content-Type": "application/json"})
	result, err := Patch(url, data, headersList...)
	if err != nil {
		return nil, err
	}
	return json.SimpleUnmarshal(result)
}

func Put(url string, headersList ...map[string]string) (string, error) {
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

func Delete(url string, headersList ...map[string]string) (string, error) {
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

func GetDiscordHeaders(token string) map[string]string {
	return map[string]string{
		"content-type":       "application/json",
		"authorization":      token,
	}
}
