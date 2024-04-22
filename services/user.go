// functionality for managing user data.

package services

import (
	"database/sql"
	"errors"
	"log"

	_ "github.com/go-sql-driver/mysql"

	// my imports
	"github.com/manish-mehra/go-todo/models"
	"github.com/manish-mehra/go-todo/utils"
)

// Package-level variables to hold prepared SQL statements for
// creating a new user and retrieving a user by email.
var (
	stmtCreateUser     *sql.Stmt
	stmtGetUserByEmail *sql.Stmt
)

// init prepares the SQL statements for creating and retrieving users.
// It assumes that mysqlDB is a valid database connection object.
func init() {

	var err error

	stmtCreateUser, err = mysqlDB.Prepare("INSERT INTO User (name, email, role, password) VALUES (?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}

	stmtGetUserByEmail, err = mysqlDB.Prepare("SELECT * FROM User WHERE email = ?")
	if err != nil {
		log.Fatal(err)
	}

}

// CreateUser inserts a new user into the database.
func CreateUser(user models.User) error {
	_, err := stmtCreateUser.Exec(user.Name, user.Email, user.Role, user.Password)
	if err != nil {
		return errors.New("failed to create user")
	}
	return nil
}

// GetUserByEmail retrieves a user from the database by email.
// Returns a "user not found" error if the user is not found.
func GetUserByEmail(userMail string) (models.User, error) {
	var user models.User

	err := stmtGetUserByEmail.QueryRow(userMail).Scan(&user.Id, &user.Name, &user.Email, &user.Role, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return user, utils.ErrNotFound
		}
		return user, errors.New("failed to get user: " + err.Error())
	}
	return user, nil
}
