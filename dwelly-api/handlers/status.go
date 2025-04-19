package handlers

import (
	"net/http"
)

func (cfg *APIConfig) HandlerAPIStatus(w http.ResponseWriter, r *http.Request) {
	respondWithSuccess(w, http.StatusOK, "API running", map[string]string{
		"status": "ok",
	})

}
