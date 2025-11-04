package product

import "context"

// ProductUseCaseInterface defines the interface for getting a product
// This allows other modules to depend on the interface rather than concrete implementation
type ProductUseCaseInterface interface {
	Execute(ctx context.Context, productID string) (*GetProductOutput, error)
}

