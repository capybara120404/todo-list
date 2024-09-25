package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func isCorrect(w http.ResponseWriter, task *task, statusCode int) bool {
	var res bool
	now := time.Now()

	if task.Title == "" {
		res = false
		writeJSONError(w, "the title field should not be empty", statusCode)
		return false
	}

	if task.Date == "" {
		task.Date = now.Format("20060102")
	}

	date, err := time.Parse("20060102", task.Date)
	if err != nil {
		res = false
		writeJSONError(w, "invalid date format", statusCode)
		return res
	}

	if date.Format("20060102") < now.Format("20060102") {
		if task.Repeat == "" {
			task.Date = now.Format("20060102")
		} else {
			nextDate, err := nextDate(now, task.Date, task.Repeat)
			if err != nil {
				res = false
				writeJSONError(w, err.Error(), statusCode)
				return res
			}

			task.Date = nextDate
		}
	}

	res = true
	return res
}

func writeJSONError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func nextDate(now time.Time, date string, repeat string) (string, error) {
	taskDate, err := time.Parse("20060102", date)
	if err != nil {
		return "", fmt.Errorf("invalid date format: %v", err)
	}

	if repeat == "" {
		return "", nil
	}

	parts := strings.Split(repeat, " ")
	if len(parts) == 0 {
		return "", fmt.Errorf("empty repeat rule")
	}

	rule := parts[0]

	if rule == "d" {
		if len(parts) != 2 {
			return "", fmt.Errorf("invalid 'd' rule format")
		}
		days, err := strconv.Atoi(parts[1])
		if err != nil || days < 1 || days > 400 {
			return "", fmt.Errorf("invalid number of days")
		}
		return nextDateByDays(taskDate, now, days), nil
	} else if rule == "y" {
		if len(parts) != 1 {
			return "", fmt.Errorf("invalid 'y' rule format")
		}
		return nextDateByYear(taskDate, now), nil
	} else {
		return "", fmt.Errorf("unsupported format")
	}
}

func nextDateByDays(taskDate, now time.Time, days int) string {
	for {
		taskDate = taskDate.AddDate(0, 0, days)
		if taskDate.After(now) {
			return taskDate.Format("20060102")
		}
	}
}

func nextDateByYear(taskDate, now time.Time) string {
	for {
		taskDate = taskDate.AddDate(1, 0, 0)
		if taskDate.After(now) {
			return taskDate.Format("20060102")
		}
	}
}
