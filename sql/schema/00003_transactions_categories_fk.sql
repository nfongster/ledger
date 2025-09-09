-- +goose Up
-- +goose StatementBegin
ALTER TABLE transactions
ADD CONSTRAINT fk_transactions_category_id
FOREIGN KEY (category_id)
REFERENCES categories (id)
ON DELETE RESTRICT;

ALTER TABLE transactions
DROP COLUMN category;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE transactions
ADD COLUMN category VARCHAR(255) NOT NULL;

UPDATE transactions
SET category = categories.name FROM categories
WHERE transactions.category_id = categories.id;

ALTER TABLE transactions
DROP CONSTRAINT fk_transactions_category_id;
-- +goose StatementEnd
