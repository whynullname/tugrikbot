package wallet

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	GetUserWalletID(ctx context.Context, userID uuid.UUID) (uuid.UUID, error)
}
