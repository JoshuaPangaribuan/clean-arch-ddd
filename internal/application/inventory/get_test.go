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

func TestGetInventoryUseCase_Execute(t *testing.T) {
	tests := []struct {
		name           string
		productID      string
		setupMocks     func(*MockInventoryRepository, *MockProductUseCase)
		wantErr        bool
		wantErrCode    apperrors.ErrorCode
		validateOutput func(*testing.T, *GetInventoryOutput)
	}{
		{
			name:      "success",
			productID: "product-123",
			setupMocks: func(mockRepo *MockInventoryRepository, mockProduct *MockProductUseCase) {
				inv := inventory.ReconstructInventory(
					"inv-1",
					"product-123",
					100,
					20,
					"Warehouse A",
					time.Now(),
					time.Now(),
				)
				mockRepo.On("GetByProductID", mock.Anything, "product-123").Return(inv, nil)
				mockProduct.On("Execute", mock.Anything, "product-123").
					Return(&product.GetProductOutput{
						ID:            "product-123",
						Name:          "Test Product",
						PriceAmount:   99.99,
						PriceCurrency: "USD",
					}, nil)
			},
			wantErr: false,
			validateOutput: func(t *testing.T, output *GetInventoryOutput) {
				assert.Equal(t, "inv-1", output.ID)
				assert.Equal(t, "product-123", output.ProductID)
				assert.Equal(t, "Test Product", output.ProductName)
				assert.Equal(t, 99.99, output.ProductPrice)
				assert.Equal(t, "USD", output.ProductCurrency)
				assert.Equal(t, 100, output.Quantity)
				assert.Equal(t, 20, output.ReservedQuantity)
				assert.Equal(t, 80, output.AvailableQuantity)
				assert.Equal(t, "Warehouse A", output.Location)
			},
		},
		{
			name:      "inventory not found",
			productID: "product-123",
			setupMocks: func(mockRepo *MockInventoryRepository, mockProduct *MockProductUseCase) {
				mockRepo.On("GetByProductID", mock.Anything, "product-123").Return(nil, nil)
			},
			wantErr:     true,
			wantErrCode: apperrors.CodeInventoryNotFound,
		},
		{
			name:      "product deleted - graceful degradation",
			productID: "product-123",
			setupMocks: func(mockRepo *MockInventoryRepository, mockProduct *MockProductUseCase) {
				inv := inventory.ReconstructInventory(
					"inv-1",
					"product-123",
					100,
					20,
					"Warehouse A",
					time.Now(),
					time.Now(),
				)
				mockRepo.On("GetByProductID", mock.Anything, "product-123").Return(inv, nil)
				// Mock product deleted - should gracefully degrade
				mockProduct.On("Execute", mock.Anything, "product-123").
					Return(nil, apperrors.New(apperrors.CodeProductNotFound, "product not found"))
			},
			wantErr: false,
			validateOutput: func(t *testing.T, output *GetInventoryOutput) {
				assert.Equal(t, "inv-1", output.ID)
				assert.Equal(t, "Unknown (Product Deleted)", output.ProductName)
				assert.Equal(t, float64(0), output.ProductPrice)
				assert.Equal(t, "", output.ProductCurrency)
				assert.Equal(t, 100, output.Quantity)
			},
		},
		{
			name:      "invalid input - empty product ID",
			productID: "",
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
			useCase := NewGetInventoryUseCase(mockInventoryRepo, mockProductUseCase)

			tt.setupMocks(mockInventoryRepo, mockProductUseCase)

			// Act
			output, err := useCase.Execute(context.Background(), tt.productID)

			// Assert
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, output)
				if tt.wantErrCode != "" {
					assert.True(t, apperrors.Is(err, tt.wantErrCode), "expected error code %s, got %s", tt.wantErrCode, apperrors.GetCode(err))
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, output)
				if tt.validateOutput != nil {
					tt.validateOutput(t, output)
				}
			}

			mockInventoryRepo.AssertExpectations(t)
			mockProductUseCase.AssertExpectations(t)
		})
	}
}
