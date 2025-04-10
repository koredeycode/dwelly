package handlers

import (
	"net/http"

	"github.com/koredeycode/dwelly/internal/database"
)

func (api *APIConfig) HandlerCreateInquiry(w http.ResponseWriter, r *http.Request, user database.User) {
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Inquiry created"))
}

func (api *APIConfig) HandlerGetInquiry(w http.ResponseWriter, r *http.Request, user database.User) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Single inquiry"))
}

func (api *APIConfig) HandlerGetInquiries(w http.ResponseWriter, r *http.Request, user database.User) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("List of inquiries"))
}

func (api *APIConfig) HandlerUpdateInquiryStatus(w http.ResponseWriter, r *http.Request, user database.User) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Inquiry status updated"))
}

func (api *APIConfig) HandlerDeleteInquiry(w http.ResponseWriter, r *http.Request, user database.User) {
	w.WriteHeader(http.StatusNoContent)
}
