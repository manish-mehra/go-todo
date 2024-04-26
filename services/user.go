// functionality for managing user data.

package services

import (
	"database/sql"
	"errors"

	_ "github.com/go-sql-driver/mysql"

	// my imports
	"github.com/manish-mehra/go-todo/models"
	"github.com/manish-mehra/go-todo/utils"
)

type UserService struct {
	db                 *sql.DB
	createUserStmt     *sql.Stmt
	getUserByEmailStmt *sql.Stmt
}

func NewUserService(db *sql.DB) (*UserService, error) {

	svc := &UserService{
		db: db,
	}

	err := svc.prepareStatements()
	if err != nil {
		return nil, err
	}

	return svc, nil
}

// method to prepare all User statements
func (u *UserService) prepareStatements() error {

	err := u.prepareCreateUserStmt()
	if err != nil {
		return err
	}

	err = u.prepareGetUserByEmailStmt()
	if err != nil {
		return err
	}

	return nil
}

func (u *UserService) prepareCreateUserStmt() error {
	stmt, err := u.db.Prepare("INSERT INTO User (name, email, role, password) VALUES (?, ?, ?, ?)")
	if err != nil {
		return err
	}
	u.createUserStmt = stmt
	return nil
}

func (u *UserService) prepareGetUserByEmailStmt() error {
	stmt, err := u.db.Prepare("SELECT * FROM User WHERE email = ?")
	if err != nil {
		return err
	}
	u.getUserByEmailStmt = stmt
	return nil
}

// CreateUser inserts a new user into the database.
func (u *UserService) CreateUser(user models.User) error {
	_, err := u.createUserStmt.Exec(user.Name, user.Email, user.Role, user.Password)
	if err != nil {
		return errors.New("failed to create user")
	}
	return nil
}

// GetUserByEmail retrieves a user from the database by email.
// Returns a "user not found" error if the user is not found.
func (u *UserService) GetUserByEmail(userMail string) (models.User, error) {
	var user models.User

	err := u.getUserByEmailStmt.QueryRow(userMail).Scan(&user.Id, &user.Name, &user.Email, &user.Role, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return user, utils.ErrNotFound
		}
		return user, errors.New("failed to get user: " + err.Error())
	}
	return user, nil
}
