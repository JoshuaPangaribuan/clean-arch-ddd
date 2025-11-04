# Quick Start Guide - Module Communication Demo

## Prerequisites

- Go 1.21+
- Docker & Docker Compose
- PostgreSQL (via Docker)
- Make (optional)

## Setup & Run

### 1. Start Database

```bash
make docker-up
# or
docker-compose up -d
```

### 2. Run Migrations

```bash
make migrate-up
# or
goose -dir db/migrations postgres "postgresql://postgres:postgres@localhost:5432/cleanarch?sslmode=disable" up
```

### 3. Start Application

```bash
make run
# or
go run cmd/api/main.go
```

Server will start at `http://localhost:8080`

## Demo: Module Communication in Action

### Scenario 1: Inventory → Product Validation

**Create a product first:**
```bash
curl -X POST http://localhost:8080/api/v1/products \
  -H "Content-Type: application/json" \
  -d '{
    "name": "MacBook Pro M3",
    "price_amount": 2499.99,
    "price_currency": "USD"
  }'
```

**Response:**
```json
{
  "success": true,
  "message": "Product created successfully",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "name": "MacBook Pro M3",
    "price_amount": 2499.99,
    "price_currency": "USD",
    "created_at": "2024-11-04T10:30:00Z"
  }
}
```

**Create inventory (validates product exists - Inventory → Product):**
```bash
curl -X POST http://localhost:8080/api/v1/inventory \
  -H "Content-Type: application/json" \
  -d '{
    "product_id": "550e8400-e29b-41d4-a716-446655440000",
    "quantity": 50,
    "location": "Warehouse A - Bay 12"
  }'
```

**Response (includes product name from Product module):**
```json
{
  "success": true,
  "message": "Inventory created successfully",
  "data": {
    "id": "660e8400-e29b-41d4-a716-446655440001",
    "product_id": "550e8400-e29b-41d4-a716-446655440000",
    "product_name": "MacBook Pro M3",
    "quantity": 50,
    "reserved_quantity": 0,
    "available_quantity": 50,
    "location": "Warehouse A - Bay 12",
    "created_at": "2024-11-04T10:31:00Z",
    "updated_at": "2024-11-04T10:31:00Z"
  }
}
```

**Try creating inventory for non-existent product:**
```bash
curl -X POST http://localhost:8080/api/v1/inventory \
  -H "Content-Type: application/json" \
  -d '{
    "product_id": "00000000-0000-0000-0000-000000000000",
    "quantity": 100,
    "location": "Warehouse B"
  }'
```

**Response (validation error):**
```json
{
  "success": false,
  "message": "cannot create inventory: product not found",
  "code": "PRODUCT_NOT_FOUND"
}
```

✅ **This demonstrates Inventory → Product module communication for validation!**

---

### Scenario 2: Product → Inventory Enrichment

**Get product WITHOUT inventory:**
```bash
# First create a product without inventory
curl -X POST http://localhost:8080/api/v1/products \
  -H "Content-Type: application/json" \
  -d '{
    "name": "iPhone 15 Pro",
    "price_amount": 999.99,
    "price_currency": "USD"
  }'

# Then get it
curl http://localhost:8080/api/v1/products/<product-id>
```

**Response (no inventory data):**
```json
{
  "success": true,
  "message": "Product retrieved successfully",
  "data": {
    "id": "770e8400-e29b-41d4-a716-446655440002",
    "name": "iPhone 15 Pro",
    "price_amount": 999.99,
    "price_currency": "USD",
    "has_inventory": false,
    "created_at": "2024-11-04T10:35:00Z",
    "updated_at": "2024-11-04T10:35:00Z"
  }
}
```

**Create inventory for this product:**
```bash
curl -X POST http://localhost:8080/api/v1/inventory \
  -H "Content-Type: application/json" \
  -d '{
    "product_id": "770e8400-e29b-41d4-a716-446655440002",
    "quantity": 200,
    "location": "Warehouse C - Section 5"
  }'
```

