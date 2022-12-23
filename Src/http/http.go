package http

import (
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

func GetOld(url string, headersList ...map[string]string) (string, error) {
	bytes, _, err := GetBytesOld(url, headersList...)
	return string(bytes), err
}

func GetJsonOld(url string, headersList ...map[string]string) (*simplejson.Json, error) {
	bytes, _, err := GetBytesOld(url, headersList...)
	if err != nil {
		return nil, err
	}
	return json.SimpleUnmarshal(string(bytes))
}

func GetJsonStatusCodeOld(url string, headersList ...map[string]string) (*simplejson.Json, int, error) {
	bytes, statusCode, err := GetBytesOld(url, headersList...)
	if err != nil {
		return nil, -1, err
	}
	jsonData, err := json.SimpleUnmarshal(string(bytes))
	return jsonData, statusCode, err
}

func GetHeadersStringOld(url string) (string, error) {
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
