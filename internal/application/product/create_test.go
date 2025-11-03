package product_test

import (
	"context"
	"errors"
	"testing"

	"github.com/JoshuaPangaribuan/clean-arch-ddd/internal/application/product"
	"github.com/JoshuaPangaribuan/clean-arch-ddd/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type CreateProductUseCaseTestSuite struct {
	suite.Suite
	mockRepo *mocks.ProductRepository
	useCase  *product.CreateProductUseCase
}

func (s *CreateProductUseCaseTestSuite) SetupTest() {
	s.mockRepo = mocks.NewProductRepository(s.T())
	s.useCase = product.NewCreateProductUseCase(s.mockRepo)
}

func (s *CreateProductUseCaseTestSuite) TestExecute() {
	testCases := []struct {
		name           string
		input          product.CreateProductInput
		setupMock      func()
		expectedError  string
		validateOutput func(*product.CreateProductOutput)
	}{
		{
			name: "Success",
			input: product.CreateProductInput{
				Name:          "Test Product",
				PriceAmount:   99.99,
				PriceCurrency: "USD",
			},
			setupMock: func() {
				s.mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*product.Product")).
					Return(nil).
					Once()
			},
			expectedError: "",
			validateOutput: func(output *product.CreateProductOutput) {
				s.NotNil(output)
				s.NotEmpty(output.ID)
				s.Equal("Test Product", output.Name)
				s.Equal(99.99, output.PriceAmount)
				s.Equal("USD", output.PriceCurrency)
				s.False(output.CreatedAt.IsZero())
			},
		},
		{
			name: "Empty Name",
			input: product.CreateProductInput{
				Name:          "",
				PriceAmount:   99.99,
				PriceCurrency: "USD",
			},
			setupMock:     func() {},
			expectedError: "name",
			validateOutput: func(output *product.CreateProductOutput) {
				s.Nil(output)
			},
		},
		{
			name: "Negative Price",
			input: product.CreateProductInput{
				Name:          "Test Product",
				PriceAmount:   -10.00,
				PriceCurrency: "USD",
			},
			setupMock:     func() {},
			expectedError: "negative",
			validateOutput: func(output *product.CreateProductOutput) {
				s.Nil(output)
			},
		},
		{
			name: "Invalid Currency",
			input: product.CreateProductInput{
				Name:          "Test Product",
				PriceAmount:   99.99,
				PriceCurrency: "US", // Invalid: not 3 characters
			},
			setupMock:     func() {},
			expectedError: "currency",
			validateOutput: func(output *product.CreateProductOutput) {
				s.Nil(output)
			},
		},
		{
			name: "Repository Error",
			input: product.CreateProductInput{
				Name:          "Test Product",
				PriceAmount:   99.99,
				PriceCurrency: "USD",
			},
			setupMock: func() {
				s.mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*product.Product")).
					Return(errors.New("database connection error")).
					Once()
			},
			expectedError: "database connection error",
			validateOutput: func(output *product.CreateProductOutput) {
				s.Nil(output)
			},
		},
		{
			name: "Zero Price",
			input: product.CreateProductInput{
				Name:          "Free Product",
				PriceAmount:   0.00,
				PriceCurrency: "USD",
			},
			setupMock: func() {
				s.mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*product.Product")).
					Return(nil).
					Once()
			},
			expectedError: "",
			validateOutput: func(output *product.CreateProductOutput) {
				s.NotNil(output)
				s.Equal(0.00, output.PriceAmount)
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			// Setup
			tc.setupMock()

			// Execute
			output, err := s.useCase.Execute(context.Background(), tc.input)

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

func TestCreateProductUseCaseTestSuite(t *testing.T) {
	suite.Run(t, new(CreateProductUseCaseTestSuite))
}
