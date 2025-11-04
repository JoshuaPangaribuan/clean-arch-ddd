# Clean Architecture + DDD Template (Golang)

A production-ready Golang project template implementing Clean Architecture and Domain-Driven Design (DDD) principles. This template provides a solid foundation for building scalable and maintainable backend applications.

## 🏗️ Architecture Overview

This template follows **Clean Architecture** principles with **Domain-Driven Design** patterns, organizing code into distinct layers with clear responsibilities:

### Layers

```
┌─────────────────────────────────────────┐
│         Delivery Layer (HTTP)           │  ← Gin Handlers, Middleware
├─────────────────────────────────────────┤
│       Application Layer (Use Cases)     │  ← Business workflows, DTOs
├─────────────────────────────────────────┤
│         Domain Layer (Entities)         │  ← Business logic, Rules
├─────────────────────────────────────────┤
│    Infrastructure (Persistence, Config) │  ← Database, External services
└─────────────────────────────────────────┘
```

#### 1. **Domain Layer** (`internal/domain/`)
- **Pure business logic** with no external dependencies
- **Entities**: Core business objects with identity and lifecycle
- **Value Objects**: Immutable objects defined by their attributes
- **Repository Interfaces**: Contracts for data persistence (not implementations)

#### 2. **Application Layer** (`internal/application/`)
- **Use Cases**: Orchestrate business workflows
- **DTOs**: Data transfer objects for input/output
- Depends on Domain layer only

#### 3. **Infrastructure Layer** (`internal/infrastructure/`)
- **Persistence**: Repository implementations using sqlc
- **Delivery**: HTTP handlers using Gin
- **Config**: Application configuration management
- Implements interfaces defined in Domain layer

#### 4. **Bootstrap** (`cmd/api/main.go`)
- Application entry point
- Dependency injection
- Wire all components together

### Key Design Principles

- **Dependency Inversion**: Inner layers define interfaces, outer layers implement them
- **Separation of Concerns**: Each layer has a single, well-defined responsibility
- **Testability**: Business logic isolated from infrastructure (easily mockable)
- **Domain-Centric**: Business rules live in the domain, not scattered across layers

## 🚀 Quick Start

### Prerequisites

- Go 1.21 or higher
- Docker & Docker Compose
- Make (optional, but recommended)

### 1. Clone and Setup

```bash
# Clone the repository
git clone <your-repo-url>
cd clean-arch-ddd

# Install dependencies
make deps
```

### 2. Start Database

```bash
# Start PostgreSQL using Docker Compose
make docker-up
```

### 3. Run Database Migrations

```bash
# Apply migrations
make migrate-up
```

### 4. Run the Application

```bash
# Run the API server
make run
```

The server will start on `http://localhost:8080`

### 5. Test the API

**Health Check:**
```bash
curl http://localhost:8080/health
```

**Create a Product:**
```bash
curl -X POST http://localhost:8080/api/v1/products \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Laptop",
    "price_amount": 999.99,
    "price_currency": "USD"
  }'
```

**Get a Product:**
```bash
curl http://localhost:8080/api/v1/products/{product-id}
```

## 📁 Project Structure

```
.
├── cmd/
│   └── api/
│       └── main.go                 # Application entry point
│
├── internal/                        # Private application code
│   ├── domain/                      # 🔵 DOMAIN LAYER
│   │   └── product/
│   │       ├── entity.go            # Product entity with business logic
│   │       ├── valueobject.go       # Price value object
│   │       └── repository.go        # Repository interface (no implementation)
│   │
│   ├── application/                 # 🟢 APPLICATION LAYER
│   │   └── product/
│   │       ├── create.go            # CreateProduct use case
│   │       ├── get.go               # GetProduct use case
│   │       ├── dto.go               # Input/Output DTOs
│   │       ├── create_test.go       # Unit tests
│   │       └── get_test.go          # Unit tests
│   │
│   ├── infrastructure/              # 🟡 INFRASTRUCTURE LAYER
│   │   ├── persistence/             # Database implementations
│   │   │   ├── product_repository.go  # Repository implementation
│   │   │   └── sqlcgen/             # Generated sqlc code
│   │   ├── delivery/                # HTTP layer
│   │   │   ├── product_handler.go   # HTTP handlers
│   │   │   └── middleware.go        # Logging, error handling, CORS
│   │   └── config/
│   │       └── config.go            # Configuration management
│   │
│   └── shared/                      # Shared utilities
│       └── model/
│           └── response.go          # API response models
│
├── db/                                # Database-related files
│   ├── migrations/                    # Database migrations (Goose)
│   │   ├── 00001_create_products_table.sql
│   │   └── 00002_create_inventory_table.sql
│   └── query/                         # SQL queries for sqlc
│       ├── product.sql
│       └── inventory.sql
│
├── mocks/                           # Generated mocks (mockery)
│   └── ProductRepository.go
│
├── pkg/                             # Public reusable code
│   └── errors/                      # Smart error handling package
│       ├── codes.go                 # Error codes and registry
│       ├── error.go                 # AppError type and functions
│       └── helpers.go               # Helper functions
├── api/                             # API documentation
├── docker-compose.yml               # Local development environment
├── Makefile                         # Common commands
├── sqlc.yaml                        # sqlc configuration
└── README.md                        # This file
```

