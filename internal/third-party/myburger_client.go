package third_party

import (
	"bytes"
	"encoding/json"
	"fmt"
	"my-burger-auth/internal/model"
	"net/http"
)

func CreateUser(request model.UserCreationRequest) error {
	url := ""
	jsonData, err := json.Marshal(request)

	if err != nil {
		return fmt.Errorf("failed to serialize request: %w", err)
	}

	response, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	defer response.Body.Close()

	return nil
}

func AuthenticateUser(request model.AuthenticationRequest) (string, error) {
	url := ""
	jsonData, err := json.Marshal(request)

	if err != nil {
		return "", fmt.Errorf("failed to serialize request: %w", err)
	}

	response, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to authenticate user: %w", err)
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("authentication failed,  status code: %d", response.StatusCode)
	}

	var authResponse model.AuthenticationResponse
	err = json.NewDecoder(response.Body).Decode(&authResponse)
	if err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	return authResponse.Token, nil
}
