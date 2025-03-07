package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/ilaumjd/pendecc/database"
	"github.com/ilaumjd/pendecc/handlers"
	_ "github.com/lib/pq"
)

func main() {
	conn, err := sql.Open("postgres", "postgres://iam:@localhost:5432/pendecc?sslmode=disable")
	if err != nil {
		log.Fatal(err)
		return
	}
	queries := database.New(conn)
	urlHandler := handlers.UrlHandler{
		Queries: queries,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /urls/{shortUrl}", urlHandler.GetDefaultUrl)
	mux.HandleFunc("POST /urls", urlHandler.CreateShortUrl)
	mux.HandleFunc("POST /urls/{customUrl}", urlHandler.CreateShortUrl)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	log.Fatal(server.ListenAndServe())
}
