package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("configs/.env")
	if err != nil {
		log.Printf("An error occurred while importing configuration files")
		return
	}

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir("web")))
	addr := fmt.Sprintf(":%s", os.Getenv("TODO_PORT"))

	err = http.ListenAndServe(addr, mux)
	if err != nil {
		log.Printf("The server could not be started due to an error")
		return
	}
}
