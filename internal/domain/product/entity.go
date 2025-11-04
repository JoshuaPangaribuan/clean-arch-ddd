package product

import (
	"time"

	"github.com/JoshuaPangaribuan/clean-arch-ddd/pkg/errors"
)

// Domain errors - using pkg/errors for consistency
var (
	ErrProductNotFound      = errors.New(errors.CodeProductNotFound, "product not found")
	ErrProductAlreadyExists = errors.New(errors.CodeProductAlreadyExists, "product already exists")
)

// Product represents a product entity in the domain
type Product struct {
	id        string
	name      string
	price     Price
	createdAt time.Time
	updatedAt time.Time
}

// NewProduct creates a new Product entity with validation
func NewProduct(id, name string, price Price) (*Product, error) {
	if id == "" {
		return nil, errors.New(errors.CodeInvalidProductID, "product id cannot be empty")
	}
	if name == "" {
		return nil, errors.New(errors.CodeInvalidProductName, "product name cannot be empty")
	}

	now := time.Now()
	return &Product{
		id:        id,
		name:      name,
		price:     price,
		createdAt: now,
		updatedAt: now,
	}, nil
}

// ReconstructProduct reconstructs a Product entity from persistence
// This is used when loading from database
func ReconstructProduct(id, name string, price Price, createdAt, updatedAt time.Time) *Product {
	return &Product{
		id:        id,
		name:      name,
		price:     price,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}
}

// ID returns the product's unique identifier
func (p *Product) ID() string {
	return p.id
}

// Name returns the product's name
func (p *Product) Name() string {
	return p.name
}

// Price returns the product's price
func (p *Product) Price() Price {
	return p.price
}

// CreatedAt returns when the product was created
func (p *Product) CreatedAt() time.Time {
	return p.createdAt
}

// UpdatedAt returns when the product was last updated
func (p *Product) UpdatedAt() time.Time {
	return p.updatedAt
}

// UpdateName updates the product's name with validation
func (p *Product) UpdateName(name string) error {
	if name == "" {
		return errors.New(errors.CodeInvalidProductName, "product name cannot be empty")
	}
	p.name = name
	p.updatedAt = time.Now()
	return nil
}

// UpdatePrice updates the product's price with validation
func (p *Product) UpdatePrice(price Price) error {
	p.price = price
	p.updatedAt = time.Now()
	return nil
}
