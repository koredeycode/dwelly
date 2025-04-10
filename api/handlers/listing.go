package handlers

import (
	"net/http"

	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
	cloudinaryutil "github.com/koredeycode/dwelly/internal/cloudinary"
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
func (api *APIConfig) HandlerUploadListingImages(w http.ResponseWriter, r *http.Request, user database.User) {

	// Get all files uploaded for "file" field (could be multiple)
	files := r.MultipartForm.File["file"]
	if len(files) == 0 {
		http.Error(w, "No files uploaded", http.StatusBadRequest)
		return
	}

	// Create Cloudinary client
	cld, err := cloudinaryutil.NewClient()
	if err != nil {
		http.Error(w, "Could not create Cloudinary client", http.StatusInternalServerError)
		return
	}

	// Loop through all files and upload them
	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			http.Error(w, "Unable to open uploaded file", http.StatusBadRequest)
			return
		}
		defer file.Close()

		// Upload to Cloudinary
		uploadResp, err := cld.Upload.Upload(r.Context(), file, uploader.UploadParams{
			PublicID: fileHeader.Filename,
			Folder:   "dwelly/listings",
		})
		if err != nil {
			http.Error(w, "Failed to upload to Cloudinary", http.StatusInternalServerError)
			return
		}

		// Save the image URL to the database
		listingID, err := uuid.Parse(chi.URLParam(r, "listingId"))
		if err != nil {
			http.Error(w, "Invalid listing ID", http.StatusBadRequest)
			return
		}
		imageURL := uploadResp.SecureURL
		_, err = api.DB.AddListingImage(r.Context(), database.AddListingImageParams{
			ID:        uuid.New(),
			ListingID: listingID,
			Url:       imageURL,
		})

		if err != nil {
			http.Error(w, "Failed to save image URL to database", http.StatusInternalServerError)
			return
		}
	}

	// Return success response
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Listing images added"))
}
