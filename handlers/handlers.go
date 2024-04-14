package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/manish-mehra/go-todo/utils"
)

var Serve http.Handler

func init() {

	log.Print("server init")
	r := http.NewServeMux()

	r.HandleFunc("GET /api/health", health)
	r.HandleFunc("POST /api/register", RegisterUser)
	r.HandleFunc("POST /api/login", LoginUser)
	r.HandleFunc("GET /api/protected", protectedRoute)
	r.HandleFunc("POST /api/todo", PostTodo)

	Serve = r
}

func health(w http.ResponseWriter, req *http.Request) {
	response, _ := utils.ParseResponse("ok 🟢")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
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
		fmt.Fprintf(w, error.Error(err))
		return
	}

	fmt.Print("owner", owner)
	message := "Welcome to protected route " + owner
	fmt.Fprintf(w, message)
}
