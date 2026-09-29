-- +goose Up
CREATE TABLE payments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    order_id UUID NOT NULL REFERENCES orders(id) ON DELETE RESTRICT,

    payment_method VARCHAR(20) NOT NULL
        CHECK (payment_method IN ('online', 'cod')),

    payment_status VARCHAR(20) NOT NULL DEFAULT 'pending'
        CHECK (payment_status IN (
            'pending',
            'processing',
            'successful',
            'failed',
            'cancelled'
        )),

    amount NUMERIC(12, 2) NOT NULL CHECK (amount >= 0),

    provider VARCHAR(50),

    transaction_id VARCHAR(255),

    provider_reference VARCHAR(255),

    failure_reason TEXT,

    paid_at TIMESTAMPTZ,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_payments_order_id
    ON payments(order_id);

CREATE INDEX idx_payments_transaction_id
    ON payments(transaction_id);

-- +goose Down
DROP TABLE IF EXISTS payments;
