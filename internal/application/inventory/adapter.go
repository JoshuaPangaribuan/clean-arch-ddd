package inventory

// InventoryAdapter adapts GetInventoryOutput to work with Product module's interface
type InventoryAdapter struct {
	output *GetInventoryOutput
}

// NewInventoryAdapter creates a new adapter
func NewInventoryAdapter(output *GetInventoryOutput) *InventoryAdapter {
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

