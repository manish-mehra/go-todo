package services

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var mysqlDB *sql.DB

func init() {
	ConnectToDB()
}

func ConnectToDB() error {
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
		return errors.New("error connecting to database")
	}

	// defer db.Close()

	// Test the connection
	pingErr := db.Ping()
	if pingErr != nil {
		log.Print("ERROR ON PING DB!")
		return errors.New("ERROR ON PING DB")
	}

	mysqlDB = db
	log.Print("Successfully connected to MySQL database!")
	return nil
}
