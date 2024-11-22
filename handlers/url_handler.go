package handlers

import (
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"strings"

	"github.com/ilaumjd/pendecc/database"
)

type UrlHandler struct {
	Queries *database.Queries
}

func (h *UrlHandler) GetLongUrl(w http.ResponseWriter, r *http.Request) {
	shortURL := r.PathValue("shortURL")
	url, err := h.Queries.GetUrl(r.Context(), shortURL)
	if err != nil {
		respondWithError(w, http.StatusNotFound)
		return
	}
	defaultURL := fmt.Sprintf("https://%s", url.DefaultUrl)
	http.Redirect(w, r, defaultURL, http.StatusMovedPermanently)
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

	shortURL := generateShortUrl(defaultURL)
	url, err := h.Queries.CreateUrl(r.Context(), database.CreateUrlParams{
		ShortUrl:   shortURL,
		DefaultUrl: defaultURL,
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, http.StatusOK, url)
}

func generateShortUrl(defaultURL string) string {
	base62Chars := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	byteSlice := []byte(defaultURL)
	bigInt := new(big.Int).SetBytes(byteSlice)
	var result []byte

	base := big.NewInt(62)
	zero := big.NewInt(0)
	mod := new(big.Int)

	for bigInt.Cmp(zero) > 0 {
		bigInt.DivMod(bigInt, base, mod) // Divide bigInt by 62, mod holds the remainder
		result = append([]byte{base62Chars[mod.Int64()]}, result...)
	}

	return string(result)[0:7]
}
