package api

import (
	"net/http"

	"github.com/koredeycode/dwelly/internal/database"
)

func (api *APIConfig) HandlerAPIStatus(w http.ResponseWriter, r *http.Request, user database.User) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("API is running"))
}
