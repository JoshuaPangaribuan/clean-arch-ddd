package inventory

import (
	"context"
	"errors"

	"github.com/JoshuaPangaribuan/clean-arch-ddd/internal/application/product"
	"github.com/JoshuaPangaribuan/clean-arch-ddd/internal/domain/inventory"
)

// GetInventoryUseCase handles the business logic for retrieving inventory
type GetInventoryUseCase struct {
	inventoryRepo  inventory.InventoryRepository
	productUseCase product.ProductUseCaseInterface
}

// NewGetInventoryUseCase creates a new instance of GetInventoryUseCase
// This demonstrates module communication: Inventory → Product
func NewGetInventoryUseCase(
	inventoryRepo inventory.InventoryRepository,
	productUseCase product.ProductUseCaseInterface,
) *GetInventoryUseCase {
	return &GetInventoryUseCase{
		inventoryRepo:  inventoryRepo,
		productUseCase: productUseCase,
	}
}

// Execute performs the get inventory operation
func (uc *GetInventoryUseCase) Execute(ctx context.Context, productID string) (*GetInventoryOutput, error) {
	// Validate input
	if productID == "" {
		return nil, errors.New("product ID is required")
	}

	// Retrieve inventory from repository
	inv, err := uc.inventoryRepo.GetByProductID(ctx, productID)
	if err != nil {
		return nil, err
	}

	// Check if inventory exists
	if inv == nil {
		return nil, inventory.ErrInventoryNotFound
	}

	// MODULE COMMUNICATION: Call Product module to get product details
	productOutput, err := uc.productUseCase.Execute(ctx, productID)
	if err != nil {
		// If product is deleted but inventory still exists, return partial data
		if err.Error() == "product not found" {
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

