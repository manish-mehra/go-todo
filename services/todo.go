package services

import (
	"database/sql"
	"errors"
	"log"

	"github.com/manish-mehra/go-todo/models"
)

var (
	stmtPostTodo *sql.Stmt
)

func init() {

	var err error
	// prepare posttodo statement
	stmtPostTodo, err = mysqlDB.Prepare("INSERT INTO Todo (title, completed, user_id) VALUES (?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
}

func PostTodo(todo models.UserTodo, userId int64) error {
	_, err := stmtPostTodo.Exec(todo.Title, todo.Completed, userId)
	if err != nil {
		return errors.New("failed to post todo")
	}
	return nil
}
