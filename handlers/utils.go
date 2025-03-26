package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
)

func respondWithError(w http.ResponseWriter, code int, message ...string) {
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json")
	errorMessage := http.StatusText(code)
	if len(message) > 0 {
		errorMessage = message[0]
	}
	response, _ := json.Marshal(map[string]string{
		"code":  strconv.Itoa(code),
		"error": errorMessage,
	})
	if _, err := w.Write(response); err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json")
	response, _ := json.Marshal(payload)
	if _, err := w.Write(response); err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}
}
