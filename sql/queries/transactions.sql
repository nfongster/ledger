-- name: GetOrCreateCategory :one
INSERT INTO categories (name)
VALUES ($1)
ON CONFLICT (name) DO UPDATE
SET name = excluded.name
RETURNING id;

-- name: CreateTransaction :one
INSERT INTO transactions (
    date, description, amount, notes, category_id
) VALUES (
    $1, $2, $3, $4, $5
)
RETURNING *;

-- name: GetTransactionById :one
SELECT t.id, t.date, t.description, t.amount, t.notes, c.name AS category
FROM transactions AS t
JOIN categories AS c ON t.category_id = c.id
WHERE t.id = $1;

-- name: GetTransactionsByCategory :many
SELECT t.id, t.date, t.description, t.amount, t.notes, c.name AS category
FROM transactions AS t
JOIN categories AS c ON t.category_id = c.id
WHERE c.name = $1;

-- name: GetAllTransactions :many
SELECT t.id, t.date, t.description, t.amount, t.notes, c.name AS category
FROM transactions AS t
JOIN categories AS c ON t.category_id = c.id;

-- name: DeleteAllTransactions :exec
DELETE FROM transactions;