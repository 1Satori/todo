package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	todo "todo-app"
)

type AuthMySql struct {
	db *sqlx.DB
}

func NewAuthMySql(db *sqlx.DB) *AuthMySql {
	return &AuthMySql{db: db}
}

func (r *AuthMySql) CreateUser(user todo.User) (int, error) {
	query := fmt.Sprintf("INSERT INTO %s (name, username, password_hash) VALUES (?, ?, ?)", userTable)
	row, err := r.db.Exec(query, user.Name, user.Username, user.Password)
	if err != nil {
		return 0, err
	}

	lastId, err := row.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(lastId), nil
}
