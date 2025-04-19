package routes

import (
	"github.com/go-chi/chi"
	"github.com/koredeycode/dwelly/dwelly-api/handlers"
)

func AuthRoutes(apiCfg *handlers.APIConfig) chi.Router {
	r := chi.NewRouter()
	r.Post("/register", apiCfg.HandlerRegisterUser)
	r.Post("/login", apiCfg.HandlerLoginUser)
	r.Post("/logout", apiCfg.Auth(apiCfg.HandlerLogoutUser))
	r.Get("/me", apiCfg.Auth(apiCfg.HandlerGetCurrentUser))
	return r
}
