-- name: GetProduct :one
SELECT id, name, description, price, created_at FROM products WHERE id = $1;

-- name: ListProducts :many
SELECT id, name, description, price, created_at FROM products ORDER BY created_at DESC;

-- name: CreateProduct :one
INSERT INTO products (id, name, description, price) VALUES ($1, $2, $3, $4) RETURNING id, name, description, price, created_at;

-- name: UpdateProduct :one
UPDATE products SET name = $2, description = $3, price = $4 WHERE id = $1 RETURNING id, name, description, price, created_at;

-- name: DeleteProduct :exec
DELETE FROM products WHERE id = $1;