## 🛠️ Available Commands

```bash
make deps            # Install dependencies
make run             # Run the application
make build           # Build the binary
make test            # Run all tests
make test-coverage   # Run tests with coverage report
make docker-up       # Start Docker containers
make docker-down     # Stop Docker containers
make migrate-up      # Apply database migrations
make migrate-down    # Rollback last migration
make migrate-create name=<name>  # Create new migration
make sqlc-generate   # Generate sqlc code
make generate-mocks  # Generate mocks for testing
make clean           # Clean build artifacts
```

## 🔧 Error Handling

This project uses a centralized error handling system located in `pkg/errors` that provides:

- **Structured error codes** with HTTP status mapping
- **Stack trace capture** for debugging
- **Easy extensibility** for adding new error types
- **Consistent error responses** across the application

### Using Error Codes

The error system uses predefined error codes that map to HTTP status codes:

```go
import apperrors "github.com/JoshuaPangaribuan/clean-arch-ddd/pkg/errors"

// Create a new error with a code
err := apperrors.New(apperrors.CodeNotFound, "Resource not found")

// Wrap an existing error
err := apperrors.Wrap(originalErr, apperrors.CodeInvalidInput, "Invalid input")

// Check error codes
if apperrors.Is(err, apperrors.CodeNotFound) {
    // Handle not found
}
```

### Available Error Codes

**Generic Errors:**
- `CodeInternalError` (500) - Internal server error
- `CodeInvalidInput` (400) - Invalid input provided
- `CodeNotFound` (404) - Resource not found
- `CodeConflict` (409) - Resource conflict
- `CodeValidation` (400) - Validation error

**Product Domain:**
- `CodeProductNotFound` (404)
- `CodeProductAlreadyExists` (409)
- `CodeInvalidProductID` (400)
- `CodeInvalidProductName` (400)
- `CodeInvalidPrice` (400)

**Inventory Domain:**
- `CodeInventoryNotFound` (404)
- `CodeInventoryExists` (409)
- `CodeInsufficientStock` (400)
- `CodeInvalidQuantity` (400)
- `CodeInvalidAdjustment` (400)

**Persistence Errors:**
- `CodeDatabaseError` (500)
- `CodeDatabaseConnection` (503)
- `CodeQueryFailed` (500)
- `CodeTransactionFailed` (500)

### Adding New Error Codes

To add a new error code, simply register it:

```go
import apperrors "github.com/JoshuaPangaribuan/clean-arch-ddd/pkg/errors"

// Define your error code
const CodeCustomError = apperrors.ErrorCode("CUSTOM_ERROR")

// Register it with HTTP status and description
apperrors.RegisterErrorCode(CodeCustomError, 400, "Custom error description")

// Use it
err := apperrors.New(CodeCustomError, "Something went wrong")
```

### Error Handling in Handlers

HTTP handlers automatically convert errors to appropriate HTTP responses:

```go
func (h *ProductHandler) Create(c *gin.Context) {
    // ... validation ...
    
    output, err := h.createUseCase.Execute(ctx, input)
    if err != nil {
        HandleError(c, err)  // Automatically maps to correct HTTP status
        return
    }
    
    // Success response
    c.JSON(http.StatusCreated, model.NewSuccessResponse(...))
}
```

