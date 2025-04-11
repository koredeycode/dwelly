package handlers

import (
	"net/http"

	"github.com/koredeycode/dwelly/dwelly-api/models"
	"github.com/koredeycode/dwelly/internal/database"
)

func (api *APIConfig) HandlerGetUser(w http.ResponseWriter, r *http.Request, user database.User) {
	respondWithJSON(w, http.StatusOK, models.DatabaseUsertoUser(user))
}
