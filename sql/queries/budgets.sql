-- name: GetBudgets :many
SELECT * FROM budgets;

-- name: GetBudgetById :one
SELECT * FROM budgets
WHERE id = $1;

-- name: CreateBudget :one
INSERT INTO budgets (
    target_amount, time_period, start_date, notes, category_id
) VALUES (
    $1, $2, $3, $4,
    (SELECT id FROM categories WHERE name = $5)
)
RETURNING *;

-- name: UpdateBudget :one
UPDATE budgets
SET target_amount = $2, time_period = $3, start_date = $4, notes = $5, category_id = (SELECT id FROM categories WHERE name = $6)
WHERE budgets.id = $1
RETURNING *;

-- name: DeleteBudget :exec
DELETE FROM budgets
WHERE id = $1;

-- name: GetBudgetStatus :one
WITH budget_info AS (
    SELECT 
        id AS budget_id,
        category_id, 
        time_period,
        target_amount,
        start_date, 
        CAST(start_date + 
            CASE time_period
                WHEN 'weekly' THEN INTERVAL '7 DAY'
                WHEN 'monthly' THEN INTERVAL '1 MONTH'
                WHEN 'bi-monthly' THEN INTERVAL '2 MONTHS'
                WHEN 'quarterly' THEN INTERVAL '3 MONTHS'
                WHEN 'yearly' THEN INTERVAL '1 YEAR'
            END AS DATE) AS end_date 
    FROM budgets
    WHERE budgets.id = $1
),
transactions_sum AS (
    SELECT 
        bi.budget_id,
        bi.category_id,
        bi.time_period,
        bi.start_date,
        bi.end_date,
        bi.target_amount,
        CAST(COALESCE(SUM(t.amount), 0) AS FLOAT) AS current_spent
    FROM transactions t
    JOIN budget_info bi ON t.category_id = bi.category_id
    WHERE t.date BETWEEN bi.start_date AND bi.end_date
    GROUP BY
        bi.budget_id,
        bi.category_id,
        bi.time_period,
        bi.start_date,
        bi.end_date,
        bi.target_amount
)
SELECT
    budget_id,
    category_id,
    time_period,
    start_date,
    end_date,
    target_amount,
    current_spent
FROM transactions_sum;

-- name: GetAllBudgetStatuses :many
WITH budget_info AS (
    SELECT 
        id AS budget_id,
        category_id, 
        time_period,
        target_amount,
        start_date, 
        CAST(start_date + 
            CASE time_period
                WHEN 'weekly' THEN INTERVAL '7 DAY'
                WHEN 'monthly' THEN INTERVAL '1 MONTH'
                WHEN 'bi-monthly' THEN INTERVAL '2 MONTHS'
                WHEN 'quarterly' THEN INTERVAL '3 MONTHS'
                WHEN 'yearly' THEN INTERVAL '1 YEAR'
            END AS DATE) AS end_date 
    FROM budgets
),
transactions_sum AS (
    SELECT 
        bi.budget_id,
        bi.category_id,
        bi.time_period,
        bi.start_date,
        bi.end_date,
        bi.target_amount,
        CAST(COALESCE(SUM(t.amount), 0) AS FLOAT) AS current_spent
    FROM transactions t
    JOIN budget_info bi ON t.category_id = bi.category_id
    WHERE t.date BETWEEN bi.start_date AND bi.end_date
    GROUP BY
        bi.budget_id,
        bi.category_id,
        bi.time_period,
        bi.start_date,
        bi.end_date,
        bi.target_amount
)
SELECT
    budget_id,
    category_id,
    time_period,
    start_date,
    end_date,
    target_amount,
    current_spent
FROM transactions_sum;