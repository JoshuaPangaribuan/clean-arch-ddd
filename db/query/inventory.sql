-- name: CreateInventory :exec
INSERT INTO inventory (
    id,
    product_id,
    quantity,
    reserved_quantity,
    location,
    created_at,
    updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
);

-- name: GetInventoryByProductID :one
SELECT * FROM inventory
WHERE product_id = $1;

-- name: UpdateInventory :exec
UPDATE inventory
SET
    quantity = $2,
    reserved_quantity = $3,
    location = $4,
    updated_at = $5
WHERE product_id = $1;

-- name: DeleteInventory :exec
DELETE FROM inventory
WHERE product_id = $1;

-- name: AdjustInventoryQuantity :exec
UPDATE inventory
SET
    quantity = quantity + $2,
    updated_at = $3
WHERE product_id = $1;

