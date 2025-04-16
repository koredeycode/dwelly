package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/koredeycode/dwelly/dwelly-api/services"
	"github.com/koredeycode/dwelly/internal/database"
)

// By url
// to do: authorization should be handled, lising owner should be able to add the image
func (cfg *APIConfig) HandlerAddListingImage(w http.ResponseWriter, r *http.Request, user database.User) {
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
		ImageURL string `json:"image_url"`
	}
	decoder := json.NewDecoder(r.Body)

	params := parameters{}

	err = decoder.Decode(&params)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("error parsing json: %v", err))
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
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error creating listing image: %v", err))
		return
	}
	respondWithJSON(w, http.StatusCreated, map[string]string{"message": "image added successfully"})
}

// to do: authorization should be handled, lising owner should be able to delete the image
func (cfg *APIConfig) HandlerDeleteListingImage(w http.ResponseWriter, r *http.Request, user database.User) {
	listingIDStr := chi.URLParam(r, "listingId")

	imageIDStr := chi.URLParam(r, "imageId")
	if imageIDStr == "" {
		respondWithError(w, http.StatusBadRequest, "imageId is required")
		return
	}

	imageId, err := uuid.Parse(listingIDStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid image id")
		return
	}

	err = cfg.DB.DeleteListingImageByID(r.Context(), imageId)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error deleting listing image")
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Image deleted"})
}
