# Library Management System - Go POC

A professional proof-of-concept library management system built with Go, following best practices and standard Go architecture patterns. Uses **SQLite** embedded database for data persistence.

## Features

### Librarian / Staff Functions
- ✅ Login as librarian/admin
- ✅ Add new book to catalog
- ✅ Update book information
- ✅ Delete or archive book
- ✅ Manage book copies / stock
- ✅ Check book availability
- ✅ Summary report with all book current information

### Database
- ✅ **SQLite embedded database** - File-based, no server needed
- ✅ Automatic schema creation via migrations
- ✅ Data persists across server restarts
- ✅ Easy to backup (single `.db` file)

## Project Structure

```
library-management/
├── cmd/
│   └── server/
│       └── main.go              # Application entry point
├── internal/
│   ├── database/
│   │   └── sqlite.go            # SQLite initialization & migrations
│   ├── handler/
│   │   └── handler.go           # HTTP request handlers
│   ├── models/
│   │   └── book.go              # Data models
│   └── repository/
│       ├── repository.go        # Interface definitions & in-memory impl
│       └── sqlite.go            # SQLite implementation
├── go.mod                       # Go module definition
├── go.sum                       # Dependency checksums
├── Dockerfile                   # Docker configuration
├── Makefile                     # Build automation
├── .env.example                 # Environment variables template
└── README.md                    # This file
```

## Architecture

### Key Design Patterns

1. **Package Organization**
   - `cmd/` - Application entry points
   - `internal/` - Private application code (not importable by external packages)
   - Clear separation of concerns

2. **Layered Architecture**
   - **Handlers** - HTTP request/response handling
   - **Repository** - Data access abstraction (swappable implementations)
   - **Models** - Domain objects
   - **Database** - Connection management and migrations

3. **Interface-Based Design**
   - Repositories defined as interfaces
   - Easy to swap implementations (SQLite ↔ PostgreSQL ↔ etc.)
   - Currently includes: SQLite & in-memory implementations
   - Facilitates testing with mocks

4. **Database Integration**
   - **SQLite** - Embedded, file-based, no server required
   - Automatic schema creation on startup
   - Thread-safe operations with proper locking
   - Easy backup (single `.db` file)

5. **Error Handling**
   - Explicit error returns
   - Proper HTTP status codes
   - User-friendly error messages

## API Endpoints

### Quick Reference

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/health` | Server health check |
| POST | `/api/v1/books` | Create new book |
| GET | `/api/v1/books` | List all books |
| GET | `/api/v1/books/{id}` | Get book details |
| PUT | `/api/v1/books/{id}` | Update book info |
| DELETE | `/api/v1/books/{id}` | Delete/archive book |
| POST | `/api/v1/books/{id}/copies` | Manage book copies/stock |
| GET | `/api/v1/books/{id}/availability` | Check book availability |
| GET | `/api/v1/reports/summary` | Get library summary report |
| POST | `/api/v1/librarians/login` | Librarian login |

**Base URL:** `http://localhost:8080`

For detailed endpoint documentation, see [API.md](API.md)

## Installation

### Prerequisites
- Go 1.21 or higher
- Make (optional, for convenience)

### Steps

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd go-sample-poc-project
   ```

2. **Download dependencies**
   ```bash
   make deps
   # or
   go mod download
   ```

3. **Build the project**
   ```bash
   make build-local
   # or
   go build -v -o bin/server ./cmd/server
   ```

4. **Run the server**
   ```bash
   make run
   # or
   ./bin/server
   ```

Server will start on `http://localhost:8080`

## API Endpoints

### Health Check
- `GET /health` - Server health status

### Books Management
- `POST /api/v1/books` - Create new book
- `GET /api/v1/books` - List all books
- `GET /api/v1/books/{id}` - Get book details
- `PUT /api/v1/books/{id}` - Update book information
- `DELETE /api/v1/books/{id}` - Archive book
- `POST /api/v1/books/{id}/copies` - Manage copies/stock
- `GET /api/v1/books/{id}/availability` - Check availability

### Reports
- `GET /api/v1/reports/summary` - Get summary report

### Authentication
- `POST /api/v1/librarians/login` - Librarian login

## API Examples

### Create a Book
```bash
curl -X POST http://localhost:8080/api/v1/books \
  -H "Content-Type: application/json" \
  -d '{
    "title": "The Go Programming Language",
    "author": "Alan Donovan, Brian Kernighan",
    "isbn": "978-0134190440",
    "description": "The definitive guide to Go",
    "published_at": "2015-10-26T00:00:00Z",
    "total_copies": 5
  }'
```

