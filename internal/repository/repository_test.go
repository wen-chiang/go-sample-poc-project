package repository

import (
	"testing"

	"github.com/example/library-management/internal/models"
)

func TestInMemoryBookRepository_Create(t *testing.T) {
	repo := NewInMemoryBookRepository()

	book := &models.Book{
		ID:          "test_1",
		Title:       "Test Book",
		Author:      "Test Author",
		ISBN:        "123-456-789",
		TotalCopies: 5,
	}

	err := repo.Create(book)
	if err != nil {
		t.Fatalf("Failed to create book: %v", err)
	}

	retrieved, err := repo.GetByID("test_1")
	if err != nil {
		t.Fatalf("Failed to retrieve book: %v", err)
	}

	if retrieved.Title != book.Title {
		t.Errorf("Expected title %s, got %s", book.Title, retrieved.Title)
	}

	if retrieved.AvailableCopies != book.TotalCopies {
		t.Errorf("Expected available copies %d, got %d", book.TotalCopies, retrieved.AvailableCopies)
	}
}

func TestInMemoryBookRepository_ManageCopies(t *testing.T) {
	repo := NewInMemoryBookRepository()

	book := &models.Book{
		ID:          "test_1",
		Title:       "Test Book",
		Author:      "Test Author",
		ISBN:        "123-456-789",
		TotalCopies: 5,
	}

	repo.Create(book)

	// Add copies
	err := repo.ManageCopies("test_1", 3)
	if err != nil {
		t.Fatalf("Failed to manage copies: %v", err)
	}

	available, _ := repo.CheckAvailability("test_1")
	if available != 8 {
		t.Errorf("Expected 8 available copies, got %d", available)
	}

	// Remove copies
	err = repo.ManageCopies("test_1", -3)
	if err != nil {
		t.Fatalf("Failed to manage copies: %v", err)
	}

	available, _ = repo.CheckAvailability("test_1")
	if available != 5 {
		t.Errorf("Expected 5 available copies, got %d", available)
	}
}

func TestInMemoryBookRepository_Delete(t *testing.T) {
	repo := NewInMemoryBookRepository()

	book := &models.Book{
		ID:          "test_1",
		Title:       "Test Book",
		Author:      "Test Author",
		ISBN:        "123-456-789",
		TotalCopies: 5,
	}

	repo.Create(book)

	err := repo.Delete("test_1")
	if err != nil {
		t.Fatalf("Failed to delete book: %v", err)
	}

	retrieved, _ := repo.GetByID("test_1")
	if !retrieved.IsArchived {
		t.Error("Expected book to be archived")
	}
}
