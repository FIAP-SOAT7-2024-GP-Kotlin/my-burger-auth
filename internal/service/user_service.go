package service

import (
	"fmt"
	"my-burger-auth/internal/model"
	myburgerClient "my-burger-auth/internal/third-party"
)

func CreateUserAndAuthenticate(request model.UserCreationRequest) (string, error) {
	err := myburgerClient.CreateUser(request)
	if err != nil {
		return "", fmt.Errorf("failed to create user: %w", err)
	}

	authenticationRequest := model.AuthenticationRequest{
		Cpf:      request.Cpf,
		Password: request.Password,
	}

	authenticationResponse, err := myburgerClient.AuthenticateUser(authenticationRequest)
	if err != nil {
		return "", fmt.Errorf("failed to authenticate user: %w", err)
	}
	return authenticationResponse, nil
}
