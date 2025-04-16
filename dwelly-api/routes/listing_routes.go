package routes

import (
	"github.com/go-chi/chi"
	"github.com/koredeycode/dwelly/dwelly-api/handlers"
)

func ListingRoutes(apiCfg *handlers.APIConfig) chi.Router {
	r := chi.NewRouter()

	r.Group(func(r chi.Router) {
		// could be merged?
		r.Get("/search", apiCfg.MiddlewareAuth(apiCfg.HandlerSearchListings))
		r.Get("/", apiCfg.MiddlewareAuth(apiCfg.HandlerGetListings))

		r.Post("/", apiCfg.MiddlewareAuth(apiCfg.HandlerCreateListing))
		r.Get("/{listingId}", apiCfg.HandlerGetListing)
		r.Delete("/{listingId}", apiCfg.MiddlewareAuth(apiCfg.HandlerDeleteListing))

		r.Put("/{listingId}", apiCfg.MiddlewareAuth(apiCfg.HandlerUpdateListing))
		r.Patch("/{listingId}/status", apiCfg.MiddlewareAuth(apiCfg.HandlerUpdateListingStatus))

		r.Post("/{listingId}/images", apiCfg.MiddlewareAuth(apiCfg.HandlerAddListingImage))
		r.Post("/{listingId}/images/upload", apiCfg.MiddlewareAuth(apiCfg.HandlerUploadListingImages))

		r.Delete("/images/{imageId}", apiCfg.MiddlewareAuth(apiCfg.HandlerDeleteListingImage))

		r.Get("{listingId}/inquiries", apiCfg.MiddlewareAuth(apiCfg.HandlerGetListingInquiries))
		r.Post("{listingId}/inquiries/", apiCfg.MiddlewareAuth(apiCfg.HandlerCreateListingInquiry))

	})

	return r
}
