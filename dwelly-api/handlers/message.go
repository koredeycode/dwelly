package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/koredeycode/dwelly/dwelly-api/models"
	"github.com/koredeycode/dwelly/dwelly-api/services"
	"github.com/koredeycode/dwelly/internal/database"
)

func (cfg *APIConfig) HandlerCreateInquiryMessage(w http.ResponseWriter, r *http.Request, user database.User) {
	listingIDStr := chi.URLParam(r, "listingId")
	if listingIDStr == "" {
		respondWithError(w, http.StatusBadRequest, "missing listing ID")
		return
	}
	listingId, err := uuid.Parse(listingIDStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid listing ID")

		return
	}

	inquiryIDStr := chi.URLParam(r, "inquiryId")
	if inquiryIDStr == "" {
		respondWithError(w, http.StatusBadRequest, "missing inquiry ID")
		return
	}

	inquiryId, err := uuid.Parse(inquiryIDStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid inquiry ID")
		return
	}

	//Permission check
	isListingOwner, err := services.IsListingOwner(cfg.DB, r, listingId, user.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error checking listing owner: %v", err))
	}

	if !isListingOwner {
		respondWithError(w, http.StatusForbidden, "user is not the owner of the listing")
		return
	}

	type parameters struct {
		Content string `json:"content"`
	}
	decoder := json.NewDecoder(r.Body)

	params := parameters{}

	err = decoder.Decode(&params)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "error parsing json")
		return
	}
	message, err := cfg.DB.CreateMessage(r.Context(), database.CreateMessageParams{
		ID:        uuid.New(),
		InquiryID: inquiryId,
		SenderID:  user.ID,
		Content:   params.Content,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error creating message to inquiry")
		return
	}
	respondWithJSON(w, http.StatusOK, models.DatabaseMessageToMessage(message))

}

// to do: authorization should be handled, inquiry owner should be able to get the messages and listing owner tied to the inquiry
func (cfg *APIConfig) HandlerGetInquiryMessages(w http.ResponseWriter, r *http.Request, user database.User) {
	inquiryIdStr := chi.URLParam(r, "inquiryId")
	if inquiryIdStr == "" {
		respondWithError(w, http.StatusBadRequest, "inquiryId is required")
		return
	}
	inquiryId, err := uuid.Parse(inquiryIdStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid inquiryId")
		return
	}
	messages, err := cfg.DB.GetMessagesByInquiry(r.Context(), inquiryId)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error getting inquiry messages")
		return
	}
	respondWithJSON(w, http.StatusOK, models.DatabaseMessagesToMessages(messages))
}

func (cfg *APIConfig) HandlerUpdateMessage(w http.ResponseWriter, r *http.Request, user database.User) {
	messageIdStr := chi.URLParam(r, "messageId")
	if messageIdStr == "" {
		respondWithError(w, http.StatusBadRequest, "messageId is required")
		return
	}
	messageId, err := uuid.Parse(messageIdStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid messageId")
		return
	}

	//Permission check
	IsMessageSender, err := services.IsMessageSender(cfg.DB, r, messageId, user.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error checking message sender: %v", err))
	}
	if !IsMessageSender {
		respondWithError(w, http.StatusForbidden, "user is not the sender of the message")
		return
	}

	type parameters struct {
		Content string `json:"content"`
	}
	decoder := json.NewDecoder(r.Body)

	params := parameters{}

	err = decoder.Decode(&params)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "error parsing json")
		return
	}
	message, err := cfg.DB.UpdateMessage(r.Context(), database.UpdateMessageParams{
		ID:        messageId,
		Content:   params.Content,
		UpdatedAt: time.Now().UTC(),
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error updating inquiry message")
		return
	}
	respondWithJSON(w, http.StatusOK, models.DatabaseMessageToMessage(message))
}

// to do: authorization should be handled, sender of mssage should be able to delete the message
func (cfg *APIConfig) HandlerDeleteMessage(w http.ResponseWriter, r *http.Request, user database.User) {
	messageIdStr := chi.URLParam(r, "messageId")
	if messageIdStr == "" {
		respondWithError(w, http.StatusBadRequest, "messageId is required")
		return
	}
	messageId, err := uuid.Parse(messageIdStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid messageId")
		return
	}
	err = cfg.DB.DeleteMessage(r.Context(), database.DeleteMessageParams{
		ID: messageId,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error deleting inquiry message")
		return
	}
	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Message deleted"})
}
