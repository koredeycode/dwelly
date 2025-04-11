package handlers

import (
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/koredeycode/dwelly/internal/database"
)

type authHandler func(http.ResponseWriter, *http.Request, database.User)

func (cfg *APIConfig) MiddlewareAuth(handler authHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			respondWithError(w, http.StatusUnauthorized, "Missing or invalid Authorization header")
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})
		if err != nil || !token.Valid {
			respondWithError(w, http.StatusUnauthorized, "Invalid or expired token")
			return
		}

		userID, ok := claims["user_id"].(string)
		if !ok {
			respondWithError(w, http.StatusUnauthorized, "Invalid token claims")
			return
		}

		userUUID, err := uuid.Parse(userID)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid user ID in token")
			return
		}

		user, err := cfg.DB.GetUserByID(r.Context(), userUUID)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "User not found")
			return
		}

		handler(w, r, user)
	}
}
