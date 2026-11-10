-- +goose Up
ALTER TABLE transactions ADD CONSTRAINT transactions_amount_positive CHECK(amount > 0);

-- +goose Down
ALTER TABLE transactions DROP CONSTRAINT transactions_amount_positive;