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
	result, err := r.db.Exec(query, user.Name, user.Username, user.Password)
	if err != nil {
		return 0, err
	}

	lastID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(lastID), nil
}

func (r *AuthMySql) GetUser(username, password string) (todo.User, error) {
	var user todo.User
	query := fmt.Sprintf("SELECT id FROM %s WHERE username=? AND password_hash=?", userTable)
	err := r.db.Get(&user, query, username, password)

	return user, err
}
