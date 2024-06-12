/* name: GetProductById :one */
SELECT id, name, image, description, price
FROM products
WHERE STRCMP(id, sqlc.arg(id)) = 0;

/* name: GetProductsPaginated :many */
SELECT id, name, image, description, price
FROM products
ORDER BY created_at desc
LIMIT ?
OFFSET ?;

/* name: GetProductTotalCount :one */
SELECT COUNT(*) AS "totalProducts"
FROM products