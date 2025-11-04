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
			expectedError: "database connection error",
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
