package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/capybara120404/todo-list/database"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("configs/.env")
	if err != nil {
		log.Printf("an error occurred while importing configuration files: %v", err)
		return
	}

	connecter, err := database.OpenOrCreate()
	if err != nil {
		log.Printf("error opening or creating database: %v", err)
		return
	}
	defer connecter.Close()
	
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir("web")))
	addr := fmt.Sprintf(":%s", os.Getenv("TODO_PORT"))

	err = http.ListenAndServe(addr, mux)
	if err != nil {
		log.Printf("the server could not be started due to an error: %v", err)
		return
	}
}
