package main

import (
	"log"
	userHandler "my-burger-auth/internal/handler"
	"net/http"
)

func main() {

	mux := http.NewServeMux()
	mux.HandleFunc("/users/create", userHandler.HandleUserCreation)

	log.Println("Starting server on :8080")

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
