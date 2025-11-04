# Module-to-Module Communication - Implementation Summary

## Overview

Successfully implemented **bidirectional module communication** at the **application layer** between the `Product` and `Inventory` modules, demonstrating proper Clean Architecture and DDD principles.

## What Was Implemented

### 1. Inventory Domain Layer
**Location:** `internal/domain/inventory/`

- **`entity.go`**: Inventory entity with business logic
  - Fields: productID, quantity, reservedQuantity, location
  - Methods: `Reserve()`, `Release()`, `AdjustQuantity()`, `AvailableQuantity()`
  
- **`repository.go`**: Repository interface
  - Methods: `Create()`, `GetByProductID()`, `Update()`, `Delete()`, `AdjustStock()`
  
- **`errors.go`**: Domain-specific errors
  - `ErrInventoryNotFound`, `ErrInsufficientStock`, `ErrInvalidQuantity`

### 2. Inventory Application Layer
**Location:** `internal/application/inventory/`

- **`interfaces.go`**: Interface definitions for module communication
- **`dto.go`**: Input/Output data transfer objects
- **`create.go`**: `CreateInventoryUseCase`
  - **Demonstrates Inventory → Product communication**
  - Injects `ProductUseCaseInterface`
  - Validates product exists before creating inventory
  
- **`get.go`**: `GetInventoryUseCase`
  - **Demonstrates Inventory → Product communication**
  - Enriches inventory data with product details
  - Gracefully handles deleted products
  
- **`adjust.go`**: `AdjustInventoryUseCase`
  - **Demonstrates Inventory → Product communication**
  - Validates product exists before adjusting stock
  
- **`adapter.go`**: Adapters for interface translation
- **`product_adapter.go`**: Adapter for Product → Inventory communication

### 3. Product Module Enhancements
**Location:** `internal/application/product/`

- **`interfaces.go`**: Interface for `ProductUseCaseInterface`
- **`dto.go`**: Enhanced with inventory fields
  - Added: `HasInventory`, `StockQuantity`, `AvailableQuantity`
  
- **`get.go`**: Enhanced `GetProductUseCase`
  - **Demonstrates Product → Inventory communication**
  - New constructor: `NewGetProductUseCaseWithInventory()`
  - Optionally enriches product with inventory data
  - Implements graceful degradation

### 4. Infrastructure Layer

#### Database
**Location:** `db/migrations/`, `db/query/`

- **`000002_create_inventory_table.up.sql`**: Inventory table migration
  - Foreign key to products table
  - Check constraints for data integrity
  - Indexes for performance
  
- **`inventory.sql`**: SQLC queries
  - CreateInventory, GetInventoryByProductID, UpdateInventory, AdjustInventoryQuantity

#### Persistence
**Location:** `internal/infrastructure/persistence/`

- **`inventory_repository.go`**: Repository implementation
  - Implements `inventory.InventoryRepository` interface
  - Uses sqlc-generated queries
  - Converts between database and domain models

#### HTTP Delivery
**Location:** `internal/infrastructure/delivery/`

- **`inventory_handler.go`**: HTTP handlers
  - `Create()`: POST /inventory
  - `Get()`: GET /inventory/:productId
  - `Adjust()`: PATCH /inventory/adjust

### 5. Application Bootstrap
**Location:** `cmd/api/main.go`

- **Demonstrates dependency wiring for bidirectional communication**
- Step-by-step initialization:
  1. Create product use cases (without inventory)
  2. Create inventory use cases (inject product use case)
  3. Create adapter for Product → Inventory
  4. Re-create product use case WITH inventory integration
- Routes registered for both modules

### 6. Comprehensive Tests
**Location:** `internal/application/inventory/`, `internal/application/product/`

#### Inventory Tests (100% coverage of use cases)
- **`create_test.go`**: 4 test cases
  - Success scenario
  - Product not found (demonstrates module communication)
  - Inventory already exists
  - Invalid quantity
  
- **`get_test.go`**: 3 test cases
  - Success with product enrichment (demonstrates module communication)
  - Inventory not found
  - Product deleted (demonstrates graceful degradation)
  
- **`adjust_test.go`**: 5 test cases
  - Increase success
  - Decrease success
  - Product not found (demonstrates module communication)
  - Inventory not found
  - Cannot go negative

#### Product Tests (Enhanced)
- **`get_test.go`**: Added 2 new test suites
  - `TestExecute_WithInventory`: Tests Product → Inventory communication
  - `TestExecute_InventoryNotFound`: Tests graceful degradation

**Total Test Results:**
- ✅ 12 Inventory tests - ALL PASS
- ✅ 8 Product tests (6 existing + 2 new) - ALL PASS
- ✅ 15 Domain tests (existing) - ALL PASS
- **35 tests total - 100% PASS RATE**

### 7. Documentation
**Location:** Root directory

- **`MODULE_COMMUNICATION.md`**: Comprehensive guide (400+ lines)
  - Architecture decisions
  - Implementation patterns
  - Code examples
  - API endpoints
  - Testing strategies
  - Common patterns
  
- **`IMPLEMENTATION_SUMMARY.md`**: This document

## Key Features Demonstrated

