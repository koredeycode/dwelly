package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func respondWithError(w http.ResponseWriter, code int, message string, errors ...string) {
	response := map[string]interface{}{
		"status":  "error",
		"message": message,
	}

	if len(errors) > 0 {
		response["errors"] = errors
	}

	respondWithJSON(w, code, response)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	dat, err := json.Marshal(payload)
	if err != nil {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		errorJSON, _ := json.Marshal(map[string]interface{}{
			"status":  "error",
			"message": "failed to marshal json response",
		})
		w.Write(errorJSON)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(dat)
}

func respondWithSuccess(w http.ResponseWriter, code int, message string, data interface{}) {
	response := map[string]interface{}{
		"status":  "success",
		"message": message,
	}

	fmt.Println(response)

	if data != nil {
		response["data"] = data
	}

	respondWithJSON(w, code, response)
}
