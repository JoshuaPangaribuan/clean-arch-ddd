package product

import "context"

// ProductQueryRepository defines the interface for product read operations
// This interface belongs to the domain layer and has no infrastructure dependencies
type ProductQueryRepository interface {
	// GetByID retrieves a product by its unique identifier
	// Returns nil if product is not found
	GetByID(ctx context.Context, id string) (*Product, error)

	// List retrieves all products with pagination
	List(ctx context.Context, limit, offset int) ([]*Product, error)
}

