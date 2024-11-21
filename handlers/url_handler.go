package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
)

type UrlHandler struct{}

func (h *UrlHandler) GetLongUrl(w http.ResponseWriter, r *http.Request) {
}

func (h *UrlHandler) CreateShortUrl(w http.ResponseWriter, r *http.Request) {

	type CreateShortUrlRequest struct {
		LongURL string `json:"long_url"`
	}

	params := CreateShortUrlRequest{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)

	if err != nil {
		respondWithError(w, http.StatusBadRequest)
		return
	}

	longURL := params.LongURL
	longURL, httpPrefixFound := strings.CutPrefix(longURL, "http://")
	longURL, httpsPrefixFound := strings.CutPrefix(longURL, "https://")

	if !httpPrefixFound && !httpsPrefixFound {
		respondWithError(w, http.StatusBadRequest)
		return
	}

	// TODO: handle valid link

	respondWithJSON(w, http.StatusOK, longURL)
}
