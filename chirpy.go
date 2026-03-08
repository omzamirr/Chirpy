package main 


import (
	"encoding/json"
	"time"
	"log" 
	"net/http"
	"github.com/google/uuid"
	"github.com/omzamirr/HttpServer/internal/database"
)


type RequestBody struct {
    Body   string    `json:"body"`
    UserID uuid.UUID `json:"user_id"`
}

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}


func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	params := RequestBody{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(400)
		return
	}

	if len(params.Body) > 140 {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Chirp is too long"}) 
		return
	} 

	cleanedBody := handlerFilterProfanity(params.Body)

	chirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
    	Body:   cleanedBody,
    	UserID: params.UserID,
	})

	if err != nil {
		log.Printf("Error creating chirp: %s", err)
    	json.NewEncoder(w).Encode(ErrorResponse{Error: "Could not create chirp"})
    	return
	}

	w.WriteHeader(201)
	json.NewEncoder(w).Encode(Chirp{
    	ID:        chirp.ID,
    	CreatedAt: chirp.CreatedAt,
    	UpdatedAt: chirp.UpdatedAt,
    	Body:      chirp.Body,
    	UserID:    chirp.UserID,
	})
	
}




