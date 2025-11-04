-- name: CreateProduct :exec
INSERT INTO products (
    id, 
    name, 
    price_amount, 
    price_currency, 
    created_at, 
    updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6
);

-- name: GetProductByID :one
SELECT 
    id, 
    name, 
    price_amount, 
    price_currency, 
    created_at, 
    updated_at
FROM products
WHERE id = $1;

-- name: UpdateProduct :exec
UPDATE products
SET 
    name = $2,
    price_amount = $3,
    price_currency = $4,
    updated_at = $5
WHERE id = $1;

-- name: DeleteProduct :exec
DELETE FROM products
WHERE id = $1;

-- name: ListProducts :many
SELECT 
    id, 
    name, 
    price_amount, 
    price_currency, 
    created_at, 
    updated_at
FROM products
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

