-- +goose Up
-- +goose StatementBegin
CREATE TABLE transactions(
    id UUID PRIMARY KEY,
    date DATE,
    description TEXT NOT NULL,
    amount FLOAT,
    category TEXT NOT NULL,
    notes TEXT
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE transactions;
-- +goose StatementEnd