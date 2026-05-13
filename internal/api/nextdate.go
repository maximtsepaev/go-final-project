package api

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var (
	// EmptyStringError возвращается, когда правило повторения пустое.
	EmptyStringError = errors.New("empty string")
	// InvalidRepeatError возвращается при некорректном формате повтора.
	InvalidRepeatError = errors.New("invalid repeat format")
	// InvalidDateError возвращается при некорректном формате даты.
	InvalidDateError = errors.New("invalid date format")
)

// afterNow проверяет, наступила ли дата после текущего времени.
func afterNow(date, now time.Time) bool {
	return date.After(now)
}

// NextDate вычисляет следующую дату по правилу повторения.
func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	var parts []string
	var date time.Time

	if repeat == "" {
		return "", EmptyStringError
	}

	date, err := time.Parse(DateFormat, dstart)
	if err != nil {
		return "", err
	}

	parts = strings.Split(repeat, " ")

	switch parts[0] {
	case "d":
		if len(parts) < 2 {
			return "", InvalidRepeatError
		}

		interval, err := strconv.Atoi(parts[1])
		if err != nil {
			return "", InvalidRepeatError
		}

		if interval > 400 {
			return "", errors.New("interval for daily repeat is too large")
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

	case "w":
		if len(parts) < 2 {
			return "", InvalidRepeatError
		}

		days := []time.Weekday{}
		slice := strings.Split(parts[1], ",")

		for _, day := range slice {
			dayInt, err := strconv.Atoi(day)
			if err != nil {
				return "", InvalidRepeatError
			}

			if dayInt < 1 || dayInt > 7 {
				return "", InvalidRepeatError
			}

			days = append(days, time.Weekday(dayInt%7))
		}

		for {
			date = date.AddDate(0, 0, 1)

			for _, day := range days {
				if date.Weekday() == day && afterNow(date, now) {
					return date.Format(DateFormat), nil
				}
			}
		}

	case "m":
		var days [32]bool
		var months [13]bool
		var lastDay, preLastDay, useMonths bool

		if len(parts) < 2 {
			return "", InvalidRepeatError
		}

		dlist := strings.Split(parts[1], ",")
		for _, d := range dlist {
			dayInt, err := strconv.Atoi(d)
			if err != nil {
				return "", InvalidRepeatError
			}

			switch {
			case dayInt == -1:
				lastDay = true
			case dayInt == -2:
				preLastDay = true
			case dayInt >= 1 && dayInt <= 31:
				days[dayInt] = true
			default:
				return "", InvalidRepeatError
			}
		}

		if len(parts) == 3 {
			useMonths = true
			mlist := strings.Split(parts[2], ",")

			for _, m := range mlist {
				monthInt, err := strconv.Atoi(m)
				if err != nil {
					return "", InvalidRepeatError
				}

				if monthInt > 0 && monthInt < 13 {
					months[monthInt] = true
				} else {
					return "", InvalidRepeatError
				}
			}
		}

		for {
			date = date.AddDate(0, 0, 1)
			day := date.Day()
			last := time.Date(date.Year(), date.Month()+1, 0, 0, 0, 0, 0, time.UTC).Day()

			if !afterNow(date, now) {
				continue
			}

			if useMonths && !months[int(date.Month())] {
				continue
			}

			if days[day] || (lastDay && day == last) || (preLastDay && day == last-1) {
				return date.Format(DateFormat), nil
			}
		}

	default:
		return "", errors.New("unknown repeat type")
	}
}

// HandleNextDate отвечает следующей датой в текстовом виде.
func HandleNextDate(w http.ResponseWriter, r *http.Request) {
	date := r.FormValue("date")
	repeat := r.FormValue("repeat")
	now := r.FormValue("now")
	if now == "" {
		now = time.Now().Format(DateFormat)
	}
	nowTime, _ := time.Parse(DateFormat, now)

	nextDate, err := NextDate(nowTime, date, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte(nextDate)); err != nil {
		http.Error(w, "write nextdate response failed", http.StatusInternalServerError)
	}
}
