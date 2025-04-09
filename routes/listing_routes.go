package routes

import (
	"github.com/go-chi/chi"
	"github.com/koredeycode/dwelly/api"
)

func ListingRoutes(apiCfg *api.APIConfig) chi.Router {
	r := chi.NewRouter()

	r.Group(func(r chi.Router) {
		r.Post("/", apiCfg.MiddlewareAuth(apiCfg.HandlerCreateListing))
		r.Get("/{listingId}", apiCfg.MiddlewareAuth(apiCfg.HandlerGetListing))
		r.Get("/", apiCfg.MiddlewareAuth(apiCfg.HandlerGetListings))
		r.Put("/{listingId}", apiCfg.MiddlewareAuth(apiCfg.HandlerUpdateListing))
		r.Delete("/{listingId}", apiCfg.MiddlewareAuth(apiCfg.HandlerDeleteListing))
		r.Patch("/{listingId}/status", apiCfg.MiddlewareAuth(apiCfg.HandlerUpdateListingStatus))
		r.Post("/{listingId}/images", apiCfg.MiddlewareAuth(apiCfg.HandlerAddListingImage))

		r.Post("/{listingId}/images/upload", apiCfg.MiddlewareAuth(apiCfg.HandlerUploadListingImages))

		r.Delete("/images/{imageId}", apiCfg.MiddlewareAuth(apiCfg.HandlerDeleteListingImage))
	})

	return r
}
