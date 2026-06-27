-- +goose Up
CREATE TABLE transactions (
    id UUID PRIMARY KEY,
    user_id BIGINT NOT NULL,
    amount BIGINT NOT NULL,
    category TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX idx_transaction_user_created ON transactions (user_id, created_at);

-- +goose Down
DROP TABLE transactions;
