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
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("error parsing json: %v", err))
		return
	}

	if err := cfg.Validate.Struct(params); err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Validation error: %v", err))
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

	respondWithJSON(w, http.StatusCreated, models.DatabaseListingToListing(listing))

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
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error getting listing: %v", err))
		return
	}

	respondWithJSON(w, http.StatusOK, models.DatabaseListingToListing(listing))
}

func (cfg *APIConfig) HandlerGetListings(w http.ResponseWriter, r *http.Request) {
	listings, err := cfg.DB.ListAllListings(r.Context(), database.ListAllListingsParams{
		Limit:  10,
		Offset: 0,
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error getting listings: %v", err))
		return
	}

	// Respond with the listings
	respondWithJSON(w, http.StatusOK, models.DatabaseListingsToListings(listings))

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
		respondWithError(w, http.StatusInternalServerError, "Error fetching listings")
		return
	}

	respondWithJSON(w, http.StatusOK, models.DatabaseListingsToListings(listings))
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
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	currentListing, err := cfg.DB.GetListingByID(r.Context(), listingID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, fmt.Sprintf("Listing not found: %v", err))
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
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to update listing: %v", err))
		return
	}

	respondWithJSON(w, http.StatusOK, models.DatabaseListingToListing(listing))
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
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error deleting listing: %v", err))
		return
	}
	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Listing deleted"})
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
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("error parsing json: %v", err))
		return
	}

	if err := cfg.Validate.Struct(params); err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Validation error: %v", err))
		return
	}

	statuses := []string{"active", "negotiation", "completed"}

	statusNotValid := !slices.Contains(statuses, params.Status)

	if statusNotValid {
		respondWithError(w, http.StatusBadRequest, "invalid status")
		return
	}

	err = cfg.DB.UpdateListingStatus(r.Context(), database.UpdateListingStatusParams{
		ID:     listingID,
		Status: params.Status,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error updating listing status: %v", err))
		return

	}
	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Listing status updated"})
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
		http.Error(w, "No files uploaded", http.StatusBadRequest)
		return
	}

	// Create Cloudinary client
	// cld, err := cloudinaryutil.NewClient()
	// if err != nil {
	// 	http.Error(w, "Could not create Cloudinary client", http.StatusInternalServerError)
	// 	return
	// }

	// Loop through all files and upload them
	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			http.Error(w, "Unable to open uploaded file", http.StatusBadRequest)
			return
		}
		defer file.Close()

		// Upload to Cloudinary
		uploadResp, err := cfg.Cloudinary.Upload.Upload(r.Context(), file, uploader.UploadParams{
			PublicID: fileHeader.Filename,
			Folder:   "dwelly/listings",
		})
		if err != nil {
			http.Error(w, "Failed to upload to Cloudinary", http.StatusInternalServerError)
			return
		}

		// // Save the image URL to the database
		// listingID, err := uuid.Parse(chi.URLParam(r, "listingId"))
		// if err != nil {
		// 	http.Error(w, "Invalid listing ID", http.StatusBadRequest)
		// 	return
		// }
		imageURL := uploadResp.SecureURL
		_, err = cfg.DB.AddListingImage(r.Context(), database.AddListingImageParams{
			ID:        uuid.New(),
			ListingID: listingID,
			Url:       imageURL,
		})

		if err != nil {
			http.Error(w, "Failed to save image URL to database", http.StatusInternalServerError)
			return
		}
	}

	respondWithJSON(w, http.StatusCreated, map[string]string{"message": "Listing images added"})
}
