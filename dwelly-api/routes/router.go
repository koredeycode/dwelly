package routes

import (
	"fmt"
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

// Print all routes implemented and method
func PrintRoutes(r chi.Router) {
	err := chi.Walk(r, func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		fmt.Printf("%-7s -> %s\n", method, route)
		return nil
	})
	if err != nil {
		fmt.Printf("Failed to walk routes: %s\n", err.Error())
	}
}
