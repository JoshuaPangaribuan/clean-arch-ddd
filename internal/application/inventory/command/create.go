package command

import (
	"context"
	"time"

	"github.com/JoshuaPangaribuan/clean-arch-ddd/internal/application/inventory/query"
	"github.com/JoshuaPangaribuan/clean-arch-ddd/internal/domain/inventory"
	apperrors "github.com/JoshuaPangaribuan/clean-arch-ddd/pkg/errors"
	"github.com/google/uuid"
)

// CreateInventoryInput represents the input for creating inventory
type CreateInventoryInput struct {
	ProductID string `json:"product_id" validate:"required"`
	Quantity  int    `json:"quantity" validate:"required,min=0"`
	Location  string `json:"location"`
}

// CreateInventoryOutput represents the output after creating inventory
type CreateInventoryOutput struct {
	ID                string    `json:"id"`
	ProductID         string    `json:"product_id"`
	ProductName       string    `json:"product_name"`
	Quantity          int       `json:"quantity"`
	ReservedQuantity  int       `json:"reserved_quantity"`
	AvailableQuantity int       `json:"available_quantity"`
	Location          string    `json:"location"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// CreateInventoryCommand handles the business logic for creating inventory
type CreateInventoryCommand struct {
	inventoryCmdRepo inventory.InventoryCommandRepository
	inventoryQueryRepo inventory.InventoryQueryRepository
	productQuery  query.ProductQueryInterface
}

// NewCreateInventoryCommand creates a new instance of CreateInventoryCommand
// This demonstrates module communication: Inventory → Product
func NewCreateInventoryCommand(
	inventoryCmdRepo inventory.InventoryCommandRepository,
	inventoryQueryRepo inventory.InventoryQueryRepository,
	productQuery query.ProductQueryInterface,
) *CreateInventoryCommand {
	return &CreateInventoryCommand{
		inventoryCmdRepo:  inventoryCmdRepo,
		inventoryQueryRepo: inventoryQueryRepo,
		productQuery:  productQuery,
	}
}

// Execute performs the create inventory operation
func (c *CreateInventoryCommand) Execute(ctx context.Context, input CreateInventoryInput) (*CreateInventoryOutput, error) {
	// Validate input
	if input.ProductID == "" {
		return nil, apperrors.New(apperrors.CodeInvalidInput, "product ID is required")
	}
	if input.Quantity < 0 {
		return nil, inventory.ErrInvalidQuantity
	}

	// MODULE COMMUNICATION: Call Product module to verify product exists
	productOutput, err := c.productQuery.Execute(ctx, input.ProductID)
	if err != nil {
		if apperrors.Is(err, apperrors.CodeProductNotFound) {
			return nil, apperrors.New(apperrors.CodeProductNotFound, "cannot create inventory: product not found")
		}
		return nil, err
	}

	// Check if inventory already exists for this product
	existingInventory, err := c.inventoryQueryRepo.GetByProductID(ctx, input.ProductID)
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
	if err := c.inventoryCmdRepo.Create(ctx, inv); err != nil {
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

