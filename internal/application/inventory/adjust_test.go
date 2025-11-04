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

func TestAdjustInventoryUseCase_Execute_IncreaseSuccess(t *testing.T) {
	// Arrange
	mockInventoryRepo := new(MockInventoryRepository)
	mockProductUseCase := new(MockProductUseCase)
	useCase := NewAdjustInventoryUseCase(mockInventoryRepo, mockProductUseCase)

	inv := inventory.ReconstructInventory(
		"inv-1",
		"product-123",
		100,
		10,
		"Warehouse A",
		time.Now(),
		time.Now(),
	)

	input := AdjustInventoryInput{
		ProductID:  "product-123",
		Adjustment: 50,
		Reason:     "Restock",
	}

	// Mock product exists
	mockProductUseCase.On("Execute", mock.Anything, "product-123").
		Return(&product.GetProductOutput{
			ID:   "product-123",
			Name: "Test Product",
		}, nil)

	// Mock inventory found
	mockInventoryRepo.On("GetByProductID", mock.Anything, "product-123").
		Return(inv, nil)

	// Mock successful update
	mockInventoryRepo.On("Update", mock.Anything, mock.AnythingOfType("*inventory.Inventory")).
		Return(nil)

	// Act
	output, err := useCase.Execute(context.Background(), input)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, output)
	assert.Equal(t, "product-123", output.ProductID)
	assert.Equal(t, 150, output.Quantity) // 100 + 50
	assert.Equal(t, 10, output.ReservedQuantity)
	assert.Equal(t, 140, output.AvailableQuantity)

	mockProductUseCase.AssertExpectations(t)
	mockInventoryRepo.AssertExpectations(t)
}

func TestAdjustInventoryUseCase_Execute_DecreaseSuccess(t *testing.T) {
	// Arrange
	mockInventoryRepo := new(MockInventoryRepository)
	mockProductUseCase := new(MockProductUseCase)
	useCase := NewAdjustInventoryUseCase(mockInventoryRepo, mockProductUseCase)

	inv := inventory.ReconstructInventory(
		"inv-1",
		"product-123",
		100,
		10,
		"Warehouse A",
		time.Now(),
		time.Now(),
	)

	input := AdjustInventoryInput{
		ProductID:  "product-123",
		Adjustment: -30,
		Reason:     "Damaged goods",
	}

	// Mock product exists
	mockProductUseCase.On("Execute", mock.Anything, "product-123").
		Return(&product.GetProductOutput{
			ID:   "product-123",
			Name: "Test Product",
		}, nil)

	// Mock inventory found
	mockInventoryRepo.On("GetByProductID", mock.Anything, "product-123").
		Return(inv, nil)

	// Mock successful update
	mockInventoryRepo.On("Update", mock.Anything, mock.AnythingOfType("*inventory.Inventory")).
		Return(nil)

	// Act
	output, err := useCase.Execute(context.Background(), input)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, output)
	assert.Equal(t, 70, output.Quantity) // 100 - 30
	assert.Equal(t, 60, output.AvailableQuantity)

	mockProductUseCase.AssertExpectations(t)
	mockInventoryRepo.AssertExpectations(t)
}

func TestAdjustInventoryUseCase_Execute_ProductNotFound(t *testing.T) {
	// Arrange
	mockInventoryRepo := new(MockInventoryRepository)
	mockProductUseCase := new(MockProductUseCase)
	useCase := NewAdjustInventoryUseCase(mockInventoryRepo, mockProductUseCase)

	input := AdjustInventoryInput{
		ProductID:  "nonexistent-product",
		Adjustment: 50,
		Reason:     "Restock",
	}

	// Mock product not found - demonstrates module communication
	mockProductUseCase.On("Execute", mock.Anything, "nonexistent-product").
		Return(nil, errors.New("product not found"))

	// Act
	output, err := useCase.Execute(context.Background(), input)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, output)
	assert.Contains(t, err.Error(), "cannot adjust inventory: product not found")

	mockProductUseCase.AssertExpectations(t)
}

func TestAdjustInventoryUseCase_Execute_InventoryNotFound(t *testing.T) {
	// Arrange
	mockInventoryRepo := new(MockInventoryRepository)
	mockProductUseCase := new(MockProductUseCase)
	useCase := NewAdjustInventoryUseCase(mockInventoryRepo, mockProductUseCase)

	input := AdjustInventoryInput{
		ProductID:  "product-123",
		Adjustment: 50,
		Reason:     "Restock",
	}

	// Mock product exists
	mockProductUseCase.On("Execute", mock.Anything, "product-123").
		Return(&product.GetProductOutput{ID: "product-123", Name: "Test Product"}, nil)

	// Mock inventory not found
	mockInventoryRepo.On("GetByProductID", mock.Anything, "product-123").
		Return(nil, nil)

	// Act
	output, err := useCase.Execute(context.Background(), input)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, output)
	assert.Equal(t, inventory.ErrInventoryNotFound, err)

	mockProductUseCase.AssertExpectations(t)
	mockInventoryRepo.AssertExpectations(t)
}

func TestAdjustInventoryUseCase_Execute_CannotGoNegative(t *testing.T) {
	// Arrange
	mockInventoryRepo := new(MockInventoryRepository)
	mockProductUseCase := new(MockProductUseCase)
	useCase := NewAdjustInventoryUseCase(mockInventoryRepo, mockProductUseCase)

	inv := inventory.ReconstructInventory(
		"inv-1",
		"product-123",
		50,
		10,
		"Warehouse A",
		time.Now(),
		time.Now(),
	)

	input := AdjustInventoryInput{
		ProductID:  "product-123",
		Adjustment: -60, // Would make quantity negative
		Reason:     "Large damage",
	}

	// Mock product exists
	mockProductUseCase.On("Execute", mock.Anything, "product-123").
		Return(&product.GetProductOutput{ID: "product-123", Name: "Test Product"}, nil)

	// Mock inventory found
	mockInventoryRepo.On("GetByProductID", mock.Anything, "product-123").
		Return(inv, nil)

	// Act
	output, err := useCase.Execute(context.Background(), input)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, output)
	assert.Equal(t, inventory.ErrInvalidQuantity, err)

	mockProductUseCase.AssertExpectations(t)
	mockInventoryRepo.AssertExpectations(t)
}

