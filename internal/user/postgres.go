package user

import (
	"context"
	"database/sql"
	"errors"

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

func (p *PostgresRepository) IsUserCreated(ctx context.Context, telegramUserId int64) (bool, error) {
	row := p.db.QueryRowContext(ctx, `SELECT id FROM users WHERE telegram_id = $1`, telegramUserId)

	var id uuid.UUID
	err := row.Scan(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func (p *PostgresRepository) CreateUser(ctx context.Context, user *domain.User, wallet *domain.Wallet) error {
	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	_, err = tx.ExecContext(ctx, `INSERT INTO users (id, telegram_id, created_at) VALUES ($1, $2, $3)`,
		user.ID, user.TelegramID, user.CreatedAt)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return ErrUserAlreadyCreated
		}
		return err
	}

	_, err = tx.ExecContext(ctx, `INSERT INTO wallets (id, title, created_at, invite_code) VALUES ($1, $2, $3, $4)`,
		wallet.ID, wallet.Title, wallet.CreatedAt, wallet.InviteCode)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `INSERT INTO wallet_members (wallet_id, user_id, created_at, role) VALUES ($1, $2, $3, $4)`,
		wallet.ID, user.ID, user.CreatedAt, domain.OwnerWalletRole)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `INSERT INTO users_settings (user_id, active_wallet_id, personal_wallet_id) VALUES ($1, $2, $3)`,
		user.ID, wallet.ID, wallet.ID)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (p *PostgresRepository) GetUserID(ctx context.Context, telegramUserID int64) (uuid.UUID, error) {
	row := p.db.QueryRowContext(ctx, `SELECT id FROM users WHERE telegram_id = $1`, telegramUserID)
	var id uuid.UUID
	err := row.Scan(&id)
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}
