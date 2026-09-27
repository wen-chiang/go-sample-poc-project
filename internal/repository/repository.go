package repository

import (
	"errors"
	"sync"
	"time"

	"github.com/example/library-management/internal/models"
)

// BookRepository defines the interface for book data operations
type BookRepository interface {
	Create(book *models.Book) error
	GetByID(id string) (*models.Book, error)
	GetAll() ([]*models.Book, error)
	Update(book *models.Book) error
	Delete(id string) error
	CheckAvailability(id string) (int, error)
	ManageCopies(id string, delta int) error
}

// LibrarianRepository defines the interface for librarian data operations
type LibrarianRepository interface {
	Create(librarian *models.Librarian) error
	GetByID(id string) (*models.Librarian, error)
	GetByUsername(username string) (*models.Librarian, error)
	GetAll() ([]*models.Librarian, error)
	Update(librarian *models.Librarian) error
	Delete(id string) error
}

// InMemoryBookRepository is an in-memory implementation of BookRepository
type InMemoryBookRepository struct {
	mu    sync.RWMutex
	books map[string]*models.Book
}

// NewInMemoryBookRepository creates a new in-memory book repository
func NewInMemoryBookRepository() BookRepository {
	return &InMemoryBookRepository{
		books: make(map[string]*models.Book),
	}
}

func (r *InMemoryBookRepository) Create(book *models.Book) error {
	if book.ID == "" {
		return errors.New("book ID is required")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.books[book.ID]; exists {
		return errors.New("book already exists")
	}

	book.CreatedAt = time.Now()
	book.UpdatedAt = time.Now()
	book.AvailableCopies = book.TotalCopies
	r.books[book.ID] = book
	return nil
}

func (r *InMemoryBookRepository) GetByID(id string) (*models.Book, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	book, exists := r.books[id]
	if !exists {
		return nil, errors.New("book not found")
	}
	return book, nil
}

func (r *InMemoryBookRepository) GetAll() ([]*models.Book, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	books := make([]*models.Book, 0, len(r.books))
	for _, book := range r.books {
		books = append(books, book)
	}
	return books, nil
}

func (r *InMemoryBookRepository) Update(book *models.Book) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.books[book.ID]; !exists {
		return errors.New("book not found")
	}

	book.UpdatedAt = time.Now()
	r.books[book.ID] = book
	return nil
}

func (r *InMemoryBookRepository) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if book, exists := r.books[id]; exists {
		book.IsArchived = true
		book.UpdatedAt = time.Now()
		r.books[id] = book
		return nil
	}
	return errors.New("book not found")
}

func (r *InMemoryBookRepository) CheckAvailability(id string) (int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	book, exists := r.books[id]
	if !exists {
		return 0, errors.New("book not found")
	}
	return book.AvailableCopies, nil
}

func (r *InMemoryBookRepository) ManageCopies(id string, delta int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	book, exists := r.books[id]
	if !exists {
		return errors.New("book not found")
	}

	newAvailable := book.AvailableCopies + delta
	if newAvailable < 0 {
		return errors.New("insufficient available copies")
	}

	book.AvailableCopies = newAvailable
	book.UpdatedAt = time.Now()
	r.books[id] = book
	return nil
}

// InMemoryLibrarianRepository is an in-memory implementation of LibrarianRepository
type InMemoryLibrarianRepository struct {
	mu         sync.RWMutex
	librarians map[string]*models.Librarian
}

// NewInMemoryLibrarianRepository creates a new in-memory librarian repository
func NewInMemoryLibrarianRepository() LibrarianRepository {
	return &InMemoryLibrarianRepository{
		librarians: make(map[string]*models.Librarian),
	}
}

func (r *InMemoryLibrarianRepository) Create(librarian *models.Librarian) error {
	if librarian.ID == "" {
		return errors.New("librarian ID is required")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.librarians[librarian.ID]; exists {
		return errors.New("librarian already exists")
	}

	librarian.CreatedAt = time.Now()
	r.librarians[librarian.ID] = librarian
	return nil
}

func (r *InMemoryLibrarianRepository) GetByID(id string) (*models.Librarian, error) {
	r.mu.RLock()
	defer r.mu.Unlock()

	librarian, exists := r.librarians[id]
	if !exists {
		return nil, errors.New("librarian not found")
	}
	return librarian, nil
}

func (r *InMemoryLibrarianRepository) GetByUsername(username string) (*models.Librarian, error) {
	r.mu.RLock()
	defer r.mu.Unlock()

	for _, librarian := range r.librarians {
		if librarian.Username == username {
			return librarian, nil
		}
	}
	return nil, errors.New("librarian not found")
}

func (r *InMemoryLibrarianRepository) GetAll() ([]*models.Librarian, error) {
	r.mu.RLock()
	defer r.mu.Unlock()

	librarians := make([]*models.Librarian, 0, len(r.librarians))
	for _, librarian := range r.librarians {
		librarians = append(librarians, librarian)
	}
	return librarians, nil
}

func (r *InMemoryLibrarianRepository) Update(librarian *models.Librarian) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.librarians[librarian.ID]; !exists {
		return errors.New("librarian not found")
	}

	r.librarians[librarian.ID] = librarian
	return nil
}

func (r *InMemoryLibrarianRepository) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.librarians[id]; !exists {
		return errors.New("librarian not found")
	}

	delete(r.librarians, id)
	return nil
}
