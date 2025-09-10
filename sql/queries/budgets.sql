-- name: GetBudgets :many
SELECT * FROM budgets;

-- name: GetBudgetById :one
SELECT * FROM budgets
WHERE id = $1;

-- name: CreateBudget :one
INSERT INTO budgets (
    target_amount, time_period, start_date, notes, category_id
) VALUES (
    $1, $2, $3, $4, $5
)
RETURNING *;

-- name: UpdateBudget :one
UPDATE budgets
SET target_amount = $2, time_period = $3, start_date = $4, notes = $5, category_id = $6
WHERE id = $1
RETURNING *;

-- name: DeleteBudget :exec
DELETE FROM budgets
WHERE id = $1;