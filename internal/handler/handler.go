package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/example/library-management/internal/models"
	"github.com/example/library-management/internal/repository"
	"github.com/go-chi/chi/v5"
)

// BookHandler handles HTTP requests related to books
type BookHandler struct {
	repo repository.BookRepository
}

// NewBookHandler creates a new BookHandler
func NewBookHandler(repo repository.BookRepository) *BookHandler {
	return &BookHandler{repo: repo}
}

// CreateBook handles POST /api/v1/books
func (h *BookHandler) CreateBook(w http.ResponseWriter, r *http.Request) {
	var req models.BookRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	book := &models.Book{
		ID:          generateID(),
		Title:       req.Title,
		Author:      req.Author,
		ISBN:        req.ISBN,
		Description: req.Description,
		PublishedAt: req.PublishedAt,
		TotalCopies: req.TotalCopies,
	}

	if err := h.repo.Create(book); err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSONResponse(w, http.StatusCreated, book)
}

// GetBook handles GET /api/v1/books/{id}
func (h *BookHandler) GetBook(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	book, err := h.repo.GetByID(id)
	if err != nil {
		writeErrorResponse(w, http.StatusNotFound, "Book not found")
		return
	}

	writeJSONResponse(w, http.StatusOK, book)
}

// ListBooks handles GET /api/v1/books
func (h *BookHandler) ListBooks(w http.ResponseWriter, r *http.Request) {
	books, err := h.repo.GetAll()
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSONResponse(w, http.StatusOK, books)
}

// UpdateBook handles PUT /api/v1/books/{id}
func (h *BookHandler) UpdateBook(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	
	// Get existing book
	book, err := h.repo.GetByID(id)
	if err != nil {
		writeErrorResponse(w, http.StatusNotFound, "Book not found")
		return
	}

	// Parse request body
	var req models.BookRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Update fields
	book.Title = req.Title
	book.Author = req.Author
	book.ISBN = req.ISBN
	book.Description = req.Description
	book.PublishedAt = req.PublishedAt

	if err := h.repo.Update(book); err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSONResponse(w, http.StatusOK, book)
}

// DeleteBook handles DELETE /api/v1/books/{id}
func (h *BookHandler) DeleteBook(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.repo.Delete(id); err != nil {
		writeErrorResponse(w, http.StatusNotFound, "Book not found")
		return
	}

	writeJSONResponse(w, http.StatusOK, map[string]string{"message": "Book archived successfully"})
}

// ManageCopies handles POST /api/v1/books/{id}/copies
func (h *BookHandler) ManageCopies(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	
	var req models.CopyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.repo.ManageCopies(id, req.Amount); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	book, _ := h.repo.GetByID(id)
	writeJSONResponse(w, http.StatusOK, book)
}

// CheckAvailability handles GET /api/v1/books/{id}/availability
func (h *BookHandler) CheckAvailability(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	available, err := h.repo.CheckAvailability(id)
	if err != nil {
		writeErrorResponse(w, http.StatusNotFound, "Book not found")
		return
	}

	writeJSONResponse(w, http.StatusOK, map[string]int{"available_copies": available})
}

// SummaryReport handles GET /api/v1/reports/summary
func (h *BookHandler) SummaryReport(w http.ResponseWriter, r *http.Request) {
	books, err := h.repo.GetAll()
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	var totalBooks, totalCopies, availableCopies, archivedBooks int
	for _, book := range books {
		totalBooks++
		totalCopies += book.TotalCopies
		availableCopies += book.AvailableCopies
		if book.IsArchived {
			archivedBooks++
		}
	}

	report := &models.SummaryReport{
		TotalBooks:      totalBooks,
		TotalCopies:     totalCopies,
		AvailableCopies: availableCopies,
		ArchivedBooks:   archivedBooks,
		GeneratedAt:     time.Now(),
	}

	writeJSONResponse(w, http.StatusOK, report)
}

// LibrarianHandler handles HTTP requests related to librarians
type LibrarianHandler struct {
	repo repository.LibrarianRepository
}

// NewLibrarianHandler creates a new LibrarianHandler
func NewLibrarianHandler(repo repository.LibrarianRepository) *LibrarianHandler {
	return &LibrarianHandler{repo: repo}
}

// Login handles POST /api/v1/librarians/login
func (h *LibrarianHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Get librarian by username
	librarian, err := h.repo.GetByUsername(req.Username)
	if err != nil {
		writeErrorResponse(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	// Verify password (in production, use bcrypt)
	if librarian.Password != req.Password {
		writeErrorResponse(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	// Create response (in production, generate JWT token)
	response := &models.LoginResponse{
		ID:       librarian.ID,
		Username: librarian.Username,
		Email:    librarian.Email,
		Role:     librarian.Role,
		Token:    generateToken(),
	}

	writeJSONResponse(w, http.StatusOK, response)
}

// Helper functions
func writeJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func writeErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func generateID() string {
	return "book_" + time.Now().Format("20060102150405")
}

func generateToken() string {
	// In production, use JWT library
	return "token_" + time.Now().Format("20060102150405")
}
