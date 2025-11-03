package product

import "time"

// CreateProductInput represents the input data for creating a product
type CreateProductInput struct {
	Name          string  `json:"name" validate:"required,min=1,max=255"`
	PriceAmount   float64 `json:"price_amount" validate:"required,gte=0"`
	PriceCurrency string  `json:"price_currency" validate:"required,len=3"`
}

// CreateProductOutput represents the output data after creating a product
type CreateProductOutput struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	PriceAmount   float64   `json:"price_amount"`
	PriceCurrency string    `json:"price_currency"`
	CreatedAt     time.Time `json:"created_at"`
}

// GetProductOutput represents the output data when retrieving a product
type GetProductOutput struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	PriceAmount   float64   `json:"price_amount"`
	PriceCurrency string    `json:"price_currency"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// UpdateProductInput represents the input data for updating a product
type UpdateProductInput struct {
	ID            string  `json:"id" validate:"required"`
	Name          string  `json:"name" validate:"required,min=1,max=255"`
	PriceAmount   float64 `json:"price_amount" validate:"required,gte=0"`
	PriceCurrency string  `json:"price_currency" validate:"required,len=3"`
}

// ListProductsInput represents the input data for listing products
type ListProductsInput struct {
	Limit  int `json:"limit" validate:"required,gte=1,lte=100"`
	Offset int `json:"offset" validate:"gte=0"`
}

// ListProductsOutput represents the output data when listing products
type ListProductsOutput struct {
	Products []GetProductOutput `json:"products"`
	Total    int                `json:"total"`
}
