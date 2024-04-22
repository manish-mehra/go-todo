package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/manish-mehra/go-todo/models"
	"github.com/manish-mehra/go-todo/services"
	"github.com/manish-mehra/go-todo/utils"
)

func GetTodo(w http.ResponseWriter, req *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	// get userid from request context
	userId, ok := req.Context().Value("userId").(string)
	if !ok {
		http.Error(w, "User ID not found", http.StatusUnauthorized)
		return
	}

	// get todo id from req param
	paramTodoID := req.PathValue("id")
	if paramTodoID == "" {
		http.Error(w, "Todo ID not found", http.StatusBadRequest)
		return
	}

	// convert todo id to int
	todoID, err := utils.StringToInt(paramTodoID)
	if err != nil {
		log.Printf("Error converting todo id %d to string", todoID)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// convert user id to int
	userID, err := utils.StringToInt(userId)
	if err != nil {
		log.Printf("Error converting user ID to %d string", userID)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	// get todo from db
	todo, err := services.GetTodo(todoID, userID)
	if err != nil {
		if errors.Is(err, utils.ErrNotFound) {
			http.Error(w, "Todo not found", http.StatusNotFound)
			return
		}
		log.Printf("Error fetching todo for  user ID %d and ID %d error %d", userID, todoID, errors.New(err.Error()))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// convert res to json
	resTodo := struct {
		Id        int    `json:"id"`
		Title     string `json:"title"`
		Completed bool   `json:"completed"`
	}{Id: todo.ID, Title: todo.Title, Completed: todo.Completed}

	response := models.Response{Message: "successful", Data: resTodo}
	err = json.NewEncoder(w).Encode(response)

	if err != nil {
		log.Printf("Error necoding response %d", errors.New(err.Error()))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func GetAllTodo(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// get userid from request context
	userId, ok := req.Context().Value("userId").(string)
	if !ok {
		http.Error(w, "User ID not found", http.StatusUnauthorized)
		return
	}

	userID, err := utils.StringToInt(userId)
	if err != nil {
		log.Printf("Error converting user ID to %d string", userID)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	todos, err := services.GetAllTodos(userID)
	if err != nil {
		if errors.Is(err, utils.ErrNotFound) {
			http.Error(w, "Todos not found", http.StatusNotFound)
			return
		}
		log.Printf("Error fetching todos for  user ID %d  error %d", userID, errors.New(err.Error()))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
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
		log.Printf("Error necoding response %d", errors.New(err.Error()))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func PostTodo(w http.ResponseWriter, req *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	// get user id from request context
	userId, ok := req.Context().Value("userId").(string)
	if !ok {
		http.Error(w, "User ID not found", http.StatusUnauthorized)
		return
	}
	// get todo from request
	var newTodo models.UserTodo
	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&newTodo); err != nil {
		log.Printf("Error converting todo to json for user ID %s", userId)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// convert user id to int
	userID, err := utils.StringToInt(userId)
	if err != nil {
		log.Printf("Error converting user ID %d to int", userID)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = services.PostTodo(newTodo, userID)
	if err != nil {
		log.Printf("Error posting  todo with user ID %d", userID)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	response := models.Response{Message: "successful", Data: nil}
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Printf("Error necoding response %d", errors.New(err.Error()))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func DeleteTodo(w http.ResponseWriter, req *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	// get userid from request context
	userId, ok := req.Context().Value("userId").(string)
	if !ok {
		http.Error(w, "User ID not found", http.StatusUnauthorized)
		return
	}

	// get todo id from req param
	paramTodoID := req.PathValue("id")
	log.Print("paramTodoID ", paramTodoID)
	if paramTodoID == "" {
		http.Error(w, "Missing Todo ID", http.StatusBadRequest)
		return
	}

	// convert todo  id to int
	todoID, err := utils.StringToInt(paramTodoID)
	if err != nil {
		log.Printf("Error converting todo id %d to string", todoID)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	userID, err := utils.StringToInt(userId)
	if err != nil {
		log.Printf("Error converting user ID to %d string", userID)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// delete  todo from db
	err = services.DeleteTodo(todoID, userID)
	if err != nil {
		if errors.Is(err, utils.ErrNotFound) {
			http.Error(w, "Todo not found", http.StatusNotFound)
			return
		}
		log.Printf("Error deleting  todo with user ID %d", userID)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	response := models.Response{Message: "successful", Data: nil}
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Printf("Error necoding response %d", errors.New(err.Error()))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func UpdateTodo(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// get userid from request context
	userId, ok := req.Context().Value("userId").(string)
	if !ok {
		http.Error(w, "User ID not found", http.StatusUnauthorized)
		return
	}

	// get todo from request
	var newTodo models.UserTodo
	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&newTodo); err != nil {
		log.Printf("Error converting todo to json for user ID %s", userId)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// get todo id from req param
	paramTodoID := req.PathValue("id")
	log.Print("paramTodoID ", paramTodoID)
	if paramTodoID == "" {
		http.Error(w, "Missing Todo ID", http.StatusBadRequest)
		return
	}

	// convert todo  id to int
	todoID, err := utils.StringToInt(paramTodoID)
	if err != nil {
		log.Printf("Error converting todo id %d to string", todoID)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	userID, err := utils.StringToInt(userId)
	if err != nil {
		log.Printf("Error converting user ID to %d string", userID)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// update the todo
	err = services.UpdateTodo(newTodo, todoID, userID)
	if err != nil {
		if errors.Is(err, utils.ErrNotFound) {
			http.Error(w, "Todo not found", http.StatusNotFound)
			return
		}
		log.Printf("Error updating  todo with user ID %d", userID)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	response := models.Response{Message: "successful", Data: nil}
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Printf("Error necoding response %d", errors.New(err.Error()))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
