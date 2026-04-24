-- name: CreateUser :exec
INSERT INTO users (id, name, email, password, role, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7);

-- name: GetUserByEmail :one
SELECT id, name, email, password, role, created_at, updated_at
FROM users
WHERE lower(btrim(email)) = lower(btrim($1))
LIMIT 1;

-- name: GetUserByID :one
SELECT id, name, email, password, role, created_at, updated_at
FROM users
WHERE id = $1
LIMIT 1;
