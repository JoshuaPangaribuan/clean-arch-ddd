package product

import (
	"context"
	"errors"

	"github.com/JoshuaPangaribuan/clean-arch-ddd/internal/domain/product"
)

// GetProductUseCase handles the business logic for retrieving a product
type GetProductUseCase struct {
	productRepo product.ProductRepository
}

// NewGetProductUseCase creates a new instance of GetProductUseCase
func NewGetProductUseCase(productRepo product.ProductRepository) *GetProductUseCase {
	return &GetProductUseCase{
		productRepo: productRepo,
	}
}

// Execute performs the get product operation
func (uc *GetProductUseCase) Execute(ctx context.Context, productID string) (*GetProductOutput, error) {
	// Validate input
	if productID == "" {
		return nil, errors.New("product ID is required")
	}

	// Retrieve product from repository
	prod, err := uc.productRepo.GetByID(ctx, productID)
	if err != nil {
		return nil, err
	}

	// Check if product exists
	if prod == nil {
		return nil, errors.New("product not found")
	}

	// Return output DTO
	return &GetProductOutput{
		ID:            prod.ID(),
		Name:          prod.Name(),
		PriceAmount:   prod.Price().Amount(),
		PriceCurrency: prod.Price().Currency(),
		CreatedAt:     prod.CreatedAt(),
		UpdatedAt:     prod.UpdatedAt(),
	}, nil
}

