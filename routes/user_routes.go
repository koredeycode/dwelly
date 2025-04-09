package routes

import (
	"github.com/go-chi/chi"
	"github.com/koredeycode/dwelly/api"
)

func UserRoutes(apiCfg *api.APIConfig) chi.Router {
	r := chi.NewRouter()
	r.Get("/me", apiCfg.MiddlewareAuth(apiCfg.HandlerGetUser))
	return r
}
