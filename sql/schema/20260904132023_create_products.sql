-- +goose Up

CREATE TABLE products (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    shop_id UUID NOT NULL
        REFERENCES shops(id)
        ON DELETE CASCADE,

    category_id UUID NOT NULL,

    subcategory_id UUID NOT NULL,

    name VARCHAR(255) NOT NULL,

    description TEXT,

    brand VARCHAR(100),

    color VARCHAR(50),

    size VARCHAR(50),

    price NUMERIC(12, 2) NOT NULL
        CHECK (price >= 0),

    stock INTEGER NOT NULL DEFAULT 0
        CHECK (stock >= 0),

    image TEXT,

    status VARCHAR(20) NOT NULL DEFAULT 'active'
        CHECK (status IN ('active', 'inactive')),

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT products_category_subcategory_fk
        FOREIGN KEY (subcategory_id, category_id)
        REFERENCES subcategories(id, category_id)
);

CREATE INDEX idx_products_shop_id
    ON products(shop_id);

CREATE INDEX idx_products_category_id
    ON products(category_id);

CREATE INDEX idx_products_subcategory_id
    ON products(subcategory_id);

CREATE INDEX idx_products_status
    ON products(status);


-- +goose Down

DROP TABLE IF EXISTS products;