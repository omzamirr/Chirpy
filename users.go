package main 


import (
    "encoding/json"
    "log"
    "net/http"
    "time"

    "github.com/google/uuid"
)


type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}


func (cfg *apiConfig) handlerRegister(w http.ResponseWriter, r *http.Request) {

	type parameters struct {
		Email  string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(500)
		return
	}
	
	user, err := cfg.db.CreateUser(r.Context(), params.Email)

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