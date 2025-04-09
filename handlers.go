package main

import (
	"net/http"

	"github.com/koredeycode/dwelly/internal/database"
)

func (apiCfg *apiConfig) handlerAPIStatus(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (apiCfg *apiConfig) handlerRegisterUser(w http.ResponseWriter, r *http.Request) {
}
func (apiCfg *apiConfig) handlerLoginUser(w http.ResponseWriter, r *http.Request) {
}

func (apiCfg *apiConfig) handlerGetUser(w http.ResponseWriter, r *http.Request, user database.User) {
}

func (apiCfg *apiConfig) handlerCreateListing(w http.ResponseWriter, r *http.Request, user database.User) {
}

func (apiCfg *apiConfig) handlerGetListing(w http.ResponseWriter, r *http.Request, user database.User) {
}

func (apiCfg *apiConfig) handlerGetListings(w http.ResponseWriter, r *http.Request, user database.User) {
}

func (apiCfg *apiConfig) handlerDeleteListing(w http.ResponseWriter, r *http.Request, user database.User) {
}

func (apiCfg *apiConfig) handlerUpdateListing(w http.ResponseWriter, r *http.Request, user database.User) {
}

func (apiCfg *apiConfig) handlerUpdateListingStatus(w http.ResponseWriter, r *http.Request, user database.User) {
}

func (apiCfg *apiConfig) handlerAddListingImage(w http.ResponseWriter, r *http.Request, user database.User) {
}

func (apiCfg *apiConfig) handlerDeleteListingImage(w http.ResponseWriter, r *http.Request, user database.User) {
}

func (apiCfg *apiConfig) handlerCreateInquiry(w http.ResponseWriter, r *http.Request, user database.User) {
}

func (apiCfg *apiConfig) handlerGetInquiry(w http.ResponseWriter, r *http.Request, user database.User) {
}

func (apiCfg *apiConfig) handlerGetInquiries(w http.ResponseWriter, r *http.Request, user database.User) {
}

func (apiCfg *apiConfig) handlerUpdateInquiryStatus(w http.ResponseWriter, r *http.Request, user database.User) {
}

func (apiCfg *apiConfig) handlerDeleteInquiry(w http.ResponseWriter, r *http.Request, user database.User) {
}

func (apiCfg *apiConfig) handlerCreateInquiryMessage(w http.ResponseWriter, r *http.Request, user database.User) {
}

func (apiCfg *apiConfig) handlerGetInquiryMessages(w http.ResponseWriter, r *http.Request, user database.User) {
}

func (apiCfg *apiConfig) handlerDeleteInquiryMessage(w http.ResponseWriter, r *http.Request, user database.User) {
}
