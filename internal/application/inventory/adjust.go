package inventory

import (
	"context"
	"errors"

	"github.com/JoshuaPangaribuan/clean-arch-ddd/internal/application/product"
	"github.com/JoshuaPangaribuan/clean-arch-ddd/internal/domain/inventory"
)

// AdjustInventoryUseCase handles the business logic for adjusting inventory quantities
type AdjustInventoryUseCase struct {
	inventoryRepo  inventory.InventoryRepository
	productUseCase product.ProductUseCaseInterface
}

// NewAdjustInventoryUseCase creates a new instance of AdjustInventoryUseCase
// This demonstrates module communication: Inventory → Product
func NewAdjustInventoryUseCase(
	inventoryRepo inventory.InventoryRepository,
	productUseCase product.ProductUseCaseInterface,
) *AdjustInventoryUseCase {
	return &AdjustInventoryUseCase{
		inventoryRepo:  inventoryRepo,
		productUseCase: productUseCase,
	}
}

// Execute performs the adjust inventory operation
func (uc *AdjustInventoryUseCase) Execute(ctx context.Context, input AdjustInventoryInput) (*AdjustInventoryOutput, error) {
	// Validate input
	if input.ProductID == "" {
		return nil, errors.New("product ID is required")
	}
	if input.Adjustment == 0 {
		return nil, errors.New("adjustment cannot be zero")
	}

	// MODULE COMMUNICATION: Verify product exists
	productOutput, err := uc.productUseCase.Execute(ctx, input.ProductID)
	if err != nil {
		if err.Error() == "product not found" {
			return nil, errors.New("cannot adjust inventory: product not found")
		}
		return nil, err
	}

	// Retrieve current inventory
	inv, err := uc.inventoryRepo.GetByProductID(ctx, input.ProductID)
	if err != nil {
		return nil, err
	}

	if inv == nil {
		return nil, inventory.ErrInventoryNotFound
	}

	// Apply adjustment to inventory entity (business logic)
	if err := inv.AdjustQuantity(input.Adjustment); err != nil {
		return nil, err
	}

	// Save updated inventory
	if err := uc.inventoryRepo.Update(ctx, inv); err != nil {
		return nil, err
	}

	// Return output DTO
	return &AdjustInventoryOutput{
		ID:                inv.ID(),
		ProductID:         inv.ProductID(),
		ProductName:       productOutput.Name,
		Quantity:          inv.Quantity(),
		ReservedQuantity:  inv.ReservedQuantity(),
		AvailableQuantity: inv.AvailableQuantity(),
		Location:          inv.Location(),
		UpdatedAt:         inv.UpdatedAt(),
	}, nil
}

