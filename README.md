# My Burger Auth

My Burger Auth is a Go-based authentication service for the My Burger application. It provides user authentication using JWT tokens and connects to a PostgreSQL database.

## Table of Contents

- [Features](#features)
- [Requirements](#requirements)
- [Installation](#installation)
- [Configuration](#configuration)
- [Usage](#usage)

## Features

- User authentication with JWT tokens
- Secure password hashing with bcrypt
- PostgreSQL database integration

## Requirements

- Go 1.20 or later
- PostgreSQL 15 or later
- Docker (for deployment)

## Installation

1. Clone the repository:
```sh
git clone https://github.com/FIAP-SOAT7-2024-GP-Kotlin/my-burger-auth.git
````
### Install dependencies:

```shell
go mod tidy
 ```
### Configuration
 
Create a .env file in the root directory with the following content:

>DATABASE_URL=postgresql://localhost:5432/my_burger</br>
DATABASE_USER=my_burger</br>
DATABASE_PASSWORD=password</br>
JWT_KEY=your_jwt_secret_key
 
### Usage
Build the application:  
```shell 
go build -o main .
```
### Run the application:  
./main
The server will start on port 8090. You can access the authentication endpoint at http://localhost:8090/authenticate.  


#### To deploy the application to DigitalOcean, use the provided `deploy.template.yaml` file.
