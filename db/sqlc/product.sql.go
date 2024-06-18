// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: product.sql

package db

import (
	"context"
)

const getProductById = `-- name: GetProductById :one
SELECT id, name, image, description, price
FROM products
WHERE STRCMP(id, ?) = 0
`

type GetProductByIdRow struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Image       string  `json:"image"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

func (q *Queries) GetProductById(ctx context.Context, id string) (GetProductByIdRow, error) {
	row := q.db.QueryRowContext(ctx, getProductById, id)
	var i GetProductByIdRow
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Image,
		&i.Description,
		&i.Price,
	)
	return i, err
}

const getProductTotalCount = `-- name: GetProductTotalCount :one
SELECT COUNT(*) AS "totalProducts"
FROM products
`

func (q *Queries) GetProductTotalCount(ctx context.Context) (int64, error) {
	row := q.db.QueryRowContext(ctx, getProductTotalCount)
	var totalproducts int64
	err := row.Scan(&totalproducts)
	return totalproducts, err
}

const getProductsPaginated = `-- name: GetProductsPaginated :many
SELECT id, name, image, description, price
FROM products
WHERE
    (? = '' OR name LIKE CONCAT('%', ?, '%'))
  AND (
    (? = 0 AND ? = 0)
        OR (price BETWEEN ? AND ?))
ORDER BY
    CASE WHEN ? = 'price' AND ? = 'ASC' THEN price END ,
    CASE WHEN ? = 'price' AND ? = 'DESC' THEN price END DESC ,
    CASE WHEN ? = 'name' AND ? = 'ASC' THEN name END,
    CASE WHEN ? = 'name' AND ? = 'DESC' THEN name END DESC
LIMIT ?
OFFSET ?
`

type GetProductsPaginatedParams struct {
	SearchTerm interface{} `json:"searchTerm"`
	MinPrice   float64     `json:"minPrice"`
	MaxPrice   float64     `json:"maxPrice"`
	SortBy     interface{} `json:"sortBy"`
	SortOrder  interface{} `json:"sortOrder"`
	Limit      int32       `json:"limit"`
	Offset     int32       `json:"offset"`
}

type GetProductsPaginatedRow struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Image       string  `json:"image"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

func (q *Queries) GetProductsPaginated(ctx context.Context, arg GetProductsPaginatedParams) ([]GetProductsPaginatedRow, error) {
	rows, err := q.db.QueryContext(ctx, getProductsPaginated,
		arg.SearchTerm,
		arg.SearchTerm,
		arg.MinPrice,
		arg.MaxPrice,
		arg.MinPrice,
		arg.MaxPrice,
		arg.SortBy,
		arg.SortOrder,
		arg.SortBy,
		arg.SortOrder,
		arg.SortBy,
		arg.SortOrder,
		arg.SortBy,
		arg.SortOrder,
		arg.Limit,
		arg.Offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetProductsPaginatedRow
	for rows.Next() {
		var i GetProductsPaginatedRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Image,
			&i.Description,
			&i.Price,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
