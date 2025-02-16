package utils

import "encoding/json"

type JsonUtils struct{}

func (JsonUtils) Escape(i string) string {
	b, err := json.Marshal(i)
	if err != nil {
		return ""
	}

	// Trim the beginning and trailing " character
	return string(b[1 : len(b)-1])
}

var Json = JsonUtils{}
