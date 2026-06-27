package transaction

import (
	"context"
	"database/sql"

	"github.com/whynullname/tugrikbot/internal/domain"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (p *PostgresRepository) Save(ctx context.Context, transaction *domain.Transaction) error {
	_, err := p.db.ExecContext(ctx, `INSERT INTO transactions (id, user_id, amount, category, created_at)
							VALUES ($1, $2, $3, $4, $5)`,
		transaction.ID, transaction.UserID, transaction.Amount, transaction.Category, transaction.CreatedAt)

	if err != nil {
		return err
	}

	return nil
}
