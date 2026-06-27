package postgres

import (
	"context"
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func New(ctx context.Context, dbDSN string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dbDSN)
	if err != nil {
		return nil, err
	}

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}
