package json

import (
	"encoding/json"

	"github.com/bitly/go-simplejson"
)

func UnmarshalJson(jsonString string) (map[string]interface{}, error) {
	var jsonMap map[string]interface{}
	err := json.Unmarshal([]byte(jsonString), &jsonMap)
	return jsonMap, err
}

func SimpleUnmarshal(jsonString string) (*simplejson.Json, error) {
	return simplejson.NewJson([]byte(jsonString))
}
