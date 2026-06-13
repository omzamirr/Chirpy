# Chirpy

Chirpy is a custom-built social media backend server written in Go. It handles user authentication, chirp management, and data persistence.

## Features

- **User Management**: Secure user creation and authentication using JWTs and Refresh Tokens.
- **Chirp API**: Create, retrieve, and delete "chirps" with built-in profanity filtering.
- **Data Persistence**: Uses a JSON-based database for local development.
- **Security**: Password hashing with bcrypt and secure token handling.

## Tech Stack

- **Language:** Go (Golang)
- **Routing:** Standard Library `net/http`
- **Authentication:** JWT (JSON Web Tokens) & Bcrypt for password hashing
- **Database:** Local JSON-based storage

## Installation

Ensure you have [Go](https://go.dev/dl/) installed on your machine.

1. Clone the repository:
```bash
   git clone https://github.com/yourusername/chirpy.git
   cd chirpy
```

2. Install dependencies:
```bash
   go mod download
```

## Usage

To start the server:

```bash
go build -o out && ./out
```

The server will start on `localhost:8080`. You can interact with it using tools like `curl` or Postman.

## API Documentation

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/users` | Create a new user |
| POST | `/api/login` | Authenticate and receive tokens |
| GET | `/api/chirps` | Retrieve all chirps |
| DELETE | `/api/chirps/{id}` | Delete a specific chirp (Authenticated) |

## What I Learned

- How to implement a secure authentication flow using Refresh Tokens to minimize the lifespan of Access Tokens.
- The importance of middleware for logging, metrics, and authentication.
- How to structure a Go project for maintainability and clear separation of concerns between data storage and API logic.
