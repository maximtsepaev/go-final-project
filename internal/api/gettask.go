package api

import (
	"net/http"

	"github.com/maximtsepaev/go-final-project/internal/db"
)

// HandleGetTaskByID возвращает задачу по идентификатору.
func HandleGetTaskByID(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "id is required"})
		return
	}

	task, err := db.GetTaskByID(id)
	if err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "task not found"})
		return
	}

	writeJSON(w, http.StatusOK, task)
}
