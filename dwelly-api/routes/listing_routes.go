package routes

import (
	"github.com/go-chi/chi"
	"github.com/koredeycode/dwelly/dwelly-api/handlers"
)

func ListingRoutes(apiCfg *handlers.APIConfig) chi.Router {
	r := chi.NewRouter()

	r.Group(func(r chi.Router) {
		// could be merged?
		r.Get("/search", apiCfg.HandlerSearchListings)
		r.Get("/", apiCfg.HandlerGetListings)

		r.Post("/", apiCfg.Auth(apiCfg.HandlerCreateListing))
		r.Get("/{listingId}", apiCfg.HandlerGetListing)
		r.Delete("/{listingId}", apiCfg.Auth(apiCfg.HandlerDeleteListing, apiCfg.ListingOwnerAuthorization))

		r.Put("/{listingId}", apiCfg.Auth(apiCfg.HandlerUpdateListing, apiCfg.ListingOwnerAuthorization))
		r.Patch("/{listingId}/status", apiCfg.Auth(apiCfg.HandlerUpdateListingStatus, apiCfg.ListingOwnerAuthorization))

		r.Post("/{listingId}/images", apiCfg.Auth(apiCfg.HandlerAddListingImage, apiCfg.ListingOwnerAuthorization))
		r.Post("/{listingId}/images/upload", apiCfg.Auth(apiCfg.HandlerUploadListingImages, apiCfg.ListingOwnerAuthorization))

		r.Delete("/images/{imageId}", apiCfg.Auth(apiCfg.HandlerDeleteListingImage, apiCfg.ListingOwnerAuthorization))
		r.Delete("/{listingId}/images/{imageId}", apiCfg.Auth(apiCfg.HandlerDeleteListingImage, apiCfg.ListingOwnerAuthorization))

		r.Get("/{listingId}/inquiries", apiCfg.Auth(apiCfg.HandlerGetListingInquiries, apiCfg.ListingOwnerAuthorization))
		r.Post("/{listingId}/inquiries/", apiCfg.Auth(apiCfg.HandlerCreateListingInquiry, apiCfg.ListingOwnerAuthorization))

	})

	return r
}
