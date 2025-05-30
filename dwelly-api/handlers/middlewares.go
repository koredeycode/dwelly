package handlers

import (
	"context"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/koredeycode/dwelly/dwelly-api/utils"
	"github.com/koredeycode/dwelly/internal/database"
)

type contextKey string

const TokenContextKey contextKey = "token"

type authHandler func(http.ResponseWriter, *http.Request, database.User)

// ChainAuth will chain the provided authorization middlewares
func (cfg *APIConfig) Auth(finalHandler authHandler, middlewares ...func(authHandler) authHandler) http.HandlerFunc {
	// Apply all middlewares to the final handler in order
	for i := len(middlewares) - 1; i >= 0; i-- {
		finalHandler = middlewares[i](finalHandler)
	}
	// Now wrap the entire chain with the authentication middleware
	return cfg.AuthenticationMiddleware(finalHandler)
}

func (cfg *APIConfig) AuthenticationMiddleware(handler authHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			respondWithError(w, http.StatusUnauthorized, "Missing authorization header")
			return
		}
		if !strings.HasPrefix(authHeader, "Bearer ") {
			respondWithError(w, http.StatusUnauthorized, "Invalid authorization header")
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		isBlacklisted, err := cfg.Redis.Exists(r.Context(), "dwelly_blacklisted_token:"+tokenString).Result()
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Token retrival errror", err.Error())
			return
		}
		if isBlacklisted == 1 {
			respondWithError(w, http.StatusUnauthorized, "Token has been revoked")
			return
		}

		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil || !token.Valid {
			respondWithError(w, http.StatusUnauthorized, "Invalid token", err.Error())
			return
		}

		expiration, err := utils.GetTokenExpiry(tokenString)

		if expiration <= 0 {
			respondWithError(w, http.StatusUnauthorized, "Token has expired")
			return
		}

		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Token expiration failure", err.Error())
			return
		}

		userID, ok := claims["user_id"].(string)
		if !ok {
			respondWithError(w, http.StatusUnauthorized, "Invalid token, user ID not found in the token")
			return
		}

		userUUID, err := uuid.Parse(userID)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid user ID in token", err.Error())
			return
		}

		user, err := cfg.DB.GetUserByID(r.Context(), userUUID)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "User not found", err.Error())
			return
		}

		ctx := context.WithValue(r.Context(), TokenContextKey, tokenString)

		handler(w, r.WithContext(ctx), user)
	}
}

func (cfg *APIConfig) ListingOwnerAuthorization(handler authHandler) authHandler {
	return func(w http.ResponseWriter, r *http.Request, user database.User) {
		listingIDStr := chi.URLParam(r, "listingId")

		listingID, errMsg := utils.GetUUIDParam(listingIDStr, "listing")

		if errMsg != "" {
			respondWithError(w, http.StatusBadRequest, errMsg)
			return
		}
		listing, err := cfg.DB.GetListingByID(r.Context(), listingID)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Error checking listing owner", err.Error())
			return
		}
		if listing.UserID != user.ID {
			respondWithError(w, http.StatusForbidden, "User is not the owner of the listing")
			return
		}

		handler(w, r, user)
	}
}

func (cfg *APIConfig) InquirySenderAuthorization(handler authHandler) authHandler {
	return func(w http.ResponseWriter, r *http.Request, user database.User) {
		inquiryIDStr := chi.URLParam(r, "inquiryId")

		inquiryID, errMsg := utils.GetUUIDParam(inquiryIDStr, "inquiry")

		if errMsg != "" {
			respondWithError(w, http.StatusBadRequest, errMsg)
			return
		}
		inquiry, err := cfg.DB.GetInquiryById(r.Context(), inquiryID)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "error checking inquiry owner", err.Error())
			return
		}
		if inquiry.SenderID != user.ID {
			respondWithError(w, http.StatusForbidden, "User is not the sender of the inquiry")
			return
		}

		handler(w, r, user)
	}
}

func (cfg *APIConfig) MessageSenderAuthorization(handler authHandler) authHandler {
	return func(w http.ResponseWriter, r *http.Request, user database.User) {
		messageIDStr := chi.URLParam(r, "messageId")

		messageID, errMsg := utils.GetUUIDParam(messageIDStr, "message")

		if errMsg != "" {
			respondWithError(w, http.StatusBadRequest, errMsg)
			return
		}
		message, err := cfg.DB.GetMessage(r.Context(), messageID)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Error checking message owner", err.Error())
			return
		}
		if message.SenderID != user.ID {
			respondWithError(w, http.StatusForbidden, "User is not the sender of the message")
			return
		}

		handler(w, r, user)
	}
}

func (cfg *APIConfig) InquirySenderOrListingOwnerAuthorization(handler authHandler) authHandler {
	return func(w http.ResponseWriter, r *http.Request, user database.User) {
		// Check if the user is the inquiry sender
		inquiryIDStr := chi.URLParam(r, "inquiryId")
		inquiryID, errMsg := utils.GetUUIDParam(inquiryIDStr, "inquiry")

		if errMsg != "" {
			respondWithError(w, http.StatusBadRequest, errMsg)
			return
		}

		inquiry, err := cfg.DB.GetInquiryById(r.Context(), inquiryID)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "error checking inquiry owner", err.Error())
			return
		}

		// If user is not the inquiry sender, check if they are the listing owner
		if inquiry.SenderID == user.ID {
			handler(w, r, user)
			return
		}

		listing, err := cfg.DB.GetListingByID(r.Context(), inquiry.ListingID)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Error checking listing owner", err.Error())
			return
		}

		if listing.UserID != user.ID {
			respondWithError(w, http.StatusForbidden, "user is neither the inquiry sender nor the listing owner")
			return
		}

		handler(w, r, user)
	}
}
