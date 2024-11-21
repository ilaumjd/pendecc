package handlers

import (
	"encoding/json"
	"net/http"
)

func respondWithError(w http.ResponseWriter, code int) {
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"error": "` + http.StatusText(code) + `"}`))
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json")
	response, _ := json.Marshal(payload)
	w.Write(response)
}
