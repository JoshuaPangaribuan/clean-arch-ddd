package product_test

import (
	"testing"

	"github.com/JoshuaPangaribuan/clean-arch-ddd/internal/domain/product"
)

func TestNewPrice(t *testing.T) {
	tests := []struct {
		name        string
		amount      float64
		currency    string
		wantErr     bool
		errContains string
	}{
		{
			name:     "valid price with USD",
			amount:   99.99,
			currency: "USD",
			wantErr:  false,
		},
		{
			name:     "valid price with EUR",
			amount:   50.00,
			currency: "EUR",
			wantErr:  false,
		},
		{
			name:     "valid zero price",
			amount:   0.00,
			currency: "USD",
			wantErr:  false,
		},
		{
			name:        "negative amount",
			amount:      -10.00,
			currency:    "USD",
			wantErr:     true,
			errContains: "negative",
		},
		{
			name:        "empty currency",
			amount:      99.99,
			currency:    "",
			wantErr:     true,
			errContains: "currency",
		},
		{
			name:        "currency too short",
			amount:      99.99,
			currency:    "US",
			wantErr:     true,
			errContains: "3-letter",
		},
		{
			name:        "currency too long",
			amount:      99.99,
			currency:    "USDD",
			wantErr:     true,
			errContains: "3-letter",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := product.NewPrice(tt.amount, tt.currency)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewPrice() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				if err == nil {
					t.Errorf("NewPrice() expected error containing %q, got nil", tt.errContains)
					return
				}
				if tt.errContains != "" && !contains(err.Error(), tt.errContains) {
					t.Errorf("NewPrice() error = %v, want error containing %q", err, tt.errContains)
				}
			} else {
				if got.Amount() != tt.amount {
					t.Errorf("NewPrice() Amount() = %v, want %v", got.Amount(), tt.amount)
				}
				if got.Currency() != tt.currency {
					t.Errorf("NewPrice() Currency() = %v, want %v", got.Currency(), tt.currency)
				}
			}
		})
	}
}

func TestPrice_Equals(t *testing.T) {
	price1, _ := product.NewPrice(100.00, "USD")
	price2, _ := product.NewPrice(100.00, "USD")
	price3, _ := product.NewPrice(100.00, "EUR")
	price4, _ := product.NewPrice(50.00, "USD")

	tests := []struct {
		name  string
		price product.Price
		other product.Price
		want  bool
	}{
		{
			name:  "equal prices",
			price: price1,
			other: price2,
			want:  true,
		},
		{
			name:  "different currencies",
			price: price1,
			other: price3,
			want:  false,
		},
		{
			name:  "different amounts",
			price: price1,
			other: price4,
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.price.Equals(tt.other); got != tt.want {
				t.Errorf("Price.Equals() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPrice_IsZero(t *testing.T) {
	zeroPrice, _ := product.NewPrice(0.00, "USD")
	nonZeroPrice, _ := product.NewPrice(100.00, "USD")

	tests := []struct {
		name  string
		price product.Price
		want  bool
	}{
		{
			name:  "zero price",
			price: zeroPrice,
			want:  true,
		},
		{
			name:  "non-zero price",
			price: nonZeroPrice,
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.price.IsZero(); got != tt.want {
				t.Errorf("Price.IsZero() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPrice_Add(t *testing.T) {
	price1, _ := product.NewPrice(100.00, "USD")
	price2, _ := product.NewPrice(50.00, "USD")
	price3, _ := product.NewPrice(50.00, "EUR")

	tests := []struct {
		name        string
		price       product.Price
		other       product.Price
		wantAmount  float64
		wantErr     bool
		errContains string
	}{
		{
			name:       "same currency",
			price:      price1,
			other:      price2,
			wantAmount: 150.00,
			wantErr:    false,
		},
		{
			name:        "different currencies",
			price:       price1,
			other:       price3,
			wantErr:     true,
			errContains: "different currencies",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.price.Add(tt.other)
			if (err != nil) != tt.wantErr {
				t.Errorf("Price.Add() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got.Amount() != tt.wantAmount {
					t.Errorf("Price.Add() Amount() = %v, want %v", got.Amount(), tt.wantAmount)
				}
				if got.Currency() != tt.price.Currency() {
					t.Errorf("Price.Add() Currency() = %v, want %v", got.Currency(), tt.price.Currency())
				}
			} else if tt.errContains != "" && !contains(err.Error(), tt.errContains) {
				t.Errorf("Price.Add() error = %v, want error containing %q", err, tt.errContains)
			}
		})
	}
}

func TestPrice_Subtract(t *testing.T) {
	price1, _ := product.NewPrice(100.00, "USD")
	price2, _ := product.NewPrice(30.00, "USD")
	price3, _ := product.NewPrice(50.00, "EUR")

	tests := []struct {
		name        string
		price       product.Price
		other       product.Price
		wantAmount  float64
		wantErr     bool
		errContains string
	}{
		{
			name:       "same currency",
			price:      price1,
			other:      price2,
			wantAmount: 70.00,
			wantErr:    false,
		},
		{
			name:        "different currencies",
			price:       price1,
			other:       price3,
			wantErr:     true,
			errContains: "different currencies",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.price.Subtract(tt.other)
			if (err != nil) != tt.wantErr {
				t.Errorf("Price.Subtract() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got.Amount() != tt.wantAmount {
					t.Errorf("Price.Subtract() Amount() = %v, want %v", got.Amount(), tt.wantAmount)
				}
				if got.Currency() != tt.price.Currency() {
					t.Errorf("Price.Subtract() Currency() = %v, want %v", got.Currency(), tt.price.Currency())
				}
			} else if tt.errContains != "" && !contains(err.Error(), tt.errContains) {
				t.Errorf("Price.Subtract() error = %v, want error containing %q", err, tt.errContains)
			}
		})
	}
}

func TestPrice_Subtract_ResultNegative(t *testing.T) {
	price1, _ := product.NewPrice(50.00, "USD")
	price2, _ := product.NewPrice(100.00, "USD")

	// Subtracting a larger amount should fail because it would result in negative price
	_, err := price1.Subtract(price2)
	if err == nil {
		t.Error("Price.Subtract() with result < 0 should return error")
	}
	if err != nil && !contains(err.Error(), "negative") {
		t.Errorf("Price.Subtract() error = %v, want error containing 'negative'", err)
	}
}

func TestPrice_String(t *testing.T) {
	price, _ := product.NewPrice(99.99, "USD")
	want := "99.99 USD"
	if got := price.String(); got != want {
		t.Errorf("Price.String() = %v, want %v", got, want)
	}
}

// Helper function
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
