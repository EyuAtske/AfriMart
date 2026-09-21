-- +goose Up
ALTER TABLE products 
ADD COLUMN gender VARCHAR(20) NOT NULL
    CHECK (gender IN ('men', 'women', 'kids'));

-- +goose Down
ALTER TABLE products 
DROP COLUMN gender;
