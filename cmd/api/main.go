package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir("web")))
	err := http.ListenAndServe(":7540" , mux)
	if err != nil {
		log.Fatalf("The server could not be started due to an error: %v", err)
		return
	}
}