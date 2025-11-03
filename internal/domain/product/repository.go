package product

import "context"

// ProductRepository defines the interface for product persistence operations
// This interface belongs to the domain layer and has no infrastructure dependencies
type ProductRepository interface {
	// Create stores a new product
	Create(ctx context.Context, product *Product) error

	// GetByID retrieves a product by its unique identifier
	// Returns nil if product is not found
	GetByID(ctx context.Context, id string) (*Product, error)

	// Update updates an existing product
	Update(ctx context.Context, product *Product) error

	// Delete removes a product by its ID
	Delete(ctx context.Context, id string) error

	// List retrieves all products with pagination
	List(ctx context.Context, limit, offset int) ([]*Product, error)
}

