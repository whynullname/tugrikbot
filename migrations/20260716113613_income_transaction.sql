-- +goose Up
ALTER TABLE transactions ADD COLUMN type TEXT CHECK(type IN ('expense', 'income'));
UPDATE transactions SET type = 'expense';
ALTER TABLE transactions ALTER COLUMN type SET NOT NULL;


-- +goose Down
ALTER TABLE transactions DROP COLUMN type;
