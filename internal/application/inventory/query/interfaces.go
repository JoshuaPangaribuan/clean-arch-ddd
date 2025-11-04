package query

import (
	"context"

	productquery "github.com/JoshuaPangaribuan/clean-arch-ddd/internal/application/product/query"
)

// ProductQueryInterface defines the interface for product query operations
// This allows Inventory module to communicate with Product module
type ProductQueryInterface interface {
	Execute(ctx context.Context, productID string) (*productquery.GetProductOutput, error)
}

// ProductQueryAdapter adapts GetProductQuery to implement ProductQueryInterface
// This enables Inventory module to call Product module
type ProductQueryAdapter struct {
	productQuery *productquery.GetProductQuery
}

// NewProductQueryAdapter creates a new adapter
func NewProductQueryAdapter(productQuery *productquery.GetProductQuery) *ProductQueryAdapter {
	return &ProductQueryAdapter{
		productQuery: productQuery,
	}
}

// Execute calls the product query and returns the result
func (a *ProductQueryAdapter) Execute(ctx context.Context, productID string) (*productquery.GetProductOutput, error) {
	return a.productQuery.Execute(ctx, productID)
}

