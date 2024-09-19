package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"os"
	"time"
)

type AuthenticationRequest struct {
	Cpf      string `json:"cpf"`
	Password string `json:"password"`
}

type AuthenticationResponse struct {
	Token string `json:"access_token"`
}

type Claims struct {
	CPF string `json:"cpf"`
	jwt.StandardClaims
}

type User struct {
	ID       uuid.UUID
	cpf      string
	password string
	role     string
}

func Main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/authenticate", handleAuthentication)

	log.Println("Starting server on :8090")
	if err := http.ListenAndServe(":8090", mux); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}

func handleAuthentication(w http.ResponseWriter, r *http.Request) {
	var request AuthenticationRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "failed to decode request", http.StatusBadRequest)
		return
	}

	db, err := setupDbConnection()
	if err != nil {
		http.Error(w, "Error connecting to database", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	user, err := findUserByCPF(db, request.Cpf)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error on finding user by cpf", http.StatusInternalServerError)
		return
	}

	if err := verifyPassword(user.password, request.Password); err != nil {
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	token, err := generateJWT(request.Cpf)
	if err != nil {
		http.Error(w, "Error generating JWT token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(AuthenticationResponse{Token: token})
}

func setupDbConnection() (*sql.DB, error) {
	databaseUrl := os.Getenv("DATABASE_URL")
	databaseUser := os.Getenv("DATABASE_USER")
	databasePassword := os.Getenv("DATABASE_PASSWORD")
	databaseName := os.Getenv("DATABASE_NAME")
	connectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		databaseUrl, "5432", databaseUser, databasePassword, databaseName)

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return db, nil
}

func findUserByCPF(db *sql.DB, cpf string) (*User, error) {
	var user User
	const findUserByCpfQuery = "SELECT id, cpf, password, role FROM users WHERE cpf = $1"
	err := db.QueryRow(findUserByCpfQuery, cpf).Scan(&user.ID, &user.cpf, &user.password, &user.role)
	if err != nil {
		return nil, fmt.Errorf("failed to query user by cpf: %w", err)
	}
	return &user, nil
}

func verifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func generateJWT(cpf string) (string, error) {
	jwtKey := os.Getenv("JWT_KEY")
	if jwtKey == "" {
		return "", fmt.Errorf("failed to get jwt key")
	}

	expirationTime := time.Now().Add(5 * time.Minute)
	claims := &Claims{
		CPF: cpf,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtKey))
}
