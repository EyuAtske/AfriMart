-- name: CreateOrder :one
INSERT INTO orders (
    user_id,
    subtotal,
    status,
    recipient_name,
    phone,
    delivery_address,
    delivery_city,
    delivery_notes
)
VALUES (
    $1,
    $2,
    'pending',
    $3,
    $4,
    $5,
    $6,
    $7
)
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

-- name: VerifyOrderSellerOwnership :one
SELECT o.status
FROM orders o
JOIN order_items oi ON oi.order_id = o.id
JOIN products p ON p.id = oi.product_id
JOIN shops s ON s.id = p.shop_id
WHERE o.id = $1
  AND s.owner_id = $2
LIMIT 1;

-- name: ListOrdersBySeller :many
SELECT DISTINCT o.*
FROM orders o
JOIN order_items oi ON oi.order_id = o.id
JOIN products p ON p.id = oi.product_id
JOIN shops s ON s.id = p.shop_id
WHERE s.owner_id = $1
ORDER BY o.created_at DESC
LIMIT $2 OFFSET $3;

-- name: GetCartByUserIDForUpdate :one
SELECT *
FROM carts
WHERE user_id = $1
FOR UPDATE;