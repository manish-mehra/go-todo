package utils

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/manish-mehra/go-todo/models"
)

type ErrMessages struct{}

func (e ErrMessages) NotFound() error {
	return errors.New("not found")
}

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
var secretKey = []byte(os.Getenv("JWT_SECRET"))

func CreateToken(userId string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"userId": userId,
			"exp":    time.Now().Add(time.Hour * 24).Unix(),
		})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// VerifyToken parses and validates a JWT token.
// - Extracts user ID from valid claims.
// - Returns userID and error if invalid or missing claim.
//   - Specific error messages for signature issues and missing claims.
func VerifyToken(tokenString string) (string, error) {
	// Parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return "", errors.New("invalid token signature")
		} else {
			return "", errors.New("invalid token")
		}
	}

	// Check token validity
	if !token.Valid {
		return "", errors.New("invalid token")
	}
	// Extract  claim
	claims := token.Claims.(jwt.MapClaims)
	userId, ok := claims["userId"].(string)
	if !ok {
		return "", errors.New("missing userId  claim")
	}

	return userId, nil
}

// convert strint to int
func StringToInt(str string) (int, error) {
	i, err := strconv.Atoi(str)
	if err != nil {
		return 0, err // Handle error (e.g., return 0 and an error)
	}
	return int(i), nil
}

// Not Found
var (
	ErrNotFound = errors.New("not found")
)
