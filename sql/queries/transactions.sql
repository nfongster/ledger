-- name: GetOrCreateCategory :one
WITH existing_category AS (
    SELECT id FROM categories WHERE categories.name = $1
),
inserted_category AS (
    INSERT INTO categories (name)
    SELECT $1
    WHERE NOT EXISTS (SELECT 1 FROM existing_category)
    RETURNING id
)
SELECT id FROM existing_category
UNION ALL
SELECT id FROM inserted_category;

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