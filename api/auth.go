package api

import (
	"net/http"
)

func (api *APIConfig) HandlerRegisterUser(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User registered"))
}

func (api *APIConfig) HandlerLoginUser(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User logged in"))
}
