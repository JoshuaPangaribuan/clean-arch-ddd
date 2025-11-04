package command

import (
	"context"
	"time"

	"github.com/JoshuaPangaribuan/clean-arch-ddd/internal/application/inventory/query"
	"github.com/JoshuaPangaribuan/clean-arch-ddd/internal/domain/inventory"
	apperrors "github.com/JoshuaPangaribuan/clean-arch-ddd/pkg/errors"
)

// AdjustInventoryInput represents the input for adjusting inventory
type AdjustInventoryInput struct {
	ProductID  string `json:"product_id" validate:"required"`
	Adjustment int    `json:"adjustment" validate:"required"`
	Reason     string `json:"reason"`
}

// AdjustInventoryOutput represents the output after adjusting inventory
type AdjustInventoryOutput struct {
	ID                string    `json:"id"`
	ProductID         string    `json:"product_id"`
	ProductName       string    `json:"product_name"`
	Quantity          int       `json:"quantity"`
	ReservedQuantity  int       `json:"reserved_quantity"`
	AvailableQuantity int       `json:"available_quantity"`
	Location          string    `json:"location"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// AdjustInventoryCommand handles the business logic for adjusting inventory quantities
type AdjustInventoryCommand struct {
	inventoryCmdRepo  inventory.InventoryCommandRepository
	inventoryQueryRepo inventory.InventoryQueryRepository
	productQuery      query.ProductQueryInterface
}

// NewAdjustInventoryCommand creates a new instance of AdjustInventoryCommand
// This demonstrates module communication: Inventory → Product
func NewAdjustInventoryCommand(
	inventoryCmdRepo inventory.InventoryCommandRepository,
	inventoryQueryRepo inventory.InventoryQueryRepository,
	productQuery query.ProductQueryInterface,
) *AdjustInventoryCommand {
	return &AdjustInventoryCommand{
		inventoryCmdRepo:  inventoryCmdRepo,
		inventoryQueryRepo: inventoryQueryRepo,
		productQuery:      productQuery,
	}
}

// Execute performs the adjust inventory operation
func (c *AdjustInventoryCommand) Execute(ctx context.Context, input AdjustInventoryInput) (*AdjustInventoryOutput, error) {
	// Validate input
	if input.ProductID == "" {
		return nil, apperrors.New(apperrors.CodeInvalidInput, "product ID is required")
	}
	if input.Adjustment == 0 {
		return nil, apperrors.New(apperrors.CodeInvalidAdjustment, "adjustment cannot be zero")
	}

	// MODULE COMMUNICATION: Verify product exists
	productOutput, err := c.productQuery.Execute(ctx, input.ProductID)
	if err != nil {
		if apperrors.Is(err, apperrors.CodeProductNotFound) {
			return nil, apperrors.New(apperrors.CodeProductNotFound, "cannot adjust inventory: product not found")
		}
		return nil, err
	}

	// First, check if inventory exists and validate business rules
	inv, err := c.inventoryQueryRepo.GetByProductID(ctx, input.ProductID)
	if err != nil {
		return nil, apperrors.WrapDatabaseError(err)
	}

	if inv == nil {
		return nil, inventory.ErrInventoryNotFound
	}

	// Validate adjustment against business rules (using in-memory entity)
	// This checks if the adjustment would result in invalid state
	if err := inv.AdjustQuantity(input.Adjustment); err != nil {
		return nil, err
	}

	// Use atomic database operation to prevent race conditions
	// AdjustStock performs: UPDATE inventory SET quantity = quantity + $adjustment
	// This is atomic and thread-safe at the database level
	if err := c.inventoryCmdRepo.AdjustStock(ctx, input.ProductID, input.Adjustment); err != nil {
		return nil, apperrors.WrapDatabaseError(err)
	}

	// Retrieve updated inventory to return accurate data
	updatedInv, err := c.inventoryQueryRepo.GetByProductID(ctx, input.ProductID)
	if err != nil {
		return nil, apperrors.WrapDatabaseError(err)
	}

	// Return output DTO
	return &AdjustInventoryOutput{
		ID:                updatedInv.ID(),
		ProductID:         updatedInv.ProductID(),
		ProductName:       productOutput.Name,
		Quantity:          updatedInv.Quantity(),
		ReservedQuantity:  updatedInv.ReservedQuantity(),
		AvailableQuantity: updatedInv.AvailableQuantity(),
		Location:          updatedInv.Location(),
		UpdatedAt:         updatedInv.UpdatedAt(),
	}, nil
}
