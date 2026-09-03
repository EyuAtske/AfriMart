-- name: CreateShop :one
INSERT INTO shops (
    id,
    owner_id,
    name,
    description
)
VALUES (
    $1,
    $2,
    $3,
    $4
)
RETURNING *;


-- name: GetShopByID :one
SELECT *
FROM shops
WHERE id = $1
LIMIT 1;


-- name: GetShopByOwnerID :one
SELECT *
FROM shops
WHERE owner_id = $1
LIMIT 1;


-- name: UpdateShopName :one
UPDATE shops
SET
    name = $2,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: UpdateShopDescription :one
UPDATE shops
SET
    description = $2,
    updated_at = NOW()
WHERE id = $1
RETURNING *;


-- name: DeactivateShop :one
UPDATE shops
SET
    status = 'inactive',
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: ActivateShop :one
UPDATE shops
SET
    status = 'active',
    updated_at = NOW()
WHERE id = $1
RETURNING *;