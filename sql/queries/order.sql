-- name: CreateOrder :one
INSERT INTO orders (
    user_id,
    subtotal,
    status
)
VALUES ($1, $2, 'pending')
RETURNING *;


-- name: CreateOrderItem :one
INSERT INTO order_items (
    order_id,
    product_id,
    quantity,
    price
)
VALUES ($1, $2, $3, $4)
RETURNING *;


-- name: GetOrderByID :one
SELECT *
FROM orders
WHERE id = $1
  AND user_id = $2;


-- name: GetOrderItems :many
SELECT
    oi.id,
    oi.order_id,
    oi.product_id,
    oi.quantity,
    oi.price,
    oi.created_at,
    p.name AS product_name,
    p.image AS product_image
FROM order_items oi
JOIN products p ON p.id = oi.product_id
WHERE oi.order_id = $1
ORDER BY oi.created_at ASC;


-- name: ListOrdersByUser :many
SELECT *
FROM orders
WHERE user_id = $1
ORDER BY created_at DESC
LIMIT $2
OFFSET $3;


-- name: UpdateOrderStatus :one
UPDATE orders
SET
    status = $2,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: ReduceProductStock :one
UPDATE products
SET
    stock = stock - $2,
    updated_at = NOW()
WHERE id = $1
  AND stock >= $2
RETURNING id, stock;