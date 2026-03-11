package main 


import (
    "encoding/json"
    "log"
    "net/http"
    "time"
	"github.com/omzamirr/HttpServer/internal/database"
    "github.com/google/uuid"
	"github.com/omzamirr/HttpServer/internal/auth"
)


type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}


func (cfg *apiConfig) handlerRegister(w http.ResponseWriter, r *http.Request) {

	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(500)
		return
	}
	
	hashedValue, err := auth.HashPassword(params.Password)
	if err != nil {
		log.Printf("Error hashing password: %s", err)
		w.WriteHeader(500)
		return
	}

	user, err := cfg.db.CreateUser(r.Context(), database.CreateUserParams{
    	Email:          params.Email,
    	HashedPassword: hashedValue,
	})

	if err != nil {
    	log.Printf("Error creating user: %s", err)
    	json.NewEncoder(w).Encode(ErrorResponse{Error: "Could not create user"})
    	return
	}

	responseUser := User{
    	ID:        user.ID,
    	CreatedAt: user.CreatedAt,
    	UpdatedAt: user.UpdatedAt,
    	Email:     user.Email,
	}

	w.WriteHeader(201)
	json.NewEncoder(w).Encode(responseUser)

}


func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {

	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(500)
		return
	}

	user, err := cfg.db.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		log.Printf("Could not get email: %s", err)
		w.WriteHeader(401)
		return
	}

	match, err := auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil || !match {
		w.WriteHeader(401)
		return
	}

	responseUser := User{
    	ID:        user.ID,
    	CreatedAt: user.CreatedAt,
    	UpdatedAt: user.UpdatedAt,
    	Email:     user.Email,
	}
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(responseUser)
	
}