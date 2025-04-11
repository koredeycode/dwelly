package handlers

import (
	"net/http"

	"github.com/koredeycode/dwelly/internal/database"
)

func (cfg *APIConfig) HandlerAddListingImage(w http.ResponseWriter, r *http.Request, user database.User) {
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Listing image added"))
}

func (cfg *APIConfig) HandlerDeleteListingImage(w http.ResponseWriter, r *http.Request, user database.User) {
	w.WriteHeader(http.StatusNoContent)
}
