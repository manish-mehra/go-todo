package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func main() {

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
		panic(err.Error())
	}

	defer db.Close()

	// Test the connection
	pingErr := db.Ping()
	if pingErr != nil {
		log.Print("Error connecting to database:", pingErr)
		return
	}

	log.Print("Successfully connected to MySQL database!")

	fileServer := http.FileServer(http.Dir("static"))
	http.Handle("/", http.StripPrefix("/", fileServer))

	http.ListenAndServe(":8080", nil)
}
