package handlers

import (
	"encoding/json"
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
		Message string `json:"message" validate:"required"`
	}

	decoder := json.NewDecoder(r.Body)

	params := parameters{}

	err := decoder.Decode(&params)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload", err.Error())
		return
	}

	if err := cfg.Validate.Struct(params); err != nil {
		errorMessages := utils.ExtractValidationErrors(err)
		respondWithError(w, http.StatusBadRequest, "Validation failed", errorMessages...)
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
		respondWithError(w, http.StatusInternalServerError, "Inquiry creation failure", err.Error())
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
		respondWithError(w, http.StatusInternalServerError, "Inquiry message creation failure", err.Error())
		return
	}

	respondWithSuccess(w, http.StatusCreated, "Inquiry created successfully", models.DatabaseInquiryToInquiry(inquiry))
}

// to do: authorization of current user should be handled, sender of inquiry and listing own should be able to see the inquiry
func (cfg *APIConfig) HandlerGetInquiry(w http.ResponseWriter, r *http.Request, user database.User) {
	inquiryIDStr := chi.URLParam(r, "inquiryId")
	inquiryId, errMsg := utils.GetUUIDParam(inquiryIDStr, "inquiry")

	if errMsg != "" {
		respondWithError(w, http.StatusBadRequest, errMsg)
		return
	}

	inquiry, err := cfg.DB.GetInquiryByIDWithMessages(r.Context(), inquiryId)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Inquiry retrieval failure", err.Error())
		return
	}

	respondWithSuccess(w, http.StatusOK, "Inquiry retrieved successfully", models.DatabaseInquiryToInquiry(inquiry))

}

// to do: authorization of current user should be handled, listing owner should be able to see the inquiries for their listing
func (cfg *APIConfig) HandlerGetListingInquiries(w http.ResponseWriter, r *http.Request, user database.User) {
	listingIDStr := chi.URLParam(r, "listingId")
	listingId, errMsg := utils.GetUUIDParam(listingIDStr, "listing")

	if errMsg != "" {
		respondWithError(w, http.StatusBadRequest, errMsg)
		return
	}

	inquiries, err := cfg.DB.GetListingInquiries(r.Context(), listingId)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Inquiries retrieval failure", err.Error())
		return
	}

	respondWithSuccess(w, http.StatusOK, "Inquiries retrieved successfully", models.DatabaseInquiriesToInquiries(inquiries))
}

// to do: authorization of current user should be handled, sender of inquiry and listing own should be able to update the inquiry
func (cfg *APIConfig) HandlerUpdateInquiryStatus(w http.ResponseWriter, r *http.Request, user database.User) {

	inquiryIDStr := chi.URLParam(r, "inquiryId")
	inquiryId, errMsg := utils.GetUUIDParam(inquiryIDStr, "inquiry")

	if errMsg != "" {
		respondWithError(w, http.StatusBadRequest, errMsg)
		return
	}

	type parameters struct {
		Status string `json:"status" validate:"required"`
	}
	decoder := json.NewDecoder(r.Body)

	params := parameters{}

	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload", err.Error())
		return
	}

	if err := cfg.Validate.Struct(params); err != nil {
		errorMessages := utils.ExtractValidationErrors(err)
		respondWithError(w, http.StatusBadRequest, "Validation failed", errorMessages...)
		return
	}
	statuses := []string{"active", "negotiation", "completed"}

	statusNotValid := !slices.Contains(statuses, params.Status)

	if statusNotValid {
		respondWithError(w, http.StatusBadRequest, "Invalid status")
		return
	}

	err = cfg.DB.UpdateInquiryStatus(r.Context(), database.UpdateInquiryStatusParams{
		ID:     inquiryId,
		Status: params.Status,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Inquiry update failure", err.Error())
		return

	}
	respondWithSuccess(w, http.StatusOK, "Inquiry status updated successfully", nil)

}

// to do: authorization of current user should be handled, sender of inquiry and listing own should be able to delete the inquiry
func (cfg *APIConfig) HandlerDeleteInquiry(w http.ResponseWriter, r *http.Request, user database.User) {

	inquiryIDStr := chi.URLParam(r, "inquiryId")
	inquiryId, errMsg := utils.GetUUIDParam(inquiryIDStr, "inquiry")

	if errMsg != "" {
		respondWithError(w, http.StatusBadRequest, errMsg)
		return
	}

	err := cfg.DB.DeleteInquiry(r.Context(), database.DeleteInquiryParams{
		ID:       inquiryId,
		SenderID: user.ID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Inquiry deletion failure", err.Error())
		return
	}

	respondWithSuccess(w, http.StatusOK, "Inquiry deleted successfully", nil)

}
