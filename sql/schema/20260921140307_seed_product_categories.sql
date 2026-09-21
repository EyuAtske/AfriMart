-- +goose Up

INSERT INTO categories (name)
VALUES
    ('Clothing'),
    ('Shoes'),
    ('Accessories')
ON CONFLICT (name) DO NOTHING;


-- Clothing
INSERT INTO subcategories (category_id, name)
SELECT categories.id, clothing.name
FROM categories
CROSS JOIN (
    VALUES
        ('T-Shirts'),
        ('Shirts'),
        ('Trousers'),
        ('Jeans'),
        ('Jackets'),
        ('Hoodies'),
        ('Dresses'),
        ('Tops'),
        ('Skirts'),
        ('Sweaters')
) AS clothing(name)
WHERE categories.name = 'Clothing'
ON CONFLICT (category_id, name) DO NOTHING;


-- Shoes
INSERT INTO subcategories (category_id, name)
SELECT categories.id, shoes.name
FROM categories
CROSS JOIN (
    VALUES
        ('Sneakers'),
        ('Formal Shoes'),
        ('Boots'),
        ('Sandals'),
        ('Heels'),
        ('Flats')
) AS shoes(name)
WHERE categories.name = 'Shoes'
ON CONFLICT (category_id, name) DO NOTHING;


-- Accessories
INSERT INTO subcategories (category_id, name)
SELECT categories.id, accessories.name
FROM categories
CROSS JOIN (
    VALUES
        ('Bags'),
        ('Watches'),
        ('Belts'),
        ('Hats'),
        ('Jewelry'),
        ('Scarves')
) AS accessories(name)
WHERE categories.name = 'Accessories'
ON CONFLICT (category_id, name) DO NOTHING;


-- +goose Down

DELETE FROM subcategories
WHERE name IN (
    'T-Shirts', 'Shirts', 'Trousers', 'Jeans', 'Jackets', 'Hoodies',
    'Dresses', 'Tops', 'Skirts', 'Sweaters', 'Sneakers', 'Formal Shoes',
    'Boots', 'Sandals', 'Heels', 'Flats', 'Bags', 'Watches', 'Belts',
    'Hats', 'Jewelry', 'Scarves'
);

DELETE FROM categories
WHERE name IN (
    'Clothing',
    'Shoes',
    'Accessories'
);