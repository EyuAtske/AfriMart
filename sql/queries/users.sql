-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, password_hash, email, first_name, last_name, username, role)
VALUES (
    gen_random_uuid(),
    Now(),
    Now(),
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
)
RETURNING *;

-- name: DeleteUsers :exec
DELETE FROM users;      

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1;

-- name: UpdateUserPassword :one
UPDATE users
SET password_hash = $1, updated_at = Now()
WHERE id = $2
RETURNING *;

-- name: UpdateUsername :one
UPDATE users
SET username = $1, updated_at = Now()
WHERE id = $2
RETURNING *;