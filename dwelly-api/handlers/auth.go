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
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Invalid request payload: %v", err))
		return
	}

	if err := cfg.Validate.Struct(params); err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Validation error: %v", err))
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Could not hash password: %v", err))
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
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Could not create user: %v", err))
		return
	}
	respondWithJSON(w, http.StatusCreated, models.DatabaseUserToUser(user))
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
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Invalid request payload: %v", err))
		return
	}

	if err := cfg.Validate.Struct(params); err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Validation error: %v", err))
		return
	}

	user, err := cfg.DB.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusForbidden, fmt.Sprintf("User email doesn't exist: %v", err))
		return
	}

	_, err = bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Could not hash password: %v", err))
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(params.Password))

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, fmt.Sprintf("Invalid password: %v", err))
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
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to generate token: %v", err))
		return
	}

	err = cfg.Redis.Set(r.Context(), redisKey, token, 72*time.Hour).Err()

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to store token in Redis: %v", err))
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		// "user":  models.DatabaseUsertoUser(user),
		"token": token,
	})

	// respondWithJSON(w, http.StatusOK, models.DatabaseUsertoUser(user))

}

func (cfg *APIConfig) HandlerLogoutUser(w http.ResponseWriter, r *http.Request, user database.User) {
	tokenString := r.Context().Value("token").(string)

	expiration, err := utils.GetTokenExpiry(tokenString)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to get token expiration: %v", err))
		return
	}

	err = cfg.Redis.Set(r.Context(), "dwelly_blacklisted_token:"+tokenString, "revoked", expiration).Err()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to blacklist token: %v", err))
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Logged out successfully"})
}
