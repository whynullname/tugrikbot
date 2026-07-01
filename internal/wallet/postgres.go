package wallet

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
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
