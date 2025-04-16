package routes

import (
	"github.com/go-chi/chi"
	"github.com/koredeycode/dwelly/dwelly-api/handlers"
)

func ImageRoutes(apiCfg *handlers.APIConfig) chi.Router {
	r := chi.NewRouter()

	r.Group(func(r chi.Router) {
		r.Delete("/{imageId}", apiCfg.MiddlewareAuth(apiCfg.HandlerDeleteListingImage))

	})

	return r
}
