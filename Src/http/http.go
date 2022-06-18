package http

import (
	"github.com/LightningDev/LB-Selfbot/json"
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

func GetJson(url string, headersList ...map[string]string) (map[string]interface{}, error) {
	bytes, _, err := GetBytes(url, headersList...)
	if err != nil {
		return nil, err
	}
	return json.UnmarshalJson(string(bytes))
}

func GetSimpleJson(url string, headersList ...map[string]string) (*simplejson.Json, error) {
	bytes, _, err := GetBytes(url, headersList...)
	if err != nil {
		return nil, err
	}
	return json.SimpleUnmarshal(string(bytes))
}

func GetSimpleJsonStatusCode(url string, headersList ...map[string]string) (*simplejson.Json, int, error) {
	bytes, statusCode, err := GetBytes(url, headersList...)
	if err != nil {
		return nil, -1, err
	}
	jsonData, err := json.SimpleUnmarshal(string(bytes))
	return jsonData, statusCode, err
}

func GetJsonStatusCode(url string, headersList ...map[string]string) (map[string]interface{}, int, error) {
	bytes, statusCode, err := GetBytes(url, headersList...)
	if err != nil {
		return nil, -1, err
	}
	jsonData, err := json.UnmarshalJson(string(bytes))
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

func PostJson(url string, data string, headersList ...map[string]string) (map[string]interface{}, error) {
	headersList = append(headersList, map[string]string{"Content-Type": "application/json"})
	result, err := Post(url, data, headersList...)
	if err != nil {
		return nil, err
	}
	return json.UnmarshalJson(result)
}

func PostSimpleJson(url string, data string, headersList ...map[string]string) (*simplejson.Json, error) {
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
		"accept":             "*/*",
		"accept-encoding":    "",
		"accept-language":    "en-GB,q=0.9",
		"cookie":             "__cfduid=d7e8d2784592da39fb3f621664b9aede51620414171; __dcfduid=24a543339247480f9b0bb95c710ce1e6",
		"referer":            "https://discord.com/",
		"origin":             "https://discord.com",
		"sec-fetch-dest":     "empty",
		"sec-fetch-mode":     "cors",
		"sec-fetch-site":     "same-origin",
		"user-agent":         "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) discord/1.0.9003 Chrome/91.0.4472.164 Electron/13.4.0 Safari/537.36",
		"x-debug-options":    "bugReporterEnabled",
		"x-discord-locale":   "en-US",
		"x-super-properties": "eyJvcyI6IldpbmRvd3MiLCJicm93c2VyIjoiRGlzY29yZCBDbGllbnQiLCJyZWxlYXNlX2NoYW5uZWwiOiJzdGFibGUiLCJjbGllbnRfdmVyc2lvbiI6IjEuMC45MDAzIiwib3NfdmVyc2lvbiI6IjEwLjAuMjIwMDAiLCJvc19hcmNoIjoieDY0Iiwic3lzdGVtX2xvY2FsZSI6ImVuLUdCIiwiY2xpZW50X2J1aWxkX251bWJlciI6MTA5MTkwLCJjbGllbnRfZXZlbnRfc291cmNlIjpudWxsfQ==",
		"content-type":       "application/json",
		"authorization":      token,
	}
}
