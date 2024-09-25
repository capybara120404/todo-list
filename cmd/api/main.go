package main

import (
	"log"
	"net/http"

	"github.com/capybara120404/todo-list/configs"
	"github.com/capybara120404/todo-list/handlers"
	"github.com/go-chi/chi/v5"
)

func main() {
	connecter, err := handlers.OpenOrCreate(configs.PathToDB)
	if err != nil {
		log.Printf("error opening or creating database: %v", err)
		return
	}
	defer connecter.Close()

		router := chi.NewRouter()
		fs := http.FileServer(http.Dir("web"))
		router.Handle("/*", http.StripPrefix("/", fs))
		router.Get("/api/nextdate", connecter.NexDateHandler)
		router.Post("/api/task", connecter.AddTaskHandler)
		router.Get("/api/tasks", connecter.GetTasksHandler)
		router.Get("/api/task", connecter.GetTaskByIdHandler)
		router.Put("/api/task", connecter.ChangeTaskHandler)

	err = http.ListenAndServe(configs.Addr, router)
	if err != nil {
		log.Printf("the server could not be started due to an error: %v", err)
		return
	}
}
