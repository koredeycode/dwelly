package handlers

import (
	"fmt"
	"net/http"

	"github.com/koredeycode/dwelly/dwelly-api/models"
	"github.com/koredeycode/dwelly/internal/database"
)

func (cfg *APIConfig) HandlerGetUser(w http.ResponseWriter, r *http.Request, user database.User) {
	respondWithJSON(w, http.StatusOK, models.DatabaseUsertoUser(user))
}

func (cfg *APIConfig) HandlerGetUserListings(w http.ResponseWriter, r *http.Request, user database.User) {
	user_listings, err := cfg.DB.ListUserListings(r.Context(), user.ID)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error getting listings: %v", err))
		return
	}

	// Respond with the listings
	respondWithJSON(w, http.StatusOK, models.DatabaseListingstoListings(user_listings))
}
