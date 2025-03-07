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
	respondWithJSON(w, http.StatusOK, map[string]string{
		"id":         url.ID.String(),
		"shortUrl":   url.ShortUrl,
		"defaultUrl": url.DefaultUrl,
	})
}

func (h *UrlHandler) CreateShortUrl(w http.ResponseWriter, r *http.Request) {

	type CreateShortUrlRequest struct {
		Url string `json:"url"`
	}

	params := CreateShortUrlRequest{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)

	if err != nil {
		respondWithError(w, http.StatusBadRequest)
		return
	}

	// remove http/https prefix
	defaultUrl := params.Url
	defaultUrl, httpPrefixFound := strings.CutPrefix(defaultUrl, "http://")
	defaultUrl, httpsPrefixFound := strings.CutPrefix(defaultUrl, "https://")

	if !httpPrefixFound && !httpsPrefixFound {
		respondWithError(w, http.StatusBadRequest)
		return
	}

	var shortUrl string

	customUrl := r.PathValue("customUrl")
	// if url custom
	if customUrl != "" {

		url, err := h.Queries.GetUrl(r.Context(), customUrl)
		if err != nil {
			shortUrl = customUrl
		} else if url.DefaultUrl == defaultUrl {
			respondWithJSON(w, http.StatusOK, map[string]string{
				"id":         url.ID.String(),
				"shortUrl":   url.ShortUrl,
				"defaultUrl": url.DefaultUrl,
			})
			return
		} else {
			respondWithError(w, http.StatusInternalServerError)
			return
		}
	} else {

		currentString := defaultUrl
		for {
			// generate short url
			currentString = encodeBase62(currentString)[0:7]

			url, err := h.Queries.GetUrl(r.Context(), currentString)

			// break if not exists
			if err != nil {
				shortUrl = currentString
				break
			}

			// if exists
			if url.DefaultUrl == defaultUrl {
				respondWithJSON(w, http.StatusOK, map[string]string{
					"id":         url.ID.String(),
					"shortUrl":   url.ShortUrl,
					"defaultUrl": url.DefaultUrl,
				})
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
		respondWithError(w, http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, http.StatusOK, url)
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
