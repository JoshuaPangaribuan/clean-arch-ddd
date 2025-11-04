package query

import (
	"context"
	"time"

	"github.com/JoshuaPangaribuan/clean-arch-ddd/internal/domain/product"
	apperrors "github.com/JoshuaPangaribuan/clean-arch-ddd/pkg/errors"
)

// GetProductQuery handles the business logic for retrieving a product
type GetProductQuery struct {
	productRepo    product.ProductQueryRepository
	inventoryQuery InventoryQueryInterface
}

// InventoryQueryInterface defines the interface for inventory query operations
// This allows Product module to communicate with Inventory module
type InventoryQueryInterface interface {
	Execute(ctx context.Context, productID string) (InventoryData, error)
}

// NewGetProductQuery creates a new instance of GetProductQuery without inventory
func NewGetProductQuery(productRepo product.ProductQueryRepository) *GetProductQuery {
	return &GetProductQuery{
		productRepo:    productRepo,
		inventoryQuery: nil,
	}
}

// NewGetProductQueryWithInventory creates a new instance with inventory integration
// This demonstrates bidirectional module communication: Product → Inventory
func NewGetProductQueryWithInventory(
	productRepo product.ProductQueryRepository,
	inventoryQuery InventoryQueryInterface,
) *GetProductQuery {
	return &GetProductQuery{
		productRepo:    productRepo,
		inventoryQuery: inventoryQuery,
	}
}

// GetProductOutput represents the output data when retrieving a product
type GetProductOutput struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	PriceAmount   float64   `json:"price_amount"`
	PriceCurrency string    `json:"price_currency"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	// Inventory fields (optional, populated when inventory service is available)
	HasInventory      bool `json:"has_inventory,omitempty"`
	StockQuantity     int  `json:"stock_quantity,omitempty"`
	AvailableQuantity int  `json:"available_quantity,omitempty"`
}

// Execute performs the get product operation
func (q *GetProductQuery) Execute(ctx context.Context, productID string) (*GetProductOutput, error) {
	// Validate input
	if productID == "" {
		return nil, apperrors.New(apperrors.CodeInvalidProductID, "product ID is required")
	}

	// Retrieve product from repository
	prod, err := q.productRepo.GetByID(ctx, productID)
	if err != nil {
		return nil, apperrors.WrapDatabaseError(err)
	}

	// Check if product exists
	if prod == nil {
		return nil, product.ErrProductNotFound
	}

	// Build base output DTO
	output := &GetProductOutput{
		ID:            prod.ID(),
		Name:          prod.Name(),
		PriceAmount:   prod.Price().Amount(),
		PriceCurrency: prod.Price().Currency(),
		CreatedAt:     prod.CreatedAt(),
		UpdatedAt:     prod.UpdatedAt(),
		HasInventory:  false,
	}

	// MODULE COMMUNICATION: Enrich with inventory data if available
	if q.inventoryQuery != nil {
		inventoryData, err := q.inventoryQuery.Execute(ctx, productID)
		if err == nil {
			// Successfully retrieved inventory
			output.HasInventory = true
			output.StockQuantity = inventoryData.GetQuantity()
			output.AvailableQuantity = inventoryData.GetAvailableQuantity()
		}
		// Gracefully handle inventory not found - product data is still valid
	}

	return output, nil
}

