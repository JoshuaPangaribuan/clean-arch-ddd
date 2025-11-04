package product

import (
	"context"

	"github.com/JoshuaPangaribuan/clean-arch-ddd/internal/domain/product"
	apperrors "github.com/JoshuaPangaribuan/clean-arch-ddd/pkg/errors"
)

// GetProductUseCase handles the business logic for retrieving a product
type GetProductUseCase struct {
	productRepo    product.ProductRepository
	inventoryUseCase InventoryUseCaseInterface
}

// InventoryUseCaseInterface defines the interface for inventory operations
// This allows Product module to communicate with Inventory module
type InventoryUseCaseInterface interface {
	Execute(ctx context.Context, productID string) (InventoryData, error)
}

// InventoryData represents minimal inventory data needed by Product module
type InventoryData interface {
	GetQuantity() int
	GetAvailableQuantity() int
}

// NewGetProductUseCase creates a new instance of GetProductUseCase without inventory
func NewGetProductUseCase(productRepo product.ProductRepository) *GetProductUseCase {
	return &GetProductUseCase{
		productRepo:    productRepo,
		inventoryUseCase: nil,
	}
}

// NewGetProductUseCaseWithInventory creates a new instance with inventory integration
// This demonstrates bidirectional module communication: Product → Inventory
func NewGetProductUseCaseWithInventory(
	productRepo product.ProductRepository,
	inventoryUseCase InventoryUseCaseInterface,
) *GetProductUseCase {
	return &GetProductUseCase{
		productRepo:    productRepo,
		inventoryUseCase: inventoryUseCase,
	}
}

// Execute performs the get product operation
func (uc *GetProductUseCase) Execute(ctx context.Context, productID string) (*GetProductOutput, error) {
	// Validate input
	if productID == "" {
		return nil, apperrors.New(apperrors.CodeInvalidProductID, "product ID is required")
	}

	// Retrieve product from repository
	prod, err := uc.productRepo.GetByID(ctx, productID)
	if err != nil {
		return nil, apperrors.WrapDatabaseError(err)
	}

	// Check if product exists
	if prod == nil {
		return nil, product.ErrProductNotFound
	}

	// Build base output DTO
	output := &GetProductOutput{
		ID:            prod.ID(),
		Name:          prod.Name(),
		PriceAmount:   prod.Price().Amount(),
		PriceCurrency: prod.Price().Currency(),
		CreatedAt:     prod.CreatedAt(),
		UpdatedAt:     prod.UpdatedAt(),
		HasInventory:  false,
	}

	// MODULE COMMUNICATION: Enrich with inventory data if available
	if uc.inventoryUseCase != nil {
		inventoryData, err := uc.inventoryUseCase.Execute(ctx, productID)
		if err == nil {
			// Successfully retrieved inventory
			output.HasInventory = true
			output.StockQuantity = inventoryData.GetQuantity()
			output.AvailableQuantity = inventoryData.GetAvailableQuantity()
		}
		// Gracefully handle inventory not found - product data is still valid
	}

	return output, nil
}

