package routes

import (
	"github.com/go-chi/chi"
	"github.com/koredeycode/dwelly/dwelly-api/handlers"
)

func UserRoutes(apiCfg *handlers.APIConfig) chi.Router {
	r := chi.NewRouter()

	r.Get("/{userId}", apiCfg.Auth(apiCfg.HandlerGetUser))
	r.Patch("/{userId}", apiCfg.Auth(apiCfg.HandlerUpdateUser))
	// r.Delete("/{userId}", apiCfg.Auth(apiCfg.HandlerDeleteUser))

	// r.Get("/{userId}/listings", apiCfg.Auth(apiCfg.HandlerGetUserListings))
	return r
}
