package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/manish-mehra/go-todo/services"
	"github.com/manish-mehra/go-todo/utils"
)

var Serve http.Handler

func init() {

	log.Print("server init")
	r := http.NewServeMux()

	r.HandleFunc("GET /api/health", health)
	r.HandleFunc("POST /api/register", registerUser)

	Serve = r
}

func health(w http.ResponseWriter, req *http.Request) {
	response, _ := utils.ParseResponse("ok 🟢")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func registerUser(w http.ResponseWriter, req *http.Request) {
	// parse user
	newUser, err := utils.DecodeUserJSON(req.Body)
	if err != nil {
		message := "Error parsing user"
		response, _ := utils.ParseResponse(message)
		w.Header().Set("Content-Type", "application/json")
		w.Write(response) // Write the JSON data to the response body
		return
	}
	log.Print(newUser)

	// add user
	err = services.CreateUser(newUser)
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}

	response, _ := utils.ParseResponse("user created")
	w.Header().Set("Content-Type", "application/json")
	w.Write(response) // Write the JSON data to the response body
}
