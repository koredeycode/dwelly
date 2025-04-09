package routes

import (
	"github.com/go-chi/chi"
	"github.com/koredeycode/dwelly/api"
)

func ListingRoutes(apiCfg *api.APIConfig) chi.Router {
	r := chi.NewRouter()

	r.Group(func(r chi.Router) {
		r.Post("/", apiCfg.HandlerCreateListing)
		r.Get("/{listingId}", apiCfg.HandlerGetListing)
		r.Get("/", apiCfg.HandlerGetListings)
		r.Put("/{listingId}", apiCfg.HandlerUpdateListing)
		r.Delete("/{listingId}", apiCfg.HandlerDeleteListing)
		r.Patch("/{listingId}/status", apiCfg.HandlerUpdateListingStatus)
		r.Post("/{listingId}/images", apiCfg.HandlerAddListingImage)
		r.Delete("/images/{imageId}", apiCfg.HandlerDeleteListingImage)
	})

	return r
}
