-- +goose Up

CREATE TABLE categories (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE subcategories (
    id UUID PRIMARY KEY,
    category_id UUID NOT NULL REFERENCES categories(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),

    CONSTRAINT subcategories_category_name_unique
        UNIQUE (category_id, name)
    
    CONSTRAINT subcategories_id_category_unique
        UNIQUE (id, category_id)
);

-- +goose Down

DROP TABLE subcategories;
DROP TABLE categories;