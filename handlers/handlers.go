package handlers

import (
	"log"
	"net/http"

	"github.com/manish-mehra/go-todo/utils"
)

var Serve http.Handler

func init() {

	log.Print("server init")
	r := http.NewServeMux()

	r.HandleFunc("GET /api/health", health)
	r.HandleFunc("POST /api/register", RegisterUser)
	r.HandleFunc("POST /api/login", LoginUser)
	r.HandleFunc("POST /api/todo", PostTodo)
	r.HandleFunc("GET /api/todo/{id}", GetTodo)
	r.HandleFunc("GET /api/todos", GetAllTodo)

	Serve = r
}

func health(w http.ResponseWriter, req *http.Request) {
	response, _ := utils.ParseResponse("ok 🟢")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
