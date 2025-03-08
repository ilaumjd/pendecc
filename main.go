package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/ilaumjd/pendecc/database"
	"github.com/ilaumjd/pendecc/handlers"
	"github.com/rs/cors"
	_ "github.com/lib/pq"
)

func main() {
	conn, err := sql.Open("postgres", "postgres://iam:postgres@localhost:5432/pendecc?sslmode=disable")
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

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},  // Be more specific in production
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"*"},
		AllowCredentials: true,
	})
	handler := c.Handler(mux)

	server := &http.Server{
		Addr:    ":5100",
		Handler: handler,
	}
	log.Fatal(server.ListenAndServe())
}
