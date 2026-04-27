package api

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var EmptyStringError = errors.New("empty string")
var InvalidRepeatError = errors.New("invalid repeat format")
var InvalidDateError = errors.New("invalid date format")

// Функция для определения, находится ли дата после текущей
func afterNow(date, now time.Time) bool {
	return date.After(now)
}

// Функция для вычисления следующей даты повторения задачи
func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	var repeatSplit []string
	var date time.Time

	if repeat == "" {
		return "", EmptyStringError
	}

	date, err := time.Parse(DateFormat, dstart)
	if err != nil {
		return "", err
	}

	repeatSplit = strings.Split(repeat, " ")

	switch repeatSplit[0] {
	case "d":
		if len(repeatSplit) < 2 {
			return "", InvalidRepeatError
		}

		interval, _ := strconv.Atoi(repeatSplit[1])

		if interval > 400 {
			return "", fmt.Errorf("interval for daily repeat is too large")
		}

		for {
			date = date.AddDate(0, 0, interval)
			if afterNow(date, now) {
				break
			}
		}

		return date.Format(DateFormat), nil

	case "y":
		for {
			date = date.AddDate(1, 0, 0)
			if afterNow(date, now) {
				break
			}
		}
		return date.Format(DateFormat), nil

	default:
		return "", fmt.Errorf("unknown repeat type: %s", repeatSplit[0])
	}
}

func HandleNextDate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	date := r.FormValue("date")
	repeat := r.FormValue("repeat")
	now := r.FormValue("now")
	if now == "" {
		now = time.Now().Format(DateFormat)
	}
	nowTime, _ := time.Parse(DateFormat, now)

	nextDate, err := NextDate(nowTime, date, repeat)
	if err != nil {
		// ДОБАВИТЬ УДАЛЕНИЕ ИЗ БД
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(nextDate))
}

// ПРОВЕРИТЬ ПОКРЫТИЕ ОШИБОК ПО ТЗ
