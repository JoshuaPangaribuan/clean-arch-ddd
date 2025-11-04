package product

// ProductRepository is deprecated. Use ProductCommandRepository and ProductQueryRepository instead.
// This interface is kept for backward compatibility during transition.
// Deprecated: Use ProductCommandRepository for write operations and ProductQueryRepository for read operations.
type ProductRepository interface {
	ProductCommandRepository
	ProductQueryRepository
}
