package inventory

import "time"

// CreateInventoryInput represents the input for creating inventory
type CreateInventoryInput struct {
	ProductID string `json:"product_id" validate:"required"`
	Quantity  int    `json:"quantity" validate:"required,min=0"`
	Location  string `json:"location"`
}

// CreateInventoryOutput represents the output after creating inventory
type CreateInventoryOutput struct {
	ID               string    `json:"id"`
	ProductID        string    `json:"product_id"`
	ProductName      string    `json:"product_name"`
	Quantity         int       `json:"quantity"`
	ReservedQuantity int       `json:"reserved_quantity"`
	AvailableQuantity int      `json:"available_quantity"`
	Location         string    `json:"location"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// GetInventoryOutput represents the output for getting inventory
type GetInventoryOutput struct {
	ID                string    `json:"id"`
	ProductID         string    `json:"product_id"`
	ProductName       string    `json:"product_name"`
	ProductPrice      float64   `json:"product_price"`
	ProductCurrency   string    `json:"product_currency"`
	Quantity          int       `json:"quantity"`
	ReservedQuantity  int       `json:"reserved_quantity"`
	AvailableQuantity int       `json:"available_quantity"`
	Location          string    `json:"location"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// AdjustInventoryInput represents the input for adjusting inventory
type AdjustInventoryInput struct {
	ProductID  string `json:"product_id" validate:"required"`
	Adjustment int    `json:"adjustment" validate:"required"`
	Reason     string `json:"reason"`
}

// AdjustInventoryOutput represents the output after adjusting inventory
type AdjustInventoryOutput struct {
	ID                string    `json:"id"`
	ProductID         string    `json:"product_id"`
	ProductName       string    `json:"product_name"`
	Quantity          int       `json:"quantity"`
	ReservedQuantity  int       `json:"reserved_quantity"`
	AvailableQuantity int       `json:"available_quantity"`
	Location          string    `json:"location"`
	UpdatedAt         time.Time `json:"updated_at"`
}

