package handlers

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/koredeycode/dwelly/dwelly-api/models"
	"github.com/koredeycode/dwelly/internal/database"
)

//	func (cfg *APIConfig) HandlerCreateInquiryMessage(w http.ResponseWriter, r *http.Request, user database.User) {
//		w.WriteHeader(http.StatusCreated)
//		w.Write([]byte("Inquiry message created"))
//	}
func (cfg *APIConfig) HandlerCreateInquiryMessage(w http.ResponseWriter, r *http.Request, user database.User) {
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
	respondWithJSON(w, http.StatusOK, models.DatabaseMessagestoMessages(messages))
}

func (cfg *APIConfig) HandlerUpdateInquiryMessage(w http.ResponseWriter, r *http.Request, user database.User) {

}

// to do: authorization should be handled, sender of mssage should be able to delete the message
func (cfg *APIConfig) HandlerDeleteInquiryMessage(w http.ResponseWriter, r *http.Request, user database.User) {
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