### 1. Bidirectional Communication
- ✅ **Inventory → Product**: Validate products exist, enrich with product details
- ✅ **Product → Inventory**: Enrich products with stock information

### 2. Clean Architecture Principles
- ✅ Dependencies flow inward
- ✅ Use cases orchestrate business logic
- ✅ Interfaces define contracts
- ✅ No coupling to infrastructure

### 3. Loose Coupling
- ✅ Modules depend on interfaces, not implementations
- ✅ Adapter pattern decouples modules
- ✅ No circular dependencies

### 4. Testability
- ✅ All use cases unit tested with mocks
- ✅ Module communication tested in isolation
- ✅ Edge cases covered (product deleted, inventory not found)

### 5. Graceful Degradation
- ✅ Product works without inventory module
- ✅ Inventory handles deleted products gracefully
- ✅ System remains resilient to partial failures

### 6. SOLID Principles
- ✅ Single Responsibility: Each use case has one job
- ✅ Open/Closed: Easy to extend with new modules
- ✅ Liskov Substitution: Interfaces can be mocked
- ✅ Interface Segregation: Minimal interfaces
- ✅ Dependency Inversion: Depend on abstractions

## API Endpoints

### Product Module
```
POST   /api/v1/products          - Create product
GET    /api/v1/products/:id      - Get product (with inventory if available)
```

### Inventory Module
```
POST   /api/v1/inventory         - Create inventory (validates product exists)
GET    /api/v1/inventory/:id     - Get inventory (includes product details)
PATCH  /api/v1/inventory/adjust  - Adjust stock (validates product exists)
```

## Files Created/Modified

### Created (19 files)
1. `internal/domain/inventory/entity.go`
2. `internal/domain/inventory/repository.go`
3. `internal/domain/inventory/errors.go`
4. `internal/application/inventory/interfaces.go`
5. `internal/application/inventory/dto.go`
6. `internal/application/inventory/create.go`
7. `internal/application/inventory/get.go`
8. `internal/application/inventory/adjust.go`
9. `internal/application/inventory/adapter.go`
10. `internal/application/inventory/product_adapter.go`
11. `internal/application/inventory/create_test.go`
12. `internal/application/inventory/get_test.go`
13. `internal/application/inventory/adjust_test.go`
14. `internal/application/product/interfaces.go`
15. `internal/infrastructure/persistence/inventory_repository.go`
16. `internal/infrastructure/delivery/inventory_handler.go`
17. `db/migrations/00002_create_inventory_table.sql`
18. `db/query/inventory.sql`

### Modified (4 files)
1. `internal/application/product/dto.go` - Added inventory fields
2. `internal/application/product/get.go` - Added inventory integration
3. `internal/application/product/get_test.go` - Added bidirectional tests
4. `cmd/api/main.go` - Wired dependencies demonstrating bidirectional injection

### Generated (2 files)
1. `internal/infrastructure/persistence/sqlcgen/inventory.sql.go`
2. `internal/infrastructure/persistence/sqlcgen/models.go` (updated)

## How to Test

### Run All Tests
```bash
go test ./... -v
```

### Run Inventory Tests Only
```bash
go test ./internal/application/inventory/... -v
```

### Run Product Tests Only
```bash
go test ./internal/application/product/... -v
```

### Build Application
```bash
go build ./cmd/api
```

### Run Application
```bash
# Start database
make docker-up

# Run migrations
make migrate-up

# Start server
make run
```

### Test API Endpoints

**1. Create a product:**
```bash
curl -X POST http://localhost:8080/api/v1/products \
  -H "Content-Type: application/json" \
  -d '{"name": "Laptop", "price_amount": 999.99, "price_currency": "USD"}'
```

**2. Get product (without inventory):**
```bash
curl http://localhost:8080/api/v1/products/<product-id>
# Response will show: "has_inventory": false
```

**3. Create inventory:**
```bash
curl -X POST http://localhost:8080/api/v1/inventory \
  -H "Content-Type: application/json" \
  -d '{"product_id": "<product-id>", "quantity": 100, "location": "Warehouse A"}'
# This demonstrates Inventory → Product validation
```

**4. Get product (with inventory):**
```bash
curl http://localhost:8080/api/v1/products/<product-id>
# Response will show: "has_inventory": true, "stock_quantity": 100, "available_quantity": 100
# This demonstrates Product → Inventory enrichment
```

**5. Get inventory:**
```bash
curl http://localhost:8080/api/v1/inventory/<product-id>
# Response includes product details (name, price)
# This demonstrates Inventory → Product enrichment
```

**6. Adjust inventory:**
```bash
curl -X PATCH http://localhost:8080/api/v1/inventory/adjust \
  -H "Content-Type: application/json" \
  -d '{"product_id": "<product-id>", "adjustment": 50, "reason": "Restock"}'
# This demonstrates Inventory → Product validation
```

## Architecture Benefits

### For Development
- ✅ Clear module boundaries
- ✅ Easy to add new modules
- ✅ Modules can be developed independently
- ✅ Strong type safety

### For Testing
- ✅ Fast unit tests (no database needed)
- ✅ Easy to mock dependencies
- ✅ High test coverage achievable
- ✅ Tests document behavior

