package inventory

import (
	"context"

	"github.com/JoshuaPangaribuan/clean-arch-ddd/internal/application/product"
)

// ProductInventoryAdapter adapts GetInventoryUseCase to implement product.InventoryUseCaseInterface
// This enables Product module to call Inventory module
type ProductInventoryAdapter struct {
	inventoryUseCase *GetInventoryUseCase
}

// NewProductInventoryAdapter creates a new adapter
func NewProductInventoryAdapter(inventoryUseCase *GetInventoryUseCase) *ProductInventoryAdapter {
	return &ProductInventoryAdapter{
		inventoryUseCase: inventoryUseCase,
	}
}

// Execute calls the inventory use case and adapts the result
func (a *ProductInventoryAdapter) Execute(ctx context.Context, productID string) (product.InventoryData, error) {
	output, err := a.inventoryUseCase.Execute(ctx, productID)
	if err != nil {
		return nil, err
	}
	return NewInventoryAdapter(output), nil
}

