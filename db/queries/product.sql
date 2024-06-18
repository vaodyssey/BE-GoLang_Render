/* name: GetProductById :one */
SELECT id, name, image, description, price
FROM products
WHERE STRCMP(id, sqlc.arg(id)) = 0;

/* name: GetProductsPaginated :many */
SELECT id, name, image, description, price
FROM products
WHERE
    (sqlc.arg(search_term) = '' OR name LIKE CONCAT('%', sqlc.arg(search_term), '%'))
  AND (
    (sqlc.arg(min_price) = 0 AND sqlc.arg(max_price) = 0)
        OR (price BETWEEN sqlc.arg(min_price) AND sqlc.arg(max_price)))
ORDER BY
    CASE WHEN sqlc.arg(sort_by) = 'price' AND sqlc.arg(sort_order) = 'ASC' THEN price END ,
    CASE WHEN sqlc.arg(sort_by) = 'price' AND sqlc.arg(sort_order) = 'DESC' THEN price END DESC ,
    CASE WHEN sqlc.arg(sort_by) = 'name' AND sqlc.arg(sort_order) = 'ASC' THEN name END,
    CASE WHEN sqlc.arg(sort_by) = 'name' AND sqlc.arg(sort_order) = 'DESC' THEN name END DESC
LIMIT ?
OFFSET ?;

/* name: GetProductTotalCount :one */
SELECT COUNT(*) AS "totalProducts"
FROM products;