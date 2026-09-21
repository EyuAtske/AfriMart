-- +goose Up

ALTER TABLE orders
ADD COLUMN recipient_name VARCHAR(255) NOT NULL,
ADD COLUMN phone VARCHAR(30) NOT NULL,
ADD COLUMN delivery_address TEXT NOT NULL,
ADD COLUMN delivery_city VARCHAR(100) NOT NULL,
ADD COLUMN delivery_notes TEXT;

-- +goose Down

ALTER TABLE orders
DROP COLUMN delivery_notes,
DROP COLUMN delivery_city,
DROP COLUMN delivery_address,
DROP COLUMN phone,
DROP COLUMN recipient_name;