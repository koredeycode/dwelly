package routes

import (
	"github.com/go-chi/chi"
	"github.com/koredeycode/dwelly/dwelly-api/handlers"
)

func StatusRoute(apiCfg *handlers.APIConfig) chi.Router {
	r := chi.NewRouter()

	r.Get("/", apiCfg.HandlerAPIStatus)

	return r
}
