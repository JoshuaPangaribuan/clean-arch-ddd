package product_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/JoshuaPangaribuan/clean-arch-ddd/internal/application/product"
	domainProduct "github.com/JoshuaPangaribuan/clean-arch-ddd/internal/domain/product"
	"github.com/JoshuaPangaribuan/clean-arch-ddd/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type GetProductUseCaseTestSuite struct {
	suite.Suite
	mockRepo *mocks.MockProductRepository
	useCase  *product.GetProductUseCase
}

func (s *GetProductUseCaseTestSuite) SetupTest() {
	s.mockRepo = mocks.NewMockProductRepository(s.T())
	s.useCase = product.NewGetProductUseCase(s.mockRepo)
}

func (s *GetProductUseCaseTestSuite) TestExecute() {
	testCases := []struct {
		name           string
		productID      string
		setupMock      func()
		expectedError  string
		validateOutput func(*product.GetProductOutput)
	}{
		{
			name:      "Success",
			productID: "test-product-id",
			setupMock: func() {
				price, _ := domainProduct.NewPrice(99.99, "USD")
				expectedProduct := domainProduct.ReconstructProduct(
					"test-product-id",
					"Test Product",
					price,
					time.Now(),
					time.Now(),
				)
				s.mockRepo.On("GetByID", mock.Anything, "test-product-id").
					Return(expectedProduct, nil).
					Once()
			},
			expectedError: "",
			validateOutput: func(output *product.GetProductOutput) {
				s.NotNil(output)
				s.Equal("test-product-id", output.ID)
				s.Equal("Test Product", output.Name)
				s.Equal(99.99, output.PriceAmount)
				s.Equal("USD", output.PriceCurrency)
				s.False(output.CreatedAt.IsZero())
				s.False(output.UpdatedAt.IsZero())
			},
		},
		{
			name:          "Empty ID",
			productID:     "",
			setupMock:     func() {},
			expectedError: "ID",
			validateOutput: func(output *product.GetProductOutput) {
				s.Nil(output)
			},
		},
		{
			name:      "Product Not Found",
			productID: "non-existent-id",
			setupMock: func() {
				s.mockRepo.On("GetByID", mock.Anything, "non-existent-id").
					Return(nil, nil).
					Once()
			},
			expectedError: "not found",
			validateOutput: func(output *product.GetProductOutput) {
				s.Nil(output)
			},
		},
		{
			name:      "Repository Error",
			productID: "test-product-id",
			setupMock: func() {
				s.mockRepo.On("GetByID", mock.Anything, "test-product-id").
					Return(nil, errors.New("database connection error")).
					Once()
			},
			expectedError: "connection", // Check for keyword instead of exact match (case-insensitive)
			validateOutput: func(output *product.GetProductOutput) {
				s.Nil(output)
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			// Setup
			tc.setupMock()

			// Execute
			output, err := s.useCase.Execute(context.Background(), tc.productID)

			// Assert
			if tc.expectedError != "" {
				s.Error(err)
				s.Contains(err.Error(), tc.expectedError)
			} else {
				s.NoError(err)
			}

			tc.validateOutput(output)
		})
	}
}

func TestGetProductUseCaseTestSuite(t *testing.T) {
	suite.Run(t, new(GetProductUseCaseTestSuite))
}

// MockInventoryUseCase is a mock for testing inventory integration
type MockInventoryUseCase struct {
	mock.Mock
}

func (m *MockInventoryUseCase) Execute(ctx context.Context, productID string) (product.InventoryData, error) {
	args := m.Called(ctx, productID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(product.InventoryData), args.Error(1)
}

// MockInventoryData is a mock implementation of InventoryData
type MockInventoryData struct {
	quantity          int
	availableQuantity int
}

func (m *MockInventoryData) GetQuantity() int {
	return m.quantity
}

func (m *MockInventoryData) GetAvailableQuantity() int {
	return m.availableQuantity
}

// Test suite for Product with Inventory integration
type GetProductWithInventoryUseCaseTestSuite struct {
	suite.Suite
	mockRepo             *mocks.MockProductRepository
	mockInventoryUseCase *MockInventoryUseCase
	useCase              *product.GetProductUseCase
}

func (s *GetProductWithInventoryUseCaseTestSuite) SetupTest() {
	s.mockRepo = mocks.NewMockProductRepository(s.T())
	s.mockInventoryUseCase = new(MockInventoryUseCase)
	// Use the constructor with inventory integration
	s.useCase = product.NewGetProductUseCaseWithInventory(s.mockRepo, s.mockInventoryUseCase)
}

func (s *GetProductWithInventoryUseCaseTestSuite) TestExecute_WithInventory() {
	// Setup product
	price, _ := domainProduct.NewPrice(99.99, "USD")
	expectedProduct := domainProduct.ReconstructProduct(
		"test-product-id",
		"Test Product",
		price,
		time.Now(),
		time.Now(),
	)

	// Setup inventory data
	inventoryData := &MockInventoryData{
		quantity:          100,
		availableQuantity: 80,
	}

	s.mockRepo.On("GetByID", mock.Anything, "test-product-id").
		Return(expectedProduct, nil).
		Once()

	// Mock inventory use case - demonstrates Product → Inventory communication
	s.mockInventoryUseCase.On("Execute", mock.Anything, "test-product-id").
		Return(inventoryData, nil).
		Once()

	// Execute
	output, err := s.useCase.Execute(context.Background(), "test-product-id")

	// Assert
	s.NoError(err)
	s.NotNil(output)
	s.Equal("test-product-id", output.ID)
	s.Equal("Test Product", output.Name)
	s.Equal(99.99, output.PriceAmount)
	s.True(output.HasInventory)
	s.Equal(100, output.StockQuantity)
	s.Equal(80, output.AvailableQuantity)

	s.mockRepo.AssertExpectations(s.T())
	s.mockInventoryUseCase.AssertExpectations(s.T())
}

func (s *GetProductWithInventoryUseCaseTestSuite) TestExecute_InventoryNotFound() {
	// Setup product
	price, _ := domainProduct.NewPrice(99.99, "USD")
	expectedProduct := domainProduct.ReconstructProduct(
		"test-product-id",
		"Test Product",
		price,
		time.Now(),
		time.Now(),
	)

	s.mockRepo.On("GetByID", mock.Anything, "test-product-id").
		Return(expectedProduct, nil).
		Once()

	// Mock inventory not found - should gracefully degrade
	s.mockInventoryUseCase.On("Execute", mock.Anything, "test-product-id").
		Return(nil, errors.New("inventory not found")).
		Once()

	// Execute
	output, err := s.useCase.Execute(context.Background(), "test-product-id")

	// Assert - Product should still be returned without inventory
	s.NoError(err)
	s.NotNil(output)
	s.Equal("test-product-id", output.ID)
	s.Equal("Test Product", output.Name)
	s.False(output.HasInventory)
	s.Equal(0, output.StockQuantity)
	s.Equal(0, output.AvailableQuantity)

	s.mockRepo.AssertExpectations(s.T())
	s.mockInventoryUseCase.AssertExpectations(s.T())
}

func TestGetProductWithInventoryUseCaseTestSuite(t *testing.T) {
	suite.Run(t, new(GetProductWithInventoryUseCaseTestSuite))
}
