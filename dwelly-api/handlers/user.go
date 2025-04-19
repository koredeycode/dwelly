package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/koredeycode/dwelly/dwelly-api/models"
	"github.com/koredeycode/dwelly/dwelly-api/utils"
	"github.com/koredeycode/dwelly/internal/database"
	"golang.org/x/crypto/bcrypt"
)

func (cfg *APIConfig) HandlerGetUser(w http.ResponseWriter, r *http.Request, user database.User) {
	userIDStr := chi.URLParam(r, "userId")

	userID, errMsg := utils.GetUUIDParam(userIDStr, "user")

	if errMsg != "" {
		respondWithError(w, http.StatusBadRequest, errMsg)
		return
	}

	if user.ID == userID {
		respondWithSuccess(w, http.StatusOK, "User retrieved successfully", models.DatabaseUserToUser(user))
		return
	}
	queriedUser, err := cfg.DB.GetUserByID(r.Context(), userID)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "User retrieval failure", err.Error())
		return
	}

	respondWithSuccess(w, http.StatusOK, "User retrieved successfully", models.DatabaseUserToUser(queriedUser))
}

func (cfg *APIConfig) HandlerGetUserListings(w http.ResponseWriter, r *http.Request, user database.User) {
	user_listings, err := cfg.DB.ListUserListings(r.Context(), user.ID)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "User's listing retrieval error", err.Error())
		return
	}

	// Respond with the listings
	respondWithSuccess(w, http.StatusOK, "User's listing retrieved successfully", models.DatabaseListingsToListings(user_listings))

}

func (cfg *APIConfig) HandlerUpdateUser(w http.ResponseWriter, r *http.Request, user database.User) {

	userIDStr := chi.URLParam(r, "userId")

	userID, errMsg := utils.GetUUIDParam(userIDStr, "listing")

	if errMsg != "" {
		respondWithError(w, http.StatusBadRequest, errMsg)
		return
	}

	if user.ID != userID {
		respondWithError(w, http.StatusForbidden, "User not authorized")
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
		respondWithError(w, http.StatusBadRequest, "Invalid request payload", err.Error())
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Password hashing error", err.Error())
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
		respondWithError(w, http.StatusInternalServerError, "User update failure", err.Error())
		return
	}
	// Respond with the updated user
	respondWithSuccess(w, http.StatusOK, "User updated successfully", models.DatabaseUserToUser(newUser))
}

func (cfg *APIConfig) HandlerDeleteUser(w http.ResponseWriter, r *http.Request, user database.User) {
	userIDStr := chi.URLParam(r, "userId")

	userID, errMsg := utils.GetUUIDParam(userIDStr, "listing")

	if errMsg != "" {
		respondWithError(w, http.StatusBadRequest, errMsg)
		return
	}

	if user.ID != userID {
		respondWithError(w, http.StatusForbidden, "User not authorized")
		return
	}

	err := cfg.DB.DeleteUser(r.Context(), user.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "User deletion failure", err.Error())
		return
	}

	// Respond with a success message
	respondWithSuccess(w, http.StatusOK, "User deleted successfully", nil)

}
