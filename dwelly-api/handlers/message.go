package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/koredeycode/dwelly/dwelly-api/models"
	"github.com/koredeycode/dwelly/dwelly-api/utils"
	"github.com/koredeycode/dwelly/internal/database"
)

func (cfg *APIConfig) HandlerCreateInquiryMessage(w http.ResponseWriter, r *http.Request, user database.User) {
	inquiryIDStr := chi.URLParam(r, "inquiryId")
	inquiryId, errMsg := utils.GetUUIDParam(inquiryIDStr, "inquiry")

	if errMsg != "" {
		respondWithError(w, http.StatusBadRequest, errMsg)
		return
	}

	type parameters struct {
		Content string `json:"content"`
	}
	decoder := json.NewDecoder(r.Body)

	params := parameters{}

	err := decoder.Decode(&params)

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
	inquiryIDStr := chi.URLParam(r, "inquiryId")
	inquiryId, errMsg := utils.GetUUIDParam(inquiryIDStr, "inquiry")

	if errMsg != "" {
		respondWithError(w, http.StatusBadRequest, errMsg)
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
	messageIDStr := chi.URLParam(r, "messageId")

	messageId, errMsg := utils.GetUUIDParam(messageIDStr, "message")

	if errMsg != "" {
		respondWithError(w, http.StatusBadRequest, errMsg)
		return
	}

	type parameters struct {
		Content string `json:"content"`
	}
	decoder := json.NewDecoder(r.Body)

	params := parameters{}

	err := decoder.Decode(&params)

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
	messageIDStr := chi.URLParam(r, "messageId")

	messageId, errMsg := utils.GetUUIDParam(messageIDStr, "message")

	if errMsg != "" {
		respondWithError(w, http.StatusBadRequest, errMsg)
		return
	}
	err := cfg.DB.DeleteMessage(r.Context(), database.DeleteMessageParams{
		ID: messageId,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error deleting inquiry message")
		return
	}
	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Message deleted"})
}
