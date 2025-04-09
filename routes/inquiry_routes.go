package routes

import (
	"github.com/go-chi/chi"
	"github.com/koredeycode/dwelly/api"
)

func InquiryRoutes(apiCfg *api.APIConfig) chi.Router {
	r := chi.NewRouter()

	r.Group(func(r chi.Router) {

		r.Post("/", apiCfg.HandlerCreateInquiry)
		r.Get("/{inquiryId}", apiCfg.HandlerGetInquiry)
		r.Get("/", apiCfg.HandlerGetInquiries)
		r.Patch("/{inquiryId}/status", apiCfg.HandlerUpdateInquiryStatus)
		r.Delete("/{inquiryId}", apiCfg.HandlerDeleteInquiry)

		r.Post("/{inquiryId}/messages", apiCfg.HandlerCreateInquiryMessage)
		r.Get("/{inquiryId}/messages", apiCfg.HandlerGetInquiryMessages)
	})

	return r
}
