package inventory

import "context"

// InventoryQueryRepository defines the interface for inventory read operations
// This interface belongs to the domain layer and has no infrastructure dependencies
type InventoryQueryRepository interface {
	// GetByProductID retrieves inventory by product ID
	// Returns nil if inventory is not found
	GetByProductID(ctx context.Context, productID string) (*Inventory, error)
}

