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

	//mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	//})

	fmt.Println("Server is startin on http://localhost:8080")

	err := server.ListenAndServe()
	if err != nil {
    	log.Fatal(err)
	}

}




