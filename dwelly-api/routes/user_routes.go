package routes

import (
	"github.com/go-chi/chi"
	"github.com/koredeycode/dwelly/dwelly-api/handlers"
)

func UserRoutes(apiCfg *handlers.APIConfig) chi.Router {
	r := chi.NewRouter()

	r.Get("/{userId}", apiCfg.MiddlewareAuth(apiCfg.HandlerGetUser))
	// r.Patch("/{userId}", apiCfg.MiddlewareAuth(apiCfg.HandlerUpdateUser))
	// r.Delete("/{userId}", apiCfg.MiddlewareAuth(apiCfg.HandlerDeleteUser))

	// r.Get("/{userId}/listings", apiCfg.MiddlewareAuth(apiCfg.HandlerGetUserListings))
	return r
}
