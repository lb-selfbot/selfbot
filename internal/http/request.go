package http

// RequestData contains the information for an HTTP request.
type RequestData struct {
	Headers map[string]string
}

// WithHeaders creates a new RequestData with the specified headers.
func WithHeaders(headers map[string]string) RequestData {
	return RequestData{
		Headers: headers,
	}
}

// WithHeader creates a new RequestData with the specified header.
func WithHeader(key, value string) RequestData {
	return WithHeaders(map[string]string{
		key: value,
	})
}
