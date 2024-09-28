package repository

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/capybara120404/todo-list/internal/utils"
)

type Task struct {
	Id      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment,omitempty"`
	Repeat  string `json:"repeat"`
}

func GetTaskFromBody(request *http.Request) (Task, error) {
	var task Task
	var buffer bytes.Buffer

	_, err := buffer.ReadFrom(request.Body)
	if err != nil {
		return Task{}, fmt.Errorf("error reading request body")
	}

	err = json.Unmarshal(buffer.Bytes(), &task)
	if err != nil {
		return Task{}, fmt.Errorf("invalid JSON format")
	}

	return task, nil
}

func convertSqlToTask(row *sql.Row) (Task, error) {
	var task Task

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

func isCorrect(task *Task) error {
	now := time.Now()

	if task.Title == "" {
		return fmt.Errorf("the title field should not be empty")
	}

	if task.Date == "" {
		task.Date = now.Format(utils.DateFormat)
	}

	date, err := time.Parse(utils.DateFormat, task.Date)
	if err != nil {
		return fmt.Errorf("invalid date format")
	}

	if date.Format(utils.DateFormat) < now.Format(utils.DateFormat) {
		if task.Repeat == "" {
			task.Date = now.Format(utils.DateFormat)
		} else {
			nextDate, err := utils.NextDate(now, task.Date, task.Repeat)
			if err != nil {
				return err
			}

			task.Date = nextDate
		}
	}

	return nil
}
