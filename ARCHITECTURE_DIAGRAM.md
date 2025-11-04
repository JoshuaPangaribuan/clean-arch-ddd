# Module Communication Architecture

## System Overview

```mermaid
flowchart TB
    subgraph HTTP["HTTP Layer (Gin)"]
        ProductsAPI["/api/v1/products<br/>POST /<br/>GET /:id"]
        InventoryAPI["/api/v1/inventory<br/>POST /<br/>GET /:productId<br/>PATCH /adjust"]
    end
    
    subgraph Delivery["Delivery Layer (Handlers)"]
        ProductHandler["ProductHandler<br/>- Create()<br/>- Get()"]
        InventoryHandler["InventoryHandler<br/>- Create()<br/>- Get()<br/>- Adjust()"]
    end
    
    subgraph Application["Application Layer (Use Cases)"]
        direction LR
        subgraph ProductModule["Product Module"]
            CreateProduct["CreateProductUseCase"]
            GetProduct["GetProductUseCase"]
        end
        
        subgraph InventoryModule["Inventory Module"]
            CreateInventory["CreateInventoryUseCase"]
            GetInventory["GetInventoryUseCase"]
            AdjustInventory["AdjustInventoryUseCase"]
        end
        
        GetProduct -.->|injects| InventoryModule
        InventoryModule -.->|validates product| ProductModule
        GetProduct -.->|enriches with inventory| InventoryModule
    end
    
    subgraph Domain["Domain Layer (Entities)"]
        direction LR
        ProductEntity["Product Entity<br/>- ID, Name, Price<br/>- UpdateName()<br/>- UpdatePrice()"]
        InventoryEntity["Inventory Entity<br/>- ID, ProductID, Quantity<br/>- Reserve()<br/>- Release()<br/>- AdjustQuantity()<br/>- AvailableQuantity()"]
        ProductRepoI["ProductRepository<br/>(interface)"]
        InventoryRepoI["InventoryRepository<br/>(interface)"]
    end
    
    subgraph Infrastructure["Infrastructure Layer (Persistence)"]
        direction LR
        ProductRepoImpl["ProductRepositoryImpl<br/>(implements interface)"]
        InventoryRepoImpl["InventoryRepositoryImpl<br/>(implements interface)"]
    end
    
    subgraph Database["Database (PostgreSQL)"]
        direction LR
        ProductsTable["products table<br/>- id (PK)<br/>- name<br/>- price_amount<br/>- price_currency<br/>- created_at<br/>- updated_at"]
        InventoryTable["inventory table<br/>- id (PK)<br/>- product_id (FK)<br/>- quantity<br/>- reserved_quantity<br/>- location<br/>- created_at<br/>- updated_at"]
        ProductsTable <-->|Foreign Key| InventoryTable
    end
    
    HTTP --> Delivery
    Delivery --> Application
    Application --> Domain
    Domain --> Infrastructure
    Infrastructure --> Database
    
    ProductsAPI --> ProductHandler
    InventoryAPI --> InventoryHandler
    
    ProductHandler --> CreateProduct
    ProductHandler --> GetProduct
    InventoryHandler --> CreateInventory
    InventoryHandler --> GetInventory
    InventoryHandler --> AdjustInventory
    
    CreateProduct --> ProductRepoI
    GetProduct --> ProductRepoI
    CreateInventory --> InventoryRepoI
    GetInventory --> InventoryRepoI
    AdjustInventory --> InventoryRepoI
    
    ProductRepoI --> ProductRepoImpl
    InventoryRepoI --> InventoryRepoImpl
    
    ProductRepoImpl --> ProductsTable
    InventoryRepoImpl --> InventoryTable
    
    style HTTP fill:#e1f5ff
    style Delivery fill:#fff4e1
    style Application fill:#e8f5e9
    style Domain fill:#f3e5f5
    style Infrastructure fill:#fff9c4
    style Database fill:#e0f2f1
```

## Module Communication Flows

### Flow 1: Create Inventory (Inventory → Product Validation)

