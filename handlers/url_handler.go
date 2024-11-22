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
	shortUrl := r.PathValue("shortUrl")
	url, err := h.Queries.GetUrl(r.Context(), shortUrl)
	if err != nil {
		respondWithError(w, http.StatusNotFound)
		return
	}
	defaultUrl := fmt.Sprintf("https://%s", url.DefaultUrl)
	http.Redirect(w, r, defaultUrl, http.StatusMovedPermanently)
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

	defaultUrl := params.Url
	defaultUrl, httpPrefixFound := strings.CutPrefix(defaultUrl, "http://")
	defaultUrl, httpsPrefixFound := strings.CutPrefix(defaultUrl, "https://")

	if !httpPrefixFound && !httpsPrefixFound {
		respondWithError(w, http.StatusBadRequest)
		return
	}

	shortUrl := generateShortUrl(defaultUrl)
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

func generateShortUrl(defaultUrl string) string {
	base62Chars := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	byteSlice := []byte(defaultUrl)
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
