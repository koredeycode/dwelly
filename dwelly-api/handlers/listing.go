package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/koredeycode/dwelly/dwelly-api/models"
	"github.com/koredeycode/dwelly/dwelly-api/utils"
	"github.com/koredeycode/dwelly/internal/database"
)

func (cfg *APIConfig) HandlerCreateListing(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		Intent      string `json:"intent" validate:"required"`
		Title       string `json:"title" validate:"required"`
		Description string `json:"description" validate:"required"`
		Price       string `json:"price" validate:"required"`
		Location    string `json:"location" validate:"required"`
		Category    string `json:"category" validate:"required"`
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
	listing, err := cfg.DB.CreateListing(r.Context(), database.CreateListingParams{
		ID:          uuid.New(),
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
		UserID:      user.ID,
		Intent:      params.Intent,
		Title:       params.Title,
		Description: params.Description,
		Price:       params.Price,
		Location:    params.Location,
		Category:    params.Category,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error creating listing: %v", err))
		return
	}

	respondWithSuccess(w, http.StatusCreated, "Listing created successfully", models.DatabaseListingToListing(listing))

}

func (cfg *APIConfig) HandlerGetListing(w http.ResponseWriter, r *http.Request) {
	listingIDStr := chi.URLParam(r, "listingId")

	listingId, errMsg := utils.GetUUIDParam(listingIDStr, "listing")

	if errMsg != "" {
		respondWithError(w, http.StatusBadRequest, errMsg)
		return
	}

	listing, err := cfg.DB.GetListingByID(r.Context(), listingId)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Listing retrieval failure", err.Error())
		return
	}

	respondWithSuccess(w, http.StatusOK, "Listing retrieved successfully", models.DatabaseListingToListing(listing))

}

func (cfg *APIConfig) HandlerGetListings(w http.ResponseWriter, r *http.Request) {
	listings, err := cfg.DB.ListAllListings(r.Context(), database.ListAllListingsParams{
		Limit:  10,
		Offset: 0,
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Listings retrieval failure", err.Error())
		return
	}

	// Respond with the listings
	respondWithSuccess(w, http.StatusOK, "Listings retrieved successfully", models.DatabaseListingsToListings(listings))

}

func (api *APIConfig) HandlerSearchListings(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	location := query.Get("location")
	category := query.Get("category")
	intent := query.Get("intent")

	listings, err := api.DB.SearchListings(r.Context(), database.SearchListingsParams{
		Column1: location,
		Column2: category,
		Column3: intent,
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Listing retrieval failure", err.Error())
		return
	}

	respondWithSuccess(w, http.StatusOK, "Listings retrieved successfully", models.DatabaseListingsToListings(listings))
}

func valueOrDefault(oldValue, newValue string) string {
	if strings.TrimSpace(newValue) == "" {
		return oldValue
	}
	return newValue
}

// to do: authorization should be handled, listing owner should be able to update the listing
func (cfg *APIConfig) HandlerUpdateListing(w http.ResponseWriter, r *http.Request, user database.User) {
	listingIDStr := chi.URLParam(r, "listingId")

	listingID, errMsg := utils.GetUUIDParam(listingIDStr, "listing")

	if errMsg != "" {
		respondWithError(w, http.StatusBadRequest, errMsg)
		return
	}

	type parameters struct {
		Intent      string `json:"intent"`
		Title       string `json:"title"`
		Description string `json:"description"`
		Price       string `json:"price"`
		Location    string `json:"location"`
		Category    string `json:"category"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}

	err := decoder.Decode(&params)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload", err.Error())
		return
	}

	currentListing, err := cfg.DB.GetListingByID(r.Context(), listingID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Listing not found", err.Error())
		return
	}

	listing, err := cfg.DB.UpdateListing(r.Context(), database.UpdateListingParams{
		ID:          listingID,
		Intent:      valueOrDefault(currentListing.Intent, params.Intent),
		Title:       valueOrDefault(currentListing.Title, params.Title),
		Description: valueOrDefault(currentListing.Description, params.Description),
		Price:       valueOrDefault(currentListing.Price, params.Price),
		Location:    valueOrDefault(currentListing.Location, params.Location),
		Category:    valueOrDefault(currentListing.Category, params.Category),
		UpdatedAt:   time.Now(),
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Listing update failure", err.Error())
		return
	}

	respondWithSuccess(w, http.StatusOK, "Listing updated successfully", models.DatabaseListingToListing(listing))
}

// to do: authorization should be handled, listing owner should be able to delete the listing
func (cfg *APIConfig) HandlerDeleteListing(w http.ResponseWriter, r *http.Request, user database.User) {
	listingIDStr := chi.URLParam(r, "listingId")

	listingID, errMsg := utils.GetUUIDParam(listingIDStr, "listing")

	if errMsg != "" {
		respondWithError(w, http.StatusBadRequest, errMsg)
		return
	}

	err := cfg.DB.DeleteListing(r.Context(), database.DeleteListingParams{
		ID:     listingID,
		UserID: user.ID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Listing deletion failure", err.Error())
		return
	}
	respondWithSuccess(w, http.StatusOK, "Listing deleted successfully", nil)

}

// to do: authorization should be handled, listing owner should be able to update the listing status
func (cfg *APIConfig) HandlerUpdateListingStatus(w http.ResponseWriter, r *http.Request, user database.User) {
	listingIDStr := chi.URLParam(r, "listingId")

	listingID, errMsg := utils.GetUUIDParam(listingIDStr, "listing")

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

	err = cfg.DB.UpdateListingStatus(r.Context(), database.UpdateListingStatusParams{
		ID:     listingID,
		Status: params.Status,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Listing status update failure", err.Error())
		return

	}
	respondWithSuccess(w, http.StatusOK, "Listing status updated successfully", nil)

}

func (cfg *APIConfig) HandlerUploadListingImages(w http.ResponseWriter, r *http.Request, user database.User) {
	listingIDStr := chi.URLParam(r, "listingId")

	listingID, errMsg := utils.GetUUIDParam(listingIDStr, "listing")

	if errMsg != "" {
		respondWithError(w, http.StatusBadRequest, errMsg)
		return
	}

	// Get all files uploaded for "file" field (could be multiple)
	files := r.MultipartForm.File["file"]
	if len(files) == 0 {
		respondWithError(w, http.StatusBadRequest, "No files uploaded")
		return
	}

	// Loop through all files and upload them
	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Unable to open uploaded file", err.Error())
			return
		}
		defer file.Close()

		// Upload to Cloudinary
		uploadResp, err := cfg.Cloudinary.Upload.Upload(r.Context(), file, uploader.UploadParams{
			PublicID: fileHeader.Filename,
			Folder:   "dwelly/listings",
		})
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Failed to upload image", err.Error())
			return
		}

		imageURL := uploadResp.SecureURL
		_, err = cfg.DB.AddListingImage(r.Context(), database.AddListingImageParams{
			ID:        uuid.New(),
			ListingID: listingID,
			Url:       imageURL,
		})

		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Failure to save iimage", err.Error())
			return
		}
	}

	respondWithSuccess(w, http.StatusOK, "Listing images added successfully", nil)

}
