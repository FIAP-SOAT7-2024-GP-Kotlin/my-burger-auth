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

type UserCreationRequest struct {
	Cpf      string   `json:"cpf"`
	Password int64    `json:"password"`
	Role     UserRole `json:"role"`
}

type UserRole string

const (
	UserRoleAdmin UserRole = "ADMIN"
	UserRoleUser  UserRole = "USER"
)

type UserCreationResponse struct {
	Id   string `json:"id"`
	Cpf  string `json:"cpf"`
	Role string `json:"role"`
}

type AuthenticationRequest struct {
	Cpf      string `json:"cpf"`
	Password int64  `json:"password"`
}

type AuthenticationResponse struct {
	Token string `json:"access_token"`
}

func main() {

	err := godotenv.Load("../../../.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/users/auth", handleUserCreation)

	log.Println("Starting server on :8090")

	if err := http.ListenAndServe(":8090", mux); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}

func handleUserCreation(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r)
	var request UserCreationRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	fmt.Println(err)
	if err != nil {
		http.Error(w, "failed to decode request", http.StatusBadRequest)
		return
	}

	token, err := createUserAndAuthenticate(request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(AuthenticationResponse{Token: token})
}

func createUserAndAuthenticate(request UserCreationRequest) (string, error) {
	user, err := createUser(request)
	if err != nil {
		return "", fmt.Errorf("failed to create user: %w", err)
	}

	authenticationRequest := AuthenticationRequest{
		Cpf:      user.Cpf,
		Password: request.Password,
	}

	authenticationResponse, err := authenticateUser(authenticationRequest)
	if err != nil {
		return "", fmt.Errorf("failed to authenticate user: %w", err)
	}
	return authenticationResponse, nil
}

func createUser(request UserCreationRequest) (*UserCreationResponse, error) {
	url := os.Getenv("MY_BURGER_APP_URL")
	jsonData, err := json.Marshal(request)

	if err != nil {
		return nil, fmt.Errorf("failed to serialize request: %w", err)
	}

	response, err := http.Post(fmt.Sprintf("%s/users", url), "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	var userResponse UserCreationResponse
	err = json.NewDecoder(response.Body).Decode(&userResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	return &userResponse, response.Body.Close()
}

func authenticateUser(request AuthenticationRequest) (string, error) {
	url := os.Getenv("MY_BURGER_APP_URL")
	jsonData, err := json.Marshal(request)

	if err != nil {
		return "", fmt.Errorf("failed to serialize request: %w", err)
	}

	response, err := http.Post(fmt.Sprintf("%s/auth", url), "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to authenticate user: %w", err)
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("authentication failed,  status code: %d", response.StatusCode)
	}

	var authResponse AuthenticationResponse
	err = json.NewDecoder(response.Body).Decode(&authResponse)
	if err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	return authResponse.Token, nil
}
