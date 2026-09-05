-- +goose Up

CREATE TABLE categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    name VARCHAR(100) NOT NULL UNIQUE,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE subcategories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    category_id UUID NOT NULL
        REFERENCES categories(id)
        ON DELETE CASCADE,

    name VARCHAR(100) NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT subcategories_category_name_unique
        UNIQUE (category_id, name),

    CONSTRAINT subcategories_id_category_unique
        UNIQUE (id, category_id)
);

CREATE INDEX idx_subcategories_category_id
    ON subcategories(category_id);


-- +goose Down

DROP TABLE IF EXISTS subcategories;
DROP TABLE IF EXISTS categories;