```mermaid
sequenceDiagram
    participant Client
    participant InventoryHandler
    participant CreateInventoryUseCase
    participant ValidateInput
    participant GetProductUseCase as GetProductUseCase<br/>(MODULE COMM)
    participant CheckExisting
    participant CreateInventory
    
    Client->>InventoryHandler: POST /inventory
    InventoryHandler->>CreateInventoryUseCase: Execute(request)
    
    par Parallel Processing
        CreateInventoryUseCase->>ValidateInput: Validate input
        ValidateInput-->>CreateInventoryUseCase: Valid
    and
        CreateInventoryUseCase->>GetProductUseCase: Execute(productID)
        GetProductUseCase-->>CreateInventoryUseCase: Product exists?
        alt Product exists
            Note over CreateInventoryUseCase: Continue
        else Product not found
            CreateInventoryUseCase-->>Client: Error: Product not found
        end
    and
        CreateInventoryUseCase->>CheckExisting: Check existing inventory
        CheckExisting-->>CreateInventoryUseCase: Not exists
    end
    
    CreateInventoryUseCase->>CreateInventory: Create inventory record
    CreateInventory-->>CreateInventoryUseCase: Success
    CreateInventoryUseCase-->>InventoryHandler: Inventory created
    InventoryHandler-->>Client: 201 Created
```

### Flow 2: Get Product (Product → Inventory Enrichment)

```mermaid
sequenceDiagram
    participant Client
    participant ProductHandler
    participant GetProductUseCase
    participant ProductRepo
    participant InventoryModule as Inventory Module<br/>(MODULE COMM)
    participant GetInventory
    participant ReturnEnriched
    
    Client->>ProductHandler: GET /products/:id
    ProductHandler->>GetProductUseCase: Execute(productID)
    GetProductUseCase->>ProductRepo: GetProductByID(id)
    ProductRepo-->>GetProductUseCase: Product data
    
    par Parallel Processing
        GetProductUseCase->>GetProductUseCase: Build base output<br/>(Product Data)
    and
        GetProductUseCase->>InventoryModule: Check if available
        InventoryModule-->>GetProductUseCase: Available?
        alt Inventory available
            GetProductUseCase->>GetInventory: Execute(productID)
            GetInventory-->>GetProductUseCase: Inventory data
            GetProductUseCase->>GetProductUseCase: Enrich output<br/>- StockQuantity<br/>- Available
        else No inventory
            Note over GetProductUseCase: has_inventory: false
        end
    end
    
    GetProductUseCase->>ReturnEnriched: Build enriched response
    ReturnEnriched-->>GetProductUseCase: Product + Stock
    GetProductUseCase-->>ProductHandler: Enriched product data
    ProductHandler-->>Client: 200 OK
```

### Flow 3: Get Inventory (Inventory → Product Enrichment)

```mermaid
sequenceDiagram
    participant Client
    participant InventoryHandler
    participant GetInventoryUseCase
    participant InventoryRepo
    participant ProductModule as GetProductUseCase<br/>(MODULE COMM)
    participant BuildOutput
    
    Client->>InventoryHandler: GET /inventory/:productId
    InventoryHandler->>GetInventoryUseCase: Execute(productID)
    GetInventoryUseCase->>InventoryRepo: GetInventoryByProductID(id)
    
    alt Inventory exists
        InventoryRepo-->>GetInventoryUseCase: Inventory data
        
        par Parallel Processing
            GetInventoryUseCase->>GetInventoryUseCase: Inventory exists?
        and
            GetInventoryUseCase->>ProductModule: Execute(productID)
            ProductModule-->>GetInventoryUseCase: Product data?
            alt Product exists
                Note over GetInventoryUseCase: Add product data
            else Product deleted
                Note over GetInventoryUseCase: Graceful fallback<br/>"Product Deleted"
            end
        end
        
        GetInventoryUseCase->>BuildOutput: Build output with Product Details
        BuildOutput-->>GetInventoryUseCase: Inventory + Product Info
        GetInventoryUseCase-->>InventoryHandler: Enriched inventory data
        InventoryHandler-->>Client: 200 OK
    else Inventory not found
        InventoryRepo-->>GetInventoryUseCase: Not found
        GetInventoryUseCase-->>InventoryHandler: Error
        InventoryHandler-->>Client: 404 Not Found
    end
```

## Dependency Injection Flow

