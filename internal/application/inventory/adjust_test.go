package inventory

import (
	"context"
	"testing"
	"time"

	"github.com/JoshuaPangaribuan/clean-arch-ddd/internal/application/product"
	"github.com/JoshuaPangaribuan/clean-arch-ddd/internal/domain/inventory"
	apperrors "github.com/JoshuaPangaribuan/clean-arch-ddd/pkg/errors"
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

func TestAdjustInventoryUseCase_Execute(t *testing.T) {
	tests := []struct {
		name           string
		input          AdjustInventoryInput
		setupMocks     func(*MockInventoryRepository, *MockProductUseCase)
		wantErr        bool
		wantErrCode    apperrors.ErrorCode
		wantErrContain string
		validateOutput func(*testing.T, *AdjustInventoryOutput)
	}{
		{
			name: "success - increase quantity",
			input: AdjustInventoryInput{
				ProductID:  "product-123",
				Adjustment: 50,
				Reason:     "Restock",
			},
			setupMocks: func(mockRepo *MockInventoryRepository, mockProduct *MockProductUseCase) {
				inv := inventory.ReconstructInventory(
					"inv-1",
					"product-123",
					100,
					10,
					"Warehouse A",
					time.Now(),
					time.Now(),
				)
				mockProduct.On("Execute", mock.Anything, "product-123").
					Return(&product.GetProductOutput{
						ID:   "product-123",
						Name: "Test Product",
					}, nil)
				mockRepo.On("GetByProductID", mock.Anything, "product-123").Return(inv, nil)
				mockRepo.On("Update", mock.Anything, mock.AnythingOfType("*inventory.Inventory")).Return(nil)
			},
			wantErr: false,
			validateOutput: func(t *testing.T, output *AdjustInventoryOutput) {
				assert.Equal(t, "product-123", output.ProductID)
				assert.Equal(t, 150, output.Quantity) // 100 + 50
				assert.Equal(t, 10, output.ReservedQuantity)
				assert.Equal(t, 140, output.AvailableQuantity)
			},
		},
		{
			name: "success - decrease quantity",
			input: AdjustInventoryInput{
				ProductID:  "product-123",
				Adjustment: -30,
				Reason:     "Damaged goods",
			},
			setupMocks: func(mockRepo *MockInventoryRepository, mockProduct *MockProductUseCase) {
				inv := inventory.ReconstructInventory(
					"inv-1",
					"product-123",
					100,
					10,
					"Warehouse A",
					time.Now(),
					time.Now(),
				)
				mockProduct.On("Execute", mock.Anything, "product-123").
					Return(&product.GetProductOutput{
						ID:   "product-123",
						Name: "Test Product",
					}, nil)
				mockRepo.On("GetByProductID", mock.Anything, "product-123").Return(inv, nil)
				mockRepo.On("Update", mock.Anything, mock.AnythingOfType("*inventory.Inventory")).Return(nil)
			},
			wantErr: false,
			validateOutput: func(t *testing.T, output *AdjustInventoryOutput) {
				assert.Equal(t, 70, output.Quantity) // 100 - 30
				assert.Equal(t, 60, output.AvailableQuantity)
			},
		},
		{
			name: "product not found",
			input: AdjustInventoryInput{
				ProductID:  "nonexistent-product",
				Adjustment: 50,
				Reason:     "Restock",
			},
			setupMocks: func(mockRepo *MockInventoryRepository, mockProduct *MockProductUseCase) {
				// Mock product not found with proper error code
				mockProduct.On("Execute", mock.Anything, "nonexistent-product").
					Return(nil, apperrors.New(apperrors.CodeProductNotFound, "product not found"))
			},
			wantErr:        true,
			wantErrCode:    apperrors.CodeProductNotFound,
			wantErrContain: "cannot adjust inventory: product not found",
		},
		{
			name: "inventory not found",
			input: AdjustInventoryInput{
				ProductID:  "product-123",
				Adjustment: 50,
				Reason:     "Restock",
			},
			setupMocks: func(mockRepo *MockInventoryRepository, mockProduct *MockProductUseCase) {
				mockProduct.On("Execute", mock.Anything, "product-123").
					Return(&product.GetProductOutput{ID: "product-123", Name: "Test Product"}, nil)
				mockRepo.On("GetByProductID", mock.Anything, "product-123").Return(nil, nil)
			},
			wantErr:     true,
			wantErrCode: apperrors.CodeInventoryNotFound,
		},
		{
			name: "cannot go negative",
			input: AdjustInventoryInput{
				ProductID:  "product-123",
				Adjustment: -60, // Would make quantity negative
				Reason:     "Large damage",
			},
			setupMocks: func(mockRepo *MockInventoryRepository, mockProduct *MockProductUseCase) {
				inv := inventory.ReconstructInventory(
					"inv-1",
					"product-123",
					50,
					10,
					"Warehouse A",
					time.Now(),
					time.Now(),
				)
				mockProduct.On("Execute", mock.Anything, "product-123").
					Return(&product.GetProductOutput{ID: "product-123", Name: "Test Product"}, nil)
				mockRepo.On("GetByProductID", mock.Anything, "product-123").Return(inv, nil)
			},
			wantErr:     true,
			wantErrCode: apperrors.CodeInvalidQuantity,
		},
		{
			name: "invalid input - empty product ID",
			input: AdjustInventoryInput{
				ProductID:  "",
				Adjustment: 50,
				Reason:     "Restock",
			},
			setupMocks: func(mockRepo *MockInventoryRepository, mockProduct *MockProductUseCase) {
				// No mocks needed
			},
			wantErr:     true,
			wantErrCode: apperrors.CodeInvalidInput,
		},
		{
			name: "invalid adjustment - zero adjustment",
			input: AdjustInventoryInput{
				ProductID:  "product-123",
				Adjustment: 0,
				Reason:     "Restock",
			},
			setupMocks: func(mockRepo *MockInventoryRepository, mockProduct *MockProductUseCase) {
				// No mocks needed
			},
			wantErr:     true,
			wantErrCode: apperrors.CodeInvalidAdjustment,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockInventoryRepo := new(MockInventoryRepository)
			mockProductUseCase := new(MockProductUseCase)
			useCase := NewAdjustInventoryUseCase(mockInventoryRepo, mockProductUseCase)

			tt.setupMocks(mockInventoryRepo, mockProductUseCase)

			// Act
			output, err := useCase.Execute(context.Background(), tt.input)

			// Assert
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, output)
				if tt.wantErrCode != "" {
					assert.True(t, apperrors.Is(err, tt.wantErrCode), "expected error code %s, got %s", tt.wantErrCode, apperrors.GetCode(err))
				}
				if tt.wantErrContain != "" {
					assert.Contains(t, err.Error(), tt.wantErrContain)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, output)
				if tt.validateOutput != nil {
					tt.validateOutput(t, output)
				}
			}

			mockProductUseCase.AssertExpectations(t)
			mockInventoryRepo.AssertExpectations(t)
		})
	}
}
