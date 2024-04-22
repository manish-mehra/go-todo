package handlers

import (
	"net/http"
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

// Serve is the HTTP handler initialized in the init function,
// containing all configured routes and middleware for the application.
var Serve http.Handler

func init() {

	r := http.NewServeMux()

	r.HandleFunc("POST /api/register", RegisterUser)
	r.HandleFunc("POST /api/login", LoginUser)

	// // add auth middleware
	r.Handle("POST /api/todo", ApplyMiddleware(http.HandlerFunc(PostTodo), AuthMiddleware))
	r.Handle("GET /api/todo/{id}", ApplyMiddleware(http.HandlerFunc(GetTodo), AuthMiddleware))
	r.Handle("GET /api/todos", ApplyMiddleware(http.HandlerFunc(GetAllTodo), AuthMiddleware))
	r.Handle("DELETE /api/todo/{id}", ApplyMiddleware(http.HandlerFunc(DeleteTodo), AuthMiddleware))
	r.Handle("PUT /api/todo/{id}", ApplyMiddleware(http.HandlerFunc(UpdateTodo), AuthMiddleware))

	Serve = r
}
