package utils

import (
	"encoding/json"
)

// ToJSON converts an object to JSON string
func ToJSON(v interface{}) (string, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// ToJSONBytes converts an object to JSON bytes
func ToJSONBytes(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

// FromJSON parses JSON string into an object
func FromJSON(data string, v interface{}) error {
	return json.Unmarshal([]byte(data), v)
}

// FromJSONBytes parses JSON bytes into an object
func FromJSONBytes(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

// PrettyJSON returns pretty-printed JSON string
func PrettyJSON(v interface{}) (string, error) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b), nil
}
