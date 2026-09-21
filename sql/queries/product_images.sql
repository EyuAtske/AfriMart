-- name: CreateProductImage :one
INSERT INTO product_images (
    product_id,
    object_key,
    display_order
)
VALUES (
    $1,
    $2,
    $3
)
RETURNING *;

-- name: GetProductImages :many
SELECT *
FROM product_images
WHERE product_id = $1
ORDER BY display_order ASC, created_at ASC;

-- name: GetProductImage :one
SELECT *
FROM product_images
WHERE id = $1
  AND product_id = $2;

-- name: DeleteProductImage :one
DELETE FROM product_images
WHERE id = $1
  AND product_id = $2
RETURNING *;

-- name: UpdateProductImage :one
UPDATE product_images
SET
    object_key = $3
WHERE id = $1
  AND product_id = $2
RETURNING *;