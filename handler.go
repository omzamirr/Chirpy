package main

import (
    "net/http"
    "encoding/json"
    "context"
    "time"
    "github.com/omzamirr/HttpServer/internal/auth"
)

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
    token, err := auth.GetBearerToken(r.Header)
    if err != nil {
        w.WriteHeader(401)
        return
    }

    userFromToken, err := cfg.db.GetUserFromRefreshToken(context.Background(), token)
    if err != nil {
        w.WriteHeader(401)
        return
    }

    makeJWT, err := auth.MakeJWT(userFromToken.ID, cfg.jwtSecret, time.Hour)
    if err != nil {
        w.WriteHeader(500)
        return
    }

    type response struct {
		Token string `json:"token"`
	}

    w.WriteHeader(200)
	json.NewEncoder(w).Encode(response{Token: makeJWT})


}

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
    token, err := auth.GetBearerToken(r.Header)
    if err != nil {
        w.WriteHeader(401)
        return
    }

    err = cfg.db.RevokeRefreshToken(context.Background(), token)
    if err != nil {
        w.WriteHeader(500)
        return
        }
    w.WriteHeader(204)
    
}