package main

import (
	"encoding/json"
	"errors"
	"io"
)

// DecodeUserJSON decodes JSON data into a USER struct
func DecodeUserJSON(data io.Reader) (USER, error) {
	decoder := json.NewDecoder(data)
	var user USER
	if err := decoder.Decode(&user); err != nil {
		return USER{}, errors.New("Error decoding user JSON: " + err.Error())
	}
	return user, nil
}

func parseResponse(message string) ([]byte, error) {
	response := Response{Message: message}
	jsonData, err := json.Marshal(response)
	if err != nil {
		return []byte{}, errors.New("error parsing response")
	}
	return jsonData, nil
}
