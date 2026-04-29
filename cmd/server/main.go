// Command server запускает веб-сервер планировщика задач.
package main

import (
	"log"
	"os"

	"github.com/maximtsepaev/go-final-project/internal/db"
	"github.com/maximtsepaev/go-final-project/internal/server"
)

// main инициализирует базу данных и стартует сервер.
func main() {
	// Инициализация базы данных
	// В будущем будет использоваться PostgreSQL, сделаю после защиты
	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		dbFile = "./scheduler.db"
	}
	if err := db.InitDB(dbFile); err != nil {
		log.Fatal("error initializing database:", err)
	}

	// Запуск HTTP-сервера
	if err := server.Run(); err != nil {
		log.Fatal("error starting server:", err)
	}
}
