package main

import (
	"log"
	"net/http"

	"github.com/ilaumjd/pendecc/handlers"
)

func main() {
	urlHandler := handlers.UrlHandler{}

	mux := http.NewServeMux()
	mux.HandleFunc("POST /urls", urlHandler.CreateShortUrl)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	log.Fatal(server.ListenAndServe())
}
