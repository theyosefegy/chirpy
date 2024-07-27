package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	myutils "github.com/theyosefegy/chriby/util"
)

func ValidateChirpHandler(w http.ResponseWriter, r *http.Request) {
	respBody := myutils.ResponseBody{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&respBody)


	if err != nil {
		log.Printf("Error decoding JSON: %s", err)
		myutils.RespondWithError(w,  http.StatusInternalServerError, "Something went wrong")
        return
	}

	// Check if the body message's length is less than 140.
	if len(respBody.Body) > 140 {
		myutils.RespondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

    cleanedbody := myutils.CleanedBody{
        Cleaned_Body: myutils.ReplaceProfaneWords(respBody.Body),
    }

    myutils.RespondWithJSON(w, http.StatusOK, cleanedbody)
	// END
}

