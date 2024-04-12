package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/golang-jwt/jwt/v5"
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

// JWT TOKEK
var secretKey = []byte("secret-key")

func CreateToken(email string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"email": email,
			"exp":   time.Now().Add(time.Hour * 24).Unix(),
		})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func VerifyToken(tokenString string) (string, error) {
	// Parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return "", fmt.Errorf("invalid token signature")
		} else {
			return "", fmt.Errorf("invalid token: %w", err)
		}
	}

	// Check token validity
	if !token.Valid {
		return "", fmt.Errorf("invalid token")
	}
	// Extract  claim
	claims := token.Claims.(jwt.MapClaims)
	email, ok := claims["email"].(string)
	if !ok {
		return "", fmt.Errorf("missing email claim")
	}

	return email, nil
}
