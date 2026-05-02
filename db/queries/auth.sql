-- REFRESH TOKENS

-- name: CreateRefreshToken :exec
INSERT INTO refresh_tokens (id, user_id, role, token_hash, expires_at, created_at)
VALUES ($1, $2, $3, $4, $5, $6);

-- name: GetRefreshTokenByHash :one
SELECT id, user_id, role, token_hash, expires_at, created_at
FROM refresh_tokens
WHERE token_hash = $1;

-- name: DeleteRefreshTokenByTokenHash :exec
DELETE FROM refresh_tokens
WHERE token_hash = $1;

-- name: DeleteRefreshTokensByUserID :exec
DELETE FROM refresh_tokens
WHERE user_id = $1;
