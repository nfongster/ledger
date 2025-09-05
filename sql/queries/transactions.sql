-- name: CreateTransaction :one
INSERT INTO transactions (
    id, date, description, amount, category, notes
) VALUES (
    $1, $2, $3, $4, $5, $6
)
RETURNING *;

-- name: GetTransactionById :one
SELECT * FROM transactions
WHERE id = $1;

-- name: GetAllTransactions :many
SELECT * FROM transactions;

-- name: DeleteAllTransactions :exec
DELETE FROM transactions;