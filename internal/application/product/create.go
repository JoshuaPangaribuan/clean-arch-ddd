package product

import (
	"context"
	"errors"

	"github.com/JoshuaPangaribuan/clean-arch-ddd/internal/domain/product"
	"github.com/google/uuid"
)

// CreateProductUseCase handles the business logic for creating a product
type CreateProductUseCase struct {
	productRepo product.ProductRepository
}

// NewCreateProductUseCase creates a new instance of CreateProductUseCase
func NewCreateProductUseCase(productRepo product.ProductRepository) *CreateProductUseCase {
	return &CreateProductUseCase{
		productRepo: productRepo,
	}
}

// Execute performs the create product operation
func (uc *CreateProductUseCase) Execute(ctx context.Context, input CreateProductInput) (*CreateProductOutput, error) {
	// Validate input
	if input.Name == "" {
		return nil, errors.New("product name is required")
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
	if err := uc.productRepo.Create(ctx, prod); err != nil {
		return nil, err
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
