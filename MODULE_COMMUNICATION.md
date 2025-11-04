# Module-to-Module Communication Guide

## Overview

This project demonstrates **bidirectional module communication** at the **application layer** in a Clean Architecture + DDD system. The implementation shows how the `Product` module and `Inventory` module communicate with each other while maintaining separation of concerns and testability.

## Communication Pattern

### Architecture Decision

**Communication happens at the Application Layer (Use Cases), NOT at the Repository Layer.**

This approach:
- ✅ Maintains proper module boundaries
- ✅ Keeps business logic in use cases
- ✅ Allows for easier testing with mocks
- ✅ Follows Clean Architecture principles
- ✅ Enables graceful degradation

## Implementation

### 1. Inventory → Product Communication

The Inventory module needs to validate that products exist before creating inventory records.

**Implementation in `internal/application/inventory/create.go`:**

```go
type CreateInventoryUseCase struct {
    inventoryRepo  inventory.InventoryRepository
    productUseCase product.ProductUseCaseInterface  // ← Injected Product module use case
}

func (uc *CreateInventoryUseCase) Execute(ctx context.Context, input CreateInventoryInput) (*CreateInventoryOutput, error) {
    // MODULE COMMUNICATION: Call Product module to verify product exists
    productOutput, err := uc.productUseCase.Execute(ctx, input.ProductID)
    if err != nil {
        return nil, errors.New("cannot create inventory: product not found")
    }
    
    // Continue with inventory creation...
}
```

**Key Points:**
- Inventory use case **depends on an interface** (`ProductUseCaseInterface`), not a concrete implementation
- Product validation happens before inventory creation
- Clear error messages indicate cross-module dependencies

### 2. Product → Inventory Communication (Bidirectional)

The Product module enriches product data with inventory information when available.

**Implementation in `internal/application/product/get.go`:**

```go
type GetProductUseCase struct {
    productRepo      product.ProductRepository
    inventoryUseCase InventoryUseCaseInterface  // ← Optional Inventory module integration
}

func (uc *GetProductUseCase) Execute(ctx context.Context, productID string) (*GetProductOutput, error) {
    // Fetch product data
    prod, err := uc.productRepo.GetByID(ctx, productID)
    // ... error handling ...

    output := &GetProductOutput{
        ID:           prod.ID(),
        Name:         prod.Name(),
        HasInventory: false,
    }

    // MODULE COMMUNICATION: Enrich with inventory data if available
    if uc.inventoryUseCase != nil {
        inventoryData, err := uc.inventoryUseCase.Execute(ctx, productID)
        if err == nil {
            output.HasInventory = true
            output.StockQuantity = inventoryData.GetQuantity()
            output.AvailableQuantity = inventoryData.GetAvailableQuantity()
        }
        // Gracefully handle inventory not found - product data is still valid
    }

    return output, nil
}
```

**Key Points:**
- Inventory integration is **optional** (nil check)
- Product module works with or without inventory
- **Graceful degradation** - if inventory service fails, product data is still returned
- Uses adapter pattern to decouple from Inventory module's concrete types

### 3. Dependency Wiring in `cmd/api/main.go`

The application bootstrap demonstrates how to wire bidirectional dependencies:

```go
func main() {
    // ... database setup ...

    // Initialize repositories
    productRepo := persistence.NewProductRepository(db)
    inventoryRepo := persistence.NewInventoryRepository(db)

    // STEP 1: Initialize product use cases (without inventory integration first)
    createProductUseCase := product.NewCreateProductUseCase(productRepo)
    getProductUseCaseBasic := product.NewGetProductUseCase(productRepo)

    // STEP 2: Initialize inventory use cases with product use case injection
    // This demonstrates Inventory → Product module communication
    createInventoryUseCase := inventory.NewCreateInventoryUseCase(inventoryRepo, getProductUseCaseBasic)
    getInventoryUseCase := inventory.NewGetInventoryUseCase(inventoryRepo, getProductUseCaseBasic)
    adjustInventoryUseCase := inventory.NewAdjustInventoryUseCase(inventoryRepo, getProductUseCaseBasic)

    // STEP 3: Create adapter for Product → Inventory communication
    inventoryAdapter := inventory.NewProductInventoryAdapter(getInventoryUseCase)

    // STEP 4: Re-initialize product get use case WITH inventory integration
    // This demonstrates Product → Inventory bidirectional module communication
    getProductUseCase := product.NewGetProductUseCaseWithInventory(productRepo, inventoryAdapter)

    // Initialize handlers
    productHandler := delivery.NewProductHandler(createProductUseCase, getProductUseCase)
    inventoryHandler := delivery.NewInventoryHandler(createInventoryUseCase, getInventoryUseCase, adjustInventoryUseCase)
    
    // ... register routes ...
}
```

