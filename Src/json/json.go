package json

import (
	"encoding/json"
)

func Escape(i string) string {
	b, err := json.Marshal(i)
	if err != nil {
		return ""
	}
	// Trim the beginning and trailing " character
	return string(b[1 : len(b)-1])
}
