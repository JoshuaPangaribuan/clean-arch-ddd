package inventory

import "github.com/JoshuaPangaribuan/clean-arch-ddd/pkg/errors"

// Domain errors - using pkg/errors for consistency
var (
	ErrInventoryNotFound = errors.New(errors.CodeInventoryNotFound, "inventory not found")
	ErrInventoryExists   = errors.New(errors.CodeInventoryExists, "inventory already exists for this product")
	ErrInsufficientStock = errors.New(errors.CodeInsufficientStock, "insufficient stock available")
	ErrInvalidQuantity   = errors.New(errors.CodeInvalidQuantity, "quantity must be non-negative")
	ErrInvalidAdjustment = errors.New(errors.CodeInvalidAdjustment, "invalid adjustment amount")
)
