package product

import (
	"errors"
	"fmt"
)

// Price is a value object that represents a monetary amount with currency
type Price struct {
	amount   float64
	currency string
}

// NewPrice creates a new Price value object with validation
func NewPrice(amount float64, currency string) (Price, error) {
	// Business rule: Price cannot be negative
	if amount < 0 {
		return Price{}, errors.New("price amount cannot be negative")
	}

	// Business rule: Currency must be specified
	if currency == "" {
		return Price{}, errors.New("currency cannot be empty")
	}

	// Business rule: Currency must be valid ISO 4217 code (simplified validation)
	if len(currency) != 3 {
		return Price{}, errors.New("currency must be a 3-letter ISO code")
	}

	return Price{
		amount:   amount,
		currency: currency,
	}, nil
}

// Amount returns the price amount
func (p Price) Amount() float64 {
	return p.amount
}

// Currency returns the currency code
func (p Price) Currency() string {
	return p.currency
}

// Equals checks if two prices are equal
func (p Price) Equals(other Price) bool {
	return p.amount == other.amount && p.currency == other.currency
}

// String returns a string representation of the price
func (p Price) String() string {
	return fmt.Sprintf("%.2f %s", p.amount, p.currency)
}

// IsZero checks if the price is zero
func (p Price) IsZero() bool {
	return p.amount == 0
}

// Add adds another price to this price (only if same currency)
func (p Price) Add(other Price) (Price, error) {
	if p.currency != other.currency {
		return Price{}, errors.New("cannot add prices with different currencies")
	}
	return NewPrice(p.amount+other.amount, p.currency)
}

// Subtract subtracts another price from this price (only if same currency)
func (p Price) Subtract(other Price) (Price, error) {
	if p.currency != other.currency {
		return Price{}, errors.New("cannot subtract prices with different currencies")
	}
	return NewPrice(p.amount-other.amount, p.currency)
}
