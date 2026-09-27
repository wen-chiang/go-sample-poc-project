# Quick Start Guide

Get the Library Management System running in minutes.

## Prerequisites
- Go 1.21+ ([Download](https://golang.org/dl/))
- Make (optional, for convenience)
- curl (for testing API)

## Option 1: Using Make (Recommended)

```bash
# 1. Download dependencies
make deps

# 2. Build the application
make build-local

# 3. Run the server
make run
```

Server runs on `http://localhost:8080`
Database file (`library.db`) is created automatically on first run.

## Option 2: Using Go Commands

```bash
# 1. Download dependencies
go mod download

# 2. Build
go build -v -o bin/server ./cmd/server

# 3. Run
./bin/server
```

## Option 3: Direct Run

```bash
go run ./cmd/server/main.go
```

## Verify Server is Running

```bash
curl http://localhost:8080/health
```

Expected response:
```json
{"status":"ok"}
```

## First API Call: Create a Book

```bash
curl -X POST http://localhost:8080/api/v1/books \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Learn Go",
    "author": "John Doe",
    "isbn": "123-456-789",
    "total_copies": 3
  }'
```

## Get All Books

```bash
curl http://localhost:8080/api/v1/books
```

## Summary Report

```bash
curl http://localhost:8080/api/v1/reports/summary
```

## Stop the Server

Press `Ctrl+C` in the terminal

## Database Information

- **Type**: SQLite (embedded, file-based)
- **File**: `library.db` (created in the working directory)
- **Data**: Persists across server restarts
- **Backup**: Simply copy the `library.db` file

### Reset Database

To start fresh, delete the `library.db` file:

```bash
rm library.db     # macOS/Linux
del library.db    # Windows
```

The database will be recreated with the schema on the next server start.

## Seed Test Data

Populate the database with sample books and librarians:

```bash
go run ./scripts/seed.go
```

This creates:
- 5 sample books
- 3 librarian accounts (admin, john_doe, jane_smith)

Then test the API:

```bash
# Get all books (should return 5)
curl http://localhost:8080/api/v1/books

# Login
curl -X POST http://localhost:8080/api/v1/librarians/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}'
```

## Next Steps

1. Read [API.md](API.md) for complete API documentation
2. Check [README.md](README.md) for architecture details
3. Review [CONTRIBUTING.md](CONTRIBUTING.md) for development guidelines

## Troubleshooting

### "Port 8080 already in use"
Change the port in `cmd/server/main.go`:
```go
Addr: ":9000",  // Change 8080 to another port
```

### "Module not found"
```bash
go mod tidy
go mod download
```

### "Build fails"
Ensure you're using Go 1.21+:
```bash
go version
```

### "Database locked error"
Multiple servers are trying to access the same `library.db` file. Stop other instances and try again.

### "Need help with API?"
Test endpoints using the curl examples above or use Postman:
- Import API calls from examples in [API.md](API.md)
- Set base URL to `http://localhost:8080`

## Common Commands Reference

| Command | Purpose |
|---------|---------|
| `make help` | Show all available commands |
| `make build-local` | Build for your OS |
| `make run` | Build and run server |
| `make test` | Run all tests |
| `make clean` | Remove build artifacts |
| `make fmt` | Format code |

## Development Tips

1. **Hot Reload**: Use tools like `air` or `CompileDaemon` for auto-recompile
   ```bash
   go install github.com/cosmtrek/air@latest
   air
   ```

2. **Debug with Prints**: Add logs to understand flow
   ```go
   log.Printf("Debug info: %v", variable)
   ```

3. **Test Endpoints**: Use curl, Postman, or VS Code REST Client
   
4. **Check Concurrency**: Run stress tests
   ```bash
   for i in {1..100}; do curl http://localhost:8080/health &; done
   wait
   ```

## What's Implemented

✅ Full CRUD for books
✅ Stock/copies management
✅ Availability checking
✅ Summary reports
✅ Librarian login (basic)
✅ **SQLite database persistence**
✅ Proper error handling
✅ Graceful shutdown
✅ Thread-safe operations
✅ Comprehensive logging

## What's Coming

- PostgreSQL/MySQL support
- JWT authentication
- Advanced search/filtering
- More comprehensive tests
- API documentation (Swagger)

Enjoy! 🎉
