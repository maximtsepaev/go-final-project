// Package db содержит слой доступа к PostgreSQL для планировщика задач.
package db

import (
	"errors"
	"log/slog"

	"github.com/jmoiron/sqlx"

	_ "github.com/lib/pq"
)

// Схема базы данных
const sqlSchema = ` 
CREATE TABLE IF NOT EXISTS scheduler (
    id SERIAL PRIMARY KEY,
    date CHAR(8) NOT NULL DEFAULT '',
    title VARCHAR(256) NOT NULL,
    comment TEXT,
    repeat VARCHAR(128)
);

CREATE INDEX IF NOT EXISTS date_index ON scheduler (date);
`

// DB хранит подключение к базе данных PostgreSQL.
var DB *sqlx.DB

// ErrDBNotInitialized возвращается, если работа с БД идет до инициализации.
var ErrDBNotInitialized = errors.New("database not initialized")

// InitDB открывает базу и создает таблицу при первом запуске.
func InitDB(dsn string) error {
	var err error

	DB, err = sqlx.Open("postgres", dsn)
	if err != nil {
		return err
	}

	if err = DB.Ping(); err != nil {
		_ = DB.Close()
		return err
	}

	if _, err = DB.Exec(sqlSchema); err != nil {
		_ = DB.Close()
		return err
	}

	slog.Info("database initialized successfully")

	return nil
}

// CloseDB закрывает подключение к базе данных.
func CloseDB() error {
	if DB == nil {
		return ErrDBNotInitialized
	}

	return DB.Close()
}
