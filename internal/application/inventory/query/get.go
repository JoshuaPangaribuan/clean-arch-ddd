package query

import (
	"context"
	"time"

	"github.com/JoshuaPangaribuan/clean-arch-ddd/internal/domain/inventory"
	apperrors "github.com/JoshuaPangaribuan/clean-arch-ddd/pkg/errors"
)

// GetInventoryQuery handles the business logic for retrieving inventory
type GetInventoryQuery struct {
	inventoryRepo inventory.InventoryQueryRepository
	productQuery  ProductQueryInterface
}

// NewGetInventoryQuery creates a new instance of GetInventoryQuery
// This demonstrates module communication: Inventory → Product
func NewGetInventoryQuery(
	inventoryRepo inventory.InventoryQueryRepository,
	productQuery ProductQueryInterface,
) *GetInventoryQuery {
	return &GetInventoryQuery{
		inventoryRepo: inventoryRepo,
		productQuery:  productQuery,
	}
}

// GetInventoryOutput represents the output for getting inventory
type GetInventoryOutput struct {
	ID                string    `json:"id"`
	ProductID         string    `json:"product_id"`
	ProductName       string    `json:"product_name"`
	ProductPrice      float64   `json:"product_price"`
	ProductCurrency   string    `json:"product_currency"`
	Quantity          int       `json:"quantity"`
	ReservedQuantity  int       `json:"reserved_quantity"`
	AvailableQuantity int       `json:"available_quantity"`
	Location          string    `json:"location"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// Execute performs the get inventory operation
func (q *GetInventoryQuery) Execute(ctx context.Context, productID string) (*GetInventoryOutput, error) {
	// Validate input
	if productID == "" {
		return nil, apperrors.New(apperrors.CodeInvalidInput, "product ID is required")
	}

	// Retrieve inventory from repository
	inv, err := q.inventoryRepo.GetByProductID(ctx, productID)
	if err != nil {
		return nil, apperrors.WrapDatabaseError(err)
	}

	// Check if inventory exists
	if inv == nil {
		return nil, inventory.ErrInventoryNotFound
	}

	// MODULE COMMUNICATION: Call Product module to get product details
	productOutput, err := q.productQuery.Execute(ctx, productID)
	if err != nil {
		// If product is deleted but inventory still exists, return partial data
		if apperrors.Is(err, apperrors.CodeProductNotFound) {
			return &GetInventoryOutput{
				ID:                inv.ID(),
				ProductID:         inv.ProductID(),
				ProductName:       "Unknown (Product Deleted)",
				ProductPrice:      0,
				ProductCurrency:   "",
				Quantity:          inv.Quantity(),
				ReservedQuantity:  inv.ReservedQuantity(),
				AvailableQuantity: inv.AvailableQuantity(),
				Location:          inv.Location(),
				CreatedAt:         inv.CreatedAt(),
				UpdatedAt:         inv.UpdatedAt(),
			}, nil
		}
		return nil, err
	}

	// Return output DTO enriched with product information
	return &GetInventoryOutput{
		ID:                inv.ID(),
		ProductID:         inv.ProductID(),
		ProductName:       productOutput.Name,
		ProductPrice:      productOutput.PriceAmount,
		ProductCurrency:   productOutput.PriceCurrency,
		Quantity:          inv.Quantity(),
		ReservedQuantity:  inv.ReservedQuantity(),
		AvailableQuantity: inv.AvailableQuantity(),
		Location:          inv.Location(),
		CreatedAt:         inv.CreatedAt(),
		UpdatedAt:         inv.UpdatedAt(),
	}, nil
}

