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

type task struct {
	Id      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment,omitempty"`
	Repeat  string `json:"repeat"`
}

func (connecter *Connecter) GetTasksHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := connecter.db.Query("SELECT * FROM scheduler")
	if err != nil {
		writeJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer rows.Close()

	var tasks []task
	for rows.Next() {
		task := task{}

		err := rows.Scan(&task.Id, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			writeJSONError(w, err.Error(), http.StatusBadRequest)
			return
		}

		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		writeJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	if tasks == nil{
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"tasks": make([]task, 0)})
	} else {
		tasks = sortTasksByDate(tasks)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"tasks": tasks})	
	}
}

func sortTasksByDate(tasks []task) []task {
	for i := 0; i < len(tasks)-1; i++ {
		for j := i + 1; j < len(tasks); j++ {
			if tasks[i].Date > tasks[j].Date {
				tasks[i], tasks[j] = tasks[j], tasks[i]
			}
		}
	}
	return tasks
}

func (connecter *Connecter) AddTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task task
	var buffer bytes.Buffer
	now := time.Now()

	_, err := buffer.ReadFrom(r.Body)
	if err != nil {
		writeJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(buffer.Bytes(), &task)
	if err != nil {
		writeJSONError(w, "invalid JSON format", http.StatusBadRequest)
		return
	}

	if task.Title == "" {
		writeJSONError(w, "title is required", http.StatusBadRequest)
		return
	}

	if task.Date == "" {
		task.Date = now.Format("20060102")
	} else {
		date, err := time.Parse("20060102", task.Date)
		if err != nil {
			writeJSONError(w, "invalid date format", http.StatusBadRequest)
			return
		}

		if date.Before(now) {
			if task.Repeat != "" {
				nextDate, err := nextDate(now, task.Date, task.Repeat)
				if err != nil {
					writeJSONError(w, err.Error(), http.StatusBadRequest)
					return
				}
				task.Date = nextDate
			} else {
				task.Date = now.Format("20060102")
			}
		}
	}

	res, err := connecter.db.Exec("INSERT INTO scheduler (date, title, comment, repeat) VALUES (:date, :title, :comment, :repeat)",
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat),
	)
	if err != nil {
		writeJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := res.LastInsertId()
	if err != nil {
		writeJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(map[string]any{"id": id})
}

func writeJSONError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func (connecter *Connecter) NexDateHandler(w http.ResponseWriter, r *http.Request) {
	nowStr := r.FormValue("now")
	dateStr := r.FormValue("date")
	repeat := r.FormValue("repeat")

	now, err := time.Parse("20060102", nowStr)
	if err != nil {
		http.Error(w, "invalid now format", http.StatusBadRequest)
		return
	}

	nextDate, err := nextDate(now, dateStr, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Fprintln(w, nextDate)
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