```mermaid
flowchart TD
    Start([main.go]) --> Step1[STEP 1: Initialize Repositories]
    
    Step1 --> Step1Code["productRepo := NewProductRepository(db)<br/>inventoryRepo := NewInventoryRepository(db)"]
    
    Step1Code --> Step2[STEP 2: Create Product Use Cases]
    
    Step2 --> Step2Code["createProductUseCase := NewCreateProductUseCase(productRepo)<br/><br/>getProductUseCaseBasic := NewGetProductUseCase(productRepo)<br/>// ← No inventory integration yet"]
    
    Step2Code --> Step3[STEP 3: Create Inventory Use Cases]
    
    Step3 --> Step3Code["createInventoryUseCase := NewCreateInventoryUseCase(<br/>  inventoryRepo,<br/>  getProductUseCaseBasic<br/>  // ← Inventory → Product<br/>)<br/><br/>getInventoryUseCase := NewGetInventoryUseCase(<br/>  inventoryRepo,<br/>  getProductUseCaseBasic<br/>  // ← Inventory → Product<br/>)"]
    
    Step3Code --> Step4[STEP 4: Create Adapter]
    
    Step4 --> Step4Code["inventoryAdapter := NewProductInventoryAdapter(<br/>  getInventoryUseCase<br/>)"]
    
    Step4Code --> Step5[STEP 5: Re-create Product Get Use Case]
    
    Step5 --> Step5Code["getProductUseCase := NewGetProductUseCaseWithInventory(<br/>  productRepo,<br/>  inventoryAdapter<br/>  // ← Product → Inventory<br/>)"]
    
    Step5Code --> Step6[STEP 6: Initialize Handlers]
    
    Step6 --> Step6Code["productHandler := NewProductHandler(<br/>  createProductUseCase,<br/>  getProductUseCase<br/>  // ← Uses version with inventory<br/>)<br/><br/>inventoryHandler := NewInventoryHandler(<br/>  createInventoryUseCase,<br/>  getInventoryUseCase,<br/>  adjustInventoryUseCase<br/>)"]
    
    Step6Code --> End([Application Ready])
    
    style Start fill:#e3f2fd
    style Step1 fill:#fff3e0
    style Step2 fill:#fff3e0
    style Step3 fill:#fff3e0
    style Step4 fill:#fff3e0
    style Step5 fill:#fff3e0
    style Step6 fill:#fff3e0
    style End fill:#e8f5e9
```

## Key Communication Interfaces

### Product Module Interface (for Inventory)
```go
type ProductUseCaseInterface interface {
    Execute(ctx context.Context, productID string) (*GetProductOutput, error)
}
```

### Inventory Module Interface (for Product)
```go
type InventoryUseCaseInterface interface {
    Execute(ctx context.Context, productID string) (InventoryData, error)
}

type InventoryData interface {
    GetQuantity() int
    GetAvailableQuantity() int
}
```

## Benefits Visualization

```mermaid
mindmap
  root((Design Benefits))
    Loose Coupling
      Modules depend on interfaces
      Not implementations
    Testability
      Each module tested in isolation
      Uses mocks
    Graceful Degradation
      System works if module unavailable
      Resilient design
    Clean Architecture
      Dependencies flow inward
      Toward domain
    SOLID Principles
      Single responsibility
      Interface segregation
    Scalability
      Ready for microservices
      Monolith to microservices
```

## Testing Strategy

```mermaid
flowchart TD
    subgraph Pyramid["Test Pyramid"]
        direction TB
        
        subgraph Integration["Integration Tests (Future: E2E tests)"]
            IntegrationLabel["Integration Tests<br/>Future: E2E tests"]
        end
        
        subgraph UseCase["Use Case Tests (35 tests - ALL PASS)"]
            UseCaseLabel["Use Case Tests with Mocks<br/>- CreateInventory<br/>- GetInventory<br/>- AdjustInventory<br/>- GetProduct w/ Inventory"]
        end
        
        subgraph Domain["Domain Entity Tests (15 tests - ALL PASS)"]
            DomainLabel["Domain Entity Tests<br/>- Product validation<br/>- Price validation<br/>- Business rules"]
        end
        
        Integration --> UseCase
        UseCase --> Domain
        
        Note1["Unit Tests (Fast, No I/O)"] -.-> Domain
        Note2["Integration Tests (With DB)"] -.-> Integration
    end
    
    style Integration fill:#fff3e0
    style UseCase fill:#e8f5e9
    style Domain fill:#e3f2fd
```

This architecture demonstrates production-ready module communication patterns that can scale from monolith to microservices while maintaining clean boundaries and testability.
