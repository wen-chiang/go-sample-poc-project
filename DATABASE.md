# Database Documentation

## Overview

The Library Management System uses **SQLite** as its embedded database. SQLite is a serverless, file-based SQL database perfect for POC projects and small-to-medium applications.

## Database File

- **Location**: `library.db` (created in the working directory)
- **Size**: Typically < 1MB for this POC
- **Backup**: Copy the `.db` file to backup
- **Restore**: Replace the current `.db` file with a backup

## Schema

### Books Table

```sql
CREATE TABLE books (
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
```

**Fields:**
- `id` - Unique book identifier (format: `book_YYYYMMDDHHMMSS`)
- `title` - Book title
- `author` - Book author name
- `isbn` - International Standard Book Number (unique)
- `description` - Book description
- `published_at` - Publication date
- `total_copies` - Total copies in inventory
- `available_copies` - Copies available for checkout
- `is_archived` - Soft delete flag (1 = archived, 0 = active)
- `created_at` - Record creation timestamp
- `updated_at` - Last update timestamp

### Librarians Table

```sql
CREATE TABLE librarians (
  id TEXT PRIMARY KEY,
  username TEXT NOT NULL UNIQUE,
  password TEXT NOT NULL,
  email TEXT UNIQUE,
  role TEXT DEFAULT 'librarian',
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

**Fields:**
- `id` - Unique librarian identifier
- `username` - Login username (unique)
- `password` - Login password (plaintext in POC, should be hashed in production)
- `email` - Email address
- `role` - User role (`admin` or `librarian`)
- `created_at` - Record creation timestamp

## Database Initialization

The database and tables are created automatically on the first server run:

```go
// From cmd/server/main.go
db, err := database.InitDB(dbPath)
```

The initialization process:
1. Opens/creates the `library.db` file
2. Runs migrations (creates tables if they don't exist)
3. Returns a ready-to-use database connection

## Accessing the Database Directly

### Using SQLite CLI

```bash
# Install sqlite3 (if not already installed)
# macOS: brew install sqlite3
# Ubuntu: sudo apt-get install sqlite3
# Windows: Download from https://www.sqlite.org/download.html

# Open the database
sqlite3 library.db

# View tables
.tables

# View books
SELECT id, title, author, total_copies, available_copies, is_archived 
FROM books;

# View librarians
SELECT id, username, email, role FROM librarians;

# Exit
.quit
```

### Using VS Code SQLite Extension

1. Install "SQLite" extension
2. Open Command Palette (Ctrl+Shift+P)
3. Search for "SQLite: Open Database"
4. Select `library.db`
5. Browse and query directly from the sidebar

### Programmatically

```go
package main

import (
    "database/sql"
    _ "github.com/mattn/go-sqlite3"
)

func main() {
    db, _ := sql.Open("sqlite3", "library.db")
    defer db.Close()
    
    rows, _ := db.Query("SELECT id, title FROM books")
    defer rows.Close()
    
    for rows.Next() {
        var id, title string
        rows.Scan(&id, &title)
        println(id, title)
    }
}
```

## Query Examples

### Get Active Books

```sql
SELECT * FROM books WHERE is_archived = 0 ORDER BY created_at DESC;
```

### Find Book by ISBN

```sql
SELECT * FROM books WHERE isbn = '978-0134190440';
```

### Count Available Copies

```sql
SELECT SUM(available_copies) as total_available FROM books WHERE is_archived = 0;
```

### Get Books by Author

```sql
SELECT * FROM books WHERE author LIKE '%Donovan%';
```

### Get All Admins

```sql
SELECT * FROM librarians WHERE role = 'admin';
```

## Data Persistence

### Automatic Backup

Create a simple backup script:

```bash
#!/bin/bash
# backup.sh
DATE=$(date +"%Y%m%d_%H%M%S")
cp library.db "backups/library_${DATE}.db"
echo "Database backed up to backups/library_${DATE}.db"
```

```bash
chmod +x backup.sh
./backup.sh
```

### Reset Database

```bash
# Delete the current database
rm library.db

# Or on Windows:
del library.db

# Restart the server - it will recreate with fresh schema
```

## Performance Considerations

### Current Setup
- In-memory implementations also available for testing
- SQLite is thread-safe with proper locking
- Suitable for concurrent read/write operations

### Optimization for Future Growth

If migrating to PostgreSQL:

```go
// Create a new PostgreSQL implementation
repo := repository.NewPostgresBookRepository(pgDB)
```

No other code changes needed due to interface-based design.

## Troubleshooting

### Database Locked Error
```
database is locked
```

**Cause**: Multiple processes accessing the database simultaneously.

**Solution**:
- Stop all running server instances
- Wait a few seconds
- Restart the server

### Corrupt Database

```bash
# Rebuild/repair the database
sqlite3 library.db "PRAGMA integrity_check;"

# Restore from backup if corrupted
cp backups/library_BACKUP.db library.db
```

### No Space on Disk

```bash
# Check database size
du -h library.db

# Vacuum to reclaim space
sqlite3 library.db "VACUUM;"
```

## Development vs Production

### Development (Current)
- Single file database
- Plain text passwords
- In-memory caching (optional)
- Perfect for POC

### Production Recommendations
- Use PostgreSQL or MySQL for better scalability
- Hash passwords with bcrypt
- Implement connection pooling
- Add database backups
- Enable WAL mode for better concurrency
- Add database indexes for performance

## Migration to PostgreSQL

When you're ready to scale:

1. **Create new repository**:
   ```go
   // internal/repository/postgres.go
   type PostgresBookRepository struct {
       db *sql.DB
   }
   ```

2. **Implement the interface**:
   ```go
   func (r *PostgresBookRepository) Create(book *models.Book) error {
       // PostgreSQL implementation
   }
   ```

3. **Update main.go**:
   ```go
   db := setupPostgresDB()
   bookRepo := repository.NewPostgresBookRepository(db)
   ```

That's it! The rest of the application remains unchanged.

## References

- SQLite Documentation: https://www.sqlite.org/docs.html
- Go sqlite3 Driver: https://github.com/mattn/go-sqlite3
- SQLite Best Practices: https://www.sqlite.org/bestpractice.html
