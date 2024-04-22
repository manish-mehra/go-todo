package services

import (
	"database/sql"
	"errors"
	"log"

	"github.com/manish-mehra/go-todo/models"
	"github.com/manish-mehra/go-todo/utils"
)

var (
	stmtPostTodo   *sql.Stmt
	stmtGetTodo    *sql.Stmt
	stmtGetTodos   *sql.Stmt
	stmtDeleteTodo *sql.Stmt
	stmtUpdateTodo *sql.Stmt
)

func init() {

	var err error

	// prepare statments

	stmtPostTodo, err = mysqlDB.Prepare("INSERT INTO Todo (title, completed, user_id) VALUES (?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	stmtGetTodo, err = mysqlDB.Prepare("SELECT id, title, completed FROM Todo WHERE id = ? AND user_id = ?")
	if err != nil {
		log.Fatal(err)
	}

	stmtGetTodos, err = mysqlDB.Prepare("SELECT id, title, completed FROM Todo WHERE user_id = ?")
	if err != nil {
		log.Fatal(err)
	}

	stmtDeleteTodo, err = mysqlDB.Prepare("DELETE FROM Todo WHERE id=? AND user_id = ?")
	if err != nil {
		log.Fatal(err)
	}

	stmtUpdateTodo, err = mysqlDB.Prepare("UPDATE Todo SET title = ?, completed = ?  WHERE id = ? AND user_id = ?")
	if err != nil {
		log.Fatal(err)
	}
}

func PostTodo(todo models.UserTodo, userId int) error {
	result, err := stmtPostTodo.Exec(todo.Title, todo.Completed, userId)
	if err != nil {
		return errors.New("failed to post todo")
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err // handle error getting rows affected
	}
	if rowsAffected == 0 {
		return utils.ErrNotFound // handle missing ID
	}
	return nil
}

// GetTodo return Todo & error
// Expect userId and id(todo) as arg
func GetTodo(id int, userId int) (models.Todo, error) {
	var todo models.Todo
	err := stmtGetTodo.QueryRow(id, userId).Scan(&todo.ID, &todo.Title, &todo.Completed)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Todo{}, utils.ErrNotFound
		}
		return models.Todo{}, errors.New("failed to get todo: " + err.Error())
	}
	return todo, nil
}

func GetAllTodos(userId int) ([]models.Todo, error) {
	var todos []models.Todo
	rows, err := stmtGetTodos.Query(userId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, utils.ErrNotFound
		}
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

func DeleteTodo(id int, userId int) error {
	result, err := stmtDeleteTodo.Exec(id, userId)
	if err != nil {
		return errors.New("failed to delete todo")
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err // handle error getting rows affected
	}
	if rowsAffected == 0 {
		return utils.ErrNotFound // handle missing ID
	}
	return nil
}

func UpdateTodo(todo models.UserTodo, todoId int, userID int) error {
	result, err := stmtUpdateTodo.Exec(todo.Title, todo.Completed, todoId, userID)
	if err != nil {
		return errors.New("failed to update todo")
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err // handle error getting rows affected
	}
	if rowsAffected == 0 {
		return utils.ErrNotFound // handle missing ID
	}
	return nil
}
