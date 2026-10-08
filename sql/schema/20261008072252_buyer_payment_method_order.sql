-- +goose Up
ALTER TABLE orders
ADD COLUMN method VARCHAR(20) NOT NULL DEFAULT 'cod'
        CHECK (method IN (
            'cod',
            'chapa'
        ));

-- +goose Down
AlTER TABLE orders
DROP COLUMN method;
