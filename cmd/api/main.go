package main

import (
	"log"
	"net/http"

	"github.com/capybara120404/todo-list/internal/configs"
	"github.com/capybara120404/todo-list/internal/database"
	"github.com/capybara120404/todo-list/internal/handlers"
	"github.com/capybara120404/todo-list/internal/repository"
	"github.com/go-chi/chi/v5"
)

func main() {
	connecter, err := database.OpenOrCreate(configs.PathToDB)
	if err != nil {
		log.Printf("%v", err)
		return
	}
	defer connecter.Close()

	repository := repository.NewTaskRepository(connecter)
	handlers := handlers.NewTaskHandler(repository)

	router := chi.NewRouter()

	fs := http.FileServer(http.Dir("web"))

	router.Handle("/*", http.StripPrefix("/", fs))
	router.Get("/api/nextdate", handlers.NexDateHandler)
	router.Get("/api/tasks", handlers.GetAllTasksHandler)
	router.Get("/api/task", handlers.GetTaskByIdHandler)
	router.Post("/api/task", handlers.AddTaskHandler)
	router.Post("/api/task/done", handlers.CompleteTaskHandler)
	router.Put("/api/task", handlers.ChangeTaskHandler)
	router.Delete("/api/task", handlers.DeleteTaskHandler)

	log.Printf("The server start at port: %s", configs.Addr)

	err = http.ListenAndServe(configs.Addr, router)
	if err != nil {
		log.Printf("The server could not be started due to an error: %v", err)
		return
	}
}
