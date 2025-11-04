package product_test

import (
	"testing"
	"time"

	"github.com/JoshuaPangaribuan/clean-arch-ddd/internal/domain/product"
)

func TestNewProduct(t *testing.T) {
	validPrice, _ := product.NewPrice(99.99, "USD")

	tests := []struct {
		name        string
		id          string
		productName string
		price       product.Price
		wantErr     bool
		errContains string
	}{
		{
			name:        "valid product",
			id:          "product-123",
			productName: "Test Product",
			price:       validPrice,
			wantErr:     false,
		},
		{
			name:        "empty ID",
			id:          "",
			productName: "Test Product",
			price:       validPrice,
			wantErr:     true,
			errContains: "id cannot be empty",
		},
		{
			name:        "empty name",
			id:          "product-123",
			productName: "",
			price:       validPrice,
			wantErr:     true,
			errContains: "name cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := product.NewProduct(tt.id, tt.productName, tt.price)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewProduct() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				if err == nil {
					t.Errorf("NewProduct() expected error containing %q, got nil", tt.errContains)
					return
				}
				if tt.errContains != "" && !contains(err.Error(), tt.errContains) {
					t.Errorf("NewProduct() error = %v, want error containing %q", err, tt.errContains)
				}
			} else {
				if got == nil {
					t.Error("NewProduct() returned nil for valid input")
					return
				}
				if got.ID() != tt.id {
					t.Errorf("NewProduct() ID() = %v, want %v", got.ID(), tt.id)
				}
				if got.Name() != tt.productName {
					t.Errorf("NewProduct() Name() = %v, want %v", got.Name(), tt.productName)
				}
				if !got.Price().Equals(tt.price) {
					t.Errorf("NewProduct() Price() = %v, want %v", got.Price(), tt.price)
				}
				if got.CreatedAt().IsZero() {
					t.Error("NewProduct() CreatedAt() should not be zero")
				}
				if got.UpdatedAt().IsZero() {
					t.Error("NewProduct() UpdatedAt() should not be zero")
				}
			}
		})
	}
}

func TestReconstructProduct(t *testing.T) {
	price, _ := product.NewPrice(99.99, "USD")
	createdAt := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	updatedAt := time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)

	prod := product.ReconstructProduct("product-123", "Test Product", price, createdAt, updatedAt)

	if prod.ID() != "product-123" {
		t.Errorf("ReconstructProduct() ID() = %v, want %v", prod.ID(), "product-123")
	}
	if prod.Name() != "Test Product" {
		t.Errorf("ReconstructProduct() Name() = %v, want %v", prod.Name(), "Test Product")
	}
	if !prod.Price().Equals(price) {
		t.Errorf("ReconstructProduct() Price() = %v, want %v", prod.Price(), price)
	}
	if !prod.CreatedAt().Equal(createdAt) {
		t.Errorf("ReconstructProduct() CreatedAt() = %v, want %v", prod.CreatedAt(), createdAt)
	}
	if !prod.UpdatedAt().Equal(updatedAt) {
		t.Errorf("ReconstructProduct() UpdatedAt() = %v, want %v", prod.UpdatedAt(), updatedAt)
	}
}

func TestProduct_UpdateName(t *testing.T) {
	price, _ := product.NewPrice(99.99, "USD")
	prod, _ := product.NewProduct("product-123", "Original Name", price)
	originalUpdatedAt := prod.UpdatedAt()

	// Wait a bit to ensure timestamp difference
	time.Sleep(10 * time.Millisecond)

	tests := []struct {
		name        string
		newName     string
		wantErr     bool
		errContains string
	}{
		{
			name:    "valid name",
			newName: "Updated Name",
			wantErr: false,
		},
		{
			name:        "empty name",
			newName:     "",
			wantErr:     true,
			errContains: "name cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := prod.UpdateName(tt.newName)
			if (err != nil) != tt.wantErr {
				t.Errorf("Product.UpdateName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if prod.Name() != tt.newName {
					t.Errorf("Product.UpdateName() Name() = %v, want %v", prod.Name(), tt.newName)
				}
				if !prod.UpdatedAt().After(originalUpdatedAt) {
					t.Error("Product.UpdateName() UpdatedAt() should be updated")
				}
			} else if tt.errContains != "" && !contains(err.Error(), tt.errContains) {
				t.Errorf("Product.UpdateName() error = %v, want error containing %q", err, tt.errContains)
			}
		})
	}
}

func TestProduct_UpdatePrice(t *testing.T) {
	originalPrice, _ := product.NewPrice(99.99, "USD")
	prod, _ := product.NewProduct("product-123", "Test Product", originalPrice)
	originalUpdatedAt := prod.UpdatedAt()

	// Wait a bit to ensure timestamp difference
	time.Sleep(10 * time.Millisecond)

	newPrice, _ := product.NewPrice(149.99, "USD")

	err := prod.UpdatePrice(newPrice)
	if err != nil {
		t.Errorf("Product.UpdatePrice() unexpected error = %v", err)
		return
	}

	if !prod.Price().Equals(newPrice) {
		t.Errorf("Product.UpdatePrice() Price() = %v, want %v", prod.Price(), newPrice)
	}
	if !prod.UpdatedAt().After(originalUpdatedAt) {
		t.Error("Product.UpdatePrice() UpdatedAt() should be updated")
	}
}

func TestProduct_Getters(t *testing.T) {
	price, _ := product.NewPrice(99.99, "USD")
	prod, _ := product.NewProduct("product-123", "Test Product", price)

	if prod.ID() != "product-123" {
		t.Errorf("Product.ID() = %v, want %v", prod.ID(), "product-123")
	}
	if prod.Name() != "Test Product" {
		t.Errorf("Product.Name() = %v, want %v", prod.Name(), "Test Product")
	}
	if !prod.Price().Equals(price) {
		t.Errorf("Product.Price() = %v, want %v", prod.Price(), price)
	}
	if prod.CreatedAt().IsZero() {
		t.Error("Product.CreatedAt() should not be zero")
	}
	if prod.UpdatedAt().IsZero() {
		t.Error("Product.UpdatedAt() should not be zero")
	}
}
