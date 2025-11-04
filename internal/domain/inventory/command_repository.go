package inventory

import "context"

// InventoryCommandRepository defines the interface for inventory write operations
// This interface belongs to the domain layer and has no infrastructure dependencies
type InventoryCommandRepository interface {
	// Create stores a new inventory record
	Create(ctx context.Context, inventory *Inventory) error

	// Update updates an existing inventory record
	Update(ctx context.Context, inventory *Inventory) error

	// Delete removes an inventory record by product ID
	Delete(ctx context.Context, productID string) error

	// AdjustStock adjusts the stock quantity for a product
	AdjustStock(ctx context.Context, productID string, adjustment int) error
}

