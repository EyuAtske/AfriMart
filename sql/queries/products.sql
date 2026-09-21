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
    $11
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
    status = $11,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: ListProducts :many

SELECT *
FROM products
WHERE status = 'active'
  AND (
      sqlc.arg(search)::text = ''
      OR name ILIKE '%' || sqlc.arg(search)::text || '%'
      OR brand ILIKE '%' || sqlc.arg(search)::text || '%'
      OR description ILIKE '%' || sqlc.arg(search)::text || '%'
  )
  AND (
      sqlc.narg(category_id)::uuid IS NULL
      OR category_id = sqlc.narg(category_id)::uuid
  )
  AND (
      sqlc.narg(subcategory_id)::uuid IS NULL
      OR subcategory_id = sqlc.narg(subcategory_id)::uuid
  )
  AND (
      sqlc.arg(brand)::text = ''
      OR brand ILIKE '%' || sqlc.arg(brand)::text || '%'
  )
  AND (
      sqlc.arg(color)::text = ''
      OR color ILIKE '%' || sqlc.arg(color)::text || '%'
  )
  AND (
      sqlc.arg(size)::text = ''
      OR size = sqlc.arg(size)::text
  )
  AND (
      sqlc.narg(min_price)::numeric IS NULL
      OR price >= sqlc.narg(min_price)::numeric
  )
  AND (
      sqlc.narg(max_price)::numeric IS NULL
      OR price <= sqlc.narg(max_price)::numeric
  )
ORDER BY created_at DESC
LIMIT sqlc.arg(page_limit)
OFFSET sqlc.arg(page_offset);

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