-- name: GetCartItem :one
SELECT id, user_id, product_id, quantity, created_at FROM carts WHERE id = $1;

-- name: ListCartByUser :many
SELECT id, user_id, product_id, quantity, created_at FROM carts WHERE user_id = $1 ORDER BY created_at DESC;

-- name: AddToCart :one
INSERT INTO carts (id, user_id, product_id, quantity) VALUES ($1, $2, $3, $4) RETURNING id, user_id, product_id, quantity, created_at;

-- name: UpdateCartItem :one
UPDATE carts SET quantity = $2 WHERE id = $1 RETURNING id, user_id, product_id, quantity, created_at;

-- name: RemoveFromCart :exec
DELETE FROM carts WHERE id = $1;
