package api

import (
	"net/http"
	"time"

	"github.com/maximtsepaev/go-final-project/internal/db"
)

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

const maxTasks = 50

func HandleGetTasks(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()

	if params.Get("search") != "" {
		param, err := time.Parse("02.01.2006", params.Get("search"))

		// Если не удалось распарсить дату, ищем по ключевому слову
		if err != nil {
			tasks, err := db.SearchByKeyword(params.Get("search"), maxTasks)
			if err != nil {
				writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
				return
			}

			writeJSON(w, http.StatusOK, TasksResp{
				Tasks: tasks,
			})

			return
		}

		// Если дата распарсилась, ищем по дате
		tasks, err := db.SearchByDate(param.Format(DateFormat), maxTasks)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}

		writeJSON(w, http.StatusOK, TasksResp{
			Tasks: tasks,
		})

		return
	}

	// Если параметр search не указан, возвращаем все задачи
	tasks, err := db.Tasks(maxTasks)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, TasksResp{
		Tasks: tasks,
	})
}
