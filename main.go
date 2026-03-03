package main

import (
	"fmt"
	"net/http"
	"log"
	"sync/atomic"
)


type apiConfig struct {
	fileserverHits atomic.Int32
}


func main() {
	mux := http.NewServeMux()

	server := &http.Server{
		Addr:  ":8080",
		Handler: mux,
	}

	apiCfg := &apiConfig{}

	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app/", http.FileServer(http.Dir(".")))))
	mux.HandleFunc("/healthz", handlerReadiness)
	mux.HandleFunc("/metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("/reset", apiCfg.handlerReset)

	fmt.Println("Server is starting on http://localhost:8080")

	err := server.ListenAndServe()
	if err != nil {
    	log.Fatal(err)
	}

}


func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}


func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(fmt.Sprintf("Hits: %v", cfg.fileserverHits.Load())))
}



