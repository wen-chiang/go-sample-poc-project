package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "modernc.org/sqlite"
)

// InitDB initializes the SQLite database and runs migrations
func InitDB(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Test connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Run migrations
	if err := runMigrations(db); err != nil {
		return nil, fmt.Errorf("migration failed: %w", err)
	}

	log.Println("Database initialized successfully")
	return db, nil
}

// runMigrations creates the necessary tables
func runMigrations(db *sql.DB) error {
	migrations := []string{
		createBooksTable,
		createLibrariansTable,
	}

	for _, migration := range migrations {
		if _, err := db.Exec(migration); err != nil {
			return err
		}
	}

	return nil
}

const (
	createBooksTable = `
	CREATE TABLE IF NOT EXISTS books (
		id TEXT PRIMARY KEY,
		title TEXT NOT NULL,
		author TEXT NOT NULL,
		isbn TEXT NOT NULL UNIQUE,
		description TEXT,
		published_at DATETIME,
		total_copies INTEGER NOT NULL DEFAULT 0,
		available_copies INTEGER NOT NULL DEFAULT 0,
		is_archived BOOLEAN DEFAULT 0,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`

	createLibrariansTable = `
	CREATE TABLE IF NOT EXISTS librarians (
		id TEXT PRIMARY KEY,
		username TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL,
		email TEXT UNIQUE,
		role TEXT DEFAULT 'librarian',
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`
)

// CloseDB closes the database connection
func CloseDB(db *sql.DB) error {
	return db.Close()
}
