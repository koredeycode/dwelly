package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/koredeycode/dwelly/dwelly-api/utils"
	"github.com/koredeycode/dwelly/internal/database"
)

// By url
// to do: authorization should be handled, lising owner should be able to add the image
func (cfg *APIConfig) HandlerAddListingImage(w http.ResponseWriter, r *http.Request, user database.User) {
	listingIDStr := chi.URLParam(r, "listingId")
	if listingIDStr == "" {
		respondWithError(w, http.StatusBadRequest, "Missing listing ID")
		return
	}

	listingId, err := uuid.Parse(listingIDStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid listing ID", err.Error())
		return
	}

	type parameters struct {
		ImageURL string `json:"image_url" validate:"required"`
	}
	decoder := json.NewDecoder(r.Body)

	params := parameters{}

	err = decoder.Decode(&params)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload", err.Error())
		return
	}

	if err := cfg.Validate.Struct(params); err != nil {
		errorMessages := utils.ExtractValidationErrors(err)
		respondWithError(w, http.StatusBadRequest, "Validation failed", errorMessages...)
		return
	}
	_, err = cfg.DB.AddListingImage(r.Context(), database.AddListingImageParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		ListingID: listingId,
		Url:       params.ImageURL,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Listing image creation failure", err.Error())
		return
	}
	respondWithSuccess(w, http.StatusOK, "Image added successfully", nil)

}

// to do: authorization should be handled, lising owner should be able to delete the image
func (cfg *APIConfig) HandlerDeleteListingImage(w http.ResponseWriter, r *http.Request, user database.User) {
	listingIDStr := chi.URLParam(r, "listingId")

	imageIDStr := chi.URLParam(r, "imageId")
	if imageIDStr == "" {
		respondWithError(w, http.StatusBadRequest, "ImageId is required")
		return
	}

	imageId, err := uuid.Parse(listingIDStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid image id", err.Error())
		return
	}

	err = cfg.DB.DeleteListingImageByID(r.Context(), imageId)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Listing image deletion failure", err.Error())
		return
	}

	respondWithSuccess(w, http.StatusOK, "Image deleted successfully", nil)
}
