package handlers

import (
	"context"
	"encoding/json"
	"errors"
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
		log.Printf("Error converting User to json")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// check if user with same email already exist
	user, _ := services.GetUserByEmail(newUser.Email)
	// if email is there, user don't exist
	if user.Email != "" {
		http.Error(w, "User already exist", http.StatusBadRequest)
		return
	}

	// add user
	err = services.CreateUser(newUser)
	if err != nil {
		log.Printf("Error creating user")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	log.Printf("New User Id %d", newUser.Id)

	response := models.Response{Message: "successful", Data: nil}
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Printf("Error necoding response %d", errors.New(err.Error()))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func LoginUser(w http.ResponseWriter, req *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	var loggedUser struct {
		Email    string
		Password string
	}

	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&loggedUser); err != nil {
		log.Printf("Error converting User to json")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// get the user from database
	dbUser, err := services.GetUserByEmail(loggedUser.Email)
	if err != nil {
		if errors.Is(err, utils.ErrNotFound) {
			http.Error(w, "User not found", http.StatusBadRequest)
			return
		}
		log.Printf("Error getting  user from database")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// match the passwords
	if loggedUser.Password != dbUser.Password {
		http.Error(w, "Incorrect Password", http.StatusBadRequest)
		return
	}

	// create jwt token
	userId := strconv.Itoa(int(dbUser.Id))
	token, err := utils.CreateToken(userId)
	if err != nil {
		log.Printf("Error generating  jwt token for user id %d", dbUser.Id)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
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
		log.Printf("Error necoding response %d", errors.New(err.Error()))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
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
			http.Error(w, "No jwt token found", http.StatusBadRequest)
			return
		}
		// verify token
		userId, err := utils.VerifyToken(token)
		if err != nil {
			log.Printf("Unauthorized user %s", userId)
			http.Error(w, "Unauthorized", http.StatusBadRequest)
			return
		}

		// Store the user ID in the request context
		ctx := context.WithValue(req.Context(), "userId", userId)
		req = req.WithContext(ctx)

		next.ServeHTTP(w, req)
	})
}
