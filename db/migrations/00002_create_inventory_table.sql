-- +goose Up
-- Create inventory table
CREATE TABLE IF NOT EXISTS inventory (
    id VARCHAR(36) PRIMARY KEY,
    product_id VARCHAR(36) NOT NULL UNIQUE,
    quantity INTEGER NOT NULL DEFAULT 0,
    reserved_quantity INTEGER NOT NULL DEFAULT 0,
    location VARCHAR(255),
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    CONSTRAINT fk_product
        FOREIGN KEY (product_id)
        REFERENCES products(id)
        ON DELETE CASCADE
);

-- Create index on product_id for faster lookups
CREATE INDEX idx_inventory_product_id ON inventory(product_id);

-- Add check constraints
ALTER TABLE inventory ADD CONSTRAINT check_quantity_positive CHECK (quantity >= 0);
ALTER TABLE inventory ADD CONSTRAINT check_reserved_positive CHECK (reserved_quantity >= 0);
ALTER TABLE inventory ADD CONSTRAINT check_reserved_lte_quantity CHECK (reserved_quantity <= quantity);

-- +goose Down
-- Drop inventory table and related constraints
DROP TABLE IF EXISTS inventory CASCADE;

