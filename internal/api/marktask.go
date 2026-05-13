package api

import (
	"net/http"
	"time"

	"github.com/maximtsepaev/go-final-project/internal/db"
)

// HandleMarkTaskAsDone отмечает задачу как выполненную и обновляет дату при повторе.
func HandleMarkTaskAsDone(w http.ResponseWriter, r *http.Request) {
	var err error
	var task *db.Task

	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "id is required"})
		return
	}

	task, err = db.GetTaskByID(id)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "task not found"})
		return
	}

	if task.Repeat == "" || task.Repeat == "NULL" {
		err = db.DeleteTask(id)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, map[string]interface{}{})
		return
	} else {
		task.Date, err = NextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
	}

	err = db.UpdateDate(task)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{})
}