### List All Books
```bash
curl http://localhost:8080/api/v1/books
```

### Get Book Details
```bash
curl http://localhost:8080/api/v1/books/{id}
```

### Check Availability
```bash
curl http://localhost:8080/api/v1/books/{id}/availability
```

### Manage Copies
```bash
curl -X POST http://localhost:8080/api/v1/books/{id}/copies \
  -H "Content-Type: application/json" \
  -d '{"amount": 2}'
```

### Get Summary Report
```bash
curl http://localhost:8080/api/v1/reports/summary
```

## Development

### Format Code
```bash
make fmt
# or
go fmt ./...
```

### Run Tests
```bash
make test
# or
go test -v -race -coverprofile=coverage.out ./...
```

### View Coverage
```bash
make test-coverage
```

### Lint Code (requires golangci-lint)
```bash
make lint
```

### Clean Build Artifacts
```bash
make clean
```

## Database

### SQLite

The project uses **SQLite** for data persistence:

- **File**: `library.db` (created on first run)
- **Type**: Embedded, serverless SQL database
- **Advantages**: 
  - No external database server needed
  - Perfect for POC/testing
  - Single file for easy backup
  - Data persists across server restarts
- **Backup**: Simply copy the `library.db` file

### Schema

Automatically created on startup:

**books table:**
```sql
CREATE TABLE books (
  id TEXT PRIMARY KEY,
  title TEXT NOT NULL,
  author TEXT NOT NULL,
  isbn TEXT UNIQUE,
  description TEXT,
  published_at DATETIME,
  total_copies INTEGER,
  available_copies INTEGER,
  is_archived BOOLEAN,
  created_at DATETIME,
  updated_at DATETIME
);
```

**librarians table:**
```sql
CREATE TABLE librarians (
  id TEXT PRIMARY KEY,
  username TEXT UNIQUE,
  password TEXT,
  email TEXT UNIQUE,
  role TEXT,
  created_at DATETIME
);
```

### Switching to Another Database

To use a different database (PostgreSQL, MySQL, etc.):

1. Create a new repository implementation (e.g., `internal/repository/postgres.go`)
2. Implement the `BookRepository` and `LibrarianRepository` interfaces
3. Update `cmd/server/main.go` to use the new repository:

```go
// Instead of:
bookRepo := repository.NewSQLiteBookRepository(db)

// Use:
bookRepo := repository.NewPostgresBookRepository(db)
```

The handlers and rest of the code remain unchanged!

## Docker Support

### Build Docker Image
```bash
make docker-build
```

### Run in Docker
```bash
make docker-run
```

Or manually:
```bash
docker build -t library-management:latest .
docker run -p 8080:8080 library-management:latest
```

## Best Practices Implemented

1. ✅ **Package Organization** - Clear separation with `cmd/`, `internal/`
2. ✅ **Error Handling** - Explicit error returns, proper HTTP status codes
3. ✅ **Concurrency** - Mutex-protected in-memory data structures
4. ✅ **Middleware** - Request ID, logging, recovery, timeouts
5. ✅ **Graceful Shutdown** - Proper server cleanup on SIGINT/SIGTERM
6. ✅ **Code Structure** - Interface-based design, dependency injection
7. ✅ **Configuration** - Environment-based setup
8. ✅ **Documentation** - Clear README and code comments
9. ✅ **Containerization** - Multi-stage Docker build
10. ✅ **Build Automation** - Makefile for common tasks

## Dependencies

- `github.com/go-chi/chi/v5` - HTTP router
- `github.com/go-chi/cors` - CORS middleware
- `github.com/mattn/go-sqlite3` - SQLite driver

## Future Enhancements

- [ ] Database integration (PostgreSQL, MySQL)
- [ ] JWT authentication tokens
- [ ] Password hashing with bcrypt
- [ ] User management system
- [ ] Book checkout/return functionality
- [ ] Fine calculation system
- [ ] Comprehensive unit tests
- [ ] Integration tests
- [ ] API documentation (Swagger/OpenAPI)
- [ ] Structured JSON logging
- [ ] Metrics and monitoring
- [ ] Rate limiting
- [ ] Advanced request validation

## License

MIT

## Author

Library Management POC Team
