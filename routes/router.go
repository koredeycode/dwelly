package routes

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

	r.Route("/api/v1", func(api chi.Router) {
		api.Mount("/auth", AuthRoutes(apiCfg))
		api.Mount("status", StatusRoute(apiCfg))
		api.Mount("/users", UserRoutes(apiCfg))
		api.Mount("/listings", ListingRoutes(apiCfg))
		api.Mount("/inquiries", InquiryRoutes(apiCfg))
		api.Mount("/messages", MessageRoutes(apiCfg))
	})

	return router
}
