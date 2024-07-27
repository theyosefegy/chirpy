package myutils

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)


type ResponseBody struct {
	Body string `json:"body"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type ValidResponse struct {
	Valid bool `json:"valid"`
}

type CleanedBody struct {
	Cleaned_Body string `json:"cleaned_body"`
}

var profaneWords = []string{"kerfuffle", "sharbert", "fornax"}

// Helper function to respond with an error
func RespondWithError(w http.ResponseWriter ,code int, errorMessage string){
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json")

	errorResp := ErrorResponse{
		Error: errorMessage,
	}

	errorJsonData, _ := json.Marshal(errorResp)
	w.Write(errorJsonData)
}

// Helper function to respond with JSON
func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(code)
    dat, err := json.Marshal(payload)
    if err != nil {
        log.Printf("Error marshalling JSON: %s", err)
        w.WriteHeader(http.StatusInternalServerError)
        w.Write([]byte(`{"error":"Something went wrong"}`))
        return
    }
    w.Write(dat)
}

// Function to replace profane words with "****"
func ReplaceProfaneWords(text string) string {
    words := strings.Split(text, " ")
    for i, word := range words {
        for _, profaneWord := range profaneWords {
            if strings.ToLower(word) == profaneWord {
                words[i] = strings.Repeat("*", len(profaneWord))
            }
        }
    }
    return strings.Join(words, " ")
}
