package main


import (
	"net/http"
	"encoding/json"
	"log"
	"strings"
	
)


type ErrorResponse struct {
    Error string `json:"error"`
}


type ValidResponse struct {
    Valid bool `json:"valid"`
}


func handlerValidateChirp(w http.ResponseWriter, r *http.Request) {

	type parameters struct {
		Body  string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(500)
		return
	}
	
	if len(params.Body) > 140 {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Chirp is too long"}) 
	} else {
		cleaned := handlerFilterProfanity(params.Body)
		json.NewEncoder(w).Encode(CleanedBody{Clean: cleaned})
	}
}


type CleanedBody struct {
	Clean string `json:"cleaned_body"`
}


var badWords = []string{"kerfuffle", "sharbert", "fornax"}


func handlerFilterProfanity(body string) string {
	words := strings.Split(body, " ")
	for i, word := range words {
		for _, badWord := range badWords {
			if strings.ToLower(word) == badWord {
				words[i] = strings.Repeat("*", 4)
			}
		}
	}

	return strings.Join(words, " ")


}