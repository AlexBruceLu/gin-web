package util

import "encoding/json"

func JsonEncode(v interface{}) (string, error) {
	json, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(json), nil
}
