package inventory

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/JoshuaPangaribuan/clean-arch-ddd/internal/application/product"
	"github.com/JoshuaPangaribuan/clean-arch-ddd/internal/domain/inventory"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetInventoryUseCase_Execute_Success(t *testing.T) {
	// Arrange
	mockInventoryRepo := new(MockInventoryRepository)
	mockProductUseCase := new(MockProductUseCase)
	useCase := NewGetInventoryUseCase(mockInventoryRepo, mockProductUseCase)

	inv := inventory.ReconstructInventory(
		"inv-1",
		"product-123",
		100,
		20,
		"Warehouse A",
		time.Now(),
		time.Now(),
	)

	// Mock inventory found
	mockInventoryRepo.On("GetByProductID", mock.Anything, "product-123").
		Return(inv, nil)

	// Mock product details - demonstrates module communication
	mockProductUseCase.On("Execute", mock.Anything, "product-123").
		Return(&product.GetProductOutput{
			ID:            "product-123",
			Name:          "Test Product",
			PriceAmount:   99.99,
			PriceCurrency: "USD",
		}, nil)

	// Act
	output, err := useCase.Execute(context.Background(), "product-123")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, output)
	assert.Equal(t, "inv-1", output.ID)
	assert.Equal(t, "product-123", output.ProductID)
	assert.Equal(t, "Test Product", output.ProductName)
	assert.Equal(t, 99.99, output.ProductPrice)
	assert.Equal(t, "USD", output.ProductCurrency)
	assert.Equal(t, 100, output.Quantity)
	assert.Equal(t, 20, output.ReservedQuantity)
	assert.Equal(t, 80, output.AvailableQuantity)
	assert.Equal(t, "Warehouse A", output.Location)

	mockInventoryRepo.AssertExpectations(t)
	mockProductUseCase.AssertExpectations(t)
}

func TestGetInventoryUseCase_Execute_InventoryNotFound(t *testing.T) {
	// Arrange
	mockInventoryRepo := new(MockInventoryRepository)
	mockProductUseCase := new(MockProductUseCase)
	useCase := NewGetInventoryUseCase(mockInventoryRepo, mockProductUseCase)

	// Mock inventory not found
	mockInventoryRepo.On("GetByProductID", mock.Anything, "product-123").
		Return(nil, nil)

	// Act
	output, err := useCase.Execute(context.Background(), "product-123")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, output)
	assert.Equal(t, inventory.ErrInventoryNotFound, err)

	mockInventoryRepo.AssertExpectations(t)
}

func TestGetInventoryUseCase_Execute_ProductDeleted(t *testing.T) {
	// Arrange
	mockInventoryRepo := new(MockInventoryRepository)
	mockProductUseCase := new(MockProductUseCase)
	useCase := NewGetInventoryUseCase(mockInventoryRepo, mockProductUseCase)

	inv := inventory.ReconstructInventory(
		"inv-1",
		"product-123",
		100,
		20,
		"Warehouse A",
		time.Now(),
		time.Now(),
	)

	// Mock inventory found
	mockInventoryRepo.On("GetByProductID", mock.Anything, "product-123").
		Return(inv, nil)

	// Mock product deleted - demonstrates graceful degradation
	mockProductUseCase.On("Execute", mock.Anything, "product-123").
		Return(nil, errors.New("product not found"))

	// Act
	output, err := useCase.Execute(context.Background(), "product-123")

	// Assert - Should still return inventory with partial data
	assert.NoError(t, err)
	assert.NotNil(t, output)
	assert.Equal(t, "inv-1", output.ID)
	assert.Equal(t, "Unknown (Product Deleted)", output.ProductName)
	assert.Equal(t, float64(0), output.ProductPrice)
	assert.Equal(t, 100, output.Quantity)

	mockInventoryRepo.AssertExpectations(t)
	mockProductUseCase.AssertExpectations(t)
}

