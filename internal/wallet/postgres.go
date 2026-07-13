package wallet

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/whynullname/tugrikbot/internal/domain"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (p *PostgresRepository) GetActiveUserWalletID(ctx context.Context, userID uuid.UUID) (uuid.UUID, error) {
	row := p.db.QueryRowContext(ctx, `
		SELECT COALESCE(active_wallet_id, personal_wallet_id) FROM users_settings WHERE user_id = $1`, userID)
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

func (p *PostgresRepository) JoinToWallet(ctx context.Context, userID uuid.UUID, walletInviteCode uuid.UUID) error {
	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	row := tx.QueryRowContext(ctx, `SELECT id FROM wallets WHERE invite_code = $1`, walletInviteCode)
	walletID := uuid.Nil
	err = row.Scan(&walletID)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `INSERT INTO wallet_members (wallet_id, user_id, created_at, role) VALUES ($1, $2, $3, $4)`,
		walletID, userID, time.Now(), domain.MemberWalletRole)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return ErrUserAlreadyInWallet
		}
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (p *PostgresRepository) GetUserWallets(ctx context.Context, userID uuid.UUID) ([]domain.WalletInfo, error) {
	rows, err := p.db.QueryContext(ctx, `
		SELECT w.title, w.id, COALESCE(w.id = us.active_wallet_id) AS is_active
		FROM wallets AS w
		JOIN wallet_members AS wm ON w.id = wm.wallet_id
		JOIN users_settings AS us ON us.user_id = wm.user_id
		WHERE wm.user_id = $1`, userID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	output := make([]domain.WalletInfo, 0)
	for rows.Next() {
		var walletInfo domain.WalletInfo
		err = rows.Scan(&walletInfo.Title, &walletInfo.ID, &walletInfo.IsActive)
		if err != nil {
			return nil, err
		}

		output = append(output, walletInfo)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return output, nil
}

func (p *PostgresRepository) SetActiveWallet(ctx context.Context, userID uuid.UUID, walletID uuid.UUID) error {
	result, err := p.db.ExecContext(ctx, `UPDATE users_settings SET active_wallet_id = $1 WHERE user_id = $2
		AND EXISTS (SELECT 1 FROM wallet_members WHERE wallet_id = $1 AND user_id = $2)`, walletID, userID)
	if err != nil {
		return err
	}

	affectedRows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affectedRows == 0 {
		return ErrUserNotInWallet
	}

	return nil
}
