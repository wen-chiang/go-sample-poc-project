package repository

import (
	"database/sql"
	"errors"
	"time"

	"github.com/example/library-management/internal/models"
)

// SQLiteBookRepository is a SQLite implementation of BookRepository
type SQLiteBookRepository struct {
	db *sql.DB
}

// NewSQLiteBookRepository creates a new SQLite book repository
func NewSQLiteBookRepository(db *sql.DB) BookRepository {
	return &SQLiteBookRepository{db: db}
}

func (r *SQLiteBookRepository) Create(book *models.Book) error {
	stmt, err := r.db.Prepare(`
		INSERT INTO books (id, title, author, isbn, description, published_at, total_copies, available_copies, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	now := time.Now()
	_, err = stmt.Exec(
		book.ID,
		book.Title,
		book.Author,
		book.ISBN,
		book.Description,
		book.PublishedAt,
		book.TotalCopies,
		book.TotalCopies,
		now,
		now,
	)
	return err
}

func (r *SQLiteBookRepository) GetByID(id string) (*models.Book, error) {
	row := r.db.QueryRow(`
		SELECT id, title, author, isbn, description, published_at, total_copies, available_copies, is_archived, created_at, updated_at
		FROM books WHERE id = ?
	`, id)

	book := &models.Book{}
	err := row.Scan(
		&book.ID,
		&book.Title,
		&book.Author,
		&book.ISBN,
		&book.Description,
		&book.PublishedAt,
		&book.TotalCopies,
		&book.AvailableCopies,
		&book.IsArchived,
		&book.CreatedAt,
		&book.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("book not found")
		}
		return nil, err
	}

	return book, nil
}

func (r *SQLiteBookRepository) GetAll() ([]*models.Book, error) {
	rows, err := r.db.Query(`
		SELECT id, title, author, isbn, description, published_at, total_copies, available_copies, is_archived, created_at, updated_at
		FROM books ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []*models.Book
	for rows.Next() {
		book := &models.Book{}
		err := rows.Scan(
			&book.ID,
			&book.Title,
			&book.Author,
			&book.ISBN,
			&book.Description,
			&book.PublishedAt,
			&book.TotalCopies,
			&book.AvailableCopies,
			&book.IsArchived,
			&book.CreatedAt,
			&book.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		books = append(books, book)
	}

	return books, rows.Err()
}

func (r *SQLiteBookRepository) Update(book *models.Book) error {
	stmt, err := r.db.Prepare(`
		UPDATE books
		SET title = ?, author = ?, isbn = ?, description = ?, published_at = ?, total_copies = ?, available_copies = ?, is_archived = ?, updated_at = ?
		WHERE id = ?
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(
		book.Title,
		book.Author,
		book.ISBN,
		book.Description,
		book.PublishedAt,
		book.TotalCopies,
		book.AvailableCopies,
		book.IsArchived,
		time.Now(),
		book.ID,
	)

	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("book not found")
	}

	return nil
}

func (r *SQLiteBookRepository) Delete(id string) error {
	stmt, err := r.db.Prepare(`
		UPDATE books SET is_archived = 1, updated_at = ? WHERE id = ?
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(time.Now(), id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("book not found")
	}

	return nil
}

func (r *SQLiteBookRepository) CheckAvailability(id string) (int, error) {
	var available int
	err := r.db.QueryRow("SELECT available_copies FROM books WHERE id = ?", id).Scan(&available)

	if err != nil {
		if err == sql.ErrNoRows {
			return 0, errors.New("book not found")
		}
		return 0, err
	}

	return available, nil
}

func (r *SQLiteBookRepository) ManageCopies(id string, delta int) error {
	// Get current available copies
	var current int
	err := r.db.QueryRow("SELECT available_copies FROM books WHERE id = ?", id).Scan(&current)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("book not found")
		}
		return err
	}

	newAvailable := current + delta
	if newAvailable < 0 {
		return errors.New("insufficient available copies")
	}

	stmt, err := r.db.Prepare("UPDATE books SET available_copies = ?, updated_at = ? WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(newAvailable, time.Now(), id)
	return err
}

// SQLiteLibrarianRepository is a SQLite implementation of LibrarianRepository
type SQLiteLibrarianRepository struct {
	db *sql.DB
}

// NewSQLiteLibrarianRepository creates a new SQLite librarian repository
func NewSQLiteLibrarianRepository(db *sql.DB) LibrarianRepository {
	return &SQLiteLibrarianRepository{db: db}
}

func (r *SQLiteLibrarianRepository) Create(librarian *models.Librarian) error {
	stmt, err := r.db.Prepare(`
		INSERT INTO librarians (id, username, password, email, role, created_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		librarian.ID,
		librarian.Username,
		librarian.Password,
		librarian.Email,
		librarian.Role,
		time.Now(),
	)
	return err
}

func (r *SQLiteLibrarianRepository) GetByID(id string) (*models.Librarian, error) {
	row := r.db.QueryRow(`
		SELECT id, username, password, email, role, created_at FROM librarians WHERE id = ?
	`, id)

	librarian := &models.Librarian{}
	err := row.Scan(
		&librarian.ID,
		&librarian.Username,
		&librarian.Password,
		&librarian.Email,
		&librarian.Role,
		&librarian.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("librarian not found")
		}
		return nil, err
	}

	return librarian, nil
}

func (r *SQLiteLibrarianRepository) GetByUsername(username string) (*models.Librarian, error) {
	row := r.db.QueryRow(`
		SELECT id, username, password, email, role, created_at FROM librarians WHERE username = ?
	`, username)

	librarian := &models.Librarian{}
	err := row.Scan(
		&librarian.ID,
		&librarian.Username,
		&librarian.Password,
		&librarian.Email,
		&librarian.Role,
		&librarian.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("librarian not found")
		}
		return nil, err
	}

	return librarian, nil
}

func (r *SQLiteLibrarianRepository) GetAll() ([]*models.Librarian, error) {
	rows, err := r.db.Query(`
		SELECT id, username, password, email, role, created_at FROM librarians ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var librarians []*models.Librarian
	for rows.Next() {
		librarian := &models.Librarian{}
		err := rows.Scan(
			&librarian.ID,
			&librarian.Username,
			&librarian.Password,
			&librarian.Email,
			&librarian.Role,
			&librarian.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		librarians = append(librarians, librarian)
	}

	return librarians, rows.Err()
}

func (r *SQLiteLibrarianRepository) Update(librarian *models.Librarian) error {
	stmt, err := r.db.Prepare(`
		UPDATE librarians
		SET username = ?, password = ?, email = ?, role = ?
		WHERE id = ?
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(
		librarian.Username,
		librarian.Password,
		librarian.Email,
		librarian.Role,
		librarian.ID,
	)

	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("librarian not found")
	}

	return nil
}

func (r *SQLiteLibrarianRepository) Delete(id string) error {
	stmt, err := r.db.Prepare("DELETE FROM librarians WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("librarian not found")
	}

	return nil
}
