package http

import (
	"bytes"
	"io"
	"unsafe"

	fhttp "github.com/bogdanfinn/fhttp"
	tls_client "github.com/bogdanfinn/tls-client"
	"github.com/bogdanfinn/tls-client/profiles"
	"github.com/goccy/go-json"
	"github.com/lb-selfbot/discordgo"
)

var tlsClientOptions = []tls_client.HttpClientOption{
	tls_client.WithTimeoutSeconds(10),
	tls_client.WithClientProfile(profiles.Chrome_133),
	tls_client.WithCookieJar(tls_client.NewCookieJar()),
	tls_client.WithRandomTLSExtensionOrder(),
	tls_client.WithNotFollowRedirects(),
}

var defaultClient *Client

type Client struct {
	client tls_client.HttpClient
}

func NewClient() (*Client, error) {
	client, err := tls_client.NewHttpClient(tls_client.NewNoopLogger(), tlsClientOptions...)
	if err != nil {
		return nil, err
	}

	return &Client{client: client}, nil
}

func DefaultClient() *Client {
	if defaultClient == nil {
		defaultClient, _ = NewClient()
	}

	return defaultClient
}

func (c *Client) Do(req *fhttp.Request) (*Response, error) {
	if req.Header.Get("User-Agent") == "" {
		req.Header.Set("User-Agent", discordgo.UserAgentWeb)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	return NewResponse(resp), nil
}

func (c *Client) applyReqData(req *fhttp.Request, data ...RequestData) {
	for _, d := range data {
		for k, v := range d.Headers {
			req.Header.Set(k, v)
		}
	}
}

func (c *Client) dataReaderFromIface(data any) (io.Reader, string, error) {
	if data == nil {
		return nil, "", nil
	}

	switch data := data.(type) {
	case string:
		if len(data) == 0 {
			return nil, "", nil
		}
		strData := unsafe.StringData(data)
		if strData == nil {
			return nil, "", nil
		}
		return bytes.NewReader(unsafe.Slice(strData, len(data))), "text/plain", nil

	case []byte:
		if len(data) == 0 {
			return nil, "", nil
		}
		return bytes.NewReader(data), "text/plain", nil

	case *bytes.Reader:

		return data, "text/plain", nil

	default:
		dataMarshaled, err := json.Marshal(data)
		if err != nil {
			return nil, "", err
		}
		return bytes.NewReader(dataMarshaled), "application/json", nil
	}
}

func (c *Client) Request(method, url string, data any, dataSlice ...RequestData) (*Response, error) {
	dataReader, contentType, err := c.dataReaderFromIface(data)
	if err != nil {
		return nil, err
	}

	req, err := fhttp.NewRequest(method, url, dataReader)
	if err != nil {
		return nil, err
	}

	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}

	c.applyReqData(req, dataSlice...)

	return c.Do(req)
}

func (c *Client) Get(url string, dataSlice ...RequestData) (*Response, error) {
	return c.Request(fhttp.MethodGet, url, nil, dataSlice...)
}

func (c *Client) Head(url string, dataSlice ...RequestData) (*Response, error) {
	return c.Request(fhttp.MethodHead, url, nil, dataSlice...)
}

func (c *Client) Post(url string, data any, dataSlice ...RequestData) (*Response, error) {
	return c.Request(fhttp.MethodPost, url, data, dataSlice...)
}

func (c *Client) Put(url string, data any, dataSlice ...RequestData) (*Response, error) {
	return c.Request(fhttp.MethodPut, url, data, dataSlice...)
}

func (c *Client) Patch(url string, data any, dataSlice ...RequestData) (*Response, error) {
	return c.Request(fhttp.MethodPatch, url, data, dataSlice...)
}

func (c *Client) Delete(url string, data any, dataSlice ...RequestData) (*Response, error) {
	return c.Request(fhttp.MethodDelete, url, data, dataSlice...)
}
