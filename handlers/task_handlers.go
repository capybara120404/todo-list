package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type task struct {
	Id      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment,omitempty"`
	Repeat  string `json:"repeat"`
}

func (connecter *connecter) DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id, err := getAndCheckId(r)
	if err != nil {
		writeJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = connecter.db.Exec("DELETE FROM scheduler WHERE id = :id",
		sql.Named("id", id))
	if err != nil {
		writeJSONError(w, "error deleting a task from the database", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{})
}

func (connecter *connecter) MarkAsCompletedHandler(w http.ResponseWriter, r *http.Request) {
	id, err := getAndCheckId(r)
	if err != nil {
		writeJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	row := connecter.db.QueryRow("SELECT id, date, title, comment, repeat FROM scheduler WHERE id = :id", sql.Named("id", id))

	task, err := convertSqlToTask(row)
	if err != nil {
		writeJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	if task.Repeat == "" {
		_, err := connecter.db.Exec("DELETE FROM scheduler WHERE id = :id",
			sql.Named("id", id))
		if err != nil {
			writeJSONError(w, "error deleting a task from the database", http.StatusBadRequest)
			return
		}
	} else {
		task.Date, err = nextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			writeJSONError(w, err.Error(), http.StatusBadRequest)
			return
		}

		result, err := connecter.db.Exec("UPDATE scheduler SET date = :date WHERE id = :id",
			sql.Named("date", task.Date),
			sql.Named("id", id))
		if err != nil {
			writeJSONError(w, "error updating task in the database", http.StatusBadRequest)
			return
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			writeJSONError(w, "error retrieving affected rows", http.StatusBadRequest)
			return
		}

		if rowsAffected == 0 {
			writeJSONError(w, "no task found with the specified Id", http.StatusNotFound)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{})
}

func (connecter *connecter) ChangeTaskHandler(w http.ResponseWriter, r *http.Request) {
	task, err := getTaskFromBody(r)
	if err != nil {
		writeJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = isCorrect(w, &task, http.StatusBadRequest)
	if err != nil {
		writeJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := getAndCheckId(r)
	if err != nil {
		writeJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := connecter.db.Exec("UPDATE scheduler SET date = :date, title = :title, comment = :comment, repeat = :repeat WHERE id= :id",
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat),
		sql.Named("id", id))
	if err != nil {
		writeJSONError(w, "error updating task in the database", http.StatusBadRequest)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		writeJSONError(w, "error retrieving affected rows", http.StatusBadRequest)
		return
	}

	if rowsAffected == 0 {
		writeJSONError(w, "no task found with the specified Id", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{})
}

func (connecter *connecter) GetTaskByIdHandler(w http.ResponseWriter, r *http.Request) {
	id, err := getAndCheckId(r)
	if err != nil {
		writeJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	row := connecter.db.QueryRow("SELECT id, date, title, comment, repeat FROM scheduler WHERE id = :id", sql.Named("id", id))

	task, err := convertSqlToTask(row)
	if err != nil {
		writeJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

func (connecter *connecter) GetTasksHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := connecter.db.Query("SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT 10")
	if err != nil {
		writeJSONError(w, "error querying tasks from the database", http.StatusBadRequest)
		return
	}
	defer rows.Close()

	var tasks []task
	for rows.Next() {
		task := task{}

		err := rows.Scan(&task.Id, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			writeJSONError(w, "error scanning task data", http.StatusBadRequest)
			return
		}

		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		writeJSONError(w, "error iterating over task rows", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if tasks == nil {
		json.NewEncoder(w).Encode(map[string]any{"tasks": make([]task, 0)})
	} else {
		json.NewEncoder(w).Encode(map[string]any{"tasks": tasks})
	}
}

func (connecter *connecter) AddTaskHandler(w http.ResponseWriter, r *http.Request) {
	task, err := getTaskFromBody(r)
	if err != nil {
		writeJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = isCorrect(w, &task, http.StatusBadRequest)
	if err != nil {
		writeJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	res, err := connecter.db.Exec("INSERT INTO scheduler (date, title, comment, repeat) VALUES (:date, :title, :comment, :repeat)",
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat),
	)
	if err != nil {
		writeJSONError(w, "error inserting data into the database", http.StatusBadRequest)
		return
	}

	id, err := res.LastInsertId()
	if err != nil {
		writeJSONError(w, "error retrieving last insert ID", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{"id": id})
}

func (connecter *connecter) NexDateHandler(w http.ResponseWriter, r *http.Request) {
	nowStr := r.FormValue("now")
	dateStr := r.FormValue("date")
	repeat := r.FormValue("repeat")

	now, err := time.Parse(dateFormat, nowStr)
	if err != nil {
		writeJSONError(w, "invalid now format", http.StatusBadRequest)
		return
	}

	nextDate, err := nextDate(now, dateStr, repeat)
	if err != nil {
		writeJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Fprintln(w, nextDate)
}
