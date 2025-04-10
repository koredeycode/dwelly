package routes

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"

	"github.com/go-chi/cors"
	"github.com/koredeycode/dwelly/dwelly-api/handlers"
)

func SetUpRouter(apiCfg *handlers.APIConfig) *chi.Mux {
	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Printf("%s %s\n", r.Method, r.URL.Path)
			next.ServeHTTP(w, r)
		})
	})

	router.Route("/api/v1", func(api chi.Router) {
		api.Mount("/auth", AuthRoutes(apiCfg))
		api.Mount("/status", StatusRoute(apiCfg))
		api.Mount("/users", UserRoutes(apiCfg))
		api.Mount("/listings", ListingRoutes(apiCfg))
		api.Mount("/inquiries", InquiryRoutes(apiCfg))
		api.Mount("/messages", MessageRoutes(apiCfg))
	})

	return router
}
