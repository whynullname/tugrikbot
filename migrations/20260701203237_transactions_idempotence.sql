-- +goose Up
DELETE FROM transactions;
ALTER TABLE transactions ADD COLUMN update_id BIGINT UNIQUE NOT NULL;

-- +goose Down
ALTER TABLE transactions DROP update_id;
