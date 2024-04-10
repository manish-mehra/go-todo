package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/manish-mehra/go-todo/services"
	"github.com/manish-mehra/go-todo/utils"
)

func RegisterUser(w http.ResponseWriter, req *http.Request) {
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
