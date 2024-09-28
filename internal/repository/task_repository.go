package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/capybara120404/todo-list/internal/database"
	"github.com/capybara120404/todo-list/internal/utils"
)

type TaskRepository struct {
	db *sql.DB
}

func NewTaskRepository(connecter *database.Connecter) *TaskRepository {
	return &TaskRepository{
		db: connecter.DB,
	}
}

func (repository *TaskRepository) Add(task *Task) (int64, error) {
	err := isCorrect(task)
	if err != nil {
		return 0, err
	}

	res, err := repository.db.Exec("INSERT INTO scheduler (date, title, comment, repeat) VALUES (:date, :title, :comment, :repeat)",
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat),
	)
	if err != nil {
		return 0, fmt.Errorf("error inserting data into the database")
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("error retrieving last insert Id")
	}

	return id, nil
}

func (repository *TaskRepository) Change(id int, task *Task) error {
	err := isCorrect(task)
	if err != nil {
		return err
	}

	res, err := repository.db.Exec("UPDATE scheduler SET date = :date, title = :title, comment = :comment, repeat = :repeat WHERE id= :id",
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat),
		sql.Named("id", id))
	if err != nil {
		return fmt.Errorf("error updating task in the database")
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("error retrieving affected rows")
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no task found with the specified Id")
	}

	return nil
}

func (repository *TaskRepository) Complete(id int) error {
	row := repository.db.QueryRow("SELECT id, date, title, comment, repeat FROM scheduler WHERE id = :id", sql.Named("id", id))

	task, err := convertSqlToTask(row)
	if err != nil {
		return err
	}

	if task.Repeat == "" {
		_, err := repository.db.Exec("DELETE FROM scheduler WHERE id = :id",
			sql.Named("id", id))
		if err != nil {
			return fmt.Errorf("error deleting a task from the database")
		}
	} else {
		task.Date, err = utils.NextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			return err
		}

		result, err := repository.db.Exec("UPDATE scheduler SET date = :date WHERE id = :id",
			sql.Named("date", task.Date),
			sql.Named("id", id))
		if err != nil {
			return fmt.Errorf("error updating task in the database")
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return fmt.Errorf("error retrieving affected rows")
		}

		if rowsAffected == 0 {
			return fmt.Errorf("no task found with the specified Id")
		}
	}

	return nil
}

func (repository *TaskRepository) Delete(id int) error {
	_, err := repository.db.Exec("DELETE FROM scheduler WHERE id = :id",
		sql.Named("id", id))
	if err != nil {
		return fmt.Errorf("error deleting a task from the database")
	}

	return nil
}

func (repository *TaskRepository) GetAll() ([]Task, error) {
	rows, err := repository.db.Query("SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT 10")
	if err != nil {
		return nil, fmt.Errorf("error querying tasks from the database")
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		task := Task{}

		err := rows.Scan(&task.Id, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return nil, fmt.Errorf("error scanning task data")
		}

		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over task rows")
	}

	return tasks, nil
}

func (repository *TaskRepository) GetById(id int) (Task, error) {
	row := repository.db.QueryRow("SELECT id, date, title, comment, repeat FROM scheduler WHERE id = :id", sql.Named("id", id))

	task, err := convertSqlToTask(row)
	if err != nil {
		return Task{}, err
	}

	return task, err
}

func (repository *TaskRepository) CalculateNextDate(nowStr, dateStr, repeat string) (string, error) {
	now, err := time.Parse(utils.DateFormat, nowStr)
	if err != nil {
		return "", fmt.Errorf("invalid now format")
	}

	nextDate, err := utils.NextDate(now, dateStr, repeat)
	if err != nil {
		return "", err
	}

	return nextDate, nil
}
