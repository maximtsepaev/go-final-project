package db

import (
	"errors"
	"strings"
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func AddTask(task *Task) (int64, error) {
	var id int64

	if DB == nil {
		return 0, ErrDBNotInitialized
	}

	query := "INSERT INTO scheduler (date, title, comment, repeat) VALUES (:date, :title, :comment, :repeat)"
	res, err := DB.NamedExec(query, map[string]interface{}{
		"date":    task.Date,
		"title":   task.Title,
		"comment": task.Comment,
		"repeat":  task.Repeat,
	})
	if err == nil {
		id, err = res.LastInsertId()
	}

	return id, err
}

func Tasks(limit int) ([]*Task, error) {
	var tasks []*Task

	if DB == nil {
		return nil, ErrDBNotInitialized
	}

	query := "SELECT * FROM scheduler ORDER BY date LIMIT ?"
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

	if tasks == nil {
		tasks = []*Task{}
	}

	return tasks, nil
}

func SearchByKeyword(keyword string, limit int) ([]*Task, error) {
	var tasks []*Task

	if DB == nil {
		return nil, ErrDBNotInitialized
	}

	keyword = strings.ToLower(strings.TrimSpace(keyword))

	// Если ключевое слово пустое (пр. search == "    "), возвращаем все задачи
	if keyword == "" {
		return Tasks(limit)
	}

	query := "SELECT * FROM scheduler ORDER BY date LIMIT ?"
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

	if tasks == nil {
		tasks = []*Task{}
	}

	return tasks, nil
}

func SearchByDate(date string, limit int) ([]*Task, error) {
	var tasks []*Task

	if DB == nil {
		return nil, ErrDBNotInitialized
	}

	query := "SELECT * FROM scheduler WHERE date = :date LIMIT :limit"

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

	if tasks == nil {
		tasks = []*Task{}
	}

	return tasks, nil
}

func GetTaskByID(id string) (*Task, error) {
	var task Task

	if DB == nil {
		return nil, ErrDBNotInitialized
	}

	query := "SELECT * FROM scheduler WHERE id = ?"

	if err := DB.QueryRowx(query, id).StructScan(&task); err != nil {
		return nil, err
	}

	return &task, nil
}

func UpdateTask(task *Task) error {
	if DB == nil {
		return ErrDBNotInitialized
	}

	query := "UPDATE scheduler SET date = :date, title = :title, comment = :comment, repeat = :repeat WHERE id = :id"
	row, err := DB.NamedExec(query, map[string]interface{}{
		"id":      task.ID,
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

func UpdateDate(task *Task) error {
	if DB == nil {
		return ErrDBNotInitialized
	}

	query := "UPDATE scheduler SET date = :date WHERE id = :id"
	row, err := DB.NamedExec(query, map[string]interface{}{
		"date": task.Date,
		"id":   task.ID,
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

func DeleteTask(id string) error {
	if DB == nil {
		return ErrDBNotInitialized
	}

	query := "DELETE FROM scheduler WHERE id = ?"
	row, err := DB.Exec(query, id)
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
