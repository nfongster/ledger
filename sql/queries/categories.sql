-- name: DeleteAllCategories :exec
DELETE FROM categories;

-- name: TruncateAllTables :exec
TRUNCATE TABLE transactions, categories RESTART IDENTITY CASCADE;

-- name: GetAllCategories :many
SELECT * FROM categories;

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

-- name: GetSpendingBetweenStartAndEnd :one
SELECT CAST(SUM(amount) AS FLOAT) FROM transactions
WHERE category_id = $1 AND date BETWEEN sqlc.arg(Start_Date) AND sqlc.arg(End_Date);

-- name: GetSpendingSinceStart :one
SELECT CAST(SUM(amount) AS FLOAT) FROM transactions
WHERE category_id = $1 AND date >= sqlc.arg(Start_Date);

-- name: GetSpendingUntilEnd :one
SELECT CAST(SUM(amount) AS FLOAT) FROM transactions
WHERE category_id = $1 AND date <= sqlc.arg(End_Date);

-- name: GetSpendingAllTime :one
SELECT CAST(SUM(amount) AS FLOAT) FROM transactions
WHERE category_id = $1;