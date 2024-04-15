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

func GetTodo(w http.ResponseWriter, req *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	// check if authenticated
	userId, err := isAuthenticated(req)
	if err != nil {
		log.Print(err)
		message := err.Error()
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, message)
		return
	}
	// get todo id from req param
	paramTodoID := req.PathValue("id")
	log.Print("paramTodoID ", paramTodoID)
	if paramTodoID == "" {
		log.Print("no todo id was found")
		message := "no todo id was found!"
		w.WriteHeader(http.StatusNotExtended)
		fmt.Fprintf(w, message)
		return
	}

	// convert user id to int
	todoID, err := utils.StringToInt64(paramTodoID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, err.Error())
		return
	}

	userID, err := utils.StringToInt64(userId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, err.Error())
		return
	}

	// get todo from db
	todo, err := services.GetTodo(userID, todoID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, err.Error())
		return
	}

	// convert res to json
	resTodo := struct {
		Id        int64  `json:"id"`
		Title     string `json:"title"`
		Completed bool   `json:"completed"`
	}{Id: todo.ID, Title: todo.Title, Completed: todo.Completed}
	response := models.Response{Message: "successful", Data: resTodo}
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Print(err)
		fmt.Fprintf(w, "internal error while parsing todo")
		return
	}
}

func GetAllTodo(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// check if authenticated
	userId, err := isAuthenticated(req)
	if err != nil {
		log.Print(err)
		message := err.Error()
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, message)
		return
	}

	userID, err := utils.StringToInt64(userId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, err.Error())
		return
	}

	todos, err := services.GetAllTodos(userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, err.Error())
		return
	}
	var resTodo []models.ResTodo
	// convert res to json
	func() {
		for _, todo := range todos {
			var newTodo models.ResTodo
			newTodo.Id = todo.ID
			newTodo.Title = todo.Title
			newTodo.Completed = todo.Completed
			resTodo = append(resTodo, newTodo)
		}
	}()

	response := models.Response{Message: "successful", Data: resTodo}
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Print(err)
		fmt.Fprintf(w, "internal error while parsing todos")
		return
	}
}

func PostTodo(w http.ResponseWriter, req *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	// check if authenticated
	userId, err := isAuthenticated(req)
	if err != nil {
		log.Print(err)
		message := err.Error()
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, message)
		return
	}
	// get todo from request
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
	w.WriteHeader(http.StatusOK)
	w.Write(response)
	return
}

func DeleteTodo(w http.ResponseWriter, req *http.Request) {}
