package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	myutils "github.com/theyosefegy/chriby/util"
)

var idCounter = 1

type Chirp struct {
	ID   int    `json:"id"`
	Body string `json:"body"`
}

type DBStructure struct {
    Chirps map[int]Chirp `json:"chirps"`
}

var chirps = []Chirp{}

// The post handler method.
func PostChirpHandler(w http.ResponseWriter, r *http.Request) {
	var chirpReq Chirp

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&chirpReq)
	if err != nil {
		log.Printf("Error decoding JSON: %s", err)
		myutils.RespondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	// Validate chirp body length
	if len(chirpReq.Body) > 140 {
		myutils.RespondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	chirpReq.Body = myutils.ReplaceProfaneWords(chirpReq.Body)

	// Lock and update ID counter and chirps slice
	chirpReq.ID = idCounter
	idCounter++
	
	chirps = append(chirps, chirpReq)

	myutils.RespondWithJSON(w, 201, chirpReq)
}

func GetChirpHandler(w http.ResponseWriter, r *http.Request) {

	myutils.RespondWithJSON(w, http.StatusOK, chirps)
}