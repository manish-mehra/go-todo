package services

import (
	"database/sql"
	"errors"
	"log"

	"github.com/manish-mehra/go-todo/models"
)

var (
	stmtPostTodo *sql.Stmt
	stmtGetTodo  *sql.Stmt
)

func init() {

	var err error
	// prepare posttodo statement
	stmtPostTodo, err = mysqlDB.Prepare("INSERT INTO Todo (title, completed, user_id) VALUES (?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	// Select * From Todo Where user_id = 1231 AND id = 4;
	stmtGetTodo, err = mysqlDB.Prepare("SELECT id, title, completed FROM Todo WHERE user_id = ? AND id = ?")
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

// GetTodo return Todo & error
// Expect userId and id(todo) as arg
func GetTodo(userId int64, id int64) (models.Todo, error) {
	var todo models.Todo
	err := stmtGetTodo.QueryRow(userId, id).Scan(&todo.ID, &todo.Title, &todo.Completed)
	if err != nil {
		if err == sql.ErrNoRows {
			// No todo found with the provided user ID and ID
			return models.Todo{}, errors.New("todo not found")
		}
		return models.Todo{}, errors.New("failed to get todo: " + err.Error())
	}
	return todo, nil
}
