package db

import (
	"database/sql"
	"errors"
	"strconv"
	"strings"
)

// Task описывает задачу планировщика.
type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

// AddTask сохраняет новую задачу и возвращает ее идентификатор.
func AddTask(task *Task) (int64, error) {
	var id int64

	query := "INSERT INTO scheduler (date, title, comment, repeat) VALUES ($1, $2, $3, $4) RETURNING id"

	err := DB.QueryRow(query, task.Date, task.Title, task.Comment, task.Repeat).Scan(&id)
	return id, err
}

// Tasks возвращает список задач, отсортированный по дате.
func Tasks(limit int) ([]*Task, error) {
	var tasks []*Task

	query := "SELECT CAST(id AS TEXT) AS id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT $1"

	rows, err := DB.Queryx(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var task Task

		if err := rows.StructScan(&task); err != nil {
			return nil, err
		}

		tasks = append(tasks, &task)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	if tasks == nil {
		tasks = []*Task{}
	}

	return tasks, nil
}

// SearchByKeyword ищет задачи по части заголовка или комментария.
func SearchByKeyword(keyword string, limit int) ([]*Task, error) {
	var tasks []*Task

	keyword = strings.ToLower(strings.TrimSpace(keyword))

	// Если ключевое слово пустое (пр. search == "    "), возвращаем все задачи
	if keyword == "" {
		return Tasks(limit)
	}

	query := "SELECT CAST(id AS TEXT) AS id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT $1"

	rows, err := DB.Queryx(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var task Task

		if err := rows.StructScan(&task); err != nil {
			return nil, err
		}

		title := strings.ToLower(task.Title)
		comment := strings.ToLower(task.Comment)

		if strings.Contains(title, keyword) || strings.Contains(comment, keyword) {
			tasks = append(tasks, &task)
		}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	if tasks == nil {
		tasks = []*Task{}
	}

	return tasks, nil
}

// SearchByDate ищет задачи по точной дате в формате DateFormat.
func SearchByDate(date string, limit int) ([]*Task, error) {
	var tasks []*Task

	query := "SELECT CAST(id AS TEXT) AS id, date, title, comment, repeat FROM scheduler WHERE date = :date LIMIT :limit"

	rows, err := DB.NamedQuery(query, map[string]interface{}{
		"date":  date,
		"limit": limit,
	})
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var task Task
		if err := rows.StructScan(&task); err != nil {
			return nil, err
		}
		tasks = append(tasks, &task)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	if tasks == nil {
		tasks = []*Task{}
	}

	return tasks, nil
}

// GetTaskByID возвращает задачу по идентификатору.
func GetTaskByID(id string) (*Task, error) {
	var task Task

	idVal, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, sql.ErrNoRows
	}

	query := "SELECT CAST(id AS TEXT) AS id, date, title, comment, repeat FROM scheduler WHERE id = $1"

	if err := DB.QueryRowx(query, idVal).StructScan(&task); err != nil {
		return nil, err
	}

	return &task, nil
}

// UpdateTask обновляет все поля задачи.
func UpdateTask(task *Task) error {
	idVal, err := strconv.ParseInt(task.ID, 10, 64)
	if err != nil {
		return errors.New("Неверный ID задачи")
	}

	query := "UPDATE scheduler SET date = :date, title = :title, comment = :comment, repeat = :repeat WHERE id = :id_int"

	row, err := DB.NamedExec(query, map[string]interface{}{
		"id_int":  idVal,
		"date":    task.Date,
		"title":   task.Title,
		"comment": task.Comment,
		"repeat":  task.Repeat,
	})
	if err != nil {
		return err
	}

	count, err := row.RowsAffected()
	if err != nil {
		return err
	}

	if count == 0 {
		return errors.New("Неверный ID задачи")
	}

	return nil
}

// UpdateDate обновляет только дату задачи.
func UpdateDate(task *Task) error {
	idVal, err := strconv.ParseInt(task.ID, 10, 64)
	if err != nil {
		return errors.New("Неверный ID задачи")
	}

	query := "UPDATE scheduler SET date = :date WHERE id = :id_int"

	row, err := DB.NamedExec(query, map[string]interface{}{
		"id_int": idVal,
		"date":   task.Date,
	})
	if err != nil {
		return err
	}

	count, err := row.RowsAffected()
	if err != nil {
		return err
	}

	if count == 0 {
		return errors.New("Неверный ID задачи")
	}

	return nil
}

// DeleteTask удаляет задачу по идентификатору.
func DeleteTask(id string) error {
	idVal, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return errors.New("Неверный ID задачи")
	}

	query := "DELETE FROM scheduler WHERE id = $1"

	row, err := DB.Exec(query, idVal)
	if err != nil {
		return err
	}

	count, err := row.RowsAffected()
	if err != nil {
		return err
	}

	if count == 0 {
		return errors.New("Неверный ID задачи")
	}

	return nil
}
