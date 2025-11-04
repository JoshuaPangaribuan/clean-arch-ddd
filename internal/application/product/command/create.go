package command

import (
	"context"
	"time"

	"github.com/JoshuaPangaribuan/clean-arch-ddd/internal/domain/product"
	apperrors "github.com/JoshuaPangaribuan/clean-arch-ddd/pkg/errors"
	"github.com/google/uuid"
)

// CreateProductInput represents the input data for creating a product
type CreateProductInput struct {
	Name          string  `json:"name" validate:"required,min=1,max=255"`
	PriceAmount   float64 `json:"price_amount" validate:"required,gte=0"`
	PriceCurrency string  `json:"price_currency" validate:"required,len=3"`
}

// CreateProductOutput represents the output data after creating a product
type CreateProductOutput struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	PriceAmount   float64   `json:"price_amount"`
	PriceCurrency string    `json:"price_currency"`
	CreatedAt     time.Time `json:"created_at"`
}

// CreateProductCommand handles the business logic for creating a product
type CreateProductCommand struct {
	productRepo product.ProductCommandRepository
}

// NewCreateProductCommand creates a new instance of CreateProductCommand
func NewCreateProductCommand(productRepo product.ProductCommandRepository) *CreateProductCommand {
	return &CreateProductCommand{
		productRepo: productRepo,
	}
}

// Execute performs the create product operation
func (c *CreateProductCommand) Execute(ctx context.Context, input CreateProductInput) (*CreateProductOutput, error) {
	// Validate input
	if input.Name == "" {
		return nil, apperrors.New(apperrors.CodeInvalidProductName, "product name is required")
	}

	// Create price value object with validation
	price, err := product.NewPrice(input.PriceAmount, input.PriceCurrency)
	if err != nil {
		return nil, err
	}

	// Generate unique ID for the product
	productID := uuid.New().String()

	// Create product entity with validation
	prod, err := product.NewProduct(productID, input.Name, price)
	if err != nil {
		return nil, err
	}

	// Persist the product
	if err := c.productRepo.Create(ctx, prod); err != nil {
		return nil, apperrors.WrapDatabaseError(err)
	}

	// Return output DTO
	return &CreateProductOutput{
		ID:            prod.ID(),
		Name:          prod.Name(),
		PriceAmount:   prod.Price().Amount(),
		PriceCurrency: prod.Price().Currency(),
		CreatedAt:     prod.CreatedAt(),
	}, nil
}

