package handlers

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/koredeycode/dwelly/dwelly-api/utils"
	"github.com/koredeycode/dwelly/internal/database"
)

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
			respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error checking listing owner: %v", err))
			return
		}
		if listing.UserID != user.ID {
			respondWithError(w, http.StatusForbidden, "user is not the owner of the listing")
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
			respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error checking inquiry owner: %v", err))
			return
		}
		if inquiry.SenderID != user.ID {
			respondWithError(w, http.StatusForbidden, "user is not the sender of the inquiry")
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
			respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error checking message owner: %v", err))
			return
		}
		if message.SenderID != user.ID {
			respondWithError(w, http.StatusForbidden, "user is not the sender of the message")
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
			respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error checking inquiry owner: %v", err))
			return
		}

		// If user is not the inquiry sender, check if they are the listing owner
		if inquiry.SenderID == user.ID {
			handler(w, r, user)
			return
		}

		listing, err := cfg.DB.GetListingByID(r.Context(), inquiry.ListingID)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error checking listing owner: %v", err))
			return
		}

		if listing.UserID == user.ID {
			handler(w, r, user)
			return
		}

		// If neither condition is met, deny access
		respondWithError(w, http.StatusForbidden, "user is neither the inquiry sender nor the listing owner")
	}
}
