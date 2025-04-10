package routes

import (
	"github.com/go-chi/chi"
	"github.com/koredeycode/dwelly/dwelly-api/handlers"
)

func StatusRoute(apiCfg *handlers.APIConfig) chi.Router {
	r := chi.NewRouter()

	r.Group(func(r chi.Router) {

		r.Delete("/", apiCfg.MiddlewareAuth(apiCfg.HandlerDeleteInquiryMessage))
	})

	return r
}
