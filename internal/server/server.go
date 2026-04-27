package server

import (
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/maximtsepaev/go-final-project/internal/api"
)

func Run() error {
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "7540"
	}

	r := chi.NewRouter() // Роутер для обработки HTTP-запросов
	api.Init(r)

	if err := http.ListenAndServe(":"+port, r); err != nil {
		return err
	}

	return nil
}
