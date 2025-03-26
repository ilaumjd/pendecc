package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ilaumjd/pendecc/database"
	"github.com/ilaumjd/pendecc/handlers"
	_ "github.com/lib/pq"
)

func main() {
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	sslMode := os.Getenv("DB_SSLMODE")

	if dbUser == "" || dbPass == "" || dbHost == "" || dbPort == "" || dbName == "" {
		log.Fatal("Database configuration environment variables are not set")
	}

	if sslMode == "" {
		sslMode = "disable"
	}

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		dbUser, dbPass, dbHost, dbPort, dbName, sslMode)
	conn, err := sql.Open("postgres", connStr)
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

	server := &http.Server{
		Addr:    ":5102",
		Handler: mux,
	}
	log.Fatal(server.ListenAndServe())
}
