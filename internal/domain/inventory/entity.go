package inventory

import (
	"errors"
	"time"
)

// Inventory represents an inventory entity in the domain
type Inventory struct {
	id               string
	productID        string
	quantity         int
	reservedQuantity int
	location         string
	createdAt        time.Time
	updatedAt        time.Time
}

// NewInventory creates a new Inventory entity with validation
func NewInventory(id, productID string, quantity int, location string) (*Inventory, error) {
	if id == "" {
		return nil, errors.New("inventory id cannot be empty")
	}
	if productID == "" {
		return nil, errors.New("product id cannot be empty")
	}
	if quantity < 0 {
		return nil, ErrInvalidQuantity
	}

	now := time.Now()
	return &Inventory{
		id:               id,
		productID:        productID,
		quantity:         quantity,
		reservedQuantity: 0,
		location:         location,
		createdAt:        now,
		updatedAt:        now,
	}, nil
}

// ReconstructInventory reconstructs an Inventory entity from persistence
// This is used when loading from database
func ReconstructInventory(id, productID string, quantity, reservedQuantity int, location string, createdAt, updatedAt time.Time) *Inventory {
	return &Inventory{
		id:               id,
		productID:        productID,
		quantity:         quantity,
		reservedQuantity: reservedQuantity,
		location:         location,
		createdAt:        createdAt,
		updatedAt:        updatedAt,
	}
}

// ID returns the inventory's unique identifier
func (i *Inventory) ID() string {
	return i.id
}

// ProductID returns the product ID this inventory is for
func (i *Inventory) ProductID() string {
	return i.productID
}

// Quantity returns the total quantity
func (i *Inventory) Quantity() int {
	return i.quantity
}

// ReservedQuantity returns the reserved quantity
func (i *Inventory) ReservedQuantity() int {
	return i.reservedQuantity
}

// Location returns the storage location
func (i *Inventory) Location() string {
	return i.location
}

// CreatedAt returns when the inventory was created
func (i *Inventory) CreatedAt() time.Time {
	return i.createdAt
}

// UpdatedAt returns when the inventory was last updated
func (i *Inventory) UpdatedAt() time.Time {
	return i.updatedAt
}

// AvailableQuantity returns the quantity available for reservation/sale
func (i *Inventory) AvailableQuantity() int {
	return i.quantity - i.reservedQuantity
}

// Reserve reserves a quantity of inventory
func (i *Inventory) Reserve(quantity int) error {
	if quantity <= 0 {
		return ErrInvalidQuantity
	}
	if i.AvailableQuantity() < quantity {
		return ErrInsufficientStock
	}
	i.reservedQuantity += quantity
	i.updatedAt = time.Now()
	return nil
}

// Release releases a reserved quantity back to available stock
func (i *Inventory) Release(quantity int) error {
	if quantity <= 0 {
		return ErrInvalidQuantity
	}
	if i.reservedQuantity < quantity {
		return errors.New("cannot release more than reserved quantity")
	}
	i.reservedQuantity -= quantity
	i.updatedAt = time.Now()
	return nil
}

// AdjustQuantity adjusts the total quantity (positive for increase, negative for decrease)
func (i *Inventory) AdjustQuantity(adjustment int) error {
	newQuantity := i.quantity + adjustment
	if newQuantity < 0 {
		return ErrInvalidQuantity
	}
	// Ensure we don't go below reserved quantity
	if newQuantity < i.reservedQuantity {
		return errors.New("cannot adjust quantity below reserved amount")
	}
	i.quantity = newQuantity
	i.updatedAt = time.Now()
	return nil
}

// UpdateLocation updates the storage location
func (i *Inventory) UpdateLocation(location string) {
	i.location = location
	i.updatedAt = time.Now()
}

