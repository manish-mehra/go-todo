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
	stmtGetTodos *sql.Stmt
)

func init() {

	var err error

	// prepare statments

	stmtPostTodo, err = mysqlDB.Prepare("INSERT INTO Todo (title, completed, user_id) VALUES (?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	stmtGetTodo, err = mysqlDB.Prepare("SELECT id, title, completed FROM Todo WHERE user_id = ? AND id = ?")
	if err != nil {
		log.Fatal(err)
	}

	stmtGetTodos, err = mysqlDB.Prepare("SELECT id, title, completed FROM Todo WHERE user_id = ?")
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

func GetAllTodos(userId int64) ([]models.Todo, error) {
	var todos []models.Todo
	rows, err := stmtGetTodos.Query(userId)
	if err != nil {
		return nil, errors.New("failed to get all todos: " + err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		var todo models.Todo
		err := rows.Scan(&todo.ID, &todo.Title, &todo.Completed)
		if err != nil {
			return nil, errors.New("failed to scan todo: " + err.Error())
		}
		todos = append(todos, todo)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.New("failed to iterate todos: " + err.Error())
	}

	return todos, nil
}
