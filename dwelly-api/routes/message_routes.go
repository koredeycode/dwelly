package routes

import (
	"github.com/go-chi/chi"
	"github.com/koredeycode/dwelly/dwelly-api/handlers"
)

func MessageRoutes(apiCfg *handlers.APIConfig) chi.Router {
	r := chi.NewRouter()

	r.Group(func(r chi.Router) {

		r.Delete("/{messageId}", apiCfg.MiddlewareAuth(apiCfg.HandlerDeleteInquiryMessage))
	})

	return r
}