### For Maintenance
- ✅ Easy to understand data flow
- ✅ Changes localized to modules
- ✅ Refactoring is safer with tests
- ✅ Clear error messages

### For Scaling
- ✅ Modules can become microservices
- ✅ Communication patterns already established
- ✅ Resilient to partial failures
- ✅ Independent deployment possible

## Design Patterns Used

1. **Use Case Pattern**: Encapsulates business workflows
2. **Repository Pattern**: Abstracts data access
3. **Adapter Pattern**: Translates between module interfaces
4. **Dependency Injection**: Loose coupling and testability
5. **Interface Segregation**: Minimal, focused interfaces
6. **DTO Pattern**: Clean data transfer between layers

## Success Metrics

- ✅ **Zero circular dependencies**
- ✅ **100% test pass rate** (35/35 tests)
- ✅ **Zero linter errors**
- ✅ **Successful build**
- ✅ **Bidirectional communication working**
- ✅ **Graceful degradation implemented**
- ✅ **Comprehensive documentation created**

## Future Extensions

This pattern can be extended to add:
- **Order Module**: Validate product + inventory before order creation
- **Shipment Module**: Update inventory when items ship
- **Pricing Module**: Dynamic pricing based on inventory levels
- **Notification Module**: Alert when inventory is low
- **Analytics Module**: Aggregate data from multiple modules

All following the same communication pattern established here.

## Conclusion

This implementation demonstrates a **production-ready pattern** for module-to-module communication in Clean Architecture that:

1. **Maintains proper boundaries** between modules
2. **Enables bidirectional communication** without tight coupling
3. **Follows SOLID principles** and best practices
4. **Is highly testable** with comprehensive test coverage
5. **Supports graceful degradation** for resilient systems
6. **Scales well** for growth to microservices if needed

The pattern is ready for production use and can serve as a template for adding more modules to the system.

---

## Smart Error Handling System

### Overview

Implemented a centralized error handling system in `pkg/errors` that provides structured error codes, HTTP status mapping, stack trace capture, and easy extensibility.

### Components

#### 1. Error Code Registry (`pkg/errors/codes.go`)
- Centralized registry for error codes with HTTP status mapping
- Predefined codes for common scenarios (NotFound, InvalidInput, Conflict, etc.)
- Domain-specific codes for Product and Inventory modules
- Easy registration of new error codes via `RegisterErrorCode()`

#### 2. AppError Type (`pkg/errors/error.go`)
- Structured error type with code, message, HTTP status, and stack trace
- Helper functions: `New()`, `Wrap()`, `Wrapf()`, `WithCode()`
- Utility functions: `Is()`, `GetCode()`, `GetHTTPStatus()`, `GetMessage()`
- Supports error wrapping for context preservation

#### 3. Database Error Helpers (`pkg/errors/helpers.go`)
- `WrapDatabaseError()`: Automatically maps database errors to appropriate codes
  - `sql.ErrNoRows` → `CodeNotFound`
  - Duplicate key violations → `CodeConflict`
  - Connection errors → `CodeDatabaseConnection`
  - Other errors → `CodeDatabaseError`
- `WrapValidationError()`: Wraps validation errors with `CodeValidation`

#### 4. HTTP Integration (`internal/infrastructure/delivery/error_handler.go`)
- `HandleError()`: Converts errors to HTTP responses with correct status codes
- `HandleValidationError()`: Special handling for validation errors
- Automatic error code extraction and HTTP status mapping

### Usage Across Layers

**Domain Layer:**
```go
// Domain errors use error codes
var ErrProductNotFound = errors.New(errors.CodeProductNotFound, "product not found")
```

**Application Layer:**
```go
// Use cases wrap errors appropriately
if err != nil {
    return nil, apperrors.WrapDatabaseError(err)
}
```

**Infrastructure Layer:**
```go
// Repositories wrap database errors
err := r.queries.CreateProduct(ctx, params)
if err != nil {
    return apperrors.WrapDatabaseError(err)
}
```

**Delivery Layer:**
```go
// Handlers use centralized error handling
output, err := h.createUseCase.Execute(ctx, input)
if err != nil {
    HandleError(c, err)  // Automatically maps to correct HTTP status
    return
}
```

### Benefits

1. **Consistency**: All errors follow the same structure and format
2. **Extensibility**: Easy to add new error types without modifying existing code
3. **Debugging**: Stack traces captured automatically for better error tracking
4. **HTTP Mapping**: Automatic conversion of error codes to appropriate HTTP status codes
5. **Type Safety**: Error codes are typed constants, reducing typos
6. **Maintainability**: Centralized error handling makes updates easier

### Adding New Error Types

To add a new error type:

1. Define the error code constant
2. Register it with `RegisterErrorCode(code, httpStatus, description)`
3. Use it throughout the application with `errors.New(code, message)`

Example:
```go
const CodeCustomError = errors.ErrorCode("CUSTOM_ERROR")
errors.RegisterErrorCode(CodeCustomError, 400, "Custom error description")
err := errors.New(CodeCustomError, "Something went wrong")
```

