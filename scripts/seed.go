package main

import (
	"log"

	"github.com/example/library-management/internal/database"
	"github.com/example/library-management/internal/models"
	"github.com/example/library-management/internal/repository"
)

// This script seeds the database with test data
// Run with: go run scripts/seed.go

func main() {
	// Initialize database
	dbPath := "./library.db"
	db, err := database.InitDB(dbPath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.CloseDB(db)

	// Initialize repositories
	bookRepo := repository.NewSQLiteBookRepository(db)
	librarianRepo := repository.NewSQLiteLibrarianRepository(db)

	// Seed books
	seedBooks(bookRepo)

	// Seed librarians
	seedLibrarians(librarianRepo)

	log.Println("✅ Database seeded successfully!")
}

func seedBooks(repo repository.BookRepository) {
	books := []*models.Book{
		{
			ID:          "book_20240927120000",
			Title:       "The Go Programming Language",
			Author:      "Alan Donovan, Brian Kernighan",
			ISBN:        "978-0134190440",
			Description: "The definitive guide to the Go programming language",
			TotalCopies: 5,
		},
		{
			ID:          "book_20240927120001",
			Title:       "Go in Action",
			Author:      "William Kennedy, Brian Ketelsen, Erik St. Martin",
			ISBN:        "978-1617291784",
			Description: "Practical patterns for building real-world applications in Go",
			TotalCopies: 3,
		},
		{
			ID:          "book_20240927120002",
			Title:       "Clean Code",
			Author:      "Robert C. Martin",
			ISBN:        "978-0132350884",
			Description: "A Handbook of Agile Software Craftsmanship",
			TotalCopies: 4,
		},
		{
			ID:          "book_20240927120003",
			Title:       "Design Patterns",
			Author:      "Gang of Four",
			ISBN:        "978-0201633610",
			Description: "Elements of Reusable Object-Oriented Software",
			TotalCopies: 2,
		},
		{
			ID:          "book_20240927120004",
			Title:       "The Pragmatic Programmer",
			Author:      "Andrew Hunt, David Thomas",
			ISBN:        "978-0201616224",
			Description: "Your Journey to Mastery in Software Development",
			TotalCopies: 6,
		},
	}

	for _, book := range books {
		if err := repo.Create(book); err != nil {
			log.Printf("⚠️  Failed to create book %s: %v", book.Title, err)
		} else {
			log.Printf("✅ Created book: %s", book.Title)
		}
	}
}

func seedLibrarians(repo repository.LibrarianRepository) {
	librarians := []*models.Librarian{
		{
			ID:       "librarian_001",
			Username: "admin",
			Password: "admin123", // ⚠️  NEVER do this in production!
			Email:    "admin@library.com",
			Role:     "admin",
		},
		{
			ID:       "librarian_002",
			Username: "john_doe",
			Password: "password123",
			Email:    "john@library.com",
			Role:     "librarian",
		},
		{
			ID:       "librarian_003",
			Username: "jane_smith",
			Password: "password123",
			Email:    "jane@library.com",
			Role:     "librarian",
		},
	}

	for _, librarian := range librarians {
		if err := repo.Create(librarian); err != nil {
			log.Printf("⚠️  Failed to create librarian %s: %v", librarian.Username, err)
		} else {
			log.Printf("✅ Created librarian: %s (%s)", librarian.Username, librarian.Role)
		}
	}
}