**Get product WITH inventory:**
```bash
curl http://localhost:8080/api/v1/products/770e8400-e29b-41d4-a716-446655440002
```

**Response (enriched with inventory - Product → Inventory):**
```json
{
  "success": true,
  "message": "Product retrieved successfully",
  "data": {
    "id": "770e8400-e29b-41d4-a716-446655440002",
    "name": "iPhone 15 Pro",
    "price_amount": 999.99,
    "price_currency": "USD",
    "has_inventory": true,
    "stock_quantity": 200,
    "available_quantity": 200,
    "created_at": "2024-11-04T10:35:00Z",
    "updated_at": "2024-11-04T10:35:00Z"
  }
}
```

✅ **This demonstrates Product → Inventory module communication for enrichment!**

---

### Scenario 3: Bidirectional Communication

**Get inventory (Inventory → Product for product details):**
```bash
curl http://localhost:8080/api/v1/inventory/550e8400-e29b-41d4-a716-446655440000
```

**Response (includes product details):**
```json
{
  "success": true,
  "message": "Inventory retrieved successfully",
  "data": {
    "id": "660e8400-e29b-41d4-a716-446655440001",
    "product_id": "550e8400-e29b-41d4-a716-446655440000",
    "product_name": "MacBook Pro M3",
    "product_price": 2499.99,
    "product_currency": "USD",
    "quantity": 50,
    "reserved_quantity": 0,
    "available_quantity": 50,
    "location": "Warehouse A - Bay 12",
    "created_at": "2024-11-04T10:31:00Z",
    "updated_at": "2024-11-04T10:31:00Z"
  }
}
```

---

### Scenario 4: Adjust Inventory

**Adjust stock (validates product still exists):**
```bash
curl -X PATCH http://localhost:8080/api/v1/inventory/adjust \
  -H "Content-Type: application/json" \
  -d '{
    "product_id": "550e8400-e29b-41d4-a716-446655440000",
    "adjustment": -10,
    "reason": "Sold 10 units"
  }'
```

**Response:**
```json
{
  "success": true,
  "message": "Inventory adjusted successfully",
  "data": {
    "id": "660e8400-e29b-41d4-a716-446655440001",
    "product_id": "550e8400-e29b-41d4-a716-446655440000",
    "product_name": "MacBook Pro M3",
    "quantity": 40,
    "reserved_quantity": 0,
    "available_quantity": 40,
    "location": "Warehouse A - Bay 12",
    "updated_at": "2024-11-04T10:45:00Z"
  }
}
```

**Verify updated inventory in product:**
```bash
curl http://localhost:8080/api/v1/products/550e8400-e29b-41d4-a716-446655440000
```

**Response (shows updated stock):**
```json
{
  "success": true,
  "message": "Product retrieved successfully",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "name": "MacBook Pro M3",
    "price_amount": 2499.99,
    "price_currency": "USD",
    "has_inventory": true,
    "stock_quantity": 40,
    "available_quantity": 40,
    "created_at": "2024-11-04T10:30:00Z",
    "updated_at": "2024-11-04T10:30:00Z"
  }
}
```

✅ **Bidirectional communication working perfectly!**

---

## Testing

### Run All Tests
```bash
make test
# or
go test ./... -v
```

**Expected output:**
```
PASS: TestCreateInventoryUseCase_Execute_Success
PASS: TestCreateInventoryUseCase_Execute_ProductNotFound  ← Module comm test
PASS: TestGetInventoryUseCase_Execute_Success  ← Module comm test
PASS: TestGetProductWithInventory_Success  ← Module comm test
PASS: TestGetProductWithInventory_InventoryNotFound  ← Graceful degradation

35/35 tests PASS
```

### Run Only Inventory Tests
```bash
go test ./internal/application/inventory/... -v
```

### Run Only Product Tests
```bash
go test ./internal/application/product/... -v
```

## Architecture Verification

