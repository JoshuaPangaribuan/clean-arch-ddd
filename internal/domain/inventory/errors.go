package inventory

import "errors"

// Domain errors
var (
	ErrInventoryNotFound   = errors.New("inventory not found")
	ErrInventoryExists     = errors.New("inventory already exists for this product")
	ErrInsufficientStock   = errors.New("insufficient stock available")
	ErrInvalidQuantity     = errors.New("quantity must be non-negative")
	ErrInvalidAdjustment   = errors.New("invalid adjustment amount")
)

