package model

type UserCreationRequest struct {
	Cpf      string   `json:"cpf"`
	Password string   `json:"password"`
	Role     UserRole `json:"role"`
}

type UserRole string

const (
	UserRoleAdmin UserRole = "ADMIN"
	UserRoleUser  UserRole = "USER"
)
