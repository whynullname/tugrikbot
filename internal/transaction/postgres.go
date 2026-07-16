package transaction

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

func (p *PostgresRepository) Save(ctx context.Context, transaction *domain.Transaction) error {
	_, err := p.db.ExecContext(ctx, `INSERT INTO transactions (id, user_id, wallet_id, amount, category, created_at, update_id, type)
							VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		transaction.ID, transaction.UserID, transaction.WalletID,
		transaction.Amount, transaction.Category, transaction.CreatedAt, transaction.UpdateID, transaction.Type)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return ErrTransactionAlreadyCreated
		}
		return err
	}

	return nil
}

func (p *PostgresRepository) GetWalletBalance(ctx context.Context, walletID uuid.UUID) (domain.Money, error) {
	row := p.db.QueryRowContext(ctx, `
		SELECT COALESCE(
			SUM(CASE WHEN type='expense' THEN -amount ELSE amount END),
			0
		) FROM transactions WHERE wallet_id = $1
	`, walletID)
	var money domain.Money
	err := row.Scan(&money)
	if err != nil {
		return 0, err
	}

	return money, nil
}
