package main

import (
	"log"
	"os"

	"github.com/maximtsepaev/go-final-project/internal/db"
	"github.com/maximtsepaev/go-final-project/internal/server"
	// "github.com/go-chi/chi/v5" в будущем, понадобится более сложная маршрутизация
)

func main() {
	// Инициализация базы данных
	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		dbFile = "./scheduler.db"
	}
	if err := db.InitDB(dbFile); err != nil {
		log.Fatal("Error initializing database:", err)
	}

	// Запуск HTTP-сервера
	if err := server.Run(); err != nil {
		log.Fatal("Error starting server:", err)
	}

	// СДЕЛАТЬ ЗАДАНИЕ СО ЗВЕЗДОЧКОЙ NEXTDATE
}
