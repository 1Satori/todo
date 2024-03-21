package mysql

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)

type Config struct {
	Username string
	Password string
	Port     string
	Host     string
	DBName   string
}

func NewMySqlDB(cfg *Config) (*sqlx.DB, error) {
	db, err := sqlx.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.DBName))
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
