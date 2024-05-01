package handlers

import (
	"net/http"

	"github.com/manish-mehra/go-todo/services"
)

// ApplyMiddleware applies a series of middleware functions to the given handler,
// chaining them in sequence to create a new composed handler with all middleware applied.
type MiddlewareFunc func(http.Handler) http.Handler

func ApplyMiddleware(handler http.Handler, middlewares ...MiddlewareFunc) http.Handler {
	for _, middleware := range middlewares {
		handler = middleware(handler)
	}
	return handler
}

var Serve http.Handler

func init() {

	userHandlers := &UserHandler{
		userService: services.UserSvc,
	}

	todoHandler := &TodoHandler{
		todoService: services.TodoSvc,
	}

	r := http.NewServeMux()

	r.HandleFunc("POST /api/register", userHandlers.RegisterUser)
	r.HandleFunc("POST /api/login", userHandlers.LoginUser)

	// // add auth middleware
	r.Handle("POST /api/todo", ApplyMiddleware(http.HandlerFunc(todoHandler.PostTodo), userHandlers.AuthMiddleware))
	r.Handle("GET /api/todo/{id}", ApplyMiddleware(http.HandlerFunc(todoHandler.GetTodo), userHandlers.AuthMiddleware))
	r.Handle("GET /api/todos", ApplyMiddleware(http.HandlerFunc(todoHandler.GetAllTodo), userHandlers.AuthMiddleware))
	r.Handle("DELETE /api/todo/{id}", ApplyMiddleware(http.HandlerFunc(todoHandler.DeleteTodo), userHandlers.AuthMiddleware))
	r.Handle("PUT /api/todo/{id}", ApplyMiddleware(http.HandlerFunc(todoHandler.UpdateTodo), userHandlers.AuthMiddleware))

	Serve = r
}
