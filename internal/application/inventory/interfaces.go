package inventory

import "context"

// InventoryUseCaseInterface defines the interface for getting inventory
// This allows other modules to depend on the interface rather than concrete implementation
type InventoryUseCaseInterface interface {
	Execute(ctx context.Context, productID string) (*GetInventoryOutput, error)
}