### Database Error Wrapping

Database errors are automatically wrapped with appropriate error codes:

```go
import apperrors "github.com/JoshuaPangaribuan/clean-arch-ddd/pkg/errors"

err := r.queries.CreateProduct(ctx, params)
if err != nil {
    return apperrors.WrapDatabaseError(err)  // Automatically maps DB errors
}
```

This handles:
- `sql.ErrNoRows` → `CodeNotFound`
- Duplicate key violations → `CodeConflict`
- Connection errors → `CodeDatabaseConnection`
- Other database errors → `CodeDatabaseError`

## 📝 How to Add New Features

### Example: Adding a "Category" Feature

#### Step 1: Create Domain Layer

**`internal/domain/category/entity.go`**
```go
package category

type Category struct {
    id   string
    name string
}

func NewCategory(id, name string) (*Category, error) {
    // Add validation
    return &Category{id: id, name: name}, nil
}
```

**`internal/domain/category/repository.go`**
```go
package category

type CategoryRepository interface {
    Create(ctx context.Context, category *Category) error
    GetByID(ctx context.Context, id string) (*Category, error)
}
```

#### Step 2: Create Application Layer

**`internal/application/category/create.go`**
```go
package category

type CreateCategoryUseCase struct {
    categoryRepo domain.CategoryRepository
}

func (uc *CreateCategoryUseCase) Execute(ctx context.Context, input CreateCategoryInput) (*CreateCategoryOutput, error) {
    // Business logic here
}
```

#### Step 3: Create Infrastructure Layer

**Database Migration:**
```bash
make migrate-create name=create_categories_table
```

**SQL Queries in `db/query/category.sql`:**
```sql
-- name: CreateCategory :exec
INSERT INTO categories (id, name) VALUES ($1, $2);
```

**Repository Implementation:**
```go
// internal/infrastructure/persistence/category_repository.go
type CategoryRepositoryImpl struct {
    queries *sqlcgen.Queries
}
```

#### Step 4: Add HTTP Handlers

**`internal/infrastructure/delivery/category_handler.go`**
```go
type CategoryHandler struct {
    createUseCase *category.CreateCategoryUseCase
}

func (h *CategoryHandler) Create(c *gin.Context) {
    // Handle HTTP request
}
```

#### Step 5: Wire in `main.go`

```go
// Initialize repository
categoryRepo := persistence.NewCategoryRepository(db)

// Initialize use case
createCategoryUseCase := category.NewCreateCategoryUseCase(categoryRepo)

// Initialize handler
categoryHandler := delivery.NewCategoryHandler(createCategoryUseCase)

// Register routes
categories := v1.Group("/categories")
{
    categories.POST("", categoryHandler.Create)
}
```

## 🧪 Testing & Test-Driven Development (TDD)

This project is configured for **Test-Driven Development (TDD)**, emphasizing writing tests before implementation to ensure high code quality and maintainability.

### TDD Workflow

The TDD cycle follows three steps:

1. **🔴 RED**: Write a failing test that describes the desired behavior
2. **🟢 GREEN**: Write the minimal code to make the test pass
3. **🔵 REFACTOR**: Improve the code while keeping tests green

### Running Tests

```bash
# Run all tests with coverage (default TDD command)
make test

# Run unit tests only (fast, no external dependencies)
make test-unit

# Run integration tests (requires database)
make test-integration

# Run tests with HTML coverage report
make test-coverage

# Run tests without coverage (faster for quick feedback)
make test-short

# Setup TDD environment (deps + generate mocks + verify)
make setup-tdd

# Verify all tests pass
make verify
```

### TDD Commands for Fast Feedback

```bash
# Quick test run (unit tests only, no coverage)
make test-unit

# Full test suite with coverage
make test-coverage

# Watch mode (requires entr - install separately)
# On Windows, consider using: https://github.com/radovskyb/watcher
make test-watch
```

### Writing Tests

#### Domain Layer Tests

Domain tests are pure unit tests with no external dependencies. They test business logic and validation:

```go
func TestNewPrice(t *testing.T) {
    tests := []struct {
        name        string
        amount      float64
        currency    string
        wantErr     bool
        errContains string
    }{
        {
            name:     "valid price",
            amount:   99.99,
            currency: "USD",
            wantErr:  false,
        },
        {
            name:        "negative amount",
            amount:      -10.00,
            currency:    "USD",
            wantErr:     true,
            errContains: "negative",
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := product.NewPrice(tt.amount, tt.currency)
            // assertions...
        })
    }
}
```

