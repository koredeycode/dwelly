package api

import (
	"net/http"

	"github.com/koredeycode/dwelly/internal/database"
)

func (api *APIConfig) HandlerCreateListing(w http.ResponseWriter, r *http.Request, user database.User) {
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Listing created"))
}

func (api *APIConfig) HandlerGetListing(w http.ResponseWriter, r *http.Request, user database.User) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Get listing"))
}

func (api *APIConfig) HandlerGetListings(w http.ResponseWriter, r *http.Request, user database.User) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("List of listings"))
}

func (api *APIConfig) HandlerUpdateListing(w http.ResponseWriter, r *http.Request, user database.User) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Listing updated"))
}

func (api *APIConfig) HandlerDeleteListing(w http.ResponseWriter, r *http.Request, user database.User) {
	w.WriteHeader(http.StatusNoContent)
}

func (api *APIConfig) HandlerUpdateListingStatus(w http.ResponseWriter, r *http.Request, user database.User) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Listing status updated"))
}
