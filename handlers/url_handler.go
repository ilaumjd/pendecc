package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/ilaumjd/pendecc/database"
)

type UrlHandler struct {
	Queries *database.Queries
}

func (h *UrlHandler) GetLongUrl(w http.ResponseWriter, r *http.Request) {
}

func (h *UrlHandler) CreateShortUrl(w http.ResponseWriter, r *http.Request) {

	type CreateShortUrlRequest struct {
		URL string `json:"url"`
	}

	params := CreateShortUrlRequest{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)

	if err != nil {
		respondWithError(w, http.StatusBadRequest)
		return
	}

	defaultURL := params.URL
	defaultURL, httpPrefixFound := strings.CutPrefix(defaultURL, "http://")
	defaultURL, httpsPrefixFound := strings.CutPrefix(defaultURL, "https://")

	if !httpPrefixFound && !httpsPrefixFound {
		respondWithError(w, http.StatusBadRequest)
		return
	}

	url, err := h.Queries.CreateUrl(r.Context(), database.CreateUrlParams{
		ShortUrl:   "ooo", // TODO: Generate short url
		DefaultUrl: defaultURL,
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, http.StatusOK, url)
}
