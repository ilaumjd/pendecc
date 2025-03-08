package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/ilaumjd/pendecc/database"
	"github.com/ilaumjd/pendecc/handlers"
	_ "github.com/lib/pq"
	"github.com/rs/cors"
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
		AllowedOrigins:   []string{"*"}, // Be more specific in production
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})
	handler := c.Handler(loggingMiddleware(mux))

	server := &http.Server{
		Addr:    ":5100",
		Handler: handler,
	}
	log.Fatal(server.ListenAndServe())
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Log the request
		logRequest(r)

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}

func logRequest(r *http.Request) {
	logEntry := fmt.Sprintf("[%s] %s %s %s\n",
		time.Now().Format("2006-01-02 15:04:05"),
		r.RemoteAddr,
		r.Method,
		r.URL.Path,
	)

	// Open the log file in append mode, create it if it doesn't exist
	f, err := os.OpenFile(os.ExpandEnv("$HOME/api_log.txt"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("Error opening log file: %v", err)
		return
	}
	defer f.Close()

	// Write the log entry to the file
	if _, err := f.WriteString(logEntry); err != nil {
		log.Printf("Error writing to log file: %v", err)
	}
}
