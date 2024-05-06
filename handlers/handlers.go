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

	middleware := Middleware{}

	r := http.NewServeMux()

	userSrvc := services.UserSvc
	todoSrvc := services.TodoSvc

	r.Handle("POST /api/register", RegisterUserHandler(userSrvc))
	r.Handle("POST /api/login", LoginUserHandler(userSrvc))

	// // add auth middleware
	r.Handle("POST /api/todo",
		ApplyMiddleware(PostTodoHandler(todoSrvc), middleware.AuthMiddleware))
	r.Handle("GET /api/todo/{id}",
		ApplyMiddleware(GetTodoHandler(todoSrvc), middleware.AuthMiddleware))
	r.Handle("GET /api/todos",
		ApplyMiddleware(GetTodosHandler(todoSrvc), middleware.AuthMiddleware))
	r.Handle("DELETE /api/todo/{id}",
		ApplyMiddleware(DeleteTodoHandler(todoSrvc), middleware.AuthMiddleware))
	r.Handle("PUT /api/todo/{id}",
		ApplyMiddleware(UpdateTodoHandler(todoSrvc), middleware.AuthMiddleware))

	Serve = r
}
