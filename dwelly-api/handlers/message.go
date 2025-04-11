package handlers

import (
	"net/http"

	"github.com/koredeycode/dwelly/internal/database"
)

func (cfg *APIConfig) HandlerCreateInquiryMessage(w http.ResponseWriter, r *http.Request, user database.User) {
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Inquiry message created"))
}

func (cfg *APIConfig) HandlerGetInquiryMessages(w http.ResponseWriter, r *http.Request, user database.User) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Messages in inquiry"))
}

func (cfg *APIConfig) HandlerDeleteInquiryMessage(w http.ResponseWriter, r *http.Request, user database.User) {
	w.WriteHeader(http.StatusNoContent)
}
