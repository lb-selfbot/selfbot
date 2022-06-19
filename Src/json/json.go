package json

import (
	"encoding/json"

	"github.com/bitly/go-simplejson"
)

func SimpleUnmarshal(jsonString string) (*simplejson.Json, error) {
	return simplejson.NewJson([]byte(jsonString))
}

func Escape(i string) string {
	b, err := json.Marshal(i)
	if err != nil {
		return ""
	}
	// Trim the beginning and trailing " character
	return string(b[1 : len(b)-1])
}
