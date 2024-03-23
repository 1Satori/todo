package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	todo "todo-app"
)

type TodoListMysql struct {
	db *sqlx.DB
}

func NewTodoListMysql(db *sqlx.DB) *TodoListMysql {
	return &TodoListMysql{db: db}
}

func (r *TodoListMysql) Create(userId int, list todo.TodoList) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	createListQuery := fmt.Sprintf("INSERT INTO %s (title, description) VALUES (?, ?)", todoListTable)
	result, err := tx.Exec(createListQuery, list.Title, list.Description)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	lastID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	createUsersListQuery := fmt.Sprintf("INSERT INTO %s (user_id, list_id) VALUES (?, ?)", usersListsTable)
	_, err = tx.Exec(createUsersListQuery, userId, lastID)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	return int(lastID), tx.Commit()
}
