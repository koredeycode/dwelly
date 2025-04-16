package routes

import (
	"github.com/go-chi/chi"
	"github.com/koredeycode/dwelly/dwelly-api/handlers"
)

func InquiryRoutes(apiCfg *handlers.APIConfig) chi.Router {
	r := chi.NewRouter()

	r.Group(func(r chi.Router) {

		r.Get("/{inquiryId}", apiCfg.MiddlewareAuth(apiCfg.HandlerGetInquiry))

		r.Patch("/{inquiryId}/status", apiCfg.MiddlewareAuth(apiCfg.HandlerUpdateInquiryStatus))
		r.Delete("/{inquiryId}", apiCfg.MiddlewareAuth(apiCfg.HandlerDeleteInquiry))

		r.Get("/{inquiryId}/messages", apiCfg.MiddlewareAuth(apiCfg.HandlerGetInquiryMessages))
		r.Post("/{inquiryId}/messages", apiCfg.MiddlewareAuth(apiCfg.HandlerCreateInquiryMessage))
	})

	return r
}
