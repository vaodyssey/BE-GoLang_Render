/* name: GetProductById :one */
SELECT id, name, image, description, price
FROM products
WHERE STRCMP(id, sqlc.arg(id)) = 0;

/* name: GetProductsPaginated :many */
SELECT id, name, image, description, price
FROM products
WHERE name LIKE CONCAT('%', CAST(sqlc.arg(search_term) AS char), '%')
ORDER BY
    CASE WHEN sqlc.arg(sort_by) = 'price' AND sqlc.arg(sort_order) = 'ASC' THEN price END,
    CASE WHEN sqlc.arg(sort_by)= 'price' AND sqlc.arg(sort_order) = 'DESC' THEN price END DESC,
    CASE WHEN sqlc.arg(sort_by) = 'name' AND sqlc.arg(sort_order) = 'ASC' THEN name END,
    CASE WHEN sqlc.arg(sort_by) = 'name' AND sqlc.arg(sort_order) = 'DESC' THEN name END DESC
LIMIT ?
OFFSET ?;

/* name: GetProductTotalCount :one */
SELECT COUNT(*) AS "totalProducts"
FROM products