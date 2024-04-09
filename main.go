package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

type Response struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"` // Can hold any type of data
}

type USER struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

var mysqlDB *sql.DB

func connectToDB() error {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	user := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	// Create connection string
	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, password, host, port, dbName)

	// Open database connection
	db, err := sql.Open("mysql", connectionString)

	if err != nil {
		return errors.New("Error connecting to database")
	}

	// defer db.Close()

	// Test the connection
	pingErr := db.Ping()
	if pingErr != nil {
		log.Print("ERROR ON PING DB!")
		return errors.New("ERROR ON PING DB")
	}

	mysqlDB = db
	return nil
}

func main() {

	err := connectToDB()
	if err != nil {
		return
	}
	defer mysqlDB.Close()
	log.Print("Successfully connected to MySQL database!")

	fileServer := http.FileServer(http.Dir("static"))
	http.Handle("/", http.StripPrefix("/", fileServer))

	http.HandleFunc("/api/register", RegisterUser)

	http.ListenAndServe(":8080", nil)
}

func CreateUser(user USER) error {
	stmt, err := mysqlDB.Prepare("INSERT INTO User (name, email, role) VALUES (?, ?, ?)")
	if err != nil {
		return errors.New("failed to create user: wrong query")
	}
	defer stmt.Close()
	// Execute the insert query with the data
	_, err = stmt.Exec(user.Name, user.Email, user.Role)
	if err != nil {
		return errors.New("failed to create user: can't execute query")
	}
	return nil
}

func RegisterUser(w http.ResponseWriter, req *http.Request) {
	// parse user
	newUser, err := DecodeUserJSON(req.Body)
	if err != nil {
		message := "Error parsing user"
		response, _ := parseResponse(message)
		w.Header().Set("Content-Type", "application/json")
		w.Write(response) // Write the JSON data to the response body
		return
	}
	log.Print(newUser)

	// add user
	err = CreateUser(newUser)
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}

	response, _ := parseResponse("user created")
	w.Header().Set("Content-Type", "application/json")
	w.Write(response) // Write the JSON data to the response body

}
