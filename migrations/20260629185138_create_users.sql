-- +goose Up
CREATE TABLE users (
    id UUID PRIMARY KEY,
    telegram_id BIGINT NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE wallets (
    id UUID PRIMARY KEY,
    title TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE wallet_members (
    wallet_id UUID REFERENCES wallets(id),
    user_id UUID REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL,
    PRIMARY KEY (wallet_id, user_id)
);

DELETE FROM transactions;
DROP INDEX idx_transaction_user_created;
ALTER TABLE transactions DROP COLUMN user_id;
ALTER TABLE transactions ADD COLUMN user_id UUID REFERENCES users(id) NOT NULL;
ALTER TABLE transactions ADD COLUMN wallet_id UUID REFERENCES wallets(id) NOT NULL;
CREATE INDEX idx_transaction_wallet_created_time ON transactions (wallet_id, created_at);


-- +goose Down
DROP INDEX idx_transaction_wallet_created_time;
ALTER TABLE transactions DROP COLUMN user_id;
ALTER TABLE transactions DROP COLUMN wallet_id;
ALTER TABLE transactions ADD COLUMN user_id BIGINT NOT NULL;
CREATE INDEX idx_transaction_user_created ON transactions (user_id, created_at);
DROP TABLE wallet_members;
DROP TABLE users;
DROP TABLE wallets;

