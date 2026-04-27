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

var DB *sqlx.DB                                                  // Глобальная переменная для доступа к базе данных
var ErrDBNotInitialized = errors.New("Database not initialized") // Ошибка для случая, когда база данных не инициализирована

// Инициализация базы данных, создание таблицы при первом запуске
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
