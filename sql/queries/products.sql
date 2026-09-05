-- name: CreateProduct :one

INSERT INTO products (
    shop_id,
    category_id,
    subcategory_id,
    name,
    description,
    brand,
    color,
    size,
    price,
    stock,
    image,
    status
)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    $8,
    $9,
    $10,
    $11,
    $12
)
RETURNING *;


-- name: GetProduct :one

SELECT *
FROM products
WHERE id = $1;


-- name: DeleteProduct :exec

DELETE FROM products
WHERE id = $1;

-- name: UpdateProduct :one

UPDATE products
SET
    category_id = $2,
    subcategory_id = $3,
    name = $4,
    description = $5,
    brand = $6,
    color = $7,
    size = $8,
    price = $9,
    stock = $10,
    image = $11,
    status = $12,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: ListProducts :many

SELECT *
FROM products
WHERE status = 'active'
ORDER BY created_at DESC
LIMIT $1
OFFSET $2;

-- name: ListProductsByShop :many

SELECT *
FROM products
WHERE shop_id = $1
ORDER BY created_at DESC
LIMIT $2
OFFSET $3;

-- name: ListProductsByCategory :many

SELECT *
FROM products
WHERE category_id = $1
  AND status = 'active'
ORDER BY created_at DESC
LIMIT $2
OFFSET $3;

-- name: ListProductsBySubcategory :many

SELECT *
FROM products
WHERE subcategory_id = $1
  AND status = 'active'
ORDER BY created_at DESC
LIMIT $2
OFFSET $3;