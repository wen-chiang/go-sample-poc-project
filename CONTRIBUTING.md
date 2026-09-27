# Contributing Guide

## Development Setup

1. **Install Go 1.21 or higher**
   - Download from https://golang.org/dl/

2. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd go-sample-poc-project
   ```

3. **Download dependencies**
   ```bash
   go mod download
   ```

4. **Build and run**
   ```bash
   make build-local
   make run
   ```

## Code Standards

### Naming Conventions
- **Package names**: lowercase, single word when possible
- **Function names**: CamelCase (exported functions start with uppercase)
- **Variable names**: camelCase (exported start with uppercase)
- **Constants**: UPPER_SNAKE_CASE

### Code Style
- Follow standard Go formatting: `go fmt ./...`
- Run linter: `golangci-lint run ./...`
- Keep functions small and focused
- Write comments for exported functions and types

### Error Handling
```go
// Good error handling
if err != nil {
    return fmt.Errorf("operation failed: %w", err)
}

// Avoid bare error returns
```

### Comments
```go
// Package models contains data structures
package models

// Book represents a book entity
type Book struct {
    ID string // Unique identifier
}

// GetBook retrieves a book by ID
func (r *Repository) GetBook(id string) (*Book, error) {
    // ...
}
```

## Project Structure Guidelines

- **cmd/** - Application entry points only
- **internal/** - Private business logic
  - **handler/** - HTTP handlers
  - **models/** - Data structures
  - **repository/** - Data access
  - Add subdirectories for larger features
- **test/** - Shared test utilities
- Do NOT use `pkg/` for public libraries in this project

## Adding New Features

### 1. Define the Model
Create types in `internal/models/` with appropriate JSON tags.

### 2. Create the Repository Interface
Define interface in `internal/repository/` for data access.

### 3. Implement Repository Methods
Add concrete implementation of repository interface.

### 4. Add HTTP Handler
Create handlers in `internal/handler/` with proper HTTP semantics.

### 5. Add Routes
Register routes in `cmd/server/main.go`.

### 6. Write Tests
Add unit tests in `*_test.go` files in the same package.

### 7. Update Documentation
- Update API.md with new endpoints
- Update README.md if needed

## Testing

### Run Tests
```bash
go test -v ./...
go test -v -race ./...  # With race detector
```

### Coverage
```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Table-Driven Tests
Preferred testing pattern:
```go
func TestFunction(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        want    string
        wantErr bool
    }{
        {
            name:    "valid input",
            input:   "test",
            want:    "result",
            wantErr: false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation
        })
    }
}
```

## Git Workflow

### Branch Naming
- `feature/description` - New features
- `bugfix/description` - Bug fixes
- `docs/description` - Documentation
- `refactor/description` - Refactoring

### Commit Messages
```
feat: Add book filtering functionality
fix: Correct availability calculation
docs: Update API documentation
refactor: Simplify repository interface
test: Add tests for book creation
```

### Pull Request Checklist
- [ ] Code follows style guidelines
- [ ] All tests pass
- [ ] New tests added for new features
- [ ] Documentation updated
- [ ] No breaking changes (unless documented)

## Performance Considerations

1. **Mutex Usage**: Current in-memory implementation uses RWMutex. Consider read/write ratios.
2. **Caching**: Future versions may need caching layer.
3. **Concurrency**: Safe for concurrent requests with current locking strategy.
4. **Database**: When migrating to database, consider connection pooling.

## Security Considerations

1. **Password Storage**: Currently plaintext - must use bcrypt in production
2. **JWT Tokens**: Implement proper JWT validation
3. **CORS**: Configure appropriately for your deployment
4. **Rate Limiting**: Should be added before production
5. **Input Validation**: Add comprehensive validation layer

## Documentation

### Code Documentation
- Every exported package should have a package comment
- Every exported type should have a comment
- Every exported function should have a comment
- Explain the "why", not just the "what"

### User Documentation
- Keep README.md updated
- Update API.md for endpoint changes
- Add examples for common operations

## Debugging

### Enable Debug Logging
Update middleware in `main.go` to add more detailed logging.

### Common Issues
- Port already in use: Check `lsof -i :8080` (macOS/Linux) or `netstat -ano | grep 8080` (Windows)
- Dependency issues: Run `go mod tidy` and `go mod download`

## Submitting Changes

1. Create a feature branch
2. Make your changes
3. Write/update tests
4. Format code: `go fmt ./...`
5. Run tests: `go test -v ./...`
6. Commit with clear message
7. Push and create pull request

## Questions?

Refer to:
- Official Go documentation: https://golang.org/doc/
- Chi router docs: https://github.com/go-chi/chi
- Project-specific documentation in this repo
