package handler

import (
	"encoding/json"
	"my-burger-auth/internal/model"
	"my-burger-auth/internal/service"
	"net/http"
)

func HandleUserCreation(w http.ResponseWriter, r *http.Request) {
	var request model.UserCreationRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "failed to decode request", http.StatusBadRequest)
		return
	}

	token, err := service.CreateUserAndAuthenticate(request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(model.AuthenticationResponse{Token: token})
}
