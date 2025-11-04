#!/bin/bash
# ci-check.sh - Run tests to verify the project is in a good state
# This script is designed to be run in CI/CD pipelines

set -e  # Exit on any error

echo "🔍 Running CI checks..."

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Error: Go is not installed"
    exit 1
fi

echo "✅ Go version: $(go version)"

# Generate sqlc code
echo "📦 Generating sqlc code..."
if ! command -v sqlc &> /dev/null; then
    echo "Installing sqlc..."
    go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
fi
sqlc generate || {
    echo "❌ Error: Failed to generate sqlc code"
    exit 1
}

# Generate mocks
echo "📦 Generating mocks..."
if ! command -v mockery &> /dev/null; then
    echo "Installing mockery..."
    go install github.com/vektra/mockery/v2@latest
fi
mockery || {
    echo "❌ Error: Failed to generate mocks"
    exit 1
}

# Run tests
echo "🧪 Running tests..."
go test -v -cover ./... || {
    echo "❌ Error: Tests failed"
    exit 1
}

# Run tests with coverage report
echo "📊 Generating coverage report..."
go test -coverprofile=coverage.out ./... || {
    echo "❌ Error: Failed to generate coverage report"
    exit 1
}

# Show coverage summary
coverage=$(go tool cover -func=coverage.out | grep total | awk '{print $3}')
echo "✅ Test coverage: $coverage"

echo "✅ All CI checks passed!"
exit 0
