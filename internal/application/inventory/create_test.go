package inventory

import (
	"context"
	"errors"
	"testing"

	"github.com/JoshuaPangaribuan/clean-arch-ddd/internal/application/product"
	"github.com/JoshuaPangaribuan/clean-arch-ddd/internal/domain/inventory"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockInventoryRepository is a mock implementation of inventory.InventoryRepository
type MockInventoryRepository struct {
	mock.Mock
}

func (m *MockInventoryRepository) Create(ctx context.Context, inv *inventory.Inventory) error {
	args := m.Called(ctx, inv)
	return args.Error(0)
}

func (m *MockInventoryRepository) GetByProductID(ctx context.Context, productID string) (*inventory.Inventory, error) {
	args := m.Called(ctx, productID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*inventory.Inventory), args.Error(1)
}

func (m *MockInventoryRepository) Update(ctx context.Context, inv *inventory.Inventory) error {
	args := m.Called(ctx, inv)
	return args.Error(0)
}

func (m *MockInventoryRepository) Delete(ctx context.Context, productID string) error {
	args := m.Called(ctx, productID)
	return args.Error(0)
}

func (m *MockInventoryRepository) AdjustStock(ctx context.Context, productID string, adjustment int) error {
	args := m.Called(ctx, productID, adjustment)
	return args.Error(0)
}

// MockProductUseCase is a mock implementation of product.ProductUseCaseInterface
type MockProductUseCase struct {
	mock.Mock
}

func (m *MockProductUseCase) Execute(ctx context.Context, productID string) (*product.GetProductOutput, error) {
	args := m.Called(ctx, productID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*product.GetProductOutput), args.Error(1)
}

func TestCreateInventoryUseCase_Execute_Success(t *testing.T) {
	// Arrange
	mockInventoryRepo := new(MockInventoryRepository)
	mockProductUseCase := new(MockProductUseCase)
	useCase := NewCreateInventoryUseCase(mockInventoryRepo, mockProductUseCase)

	input := CreateInventoryInput{
		ProductID: "product-123",
		Quantity:  100,
		Location:  "Warehouse A",
	}

	// Mock product exists
	mockProductUseCase.On("Execute", mock.Anything, "product-123").
		Return(&product.GetProductOutput{
			ID:   "product-123",
			Name: "Test Product",
		}, nil)

	// Mock no existing inventory
	mockInventoryRepo.On("GetByProductID", mock.Anything, "product-123").
		Return(nil, nil)

	// Mock successful create
	mockInventoryRepo.On("Create", mock.Anything, mock.AnythingOfType("*inventory.Inventory")).
		Return(nil)

	// Act
	output, err := useCase.Execute(context.Background(), input)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, output)
	assert.Equal(t, "product-123", output.ProductID)
	assert.Equal(t, "Test Product", output.ProductName)
	assert.Equal(t, 100, output.Quantity)
	assert.Equal(t, 0, output.ReservedQuantity)
	assert.Equal(t, 100, output.AvailableQuantity)
	assert.Equal(t, "Warehouse A", output.Location)

	mockProductUseCase.AssertExpectations(t)
	mockInventoryRepo.AssertExpectations(t)
}

func TestCreateInventoryUseCase_Execute_ProductNotFound(t *testing.T) {
	// Arrange
	mockInventoryRepo := new(MockInventoryRepository)
	mockProductUseCase := new(MockProductUseCase)
	useCase := NewCreateInventoryUseCase(mockInventoryRepo, mockProductUseCase)

	input := CreateInventoryInput{
		ProductID: "nonexistent-product",
		Quantity:  100,
		Location:  "Warehouse A",
	}

	// Mock product not found - demonstrates module communication
	mockProductUseCase.On("Execute", mock.Anything, "nonexistent-product").
		Return(nil, errors.New("product not found"))

	// Act
	output, err := useCase.Execute(context.Background(), input)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, output)
	assert.Contains(t, err.Error(), "cannot create inventory: product not found")

	mockProductUseCase.AssertExpectations(t)
}

func TestCreateInventoryUseCase_Execute_InventoryAlreadyExists(t *testing.T) {
	// Arrange
	mockInventoryRepo := new(MockInventoryRepository)
	mockProductUseCase := new(MockProductUseCase)
	useCase := NewCreateInventoryUseCase(mockInventoryRepo, mockProductUseCase)

	input := CreateInventoryInput{
		ProductID: "product-123",
		Quantity:  100,
		Location:  "Warehouse A",
	}

	existingInventory, _ := inventory.NewInventory("inv-1", "product-123", 50, "Warehouse B")

	// Mock product exists
	mockProductUseCase.On("Execute", mock.Anything, "product-123").
		Return(&product.GetProductOutput{ID: "product-123", Name: "Test Product"}, nil)

	// Mock inventory already exists
	mockInventoryRepo.On("GetByProductID", mock.Anything, "product-123").
		Return(existingInventory, nil)

	// Act
	output, err := useCase.Execute(context.Background(), input)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, output)
	assert.Equal(t, inventory.ErrInventoryExists, err)

	mockProductUseCase.AssertExpectations(t)
	mockInventoryRepo.AssertExpectations(t)
}

func TestCreateInventoryUseCase_Execute_InvalidQuantity(t *testing.T) {
	// Arrange
	mockInventoryRepo := new(MockInventoryRepository)
	mockProductUseCase := new(MockProductUseCase)
	useCase := NewCreateInventoryUseCase(mockInventoryRepo, mockProductUseCase)

	input := CreateInventoryInput{
		ProductID: "product-123",
		Quantity:  -10,
		Location:  "Warehouse A",
	}

	// Act
	output, err := useCase.Execute(context.Background(), input)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, output)
	assert.Equal(t, inventory.ErrInvalidQuantity, err)
}

