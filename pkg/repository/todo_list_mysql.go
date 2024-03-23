package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"strings"
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

func (r *TodoListMysql) GetAll(userId int) ([]todo.TodoList, error) {
	var lists []todo.TodoList
	query := fmt.Sprintf("SELECT tl.id, tl.title, tl.description FROM %s AS tl INNER JOIN %s AS ul ON tl.id = ul.list_id WHERE ul.user_id = ?", todoListTable, usersListsTable)
	err := r.db.Select(&lists, query, userId)

	return lists, err
}

func (r *TodoListMysql) GetById(userId, listId int) (todo.TodoList, error) {
	var list todo.TodoList

	query := fmt.Sprintf("SELECT tl.id, tl.title, tl.description FROM %s AS tl INNER JOIN %s "+
		"AS ul ON tl.id = ul.list_id WHERE ul.user_id = ? AND ul.list_id = ?", todoListTable, usersListsTable)
	err := r.db.Get(&list, query, userId, listId)

	return list, err
}

func (r *TodoListMysql) Delete(userId, listId int) error {
	query := fmt.Sprintf("DELETE tl FROM %s tl INNER JOIN %s ul ON tl.id = ul.list_id WHERE ul.user_id=? AND ul.list_id=?", todoListTable, usersListsTable)
	_, err := r.db.Exec(query, userId, listId)

	return err
}

func (r *TodoListMysql) Update(userId, listId int, input todo.UpdateListInput) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if input.Title != nil {
		setValues = append(setValues, fmt.Sprintf("title=$%d", argId))
		args = append(args, *input.Title)
		argId++
	}

	if input.Description != nil {
		setValues = append(setValues, fmt.Sprintf("description=$%d", argId))
		args = append(args, *input.Description)
		argId++
	}

	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf("UPDATE %s tl INNER JOIN %s ul ON tl.id = ul.list_id SET %s WHERE ul.list_id=%d AND ul.user_id=%d", todoListTable, usersListsTable, setQuery, argId, argId+1)

	args = append(args, listId, userId)

	logrus.Debugf("updateQuery: %s", query)
	logrus.Debugf("args: %s", args)

	_, err := r.db.Exec(query, args...)
	return err
}
