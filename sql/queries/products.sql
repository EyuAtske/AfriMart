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
    gender,
    price,
    status,
    stock
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
SELECT
    p.*,
    s.name AS shop_name
FROM products p
JOIN shops s ON s.id = p.shop_id
WHERE p.id = $1;


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
    status = $11,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: ListProducts :many

SELECT
    p.*,
    s.name AS shop_name
FROM products p
JOIN shops s ON s.id = p.shop_id
WHERE p.status = 'active'
  AND s.status = 'active'
  AND (
      sqlc.arg(search)::text = ''
      OR p.name ILIKE '%' || sqlc.arg(search)::text || '%'
      OR p.brand ILIKE '%' || sqlc.arg(search)::text || '%'
      OR p.description ILIKE '%' || sqlc.arg(search)::text || '%'
  )
  AND (
      sqlc.narg(category_id)::uuid IS NULL
      OR p.category_id = sqlc.narg(category_id)::uuid
  )
  AND (
      sqlc.narg(subcategory_id)::uuid IS NULL
      OR p.subcategory_id = sqlc.narg(subcategory_id)::uuid
  )
  AND (
      sqlc.arg(brand)::text = ''
      OR p.brand ILIKE '%' || sqlc.arg(brand)::text || '%'
  )
  AND (
      sqlc.arg(color)::text = ''
      OR p.color ILIKE '%' || sqlc.arg(color)::text || '%'
  )
  AND (
      sqlc.arg(size)::text = ''
      OR p.size = sqlc.arg(size)::text
  )
  AND (
      sqlc.narg(min_price)::numeric IS NULL
      OR p.price >= sqlc.narg(min_price)::numeric
  )
  AND (
      sqlc.narg(max_price)::numeric IS NULL
      OR p.price <= sqlc.narg(max_price)::numeric
  )
ORDER BY p.created_at DESC
LIMIT sqlc.arg(page_limit)
OFFSET sqlc.arg(page_offset);

-- name: ListProductsByShop :many

SELECT
    p.*,
    s.name AS shop_name
FROM products p
JOIN shops s ON s.id = p.shop_id
WHERE p.shop_id = $1
  AND p.status = 'active'
  AND s.status = 'active'
ORDER BY p.created_at DESC
LIMIT $2
OFFSET $3;

-- name: ListProductsByCategory :many

SELECT
    p.*,
    s.name AS shop_name
FROM products p
JOIN shops s ON s.id = p.shop_id
WHERE p.category_id = $1
  AND p.status = 'active'
  AND s.status = 'active'
ORDER BY p.created_at DESC
LIMIT $2
OFFSET $3;

-- name: ListProductsBySubcategory :many

SELECT
    p.*,
    s.name AS shop_name
FROM products p
JOIN shops s ON s.id = p.shop_id
WHERE p.subcategory_id = $1
  AND p.status = 'active'
  AND s.status = 'active'
ORDER BY p.created_at DESC
LIMIT $2
OFFSET $3;

-- name: ListCategories :many

SELECT *
FROM categories
ORDER BY name ASC;


-- name: GetCategory :one

SELECT *
FROM categories
WHERE id = $1;


-- name: ListSubcategoriesByCategory :many

SELECT *
FROM subcategories
WHERE category_id = $1
ORDER BY name ASC;