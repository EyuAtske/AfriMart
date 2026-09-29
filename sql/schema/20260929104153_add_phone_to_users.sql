-- +goose Up
ALTER TABLE users 
ADD COLUMN phone_number VARCHAR(20);

CREATE INDEX idx_users_phone_number ON users(phone_number);

-- +goose Down
DROP INDEX IF EXISTS idx_users_phone_number;
ALTER TABLE users DROP COLUMN phone_number;