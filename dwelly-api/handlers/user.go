package handlers

import (
	"fmt"
	"net/http"

	"github.com/koredeycode/dwelly/dwelly-api/models"
	"github.com/koredeycode/dwelly/internal/database"
)

func (cfg *APIConfig) HandlerGetUser(w http.ResponseWriter, r *http.Request, user database.User) {
	respondWithJSON(w, http.StatusOK, models.DatabaseUserToUser(user))
}

func (cfg *APIConfig) HandlerGetUserListings(w http.ResponseWriter, r *http.Request, user database.User) {
	user_listings, err := cfg.DB.ListUserListings(r.Context(), user.ID)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error getting listings: %v", err))
		return
	}

	// Respond with the listings
	respondWithJSON(w, http.StatusOK, models.DatabaseListingsToListings(user_listings))
}

func (cfg *APIConfig) HandlerUpdateUser(w http.ResponseWriter, r *http.Request, user database.User) {

	// Respond with the updated user
	respondWithJSON(w, http.StatusOK, models.DatabaseUserToUser(user))
}

func (cfg *APIConfig) HandlerDeleteUser(w http.ResponseWriter, r *http.Request, user database.User) {
	// Delete the user
	err := cfg.DB.DeleteUser(r.Context(), user.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error deleting user: %v", err))
		return
	}

	// Respond with a success message
	respondWithJSON(w, http.StatusOK, map[string]string{"message": "user deleted"})
}
