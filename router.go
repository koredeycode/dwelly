package main

import ("github.com/go-chi/chi"

	"github.com/go-chi/cors")

router := chi.NewRouter()

router.Use(cors.Handler(cors.Options{
	AllowedOrigins:   []string{"https://*", "http://*"},
	AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	AllowedHeaders:   []string{"*"},
	ExposedHeaders:   []string{"Link"},
	AllowCredentials: false,
	MaxAge:           300,
}))

v1Router := chi.NewRouter()

// v1Router.Get("/ready", handlerReadiness)
// v1Router.Get("/err", handlerErr)

v1Router.Post("/users/register", apiCfg.handlerRegisterUser)
v1Router.Post("/users/login", apiCfg.handlerLoginUser)
v1Router.Get("/users/me", apiCfg.middlewareAuth(apiCfg.handlerGetUser))

v1Router.Post("/listings", apiCfg.middlewareAuth(apiCfg.handlerCreateListing))
v1Router.Get("/listings/{listingId}", apiCfg.handlerGetListing)
v1Router.Get("/listings}", apiCfg.handlerGetListings)
v1Router.Put("/listings/{listingId}", apiCfg.handlerUpdateListing)
v1Router.Delete("/listings/{listingId}", apiCfg.handlerDeleteListing)
v1Router.Patch("/listings/{listingId}/status", apiCfg.handlerUpdateListingStatus)

v1Router.Post("/listings/{listingId}/images", apiCfg.middlewareAuth(apiCfg.handlerCreateListing))
v1Router.Delete("/listings/images/{imageId}", apiCfg.handlerDeleteListing)


router.Mount("/v1", v1Router)