package routes

import (
	"github.com/go-chi/chi"
	"github.com/koredeycode/dwelly/dwelly-api/handlers"
)

func InquiryRoutes(apiCfg *handlers.APIConfig) chi.Router {
	r := chi.NewRouter()

	r.Group(func(r chi.Router) {

		r.Get("/{inquiryId}", apiCfg.Auth(apiCfg.HandlerGetInquiry, apiCfg.InquirySenderOrListingOwnerAuthorization))

		r.Patch("/{inquiryId}/status", apiCfg.Auth(apiCfg.HandlerUpdateInquiryStatus, apiCfg.ListingOwnerAuthorization))
		r.Delete("/{inquiryId}", apiCfg.Auth(apiCfg.HandlerDeleteInquiry, apiCfg.InquirySenderAuthorization))

		r.Get("/{inquiryId}/messages", apiCfg.Auth(apiCfg.HandlerGetInquiryMessages, apiCfg.InquirySenderOrListingOwnerAuthorization))
		r.Post("/{inquiryId}/messages", apiCfg.Auth(apiCfg.HandlerCreateInquiryMessage, apiCfg.InquirySenderOrListingOwnerAuthorization))
	})

	return r
}
