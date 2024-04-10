package utils

import (
	"encoding/json"
	"errors"
	"io"

	"github.com/manish-mehra/go-todo/models"
)

// DecodeUserJSON decodes JSON data into a USER struct
func DecodeUserJSON(data io.Reader) (models.User, error) {
	decoder := json.NewDecoder(data)
	var user models.User
	if err := decoder.Decode(&user); err != nil {
		return models.User{}, errors.New("Error decoding user JSON: " + err.Error())
	}
	return user, nil
}

func ParseResponse(message string) ([]byte, error) {
	response := models.Response{Message: message}
	jsonData, err := json.Marshal(response)
	if err != nil {
		return []byte{}, errors.New("error parsing response")
	}
	return jsonData, nil
}
