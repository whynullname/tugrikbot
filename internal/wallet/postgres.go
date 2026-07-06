package wallet

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/whynullname/tugrikbot/internal/domain"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (p *PostgresRepository) GetUserWalletID(ctx context.Context, userID uuid.UUID) (uuid.UUID, error) {
	row := p.db.QueryRowContext(ctx, `SELECT wallet_id FROM wallet_members WHERE user_id = $1`, userID)
	var id uuid.UUID
	err := row.Scan(&id)
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

func (p *PostgresRepository) CreateNewWallet(ctx context.Context, userID uuid.UUID, wallet *domain.Wallet) error {
	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	_, err = tx.ExecContext(ctx, `INSERT INTO wallets (id, title, created_at, invite_code) VALUES ($1, $2, $3, $4)`,
		wallet.ID, wallet.Title, wallet.CreatedAt, wallet.InviteCode)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `INSERT INTO wallet_members (wallet_id, user_id, created_at, role) VALUES ($1, $2, $3, $4)`,
		wallet.ID, userID, wallet.CreatedAt, domain.OwnerWalletRole)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
