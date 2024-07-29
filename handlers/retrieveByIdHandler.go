package handlers

import (
	"net/http"
	"strconv"
	"strings"

	myutils "github.com/theyosefegy/chriby/util"
)

func GetChripByIDHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/chirp/")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		myutils.RespondWithError(w, http.StatusBadRequest, "Invalid chirp ID")
		return
	}

	for _, chirp := range chirps {
		if chirp.ID == id {
			myutils.RespondWithJSON(w, http.StatusOK, chirp)
			return
		}
	}

	myutils.RespondWithError(w, http.StatusNotFound, "Chirp not found")
}