-- name: GetUser :one
SELECT id, name, email, created_at FROM users WHERE id = $1;

-- name: ListUsers :many
SELECT id, name, email, created_at FROM users ORDER BY created_at DESC;

-- name: CreateUser :one
INSERT INTO users (id, name, email) VALUES ($1, $2, $3) RETURNING id, name, email, created_at;

-- name: UpdateUser :one
UPDATE users SET name = $2, email = $3 WHERE id = $1 RETURNING id, name, email, created_at;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1;
