package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/koredeycode/dwelly/dwelly-api/models"
	"github.com/koredeycode/dwelly/dwelly-api/utils"
	"github.com/koredeycode/dwelly/internal/database"
	"golang.org/x/crypto/bcrypt"
)

func (cfg *APIConfig) HandlerRegisterUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		FirstName string `json:"first_name" validate:"required"`
		LastName  string `json:"last_name" validate:"required"`
		Password  string `json:"password" validate:"required,min=6"`
		Email     string `json:"email" validate:"required,email"`
	}
	decoder := json.NewDecoder(r.Body)

	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload", err.Error())
		return
	}

	// if err := cfg.Validate.Struct(params); err != nil {
	// 	respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Validation error: %v", err))
	// 	return
	// }
	if err := cfg.Validate.Struct(params); err != nil {
		errorMessages := utils.ExtractValidationErrors(err)
		respondWithError(w, http.StatusBadRequest, "Validation failed", errorMessages...)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Password hashing failure", err.Error())
		return
	}

	user, err := cfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:           uuid.New(),
		FirstName:    params.FirstName,
		LastName:     params.LastName,
		Email:        params.Email,
		PasswordHash: string(hash),
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
	})

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "User creation failed", err.Error())
		return
	}

	respondWithSuccess(w, http.StatusCreated, "User created successfully", models.DatabaseUserToUser(user))
}

// Handle user logging in
func (cfg *APIConfig) HandlerLoginUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password" validate:"required"`
		Email    string `json:"email" validate:"required"`
	}
	decoder := json.NewDecoder(r.Body)

	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload", err.Error())
		return
	}

	if err := cfg.Validate.Struct(params); err != nil {
		errorMessages := utils.ExtractValidationErrors(err)
		respondWithError(w, http.StatusBadRequest, "Validation failed", errorMessages...)
		return
	}

	user, err := cfg.DB.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusForbidden, "User email not found", err.Error())
		return
	}

	_, err = bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Password hashing failure", err.Error())
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(params.Password))

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Inavalid pawword", err.Error())
		return
	}

	redisKey := fmt.Sprintf("dwelly-user-token:%s", user.ID.String())

	existingToken, err := cfg.Redis.Get(r.Context(), redisKey).Result()

	if err == nil && existingToken != "" {
		respondWithJSON(w, http.StatusOK, map[string]interface{}{
			"token": existingToken,
		})
		return
	}

	//handle authentication
	token, err := utils.GenerateJWT(user.ID.String())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Token generation failure", err.Error())
		return
	}

	err = cfg.Redis.Set(r.Context(), redisKey, token, 72*time.Hour).Err()

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to save token", err.Error())
		return
	}

	// "user":  models.DatabaseUsertoUser(user),
	respondWithSuccess(w, http.StatusOK, "User logged in successfully", map[string]interface{}{
		"token": token,
	})

}

func (cfg *APIConfig) HandlerGetCurrentUser(w http.ResponseWriter, r *http.Request, user database.User) {
	respondWithSuccess(w, http.StatusOK, "Current user retrieved successfully", models.DatabaseUserToUser(user))
}

func (cfg *APIConfig) HandlerLogoutUser(w http.ResponseWriter, r *http.Request, user database.User) {
	tokenString := r.Context().Value(TokenContextKey).(string)

	expiration, err := utils.GetTokenExpiry(tokenString)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to get token expiration", err.Error())
		return
	}

	err = cfg.Redis.Set(r.Context(), "dwelly_blacklisted_token:"+tokenString, "revoked", expiration).Err()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Blacklisting token failure", err.Error())
		return
	}

	respondWithSuccess(w, http.StatusOK, "Logged out successfully", nil)
}
