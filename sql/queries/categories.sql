-- name: DeleteAllCategories :exec
DELETE FROM categories;

-- name: TruncateAllTables :exec
TRUNCATE TABLE transactions, categories RESTART IDENTITY CASCADE;