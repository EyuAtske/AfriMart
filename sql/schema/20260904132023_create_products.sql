-- +goose Up

CREATE TABLE products (
    id UUID PRIMARY KEY,
    shop_id UUID NOT NULL REFERENCES shops(id) ON DELETE CASCADE,

    category_id UUID NOT NULL,
    subcategory_id UUID NOT NULL,

    name TEXT NOT NULL,
    description TEXT,
    brand TEXT,
    color TEXT,
    size TEXT,

    price NUMERIC(12, 2) NOT NULL,
    stock_quantity INTEGER NOT NULL DEFAULT 0,

    gender TEXT NOT NULL,

    status TEXT NOT NULL DEFAULT 'active',

    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),

    CONSTRAINT products_price_check
        CHECK (price >= 0),

    CONSTRAINT products_stock_check
        CHECK (stock_quantity >= 0),

    CONSTRAINT products_gender_check
        CHECK (gender IN ('men', 'women', 'kids')),

    CONSTRAINT products_status_check
        CHECK (status IN ('active', 'inactive'))

    CONSTRAINT products_category_subcategory_fk
        FOREIGN KEY (subcategory_id, category_id)
        REFERENCES subcategories(id, category_id)
);

-- +goose Down

DROP TABLE products;