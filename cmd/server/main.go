package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/example/library-management/internal/database"
	"github.com/example/library-management/internal/handler"
	"github.com/example/library-management/internal/repository"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	// Initialize SQLite database
	dbPath := "./library.db"
	db, err := database.InitDB(dbPath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.CloseDB(db)

	// Initialize repositories with SQLite
	bookRepo := repository.NewSQLiteBookRepository(db)
	librarianRepo := repository.NewSQLiteLibrarianRepository(db)

	// Initialize handlers
	bookHandler := handler.NewBookHandler(bookRepo)
	librarianHandler := handler.NewLibrarianHandler(librarianRepo)

	// Setup router
	router := setupRouter(bookHandler, librarianHandler)

	// Create HTTP server
	srv := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Starting server on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}

func setupRouter(bookHandler *handler.BookHandler, librarianHandler *handler.LibrarianHandler) *chi.Mux {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(30 * time.Second))

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"status":"ok"}`)
	})

	// API routes
	r.Route("/api/v1", func(r chi.Router) {
		// Librarian routes
		r.Route("/librarians", func(r chi.Router) {
			r.Post("/login", librarianHandler.Login)
		})

		// Book routes
		r.Route("/books", func(r chi.Router) {
			r.Post("/", bookHandler.CreateBook)           // Add new book
			r.Get("/", bookHandler.ListBooks)             // Get all books
			r.Get("/{id}", bookHandler.GetBook)            // Get book details
			r.Put("/{id}", bookHandler.UpdateBook)         // Update book info
			r.Delete("/{id}", bookHandler.DeleteBook)      // Delete/archive book
			r.Post("/{id}/copies", bookHandler.ManageCopies) // Manage copies/stock
			r.Get("/{id}/availability", bookHandler.CheckAvailability) // Check availability
		})

		// Summary/Report routes
		r.Route("/reports", func(r chi.Router) {
			r.Get("/summary", bookHandler.SummaryReport) // Get summary report
		})
	})

	return r
}
