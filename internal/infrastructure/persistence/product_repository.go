package persistence

import (
	"context"
	"database/sql"
	"errors"
	"strconv"

	"github.com/JoshuaPangaribuan/clean-arch-ddd/internal/domain/product"
	"github.com/JoshuaPangaribuan/clean-arch-ddd/internal/infrastructure/persistence/sqlcgen"
	apperrors "github.com/JoshuaPangaribuan/clean-arch-ddd/pkg/errors"
)

// ProductRepositoryImpl implements the product.ProductRepository interface
// It also satisfies both ProductCommandRepository and ProductQueryRepository
type ProductRepositoryImpl struct {
	queries *sqlcgen.Queries
}

// NewProductRepository creates a new instance of ProductRepositoryImpl
// Deprecated: Use NewProductCommandRepository and NewProductQueryRepository instead
func NewProductRepository(db *sql.DB) product.ProductRepository {
	return &ProductRepositoryImpl{
		queries: sqlcgen.New(db),
	}
}

// NewProductCommandRepository creates a new instance for command operations
func NewProductCommandRepository(db *sql.DB) product.ProductCommandRepository {
	return &ProductRepositoryImpl{
		queries: sqlcgen.New(db),
	}
}

// NewProductQueryRepository creates a new instance for query operations
func NewProductQueryRepository(db *sql.DB) product.ProductQueryRepository {
	return &ProductRepositoryImpl{
		queries: sqlcgen.New(db),
	}
}

// Create stores a new product in the database
func (r *ProductRepositoryImpl) Create(ctx context.Context, prod *product.Product) error {
	params := sqlcgen.CreateProductParams{
		ID:            prod.ID(),
		Name:          prod.Name(),
		PriceAmount:   strconv.FormatFloat(prod.Price().Amount(), 'f', -1, 64),
		PriceCurrency: prod.Price().Currency(),
		CreatedAt:     prod.CreatedAt(),
		UpdatedAt:     prod.UpdatedAt(),
	}

	err := r.queries.CreateProduct(ctx, params)
	if err != nil {
		return apperrors.WrapDatabaseError(err)
	}
	return nil
}

// GetByID retrieves a product by its ID from the database
func (r *ProductRepositoryImpl) GetByID(ctx context.Context, id string) (*product.Product, error) {
	dbProduct, err := r.queries.GetProductByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // Product not found
		}
		return nil, apperrors.WrapDatabaseError(err)
	}

	return r.toDomainProduct(dbProduct)
}

// Update updates an existing product in the database
func (r *ProductRepositoryImpl) Update(ctx context.Context, prod *product.Product) error {
	params := sqlcgen.UpdateProductParams{
		ID:            prod.ID(),
		Name:          prod.Name(),
		PriceAmount:   strconv.FormatFloat(prod.Price().Amount(), 'f', -1, 64),
		PriceCurrency: prod.Price().Currency(),
		UpdatedAt:     prod.UpdatedAt(),
	}

	err := r.queries.UpdateProduct(ctx, params)
	if err != nil {
		return apperrors.WrapDatabaseError(err)
	}
	return nil
}

// Delete removes a product from the database
func (r *ProductRepositoryImpl) Delete(ctx context.Context, id string) error {
	err := r.queries.DeleteProduct(ctx, id)
	if err != nil {
		return apperrors.WrapDatabaseError(err)
	}
	return nil
}

// List retrieves all products with pagination
func (r *ProductRepositoryImpl) List(ctx context.Context, limit, offset int) ([]*product.Product, error) {
	params := sqlcgen.ListProductsParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	}

	dbProducts, err := r.queries.ListProducts(ctx, params)
	if err != nil {
		return nil, apperrors.WrapDatabaseError(err)
	}

	products := make([]*product.Product, 0, len(dbProducts))
	for _, dbProduct := range dbProducts {
		domainProduct, err := r.toDomainProduct(dbProduct)
		if err != nil {
			return nil, err
		}
		products = append(products, domainProduct)
	}

	return products, nil
}

// toDomainProduct converts a database product model to a domain product entity
func (r *ProductRepositoryImpl) toDomainProduct(dbProduct sqlcgen.Product) (*product.Product, error) {
	priceAmount, err := strconv.ParseFloat(dbProduct.PriceAmount, 64)
	if err != nil {
		return nil, err
	}

	price, err := product.NewPrice(priceAmount, dbProduct.PriceCurrency)
	if err != nil {
		return nil, err
	}

	return product.ReconstructProduct(
		dbProduct.ID,
		dbProduct.Name,
		price,
		dbProduct.CreatedAt,
		dbProduct.UpdatedAt,
	), nil
}
