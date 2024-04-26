package services

import (
	"database/sql"
	"errors"

	"github.com/manish-mehra/go-todo/models"
	"github.com/manish-mehra/go-todo/utils"
)

type TodoService struct {
	db             *sql.DB
	getTodoStmt    *sql.Stmt
	postTodoStmt   *sql.Stmt
	getAllTodoStmt *sql.Stmt
	deleteTodoStmt *sql.Stmt
	updateTodoStmt *sql.Stmt
}

// function to prepare and set TodoService statements
func prepareStmts(svc *TodoService) error {

	err := svc.prepareGetTodoStmt()
	if err != nil {
		return err
	}

	err = svc.preparePostTodoStmt()
	if err != nil {
		return err
	}

	err = svc.prepareGetAllTodoStmt()
	if err != nil {
		return err
	}

	err = svc.prepareDeleteTodoStmt()
	if err != nil {
		return err
	}

	err = svc.prepareUpdateTodoStmt()
	if err != nil {
		return err
	}

	return nil
}

func NewTodoService(db *sql.DB) (*TodoService, error) {
	svc := &TodoService{
		db: db,
	}
	// prepare statements
	err := prepareStmts(svc)
	if err != nil {
		return nil, err
	}

	return svc, nil
}

func (t *TodoService) prepareGetTodoStmt() error {
	stmt, err := t.db.Prepare("SELECT id, title, completed FROM Todo WHERE id = ? AND user_id = ?")
	if err != nil {
		return err
	}
	t.getTodoStmt = stmt
	return nil
}
func (t *TodoService) preparePostTodoStmt() error {
	stmt, err := t.db.Prepare("INSERT INTO Todo (title, completed, user_id) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}
	t.postTodoStmt = stmt
	return nil
}
func (t *TodoService) prepareGetAllTodoStmt() error {
	stmt, err := t.db.Prepare("SELECT id, title, completed FROM Todo WHERE user_id = ?")
	if err != nil {
		return err
	}
	t.getAllTodoStmt = stmt
	return nil
}
func (t *TodoService) prepareDeleteTodoStmt() error {
	stmt, err := t.db.Prepare("DELETE FROM Todo WHERE id=? AND user_id = ?")
	if err != nil {
		return err
	}
	t.deleteTodoStmt = stmt
	return nil
}
func (t *TodoService) prepareUpdateTodoStmt() error {
	stmt, err := t.db.Prepare("UPDATE Todo SET title = ?, completed = ?  WHERE id = ? AND user_id = ?")
	if err != nil {
		return err
	}
	t.updateTodoStmt = stmt
	return nil
}

func (t *TodoService) PostTodo(todo models.UserTodo, userId int) error {
	result, err := t.postTodoStmt.Exec(todo.Title, todo.Completed, userId)
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
func (t *TodoService) GetTodo(id int, userId int) (models.Todo, error) {
	var todo models.Todo
	err := t.getTodoStmt.QueryRow(id, userId).Scan(&todo.ID, &todo.Title, &todo.Completed)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Todo{}, utils.ErrNotFound
		}
		return models.Todo{}, errors.New("failed to get todo: " + err.Error())
	}
	return todo, nil
}
func (t *TodoService) GetAllTodos(userId int) ([]models.Todo, error) {
	var todos []models.Todo
	rows, err := t.getAllTodoStmt.Query(userId)
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
func (t *TodoService) DeleteTodo(id int, userId int) error {
	result, err := t.deleteTodoStmt.Exec(id, userId)
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
func (t *TodoService) UpdateTodo(todo models.UserTodo, todoId int, userID int) error {
	result, err := t.updateTodoStmt.Exec(todo.Title, todo.Completed, todoId, userID)
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
