-- +goose Up
ALTER TABLE users_settings ADD COLUMN personal_wallet_id UUID REFERENCES wallets(id) ON DELETE RESTRICT;
INSERT INTO users_settings (user_id, active_wallet_id, personal_wallet_id) SELECT user_id, wallet_id, wallet_id FROM wallet_members 
    ON CONFLICT (user_id) DO UPDATE SET personal_wallet_id = EXCLUDED.personal_wallet_id;
ALTER TABLE users_settings ALTER COLUMN personal_wallet_id SET NOT NULL;

-- +goose Down
ALTER TABLE users_settings DROP COLUMN personal_wallet_id;
