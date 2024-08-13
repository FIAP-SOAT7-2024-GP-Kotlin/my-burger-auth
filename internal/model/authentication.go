package model

type AuthenticationRequest struct {
	Cpf      string `json:"cpf"`
	Password string `json:"password"`
}

type AuthenticationResponse struct {
	Token string `json:"token"`
}
