package inventory

// InventoryRepository is deprecated. Use InventoryCommandRepository and InventoryQueryRepository instead.
// This interface is kept for backward compatibility during transition.
// Deprecated: Use InventoryCommandRepository for write operations and InventoryQueryRepository for read operations.
type InventoryRepository interface {
	InventoryCommandRepository
	InventoryQueryRepository
}

