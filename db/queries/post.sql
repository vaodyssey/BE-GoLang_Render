-- name: GetProductById :one
SELECT id, name, image, description, price
FROM products
WHERE STRCMP(id, sqlc.arg(id)) = 0;