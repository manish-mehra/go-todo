package main

import (
	"net/http"

	"github.com/manish-mehra/go-todo/handlers"
)

func main() {

	fileServer := http.FileServer(http.Dir("static"))
	http.Handle("/", http.StripPrefix("/", fileServer))

	http.HandleFunc("/api/register", handlers.RegisterUser)

	http.ListenAndServe(":8080", nil)
}
