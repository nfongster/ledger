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

-- name: DeleteTransaction :exec
DELETE FROM transactions
WHERE id = $1;

-- name: UpdateTransaction :one
UPDATE transactions
SET date = $2, description = $3, amount = $4, category_id = $5, notes = $6
WHERE id = $1
RETURNING *;