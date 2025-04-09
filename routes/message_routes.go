package routes

import (
	"github.com/go-chi/chi"
	"github.com/koredeycode/dwelly/api"
)

func MessageRoutes(apiCfg *api.APIConfig) chi.Router {
	r := chi.NewRouter()

	r.Group(func(r chi.Router) {

		r.Delete("/{messageId}", apiCfg.MiddlewareAuth(apiCfg.HandlerDeleteInquiryMessage))
	})

	return r
}
