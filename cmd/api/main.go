package main

import (
	"log"
	"net/http"

	"github.com/capybara120404/todo-list/configs"
	"github.com/capybara120404/todo-list/database"
)

func main() {
	connecter, err := database.OpenOrCreate(configs.PathToDB)
	if err != nil {
		log.Printf("error opening or creating database: %v", err)
		return
	}
	defer connecter.Close()

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir("web")))

	err = http.ListenAndServe(configs.PathToDB, mux)
	if err != nil {
		log.Printf("the server could not be started due to an error: %v", err)
		return
	}
}
