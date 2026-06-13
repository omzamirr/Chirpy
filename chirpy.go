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
	"github.com/omzamirr/HttpServer/internal/auth"
	"sort"
)


type RequestBody struct {
    Body   string    `json:"body"`
    
}

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}


func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("Error getting JWT: %s", err)
		w.WriteHeader(401)
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		log.Printf("Error validating JWT: %s", err)
		w.WriteHeader(401)
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := RequestBody{}
	err = decoder.Decode(&params)
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
    	UserID: userID,
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

	authorIDString := r.URL.Query().Get("author_id")
	// 1. Declare a variable to hold the slice of database chirps
	var dbChirps []database.Chirp
	var err error

// 2. Check if the query parameter was provided
	if authorIDString != "" {
    // Try to parse the string into a UUID
    	authorID, parseErr := uuid.Parse(authorIDString)
    	if parseErr != nil {
        	w.WriteHeader(http.StatusBadRequest)
        	json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid author ID"})
        	return
    }
    
    // Fetch chirps only for this specific author
    	dbChirps, err = cfg.db.GetChirpsByAuthor(r.Context(), authorID)
		} else {
    	// Fetch all chirps since no author was specified
    	dbChirps, err = cfg.db.GetAllChirps(r.Context())
	}

	if err != nil {
    	log.Printf("Error getting chirps: %s", err)
    	w.WriteHeader(http.StatusInternalServerError)
    	json.NewEncoder(w).Encode(ErrorResponse{Error: "Could not get chirps"})
    	return
}

	responseChirps := []Chirp{}

	for _, dbChirp := range dbChirps {
		responseChirps = append(responseChirps, Chirp{
			ID:        dbChirp.ID,
			CreatedAt: dbChirp.CreatedAt,
			UpdatedAt: dbChirp.UpdatedAt,
			Body:      dbChirp.Body,
			UserID:    dbChirp.UserID,
		})
	}

	// Default to our fallback value
	direction := "asc"

	// Extract the query parameter
	param := r.URL.Query().Get("sort")
		if param == "desc" {
    	direction = "desc"
	}

	sort.Slice(responseChirps, func(i, j int) bool {
    	if direction == "desc" {
        	return responseChirps[i].CreatedAt.After(responseChirps[j].CreatedAt)
    	}
    	return responseChirps[i].CreatedAt.Before(responseChirps[j].CreatedAt)
	})
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


func (cfg *apiConfig) handlerDelete(w http.ResponseWriter, r *http.Request) {

	idString := r.PathValue("chirpID")

	parsedID, err := uuid.Parse(idString)
	if err != nil {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Could not parse"})
		return
		}
	
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		w.WriteHeader(401)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Could not get bearer token"})
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		log.Printf("Error validating JWT: %s", err)
		w.WriteHeader(401)
		return
	}

	chirp, err := cfg.db.GetOneChirp(r.Context(), parsedID)
	if err != nil {
    	w.WriteHeader(404)
    	return
	}

	if chirp.UserID != userID {
    	w.WriteHeader(http.StatusForbidden) // 403
    	return
		}

	err = cfg.db.DeleteChirp(r.Context(), parsedID)
	if err != nil {
    	w.WriteHeader(http.StatusInternalServerError) // 500
    	return
	}

	w.WriteHeader(http.StatusNoContent) // 204
	
}



