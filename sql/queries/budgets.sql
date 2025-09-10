-- name: GetBudgets :many
SELECT * FROM budgets;

-- name: CreateBudget :one
INSERT INTO budgets (
    target_amount, time_period, start_date, notes, category_id
) VALUES (
    $1, $2, $3, $4, $5
)
RETURNING *;