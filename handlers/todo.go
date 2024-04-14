package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/manish-mehra/go-todo/models"
	"github.com/manish-mehra/go-todo/services"
	"github.com/manish-mehra/go-todo/utils"
)

func isAuthenticated(req *http.Request) (string, error) {
	// get the token from request cookie
	cookies := req.Cookies()
	var token string
	for _, cookie := range cookies {
		if cookie.Name == "token" {
			token = cookie.Value
		}
	}

	if token == "" {
		return "", errors.New("No jwt found")
	}

	// verify token
	userId, err := utils.VerifyToken(token)
	if err != nil {
		return "", err
	}

	return userId, nil
}

func GetTodo(w http.ResponseWriter, req *http.Request)    {}
func GetAllTodo(w http.ResponseWriter, req *http.Request) {}

func PostTodo(w http.ResponseWriter, req *http.Request) {
	// check if authenticated
	userId, err := isAuthenticated(req)
	if err != nil {
		log.Print(err)
		message := err.Error()
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, message)
		return
	}
	// get todo from response
	var newTodo models.UserTodo
	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&newTodo); err != nil {
		fmt.Fprintf(w, "error decoding json data")
		return
	}

	// convert user id to int
	uId, err := utils.StringToInt64(userId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, err.Error())
		return
	}

	err = services.PostTodo(newTodo, uId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, err.Error())
		return
	}

	message := "todo created"
	response, _ := utils.ParseResponse(message)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
	return
}

func DeleteTodo(w http.ResponseWriter, req *http.Request) {}
