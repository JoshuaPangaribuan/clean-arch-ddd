package query

import (
	"context"
)

// InventoryData represents minimal inventory data needed by Product module
type InventoryData interface {
	GetQuantity() int
	GetAvailableQuantity() int
}

// InventoryOutput represents inventory output data
// This type is defined here to avoid circular imports
type InventoryOutput struct {
	Quantity          int
	AvailableQuantity int
}

// InventoryAdapter adapts InventoryOutput to work with Product module's interface
type InventoryAdapter struct {
	output *InventoryOutput
}

// NewInventoryAdapter creates a new adapter
func NewInventoryAdapter(output *InventoryOutput) *InventoryAdapter {
	return &InventoryAdapter{output: output}
}

// GetQuantity returns the total quantity
func (a *InventoryAdapter) GetQuantity() int {
	if a.output == nil {
		return 0
	}
	return a.output.Quantity
}

// GetAvailableQuantity returns the available quantity
func (a *InventoryAdapter) GetAvailableQuantity() int {
	if a.output == nil {
		return 0
	}
	return a.output.AvailableQuantity
}

// InventoryQueryFunc is a function type that executes an inventory query
type InventoryQueryFunc func(ctx context.Context, productID string) (*InventoryOutput, error)

// ProductInventoryAdapter adapts an inventory query function to implement InventoryQueryInterface
// This enables Product module to call Inventory module without circular imports
type ProductInventoryAdapter struct {
	inventoryQuery InventoryQueryFunc
}

// NewProductInventoryAdapter creates a new adapter from a function
func NewProductInventoryAdapter(inventoryQuery InventoryQueryFunc) *ProductInventoryAdapter {
	return &ProductInventoryAdapter{
		inventoryQuery: inventoryQuery,
	}
}

// Execute calls the inventory query function and adapts the result
func (a *ProductInventoryAdapter) Execute(ctx context.Context, productID string) (InventoryData, error) {
	output, err := a.inventoryQuery(ctx, productID)
	if err != nil {
		return nil, err
	}
	return NewInventoryAdapter(output), nil
}
