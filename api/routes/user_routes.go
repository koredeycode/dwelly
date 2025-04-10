package routes

import (
	"github.com/go-chi/chi"
	"github.com/koredeycode/dwelly/api/handlers"
)

func UserRoutes(apiCfg *handlers.APIConfig) chi.Router {
	r := chi.NewRouter()
	r.Get("/me", apiCfg.MiddlewareAuth(apiCfg.HandlerGetUser))
	return r
}
