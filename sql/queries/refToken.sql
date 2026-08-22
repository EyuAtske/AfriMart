-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (token_hash, created_at, user_id, expires_at, revoked_at)
VALUES (
    $1,
    Now(),
    $2,
    Now() + INTERVAL '30 days',
    NULL
)
RETURNING *;

-- name: GetRefreshToken :one
SELECT * FROM refresh_tokens
WHERE token_hash = $1;

-- name: RevokeRefreshToken :exec
UPDATE refresh_tokens
SET revoked_at = Now(), updated_at = Now()
WHERE token_hash = $1;x