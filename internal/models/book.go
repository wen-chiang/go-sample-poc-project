package models

import "time"

// Book represents a book entity in the library
type Book struct {
	ID           string    `json:"id"`
	Title        string    `json:"title"`
	Author       string    `json:"author"`
	ISBN         string    `json:"isbn"`
	Description  string    `json:"description"`
	PublishedAt  time.Time `json:"published_at"`
	TotalCopies  int       `json:"total_copies"`
	AvailableCopies int    `json:"available_copies"`
	IsArchived   bool      `json:"is_archived"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// BookRequest represents the request payload for creating/updating a book
type BookRequest struct {
	Title       string    `json:"title" validate:"required"`
	Author      string    `json:"author" validate:"required"`
	ISBN        string    `json:"isbn" validate:"required"`
	Description string    `json:"description"`
	PublishedAt time.Time `json:"published_at"`
	TotalCopies int       `json:"total_copies" validate:"min=1"`
}

// CopyRequest represents the request for managing book copies
type CopyRequest struct {
	Amount int `json:"amount" validate:"required"`
}

// Librarian represents a librarian/staff member
type Librarian struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Password  string    `json:"-"` // Never expose password in JSON
	Email     string    `json:"email"`
	Role      string    `json:"role"` // "librarian" or "admin"
	CreatedAt time.Time `json:"created_at"`
}

// LoginRequest represents the login request payload
type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// LoginResponse represents the login response
type LoginResponse struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	Token    string `json:"token"`
}

// SummaryReport represents the library summary report
type SummaryReport struct {
	TotalBooks      int       `json:"total_books"`
	TotalCopies     int       `json:"total_copies"`
	AvailableCopies int       `json:"available_copies"`
	ArchivedBooks   int       `json:"archived_books"`
	GeneratedAt     time.Time `json:"generated_at"`
}
