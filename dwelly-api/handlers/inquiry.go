package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"slices"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/koredeycode/dwelly/dwelly-api/models"
	"github.com/koredeycode/dwelly/dwelly-api/utils"
	"github.com/koredeycode/dwelly/internal/database"
)

func (cfg *APIConfig) HandlerCreateListingInquiry(w http.ResponseWriter, r *http.Request, user database.User) {
	listingIDStr := chi.URLParam(r, "listingId")

	listingId, errMsg := utils.GetUUIDParam(listingIDStr, "listing")

	if errMsg != "" {
		respondWithError(w, http.StatusBadRequest, errMsg)
		return
	}

	type parameters struct {
		Message string `json:"message"`
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
		ListingID: listingId,
		// Message: params.Message,
	})

	if err != nil {
		respondWithError(w, http.StatusBadGateway, "error creating inquiry")
		return
	}

	_, err = cfg.DB.CreateMessage(r.Context(), database.CreateMessageParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		SenderID:  user.ID,
		InquiryID: inquiry.ID,
		Content:   params.Message,
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error creating message to inquiry")
		return
	}

	respondWithJSON(w, http.StatusCreated, models.DatabaseInquiryToInquiry(inquiry))
}

// to do: authorization of current user should be handled, sender of inquiry and listing own should be able to see the inquiry
func (cfg *APIConfig) HandlerGetInquiry(w http.ResponseWriter, r *http.Request, user database.User) {
	listingIDStr := chi.URLParam(r, "listingId")

	listingId, errMsg := utils.GetUUIDParam(listingIDStr, "listing")

	if errMsg != "" {
		respondWithError(w, http.StatusBadRequest, errMsg)
		return
	}

	inquiryIDStr := chi.URLParam(r, "inquiryId")
	inquiryId, errMsg := utils.GetUUIDParam(inquiryIDStr, "inquiry")

	if errMsg != "" {
		respondWithError(w, http.StatusBadRequest, errMsg)
		return
	}

	//Permission check
	isInquirySender, err := utils.IsInquirySender(cfg.DB, r, inquiryId, user.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error checking inquiry sender: %v", err))
	}

	isListingOwner, err := utils.IsListingOwner(cfg.DB, r, listingId, user.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error checking listing owner: %v", err))
	}

	if !isInquirySender || !isListingOwner {
		respondWithError(w, http.StatusForbidden, "user is not the owner of the listing or sender of the inquiry")
		return
	}

	inquiry, err := cfg.DB.GetInquiryByIDWithMessages(r.Context(), inquiryId)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error getting inquiry")
		return
	}

	respondWithJSON(w, http.StatusOK, models.DatabaseInquiryToInquiry(inquiry))

}

// to do: authorization of current user should be handled, listing owner should be able to see the inquiries for their listing
func (cfg *APIConfig) HandlerGetListingInquiries(w http.ResponseWriter, r *http.Request, user database.User) {
	listingIDStr := chi.URLParam(r, "listingId")
	listingId, errMsg := utils.GetUUIDParam(listingIDStr, "listing")

	if errMsg != "" {
		respondWithError(w, http.StatusBadRequest, errMsg)
		return
	}
	//Permission check
	isListingOwner, err := utils.IsListingOwner(cfg.DB, r, listingId, user.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error checking listing owner: %v", err))
	}

	if !isListingOwner {
		respondWithError(w, http.StatusForbidden, "user is not the owner of the listing")
		return
	}

	inquiries, err := cfg.DB.GetListingInquiries(r.Context(), listingId)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error getting inquiries")
		return
	}

	respondWithJSON(w, http.StatusOK, models.DatabaseInquiriesToInquiries(inquiries))
}

// to do: authorization of current user should be handled, sender of inquiry and listing own should be able to update the inquiry
func (cfg *APIConfig) HandlerUpdateInquiryStatus(w http.ResponseWriter, r *http.Request, user database.User) {
	listingIDStr := chi.URLParam(r, "listingId")

	listingId, errMsg := utils.GetUUIDParam(listingIDStr, "listing")

	if errMsg != "" {
		respondWithError(w, http.StatusBadRequest, errMsg)
		return
	}

	inquiryIDStr := chi.URLParam(r, "inquiryId")
	inquiryId, errMsg := utils.GetUUIDParam(inquiryIDStr, "inquiry")

	if errMsg != "" {
		respondWithError(w, http.StatusBadRequest, errMsg)
		return
	}

	//Permission check
	isListingOwner, err := utils.IsListingOwner(cfg.DB, r, listingId, user.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error checking listing owner: %v", err))
	}

	if !isListingOwner {
		respondWithError(w, http.StatusForbidden, "user is not the owner of the listing")
		return
	}

	type parameters struct {
		Status string `json:"status"`
	}
	decoder := json.NewDecoder(r.Body)

	params := parameters{}

	err = decoder.Decode(&params)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("error parsing json: %v", err))
		return
	}
	statuses := []string{"active", "negotiation", "completed"}

	statusNotValid := !slices.Contains(statuses, params.Status)

	if statusNotValid {
		respondWithError(w, http.StatusBadRequest, "invalid status")
		return
	}

	err = cfg.DB.UpdateInquiryStatus(r.Context(), database.UpdateInquiryStatusParams{
		ID:     inquiryId,
		Status: params.Status,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error updating inquiry status: %v", err))
		return

	}
	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Inquiry status updated"})
}

// to do: authorization of current user should be handled, sender of inquiry and listing own should be able to delete the inquiry
func (cfg *APIConfig) HandlerDeleteInquiry(w http.ResponseWriter, r *http.Request, user database.User) {

	inquiryIDStr := chi.URLParam(r, "inquiryId")
	inquiryId, errMsg := utils.GetUUIDParam(inquiryIDStr, "inquiry")

	if errMsg != "" {
		respondWithError(w, http.StatusBadRequest, errMsg)
		return
	}

	//Permission check
	isInquirySender, err := utils.IsInquirySender(cfg.DB, r, inquiryId, user.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error checking inquiry sender: %v", err))
	}

	if !isInquirySender {
		respondWithError(w, http.StatusForbidden, "user is not the sender of the inquiry")
		return
	}

	err = cfg.DB.DeleteInquiry(r.Context(), database.DeleteInquiryParams{
		ID:       inquiryId,
		SenderID: user.ID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error deleting inquiry")
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{
		"message": "Inquiry deleted",
	})
}
