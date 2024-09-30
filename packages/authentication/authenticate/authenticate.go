package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

const (
	UserRoleAdmin = "ADMIN"
	UserRoleUser  = "USER"
)

type Request struct {
	Cpf      string      `json:"cpf"`
	Password string      `json:"password"`
	Type     RequestType `json:"type"`
}

type RequestType string

const (
	USER_CREATION       RequestType = "USER_CREATION"
	USER_AUTHENTICATION RequestType = "USER_AUTHENTICATION"
)

type Response struct {
	StatusCode int               `json:"statusCode,omitempty"`
	Headers    map[string]string `json:"headers,omitempty"`
	Token      string            `json:"access_token"`
}

type Claims struct {
	CPF string `json:"cpf"`
	jwt.StandardClaims
}

type User struct {
	ID       uuid.UUID
	Cpf      string
	Password string
	Role     string
}

var (
	config Config
	ErrNoRequest = errors.New("no request type provided")
)

type Config struct {
	DatabaseUrl      string
	DatabaseUsername string
	DatabasePassword string
	DatabaseName     string
	DatabaseSchema   string
	JwtKey           string
	Port						 string
}

func init() {
	config = Config{
		DatabaseUrl:      getEnv("DATABASE_URL"),
		DatabaseUsername: getEnv("DATABASE_USERNAME"),
		DatabasePassword: getEnv("DATABASE_PASSWORD"),
		DatabaseName:     getEnv("DATABASE_NAME"),
		DatabaseSchema:   getEnv("DATABASE_SCHEMA"),
		JwtKey:           getEnv("JWT_KEY"),
		Port:						  getEnv("DATABASE_PORT"),
	}
}

func getEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("%s environment variable is required", key)
	}
	return value
}

func Main(input Request) (*Response, error) {
	switch input.Type {
	case USER_CREATION:
		err := handleUserCreation(input)
		if err != nil {
			return &Response{StatusCode: http.StatusInternalServerError}, err
		}
		return &Response{StatusCode: http.StatusCreated}, nil
	case USER_AUTHENTICATION:
		response, err := handleAuthentication(input)
		if err != nil {
			return &Response{StatusCode: http.StatusUnauthorized}, err
		}
		return &Response{StatusCode: http.StatusOK, Token: response}, nil
	default:
		return &Response{StatusCode: http.StatusBadRequest}, ErrNoRequest
	}
}

func handleAuthentication(request Request) (string, error) {
	db, err := setupDbConnection()
	if err != nil {
		log.Println("Error connecting to database:", err)
		return "", err
	}
	defer db.Close()

	user, err := findUserByCPF(db, request.Cpf)
	if err != nil {
		log.Println("Error finding user by CPF:", err)
		return "", err
	}

	if err := verifyPassword(user.Password, request.Password); err != nil {
		return "", fmt.Errorf("invalid password: %v", http.StatusUnauthorized)
	}

	token, err := generateJWT(request.Cpf)
	if err != nil {
		log.Println("Error generating JWT token:", err)
		return "", err
	}
	return token, nil
}

func handleUserCreation(request Request) error {
	db, err := setupDbConnection()
	if err != nil {
		log.Println("Error connecting to database:", err)
		return err
	}
	defer db.Close()

	user, err := findUserByCPF(db, request.Cpf)
	if err != nil && err != sql.ErrNoRows {
		log.Println("Error finding user by CPF:", err)
		return err
	}

	if user != nil {
		return fmt.Errorf("user already exists")
	}

	if err := createUser(db, request); err != nil {
		log.Println("Error creating user:", err)
		return err
	}

	return nil
}

func setupDbConnection() (*sql.DB, error) {
	connectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s search_path=%s sslmode=require",
		config.DatabaseUrl, config.Port, config.DatabaseUsername, config.DatabasePassword, config.DatabaseName, config.DatabaseSchema)
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		log.Println("Error opening database connection:", err)
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return db, nil
}

func findUserByCPF(db *sql.DB, cpf string) (*User, error) {
	var user User
	const findUserByCpfQuery = "SELECT id, cpf, password, role FROM \"user\" WHERE cpf = $1"
	err := db.QueryRow(findUserByCpfQuery, cpf).Scan(&user.ID, &user.Cpf, &user.Password, &user.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to query user by cpf: %w", err)
	}
	return &user, nil
}

func createUser(db *sql.DB, request Request) error {
	const createUserQuery = "INSERT INTO \"user\" (id, cpf, password, role) VALUES ($1, $2, $3, $4)"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}
	_, err = db.Exec(createUserQuery, uuid.New(), request.Cpf, hashedPassword, UserRoleUser)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	fmt.Printf("User successfully created with cpf: %s\n", request.Cpf)
	return nil
}

func verifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func generateJWT(cpf string) (string, error) {
	expirationTime := time.Now().Add(5 * time.Minute)
	claims := &Claims{
		CPF: cpf,
		StandardClaims: jwt.StandardClaims{
			Subject:   cpf,
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.JwtKey))
}
