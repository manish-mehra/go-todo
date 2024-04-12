package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/manish-mehra/go-todo/models"
	"github.com/manish-mehra/go-todo/services"
	"github.com/manish-mehra/go-todo/utils"
)

var Serve http.Handler

func init() {

	log.Print("server init")
	r := http.NewServeMux()

	r.HandleFunc("GET /api/health", health)
	r.HandleFunc("POST /api/register", registerUser)
	r.HandleFunc("POST /api/login", loginUser)
	r.HandleFunc("GET /api/protected", protectedRoute)

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

func loginUser(w http.ResponseWriter, req *http.Request) {

	var loggedUser struct {
		Email    string
		Password string
	}

	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&loggedUser); err != nil {
		message := "Error parsing user"
		response, _ := utils.ParseResponse(message)
		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
		return
	}

	// get the user from database
	dbUser, err := services.GetUserByEmail(loggedUser.Email)
	if err != nil {
		message := string(err.Error())
		response, _ := utils.ParseResponse(message)
		w.Header().Set("Content-Type", "application/json")
		w.Write(response) // Write the JSON data to the response body
		return
	}

	// match the passwords
	if loggedUser.Password != dbUser.Password {
		fmt.Print(dbUser, loggedUser)
		message := "incorrect password"
		response, _ := utils.ParseResponse(message)
		w.Header().Set("Content-Type", "application/json")
		w.Write(response) // Write the JSON data to the response body
		return
	}

	// create jwt token
	token, err := utils.CreateToken(dbUser.Email)
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
		Domain: "http://localhost:8080",
		Path:   "/",
		Secure: true, // Only send cookie over HTTPS connections (if applicable)
	}
	// Set the cookie in the response
	http.SetCookie(w, &cookie)
	log.Print(cookie)
	resUser := struct {
		Id    string `json:"id"`
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

// FOR TESTING
func protectedRoute(w http.ResponseWriter, req *http.Request) {
	// get the token from request cookie
	cookies := req.Cookies()
	var token string
	for _, cookie := range cookies {
		if cookie.Name == "token" { // Replace "myCookie" with your actual cookie name
			token = cookie.Value
		}
	}

	if token == "" {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, "JWT token not found")
		fmt.Print("JWT not found")
		return
	}

	// verify token
	owner, err := utils.VerifyToken(token)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, "Invalid Token!")
		return
	}

	fmt.Print("owner", owner)
	message := "Welcome to protected route " + owner
	fmt.Fprintf(w, message)
}
