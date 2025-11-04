package persistence

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/JoshuaPangaribuan/clean-arch-ddd/internal/domain/inventory"
	"github.com/JoshuaPangaribuan/clean-arch-ddd/internal/infrastructure/persistence/sqlcgen"
	apperrors "github.com/JoshuaPangaribuan/clean-arch-ddd/pkg/errors"
)

// InventoryRepositoryImpl implements the inventory.InventoryRepository interface
type InventoryRepositoryImpl struct {
	queries *sqlcgen.Queries
}

// NewInventoryRepository creates a new instance of InventoryRepositoryImpl
func NewInventoryRepository(db *sql.DB) inventory.InventoryRepository {
	return &InventoryRepositoryImpl{
		queries: sqlcgen.New(db),
	}
}

// Create stores a new inventory record in the database
func (r *InventoryRepositoryImpl) Create(ctx context.Context, inv *inventory.Inventory) error {
	params := sqlcgen.CreateInventoryParams{
		ID:               inv.ID(),
		ProductID:        inv.ProductID(),
		Quantity:         int32(inv.Quantity()),
		ReservedQuantity: int32(inv.ReservedQuantity()),
		Location:         toNullString(inv.Location()),
		CreatedAt:        inv.CreatedAt(),
		UpdatedAt:        inv.UpdatedAt(),
	}

	err := r.queries.CreateInventory(ctx, params)
	if err != nil {
		return apperrors.WrapDatabaseError(err)
	}
	return nil
}

// GetByProductID retrieves inventory by product ID from the database
func (r *InventoryRepositoryImpl) GetByProductID(ctx context.Context, productID string) (*inventory.Inventory, error) {
	dbInventory, err := r.queries.GetInventoryByProductID(ctx, productID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // Inventory not found
		}
		return nil, apperrors.WrapDatabaseError(err)
	}

	return r.toDomainInventory(dbInventory), nil
}

// Update updates an existing inventory record in the database
func (r *InventoryRepositoryImpl) Update(ctx context.Context, inv *inventory.Inventory) error {
	params := sqlcgen.UpdateInventoryParams{
		ProductID:        inv.ProductID(),
		Quantity:         int32(inv.Quantity()),
		ReservedQuantity: int32(inv.ReservedQuantity()),
		Location:         toNullString(inv.Location()),
		UpdatedAt:        inv.UpdatedAt(),
	}

	err := r.queries.UpdateInventory(ctx, params)
	if err != nil {
		return apperrors.WrapDatabaseError(err)
	}
	return nil
}

// Delete removes an inventory record from the database
func (r *InventoryRepositoryImpl) Delete(ctx context.Context, productID string) error {
	err := r.queries.DeleteInventory(ctx, productID)
	if err != nil {
		return apperrors.WrapDatabaseError(err)
	}
	return nil
}

// AdjustStock adjusts the stock quantity for a product
func (r *InventoryRepositoryImpl) AdjustStock(ctx context.Context, productID string, adjustment int) error {
	params := sqlcgen.AdjustInventoryQuantityParams{
		ProductID: productID,
		Quantity:  int32(adjustment),
		UpdatedAt: time.Now(),
	}

	err := r.queries.AdjustInventoryQuantity(ctx, params)
	if err != nil {
		return apperrors.WrapDatabaseError(err)
	}
	return nil
}

// toDomainInventory converts a database inventory model to a domain inventory entity
func (r *InventoryRepositoryImpl) toDomainInventory(dbInventory sqlcgen.Inventory) *inventory.Inventory {
	return inventory.ReconstructInventory(
		dbInventory.ID,
		dbInventory.ProductID,
		int(dbInventory.Quantity),
		int(dbInventory.ReservedQuantity),
		fromNullString(dbInventory.Location),
		dbInventory.CreatedAt,
		dbInventory.UpdatedAt,
	)
}

// toNullString converts a string to sql.NullString
func toNullString(s string) sql.NullString {
	return sql.NullString{
		String: s,
		Valid:  s != "",
	}
}

// fromNullString converts sql.NullString to string
func fromNullString(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}
