package api

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/maximtsepaev/go-final-project/internal/db"
)

// Функция для проверки и корректировки даты задачи
func checkDate(task *db.Task) error {
	var next string
	now := time.Now()
	today, err := time.Parse(DateFormat, now.Format(DateFormat))
	if err != nil {
		return err
	}

	if task.Date == "" {
		task.Date = today.Format(DateFormat)
		return nil
	}

	taskDate, err := time.Parse(DateFormat, task.Date)
	if err != nil {
		return InvalidDateError
	}

	if len(task.Repeat) != 0 {
		next, err = NextDate(today, task.Date, task.Repeat)
		if err != nil {
			return InvalidRepeatError
		}
	}

	if afterNow(today, taskDate) {
		if len(task.Repeat) == 0 {
			task.Date = today.Format(DateFormat)
		} else {
			task.Date = next
		}
	}

	return nil
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func HandleAddTask(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	var id int64
	var err error

	err = json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid JSON"})
		return
	}

	if task.Title == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Title is required"})
		return
	}

	err = checkDate(&task)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	id, err = db.AddTask(&task)
	if err != nil {
		log.Println(err.Error())
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Internal Server Error"})
		return
	}

	writeJSON(w, http.StatusCreated, map[string]interface{}{"id": id})
}
