package routes

import (
	"github.com/go-chi/chi"
	"github.com/koredeycode/dwelly/api"
)

func InquiryRoutes(apiCfg *api.APIConfig) chi.Router {
	r := chi.NewRouter()

	r.Group(func(r chi.Router) {

		r.Post("/", apiCfg.MiddlewareAuth(apiCfg.HandlerCreateInquiry))
		r.Get("/{inquiryId}", apiCfg.MiddlewareAuth(apiCfg.HandlerGetInquiry))
		r.Get("/", apiCfg.MiddlewareAuth(apiCfg.HandlerGetInquiries))
		r.Patch("/{inquiryId}/status", apiCfg.MiddlewareAuth(apiCfg.HandlerUpdateInquiryStatus))
		r.Delete("/{inquiryId}", apiCfg.MiddlewareAuth(apiCfg.HandlerDeleteInquiry))

		r.Post("/{inquiryId}/messages", apiCfg.MiddlewareAuth(apiCfg.HandlerCreateInquiryMessage))
		r.Get("/{inquiryId}/messages", apiCfg.MiddlewareAuth(apiCfg.HandlerGetInquiryMessages))
	})

	return r
}
