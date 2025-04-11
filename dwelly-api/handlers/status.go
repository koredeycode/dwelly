package handlers

import (
	"net/http"
)

func (cfg *APIConfig) HandlerAPIStatus(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, http.StatusOK, map[string]string{
		"status": "ok",
	})
}
