-- +goose Up
ALTER TABLE wallet_members ADD COLUMN role TEXT CHECK (role IN ('owner', 'member'));
UPDATE wallet_members SET role = 'owner';
ALTER TABLE wallet_members ALTER COLUMN role SET NOT NULL;

ALTER TABLE wallets ADD COLUMN invite_code UUID UNIQUE;
UPDATE wallets SET invite_code = gen_random_uuid();
ALTER TABLE wallets ALTER COLUMN invite_code SET NOT NULL;

CREATE TABLE users_settings (
    user_id UUID PRIMARY KEY REFERENCES users(id) NOT NULL,
    active_wallet_id UUID REFERENCES wallets(id) ON DELETE SET NULL
);

-- +goose Down
DROP TABLE users_settings;
ALTER TABLE wallets DROP invite_code; 
ALTER TABLE wallet_members DROP role; 
