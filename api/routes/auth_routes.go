package routes

import (
	"github.com/go-chi/chi"
	"github.com/koredeycode/dwelly/api/handlers"
)

func AuthRoutes(apiCfg *handlers.APIConfig) chi.Router {
	r := chi.NewRouter()
	r.Post("/register", apiCfg.HandlerRegisterUser)
	r.Post("/login", apiCfg.HandlerLoginUser)
	return r
}
