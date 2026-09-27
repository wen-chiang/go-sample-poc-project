# API Documentation

## Base URL
```
http://localhost:8080
```

## Authentication
Currently uses basic credentials in POST /api/v1/librarians/login.
Future implementation will use JWT tokens.

## Response Format

### Success Response
```json
{
  "field1": "value1",
  "field2": "value2"
}
```

### Error Response
```json
{
  "error": "Error message describing what went wrong"
}
```

---

## Endpoints

### Health Check

#### GET /health
Check if the server is running.

**Response (200 OK):**
```json
{
  "status": "ok"
}
```

---

### Books

#### POST /api/v1/books
Create a new book in the library.

**Request Body:**
```json
{
  "title": "The Go Programming Language",
  "author": "Alan Donovan, Brian Kernighan",
  "isbn": "978-0134190440",
  "description": "The definitive guide to Go programming language",
  "published_at": "2015-10-26T00:00:00Z",
  "total_copies": 5
}
```

**Response (201 Created):**
```json
{
  "id": "book_20240927120000",
  "title": "The Go Programming Language",
  "author": "Alan Donovan, Brian Kernighan",
  "isbn": "978-0134190440",
  "description": "The definitive guide to Go programming language",
  "published_at": "2015-10-26T00:00:00Z",
  "total_copies": 5,
  "available_copies": 5,
  "is_archived": false,
  "created_at": "2024-09-27T12:00:00Z",
  "updated_at": "2024-09-27T12:00:00Z"
}
```

---

#### GET /api/v1/books
Retrieve all books in the library.

**Response (200 OK):**
```json
[
  {
    "id": "book_20240927120000",
    "title": "The Go Programming Language",
    "author": "Alan Donovan, Brian Kernighan",
    "isbn": "978-0134190440",
    "description": "The definitive guide to Go programming language",
    "published_at": "2015-10-26T00:00:00Z",
    "total_copies": 5,
    "available_copies": 5,
    "is_archived": false,
    "created_at": "2024-09-27T12:00:00Z",
    "updated_at": "2024-09-27T12:00:00Z"
  }
]
```

---

#### GET /api/v1/books/{id}
Get details of a specific book.

**Path Parameters:**
- `id` (string, required) - Book ID

**Response (200 OK):**
```json
{
  "id": "book_20240927120000",
  "title": "The Go Programming Language",
  "author": "Alan Donovan, Brian Kernighan",
  "isbn": "978-0134190440",
  "description": "The definitive guide to Go programming language",
  "published_at": "2015-10-26T00:00:00Z",
  "total_copies": 5,
  "available_copies": 5,
  "is_archived": false,
  "created_at": "2024-09-27T12:00:00Z",
  "updated_at": "2024-09-27T12:00:00Z"
}
```

**Error Response (404 Not Found):**
```json
{
  "error": "Book not found"
}
```

---

#### PUT /api/v1/books/{id}
Update book information.

**Path Parameters:**
- `id` (string, required) - Book ID

**Request Body:**
```json
{
  "title": "The Go Programming Language (2nd Edition)",
  "author": "Alan Donovan, Brian Kernighan",
  "isbn": "978-0134190440",
  "description": "Updated description",
  "published_at": "2015-10-26T00:00:00Z",
  "total_copies": 5
}
```

**Response (200 OK):**
```json
{
  "id": "book_20240927120000",
  "title": "The Go Programming Language (2nd Edition)",
  "author": "Alan Donovan, Brian Kernighan",
  "isbn": "978-0134190440",
  "description": "Updated description",
  "published_at": "2015-10-26T00:00:00Z",
  "total_copies": 5,
  "available_copies": 5,
  "is_archived": false,
  "created_at": "2024-09-27T12:00:00Z",
  "updated_at": "2024-09-27T12:00:05Z"
}
```

---

#### DELETE /api/v1/books/{id}
Archive (soft delete) a book.

**Path Parameters:**
- `id` (string, required) - Book ID

**Response (200 OK):**
```json
{
  "message": "Book archived successfully"
}
```

**Error Response (404 Not Found):**
```json
{
  "error": "Book not found"
}
```

---

#### POST /api/v1/books/{id}/copies
Manage book copies (add or remove stock).

**Path Parameters:**
- `id` (string, required) - Book ID

**Request Body:**
```json
{
  "amount": 3
}
```

Note: Use negative amount to decrease copies:
```json
{
  "amount": -2
}
```

**Response (200 OK):**
```json
{
  "id": "book_20240927120000",
  "title": "The Go Programming Language",
  "author": "Alan Donovan, Brian Kernighan",
  "isbn": "978-0134190440",
  "description": "The definitive guide to Go programming language",
  "published_at": "2015-10-26T00:00:00Z",
  "total_copies": 5,
  "available_copies": 8,
  "is_archived": false,
  "created_at": "2024-09-27T12:00:00Z",
  "updated_at": "2024-09-27T12:00:10Z"
}
```

**Error Response (400 Bad Request):**
```json
{
  "error": "insufficient available copies"
}
```

---

#### GET /api/v1/books/{id}/availability
Check the availability of a specific book.

**Path Parameters:**
- `id` (string, required) - Book ID

**Response (200 OK):**
```json
{
  "available_copies": 5
}
```

**Error Response (404 Not Found):**
```json
{
  "error": "Book not found"
}
```

---

### Reports

#### GET /api/v1/reports/summary
Get a summary report of all books in the library.

**Response (200 OK):**
```json
{
  "total_books": 10,
  "total_copies": 50,
  "available_copies": 42,
  "archived_books": 1,
  "generated_at": "2024-09-27T12:00:15Z"
}
```

---

### Authentication

#### POST /api/v1/librarians/login
Authenticate a librarian/staff member.

**Request Body:**
```json
{
  "username": "admin",
  "password": "password123"
}
```

**Response (200 OK):**
```json
{
  "id": "librarian_001",
  "username": "admin",
  "email": "admin@library.com",
  "role": "admin",
  "token": "token_20240927120000"
}
```

**Error Response (401 Unauthorized):**
```json
{
  "error": "Invalid credentials"
}
```

---

## HTTP Status Codes

| Code | Meaning |
|------|---------|
| 200 | OK - Request successful |
| 201 | Created - Resource created successfully |
| 400 | Bad Request - Invalid request data |
| 401 | Unauthorized - Authentication failed |
| 404 | Not Found - Resource not found |
| 500 | Internal Server Error - Server error |

---

## Rate Limiting

Not yet implemented. Future version will include:
- Rate limiting per IP
- Rate limiting per authenticated user
- Configurable limits

---

## Pagination

Not yet implemented. Future versions will support:
- Offset/limit pagination
- Cursor-based pagination

---

## Filtering & Sorting

Not yet implemented. Future versions will support:
- Filtering by title, author, ISBN
- Sorting by various fields
- Full-text search
