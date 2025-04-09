package main

import (
	"github.com/go-chi/chi"

	"github.com/go-chi/cors"
)

func setUpRouter(apiCfg *apiConfig) *chi.Mux {
	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	v1Router := chi.NewRouter()

	v1Router.Get("/status", apiCfg.handlerAPIStatus)

	v1Router.Post("/auth/register", apiCfg.handlerRegisterUser)
	v1Router.Post("/auth/login", apiCfg.handlerLoginUser)

	// v1Router.Post("/users/register", apiCfg.handlerRegisterUser)
	// v1Router.Post("/users/login", apiCfg.handlerLoginUser)
	v1Router.Get("/users/me", apiCfg.middlewareAuth(apiCfg.handlerGetUser))

	v1Router.Post("/listings", apiCfg.middlewareAuth(apiCfg.handlerCreateListing))
	v1Router.Get("/listings/{listingId}", apiCfg.middlewareAuth(apiCfg.handlerGetListing))
	v1Router.Get("/listings}", apiCfg.middlewareAuth(apiCfg.handlerGetListings))
	v1Router.Put("/listings/{listingId}", apiCfg.middlewareAuth(apiCfg.handlerUpdateListing))
	v1Router.Delete("/listings/{listingId}", apiCfg.middlewareAuth(apiCfg.handlerDeleteListing))
	v1Router.Patch("/listings/{listingId}/status", apiCfg.middlewareAuth(apiCfg.handlerUpdateListingStatus))

	v1Router.Post("/listings/{listingId}/images", apiCfg.middlewareAuth(apiCfg.handlerAddListingImage))
	v1Router.Delete("/listings/images/{imageId}", apiCfg.middlewareAuth(apiCfg.handlerDeleteListingImage))

	v1Router.Post("/inquiries", apiCfg.middlewareAuth(apiCfg.handlerCreateInquiry))
	v1Router.Get("/inquiries/{inquiryId}", apiCfg.middlewareAuth(apiCfg.handlerGetInquiry))
	v1Router.Get("/inquiries}", apiCfg.middlewareAuth(apiCfg.handlerGetInquiries))
	v1Router.Patch("/inquiries/{inquiryId}/status", apiCfg.middlewareAuth(apiCfg.handlerUpdateInquiryStatus))
	v1Router.Delete("/inquiries/{inquiryId}", apiCfg.middlewareAuth(apiCfg.handlerDeleteInquiry))

	v1Router.Post("/inquiries/{inquiryId}/messages", apiCfg.middlewareAuth(apiCfg.handlerCreateInquiryMessage))
	v1Router.Get("/inquiries/{inquiryId}/messages", apiCfg.middlewareAuth(apiCfg.handlerGetInquiryMessages))
	v1Router.Delete("/messages/{messageId}", apiCfg.middlewareAuth(apiCfg.handlerDeleteInquiryMessage))

	router.Mount("/v1", v1Router)

	router.Mount("/v1", v1Router)
	return router
}
