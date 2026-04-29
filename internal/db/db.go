// Package db содержит слой доступа к SQLite для планировщика задач.
package db

import (
	"errors"

	"github.com/jmoiron/sqlx"

	_ "modernc.org/sqlite"
)

// Схема базы данных
const sqlSchema = ` 
CREATE TABLE IF NOT EXISTS scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL DEFAULT "",
    title VARCHAR(256) NOT NULL,
    comment TEXT,
    repeat VARCHAR(128)
);

CREATE INDEX IF NOT EXISTS date_index ON scheduler (date);
`

// DB хранит подключение к базе данных SQLite.
var DB *sqlx.DB

// ErrDBNotInitialized возвращается, если работа с БД идет до инициализации.
var ErrDBNotInitialized = errors.New("database not initialized")

// InitDB открывает базу и создает таблицу при первом запуске.
func InitDB(dbFile string) error {
	var err error

	DB, err = sqlx.Open("sqlite", dbFile)
	if err != nil {
		return err
	}

	if _, err = DB.Exec(sqlSchema); err != nil {
		return err
	}

	return nil
}
