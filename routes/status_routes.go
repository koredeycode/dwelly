package routes

import (
	"github.com/go-chi/chi"
	"github.com/koredeycode/dwelly/api"
)

func StatusRoute(apiCfg *api.APIConfig) chi.Router {
	r := chi.NewRouter()

	r.Group(func(r chi.Router) {

		r.Delete("/status", apiCfg.HandlerDeleteInquiryMessage)
	})

	return r
}