**Key Points:**
- Product use case created **twice**: once standalone, once with inventory
- Inventory use cases get the standalone product use case
- Final product use case gets inventory integration via adapter
- This pattern avoids circular dependencies

## Adapter Pattern

To maintain loose coupling, we use adapters to translate between module interfaces:

**`internal/application/inventory/product_adapter.go`:**

```go
type ProductInventoryAdapter struct {
    inventoryUseCase *GetInventoryUseCase
}

func (a *ProductInventoryAdapter) Execute(ctx context.Context, productID string) (product.InventoryData, error) {
    output, err := a.inventoryUseCase.Execute(ctx, productID)
    if err != nil {
        return nil, err
    }
    return NewInventoryAdapter(output), nil
}
```

This adapter:
- Implements `product.InventoryUseCaseInterface`
- Translates Inventory module's `GetInventoryOutput` to Product module's `InventoryData` interface
- Prevents direct coupling between modules

## Testing

### Testing Inventory → Product Communication

**`internal/application/inventory/create_test.go`:**

```go
func TestCreateInventoryUseCase_Execute_ProductNotFound(t *testing.T) {
    mockInventoryRepo := new(MockInventoryRepository)
    mockProductUseCase := new(MockProductUseCase)
    useCase := NewCreateInventoryUseCase(mockInventoryRepo, mockProductUseCase)

    // Mock product not found - demonstrates module communication
    mockProductUseCase.On("Execute", mock.Anything, "nonexistent-product").
        Return(nil, errors.New("product not found"))

    // Execute
    output, err := useCase.Execute(context.Background(), input)

    // Assert
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "cannot create inventory: product not found")
}
```

### Testing Product → Inventory Communication

**`internal/application/product/get_test.go`:**

```go
func (s *GetProductWithInventoryUseCaseTestSuite) TestExecute_WithInventory() {
    // Setup product
    expectedProduct := domainProduct.ReconstructProduct(...)
    s.mockRepo.On("GetByID", mock.Anything, "test-product-id").
        Return(expectedProduct, nil)

    // Setup inventory data
    inventoryData := &MockInventoryData{
        quantity:          100,
        availableQuantity: 80,
    }

    // Mock inventory use case - demonstrates Product → Inventory communication
    s.mockInventoryUseCase.On("Execute", mock.Anything, "test-product-id").
        Return(inventoryData, nil)

    // Execute
    output, err := s.useCase.Execute(context.Background(), "test-product-id")

    // Assert - Product enriched with inventory
    s.NoError(err)
    s.True(output.HasInventory)
    s.Equal(100, output.StockQuantity)
    s.Equal(80, output.AvailableQuantity)
}
```

## API Endpoints

### Product Endpoints

**GET `/api/v1/products/:id`** - Returns product with inventory if available

```json
{
  "success": true,
  "message": "Product retrieved successfully",
  "data": {
    "id": "product-123",
    "name": "Laptop",
    "price_amount": 999.99,
    "price_currency": "USD",
    "has_inventory": true,
    "stock_quantity": 100,
    "available_quantity": 80
  }
}
```

### Inventory Endpoints

**POST `/api/v1/inventory`** - Create inventory (validates product exists)

```json
{
  "product_id": "product-123",
  "quantity": 100,
  "location": "Warehouse A"
}
```

