package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

type AuthenticationRequest struct {
	Cpf      string `json:"cpf"`
	Password int64  `json:"password"`
}

type AuthenticationResponse struct {
	Token string `json:"access_token"`
}

func Main() {

	err := godotenv.Load("../../../.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/users/auth", handleAuthentication)

	log.Println("Starting server on :8090")

	if err := http.ListenAndServe(":8090", mux); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}

func handleAuthentication(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r)
	var request AuthenticationRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	fmt.Println(err)
	if err != nil {
		http.Error(w, "failed to decode request", http.StatusBadRequest)
		return
	}

	response, err := authenticate(request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func authenticate(request AuthenticationRequest) (AuthenticationResponse, error) {
	url := os.Getenv("MY_BURGER_APP_URL")
	jsonData, err := json.Marshal(request)

	if err != nil {
		return AuthenticationResponse{}, fmt.Errorf("failed to serialize request: %w", err)
	}

	response, err := http.Post(fmt.Sprintf("%s/auth", url), "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return AuthenticationResponse{}, fmt.Errorf("failed to authenticate user: %w", err)
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return AuthenticationResponse{}, fmt.Errorf("authentication failed,  status code: %d", response.StatusCode)
	}

	var authResponse AuthenticationResponse
	err = json.NewDecoder(response.Body).Decode(&authResponse)
	if err != nil {
		return AuthenticationResponse{}, fmt.Errorf("failed to decode response: %w", err)
	}

	return authResponse, nil
}
