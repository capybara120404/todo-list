package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const dateFormat = "20060102"

func convertSqlToTask(row *sql.Row) (task, error) {
	var task task
	
	err := row.Scan(&task.Id, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		if err == sql.ErrNoRows {
			return task, fmt.Errorf("task not found")
		} else {
			return task, fmt.Errorf("error retrieving task from database")
		}
	}

	return task, nil
}

func getTaskFromBody(request *http.Request) (task, error) {
	var task task
	var buffer bytes.Buffer

	_, err := buffer.ReadFrom(request.Body)
	if err != nil {
		return task, fmt.Errorf("error reading request body")
	}

	err = json.Unmarshal(buffer.Bytes(), &task)
	if err != nil {
		return task, fmt.Errorf("invalid JSON format")
	}

	return task, nil
}

func getAndCheckId(request *http.Request) (int, error) {
	id, err := strconv.Atoi(request.URL.Query().Get("id"))
	if err != nil {
		return -1, fmt.Errorf("invalid task Id format")
	}

	if id <= 0 {
		return -1, fmt.Errorf("task Id must be greater than zero")
	}

	return id, nil
}

func isCorrect(w http.ResponseWriter, task *task, statusCode int) error {
	now := time.Now()

	if task.Title == "" {
		return fmt.Errorf("the title field should not be empty")
	}

	if task.Date == "" {
		task.Date = now.Format(dateFormat)
	}

	date, err := time.Parse(dateFormat, task.Date)
	if err != nil {
		writeJSONError(w, "invalid date format", statusCode)
		return fmt.Errorf("invalid date format")
	}

	if date.Format(dateFormat) < now.Format(dateFormat) {
		if task.Repeat == "" {
			task.Date = now.Format(dateFormat)
		} else {
			nextDate, err := nextDate(now, task.Date, task.Repeat)
			if err != nil {
				writeJSONError(w, err.Error(), statusCode)
				return err
			}

			task.Date = nextDate
		}
	}

	return nil
}

func writeJSONError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func nextDate(now time.Time, date string, repeat string) (string, error) {
	taskDate, err := time.Parse(dateFormat, date)
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
			return taskDate.Format(dateFormat)
		}
	}
}

func nextDateByYear(taskDate, now time.Time) string {
	for {
		taskDate = taskDate.AddDate(1, 0, 0)
		if taskDate.After(now) {
			return taskDate.Format(dateFormat)
		}
	}
}
