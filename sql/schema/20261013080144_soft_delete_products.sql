-- +goose Up
ALTER TABLE products
    ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMPTZ NULL;

-- Hot path for every catalog query
CREATE INDEX IF NOT EXISTS idx_products_not_deleted
    ON products (id)
    WHERE deleted_at IS NULL;

-- Shop dashboard / listing queries
CREATE INDEX IF NOT EXISTS idx_products_shop_status_active
    ON products (shop_id, status)
    WHERE deleted_at IS NULL;

ALTER TABLE products DROP CONSTRAINT IF EXISTS products_shop_id_name_key;
CREATE UNIQUE INDEX IF NOT EXISTS uq_products_shop_name_active
    ON products (shop_id, name)
    WHERE deleted_at IS NULL;

-- +goose Down
DROP INDEX IF EXISTS idx_products_shop_status_active;
DROP INDEX IF EXISTS idx_products_not_deleted;
ALTER TABLE products DROP COLUMN IF EXISTS deleted_at;