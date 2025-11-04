package inventory

import (
	"context"
	"testing"

	"github.com/JoshuaPangaribuan/clean-arch-ddd/internal/application/product"
	"github.com/JoshuaPangaribuan/clean-arch-ddd/internal/domain/inventory"
	apperrors "github.com/JoshuaPangaribuan/clean-arch-ddd/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateInventoryUseCase_Execute(t *testing.T) {
	tests := []struct {
		name           string
		input          CreateInventoryInput
		setupMocks     func(*MockInventoryRepository, *MockProductUseCase)
		wantErr        bool
		wantErrCode    apperrors.ErrorCode
		wantErrContain string
		validateOutput func(*testing.T, *CreateInventoryOutput)
	}{
		{
			name: "success",
			input: CreateInventoryInput{
				ProductID: "product-123",
				Quantity:  100,
				Location:  "Warehouse A",
			},
			setupMocks: func(mockRepo *MockInventoryRepository, mockProduct *MockProductUseCase) {
				mockProduct.On("Execute", mock.Anything, "product-123").
					Return(&product.GetProductOutput{
						ID:   "product-123",
						Name: "Test Product",
					}, nil)
				mockRepo.On("GetByProductID", mock.Anything, "product-123").Return(nil, nil)
				mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*inventory.Inventory")).Return(nil)
			},
			wantErr: false,
			validateOutput: func(t *testing.T, output *CreateInventoryOutput) {
				assert.Equal(t, "product-123", output.ProductID)
				assert.Equal(t, "Test Product", output.ProductName)
				assert.Equal(t, 100, output.Quantity)
				assert.Equal(t, 0, output.ReservedQuantity)
				assert.Equal(t, 100, output.AvailableQuantity)
				assert.Equal(t, "Warehouse A", output.Location)
			},
		},
		{
			name: "product not found",
			input: CreateInventoryInput{
				ProductID: "nonexistent-product",
				Quantity:  100,
				Location:  "Warehouse A",
			},
			setupMocks: func(mockRepo *MockInventoryRepository, mockProduct *MockProductUseCase) {
				// Mock product not found with proper error code
				mockProduct.On("Execute", mock.Anything, "nonexistent-product").
					Return(nil, apperrors.New(apperrors.CodeProductNotFound, "product not found"))
			},
			wantErr:        true,
			wantErrCode:    apperrors.CodeProductNotFound,
			wantErrContain: "cannot create inventory: product not found",
		},
		{
			name: "inventory already exists",
			input: CreateInventoryInput{
				ProductID: "product-123",
				Quantity:  100,
				Location:  "Warehouse A",
			},
			setupMocks: func(mockRepo *MockInventoryRepository, mockProduct *MockProductUseCase) {
				existingInventory, _ := inventory.NewInventory("inv-1", "product-123", 50, "Warehouse B")
				mockProduct.On("Execute", mock.Anything, "product-123").
					Return(&product.GetProductOutput{ID: "product-123", Name: "Test Product"}, nil)
				mockRepo.On("GetByProductID", mock.Anything, "product-123").Return(existingInventory, nil)
			},
			wantErr:     true,
			wantErrCode: apperrors.CodeInventoryExists,
		},
		{
			name: "invalid quantity - negative",
			input: CreateInventoryInput{
				ProductID: "product-123",
				Quantity:  -10,
				Location:  "Warehouse A",
			},
			setupMocks: func(mockRepo *MockInventoryRepository, mockProduct *MockProductUseCase) {
				// No mocks needed
			},
			wantErr:     true,
			wantErrCode: apperrors.CodeInvalidQuantity,
		},
		{
			name: "invalid input - empty product ID",
			input: CreateInventoryInput{
				ProductID: "",
				Quantity:  100,
				Location:  "Warehouse A",
			},
			setupMocks: func(mockRepo *MockInventoryRepository, mockProduct *MockProductUseCase) {
				// No mocks needed
			},
			wantErr:     true,
			wantErrCode: apperrors.CodeInvalidInput,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockInventoryRepo := new(MockInventoryRepository)
			mockProductUseCase := new(MockProductUseCase)
			useCase := NewCreateInventoryUseCase(mockInventoryRepo, mockProductUseCase)

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
