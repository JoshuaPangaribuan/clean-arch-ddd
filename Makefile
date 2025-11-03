.PHONY: help deps run build test clean docker-up docker-down migrate-up migrate-down migrate-create sqlc-generate generate-mocks

# Default target
help:
	@echo "Available targets:"
	@echo "  deps            - Download and install dependencies"
	@echo "  run             - Run the application"
	@echo "  build           - Build the application"
	@echo "  test            - Run tests"
	@echo "  clean           - Clean build artifacts"
	@echo "  docker-up       - Start Docker containers"
	@echo "  docker-down     - Stop Docker containers"
	@echo "  migrate-up      - Run database migrations"
	@echo "  migrate-down    - Rollback database migrations"
	@echo "  migrate-create  - Create a new migration (usage: make migrate-create name=migration_name)"
	@echo "  sqlc-generate   - Generate sqlc code"
	@echo "  generate-mocks  - Generate mocks for interfaces"

# Download and install dependencies
deps:
	go mod download
	go mod tidy
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	go install github.com/vektra/mockery/v2@latest

# Run the application
run:
	go run cmd/api/main.go

# Build the application
build:
	go build -o bin/api cmd/api/main.go

# Run tests
test:
	go test -v -cover ./...

# Run tests with coverage report
test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

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
	migrate -path migrations -database "postgresql://postgres:postgres@localhost:5432/cleanarch?sslmode=disable" up

# Rollback database migrations
migrate-down:
	migrate -path migrations -database "postgresql://postgres:postgres@localhost:5432/cleanarch?sslmode=disable" down 1

# Create a new migration
migrate-create:
	@if [ -z "$(name)" ]; then \
		echo "Error: name parameter is required. Usage: make migrate-create name=migration_name"; \
		exit 1; \
	fi
	migrate create -ext sql -dir migrations -seq $(name)

# Generate sqlc code
sqlc-generate:
	sqlc generate

# Generate mocks for interfaces
generate-mocks:
	mockery

