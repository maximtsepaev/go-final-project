package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

const webDir = "./web" // Путь к директории с веб-ресурсами (HTML, CSS, JS)

const DateFormat = "20060102" // Формат даты для хранения и обработки (ГГГГММДД)

func Init(r chi.Router) {
	r.Get("/api/nextdate", HandleNextDate)
	r.Get("/api/tasks", HandleGetTasks)

	// Роут для разных методов на одном пути
	r.Route("/api/task", func(r chi.Router) {
		r.Post("/", HandleAddTask)
		r.Get("/", HandleGetTaskByID)
		r.Put("/", HandleUpdateTask)
		r.Post("/done", HandleMarkTaskAsDone)
		r.Delete("/", HandleDeleteTask)
	})

	r.Handle("/*", http.FileServer(http.Dir(webDir)))
}