**GET `/api/v1/inventory/:productId`** - Get inventory (includes product details)

```json
{
  "success": true,
  "message": "Inventory retrieved successfully",
  "data": {
    "id": "inv-1",
    "product_id": "product-123",
    "product_name": "Laptop",
    "product_price": 999.99,
    "product_currency": "USD",
    "quantity": 100,
    "reserved_quantity": 20,
    "available_quantity": 80,
    "location": "Warehouse A"
  }
}
```

**PATCH `/api/v1/inventory/adjust`** - Adjust inventory (validates product exists)

```json
{
  "product_id": "product-123",
  "adjustment": 50,
  "reason": "Restock"
}
```

## Benefits of This Approach

### 1. **Loose Coupling**
- Modules depend on interfaces, not implementations
- Easy to swap implementations for testing or different environments

### 2. **Testability**
- Each module can be tested in isolation with mocks
- No need for complex integration test setup

### 3. **Graceful Degradation**
- Product module works even if inventory module is unavailable
- System remains resilient to partial failures

### 4. **Clear Boundaries**
- Module responsibilities are well-defined
- Easy to understand data flow

### 5. **Scalability**
- Modules can be separated into microservices later
- Communication pattern already established

## Common Patterns

### Pattern 1: Validation Across Modules

When creating a resource in Module A that depends on Module B:

```go
// Validate dependency exists
dependencyData, err := moduleB.UseCase.Execute(ctx, dependencyID)
if err != nil {
    return nil, fmt.Errorf("cannot create resource: dependency not found")
}

// Proceed with creation
```

### Pattern 2: Enrichment Across Modules

When fetching a resource that can be enriched with data from another module:

```go
// Fetch primary data
primaryData := ...

// Optionally enrich with secondary data
if secondaryUseCase != nil {
    secondaryData, err := secondaryUseCase.Execute(ctx, id)
    if err == nil {
        primaryData.EnrichWith(secondaryData)
    }
}
```

### Pattern 3: Adapter for Decoupling

When Module A needs data from Module B but shouldn't depend on B's types:

```go
// Module A defines what it needs
type DataInterface interface {
    GetValue() string
}

// Adapter translates Module B's output to Module A's interface
type ModuleBAdapter struct {
    useCase *ModuleBUseCase
}

func (a *ModuleBAdapter) Execute(ctx context.Context, id string) (DataInterface, error) {
    output, err := a.useCase.Execute(ctx, id)
    return &Adapter{output}, err
}
```

## Running the Application

### 1. Start Database
```bash
make docker-up
```

### 2. Run Migrations
```bash
make migrate-up
```

### 3. Start Application
```bash
make run
```

### 4. Test Module Communication

**Create a product:**
```bash
curl -X POST http://localhost:8080/api/v1/products \
  -H "Content-Type: application/json" \
  -d '{"name": "Laptop", "price_amount": 999.99, "price_currency": "USD"}'
```

**Create inventory (Inventory → Product communication):**
```bash
curl -X POST http://localhost:8080/api/v1/inventory \
  -H "Content-Type: application/json" \
  -d '{"product_id": "<product-id>", "quantity": 100, "location": "Warehouse A"}'
```

**Get product with inventory (Product → Inventory communication):**
```bash
curl http://localhost:8080/api/v1/products/<product-id>
```

**Get inventory with product details (Inventory → Product communication):**
```bash
curl http://localhost:8080/api/v1/inventory/<product-id>
```

## Summary

This implementation demonstrates a **production-ready pattern for module-to-module communication** in Clean Architecture:

- ✅ Communication at the **application layer** via use case injection
- ✅ **Bidirectional** communication (Inventory ↔ Product)
- ✅ **Loose coupling** through interfaces and adapters
- ✅ **Highly testable** with mocks
- ✅ **Graceful degradation** when modules are unavailable
- ✅ Follows **SOLID principles** and Clean Architecture
- ✅ Ready for **scaling** to microservices if needed

This pattern can be extended to add more modules (e.g., Orders, Shipments) that need to communicate with Product and Inventory modules.

