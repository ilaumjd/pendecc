package handlers

import (
	"encoding/json"
	"math/big"
	"net/http"
	"strings"

	"github.com/ilaumjd/pendecc/database"
)

type UrlHandler struct {
	Queries *database.Queries
}

func (h *UrlHandler) GetDefaultUrl(w http.ResponseWriter, r *http.Request) {
	shortUrl := r.PathValue("shortUrl")
	url, err := h.Queries.GetUrl(r.Context(), shortUrl)
	if err != nil {
		respondWithError(w, http.StatusNotFound)
		return
	}
	respondWithURL(w, url)
}

func (h *UrlHandler) CreateShortUrl(w http.ResponseWriter, r *http.Request) {

	type CreateShortUrlRequest struct {
		Url       string `json:"url"`
		CustomUrl string `json:"customUrl"`
	}

	params := CreateShortUrlRequest{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)

	if err != nil {
		respondWithError(w, http.StatusBadRequest)
		return
	}

	// reject if no protocol
	defaultUrl := params.Url
	if !strings.Contains(defaultUrl, "://") {
		respondWithError(w, http.StatusBadRequest, "URL must contains protocol (ex. http / https)")
		return
	}

	var shortUrl string

	customUrl := params.CustomUrl
	// if url custom
	if customUrl != "" {

		url, err := h.Queries.GetUrl(r.Context(), customUrl)
		if err != nil {
			shortUrl = customUrl
		} else if url.DefaultUrl == defaultUrl {
			respondWithURL(w, url)
			return
		} else {
			respondWithError(w, http.StatusInternalServerError, "Custom URL is already used")
			return
		}
	} else {

		currentString := defaultUrl
		for {
			// generate short url
			currentString = encodeBase62(currentString)[0:4]

			url, err := h.Queries.GetUrl(r.Context(), currentString)

			// break if not exists
			if err != nil {
				shortUrl = currentString
				break
			}

			// if exists
			if url.DefaultUrl == defaultUrl {
				respondWithURL(w, url)
				return
			}
		}
	}

	// create short url
	url, err := h.Queries.CreateUrl(r.Context(), database.CreateUrlParams{
		ShortUrl:   shortUrl,
		DefaultUrl: defaultUrl,
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create short url")
		return
	}

	respondWithURL(w, url)
}

func respondWithURL(w http.ResponseWriter, url database.Url) {
	respondWithJSON(w, http.StatusOK, map[string]string{
		"id":         url.ID.String(),
		"shortUrl":   url.ShortUrl,
		"defaultUrl": url.DefaultUrl,
	})
}

func encodeBase62(defaultString string) string {
	base62Chars := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	byteSlice := []byte(defaultString)
	bigInt := new(big.Int).SetBytes(byteSlice)
	var result []byte

	base := big.NewInt(62)
	zero := big.NewInt(0)
	mod := new(big.Int)

	for bigInt.Cmp(zero) > 0 {
		bigInt.DivMod(bigInt, base, mod) // Divide bigInt by 62, mod holds the remainder
		result = append([]byte{base62Chars[mod.Int64()]}, result...)
	}

	return string(result)
}
