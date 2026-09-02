-- +goose Up

CREATE TABLE shops (
    id UUID PRIMARY KEY,
    owner_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    description TEXT,
    status TEXT NOT NULL DEFAULT 'active',
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),

    CONSTRAINT shops_status_check
        CHECK (status IN ('active', 'inactive')),

    CONSTRAINT shops_owner_unique
        UNIQUE (owner_id)
);

-- +goose Down

DROP TABLE shops;