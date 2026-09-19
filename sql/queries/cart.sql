-- name: GetCartByUserID :one
SELECT *
FROM carts
WHERE user_id = $1;

-- name: CreateCart :one
INSERT INTO carts (
    user_id
)
VALUES ($1)
RETURNING *;

-- name: GetCartItems :many
SELECT
    ci.id,
    ci.cart_id,
    ci.product_id,
    ci.quantity,
    ci.created_at,
    ci.updated_at,
    p.name AS product_name,
    p.price,
    p.stock,
    p.image,
    p.status
FROM cart_items ci
JOIN products p ON p.id = ci.product_id
WHERE ci.cart_id = $1
ORDER BY ci.created_at DESC;

-- name: AddCartItem :one
INSERT INTO cart_items (
    cart_id,
    product_id,
    quantity
)
VALUES ($1, $2, $3)
ON CONFLICT (cart_id, product_id)
DO UPDATE SET
    quantity = cart_items.quantity + EXCLUDED.quantity,
    updated_at = NOW()
RETURNING *;

-- name: UpdateCartItemQuantity :one
UPDATE cart_items
SET
    quantity = $3,
    updated_at = NOW()
WHERE id = $1
  AND cart_id = $2
RETURNING *;

-- name: DeleteCartItem :exec
DELETE FROM cart_items
WHERE id = $1
  AND cart_id = $2;

-- name: ClearCart :exec
DELETE FROM cart_items
WHERE cart_id = $1;

-- name: UpdateCartTimestamp :exec
UPDATE carts
SET updated_at = NOW()
WHERE id = $1;