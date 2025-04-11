package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/koredeycode/dwelly/dwelly-api/models"
	"github.com/koredeycode/dwelly/internal/database"
)

func (cfg *APIConfig) HandlerCreateInquiry(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		ListingID string `json:"listing_id"`
		Message   string `json:"message"`
	}

	decoder := json.NewDecoder(r.Body)

	params := parameters{}

	err := decoder.Decode(&params)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "error parsing json")
		return
	}

	inquiry, err := cfg.DB.CreateInquiry(r.Context(), database.CreateInquiryParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		SenderID:  user.ID,
		ListingID: func() uuid.UUID {
			id, err := uuid.Parse(params.ListingID)
			if err != nil {
				respondWithError(w, http.StatusBadRequest, "invalid listing_id format")
				return uuid.Nil
			}
			return id
		}(),
		Message: params.Message,
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error creating inquiry")
		return
	}

	respondWithJSON(w, http.StatusCreated, models.DatabaseInquirytoInquiry(inquiry))
}

func (cfg *APIConfig) HandlerGetInquiry(w http.ResponseWriter, r *http.Request, user database.User) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Single inquiry"))
}

func (cfg *APIConfig) HandlerGetInquiries(w http.ResponseWriter, r *http.Request, user database.User) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("List of inquiries"))
}

func (cfg *APIConfig) HandlerUpdateInquiryStatus(w http.ResponseWriter, r *http.Request, user database.User) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Inquiry status updated"))
}

func (cfg *APIConfig) HandlerDeleteInquiry(w http.ResponseWriter, r *http.Request, user database.User) {
	w.WriteHeader(http.StatusNoContent)
}
