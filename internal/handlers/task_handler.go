package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/capybara120404/todo-list/internal/repository"
	"github.com/capybara120404/todo-list/internal/utils"
)

type taskHandler struct {
	repository *repository.TaskRepository
}

func NewTaskHandler(repository *repository.TaskRepository) *taskHandler {
	return &taskHandler{
		repository: repository,
	}
}

func (handler *taskHandler) AddTaskHandler(w http.ResponseWriter, r *http.Request) {
	task, err := repository.GetTaskFromBody(r)
	if err != nil {
		utils.WriteJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := handler.repository.Add(&task)
	if err != nil {
		utils.WriteJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{"id": id})
}

func (handler *taskHandler) ChangeTaskHandler(w http.ResponseWriter, r *http.Request) {
	task, err := repository.GetTaskFromBody(r)
	if err != nil {
		utils.WriteJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(task.Id)
	if err != nil {
		utils.WriteJSONError(w, "invalid task Id format", http.StatusBadRequest)
		return
	}

	err = handler.repository.Change(id, &task)
	if err != nil {
		utils.WriteJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{})
}

func (handler *taskHandler) CompleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id, err := utils.GetAndCheckId(r)
	if err != nil {
		utils.WriteJSONError(w, "invalid task Id format", http.StatusBadRequest)
		return
	}

	err = handler.repository.Complete(id)
	if err != nil {
		utils.WriteJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{})
}

func (handler *taskHandler) DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id, err := utils.GetAndCheckId(r)
	if err != nil {
		utils.WriteJSONError(w, "invalid task Id format", http.StatusBadRequest)
		return
	}

	err = handler.repository.Delete(id)
	if err != nil {
		utils.WriteJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{})
}

func (handler *taskHandler) GetAllTasksHandler(w http.ResponseWriter, r *http.Request) {
	tasks, err := handler.repository.GetAll()
	if err != nil {
		utils.WriteJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if tasks == nil {
		json.NewEncoder(w).Encode(map[string]any{"tasks": make([]repository.Task, 0)})
	} else {
		json.NewEncoder(w).Encode(map[string]any{"tasks": tasks})
	}
}

func (handler *taskHandler) GetTaskByIdHandler(w http.ResponseWriter, r *http.Request) {
	id, err := utils.GetAndCheckId(r)
	if err != nil {
		utils.WriteJSONError(w, "invalid task Id format", http.StatusBadRequest)
		return
	}

	task, err := handler.repository.GetById(id)
	if err != nil {
		utils.WriteJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

func (handler *taskHandler) NexDateHandler(w http.ResponseWriter, r *http.Request) {
	nowStr := r.FormValue("now")
	dateStr := r.FormValue("date")
	repeat := r.FormValue("repeat")

	nextDate, err := handler.repository.CalculateNextDate(nowStr, dateStr, repeat)
	if err != nil {
		utils.WriteJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Fprintln(w, nextDate)
}
