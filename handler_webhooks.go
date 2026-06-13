package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"log"
	"net/http"
	"github.com/omzamirr/HttpServer/internal/auth"
)

func (cfg *apiConfig) handlerPolkaWebhooks(w http.ResponseWriter, r *http.Request) {

    apiKey, err := auth.GetAPIKey(r.Header)
    if err != nil {
        // If there's an error getting the key (e.g., missing header), return 401
        w.WriteHeader(http.StatusUnauthorized)
        return
    }

    // 2. Validate the key matches our configured polkaKey
    if apiKey != cfg.polkaKey {
        // If the key doesn't match, return 401
        w.WriteHeader(http.StatusUnauthorized)
        return
    }


	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserID uuid.UUID `json:"user_id"`
		} `json:"data"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(500)
		return
	}

	if params.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	_, err = cfg.db.UpgradeUserToChirpyRed(r.Context(), params.Data.UserID)

	if errors.Is(err, sql.ErrNoRows) {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if err != nil {
		log.Printf("Error upgrading user: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)

}
