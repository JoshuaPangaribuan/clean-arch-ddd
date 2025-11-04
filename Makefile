.PHONY: help deps run build test test-unit test-integration test-coverage test-watch test-short clean docker-up docker-down migrate-up migrate-down migrate-create sqlc-generate generate-mocks setup-tdd verify

# Default target
help:
	@echo "Available targets:"
	@echo "  deps            - Download and install dependencies"
	@echo "  run             - Run the application"
	@echo "  build           - Build the application"
	@echo "  test            - Run all tests with coverage"
	@echo "  test-unit       - Run unit tests only (fast, no external dependencies)"
	@echo "  test-integration - Run integration tests (requires database)"
	@echo "  test-coverage   - Run tests with HTML coverage report"
	@echo "  test-watch      - Run tests in watch mode (requires entr or similar)"
	@echo "  test-short      - Run tests without coverage (faster)"
	@echo "  clean           - Clean build artifacts"
	@echo "  docker-up       - Start Docker containers"
	@echo "  docker-down     - Stop Docker containers"
	@echo "  migrate-up      - Run database migrations"
	@echo "  migrate-down    - Rollback database migrations"
	@echo "  migrate-create  - Create a new migration (usage: make migrate-create name=migration_name)"
	@echo "  sqlc-generate   - Generate sqlc code"
	@echo "  generate-mocks  - Generate mocks for interfaces"
	@echo "  setup-tdd       - Setup TDD environment (deps + generate mocks + verify)"
	@echo "  verify          - Verify all tests pass"

# Download and install dependencies
deps:
	go mod download
	go mod tidy
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	go install github.com/pressly/goose/v3/cmd/goose@latest
	go install github.com/vektra/mockery/v2@latest

# Run the application
run:
	go run cmd/api/main.go

# Build the application
build:
	go build -o bin/api cmd/api/main.go

# Run all tests with coverage (default TDD command)
test: generate-mocks
	go test -v -cover ./...

# Run unit tests only (fast, no external dependencies)
test-unit: generate-mocks
	go test -v -cover -short ./internal/domain/... ./internal/application/...

# Run integration tests (requires database)
test-integration: docker-up migrate-up
	@echo "Waiting for database to be ready..."
	@sleep 2
	go test -v -cover ./internal/infrastructure/... ./cmd/...
	@$(MAKE) docker-down

# Run tests with coverage report
test-coverage: generate-mocks
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run tests in watch mode (requires entr or similar tool)
# On Windows, consider using: https://github.com/radovskyb/watcher
test-watch: generate-mocks
	@echo "Watching for changes... (Press Ctrl+C to stop)"
	@if command -v entr > /dev/null; then \
		find . -name "*.go" -not -path "./vendor/*" | entr -c go test -v -short ./...; \
	else \
		echo "Error: entr not found. Install it or use: go test -v -short ./..."; \
	fi

# Run tests without coverage (faster for quick feedback)
test-short: generate-mocks
	go test -v -short ./...

# Clean build artifacts
clean:
	rm -rf bin/
	rm -f coverage.out coverage.html

# Start Docker containers
docker-up:
	docker-compose up -d

# Stop Docker containers
docker-down:
	docker-compose down

# Run database migrations
migrate-up:
	goose -dir db/migrations postgres "postgresql://postgres:postgres@localhost:5432/cleanarch?sslmode=disable" up

# Rollback database migrations
migrate-down:
	goose -dir db/migrations postgres "postgresql://postgres:postgres@localhost:5432/cleanarch?sslmode=disable" down

# Create a new migration
migrate-create:
	@if [ -z "$(name)" ]; then \
		echo "Error: name parameter is required. Usage: make migrate-create name=migration_name"; \
		exit 1; \
	fi
	goose -dir db/migrations create $(name) sql

# Generate sqlc code
sqlc-generate:
	sqlc generate

# Generate mocks for interfaces
generate-mocks:
	@if ! command -v mockery > /dev/null; then \
		echo "Installing mockery..."; \
		go install github.com/vektra/mockery/v2@latest; \
	fi
	mockery

# Setup TDD environment
setup-tdd: deps generate-mocks verify
	@echo "TDD environment ready!"

# Verify all tests pass
verify: generate-mocks
	@echo "Running tests to verify everything works..."
	go test -v -cover ./...
	@echo "All tests passed!"

