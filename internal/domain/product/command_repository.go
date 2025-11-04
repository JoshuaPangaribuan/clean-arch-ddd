package product

import "context"

// ProductCommandRepository defines the interface for product write operations
// This interface belongs to the domain layer and has no infrastructure dependencies
type ProductCommandRepository interface {
	// Create stores a new product
	Create(ctx context.Context, product *Product) error

	// Update updates an existing product
	Update(ctx context.Context, product *Product) error

	// Delete removes a product by its ID
	Delete(ctx context.Context, id string) error
}