#### Application Layer Tests

Application tests use mocks to isolate use cases from infrastructure:

```go
func TestCreateProductUseCase_Execute_Success(t *testing.T) {
    mockRepo := mocks.NewProductRepository(t)
    useCase := product.NewCreateProductUseCase(mockRepo)
    
    mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*product.Product")).
        Return(nil).
        Once()
    
    output, err := useCase.Execute(context.Background(), input)
    
    assert.NoError(t, err)
    assert.NotNil(t, output)
}
```

### Generating Mocks

Mocks are automatically generated before tests run. To manually generate:

```bash
make generate-mocks
```

Mocks are generated using `mockery` based on `.mockery.yaml` configuration.

### CI/CD Integration

Tests are automatically run in CI/CD pipelines:

```bash
# Run CI checks locally
./scripts/ci-check.sh  # Linux/Mac
scripts\ci-check.bat   # Windows
```

The GitHub Actions workflow (`.github/workflows/ci.yml`) runs tests on every push and pull request.

### TDD Best Practices

1. **Start with Domain**: Write domain tests first - they're fastest and validate core business logic
2. **Use Table-Driven Tests**: Write comprehensive test cases covering success and error scenarios
3. **Mock External Dependencies**: Use mocks for repositories and external services
4. **Keep Tests Fast**: Use `make test-unit` for quick feedback during development
5. **Test Behavior, Not Implementation**: Focus on what the code does, not how it does it
6. **Refactor Regularly**: Use green tests as a safety net to refactor and improve code

## 🔧 Configuration

Configuration is managed via environment variables using Viper. See `.env.example` for available options:

```env
# Server
SERVER_HOST=0.0.0.0
SERVER_PORT=8080

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=cleanarch
DB_SSLMODE=disable

# Application
APP_ENV=development
LOG_LEVEL=debug
```

Copy `.env.example` to `.env` and adjust values as needed.

## 📚 Tech Stack

- **Web Framework**: [Gin](https://github.com/gin-gonic/gin) - High-performance HTTP framework
- **Database**: PostgreSQL
- **Query Builder**: [sqlc](https://sqlc.dev/) - Type-safe SQL code generation
- **Migrations**: [goose](https://github.com/pressly/goose)
- **Configuration**: [Viper](https://github.com/spf13/viper)
- **Validation**: [go-playground/validator](https://github.com/go-playground/validator)
- **Mocking**: [mockery](https://github.com/vektra/mockery)
- **Testing**: [testify](https://github.com/stretchr/testify)

## 🎯 Design Decisions

### Why sqlc?
- **Type Safety**: Generates type-safe Go code from SQL
- **Performance**: Uses pure SQL without ORM overhead
- **Control**: Full control over queries while maintaining safety

### Why Clean Architecture + DDD?
- **Maintainability**: Clear separation makes code easier to understand and modify
- **Testability**: Business logic can be tested without infrastructure
- **Flexibility**: Easy to swap implementations (e.g., change database)
- **Scalability**: Structure supports growing complexity

### Why Value Objects?
- **Validation**: Encapsulate validation logic (e.g., Price cannot be negative)
- **Immutability**: Prevent accidental modifications
- **Domain Modeling**: Better represent business concepts

## 🚧 Out of Scope (Future Enhancements)

The following features are not included in v1.0 but may be added later:

- Authentication & Authorization (JWT, OAuth2)
- Message Broker integration (RabbitMQ, Kafka)
- Docker deployment configuration
- API documentation (Swagger/OpenAPI)
- Caching (Redis)
- gRPC support

## 📖 Additional Resources

- [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html) by Robert C. Martin
- [Domain-Driven Design](https://martinfowler.com/tags/domain%20driven%20design.html) resources
- [sqlc Documentation](https://docs.sqlc.dev/)
- [Gin Documentation](https://gin-gonic.com/docs/)

## 📄 License

This template is provided as-is for educational and commercial use.

## 🤝 Contributing

Feel free to submit issues and enhancement requests!

---

**Happy Coding! 🚀**

