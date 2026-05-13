package api

import (
	"net/http"

	"github.com/maximtsepaev/go-final-project/internal/db"
)

// HandleDeleteTask удаляет задачу по идентификатору.
func HandleDeleteTask(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "id is required"})
		return
	}

	err := db.DeleteTask(id)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{})
}
