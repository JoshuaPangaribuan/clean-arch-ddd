package inventory

import (
	"context"

	"github.com/JoshuaPangaribuan/clean-arch-ddd/internal/application/product"
	"github.com/JoshuaPangaribuan/clean-arch-ddd/internal/domain/inventory"
	apperrors "github.com/JoshuaPangaribuan/clean-arch-ddd/pkg/errors"
	"github.com/google/uuid"
)

// CreateInventoryUseCase handles the business logic for creating inventory
type CreateInventoryUseCase struct {
	inventoryRepo inventory.InventoryRepository
	productUseCase product.ProductUseCaseInterface
}

// NewCreateInventoryUseCase creates a new instance of CreateInventoryUseCase
// This demonstrates module communication: Inventory → Product
func NewCreateInventoryUseCase(
	inventoryRepo inventory.InventoryRepository,
	productUseCase product.ProductUseCaseInterface,
) *CreateInventoryUseCase {
	return &CreateInventoryUseCase{
		inventoryRepo:  inventoryRepo,
		productUseCase: productUseCase,
	}
}

// Execute performs the create inventory operation
func (uc *CreateInventoryUseCase) Execute(ctx context.Context, input CreateInventoryInput) (*CreateInventoryOutput, error) {
	// Validate input
	if input.ProductID == "" {
		return nil, apperrors.New(apperrors.CodeInvalidInput, "product ID is required")
	}
	if input.Quantity < 0 {
		return nil, inventory.ErrInvalidQuantity
	}

	// MODULE COMMUNICATION: Call Product module to verify product exists
	productOutput, err := uc.productUseCase.Execute(ctx, input.ProductID)
	if err != nil {
		if apperrors.Is(err, apperrors.CodeProductNotFound) {
			return nil, apperrors.New(apperrors.CodeProductNotFound, "cannot create inventory: product not found")
		}
		return nil, err
	}

	// Check if inventory already exists for this product
	existingInventory, err := uc.inventoryRepo.GetByProductID(ctx, input.ProductID)
	if err != nil {
		return nil, apperrors.WrapDatabaseError(err)
	}
	if existingInventory != nil {
		return nil, inventory.ErrInventoryExists
	}

	// Create new inventory entity
	inv, err := inventory.NewInventory(
		uuid.New().String(),
		input.ProductID,
		input.Quantity,
		input.Location,
	)
	if err != nil {
		return nil, err
	}

	// Save to repository
	if err := uc.inventoryRepo.Create(ctx, inv); err != nil {
		return nil, apperrors.WrapDatabaseError(err)
	}

	// Return output DTO with product information
	return &CreateInventoryOutput{
		ID:                inv.ID(),
		ProductID:         inv.ProductID(),
		ProductName:       productOutput.Name,
		Quantity:          inv.Quantity(),
		ReservedQuantity:  inv.ReservedQuantity(),
		AvailableQuantity: inv.AvailableQuantity(),
		Location:          inv.Location(),
		CreatedAt:         inv.CreatedAt(),
		UpdatedAt:         inv.UpdatedAt(),
	}, nil
}

