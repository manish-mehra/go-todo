package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/manish-mehra/go-todo/models"
	"github.com/manish-mehra/go-todo/services"
	"github.com/manish-mehra/go-todo/utils"
)

func RegisterUser(w http.ResponseWriter, req *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	// parse user
	newUser, err := utils.DecodeUserJSON(req.Body)
	if err != nil {
		message := "Error parsing user"
		response, _ := utils.ParseResponse(message)
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
	w.Write(response) // Write the JSON data to the response body
}

func LoginUser(w http.ResponseWriter, req *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	var loggedUser struct {
		Email    string
		Password string
	}

	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&loggedUser); err != nil {
		message := "Error parsing user"
		response, _ := utils.ParseResponse(message)
		w.Write(response)
		return
	}

	// get the user from database
	dbUser, err := services.GetUserByEmail(loggedUser.Email)
	if err != nil {
		message := string(err.Error())
		response, _ := utils.ParseResponse(message)
		w.Write(response) // Write the JSON data to the response body
		return
	}

	// match the passwords
	if loggedUser.Password != dbUser.Password {
		fmt.Print(dbUser, loggedUser)
		message := "incorrect password"
		response, _ := utils.ParseResponse(message)
		w.Write(response) // Write the JSON data to the response body
		return
	}

	// create jwt token
	userId := strconv.Itoa(int(dbUser.Id))
	token, err := utils.CreateToken(userId)
	if err != nil {
		message := "error generating token"
		response, _ := utils.ParseResponse(message)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(response)
		fmt.Print(message)
		return
	}

	// attach JWT TOKEN TO Response
	cookie := http.Cookie{
		Name:   "token",
		Value:  token,
		MaxAge: 3600,
		Path:   "/",
		Secure: true, // Only send cookie over HTTPS connections (if applicable)
	}
	// Set the cookie in the response
	http.SetCookie(w, &cookie)
	resUser := struct {
		Id    int    `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
		Role  string `json:"role"`
	}{Id: dbUser.Id, Name: dbUser.Name, Email: dbUser.Email, Role: dbUser.Role}
	response := models.Response{Message: "successful!", Data: resUser}
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Print(err)
		fmt.Fprintf(w, "internal error")
		return
	}
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

		w.Header().Set("Content-Type", "application/json")

		// get the token from request cookie
		cookies := req.Cookies()
		var token string
		for _, cookie := range cookies {
			if cookie.Name == "token" {
				token = cookie.Value
			}
		}
		if token == "" {
			fmt.Fprintf(w, "Not JWT Token Found!")
			return
		}
		// verify token
		userId, err := utils.VerifyToken(token)
		if err != nil {
			fmt.Fprintf(w, "Unauthorized")
			return
		}
		log.Print("authorized: " + userId)

		// Store the user ID in the request context
		ctx := context.WithValue(req.Context(), "userId", userId)
		req = req.WithContext(ctx)

		next.ServeHTTP(w, req)
	})
}
