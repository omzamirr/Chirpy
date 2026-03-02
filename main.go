package main

import (
	"fmt"
	"net/http"
	"log"
)


func main() {
	mux := http.NewServeMux()

	server := &http.Server{
		Addr:  ":8080",
		Handler: mux,
	}

	mux.Handle("/app/", http.StripPrefix("/app/", http.FileServer(http.Dir("."))))
	mux.HandleFunc("/healthz", handlerReadiness)

	fmt.Println("Server is starting on http://localhost:8080")

	err := server.ListenAndServe()
	if err != nil {
    	log.Fatal(err)
	}

}

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}




