.PHONY: help build run test clean deps

help:
	@echo "Available commands:"
	@echo "  make build        - Build the application"
	@echo "  make run          - Run the application"
	@echo "  make test         - Run tests"
	@echo "  make clean        - Clean build artifacts"
	@echo "  make deps         - Download dependencies"
	@echo "  make fmt          - Format code"
	@echo "  make lint         - Lint code"

deps:
	go mod download
	go mod tidy

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o bin/server ./cmd/server

build-local:
	go build -v -o bin/server ./cmd/server

run: build-local
	./bin/server

test:
	go test -v -race -coverprofile=coverage.out ./...

test-coverage: test
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

clean:
	rm -rf bin/
	rm -f coverage.out coverage.html

fmt:
	go fmt ./...

lint:
	golangci-lint run ./...

docker-build:
	docker build -t library-management:latest .

docker-run:
	docker run -p 8080:8080 library-management:latest
