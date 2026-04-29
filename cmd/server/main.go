// Command server запускает веб-сервер планировщика задач.
package main

import (
	"log/slog"
	"os"

	"github.com/maximtsepaev/go-final-project/internal/db"
	"github.com/maximtsepaev/go-final-project/internal/logger"
	"github.com/maximtsepaev/go-final-project/internal/server"
)

// main инициализирует логгер, базу данных и стартует сервер.
func main() {
	// Инициализация логгера.
    // Сейчас используется простая реализация с выводом в консоль,
    // этого достаточно для моего проекта.
    // Код уже организован так, что в будущем можно легко и без проблем
    // расширить логирование (например, добавить запись в файл и настройки),
    // а также ввести больше уровней логов и логирование бизнес-логики.
	appLogger := logger.InitLogger()
	slog.SetDefault(appLogger)

	slog.Info("app starting")

	// Инициализация базы данных
	// В будущем будет использоваться PostgreSQL, сделаю после защиты
	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		dbFile = "./scheduler.db"
	}

	slog.Info("initializing", "db_file", dbFile) // Нет смысла использовать slog.With
	if err := db.InitDB(dbFile); err != nil {
		slog.Error("database init failed", "error", err)
		os.Exit(1)
	}

	// Запуск HTTP-сервера
	if err := server.Run(); err != nil {
		slog.Error("server start failed", "error", err)
		os.Exit(1)
	}
}
