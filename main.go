package main

import (
	"net/http"

	"github.com/manish-mehra/go-todo/handlers"
)

func main() {
	http.ListenAndServe(":8080", handlers.Serve)
}