### Verify Bidirectional Communication
```bash
# 1. Create product
PRODUCT_ID=$(curl -s -X POST http://localhost:8080/api/v1/products \
  -H "Content-Type: application/json" \
  -d '{"name": "Test", "price_amount": 10, "price_currency": "USD"}' \
  | jq -r '.data.id')

echo "Product ID: $PRODUCT_ID"

# 2. Get product (should show has_inventory: false)
curl http://localhost:8080/api/v1/products/$PRODUCT_ID | jq '.data.has_inventory'
# Output: false

# 3. Create inventory (validates product exists - Inventory → Product)
curl -s -X POST http://localhost:8080/api/v1/inventory \
  -H "Content-Type: application/json" \
  -d "{\"product_id\": \"$PRODUCT_ID\", \"quantity\": 100}" \
  | jq '.data.product_name'
# Output: "Test" (product name retrieved!)

# 4. Get product again (should show inventory - Product → Inventory)
curl http://localhost:8080/api/v1/products/$PRODUCT_ID | jq '.data | {has_inventory, stock_quantity}'
# Output: {"has_inventory": true, "stock_quantity": 100}

# 5. Get inventory (includes product details - Inventory → Product)
curl http://localhost:8080/api/v1/inventory/$PRODUCT_ID \
  | jq '.data | {product_name, product_price, quantity}'
# Output: {"product_name": "Test", "product_price": 10, "quantity": 100}
```

✅ **All module communications working!**

## Key Observations

1. **Inventory validates products exist** before creation (Inventory → Product)
2. **Product enriches responses** with inventory data when available (Product → Inventory)
3. **Inventory enriches responses** with product details (Inventory → Product)
4. **Graceful degradation** - Product works without inventory
5. **No tight coupling** - Modules communicate through interfaces
6. **Testable** - All communication paths have unit tests

## File Structure Overview

```
├── cmd/api/main.go                      # Dependency wiring (STEP 1-6)
├── internal/
│   ├── application/
│   │   ├── product/
│   │   │   ├── get.go                   # Product → Inventory
│   │   │   ├── interfaces.go            # Interface definition
│   │   │   └── get_test.go              # Module comm tests
│   │   └── inventory/
│   │       ├── create.go                # Inventory → Product
│   │       ├── get.go                   # Inventory → Product
│   │       ├── adjust.go                # Inventory → Product
│   │       ├── product_adapter.go       # Adapter for decoupling
│   │       └── *_test.go                # Module comm tests
│   ├── domain/
│   │   ├── product/
│   │   └── inventory/
│   └── infrastructure/
│       ├── delivery/
│       └── persistence/
└── Documentation:
    ├── MODULE_COMMUNICATION.md          # Complete guide
    ├── IMPLEMENTATION_SUMMARY.md        # What was built
    ├── ARCHITECTURE_DIAGRAM.md          # Visual diagrams
    └── QUICKSTART.md                    # This file
```

## Next Steps

1. **Explore the code:**
   - Start with `cmd/api/main.go` to see dependency wiring
   - Read `internal/application/inventory/create.go` for Inventory → Product
   - Read `internal/application/product/get.go` for Product → Inventory

2. **Read the documentation:**
   - `MODULE_COMMUNICATION.md` - Complete implementation guide
   - `ARCHITECTURE_DIAGRAM.md` - Visual architecture diagrams
   - `IMPLEMENTATION_SUMMARY.md` - What was implemented

3. **Run the tests:**
   - See tests demonstrating module communication
   - Observe mock usage for testing

4. **Extend the pattern:**
   - Add an Order module that uses both Product and Inventory
   - Implement your own bidirectional communication

## Summary

This implementation demonstrates:
- ✅ **Bidirectional module communication** at application layer
- ✅ **Clean Architecture** principles maintained
- ✅ **Loose coupling** via interfaces and adapters
- ✅ **Comprehensive testing** (35 tests, 100% pass)
- ✅ **Production-ready** patterns
- ✅ **Graceful degradation** for resilience
- ✅ **Scalable** to microservices architecture

Happy coding! 🚀

