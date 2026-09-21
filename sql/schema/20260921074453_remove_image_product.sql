-- +goose Up
ALTER TABLE products DROP COLUMN image;

-- +goose Down
ALTER TABLE products ADD COLUMN image TEXT;