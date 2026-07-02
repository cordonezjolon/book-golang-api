package main

import (
	"api-golang/internal/service"
	"api-golang/internal/store"
	"api-golang/internal/transport"
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {

	// Database connection
	// read from environment variables
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL environment variable is required")
	}

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	// Initialize database schema
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS books (
		id SERIAL PRIMARY KEY,
		title TEXT NOT NULL,
		author TEXT NOT NULL
	);
	`)
	if err != nil {
		log.Fatal(err)
	}

	// Inyect dependencies
	bookStore := store.NewStore(db)
	bookService := service.NewService(bookStore)
	bookHandler := transport.NewBookHandler(bookService)

	// Start HTTP server
	mux := http.NewServeMux()

	mux.HandleFunc("/books", bookHandler.HandleBooks)
	mux.HandleFunc("/books/", bookHandler.HandleBookByID)

	log.Println("Server is running on port 8080...")
	log.Println("API endpoints:")
	log.Println("GET /books - List all books")
	log.Println("POST /books - Create a new book")
	log.Println("GET /books/{id} - Get a book by ID")
	log.Println("PUT /books/{id} - Update a book by ID")
	log.Println("DELETE /books/{id} - Delete a book by ID")

	srv := &http.Server{
		Addr:              ":8080",
		Handler:           mux,
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       120 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}

}
