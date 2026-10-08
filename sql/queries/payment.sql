-- name: CreatePayment :one
INSERT INTO payments (
    order_id,
    payment_method,
    payment_status,
    amount,
    provider,
    transaction_id
) VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
) RETURNING *;

-- name: GetPaymentByProviderRef :one
SELECT * FROM payments 
WHERE provider_reference = $1 
LIMIT 1;

-- name: UpdatePaymentStatusConditional :execrows
UPDATE payments
SET
    payment_status = sqlc.arg(new_status),
    paid_at = COALESCE(sqlc.narg(paid_at), paid_at),
    provider_reference = COALESCE(sqlc.narg(provider_reference), provider_reference),
    transaction_id = COALESCE(sqlc.narg(transaction_id), transaction_id),
    failure_reason = COALESCE(sqlc.narg(failure_reason), failure_reason),
    updated_at = NOW()
WHERE id = sqlc.arg(id) AND payment_status = sqlc.arg(current_status);

-- name: GetExpiredPendingPayments :many
SELECT id, order_id 
FROM payments 
WHERE payment_method = 'online' 
  AND payment_status = 'pending' 
  AND created_at < $1;

-- name: GetPaymentByOrderID :one
SELECT * FROM payments
WHERE order_id = $1
LIMIT 1;