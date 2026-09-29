-- name: CreateUser :one
INSERT INTO users (
    id, 
    created_at, 
    updated_at, 
    password_hash, 
    email, 
    first_name, 
    last_name, 
    username,
    phone_number 
)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
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

-- name: GetUserByID :one
SELECT email, username
FROM users
WHERE id = $1;

-- name: GetUserByIDFull :one
SELECT * FROM users WHERE id = $1;