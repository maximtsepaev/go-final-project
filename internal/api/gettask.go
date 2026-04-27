package api

import (
	"net/http"

	"github.com/maximtsepaev/go-final-project/internal/db"
)

func HandleGetTaskByID(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Не указан идентификатор"})
		return
	}

	task, err := db.GetTaskByID(id)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Задача не найдена"})
		return
	}

	writeJSON(w, http.StatusOK, task)
}
