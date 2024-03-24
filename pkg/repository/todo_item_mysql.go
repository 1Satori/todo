package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	todo "todo-app"
)

type TodoItemMysql struct {
	db *sqlx.DB
}

func NewTodoItemMysql(db *sqlx.DB) *TodoItemMysql {
	return &TodoItemMysql{db: db}
}

func (r *TodoItemMysql) Create(listId int, item todo.TodoItem) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	createItemQuery := fmt.Sprintf("INSERT INTO %s (title, description) VALUES (?, ?)", todoItemsTable)
	res, err := tx.Exec(createItemQuery, item.Title, item.Description)
	if err != nil {
		return 0, err
	}
	lastID, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	createListItemsQuery := fmt.Sprintf("INSERT INTO %s (list_id, item_id) values (? ?)", listsItemsTable)
	_, err = tx.Exec(createListItemsQuery, listId, lastID)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	return int(lastID), tx.Commit()
}

func (r *TodoItemMysql) GetAll(userId, listId int) ([]todo.TodoItem, error) {
	var items []todo.TodoItem

	query := fmt.Sprintf("SELECT ti.id, ti.title, ti.description, ti.done FROM %s ti INNER JOIN %s li ON li.item_id = ti.id INNER JOIN %s ul"+
		" ON ul.list_id = li.list_id WHERE li.list_id = ? AND ul.user_id = ?",
		todoItemsTable, listsItemsTable, usersListsTable)

	if err := r.db.Select(&items, query, listId, userId); err != nil {
		return nil, err
	}

	return items, nil
}
