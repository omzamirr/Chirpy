package main 


import (
	"encoding/json"
	"time"
	"log" 
	"net/http"
	"github.com/google/uuid"
	"github.com/omzamirr/HttpServer/internal/database"
	"database/sql"
    "errors"
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


func (cfg *apiConfig) handlerGetAllChirps(w http.ResponseWriter, r *http.Request) {

	allChirps, err := cfg.db.GetAllChirps(r.Context())

	if err != nil {
		log.Printf("Error getting all chirps: %s", err)
    	json.NewEncoder(w).Encode(ErrorResponse{Error: "Could not not get all chirps"})
    	return
	}

	responseChirps := []Chirp{}

	for _, dbChirp := range allChirps {
		responseChirps = append(responseChirps, Chirp{
			ID:        dbChirp.ID,
			CreatedAt: dbChirp.CreatedAt,
			UpdatedAt: dbChirp.UpdatedAt,
			Body:      dbChirp.Body,
			UserID:    dbChirp.UserID,
		})
	}

	w.WriteHeader(200)
	json.NewEncoder(w).Encode(responseChirps)


}


func (cfg *apiConfig) handlerGetOneChirp(w http.ResponseWriter, r *http.Request) {

	chirpIDString := r.PathValue("chirpID")
	id, err := uuid.Parse(chirpIDString)
	if err != nil {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Could not retrieve the value of the 'chirpID' path parameter"})
		return
	}

	chirp, err := cfg.db.GetOneChirp(r.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Could not get the chirp"})
		return
	}

	w.WriteHeader(200)
	json.NewEncoder(w).Encode(Chirp{
    	ID:        chirp.ID,
    	CreatedAt: chirp.CreatedAt,
    	UpdatedAt: chirp.UpdatedAt,
    	Body:      chirp.Body,
    	UserID:    chirp.UserID,
	})
}


