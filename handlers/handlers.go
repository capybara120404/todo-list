package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func NexDateHandler(w http.ResponseWriter, r *http.Request) {
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
