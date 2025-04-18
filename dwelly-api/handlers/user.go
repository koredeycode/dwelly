package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/koredeycode/dwelly/dwelly-api/models"
	"github.com/koredeycode/dwelly/dwelly-api/utils"
	"github.com/koredeycode/dwelly/internal/database"
	"golang.org/x/crypto/bcrypt"
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

	userIDStr := chi.URLParam(r, "userId")

	userID, errMsg := utils.GetUUIDParam(userIDStr, "listing")

	if errMsg != "" {
		respondWithError(w, http.StatusBadRequest, errMsg)
		return
	}

	if user.ID != userID {
		respondWithError(w, http.StatusForbidden, "user not authorized to update this profile")
		return
	}
	type parameters struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Password  string `json:"password"`
		Email     string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}

	err := decoder.Decode(&params)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Could not hash password: %v", err))
		return
	}

	newUser, err := cfg.DB.UpdateUser(r.Context(), database.UpdateUserParams{
		ID:           user.ID,
		FirstName:    valueOrDefault(user.FirstName, params.FirstName),
		LastName:     valueOrDefault(user.LastName, params.LastName),
		Email:        valueOrDefault(user.Email, params.Email),
		PasswordHash: valueOrDefault(user.PasswordHash, string(hash)),
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to update user: %v", err))
		return
	}
	// Respond with the updated user
	respondWithJSON(w, http.StatusOK, models.DatabaseUserToUser(newUser))
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
